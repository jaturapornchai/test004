package trading

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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

// GetFuturesContracts ดึงรายชื่อ contracts ทั้งหมด (USDT pairs เท่านั้น)
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

	fmt.Printf("📋 กรองแล้วได้ %d contracts ที่ active\n", len(usdtContracts))
	return usdtContracts, nil
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

// OpenPositionWithSize เปิด position ใหม่ด้วย Position Size ที่คำนวณแล้ว (ทศนิยม)
func (gc *GateClient) OpenPositionWithSize(contract, side string, leverage float64) (bool, error) {
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

	// ดึงข้อมูล contract multiplier (สำคัญมาก!)
	contractMultiplier, _ := strconv.ParseFloat(contractInfo.QuantoMultiplier, 64)
	if contractMultiplier == 0 {
		contractMultiplier = 1 // default ถ้าไม่มีค่า
	}
	fmt.Printf("📐 Contract Multiplier: %.4f\n", contractMultiplier)

	currentPrice, _ := strconv.ParseFloat(contractInfo.LastPrice, 64)
	minOrderSize := contractInfo.OrderSizeMin
	maxOrderSize := contractInfo.OrderSizeMax

	// สำคัญ! คำนวณ Position Size โดยใช้สูตรใหม่
	// positionSize = (50 * leverage) / (currentPrice * contractMultiplier)
	positionSize := (50 * leverage) / (currentPrice * contractMultiplier)

	fmt.Printf("📊 ข้อมูล %s:\n", contract)
	fmt.Printf("   ราคาปัจจุบัน: %.6f\n", currentPrice)
	fmt.Printf("   Min Order Size: %d\n", minOrderSize)
	fmt.Printf("   Max Order Size: %d\n", maxOrderSize)

	// ใช้ Position Size ที่คำนวณแล้ว หรือใช้สูตร fallback
	var size float64
	if positionSize > 0 {
		size = positionSize
		fmt.Printf("📐 ใช้ Position Size ที่คำนวณแล้ว: %.6f contracts\n", size)
	} else {
		// คำนวณ Position Size โดยใช้สูตรใหม่: size = 10 / current_price
		targetSize := 10.0 / currentPrice
		size = targetSize
		if size < 1 {
			size = 1
		}
		fmt.Printf("📐 คำนวณ Position Size: 10 / %.6f = %.6f → %.6f contracts\n", currentPrice, targetSize, size)
	}

	// ตรวจสอบขอบเขต
	minOrderSizeFloat := float64(minOrderSize)
	maxOrderSizeFloat := float64(maxOrderSize)

	if size < minOrderSizeFloat {
		fmt.Printf("⚠️ Size (%.6f) น้อยกว่า Min (%.0f) - ปรับเป็น %.0f\n", size, minOrderSizeFloat, minOrderSizeFloat)
		size = minOrderSizeFloat
	}
	if maxOrderSizeFloat > 0 && size > maxOrderSizeFloat {
		fmt.Printf("⚠️ Size (%.6f) เกิน Max (%.0f) - ปรับเป็น %.0f\n", size, maxOrderSizeFloat, maxOrderSizeFloat)
		size = maxOrderSizeFloat
	}

	// คำนวณ notional value จริง
	actualNotional := size * currentPrice

	fmt.Printf("📐 Final Position Size: %.6f contracts (Notional: $%.6f)\n", size, actualNotional)

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
		order.Size = -int64(size + 0.5) // ปัดขึ้น
	} else {
		order.Size = int64(size + 0.5) // ปัดขึ้น
	}
	fmt.Printf("ปัดเป็น: %d contracts\n", order.Size)

	fmt.Printf("🚀 ส่งคำสั่งเปิด position: %s %.6f contracts (rounded to %d)\n", contract, size, order.Size)

	createdOrder, _, err := futuresApi.CreateFuturesOrder(gc.ctx, "usdt", order)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเปิด position ได้: %v\n", err)
		return false, err
	}

	fmt.Printf("📋 Order ID: %d, Status: %s\n", createdOrder.Id, createdOrder.Status)

	if createdOrder.Status == "finished" {
		fmt.Printf("✅ เปิด position %s สำเร็จ!\n", contract)
		fmt.Printf("💵 Notional Value: $%.6f (สูตรใหม่: 10/ราคา)\n", actualNotional)

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
func (gc *GateClient) OpenPosition(contract, side string, margin, leverage float64) (bool, error) {
	futuresApi := gc.client.FuturesApi

	fmt.Printf("💰 เตรียมเปิด position %s %s (Margin: %.2f USDT, Leverage: %.0fx)\n",
		side, contract, margin, leverage)

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

	// ไม่คำนวณ position size ที่นี่ - ให้ AI คำนวณแทน
	// เก็บข้อมูลพื้นฐานสำหรับ AI
	minOrderSize := contractInfo.OrderSizeMin
	maxOrderSize := contractInfo.OrderSizeMax

	fmt.Printf("📊 ข้อมูล %s:\n", contract)
	fmt.Printf("   ราคาปัจจุบัน: %.6f\n", currentPrice)
	fmt.Printf("   Min Order Size: %d\n", minOrderSize)
	fmt.Printf("   Max Order Size: %d\n", maxOrderSize)
	fmt.Printf("   💡 AI จะคำนวณ position size ที่เหมาะสม\n")

	// ใช้สูตรใหม่: position_size = 10 / current_price
	targetSize := 10.0 / currentPrice
	size := targetSize

	// ตรวจสอบขอบเขต
	minOrderSizeFloat := float64(minOrderSize)
	maxOrderSizeFloat := float64(maxOrderSize)

	if size < minOrderSizeFloat {
		size = minOrderSizeFloat
	}
	if maxOrderSizeFloat > 0 && size > maxOrderSizeFloat {
		size = maxOrderSizeFloat
	}
	if size == 0 {
		size = 1
	}

	// ป้องกัน position ใหญ่เกินไป
	actualNotional := size * currentPrice
	if actualNotional > 50.0 { // ลดจาก 100 เป็น 50
		size = 50.0 / currentPrice
		if size < 1 {
			size = 1
		}
		fmt.Printf("📐 ปรับ position size: %.6f → %.6f (ลด notional)\n", actualNotional/currentPrice, size)
	}

	fmt.Printf("📐 Position Size: %.6f contracts (Formula: 10/%.6f = %.6f)\n", size, currentPrice, targetSize)

	// สร้าง order (แปลงเป็น int64 โดยปัดทศนิยม)
	order := gateapi.FuturesOrder{
		Contract: contract,
		Price:    "0",   // market order
		Tif:      "ioc", // immediate or cancel
		Text:     "t-bot-auto",
	}

	// เปิด SHORT position
	if side == "short" {
		order.Size = -int64(size + 0.5) // ปัดขึ้น
	} else {
		order.Size = int64(size + 0.5) // ปัดขึ้น
	}

	fmt.Printf("🚀 ส่งคำสั่งปิด position: %s %d contracts\n", contract, order.Size)

	createdOrder, _, err := futuresApi.CreateFuturesOrder(gc.ctx, "usdt", order)
	if err != nil {
		fmt.Printf("❌ ไม่สามารถเปิด position ได้: %v\n", err)
		return false, err
	}

	fmt.Printf("📋 Order ID: %d, Status: %s\n", createdOrder.Id, createdOrder.Status)

	if createdOrder.Status == "finished" {
		fmt.Printf("✅ เปิด position %s สำเร็จ!\n", contract)
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

// CalculateEMA คำนวณ Exponential Moving Average
func (gc *GateClient) CalculateEMA(ohlcv []OHLCV, period int) float64 {
	if len(ohlcv) < period {
		return 0
	}

	// คำนวณ EMA โดยใช้ Close price
	k := 2.0 / float64(period+1) // multiplier
	ema := ohlcv[0].Close        // เริ่มต้นด้วยราคาปิดแรก

	for i := 1; i < len(ohlcv); i++ {
		ema = ohlcv[i].Close*k + ema*(1-k)
	}

	return ema
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

// GetHighVolumeContracts ดึง contracts ที่มี volume มากกว่าที่กำหนด (USDT pairs เท่านั้น)
func (gc *GateClient) GetHighVolumeContracts(minVolumeUSDT float64) ([]string, error) {
	fmt.Printf("🔍 กรอง contracts ที่มี volume มากกว่า $%.0f ต่อวัน...\n", minVolumeUSDT)

	// ดึงรายชื่อ contracts ทั้งหมดก่อน
	contracts, err := gc.GetFuturesContracts()
	if err != nil {
		return nil, err
	}

	var highVolumeContracts []string

	// ตรวจสอบ volume ของแต่ละ contract
	for _, contract := range contracts {
		// ดึงข้อมูล OHLCV ย้อนหลัง 24 ชั่วโมง (24 candles ของ 1h)
		ohlcv, err := gc.GetOHLCV(contract, "1h", 24)
		if err != nil {
			fmt.Printf("⚠️ ไม่สามารถดึงข้อมูล %s: %v\n", contract, err)
			continue
		}

		if len(ohlcv) == 0 {
			continue
		}

		// คำนวณ volume รวม 24 ชั่วโมง
		totalVolume := 0.0
		avgPrice := 0.0
		for _, candle := range ohlcv {
			totalVolume += candle.Volume
			avgPrice += (candle.High + candle.Low + candle.Close) / 3.0
		}

		if len(ohlcv) > 0 {
			avgPrice /= float64(len(ohlcv))
		}

		// คำนวณ volume ในหน่วย USDT (volume * ราคาเฉลี่ย)
		volumeUSDT := totalVolume * avgPrice

		if volumeUSDT >= minVolumeUSDT {
			highVolumeContracts = append(highVolumeContracts, contract)
			fmt.Printf("✅ %s - Volume: $%.0f (ผ่านเกณฑ์)\n", contract, volumeUSDT)
		} else {
			fmt.Printf("❌ %s - Volume: $%.0f (ต่ำเกินไป)\n", contract, volumeUSDT)
		}

		// พักเล็กน้อยระหว่างเหรียญ (หลีกเลี่ยง rate limit)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("📊 กรองแล้วได้ %d contracts ที่มี volume เพียงพอจาก %d contracts ทั้งหมด\n",
		len(highVolumeContracts), len(contracts))

	return highVolumeContracts, nil
}
