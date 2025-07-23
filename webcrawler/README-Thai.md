# Thai Web Crawler for Email Extraction

โปรแกรม Go สำหรับ crawl เว็บไซต์ไทยและค้นหาอีเมลแล้วบันทึกลงไฟล์ JSON

## คุณสมบัติ

- **เฉพาะเว็บไซต์ไทย**: crawl เฉพาะเว็บไซต์ที่มีโดเมน .th และเว็บไซต์ไทยที่รู้จัก
- **ตรวจสอบเนื้อหาภาษาไทย**: ข้ามหน้าที่ไม่มีเนื้อหาภาษาไทย
- **ค้นหาอีเมล**: ดึงอีเมลจากเนื้อหาเว็บไซต์ไทย
- **บันทึก JSON**: เก็บอีเมลที่พบลงไฟล์ JSON
- **ไม่ crawl ซ้ำ**: ข้ามเว็บไซต์ที่เคย crawl แล้ว
- **ปลอดภัย**: มีการ delay ระหว่างการ request เพื่อเคารพเซิร์ฟเวอร์

## การติดตั้ง

1. ตรวจสอบให้แน่ใจว่ามี Go ติดตั้งแล้ว (version 1.21+)
2. เข้าไปยังโฟลเดอร์ webcrawler
3. รันคำสั่ง:

```bash
go mod tidy
go build -o webcrawler-thai.exe main.go
```

## การใช้งาน

### รันโปรแกรม

```bash
# รันด้วย URL เริ่มต้น (เว็บไซต์ไทยที่กำหนดไว้)
go run main.go

# หรือใช้ไฟล์ที่ build แล้ว
./webcrawler-thai.exe
```

### ระบุ URL เว็บไซต์ไทยเอง

```bash
go run main.go https://pantip.com https://sanook.com
```

## เว็บไซต์ไทยที่รองรับ

### โดเมนไทย
- .th (ทุกโดเมน .th)
- .co.th
- .ac.th  
- .go.th
- .or.th
- .net.th

### เว็บไซต์ไทยที่นิยม
- pantip.com
- sanook.com
- kapook.com
- mthai.com
- thairath.co.th
- manager.co.th
- siamzone.com
- dek-d.com
- thaipbs.or.th
- และอีกมากมาย

## การตั้งค่า

แก้ไขไฟล์ `config.json`:

- `max_depth`: ความลึกสูงสุดในการ crawl (default: 3)
- `delay_seconds`: ระยะเวลา delay ระหว่างการ request (default: 2 วินาที)
- `min_thai_ratio`: สัดส่วนเนื้อหาภาษาไทยขั้นต่ำ (default: 0.2 = 20%)
- `min_thai_chars`: จำนวนตัวอักษรไทยขั้นต่ำ (default: 10)

## ผลลัพธ์

โปรแกรมจะสร้างไฟล์ `found_emails.json`:

```json
[
  {
    "email": "info@thaiwebsite.co.th",
    "url": "https://thaiwebsite.co.th/contact",
    "date": "2025-07-10 10:30:45"
  },
  {
    "email": "admin@pantip.com",
    "url": "https://pantip.com/about",
    "date": "2025-07-10 10:31:20"
  }
]
```

## ข้อมูลเพิ่มเติม

- โปรแกรมจะแสดงสถานะทุก 30 วินาที
- ข้ามหน้าที่ไม่มีเนื้อหาภาษาไทยเพียงพอ
- ใช้ Ctrl+C เพื่อหยุดโปรแกรม
- รองรับการ crawl แบบ concurrent
- เคารพเซิร์ฟเวอร์ด้วยการ delay

## ข้อควรระวัง

⚠️ **ใช้อย่างมีความรับผิดชอบ**:
- อย่า crawl เว็บไซต์ที่ไม่ได้รับอนุญาต
- เคารพไฟล์ robots.txt ของเว็บไซต์
- ไม่ควรใช้ delay ที่เร็วเกินไป
- ปฏิบัติตามกฎหมายและจริยธรรม

## ตัวอย่างการรัน

```bash
# รันครั้งแรก
go run main.go

# ผลลัพธ์
Web Crawler Started!
Press Ctrl+C to stop the crawler
Starting web crawler with 8 initial URLs
Max depth: 3, Delay: 2s
Output file: found_emails.json
Crawling: https://pantip.com (depth: 0)
Found 3 emails on https://pantip.com
Crawling: https://pantip.com/forum/... (depth: 1)
Skipping https://example.com - No Thai content detected
Status: Found 15 emails, Visited 25 URLs
```
