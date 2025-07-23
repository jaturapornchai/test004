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

// MultiSymbolResult ผลลัพธ์การทดสอบหลายเหรียญ
type MultiSymbolResult struct {
	Symbol         string                  `json:"symbol"`
	InitialCapital float64                 `json:"initial_capital"`
	FinalCapital   float64                 `json:"final_capital"`
	TotalReturn    float64                 `json:"total_return"`
	TotalReturnPct float64                 `json:"total_return_pct"`
	TotalTrades    int                     `json:"total_trades"`
	WinRate        float64                 `json:"win_rate"`
	MaxDrawdown    float64                 `json:"max_drawdown"`
	Result         *trading.BacktestResult `json:"result"`
}

// ComparisonSummary สรุปการเปรียบเทียบ
type ComparisonSummary struct {
	TestDate    time.Time           `json:"test_date"`
	Strategy    string              `json:"strategy"`
	Period      string              `json:"period"`
	Results     []MultiSymbolResult `json:"results"`
	BestSymbol  string              `json:"best_symbol"`
	WorstSymbol string              `json:"worst_symbol"`
	AvgReturn   float64             `json:"avg_return"`
}

func main() {
	runMultiSymbolComparison()
}

func runMultiSymbolComparison() {
	fmt.Println("🚀 เริ่มต้น Multi-Symbol Triple EMA Strategy Comparison")
	fmt.Println("📈 ทดสอบ SOL, BTC, ETH ด้วย Triple EMA Strategy")
	fmt.Println("⏰ ข้อมูลย้อนหลัง 12 เดือน | Time Frame: 15 นาที")
	fmt.Println("🎯 เป้าหมาย: หาเหรียญที่ให้ผลตอบแทนดีที่สุด")

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("⚠️ ไม่สามารถโหลดไฟล์ .env: %v (ใช้ค่าเริ่มต้น)", err)
	}

	symbols := []string{"SOL_USDT", "BTC_USDT", "ETH_USDT"}
	initialCapital := 1000.0
	var results []MultiSymbolResult

	fmt.Printf("\n%s\n", strings.Repeat("=", 80))

	for i, symbol := range symbols {
		fmt.Printf("\n🔍 [%d/%d] กำลังทดสอบ %s...\n", i+1, len(symbols), symbol)

		result := testSingleSymbol(symbol, initialCapital)
		if result != nil {
			results = append(results, *result)
		}

		// พักระหว่างการเรียก API
		if i < len(symbols)-1 {
			fmt.Printf("⏳ รอ 2 วินาที เพื่อป้องกัน rate limit...\n")
			time.Sleep(2 * time.Second)
		}
	}

	// สรุปผลการเปรียบเทียบ
	printComparisonResults(results)

	// บันทึกผลลัพธ์
	saveComparisonResults(results)
}

func testSingleSymbol(symbol string, initialCapital float64) *MultiSymbolResult {
	fmt.Printf("📊 ดำเนินการทดสอบ %s...\n", symbol)

	// สร้าง data fetcher
	fetcher := trading.NewDataFetcher(symbol, 365, "15m")

	// ดึงหรือโหลดข้อมูล
	historicalData, err := fetcher.FetchOrLoadData()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถโหลดข้อมูล %s: %v\n", symbol, err)
		return nil
	}

	// แสดงข้อมูลสรุป
	fetcher.GetSymbolInfo(historicalData)

	// สร้าง backtester
	backtester, err := trading.NewBacktester(symbol, 365, initialCapital, "")
	if err != nil {
		fmt.Printf("❌ ไม่สามารถสร้าง backtester สำหรับ %s: %v\n", symbol, err)
		return nil
	}

	// โหลดข้อมูล
	backtester.LoadHistoricalData(historicalData)

	// รัน Triple EMA strategy
	fmt.Printf("🎯 เริ่มต้น Triple EMA Strategy สำหรับ %s\n", symbol)
	result := backtester.RunNew15mStrategy()

	// สร้างผลลัพธ์
	return &MultiSymbolResult{
		Symbol:         symbol,
		InitialCapital: initialCapital,
		FinalCapital:   result.FinalCapital,
		TotalReturn:    result.TotalReturn,
		TotalReturnPct: result.TotalReturnPct,
		TotalTrades:    result.TotalTrades,
		WinRate:        result.WinRate,
		MaxDrawdown:    result.MaxDrawdownPct,
		Result:         result,
	}
}

func printComparisonResults(results []MultiSymbolResult) {
	separator := "=" * 100

	fmt.Println("\n" + separator)
	fmt.Println("📊 ผลลัพธ์การเปรียบเทียบ Triple EMA Strategy (12 เดือน)")
	fmt.Println(separator)

	fmt.Printf("%-10s | %-12s | %-12s | %-10s | %-8s | %-8s | %-10s\n",
		"Symbol", "Initial($)", "Final($)", "Return(%)", "Trades", "Win(%)", "Drawdown(%)")
	fmt.Println("-" * 100)

	var bestReturn float64 = -999999
	var worstReturn float64 = 999999
	var bestSymbol, worstSymbol string
	var totalReturn float64

	for _, result := range results {
		fmt.Printf("%-10s | %-12.2f | %-12.2f | %-10.2f | %-8d | %-8.2f | %-10.2f\n",
			result.Symbol,
			result.InitialCapital,
			result.FinalCapital,
			result.TotalReturnPct,
			result.TotalTrades,
			result.WinRate,
			result.MaxDrawdown)

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

	fmt.Println("-" * 100)
	fmt.Printf("🏆 ผลงานดีที่สุด: %s (%.2f%%)\n", bestSymbol, bestReturn)
	fmt.Printf("🔻 ผลงานต่ำสุด: %s (%.2f%%)\n", worstSymbol, worstReturn)
	fmt.Printf("📊 ผลตอบแทนเฉลี่ย: %.2f%%\n", avgReturn)

	fmt.Println("\n" + separator)
	fmt.Println("🚀 Triple EMA Strategy Features:")
	fmt.Println("• Time Frame: 15 นาที")
	fmt.Println("• Test Period: 12 เดือน")
	fmt.Println("• Strategy: Triple EMA + RSI + Volume Momentum")
	fmt.Println("• Leverage: 5x (Future Trading)")
	fmt.Println("• Risk per Trade: 1%")
	fmt.Println("• Risk:Reward Ratio: 1:2")
	fmt.Println(separator)
}

func saveComparisonResults(results []MultiSymbolResult) {
	summary := ComparisonSummary{
		TestDate: time.Now(),
		Strategy: "Triple EMA Momentum",
		Period:   "12 months",
		Results:  results,
	}

	// หาผลงานดีที่สุดและต่ำสุด
	var bestReturn float64 = -999999
	var worstReturn float64 = 999999
	var totalReturn float64

	for _, result := range results {
		totalReturn += result.TotalReturnPct

		if result.TotalReturnPct > bestReturn {
			bestReturn = result.TotalReturnPct
			summary.BestSymbol = result.Symbol
		}

		if result.TotalReturnPct < worstReturn {
			worstReturn = result.TotalReturnPct
			summary.WorstSymbol = result.Symbol
		}
	}

	summary.AvgReturn = totalReturn / float64(len(results))

	// บันทึกไฟล์
	filename := fmt.Sprintf("triple_ema_comparison_%s.json",
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

func runSingleBacktest() {
	fmt.Println("🚀 เริ่มต้น Single Symbol Triple EMA Strategy Backtesting")
	fmt.Println("📈 ทดสอบกลยุทธ์การเทรดใหม่ ย้อนหลัง 12 เดือน")
	fmt.Println("⏰ Time Frame: 15 นาที")
	fmt.Println("🎯 เป้าหมาย: กลยุทธ์ใหม่ที่เน้นความแม่นยำสูง")

	// โหลด environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("⚠️ ไม่สามารถโหลดไฟล์ .env: %v", err)
	}

	fmt.Println("✅ โหลด configuration สำเร็จ")

	// ทดสอบกับ SOL (default)
	result := testSingleSymbol("SOL_USDT", 1000.0)
	if result != nil {
		printSingleResult(result)
		saveSingleResult(result)
	}
}

func printSingleResult(result *MultiSymbolResult) {
	separator := "=" * 80

	fmt.Println("\n" + separator)
	fmt.Printf("📊 ผลลัพธ์ Triple EMA Strategy: %s (12 เดือน)\n", result.Symbol)
	fmt.Println(separator)

	fmt.Printf("💰 เงินทุนเริ่มต้น: $%.2f\n", result.InitialCapital)
	fmt.Printf("💰 เงินทุนสุดท้าย: $%.2f\n", result.FinalCapital)
	fmt.Printf("📈 กำไร/ขาดทุนรวม: $%.2f (%.2f%%)\n", result.TotalReturn, result.TotalReturnPct)
	fmt.Printf("📊 จำนวนการเทรดทั้งหมด: %d\n", result.TotalTrades)
	fmt.Printf("🎯 อัตราชนะ: %.2f%%\n", result.WinRate)
	fmt.Printf("📉 Drawdown สูงสุด: %.2f%%\n", result.MaxDrawdown)

	fmt.Println(separator)
}

func saveSingleResult(result *MultiSymbolResult) {
	filename := fmt.Sprintf("triple_ema_%s_result_%s.json",
		result.Symbol, time.Now().Format("20060102_150405"))

	data, err := json.MarshalIndent(result.Result, "", "  ")
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
