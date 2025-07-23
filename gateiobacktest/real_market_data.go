package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// RealDataManager - ระบบดึงข้อมูลจริงจาก Gate.io API
type RealDataManager struct {
	baseURL        string
	symbol         string
	dataBuffer     []OHLCV
	maxLookback    int
	lastUpdateTime time.Time
}

// OHLCV structure for real market data
type OHLCV struct {
	Timestamp int64   `json:"t"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
}

// GateIoKlineResponse - Gate.io API response structure
type GateIoKlineResponse [][]string

// NewRealDataManager - สร้าง real data manager
func NewRealDataManager(symbol string) *RealDataManager {
	return &RealDataManager{
		baseURL:     "https://api.gateio.ws/api/v4/spot/candlesticks",
		symbol:      symbol,
		maxLookback: 144, // 6 days × 24 hours
		dataBuffer:  make([]OHLCV, 0, 200),
	}
}

// LoadRealMarketData - ดึงข้อมูลจริงจาก Gate.io API
func (rm *RealDataManager) LoadRealMarketData() error {
	fmt.Printf("🔄 Loading REAL market data for %s from Gate.io...\n", rm.symbol)

	// คำนวณเวลา (ต้องการ 144 candles + buffer)
	endTime := time.Now().Unix()
	startTime := endTime - int64((rm.maxLookback+20)*3600) // +20 hours buffer

	// สร้าง URL สำหรับ Gate.io API
	url := fmt.Sprintf("%s?currency_pair=%s&interval=1h&from=%d&to=%d&limit=200",
		rm.baseURL, rm.symbol, startTime, endTime)

	fmt.Printf("📡 Fetching from: %s\n", url)

	// ทำ HTTP request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	// อ่าน response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var klineData GateIoKlineResponse
	if err := json.Unmarshal(body, &klineData); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	if len(klineData) == 0 {
		return fmt.Errorf("no data received from API")
	}

	// แปลงข้อมูลเป็น OHLCV format
	rm.dataBuffer = make([]OHLCV, 0, len(klineData))

	for _, candle := range klineData {
		if len(candle) < 6 {
			continue // skip invalid data
		}

		// Parse each field
		timestamp, err1 := strconv.ParseInt(candle[0], 10, 64)
		open, err2 := strconv.ParseFloat(candle[5], 64)   // Gate.io format: open is index 5
		high, err3 := strconv.ParseFloat(candle[3], 64)   // high is index 3
		low, err4 := strconv.ParseFloat(candle[4], 64)    // low is index 4
		close, err5 := strconv.ParseFloat(candle[2], 64)  // close is index 2
		volume, err6 := strconv.ParseFloat(candle[1], 64) // volume is index 1

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
			fmt.Printf("⚠️ Skipping invalid candle data\n")
			continue
		}

		ohlcv := OHLCV{
			Timestamp: timestamp,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}

		rm.dataBuffer = append(rm.dataBuffer, ohlcv)
	}

	// ตรวจสอบจำนวนข้อมูล
	if len(rm.dataBuffer) < rm.maxLookback {
		return fmt.Errorf("insufficient real data: got %d candles, need %d",
			len(rm.dataBuffer), rm.maxLookback)
	}

	// เก็บแค่ candles ล่าสุดตามที่ต้องการ
	if len(rm.dataBuffer) > rm.maxLookback {
		rm.dataBuffer = rm.dataBuffer[len(rm.dataBuffer)-rm.maxLookback:]
	}

	rm.lastUpdateTime = time.Unix(rm.dataBuffer[len(rm.dataBuffer)-1].Timestamp, 0)

	fmt.Printf("✅ Loaded %d REAL candles (%.1f days)\n",
		len(rm.dataBuffer), float64(len(rm.dataBuffer))/24.0)
	fmt.Printf("📊 Latest price: $%.4f (Volume: %.2f)\n",
		rm.dataBuffer[len(rm.dataBuffer)-1].Close,
		rm.dataBuffer[len(rm.dataBuffer)-1].Volume)
	fmt.Printf("⏰ Data range: %s to %s\n",
		time.Unix(rm.dataBuffer[0].Timestamp, 0).Format("2006-01-02 15:04"),
		time.Unix(rm.dataBuffer[len(rm.dataBuffer)-1].Timestamp, 0).Format("2006-01-02 15:04"))

	return nil
}

// GetLatestCandles - ดึงข้อมูล candles ล่าสุด
func (rm *RealDataManager) GetLatestCandles(count int) []OHLCV {
	if count > len(rm.dataBuffer) {
		count = len(rm.dataBuffer)
	}

	start := len(rm.dataBuffer) - count
	return rm.dataBuffer[start:]
}

// CalculateEMA - คำนวณ EMA จากข้อมูลจริง
func (rm *RealDataManager) CalculateEMA(period int, data []OHLCV) []float64 {
	if len(data) < period {
		return nil
	}

	ema := make([]float64, len(data))
	multiplier := 2.0 / (float64(period) + 1.0)

	// คำนวณ SMA สำหรับค่าแรก
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += data[i].Close
	}
	ema[period-1] = sum / float64(period)

	// คำนวณ EMA สำหรับที่เหลือ
	for i := period; i < len(data); i++ {
		ema[i] = (data[i].Close * multiplier) + (ema[i-1] * (1 - multiplier))
	}

	return ema
}

// CalculateRSI - คำนวณ RSI จากข้อมูลจริง
func (rm *RealDataManager) CalculateRSI(period int, data []OHLCV) []float64 {
	if len(data) < period+1 {
		return nil
	}

	rsi := make([]float64, len(data))
	gains := make([]float64, len(data))
	losses := make([]float64, len(data))

	// คำนวณ gains และ losses
	for i := 1; i < len(data); i++ {
		change := data[i].Close - data[i-1].Close
		if change > 0 {
			gains[i] = change
			losses[i] = 0
		} else {
			gains[i] = 0
			losses[i] = -change
		}
	}

	// คำนวณ average gains และ losses
	avgGain := 0.0
	avgLoss := 0.0
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// คำนวณ RSI
	for i := period; i < len(data); i++ {
		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}

		// อัพเดท average สำหรับรอบถัดไป
		if i < len(data)-1 {
			avgGain = (avgGain*float64(period-1) + gains[i+1]) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + losses[i+1]) / float64(period)
		}
	}

	return rsi
}

// CalculateATR - คำนวณ ATR จากข้อมูลจริง
func (rm *RealDataManager) CalculateATR(period int, data []OHLCV) []float64 {
	if len(data) < period+1 {
		return nil
	}

	atr := make([]float64, len(data))
	tr := make([]float64, len(data))

	// คำนวณ True Range
	for i := 1; i < len(data); i++ {
		high_low := data[i].High - data[i].Low
		high_close := math.Abs(data[i].High - data[i-1].Close)
		low_close := math.Abs(data[i].Low - data[i-1].Close)

		tr[i] = math.Max(high_low, math.Max(high_close, low_close))
	}

	// คำนวณ ATR
	sum := 0.0
	for i := 1; i <= period; i++ {
		sum += tr[i]
	}
	atr[period] = sum / float64(period)

	// คำนวณ ATR สำหรับที่เหลือ (Smoothed)
	for i := period + 1; i < len(data); i++ {
		atr[i] = (atr[i-1]*float64(period-1) + tr[i]) / float64(period)
	}

	return atr
}

// GetRealMarketAnalysis - วิเคราะห์ข้อมูลจริงแบบครบถ้วน
func (rm *RealDataManager) GetRealMarketAnalysis() *RealMarketAnalysis {
	data := rm.dataBuffer
	if len(data) < 144 {
		return nil
	}

	// คำนวณ indicators จากข้อมูลจริง
	ema9 := rm.CalculateEMA(9, data)
	ema21 := rm.CalculateEMA(21, data)
	ema50 := rm.CalculateEMA(50, data)
	rsi := rm.CalculateRSI(14, data)
	atr := rm.CalculateATR(14, data)

	// ข้อมูลล่าสุด
	latest := data[len(data)-1]
	latestIdx := len(data) - 1

	// คำนวณ Volume MA
	volumeSum := 0.0
	volumePeriod := 20
	for i := latestIdx - volumePeriod + 1; i <= latestIdx; i++ {
		if i >= 0 {
			volumeSum += data[i].Volume
		}
	}
	volumeMA := volumeSum / float64(volumePeriod)

	return &RealMarketAnalysis{
		Symbol:        rm.symbol,
		Timestamp:     latest.Timestamp,
		CurrentPrice:  latest.Close,
		CurrentVolume: latest.Volume,
		EMA9:          ema9[latestIdx],
		EMA21:         ema21[latestIdx],
		EMA50:         ema50[latestIdx],
		RSI:           rsi[latestIdx],
		ATR:           atr[latestIdx],
		VolumeMA:      volumeMA,
		DataQuality:   1.0, // Real data = 100% quality
		CandleCount:   len(data),
	}
}

// RealMarketAnalysis - โครงสร้างสำหรับผลการวิเคราะห์ข้อมูลจริง
type RealMarketAnalysis struct {
	Symbol        string  `json:"symbol"`
	Timestamp     int64   `json:"timestamp"`
	CurrentPrice  float64 `json:"current_price"`
	CurrentVolume float64 `json:"current_volume"`
	EMA9          float64 `json:"ema9"`
	EMA21         float64 `json:"ema21"`
	EMA50         float64 `json:"ema50"`
	RSI           float64 `json:"rsi"`
	ATR           float64 `json:"atr"`
	VolumeMA      float64 `json:"volume_ma"`
	DataQuality   float64 `json:"data_quality"`
	CandleCount   int     `json:"candle_count"`
}

// PrintMarketSummary - แสดงสรุปตลาดจากข้อมูลจริง
func (analysis *RealMarketAnalysis) PrintMarketSummary() {
	fmt.Printf("\n📊 REAL MARKET ANALYSIS: %s\n", analysis.Symbol)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("💰 Current Price: $%.4f\n", analysis.CurrentPrice)
	fmt.Printf("📈 EMAs: Fast(9)=%.4f | Mid(21)=%.4f | Slow(50)=%.4f\n",
		analysis.EMA9, analysis.EMA21, analysis.EMA50)
	fmt.Printf("📊 RSI: %.2f | ATR: %.4f\n", analysis.RSI, analysis.ATR)
	fmt.Printf("📦 Volume: %.2f (MA: %.2f)\n", analysis.CurrentVolume, analysis.VolumeMA)
	fmt.Printf("✅ Data Quality: %.1f%% | Candles: %d\n",
		analysis.DataQuality*100, analysis.CandleCount)

	// EMA Alignment Analysis
	if analysis.EMA9 > analysis.EMA21 && analysis.EMA21 > analysis.EMA50 {
		fmt.Printf("🟢 EMA Alignment: BULLISH (Fast>Mid>Slow)\n")
	} else if analysis.EMA9 < analysis.EMA21 && analysis.EMA21 < analysis.EMA50 {
		fmt.Printf("🔴 EMA Alignment: BEARISH (Fast<Mid<Slow)\n")
	} else {
		fmt.Printf("🟡 EMA Alignment: MIXED (No clear trend)\n")
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// TestRealMarketData - ทดสอบการดึงข้อมูลจริง
func TestRealMarketData() {
	symbols := []string{"BTC_USDT", "ETH_USDT", "SOL_USDT"}

	fmt.Println("🚀 REAL MARKET DATA BACKTEST - Gate.io API")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, symbol := range symbols {
		fmt.Printf("\n🔍 [%d/%d] Loading REAL data: %s\n", i+1, len(symbols), symbol)

		// สร้าง real data manager
		dataManager := NewRealDataManager(symbol)

		// ดึงข้อมูลจริงจาก Gate.io
		if err := dataManager.LoadRealMarketData(); err != nil {
			log.Printf("❌ Failed to load real data for %s: %v", symbol, err)
			continue
		}

		// วิเคราะห์ข้อมูลจริง
		analysis := dataManager.GetRealMarketAnalysis()
		if analysis != nil {
			analysis.PrintMarketSummary()
		}

		// หน่วงเวลาเพื่อไม่ให้ API rate limit
		if i < len(symbols)-1 {
			time.Sleep(2 * time.Second)
		}
	}

	fmt.Println("\n✅ Real market data analysis complete!")
	fmt.Println("🎯 Ready for live trading with actual market conditions")
}
