package main

import (
	"fmt"
	"strings"
	"time"

	"gateio-trading-bot/internal/gateio"
)

func main() {
	fmt.Println("🚀 Enhanced Triple EMA 1H Strategy with 144 Timeframes")
	fmt.Println(strings.Repeat("=", 70))

	// สร้าง client
	client := gateio.NewClient("", "", "https://api.gateio.ws")

	// ทดสอบ Enhanced Strategy
	symbols := []string{"SOL_USDT", "BTC_USDT", "ETH_USDT"}

	for i, symbol := range symbols {
		fmt.Printf("\n🔍 [%d/%d] Enhanced Triple EMA Strategy: %s\n", i+1, len(symbols), symbol)
		runEnhancedStrategy(symbol, client)

		if i < len(symbols)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Println("\n✅ Enhanced Strategy Testing Complete!")
}

func runEnhancedStrategy(symbol string, client *gateio.Client) {
	// สร้าง Enhanced Data Manager
	dm := NewEnhanced1HDataManager(symbol, client)

	// โหลดข้อมูล 144 timeframes
	if err := dm.LoadHistoricalData(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	// สร้าง Enhanced Backtester
	backtester := NewEnhancedBacktester(symbol, dm, 10000.0)

	// รัน Enhanced Strategy
	result := backtester.RunEnhancedTripleEMA1H()

	// แสดงผลลัพธ์
	displayEnhancedResults(result)
}

// EnhancedBacktester - Backtester ที่ใช้ Enhanced Data Manager
type EnhancedBacktester struct {
	symbol         string
	dataManager    *Enhanced1HDataManager
	initialCapital float64
	currentCapital float64
	position       *EnhancedPosition
	trades         []EnhancedTrade
	tradeID        int
}

type EnhancedPosition struct {
	ID          int       `json:"id"`
	Symbol      string    `json:"symbol"`
	Side        string    `json:"side"`
	EntryTime   time.Time `json:"entry_time"`
	EntryPrice  float64   `json:"entry_price"`
	Quantity    float64   `json:"quantity"`
	StopLoss    float64   `json:"stop_loss"`
	TakeProfit  float64   `json:"take_profit"`
	Leverage    float64   `json:"leverage"`
	Confidence  float64   `json:"confidence"`
	EntryReason string    `json:"entry_reason"`
}

type EnhancedTrade struct {
	ID          int           `json:"id"`
	Symbol      string        `json:"symbol"`
	Side        string        `json:"side"`
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
	Confidence  float64       `json:"confidence"`
	DataQuality float64       `json:"data_quality"`
}

type EnhancedResult struct {
	Symbol         string          `json:"symbol"`
	Strategy       string          `json:"strategy"`
	TimeFrame      string          `json:"timeframe"`
	LookbackHours  int             `json:"lookback_hours"`
	InitialCapital float64         `json:"initial_capital"`
	FinalCapital   float64         `json:"final_capital"`
	TotalReturn    float64         `json:"total_return"`
	TotalReturnPct float64         `json:"total_return_pct"`
	TotalTrades    int             `json:"total_trades"`
	WinningTrades  int             `json:"winning_trades"`
	LosingTrades   int             `json:"losing_trades"`
	WinRate        float64         `json:"win_rate"`
	AvgTradeHours  float64         `json:"avg_trade_hours"`
	MaxDrawdown    float64         `json:"max_drawdown"`
	AvgConfidence  float64         `json:"avg_confidence"`
	AvgDataQuality float64         `json:"avg_data_quality"`
	Trades         []EnhancedTrade `json:"trades"`
	CacheHitRate   float64         `json:"cache_hit_rate"`
}

func NewEnhancedBacktester(symbol string, dm *Enhanced1HDataManager, capital float64) *EnhancedBacktester {
	return &EnhancedBacktester{
		symbol:         symbol,
		dataManager:    dm,
		initialCapital: capital,
		currentCapital: capital,
		trades:         make([]EnhancedTrade, 0),
	}
}

func (eb *EnhancedBacktester) RunEnhancedTripleEMA1H() *EnhancedResult {
	fmt.Printf("🚀 Enhanced Triple EMA 1H Strategy: %s\n", eb.symbol)
	fmt.Printf("💰 Initial Capital: $%.2f\n", eb.initialCapital)
	fmt.Printf("⏰ Lookback: 144 timeframes (6 days)\n")

	// ใช้ข้อมูลจาก data manager
	dataBuffer := eb.dataManager.dataBuffer

	// จำลองการเทรดจากแท่งที่ 50 (ให้ indicators พร้อม)
	for i := 50; i < len(dataBuffer); i++ {
		current := dataBuffer[i]
		currentTime := time.Unix(current.Timestamp, 0)
		currentPrice := current.Close

		// อัปเดต data manager index (จำลอง real-time)
		eb.dataManager.currentIndex = i

		// วิเคราะห์สัญญาณ
		analysis := eb.dataManager.GetLatestAnalysis()
		if analysis == nil {
			continue
		}

		// จัดการ position ที่มีอยู่
		eb.checkPosition(currentPrice, currentTime, analysis)

		// หาจุดเข้าใหม่
		eb.lookForEntry(currentPrice, currentTime, analysis)
	}

	// ปิด position สุดท้าย
	if eb.position != nil {
		eb.closePosition("End of Backtest", time.Unix(dataBuffer[len(dataBuffer)-1].Timestamp, 0))
	}

	return eb.generateResult()
}

func (eb *EnhancedBacktester) lookForEntry(price float64, currentTime time.Time, analysis *Enhanced1HAnalysis) {
	if eb.position != nil {
		return // มี position อยู่แล้ว
	}

	// เข้า LONG
	if analysis.Signal == "LONG" && analysis.Confidence >= 80.0 && analysis.DataQuality >= 90.0 {
		eb.openPosition("LONG", price, currentTime, analysis)
	}

	// เข้า SHORT
	if analysis.Signal == "SHORT" && analysis.Confidence >= 80.0 && analysis.DataQuality >= 90.0 {
		eb.openPosition("SHORT", price, currentTime, analysis)
	}
}

func (eb *EnhancedBacktester) openPosition(side string, price float64, currentTime time.Time, analysis *Enhanced1HAnalysis) {
	leverage := 5.0
	atrMultiplier := 2.0

	var stopLoss, takeProfit float64
	if side == "LONG" {
		stopLoss = price - (analysis.ATR * atrMultiplier)
		takeProfit = price + (analysis.ATR * atrMultiplier * 2.0) // 1:2 R/R
	} else {
		stopLoss = price + (analysis.ATR * atrMultiplier)
		takeProfit = price - (analysis.ATR * atrMultiplier * 2.0)
	}

	// Risk management: 1% per trade
	riskAmount := eb.currentCapital * 0.01
	var positionSize float64
	if side == "LONG" {
		positionSize = (riskAmount * leverage) / (price - stopLoss)
	} else {
		positionSize = (riskAmount * leverage) / (stopLoss - price)
	}

	// ไม่ให้เกิน 15% ของ capital
	maxValue := eb.currentCapital * 0.15 * leverage
	if positionSize*price > maxValue {
		positionSize = maxValue / price
	}

	commission := positionSize * price * 0.0005
	eb.currentCapital -= commission

	eb.tradeID++
	eb.position = &EnhancedPosition{
		ID:          eb.tradeID,
		Symbol:      eb.symbol,
		Side:        side,
		EntryTime:   currentTime,
		EntryPrice:  price,
		Quantity:    positionSize,
		StopLoss:    stopLoss,
		TakeProfit:  takeProfit,
		Leverage:    leverage,
		Confidence:  analysis.Confidence,
		EntryReason: fmt.Sprintf("Enhanced %s (Conf: %.1f%%, Quality: %.1f%%)", side, analysis.Confidence, analysis.DataQuality),
	}

	fmt.Printf("🎯 OPEN %s: $%.2f, SL: $%.2f, TP: $%.2f, Conf: %.1f%%, Quality: %.1f%%\n",
		side, price, stopLoss, takeProfit, analysis.Confidence, analysis.DataQuality)
}

func (eb *EnhancedBacktester) checkPosition(price float64, currentTime time.Time, analysis *Enhanced1HAnalysis) {
	if eb.position == nil {
		return
	}

	shouldExit := false
	exitReason := ""

	// Stop Loss / Take Profit
	if eb.position.Side == "LONG" {
		if price <= eb.position.StopLoss {
			shouldExit = true
			exitReason = "Stop Loss"
		} else if price >= eb.position.TakeProfit {
			shouldExit = true
			exitReason = "Take Profit"
		}
	} else {
		if price >= eb.position.StopLoss {
			shouldExit = true
			exitReason = "Stop Loss"
		} else if price <= eb.position.TakeProfit {
			shouldExit = true
			exitReason = "Take Profit"
		}
	}

	// Signal reversal (higher confidence required)
	if analysis.Confidence >= 85.0 {
		if eb.position.Side == "LONG" && analysis.Signal == "SHORT" {
			shouldExit = true
			exitReason = "Signal Reversal"
		} else if eb.position.Side == "SHORT" && analysis.Signal == "LONG" {
			shouldExit = true
			exitReason = "Signal Reversal"
		}
	}

	// Time limit: 24 hours
	if currentTime.Sub(eb.position.EntryTime).Hours() >= 24 {
		shouldExit = true
		exitReason = "Time Limit (24h)"
	}

	// Data quality deterioration
	if analysis.DataQuality < 80.0 {
		shouldExit = true
		exitReason = "Data Quality Drop"
	}

	if shouldExit {
		eb.closePosition(exitReason, currentTime)
	}
}

func (eb *EnhancedBacktester) closePosition(reason string, currentTime time.Time) {
	if eb.position == nil {
		return
	}

	currentPrice := eb.dataManager.GetLatestAnalysis().Price

	var pnl float64
	if eb.position.Side == "LONG" {
		pnl = (currentPrice - eb.position.EntryPrice) * eb.position.Quantity
	} else {
		pnl = (eb.position.EntryPrice - currentPrice) * eb.position.Quantity
	}

	entryCommission := eb.position.EntryPrice * eb.position.Quantity * 0.0005
	exitCommission := currentPrice * eb.position.Quantity * 0.0005
	totalCommission := entryCommission + exitCommission
	netPnL := pnl - totalCommission

	eb.currentCapital += netPnL

	trade := EnhancedTrade{
		ID:          eb.position.ID,
		Symbol:      eb.position.Symbol,
		Side:        eb.position.Side,
		EntryTime:   eb.position.EntryTime,
		ExitTime:    currentTime,
		EntryPrice:  eb.position.EntryPrice,
		ExitPrice:   currentPrice,
		Quantity:    eb.position.Quantity,
		PnL:         pnl,
		PnLPct:      (pnl / (eb.position.EntryPrice * eb.position.Quantity)) * 100,
		Commission:  totalCommission,
		NetPnL:      netPnL,
		Duration:    currentTime.Sub(eb.position.EntryTime),
		EntryReason: eb.position.EntryReason,
		ExitReason:  reason,
		Confidence:  eb.position.Confidence,
		DataQuality: eb.dataManager.GetLatestAnalysis().DataQuality,
	}

	eb.trades = append(eb.trades, trade)

	pnlPct := (netPnL / (eb.position.EntryPrice * eb.position.Quantity)) * 100
	fmt.Printf("🔹 CLOSE %s: PnL $%.2f (%.2f%%) - %s\n",
		eb.position.Side, netPnL, pnlPct, reason)

	eb.position = nil
}

func (eb *EnhancedBacktester) generateResult() *EnhancedResult {
	totalReturn := eb.currentCapital - eb.initialCapital
	totalReturnPct := (totalReturn / eb.initialCapital) * 100

	var winTrades, lossTrades int
	var totalConfidence, totalQuality float64
	var totalDuration time.Duration

	for _, trade := range eb.trades {
		if trade.NetPnL > 0 {
			winTrades++
		} else {
			lossTrades++
		}
		totalConfidence += trade.Confidence
		totalQuality += trade.DataQuality
		totalDuration += trade.Duration
	}

	totalTrades := len(eb.trades)
	var winRate, avgConfidence, avgQuality, avgTradeHours float64

	if totalTrades > 0 {
		winRate = (float64(winTrades) / float64(totalTrades)) * 100
		avgConfidence = totalConfidence / float64(totalTrades)
		avgQuality = totalQuality / float64(totalTrades)
		avgTradeHours = totalDuration.Hours() / float64(totalTrades)
	}

	_, _, cacheHitRate := eb.dataManager.GetCacheStats()

	return &EnhancedResult{
		Symbol:         eb.symbol,
		Strategy:       "Enhanced Triple EMA 1H",
		TimeFrame:      "1H",
		LookbackHours:  144,
		InitialCapital: eb.initialCapital,
		FinalCapital:   eb.currentCapital,
		TotalReturn:    totalReturn,
		TotalReturnPct: totalReturnPct,
		TotalTrades:    totalTrades,
		WinningTrades:  winTrades,
		LosingTrades:   lossTrades,
		WinRate:        winRate,
		AvgTradeHours:  avgTradeHours,
		AvgConfidence:  avgConfidence,
		AvgDataQuality: avgQuality,
		Trades:         eb.trades,
		CacheHitRate:   cacheHitRate,
	}
}

func displayEnhancedResults(result *EnhancedResult) {
	fmt.Printf("\n📊 Enhanced Results: %s\n", result.Symbol)
	fmt.Printf("💰 Return: $%.2f → $%.2f (%+.2f%%)\n",
		result.InitialCapital, result.FinalCapital, result.TotalReturnPct)
	fmt.Printf("📈 Trades: %d (Win: %d, Loss: %d, Rate: %.1f%%)\n",
		result.TotalTrades, result.WinningTrades, result.LosingTrades, result.WinRate)
	fmt.Printf("⏱️  Avg Trade: %.1f hours\n", result.AvgTradeHours)
	fmt.Printf("🎯 Avg Confidence: %.1f%%\n", result.AvgConfidence)
	fmt.Printf("✅ Avg Data Quality: %.1f%%\n", result.AvgDataQuality)
	fmt.Printf("🚀 Cache Hit Rate: %.1f%%\n", result.CacheHitRate)

	// แสดง trades ล่าสุด
	if len(result.Trades) > 0 {
		fmt.Printf("\n📋 Recent Trades:\n")
		start := len(result.Trades) - 3
		if start < 0 {
			start = 0
		}

		for i := start; i < len(result.Trades); i++ {
			trade := result.Trades[i]
			fmt.Printf("   %d. %s: %+.2f%% (%.1fh, Conf: %.1f%%) - %s\n",
				trade.ID, trade.Side, trade.PnLPct, trade.Duration.Hours(),
				trade.Confidence, trade.ExitReason)
		}
	}

	// Performance assessment
	fmt.Printf("\n🎯 Performance Assessment:\n")
	if result.TotalReturnPct > 20 {
		fmt.Printf("   🚀 EXCELLENT Performance!\n")
	} else if result.TotalReturnPct > 10 {
		fmt.Printf("   📈 GOOD Performance\n")
	} else if result.TotalReturnPct > 0 {
		fmt.Printf("   📊 POSITIVE Performance\n")
	} else {
		fmt.Printf("   📉 NEEDS IMPROVEMENT\n")
	}

	if result.WinRate > 60 {
		fmt.Printf("   ✅ Strong Win Rate (%.1f%%)\n", result.WinRate)
	} else if result.WinRate > 40 {
		fmt.Printf("   📊 Decent Win Rate (%.1f%%)\n", result.WinRate)
	} else {
		fmt.Printf("   ⚠️  Low Win Rate (%.1f%%)\n", result.WinRate)
	}
}
