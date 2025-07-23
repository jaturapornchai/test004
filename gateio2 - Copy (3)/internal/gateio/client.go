package gateio

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Contract      string  `json:"contract"`
	Size          float64 `json:"size"`
	Mode          string  `json:"mode"`
	Value         string  `json:"value"`
	EntryPrice    string  `json:"entry_price"`
	MarkPrice     string  `json:"mark_price"`
	UnrealizedPnl string  `json:"unrealized_pnl"`
	RealizedPnl   string  `json:"realized_pnl"`
	Side          string  `json:"side"`
	Leverage      string  `json:"leverage"`
}

type Balance struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Total     string `json:"total"`
}

type Candlestick struct {
	Timestamp int64   `json:"t"`
	Volume    float64 `json:"v"`
	Close     float64 `json:"c"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Open      float64 `json:"o"`
}

type Contract struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	QuantoMultiplier string `json:"quanto_multiplier"`
	LeverageMin      string `json:"leverage_min"`
	LeverageMax      string `json:"leverage_max"`
	MaintenanceRate  string `json:"maintenance_rate"`
	MarkType         string `json:"mark_type"`
	LastPrice        string `json:"last_price"`
	MarkPrice        string `json:"mark_price"`
	IndexPrice       string `json:"index_price"`
	FundingRate      string `json:"funding_rate"`
	TakerFeeRate     string `json:"taker_fee_rate"`
	MakerFeeRate     string `json:"maker_fee_rate"`
	OrderSizeMin     string `json:"order_size_min"`
	OrderSizeMax     string `json:"order_size_max"`
	OrderPriceRound  string `json:"order_price_round"`
	MarkPriceRound   string `json:"mark_price_round"`
	TradeSize        string `json:"trade_size"`
	PositionSize     string `json:"position_size"`
	ConfigChangeTime int64  `json:"config_change_time"`
	InDelisting      bool   `json:"in_delisting"`
	OrdersLimit      int    `json:"orders_limit"`
}

type OrderRequest struct {
	Contract      string `json:"contract"`
	Size          int64  `json:"size"`
	Price         string `json:"price,omitempty"`
	TIF           string `json:"tif,omitempty"`
	Text          string `json:"text,omitempty"`
	ReduceOnly    bool   `json:"reduce_only,omitempty"`
	ClosePosition bool   `json:"close_position,omitempty"`
}

type OrderResponse struct {
	ID           int64  `json:"id"`
	User         int    `json:"user"`
	CreateTime   int64  `json:"create_time"`
	FinishTime   int64  `json:"finish_time"`
	FinishAs     string `json:"finish_as"`
	Status       string `json:"status"`
	Contract     string `json:"contract"`
	Size         int64  `json:"size"`
	Price        string `json:"price"`
	FillPrice    string `json:"fill_price"`
	TIF          string `json:"tif"`
	Left         int64  `json:"left"`
	Text         string `json:"text"`
	TkfmFee      string `json:"tkfm_fee"`
	ReduceOnly   bool   `json:"reduce_only"`
	IsReduceOnly bool   `json:"is_reduce_only"`
	IsClose      bool   `json:"is_close"`
	IsLiq        bool   `json:"is_liq"`
	STP          string `json:"stp"`
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

func (c *Client) sign(method, uri, body string, timestamp int64) string {
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n%d", method, uri, "", body, timestamp)
	h := hmac.New(sha512.New, []byte(c.APISecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) request(method, endpoint string, body []byte) ([]byte, error) {
	url := c.BaseURL + endpoint
	timestamp := time.Now().Unix()

	bodyStr := ""
	if body != nil {
		bodyStr = string(body)
	}

	signature := c.sign(method, endpoint, bodyStr, timestamp)

	req, err := http.NewRequest(method, url, strings.NewReader(bodyStr))
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง request ได้: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", c.APIKey)
	req.Header.Set("Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("SIGN", signature)

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
	respBody, err := c.request("GET", "/api/v4/futures/usdt/positions", nil)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึง positions ได้: %v", err)
	}

	var positions []Position
	err = json.Unmarshal(respBody, &positions)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse positions data ได้: %v", err)
	}

	// กรองเฉพาะ positions ที่เปิดอยู่ (size != 0)
	var openPositions []Position
	for _, pos := range positions {
		if pos.Size != 0 {
			openPositions = append(openPositions, pos)
		}
	}

	return openPositions, nil
}

// GetBalance ดึงยอด balance ของ futures account
func (c *Client) GetBalance() (Balance, error) {
	respBody, err := c.request("GET", "/api/v4/futures/usdt/accounts", nil)
	if err != nil {
		return Balance{}, fmt.Errorf("ไม่สามารถดึง balance ได้: %v", err)
	}

	var balance Balance
	err = json.Unmarshal(respBody, &balance)
	if err != nil {
		return Balance{}, fmt.Errorf("ไม่สามารถ parse balance data ได้: %v", err)
	}

	return balance, nil
}

// GetCandlesticks ดึงข้อมูล OHLCV candlesticks
func (c *Client) GetCandlesticks(contract string, interval string, limit int) ([]Candlestick, error) {
	endpoint := fmt.Sprintf("/api/v4/futures/usdt/candlesticks?contract=%s&interval=%s&limit=%d",
		contract, interval, limit)

	respBody, err := c.request("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึง candlesticks ได้: %v", err)
	}

	var rawData [][]interface{}
	err = json.Unmarshal(respBody, &rawData)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse candlesticks data ได้: %v", err)
	}

	var candlesticks []Candlestick
	for _, item := range rawData {
		if len(item) >= 6 {
			timestamp, _ := strconv.ParseInt(fmt.Sprintf("%.0f", item[0]), 10, 64)
			volume, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[1]), 64)
			close, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[2]), 64)
			high, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[3]), 64)
			low, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[4]), 64)
			open, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[5]), 64)

			candlesticks = append(candlesticks, Candlestick{
				Timestamp: timestamp,
				Volume:    volume,
				Close:     close,
				High:      high,
				Low:       low,
				Open:      open,
			})
		}
	}

	return candlesticks, nil
}

// GetContracts ดึงรายการ contracts ทั้งหมด
func (c *Client) GetContracts() ([]Contract, error) {
	respBody, err := c.request("GET", "/api/v4/futures/usdt/contracts", nil)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึง contracts ได้: %v", err)
	}

	var contracts []Contract
	err = json.Unmarshal(respBody, &contracts)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse contracts data ได้: %v", err)
	}

	// กรองเฉพาะ USDT pairs และไม่อยู่ใน delisting
	var validContracts []Contract
	for _, contract := range contracts {
		if strings.HasSuffix(contract.Name, "_USDT") && !contract.InDelisting {
			validContracts = append(validContracts, contract)
		}
	}

	return validContracts, nil
}

// CreateOrder สร้าง order ใหม่
func (c *Client) CreateOrder(req OrderRequest) (*OrderResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ marshal order request ได้: %v", err)
	}

	respBody, err := c.request("POST", "/api/v4/futures/usdt/orders", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง order ได้: %v", err)
	}

	var order OrderResponse
	err = json.Unmarshal(respBody, &order)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถ parse order response ได้: %v", err)
	}

	return &order, nil
}

// ClosePosition ปิด position โดยการสร้าง market order ฝั่งตรงข้าม
func (c *Client) ClosePosition(contract string) error {
	// ดึงข้อมูล position ปัจจุบัน
	positions, err := c.GetPositions()
	if err != nil {
		return fmt.Errorf("ไม่สามารถดึง positions ได้: %v", err)
	}

	var targetPosition *Position
	for _, pos := range positions {
		if pos.Contract == contract {
			targetPosition = &pos
			break
		}
	}

	if targetPosition == nil {
		return fmt.Errorf("ไม่พบ position สำหรับ %s", contract)
	}

	// สร้าง close order
	orderReq := OrderRequest{
		Contract:      contract,
		Size:          int64(-targetPosition.Size), // ฝั่งตรงข้าม
		ClosePosition: true,
		ReduceOnly:    true,
		Text:          "AI_CLOSE",
	}

	_, err = c.CreateOrder(orderReq)
	if err != nil {
		return fmt.Errorf("ไม่สามารถปิด position ได้: %v", err)
	}

	return nil
}

// SetLeverage ตั้งค่า leverage สำหรับ contract
func (c *Client) SetLeverage(contract string, leverage int) error {
	body := map[string]interface{}{
		"contract": contract,
		"leverage": strconv.Itoa(leverage),
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("ไม่สามารถ marshal leverage request ได้: %v", err)
	}

	_, err = c.request("POST", "/api/v4/futures/usdt/positions/"+contract+"/leverage", bodyBytes)
	if err != nil {
		return fmt.Errorf("ไม่สามารถตั้งค่า leverage ได้: %v", err)
	}

	return nil
}
