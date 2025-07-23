package main

import (
	"encoding/json"
	"fmt"
	"time"

	"gateio-trading-bot/internal/gateio"
)

func main() {
	fmt.Println("🚀 Enhanced 1H Data Manager - 144 Timeframes")
	fmt.Println("=" + string(make([]byte, 60)))

	// สร้าง client (demo mode)
	client := gateio.NewClient("", "", "https://api.gateio.ws")

	// ทดสอบหลาย symbols
	symbols := []string{"SOL_USDT", "BTC_USDT", "ETH_USDT"}

	for i, symbol := range symbols {
		fmt.Printf("\n🔍 [%d/%d] Testing Enhanced Data Manager: %s\n", i+1, len(symbols), symbol)
		testEnhancedDataManager(symbol, client)

		if i < len(symbols)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Println("\n✅ Enhanced Data Manager Testing Complete!")
}

func testEnhancedDataManager(symbol string, client *gateio.Client) {
	// สร้าง Enhanced Data Manager
	dm := NewEnhanced1HDataManager(symbol, client)

	// โหลดข้อมูลย้อนหลัง 144 timeframes
	if err := dm.LoadHistoricalData(); err != nil {
		fmt.Printf("❌ Error loading data: %v\n", err)
		return
	}

	// ดึงการวิเคราะห์ล่าสุด
	analysis := dm.GetLatestAnalysis()
	if analysis == nil {
		fmt.Printf("❌ No analysis available\n")
		return
	}

	// แสดงผลลัพธ์
	fmt.Printf("📊 Latest Analysis:\n")
	fmt.Printf("   💰 Price: $%.2f\n", analysis.Price)
	fmt.Printf("   📈 EMA: Fast=%.2f, Mid=%.2f, Slow=%.2f\n",
		analysis.EMA9, analysis.EMA21, analysis.EMA50)
	fmt.Printf("   📊 RSI: %.1f, ATR: %.2f\n", analysis.RSI, analysis.ATR)
	fmt.Printf("   🔊 Volume: %.0f (MA: %.0f)\n", analysis.CurrentVol, analysis.VolumeMA)
	fmt.Printf("   🎯 Signal: %s (%.1f%%)\n", analysis.Signal, analysis.Confidence)
	fmt.Printf("   ✅ Data Quality: %.1f%%\n", analysis.DataQuality)

	// Cache performance
	hits, misses, hitRate := dm.GetCacheStats()
	fmt.Printf("   🚀 Cache: %d hits, %d misses (%.1f%% hit rate)\n",
		hits, misses, hitRate)

	// Data summary
	summary := dm.GetDataSummary()
	fmt.Printf("📋 Data Summary:\n")
	summaryJSON, _ := json.MarshalIndent(summary, "   ", "  ")
	fmt.Printf("   %s\n", string(summaryJSON))

	// Signal analysis
	analyzeSignalStrength(analysis)
}

func analyzeSignalStrength(analysis *Enhanced1HAnalysis) {
	fmt.Printf("🎯 Signal Analysis:\n")

	// EMA Alignment
	emaAligned := false
	var alignment string
	if analysis.EMA9 > analysis.EMA21 && analysis.EMA21 > analysis.EMA50 {
		emaAligned = true
		alignment = "BULLISH (Fast > Mid > Slow)"
	} else if analysis.EMA9 < analysis.EMA21 && analysis.EMA21 < analysis.EMA50 {
		emaAligned = true
		alignment = "BEARISH (Fast < Mid < Slow)"
	} else {
		alignment = "MIXED (No clear alignment)"
	}

	fmt.Printf("   📈 EMA Alignment: %s %s\n", getIcon(emaAligned), alignment)

	// Price vs EMA
	priceVsEMA := ""
	if analysis.Price > analysis.EMA9 {
		priceVsEMA = "Above Fast EMA ↗️"
	} else if analysis.Price < analysis.EMA9 {
		priceVsEMA = "Below Fast EMA ↘️"
	} else {
		priceVsEMA = "At Fast EMA ➡️"
	}
	fmt.Printf("   💰 Price Position: %s\n", priceVsEMA)

	// RSI Analysis
	rsiAnalysis := ""
	if analysis.RSI > 70 {
		rsiAnalysis = "Overbought (>70) ⚠️"
	} else if analysis.RSI < 30 {
		rsiAnalysis = "Oversold (<30) ⚠️"
	} else if analysis.RSI >= 50 && analysis.RSI <= 70 {
		rsiAnalysis = "Bullish Zone (50-70) ✅"
	} else if analysis.RSI >= 30 && analysis.RSI < 50 {
		rsiAnalysis = "Bearish Zone (30-50) ✅"
	} else {
		rsiAnalysis = "Neutral Zone 📊"
	}
	fmt.Printf("   📊 RSI: %.1f - %s\n", analysis.RSI, rsiAnalysis)

	// Volume Analysis
	volumeRatio := analysis.CurrentVol / analysis.VolumeMA
	volumeAnalysis := ""
	if volumeRatio > 2.0 {
		volumeAnalysis = "Very High Volume (>2x) 🔥"
	} else if volumeRatio > 1.5 {
		volumeAnalysis = "High Volume (>1.5x) 📈"
	} else if volumeRatio > 1.2 {
		volumeAnalysis = "Above Average (>1.2x) ✅"
	} else if volumeRatio < 0.8 {
		volumeAnalysis = "Low Volume (<0.8x) 📉"
	} else {
		volumeAnalysis = "Normal Volume 📊"
	}
	fmt.Printf("   🔊 Volume: %.1fx average - %s\n", volumeRatio, volumeAnalysis)

	// Overall Assessment
	fmt.Printf("🎯 Overall Assessment:\n")
	if analysis.Confidence >= 80 {
		fmt.Printf("   🚀 STRONG %s Signal (%.1f%%)\n", analysis.Signal, analysis.Confidence)
	} else if analysis.Confidence >= 75 {
		fmt.Printf("   📈 GOOD %s Signal (%.1f%%)\n", analysis.Signal, analysis.Confidence)
	} else if analysis.Confidence >= 60 {
		fmt.Printf("   📊 WEAK %s Signal (%.1f%%)\n", analysis.Signal, analysis.Confidence)
	} else {
		fmt.Printf("   ⚠️  NO CLEAR Signal (%.1f%%)\n", analysis.Confidence)
	}

	// Trading Recommendation
	fmt.Printf("💡 Trading Recommendation:\n")
	if analysis.Signal == "LONG" && analysis.Confidence >= 75 {
		fmt.Printf("   ✅ Consider LONG position\n")
		fmt.Printf("   🛑 Stop Loss: Below $%.2f (EMA21)\n", analysis.EMA21*0.995)
		fmt.Printf("   🎯 Take Profit: Above $%.2f (ATR based)\n", analysis.Price+analysis.ATR*2)
	} else if analysis.Signal == "SHORT" && analysis.Confidence >= 75 {
		fmt.Printf("   ✅ Consider SHORT position\n")
		fmt.Printf("   🛑 Stop Loss: Above $%.2f (EMA21)\n", analysis.EMA21*1.005)
		fmt.Printf("   🎯 Take Profit: Below $%.2f (ATR based)\n", analysis.Price-analysis.ATR*2)
	} else {
		fmt.Printf("   ⏳ Wait for better opportunity\n")
		fmt.Printf("   👀 Monitor EMA alignment and volume\n")
	}
}

func getIcon(condition bool) string {
	if condition {
		return "✅"
	}
	return "❌"
}
