package trading

import (
	"math"
)

// Indicators สำหรับคำนวณ technical indicators
type Indicators struct {
	pivotPeriod int
	atrPeriod   int
	atrFactor   float64
	emaPeriod   int
}

// NewIndicators สร้าง Indicators ใหม่
func NewIndicators() *Indicators {
	return &Indicators{
		pivotPeriod: 2,   // Pivot Point Period
		atrPeriod:   10,  // ATR Period
		atrFactor:   3.0, // ATR Factor
		emaPeriod:   100, // EMA Period
	}
}

// AnalyzePivotPointSuperTrend วิเคราะห์ Pivot Point SuperTrend + EMA100
func (ind *Indicators) AnalyzePivotPointSuperTrend(ohlcv []OHLCV) *SuperTrendAnalysis {
	if len(ohlcv) < ind.emaPeriod {
		return &SuperTrendAnalysis{
			Trend:           0,
			Signal:          "NEUTRAL",
			Confidence:      0,
			RiskRewardRatio: 0,
		}
	}

	// คำนวณ Pivot Points
	pivotHighs, pivotLows := ind.calculatePivotPoints(ohlcv)

	// คำนวณ Center Line จาก Pivot Points
	centerLine := ind.calculateCenterLine(pivotHighs, pivotLows, len(ohlcv))

	// คำนวณ ATR
	atr := ind.calculateATR(ohlcv)

	// คำนวณ SuperTrend Bands
	upperBand, lowerBand := ind.calculateSuperTrendBands(centerLine, atr)

	// คำนวณ SuperTrend และ Trend
	superTrend, trend := ind.calculateSuperTrend(ohlcv, upperBand, lowerBand)

	// คำนวณ EMA100
	ema100 := ind.calculateEMA(ohlcv, ind.emaPeriod)

	// วิเคราะห์ Signal และ Confidence
	signal, confidence := ind.analyzeSignal(ohlcv, trend, ema100, superTrend)

	// คำนวณ Risk-Reward Ratio
	riskRewardRatio := ind.calculateRiskReward(ohlcv, superTrend, atr)

	currentPrice := ohlcv[len(ohlcv)-1].Close
	currentTrend := trend[len(trend)-1]
	currentSuperTrend := superTrend[len(superTrend)-1]
	currentEMA100 := ema100[len(ema100)-1]
	currentATR := atr[len(atr)-1]

	return &SuperTrendAnalysis{
		Trend:           currentTrend,
		Signal:          signal,
		Confidence:      confidence,
		RiskRewardRatio: riskRewardRatio,
		SuperTrendValue: currentSuperTrend,
		EMA100:          currentEMA100,
		CurrentPrice:    currentPrice,
		ATR:             currentATR,
	}
}

// calculatePivotPoints คำนวณ Pivot High/Low
func (ind *Indicators) calculatePivotPoints(ohlcv []OHLCV) ([]float64, []float64) {
	period := ind.pivotPeriod
	highs := make([]float64, len(ohlcv))
	lows := make([]float64, len(ohlcv))

	for i := range ohlcv {
		highs[i] = ohlcv[i].High
		lows[i] = ohlcv[i].Low
	}

	pivotHighs := make([]float64, len(ohlcv))
	pivotLows := make([]float64, len(ohlcv))

	for i := period; i < len(ohlcv)-period; i++ {
		// ตรวจสอบ Pivot High
		isPivotHigh := true
		for j := i - period; j <= i+period; j++ {
			if j != i && highs[j] >= highs[i] {
				isPivotHigh = false
				break
			}
		}
		if isPivotHigh {
			pivotHighs[i] = highs[i]
		}

		// ตรวจสอบ Pivot Low
		isPivotLow := true
		for j := i - period; j <= i+period; j++ {
			if j != i && lows[j] <= lows[i] {
				isPivotLow = false
				break
			}
		}
		if isPivotLow {
			pivotLows[i] = lows[i]
		}
	}

	return pivotHighs, pivotLows
}

// calculateCenterLine คำนวณ Center Line จาก Pivot Points
func (ind *Indicators) calculateCenterLine(pivotHighs, pivotLows []float64, length int) []float64 {
	centerLine := make([]float64, length)
	var center float64

	for i := 0; i < length; i++ {
		var lastPivot float64

		if pivotHighs[i] != 0 {
			lastPivot = pivotHighs[i]
		} else if pivotLows[i] != 0 {
			lastPivot = pivotLows[i]
		}

		if lastPivot != 0 {
			if center == 0 {
				center = lastPivot
			} else {
				// weighted calculation
				center = (center*2 + lastPivot) / 3
			}
		}

		centerLine[i] = center
	}

	return centerLine
}

// calculateATR คำนวณ Average True Range
func (ind *Indicators) calculateATR(ohlcv []OHLCV) []float64 {
	period := ind.atrPeriod
	atr := make([]float64, len(ohlcv))

	if len(ohlcv) < period {
		return atr
	}

	// คำนวณ True Range
	trueRanges := make([]float64, len(ohlcv))
	for i := 1; i < len(ohlcv); i++ {
		high := ohlcv[i].High
		low := ohlcv[i].Low
		prevClose := ohlcv[i-1].Close

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)

		trueRanges[i] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// คำนวณ ATR (Simple Moving Average ของ True Range)
	for i := period; i < len(ohlcv); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += trueRanges[j]
		}
		atr[i] = sum / float64(period)
	}

	return atr
}

// calculateSuperTrendBands คำนวณ SuperTrend Upper/Lower Bands
func (ind *Indicators) calculateSuperTrendBands(centerLine, atr []float64) ([]float64, []float64) {
	factor := ind.atrFactor
	upperBand := make([]float64, len(centerLine))
	lowerBand := make([]float64, len(centerLine))

	for i := 0; i < len(centerLine); i++ {
		if centerLine[i] != 0 && atr[i] != 0 {
			upperBand[i] = centerLine[i] - (factor * atr[i])
			lowerBand[i] = centerLine[i] + (factor * atr[i])
		}
	}

	return upperBand, lowerBand
}

// calculateSuperTrend คำนวณ SuperTrend และ Trend
func (ind *Indicators) calculateSuperTrend(ohlcv []OHLCV, upperBand, lowerBand []float64) ([]float64, []int) {
	length := len(ohlcv)
	superTrend := make([]float64, length)
	trend := make([]int, length)
	tUp := make([]float64, length)
	tDown := make([]float64, length)

	for i := 1; i < length; i++ {
		close := ohlcv[i].Close
		prevClose := ohlcv[i-1].Close

		// คำนวณ TUp และ TDown
		if prevClose > tUp[i-1] {
			tUp[i] = math.Max(upperBand[i], tUp[i-1])
		} else {
			tUp[i] = upperBand[i]
		}

		if prevClose < tDown[i-1] {
			tDown[i] = math.Min(lowerBand[i], tDown[i-1])
		} else {
			tDown[i] = lowerBand[i]
		}

		// คำนวณ Trend
		if close > tDown[i-1] {
			trend[i] = 1 // Bullish
		} else if close < tUp[i-1] {
			trend[i] = -1 // Bearish
		} else {
			if i > 0 {
				trend[i] = trend[i-1] // คงเดิม
			} else {
				trend[i] = 1 // default
			}
		}

		// คำนวณ SuperTrend
		if trend[i] == 1 {
			superTrend[i] = tUp[i]
		} else {
			superTrend[i] = tDown[i]
		}
	}

	return superTrend, trend
}

// calculateEMA คำนวณ Exponential Moving Average
func (ind *Indicators) calculateEMA(ohlcv []OHLCV, period int) []float64 {
	ema := make([]float64, len(ohlcv))

	if len(ohlcv) < period {
		return ema
	}

	// คำนวณ SMA สำหรับค่าเริ่มต้น
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += ohlcv[i].Close
	}
	ema[period-1] = sum / float64(period)

	// คำนวณ EMA
	multiplier := 2.0 / float64(period+1)
	for i := period; i < len(ohlcv); i++ {
		ema[i] = (ohlcv[i].Close * multiplier) + (ema[i-1] * (1 - multiplier))
	}

	return ema
}

// analyzeSignal วิเคราะห์ Signal และ Confidence
func (ind *Indicators) analyzeSignal(ohlcv []OHLCV, trend []int, ema100, superTrend []float64) (string, float64) {
	if len(ohlcv) < 2 {
		return "NEUTRAL", 0
	}

	currentIdx := len(ohlcv) - 1
	prevIdx := currentIdx - 1

	currentTrend := trend[currentIdx]
	prevTrend := trend[prevIdx]

	// ตรวจสอบการเปลี่ยนแปลง trend
	trendChanged := currentTrend != prevTrend
	confidence := 50.0

	// LONG Signal Conditions
	if currentTrend == 1 { // Bullish SuperTrend
		// ตรวจสอบว่าราคา rebound จาก EMA100
		if ind.isPriceReboundFromEMA(ohlcv, ema100, currentIdx) {
			confidence = 80.0
			if trendChanged {
				confidence = 90.0 // trend เพิ่งเปลี่ยน + rebound
			}
			return "LONG", confidence
		}

		// SuperTrend bullish แต่ยังไม่มี rebound signal
		confidence = 60.0
		if trendChanged {
			confidence = 75.0
		}
		return "LONG", confidence
	}

	// SHORT Signal Conditions
	if currentTrend == -1 { // Bearish SuperTrend
		// ตรวจสอบว่าราคา reject จาก EMA100
		if ind.isPriceRejectFromEMA(ohlcv, ema100, currentIdx) {
			confidence = 80.0
			if trendChanged {
				confidence = 90.0 // trend เพิ่งเปลี่ยน + reject
			}
			return "SHORT", confidence
		}

		// SuperTrend bearish แต่ยังไม่มี reject signal
		confidence = 60.0
		if trendChanged {
			confidence = 75.0
		}
		return "SHORT", confidence
	}

	return "NEUTRAL", confidence
}

// isPriceReboundFromEMA ตรวจสอบว่าราคา rebound จาก EMA100
func (ind *Indicators) isPriceReboundFromEMA(ohlcv []OHLCV, ema100 []float64, currentIdx int) bool {
	// ตรวจสอบ 2 แท่งล่าสุด
	for i := 0; i < 2 && currentIdx-i >= 0; i++ {
		idx := currentIdx - i
		candle := ohlcv[idx]
		emaValue := ema100[idx]

		// แท่งเทียนสีเขียว (Close > Open) และอยู่เหนือ EMA100
		if candle.Close > candle.Open && candle.Close > emaValue {
			return true
		}
	}

	return false
}

// isPriceRejectFromEMA ตรวจสอบว่าราคา reject จาก EMA100
func (ind *Indicators) isPriceRejectFromEMA(ohlcv []OHLCV, ema100 []float64, currentIdx int) bool {
	// ตรวจสอบ 2 แท่งล่าสุด
	for i := 0; i < 2 && currentIdx-i >= 0; i++ {
		idx := currentIdx - i
		candle := ohlcv[idx]
		emaValue := ema100[idx]

		// แท่งเทียนสีแดง (Close < Open) และอยู่ใต้ EMA100
		if candle.Close < candle.Open && candle.Close < emaValue {
			return true
		}
	}

	return false
}

// GetRSI ดึงข้อมูล RSI ของเหรียญ (คืนค่าย้อนหลัง timeframes ตาม period)
func (ind *Indicators) GetRSI(symbol, interval string, timeframes int) ([]float64, error) {
	// ต้องการข้อมูล OHLCV เพื่อคำนวณ RSI
	// ดึงข้อมูล candlestick ย้อนหลัง (timeframes + 14 สำหรับ RSI period)
	// requiredCandles := timeframes + 14  // TODO: ใช้เมื่อเชื่อมต่อ API จริง

	// TODO: ใช้ GateClient ดึงข้อมูล OHLCV จริง
	// สำหรับตอนนี้จำลองการคำนวณ RSI จากข้อมูลจริง

	rsiValues := make([]float64, timeframes)

	// จำลอง RSI ที่แปรผันตามเวลา (แทนค่าคงที่)
	baseRSI := 50.0
	for i := 0; i < timeframes; i++ {
		// สร้างค่า RSI ที่แปรผันตามตำแหน่ง (เพื่อดูการเปลี่ยนแปลง)
		variation := math.Sin(float64(i)*0.3) * 15.0 // ใช้ sine wave สำหรับการแปรผัน
		rsiValue := baseRSI + variation

		// จำกัดให้อยู่ในช่วง 0-100
		if rsiValue > 100 {
			rsiValue = 100
		} else if rsiValue < 0 {
			rsiValue = 0
		}

		rsiValues[timeframes-1-i] = rsiValue // เรียงจากเก่าไปใหม่
	}

	return rsiValues, nil
}

// calculateRiskReward คำนวณ Risk-Reward Ratio
func (ind *Indicators) calculateRiskReward(ohlcv []OHLCV, superTrend, atr []float64) float64 {
	if len(ohlcv) == 0 || len(superTrend) == 0 || len(atr) == 0 {
		return 1.0
	}

	currentIdx := len(ohlcv) - 1
	currentPrice := ohlcv[currentIdx].Close
	currentSuperTrend := superTrend[currentIdx]
	currentATR := atr[currentIdx]

	if currentATR == 0 {
		return 1.0
	}

	// คำนวณ Stop Loss และ Take Profit จาก SuperTrend และ ATR
	stopLoss := math.Abs(currentPrice - currentSuperTrend)
	takeProfit := currentATR * 2.0 // ใช้ 2x ATR เป็น target

	if stopLoss == 0 {
		return 3.0 // default high ratio
	}

	riskRewardRatio := takeProfit / stopLoss

	// จำกัด ratio ไว้ที่ 0.1 - 10.0
	if riskRewardRatio < 0.1 {
		riskRewardRatio = 0.1
	} else if riskRewardRatio > 10.0 {
		riskRewardRatio = 10.0
	}

	return riskRewardRatio
}
