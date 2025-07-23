package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// Pivot Point SuperTrend + EMA100 Strategy
// Updated: Previous candle must be above/below EMA100 AND correct color

type PivotSupertrendEMAStrategy struct {
	Symbol      string
	Data        []CandleData
	Positions   []Position
	StartCash   float64
	CurrentCash float64
}

type Position struct {
	Symbol     string
	Side       string // "LONG" or "SHORT"
	EntryPrice float64
	Size       float64
	EntryTime  time.Time
	ExitPrice  float64
	ExitTime   time.Time
	PnL        float64
	Reason     string
	Duration   time.Duration
}

type CandleData struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	// Indicators
	EMA100 float64 // EMA 100 filter
	RSI    float64 // RSI for exits
	ATR    float64 // ATR for SuperTrend
	// Pivot Point SuperTrend
	PivotHigh  float64
	PivotLow   float64
	Center     float64
	UpperBand  float64
	LowerBand  float64
	SuperTrend float64
	Trend      int // 1 = Buy Trend, -1 = Sell Trend
	// Signal flags
	BuySignal       bool
	SellSignal      bool
	IsGreenCandle   bool // Current candle is green
	IsRedCandle     bool // Current candle is red
	PrevGreenCandle bool // Previous candle is green
	PrevRedCandle   bool // Previous candle is red
	PrevAboveEMA100 bool // Previous candle close above EMA100
	PrevBelowEMA100 bool // Previous candle close below EMA100
}

func NewPivotSupertrendEMAStrategy(symbol string, startCash float64) *PivotSupertrendEMAStrategy {
	return &PivotSupertrendEMAStrategy{
		Symbol:      symbol,
		StartCash:   startCash,
		CurrentCash: startCash,
		Positions:   make([]Position, 0),
	}
}

func (s *PivotSupertrendEMAStrategy) LoadRealMarketData() error {
	fmt.Printf("🔄 Loading REAL market data for %s from Gate.io...\n", s.Symbol)

	// Current time and 6 days ago (144 hours)
	endTime := time.Now().Unix()
	startTime := endTime - (144 * 3600) // 144 hours ago

	url := fmt.Sprintf("https://api.gateio.ws/api/v4/spot/candlesticks?currency_pair=%s&interval=1h&from=%d&to=%d&limit=200",
		s.Symbol, startTime, endTime)

	fmt.Printf("📡 Fetching from: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var rawData [][]string
	if err := json.Unmarshal(body, &rawData); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	s.Data = make([]CandleData, 0, len(rawData))

	for _, item := range rawData {
		if len(item) < 6 {
			continue
		}

		timestamp, _ := strconv.ParseInt(item[0], 10, 64)
		volume, _ := strconv.ParseFloat(item[1], 64)
		close, _ := strconv.ParseFloat(item[2], 64)
		high, _ := strconv.ParseFloat(item[3], 64)
		low, _ := strconv.ParseFloat(item[4], 64)
		open, _ := strconv.ParseFloat(item[5], 64)

		candle := CandleData{
			Timestamp: time.Unix(timestamp, 0),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}

		s.Data = append(s.Data, candle)
	}

	// Sort by timestamp (oldest first)
	for i := 0; i < len(s.Data)-1; i++ {
		for j := i + 1; j < len(s.Data); j++ {
			if s.Data[i].Timestamp.After(s.Data[j].Timestamp) {
				s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
			}
		}
	}

	fmt.Printf("✅ Loaded %d REAL candles (%.1f days)\n", len(s.Data), float64(len(s.Data))/24.0)
	if len(s.Data) > 0 {
		fmt.Printf("📊 Latest price: $%.4f (Volume: %.2f)\n", s.Data[len(s.Data)-1].Close, s.Data[len(s.Data)-1].Volume)
		fmt.Printf("⏰ Data range: %s to %s\n",
			s.Data[0].Timestamp.Format("2006-01-02 15:04"),
			s.Data[len(s.Data)-1].Timestamp.Format("2006-01-02 15:04"))
	}

	return nil
}

func (s *PivotSupertrendEMAStrategy) CalculateIndicators() {
	if len(s.Data) < 100 {
		return
	}

	prd := 2      // Pivot Point Period
	factor := 3.0 // ATR Factor
	atrPd := 10   // ATR Period

	// Calculate basic indicators
	for i := range s.Data {
		// EMA100
		if i >= 99 {
			s.Data[i].EMA100 = s.calculateEMA(i, 100)
		}

		// RSI
		if i >= 14 {
			s.Data[i].RSI = s.calculateRSI(i, 14)
		}

		// ATR
		if i >= atrPd-1 {
			s.Data[i].ATR = s.calculateATR(i, atrPd)
		}

		// Current candle color
		s.Data[i].IsGreenCandle = s.Data[i].Close > s.Data[i].Open
		s.Data[i].IsRedCandle = s.Data[i].Close < s.Data[i].Open

		// Previous candle analysis
		if i > 0 {
			prev := s.Data[i-1]
			s.Data[i].PrevGreenCandle = prev.Close > prev.Open
			s.Data[i].PrevRedCandle = prev.Close < prev.Open
			if prev.EMA100 > 0 {
				s.Data[i].PrevAboveEMA100 = prev.Close > prev.EMA100
				s.Data[i].PrevBelowEMA100 = prev.Close < prev.EMA100
			}
		}

		// Pivot Points
		if i >= prd {
			s.calculatePivotPoints(i, prd)
		}
	}

	// Calculate Pivot Point SuperTrend
	s.calculatePivotSuperTrend(factor)
}

func (s *PivotSupertrendEMAStrategy) calculateEMA(index, period int) float64 {
	if index < period-1 {
		return 0
	}

	multiplier := 2.0 / (float64(period) + 1.0)

	if index == period-1 {
		// First EMA is SMA
		sum := 0.0
		for i := index - period + 1; i <= index; i++ {
			sum += s.Data[i].Close
		}
		return sum / float64(period)
	}

	// EMA = (Close * multiplier) + (Previous EMA * (1 - multiplier))
	prevEMA := s.Data[index-1].EMA100
	return (s.Data[index].Close * multiplier) + (prevEMA * (1 - multiplier))
}

func (s *PivotSupertrendEMAStrategy) calculateRSI(index, period int) float64 {
	if index < period {
		return 50.0
	}

	gains := 0.0
	losses := 0.0

	for i := index - period + 1; i <= index; i++ {
		change := s.Data[i].Close - s.Data[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += -change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return rsi
}

func (s *PivotSupertrendEMAStrategy) calculateATR(index, period int) float64 {
	if index < period {
		return 0
	}

	sum := 0.0
	for i := index - period + 1; i <= index; i++ {
		high := s.Data[i].High
		low := s.Data[i].Low
		prevClose := s.Data[i-1].Close

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)

		tr := math.Max(tr1, math.Max(tr2, tr3))
		sum += tr
	}

	return sum / float64(period)
}

func (s *PivotSupertrendEMAStrategy) calculatePivotPoints(index, period int) {
	// Pivot High
	if index >= period*2 {
		isHigh := true
		centerHigh := s.Data[index-period].High

		// Check if center is highest in the window
		for i := index - period*2; i <= index; i++ {
			if i != index-period && s.Data[i].High >= centerHigh {
				isHigh = false
				break
			}
		}

		if isHigh {
			s.Data[index-period].PivotHigh = centerHigh
		}
	}

	// Pivot Low
	if index >= period*2 {
		isLow := true
		centerLow := s.Data[index-period].Low

		// Check if center is lowest in the window
		for i := index - period*2; i <= index; i++ {
			if i != index-period && s.Data[i].Low <= centerLow {
				isLow = false
				break
			}
		}

		if isLow {
			s.Data[index-period].PivotLow = centerLow
		}
	}
}

func (s *PivotSupertrendEMAStrategy) calculatePivotSuperTrend(factor float64) {
	var center float64

	for i := range s.Data {
		// Update center line using pivot points
		if s.Data[i].PivotHigh > 0 || s.Data[i].PivotLow > 0 {
			var lastpp float64
			if s.Data[i].PivotHigh > 0 {
				lastpp = s.Data[i].PivotHigh
			} else {
				lastpp = s.Data[i].PivotLow
			}

			if center == 0 {
				center = lastpp
			} else {
				// Weighted calculation
				center = (center*2 + lastpp) / 3
			}
		}

		s.Data[i].Center = center

		if center > 0 && s.Data[i].ATR > 0 {
			// Calculate bands
			s.Data[i].UpperBand = center - (factor * s.Data[i].ATR)
			s.Data[i].LowerBand = center + (factor * s.Data[i].ATR)

			// Calculate SuperTrend
			if i > 0 {
				// TUp calculation
				tUp := s.Data[i].UpperBand
				if s.Data[i-1].Close > s.Data[i-1].SuperTrend && s.Data[i-1].Trend == 1 {
					tUp = math.Max(s.Data[i].UpperBand, s.Data[i-1].SuperTrend)
				}

				// TDown calculation
				tDown := s.Data[i].LowerBand
				if s.Data[i-1].Close < s.Data[i-1].SuperTrend && s.Data[i-1].Trend == -1 {
					tDown = math.Min(s.Data[i].LowerBand, s.Data[i-1].SuperTrend)
				}

				// Determine trend
				if s.Data[i].Close > s.Data[i-1].SuperTrend && s.Data[i-1].Trend == -1 {
					s.Data[i].Trend = 1
				} else if s.Data[i].Close < s.Data[i-1].SuperTrend && s.Data[i-1].Trend == 1 {
					s.Data[i].Trend = -1
				} else {
					s.Data[i].Trend = s.Data[i-1].Trend
				}

				// Set SuperTrend line
				if s.Data[i].Trend == 1 {
					s.Data[i].SuperTrend = tUp
				} else {
					s.Data[i].SuperTrend = tDown
				}

				// Generate signals
				s.Data[i].BuySignal = s.Data[i].Trend == 1 && s.Data[i-1].Trend == -1
				s.Data[i].SellSignal = s.Data[i].Trend == -1 && s.Data[i-1].Trend == 1
			}
		}
	}
}

func (s *PivotSupertrendEMAStrategy) GenerateSignal(index int) (string, float64) {
	if index < 100 || index >= len(s.Data) {
		return "HOLD", 0.0
	}

	current := s.Data[index]

	// Entry conditions (UPDATED RULES)

	// LONG: In Buy Trend + Previous candle is green + Previous candle above EMA100
	if current.Trend == 1 && current.PrevGreenCandle && current.PrevAboveEMA100 {
		confidence := 85.0
		return "LONG", confidence
	}

	// SHORT: In Sell Trend + Previous candle is red + Previous candle below EMA100
	if current.Trend == -1 && current.PrevRedCandle && current.PrevBelowEMA100 {
		confidence := 85.0
		return "SHORT", confidence
	}

	// Exit conditions
	if len(s.Positions) > 0 && s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
		lastPos := s.Positions[len(s.Positions)-1]

		// LONG exit: RSI > 70
		if lastPos.Side == "LONG" && current.RSI > 70 {
			return "EXIT", 90.0
		}

		// SHORT exit: RSI < 30
		if lastPos.Side == "SHORT" && current.RSI < 30 {
			return "EXIT", 90.0
		}

		// Time limit: 24 hours
		if time.Since(lastPos.EntryTime) > 24*time.Hour {
			return "EXIT", 80.0
		}
	}

	return "HOLD", 0.0
}

func (s *PivotSupertrendEMAStrategy) ExecuteBacktest() {
	fmt.Printf("🎯 Testing: Pivot Point SuperTrend + EMA100 Strategy (Updated) - %s\n", s.Symbol)

	if err := s.LoadRealMarketData(); err != nil {
		log.Printf("Error loading data: %v", err)
		return
	}

	s.CalculateIndicators()

	for i := 100; i < len(s.Data); i++ {
		signal, confidence := s.GenerateSignal(i)
		current := s.Data[i]

		if signal == "LONG" || signal == "SHORT" {
			// Only trade if no active position
			if len(s.Positions) == 0 || !s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
				if confidence >= 80.0 {
					reason := fmt.Sprintf("%s Signal (Trend: %d, PrevGreen: %t, PrevAboveEMA: %t, RSI: %.1f)",
						signal, current.Trend, current.PrevGreenCandle, current.PrevAboveEMA100, current.RSI)
					s.openPosition(signal, current.Close, current.Timestamp, reason)
				}
			}
		} else if signal == "EXIT" {
			// Close active position
			if len(s.Positions) > 0 && s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
				reason := fmt.Sprintf("RSI Exit (RSI: %.1f)", current.RSI)
				s.closePosition(current.Close, current.Timestamp, reason)
			}
		}

		// Risk management: Stop loss and take profit
		if len(s.Positions) > 0 && s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
			lastPos := &s.Positions[len(s.Positions)-1]

			// Calculate P&L
			var pnlPercent float64
			if lastPos.Side == "LONG" {
				pnlPercent = (current.Close - lastPos.EntryPrice) / lastPos.EntryPrice * 100
			} else {
				pnlPercent = (lastPos.EntryPrice - current.Close) / lastPos.EntryPrice * 100
			}

			// Stop loss: -3%
			if pnlPercent <= -3.0 {
				s.closePosition(current.Close, current.Timestamp, "Stop Loss (-3%)")
			}
			// Take profit: +5%
			if pnlPercent >= 5.0 {
				s.closePosition(current.Close, current.Timestamp, "Take Profit (+5%)")
			}
		}
	}

	// Close any remaining open position
	if len(s.Positions) > 0 && s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
		lastData := s.Data[len(s.Data)-1]
		s.closePosition(lastData.Close, lastData.Timestamp, "End of Data")
	}

	s.printResults()
}

func (s *PivotSupertrendEMAStrategy) openPosition(side string, price float64, timestamp time.Time, reason string) {
	riskPercent := 1.0 // 1% risk per trade
	positionSize := s.CurrentCash * riskPercent / 100

	position := Position{
		Symbol:     s.Symbol,
		Side:       side,
		EntryPrice: price,
		Size:       positionSize,
		EntryTime:  timestamp,
		Reason:     reason,
	}

	s.Positions = append(s.Positions, position)
	fmt.Printf("🎯 PIVOT SUPERTREND ENTRY %s: %.4f (Reason: %s)\n", side, price, reason)
}

func (s *PivotSupertrendEMAStrategy) closePosition(price float64, timestamp time.Time, reason string) {
	if len(s.Positions) == 0 {
		return
	}

	position := &s.Positions[len(s.Positions)-1]
	if !position.ExitTime.IsZero() {
		return
	}

	position.ExitPrice = price
	position.ExitTime = timestamp
	position.Duration = timestamp.Sub(position.EntryTime)

	// Calculate P&L
	if position.Side == "LONG" {
		position.PnL = (price - position.EntryPrice) / position.EntryPrice * position.Size
	} else {
		position.PnL = (position.EntryPrice - price) / position.EntryPrice * position.Size
	}

	s.CurrentCash += position.PnL

	pnlPercent := position.PnL / position.Size * 100
	fmt.Printf("🔹 PIVOT SUPERTREND EXIT %s: PnL %.2f%% - %s\n", position.Side, pnlPercent, reason)
}

func (s *PivotSupertrendEMAStrategy) printResults() {
	totalReturn := (s.CurrentCash - s.StartCash) / s.StartCash * 100

	winTrades := 0
	totalTrades := 0
	totalDuration := time.Duration(0)

	for _, pos := range s.Positions {
		if !pos.ExitTime.IsZero() {
			totalTrades++
			totalDuration += pos.Duration
			if pos.PnL > 0 {
				winTrades++
			}
		}
	}

	fmt.Printf("📊 Pivot Point SuperTrend + EMA100 Strategy RESULTS: %s\n", s.Symbol)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("💰 Return: $%.2f → $%.2f (%.2f%%)\n", s.StartCash, s.CurrentCash, totalReturn)

	if totalTrades > 0 {
		winRate := float64(winTrades) / float64(totalTrades) * 100
		avgDuration := totalDuration / time.Duration(totalTrades)
		fmt.Printf("📈 Trades: %d (Win: %d, Loss: %d, Rate: %.1f%%)\n", totalTrades, winTrades, totalTrades-winTrades, winRate)
		fmt.Printf("⏱️ Avg Duration: %v\n", avgDuration.Truncate(time.Minute))
	}

	if totalReturn > 0 {
		fmt.Printf("✅ PROFITABLE Pivot SuperTrend Strategy!\n")
	} else {
		fmt.Printf("📉 Needs optimization\n")
	}
}

func main() {
	fmt.Println("🚀 PIVOT POINT SUPERTREND + EMA100 STRATEGY TEST (UPDATED)")
	fmt.Println("📊 Updated: Previous candle must be above/below EMA100 AND correct color")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	symbols := []string{"BTC_USDT", "ETH_USDT", "SOL_USDT"}

	for _, symbol := range symbols {
		fmt.Printf("\n🔍 Testing Symbol: %s\n", symbol)
		fmt.Println("--------------------------------------------------")

		strategy := NewPivotSupertrendEMAStrategy(symbol, 10000.0)
		strategy.ExecuteBacktest()

		time.Sleep(1 * time.Second) // Rate limiting
	}

	fmt.Println("\n✅ Pivot Point SuperTrend + EMA100 strategy testing complete!")
}
