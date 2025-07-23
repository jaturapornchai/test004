# 📋 Gate.io Trading Bot - Pivot Point SuperTrend with EMA100 and AI

## 🎯 ภาพรวมระบบ
ระบบเทรดอัตโนมัติที่ใช้ Pivot Point SuperTrend และ EMA100 ร่วมกับ AI ในการตัดสินใจเปิด/ปิด positions บน Gate.io Futures

## ⚙️ กฎการพัฒนาและบำรุงรักษา
- ใช้ Python เป็นหลัก (ลบ Golang files ทั้งหมด)
- กระจาย file python ใช้ import function ระหว่าง file ให้มากที่สุด
- ลบ file และ function ที่ไม่จำเป็นออกเสมอ
- การทดสอบทั้งหมดให้อยู่ใน folder test เท่านั้น
- ใช้ function จริงและข้อมูลจริงจาก Gate.io ในการทดสอบ (ไม่ใช้ mock/stub)
- ตรวจสอบความสมบูรณ์ของระบบก่อนทุกครั้ง ว่าไม่มี error
- ตอนจะทดสอบระบบ ให้ปิดทุกอย่างก่อน แล้วค่อยเปิดใหม่
- ข้อความ output (terminal) ให้เป็นภาษาไทยทั้งหมด
- ไม่ใช้ f-string formatting ที่ซับซ้อน ใช้ concatenate string ธรรมดา
- ห้ามแก้ไข folder markdown (เป็นต้นฉบับวิธีการทำงาน)
- ห้ามแก้ไข folder prompts (เป็น manual prompt files)
- ห้ามแก้ไขไฟล์ .env (config สำหรับการเชื่อมต่อ API)
- ไม่ต้องสร้าง .md file ใหม่ ใช้ไฟล์ markdown/step.md เป็นหลัก
- ใช้ try-except จัดการ error ทุกจุด ไม่ให้ระบบล่ม
- แสดง error message เป็นภาษาไทย
- ไม่ใช้ RSI ให้ตามลบ function และ class ที่เกี่ยวข้อง และ file ที่ไม่จำเป็นออก ทั้งหมด

## 📊 กฎสำคัญของระบบเทรด
### ❌ ข้อห้าม
- **ห้ามเปิด position ซ้ำในเหรียญเดียวกัน**
- **ห้ามเปิด position เพิ่ม หรือเปลี่ยนขนาด position** ถ้า position เปิดอยู่แล้ว
- **ห้ามตั้ง Stop Loss/Take Profit อัตโนมัติ** - เป็น AI-Only Mode

### ✅ ข้อบังคับ
- **ใช้ Pivot Point SuperTrend + EMA100 เท่านั้น** - ยกเลิก Chart Patterns อื่นทั้งหมด
- **Time Frame 1H เท่านั้น** - ไม่ใช้ timeframe อื่น
- **SuperTrend กำหนดทิศทาง** - Trend = 1 เล่นฝั่ง LONG, Trend = -1 เล่นฝั่ง SHORT
- **AUTO OPEN Mode**: SuperTrend Confidence ≥ 75% + Signal ≠ NEUTRAL = เปิด position อัตโนมัติ
- **AI Mode**: SuperTrend Confidence < 75% หรือ Signal = NEUTRAL ให้ AI ตัดสินใจ (ต้อง AI Confidence ≥ 85%)

## 🔧 Pivot Point SuperTrend Parameters
- **Pivot Point Period**: 2 periods (สำหรับหา Pivot High/Low)
- **ATR Factor**: 3.0 (สำหรับคำนวณ SuperTrend bands)
- **ATR Period**: 10 periods (สำหรับคำนวณ ATR)
- **EMA Period**: 100 periods (สำหรับยืนยันเทรนด์)
- **Time Frame**: 1H เท่านั้น

## 📊 Trading Signal Conditions

### 🟢 Long Position Strategy
**เงื่อนไข:** SuperTrend Trend = 1 (Bullish Trend)
1. **ติดตาม Trend**: เมื่อ SuperTrend แสดง Trend = 1 ให้เตรียมเล่นฝั่ง LONG
2. **เปิด Long**: เมื่อราคา rebound จาก EMA100 (แท่งเทียนสีเขียว Close > Open) 2 time frames อันไหนก็ได้

### 🔴 Short Position Strategy  
**เงื่อนไข:** SuperTrend Trend = -1 (Bearish Trend)
1. **ติดตาม Trend**: เมื่อ SuperTrend แสดง Trend = -1 ให้เตรียมเล่นฝั่ง SHORT
2. **เปิด Short**: เมื่อราคา reject จาก EMA100 (แท่งเทียนสีแดง Close < Open) 2 time frames อันไหนก็ได้

## 📈 Position Parameters
- **Position Size**: 15 USDT ต่อ position
- **Leverage**: 5x (ทดสอบแล้ว ✅ ทำงานได้สมบูรณ์)
- **Margin Mode**: isolated (ใช้ค่า default จาก Gate.io ✅)
- **Total Position Value**: 75 USDT (15 USDT x 5x leverage)
- **ต้องมี Balance**: อย่างน้อย 15 USDT

## 🔄 ขั้นตอนการทำงานของระบบ

### เริ่มต้นโปรแกรม
```bash
python app.py
```

### LOOP1: ตรวจสอบและจัดการ Positions (ทุก 1 ชั่วโมง)
1. ตรวจสอบ positions ที่เปิดอยู่ทั้งหมด
2. ดึง OHLCV 100 แท่งล่าสุด (1H) สำหรับแต่ละ position
3. ส่งข้อมูล position + OHLCV ให้ AI วิเคราะห์
4. AI ตัดสินใจ: CLOSE หรือ HOLD
5. ปิด position ถ้า AI แนะนำ CLOSE และ confidence ≥ 85%
6. ดึงรายชื่อเหรียญทั้งหมดจาก Gate.io Futures (เฉพาะ USDT pairs)

### LOOP2: สแกนหาโอกาสใหม่
1. วิเคราะห์เหรียญทีละตัวด้วย Pivot Point SuperTrend + EMA100
   - ดึง OHLCV 120 แท่ง (ใช้ 100 แท่งสุดท้ายวิเคราะห์)
   - คำนวณ Pivot Points (High/Low) ด้วย period 2
   - คำนวณ ATR ด้วย period 10
   - คำนวณ SuperTrend bands (Factor 3.0) และ Trend value
   - คำนวณ EMA100
   - ดูทิศทาง: Trend = 1 (LONG), Trend = -1 (SHORT)
   - คำนวณ Signal และ Confidence % จาก SuperTrend

2. **AUTO OPEN Mode** (SuperTrend Confidence ≥ 75% + Signal ≠ NEUTRAL):
   - เปิด position อัตโนมัติทันทีโดยไม่ต้องถาม AI
   - ใช้ Stop Loss/Take Profit ที่คำนวณจาก SuperTrend
   - ตรวจสอบ Risk-Reward Ratio ≥ 3.0
   - เปิด position ทันทีถ้าเงื่อนไขครบ

3. **AI Mode** (SuperTrend Confidence < 75% หรือ Signal = NEUTRAL):
   - ส่งข้อมูลให้ AI วิเคราะห์เพิ่มเติม
   - AI ตัดสินใจ: LONG/SHORT/HOLD พร้อม confidence %
   - ตรวจสอบ Risk-Reward Ratio ≥ 3.0 (AI คำนวณให้)
   - เปิด position ถ้า AI confidence ≥ 85%

4. เงื่อนไขการเปิด position (ทั้ง AUTO และ AI Mode):
   - ไม่มี position ซ้ำในเหรียญนั้น
   - มี balance เพียงพอ (≥ 15 USDT)
   - Risk-Reward Ratio ≥ 3.0

5. ใช้ batch processing (20 เหรียญ/batch, พัก 5 วินาที)
6. วนจนครบทุกเหรียญหรือ balance หมด
7. รอจนถึงชั่วโมงถัดไป แล้วกลับไป LOOP1

## 🤖 AI Integration

### Dual Mode System: AUTO OPEN + AI
ระบบมี 2 โหมดการทำงาน:

#### 🚀 AUTO OPEN Mode
- **เงื่อนไข**: SuperTrend Confidence ≥ 75% + Signal ≠ NEUTRAL
- **การทำงาน**: เปิด position อัตโนมัติทันทีโดยไม่ต้องถาม AI
- **ข้อดี**: ความเร็วสูง, ลด API calls, เหมาะกับสัญญาณที่ชัดเจน
- **Stop Loss/Take Profit**: คำนวณจาก SuperTrend bands และ ATR

#### 🤖 AI Mode
- **เงื่อนไข**: SuperTrend Confidence < 75% หรือ Signal = NEUTRAL
- **การทำงาน**: ส่งข้อมูลให้ AI วิเคราะห์เพิ่มเติม
- **ข้อดี**: การวิเคราะห์ลึก, เหมาะกับสัญญาณที่ไม่ชัดเจน
- **AI Confidence ต้อง ≥ 85%** ถึงจะเปิด position

### Manual Prompt Mode
- **ไม่สร้าง prompt อัตโนมัติ** - ใช้ manual prompt ที่เขียนไว้
- **Prompt files**:
  - `prompts/open_position_prompt.txt` - สำหรับเปิด position
  - `prompts/close_position_prompt.txt` - สำหรับปิด position
- **ภาษา**: Prompt ภาษาไทย, Response ภาษาไทย
- **ไม่ต้องชี้นำ AI มาก** - ให้ข้อมูลแล้วให้ AI วิเคราะห์เอง

### AI Response Handling
- Response อาจไม่เป็น JSON เสมอไป
- ใช้ regex หา JSON format ใน response
- มี fallback เมื่อ parse JSON ไม่ได้
- แสดง AI confidence แม้ action จะเป็น HOLD

### AI Decision Points
1. **AUTO OPEN Mode (SuperTrend Confidence ≥ 75% + Signal ≠ NEUTRAL)**:
   - รับข้อมูล: SuperTrend Trend + EMA100 + Signal + Confidence
   - ตัดสินใจ: LONG/SHORT ตาม Signal อัตโนมัติ
   - คำนวณ: Stop Loss/Take Profit จาก SuperTrend bands
   - ไม่ต้องใช้ AI

2. **AI Mode (SuperTrend Confidence < 75% หรือ Signal = NEUTRAL)**:
   - รับข้อมูล: OHLCV 100 แท่ง + SuperTrend Trend + EMA100 + ราคา proximity
   - ตัดสินใจ: LONG/SHORT/HOLD (ต้อง AI Confidence ≥ 85%)
   - คำนวณ: Risk-Reward Ratio
   - แนะนำ: SL/TP (แต่ระบบไม่ตั้งอัตโนมัติ)

3. **ปิด Position (ใช้ AI เท่านั้น)**:
   - รับข้อมูล: OHLCV 100 แท่ง + position data
   - ตัดสินใจ: CLOSE/HOLD (ต้อง AI Confidence ≥ 85%)
   - วิเคราะห์: โอกาสขาดทุน, รักษากำไร, หรือถือต่อ

## 📁 โครงสร้างไฟล์หลัก
1. **Core System**:
   - `app.py` - ไฟล์หลักของระบบ
   - `enhanced_position_manager.py` - จัดการ position และ workflow

2. **Analysis**:
   - `pivot_point_supertrend_detector.py` - Pivot Point SuperTrend Detection algorithm
   - `ai_analyzer.py` - AI integration สำหรับเปิด position
   - `position_analyzer.py` - AI integration สำหรับปิด position

3. **Data & Exchange**:
   - `exchange_client.py` - Gate.io connection via https://github.com/gateio/gateapi-python
   - `market_data_provider.py` - ดึงข้อมูลตลาด

4. **Configuration**:
   - `trading_rules.py` - กฎการเทรดและ parameters
   - `.env` - API keys (ห้ามแก้ไข)

## 📂 โฟลเดอร์สำคัญ
- `test/` - ไฟล์ทดสอบทั้งหมด
- `markdown/` - เอกสารต้นฉบับ (ห้ามแก้ไข)
- `prompts/` - Manual prompt files

## ⚠️ ข้อจำกัดที่ทราบ
- **Leverage Setting**: ✅ ทำงานได้สมบูรณ์ (ทดสอบแล้ว)
- **API Signature Error**: เกิดขึ้นบางครั้งตอนตั้ง margin mode แต่ไม่กระทบการทำงาน
- **AI response**: อาจช้าหรือ timeout (มี retry mechanism)
- **Balance Requirement**: ต้องมี balance อย่างน้อย 15 USDT ต่อ position
- **Position Value**: 15 USDT x 5x leverage = 75 USDT total exposure

## 📊 การติดตามผล
- **Terminal Output**:
  - Balance และ PnL realtime
  - Position status ทุกชั่วโมง
  - AI confidence % ทุกการตัดสินใจ
  - Pivot Point SuperTrend detection results
- **Log Files**: ไม่มี (ดูจาก terminal เท่านั้น)
- **Position Tracking**: แสดง entry price, size, PnL

## 🔐 ความปลอดภัย
- ใช้ isolated margin ลดความเสี่ยง
- Position size เล็ก (15 USDT)
- ไม่มี martingale หรือ grid
- AI ตัดสินใจทุกอย่าง ไม่มี hard rule

# Pine Script
// This source code is subject to the terms of the Mozilla Public License 2.0 at https://mozilla.org/MPL/2.0/
// © LonesomeTheBlue

//@version=4
study("Pivot Point SuperTrend", overlay = true)
prd = input(defval = 2, title="Pivot Point Period", minval = 1, maxval = 50)
Factor=input(defval = 3, title = "ATR Factor", minval = 1, step = 0.1)
Pd=input(defval = 10, title = "ATR Period", minval=1)
showpivot = input(defval = false, title="Show Pivot Points")
showlabel = input(defval = true, title="Show Buy/Sell Labels")
showcl = input(defval = false, title="Show PP Center Line")
showsr = input(defval = false, title="Show Support/Resistance")

// get Pivot High/Low
float ph = pivothigh(prd, prd)
float pl = pivotlow(prd, prd)

// drawl Pivot Points if "showpivot" is enabled
plotshape(ph and showpivot, text="H",  style=shape.labeldown, color=na, textcolor=color.red, location=location.abovebar, transp=0, offset = -prd)
plotshape(pl and showpivot, text="L",  style=shape.labeldown, color=na, textcolor=color.lime, location=location.belowbar, transp=0, offset = -prd)

// calculate the Center line using pivot points
var float center = na
float lastpp = ph ? ph : pl ? pl : na
if lastpp
    if na(center)
        center := lastpp
    else
        //weighted calculation
        center := (center * 2 + lastpp) / 3

// upper/lower bands calculation
Up = center - (Factor * atr(Pd))
Dn = center + (Factor * atr(Pd))

// get the trend
float TUp = na
float TDown = na
Trend = 0
TUp := close[1] > TUp[1] ? max(Up, TUp[1]) : Up
TDown := close[1] < TDown[1] ? min(Dn, TDown[1]) : Dn
Trend := close > TDown[1] ? 1: close < TUp[1]? -1: nz(Trend[1], 1)
Trailingsl = Trend == 1 ? TUp : TDown

// plot the trend
linecolor = Trend == 1 and nz(Trend[1]) == 1 ? color.lime : Trend == -1 and nz(Trend[1]) == -1 ? color.red : na
plot(Trailingsl, color = linecolor ,  linewidth = 2, title = "PP SuperTrend")
 
plot(showcl ? center : na, color = showcl ? center < hl2 ? color.blue : color.red : na)

// check and plot the signals
bsignal = Trend == 1 and Trend[1] == -1
ssignal = Trend == -1 and Trend[1] == 1
plotshape(bsignal and showlabel ? Trailingsl : na, title="Buy", text="Buy", location = location.absolute, style = shape.labelup, size = size.tiny, color = color.lime, textcolor = color.black, transp = 0)
plotshape(ssignal and showlabel ? Trailingsl : na, title="Sell", text="Sell", location = location.absolute, style = shape.labeldown, size = size.tiny, color = color.red, textcolor = color.white, transp = 0)

//get S/R levels using Pivot Points
float resistance = na
float support = na
support := pl ? pl : support[1]
resistance := ph ? ph : resistance[1]

// if enabled then show S/R levels
plot(showsr and support ? support : na, color = showsr and support ? color.lime : na, style = plot.style_circles, offset = -prd)
plot(showsr and resistance ? resistance : na, color = showsr and resistance ? color.red : na, style = plot.style_circles, offset = -prd)

// alerts
alertcondition(Trend == 1 and Trend[1] == -1, title='Buy Signal', message='Buy Signal')
alertcondition(Trend == -1 and Trend[1] == 1, title='Sell Signal', message='Sell Signal')
alertcondition(change(Trend), title='Trend Changed', message='Trend Changed')
