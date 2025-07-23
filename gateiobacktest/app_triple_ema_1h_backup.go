package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gateio-trading-bot/internal/trading"

	"github.com/joho/godotenv"
)

// TripleEMA1HResult ผลลัพธ์การทดสอบ Triple EMA 1H
type TripleEMA1HResult struct {
	Symbol         string                      `json:"symbol"`
	Timeframe      string                      `json:"timeframe"`
	DataPoints     int                         `json:"data_points"`
	StartDate      time.Time                   `json:"start_date"`
	EndDate        time.Time                   `json:"end_date"`
	InitialCapital float64                     `json:"initial_capital"`
	FinalCapital   float64                     `json:"final_capital"`
	TotalReturn    float64                     `json:"total_return"`
	TotalReturnPct float64                     `json:"total_return_pct"`
	TotalTrades    int                         `json:"total_trades"`
	WinRate        float64                     `json:"win_rate"`
	MaxDrawdown    float64                     `json:"max_drawdown"`
	AvgTradeDuration time.Duration             `json:"avg_trade_duration"`
	Strategy       string                      `json:"strategy"`
	Result         *trading.BacktestResult     `json:"result"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "multi" {
		runMultiSymbol1HBacktest()
		return
	}
	
	runSingle1HBacktest()
}

func runSingle1HBacktest() {
	fmt.Println("🚀 Triple EMA 1H Strategy Backtest")
	fmt.Println("📊 Time Frame: 1 Hour (144 candles lookback)")
	fmt.Println("🎯 Strategy: Triple EMA Momentum with Real Data")

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("⚠️ ไม่สามารถโหลดไฟล์ .env: %v", err)
		return
	}

	// ทดสอบ SOL_USDT
	symbol := "SOL_USDT"
	initialCapital := 1000.0
	
	result := test1HSymbol(symbol, initialCapital)
	if result != nil {
		print1HResults(result)
		save1HResults(result)
	}
}

func runMultiSymbol1HBacktest() {
	fmt.Println("🚀 Multi-Symbol Triple EMA 1H Strategy Comparison")
	fmt.Println("📊 Time Frame: 1 Hour (144 candles each)")
	fmt.Println("🎯 Testing: SOL, BTC, ETH with Real 1H Data")

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("⚠️ ไม่สามารถโหลดไฟล์ .env: %v", err)
		return
	}

	symbols := []string{"SOL_USDT", "BTC_USDT", "ETH_USDT"}
	initialCapital := 1000.0
	var results []TripleEMA1HResult

	fmt.Printf("\n%s\n", strings.Repeat("=", 80))

	for i, symbol := range symbols {
		fmt.Printf("\n🔍 [%d/%d] กำลังทดสอบ %s (1H timeframe)...\n", i+1, len(symbols), symbol)
		
		result := test1HSymbol(symbol, initialCapital)
		if result != nil {
			results = append(results, *result)
		}
		
		// พักระหว่างการเรียก API
		if i < len(symbols)-1 {
			fmt.Printf("⏳ รอ 3 วินาที เพื่อป้องกัน rate limit...\n")
			time.Sleep(3 * time.Second)
		}
	}

	// สรุปผลการเปรียบเทียบ
	print1HComparison(results)
	save1HComparison(results)
}

func test1HSymbol(symbol string, initialCapital float64) *TripleEMA1HResult {
	fmt.Printf("📊 ดำเนินการทดสอบ %s (1H timeframe)...\n", symbol)
	
	// สร้าง data fetcher สำหรับ 1H data (144 candles = 6 วัน)
	fetcher := trading.NewDataFetcher(symbol, 6, "1h")
	
	// ดึงหรือโหลดข้อมูล
	historicalData, err := fetcher.FetchOrLoadData()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถโหลดข้อมูล %s: %v\n", symbol, err)
		return nil
	}
	
	// ตรวจสอบว่ามีข้อมูลครบ 144 candles
	if len(historicalData) < 144 {
		fmt.Printf("⚠️ ข้อมูล %s ไม่ครบ 144 candles (มี %d candles)\n", symbol, len(historicalData))
		return nil
	}
	
	// ใช้เฉพาะ 144 candles ล่าสุด
	if len(historicalData) > 144 {
		historicalData = historicalData[len(historicalData)-144:]
	}
	
	// แสดงข้อมูลสรุป
	startPrice := historicalData[0].Close
	endPrice := historicalData[len(historicalData)-1].Close
	priceChange := ((endPrice - startPrice) / startPrice) * 100
	
	startTime := time.Unix(historicalData[0].Timestamp, 0)
	endTime := time.Unix(historicalData[len(historicalData)-1].Timestamp, 0)
	
	fmt.Printf("📈 %s (1H): $%.2f -> $%.2f (%.2f%%) | %s ถึง %s\n", 
		symbol, startPrice, endPrice, priceChange,
		startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))
	fmt.Printf("📊 ข้อมูล: %d แท่งเทียน 1H (%.1f วัน)\n", len(historicalData), float64(len(historicalData))/24)
	
	// สร้าง backtester
	backtester, err := trading.NewBacktesterSimple(symbol, 6, initialCapital)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถสร้าง backtester สำหรับ %s: %v\n", symbol, err)
		return nil
	}
	
	// โหลดข้อมูล
	backtester.LoadHistoricalData(historicalData)
	
	// รัน Triple EMA 1H strategy
	fmt.Printf("🎯 เริ่มต้น Triple EMA 1H Strategy สำหรับ %s\n", symbol)
	result := backtester.RunTripleEMA1HStrategy()
	
	// คำนวณ average trade duration
	var totalDuration time.Duration
	completedTrades := 0
	for _, trade := range result.Trades {
		if !trade.ExitTime.IsZero() {
			totalDuration += trade.Duration
			completedTrades++
		}
	}
	
	var avgDuration time.Duration
	if completedTrades > 0 {
		avgDuration = totalDuration / time.Duration(completedTrades)
	}
	
	// สร้างผลลัพธ์
	return &TripleEMA1HResult{
		Symbol:           symbol,
		Timeframe:        "1H",
		DataPoints:       len(historicalData),
		StartDate:        startTime,
		EndDate:          endTime,
		InitialCapital:   initialCapital,
		FinalCapital:     result.FinalCapital,
		TotalReturn:      result.TotalReturn,
		TotalReturnPct:   result.TotalReturnPct,
		TotalTrades:      result.TotalTrades,
		WinRate:          result.WinRate,
		MaxDrawdown:      result.MaxDrawdownPct,
		AvgTradeDuration: avgDuration,
		Strategy:         "Triple EMA 1H Momentum",
		Result:           result,
	}
}

func print1HResults(result *TripleEMA1HResult) {
	separator := strings.Repeat("=", 80)
	
	fmt.Println("\n" + separator)
	fmt.Printf("📊 ผลลัพธ์ Triple EMA 1H Strategy: %s\n", result.Symbol)
	fmt.Println(separator)
	
	fmt.Printf("⏰ Time Frame: %s (%d candles, %.1f วัน)\n", 
		result.Timeframe, result.DataPoints, float64(result.DataPoints)/24)
	fmt.Printf("📅 ช่วงเวลา: %s ถึง %s\n", 
		result.StartDate.Format("2006-01-02 15:04"), 
		result.EndDate.Format("2006-01-02 15:04"))
	fmt.Printf("💰 เงินทุนเริ่มต้น: $%.2f\n", result.InitialCapital)
	fmt.Printf("💰 เงินทุนสุดท้าย: $%.2f\n", result.FinalCapital)
	fmt.Printf("📈 กำไร/ขาดทุนรวม: $%.2f (%.2f%%)\n", result.TotalReturn, result.TotalReturnPct)
	fmt.Printf("📊 จำนวนการเทรดทั้งหมด: %d\n", result.TotalTrades)
	fmt.Printf("🎯 อัตราชนะ: %.2f%%\n", result.WinRate)
	fmt.Printf("📉 Drawdown สูงสุด: %.2f%%\n", result.MaxDrawdown)
	fmt.Printf("⏱️ ระยะเวลาเทรดเฉลี่ย: %v\n", result.AvgTradeDuration.Round(time.Minute))
	
	fmt.Println(separator)
	fmt.Println("🚀 Triple EMA 1H Strategy Features:")
	fmt.Println("• Time Frame: 1 ชั่วโมง (144 candles lookback)")
	fmt.Println("• EMAs: Fast(9), Mid(21), Slow(50)")
	fmt.Println("• Indicators: RSI(14), ATR(14), Volume(20)")
	fmt.Println("• Risk Management: 1% per trade, 1:2 R/R")
	fmt.Println("• Leverage: 5x Future Trading")
	fmt.Println(separator)
}

func print1HComparison(results []TripleEMA1HResult) {
	separator := strings.Repeat("=", 100)
	
	fmt.Println("\n" + separator)
	fmt.Println("📊 ผลลัพธ์การเปรียบเทียบ Triple EMA 1H Strategy")
	fmt.Println(separator)
	
	fmt.Printf("%-10s | %-8s | %-12s | %-12s | %-10s | %-8s | %-8s | %-12s\n",
		"Symbol", "Candles", "Initial($)", "Final($)", "Return(%)", "Trades", "Win(%)", "Avg Duration")
	fmt.Println(strings.Repeat("-", 100))
	
	var bestReturn float64 = -999999
	var worstReturn float64 = 999999
	var bestSymbol, worstSymbol string
	var totalReturn float64
	
	for _, result := range results {
		fmt.Printf("%-10s | %-8d | %-12.2f | %-12.2f | %-10.2f | %-8d | %-8.2f | %-12v\n",
			result.Symbol,
			result.DataPoints,
			result.InitialCapital,
			result.FinalCapital,
			result.TotalReturnPct,
			result.TotalTrades,
			result.WinRate,
			result.AvgTradeDuration.Round(time.Minute))
		
		totalReturn += result.TotalReturnPct
		
		if result.TotalReturnPct > bestReturn {
			bestReturn = result.TotalReturnPct
			bestSymbol = result.Symbol
		}
		
		if result.TotalReturnPct < worstReturn {
			worstReturn = result.TotalReturnPct
			worstSymbol = result.Symbol
		}
	}
	
	avgReturn := totalReturn / float64(len(results))
	
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("🏆 ผลงานดีที่สุด: %s (%.2f%%)\n", bestSymbol, bestReturn)
	fmt.Printf("🔻 ผลงานต่ำสุด: %s (%.2f%%)\n", worstSymbol, worstReturn)
	fmt.Printf("📊 ผลตอบแทนเฉลี่ย: %.2f%%\n", avgReturn)
	
	fmt.Println("\n" + separator)
	fmt.Println("🚀 Triple EMA 1H Strategy Summary:")
	fmt.Println("• Time Frame: 1 ชั่วโมง (144 candles = 6 วัน)")
	fmt.Println("• Strategy: Triple EMA + RSI + Volume + ATR")
	fmt.Println("• Risk Management: 1% per trade, 5x leverage")
	fmt.Println("• Data Source: Real market data via API")
	fmt.Println(separator)
}

func save1HResults(result *TripleEMA1HResult) {
	filename := fmt.Sprintf("triple_ema_1h_%s_result_%s.json", 
		result.Symbol, time.Now().Format("20060102_150405"))
	
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("❌ ไม่สามารถสร้าง JSON: %v\n", err)
		return
	}
	
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถบันทึกไฟล์: %v\n", err)
		return
	}
	
	fmt.Printf("💾 บันทึกผลลัพธ์: %s\n", filename)
}

func save1HComparison(results []TripleEMA1HResult) {
	summary := map[string]interface{}{
		"test_date":    time.Now(),
		"strategy":     "Triple EMA 1H Momentum",
		"timeframe":    "1H",
		"data_points":  144,
		"period_days":  6,
		"results":      results,
	}
	
	// หาผลงานดีที่สุดและต่ำสุด
	var bestReturn float64 = -999999
	var worstReturn float64 = 999999
	var totalReturn float64
	
	for _, result := range results {
		totalReturn += result.TotalReturnPct
		
		if result.TotalReturnPct > bestReturn {
			bestReturn = result.TotalReturnPct
			summary["best_symbol"] = result.Symbol
			summary["best_return"] = bestReturn
		}
		
		if result.TotalReturnPct < worstReturn {
			worstReturn = result.TotalReturnPct
			summary["worst_symbol"] = result.Symbol
			summary["worst_return"] = worstReturn
		}
	}
	
	summary["avg_return"] = totalReturn / float64(len(results))
	
	// บันทึกไฟล์
	filename := fmt.Sprintf("triple_ema_1h_comparison_%s.json", 
		time.Now().Format("20060102_150405"))
	
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Printf("❌ ไม่สามารถสร้าง JSON: %v\n", err)
		return
	}
	
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถบันทึกไฟล์: %v\n", err)
		return
	}
	
	fmt.Printf("💾 บันทึกผลลัพธ์การเปรียบเทียบ: %s\n", filename)
}
