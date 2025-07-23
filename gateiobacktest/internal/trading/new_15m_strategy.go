package trading

import (
	"fmt"
	"math"
	"time"
)

// RunNew15mStrategy - กลยุทธ์ Triple Moving Average Momentum สำหรับ 15m timeframe
func (bt *Backtester) RunNew15mStrategy() *BacktestResult {
	fmt.Printf("\n🚀 เริ่มต้น Triple EMA Momentum Strategy %s\n", bt.symbol)
	fmt.Printf("💰 เงินทุนเริ่มต้น: $%.2f\n", bt.initialCapital)
	fmt.Printf("⏰ Time Frame: 15 นาที\n")
	fmt.Printf("🎯 กลยุทธ์: Triple EMA + RSI + Volume Momentum (5x Leverage)\n")
	fmt.Printf("📊 รองรับ: LONG & SHORT Positions with Smart Momentum Detection\n")

	for bt.currentIndex = 50; bt.currentIndex < len(bt.ohlcvData); bt.currentIndex++ {
		currentCandle := bt.ohlcvData[bt.currentIndex]
		bt.currentTime = time.Unix(currentCandle.Timestamp, 0)
		bt.currentPrice = currentCandle.Close

		// วิเคราะห์ตลาดด้วยกลยุทธ์ใหม่
		analysis := bt.analyze15mMarket()
		if analysis == nil {
			continue
		}

		// จัดการ position ที่มีอยู่
		bt.check15mPosition()

		// หาจุดเข้าใหม่
		bt.lookFor15mEntry(analysis)

		// Debug: แสดงข้อมูลทุก 1000 candles
		if bt.currentIndex%1000 == 0 {
			fmt.Printf("🔍 Debug [%d]: Price=%.2f, Signal=%s, Conf=%.1f%%, Volume=%.0f\n",
				bt.currentIndex, bt.currentPrice, analysis.Signal, analysis.Confidence, currentCandle.Volume)
		}
	}

	// ปิด position สุดท้าย
	if bt.position != nil {
		bt.close15mPosition("End of Backtest")
	}

	return bt.generate15mResults()
}

// analyze15mMarket วิเคราะห์ตลาดด้วย Triple EMA Momentum Strategy
func (bt *Backtester) analyze15mMarket() *SuperTrendAnalysis {
	if bt.currentIndex < 50 {
		return nil
	}

	// เอาข้อมูล 50 แท่งล่าสุดสำหรับการคำนวณ indicators
	data := bt.ohlcvData[bt.currentIndex-49 : bt.currentIndex+1]

	// คำนวณ Triple EMA: Fast=5, Mid=13, Slow=21 (จากกลยุทธ์ Scalping ที่ดีที่สุด)
	emaFast := bt.calculateEMA(data, 5)  // Fast EMA for quick signals
	emaMid := bt.calculateEMA(data, 13)  // Mid EMA for confirmation
	emaSlow := bt.calculateEMA(data, 21) // Slow EMA for trend direction

	// คำนวณ RSI และ indicators อื่นๆ
	rsi := bt.calculateRSI(data, 14)
	atr := bt.calculateATR(data, 14)
	avgVolume := bt.calculateAvgVolume(data, 20)

	currentCandle := data[len(data)-1]

	// กำหนดสัญญาณ
	signal := "NEUTRAL"
	confidence := 50.0

	// Bullish Momentum: Fast > Mid > Slow (Triple EMA Alignment) - เงื่อนไขง่าๆ
	if emaFast > emaMid && emaMid > emaSlow &&
		bt.currentPrice > emaFast*0.999 && // Price near or above fast EMA (relaxed)
		rsi > 35 && rsi < 80 { // RSI in reasonable range (relaxed)

		signal = "LONG"
		confidence = 70.0

		// เพิ่ม confidence สำหรับ momentum แข็งแกร่ง
		if currentCandle.Close > currentCandle.Open { // Green candle
			confidence += 10.0
		}

		// เพิ่ม confidence ถ้า EMA spread กว้าง (momentum แข็งแกร่ง)
		emaSpread := ((emaFast - emaSlow) / emaSlow) * 100
		if emaSpread > 0.1 { // EMA spread > 0.1% (relaxed)
			confidence += 10.0
		}

		// Volume confirmation (relaxed)
		if currentCandle.Volume > avgVolume*1.1 {
			confidence += 5.0
		}
	}

	// Bearish Momentum: Fast < Mid < Slow (Triple EMA Alignment) - เงื่อนไขง่าๆ
	if emaFast < emaMid && emaMid < emaSlow &&
		bt.currentPrice < emaFast*1.001 && // Price near or below fast EMA (relaxed)
		rsi > 20 && rsi < 65 { // RSI in reasonable range (relaxed)

		signal = "SHORT"
		confidence = 70.0

		// เพิ่ม confidence สำหรับ momentum แข็งแกร่ง
		if currentCandle.Close < currentCandle.Open { // Red candle
			confidence += 10.0
		}

		// เพิ่ม confidence ถ้า EMA spread กว้าง (momentum แข็งแกร่ง)
		emaSpread := ((emaSlow - emaFast) / emaSlow) * 100
		if emaSpread > 0.1 { // EMA spread > 0.1% (relaxed)
			confidence += 10.0
		}

		// Volume confirmation (relaxed)
		if currentCandle.Volume > avgVolume*1.1 {
			confidence += 5.0
		}
	}

	// Scalping opportunities: Quick reversal signals (relaxed conditions)
	if signal == "NEUTRAL" {
		// Bullish scalping: RSI oversold + price near fast EMA
		if rsi < 40 && bt.currentPrice >= emaFast*0.995 && bt.currentPrice <= emaFast*1.005 &&
			emaFast > emaSlow {
			signal = "LONG"
			confidence = 65.0
		}

		// Bearish scalping: RSI overbought + price near fast EMA
		if rsi > 60 && bt.currentPrice >= emaFast*0.995 && bt.currentPrice <= emaFast*1.005 &&
			emaFast < emaSlow {
			signal = "SHORT"
			confidence = 65.0
		}
	}

	return &SuperTrendAnalysis{
		Trend:           getTrendFromSignal(signal),
		Signal:          signal,
		Confidence:      confidence,
		RiskRewardRatio: 2.0,     // 1:2 Risk:Reward
		SuperTrendValue: emaMid,  // ใช้ Mid EMA เป็น reference
		EMA100:          emaSlow, // Long-term trend
		CurrentPrice:    bt.currentPrice,
		ATR:             atr,
	}
}

// lookFor15mEntry หาจุดเข้าด้วยกลยุทธ์ Triple EMA Momentum
func (bt *Backtester) lookFor15mEntry(analysis *SuperTrendAnalysis) {
	if bt.position != nil {
		return
	}

	// เข้า LONG position (ใช้ confidence threshold 75% สำหรับความแม่นยำ)
	if analysis.Signal == "LONG" && analysis.Confidence >= 75.0 {
		bt.open15mPosition("LONG", analysis)
	}

	// เข้า SHORT position (Future trading)
	if analysis.Signal == "SHORT" && analysis.Confidence >= 75.0 {
		bt.open15mPosition("SHORT", analysis)
	}
}

// open15mPosition เปิด position ด้วย Triple EMA Momentum Strategy
func (bt *Backtester) open15mPosition(side string, analysis *SuperTrendAnalysis) {
	// Future Trading Parameters
	leverage := 5.0 // ใช้ Leverage 5x

	var stopLoss, takeProfit float64

	// คำนวณ Stop Loss และ Take Profit โดยใช้ ATR
	atrMultiplier := 1.5 // ใช้ 1.5x ATR สำหรับ stop loss

	if side == "LONG" {
		// LONG Position - ใช้ ATR-based stops
		stopLoss = bt.currentPrice - (analysis.ATR * atrMultiplier)
		takeProfit = bt.currentPrice + (analysis.ATR * atrMultiplier * analysis.RiskRewardRatio)
	} else {
		// SHORT Position - ใช้ ATR-based stops
		stopLoss = bt.currentPrice + (analysis.ATR * atrMultiplier)
		takeProfit = bt.currentPrice - (analysis.ATR * atrMultiplier * analysis.RiskRewardRatio)
	}

	// Risk per trade 1.0% (ลดลงจาก 1.5% เพื่อความปลอดภัย)
	riskAmount := bt.currentCapital * 0.01

	var positionSize float64
	if side == "LONG" {
		positionSize = (riskAmount * leverage) / (bt.currentPrice - stopLoss)
	} else {
		positionSize = (riskAmount * leverage) / (stopLoss - bt.currentPrice)
	}

	// ป้องกันใช้เงินเกิน 70% (Future trading ต้องระวัง margin call)
	maxValue := bt.currentCapital * 0.7 * leverage
	if positionSize*bt.currentPrice > maxValue {
		positionSize = maxValue / bt.currentPrice
	}

	// คำนวณค่าธรรมเนียม Future (ต่ำกว่า Spot)
	entryCommission := bt.currentPrice * positionSize * 0.0005 // 0.05% สำหรับ Future

	// สร้าง position
	bt.position = &BacktestPosition{
		Symbol:      bt.symbol,
		Side:        side,
		EntryTime:   bt.currentTime,
		EntryPrice:  bt.currentPrice,
		Quantity:    positionSize,
		StopLoss:    stopLoss,
		TakeProfit:  takeProfit,
		EntryReason: fmt.Sprintf("15m %s: Triple EMA Momentum, Conf=%.1f%%, Lev=%.1fx", side, analysis.Confidence, leverage),
	}

	// หักค่าธรรมเนียมเปิด position
	bt.currentCapital -= entryCommission

	// บันทึก trade
	bt.record15mTrade("OPEN", fmt.Sprintf("15m Entry: %s (Triple EMA, %dx Lev)", side, int(leverage)))

	var slPct, tpPct float64
	if side == "LONG" {
		slPct = ((stopLoss - bt.currentPrice) / bt.currentPrice) * 100
		tpPct = ((takeProfit - bt.currentPrice) / bt.currentPrice) * 100
	} else {
		slPct = ((bt.currentPrice - stopLoss) / bt.currentPrice) * 100
		tpPct = ((bt.currentPrice - takeProfit) / bt.currentPrice) * 100
	}

	fmt.Printf("🎯 เปิด 15m Triple EMA %s: @ $%.2f, SL: $%.2f (%.2f%%), TP: $%.2f (%.2f%%), Size: %.6f, Lev: %.1fx, Fee: $%.2f\n",
		side, bt.currentPrice, stopLoss, slPct, takeProfit, tpPct, positionSize, leverage, entryCommission)
}

// check15mPosition ตรวจสอบ position สำหรับ Triple EMA Strategy
func (bt *Backtester) check15mPosition() {
	if bt.position == nil {
		return
	}

	shouldExit := false
	exitReason := ""

	if bt.position.Side == "LONG" {
		// LONG Position checks
		if bt.currentPrice <= bt.position.StopLoss {
			shouldExit = true
			exitReason = "Stop Loss Hit"
		}
		if bt.currentPrice >= bt.position.TakeProfit {
			shouldExit = true
			exitReason = "Take Profit Hit"
		}

		// Trailing Stop Logic - ใช้ EMA เป็น trailing stop
		analysis := bt.analyze15mMarket()
		if analysis != nil {
			// ถ้าราคาเหนือ fast EMA มาก และมี profit > 2%
			currentPnL := ((bt.currentPrice - bt.position.EntryPrice) / bt.position.EntryPrice) * 100
			if currentPnL > 2.0 && bt.currentPrice > analysis.SuperTrendValue {
				// ใช้ Mid EMA เป็น trailing stop
				newStopLoss := analysis.SuperTrendValue * 0.995 // 0.5% buffer
				if newStopLoss > bt.position.StopLoss {
					bt.position.StopLoss = newStopLoss
				}
			}
		}

		// EMA Signal Reversal Exit
		if analysis != nil && analysis.Signal == "SHORT" && analysis.Confidence >= 80.0 {
			shouldExit = true
			exitReason = "EMA Signal Reversal"
		}
	} else {
		// SHORT Position checks
		if bt.currentPrice >= bt.position.StopLoss {
			shouldExit = true
			exitReason = "Stop Loss Hit"
		}
		if bt.currentPrice <= bt.position.TakeProfit {
			shouldExit = true
			exitReason = "Take Profit Hit"
		}

		// Trailing Stop Logic - ใช้ EMA เป็น trailing stop
		analysis := bt.analyze15mMarket()
		if analysis != nil {
			// ถ้าราคาต่ำกว่า fast EMA มาก และมี profit > 2%
			currentPnL := ((bt.position.EntryPrice - bt.currentPrice) / bt.position.EntryPrice) * 100
			if currentPnL > 2.0 && bt.currentPrice < analysis.SuperTrendValue {
				// ใช้ Mid EMA เป็น trailing stop
				newStopLoss := analysis.SuperTrendValue * 1.005 // 0.5% buffer
				if newStopLoss < bt.position.StopLoss {
					bt.position.StopLoss = newStopLoss
				}
			}
		}

		// EMA Signal Reversal Exit
		if analysis != nil && analysis.Signal == "LONG" && analysis.Confidence >= 80.0 {
			shouldExit = true
			exitReason = "EMA Signal Reversal"
		}
	}

	// Time exit - ลดเวลาการถือ position
	if bt.currentTime.Sub(bt.position.EntryTime).Hours() > 4 {
		shouldExit = true
		exitReason = "Time Exit (4h)"
	}

	// Emergency exit based on drawdown
	currentPnL := 0.0
	if bt.position.Side == "LONG" {
		currentPnL = ((bt.currentPrice - bt.position.EntryPrice) / bt.position.EntryPrice) * 100
	} else {
		currentPnL = ((bt.position.EntryPrice - bt.currentPrice) / bt.position.EntryPrice) * 100
	}

	if currentPnL <= -3.0 { // Emergency exit ถ้าขาดทุนเกิน 3%
		shouldExit = true
		exitReason = "Emergency Exit (-3%)"
	}

	if shouldExit {
		bt.close15mPosition(exitReason)
	}
}

// close15mPosition ปิด position สำหรับ Triple EMA Strategy
func (bt *Backtester) close15mPosition(reason string) {
	if bt.position == nil {
		return
	}

	var pnl float64

	if bt.position.Side == "LONG" {
		// คำนวณ PnL สำหรับ LONG
		pnl = (bt.currentPrice - bt.position.EntryPrice) * bt.position.Quantity
	} else {
		// คำนวณ PnL สำหรับ SHORT
		pnl = (bt.position.EntryPrice - bt.currentPrice) * bt.position.Quantity
	}

	// คำนวณค่าธรรมเนียม Future (เปิด + ปิด position = 0.05% * 2 = 0.1%)
	entryCommission := bt.position.EntryPrice * bt.position.Quantity * 0.0005
	exitCommission := bt.currentPrice * bt.position.Quantity * 0.0005
	totalCommission := entryCommission + exitCommission

	netPnL := pnl - totalCommission

	// อัปเดตเงินทุน
	bt.currentCapital += netPnL

	// บันทึก trade
	bt.record15mTrade("CLOSE", reason)

	// แสดงผล
	pnlPct := (netPnL / (bt.position.EntryPrice * bt.position.Quantity)) * 100
	fmt.Printf("🔹 ปิด 15m Triple EMA %s: PnL $%.2f (%.2f%%) - Commission: $%.2f - %s\n",
		bt.position.Side, netPnL, pnlPct, totalCommission, reason)

	// ลบ position
	bt.position = nil
}

// record15mTrade บันทึก trade สำหรับ Future Trading
func (bt *Backtester) record15mTrade(action, reason string) {
	bt.tradeID++

	var exitTime time.Time
	var exitPrice, pnl, netPnL, totalCommission float64

	if action == "CLOSE" && bt.position != nil {
		exitTime = bt.currentTime
		exitPrice = bt.currentPrice

		if bt.position.Side == "LONG" {
			pnl = (bt.currentPrice - bt.position.EntryPrice) * bt.position.Quantity
		} else {
			pnl = (bt.position.EntryPrice - bt.currentPrice) * bt.position.Quantity
		}

		// คำนวณค่าธรรมเนียม Future (0.05% ต่อรอบ)
		entryCommission := bt.position.EntryPrice * bt.position.Quantity * 0.0005
		exitCommission := bt.currentPrice * bt.position.Quantity * 0.0005
		totalCommission = entryCommission + exitCommission

		netPnL = pnl - totalCommission
	}

	trade := BacktestTrade{
		ID:          bt.tradeID,
		Symbol:      bt.symbol,
		Side:        bt.position.Side,
		EntryTime:   bt.position.EntryTime,
		ExitTime:    exitTime,
		EntryPrice:  bt.position.EntryPrice,
		ExitPrice:   exitPrice,
		Quantity:    bt.position.Quantity,
		PnL:         pnl,
		PnLPct:      (pnl / (bt.position.EntryPrice * bt.position.Quantity)) * 100,
		Commission:  totalCommission,
		NetPnL:      netPnL,
		EntryReason: bt.position.EntryReason,
		ExitReason:  reason,
	}

	if action == "CLOSE" {
		trade.Duration = bt.currentTime.Sub(bt.position.EntryTime)
	}

	bt.trades = append(bt.trades, trade)
}

// generate15mResults สร้างผลลัพธ์
func (bt *Backtester) generate15mResults() *BacktestResult {
	// คำนวณสถิติ
	totalReturn := ((bt.currentCapital - bt.initialCapital) / bt.initialCapital) * 100

	// นับ trades
	var winTrades, lossTrades int
	var totalPnL float64

	for _, trade := range bt.trades {
		if trade.NetPnL != 0 { // Trade ที่ปิดแล้ว
			totalPnL += trade.NetPnL
			if trade.NetPnL > 0 {
				winTrades++
			} else {
				lossTrades++
			}
		}
	}

	totalTrades := winTrades + lossTrades
	var winRate float64
	if totalTrades > 0 {
		winRate = (float64(winTrades) / float64(totalTrades)) * 100
	}

	// คำนวณ max drawdown (simplified)
	maxDrawdown := 0.0

	return &BacktestResult{
		Symbol:         bt.symbol,
		StartDate:      bt.startDate,
		EndDate:        bt.endDate,
		InitialCapital: bt.initialCapital,
		FinalCapital:   bt.currentCapital,
		TotalReturn:    bt.currentCapital - bt.initialCapital,
		TotalReturnPct: totalReturn,
		TotalTrades:    totalTrades,
		WinningTrades:  winTrades,
		LosingTrades:   lossTrades,
		WinRate:        winRate,
		MaxDrawdown:    maxDrawdown,
		MaxDrawdownPct: maxDrawdown,
		Trades:         bt.trades,
		DailyReturns:   bt.dailyReturns,
	}
}

// Helper functions สำหรับการคำนวณ indicators
func (bt *Backtester) calculateRSI(data []OHLCV, period int) float64 {
	if len(data) < period+1 {
		return 50.0
	}

	var gains, losses float64
	for i := len(data) - period; i < len(data); i++ {
		change := data[i].Close - data[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += -change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return rsi
}

func (bt *Backtester) calculateMACD(data []OHLCV, fast, slow, signal int) (float64, float64, float64) {
	if len(data) < slow {
		return 0, 0, 0
	}

	emaFast := bt.calculateEMA(data, fast)
	emaSlow := bt.calculateEMA(data, slow)
	macd := emaFast - emaSlow

	// Simplified signal line (would need more complex calculation in reality)
	macdSignal := macd * 0.9
	macdHist := macd - macdSignal

	return macd, macdSignal, macdHist
}

func (bt *Backtester) calculateAvgVolume(data []OHLCV, period int) float64 {
	if len(data) < period {
		return data[len(data)-1].Volume
	}

	var sum float64
	for i := len(data) - period; i < len(data); i++ {
		sum += data[i].Volume
	}
	return sum / float64(period)
}

func (bt *Backtester) calculateEMA(data []OHLCV, period int) float64 {
	if len(data) < period {
		return data[len(data)-1].Close
	}

	multiplier := 2.0 / (float64(period) + 1.0)
	ema := data[len(data)-period].Close

	for i := len(data) - period + 1; i < len(data); i++ {
		ema = (data[i].Close * multiplier) + (ema * (1 - multiplier))
	}

	return ema
}

func (bt *Backtester) calculateATR(data []OHLCV, period int) float64 {
	if len(data) < period+1 {
		return (data[len(data)-1].High - data[len(data)-1].Low)
	}

	var trSum float64
	for i := len(data) - period; i < len(data); i++ {
		high := data[i].High
		low := data[i].Low
		prevClose := data[i-1].Close

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)

		tr := math.Max(tr1, math.Max(tr2, tr3))
		trSum += tr
	}

	return trSum / float64(period)
}

func getTrendFromSignal(signal string) int {
	if signal == "LONG" {
		return 1
	} else if signal == "SHORT" {
		return -1
	}
	return 0
}
