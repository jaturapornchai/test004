package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

// RealMarketBacktester - ระบบ backtest ด้วยข้อมูลจริง
type RealMarketBacktester struct {
	dataManager     *RealDataManager
	initialCapital  float64
	currentCapital  float64
	positions       []*RealPosition
	completedTrades []*RealTrade
	maxHoldHours    int
}

// RealPosition - โพซิชั่นการเทรดจริง
type RealPosition struct {
	Symbol     string    `json:"symbol"`
	Direction  string    `json:"direction"`
	EntryPrice float64   `json:"entry_price"`
	EntryTime  time.Time `json:"entry_time"`
	Size       float64   `json:"size"`
	StopLoss   float64   `json:"stop_loss"`
	TakeProfit float64   `json:"take_profit"`
	Confidence float64   `json:"confidence"`
	ATR        float64   `json:"atr"`
	TimeLimit  time.Time `json:"time_limit"`
	IsActive   bool      `json:"is_active"`
}

// RealTrade - เทรดที่เสร็จสิ้นแล้ว
type RealTrade struct {
	Symbol     string        `json:"symbol"`
	Direction  string        `json:"direction"`
	EntryPrice float64       `json:"entry_price"`
	ExitPrice  float64       `json:"exit_price"`
	EntryTime  time.Time     `json:"entry_time"`
	ExitTime   time.Time     `json:"exit_time"`
	PnL        float64       `json:"pnl"`
	PnLPercent float64       `json:"pnl_percent"`
	Confidence float64       `json:"confidence"`
	ExitReason string        `json:"exit_reason"`
	Duration   time.Duration `json:"duration"`
}

// NewRealMarketBacktester - สร้าง real market backtester
func NewRealMarketBacktester(symbol string, initialCapital float64) *RealMarketBacktester {
	return &RealMarketBacktester{
		dataManager:     NewRealDataManager(symbol),
		initialCapital:  initialCapital,
		currentCapital:  initialCapital,
		positions:       make([]*RealPosition, 0),
		completedTrades: make([]*RealTrade, 0),
		maxHoldHours:    24,
	}
}

// RunRealMarketBacktest - รัน backtest ด้วยข้อมูลจริง
func (rb *RealMarketBacktester) RunRealMarketBacktest() error {
	symbol := rb.dataManager.symbol
	fmt.Printf("🚀 Real Market Backtest: %s\n", symbol)
	fmt.Printf("💰 Initial Capital: $%.2f\n", rb.initialCapital)

	// ดึงข้อมูลจริงจาก Gate.io
	if err := rb.dataManager.LoadRealMarketData(); err != nil {
		return fmt.Errorf("failed to load real market data: %v", err)
	}

	// วิเคราะห์ข้อมูลจริง
	analysis := rb.dataManager.GetRealMarketAnalysis()
	if analysis == nil {
		return fmt.Errorf("failed to analyze real market data")
	}

	analysis.PrintMarketSummary()

	// จำลองการเทรดย้อนหลัง (ใช้ข้อมูล 144 candles)
	data := rb.dataManager.dataBuffer

	// เริ่มการ backtest จากตำแหน่งที่ 50 (ให้ indicator มีข้อมูลเพียงพอ)
	for i := 50; i < len(data)-1; i++ {
		currentTime := time.Unix(data[i].Timestamp, 0)
		nextCandle := data[i+1]

		// ตรวจสอบสัญญาณเข้าซื้อ/ขาย
		if len(rb.positions) == 0 { // ไม่มีโพซิชั่นอยู่
			signal := rb.analyzeRealEntrySignal(data, i)
			if signal != nil {
				rb.executeRealEntry(signal, data[i], currentTime)
			}
		}

		// จัดการโพซิชั่นที่มีอยู่
		for _, pos := range rb.positions {
			if pos.IsActive {
				exitReason := rb.checkRealExitConditions(pos, nextCandle, currentTime)
				if exitReason != "" {
					rb.executeRealExit(pos, nextCandle, currentTime, exitReason)
				}
			}
		}
	}

	// ปิดโพซิชั่นที่เหลือ (ถ้ามี)
	lastCandle := data[len(data)-1]
	lastTime := time.Unix(lastCandle.Timestamp, 0)
	for _, pos := range rb.positions {
		if pos.IsActive {
			rb.executeRealExit(pos, lastCandle, lastTime, "End of Backtest")
		}
	}

	// แสดงผลลัพธ์
	rb.printRealBacktestResults()

	return nil
}

// analyzeRealEntrySignal - วิเคราะห์สัญญาณเข้าจากข้อมูลจริง
func (rb *RealMarketBacktester) analyzeRealEntrySignal(data []OHLCV, index int) *RealEntrySignal {
	if index < 50 || index >= len(data)-1 {
		return nil
	}

	// คำนวณ indicators ที่จำเป็น
	recentData := data[index-49 : index+1] // ใช้ 50 candles สำหรับการคำนวณ

	ema9 := rb.dataManager.CalculateEMA(9, recentData)
	ema21 := rb.dataManager.CalculateEMA(21, recentData)
	ema50 := rb.dataManager.CalculateEMA(50, recentData)
	rsi := rb.dataManager.CalculateRSI(14, recentData)
	atr := rb.dataManager.CalculateATR(14, recentData)

	if ema9 == nil || ema21 == nil || ema50 == nil || rsi == nil || atr == nil {
		return nil
	}

	// ใช้ค่าล่าสุด
	latestIdx := len(recentData) - 1
	currentPrice := recentData[latestIdx].Close
	currentVolume := recentData[latestIdx].Volume

	// คำนวณ Volume MA
	volumeSum := 0.0
	volumePeriod := 20
	for i := latestIdx - volumePeriod + 1; i <= latestIdx; i++ {
		if i >= 0 {
			volumeSum += recentData[i].Volume
		}
	}
	volumeMA := volumeSum / float64(volumePeriod)

	// ตรวจสอบเงื่อนไข LONG
	longConditions := []bool{
		ema9[latestIdx] > ema21[latestIdx] && ema21[latestIdx] > ema50[latestIdx], // EMA alignment
		currentPrice >= ema9[latestIdx]*0.998,                                     // Price position
		rsi[latestIdx] >= 40 && rsi[latestIdx] <= 80,                              // RSI range
		currentVolume > volumeMA*1.2,                                              // Volume confirmation
		math.Abs(ema9[latestIdx]-ema50[latestIdx])/ema50[latestIdx]*100 > 0.15,    // EMA spread
	}

	// ตรวจสอบเงื่อนไข SHORT
	shortConditions := []bool{
		ema9[latestIdx] < ema21[latestIdx] && ema21[latestIdx] < ema50[latestIdx], // EMA alignment
		currentPrice <= ema9[latestIdx]*1.002,                                     // Price position
		rsi[latestIdx] >= 20 && rsi[latestIdx] <= 60,                              // RSI range
		currentVolume > volumeMA*1.2,                                              // Volume confirmation
		math.Abs(ema9[latestIdx]-ema50[latestIdx])/ema50[latestIdx]*100 > 0.15,    // EMA spread
	}

	// คำนวณ confidence
	var direction string
	var confidence float64

	longScore := 0
	for _, condition := range longConditions {
		if condition {
			longScore++
		}
	}

	shortScore := 0
	for _, condition := range shortConditions {
		if condition {
			shortScore++
		}
	}

	if longScore >= 4 { // ต้องผ่านอย่างน้อย 4 เงื่อนไข
		direction = "LONG"
		confidence = float64(longScore) / float64(len(longConditions))
	} else if shortScore >= 4 {
		direction = "SHORT"
		confidence = float64(shortScore) / float64(len(shortConditions))
	}

	// ต้องมี confidence >= 80%
	if confidence < 0.80 {
		return nil
	}

	return &RealEntrySignal{
		Symbol:     rb.dataManager.symbol,
		Direction:  direction,
		Price:      currentPrice,
		Confidence: confidence,
		ATR:        atr[latestIdx],
		EMA9:       ema9[latestIdx],
		EMA21:      ema21[latestIdx],
		EMA50:      ema50[latestIdx],
		RSI:        rsi[latestIdx],
	}
}

// RealEntrySignal - สัญญาณเข้าจากข้อมูลจริง
type RealEntrySignal struct {
	Symbol     string
	Direction  string
	Price      float64
	Confidence float64
	ATR        float64
	EMA9       float64
	EMA21      float64
	EMA50      float64
	RSI        float64
}

// executeRealEntry - ดำเนินการเข้าเทรดจากข้อมูลจริง
func (rb *RealMarketBacktester) executeRealEntry(signal *RealEntrySignal, candle OHLCV, currentTime time.Time) {
	// คำนวณ position size (1% risk)
	riskAmount := rb.currentCapital * 0.01
	stopDistance := signal.ATR * 2.0
	positionSize := riskAmount / stopDistance

	// กำหนด stop loss และ take profit
	var stopLoss, takeProfit float64
	if signal.Direction == "LONG" {
		stopLoss = signal.Price - stopDistance
		takeProfit = signal.Price + (stopDistance * 2.0) // 1:2 R/R
	} else {
		stopLoss = signal.Price + stopDistance
		takeProfit = signal.Price - (stopDistance * 2.0)
	}

	position := &RealPosition{
		Symbol:     signal.Symbol,
		Direction:  signal.Direction,
		EntryPrice: signal.Price,
		EntryTime:  currentTime,
		Size:       positionSize,
		StopLoss:   stopLoss,
		TakeProfit: takeProfit,
		Confidence: signal.Confidence,
		ATR:        signal.ATR,
		TimeLimit:  currentTime.Add(24 * time.Hour),
		IsActive:   true,
	}

	rb.positions = append(rb.positions, position)

	fmt.Printf("🎯 REAL ENTRY %s: %.4f, SL: %.4f, TP: %.4f, Conf: %.1f%%\n",
		signal.Direction, signal.Price, stopLoss, takeProfit, signal.Confidence*100)
}

// checkRealExitConditions - ตรวจสอบเงื่อนไขออกจากข้อมูลจริง
func (rb *RealMarketBacktester) checkRealExitConditions(pos *RealPosition, candle OHLCV, currentTime time.Time) string {
	currentPrice := candle.Close

	// Stop Loss
	if (pos.Direction == "LONG" && currentPrice <= pos.StopLoss) ||
		(pos.Direction == "SHORT" && currentPrice >= pos.StopLoss) {
		return "Stop Loss"
	}

	// Take Profit
	if (pos.Direction == "LONG" && currentPrice >= pos.TakeProfit) ||
		(pos.Direction == "SHORT" && currentPrice <= pos.TakeProfit) {
		return "Take Profit"
	}

	// Time Limit (24 hours)
	if currentTime.After(pos.TimeLimit) {
		return "Time Limit (24h)"
	}

	return "" // ไม่มีเงื่อนไขออก
}

// executeRealExit - ดำเนินการออกจากเทรด
func (rb *RealMarketBacktester) executeRealExit(pos *RealPosition, candle OHLCV, currentTime time.Time, reason string) {
	exitPrice := candle.Close
	duration := currentTime.Sub(pos.EntryTime)

	// คำนวณ P&L
	var pnl, pnlPercent float64
	if pos.Direction == "LONG" {
		pnl = (exitPrice - pos.EntryPrice) * pos.Size
		pnlPercent = (exitPrice - pos.EntryPrice) / pos.EntryPrice * 100
	} else {
		pnl = (pos.EntryPrice - exitPrice) * pos.Size
		pnlPercent = (pos.EntryPrice - exitPrice) / pos.EntryPrice * 100
	}

	// อัพเดท capital
	rb.currentCapital += pnl

	// สร้าง trade record
	trade := &RealTrade{
		Symbol:     pos.Symbol,
		Direction:  pos.Direction,
		EntryPrice: pos.EntryPrice,
		ExitPrice:  exitPrice,
		EntryTime:  pos.EntryTime,
		ExitTime:   currentTime,
		PnL:        pnl,
		PnLPercent: pnlPercent,
		Confidence: pos.Confidence,
		ExitReason: reason,
		Duration:   duration,
	}

	rb.completedTrades = append(rb.completedTrades, trade)
	pos.IsActive = false

	fmt.Printf("🔹 REAL EXIT %s: PnL $%.2f (%.2f%%) - %s\n",
		pos.Direction, pnl, pnlPercent, reason)
}

// printRealBacktestResults - แสดงผลลัพธ์ backtest จากข้อมูลจริง
func (rb *RealMarketBacktester) printRealBacktestResults() {
	totalReturn := ((rb.currentCapital - rb.initialCapital) / rb.initialCapital) * 100
	totalTrades := len(rb.completedTrades)

	winTrades := 0
	totalDuration := time.Duration(0)
	for _, trade := range rb.completedTrades {
		if trade.PnL > 0 {
			winTrades++
		}
		totalDuration += trade.Duration
	}

	winRate := 0.0
	avgDuration := time.Duration(0)
	if totalTrades > 0 {
		winRate = float64(winTrades) / float64(totalTrades) * 100
		avgDuration = totalDuration / time.Duration(totalTrades)
	}

	fmt.Printf("\n📊 REAL MARKET BACKTEST RESULTS: %s\n", rb.dataManager.symbol)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("💰 Return: $%.2f → $%.2f (%.2f%%)\n",
		rb.initialCapital, rb.currentCapital, totalReturn)
	fmt.Printf("📈 Trades: %d (Win: %d, Loss: %d, Rate: %.1f%%)\n",
		totalTrades, winTrades, totalTrades-winTrades, winRate)
	fmt.Printf("⏱️  Avg Trade Duration: %s\n", avgDuration.Truncate(time.Hour))

	// แสดงเทรดล่าสุด 3 รายการ
	if len(rb.completedTrades) > 0 {
		fmt.Printf("\n📋 Recent Trades:\n")
		start := len(rb.completedTrades) - 3
		if start < 0 {
			start = 0
		}
		for i, trade := range rb.completedTrades[start:] {
			fmt.Printf("   %d. %s: %+.2f%% (%.1fh, Conf: %.1f%%) - %s\n",
				start+i+1, trade.Direction, trade.PnLPercent,
				trade.Duration.Hours(), trade.Confidence*100, trade.ExitReason)
		}
	}

	// ประเมินผลลัพธ์
	fmt.Printf("\n🎯 Performance Assessment:\n")
	if totalReturn > 50 {
		fmt.Printf("   🚀 EXCELLENT Performance with REAL data!\n")
	} else if totalReturn > 20 {
		fmt.Printf("   ✅ Good Performance with REAL data\n")
	} else if totalReturn > 0 {
		fmt.Printf("   📈 Positive Performance with REAL data\n")
	} else {
		fmt.Printf("   📉 Strategy needs optimization with REAL data\n")
	}

	if winRate >= 70 {
		fmt.Printf("   ✅ Strong Win Rate (%.1f%%)\n", winRate)
	} else if winRate >= 50 {
		fmt.Printf("   📊 Decent Win Rate (%.1f%%)\n", winRate)
	} else {
		fmt.Printf("   ⚠️ Low Win Rate (%.1f%%) - Review strategy\n", winRate)
	}
}

// RunRealMarketBacktestSuite - รัน backtest ข้อมูลจริงหลายเหรียญ
func RunRealMarketBacktestSuite() {
	symbols := []string{"BTC_USDT", "ETH_USDT", "SOL_USDT"}
	initialCapital := 10000.0

	fmt.Println("🚀 REAL MARKET BACKTEST SUITE - Enhanced Triple EMA 1H")
	fmt.Println("📊 Using REAL data from Gate.io API")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	totalReturn := 0.0
	totalTrades := 0
	totalWins := 0

	for i, symbol := range symbols {
		fmt.Printf("\n🔍 [%d/%d] Real Market Backtest: %s\n", i+1, len(symbols), symbol)

		backtester := NewRealMarketBacktester(symbol, initialCapital)
		if err := backtester.RunRealMarketBacktest(); err != nil {
			log.Printf("❌ Backtest failed for %s: %v", symbol, err)
			continue
		}

		// สะสมสถิติ
		symbolReturn := ((backtester.currentCapital - initialCapital) / initialCapital) * 100
		totalReturn += symbolReturn
		totalTrades += len(backtester.completedTrades)

		for _, trade := range backtester.completedTrades {
			if trade.PnL > 0 {
				totalWins++
			}
		}

		// หน่วงเวลาระหว่างการดึงข้อมูล
		if i < len(symbols)-1 {
			time.Sleep(3 * time.Second)
		}
	}

	// สรุปผลรวม
	avgReturn := totalReturn / float64(len(symbols))
	overallWinRate := 0.0
	if totalTrades > 0 {
		overallWinRate = float64(totalWins) / float64(totalTrades) * 100
	}

	fmt.Printf("\n🎯 REAL MARKET BACKTEST SUMMARY\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 Average Return: %.2f%% (across %d symbols)\n", avgReturn, len(symbols))
	fmt.Printf("📈 Total Trades: %d (Win Rate: %.1f%%)\n", totalTrades, overallWinRate)
	fmt.Printf("🎪 Strategy Performance: ")

	if avgReturn > 30 && overallWinRate > 70 {
		fmt.Printf("🚀 EXCELLENT with real market data!\n")
	} else if avgReturn > 10 && overallWinRate > 60 {
		fmt.Printf("✅ GOOD with real market data\n")
	} else if avgReturn > 0 {
		fmt.Printf("📈 POSITIVE but can be improved\n")
	} else {
		fmt.Printf("📉 NEEDS OPTIMIZATION\n")
	}

	fmt.Printf("✅ Real market backtest complete!\n")
}
