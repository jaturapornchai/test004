package trading

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// BacktestResult ผลการ backtest
type BacktestResult struct {
	Symbol         string          `json:"symbol"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	InitialCapital float64         `json:"initial_capital"`
	FinalCapital   float64         `json:"final_capital"`
	TotalReturn    float64         `json:"total_return"`
	TotalReturnPct float64         `json:"total_return_pct"`
	TotalTrades    int             `json:"total_trades"`
	WinningTrades  int             `json:"winning_trades"`
	LosingTrades   int             `json:"losing_trades"`
	WinRate        float64         `json:"win_rate"`
	MaxDrawdown    float64         `json:"max_drawdown"`
	MaxDrawdownPct float64         `json:"max_drawdown_pct"`
	Trades         []BacktestTrade `json:"trades"`
	DailyReturns   []DailyReturn   `json:"daily_returns"`
}

// BacktestTrade การเทรดใน backtest
type BacktestTrade struct {
	ID          int           `json:"id"`
	Symbol      string        `json:"symbol"`
	Side        string        `json:"side"` // "LONG" or "SHORT"
	EntryTime   time.Time     `json:"entry_time"`
	ExitTime    time.Time     `json:"exit_time"`
	EntryPrice  float64       `json:"entry_price"`
	ExitPrice   float64       `json:"exit_price"`
	Quantity    float64       `json:"quantity"`
	PnL         float64       `json:"pnl"`
	PnLPct      float64       `json:"pnl_pct"`
	Commission  float64       `json:"commission"`
	NetPnL      float64       `json:"net_pnl"`
	Duration    time.Duration `json:"duration"`
	EntryReason string        `json:"entry_reason"`
	ExitReason  string        `json:"exit_reason"`
	StopLoss    float64       `json:"stop_loss"`
	TakeProfit  float64       `json:"take_profit"`
}

// DailyReturn ผลตอบแทนรายวัน
type DailyReturn struct {
	Date   time.Time `json:"date"`
	Return float64   `json:"return"`
	Equity float64   `json:"equity"`
}

// Position สำหรับ backtest
type BacktestPosition struct {
	Symbol      string    `json:"symbol"`
	Side        string    `json:"side"`
	EntryTime   time.Time `json:"entry_time"`
	EntryPrice  float64   `json:"entry_price"`
	Quantity    float64   `json:"quantity"`
	StopLoss    float64   `json:"stop_loss"`
	TakeProfit  float64   `json:"take_profit"`
	EntryReason string    `json:"entry_reason"`
}

// Position management แบบใหม่
type PositionManager struct {
	maxPositions  int     // จำนวน position สูงสุดที่สามารถเพิ่มได้
	pyramidPnLMin float64 // กำไรขั้นต่ำก่อนเพิ่ม position (%)
	totalRisk     float64 // ความเสี่ยงรวมทั้งหมด
}

// StructuredPosition position ที่มีโครงสร้าง
type StructuredPosition struct {
	*BacktestPosition
	PositionLevel int     // ระดับของ position (1, 2, 3...)
	BaseQuantity  float64 // ขนาด position พื้นฐาน
	ProfitTarget  float64 // เป้าหมายกำไรของ position นี้
}

// Backtester หลักของระบบ backtest
type Backtester struct {
	symbol         string
	startDate      time.Time
	endDate        time.Time
	initialCapital float64
	commission     float64 // อัตราค่าคอมมิชชั่น (0.001 = 0.1%)

	// ข้อมูลปัจจุบัน
	currentTime    time.Time
	currentCapital float64
	currentPrice   float64
	position       *BacktestPosition

	// ข้อมูลราคา
	ohlcvData    []OHLCV
	currentIndex int

	// ผลลัพธ์
	trades       []BacktestTrade
	dailyReturns []DailyReturn
	tradeID      int

	// ตัววิเคราะห์
	indicators *Indicators
	aiClient   *AIClient

	// ตัวจัดการ Position ใหม่
	positionManager *PositionManager
}

// NewBacktester สร้าง backtester ใหม่
func NewBacktester(symbol string, daysBack int, initialCapital float64, deepseekKey string) (*Backtester, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -daysBack)

	// สร้าง AI client
	aiClient, err := NewAIClient(deepseekKey)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง AI client ได้: %v", err)
	}

	// สร้าง Position Manager
	positionManager := NewPositionManager()

	return &Backtester{
		symbol:          symbol,
		startDate:       startDate,
		endDate:         endDate,
		initialCapital:  initialCapital,
		currentCapital:  initialCapital,
		commission:      0.001, // 0.1% commission
		trades:          make([]BacktestTrade, 0),
		dailyReturns:    make([]DailyReturn, 0),
		indicators:      NewIndicators(),
		aiClient:        aiClient,
		positionManager: positionManager,
	}, nil
}

// NewBacktesterSimple สร้าง backtester ใหม่โดยไม่ใช้ AI
func NewBacktesterSimple(symbol string, daysBack int, initialCapital float64) (*Backtester, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -daysBack)

	// สร้าง Position Manager
	positionManager := NewPositionManager()

	return &Backtester{
		symbol:          symbol,
		startDate:       startDate,
		endDate:         endDate,
		initialCapital:  initialCapital,
		currentCapital:  initialCapital,
		commission:      0.0005, // 0.05% commission for futures
		trades:          make([]BacktestTrade, 0),
		dailyReturns:    make([]DailyReturn, 0),
		indicators:      NewIndicators(),
		positionManager: positionManager,
	}, nil
}

// สร้าง PositionManager ใหม่
func NewPositionManager() *PositionManager {
	return &PositionManager{
		maxPositions:  3,    // สูงสุด 3 level
		pyramidPnLMin: 1.0,  // ต้องมีกำไร 1% ก่อนเพิ่ม
		totalRisk:     0.06, // ความเสี่ยงรวม 6%
	}
}

// LoadHistoricalData โหลดข้อมูลราคาย้อนหลัง
func (bt *Backtester) LoadHistoricalData(ohlcvData []OHLCV) {
	// เรียงข้อมูลตามเวลา
	sort.Slice(ohlcvData, func(i, j int) bool {
		return ohlcvData[i].Timestamp < ohlcvData[j].Timestamp
	})

	bt.ohlcvData = ohlcvData
	fmt.Printf("📊 โหลดข้อมูลราคา %s จำนวน %d แท่งเทียน\n", bt.symbol, len(ohlcvData))
	fmt.Printf("📅 ช่วงเวลา: %s ถึง %s\n",
		time.Unix(ohlcvData[0].Timestamp, 0).Format("2006-01-02 15:04:05"),
		time.Unix(ohlcvData[len(ohlcvData)-1].Timestamp, 0).Format("2006-01-02 15:04:05"))
}

// LoadOHLCVData โหลดข้อมูล OHLCV
func (bt *Backtester) LoadOHLCVData(data []OHLCV, startDate, endDate time.Time) {
	bt.ohlcvData = data
	bt.startDate = startDate
	bt.endDate = endDate
	bt.currentCapital = bt.initialCapital
	bt.currentIndex = 0
}

// RunBacktest รัน backtest
func (bt *Backtester) RunBacktest() (*BacktestResult, error) {
	fmt.Printf("🚀 เริ่มต้น Backtest สำหรับ %s\n", bt.symbol)
	fmt.Printf("💰 เงินทุนเริ่มต้น: $%.2f\n", bt.initialCapital)
	fmt.Printf("📅 ระยะเวลา: %d วัน (%s ถึง %s)\n",
		int(bt.endDate.Sub(bt.startDate).Hours()/24),
		bt.startDate.Format("2006-01-02"),
		bt.endDate.Format("2006-01-02"))

	if len(bt.ohlcvData) == 0 {
		return nil, fmt.Errorf("ไม่มีข้อมูลราคาสำหรับ backtest")
	}

	maxDrawdown := 0.0
	peak := bt.initialCapital

	// วนลูปผ่านข้อมูลราคาทั้งหมด
	for bt.currentIndex = 50; bt.currentIndex < len(bt.ohlcvData); bt.currentIndex++ {
		candle := bt.ohlcvData[bt.currentIndex]
		bt.currentTime = time.Unix(candle.Timestamp, 0)
		bt.currentPrice = candle.Close

		// ตรวจสอบ position ที่เปิดอยู่
		bt.checkPosition()

		// หาโอกาสเทรดใหม่ (ถ้าไม่มี position)
		if bt.position == nil {
			bt.lookForEntry()
		}

		// คำนวณ drawdown
		if bt.currentCapital > peak {
			peak = bt.currentCapital
		}
		drawdown := peak - bt.currentCapital
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}

		// บันทึกผลตอบแทนรายวัน (ทุก 24 ชั่วโมง)
		if bt.currentIndex%24 == 0 {
			bt.dailyReturns = append(bt.dailyReturns, DailyReturn{
				Date:   bt.currentTime,
				Return: (bt.currentCapital - bt.initialCapital) / bt.initialCapital * 100,
				Equity: bt.currentCapital,
			})
		}
	}

	// ปิด position ถ้ามี (ตอนจบ backtest)
	if bt.position != nil {
		bt.closePosition("END_OF_BACKTEST")
	}

	// คำนวณสถิติ
	result := bt.calculateResults(maxDrawdown)

	return result, nil
}

// checkPosition ตรวจสอบ position ที่เปิดอยู่ - ปรับปรุงใหม่
func (bt *Backtester) checkPosition() {
	if bt.position == nil {
		return
	}

	// ตรวจสอบ Stop Loss (แบบ hard stop)
	if bt.position.Side == "LONG" && bt.currentPrice <= bt.position.StopLoss {
		bt.closePosition("STOP_LOSS")
		return
	}
	if bt.position.Side == "SHORT" && bt.currentPrice >= bt.position.StopLoss {
		bt.closePosition("STOP_LOSS")
		return
	}

	// ตรวจสอบ Take Profit (แบบ hard target)
	if bt.position.Side == "LONG" && bt.currentPrice >= bt.position.TakeProfit {
		bt.closePosition("TAKE_PROFIT")
		return
	}
	if bt.position.Side == "SHORT" && bt.currentPrice <= bt.position.TakeProfit {
		bt.closePosition("TAKE_PROFIT")
		return
	}

	// ตรวจสอบ Trailing Stop (ใหม่)
	bt.updateTrailingStop()

	// ตรวจสอบ early exit signals (ทุก 2 ชั่วโมง)
	if bt.currentIndex%2 == 0 {
		bt.checkEarlyExit()
	}
}

// lookForEntry หาโอกาสเข้า position ใหม่ - กลยุทธ์ใหม่ + debug
func (bt *Backtester) lookForEntry() {
	// ต้องมีข้อมูลย้อนหลังเพียงพอ
	if bt.currentIndex < 100 {
		return
	}

	// วิเคราะห์ทางเทคนิค
	analysis := bt.analyzeMarket()
	if analysis == nil {
		return
	}

	// Debug: แสดงข้อมูลการวิเคราะห์
	if bt.currentIndex%24 == 0 { // แสดงทุก 24 ชั่วโมง
		fmt.Printf("🔍 Debug [%s]: Trend=%d, Signal=%s, Confidence=%.1f%%, RR=%.2f\n",
			bt.currentTime.Format("01-02 15:04"), analysis.Trend, analysis.Signal,
			analysis.Confidence, analysis.RiskRewardRatio)
	}

	// กลยุทธ์ใหม่: ตรวจสอบสัญญาณแข็งแกร่งก่อนเข้าเทรด
	if bt.isStrongSignal(analysis) {
		direction := bt.determineDirection(analysis)
		if direction != "" {
			stopLoss, takeProfit := bt.calculateRiskReward(direction, bt.currentPrice, analysis.ATR)
			fmt.Printf("📊 Strong signal detected: %s at $%.2f\n", direction, bt.currentPrice)
			bt.openPosition(direction, stopLoss, takeProfit, bt.getSignalReason(analysis))
		}
	} else {
		// เพิ่มการตรวจสอบสัญญาณง่าย ๆ เป็น fallback
		if bt.currentIndex%6 == 0 { // ทุก 6 ชั่วโมง
			bt.checkSimpleSignals(analysis)
		}
	}
}

// analyzeMarket วิเคราะห์ตลาด
func (bt *Backtester) analyzeMarket() *SuperTrendAnalysis {
	// เตรียมข้อมูล
	endIdx := bt.currentIndex + 1
	startIdx := endIdx - 100
	if startIdx < 0 {
		startIdx = 0
	}

	recentData := bt.ohlcvData[startIdx:endIdx]

	// คำนวณตัวชี้วัด
	analysis := bt.indicators.AnalyzePivotPointSuperTrend(recentData)
	analysis.CurrentPrice = bt.currentPrice

	return analysis
}

// getAIDecision ขอคำแนะนำจาก AI
func (bt *Backtester) getAIDecision(analysis *SuperTrendAnalysis) *AIDecision {
	startIndex := bt.currentIndex - 20
	if startIndex < 0 {
		startIndex = 0
	}
	decision, err := bt.aiClient.AnalyzeOpenPosition(bt.symbol, analysis, bt.ohlcvData[startIndex:bt.currentIndex+1])
	if err != nil {
		fmt.Printf("⚠️ ไม่สามารถขอคำแนะนำจาก AI ได้: %v\n", err)
		return nil
	}

	return decision
}

// openPosition เปิด position ใหม่ - ปรับปรุงการจัดการความเสี่ยง
func (bt *Backtester) openPosition(side string, stopLoss, takeProfit float64, reason string) {
	// คำนวณ risk per trade (2% ของเงินทุน)
	riskAmount := bt.currentCapital * 0.02

	// คำนวณระยะห่าง stop loss
	var stopDistance float64
	if side == "LONG" {
		stopDistance = bt.currentPrice - stopLoss
	} else {
		stopDistance = stopLoss - bt.currentPrice
	}

	// ป้องกันการหาร 0
	if stopDistance <= 0 {
		fmt.Printf("⚠️ Stop distance ไม่ถูกต้อง: %.2f\n", stopDistance)
		return
	}

	// คำนวณขนาด position ตาม risk
	quantity := riskAmount / stopDistance

	// จำกัดขนาด position ไม่เกิน 50% ของเงินทุน
	maxCapital := bt.currentCapital * 0.5
	maxQuantity := maxCapital / bt.currentPrice
	if quantity > maxQuantity {
		quantity = maxQuantity
	}

	// ตรวจสอบว่าเงินทุนเพียงพอ
	requiredCapital := quantity * bt.currentPrice
	if requiredCapital > bt.currentCapital*0.8 {
		fmt.Printf("⚠️ เงินทุนไม่เพียงพอสำหรับการเทรด\n")
		return
	}

	bt.position = &BacktestPosition{
		Symbol:      bt.symbol,
		Side:        side,
		EntryTime:   bt.currentTime,
		EntryPrice:  bt.currentPrice,
		Quantity:    quantity,
		StopLoss:    stopLoss,
		TakeProfit:  takeProfit,
		EntryReason: reason,
	}

	fmt.Printf("📈 เปิด %s: ราคา $%.2f, ปริมาณ %.6f, SL: $%.2f, TP: $%.2f\n",
		side, bt.currentPrice, quantity, stopLoss, takeProfit)
	fmt.Printf("🎯 เหตุผล: %s\n", reason)
	fmt.Printf("💰 ความเสี่ยง: $%.2f (%.2f%%)\n", riskAmount, 2.0)
}

// closePosition ปิด position
func (bt *Backtester) closePosition(reason string) {
	if bt.position == nil {
		return
	}

	// คำนวณ PnL
	var pnl float64
	if bt.position.Side == "LONG" {
		pnl = (bt.currentPrice - bt.position.EntryPrice) * bt.position.Quantity
	} else {
		pnl = (bt.position.EntryPrice - bt.currentPrice) * bt.position.Quantity
	}

	// คำนวณค่าคอมมิชชั่น
	commission := (bt.position.EntryPrice + bt.currentPrice) * bt.position.Quantity * bt.commission
	netPnL := pnl - commission

	// อัปเดตเงินทุน
	bt.currentCapital += netPnL

	// บันทึกการเทรด
	bt.tradeID++
	trade := BacktestTrade{
		ID:          bt.tradeID,
		Symbol:      bt.position.Symbol,
		Side:        bt.position.Side,
		EntryTime:   bt.position.EntryTime,
		ExitTime:    bt.currentTime,
		EntryPrice:  bt.position.EntryPrice,
		ExitPrice:   bt.currentPrice,
		Quantity:    bt.position.Quantity,
		PnL:         pnl,
		PnLPct:      pnl / (bt.position.EntryPrice * bt.position.Quantity) * 100,
		Commission:  commission,
		NetPnL:      netPnL,
		Duration:    bt.currentTime.Sub(bt.position.EntryTime),
		EntryReason: bt.position.EntryReason,
		ExitReason:  reason,
		StopLoss:    bt.position.StopLoss,
		TakeProfit:  bt.position.TakeProfit,
	}

	bt.trades = append(bt.trades, trade)

	status := "📉"
	if netPnL > 0 {
		status = "📈"
	}

	fmt.Printf("%s ปิด %s: ราคา $%.2f, PnL: $%.2f (%.2f%%), เวลา: %v\n",
		status, bt.position.Side, bt.currentPrice, netPnL, trade.PnLPct, trade.Duration)
	fmt.Printf("🎯 เหตุผล: %s\n", reason)
	fmt.Printf("💰 เงินทุนปัจจุบัน: $%.2f\n", bt.currentCapital)

	// ล้าง position
	bt.position = nil
}

// checkExitSignal ตรวจสอบสัญญาณออกจาก AI
func (bt *Backtester) checkExitSignal() {
	analysis := bt.analyzeMarket()
	if analysis == nil {
		return
	}

	// สำหรับการปิด position ให้ใช้วิธีอื่น เช่น ตรวจสอบจากสัญญาณเทคนิค
	// หรือสร้าง position จำลองเพื่อส่งให้ AI
	if analysis.Signal == "NEUTRAL" ||
		(bt.position.Side == "LONG" && analysis.Trend == -1) ||
		(bt.position.Side == "SHORT" && analysis.Trend == 1) {
		bt.closePosition("TECHNICAL_EXIT_SIGNAL")
	}
}

// calculateResults คำนวณผลลัพธ์สุดท้าย
func (bt *Backtester) calculateResults(maxDrawdown float64) *BacktestResult {
	totalReturn := bt.currentCapital - bt.initialCapital
	totalReturnPct := totalReturn / bt.initialCapital * 100

	winningTrades := 0
	losingTrades := 0

	for _, trade := range bt.trades {
		if trade.NetPnL > 0 {
			winningTrades++
		} else {
			losingTrades++
		}
	}

	winRate := 0.0
	if len(bt.trades) > 0 {
		winRate = float64(winningTrades) / float64(len(bt.trades)) * 100
	}

	return &BacktestResult{
		Symbol:         bt.symbol,
		StartDate:      bt.startDate,
		EndDate:        bt.endDate,
		InitialCapital: bt.initialCapital,
		FinalCapital:   bt.currentCapital,
		TotalReturn:    totalReturn,
		TotalReturnPct: totalReturnPct,
		TotalTrades:    len(bt.trades),
		WinningTrades:  winningTrades,
		LosingTrades:   losingTrades,
		WinRate:        winRate,
		MaxDrawdown:    maxDrawdown,
		MaxDrawdownPct: maxDrawdown / bt.initialCapital * 100,
		Trades:         bt.trades,
		DailyReturns:   bt.dailyReturns,
	}
}

// isStrongSignal ตรวจสอบว่าเป็นสัญญาณที่แข็งแกร่งหรือไม่ - ปรับให้ผ่อนคลาย
func (bt *Backtester) isStrongSignal(analysis *SuperTrendAnalysis) bool {
	// ตรวจสอบ momentum และ trend alignment (ลดความเข้มงวด)
	trendAligned := analysis.Trend != 0 && analysis.Signal != "NEUTRAL"
	confidenceOK := analysis.Confidence >= 70.0 // ลดจาก 80 เป็น 70

	// ตรวจสอบการยืนยันจากหลายตัวชี้วัด (ผ่อนคลาย)
	priceAligned := (analysis.Trend == 1 && analysis.CurrentPrice > analysis.SuperTrendValue) ||
		(analysis.Trend == -1 && analysis.CurrentPrice < analysis.SuperTrendValue)

	// ตรวจสอบ risk-reward ratio (ลดเงื่อนไข)
	goodRiskReward := analysis.RiskRewardRatio >= 1.5 // ลดจาก 2.0 เป็น 1.5

	// ตรวจสอบ volume confirmation (ผ่อนคลาย)
	volumeConfirm := bt.hasVolumeConfirmation()

	// ต้องผ่านอย่างน้อย 4 จาก 5 เงื่อนไข
	conditions := []bool{trendAligned, confidenceOK, priceAligned, goodRiskReward, volumeConfirm}
	passedCount := 0
	for _, condition := range conditions {
		if condition {
			passedCount++
		}
	}

	return passedCount >= 4
}

// determineDirection กำหนดทิศทางการเทรด
func (bt *Backtester) determineDirection(analysis *SuperTrendAnalysis) string {
	// ตรวจสอบ SuperTrend + EMA alignment
	if analysis.Trend == 1 &&
		analysis.CurrentPrice > analysis.SuperTrendValue &&
		analysis.CurrentPrice > analysis.EMA100 &&
		analysis.Signal == "LONG" {
		return "LONG"
	}

	if analysis.Trend == -1 &&
		analysis.CurrentPrice < analysis.SuperTrendValue &&
		analysis.CurrentPrice < analysis.EMA100 &&
		analysis.Signal == "SHORT" {
		return "SHORT"
	}

	return ""
}

// calculateRiskReward คำนวณ Stop Loss และ Take Profit - ปรับให้เหมาะสม
func (bt *Backtester) calculateRiskReward(direction string, currentPrice, atr float64) (stopLoss, takeProfit float64) {
	// ใช้ ATR สำหรับการคำนวณ SL/TP (ปรับให้เหมาะสมกับ 1H timeframe)
	atrMultiplier := 1.5   // ลดจาก 2.0 เป็น 1.5
	riskRewardRatio := 2.5 // ลดจาก 3.0 เป็น 2.5

	if direction == "LONG" {
		stopLoss = currentPrice - (atr * atrMultiplier)
		takeProfit = currentPrice + (atr * atrMultiplier * riskRewardRatio)
	} else if direction == "SHORT" {
		stopLoss = currentPrice + (atr * atrMultiplier)
		takeProfit = currentPrice - (atr * atrMultiplier * riskRewardRatio)
	}

	return stopLoss, takeProfit
}

// getSignalReason สร้างเหตุผลของสัญญาณ
func (bt *Backtester) getSignalReason(analysis *SuperTrendAnalysis) string {
	return fmt.Sprintf("Strong %s signal: Trend=%d, Confidence=%.1f%%, Price=%.2f, SuperTrend=%.2f, EMA100=%.2f, RR=%.2f",
		analysis.Signal, analysis.Trend, analysis.Confidence,
		analysis.CurrentPrice, analysis.SuperTrendValue, analysis.EMA100, analysis.RiskRewardRatio)
}

// hasVolumeConfirmation ตรวจสอบการยืนยันจาก volume - ปรับให้ผ่อนคลาย
func (bt *Backtester) hasVolumeConfirmation() bool {
	if bt.currentIndex < 3 { // ลดจาก 5 เป็น 3
		return true // ให้ผ่านถ้าข้อมูลไม่เพียงพอ
	}

	// เปรียบเทียบ volume ล่าสุดกับค่าเฉลี่ย 3 แท่งก่อนหน้า
	currentVolume := bt.ohlcvData[bt.currentIndex].Volume
	avgVolume := 0.0

	for i := bt.currentIndex - 3; i < bt.currentIndex; i++ {
		avgVolume += bt.ohlcvData[i].Volume
	}
	avgVolume /= 3.0

	// volume ต้องสูงกว่าค่าเฉลี่ย 10% (ลดจาก 20%)
	return currentVolume > avgVolume*1.1
}

// isPriceNearSupport ตรวจสอบราคาใกล้ support/resistance
func (bt *Backtester) isPriceNearSupport(direction string) bool {
	if bt.currentIndex < 50 {
		return false
	}

	// หา support/resistance level ใน 50 แท่งล่าสุด
	var levels []float64
	for i := bt.currentIndex - 50; i < bt.currentIndex; i++ {
		levels = append(levels, bt.ohlcvData[i].High, bt.ohlcvData[i].Low)
	}

	// ตรวจสอบระยะห่างจาก key levels
	tolerance := bt.currentPrice * 0.01 // 1% tolerance

	for _, level := range levels {
		if math.Abs(bt.currentPrice-level) < tolerance {
			if direction == "LONG" && level < bt.currentPrice {
				return true // ราคาอยู่เหนือ support
			}
			if direction == "SHORT" && level > bt.currentPrice {
				return true // ราคาอยู่ใต้ resistance
			}
		}
	}

	return false
}

// updateTrailingStop อัปเดต trailing stop loss
func (bt *Backtester) updateTrailingStop() {
	if bt.position == nil {
		return
	}

	// คำนวณ trailing stop ที่ 1.5x ATR
	analysis := bt.analyzeMarket()
	if analysis == nil {
		return
	}

	trailingDistance := analysis.ATR * 1.5

	if bt.position.Side == "LONG" {
		// สำหรับ LONG position
		newStopLoss := bt.currentPrice - trailingDistance
		if newStopLoss > bt.position.StopLoss {
			bt.position.StopLoss = newStopLoss
		}
	} else if bt.position.Side == "SHORT" {
		// สำหรับ SHORT position
		newStopLoss := bt.currentPrice + trailingDistance
		if newStopLoss < bt.position.StopLoss {
			bt.position.StopLoss = newStopLoss
		}
	}
}

// checkEarlyExit ตรวจสอบสัญญาณออกก่อนกำหนด
func (bt *Backtester) checkEarlyExit() {
	if bt.position == nil {
		return
	}

	analysis := bt.analyzeMarket()
	if analysis == nil {
		return
	}

	// ตรวจสอบ momentum reversal
	if bt.isMomentumReversal(analysis) {
		bt.closePosition("MOMENTUM_REVERSAL")
		return
	}

	// ตรวจสอบ overbought/oversold extreme
	if bt.isExtremeCondition(analysis) {
		bt.closePosition("EXTREME_CONDITION")
		return
	}

	// ตรวจสอบ time-based exit (ถือครองเกิน 8 ชั่วโมง)
	holdingTime := bt.currentTime.Sub(bt.position.EntryTime)
	if holdingTime > 8*time.Hour {
		bt.closePosition("TIME_EXIT")
		return
	}
}

// isMomentumReversal ตรวจสอบการกลับทิศของ momentum
func (bt *Backtester) isMomentumReversal(analysis *SuperTrendAnalysis) bool {
	// ตรวจสอบว่า trend เปลี่ยนทิศหรือไม่
	if bt.position.Side == "LONG" && analysis.Trend == -1 {
		return true
	}
	if bt.position.Side == "SHORT" && analysis.Trend == 1 {
		return true
	}

	// ตรวจสอบการ break ของ SuperTrend
	if bt.position.Side == "LONG" && analysis.CurrentPrice < analysis.SuperTrendValue {
		return true
	}
	if bt.position.Side == "SHORT" && analysis.CurrentPrice > analysis.SuperTrendValue {
		return true
	}

	return false
}

// isExtremeCondition ตรวจสอบสภาวะ extreme
func (bt *Backtester) isExtremeCondition(analysis *SuperTrendAnalysis) bool {
	// จำลองการตรวจสอบ RSI extreme
	// ในการใช้งานจริงควรมีการคำนวณ RSI

	// ตรวจสอบระยะห่างจาก EMA100 เกินไป (มากกว่า 5%)
	emaDistance := (analysis.CurrentPrice - analysis.EMA100) / analysis.EMA100 * 100

	if bt.position.Side == "LONG" && emaDistance > 5.0 {
		return true // ราคาสูงเกินไปจาก EMA
	}
	if bt.position.Side == "SHORT" && emaDistance < -5.0 {
		return true // ราคาต่ำเกินไปจาก EMA
	}

	return false
}

// checkSimpleSignals ตรวจสอบสัญญาณง่าย ๆ เป็น fallback
func (bt *Backtester) checkSimpleSignals(analysis *SuperTrendAnalysis) {
	// สัญญาณง่าย: SuperTrend + EMA crossover
	if analysis.Trend == 1 && analysis.CurrentPrice > analysis.SuperTrendValue &&
		analysis.CurrentPrice > analysis.EMA100 {
		stopLoss := analysis.SuperTrendValue
		takeProfit := bt.currentPrice + (bt.currentPrice-stopLoss)*2.0
		fmt.Printf("📈 Simple LONG signal: Price %.2f > SuperTrend %.2f > EMA100 %.2f\n",
			analysis.CurrentPrice, analysis.SuperTrendValue, analysis.EMA100)
		bt.openPosition("LONG", stopLoss, takeProfit, "Simple SuperTrend + EMA100 LONG")
	} else if analysis.Trend == -1 && analysis.CurrentPrice < analysis.SuperTrendValue &&
		analysis.CurrentPrice < analysis.EMA100 {
		stopLoss := analysis.SuperTrendValue
		takeProfit := bt.currentPrice - (stopLoss-bt.currentPrice)*2.0
		fmt.Printf("📉 Simple SHORT signal: Price %.2f < SuperTrend %.2f < EMA100 %.2f\n",
			analysis.CurrentPrice, analysis.SuperTrendValue, analysis.EMA100)
		bt.openPosition("SHORT", stopLoss, takeProfit, "Simple SuperTrend + EMA100 SHORT")
	}
}

// RunTripleEMA1HStrategy - กลยุทธ์ Triple EMA สำหรับ 1H timeframe (144 candles)
func (bt *Backtester) RunTripleEMA1HStrategy() *BacktestResult {
	fmt.Printf("\n🚀 เริ่มต้น Triple EMA 1H Strategy %s\n", bt.symbol)
	fmt.Printf("💰 เงินทุนเริ่มต้น: $%.2f\n", bt.initialCapital)
	fmt.Printf("⏰ Time Frame: 1 ชั่วโมง (144 candles)\n")
	fmt.Printf("🎯 กลยุทธ์: Triple EMA(9,21,50) + RSI + Volume + ATR\n")
	fmt.Printf("📊 รองรับ: LONG & SHORT Positions (5x Leverage)\n")

	// ตรวจสอบข้อมูลครบ 144 candles
	if len(bt.ohlcvData) < 144 {
		fmt.Printf("❌ ข้อมูลไม่ครบ 144 candles (มี %d candles)\n", len(bt.ohlcvData))
		return bt.generateTripleEMA1HResults()
	}

	// เริ่มจากแท่งที่ 50 เพื่อให้ indicators พร้อม
	for bt.currentIndex = 50; bt.currentIndex < len(bt.ohlcvData); bt.currentIndex++ {
		currentCandle := bt.ohlcvData[bt.currentIndex]
		bt.currentTime = time.Unix(currentCandle.Timestamp, 0)
		bt.currentPrice = currentCandle.Close

		// วิเคราะห์ตลาดด้วย Triple EMA 1H
		analysis := bt.analyzeTripleEMA1H()
		if analysis == nil {
			continue
		}

		// จัดการ position ที่มีอยู่
		bt.checkTripleEMA1HPosition()

		// หาจุดเข้าใหม่
		bt.lookForTripleEMA1HEntry(analysis)

		// Debug: แสดงข้อมูลทุก 24 candles (1 วัน)
		if bt.currentIndex%24 == 0 {
			fmt.Printf("🔍 Debug [%d]: Price=%.2f, Signal=%s, Conf=%.1f%%, 1H\n",
				bt.currentIndex, bt.currentPrice, analysis.Signal, analysis.Confidence)
		}
	}

	// ปิด position สุดท้าย
	if bt.position != nil {
		bt.closeTripleEMA1HPosition("End of Backtest")
	}

	return bt.generateTripleEMA1HResults()
}

// analyzeTripleEMA1H วิเคราะห์ตลาดด้วย Triple EMA สำหรับ 1H
func (bt *Backtester) analyzeTripleEMA1H() *SuperTrendAnalysis {
	if bt.currentIndex < 50 {
		return nil
	}

	// เอาข้อมูล 50 แท่งล่าสุดสำหรับการคำนวณ indicators
	data := bt.ohlcvData[bt.currentIndex-49 : bt.currentIndex+1]

	// คำนวณ Triple EMA สำหรับ 1H: Fast=9, Mid=21, Slow=50
	emaFast := bt.calculateEMA(data, 9)  // Fast EMA
	emaMid := bt.calculateEMA(data, 21)  // Mid EMA
	emaSlow := bt.calculateEMA(data, 50) // Slow EMA (ใช้ทุกข้อมูล 50 แท่ง)

	// คำนวณ indicators อื่นๆ
	rsi := bt.calculateRSI(data, 14)
	atr := bt.calculateATR(data, 14)
	avgVolume := bt.calculateAvgVolume(data, 20)

	currentCandle := data[len(data)-1]

	// กำหนดสัญญาณ
	signal := "NEUTRAL"
	confidence := 60.0

	// Bullish Signal: Fast > Mid > Slow (Perfect alignment)
	if emaFast > emaMid && emaMid > emaSlow &&
		bt.currentPrice >= emaFast*0.998 && // Price near/above fast EMA
		rsi > 40 && rsi < 80 { // RSI in healthy range

		signal = "LONG"
		confidence = 70.0

		// เพิ่ม confidence สำหรับ momentum แข็งแกร่ง
		if currentCandle.Close > currentCandle.Open { // Green candle
			confidence += 10.0
		}

		// เพิ่ม confidence ถ้า EMA spread กว้าง
		emaSpread := ((emaFast - emaSlow) / emaSlow) * 100
		if emaSpread > 0.15 { // EMA spread > 0.15%
			confidence += 10.0
		}

		// Volume confirmation
		if currentCandle.Volume > avgVolume*1.2 {
			confidence += 5.0
		}

		// RSI optimal range
		if rsi >= 50 && rsi <= 70 {
			confidence += 10.0
		}
	}

	// Bearish Signal: Fast < Mid < Slow (Perfect alignment)
	if emaFast < emaMid && emaMid < emaSlow &&
		bt.currentPrice <= emaFast*1.002 && // Price near/below fast EMA
		rsi > 20 && rsi < 60 { // RSI in healthy range

		signal = "SHORT"
		confidence = 70.0

		// เพิ่ม confidence สำหรับ momentum แข็งแกร่ง
		if currentCandle.Close < currentCandle.Open { // Red candle
			confidence += 10.0
		}

		// เพิ่ม confidence ถ้า EMA spread กว้าง
		emaSpread := ((emaSlow - emaFast) / emaSlow) * 100
		if emaSpread > 0.15 { // EMA spread > 0.15%
			confidence += 10.0
		}

		// Volume confirmation
		if currentCandle.Volume > avgVolume*1.2 {
			confidence += 5.0
		}

		// RSI optimal range
		if rsi >= 30 && rsi <= 50 {
			confidence += 10.0
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

// lookForTripleEMA1HEntry หาจุดเข้าด้วย Triple EMA 1H
func (bt *Backtester) lookForTripleEMA1HEntry(analysis *SuperTrendAnalysis) {
	if bt.position != nil {
		return
	}

	// เข้า LONG position (ใช้ confidence threshold 75%)
	if analysis.Signal == "LONG" && analysis.Confidence >= 75.0 {
		bt.openTripleEMA1HPosition("LONG", analysis)
	}

	// เข้า SHORT position
	if analysis.Signal == "SHORT" && analysis.Confidence >= 75.0 {
		bt.openTripleEMA1HPosition("SHORT", analysis)
	}
}

// openTripleEMA1HPosition เปิด position ด้วย Triple EMA 1H
func (bt *Backtester) openTripleEMA1HPosition(side string, analysis *SuperTrendAnalysis) {
	// Future Trading Parameters
	leverage := 5.0

	var stopLoss, takeProfit float64

	// คำนวณ Stop Loss และ Take Profit โดยใช้ ATR
	atrMultiplier := 2.0 // ใช้ 2.0x ATR สำหรับ 1H timeframe

	if side == "LONG" {
		stopLoss = bt.currentPrice - (analysis.ATR * atrMultiplier)
		takeProfit = bt.currentPrice + (analysis.ATR * atrMultiplier * analysis.RiskRewardRatio)
	} else {
		stopLoss = bt.currentPrice + (analysis.ATR * atrMultiplier)
		takeProfit = bt.currentPrice - (analysis.ATR * atrMultiplier * analysis.RiskRewardRatio)
	}

	// Risk per trade 1.0%
	riskAmount := bt.currentCapital * 0.01

	var positionSize float64
	if side == "LONG" {
		positionSize = (riskAmount * leverage) / (bt.currentPrice - stopLoss)
	} else {
		positionSize = (riskAmount * leverage) / (stopLoss - bt.currentPrice)
	}

	// ป้องกันใช้เงินเกิน 15%
	maxValue := bt.currentCapital * 0.15 * leverage
	if positionSize*bt.currentPrice > maxValue {
		positionSize = maxValue / bt.currentPrice
	}

	// คำนวณค่าธรรมเนียม Future
	entryCommission := bt.currentPrice * positionSize * 0.0005

	// สร้าง position
	bt.position = &BacktestPosition{
		Symbol:      bt.symbol,
		Side:        side,
		EntryTime:   bt.currentTime,
		EntryPrice:  bt.currentPrice,
		Quantity:    positionSize,
		StopLoss:    stopLoss,
		TakeProfit:  takeProfit,
		EntryReason: fmt.Sprintf("1H %s: Triple EMA, Conf=%.1f%%, Lev=%.1fx", side, analysis.Confidence, leverage),
	}

	// หักค่าธรรมเนียม
	bt.currentCapital -= entryCommission

	// บันทึก trade
	bt.recordTripleEMA1HTrade("OPEN", fmt.Sprintf("1H Entry: %s (Triple EMA)", side))

	var slPct, tpPct float64
	if side == "LONG" {
		slPct = ((stopLoss - bt.currentPrice) / bt.currentPrice) * 100
		tpPct = ((takeProfit - bt.currentPrice) / bt.currentPrice) * 100
	} else {
		slPct = ((bt.currentPrice - stopLoss) / bt.currentPrice) * 100
		tpPct = ((bt.currentPrice - takeProfit) / bt.currentPrice) * 100
	}

	fmt.Printf("🎯 เปิด 1H Triple EMA %s: @ $%.2f, SL: $%.2f (%.2f%%), TP: $%.2f (%.2f%%), Size: %.6f, Fee: $%.2f\n",
		side, bt.currentPrice, stopLoss, slPct, takeProfit, tpPct, positionSize, entryCommission)
}

// checkTripleEMA1HPosition ตรวจสอบ position สำหรับ 1H
func (bt *Backtester) checkTripleEMA1HPosition() {
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
	}

	// Signal reversal exit
	analysis := bt.analyzeTripleEMA1H()
	if analysis != nil {
		if bt.position.Side == "LONG" && analysis.Signal == "SHORT" && analysis.Confidence >= 80.0 {
			shouldExit = true
			exitReason = "Signal Reversal"
		}
		if bt.position.Side == "SHORT" && analysis.Signal == "LONG" && analysis.Confidence >= 80.0 {
			shouldExit = true
			exitReason = "Signal Reversal"
		}
	}

	// Time exit - 24 hours (1H × 24 = 1 วัน)
	if bt.currentTime.Sub(bt.position.EntryTime).Hours() > 24 {
		shouldExit = true
		exitReason = "Time Exit (24h)"
	}

	if shouldExit {
		bt.closeTripleEMA1HPosition(exitReason)
	}
}

// closeTripleEMA1HPosition ปิด position สำหรับ 1H
func (bt *Backtester) closeTripleEMA1HPosition(reason string) {
	if bt.position == nil {
		return
	}

	var pnl float64

	if bt.position.Side == "LONG" {
		pnl = (bt.currentPrice - bt.position.EntryPrice) * bt.position.Quantity
	} else {
		pnl = (bt.position.EntryPrice - bt.currentPrice) * bt.position.Quantity
	}

	// คำนวณค่าธรรมเนียม
	entryCommission := bt.position.EntryPrice * bt.position.Quantity * 0.0005
	exitCommission := bt.currentPrice * bt.position.Quantity * 0.0005
	totalCommission := entryCommission + exitCommission

	netPnL := pnl - totalCommission

	// อัปเดตเงินทุน
	bt.currentCapital += netPnL

	// บันทึก trade
	bt.recordTripleEMA1HTrade("CLOSE", reason)

	// แสดงผล
	pnlPct := (netPnL / (bt.position.EntryPrice * bt.position.Quantity)) * 100
	fmt.Printf("🔹 ปิด 1H Triple EMA %s: PnL $%.2f (%.2f%%) - Commission: $%.2f - %s\n",
		bt.position.Side, netPnL, pnlPct, totalCommission, reason)

	// ลบ position
	bt.position = nil
}

// recordTripleEMA1HTrade บันทึก trade สำหรับ 1H
func (bt *Backtester) recordTripleEMA1HTrade(action, reason string) {
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

// generateTripleEMA1HResults สร้างผลลัพธ์สำหรับ 1H
func (bt *Backtester) generateTripleEMA1HResults() *BacktestResult {
	totalReturn := ((bt.currentCapital - bt.initialCapital) / bt.initialCapital) * 100

	// นับ trades
	var winTrades, lossTrades int
	var totalPnL float64

	for _, trade := range bt.trades {
		if trade.NetPnL != 0 {
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
		MaxDrawdown:    0.0,
		MaxDrawdownPct: 0.0,
		Trades:         bt.trades,
		DailyReturns:   bt.dailyReturns,
	}
}
