package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// SimpleMarketStrategy - กลยุทธ์เรียบง่ายสำหรับตลาดจริง
type SimpleMarketStrategy struct {
	dataManager    *RealDataManager
	initialCapital float64
	currentCapital float64
	positions      []*SimplePosition
	trades         []*SimpleTrade
	name           string
}

// SimplePosition - โพซิชั่นง่ายๆ
type SimplePosition struct {
	Symbol     string
	Direction  string
	EntryPrice float64
	EntryTime  time.Time
	Size       float64
	StopLoss   float64
	IsActive   bool
}

// SimpleTrade - เทรดที่เสร็จสิ้น
type SimpleTrade struct {
	Symbol     string
	Direction  string
	EntryPrice float64
	ExitPrice  float64
	PnL        float64
	PnLPercent float64
	ExitReason string
	Duration   time.Duration
}

// NewSimpleMarketStrategy - สร้างกลยุทธ์ใหม่
func NewSimpleMarketStrategy(symbol, strategyName string, capital float64) *SimpleMarketStrategy {
	return &SimpleMarketStrategy{
		dataManager:    NewRealDataManager(symbol),
		initialCapital: capital,
		currentCapital: capital,
		positions:      make([]*SimplePosition, 0),
		trades:         make([]*SimpleTrade, 0),
		name:           strategyName,
	}
}

// Strategy 1: RSI Bounce Strategy - ตีกลับจากแนวรับ/แนวต้าน
func (s *SimpleMarketStrategy) RSIBounceStrategy() error {
	fmt.Printf("🎯 Testing: RSI Bounce Strategy - %s\n", s.dataManager.symbol)

	if err := s.dataManager.LoadRealMarketData(); err != nil {
		return err
	}

	data := s.dataManager.dataBuffer

	for i := 20; i < len(data)-1; i++ {
		currentTime := time.Unix(data[i].Timestamp, 0)
		nextCandle := data[i+1]

		// ตรวจสอบสัญญาณเข้า
		if len(s.positions) == 0 {
			signal := s.analyzeRSIBounce(data, i)
			if signal != nil {
				s.executeSimpleEntry(signal, data[i], currentTime)
			}
		}

		// จัดการโพซิชั่นที่มี
		for _, pos := range s.positions {
			if pos.IsActive {
				if s.checkSimpleExit(pos, nextCandle, currentTime) {
					s.executeSimpleExit(pos, nextCandle, currentTime)
				}
			}
		}
	}

	s.printResults("RSI Bounce Strategy")
	return nil
}

// analyzeRSIBounce - วิเคราะห์สัญญาณ RSI Bounce
func (s *SimpleMarketStrategy) analyzeRSIBounce(data []OHLCV, index int) *SimpleSignal {
	if index < 20 {
		return nil
	}

	// คำนวณ RSI
	recentData := data[index-19 : index+1]
	rsi := s.dataManager.CalculateRSI(14, recentData)
	if rsi == nil {
		return nil
	}

	currentRSI := rsi[len(rsi)-1]
	prevRSI := rsi[len(rsi)-2]
	currentPrice := data[index].Close

	// LONG: RSI ต่ำกว่า 30 แล้วเริ่มขึ้น (Oversold Bounce)
	if prevRSI < 30 && currentRSI > prevRSI && currentRSI < 40 {
		return &SimpleSignal{
			Direction: "LONG",
			Price:     currentPrice,
			Reason:    "RSI Oversold Bounce",
		}
	}

	// SHORT: RSI สูงกว่า 70 แล้วเริ่มลง (Overbought Decline)
	if prevRSI > 70 && currentRSI < prevRSI && currentRSI > 60 {
		return &SimpleSignal{
			Direction: "SHORT",
			Price:     currentPrice,
			Reason:    "RSI Overbought Decline",
		}
	}

	return nil
}

// Strategy 2: Volume Breakout Strategy - ปริมาณการซื้อขายสูง
func (s *SimpleMarketStrategy) VolumeBreakoutStrategy() error {
	fmt.Printf("🎯 Testing: Volume Breakout Strategy - %s\n", s.dataManager.symbol)

	if err := s.dataManager.LoadRealMarketData(); err != nil {
		return err
	}

	data := s.dataManager.dataBuffer

	for i := 20; i < len(data)-1; i++ {
		currentTime := time.Unix(data[i].Timestamp, 0)
		nextCandle := data[i+1]

		if len(s.positions) == 0 {
			signal := s.analyzeVolumeBreakout(data, i)
			if signal != nil {
				s.executeSimpleEntry(signal, data[i], currentTime)
			}
		}

		for _, pos := range s.positions {
			if pos.IsActive {
				if s.checkSimpleExit(pos, nextCandle, currentTime) {
					s.executeSimpleExit(pos, nextCandle, currentTime)
				}
			}
		}
	}

	s.printResults("Volume Breakout Strategy")
	return nil
}

// analyzeVolumeBreakout - วิเคราะห์สัญญาณ Volume Breakout
func (s *SimpleMarketStrategy) analyzeVolumeBreakout(data []OHLCV, index int) *SimpleSignal {
	if index < 20 {
		return nil
	}

	currentCandle := data[index]
	prevCandle := data[index-1]

	// คำนวณ Volume average 10 periods
	volumeSum := 0.0
	for i := index - 9; i <= index; i++ {
		volumeSum += data[i].Volume
	}
	avgVolume := volumeSum / 10.0

	// เงื่อนไข: Volume มากกว่า 2x average + ราคาขึ้น > 1%
	volumeSpike := currentCandle.Volume > avgVolume*2.0
	priceUp := (currentCandle.Close-prevCandle.Close)/prevCandle.Close*100 > 1.0
	priceDown := (prevCandle.Close-currentCandle.Close)/prevCandle.Close*100 > 1.0

	if volumeSpike && priceUp {
		return &SimpleSignal{
			Direction: "LONG",
			Price:     currentCandle.Close,
			Reason:    "Volume Breakout Up",
		}
	}

	if volumeSpike && priceDown {
		return &SimpleSignal{
			Direction: "SHORT",
			Price:     currentCandle.Close,
			Reason:    "Volume Breakout Down",
		}
	}

	return nil
}

// Strategy 3: Price Action Strategy - การเคลื่อนไหวของราคา
func (s *SimpleMarketStrategy) PriceActionStrategy() error {
	fmt.Printf("🎯 Testing: Price Action Strategy - %s\n", s.dataManager.symbol)

	if err := s.dataManager.LoadRealMarketData(); err != nil {
		return err
	}

	data := s.dataManager.dataBuffer

	for i := 10; i < len(data)-1; i++ {
		currentTime := time.Unix(data[i].Timestamp, 0)
		nextCandle := data[i+1]

		if len(s.positions) == 0 {
			signal := s.analyzePriceAction(data, i)
			if signal != nil {
				s.executeSimpleEntry(signal, data[i], currentTime)
			}
		}

		for _, pos := range s.positions {
			if pos.IsActive {
				if s.checkSimpleExit(pos, nextCandle, currentTime) {
					s.executeSimpleExit(pos, nextCandle, currentTime)
				}
			}
		}
	}

	s.printResults("Price Action Strategy")
	return nil
}

// analyzePriceAction - วิเคราะห์ Price Action
func (s *SimpleMarketStrategy) analyzePriceAction(data []OHLCV, index int) *SimpleSignal {
	if index < 5 {
		return nil
	}

	current := data[index]
	prev1 := data[index-1]
	prev2 := data[index-2]

	// Bullish Engulfing Pattern
	if prev1.Close < prev1.Open && // Previous red candle
		current.Close > current.Open && // Current green candle
		current.Open < prev1.Close && // Gap down open
		current.Close > prev1.Open { // Close above previous open

		return &SimpleSignal{
			Direction: "LONG",
			Price:     current.Close,
			Reason:    "Bullish Engulfing",
		}
	}

	// Bearish Engulfing Pattern
	if prev1.Close > prev1.Open && // Previous green candle
		current.Close < current.Open && // Current red candle
		current.Open > prev1.Close && // Gap up open
		current.Close < prev1.Open { // Close below previous open

		return &SimpleSignal{
			Direction: "SHORT",
			Price:     current.Close,
			Reason:    "Bearish Engulfing",
		}
	}

	// Hammer/Doji patterns
	bodySize := math.Abs(current.Close - current.Open)
	totalSize := current.High - current.Low

	if bodySize/totalSize < 0.3 && // Small body
		(current.Low < prev1.Low && current.Low < prev2.Low) { // Lower low

		return &SimpleSignal{
			Direction: "LONG",
			Price:     current.Close,
			Reason:    "Hammer/Doji Reversal",
		}
	}

	return nil
}

// Strategy 4: Mean Reversion Strategy - กลับสู่ค่าเฉลี่ย
func (s *SimpleMarketStrategy) MeanReversionStrategy() error {
	fmt.Printf("🎯 Testing: Mean Reversion Strategy - %s\n", s.dataManager.symbol)

	if err := s.dataManager.LoadRealMarketData(); err != nil {
		return err
	}

	data := s.dataManager.dataBuffer

	for i := 20; i < len(data)-1; i++ {
		currentTime := time.Unix(data[i].Timestamp, 0)
		nextCandle := data[i+1]

		if len(s.positions) == 0 {
			signal := s.analyzeMeanReversion(data, i)
			if signal != nil {
				s.executeSimpleEntry(signal, data[i], currentTime)
			}
		}

		for _, pos := range s.positions {
			if pos.IsActive {
				if s.checkSimpleExit(pos, nextCandle, currentTime) {
					s.executeSimpleExit(pos, nextCandle, currentTime)
				}
			}
		}
	}

	s.printResults("Mean Reversion Strategy")
	return nil
}

// analyzeMeanReversion - วิเคราะห์ Mean Reversion
func (s *SimpleMarketStrategy) analyzeMeanReversion(data []OHLCV, index int) *SimpleSignal {
	if index < 20 {
		return nil
	}

	// คำนวณ Moving Average 20 periods
	sum := 0.0
	for i := index - 19; i <= index; i++ {
		sum += data[i].Close
	}
	ma20 := sum / 20.0

	currentPrice := data[index].Close
	deviation := (currentPrice - ma20) / ma20 * 100

	// ถ้าราคาอยู่ต่ำกว่า MA 2%+ → LONG (คาดว่าจะกลับขึ้น)
	if deviation < -2.0 {
		return &SimpleSignal{
			Direction: "LONG",
			Price:     currentPrice,
			Reason:    "Mean Reversion Up",
		}
	}

	// ถ้าราคาอยู่สูงกว่า MA 2%+ → SHORT (คาดว่าจะกลับลง)
	if deviation > 2.0 {
		return &SimpleSignal{
			Direction: "SHORT",
			Price:     currentPrice,
			Reason:    "Mean Reversion Down",
		}
	}

	return nil
}

// SimpleSignal - สัญญาณง่ายๆ
type SimpleSignal struct {
	Direction string
	Price     float64
	Reason    string
}

// executeSimpleEntry - เข้าเทรดง่ายๆ
func (s *SimpleMarketStrategy) executeSimpleEntry(signal *SimpleSignal, candle OHLCV, currentTime time.Time) {
	// Risk 0.5% per trade (Conservative)
	riskAmount := s.currentCapital * 0.005

	// Simple stop loss: 1% from entry
	var stopLoss float64
	if signal.Direction == "LONG" {
		stopLoss = signal.Price * 0.99 // 1% below
	} else {
		stopLoss = signal.Price * 1.01 // 1% above
	}

	stopDistance := math.Abs(signal.Price - stopLoss)
	positionSize := riskAmount / stopDistance

	position := &SimplePosition{
		Symbol:     s.dataManager.symbol,
		Direction:  signal.Direction,
		EntryPrice: signal.Price,
		EntryTime:  currentTime,
		Size:       positionSize,
		StopLoss:   stopLoss,
		IsActive:   true,
	}

	s.positions = append(s.positions, position)
	fmt.Printf("🎯 SIMPLE ENTRY %s: %.4f (Reason: %s)\n",
		signal.Direction, signal.Price, signal.Reason)
}

// checkSimpleExit - ตรวจสอบการออกง่ายๆ
func (s *SimpleMarketStrategy) checkSimpleExit(pos *SimplePosition, candle OHLCV, currentTime time.Time) bool {
	currentPrice := candle.Close

	// Stop Loss
	if (pos.Direction == "LONG" && currentPrice <= pos.StopLoss) ||
		(pos.Direction == "SHORT" && currentPrice >= pos.StopLoss) {
		return true
	}

	// Take Profit: 2% (2:1 R/R ratio)
	profitTarget := 0.0
	if pos.Direction == "LONG" {
		profitTarget = pos.EntryPrice * 1.02 // 2% above
		if currentPrice >= profitTarget {
			return true
		}
	} else {
		profitTarget = pos.EntryPrice * 0.98 // 2% below
		if currentPrice <= profitTarget {
			return true
		}
	}

	// Time limit: 12 hours (shorter than before)
	if currentTime.Sub(pos.EntryTime) > 12*time.Hour {
		return true
	}

	return false
}

// executeSimpleExit - ออกจากเทรดง่ายๆ
func (s *SimpleMarketStrategy) executeSimpleExit(pos *SimplePosition, candle OHLCV, currentTime time.Time) {
	exitPrice := candle.Close
	duration := currentTime.Sub(pos.EntryTime)

	var pnl, pnlPercent float64
	var reason string

	if pos.Direction == "LONG" {
		pnl = (exitPrice - pos.EntryPrice) * pos.Size
		pnlPercent = (exitPrice - pos.EntryPrice) / pos.EntryPrice * 100
	} else {
		pnl = (pos.EntryPrice - exitPrice) * pos.Size
		pnlPercent = (pos.EntryPrice - exitPrice) / pos.EntryPrice * 100
	}

	// กำหนดเหตุผลการออก
	if (pos.Direction == "LONG" && exitPrice <= pos.StopLoss) ||
		(pos.Direction == "SHORT" && exitPrice >= pos.StopLoss) {
		reason = "Stop Loss"
	} else if pnlPercent >= 1.8 { // Close to 2% target
		reason = "Take Profit"
	} else {
		reason = "Time Limit"
	}

	s.currentCapital += pnl

	trade := &SimpleTrade{
		Symbol:     pos.Symbol,
		Direction:  pos.Direction,
		EntryPrice: pos.EntryPrice,
		ExitPrice:  exitPrice,
		PnL:        pnl,
		PnLPercent: pnlPercent,
		ExitReason: reason,
		Duration:   duration,
	}

	s.trades = append(s.trades, trade)
	pos.IsActive = false

	fmt.Printf("🔹 SIMPLE EXIT %s: PnL %.2f%% - %s\n",
		pos.Direction, pnlPercent, reason)
}

// printResults - แสดงผลลัพธ์
func (s *SimpleMarketStrategy) printResults(strategyName string) {
	totalReturn := ((s.currentCapital - s.initialCapital) / s.initialCapital) * 100
	totalTrades := len(s.trades)

	winTrades := 0
	totalDuration := time.Duration(0)
	for _, trade := range s.trades {
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

	fmt.Printf("\n📊 %s RESULTS: %s\n", strategyName, s.dataManager.symbol)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("💰 Return: $%.2f → $%.2f (%.2f%%)\n",
		s.initialCapital, s.currentCapital, totalReturn)
	fmt.Printf("📈 Trades: %d (Win: %d, Loss: %d, Rate: %.1f%%)\n",
		totalTrades, winTrades, totalTrades-winTrades, winRate)
	fmt.Printf("⏱️ Avg Duration: %s\n", avgDuration.Truncate(time.Hour))

	if totalReturn > 0 && winRate >= 50 {
		fmt.Printf("✅ PROFITABLE Strategy!\n")
	} else {
		fmt.Printf("📉 Needs improvement\n")
	}
}

// TestAllSimpleStrategies - ทดสอบกลยุทธ์ง่ายๆ ทั้งหมด
func main() {
	symbols := []string{"BTC_USDT", "ETH_USDT", "SOL_USDT"}
	strategies := []string{"RSI Bounce", "Volume Breakout", "Price Action", "Mean Reversion"}

	fmt.Println("🚀 SIMPLE PROFITABLE STRATEGIES TEST")
	fmt.Println("📊 Testing with REAL market data")
	fmt.Println(strings.Repeat("━", 80))

	bestStrategy := ""
	bestReturn := -999.0

	for _, symbol := range symbols {
		fmt.Printf("\n🔍 Testing Symbol: %s\n", symbol)
		fmt.Println(strings.Repeat("-", 50))

		// Test RSI Bounce
		strategy1 := NewSimpleMarketStrategy(symbol, "RSI Bounce", 10000.0)
		strategy1.RSIBounceStrategy()
		return1 := ((strategy1.currentCapital - strategy1.initialCapital) / strategy1.initialCapital) * 100

		time.Sleep(1 * time.Second)

		// Test Volume Breakout
		strategy2 := NewSimpleMarketStrategy(symbol, "Volume Breakout", 10000.0)
		strategy2.VolumeBreakoutStrategy()
		return2 := ((strategy2.currentCapital - strategy2.initialCapital) / strategy2.initialCapital) * 100

		time.Sleep(1 * time.Second)

		// Test Price Action
		strategy3 := NewSimpleMarketStrategy(symbol, "Price Action", 10000.0)
		strategy3.PriceActionStrategy()
		return3 := ((strategy3.currentCapital - strategy3.initialCapital) / strategy3.initialCapital) * 100

		time.Sleep(1 * time.Second)

		// Test Mean Reversion
		strategy4 := NewSimpleMarketStrategy(symbol, "Mean Reversion", 10000.0)
		strategy4.MeanReversionStrategy()
		return4 := ((strategy4.currentCapital - strategy4.initialCapital) / strategy4.initialCapital) * 100

		// หาผลตอบแทนที่ดีที่สุด
		returns := []float64{return1, return2, return3, return4}
		for i, ret := range returns {
			if ret > bestReturn {
				bestReturn = ret
				bestStrategy = fmt.Sprintf("%s - %s", strategies[i], symbol)
			}
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Printf("\n🏆 BEST SIMPLE STRATEGY: %s\n", bestStrategy)
	fmt.Printf("📈 Best Return: %.2f%%\n", bestReturn)
	fmt.Printf("✅ Simple strategies testing complete!\n")
}
