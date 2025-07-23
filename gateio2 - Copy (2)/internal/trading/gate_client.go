package trading

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v5"
)

// GateClient wrapper สำหรับ Gate.io API
type GateClient struct {
	client *gateapi.APIClient
	ctx    context.Context
}

// NewGateClient สร้าง GateClient ใหม่
func NewGateClient(client *gateapi.APIClient, ctx context.Context) *GateClient {
	return &GateClient{
		client: client,
		ctx:    ctx,
	}
}

// TestConnection ทดสอบการเชื่อมต่อ
func (gc *GateClient) TestConnection() bool {
	futuresApi := gc.client.FuturesApi

	account, _, err := futuresApi.ListFuturesAccounts(gc.ctx, "usdt")
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเชื่อมต่อ Gate.io API ได้: %v\n", err)
		return false
	}

	fmt.Printf("✅ เชื่อมต่อ Gate.io สำเร็จ - Balance: %s USDT\n", account.Available)
	return true
}

// GetBalance ดึง balance
func (gc *GateClient) GetBalance() (string, error) {
	futuresApi := gc.client.FuturesApi

	account, _, err := futuresApi.ListFuturesAccounts(gc.ctx, "usdt")
	if err != nil {
		return "0", err
	}

	return account.Available, nil
}

// GetFuturesContracts ดึงรายชื่อ contracts ที่มี volume > $1,000,000 ใน 24h
func (gc *GateClient) GetFuturesContracts() ([]string, error) {
	futuresApi := gc.client.FuturesApi

	contracts, _, err := futuresApi.ListFuturesContracts(gc.ctx, "usdt")
	if err != nil {
		return nil, err
	}

	var usdtContracts []string

	for _, contract := range contracts {
		// เอาเฉพาะ USDT pairs ที่ active
		if len(contract.Name) > 5 &&
			contract.Name[len(contract.Name)-4:] == "USDT" &&
			!contract.InDelisting { // เอาเฉพาะที่ไม่ถูก delisting
			usdtContracts = append(usdtContracts, contract.Name)
		}
	}

	fmt.Printf("📋 พบ %d contracts ทั้งหมด, กำลังกรอง volume...\n", len(usdtContracts))

	// ดึงข้อมูล tickers ทั้งหมดพร้อมกัน
	spotApi := gc.client.SpotApi
	allTickers, _, err := spotApi.ListTickers(gc.ctx, &gateapi.ListTickersOpts{})
	if err != nil {
		return nil, err
	}

	// สร้าง map เพื่อหา ticker ได้เร็ว
	tickerMap := make(map[string]gateapi.Ticker)
	for _, ticker := range allTickers {
		tickerMap[ticker.CurrencyPair] = ticker
	}
	var highVolumeContracts []string
	// กรองตาม volume 24h > $1,000,000
	for _, contract := range usdtContracts {
		if ticker, exists := tickerMap[contract]; exists {
			// ใช้ QuoteVolume ซึ่งเป็น volume ใน USDT โดยตรง
			volumeUSD, err := strconv.ParseFloat(ticker.QuoteVolume, 64)

			if err == nil && volumeUSD > 1000000 { // > $1,000,000
				highVolumeContracts = append(highVolumeContracts, contract)
			}
		}
	}

	fmt.Printf("💰 กรอง volume > $1,000,000 แล้วได้ %d contracts\n", len(highVolumeContracts))
	return highVolumeContracts, nil
}

// GetOHLCV ดึงข้อมูล OHLCV (ใช้ 1h timeframe เสมอ)
func (gc *GateClient) GetOHLCV(contract, interval string, limit int) ([]OHLCV, error) {
	futuresApi := gc.client.FuturesApi

	fmt.Printf("📊 ดึงข้อมูล OHLCV %s (1h timeframe, %d candles)\n", contract, limit)

	// บังคับใช้ 1h timeframe เสมอ (ignore interval parameter)
	candles, _, err := futuresApi.ListFuturesCandlesticks(gc.ctx, "usdt", contract, &gateapi.ListFuturesCandlesticksOpts{
		Interval: optional.NewString("1h"),             // บังคับ 1h สำหรับการวิเคราะห์ทั้งหมด
		Limit:    optional.NewInt32(int32(limit + 20)), // ขอข้อมูลเยอะหน่อยเผื่อตัด
	})
	if err != nil {
		return nil, err
	}

	var ohlcv []OHLCV
	for _, candle := range candles {
		timestamp := int64(candle.T)
		open, _ := strconv.ParseFloat(candle.O, 64)
		high, _ := strconv.ParseFloat(candle.H, 64)
		low, _ := strconv.ParseFloat(candle.L, 64)
		close, _ := strconv.ParseFloat(candle.C, 64)
		volume := float64(candle.V)

		ohlcv = append(ohlcv, OHLCV{
			Timestamp: timestamp,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		})
	}

	// ตัดให้เหลือจำนวนที่ต้องการ
	if len(ohlcv) > limit {
		ohlcv = ohlcv[len(ohlcv)-limit:]
	}

	return ohlcv, nil
}

// GetOpenPositions ดึง positions ที่เปิดอยู่ทั้งหมด
func (gc *GateClient) GetOpenPositions() ([]*Position, error) {
	futuresApi := gc.client.FuturesApi

	positions, _, err := futuresApi.ListPositions(gc.ctx, "usdt")
	if err != nil {
		return nil, err
	}

	var openPositions []*Position
	for _, pos := range positions {
		if pos.Size != 0 {
			size := pos.Size
			entryPrice, _ := strconv.ParseFloat(pos.EntryPrice, 64)
			markPrice, _ := strconv.ParseFloat(pos.MarkPrice, 64)
			unrealizedPnl, _ := strconv.ParseFloat(pos.UnrealisedPnl, 64)
			margin, _ := strconv.ParseFloat(pos.Margin, 64)
			leverage, _ := strconv.ParseFloat(pos.Leverage, 64)

			position := &Position{
				Contract:      pos.Contract,
				Size:          size,
				EntryPrice:    entryPrice,
				MarkPrice:     markPrice,
				UnrealizedPnl: unrealizedPnl,
				Margin:        margin,
				Leverage:      leverage,
				Mode:          pos.Mode,
			}

			openPositions = append(openPositions, position)
		}
	}

	return openPositions, nil
}

// HasOpenPosition ตรวจสอบว่ามี position เปิดอยู่หรือไม่
func (gc *GateClient) HasOpenPosition(contract string) (bool, error) {
	futuresApi := gc.client.FuturesApi

	position, _, err := futuresApi.GetPosition(gc.ctx, "usdt", contract)
	if err != nil {
		// ถ้า contract ไม่มีอยู่ ให้ถือว่าไม่มี position
		errorStr := err.Error()
		if strings.Contains(errorStr, "POSITION_NOT_FOUND") ||
			strings.Contains(errorStr, "CONTRACT_NOT_EXISTS") ||
			strings.Contains(errorStr, "INVALID_CONTRACT") {
			return false, nil // ไม่มี position (เพราะ contract ไม่มี)
		}
		return false, err // error อื่นๆ ให้ return error
	}

	return position.Size != 0, nil
}

// OpenPosition เปิด position ใหม่ด้วย 15 USDT target margin
func (gc *GateClient) OpenPosition(contract, side string, leverage float64) (bool, error) {
	// ตั้ง size โดยปัดเป็น int64
	futuresApi := gc.client.FuturesApi

	// 1️⃣ ตั้งค่า Leverage และ Margin Mode ก่อนเปิด position
	err := gc.SetLeverageAndMarginMode(contract, leverage)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถตั้งค่า leverage/margin mode ได้: %v\n", err)
		return false, err
	}

	// 2️⃣ ดึงข้อมูล contract
	contractInfo, _, err := futuresApi.GetFuturesContract(gc.ctx, "usdt", contract)
	if err != nil {
		return false, fmt.Errorf("ไม่สามารถดึงข้อมูล contract ได้: %v", err)
	}

	currentPrice, _ := strconv.ParseFloat(contractInfo.LastPrice, 64)
	fmt.Printf("📊 ราคาปัจจุบัน %s: $%.3f\\n", contract, currentPrice)

	// ดึงข้อมูล contract multiplier (สำคัญมาก!)
	contractMultiplier, _ := strconv.ParseFloat(contractInfo.QuantoMultiplier, 64)
	if contractMultiplier == 0 {
		contractMultiplier = 1 // default ถ้าไม่มีค่า
	}
	fmt.Printf("📐 Contract Multiplier: %.4f\n", contractMultiplier)

	minOrderSize := contractInfo.OrderSizeMin
	maxOrderSize := contractInfo.OrderSizeMax

	fmt.Printf("📊 ข้อมูล %s:\n", contract)
	fmt.Printf("   ราคาปัจจุบัน: %.6f\n", currentPrice)
	fmt.Printf("   Min Order Size: %d\n", minOrderSize)
	fmt.Printf("   Max Order Size: %d\n", maxOrderSize)

	targetMargin := 15.0 // $15 USDT margin target
	leverage = 5.0

	targetSize := (targetMargin * leverage) / (currentPrice * contractMultiplier)
	size := int64(math.Round(targetSize))

	// ตรวจสอบว่า size ต้องไม่เป็น 0
	if size == 0 {
		size = 1
	}

	// คำนวณ stop loss 5% เพื่อแจ้งข้อมูล (AI จะดูแลแทน)
	var stopLossPrice float64
	if side == "long" {
		// Long position: stop loss = ราคาปัจจุบัน - 5%
		stopLossPrice = currentPrice * 0.95
		fmt.Printf("📊 Long Stop Loss Reference: %.6f (ราคาปัจจุบัน - 5%% - AI monitoring)\n", stopLossPrice)
	} else {
		// Short position: stop loss = ราคาปัจจุบัน + 5%
		stopLossPrice = currentPrice * 1.05
		fmt.Printf("📊 Short Stop Loss Reference: %.6f (ราคาปัจจุบัน + 5%% - AI monitoring)\n", stopLossPrice)
	}

	// สร้าง order (แปลงเป็น int64 โดยปัดทศนิยม)
	order := gateapi.FuturesOrder{
		Contract: contract,
		Price:    "0",   // market order
		Tif:      "ioc", // immediate or cancel
		Text:     "t-ai-bot",
	}

	// ไม่ตั้ง stop loss order อัตโนมัติ (ให้ AI ดูแลแทน)
	if stopLossPrice > 0 {
		fmt.Printf("🤖 AI จะดูแล position และตัดสินใจปิดตามสถานการณ์\n")
	}

	if side == "short" {
		order.Size = -size
	} else {
		order.Size = size
	}
	fmt.Printf("ปัดเป็น: %d contracts\n", order.Size)

	fmt.Printf("🚀 ส่งคำสั่งเปิด position: %s %d contracts\n", contract, order.Size)

	createdOrder, _, err := futuresApi.CreateFuturesOrder(gc.ctx, "usdt", order)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเปิด position ได้: %v\n", err)
		return false, err
	}

	fmt.Printf("📋 Order ID: %d, Status: %s\n", createdOrder.Id, createdOrder.Status)

	if createdOrder.Status == "finished" {
		fmt.Printf("✅ เปิด position %s สำเร็จ!\n", contract)

		// ไม่ตั้ง stop loss order อัตโนมัติ - ให้ AI ตัดสินใจปิด position เอง
		if stopLossPrice > 0 {
			fmt.Printf("🤖 AI จะดูแล position และตัดสินใจปิดเมื่อจำเป็น\n")
			fmt.Printf("📝 ไม่ตั้ง stop loss order อัตโนมัติ - ให้ AI วิเคราะห์แทน\n")
		}

		return true, nil
	} else {
		fmt.Printf("⚠️ Order status: %s\n", createdOrder.Status)
		return false, nil
	}
}

// ClosePosition ปิด position
func (gc *GateClient) ClosePosition(contract string) (bool, error) {
	futuresApi := gc.client.FuturesApi

	// ดึงข้อมูล position ปัจจุบัน
	position, _, err := futuresApi.GetPosition(gc.ctx, "usdt", contract)
	if err != nil {
		return false, err
	}

	if position.Size == 0 {
		return true, nil // ไม่มี position ที่ต้องปิด
	}

	// สร้าง order ปิด position (ใช้ size ตรงข้าม)
	closeOrder := gateapi.FuturesOrder{
		Contract: contract,
		Size:     -position.Size,
		Price:    "0",   // market order
		Tif:      "ioc", // immediate or cancel
		Text:     "t-bot-close",
	}

	createdOrder, _, err := futuresApi.CreateFuturesOrder(gc.ctx, "usdt", closeOrder)
	if err != nil {
		return false, err
	}

	return createdOrder.Status == "finished", nil
}

// SetLeverageAndMarginMode ตั้งค่า leverage และ margin mode สำหรับ contract
func (gc *GateClient) SetLeverageAndMarginMode(contract string, leverage float64) error {
	futuresApi := gc.client.FuturesApi

	fmt.Printf("🔧 ตั้งค่า Leverage และ Margin Mode สำหรับ %s...\n", contract)

	// ตรวจสอบ position ปัจจุบันเพื่อเช็คการตั้งค่า
	position, _, err := futuresApi.GetPosition(gc.ctx, "usdt", contract)
	if err != nil {
		fmt.Printf("⚠️ ไม่สามารถดึงข้อมูล position ได้: %v\n", err)
		// ถ้าไม่สามารถดึงได้ ลองตั้งค่าต่อไป
	}

	// ตั้งค่า leverage ถ้าต่างจากที่ต้องการ
	if err == nil {
		currentLeverage, _ := strconv.ParseFloat(position.Leverage, 64)

		if currentLeverage != leverage {
			fmt.Printf("🔧 กำลังตั้งค่า Leverage = %.0fx (จาก %.0fx)...\n", leverage, currentLeverage)
			leverageStr := fmt.Sprintf("%.0f", leverage)
			_, _, err = futuresApi.UpdatePositionLeverage(gc.ctx, "usdt", contract, leverageStr)
			if err != nil {
				fmt.Printf("⚠️ การตั้งค่า leverage มีปัญหา: %v\n", err)
				// ไม่ return error เพราะบางครั้งอาจตั้งค่าไม่ได้แต่ใช้งานได้
			} else {
				fmt.Printf("✅ ตั้งค่า Leverage = %.0fx สำเร็จ\n", leverage)
			}
		} else {
			fmt.Printf("✅ Leverage อยู่ที่ %.0fx อยู่แล้ว\n", leverage)
		}

		// ตั้งค่า margin mode = isolated (single) ถ้าต่างจากที่ต้องการ
		if position.Mode != "single" {
			fmt.Printf("🔧 กำลังตั้งค่า Margin Mode = isolated (จาก %s)...\n", position.Mode)
			// ใช้ UpdatePositionMargin กับ amount = "0" เพื่อเปลี่ยนเป็น isolated
			_, _, err = futuresApi.UpdatePositionMargin(gc.ctx, "usdt", contract, "0")
			if err != nil {
				fmt.Printf("⚠️ การตั้งค่า margin mode มีปัญหา: %v\n", err)
				// ไม่ return error เพราะบางครั้งอาจตั้งค่าไม่ได้แต่ใช้งานได้
			} else {
				fmt.Printf("✅ ตั้งค่า Margin Mode = isolated สำเร็จ\n")
			}
		} else {
			fmt.Printf("✅ Margin Mode อยู่ที่ isolated อยู่แล้ว\n")
		}
	} else {
		// ถ้าไม่มี position ให้ตั้งค่าไปก่อน
		fmt.Printf("🔧 ตั้งค่า Leverage = %.0fx...\n", leverage)
		leverageStr := fmt.Sprintf("%.0f", leverage)
		_, _, err = futuresApi.UpdatePositionLeverage(gc.ctx, "usdt", contract, leverageStr)
		if err != nil {
			fmt.Printf("⚠️ การตั้งค่า leverage มีปัญหา: %v\n", err)
		} else {
			fmt.Printf("✅ ตั้งค่า Leverage = %.0fx สำเร็จ\n", leverage)
		}
	}

	return nil
}

// CreateStopLossOrder สร้าง stop loss order แบบ limit order ที่ราคา 5%
func (gc *GateClient) CreateStopLossOrder(contract string, stopPrice float64, positionSize int64) error {
	futuresApi := gc.client.FuturesApi

	fmt.Printf("🛡️ กำลังสร้าง Stop Loss Order...\n")
	fmt.Printf("   Contract: %s\n", contract)
	fmt.Printf("   Stop Price: %.6f (ราคา 5%%)\n", stopPrice)
	fmt.Printf("   Position Size: %d\n", positionSize)

	// สร้าง stop loss โดยใช้ reduce-only limit order ที่ราคา 1%
	var orderSide int64

	if positionSize > 0 {
		// Long position: ขายเพื่อปิด position
		orderSide = -positionSize
		fmt.Printf("   📉 Long Position: ตั้ง limit sell ที่ %.6f (ราคา - 5%%)\n", stopPrice)
	} else {
		// Short position: ซื้อเพื่อปิด position
		orderSide = -positionSize
		fmt.Printf("   📈 Short Position: ตั้ง limit buy ที่ %.6f (ราคา + 5%%)\n", stopPrice)
	}

	// สร้าง limit order ที่ราคา 5% เป็น reduce-only
	stopPriceStr := fmt.Sprintf("%.6f", stopPrice)
	stopOrder := gateapi.FuturesOrder{
		Contract:   contract,
		Size:       orderSide,
		Price:      stopPriceStr,  // limit order ที่ราคา 5%
		Tif:        "gtc",         // good till canceled
		Text:       "t-stop-5pct", // ชื่อใหม่
		ReduceOnly: true,          // ปิด position เท่านั้น
	}

	// ลองสร้าง order
	createdOrder, _, err := futuresApi.CreateFuturesOrder(gc.ctx, "usdt", stopOrder)
	if err != nil {
		fmt.Printf("⚠️ Gate.io ไม่รองรับ conditional order - ใช้ manual monitoring\n")
		fmt.Printf("📝 จะตรวจสอบ stop loss ใน trading loop แทน\n")
		return nil
	}

	fmt.Printf("✅ สร้าง Stop Loss Order สำเร็จ!\n")
	fmt.Printf("   Order ID: %d\n", createdOrder.Id)
	fmt.Printf("   Status: %s\n", createdOrder.Status)

	return nil
}

// CheckStopLoss ตรวจสอบและทำ stop loss แบบ manual (ใช้ราคา 5%)
func (gc *GateClient) CheckStopLoss(contract string, stopPrice float64, isLong bool) (bool, error) {
	futuresApi := gc.client.FuturesApi

	// ดึงข้อมูล contract เพื่อดูราคาปัจจุบัน
	contractInfo, _, err := futuresApi.GetFuturesContract(gc.ctx, "usdt", contract)
	if err != nil {
		return false, err
	}

	currentPrice, _ := strconv.ParseFloat(contractInfo.LastPrice, 64)

	fmt.Printf("🔍 ตรวจสอบ Stop Loss: %s\n", contract)
	fmt.Printf("   ราคาปัจจุบัน: %.6f\n", currentPrice)
	fmt.Printf("   Stop Price: %.6f (ราคา 5%%)\n", stopPrice)

	// ตรวจสอบเงื่อนไข stop loss
	shouldStop := false
	if isLong && currentPrice <= stopPrice {
		// Long position: ถ้าราคาลงต่ำกว่าราคา - 5%
		shouldStop = true
		fmt.Printf("🚨 Long Position: ราคาลงต่ำกว่าราคา - 5%% - ต้อง STOP LOSS!\n")
	} else if !isLong && currentPrice >= stopPrice {
		// Short position: ถ้าราคาขึ้นสูงกว่าราคา + 5%
		shouldStop = true
		fmt.Printf("🚨 Short Position: ราคาขึ้นสูงกว่าราคา + 5%% - ต้อง STOP LOSS!\n")
	}

	if shouldStop {
		// ปิด position ทันที
		fmt.Printf("⚡ ปิด position เพื่อ Stop Loss ทันที!\n")
		success, err := gc.ClosePosition(contract)
		if err != nil {
			fmt.Printf("❌ ไม่สามารถปิด position สำหรับ stop loss ได้: %v\n", err)
			return false, err
		}

		if success {
			fmt.Printf("✅ Stop Loss สำเร็จ: ปิด position %s แล้ว\n", contract)
			return true, nil
		}
	}

	return false, nil
}
