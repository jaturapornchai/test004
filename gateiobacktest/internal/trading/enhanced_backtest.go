package trading

import (
	"fmt"
	"time"
)

// AggressiveStrategy - เมธอดใหม่สำหรับกลยุทธ์ที่เน้นกำไรสูง
func (bt *Backtester) RunAggressiveBacktest() *BacktestResult {
	fmt.Printf("\n🚀 เริ่มต้น Aggressive Profit Strategy %s\n", bt.symbol)
	fmt.Printf("💰 เงินทุนเริ่มต้น: $%.2f\n", bt.initialCapital)
	fmt.Printf("🎯 เป้าหมาย: เพิ่มกำไรให้สูงกว่าเดิม\n")

	for bt.currentIndex = 100; bt.currentIndex < len(bt.ohlcvData); bt.currentIndex++ {
		currentCandle := bt.ohlcvData[bt.currentIndex]
		bt.currentTime = time.Unix(currentCandle.Timestamp, 0)
		bt.currentPrice = currentCandle.Close

		// วิเคราะห์ตลาด
		analysis := bt.analyzeMarket()
		if analysis == nil {
			continue
		}

		// จัดการ position ที่มีอยู่ด้วยกลยุทธ์ aggressive
		bt.checkAggressivePosition()

		// หาจุดเข้าใหม่ด้วยกลยุทธ์ที่เน้นกำไร
		bt.lookForAggressiveEntry(analysis)
	}

	// ปิด position สุดท้าย
	if bt.position != nil {
		bt.closeAggressivePosition("End of Backtest")
	}

	return bt.generateAggressiveResults()
}

// lookForAggressiveEntry หาจุดเข้าที่ให้กำไรสูง
func (bt *Backtester) lookForAggressiveEntry(analysis *SuperTrendAnalysis) {
	// ไม่เข้าถ้ามี position อยู่แล้ว
	if bt.position != nil {
		return
	}

	currentCandle := bt.ohlcvData[bt.currentIndex]

	// เงื่อนไขเข้า Long ที่เข้มงวดขึ้น
	if analysis.Trend == 1 && // SuperTrend bullish
		currentCandle.Close > analysis.EMA100 && // ราคาเหนือ EMA100
		analysis.Confidence >= 75 && // Confidence สูง
		bt.isVolumeStrong(currentCandle) && // Volume แข็งแกร่ง
		bt.isPriceMomentumStrong() { // Momentum แข็งแกร่ง

		bt.openAggressivePosition("LONG", analysis)
	}
}

// isVolumeStrong ตรวจสอบ Volume แข็งแกร่ง
func (bt *Backtester) isVolumeStrong(currentCandle OHLCV) bool {
	if bt.currentIndex < 20 {
		return false
	}

	// คำนวณ Volume เฉลี่ย 20 periods
	var avgVolume float64
	for i := bt.currentIndex - 19; i <= bt.currentIndex; i++ {
		avgVolume += bt.ohlcvData[i].Volume
	}
	avgVolume /= 20

	// Volume ต้องมากกว่าเฉลี่ย 30%
	return currentCandle.Volume > avgVolume*1.3
}

// isPriceMomentumStrong ตรวจสอบ Momentum ราคา
func (bt *Backtester) isPriceMomentumStrong() bool {
	if bt.currentIndex < 5 {
		return false
	}

	// ตรวจสอบว่าราคาปิดสูงขึ้นใน 3 ใน 5 periods สุดท้าย
	upCount := 0
	for i := bt.currentIndex - 4; i < bt.currentIndex; i++ {
		if bt.ohlcvData[i+1].Close > bt.ohlcvData[i].Close {
			upCount++
		}
	}

	return upCount >= 3
}

// openAggressivePosition เปิด position ที่เน้นกำไร
func (bt *Backtester) openAggressivePosition(side string, analysis *SuperTrendAnalysis) {
	// Stop Loss ใกล้ SuperTrend มากขึ้น (เสี่ยงสูงแต่ Risk:Reward ดีขึ้น)
	stopLoss := analysis.SuperTrendValue * 0.997 // ใกล้ SuperTrend 0.3%

	// Take Profit ที่ 1:4 Risk:Reward
	riskDistance := bt.currentPrice - stopLoss
	takeProfit := bt.currentPrice + (riskDistance * 4.0) // 1:4 R:R

	// เพิ่มขนาด position เป็น 3% (จากเดิม 2%)
	riskAmount := bt.currentCapital * 0.03
	positionSize := riskAmount / riskDistance

	// ป้องกันใช้เงินเกิน 95%
	maxValue := bt.currentCapital * 0.95
	if positionSize*bt.currentPrice > maxValue {
		positionSize = maxValue / bt.currentPrice
	}

	// สร้าง position
	bt.position = &BacktestPosition{
		Symbol:      bt.symbol,
		Side:        side,
		EntryTime:   bt.currentTime,
		EntryPrice:  bt.currentPrice,
		Quantity:    positionSize,
		StopLoss:    stopLoss,
		TakeProfit:  takeProfit,
		EntryReason: fmt.Sprintf("Aggressive %s: Conf=%.1f%%, R:R=1:4", side, analysis.Confidence),
	}

	// บันทึก trade
	bt.recordAggressiveTrade("OPEN", fmt.Sprintf("Aggressive Entry: %s (R:R 1:4)", side))

	fmt.Printf("📈 เปิด Aggressive Position: %s @ $%.2f, SL: $%.2f (%.2f%%), TP: $%.2f (%.2f%%), Size: %.6f\n",
		side, bt.currentPrice, stopLoss,
		((stopLoss-bt.currentPrice)/bt.currentPrice)*100,
		takeProfit,
		((takeProfit-bt.currentPrice)/bt.currentPrice)*100,
		positionSize)
}

// recordAggressiveTrade บันทึก trade
func (bt *Backtester) recordAggressiveTrade(action, reason string) {
	bt.tradeID++

	var exitTime time.Time
	var exitPrice, pnl, netPnL float64

	if action == "CLOSE" && bt.position != nil {
		exitTime = bt.currentTime
		exitPrice = bt.currentPrice
		pnl = (bt.currentPrice - bt.position.EntryPrice) * bt.position.Quantity
		commission := bt.currentPrice * bt.position.Quantity * bt.commission
		netPnL = pnl - commission
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
		Commission:  bt.currentPrice * bt.position.Quantity * bt.commission,
		NetPnL:      netPnL,
		EntryReason: bt.position.EntryReason,
		ExitReason:  reason,
	}

	if action == "CLOSE" {
		trade.Duration = bt.currentTime.Sub(bt.position.EntryTime)
	}

	bt.trades = append(bt.trades, trade)
}

// checkAggressivePosition ตรวจสอบ position ด้วยกลยุทธ์เน้นกำไร
func (bt *Backtester) checkAggressivePosition() {
	if bt.position == nil {
		return
	}

	// อัปเดต trailing stop ที่เข้มงวดขึ้น
	bt.updateAggressiveTrailingStop()

	// เงื่อนไขออก
	shouldExit := false
	exitReason := ""

	// Stop Loss
	if bt.currentPrice <= bt.position.StopLoss {
		shouldExit = true
		exitReason = "Stop Loss"
	}

	// Take Profit
	if bt.currentPrice >= bt.position.TakeProfit {
		shouldExit = true
		exitReason = "Take Profit"
	}

	// ออกเมื่อมีกำไร 3% และ trend เริ่มอ่อนแอ
	currentPnL := ((bt.currentPrice - bt.position.EntryPrice) / bt.position.EntryPrice) * 100
	if currentPnL >= 3.0 {
		analysis := bt.analyzeMarket()
		if analysis != nil && analysis.Trend != 1 {
			shouldExit = true
			exitReason = "Profit Protection (3%+ with trend weakness)"
		}
	}

	// Time exit (ถือไม่เกิน 18 ชั่วโมง)
	if bt.currentTime.Sub(bt.position.EntryTime).Hours() > 18 {
		shouldExit = true
		exitReason = "Time Exit (18h)"
	}

	if shouldExit {
		bt.closeAggressivePosition(exitReason)
	}
}

// updateAggressiveTrailingStop อัปเดต trailing stop แบบ aggressive
func (bt *Backtester) updateAggressiveTrailingStop() {
	if bt.position == nil || bt.position.Side != "LONG" {
		return
	}

	// Trailing stop ที่ 1% (เข้มงวดขึ้น)
	newTrailingStop := bt.currentPrice * 0.99

	if newTrailingStop > bt.position.StopLoss {
		bt.position.StopLoss = newTrailingStop
	}
}

// closeAggressivePosition ปิด position
func (bt *Backtester) closeAggressivePosition(reason string) {
	if bt.position == nil {
		return
	}

	// คำนวณ PnL
	pnl := (bt.currentPrice - bt.position.EntryPrice) * bt.position.Quantity
	commission := bt.currentPrice * bt.position.Quantity * bt.commission
	netPnL := pnl - commission

	// อัปเดตเงินทุน
	bt.currentCapital += netPnL

	// บันทึก trade
	bt.recordAggressiveTrade("CLOSE", reason)

	// แสดงผล
	pnlPct := (netPnL / (bt.position.EntryPrice * bt.position.Quantity)) * 100
	fmt.Printf("📉 ปิด Position: PnL $%.2f (%.2f%%) - %s\n", netPnL, pnlPct, reason)

	// ลบ position
	bt.position = nil
}

// generateAggressiveResults สร้างผลลัพธ์สำหรับ aggressive strategy
func (bt *Backtester) generateAggressiveResults() *BacktestResult {
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
	maxDrawdown := 0.0 // ใช้ค่า 0 ก่อน

	return &BacktestResult{
		Symbol:         bt.symbol,
		StartDate:      bt.startDate,
		EndDate:        bt.endDate,
		InitialCapital: bt.initialCapital,
		FinalCapital:   bt.currentCapital,
		TotalReturn:    totalReturn,
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
