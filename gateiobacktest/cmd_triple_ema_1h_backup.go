package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gateio-trading-bot/internal/gateio"
	"gateio-trading-bot/internal/trading"
)

func main() {
	fmt.Println("🚀 Gate.io Triple EMA 1H Strategy Backtester (144 Candles)")
	fmt.Println(strings.Repeat("=", 80))

	// โหลดการตั้งค่า
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Fatalf("❌ .env file not found")
	}

	// สร้าง Gate.io client - ใช้ demo data เท่านั้น
	client := gateio.NewClient("", "", "https://api.gateio.ws")

	// ทดสอบ Triple EMA 1H Strategy
	if len(os.Args) > 1 && os.Args[1] == "multi" {
		runMultiSymbol1H(client)
	} else {
		runSingle1H(client)
	}

	fmt.Println("\n✅ Triple EMA 1H Backtesting สำเร็จ!")
}

// runSingle1H ทดสอบ symbol เดียว
func runSingle1H(client *gateio.Client) {
	symbol := "SOL_USDT"
	fmt.Printf("📊 Testing %s with Triple EMA 1H Strategy...\n\n", symbol)

	result := testSymbol1H(symbol, client, 10000.0)
	if result != nil {
		printSingleResult1H(result)
		saveSingleResult1H(result)
	}
}

// runMultiSymbol1H ทดสอบหลาย symbols
func runMultiSymbol1H(client *gateio.Client) {
	symbols := []string{"SOL_USDT", "BTC_USDT", "ETH_USDT"}
	fmt.Printf("📊 Testing %d symbols with Triple EMA 1H Strategy...\n\n", len(symbols))

	var results []trading.BacktestResult

	for i, symbol := range symbols {
		fmt.Printf("🔍 [%d/%d] Testing %s...\n", i+1, len(symbols), symbol)
		
		if result := testSymbol1H(symbol, client, 10000.0); result != nil {
			results = append(results, *result)
			time.Sleep(500 * time.Millisecond) // หน่วงเวลา
		}
		fmt.Println()
	}

	printMultiResults1H(results)
	saveMultiResults1H(results)
}

// testSymbol1H ทดสอบ symbol เดียวด้วย Triple EMA 1H (ใช้ข้อมูล demo)
func testSymbol1H(symbol string, client *gateio.Client, initialCapital float64) *trading.BacktestResult {
	fmt.Printf("  📈 สร้างข้อมูล Demo %s (1H, 150 candles)...\n", symbol)

	// สร้างข้อมูล demo สำหรับทดสอบ
	candlesticks := generateDemoData1H(symbol, 150)

	if len(candlesticks) < 144 {
		fmt.Printf("  ❌ %s: ข้อมูลไม่ครบ (ต้องการ 144, ได้ %d)\n", symbol, len(candlesticks))
		return nil
	}

	fmt.Printf("  ✅ ได้ข้อมูล %d candles\n", len(candlesticks))

	// แปลงเป็น OHLCV format
	ohlcvData := make([]trading.OHLCV, len(candlesticks))
	for i, candle := range candlesticks {
		ohlcvData[i] = trading.OHLCV{
			Timestamp: candle.Timestamp,
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
			Volume:    candle.Volume,
		}
	}

	// สร้าง backtester
	backtester, err := trading.NewBacktesterSimple(symbol, 6, initialCapital) // 6 วัน
	if err != nil {
		fmt.Printf("  ❌ Error creating backtester: %v\n", err)
		return nil
	}

	// ตั้งข้อมูลโดยตรง
	backtester.LoadOHLCVData(ohlcvData, 
		time.Unix(candlesticks[0].Timestamp, 0),
		time.Unix(candlesticks[len(candlesticks)-1].Timestamp, 0))

	// รัน Triple EMA 1H strategy
	fmt.Printf("  🔄 กำลังรัน Triple EMA 1H Strategy...\n")
	result := backtester.RunTripleEMA1HStrategy()

	// แสดงผลลัพธ์
	fmt.Printf("  📊 %s Results:\n", symbol)
	fmt.Printf("    💰 Return: $%.2f → $%.2f (%.2f%%)\n", 
		result.InitialCapital, result.FinalCapital, result.TotalReturnPct)
	fmt.Printf("    📈 Trades: %d (Win: %d, Loss: %d, Rate: %.1f%%)\n", 
		result.TotalTrades, result.WinningTrades, result.LosingTrades, result.WinRate)

	return result
}

// generateDemoData1H สร้างข้อมูล demo สำหรับ 1H timeframe
func generateDemoData1H(symbol string, count int) []gateio.Candlestick {
	candlesticks := make([]gateio.Candlestick, count)
	
	// ราคาเริ่มต้นตาม symbol
	var basePrice float64
	switch symbol {
	case "SOL_USDT":
		basePrice = 200.0
	case "BTC_USDT":
		basePrice = 90000.0
	case "ETH_USDT":
		basePrice = 3200.0
	default:
		basePrice = 100.0
	}
	
	// เริ่มจาก 150 ชั่วโมงที่แล้ว
	startTime := time.Now().Add(-time.Duration(count) * time.Hour)
	
	currentPrice := basePrice
	
	for i := 0; i < count; i++ {
		timestamp := startTime.Add(time.Duration(i) * time.Hour).Unix()
		
		// สร้างการเปลี่ยนแปลงราคาแบบสุ่ม
		change := (float64(i%7-3) * 0.02) + (float64(i%13-6) * 0.01) // รูปแบบ trend
		currentPrice = currentPrice * (1 + change)
		
		// สร้าง OHLC
		open := currentPrice
		volatility := currentPrice * 0.015 // 1.5% volatility
		
		high := open + volatility*(float64(i%3)/2.0)
		low := open - volatility*(float64((i+1)%3)/2.0)
		close := low + (high-low)*(0.3+float64(i%7)/10.0)
		
		// ปรับให้ close เป็นราคาต่อไป
		currentPrice = close
		
		// Volume แบบสุ่ม
		volume := 1000000 + float64(i%100)*10000
		
		candlesticks[i] = gateio.Candlestick{
			Timestamp: timestamp,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}
	}
	
	return candlesticks
}

// printSingleResult1H แสดงผลลัพธ์เดียว
func printSingleResult1H(result *trading.BacktestResult) {
	fmt.Printf("\n🎯 Triple EMA 1H Strategy Results: %s\n", result.Symbol)
	fmt.Println(strings.Repeat("=", 60))
	
	fmt.Printf("⏰ Period: %s → %s (%.1f days)\n", 
		result.StartDate.Format("2006-01-02"), 
		result.EndDate.Format("2006-01-02"),
		result.EndDate.Sub(result.StartDate).Hours()/24)
	
	fmt.Printf("💰 Capital: $%.2f → $%.2f\n", result.InitialCapital, result.FinalCapital)
	fmt.Printf("📈 Return: $%.2f (%.2f%%)\n", result.TotalReturn, result.TotalReturnPct)
	
	if result.TotalTrades > 0 {
		fmt.Printf("🎯 Trades: %d total\n", result.TotalTrades)
		fmt.Printf("   ✅ Winning: %d (%.1f%%)\n", result.WinningTrades, result.WinRate)
		fmt.Printf("   ❌ Losing: %d\n", result.LosingTrades)
		
		// แสดง trades ล่าสุด 5 อัน
		if len(result.Trades) > 0 {
			fmt.Printf("\n📋 Recent Trades:\n")
			start := len(result.Trades) - 5
			if start < 0 {
				start = 0
			}
			
			for i := start; i < len(result.Trades); i++ {
				trade := result.Trades[i]
				if !trade.ExitTime.IsZero() {
					fmt.Printf("   %d. %s: %+.2f%% (%s)\n", 
						trade.ID, trade.Side, trade.PnLPct, trade.ExitReason)
				}
			}
		}
	}
}

// printMultiResults1H แสดงผลลัพธ์หลาย symbols
func printMultiResults1H(results []trading.BacktestResult) {
	fmt.Printf("\n🎯 สรุปผลลัพธ์ Triple EMA 1H Strategy\n")
	fmt.Println(strings.Repeat("=", 80))

	if len(results) == 0 {
		fmt.Println("❌ ไม่มีข้อมูลผลลัพธ์")
		return
	}

	// Header
	fmt.Printf("%-12s %-12s %-12s %-8s %-8s %-10s\n", 
		"Symbol", "Return", "Return%", "Trades", "Win%", "Performance")
	fmt.Println(strings.Repeat("-", 70))

	var totalInitial, totalFinal float64
	var bestReturn, worstReturn float64
	var bestSymbol, worstSymbol string

	for i, result := range results {
		totalInitial += result.InitialCapital
		totalFinal += result.FinalCapital

		performance := "📊"
		if result.TotalReturnPct > 50 {
			performance = "🚀"
		} else if result.TotalReturnPct > 20 {
			performance = "📈"
		} else if result.TotalReturnPct < -20 {
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

		fmt.Printf("%-12s $%-10.2f %+8.2f%% %-8d %6.1f%% %-10s\n",
			result.Symbol,
			result.TotalReturn,
			result.TotalReturnPct,
			result.TotalTrades,
			result.WinRate,
			performance)
	}

	fmt.Println(strings.Repeat("-", 70))

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
	if avgReturn > 30 {
		fmt.Println("   🚀 กลยุทธ์ Triple EMA 1H ให้ผลลัพธ์ยอดเยี่ยม!")
	} else if avgReturn > 15 {
		fmt.Println("   📈 กลยุทธ์ Triple EMA 1H ให้ผลลัพธ์ดี")
	} else if avgReturn > 0 {
		fmt.Println("   📊 กลยุทธ์ Triple EMA 1H ให้ผลลัพธ์เป็นบวก")
	} else {
		fmt.Println("   ⚠️  กลยุทธ์ Triple EMA 1H ควรปรับปรุง")
	}

	fmt.Printf("   ⏰ Time Frame: 1H (144 candles = 6 days)\n")
	fmt.Printf("   🎯 Strategy: Triple EMA (9,21,50) + RSI + Volume + ATR\n")
	fmt.Printf("   💰 Risk: 1%% per trade, 1:2 R/R, 5x leverage\n")
}

// saveSingleResult1H บันทึกผลลัพธ์เดียว
func saveSingleResult1H(result *trading.BacktestResult) {
	fmt.Printf("💾 ผลลัพธ์: %s ได้ %.2f%% จาก %d trades\n", 
		result.Symbol, result.TotalReturnPct, result.TotalTrades)
}

// saveMultiResults1H บันทึกผลลัพธ์หลาย symbols
func saveMultiResults1H(results []trading.BacktestResult) {
	fmt.Printf("💾 ผลลัพธ์รวม: %d symbols tested\n", len(results))
	
	var totalReturn float64
	for _, result := range results {
		totalReturn += result.TotalReturnPct
	}
	
	if len(results) > 0 {
		avgReturn := totalReturn / float64(len(results))
		fmt.Printf("� ผลตอบแทนเฉลี่ย: %.2f%%\n", avgReturn)
	}
}
