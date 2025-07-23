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

	// สร้าง gate client wrapper
	gateClient := NewGateClient(client, ctx)

	return &TradingBot{
		client:     client,
		ctx:        ctx,
		aiClient:   aiClient,
		gateClient: gateClient,
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

	// LOOP1: ตรวจสอบและจัดการ positions ที่เปิดอยู่
	bot.manageExistingPositions()

	// LOOP2: สแกนหาโอกาสใหม่
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
	fmt.Printf("🤖 ให้ AI วิเคราะห์และตัดสินใจเกี่ยวกับ position นี้\n")

	// ดึง OHLCV 144 แท่งสำหรับ AI analysis (1 ชั่วโมง)
	ohlcv, err := bot.gateClient.GetOHLCV(contract, "1h", 144)
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
	fmt.Printf("🤖 เหตุผล: %s\n", decision.Reason)

	// ปิด position ทันทีถ้า AI แนะนำ CLOSE (ไม่ดู confidence)
	if decision.Action == "CLOSE" {
		fmt.Printf("🔥 ปิด position ทันที - AI สัญญาณ CLOSE!\n")
		bot.closePosition(position)
	}
}

// scanForNewOpportunities สแกนหาโอกาสใหม่
func (bot *TradingBot) scanForNewOpportunities() {
	fmt.Println("🔍 LOOP2: สแกนหาโอกาสใหม่...")

	// ดึงรายชื่อเหรียญทั้งหมด
	contracts, err := bot.gateClient.GetFuturesContracts()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงรายชื่อ contracts ได้: %v\n", err)
		return
	}

	fmt.Printf("🔎 พบ %d contracts ทั้งหมด\n", len(contracts))

	// randomize order of contracts
	rand.Shuffle(len(contracts), func(i, j int) {
		contracts[i], contracts[j] = contracts[j], contracts[i]
	})

	// ตรวจสอบ balance
	balance, err := bot.gateClient.GetBalance()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึง balance ได้: %v\n", err)
		return
	}

	availableBalance, _ := strconv.ParseFloat(balance, 64)
	fmt.Printf("💰 Balance ปัจจุบัน: %.2f USDT\n", availableBalance)

	if availableBalance < 15.0 {
		fmt.Println("❌ Balance ไม่เพียงพอ (ต้องการอย่างน้อย 15 USDT)")
		return
	}

	// วิเคราะห์เหรียญทีละตัว (ไม่ใช้ batch)
	contractsOpened := 0
	for i, contract := range contracts {
		// ตรวจสอบ balance แบบ real-time ก่อนวิเคราะห์แต่ละเหรียญ
		currentBalance, err := bot.gateClient.GetBalance()
		if err != nil {
			fmt.Printf("❌ ไม่สามารถดึง balance ได้: %v\n", err)
			break
		}

		currentBalanceFloat, _ := strconv.ParseFloat(currentBalance, 64)

		if currentBalanceFloat < 15.0 { // ต้องการ 15 USDT สำหรับ position ใหม่
			fmt.Printf("💰 Balance ปัจจุบัน: %.2f USDT\n", currentBalanceFloat)
			fmt.Printf("⚠️ เงินไม่เพียงพอสำหรับ position ใหม่ (ต้องการ 15 USDT)\n")
			fmt.Printf("⏸️ หยุดสแกนและรอชั่วโมงถัดไป...\n")
			break
		}

		fmt.Printf("🔍 วิเคราะห์ %d/%d: %s (Balance: %.2f USDT)\n", i+1, len(contracts), contract, currentBalanceFloat)

		opened := bot.analyzeContract(contract)
		if opened {
			contractsOpened++
			fmt.Printf("✅ เปิด position สำเร็จ: %s (รวม %d positions)\n", contract, contractsOpened)
		}

		// พักเล็กน้อยระหว่างเหรียญ (0.1 วินาที)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("📊 สิ้นสุดการสแกน: เปิด %d positions จาก %d contracts\n", contractsOpened, len(contracts))
}

// analyzeContract วิเคราะห์เหรียญ 1 ตัว (Pure AI)
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

	// ดึง OHLCV 244 แท่ง (สำหรับ EMA200 filter)
	fmt.Printf("2️⃣ ดึงข้อมูล OHLCV (244 แท่ง, 1 ชั่วโมง)...\n")
	ohlcvFull, err := bot.gateClient.GetOHLCV(contract, "1h", 244)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล OHLCV: %v\n", err)
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (ไม่มีข้อมูล)\n")
		return false
	}

	if len(ohlcvFull) < 244 {
		fmt.Printf("❌ ข้อมูลไม่พอ (มี %d แท่ง, ต้องการ 244)\n", len(ohlcvFull))
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (ข้อมูลไม่เพียงพอ)\n")
		return false
	}
	fmt.Printf("✅ ดึงข้อมูล OHLCV ได้ %d แท่ง\n", len(ohlcvFull))

	// ส่งข้อมูลครบทั้งหมด 244 แท่งให้ AI
	ohlcv := ohlcvFull

	lastCandle := ohlcv[len(ohlcv)-1]
	fmt.Printf("📊 ราคาล่าสุด: O=%.6f H=%.6f L=%.6f C=%.6f V=%.0f\n",
		lastCandle.Open, lastCandle.High, lastCandle.Low, lastCandle.Close, lastCandle.Volume)

	// Pre-filter ด้วยการตรวจสอบ EMA200 และ EMA Alignment เพื่อประหยัด AI calls
	fmt.Printf("3️⃣ ตรวจสอบ EMA200 + EMA Alignment pre-filter...\n")
	if !bot.checkEMA200Condition(ohlcvFull) {
		fmt.Printf("❌ EMA200 ไม่เข้าเงื่อนไข - ข้ามการเรียก AI\n")
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (EMA200 pre-filter)\n")
		return false
	}

	// ตรวจสอบ EMA Alignment
	if !bot.checkEMAAlignment(ohlcv) {
		fmt.Printf("❌ EMA Alignment ไม่เข้าเงื่อนไข - ข้ามการเรียก AI\n")
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (EMA Alignment pre-filter)\n")
		return false
	}

	fmt.Printf("✅ EMA200 + EMA Alignment เข้าเงื่อนไข - ดำเนินการเรียก AI\n")

	// ให้ AI ตัดสินใจเอง
	fmt.Printf("4️⃣ ส่งข้อมูลให้ AI วิเคราะห์...\n")
	decision, err := bot.aiClient.AnalyzeOpenPosition(contract, ohlcv)
	if err != nil {
		fmt.Printf("❌ AI วิเคราะห์ไม่สำเร็จ: %v\n", err)
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (AI error)\n")
		return false
	}

	fmt.Printf("🤖 AI แนะนำ: %s (Confidence: %.1f%%)\n", decision.Action, decision.Confidence)
	fmt.Printf("🤖 เหตุผล: %s\n", decision.Reason)

	// ตัดสินใจเปิด position
	fmt.Printf("5️⃣ ตัดสินใจเปิด position...\n")
	result := bot.executeAIDecision(contract, decision)

	if result {
		fmt.Printf("🎯 ผลลัพธ์: เปิด position สำเร็จ ✅\n")
	} else {
		fmt.Printf("🎯 ผลลัพธ์: ไม่เปิด position ❌\n")
	}

	return result
}

// executeAIDecision ดำเนินการตามคำแนะนำของ AI
func (bot *TradingBot) executeAIDecision(contract string, decision *AIDecision) bool {
	fmt.Printf("🎯 ตัดสินใจ AI: %s (Confidence: %.1f%%)\n", decision.Action, decision.Confidence)

	// ตรวจสอบ action เท่านั้น (ไม่ดู confidence)
	if decision.Action == "HOLD" {
		fmt.Printf("⏸️ AI แนะนำ HOLD - ไม่เปิด position\n")
		return false
	}

	if decision.Action != "LONG" && decision.Action != "SHORT" {
		fmt.Printf("❌ AI action ไม่ถูกต้อง: %s (ต้องเป็น LONG/SHORT/HOLD)\n", decision.Action)
		return false
	}

	// เปิด position ตาม AI (ไม่สนใจ confidence)
	fmt.Printf("✅ AI แนะนำ %s - เปิด position ทันที!\n", decision.Action)

	side := "long"
	if decision.Action == "SHORT" {
		side = "short"
	}

	success, err := bot.gateClient.OpenPosition(contract, side, 5.0) // 5x leverage
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเปิด position ได้: %v\n", err)
		return false
	}

	if success {
		fmt.Printf("🎉 เปิด position %s %s สำเร็จ!\n", side, contract)
		return true
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

// checkEMA200Condition ตรวจสอบ EMA200 เพื่อประหยัด AI calls
func (bot *TradingBot) checkEMA200Condition(ohlcv []OHLCV) bool {
	if len(ohlcv) < 244 {
		return false
	}

	// ตรวจสอบว่าใน 48 แท่งย้อนหลัง มีแท่งใดทับ EMA200 หรือไม่
	ema200Condition := bot.checkEMA200TouchAny24Candles(ohlcv)

	if ema200Condition {
		fmt.Printf("   📊 พบแท่งเทียนใน 48 แท่งย้อนหลัง ทับเส้น EMA200 (EMA200 CONDITION)\n")
	} else {
		fmt.Printf("   📊 ไม่พบแท่งเทียนใน 48 แท่งย้อนหลัง ทับเส้น EMA200\n")
	}

	if !ema200Condition {
		fmt.Printf("   ❌ ไม่ผ่านเงื่อนไข EMA200\n")
	}

	return ema200Condition
}

// checkEMA200TouchAny24Candles ตรวจสอบว่าใน 48 แท่งย้อนหลัง มีแท่งใดทับ EMA200 หรือไม่
func (bot *TradingBot) checkEMA200TouchAny24Candles(ohlcv []OHLCV) bool {
	if len(ohlcv) < 244 {
		return false
	}

	// คำนวณ EMA200 ทั้งหมด
	ema200Values := bot.calculateEMA200Series(ohlcv)
	if len(ema200Values) == 0 {
		return false
	}

	// ตรวจสอบ 48 แท่งย้อนหลัง
	startIdx := len(ohlcv) - 48
	if startIdx < 0 {
		startIdx = 0
	}

	touchCount := 0
	for i := startIdx; i < len(ohlcv); i++ {
		if i >= len(ema200Values) {
			continue
		}

		candle := ohlcv[i]
		ema200 := ema200Values[i]

		// ตรวจสอบว่าแท่งเทียนทับเส้น EMA200 (EMA200 อยู่ระหว่าง Low และ High)
		if candle.Low <= ema200 && ema200 <= candle.High {
			touchCount++
			fmt.Printf("   ✅ แท่งที่ %d ทับ EMA200: EMA200=%.6f, Low=%.6f, High=%.6f\n",
				i-startIdx+1, ema200, candle.Low, candle.High)
			return true // พอพบแท่งแรกที่ทับก็ return true ทันที
		}
	}

	fmt.Printf("   ❌ ไม่พบแท่งใดใน 48 แท่งย้อนหลัง ที่ทับ EMA200\n")
	return false
}

// calculateEMA200Series คำนวณ EMA200 สำหรับทุกแท่งเทียน
func (bot *TradingBot) calculateEMA200Series(ohlcv []OHLCV) []float64 {
	if len(ohlcv) < 200 {
		return []float64{}
	}

	result := make([]float64, len(ohlcv))
	k := 2.0 / float64(200+1) // multiplier สำหรับ EMA200

	// เริ่มต้นด้วยราคาปิดแรก
	result[0] = ohlcv[0].Close

	// คำนวณ EMA200 สำหรับทุกแท่ง
	for i := 1; i < len(ohlcv); i++ {
		result[i] = ohlcv[i].Close*k + result[i-1]*(1-k)
	}

	return result
}

// checkEMAAlignment ตรวจสอบ EMA Alignment สำหรับ LONG/SHORT
func (bot *TradingBot) checkEMAAlignment(ohlcv []OHLCV) bool {
	if len(ohlcv) < 200 {
		return false
	}

	// คำนวณ EMA20, EMA50, EMA100, EMA200
	ema20 := bot.calculateEMA(ohlcv, 20)
	ema50 := bot.calculateEMA(ohlcv, 50)
	ema100 := bot.calculateEMA(ohlcv, 100)
	ema200 := bot.calculateEMA(ohlcv, 200)

	// ตรวจสอบ LONG alignment: EMA20 > EMA50 > EMA100 > EMA200
	longAlignment := ema20 > ema50 && ema50 > ema100 && ema100 > ema200

	// ตรวจสอบ SHORT alignment: EMA20 < EMA50 < EMA100 < EMA200
	shortAlignment := ema20 < ema50 && ema50 < ema100 && ema100 < ema200

	if longAlignment {
		fmt.Printf("   📈 LONG EMA Alignment: EMA20>EMA50>EMA100>EMA200 ✅\n")
		return true
	} else if shortAlignment {
		fmt.Printf("   📉 SHORT EMA Alignment: EMA20<EMA50<EMA100<EMA200 ✅\n")
		return true
	} else {
		fmt.Printf("   📊 EMA Alignment ไม่เข้าเงื่อนไข (ไม่ใช่ LONG หรือ SHORT alignment)\n")
		return false
	}
}

// calculateEMA คำนวณ Exponential Moving Average
func (bot *TradingBot) calculateEMA(ohlcv []OHLCV, period int) float64 {
	if len(ohlcv) < period {
		return 0
	}

	// คำนวณ EMA โดยใช้ Close price
	k := 2.0 / float64(period+1) // multiplier
	ema := ohlcv[0].Close        // เริ่มต้นด้วยราคาปิดแรก

	for i := 1; i < len(ohlcv); i++ {
		ema = ohlcv[i].Close*k + ema*(1-k)
	}

	return ema
}
