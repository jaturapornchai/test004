package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("🚀 VOLUME BREAKOUT + EMA STRATEGY TEST")
	fmt.Println("📊 Enhanced strategy combining volume breakouts with EMA trend confirmation")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	symbols := []string{"BTC_USDT", "ETH_USDT", "SOL_USDT"}

	for _, symbol := range symbols {
		fmt.Printf("\n🔍 Testing Symbol: %s\n", symbol)
		fmt.Println("--------------------------------------------------")

		strategy := NewVolumeBreakoutEMAStrategy(symbol, 10000.0)
		strategy.ExecuteBacktest()

		time.Sleep(1 * time.Second) // Rate limiting
	}

	fmt.Println("\n✅ Volume Breakout + EMA strategy testing complete!")
}
