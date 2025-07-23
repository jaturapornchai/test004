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

	// ดึง OHLCV 288 แท่งสำหรับ AI analysis (1 ชั่วโมง)
	ohlcv, err := bot.gateClient.GetOHLCV(contract, "1h", 288)
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
		if availableBalance < 15.0 {
			fmt.Printf("⚠️ Balance ไม่เพียงพอสำหรับ position ถัดไป (< 15 USDT)\n")
			break
		}

		fmt.Printf("🔍 วิเคราะห์ %d/%d: %s\n", i+1, len(contracts), contract)

		opened := bot.analyzeContract(contract)
		if opened {
			contractsOpened++
			availableBalance -= 15.0 // ลด balance เมื่อเปิด position (15 USDT margin)
			fmt.Printf("✅ เปิด position สำเร็จ: %s (รวม %d positions)\n", contract, contractsOpened)
			fmt.Printf("💰 Balance คงเหลือ: %.2f USDT\n", availableBalance)
		}

		// พักเล็กน้อยระหว่างเหรียญ (1 วินาที)
		time.Sleep(1 * time.Second)
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

	// ดึง OHLCV 288 แท่ง (1 ชั่วโมง)
	fmt.Printf("2️⃣ ดึงข้อมูล OHLCV (288 แท่ง, 1 ชั่วโมง)...\n")
	ohlcv, err := bot.gateClient.GetOHLCV(contract, "1h", 288)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล OHLCV: %v\n", err)
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (ไม่มีข้อมูล)\n")
		return false
	}

	if len(ohlcv) < 288 {
		fmt.Printf("❌ ข้อมูลไม่พอ (มี %d แท่ง, ต้องการ 288)\n", len(ohlcv))
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (ข้อมูลไม่เพียงพอ)\n")
		return false
	}
	fmt.Printf("✅ ดึงข้อมูล OHLCV ได้ %d แท่ง\n", len(ohlcv))

	lastCandle := ohlcv[len(ohlcv)-1]
	fmt.Printf("📊 ราคาล่าสุด: O=%.6f H=%.6f L=%.6f C=%.6f V=%.0f\n",
		lastCandle.Open, lastCandle.High, lastCandle.Low, lastCandle.Close, lastCandle.Volume)

	// Pre-filter ด้วย EMA เพื่อประหยัด AI calls
	fmt.Printf("3️⃣ ตรวจสอบ EMA pre-filter...\n")
	if !bot.checkEMAConditions(ohlcv) {
		fmt.Printf("❌ EMA ไม่เข้าเงื่อนไข - ข้ามการเรียก AI\n")
		fmt.Printf("⏭️ ข้ามเหรียญนี้ (EMA pre-filter)\n")
		return false
	}
	fmt.Printf("✅ EMA เข้าเงื่อนไข - ดำเนินการเรียก AI\n")

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

// checkEMAConditions ตรวจสอบ EMA เพื่อประหยัด AI calls
func (bot *TradingBot) checkEMAConditions(ohlcv []OHLCV) bool {
	if len(ohlcv) < 200 {
		return false
	}

	// คำนวณ EMA 20, 50, 100, 200
	ema20 := bot.gateClient.CalculateEMA(ohlcv, 20)
	ema50 := bot.gateClient.CalculateEMA(ohlcv, 50)
	ema100 := bot.gateClient.CalculateEMA(ohlcv, 100)
	ema200 := bot.gateClient.CalculateEMA(ohlcv, 200)

	prevCandle := ohlcv[len(ohlcv)-2]

	// เงื่อนไข LONG: EMA เรียง 20>50>100>200 + แท่งก่อนหน้าเขียว + ทับ EMA100
	longCondition := ema20 > ema50 && ema50 > ema100 && ema100 > ema200 &&
		prevCandle.Close > prevCandle.Open && // แท่งก่อนหน้าเขียว
		prevCandle.Close > ema100 // ทับ EMA100

	// เงื่อนไข SHORT: EMA เรียง 20<50<100<200 + แท่งก่อนหน้าแดง + ทับ EMA100
	shortCondition := ema20 < ema50 && ema50 < ema100 && ema100 < ema200 &&
		prevCandle.Close < prevCandle.Open && // แท่งก่อนหน้าแดง
		prevCandle.Close < ema100 // ทับ EMA100

	return longCondition || shortCondition
}
