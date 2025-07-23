package trading

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gateio/gateapi-go/v5"
)

// TradingBot หลักของระบบ
type TradingBot struct {
	client     *gateapi.APIClient
	ctx        context.Context
	aiClient   *AIClient
	indicators *Indicators
	gateClient *GateClient
}

// NewTradingBot สร้าง instance ใหม่
func NewTradingBot(apiKey, apiSecret, deepseekKey string) (*TradingBot, error) {
	// สร้าง Gate.io client
	client := gateapi.NewAPIClient(gateapi.NewConfiguration())
	ctx := context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
		Key:    apiKey,
		Secret: apiSecret,
	})

	// สร้าง AI client
	aiClient, err := NewAIClient(deepseekKey)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง AI client ได้: %v", err)
	}

	// สร้าง indicators
	indicators := NewIndicators()

	// สร้าง gate client wrapper
	gateClient := NewGateClient(client, ctx)

	return &TradingBot{
		client:     client,
		ctx:        ctx,
		aiClient:   aiClient,
		indicators: indicators,
		gateClient: gateClient,
	}, nil
}

// NewBacktesterBot สร้าง bot สำหรับ backtesting โดยไม่ใช้ AI
func NewBacktesterBot() (*TradingBot, error) {
	return &TradingBot{
		indicators: NewIndicators(),
	}, nil
}

// TestConnections ทดสอบการเชื่อมต่อทั้งหมด
func (bot *TradingBot) TestConnections() bool {
	fmt.Println("🔍 ทดสอบการเชื่อมต่อ Gate.io...")

	// ทดสอบ Gate.io
	if !bot.gateClient.TestConnection() {
		fmt.Println("❌ การเชื่อมต่อ Gate.io ไม่สำเร็จ")
		return false
	}

	fmt.Println("🔍 ทดสอบการเชื่อมต่อ AI...")

	// ทดสอบ AI
	if !bot.aiClient.TestConnection() {
		fmt.Println("❌ การเชื่อมต่อ AI ไม่สำเร็จ")
		return false
	}

	return true
}

// StartTradingLoop เริ่ม main trading loop
func (bot *TradingBot) StartTradingLoop() {
	fmt.Println("🔄 เริ่มต้น Main Trading Loop...")

	for {
		bot.runTradingCycle()

		// รอจนถึงชั่วโมงถัดไป
		bot.waitUntilNextHour()
	}
}

// runTradingCycle รัน 1 รอบของการเทรด
func (bot *TradingBot) runTradingCycle() {
	fmt.Println("📊 === เริ่มรอบการเทรดใหม่ ===")

	// ยกเลิกการตรวจสอบ BTC RSI - ให้ AI วิเคราะห์แต่ละเหรียญอิสระ
	fmt.Println("� LOOP1: ตรวจสอบและจัดการ Positions ที่เปิดอยู่...")
	bot.manageExistingPositions()

	fmt.Println("🔍 LOOP2: สแกนหาโอกาสใหม่...")
	bot.scanForNewOpportunities()
}

// manageExistingPositions จัดการ positions ที่เปิดอยู่
func (bot *TradingBot) manageExistingPositions() {
	fmt.Println("🔍 LOOP1: ตรวจสอบและจัดการ Positions ที่เปิดอยู่...")

	positions, err := bot.gateClient.GetOpenPositions()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล positions ได้: %v\n", err)
		return
	}

	if len(positions) == 0 {
		fmt.Println("📄 ไม่มี positions ที่เปิดอยู่")
		return
	}

	fmt.Printf("📊 พบ %d positions ที่เปิดอยู่\n", len(positions))

	for _, position := range positions {
		bot.analyzeExistingPosition(position)
	}
}

// analyzeExistingPosition วิเคราะห์ position ที่เปิดอยู่
func (bot *TradingBot) analyzeExistingPosition(position *Position) {
	contract := position.Contract
	fmt.Printf("🔍 วิเคราะห์ position: %s (Size: %d)\n", contract, position.Size)

	// ให้ AI ตัดสินใจเองว่าควรปิด position หรือไม่
	// ไม่ใช้ stop loss แบบเก่าอีกต่อไป
	fmt.Printf("🤖 ให้ AI วิเคราะห์และตัดสินใจเกี่ยวกับ position นี้\n")

	// ดึง OHLCV 100 แท่งสำหรับ AI analysis
	ohlcv, err := bot.gateClient.GetOHLCV(contract, "1h", 100)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล OHLCV สำหรับ AI: %v\n", err)
		return
	}

	// ส่งให้ AI วิเคราะห์
	decision, err := bot.aiClient.AnalyzeClosePosition(position, ohlcv)
	if err != nil {
		fmt.Printf("❌ AI วิเคราะห์ไม่สำเร็จสำหรับ %s: %v\n", contract, err)
		return
	}

	fmt.Printf("🤖 AI แนะนำ: %s (Confidence: %.1f%%)\n", decision.Action, decision.Confidence)

	// ปิด position ทันทีถ้า AI แนะนำ CLOSE (ไม่ดู confidence)
	if decision.Action == "CLOSE" {
		fmt.Printf("🔥 ปิด position ทันที - AI สัญญาณ CLOSE!\n")
		bot.closePosition(position)
	}
}

// scanForNewOpportunities สแกนหาโอกาสใหม่ (วิเคราะห์ทีละเหรียญทันที)
func (bot *TradingBot) scanForNewOpportunities() {
	fmt.Println("🔍 LOOP2: สแกนหาโอกาสใหม่...")

	// ดึงรายชื่อเหรียญทั้งหมดก่อน
	allContracts, err := bot.gateClient.GetFuturesContracts()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงรายชื่อ contracts ได้: %v\n", err)
		return
	}

	fmt.Printf("🔍 พบ %d contracts ทั้งหมด - เริ่ม scan ทีละเหรียญ\n", len(allContracts))

	// randomize order of contracts
	rand.Shuffle(len(allContracts), func(i, j int) {
		allContracts[i], allContracts[j] = allContracts[j], allContracts[i]
	})

	// ตรวจสอบ balance
	balance, err := bot.gateClient.GetBalance()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึง balance ได้: %v\n", err)
		return
	}

	availableBalance, _ := strconv.ParseFloat(balance, 64)
	fmt.Printf("💰 Balance ปัจจุบัน: %.2f USDT\n", availableBalance)

	if availableBalance < 10.0 {
		fmt.Println("❌ Balance ไม่เพียงพอ (ต้องการอย่างน้อย 10 USDT)")
		return
	}

	if availableBalance < 50.0 {
		fmt.Printf("⚠️ Balance เหลือน้อย (%.2f USDT < $50) - หยุด scan และรอจนถึงชั่วโมงถัดไป\n", availableBalance)
		fmt.Println("⏰ รอจนถึงชั่วโมงถัดไปก่อนจะ scan ใหม่...")
		bot.waitUntilNextHour()
		return
	}

	// วิเคราะห์เหรียญทีละตัว - ตรวจสอบ volume และวิเคราะห์ทันที
	contractsOpened := 0
	minVolumeUSDT := 1000000.0 // $1,000,000

	for i, contract := range allContracts {
		if availableBalance < 10.0 {
			fmt.Printf("⚠️ Balance ไม่เพียงพอสำหรับ position ถัดไป (< 10 USDT)\n")
			break
		}

		if availableBalance < 50.0 {
			fmt.Printf("⚠️ Balance เหลือน้อย (%.2f USDT < $50) - หยุด scan และรอจนถึงชั่วโมงถัดไป\n", availableBalance)
			fmt.Println("⏰ รอจนถึงชั่วโมงถัดไปก่อนจะ scan ใหม่...")
			bot.waitUntilNextHour()
			break
		}

		fmt.Printf("🔍 ตรวจสอบ %d/%d: %s\n", i+1, len(allContracts), contract)

		// ตรวจสอบ volume ก่อน
		ohlcv24h, err := bot.gateClient.GetOHLCV(contract, "1h", 24)
		if err != nil {
			fmt.Printf("⚠️ ไม่สามารถดึงข้อมูล %s: %v\n", contract, err)
			continue
		}

		if len(ohlcv24h) == 0 {
			fmt.Printf("❌ %s - ไม่มีข้อมูล\n", contract)
			continue
		}

		// คำนวณ volume 24 ชั่วโมง
		totalVolume := 0.0
		avgPrice := 0.0
		for _, candle := range ohlcv24h {
			totalVolume += candle.Volume
			avgPrice += (candle.High + candle.Low + candle.Close) / 3.0
		}

		if len(ohlcv24h) > 0 {
			avgPrice /= float64(len(ohlcv24h))
		}

		volumeUSDT := totalVolume * avgPrice

		if volumeUSDT < minVolumeUSDT {
			fmt.Printf("❌ %s - Volume: $%.0f (ต่ำเกินไป)\n", contract, volumeUSDT)
			time.Sleep(50 * time.Millisecond) // พักสั้นๆ
			continue
		}

		fmt.Printf("✅ %s - Volume: $%.0f (ผ่านเกณฑ์) - วิเคราะห์ทันที\n", contract, volumeUSDT)

		// วิเคราะห์ทันทีที่ผ่าน volume filter
		opened := bot.analyzeContract(contract)
		if opened {
			contractsOpened++
			availableBalance -= 10.0 // ลด balance เมื่อเปิด position (ประมาณ ~$10 ต่อ position)
			fmt.Printf("✅ เปิด position สำเร็จ: %s (รวม %d positions)\n", contract, contractsOpened)
			fmt.Printf("💰 Balance คงเหลือ: %.2f USDT\n", availableBalance)
		}

		// พักเล็กน้อยระหว่างเหรียญ
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("📊 สิ้นสุดการสแกน: เปิด %d positions\n", contractsOpened)
}

// analyzeContract วิเคราะห์เหรียญ 1 ตัว
func (bot *TradingBot) analyzeContract(contract string) bool {
	fmt.Printf("\n🔍 === วิเคราะห์ %s ===\n", contract)

	// ตรวจสอบว่ามี position ซ้ำหรือไม่
	fmt.Printf("1️⃣ ตรวจสอบ position ซ้ำ...\n")
	hasPosition, err := bot.gateClient.HasOpenPosition(contract)
	if err != nil {
		// แสดง error เฉพาะกรณีที่ไม่ใช่ contract ไม่มี
		if !strings.Contains(err.Error(), "POSITION_NOT_FOUND") {
			fmt.Printf("❌ ไม่สามารถตรวจสอบ position สำหรับ %s: %v\n", contract, err)
		}
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (ไม่สามารถตรวจสอบ position)\n")
		return false
	}

	if hasPosition {
		fmt.Printf("⚠️ มี position อยู่แล้ว - ข้ามไป\n")
		return false
	}
	fmt.Printf("✅ ไม่มี position ซ้ำ\n")

	// ดึง OHLCV 120 แท่ง (ใช้ 100 แท่งสุดท้าย)
	fmt.Printf("2️⃣ ดึงข้อมูล OHLCV (120 แท่ง)...\n")
	ohlcv, err := bot.gateClient.GetOHLCV(contract, "1h", 120)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล OHLCV: %v\n", err)
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (ไม่มีข้อมูล)\n")
		return false
	}

	if len(ohlcv) < 100 {
		fmt.Printf("❌ ข้อมูลไม่พอ (มี %d แท่ง, ต้องการ 100)\n", len(ohlcv))
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (ข้อมูลไม่เพียงพอ)\n")
		return false
	}
	fmt.Printf("✅ ดึงข้อมูล OHLCV ได้ %d แท่ง\n", len(ohlcv))

	// ใช้ 100 แท่งสุดท้าย
	analysisData := ohlcv[len(ohlcv)-100:]
	lastCandle := analysisData[len(analysisData)-1]
	fmt.Printf("📊 ราคาล่าสุด: O=%.6f H=%.6f L=%.6f C=%.6f V=%.0f\n",
		lastCandle.Open, lastCandle.High, lastCandle.Low, lastCandle.Close, lastCandle.Volume)

	// คำนวณ Technical Indicators
	fmt.Printf("3️⃣ คำนวณ Technical Indicators...\n")
	analysis := bot.indicators.AnalyzePivotPointSuperTrend(analysisData)

	// ข้ามการแสดง SuperTrend Analysis details เพราะให้ AI วิเคราะห์เลย
	fmt.Printf("📈 SuperTrend วิเคราะห์เสร็จ - ส่งให้ AI ตัดสินใจ\n")

	// ตัดสินใจเปิด position
	fmt.Printf("4️⃣ ตัดสินใจเปิด position...\n")
	result := bot.makePositionDecision(contract, analysis, analysisData)

	if result {
		fmt.Printf("🎯 ผลลัพธ์: เปิด position สำเร็จ ✅\n")
	} else {
		fmt.Printf("🎯 ผลลัพธ์: ไม่เปิด position ❌\n")
	}

	return result
}

// shouldConsiderForTrading ตรวจสอบว่าเหรียญควรพิจารณาสำหรับเทรดหรือไม่
func (bot *TradingBot) shouldConsiderForTrading(analysis *SuperTrendAnalysis, ohlcv []OHLCV) bool {
	fmt.Printf("🔍 === การคัดกรองเงื่อนไขแบบง่าย ===\n")

	if len(ohlcv) < 2 {
		fmt.Printf("❌ ข้อมูลไม่เพียงพอ (มี %d แท่ง, ต้องการ 2)\n", len(ohlcv))
		return false
	}

	// เงื่อนไข LONG: เทรนขาขึ้น + แท่งเทียนก่อนหน้าทับเส้น EMA100 + แท่งเขียว
	if analysis.Trend == 1 {
		prevCandle := ohlcv[len(ohlcv)-2]
		ema100 := analysis.EMA100

		fmt.Printf("📈 ตรวจสอบเงื่อนไข LONG:\n")
		fmt.Printf("   - Trend: %d (ขาขึ้น ✅)\n", analysis.Trend)
		fmt.Printf("   - EMA100: %.6f\n", ema100)
		fmt.Printf("   - แท่งก่อนหน้า: High=%.6f, Low=%.6f, Close=%.6f\n",
			prevCandle.High, prevCandle.Low, prevCandle.Close)

		// ตรวจสอบแท่งเขียว
		isGreenCandle := prevCandle.Close > prevCandle.Open
		fmt.Printf("   - แท่งเขียว: %v\n", isGreenCandle)

		// ตรวจสอบการทับเส้น EMA100 (High >= EMA100 >= Low)
		touchesEMA := prevCandle.Low <= ema100 && prevCandle.High >= ema100
		fmt.Printf("   - ทับเส้น EMA100: %v\n", touchesEMA)

		if isGreenCandle && touchesEMA {
			fmt.Printf("✅ ผ่านเงื่อนไข LONG\n")
			return true
		}
		fmt.Printf("❌ ไม่ผ่านเงื่อนไข LONG\n")
		return false
	}

	// เงื่อนไข SHORT: เทรนขาลง + แท่งเทียนก่อนหน้าทับเส้น EMA100 + แท่งแดง
	if analysis.Trend == -1 {
		prevCandle := ohlcv[len(ohlcv)-2]
		ema100 := analysis.EMA100

		fmt.Printf("📉 ตรวจสอบเงื่อนไข SHORT:\n")
		fmt.Printf("   - Trend: %d (ขาลง ✅)\n", analysis.Trend)
		fmt.Printf("   - EMA100: %.6f\n", ema100)
		fmt.Printf("   - แท่งก่อนหน้า: High=%.6f, Low=%.6f, Close=%.6f\n",
			prevCandle.High, prevCandle.Low, prevCandle.Close)

		// ตรวจสอบแท่งแดง
		isRedCandle := prevCandle.Close < prevCandle.Open
		fmt.Printf("   - แท่งแดง: %v\n", isRedCandle)

		// ตรวจสอบการทับเส้น EMA100 (High >= EMA100 >= Low)
		touchesEMA := prevCandle.Low <= ema100 && prevCandle.High >= ema100
		fmt.Printf("   - ทับเส้น EMA100: %v\n", touchesEMA)

		if isRedCandle && touchesEMA {
			fmt.Printf("✅ ผ่านเงื่อนไข SHORT\n")
			return true
		}
		fmt.Printf("❌ ไม่ผ่านเงื่อนไข SHORT\n")
		return false
	}

	fmt.Printf("❌ Trend ไม่ชัดเจน (%d)\n", analysis.Trend)
	return false
}

// makePositionDecision ตัดสินใจเปิด position
func (bot *TradingBot) makePositionDecision(contract string, analysis *SuperTrendAnalysis, ohlcv []OHLCV) bool {
	// เนื่องจากได้คัดเหรียญที่มี volume สูงมาแล้ว ให้ส่งไป AI วิเคราะห์เลย
	// ยกเลิกการตรวจสอบ shouldConsiderForTrading

	fmt.Printf("\n🎯 === ตัดสินใจเปิด Position ===\n")
	fmt.Printf("📊 %s: ส่งให้ AI วิเคราะห์ (ข้าม SuperTrend filtering)\n", contract)

	// ส่งให้ AI วิเคราะห์ทุกกรณี
	fmt.Printf("🤖 AI Mode - ส่งให้ AI วิเคราะห์\n")
	return bot.openPositionAI(contract, analysis, ohlcv)
}

// openPositionAuto เปิด position ในโหมด AUTO
// openPositionAI เปิด position ในโหมด AI
func (bot *TradingBot) openPositionAI(contract string, analysis *SuperTrendAnalysis, ohlcv []OHLCV) bool {
	fmt.Printf("\n🤖 === AI POSITION ANALYSIS ===\n")
	fmt.Printf("🤖 ส่งให้ AI วิเคราะห์: %s\n", contract)
	fmt.Printf("📤 กำลังส่งข้อมูลไปยัง AI...\n")

	// ส่งให้ AI วิเคราะห์
	decision, err := bot.aiClient.AnalyzeOpenPosition(contract, analysis, ohlcv)
	if err != nil {
		fmt.Printf("❌ AI วิเคราะห์ไม่สำเร็จสำหรับ %s: %v\n", contract, err)
		return false
	}

	fmt.Printf("📥 AI ตอบกลับแล้ว!\n")
	fmt.Printf("🎯 AI แนะนำ: %s (Confidence: %.1f%%)\n", decision.Action, decision.Confidence)
	fmt.Printf("📊 AI Risk-Reward Ratio: %.2f\n", decision.RiskRewardRatio)
	fmt.Printf("�💭 AI เหตุผล: %s\n", decision.Reason)

	// เปิด position ถ้า AI แนะนำ LONG หรือ SHORT
	if decision.Action == "LONG" || decision.Action == "SHORT" {
		fmt.Printf("✅ AI แนะนำเปิด position!\n")
		fmt.Printf("💰 กำลังเปิด position ตาม AI...\n")

		side := "long"
		if decision.Action == "SHORT" {
			side = "short"
		}

		fmt.Printf("📈 ทิศทาง: %s\n", decision.Action)

		// เปิด position
		success, err := bot.gateClient.OpenPositionWithSize(contract, side, 5.0)

		if err != nil {
			fmt.Printf("❌ ไม่สามารถเปิด position %s: %v\n", contract, err)
			return false
		}

		if success {
			fmt.Printf("✅ เปิด position %s สำเร็จ (%s)\n", contract, side)
			fmt.Printf("🎉 AI MODE สำเร็จ!\n")
			return true
		}
	} else {
		fmt.Printf("❌ AI แนะนำ %s - ไม่เปิด position\n", decision.Action)
	}

	return false
}

// closePosition ปิด position
func (bot *TradingBot) closePosition(position *Position) {
	fmt.Printf("🔚 ปิด position: %s\n", position.Contract)

	success, err := bot.gateClient.ClosePosition(position.Contract)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถปิด position %s: %v\n", position.Contract, err)
		return
	}

	if success {
		fmt.Printf("✅ ปิด position %s สำเร็จ\n", position.Contract)
	}
}

// waitUntilNextHour รอจนถึงชั่วโมงถัดไป
func (bot *TradingBot) waitUntilNextHour() {
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(time.Hour)
	duration := nextHour.Sub(now)

	fmt.Printf("⏰ รอจนถึงชั่วโมงถัดไป: %v\n", duration)
	time.Sleep(duration)
}
