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

// Volume Breakout + EMA Strategy
// Combines volume breakout signals with EMA trend confirmation
// Enhanced version of the successful Volume Breakout Strategy

type VolumeBreakoutEMAStrategy struct {
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
	// EMA indicators
	EMA9  float64 // Fast EMA
	EMA21 float64 // Mid EMA
	EMA50 float64 // Slow EMA
	// Volume indicators
	VolumeMA20  float64 // 20-period volume moving average
	VolumeRatio float64 // Current volume / VolumeMA20
	// Combined signals
	TrendStrength float64 // EMA alignment strength
	VolumeSignal  string  // "HIGH_BULL", "HIGH_BEAR", "NORMAL"
}

func NewVolumeBreakoutEMAStrategy(symbol string, startCash float64) *VolumeBreakoutEMAStrategy {
	return &VolumeBreakoutEMAStrategy{
		Symbol:      symbol,
		StartCash:   startCash,
		CurrentCash: startCash,
		Positions:   make([]Position, 0),
	}
}

func (s *VolumeBreakoutEMAStrategy) LoadRealMarketData() error {
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

func (s *VolumeBreakoutEMAStrategy) CalculateIndicators() {
	if len(s.Data) < 50 {
		return
	}

	// Calculate EMAs
	for i := range s.Data {
		if i >= 8 {
			s.Data[i].EMA9 = s.calculateEMA(i, 9)
		}
		if i >= 20 {
			s.Data[i].EMA21 = s.calculateEMA(i, 21)
		}
		if i >= 49 {
			s.Data[i].EMA50 = s.calculateEMA(i, 50)
		}
		if i >= 19 {
			s.Data[i].VolumeMA20 = s.calculateVolumeMA(i, 20)
			if s.Data[i].VolumeMA20 > 0 {
				s.Data[i].VolumeRatio = s.Data[i].Volume / s.Data[i].VolumeMA20
			}
		}

		// Calculate trend strength and volume signals
		s.calculateCombinedSignals(i)
	}
}

func (s *VolumeBreakoutEMAStrategy) calculateEMA(index, period int) float64 {
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
	var prevEMA float64
	switch period {
	case 9:
		prevEMA = s.Data[index-1].EMA9
	case 21:
		prevEMA = s.Data[index-1].EMA21
	case 50:
		prevEMA = s.Data[index-1].EMA50
	}

	return (s.Data[index].Close * multiplier) + (prevEMA * (1 - multiplier))
}

func (s *VolumeBreakoutEMAStrategy) calculateVolumeMA(index, period int) float64 {
	if index < period-1 {
		return 0
	}

	sum := 0.0
	for i := index - period + 1; i <= index; i++ {
		sum += s.Data[i].Volume
	}
	return sum / float64(period)
}

func (s *VolumeBreakoutEMAStrategy) calculateCombinedSignals(index int) {
	if index < 49 { // Need at least 50 candles for all EMAs
		return
	}

	current := &s.Data[index]

	// Calculate EMA trend strength
	ema9 := current.EMA9
	ema21 := current.EMA21
	ema50 := current.EMA50

	// Trend strength: all EMAs aligned = strong trend
	trendStrength := 0.0
	if ema9 > ema21 && ema21 > ema50 {
		// Bullish alignment
		trendStrength = (ema9 - ema50) / ema50 * 100
	} else if ema9 < ema21 && ema21 < ema50 {
		// Bearish alignment
		trendStrength = (ema50 - ema9) / ema50 * 100 * -1
	}
	current.TrendStrength = trendStrength

	// Volume signal classification
	if current.VolumeRatio >= 2.0 {
		// High volume breakout
		priceChange := (current.Close - current.Open) / current.Open * 100
		if priceChange >= 1.0 {
			current.VolumeSignal = "HIGH_BULL"
		} else if priceChange <= -1.0 {
			current.VolumeSignal = "HIGH_BEAR"
		} else {
			current.VolumeSignal = "HIGH_NEUTRAL"
		}
	} else {
		current.VolumeSignal = "NORMAL"
	}
}

func (s *VolumeBreakoutEMAStrategy) GenerateSignal(index int) (string, float64) {
	if index < 50 || index >= len(s.Data) {
		return "HOLD", 0.0
	}

	current := s.Data[index]
	previous := s.Data[index-1]

	// Enhanced entry conditions: Volume Breakout + EMA Trend Confirmation

	// LONG Entry: High volume bullish breakout + bullish EMA trend
	if current.VolumeSignal == "HIGH_BULL" && current.TrendStrength > 1.0 {
		// Additional confirmation: price above EMA21
		if current.Close > current.EMA21 {
			confidence := math.Min(90.0, 60.0+current.TrendStrength+(current.VolumeRatio-2.0)*10)
			return "LONG", confidence
		}
	}

	// SHORT Entry: High volume bearish breakout + bearish EMA trend
	if current.VolumeSignal == "HIGH_BEAR" && current.TrendStrength < -1.0 {
		// Additional confirmation: price below EMA21
		if current.Close < current.EMA21 {
			confidence := math.Min(90.0, 60.0+math.Abs(current.TrendStrength)+(current.VolumeRatio-2.0)*10)
			return "SHORT", confidence
		}
	}

	// Exit conditions for existing positions
	if len(s.Positions) > 0 && s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
		lastPos := s.Positions[len(s.Positions)-1]

		// Check time limit (12 hours for faster strategy)
		if time.Since(lastPos.EntryTime) > 12*time.Hour {
			return "EXIT", 80.0
		}

		// Trend reversal exits
		if lastPos.Side == "LONG" {
			// Exit LONG: EMA bearish crossover or price below EMA21
			if (current.EMA9 < current.EMA21 && previous.EMA9 >= previous.EMA21) ||
				(current.Close < current.EMA21 && current.TrendStrength < 0) {
				return "EXIT", 85.0
			}
		} else if lastPos.Side == "SHORT" {
			// Exit SHORT: EMA bullish crossover or price above EMA21
			if (current.EMA9 > current.EMA21 && previous.EMA9 <= previous.EMA21) ||
				(current.Close > current.EMA21 && current.TrendStrength > 0) {
				return "EXIT", 85.0
			}
		}
	}

	return "HOLD", 0.0
}

func (s *VolumeBreakoutEMAStrategy) ExecuteBacktest() {
	fmt.Printf("🎯 Testing: Volume Breakout + EMA Strategy - %s\n", s.Symbol)

	if err := s.LoadRealMarketData(); err != nil {
		log.Printf("Error loading data: %v", err)
		return
	}

	s.CalculateIndicators()

	for i := 50; i < len(s.Data); i++ {
		signal, confidence := s.GenerateSignal(i)
		current := s.Data[i]

		if signal == "LONG" || signal == "SHORT" {
			// Only trade if no active position
			if len(s.Positions) == 0 || !s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
				if confidence >= 75.0 { // Higher threshold for combined strategy
					s.openPosition(signal, current.Close, current.Timestamp, fmt.Sprintf("Volume Breakout + EMA Trend (%.1f%% confidence)", confidence))
				}
			}
		} else if signal == "EXIT" {
			// Close active position
			if len(s.Positions) > 0 && s.Positions[len(s.Positions)-1].ExitTime.IsZero() {
				s.closePosition(current.Close, current.Timestamp, "Strategy Exit Signal")
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

			// Stop loss: -1%
			if pnlPercent <= -1.0 {
				s.closePosition(current.Close, current.Timestamp, "Stop Loss")
			}
			// Take profit: +2%
			if pnlPercent >= 2.0 {
				s.closePosition(current.Close, current.Timestamp, "Take Profit")
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

func (s *VolumeBreakoutEMAStrategy) openPosition(side string, price float64, timestamp time.Time, reason string) {
	riskPercent := 0.5 // 0.5% risk per trade
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
	fmt.Printf("🎯 VOLUME+EMA ENTRY %s: %.4f (Reason: %s)\n", side, price, reason)
}

func (s *VolumeBreakoutEMAStrategy) closePosition(price float64, timestamp time.Time, reason string) {
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
	fmt.Printf("🔹 VOLUME+EMA EXIT %s: PnL %.2f%% - %s\n", position.Side, pnlPercent, reason)
}

func (s *VolumeBreakoutEMAStrategy) printResults() {
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

	fmt.Printf("📊 Volume Breakout + EMA Strategy RESULTS: %s\n", s.Symbol)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("💰 Return: $%.2f → $%.2f (%.2f%%)\n", s.StartCash, s.CurrentCash, totalReturn)

	if totalTrades > 0 {
		winRate := float64(winTrades) / float64(totalTrades) * 100
		avgDuration := totalDuration / time.Duration(totalTrades)
		fmt.Printf("📈 Trades: %d (Win: %d, Loss: %d, Rate: %.1f%%)\n", totalTrades, winTrades, totalTrades-winTrades, winRate)
		fmt.Printf("⏱️ Avg Duration: %v\n", avgDuration.Truncate(time.Minute))
	}

	if totalReturn > 0 {
		fmt.Printf("✅ PROFITABLE Enhanced Strategy!\n")
	} else {
		fmt.Printf("📉 Needs optimization\n")
	}
}

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
