package trading

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"binance-trading-bot/internal/binance"
)

// TradingBot หลักของระบบ
type TradingBot struct {
	binanceClient *BinanceClient
	aiClient      *AIClient
}

// NewTradingBot สร้าง instance ใหม่
func NewTradingBot(apiKey, apiSecret, deepseekKey string) (*TradingBot, error) {
	// สร้าง Binance client
	client := binance.NewClient(apiKey, apiSecret, "https://fapi.binance.com")

	// สร้าง AI client
	aiClient, err := NewAIClient(deepseekKey)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง AI client ได้: %v", err)
	}

	// สร้าง binance client wrapper
	binanceClient := NewBinanceClient(client)

	return &TradingBot{
		binanceClient: binanceClient,
		aiClient:      aiClient,
	}, nil
}

// TestConnections ทดสอบการเชื่อมต่อทั้งหมด
func (bot *TradingBot) TestConnections() bool {
	fmt.Println("🔍 ทดสอบการเชื่อมต่อ Binance...")
	if !bot.binanceClient.TestConnection() {
		fmt.Println("❌ การเชื่อมต่อ Binance ไม่สำเร็จ")
		return false
	}

	fmt.Println("🔍 ทดสอบการเชื่อมต่อ AI...")
	if !bot.aiClient.TestConnection() {
		fmt.Println("❌ การเชื่อมต่อ AI ไม่สำเร็จ")
		return false
	}

	fmt.Println("✅ การเชื่อมต่อทั้งหมดสำเร็จ!")
	return true
}

// Start เริ่มระบบเทรด
func (bot *TradingBot) Start() {
	fmt.Println("🚀 เริ่มระบบเทรด Binance AI Bot...")

	// ทดสอบการเชื่อมต่อก่อน
	if !bot.TestConnections() {
		fmt.Println("❌ ไม่สามารถเริ่มระบบได้เนื่องจากการเชื่อมต่อไม่สำเร็จ")
		return
	}

	// เริ่ม loop หลัก
	for {
		bot.runTradingCycle()
		// ไม่ต้องรอในลูปหลัก เพราะ findTradingOpportunity จะรอชั่วโมงถัดไปเอง
	}
}

// runTradingCycle รันหนึ่งรอบของการเทรด
func (bot *TradingBot) runTradingCycle() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🔄 เริ่มรอบการวิเคราะห์ใหม่...")
	fmt.Println(strings.Repeat("=", 50))

	// ตรวจสอบ positions ที่เปิดอยู่
	positions, err := bot.binanceClient.GetOpenPositions()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล positions ได้: %v\n", err)
		return
	}

	// กรองเฉพาะ positions ที่เปิดอยู่จริง (PositionAmt != "0")
	var activePositions []binance.Position
	for _, position := range positions {
		// แปลง PositionAmt เป็น float เพื่อเช็คว่าเป็น 0 หรือไม่
		positionAmt, err := strconv.ParseFloat(position.PositionAmt, 64)
		if err == nil && positionAmt != 0.0 {
			activePositions = append(activePositions, position)
		}
	}

	// ถ้ามี position เปิดอยู่จริง ให้วิเคราะห์ว่าควรปิดหรือไม่
	if len(activePositions) > 0 {
		fmt.Printf("📊 พบ %d position(s) ที่เปิดอยู่จริง\n", len(activePositions))
		for _, position := range activePositions {
			fmt.Printf("📊 กำลังวิเคราะห์ position: %s (ปริมาณ: %s)\n", position.Symbol, position.PositionAmt)
			bot.analyzeClosePosition(position)
		}
	}

	// หาโอกาสใหม่ทุกครั้ง (ไม่ว่าจะมี position หรือไม่)
	fmt.Println("📈 กำลังหาโอกาสใหม่...")
	bot.findTradingOpportunity()
}

// findTradingOpportunity หาโอกาสในการเทรด
func (bot *TradingBot) findTradingOpportunity() {
	fmt.Println("🔍 กำลังหาโอกาสในการเทรด...")
	fmt.Println(strings.Repeat("=", 30) + " SYMBOL ANALYSIS " + strings.Repeat("=", 30))

	// ดึงรายการ contracts
	contracts, err := bot.binanceClient.GetFuturesContracts()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงรายการ symbols ได้: %v\n", err)
		return
	}

	fmt.Printf("� จำนวน symbols ทั้งหมดที่จะวิเคราะห์: %d symbols\n", len(contracts))

	// ตรวจสอบ balance ก่อน
	balance, err := bot.binanceClient.GetBalance()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึง balance ได้: %v\n", err)
		return
	}

	balanceFloat, _ := strconv.ParseFloat(balance, 64)
	maxPositions := int(balanceFloat / 21) // 20 USDT margin + 1 USDT buffer per position
	fmt.Printf("💰 Balance เริ่มต้น: %s USDT\n", balance)
	fmt.Printf("🎯 สามารถเทรดได้สูงสุด: %d positions (20 USDT margin + 1 USDT buffer ต่อ position)\n", maxPositions)

	// วิเคราะห์แต่ละ contract จนกว่าเงินจะหมด
	analyzedCount := 0
	tradedCount := 0
	skippedCount := 0

	for i, contract := range contracts {
		fmt.Printf("\n🔍 [%d/%d] กำลังวิเคราะห์ %s\n", i+1, len(contracts), contract)
		fmt.Printf("📊 สถานะ: วิเคราะห์แล้ว %d | เทรดแล้ว %d | ข้ามแล้ว %d\n", analyzedCount, tradedCount, skippedCount)

		// ตรวจสอบ balance ก่อนการวิเคราะห์แต่ละครั้ง
		currentBalance, err := bot.binanceClient.GetBalance()
		if err != nil {
			fmt.Printf("❌ ไม่สามารถดึง balance ได้: %v\n", err)
			continue
		}

		balanceFloat, err := strconv.ParseFloat(currentBalance, 64)
		if err != nil || balanceFloat < 21 { // ต้องมี balance อย่างน้อย 21 USDT (20 USDT margin + 1 USDT buffer)
			fmt.Printf("💰 Balance ไม่เพียงพอแล้ว! (มี: %s USDT, ต้องการ: 21 USDT)\n", currentBalance)
			fmt.Printf("🛑 หยุดการวิเคราะห์ symbols เพิ่มเติม\n")
			break // หยุดเมื่อเงินไม่พอ แล้วจะไปรอชั่วโมงถัดไป
		}

		fmt.Printf("💰 Balance ปัจจุบัน: %s USDT (เพียงพอสำหรับ %.0f positions)\n",
			currentBalance, balanceFloat/21)

		// วิเคราะห์ว่าควรเทรดหรือไม่
		shouldTrade, side := bot.shouldOpenPosition(contract)
		analyzedCount++

		if shouldTrade {
			fmt.Printf("🎯 พบโอกาส! เทรด %s ทิศทาง %s\n", contract, side)
			bot.openPosition(contract, side)
			tradedCount++
			fmt.Printf("✅ เทรด %s เสร็จแล้ว ต่อไปวิเคราะห์ symbol ถัดไป\n", contract)
			// ไม่ break ให้ทำต่อไปจนกว่าเงินจะหมด
		} else {
			skippedCount++
			fmt.Printf("⏭️ ข้าม %s - ไม่พบสัญญาณที่ชัดเจน\n", contract)
		}

		// หยุดพักเล็กน้อยระหว่างการวิเคราะห์
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n📈 สรุปผลการวิเคราะห์ครั้งนี้:\n")
	fmt.Printf("🔍 ทั้งหมด: %d symbols\n", analyzedCount)
	fmt.Printf("✅ เทรดแล้ว: %d positions\n", tradedCount)
	fmt.Printf("⏭️ ข้ามไป: %d symbols\n", skippedCount)
	if analyzedCount > 0 {
		fmt.Printf("📊 อัตราความสำเร็จ: %.1f%%\n",
			float64(tradedCount)/float64(analyzedCount)*100)
	}
	fmt.Printf("🏁 เสร็จสิ้นการวิเคราะห์ทั้งหมด\n")

	// รอจนถึง 1 ชั่วโมงถัดไป
	bot.waitForNext1Hour()
}

// waitForNext1Hour รอจนถึง 1 ชั่วโมงถัดไป (0:00)
func (bot *TradingBot) waitForNext1Hour() {
	now := time.Now()
	// หาช่วง 1 ชั่วโมงถัดไป (0:00)
	minutesToNext := 60 - now.Minute()
	next1Hour := now.Add(time.Duration(minutesToNext) * time.Minute).Truncate(time.Minute)
	waitTime := next1Hour.Sub(now)

	fmt.Printf("⏰ รอจนถึง 1 ชั่วโมงถัดไป (%s) - เหลือเวลา: %v\n",
		next1Hour.Format("15:04"), waitTime.Round(time.Second))

	time.Sleep(waitTime)
}

// shouldOpenPosition วิเคราะห์ว่าควรเปิด position หรือไม่
func (bot *TradingBot) shouldOpenPosition(contract string) (bool, string) {
	// ตรวจสอบว่ามี position เปิดอยู่หรือไม่
	positions, err := bot.binanceClient.GetOpenPositions()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถตรวจสอบ position ได้: %v\n", err)
		return false, ""
	}

	// Debug: แสดง positions ทั้งหมดที่ดึงมา
	if len(positions) > 0 {
		fmt.Printf("🔍 พบ %d positions เปิดอยู่:\n", len(positions))
		for _, pos := range positions {
			posAmt, _ := strconv.ParseFloat(pos.PositionAmt, 64)
			fmt.Printf("   %s: amount=%s (%.8f)\n", pos.Symbol, pos.PositionAmt, posAmt)
		}
	}

	for _, pos := range positions {
		// แปลงเป็น float เพื่อตรวจสอบว่าเป็น 0 หรือไม่อย่างแม่นยำ
		posAmt, err := strconv.ParseFloat(pos.PositionAmt, 64)
		if err == nil && pos.Symbol == contract && math.Abs(posAmt) > 0.0001 {
			fmt.Printf("⚠️ มี position %s เปิดอยู่แล้ว (amount: %s)\n", contract, pos.PositionAmt)
			return false, ""
		}
	}

	// ดึงข้อมูล candlestick
	candlesticks, err := bot.binanceClient.GetCandlesticks(contract, "1h", 288)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล candlestick สำหรับ %s ได้: %v\n", contract, err)
		return false, ""
	}

	if len(candlesticks) < 288 {
		fmt.Printf("⚠️ ข้อมูล candlestick ไม่เพียงพอสำหรับ %s (มี %d แท่ง ต้องการ 288 แท่ง)\n", contract, len(candlesticks))
		return false, ""
	}

	// ตรวจสอบ EMA ก่อนส่งไป AI
	fmt.Printf("📊 ตรวจสอบ EMA สำหรับ %s...\n", contract)
	emaSignal := bot.checkEMAFilter(candlesticks)
	if emaSignal == "HOLD" {
		fmt.Printf("⚠️ EMA Filter: ข้าม %s - ไม่ผ่านเกณฑ์ EMA\n", contract)
		return false, ""
	}
	fmt.Printf("✅ EMA Filter: %s ผ่านเกณฑ์ EMA - สัญญาณ %s\n", contract, emaSignal)

	// แปลงข้อมูลให้อยู่ในรูปแบบที่ AI ต้องการ
	ohlcvSlice := bot.convertToOHLCV(candlesticks)

	// ใช้ AI วิเคราะห์
	decision, err := bot.aiClient.AnalyzeOpenPosition(contract, ohlcvSlice)
	if err != nil {
		fmt.Printf("❌ AI analysis error: %v\n", err)
		return false, ""
	}

	fmt.Printf("🤖 AI Analysis: %s\n", decision.Action)

	// ตัดสินใจตาม AI (ไม่ใช้ confidence แล้ว)
	action := strings.ToLower(decision.Action)
	if action == "long" || strings.Contains(action, "buy") {
		return true, "BUY"
	} else if action == "short" || strings.Contains(action, "sell") {
		return true, "SELL"
	}

	return false, ""
}

// convertToOHLCV แปลงข้อมูล Binance candlestick เป็น OHLCV
func (bot *TradingBot) convertToOHLCV(candlesticks []binance.Candlestick) []OHLCV {
	var ohlcvData []OHLCV

	for _, candle := range candlesticks {
		open, _ := strconv.ParseFloat(candle.Open, 64)
		high, _ := strconv.ParseFloat(candle.High, 64)
		low, _ := strconv.ParseFloat(candle.Low, 64)
		close, _ := strconv.ParseFloat(candle.Close, 64)
		volume, _ := strconv.ParseFloat(candle.Volume, 64)

		ohlcv := OHLCV{
			Timestamp: candle.OpenTime,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}
		ohlcvData = append(ohlcvData, ohlcv)
	}

	return ohlcvData
}

// openPosition เปิด position ใหม่
func (bot *TradingBot) openPosition(contract string, side string) {
	fmt.Printf("🔥 กำลังเปิด position %s ทิศทาง %s...\n", contract, side)

	// ตั้งค่า leverage และ margin type ก่อนเปิด position
	fmt.Printf("⚙️ กำลังตรวจสอบและตั้งค่า leverage และ margin type สำหรับ %s...\n", contract)

	// ตั้งค่า leverage เป็น 10x
	err := bot.binanceClient.SetLeverage(contract, 10)
	if err != nil {
		fmt.Printf("⚠️ ไม่สามารถตั้งค่า leverage ได้: %v (อาจตั้งไว้แล้ว)\n", err)
	} else {
		fmt.Printf("✅ ตั้งค่า leverage เป็น 10x สำเร็จ\n")
	}

	// ตั้งค่า margin type เป็น isolated
	err = bot.binanceClient.SetMarginType(contract, "ISOLATED")
	if err != nil {
		fmt.Printf("⚠️ ไม่สามารถตั้งค่า margin type ได้: %v (อาจตั้งไว้แล้ว)\n", err)
	} else {
		fmt.Printf("✅ ตั้งค่า margin type เป็น ISOLATED สำเร็จ\n")
	}

	// คำนวณขนาด position (200 USDT position value จาก 20 USDT margin × 10x leverage)
	marginSize := 20.0
	leverage := 10.0
	positionSize := marginSize * leverage // 20 × 10 = 200 USDT

	// ดึงราคาปัจจุบัน
	candlesticks, err := bot.binanceClient.GetCandlesticks(contract, "1m", 1)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงราคาปัจจุบัน: %v\n", err)
		return
	}

	if len(candlesticks) == 0 {
		fmt.Printf("❌ ไม่มีข้อมูลราคา\n")
		return
	}

	currentPrice, _ := strconv.ParseFloat(candlesticks[0].Close, 64)
	quantity := positionSize / currentPrice

	// ปรับ quantity precision ให้เหมาะสมกับ symbol
	adjustedQuantity, err := bot.binanceClient.AdjustQuantityPrecision(contract, quantity)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถปรับ quantity precision: %v\n", err)
		return
	}

	fmt.Printf("📊 ราคา: %.6f, Quantity ต้นฉบับ: %.6f, ปรับแล้ว: %s\n", currentPrice, quantity, adjustedQuantity)

	// แปลง string กลับเป็น float64 สำหรับ CreateMarketOrder
	finalQuantity, _ := strconv.ParseFloat(adjustedQuantity, 64)

	// สร้าง market order
	orderResponse, err := bot.binanceClient.CreateMarketOrder(contract, side, finalQuantity)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเปิด position ได้: %v\n", err)
		return
	}

	fmt.Printf("✅ เปิด position สำเร็จ! Order ID: %d\n", orderResponse.OrderId)
}

// analyzeClosePosition วิเคราะห์ว่าควรปิด position หรือไม่
func (bot *TradingBot) analyzeClosePosition(position binance.Position) {
	fmt.Printf("🔍 วิเคราะห์การปิด position: %s\n", position.Symbol)

	// ตรวจสอบว่า position มี size จริงหรือไม่
	positionAmt, err := strconv.ParseFloat(position.PositionAmt, 64)
	if err != nil || positionAmt == 0.0 {
		fmt.Printf("⚠️ ข้าม position %s - ไม่มี size จริง (amount: %s)\n", position.Symbol, position.PositionAmt)
		return
	}

	// ดึงข้อมูล candlestick
	candlesticks, err := bot.binanceClient.GetCandlesticks(position.Symbol, "1h", 288)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึงข้อมูล candlestick ได้: %v\n", err)
		return
	}

	// แปลงข้อมูลสำหรับ AI
	ohlcvData := bot.convertToOHLCV(candlesticks)

	// สร้าง Position object สำหรับ AI client
	tradingPosition := &Position{
		Contract:      position.Symbol,
		Size:          int64(bot.parseFloat(position.PositionAmt)),
		EntryPrice:    bot.parseFloat(position.EntryPrice),
		MarkPrice:     bot.parseFloat(position.MarkPrice),
		UnrealizedPnl: bot.parseFloat(position.UnrealizedProfit),
		Leverage:      bot.parseFloat(position.Leverage),
	}

	// ใช้ AI วิเคราะห์
	decision, err := bot.aiClient.AnalyzeClosePosition(tradingPosition, ohlcvData)
	if err != nil {
		fmt.Printf("❌ AI analysis error: %v\n", err)
		return
	}

	fmt.Printf("🤖 AI Close Analysis: %s\n", decision.Action)

	// ตัดสินใจปิด position (ไม่ใช้ confidence แล้ว)
	if strings.ToLower(decision.Action) == "close" {
		fmt.Printf("🔥 AI แนะนำให้ปิด position %s\n", position.Symbol)
		bot.closePosition(position)
	} else {
		fmt.Printf("⏳ รอสัญญาณปิด position %s ต่อไป\n", position.Symbol)
	}
}

// closePosition ปิด position
func (bot *TradingBot) closePosition(position binance.Position) {
	fmt.Printf("🔥 กำลังปิด position: %s\n", position.Symbol)

	orderResponse, err := bot.binanceClient.ClosePosition(position.Symbol)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถปิด position ได้: %v\n", err)
		return
	}

	fmt.Printf("✅ ปิด position สำเร็จ! Order ID: %d\n", orderResponse.OrderId)
}

// parseFloat helper method to convert string to float64
func (bot *TradingBot) parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

// checkEMAFilter ตรวจสอบสัญญาณ EMA เพื่อกรองเหรียญ
func (bot *TradingBot) checkEMAFilter(candlesticks []binance.Candlestick) string {
	if len(candlesticks) < 200 { // ต้องมีข้อมูลพอสำหรับคำนวณ EMA200
		return "HOLD"
	}

	// แปลงราคาปิด
	closes := make([]float64, len(candlesticks))
	for i, candle := range candlesticks {
		closes[i], _ = strconv.ParseFloat(candle.Close, 64)
	}

	// คำนวณ EMA
	ema20 := bot.calculateEMA(closes, 20)
	ema50 := bot.calculateEMA(closes, 50)
	ema100 := bot.calculateEMA(closes, 100)
	ema200 := bot.calculateEMA(closes, 200)

	if len(ema20) == 0 || len(ema50) == 0 || len(ema100) == 0 || len(ema200) == 0 {
		return "HOLD"
	}

	// ใช้ค่า EMA ล่าสุด
	currentEMA20 := ema20[len(ema20)-1]
	currentEMA50 := ema50[len(ema50)-1]
	currentEMA100 := ema100[len(ema100)-1]
	currentEMA200 := ema200[len(ema200)-1]

	// ตรวจสอบ BULLISH alignment: EMA20 > EMA50 > EMA100 > EMA200
	if currentEMA20 > currentEMA50 && currentEMA50 > currentEMA100 && currentEMA100 > currentEMA200 {
		return "LONG"
	}

	// ตรวจสอบ BEARISH alignment: EMA20 < EMA50 < EMA100 < EMA200
	if currentEMA20 < currentEMA50 && currentEMA50 < currentEMA100 && currentEMA100 < currentEMA200 {
		return "SHORT"
	}

	// SIDEWAYS: ไม่เรียงตัวกัน → ข้าม AI
	return "HOLD"
}

// calculateEMA คำนวณ Exponential Moving Average
func (bot *TradingBot) calculateEMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	var ema []float64
	multiplier := 2.0 / float64(period+1)

	// เริ่มต้นด้วยราคาแรก
	ema = append(ema, prices[0])

	// คำนวณ EMA สำหรับจุดต่อไป
	for i := 1; i < len(prices); i++ {
		currentEMA := (prices[i] * multiplier) + (ema[len(ema)-1] * (1 - multiplier))
		ema = append(ema, currentEMA)
	}

	return ema
}
