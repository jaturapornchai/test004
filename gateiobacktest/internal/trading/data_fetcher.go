package trading

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// CoinGeckoCandle ข้อมูล OHLCV จาก CoinGecko API
type CoinGeckoCandle struct {
	Timestamp float64
	Open      float64
	High      float64
	Low       float64
	Close     float64
}

// DataFetcher สำหรับดึงข้อมูลราคาจาก API
type DataFetcher struct {
	symbol   string
	days     int
	interval string
}

// NewDataFetcher สร้าง data fetcher ใหม่
func NewDataFetcher(symbol string, days int, interval string) *DataFetcher {
	return &DataFetcher{
		symbol:   symbol,
		days:     days,
		interval: interval,
	}
}

// FetchOrLoadData ดึงข้อมูลจาก API หรือโหลดจากไฟล์
func (df *DataFetcher) FetchOrLoadData() ([]OHLCV, error) {
	filename := fmt.Sprintf("data_%s_%dm_%dd.json", df.symbol, 15, df.days)

	// ตรวจสอบว่ามีไฟล์อยู่แล้วหรือไม่
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("📄 โหลดข้อมูลจากไฟล์: %s\n", filename)
		return df.loadFromFile(filename)
	}

	// ดึงข้อมูลจาก API
	fmt.Printf("🌐 ดึงข้อมูลจาก CoinGecko API: %s (%d วัน)\n", df.symbol, df.days)
	data, err := df.fetchFromAPI()
	if err != nil {
		return nil, err
	}

	// บันทึกลงไฟล์
	err = df.saveToFile(filename, data)
	if err != nil {
		fmt.Printf("⚠️ ไม่สามารถบันทึกไฟล์ได้: %v\n", err)
	}

	return data, nil
}

// fetchFromAPI ดึงข้อมูลจาก CoinGecko API
func (df *DataFetcher) fetchFromAPI() ([]OHLCV, error) {
	coinID := df.getCoinGeckoID()
	if coinID == "" {
		return nil, fmt.Errorf("ไม่สนับสนุนเหรียญ: %s", df.symbol)
	}

	// CoinGecko API URL
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/ohlc?vs_currency=usd&days=%d",
		coinID, df.days)

	fmt.Printf("🔗 API URL: %s\n", url)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถเชื่อมต่อ API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถอ่านข้อมูล: %v", err)
	}

	var candles [][]float64
	err = json.Unmarshal(body, &candles)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse JSON: %v", err)
	}

	return df.convertToOHLCV(candles), nil
}

// getCoinGeckoID แปลง symbol เป็น CoinGecko ID
func (df *DataFetcher) getCoinGeckoID() string {
	switch df.symbol {
	case "SOL_USDT", "SOL":
		return "solana"
	case "BTC_USDT", "BTC":
		return "bitcoin"
	case "ETH_USDT", "ETH":
		return "ethereum"
	default:
		return ""
	}
}

// convertToOHLCV แปลงข้อมูลจาก CoinGecko เป็น OHLCV
func (df *DataFetcher) convertToOHLCV(candles [][]float64) []OHLCV {
	var result []OHLCV

	for _, candle := range candles {
		if len(candle) >= 5 {
			// CoinGecko ให้ข้อมูลแบบ hourly หรือ daily
			// เราจะจำลอง 15m data โดยการแบ่งข้อมูล
			timestamp := int64(candle[0] / 1000) // Convert from milliseconds
			open := candle[1]
			high := candle[2]
			low := candle[3]
			close := candle[4]

			// สร้าง 15m candles โดยการแบ่งข้อมูล daily
			df.generate15mCandles(&result, timestamp, open, high, low, close)
		}
	}

	fmt.Printf("📊 แปลงข้อมูลเป็น %d แท่งเทียน 15m\n", len(result))
	return result
}

// generate15mCandles สร้าง 15m candles จากข้อมูล daily
func (df *DataFetcher) generate15mCandles(result *[]OHLCV, timestamp int64, open, high, low, close float64) {
	// สร้าง 96 candles ต่อวัน (24h * 4 candles per hour)
	baseTime := timestamp

	for i := 0; i < 96; i++ {
		candleTime := baseTime + int64(i*15*60) // เพิ่ม 15 นาที

		// จำลองการเคลื่อนไหวราคาภายในวัน
		progress := float64(i) / 96.0

		// สร้างราคาที่เป็นธรรมชาติ
		currentPrice := open + (close-open)*progress

		// เพิ่มความผันผวนแบบสุ่ม
		volatility := (high - low) * 0.1 // 10% ของ daily range
		variation := (float64(i%7) - 3) / 10.0 * volatility

		candleOpen := currentPrice + variation
		candleClose := candleOpen + (float64(i%5)-2)/20.0*volatility
		candleHigh := candleOpen + (high-low)*0.05
		candleLow := candleOpen - (high-low)*0.05

		// ปรับให้อยู่ในช่วงที่เหมาะสม
		if candleHigh > high {
			candleHigh = high
		}
		if candleLow < low {
			candleLow = low
		}

		// Volume จำลอง
		baseVolume := 500000.0
		volumeVariation := float64(i%20) * 50000.0
		volume := baseVolume + volumeVariation

		*result = append(*result, OHLCV{
			Timestamp: candleTime,
			Open:      candleOpen,
			High:      candleHigh,
			Low:       candleLow,
			Close:     candleClose,
			Volume:    volume,
		})
	}
}

// saveToFile บันทึกข้อมูลลงไฟล์
func (df *DataFetcher) saveToFile(filename string, data []OHLCV) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("💾 บันทึกข้อมูลลงไฟล์: %s (%d candles)\n", filename, len(data))
	return nil
}

// loadFromFile โหลดข้อมูลจากไฟล์
func (df *DataFetcher) loadFromFile(filename string) ([]OHLCV, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var ohlcvData []OHLCV
	err = json.Unmarshal(data, &ohlcvData)
	if err != nil {
		return nil, err
	}

	fmt.Printf("📊 โหลดข้อมูล %d แท่งเทียนจากไฟล์\n", len(ohlcvData))
	return ohlcvData, nil
}

// GetSymbolInfo แสดงข้อมูลสรุปของเหรียญ
func (df *DataFetcher) GetSymbolInfo(data []OHLCV) {
	if len(data) == 0 {
		return
	}

	firstCandle := data[0]
	lastCandle := data[len(data)-1]

	startTime := time.Unix(firstCandle.Timestamp, 0)
	endTime := time.Unix(lastCandle.Timestamp, 0)

	priceChange := lastCandle.Close - firstCandle.Close
	priceChangePct := (priceChange / firstCandle.Close) * 100

	fmt.Printf("📈 %s: $%.2f -> $%.2f (%.2f%%) | %s ถึง %s\n",
		df.symbol,
		firstCandle.Close,
		lastCandle.Close,
		priceChangePct,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"))
}
