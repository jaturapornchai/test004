package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"gateio-trading-bot/internal/trading"

	"github.com/joho/godotenv"
)

func main() {
	// ตรวจสอบ argument สำหรับเลือกโหมด
	if len(os.Args) > 1 && os.Args[1] == "backtest" {
		runBacktest()
		return
	}

	// โหมดปกติ (live trading)
	runLiveTrading()
}

func runLiveTrading() {
	fmt.Println("🚀 เริ่มต้น Gate.io Trading Bot - Triple EMA Momentum + AI")

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("❌ ไม่สามารถโหลดไฟล์ .env ได้:", err)
	}

	// TODO: implement live trading
	fmt.Println("🚧 Live trading ยังไม่ได้ implement")
}

func runBacktest() {
	fmt.Println("🚀 เริ่มต้น New 15m Strategy Backtesting System")
	fmt.Println("📈 ทดสอบกลยุทธ์การเทรดใหม่ ย้อนหลัง 12 เดือน")
	fmt.Println("⏰ Time Frame: 15 นาที")
	fmt.Println("🎯 เป้าหมาย: กลยุทธ์ใหม่ที่เน้นความแม่นยำสูง")

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("❌ ไม่สามารถโหลดไฟล์ .env ได้:", err)
	}

	// ตรวจสอบ API credentials
	deepseekKey := os.Getenv("DEEPSEEK_API_KEY")
	if deepseekKey == "" {
		log.Fatal("❌ DEEPSEEK_API_KEY ไม่ได้ตั้งค่าใน .env")
	}

	fmt.Println("✅ โหลด configuration สำเร็จ")

	// สร้าง backtester สำหรับ 12 เดือน SOL_USDT
	backtester, err := trading.NewBacktester("SOL_USDT", 365, 1000.0, deepseekKey)
	if err != nil {
		log.Fatal("❌ ไม่สามารถสร้าง backtester ได้:", err)
	}

	// โหลดข้อมูลราคาจำลอง SOL 15m timeframe
	fmt.Println("📊 กำลังโหลดข้อมูลราคาย้อนหลัง 12 เดือน (15m timeframe)...")
	historicalData := generateSOL15mData(365)
	backtester.LoadHistoricalData(historicalData)

	// รัน new strategy backtest
	fmt.Println("🎯 เริ่มต้น New 15m Strategy Backtesting...")
	result := backtester.RunNew15mStrategy()

	// แสดงผลลัพธ์
	printNewStrategyResults(result)

	// บันทึกผลลัพธ์เป็นไฟล์ JSON
	saveResultsToFile(result)
}

// generateSOL15mData สร้างข้อมูลราคาจำลอง SOL_USDT 15m timeframe
func generateSOL15mData(days int) []trading.OHLCV {
	var data []trading.OHLCV

	// เริ่มต้นจากวันที่ที่กำหนด
	startTime := time.Now().AddDate(0, 0, -days)
	basePrice := 200.0 // ราคา SOL เริ่มต้น

	// สร้างข้อมูล 15m timeframe (วัน x 96 candles per day)
	totalCandles := days * 96 // 24h * 4 candles per hour
	for i := 0; i < totalCandles; i++ {
		timestamp := startTime.Add(time.Duration(i) * 15 * time.Minute)

		// สร้างความผันผวนแบบสุ่มสำหรับ SOL (ผันผวนมากกว่า BTC)
		volatility := 0.03 // 3% volatility per 15m (SOL มีความผันผวนสูง)
		change := (float64(i%96-48) / 48.0) * volatility * basePrice

		// เพิ่มรูปแบบความผันผวนของ crypto
		cycleTrend := math.Sin(float64(i)/500.0) * 50.0 // Long-term cycles
		shortTrend := math.Sin(float64(i)/50.0) * 15.0  // Short-term cycles

		// Random spikes for crypto volatility
		randomSpike := 0.0
		if i%100 == 0 {
			randomSpike = (float64(i%2)*2.0 - 1.0) * 30.0 // +/- $30 price spikes
		}

		// คำนวณราคา
		open := basePrice + change + cycleTrend + shortTrend + randomSpike
		high := open * (1.0 + volatility/2)
		low := open * (1.0 - volatility/2)
		close := open + (float64(i%15-7) / 15.0 * volatility * basePrice)

		// Volume แปรผันตามความผันผวน
		baseVolume := 800000.0
		volumeVariation := float64(i%100) * 10000.0
		volume := baseVolume + volumeVariation

		// เพิ่มให้มีช่วงราคาที่หลากหลาย
		if close < 0 {
			close = basePrice * 0.8
		}
		if open < 0 {
			open = basePrice * 0.8
		}

		data = append(data, trading.OHLCV{
			Timestamp: timestamp.Unix(),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		})
	}

	return data
}

// printNewStrategyResults แสดงผลลัพธ์กลยุทธ์ใหม่
func printNewStrategyResults(result *trading.BacktestResult) {
	separator := strings.Repeat("=", 80)

	fmt.Println("\n" + separator)
	fmt.Println("📊 ผลลัพธ์ New 15m Strategy Backtesting (12 เดือน)")
	fmt.Println(separator)

	fmt.Printf("💰 เงินทุนเริ่มต้น: $%.2f\n", result.InitialCapital)
	fmt.Printf("💰 เงินทุนสุดท้าย: $%.2f\n", result.FinalCapital)
	fmt.Printf("📈 กำไร/ขาดทุนรวม: $%.2f (%.2f%%)\n",
		result.FinalCapital-result.InitialCapital, result.TotalReturnPct)
	fmt.Printf("📊 จำนวนการเทรดทั้งหมด: %d\n", result.TotalTrades)
	fmt.Printf("✅ การเทรดที่ได้กำไร: %d\n", result.WinningTrades)
	fmt.Printf("❌ การเทรดที่ขาดทุน: %d\n", result.LosingTrades)
	fmt.Printf("🎯 อัตราชนะ: %.2f%%\n", result.WinRate)
	fmt.Printf("📉 Drawdown สูงสุด: %.2f%%\n", result.MaxDrawdownPct)

	fmt.Println("\n" + separator)
	fmt.Println("🔍 รายละเอียดการเทรด:")
	fmt.Println(separator)

	for i, trade := range result.Trades {
		if trade.NetPnL != 0 { // แสดงเฉพาะ trade ที่ปิดแล้ว
			status := "✅"
			if trade.NetPnL < 0 {
				status = "❌"
			}
			fmt.Printf("%s Trade %d: %s %.2f -> %.2f | PnL: $%.2f (%.2f%%)\n",
				status, i+1, trade.Side, trade.EntryPrice, trade.ExitPrice,
				trade.NetPnL, trade.PnLPct)
		}
	}

	if len(result.Trades) == 0 {
		fmt.Println("ไม่มีการเทรดในช่วงเวลาที่ทดสอบ")
	}

	fmt.Println("\n" + separator)
	fmt.Println("🚀 New 15m Strategy Features:")
	fmt.Println("• Time Frame: 15 นาที")
	fmt.Println("• Test Period: 12 เดือน")
	fmt.Println("• Triple EMA + RSI + Volume Strategy")
	fmt.Println("• Risk per Trade: 1%")
	fmt.Println("• Risk:Reward Ratio: 1:2")
	fmt.Println("• Smart Position Sizing")
	fmt.Println("• Multi-timeframe Confirmation")
	fmt.Println(separator)
}

// saveResultsToFile บันทึกผลลัพธ์เป็นไฟล์ JSON
func saveResultsToFile(result *trading.BacktestResult) {
	filename := fmt.Sprintf("backtest_15m_result_%s.json", time.Now().Format("20060102_150405"))

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("❌ ไม่สามารถ marshal JSON ได้: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถบันทึกไฟล์ได้: %v\n", err)
		return
	}

	fmt.Printf("💾 บันทึกผลลัพธ์ลงไฟล์: %s\n", filename)
}
