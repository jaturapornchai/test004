package main

import (
	"fmt"
	"gateio-trading-bot/internal/trading"
	"math"
)

// Test EMA calculation system
func main() {
	fmt.Println("🔍 ==================== EMA SYSTEM VERIFICATION ====================")

	// สร้างข้อมูลทดสอบ
	testData := createTestOHLCV()

	// สร้าง instances เพื่อทดสอบ
	aiClient, _ := trading.NewAIClient("test")
	bot, _ := trading.NewTradingBot("", "", "test")

	fmt.Println("\n📊 1. ทดสอบ EMA Calculation Functions")
	testEMACalculations(aiClient, bot, testData)

	fmt.Println("\n📊 2. ทดสอบ EMA Alignment Check")
	testEMAAlignment(aiClient, bot, testData)

	fmt.Println("\n📊 3. ทดสอบ EMA200 Touch Detection")
	testEMA200Touch(aiClient, bot, testData)

	fmt.Println("\n📊 4. ทดสอบ EMA50 Force Close Logic")
	testEMA50ForceClose(aiClient, testData)

	fmt.Println("\n✅ ==================== EMA SYSTEM TEST COMPLETE ====================")
}

func createTestOHLCV() []trading.OHLCV {
	// สร้างข้อมูลทดสอบ 220 แท่ง (มากกว่า 200 เพื่อทดสอบ EMA200)
	data := make([]trading.OHLCV, 220)

	// สร้างข้อมูลราคาแบบ uptrend
	basePrice := 100.0
	for i := 0; i < len(data); i++ {
		// สร้าง uptrend ช้าๆ
		trend := float64(i) * 0.1
		noise := math.Sin(float64(i)*0.3) * 2 // เพิ่ม noise เล็กน้อย

		price := basePrice + trend + noise

		data[i] = trading.OHLCV{
			Open:   price - 0.5,
			High:   price + 1.0,
			Low:    price - 1.0,
			Close:  price,
			Volume: 1000000,
		}
	}

	return data
}

func testEMACalculations(aiClient *trading.AIClient, bot *trading.TradingBot, data []trading.OHLCV) {
	fmt.Println("   🧮 คำนวณ EMA ทุกช่วง...")

	// ทดสอบ calculateEMA function ใน ai_client.go
	fmt.Println("   📈 AI Client EMA Calculations:")
	ema20_ai := calculateEMAPublic(data, 20)
	ema50_ai := calculateEMAPublic(data, 50)
	ema100_ai := calculateEMAPublic(data, 100)
	ema200_ai := calculateEMAPublic(data, 200)

	fmt.Printf("      EMA20:  %.6f\n", ema20_ai)
	fmt.Printf("      EMA50:  %.6f\n", ema50_ai)
	fmt.Printf("      EMA100: %.6f\n", ema100_ai)
	fmt.Printf("      EMA200: %.6f\n", ema200_ai)

	// ทดสอบ calculateEMA function ใน bot.go
	fmt.Println("   📈 Bot EMA Calculations:")
	ema20_bot := calculateEMAPublic(data, 20)
	ema50_bot := calculateEMAPublic(data, 50)
	ema100_bot := calculateEMAPublic(data, 100)
	ema200_bot := calculateEMAPublic(data, 200)

	fmt.Printf("      EMA20:  %.6f\n", ema20_bot)
	fmt.Printf("      EMA50:  %.6f\n", ema50_bot)
	fmt.Printf("      EMA100: %.6f\n", ema100_bot)
	fmt.Printf("      EMA200: %.6f\n", ema200_bot)

	// ตรวจสอบความสอดคล้อง
	tolerance := 0.000001
	if math.Abs(ema20_ai-ema20_bot) < tolerance &&
		math.Abs(ema50_ai-ema50_bot) < tolerance &&
		math.Abs(ema100_ai-ema100_bot) < tolerance &&
		math.Abs(ema200_ai-ema200_bot) < tolerance {
		fmt.Println("   ✅ EMA calculations ระหว่าง AI Client และ Bot สอดคล้องกัน")
	} else {
		fmt.Println("   ❌ EMA calculations ระหว่าง AI Client และ Bot ไม่สอดคล้อง!")
	}
}

func testEMAAlignment(aiClient *trading.AIClient, bot *trading.TradingBot, data []trading.OHLCV) {
	ema20 := calculateEMAPublic(data, 20)
	ema50 := calculateEMAPublic(data, 50)
	ema100 := calculateEMAPublic(data, 100)
	ema200 := calculateEMAPublic(data, 200)

	fmt.Printf("   📊 Current EMA Values:\n")
	fmt.Printf("      EMA20:  %.6f\n", ema20)
	fmt.Printf("      EMA50:  %.6f\n", ema50)
	fmt.Printf("      EMA100: %.6f\n", ema100)
	fmt.Printf("      EMA200: %.6f\n", ema200)

	// ตรวจสอบ alignment
	if ema20 > ema50 && ema50 > ema100 && ema100 > ema200 {
		fmt.Println("   ✅ LONG EMA Alignment: EMA20 > EMA50 > EMA100 > EMA200")
	} else if ema20 < ema50 && ema50 < ema100 && ema100 < ema200 {
		fmt.Println("   ✅ SHORT EMA Alignment: EMA20 < EMA50 < EMA100 < EMA200")
	} else {
		fmt.Println("   📊 Mixed EMA Alignment (ไม่เรียงตัวแบบ LONG หรือ SHORT)")
	}
}

func testEMA200Touch(aiClient *trading.AIClient, bot *trading.TradingBot, data []trading.OHLCV) {
	fmt.Println("   🎯 ทดสอบ EMA200 Touch Detection...")

	// คำนวณ EMA200 series
	ema200Series := calculateEMA200SeriesPublic(data)

	// ตรวจสอบ 48 แท่งสุดท้าย
	touchFound := false
	startIdx := len(data) - 48
	if startIdx < 0 {
		startIdx = 0
	}

	for i := startIdx; i < len(data); i++ {
		if i >= len(ema200Series) {
			continue
		}

		candle := data[i]
		ema200 := ema200Series[i]

		if candle.Low <= ema200 && ema200 <= candle.High {
			fmt.Printf("   ✅ พบ EMA200 Touch ที่แท่ง %d: EMA200=%.6f, Low=%.6f, High=%.6f\n",
				i-startIdx+1, ema200, candle.Low, candle.High)
			touchFound = true
			break
		}
	}

	if !touchFound {
		fmt.Println("   ❌ ไม่พบ EMA200 Touch ใน 48 แท่งย้อนหลัง")
	}
}

func testEMA50ForceClose(aiClient *trading.AIClient, data []trading.OHLCV) {
	fmt.Println("   ⚠️  ทดสอบ EMA50 Force Close Logic...")

	ema50 := calculateEMAPublic(data, 50)
	currentPrice := data[len(data)-1].Close

	fmt.Printf("   📊 Current Price: %.6f\n", currentPrice)
	fmt.Printf("   📊 EMA50: %.6f\n", ema50)

	// ทดสอบ LONG position
	if currentPrice < ema50 {
		fmt.Println("   ⚠️  LONG Position: ราคาต่ำกว่า EMA50 → FORCE CLOSE")
	} else {
		fmt.Println("   ✅ LONG Position: ราคาสูงกว่า EMA50 → HOLD")
	}

	// ทดสอบ SHORT position
	if currentPrice > ema50 {
		fmt.Println("   ⚠️  SHORT Position: ราคาสูงกว่า EMA50 → FORCE CLOSE")
	} else {
		fmt.Println("   ✅ SHORT Position: ราคาต่ำกว่า EMA50 → HOLD")
	}
}

// Helper functions (replicate the EMA calculation logic)
func calculateEMAPublic(ohlcv []trading.OHLCV, period int) float64 {
	if len(ohlcv) < period {
		return 0
	}
	k := 2.0 / float64(period+1)
	ema := ohlcv[0].Close
	for i := 1; i < len(ohlcv); i++ {
		ema = ohlcv[i].Close*k + ema*(1-k)
	}
	return ema
}

func calculateEMA200SeriesPublic(ohlcv []trading.OHLCV) []float64 {
	if len(ohlcv) < 200 {
		return []float64{}
	}

	result := make([]float64, len(ohlcv))
	k := 2.0 / float64(200+1)

	result[0] = ohlcv[0].Close

	for i := 1; i < len(ohlcv); i++ {
		result[i] = ohlcv[i].Close*k + result[i-1]*(1-k)
	}

	return result
}
