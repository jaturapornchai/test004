package main

import (
	"fmt"
	"math"
	"time"

	"gateio-trading-bot/internal/gateio"
	"gateio-trading-bot/internal/trading"
)

// Enhanced1HDataManager - ระบบจัดการข้อมูล 1H ขั้นสูง
type Enhanced1HDataManager struct {
	symbol         string
	client         *gateio.Client
	dataBuffer     []trading.OHLCV
	currentIndex   int
	maxLookback    int // 144 timeframes
	isRealTime     bool
	lastUpdateTime time.Time

	// Cache สำหรับ indicators
	emaCache    map[int][]float64 // EMA cache
	rsiCache    []float64
	atrCache    []float64
	volumeCache []float64

	// Performance tracking
	cacheHits   int
	cacheMisses int
}

// NewEnhanced1HDataManager สร้าง data manager ใหม่
func NewEnhanced1HDataManager(symbol string, client *gateio.Client) *Enhanced1HDataManager {
	return &Enhanced1HDataManager{
		symbol:      symbol,
		client:      client,
		maxLookback: 144,                           // 6 วัน × 24 ชั่วโมง
		dataBuffer:  make([]trading.OHLCV, 0, 200), // เผื่อพื้นที่เพิ่ม
		emaCache:    make(map[int][]float64),
		isRealTime:  false,
	}
}

// LoadHistoricalData โหลดข้อมูลย้อนหลัง 144 timeframes
func (dm *Enhanced1HDataManager) LoadHistoricalData() error {
	fmt.Printf("🔄 Loading 144H historical data for %s...\n", dm.symbol)

	// กำหนดเวลาย้อนหลัง 150 ชั่วโมง (เผื่อไว้)
	endTime := time.Now()
	startTime := endTime.Add(-150 * time.Hour)

	// วิธีที่ 1: ใช้ Demo Data แบบ Smart
	data := dm.generateSmartDemoData(startTime, endTime)

	if len(data) < dm.maxLookback {
		return fmt.Errorf("insufficient data: got %d, need %d", len(data), dm.maxLookback)
	}

	// เก็บแค่ 144 timeframes ล่าสุด
	if len(data) > dm.maxLookback {
		data = data[len(data)-dm.maxLookback:]
	}

	dm.dataBuffer = data
	dm.currentIndex = len(data) - 1
	dm.lastUpdateTime = time.Unix(data[len(data)-1].Timestamp, 0)

	fmt.Printf("✅ Loaded %d candles (%.1f days)\n",
		len(dm.dataBuffer), float64(len(dm.dataBuffer))/24.0)

	// Pre-calculate indicators
	dm.preCalculateIndicators()

	return nil
}

// generateSmartDemoData สร้างข้อมูล demo ที่มีรูปแบบการเทรดจริง
func (dm *Enhanced1HDataManager) generateSmartDemoData(start, end time.Time) []trading.OHLCV {
	duration := end.Sub(start)
	hours := int(duration.Hours())

	data := make([]trading.OHLCV, hours)

	// ราคาเริ่มต้นตาม symbol
	var basePrice float64
	var volatility float64

	switch dm.symbol {
	case "SOL_USDT":
		basePrice = 180.0
		volatility = 0.025 // 2.5%
	case "BTC_USDT":
		basePrice = 92000.0
		volatility = 0.02 // 2%
	case "ETH_USDT":
		basePrice = 3300.0
		volatility = 0.022 // 2.2%
	default:
		basePrice = 100.0
		volatility = 0.02
	}

	// สร้าง trend patterns ที่เหมือนจริง
	trendPatterns := []TrendPattern{
		{Type: "sideways", Duration: 20, Strength: 0.1},
		{Type: "bullish", Duration: 30, Strength: 0.4},
		{Type: "bearish", Duration: 25, Strength: -0.3},
		{Type: "sideways", Duration: 15, Strength: 0.05},
		{Type: "bullish", Duration: 40, Strength: 0.6},
		{Type: "correction", Duration: 20, Strength: -0.2},
	}

	currentPrice := basePrice
	patternIndex := 0
	patternProgress := 0

	for i := 0; i < hours; i++ {
		timestamp := start.Add(time.Duration(i) * time.Hour)

		// เลือก pattern ปัจจุบัน
		if patternIndex < len(trendPatterns) {
			pattern := trendPatterns[patternIndex]

			// คำนวณการเปลี่ยนแปลงราคาตาม pattern
			priceChange := dm.calculatePatternChange(pattern, patternProgress, volatility)

			// เพิ่ม noise แบบสุ่ม
			noise := (float64(i%17-8) / 100.0) * volatility * currentPrice

			// Market hours effect (เพิ่ม volatility ใน US/EU hours)
			hour := timestamp.Hour()
			if (hour >= 8 && hour <= 16) || (hour >= 22 || hour <= 2) {
				noise *= 1.5 // เพิ่ม volatility
			}

			currentPrice = currentPrice*(1+priceChange) + noise

			// สร้าง OHLC
			open := currentPrice
			high, low, close := dm.generateOHLC(open, volatility, i)
			currentPrice = close

			// Volume ที่เหมือนจริง
			volume := dm.generateRealisticVolume(hour, pattern.Type, i)

			data[i] = trading.OHLCV{
				Timestamp: timestamp.Unix(),
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    volume,
			}

			patternProgress++
			if patternProgress >= pattern.Duration {
				patternIndex++
				patternProgress = 0
			}
		}
	}

	return data
}

type TrendPattern struct {
	Type     string  // "bullish", "bearish", "sideways", "correction"
	Duration int     // จำนวน hours
	Strength float64 // แรงของ trend (-1 ถึง 1)
}

// calculatePatternChange คำนวณการเปลี่ยนแปลงราคาตาม pattern
func (dm *Enhanced1HDataManager) calculatePatternChange(pattern TrendPattern, progress int, baseVol float64) float64 {
	progressRatio := float64(progress) / float64(pattern.Duration)

	switch pattern.Type {
	case "bullish":
		// Exponential growth แล้วค่อยๆ ลด
		intensity := math.Pow(1-progressRatio, 2)
		return pattern.Strength * baseVol * intensity

	case "bearish":
		// Sharp drop แล้วค่อยๆ recover
		intensity := math.Pow(progressRatio, 0.5)
		return pattern.Strength * baseVol * intensity

	case "sideways":
		// Oscillation รอบราคา
		oscillation := math.Sin(progressRatio * 4 * math.Pi)
		return pattern.Strength * baseVol * oscillation

	case "correction":
		// Quick drop แล้ว consolidate
		if progressRatio < 0.3 {
			return pattern.Strength * baseVol * 2
		} else {
			return pattern.Strength * baseVol * 0.2
		}

	default:
		return 0
	}
}

// generateOHLC สร้าง OHLC จาก price และ volatility
func (dm *Enhanced1HDataManager) generateOHLC(price, vol float64, index int) (high, low, close float64) {
	maxMove := price * vol

	// สุ่มการเคลื่อนไหวแบบ realistic
	highMove := maxMove * (0.3 + float64(index%7)/10.0)
	lowMove := maxMove * (0.3 + float64((index+3)%7)/10.0)

	high = price + highMove
	low = price - lowMove

	// Close ระหว่าง low กับ high โดยมี bias ตาม trend
	closeRatio := 0.3 + float64(index%13)/20.0 // 0.3 - 0.95
	close = low + (high-low)*closeRatio

	return high, low, close
}

// generateRealisticVolume สร้าง volume ที่เหมือนจริง
func (dm *Enhanced1HDataManager) generateRealisticVolume(hour int, trendType string, index int) float64 {
	baseVolume := 1000000.0

	// Market hours effect
	var timeMultiplier float64
	if hour >= 8 && hour <= 16 { // EU/US morning
		timeMultiplier = 1.8
	} else if hour >= 22 || hour <= 2 { // US evening
		timeMultiplier = 1.5
	} else if hour >= 3 && hour <= 7 { // Asian morning
		timeMultiplier = 1.2
	} else {
		timeMultiplier = 0.8 // Low activity
	}

	// Trend effect
	var trendMultiplier float64
	switch trendType {
	case "bullish":
		trendMultiplier = 1.4
	case "bearish":
		trendMultiplier = 1.6 // Higher volume in sell-offs
	case "correction":
		trendMultiplier = 2.0 // Panic selling
	default:
		trendMultiplier = 1.0
	}

	// Random variation
	randomFactor := 0.7 + float64(index%31)/50.0 // 0.7 - 1.32

	return baseVolume * timeMultiplier * trendMultiplier * randomFactor
}

// preCalculateIndicators คำนวณ indicators ทั้งหมดล่วงหน้า
func (dm *Enhanced1HDataManager) preCalculateIndicators() {
	fmt.Printf("🔧 Pre-calculating indicators...\n")

	if len(dm.dataBuffer) < 50 {
		return
	}

	// Calculate EMAs
	periods := []int{9, 21, 50}
	for _, period := range periods {
		ema := dm.calculateEMAFast(period)
		dm.emaCache[period] = ema
	}

	// Calculate RSI
	dm.rsiCache = dm.calculateRSIFast(14)

	// Calculate ATR
	dm.atrCache = dm.calculateATRFast(14)

	// Calculate Volume MA
	dm.volumeCache = dm.calculateVolumeMAFast(20)

	fmt.Printf("✅ Indicators ready (EMA: %d, RSI: %d, ATR: %d)\n",
		len(dm.emaCache), len(dm.rsiCache), len(dm.atrCache))
}

// calculateEMAFast คำนวณ EMA แบบเร็ว
func (dm *Enhanced1HDataManager) calculateEMAFast(period int) []float64 {
	data := dm.dataBuffer
	result := make([]float64, len(data))

	if len(data) == 0 {
		return result
	}

	multiplier := 2.0 / float64(period+1)

	// เริ่มต้นด้วย SMA
	sum := 0.0
	for i := 0; i < period && i < len(data); i++ {
		sum += data[i].Close
	}
	result[period-1] = sum / float64(period)

	// คำนวณ EMA
	for i := period; i < len(data); i++ {
		result[i] = (data[i].Close * multiplier) + (result[i-1] * (1 - multiplier))
	}

	return result
}

// calculateRSIFast คำนวณ RSI แบบเร็ว
func (dm *Enhanced1HDataManager) calculateRSIFast(period int) []float64 {
	data := dm.dataBuffer
	result := make([]float64, len(data))

	if len(data) <= period {
		return result
	}

	gains := make([]float64, 0)
	losses := make([]float64, 0)

	// คำนวณ gains/losses
	for i := 1; i < len(data); i++ {
		change := data[i].Close - data[i-1].Close
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -change)
		}
	}

	// คำนวณ RSI
	for i := period; i < len(data); i++ {
		avgGain := average(gains[i-period : i])
		avgLoss := average(losses[i-period : i])

		if avgLoss == 0 {
			result[i] = 100
		} else {
			rs := avgGain / avgLoss
			result[i] = 100 - (100 / (1 + rs))
		}
	}

	return result
}

// calculateATRFast คำนวณ ATR แบบเร็ว
func (dm *Enhanced1HDataManager) calculateATRFast(period int) []float64 {
	data := dm.dataBuffer
	result := make([]float64, len(data))

	if len(data) <= period {
		return result
	}

	// คำนวณ True Range
	tr := make([]float64, len(data))
	for i := 1; i < len(data); i++ {
		current := data[i]
		previous := data[i-1]

		tr1 := current.High - current.Low
		tr2 := math.Abs(current.High - previous.Close)
		tr3 := math.Abs(current.Low - previous.Close)

		tr[i] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// คำนวณ ATR (Simple Moving Average ของ TR)
	for i := period; i < len(data); i++ {
		result[i] = average(tr[i-period+1 : i+1])
	}

	return result
}

// calculateVolumeMAFast คำนวณ Volume Moving Average แบบเร็ว
func (dm *Enhanced1HDataManager) calculateVolumeMAFast(period int) []float64 {
	data := dm.dataBuffer
	result := make([]float64, len(data))

	if len(data) <= period {
		return result
	}

	for i := period; i < len(data); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += data[j].Volume
		}
		result[i] = sum / float64(period)
	}

	return result
}

// GetLatestAnalysis ดึงการวิเคราะห์ล่าสุด
func (dm *Enhanced1HDataManager) GetLatestAnalysis() *Enhanced1HAnalysis {
	if len(dm.dataBuffer) < 50 {
		return nil
	}

	idx := len(dm.dataBuffer) - 1
	current := dm.dataBuffer[idx]

	analysis := &Enhanced1HAnalysis{
		Timestamp:   current.Timestamp,
		Price:       current.Close,
		EMA9:        dm.getLatestEMA(9),
		EMA21:       dm.getLatestEMA(21),
		EMA50:       dm.getLatestEMA(50),
		RSI:         dm.getLatestRSI(),
		ATR:         dm.getLatestATR(),
		VolumeMA:    dm.getLatestVolumeMA(),
		CurrentVol:  current.Volume,
		DataQuality: dm.assessDataQuality(),
	}

	// วิเคราะห์ signal
	analysis.Signal = dm.analyzeSignal(analysis)
	analysis.Confidence = dm.calculateConfidence(analysis)

	return analysis
}

type Enhanced1HAnalysis struct {
	Timestamp   int64   `json:"timestamp"`
	Price       float64 `json:"price"`
	EMA9        float64 `json:"ema9"`
	EMA21       float64 `json:"ema21"`
	EMA50       float64 `json:"ema50"`
	RSI         float64 `json:"rsi"`
	ATR         float64 `json:"atr"`
	VolumeMA    float64 `json:"volume_ma"`
	CurrentVol  float64 `json:"current_volume"`
	Signal      string  `json:"signal"`
	Confidence  float64 `json:"confidence"`
	DataQuality float64 `json:"data_quality"`
}

// Helper functions
func (dm *Enhanced1HDataManager) getLatestEMA(period int) float64 {
	if ema, exists := dm.emaCache[period]; exists && len(ema) > 0 {
		dm.cacheHits++
		return ema[len(ema)-1]
	}
	dm.cacheMisses++
	return 0
}

func (dm *Enhanced1HDataManager) getLatestRSI() float64 {
	if len(dm.rsiCache) > 0 {
		return dm.rsiCache[len(dm.rsiCache)-1]
	}
	return 50 // Neutral
}

func (dm *Enhanced1HDataManager) getLatestATR() float64 {
	if len(dm.atrCache) > 0 {
		return dm.atrCache[len(dm.atrCache)-1]
	}
	return dm.dataBuffer[len(dm.dataBuffer)-1].Close * 0.02 // 2% fallback
}

func (dm *Enhanced1HDataManager) getLatestVolumeMA() float64 {
	if len(dm.volumeCache) > 0 {
		return dm.volumeCache[len(dm.volumeCache)-1]
	}
	return 1000000 // Default
}

func (dm *Enhanced1HDataManager) analyzeSignal(analysis *Enhanced1HAnalysis) string {
	// Triple EMA alignment
	if analysis.EMA9 > analysis.EMA21 && analysis.EMA21 > analysis.EMA50 {
		if analysis.Price >= analysis.EMA9*0.998 && analysis.RSI > 40 && analysis.RSI < 80 {
			return "LONG"
		}
	}

	if analysis.EMA9 < analysis.EMA21 && analysis.EMA21 < analysis.EMA50 {
		if analysis.Price <= analysis.EMA9*1.002 && analysis.RSI > 20 && analysis.RSI < 60 {
			return "SHORT"
		}
	}

	return "NEUTRAL"
}

func (dm *Enhanced1HDataManager) calculateConfidence(analysis *Enhanced1HAnalysis) float64 {
	confidence := 60.0

	// EMA alignment strength
	emaSpread := math.Abs(analysis.EMA9-analysis.EMA50) / analysis.EMA50 * 100
	if emaSpread > 0.15 {
		confidence += 15.0
	}

	// Volume confirmation
	if analysis.CurrentVol > analysis.VolumeMA*1.2 {
		confidence += 10.0
	}

	// RSI confirmation
	if analysis.Signal == "LONG" && analysis.RSI >= 50 && analysis.RSI <= 70 {
		confidence += 15.0
	}
	if analysis.Signal == "SHORT" && analysis.RSI >= 30 && analysis.RSI <= 50 {
		confidence += 15.0
	}

	return math.Min(confidence, 100.0)
}

func (dm *Enhanced1HDataManager) assessDataQuality() float64 {
	// ประเมินคุณภาพข้อมูล
	quality := 100.0

	// Check data completeness
	if len(dm.dataBuffer) < dm.maxLookback {
		quality -= 20.0
	}

	// Check cache efficiency
	totalAccess := dm.cacheHits + dm.cacheMisses
	if totalAccess > 0 {
		hitRate := float64(dm.cacheHits) / float64(totalAccess)
		if hitRate < 0.8 {
			quality -= 10.0
		}
	}

	return math.Max(quality, 0.0)
}

func average(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// GetCacheStats ดูสถิติ cache
func (dm *Enhanced1HDataManager) GetCacheStats() (hits, misses int, hitRate float64) {
	total := dm.cacheHits + dm.cacheMisses
	if total > 0 {
		hitRate = float64(dm.cacheHits) / float64(total) * 100
	}
	return dm.cacheHits, dm.cacheMisses, hitRate
}

// GetDataSummary ดูสรุปข้อมูล
func (dm *Enhanced1HDataManager) GetDataSummary() map[string]interface{} {
	latest := dm.GetLatestAnalysis()
	hits, misses, hitRate := dm.GetCacheStats()

	return map[string]interface{}{
		"symbol":         dm.symbol,
		"data_points":    len(dm.dataBuffer),
		"lookback_hours": dm.maxLookback,
		"lookback_days":  float64(dm.maxLookback) / 24.0,
		"cache_hits":     hits,
		"cache_misses":   misses,
		"cache_hit_rate": hitRate,
		"data_quality":   latest.DataQuality,
		"latest_price":   latest.Price,
		"latest_signal":  latest.Signal,
		"confidence":     latest.Confidence,
		"last_update":    time.Unix(latest.Timestamp, 0).Format("2006-01-02 15:04:05"),
	}
}
