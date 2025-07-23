package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	httpClient *http.Client
}

type Position struct {
	Symbol           string `json:"symbol"`
	PositionAmt      string `json:"positionAmt"`
	EntryPrice       string `json:"entryPrice"`
	MarkPrice        string `json:"markPrice"`
	UnrealizedProfit string `json:"unRealizedProfit"`
	Leverage         string `json:"leverage"`
	PositionSide     string `json:"positionSide"`
	MarginType       string `json:"marginType"`
}

type Balance struct {
	Asset                  string `json:"asset"`
	WalletBalance          string `json:"walletBalance"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	MarginBalance          string `json:"marginBalance"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	AvailableBalance       string `json:"availableBalance"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
}

type AccountInfo struct {
	Assets    []Balance  `json:"assets"`
	Positions []Position `json:"positions"`
}

type Candlestick struct {
	OpenTime      int64  `json:"openTime"`
	Open          string `json:"open"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Close         string `json:"close"`
	Volume        string `json:"volume"`
	CloseTime     int64  `json:"closeTime"`
	QuoteVolume   string `json:"quoteAssetVolume"`
	TradeCount    int    `json:"count"`
	TakerBuyBase  string `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuote string `json:"takerBuyQuoteAssetVolume"`
}

type OrderRequest struct {
	Symbol           string `json:"symbol"`
	Side             string `json:"side"`
	Type             string `json:"type"`
	Quantity         string `json:"quantity,omitempty"`
	Price            string `json:"price,omitempty"`
	TimeInForce      string `json:"timeInForce,omitempty"`
	ReduceOnly       bool   `json:"reduceOnly,omitempty"`
	NewClientOrderId string `json:"newClientOrderId,omitempty"`
}

type OrderResponse struct {
	OrderId       int64  `json:"orderId"`
	Symbol        string `json:"symbol"`
	Status        string `json:"status"`
	ClientOrderId string `json:"clientOrderId"`
	Price         string `json:"price"`
	AvgPrice      string `json:"avgPrice"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	CumQty        string `json:"cumQty"`
	CumQuote      string `json:"cumQuote"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	ReduceOnly    bool   `json:"reduceOnly"`
	ClosePosition bool   `json:"closePosition"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	StopPrice     string `json:"stopPrice"`
	WorkingType   string `json:"workingType"`
	PriceProtect  bool   `json:"priceProtect"`
	OrigType      string `json:"origType"`
	UpdateTime    int64  `json:"updateTime"`
}

func NewClient(apiKey, apiSecret, baseURL string) *Client {
	return &Client{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) sign(queryString string) string {
	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) request(method, endpoint string, params map[string]string) ([]byte, error) {
	fullURL := c.BaseURL + endpoint

	// สร้าง query string
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	// เพิ่ม timestamp
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	values.Set("timestamp", timestamp)

	queryString := values.Encode()
	signature := c.sign(queryString)
	queryString += "&signature=" + signature

	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequest(method, fullURL+"?"+queryString, nil)
	} else {
		req, err = http.NewRequest(method, fullURL, strings.NewReader(queryString))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง request ได้: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", c.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถส่ง request ได้: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถอ่าน response ได้: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetPositions ดึงรายการ positions ที่เปิดอยู่ทั้งหมด
func (c *Client) GetPositions() ([]Position, error) {
	params := map[string]string{}

	respBody, err := c.request("GET", "/fapi/v2/positionRisk", params)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึง positions ได้: %v", err)
	}

	var positions []Position
	err = json.Unmarshal(respBody, &positions)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse positions data ได้: %v", err)
	}

	// กรองเฉพาะ positions ที่เปิดอยู่ (positionAmt != "0")
	var openPositions []Position
	for _, pos := range positions {
		// แปลงเป็น float เพื่อตรวจสอบว่าเป็น 0 หรือไม่
		posAmt, err := strconv.ParseFloat(pos.PositionAmt, 64)
		if err == nil && math.Abs(posAmt) > 0.0001 { // ใช้ absolute value และ threshold เล็กน้อย
			openPositions = append(openPositions, pos)
		}
	}

	return openPositions, nil
}

// GetBalance ดึงยอด balance ของ futures account
func (c *Client) GetBalance() ([]Balance, error) {
	params := map[string]string{}

	respBody, err := c.request("GET", "/fapi/v2/account", params)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึง account info ได้: %v", err)
	}

	var accountInfo AccountInfo
	err = json.Unmarshal(respBody, &accountInfo)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse account data ได้: %v", err)
	}

	return accountInfo.Assets, nil
}

// GetCandlesticks ดึงข้อมูล candlestick สำหรับการวิเคราะห์เทคนิค
func (c *Client) GetCandlesticks(symbol, interval string, limit int) ([]Candlestick, error) {
	params := map[string]string{
		"symbol":   symbol,
		"interval": interval,
		"limit":    strconv.Itoa(limit),
	}

	respBody, err := c.request("GET", "/fapi/v1/klines", params)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึง candlesticks ได้: %v", err)
	}

	var rawData [][]interface{}
	err = json.Unmarshal(respBody, &rawData)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse candlesticks data ได้: %v", err)
	}

	var candlesticks []Candlestick
	for _, data := range rawData {
		if len(data) >= 12 {
			openTime, _ := data[0].(float64)
			open, _ := data[1].(string)
			high, _ := data[2].(string)
			low, _ := data[3].(string)
			close, _ := data[4].(string)
			volume, _ := data[5].(string)
			closeTime, _ := data[6].(float64)
			quoteVolume, _ := data[7].(string)
			tradeCount, _ := data[8].(float64)
			takerBuyBase, _ := data[9].(string)
			takerBuyQuote, _ := data[10].(string)

			candlestick := Candlestick{
				OpenTime:      int64(openTime),
				Open:          open,
				High:          high,
				Low:           low,
				Close:         close,
				Volume:        volume,
				CloseTime:     int64(closeTime),
				QuoteVolume:   quoteVolume,
				TradeCount:    int(tradeCount),
				TakerBuyBase:  takerBuyBase,
				TakerBuyQuote: takerBuyQuote,
			}
			candlesticks = append(candlesticks, candlestick)
		}
	}

	return candlesticks, nil
}

// CreateOrder สร้าง order ใหม่
func (c *Client) CreateOrder(order OrderRequest) (*OrderResponse, error) {
	params := map[string]string{
		"symbol": order.Symbol,
		"side":   order.Side,
		"type":   order.Type,
	}

	if order.Quantity != "" {
		params["quantity"] = order.Quantity
	}
	if order.Price != "" {
		params["price"] = order.Price
	}
	if order.TimeInForce != "" {
		params["timeInForce"] = order.TimeInForce
	}
	if order.ReduceOnly {
		params["reduceOnly"] = "true"
	}
	if order.NewClientOrderId != "" {
		params["newClientOrderId"] = order.NewClientOrderId
	}

	respBody, err := c.request("POST", "/fapi/v1/order", params)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง order ได้: %v", err)
	}

	var orderResponse OrderResponse
	err = json.Unmarshal(respBody, &orderResponse)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse order response ได้: %v", err)
	}

	return &orderResponse, nil
}

// CancelOrder ยกเลิก order
func (c *Client) CancelOrder(symbol string, orderId int64) (*OrderResponse, error) {
	params := map[string]string{
		"symbol":  symbol,
		"orderId": strconv.FormatInt(orderId, 10),
	}

	respBody, err := c.request("DELETE", "/fapi/v1/order", params)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถยกเลิก order ได้: %v", err)
	}

	var orderResponse OrderResponse
	err = json.Unmarshal(respBody, &orderResponse)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse cancel order response ได้: %v", err)
	}

	return &orderResponse, nil
}

// ClosePosition ปิด position ทั้งหมดของ symbol
func (c *Client) ClosePosition(symbol, positionSide string) (*OrderResponse, error) {
	// ดึงข้อมูล position ปัจจุบัน
	positions, err := c.GetPositions()
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูล position ได้: %v", err)
	}

	var targetPosition *Position
	for _, pos := range positions {
		if pos.Symbol == symbol && pos.PositionSide == positionSide {
			targetPosition = &pos
			break
		}
	}

	if targetPosition == nil {
		return nil, fmt.Errorf("ไม่พบ position สำหรับ %s", symbol)
	}

	// สร้าง order เพื่อปิด position
	side := "SELL"
	if strings.Contains(targetPosition.PositionAmt, "-") {
		side = "BUY"
	}

	order := OrderRequest{
		Symbol:     symbol,
		Side:       side,
		Type:       "MARKET",
		Quantity:   strings.Replace(targetPosition.PositionAmt, "-", "", 1), // เอา - ออก
		ReduceOnly: true,
	}

	return c.CreateOrder(order)
}

// ExchangeInfo structure for Binance exchange information
type ExchangeInfo struct {
	Symbols []SymbolInfo `json:"symbols"`
}

type SymbolInfo struct {
	Symbol     string   `json:"symbol"`
	Status     string   `json:"status"`
	BaseAsset  string   `json:"baseAsset"`
	QuoteAsset string   `json:"quoteAsset"`
	Filters    []Filter `json:"filters"`
}

type Filter struct {
	FilterType  string `json:"filterType"`
	MinQty      string `json:"minQty,omitempty"`
	MaxQty      string `json:"maxQty,omitempty"`
	StepSize    string `json:"stepSize,omitempty"`
	MinPrice    string `json:"minPrice,omitempty"`
	MaxPrice    string `json:"maxPrice,omitempty"`
	TickSize    string `json:"tickSize,omitempty"`
	MinNotional string `json:"minNotional,omitempty"`
}

// GetExchangeInfo ดึงข้อมูล exchange info
func (c *Client) GetExchangeInfo() (*ExchangeInfo, error) {
	params := map[string]string{}

	respBody, err := c.request("GET", "/fapi/v1/exchangeInfo", params)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึง exchange info ได้: %v", err)
	}

	var exchangeInfo ExchangeInfo
	err = json.Unmarshal(respBody, &exchangeInfo)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse exchange info ได้: %v", err)
	}

	return &exchangeInfo, nil
}

// AdjustQuantityPrecision ปรับ quantity ให้เหมาะสมกับ precision ของ symbol
func (c *Client) AdjustQuantityPrecision(symbol string, quantity float64) (string, error) {
	exchangeInfo, err := c.GetExchangeInfo()
	if err != nil {
		return "", err
	}

	// หา symbol info
	var symbolInfo *SymbolInfo
	for _, s := range exchangeInfo.Symbols {
		if s.Symbol == symbol {
			symbolInfo = &s
			break
		}
	}

	if symbolInfo == nil {
		// ถ้าไม่เจอ symbol ให้ใช้ default precision
		return fmt.Sprintf("%.3f", quantity), nil
	}

	// หา LOT_SIZE filter
	var stepSize float64 = 0.001 // default
	for _, filter := range symbolInfo.Filters {
		if filter.FilterType == "LOT_SIZE" && filter.StepSize != "" {
			stepSize, _ = strconv.ParseFloat(filter.StepSize, 64)
			break
		}
	}

	// ปรับ quantity ให้เป็นไปตาม step size
	adjustedQty := math.Floor(quantity/stepSize) * stepSize

	// นับจำนวนทศนิยมจาก step size
	stepStr := fmt.Sprintf("%.10f", stepSize)
	stepStr = strings.TrimRight(stepStr, "0")
	stepStr = strings.TrimRight(stepStr, ".")

	decimalPlaces := 0
	if dotIndex := strings.Index(stepStr, "."); dotIndex != -1 {
		decimalPlaces = len(stepStr) - dotIndex - 1
	}

	return fmt.Sprintf("%."+fmt.Sprintf("%d", decimalPlaces)+"f", adjustedQty), nil
}

// SetLeverage ตั้งค่า leverage สำหรับ symbol
func (c *Client) SetLeverage(symbol string, leverage int) error {
	params := map[string]string{
		"symbol":   symbol,
		"leverage": strconv.Itoa(leverage),
	}

	_, err := c.request("POST", "/fapi/v1/leverage", params)
	if err != nil {
		return fmt.Errorf("ไม่สามารถตั้งค่า leverage ได้: %v", err)
	}

	return nil
}

// SetMarginType ตั้งค่า margin type สำหรับ symbol
func (c *Client) SetMarginType(symbol string, marginType string) error {
	params := map[string]string{
		"symbol":     symbol,
		"marginType": marginType,
	}

	_, err := c.request("POST", "/fapi/v1/marginType", params)
	if err != nil {
		return fmt.Errorf("ไม่สามารถตั้งค่า margin type ได้: %v", err)
	}

	return nil
}
