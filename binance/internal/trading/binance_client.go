package trading

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"binance-trading-bot/internal/binance"
)

// BinanceClient wrapper สำหรับ Binance API
type BinanceClient struct {
	client *binance.Client
}

// NewBinanceClient สร้าง BinanceClient ใหม่
func NewBinanceClient(client *binance.Client) *BinanceClient {
	return &BinanceClient{
		client: client,
	}
}

// TestConnection ทดสอบการเชื่อมต่อ
func (bc *BinanceClient) TestConnection() bool {
	balances, err := bc.client.GetBalance()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเชื่อมต่อ Binance API ได้: %v\n", err)
		return false
	}

	for _, balance := range balances {
		if balance.Asset == "USDT" {
			fmt.Printf("✅ เชื่อมต่อ Binance สำเร็จ - Balance: %s USDT\n", balance.AvailableBalance)
			break
		}
	}
	return true
}

// GetBalance ดึง USDT balance
func (bc *BinanceClient) GetBalance() (string, error) {
	balances, err := bc.client.GetBalance()
	if err != nil {
		return "0", err
	}

	for _, balance := range balances {
		if balance.Asset == "USDT" {
			return balance.AvailableBalance, nil
		}
	}

	return "0", fmt.Errorf("ไม่พบ USDT balance")
}

// GetFuturesContracts ดึงรายชื่อ symbols ทั้งหมดที่ลงท้ายด้วย USDT จาก Binance Futures
func (bc *BinanceClient) GetFuturesContracts() ([]string, error) {
	fmt.Printf("📋 กำลังดึงรายชื่อ symbols จาก Binance Futures API...\n")

	// ดึงข้อมูล exchange info
	exchangeInfo, err := bc.client.GetExchangeInfo()
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดึง exchange info ได้: %v\n", err)
		// ถ้าดึงไม่ได้ ใช้ symbols ที่นิยมแทน
		return bc.getFallbackSymbols(), nil
	}

	var usdtSymbols []string
	activeCount := 0
	totalCount := 0

	// กรองเฉพาะ symbols ที่ลงท้ายด้วย USDT และ status = TRADING
	for _, symbol := range exchangeInfo.Symbols {
		totalCount++
		if strings.HasSuffix(symbol.Symbol, "USDT") && symbol.Status == "TRADING" {
			usdtSymbols = append(usdtSymbols, symbol.Symbol)
			activeCount++
		}
	}

	fmt.Printf("📊 พบ %d symbols ทั้งหมด\n", totalCount)
	fmt.Printf("✅ กรองได้ %d symbols ที่ลงท้ายด้วย USDT และ active\n", activeCount)

	// สับไพ่
	for i := len(usdtSymbols) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		usdtSymbols[i], usdtSymbols[j] = usdtSymbols[j], usdtSymbols[i]
	}

	// แสดงรายชื่อ symbols (แค่ 10 ตัวแรกเพื่อไม่ให้ยาวเกินไป)
	fmt.Printf("📋 ตัวอย่าง symbols ที่จะวิเคราะห์: ")
	for i, symbol := range usdtSymbols {
		if i >= 10 {
			fmt.Printf("... และอีก %d symbols", len(usdtSymbols)-10)
			break
		}
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", symbol)
	}
	fmt.Println()

	return usdtSymbols, nil
}

// getFallbackSymbols ส่งคืน symbols ที่นิยมในกรณี API ล้มเหลว
func (bc *BinanceClient) getFallbackSymbols() []string {
	symbols := []string{
		"BTCUSDT", "ETHUSDT", "ADAUSDT", "DOTUSDT", "LINKUSDT",
		"LTCUSDT", "BCHUSDT", "XLMUSDT", "EOSUSDT", "TRXUSDT",
		"XRPUSDT", "BNBUSDT", "AVAXUSDT", "SOLUSDT", "MATICUSDT",
		"ATOMUSDT", "VETUSDT", "FILUSDT", "UNIUSDT", "AAVEUSDT",
	}
	fmt.Printf("⚠️ ใช้ %d symbols สำรองแทน\n", len(symbols))
	return symbols
}

// GetOpenPositions ดึงรายการ positions ที่เปิดอยู่
func (bc *BinanceClient) GetOpenPositions() ([]binance.Position, error) {
	return bc.client.GetPositions()
}

// CreateMarketOrder สร้าง market order
func (bc *BinanceClient) CreateMarketOrder(symbol, side string, quantity float64) (*binance.OrderResponse, error) {
	order := binance.OrderRequest{
		Symbol:   symbol,
		Side:     side,
		Type:     "MARKET",
		Quantity: fmt.Sprintf("%.8f", quantity),
	}

	return bc.client.CreateOrder(order)
}

// CreateLimitOrder สร้าง limit order
func (bc *BinanceClient) CreateLimitOrder(symbol, side string, quantity, price float64) (*binance.OrderResponse, error) {
	order := binance.OrderRequest{
		Symbol:      symbol,
		Side:        side,
		Type:        "LIMIT",
		Quantity:    fmt.Sprintf("%.8f", quantity),
		Price:       fmt.Sprintf("%.8f", price),
		TimeInForce: "GTC",
	}

	return bc.client.CreateOrder(order)
}

// ClosePosition ปิด position
func (bc *BinanceClient) ClosePosition(symbol string) (*binance.OrderResponse, error) {
	return bc.client.ClosePosition(symbol, "BOTH")
}

// GetCandlesticks ดึงข้อมูล candlestick
func (bc *BinanceClient) GetCandlesticks(symbol, interval string, limit int) ([]binance.Candlestick, error) {
	return bc.client.GetCandlesticks(symbol, interval, limit)
}

// ParseFloat helper function สำหรับแปลง string เป็น float64
func (bc *BinanceClient) ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

// FormatFloat helper function สำหรับแปลง float64 เป็น string
func (bc *BinanceClient) FormatFloat(f float64, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", f)
}

// CalculateQuantity คำนวณ quantity สำหรับ order
func (bc *BinanceClient) CalculateQuantity(symbol string, usdtAmount float64, price float64) float64 {
	// สำหรับ Binance Futures quantity = USDT amount / price
	quantity := usdtAmount / price

	// ปรับให้เหมาะกับ step size (ปกติ Binance ใช้ precision 3 ตำแหน่ง)
	return quantity
}

// GetSymbolInfo ดึงข้อมูล symbol (สำหรับอนาคตถ้าต้องการ)
func (bc *BinanceClient) GetSymbolInfo(symbol string) map[string]interface{} {
	// สำหรับตอนนี้ return ข้อมูลพื้นฐาน
	return map[string]interface{}{
		"symbol":     symbol,
		"status":     "TRADING",
		"baseAsset":  strings.Replace(symbol, "USDT", "", 1),
		"quoteAsset": "USDT",
		"minQty":     "0.001",
		"maxQty":     "1000000",
		"stepSize":   "0.001",
	}
}

// AdjustQuantityPrecision ปรับ quantity precision
func (bc *BinanceClient) AdjustQuantityPrecision(symbol string, quantity float64) (string, error) {
	return bc.client.AdjustQuantityPrecision(symbol, quantity)
}

// SetLeverage ตั้งค่า leverage สำหรับ symbol
func (bc *BinanceClient) SetLeverage(symbol string, leverage int) error {
	return bc.client.SetLeverage(symbol, leverage)
}

// SetMarginType ตั้งค่า margin type สำหรับ symbol
func (bc *BinanceClient) SetMarginType(symbol string, marginType string) error {
	return bc.client.SetMarginType(symbol, marginType)
}
