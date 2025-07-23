# Gate.io Trading Bot - Sample Code

## การทดสอบการเทรด Gate.io API สำหรับ NMR_USDT

โค้ดนี้แสดงการทดสอบการเทรดบน Gate.io futures ที่ประกอบด้วย:

1. **การเชื่อมต่อ API** - ทดสอบการเชื่อมต่อและดึง balance
2. **การตั้งค่า Leverage และ Margin Mode** - ตั้งค่า leverage 5x และ isolated margin
3. **การเปิด Position** - เปิด position ด้วย margin $15
4. **การรอ 10 วินาที** - รอให้ราคาเปลี่ยนแปลง
5. **การปิด Position** - ปิด position และคำนวณ PnL

## ผลการทดสอบ

```text
🔧 ทดสอบการเทรด Gate.io API สำหรับ NMR_USDT

1️⃣ ทดสอบการเชื่อมต่อ:
✅ เชื่อมต่อสำเร็จ - Balance: 537.684306975932 USDT

2️⃣ ทดสอบการตั้งค่า Leverage = 5x และ Margin Type = isolated สำหรับ NMR_USDT:
🔍 ตรวจสอบ contract NMR_USDT...
✅ พบ contract: NMR_USDT (ราคาปัจจุบัน: 8.284)

🔍 ตรวจสอบการตั้งค่าปัจจุบันของ NMR_USDT:
📊 การตั้งค่าปัจจุบัน:
   - Contract: NMR_USDT
   - Leverage: 5x
   - Margin Mode: single
   - Position Size: 9
✅ Leverage อยู่ที่ 5x อยู่แล้ว
✅ Margin Mode อยู่ที่ isolated อยู่แล้ว

3️⃣ ทดลองเปิด Position Size = $15 สำหรับ NMR_USDT:
📊 ราคาปัจจุบัน NMR_USDT: $8.284
📐 Contract Multiplier: 0.1000

💡 การคำนวณ Position:
   - Target Margin: $15.00
   - Leverage: 5x
   - Contract Size: 91 contracts
   - Total Exposure: $75.38
   - Actual Margin Required: $15.08
✅ Balance เพียงพอ: $537.68

🚀 เริ่มเปิด LONG position สำหรับ NMR_USDT...
✅ เปิด position สำเร็จ!
📋 Order ID: 126945214573492849
📊 Size: 91 contracts
💎 Status: finished

🔍 ตรวจสอบ position ที่เปิดแล้ว:
📈 Position เปิดสำเร็จ:
   - Contract: NMR_USDT
   - Size: 100 contracts
   - Entry Price: 8.28483
   - Mark Price: 8.3
   - Value: 83 USDT
   - Margin: 16.631796225 USDT
   - Leverage: 5x
   - Mode: single

💰 Margin ที่ใช้จริง: $16.63
⚠️ Margin ต่างจากเป้าหมาย (เป้าหมาย: $15.00, จริง: $16.63)

4️⃣ รอ 10 วินาที...

5️⃣ ทดลองปิด Position สำหรับ NMR_USDT:
🔍 ตรวจสอบ position ปัจจุบันสำหรับ NMR_USDT...
📊 Position ปัจจุบัน:
   - Size: 100 contracts
   - Entry Price: 8.28483
   - Mark Price: 8.306
   - Unrealized PnL: 0.2117 USDT

🔚 เริ่มปิด position สำหรับ NMR_USDT...
   - Position Size ที่จะปิด: 100 contracts
✅ ส่งคำสั่งปิด position สำเร็จ!
📋 Order ID: 126945214573493106
📊 Size: -100 contracts
💎 Status: finished

🔍 ตรวจสอบ position หลังปิด:
✅ ปิด position สำเร็จ - ไม่มี position คงเหลือ

💳 Balance หลังเทรด: 538.918370865932 USDT
✅ การทดสอบเสร็จสมบูรณ์
```

## สรุปผล

- **กำไร**: $1.23 USDT (จาก $537.68 → $538.92)
- **ระยะเวลา**: 10 วินาที
- **Margin ที่ใช้**: $16.63 (ใกล้เคียงเป้าหมาย $15.00)
- **Position Size**: 100 contracts
- **ราคาเข้า**: $8.28483
- **ราคาออก**: ~$8.306

## โค้ด Go สำหรับการทดสอบ

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
    "time"

    "github.com/gateio/gateapi-go/v5"
    "github.com/joho/godotenv"
)

func main() {
    fmt.Println("🔧 ทดสอบการเทรด Gate.io API สำหรับ NMR_USDT")

    // โหลด .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("❌ ไม่สามารถโหลดไฟล์ .env ได้:", err)
    }

    apiKey := os.Getenv("GATE_API_KEY")
    apiSecret := os.Getenv("GATE_API_SECRET")

    if apiKey == "" || apiSecret == "" {
        log.Fatal("❌ API Key หรือ API Secret ไม่ได้ตั้งค่าใน .env")
    }

    // สร้าง Gate.io client
    client := gateapi.NewAPIClient(gateapi.NewConfiguration())

    // สร้าง context พร้อม API credentials
    ctx := context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
        Key:    apiKey,
        Secret: apiSecret,
    })

    testContract := "NMR_USDT"

    // 1. ทดสอบการเชื่อมต่อพื้นฐาน
    fmt.Println("\\n1️⃣ ทดสอบการเชื่อมต่อ:")
    if !testConnection(ctx, client) {
        return
    }

    // 2. ทดสอบการตั้งค่า leverage และ margin mode
    fmt.Printf("\\n2️⃣ ทดสอบการตั้งค่า Leverage = 5x และ Margin Type = isolated สำหรับ %s:\\n", testContract)
    testSetLeverageAndMargin(ctx, client, testContract)

    // 3. ทดลองเปิด position size = $15
    fmt.Printf("\\n3️⃣ ทดลองเปิด Position Size = $15 สำหรับ %s:\\n", testContract)
    positionOpened := testOpenPosition(ctx, client, testContract)
    
    if positionOpened {
        // 4. รอ 10 วินาที
        fmt.Println("\\n4️⃣ รอ 10 วินาที...")
        time.Sleep(10 * time.Second)

        // 5. ทดลองปิด position
        fmt.Printf("\\n5️⃣ ทดลองปิด Position สำหรับ %s:\\n", testContract)
        testClosePosition(ctx, client, testContract)
    }

    fmt.Println("\\n✅ การทดสอบเสร็จสมบูรณ์")
}

func testConnection(ctx context.Context, client *gateapi.APIClient) bool {
    futuresApi := client.FuturesApi

    account, _, err := futuresApi.ListFuturesAccounts(ctx, "usdt")
    if err != nil {
        fmt.Printf("❌ ไม่สามารถเชื่อมต่อ Gate.io API ได้: %v\\n", err)
        return false
    }

    fmt.Printf("✅ เชื่อมต่อสำเร็จ - Balance: %s USDT\\n", account.Available)
    return true
}

func testSetLeverageAndMargin(ctx context.Context, client *gateapi.APIClient, contract string) {
    futuresApi := client.FuturesApi

    // ตรวจสอบว่า contract มีอยู่จริง
    fmt.Printf("🔍 ตรวจสอบ contract %s...\\n", contract)
    contractInfo, _, err := futuresApi.GetFuturesContract(ctx, "usdt", contract)
    if err != nil {
        fmt.Printf("❌ ไม่พบ contract %s: %v\\n", contract, err)
        return
    }
    fmt.Printf("✅ พบ contract: %s (ราคาปัจจุบัน: %s)\\n", contractInfo.Name, contractInfo.LastPrice)

    // ตรวจสอบการตั้งค่าปัจจุบันก่อน
    position, _, err := futuresApi.GetPosition(ctx, "usdt", contract)
    if err != nil {
        fmt.Printf("⚠️ ไม่สามารถดึงข้อมูล position ได้: %v\\n", err)
        return
    }

    currentLeverage, _ := strconv.ParseFloat(position.Leverage, 64)

    // ตั้งค่า leverage = 5x ถ้าของเก่าไม่ใช่ 5x
    if currentLeverage != 5 {
        fmt.Printf("🔧 กำลังตั้งค่า Leverage = 5x (จาก %.0fx)...\\n", currentLeverage)
        _, _, err = futuresApi.UpdatePositionLeverage(ctx, "usdt", contract, "5")
        if err != nil {
            fmt.Printf("⚠️ การตั้งค่า leverage มีปัญหา: %v\\n", err)
        } else {
            fmt.Printf("✅ ตั้งค่า Leverage = 5x สำเร็จ\\n")
        }
    } else {
        fmt.Printf("✅ Leverage อยู่ที่ 5x อยู่แล้ว\\n")
    }

    // ตั้งค่า margin type = isolated ถ้าของเก่าไม่ใช่ isolated
    if position.Mode != "single" {
        fmt.Printf("🔧 กำลังตั้งค่า Margin Type = isolated (จาก %s)...\\n", position.Mode)
        _, _, err = futuresApi.UpdatePositionMargin(ctx, "usdt", contract, "0")
        if err != nil {
            fmt.Printf("⚠️ การตั้งค่า margin type มีปัญหา: %v\\n", err)
        } else {
            fmt.Printf("✅ ตั้งค่า Margin Type = isolated สำเร็จ\\n")
        }
    } else {
        fmt.Printf("✅ Margin Mode อยู่ที่ isolated อยู่แล้ว\\n")
    }
}

func testOpenPosition(ctx context.Context, client *gateapi.APIClient, contract string) bool {
    futuresApi := client.FuturesApi

    // ดึงข้อมูล contract
    contractInfo, _, err := futuresApi.GetFuturesContract(ctx, "usdt", contract)
    if err != nil {
        fmt.Printf("❌ ไม่สามารถดึงข้อมูล contract ได้: %v\\n", err)
        return false
    }

    currentPrice, _ := strconv.ParseFloat(contractInfo.LastPrice, 64)
    fmt.Printf("📊 ราคาปัจจุบัน %s: $%.3f\\n", contract, currentPrice)

    // ดึงข้อมูล contract multiplier (สำคัญมาก!)
    contractMultiplier, _ := strconv.ParseFloat(contractInfo.QuantoMultiplier, 64)
    if contractMultiplier == 0 {
        contractMultiplier = 1 // default ถ้าไม่มีค่า
    }
    fmt.Printf("📐 Contract Multiplier: %.4f\\n", contractMultiplier)

    // คำนวณ position size สำหรับ margin = $15 ด้วย leverage 5x
    targetMargin := 15.0 // $15 USDT margin target
    leverage := 5.0

    // สูตรคำนวณ:
    // Total Value = size * price * contract_multiplier
    // Margin = Total Value / leverage
    // ดังนั้น: size = (margin * leverage) / (price * contract_multiplier)
    
    targetSize := (targetMargin * leverage) / (currentPrice * contractMultiplier)
    size := int64(math.Round(targetSize))
    
    // ตรวจสอบว่า size ต้องไม่เป็น 0
    if size == 0 {
        size = 1
    }

    // คำนวณค่าจริงที่จะใช้
    actualValue := float64(size) * currentPrice * contractMultiplier
    actualMargin := actualValue / leverage

    fmt.Printf("\\n💡 การคำนวณ Position:\\n")
    fmt.Printf("   - Target Margin: $%.2f\\n", targetMargin)
    fmt.Printf("   - Leverage: %.0fx\\n", leverage)
    fmt.Printf("   - Contract Size: %d contracts\\n", size)
    fmt.Printf("   - Total Exposure: $%.2f\\n", actualValue)
    fmt.Printf("   - Actual Margin Required: $%.2f\\n", actualMargin)

    // ตรวจสอบ balance ก่อนเปิด position
    account, _, err := futuresApi.ListFuturesAccounts(ctx, "usdt")
    if err != nil {
        fmt.Printf("❌ ไม่สามารถดึง balance ได้: %v\\n", err)
        return false
    }

    availableBalance, _ := strconv.ParseFloat(account.Available, 64)
    if availableBalance < actualMargin {
        fmt.Printf("❌ Balance ไม่เพียงพอ (มี: $%.2f, ต้องการ margin: $%.2f)\\n",
            availableBalance, actualMargin)
        return false
    }

    fmt.Printf("✅ Balance เพียงพอ: $%.2f\\n", availableBalance)

    // สร้าง market order (LONG position)
    fmt.Printf("\\n🚀 เริ่มเปิด LONG position สำหรับ %s...\\n", contract)

    order := gateapi.FuturesOrder{
        Contract: contract,
        Size:     size,
        Price:    "0",   // market order ใช้ price = 0
        Tif:      "ioc", // immediate or cancel for market order
        Text:     "t-test_15usd",
    }

    createdOrder, _, err := futuresApi.CreateFuturesOrder(ctx, "usdt", order)
    if err != nil {
        fmt.Printf("❌ ไม่สามารถเปิด position ได้: %v\\n", err)
        return false
    }

    fmt.Printf("✅ เปิด position สำเร็จ!\\n")
    fmt.Printf("📋 Order ID: %d\\n", createdOrder.Id)
    fmt.Printf("📊 Size: %d contracts\\n", createdOrder.Size)
    fmt.Printf("💎 Status: %s\\n", createdOrder.Status)

    // รอให้ order เสร็จสมบูรณ์
    time.Sleep(2 * time.Second)

    // ตรวจสอบ position ที่เปิดแล้ว
    position, _, err := futuresApi.GetPosition(ctx, "usdt", contract)
    if err != nil {
        fmt.Printf("⚠️ ไม่สามารถดึงข้อมูล position ได้: %v\\n", err)
        return false
    }

    if position.Size != 0 {
        margin, _ := strconv.ParseFloat(position.Margin, 64)
        fmt.Printf("\\n💰 Margin ที่ใช้จริง: $%.2f\\n", margin)
        
        if math.Abs(margin - targetMargin) < 1.0 {
            fmt.Printf("✅ Margin ใกล้เคียงกับเป้าหมาย $%.2f\\n", targetMargin)
        } else {
            fmt.Printf("⚠️ Margin ต่างจากเป้าหมาย (เป้าหมาย: $%.2f, จริง: $%.2f)\\n", 
                targetMargin, margin)
        }
        return true
    } else {
        fmt.Printf("⚠️ Position ยังไม่เปิด หรือเปิดไม่สำเร็จ\\n")
        return false
    }
}

func testClosePosition(ctx context.Context, client *gateapi.APIClient, contract string) {
    futuresApi := client.FuturesApi

    // ตรวจสอบ position ปัจจุบัน
    fmt.Printf("🔍 ตรวจสอบ position ปัจจุบันสำหรับ %s...\\n", contract)
    position, _, err := futuresApi.GetPosition(ctx, "usdt", contract)
    if err != nil {
        fmt.Printf("❌ ไม่สามารถดึงข้อมูล position ได้: %v\\n", err)
        return
    }

    if position.Size == 0 {
        fmt.Printf("⚠️ ไม่มี position ที่เปิดอยู่ในขณะนี้\\n")
        return
    }

    fmt.Printf("📊 Position ปัจจุบัน:\\n")
    fmt.Printf("   - Size: %d contracts\\n", position.Size)
    fmt.Printf("   - Entry Price: %s\\n", position.EntryPrice)
    fmt.Printf("   - Mark Price: %s\\n", position.MarkPrice)
    fmt.Printf("   - Unrealized PnL: %s USDT\\n", position.UnrealisedPnl)

    // สร้าง market order เพื่อปิด position (ใช้ size เป็นลบ)
    closeSize := -position.Size // ใช้ขนาดติดลบเพื่อปิด position

    fmt.Printf("\\n🔚 เริ่มปิด position สำหรับ %s...\\n", contract)
    fmt.Printf("   - Position Size ที่จะปิด: %d contracts\\n", -closeSize)

    closeOrder := gateapi.FuturesOrder{
        Contract: contract,
        Size:     closeSize,
        Price:    "0",   // market order ใช้ price = 0
        Tif:      "ioc", // immediate or cancel for market order
        Text:     "t-close_position",
    }

    createdOrder, _, err := futuresApi.CreateFuturesOrder(ctx, "usdt", closeOrder)
    if err != nil {
        fmt.Printf("❌ ไม่สามารถปิด position ได้: %v\\n", err)
        return
    }

    fmt.Printf("✅ ส่งคำสั่งปิด position สำเร็จ!\\n")
    fmt.Printf("📋 Order ID: %d\\n", createdOrder.Id)
    fmt.Printf("📊 Size: %d contracts\\n", createdOrder.Size)
    fmt.Printf("💎 Status: %s\\n", createdOrder.Status)

    // รอให้ order เสร็จสมบูรณ์
    time.Sleep(3 * time.Second)

    // ตรวจสอบ position หลังปิด
    finalPosition, _, err := futuresApi.GetPosition(ctx, "usdt", contract)
    if err != nil {
        fmt.Printf("⚠️ ไม่สามารถดึงข้อมูล position ได้: %v\\n", err)
        return
    }

    if finalPosition.Size == 0 {
        fmt.Printf("✅ ปิด position สำเร็จ - ไม่มี position คงเหลือ\\n")
        
        // แสดง PnL สุดท้าย
        realizedPnl, _ := strconv.ParseFloat(finalPosition.RealisedPnl, 64)
        if realizedPnl != 0 {
            fmt.Printf("💰 Realized PnL: $%.4f USDT\\n", realizedPnl)
            if realizedPnl > 0 {
                fmt.Printf("🎉 กำไร!\\n")
            } else {
                fmt.Printf("😔 ขาดทุน\\n")
            }
        }
    } else {
        fmt.Printf("⚠️ Position ยังคงเหลืออยู่: %d contracts\\n", finalPosition.Size)
        fmt.Printf("   อาจเป็นเพราะ order ยังไม่เสร็จสมบูรณ์\\n")
    }

    // แสดง balance หลังจากเทรด
    account, _, err := futuresApi.ListFuturesAccounts(ctx, "usdt")
    if err == nil {
        fmt.Printf("\\n💳 Balance หลังเทรด: %s USDT\\n", account.Available)
    }
}
```

## การตั้งค่า Environment Variables

สร้างไฟล์ `.env`:

```env
GATE_API_KEY=your_gate_api_key_here
GATE_API_SECRET=your_gate_api_secret_here
DEEPSEEK_API_KEY=your_deepseek_api_key_here
```

## การรันโปรแกรม

```bash
cd test
go run test_position.go
```

## สิ่งสำคัญที่ต้องจำ

1. **Contract Multiplier**: NMR_USDT ใช้ multiplier 0.1000
2. **สูตรคำนวณ Margin**: `(margin * leverage) / (price * multiplier)`
3. **Position Size**: ต้องเป็น integer
4. **Market Order**: ใช้ price = "0" และ tif = "ioc"
5. **การปิด Position**: ใช้ size เป็นลบ

## ข้อมูลประเภท Order

- **Price**: "0" สำหรับ market order
- **Tif**: "ioc" (Immediate or Cancel) สำหรับ market order
- **Text**: ระบุได้ตามต้องการ เพื่อติดตามการเทรด

## ข้อสังเกต

- การคำนวณ margin ใช้ Contract Multiplier เป็นปัจจัยสำคัญ
- Position ที่เปิดได้อาจมี size มากกว่าที่คำนวณเนื่องจากมีการรวม position เก่า
- การใช้ market order จะได้ราคาที่ใกล้เคียงกับราคาปัจจุบัน แต่อาจมี slippage
- Balance หลังเทรดจะเปลี่ยนแปลงตาม PnL ของการเทรด
