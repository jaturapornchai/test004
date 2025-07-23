package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gateio-trading-bot/internal/config"
	"gateio-trading-bot/internal/gateio"
	"gateio-trading-bot/internal/trading"
)

// TripleEMA1HResult ผลลัพธ์ Triple EMA 1H
type TripleEMA1HResult struct {
	Symbol         string    `json:"symbol"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Duration       string    `json:"duration"`
	InitialCapital float64   `json:"initial_capital"`
	FinalCapital   float64   `json:"final_capital"`
	TotalReturn    float64   `json:"total_return"`
	TotalReturnPct float64   `json:"total_return_pct"`
	TotalTrades    int       `json:"total_trades"`
	WinningTrades  int       `json:"winning_trades"`
	LosingTrades   int       `json:"losing_trades"`
	WinRate        float64   `json:"win_rate"`
	AvgTradeDays   float64   `json:"avg_trade_days"`
	Strategy       string    `json:"strategy"`
	TimeFrame      string    `json:"timeframe"`
	LookbackDays   int       `json:"lookback_days"`
	Candles        int       `json:"candles"`
}

func main() {
	fmt.Println("🚀 Gate.io Triple EMA 1H Strategy Backtester (144 Candles)")
	fmt.Println("=" * 80)

	// โหลดการตั้งค่า
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("❌ Error loading .env: %v", err)
	}

	// สร้าง Gate.io client
	client := gateio.NewClient(cfg.APIKey, cfg.APISecret, cfg.BaseURL)

	// สัญลักษณ์ที่จะทดสอบ
	symbols := []string{"SOL_USDT", "BTC_USDT", "ETH_USDT"}
	results := make([]TripleEMA1HResult, 0)

	fmt.Printf("📊 Testing %d symbols with 1H timeframe (144 candles = 6 days)\n\n", len(symbols))

	// ทดสอบแต่ละสัญลักษณ์
	for i, symbol := range symbols {
		fmt.Printf("🔍 [%d/%d] Testing %s...\n", i+1, len(symbols), symbol)

		result := test1HSymbol(symbol, client)
		if result != nil {
			results = append(results, *result)
			time.Sleep(500 * time.Millisecond) // หน่วงเวลาเพื่อไม่ให้ API overload
		}
	}

	// แสดงผลลัพธ์สรุป
	printTripleEMA1HSummary(results)

	// บันทึกผลลัพธ์เป็น JSON
	saveTripleEMA1HResults(results)

	fmt.Println("\n✅ Triple EMA 1H Backtesting สำเร็จ!")
}

// test1HSymbol ทดสอบสัญลักษณ์เดียวด้วย Triple EMA 1H
func test1HSymbol(symbol string, client *gateio.Client) *TripleEMA1HResult {
	fmt.Printf("  📈 กำลังดึงข้อมูล %s (1H, 150 candles)...\n", symbol)

	// ดึงข้อมูล 150 candles (เผื่อไว้สำหรับ indicators)
	endTime := time.Now()
	startTime := endTime.Add(-150 * time.Hour) // 150 ชั่วโมงย้อนหลัง

	klines, err := client.GetKlines(symbol, "1h", int(startTime.Unix()), int(endTime.Unix()), 150)
	if err != nil {
		fmt.Printf("  ❌ Error fetching %s: %v\n", symbol, err)
		return nil
	}

	if len(klines) < 144 {
		fmt.Printf("  ❌ %s: ข้อมูลไม่ครบ (ต้องการ 144, ได้ %d)\n", symbol, len(klines))
		return nil
	}

	fmt.Printf("  ✅ ได้ข้อมูล %d candles\n", len(klines))

	// แปลงเป็น OHLCV format
	ohlcvData := make([]trading.OHLCV, len(klines))
	for i, kline := range klines {
		ohlcvData[i] = trading.OHLCV{
			Timestamp: kline.Timestamp,
			Open:      kline.Open,
			High:      kline.High,
			Low:       kline.Low,
			Close:     kline.Close,
			Volume:    kline.Volume,
		}
	}

	// สร้าง backtester
	backtester := trading.NewBacktester(
		symbol,
		10000.0, // เงินทุนเริ่มต้น $10,000
		ohlcvData,
		time.Unix(klines[0].Timestamp, 0),
		time.Unix(klines[len(klines)-1].Timestamp, 0),
	)

	// รัน Triple EMA 1H strategy
	fmt.Printf("  🔄 กำลังรัน Triple EMA 1H Strategy...\n")
	result := backtester.RunTripleEMA1HStrategy()

	// คำนวณระยะเวลาเฉลี่ยของ trade
	var totalDuration time.Duration
	completedTrades := 0

	for _, trade := range result.Trades {
		if !trade.ExitTime.IsZero() {
			totalDuration += trade.Duration
			completedTrades++
		}
	}

	var avgTradeDays float64
	if completedTrades > 0 {
		avgTradeDays = totalDuration.Hours() / float64(completedTrades) / 24
	}

	// สร้างผลลัพธ์
	triple1HResult := &TripleEMA1HResult{
		Symbol:         symbol,
		StartDate:      result.StartDate,
		EndDate:        result.EndDate,
		Duration:       fmt.Sprintf("%.1f days", result.EndDate.Sub(result.StartDate).Hours()/24),
		InitialCapital: result.InitialCapital,
		FinalCapital:   result.FinalCapital,
		TotalReturn:    result.TotalReturn,
		TotalReturnPct: result.TotalReturnPct,
		TotalTrades:    result.TotalTrades,
		WinningTrades:  result.WinningTrades,
		LosingTrades:   result.LosingTrades,
		WinRate:        result.WinRate,
		AvgTradeDays:   avgTradeDays,
		Strategy:       "Triple EMA (Fast=9, Mid=21, Slow=50)",
		TimeFrame:      "1H",
		LookbackDays:   6, // 144 candles × 1H = 6 days
		Candles:        144,
	}

	// แสดงผลลัพธ์
	fmt.Printf("  📊 %s Results:\n", symbol)
	fmt.Printf("    💰 Return: $%.2f → $%.2f (%.2f%%)\n",
		result.InitialCapital, result.FinalCapital, result.TotalReturnPct)
	fmt.Printf("    📈 Trades: %d (Win: %d, Loss: %d, Rate: %.1f%%)\n",
		result.TotalTrades, result.WinningTrades, result.LosingTrades, result.WinRate)
	fmt.Printf("    ⏱️  Avg Trade: %.1f days\n", avgTradeDays)

	return triple1HResult
}

// printTripleEMA1HSummary แสดงสรุปผลลัพธ์ Triple EMA 1H
func printTripleEMA1HSummary(results []TripleEMA1HResult) {
	fmt.Println("\n🎯 สรุปผลลัพธ์ Triple EMA 1H Strategy (144 Candles)")
	fmt.Println("=" * 80)

	if len(results) == 0 {
		fmt.Println("❌ ไม่มีข้อมูลผลลัพธ์")
		return
	}

	// แสดงผลแต่ละสัญลักษณ์
	fmt.Printf("%-12s %-12s %-12s %-8s %-8s %-10s %-8s\n",
		"Symbol", "Return", "Return%", "Trades", "Win%", "Avg Days", "Performance")
	fmt.Println("-" * 80)

	var totalInitial, totalFinal float64
	var bestReturn, worstReturn float64
	var bestSymbol, worstSymbol string

	for i, result := range results {
		totalInitial += result.InitialCapital
		totalFinal += result.FinalCapital

		performance := "📊"
		if result.TotalReturnPct > 20 {
			performance = "🚀"
		} else if result.TotalReturnPct > 10 {
			performance = "📈"
		} else if result.TotalReturnPct < -10 {
			performance = "📉"
		} else if result.TotalReturnPct < 0 {
			performance = "⚠️"
		}

		if i == 0 || result.TotalReturnPct > bestReturn {
			bestReturn = result.TotalReturnPct
			bestSymbol = result.Symbol
		}
		if i == 0 || result.TotalReturnPct < worstReturn {
			worstReturn = result.TotalReturnPct
			worstSymbol = result.Symbol
		}

		fmt.Printf("%-12s $%-10.2f %+8.2f%% %-8d %6.1f%% %8.1f   %-8s\n",
			result.Symbol,
			result.TotalReturn,
			result.TotalReturnPct,
			result.TotalTrades,
			result.WinRate,
			result.AvgTradeDays,
			performance)
	}

	fmt.Println("-" * 80)

	// สรุปรวม
	totalReturnPct := ((totalFinal - totalInitial) / totalInitial) * 100
	avgReturn := totalReturnPct / float64(len(results))

	fmt.Printf("📊 สรุปรวม: $%.2f → $%.2f (%+.2f%%)\n",
		totalInitial, totalFinal, totalReturnPct)
	fmt.Printf("📈 เฉลี่ย: %.2f%% per symbol\n", avgReturn)
	fmt.Printf("🥇 ดีที่สุด: %s (%+.2f%%)\n", bestSymbol, bestReturn)
	fmt.Printf("🥉 แย่ที่สุด: %s (%+.2f%%)\n", worstSymbol, worstReturn)

	// คำแนะนำ
	fmt.Println("\n💡 คำแนะนำ:")
	if avgReturn > 15 {
		fmt.Println("   🚀 กลยุทธ์ Triple EMA 1H ให้ผลลัพธ์ยอดเยี่ยม!")
	} else if avgReturn > 5 {
		fmt.Println("   📈 กลยุทธ์ Triple EMA 1H ให้ผลลัพธ์ดี")
	} else if avgReturn > 0 {
		fmt.Println("   📊 กลยุทธ์ Triple EMA 1H ให้ผลลัพธ์เป็นบวก")
	} else {
		fmt.Println("   ⚠️  กลยุทธ์ Triple EMA 1H ควรปรับปรุง")
	}

	fmt.Printf("   ⏰ Time Frame: 1H (144 candles = 6 days lookback)\n")
	fmt.Printf("   🎯 Strategy: Triple EMA (9,21,50) + RSI + Volume + ATR\n")
	fmt.Printf("   💰 Risk: 1%% per trade, 1:2 R/R, 5x leverage\n")
}

// saveTripleEMA1HResults บันทึกผลลัพธ์เป็น JSON
func saveTripleEMA1HResults(results []TripleEMA1HResult) {
	filename := fmt.Sprintf("triple_ema_1h_results_%s.json",
		time.Now().Format("2006-01-02_15-04-05"))

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error marshaling results: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("❌ Error saving results: %v\n", err)
		return
	}

	fmt.Printf("💾 ผลลัพธ์บันทึกแล้ว: %s\n", filename)
}
