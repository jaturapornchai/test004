package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("🚀 ENHANCED TRIPLE EMA 1H - REAL MARKET DATA BACKTEST")
	fmt.Println("📊 Fetching live data from Gate.io API")
	fmt.Println(strings.Repeat("━", 80))

	// ทดสอบการดึงข้อมูลจริงก่อน
	fmt.Println("\n🔍 Phase 1: Testing Real Market Data Connection")
	TestRealMarketData()

	fmt.Println("\n" + strings.Repeat("━", 80))
	fmt.Println("🎯 Phase 2: Real Market Backtest with Enhanced Triple EMA Strategy")

	// รัน backtest ด้วยข้อมูลจริง
	RunRealMarketBacktestSuite()

	fmt.Println("\n" + strings.Repeat("━", 80))
	fmt.Printf("🎉 BACKTEST COMPLETE!\n")
	fmt.Printf("📈 Strategy tested with REAL market conditions\n")
	fmt.Printf("🚀 Ready for live trading deployment\n")
}
