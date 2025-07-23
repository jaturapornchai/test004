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

// Multi-EMA SuperTrend Strategy
// Using EMA 20/50/100/200 for comprehensive trend analysis

type MultiEmaSupertrendStrategy struct {
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
	// Multiple EMA Indicators
	EMA20  float64 // Shortest EMA - immediate trend
	EMA50  float64 // Short EMA - short-term trend
	EMA100 float64 // Longer EMA - medium-term trend
	EMA200 float64 // Longest EMA - long-term trend
	RSI    float64 // RSI for exits
	ATR    float64 // ATR for SuperTrend
	// Simple SuperTrend
	UpperBand  float64
	LowerBand  float64
	Trend      int // 1 for uptrend, -1 for downtrend
	SuperTrend float64
	PrevClose  float64 // Previous candle close for color check
}

type GateIOResponse struct {
	Result [][]string `json:"result"`
}

func NewMultiEmaSupertrendStrategy(symbol string) *MultiEmaSupertrendStrategy {
	return &MultiEmaSupertrendStrategy{
		Symbol:      symbol,
		StartCash:   10000.0,
		CurrentCash: 10000.0,
	}
}

// Calculate EMA
func calculateEMA(values []float64, period int) []float64 {
	if len(values) < period {
		return make([]float64, len(values))
	}

	ema := make([]float64, len(values))
	multiplier := 2.0 / (float64(period) + 1.0)

	// Initialize first EMA as simple average
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += values[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate subsequent EMAs
	for i := period; i < len(values); i++ {
		ema[i] = (values[i] * multiplier) + (ema[i-1] * (1 - multiplier))
	}

	return ema
}

// Calculate RSI
func calculateRSI(closes []float64, period int) []float64 {
	if len(closes) < period+1 {
		return make([]float64, len(closes))
	}

	rsi := make([]float64, len(closes))
	gains := make([]float64, 0)
	losses := make([]float64, 0)

	// Calculate initial gains and losses
	for i := 1; i <= period; i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -change)
		}
	}

	// Calculate initial averages
	avgGain := 0.0
	avgLoss := 0.0
	for i := 0; i < period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI for initial period
	if avgLoss == 0 {
		rsi[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// Calculate subsequent RSI values
	for i := period + 1; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		var gain, loss float64

		if change > 0 {
			gain = change
			loss = 0
		} else {
			gain = 0
			loss = -change
		}

		avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi
}

// Calculate ATR
func calculateATR(highs, lows, closes []float64, period int) []float64 {
	if len(highs) < period {
		return make([]float64, len(highs))
	}

	atr := make([]float64, len(highs))
	tr := make([]float64, len(highs))

	// Calculate True Range
	for i := 1; i < len(highs); i++ {
		hl := highs[i] - lows[i]
		hc := math.Abs(highs[i] - closes[i-1])
		lc := math.Abs(lows[i] - closes[i-1])
		tr[i] = math.Max(hl, math.Max(hc, lc))
	}

	// Calculate initial ATR
	sum := 0.0
	for i := 1; i <= period; i++ {
		sum += tr[i]
	}
	atr[period] = sum / float64(period)

	// Calculate subsequent ATR values using smoothing
	for i := period + 1; i < len(highs); i++ {
		atr[i] = ((atr[i-1] * float64(period-1)) + tr[i]) / float64(period)
	}

	return atr
}

// Calculate Simple SuperTrend
func calculateSuperTrend(highs, lows, closes, atr []float64, factor float64) ([]float64, []float64, []int, []float64) {
	upperBands := make([]float64, len(closes))
	lowerBands := make([]float64, len(closes))
	trends := make([]int, len(closes))
	superTrends := make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		if atr[i] == 0 {
			continue
		}

		hl2 := (highs[i] + lows[i]) / 2
		upperBands[i] = hl2 + (factor * atr[i])
		lowerBands[i] = hl2 - (factor * atr[i])

		// SuperTrend calculation
		if i == 1 {
			trends[i] = 1
			superTrends[i] = lowerBands[i]
		} else {
			// Basic trend determination
			if closes[i] <= lowerBands[i-1] {
				trends[i] = -1
				superTrends[i] = upperBands[i]
			} else if closes[i] >= upperBands[i-1] {
				trends[i] = 1
				superTrends[i] = lowerBands[i]
			} else {
				trends[i] = trends[i-1]
				if trends[i] == 1 {
					superTrends[i] = lowerBands[i]
				} else {
					superTrends[i] = upperBands[i]
				}
			}
		}
	}

	return upperBands, lowerBands, trends, superTrends
}

// Fetch market data from Gate.io
func (s *MultiEmaSupertrendStrategy) fetchMarketData() error {
	// Fetch 20 days of hourly data (480 hours) for EMA 200
	url := fmt.Sprintf("https://api.gateio.ws/api/v4/spot/candlesticks?currency_pair=%s&interval=1h&limit=480", s.Symbol)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var gateResp [][]string
	if err := json.Unmarshal(body, &gateResp); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Convert to CandleData
	s.Data = make([]CandleData, 0, len(gateResp))
	for i := len(gateResp) - 1; i >= 0; i-- { // Reverse to get chronological order
		candle := gateResp[i]
		if len(candle) < 6 {
			continue
		}

		timestamp, _ := strconv.ParseInt(candle[0], 10, 64)
		volume, _ := strconv.ParseFloat(candle[1], 64)
		close, _ := strconv.ParseFloat(candle[2], 64)
		high, _ := strconv.ParseFloat(candle[3], 64)
		low, _ := strconv.ParseFloat(candle[4], 64)
		open, _ := strconv.ParseFloat(candle[5], 64)

		candleData := CandleData{
			Timestamp: time.Unix(timestamp, 0),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}

		s.Data = append(s.Data, candleData)
	}

	return nil
}

// Calculate all indicators
func (s *MultiEmaSupertrendStrategy) calculateIndicators() {
	if len(s.Data) < 250 {
		log.Printf("Insufficient data for calculations: %d candles", len(s.Data))
		return
	}

	// Extract price arrays
	closes := make([]float64, len(s.Data))
	highs := make([]float64, len(s.Data))
	lows := make([]float64, len(s.Data))

	for i, candle := range s.Data {
		closes[i] = candle.Close
		highs[i] = candle.High
		lows[i] = candle.Low
	}

	// Calculate multiple EMAs
	ema20Values := calculateEMA(closes, 20)
	ema50Values := calculateEMA(closes, 50)
	ema100Values := calculateEMA(closes, 100)
	ema200Values := calculateEMA(closes, 200)

	// Calculate RSI and ATR
	rsiValues := calculateRSI(closes, 14)
	atrValues := calculateATR(highs, lows, closes, 14)

	// Calculate SuperTrend
	upperBands, lowerBands, trends, superTrends := calculateSuperTrend(highs, lows, closes, atrValues, 2.5)

	// Assign values to candles
	for i := range s.Data {
		s.Data[i].EMA20 = ema20Values[i]
		s.Data[i].EMA50 = ema50Values[i]
		s.Data[i].EMA100 = ema100Values[i]
		s.Data[i].EMA200 = ema200Values[i]
		s.Data[i].RSI = rsiValues[i]
		s.Data[i].ATR = atrValues[i]
		s.Data[i].UpperBand = upperBands[i]
		s.Data[i].LowerBand = lowerBands[i]
		s.Data[i].Trend = trends[i]
		s.Data[i].SuperTrend = superTrends[i]

		if i > 0 {
			s.Data[i].PrevClose = s.Data[i-1].Close
		}
	}
}

// Check entry conditions with Multi-EMA filter
func (s *MultiEmaSupertrendStrategy) checkEntry(i int) (bool, string) {
	if i < 1 {
		return false, ""
	}

	current := s.Data[i]
	prev := s.Data[i-1]

	// LONG Entry Conditions (Simplified)
	if current.Trend == 1 && // SuperTrend bullish
		prev.Close > prev.Open && // Previous candle was green
		current.Close > current.EMA100 { // Price above EMA100
		return true, "LONG"
	}

	// SHORT Entry Conditions (Simplified)
	if current.Trend == -1 && // SuperTrend bearish
		prev.Close < prev.Open && // Previous candle was red
		current.Close < current.EMA100 { // Price below EMA100
		return true, "SHORT"
	}

	return false, ""
}

// Analyze EMA Alignment
func (s *MultiEmaSupertrendStrategy) analyzeEMAAlignment(candle CandleData) string {
	// For bullish alignment: EMA20 > EMA50 > EMA100 > EMA200
	// For bearish alignment: EMA20 < EMA50 < EMA100 < EMA200

	if candle.EMA20 > candle.EMA50 && candle.EMA50 > candle.EMA100 && candle.EMA100 > candle.EMA200 {
		return "BULLISH"
	}

	if candle.EMA20 < candle.EMA50 && candle.EMA50 < candle.EMA100 && candle.EMA100 < candle.EMA200 {
		return "BEARISH"
	}

	return "NEUTRAL"
}

// Check exit conditions
func (s *MultiEmaSupertrendStrategy) checkExit(position *Position, i int) (bool, string) {
	current := s.Data[i]

	// Calculate position duration
	duration := current.Timestamp.Sub(position.EntryTime)

	// Calculate current P&L
	var currentPrice float64
	if position.Side == "LONG" {
		currentPrice = current.Close
	} else {
		currentPrice = current.Close
	}

	pnlPercent := ((currentPrice - position.EntryPrice) / position.EntryPrice) * 100
	if position.Side == "SHORT" {
		pnlPercent = -pnlPercent
	}

	// RSI Exit Priority (Primary Exit Method)
	if position.Side == "LONG" && current.RSI > 70 {
		return true, "RSI_OVERBOUGHT"
	}
	if position.Side == "SHORT" && current.RSI < 30 {
		return true, "RSI_OVERSOLD"
	}

	// Stop Loss
	if pnlPercent <= -3.0 {
		return true, "STOP_LOSS"
	}

	// Take Profit
	if pnlPercent >= 5.0 {
		return true, "TAKE_PROFIT"
	}

	// Time Limit (24 hours = 24 * 60 * 60 seconds)
	if duration.Hours() >= 24 {
		return true, "TIME_LIMIT"
	}

	// SuperTrend Reversal
	if position.Side == "LONG" && current.Trend == -1 {
		return true, "SUPERTREND_REVERSAL"
	}
	if position.Side == "SHORT" && current.Trend == 1 {
		return true, "SUPERTREND_REVERSAL"
	}

	// EMA100 Break
	if position.Side == "LONG" && current.Close < current.EMA100*0.98 { // 2% below EMA100
		return true, "EMA100_BREAK"
	}
	if position.Side == "SHORT" && current.Close > current.EMA100*1.02 { // 2% above EMA100
		return true, "EMA100_BREAK"
	}

	return false, ""
}

// Execute trades
func (s *MultiEmaSupertrendStrategy) executeStrategy() {
	for i := 250; i < len(s.Data); i++ { // Start after enough data for indicators
		current := s.Data[i]

		// Check exits for existing positions
		for j := len(s.Positions) - 1; j >= 0; j-- {
			position := &s.Positions[j]
			if position.ExitPrice == 0 { // Position still open
				shouldExit, reason := s.checkExit(position, i)
				if shouldExit {
					// Close position
					position.ExitPrice = current.Close
					position.ExitTime = current.Timestamp
					position.Reason = reason
					position.Duration = position.ExitTime.Sub(position.EntryTime)

					// Calculate P&L
					if position.Side == "LONG" {
						position.PnL = ((position.ExitPrice - position.EntryPrice) / position.EntryPrice) * position.Size
					} else {
						position.PnL = ((position.EntryPrice - position.ExitPrice) / position.EntryPrice) * position.Size
					}

					s.CurrentCash += position.Size + position.PnL
				}
			}
		}

		// Check for new entries (only if no open positions)
		hasOpenPosition := false
		for _, pos := range s.Positions {
			if pos.ExitPrice == 0 {
				hasOpenPosition = true
				break
			}
		}

		if !hasOpenPosition {
			shouldEnter, side := s.checkEntry(i)
			if shouldEnter {
				// Calculate position size (1% of current cash)
				positionSize := s.CurrentCash * 0.01

				// Create new position
				position := Position{
					Symbol:     s.Symbol,
					Side:       side,
					EntryPrice: current.Close,
					Size:       positionSize,
					EntryTime:  current.Timestamp,
				}

				s.Positions = append(s.Positions, position)
				s.CurrentCash -= positionSize
			}
		}
	}
}

// Print results
func (s *MultiEmaSupertrendStrategy) printResults() {
	if len(s.Positions) == 0 {
		fmt.Printf("📊 Multi-EMA SuperTrend Strategy Results - %s\n", s.Symbol)
		fmt.Println("❌ No trades executed")
		return
	}

	totalPnL := 0.0
	winCount := 0
	lossCount := 0
	totalDuration := time.Duration(0)

	for _, pos := range s.Positions {
		if pos.ExitPrice > 0 {
			totalPnL += pos.PnL
			totalDuration += pos.Duration
			if pos.PnL > 0 {
				winCount++
			} else {
				lossCount++
			}
		}
	}

	totalTrades := winCount + lossCount
	winRate := 0.0
	if totalTrades > 0 {
		winRate = (float64(winCount) / float64(totalTrades)) * 100
	}

	avgDuration := time.Duration(0)
	if totalTrades > 0 {
		avgDuration = totalDuration / time.Duration(totalTrades)
	}

	finalBalance := s.CurrentCash + totalPnL
	totalReturn := ((finalBalance - s.StartCash) / s.StartCash) * 100

	fmt.Printf("\n🚀 Multi-EMA SuperTrend Strategy Results - %s\n", s.Symbol)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📈 Total Return: %.2f%%\n", totalReturn)
	fmt.Printf("💰 Final Balance: $%.2f (Start: $%.2f)\n", finalBalance, s.StartCash)
	fmt.Printf("📊 Total Trades: %d\n", totalTrades)
	fmt.Printf("✅ Win Rate: %.1f%% (%d wins, %d losses)\n", winRate, winCount, lossCount)
	fmt.Printf("⏱️ Average Duration: %v\n", avgDuration)
	fmt.Printf("💎 Total P&L: $%.2f\n", totalPnL)

	if totalTrades > 0 {
		fmt.Printf("\n📋 EMA Analysis:\n")
		lastCandle := s.Data[len(s.Data)-1]
		fmt.Printf("EMA 20: $%.2f\n", lastCandle.EMA20)
		fmt.Printf("EMA 50: $%.2f\n", lastCandle.EMA50)
		fmt.Printf("EMA 100: $%.2f\n", lastCandle.EMA100)
		fmt.Printf("EMA 200: $%.2f\n", lastCandle.EMA200)
		fmt.Printf("Alignment: %s\n", s.analyzeEMAAlignment(lastCandle))

		fmt.Printf("\n🔍 Recent Trades:\n")
		for i, pos := range s.Positions {
			if i >= len(s.Positions)-5 && pos.ExitPrice > 0 { // Show last 5 trades
				pnlPercent := (pos.PnL / pos.Size) * 100
				fmt.Printf("%d. %s: Entry $%.2f → Exit $%.2f (%.2f%%) - %s - %v\n",
					i+1, pos.Side, pos.EntryPrice, pos.ExitPrice, pnlPercent, pos.Reason, pos.Duration)
			}
		}
	}
	fmt.Println()
}

func main() {
	symbols := []string{"BTC_USDT", "ETH_USDT", "SOL_USDT"}

	fmt.Println("🚀 Starting Multi-EMA SuperTrend Strategy Analysis...")
	fmt.Println("📊 Using EMA 20/50/100/200 for comprehensive trend analysis")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, symbol := range symbols {
		fmt.Printf("\n🔍 Analyzing %s...\n", symbol)

		strategy := NewMultiEmaSupertrendStrategy(symbol)

		// Fetch market data
		if err := strategy.fetchMarketData(); err != nil {
			log.Printf("Error fetching data for %s: %v", symbol, err)
			continue
		}

		// Calculate indicators
		strategy.calculateIndicators()

		// Execute strategy
		strategy.executeStrategy()

		// Print results
		strategy.printResults()

		// Wait between API calls
		time.Sleep(1 * time.Second)
	}

	fmt.Println("✅ Multi-EMA SuperTrend Analysis Complete!")
}
