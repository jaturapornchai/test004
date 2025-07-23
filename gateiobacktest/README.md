# 🚀 Gate.io Trading Bot - ระบบเทรดอัตโนมัติ

ระบบเทรดอัตโนมัติที่ใช้ **Pivot Point SuperTrend + EMA100 + AI** สำหรับ Gate.io Futures Trading

## 📋 โครงสร้างโปรเจค

```
gateio-trading-bot/
├── app.go                           # ไฟล์หลักของโปรแกรม
├── go.mod                          # Go module dependencies
├── go.sum                          # Go module checksums
├── internal/trading/               # โค้ดหลักของระบบเทรด
│   ├── bot.go                      # Trading Bot หลัก
│   ├── types.go                    # โครงสร้างข้อมูล
│   ├── gate_client.go              # Gate.io API client
│   ├── ai_client.go                # AI client (DeepSeek)
│   └── indicators.go               # Technical indicators
├── prompts/                        # Manual prompt files
│   ├── open_position_prompt.txt    # Prompt สำหรับเปิด position
│   └── close_position_prompt.txt   # Prompt สำหรับปิด position
├── test/                          # ไฟล์ทดสอบ
│   ├── connection_test.go         # ทดสอบการเชื่อมต่อ
│   ├── test_position.go           # ทดสอบการเปิด/ปิด position
│   └── ...                        # ไฟล์ทดสอบอื่นๆ
└── markdown/                      # เอกสารต้นฉบับ (ห้ามแก้ไข)
    ├── step.md                    # ข้อกำหนดระบบ
    └── sample.md                  # ตัวอย่างการทดสอบ
```

## 🔧 การติดตั้งและรัน

### 1. ตั้งค่า Environment Variables
สร้างไฟล์ `.env` ในโฟลเดอร์หลัก:
```env
GATE_API_KEY=your_gate_api_key_here
GATE_API_SECRET=your_gate_api_secret_here  
DEEPSEEK_API_KEY=your_deepseek_api_key_here
```

### 2. ติดตั้ง Dependencies
```bash
go mod tidy
```

### 3. รันโปรแกรม
```bash
go run app.go
```

## 🤖 คุณสมบัติหลัก

### ✨ Dual Mode System
- **AUTO OPEN Mode**: SuperTrend Confidence ≥ 75% + Signal ≠ NEUTRAL
- **AI Mode**: SuperTrend Confidence < 75% หรือ Signal = NEUTRAL

### 📊 Technical Analysis
- **Pivot Point SuperTrend**: Period 2, ATR Factor 3.0, ATR Period 10
- **EMA100**: Exponential Moving Average 100 periods
- **Time Frame**: 1H เท่านั้น

### 💰 Position Management
- **Position Size**: 15 USDT ต่อ position
- **Leverage**: 5x (isolated margin)
- **Risk-Reward Ratio**: ≥ 3.0
- **AI Confidence**: ≥ 85% สำหรับการตัดสินใจ

### 🔄 Trading Loop
1. **LOOP1**: ตรวจสอบและจัดการ positions ที่เปิดอยู่ (ทุก 1 ชั่วโมง)
2. **LOOP2**: สแกนหาโอกาสใหม่ (batch processing 20 เหรียญ/batch)

## 📈 กฎการเทรด

### ✅ ข้อบังคับ
- ใช้ Pivot Point SuperTrend + EMA100 เท่านั้น
- Time Frame 1H เท่านั้น
- SuperTrend กำหนดทิศทาง (Trend = 1 → LONG, Trend = -1 → SHORT)
- AI ตัดสินใจทั้งการเปิดและปิด position

### ❌ ข้อห้าม
- ห้ามเปิด position ซ้ำในเหรียญเดียวกัน
- ห้ามเปิด position เพิ่มหรือเปลี่ยนขนาด position
- ห้ามตั้ง Stop Loss/Take Profit อัตโนมัติ (เป็น AI-Only Mode)

## 🟢 Long Position Strategy
**เงื่อนไข**: SuperTrend Trend = 1 (Bullish)
1. ติดตาม Trend เมื่อ SuperTrend แสดง Trend = 1
2. เปิด Long เมื่อราคา rebound จาก EMA100 (แท่งเทียนสีเขียว Close > Open)

## 🔴 Short Position Strategy  
**เงื่อนไข**: SuperTrend Trend = -1 (Bearish)
1. ติดตาม Trend เมื่อ SuperTrend แสดง Trend = -1
2. เปิด Short เมื่อราคา reject จาก EMA100 (แท่งเทียนสีแดง Close < Open)

## 🧪 การทดสอบ

### ทดสอบการเชื่อมต่อ
```bash
cd test
go run connection_test.go
```

### ทดสอบการเปิด/ปิด Position
```bash
cd test  
go run test_position.go
```

## 🔐 ความปลอดภัย
- ใช้ isolated margin ลดความเสี่ยง
- Position size เล็ก (15 USDT)
- ไม่มี martingale หรือ grid
- AI ตัดสินใจทุกอย่าง ไม่มี hard rule

## 📊 การติดตามผล
- Balance และ PnL realtime
- Position status ทุกชั่วโมง  
- AI confidence % ทุกการตัดสินใจ
- Pivot Point SuperTrend detection results

## ⚠️ ข้อจำกัดที่ทราบ
- AI response อาจช้าหรือ timeout (มี retry mechanism)
- Balance ต้องมีอย่างน้อย 15 USDT ต่อ position
- Position Value: 15 USDT x 5x leverage = 75 USDT total exposure

---

**หมายเหตุ**: ระบบนี้เป็นการเทรดที่มีความเสี่ยง โปรดใช้เงินที่พร้อมเสียได้และศึกษาทำความเข้าใจก่อนใช้งาน
