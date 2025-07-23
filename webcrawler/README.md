# Web Crawler for Email Extraction

โปรแกรม Go สำหรับ crawl เว็บไซต์ตามลิงก์และค้นหาอีเมลแล้วบันทึกลงไฟล์ JSON

## Features

- Crawl เว็บไซต์ตามลิงก์อย่างต่อเนื่อง
- ค้นหาและดึงอีเมลจากเนื้อหาเว็บไซต์
- บันทึกอีเมลที่พบลงไฟล์ JSON
- ข้ามเว็บไซต์ที่เคย crawl แล้ว (ไม่ crawl ซ้ำ)
- รองรับการกำหนดความลึกในการ crawl
- มีการ delay ระหว่างการ request เพื่อเคารพเซิร์ฟเวอร์

## Installation

1. ตรวจสอบให้แน่ใจว่ามี Go ติดตั้งแล้ว (version 1.21 หรือใหม่กว่า)
2. เข้าไปยังโฟลเดอร์ webcrawler
3. รันคำสั่ง:

\`\`\`bash
go mod tidy
\`\`\`

## Usage

### วิธีการใช้งานพื้นฐาน:

\`\`\`bash
go run main.go
\`\`\`

### ระบุ URL เริ่มต้นเอง:

\`\`\`bash
go run main.go https://example.com https://google.com
\`\`\`

### Build และรันไฟล์ executable:

\`\`\`bash
go build -o webcrawler.exe main.go
./webcrawler.exe
\`\`\`

## Configuration

แก้ไขไฟล์ `config.json` เพื่อปรับแต่งการทำงาน:

- `max_depth`: ความลึกสูงสุดในการ crawl (default: 3)
- `delay_seconds`: ระยะเวลา delay ระหว่างการ request (default: 2 วินาที)
- `output_file`: ชื่อไฟล์ JSON สำหรับบันทึกอีเมล (default: "found_emails.json")
- `start_urls`: URL เริ่มต้นสำหรับการ crawl
- `skip_extensions`: นามสกุลไฟล์ที่ต้องการข้าม

## Output

โปรแกรมจะสร้างไฟล์ `found_emails.json` ที่มีโครงสร้างดังนี้:

\`\`\`json
[
  {
    "email": "example@example.com",
    "url": "https://example.com/contact",
    "date": "2025-07-10 10:30:45"
  },
  {
    "email": "info@test.com",
    "url": "https://test.com/about",
    "date": "2025-07-10 10:31:20"
  }
]
\`\`\`

## ข้อมูลเพิ่มเติม

- โปรแกรมจะแสดงสถานะการทำงานทุก 30 วินาที
- การ crawl จะดำเนินการต่อไปในพื้นหลังจนกว่าจะหยุดโปรแกรม
- ใช้ Ctrl+C เพื่อหยุดโปรแกรม
- โปรแกรมจะเคารพเซิร์ฟเวอร์โดยการมี delay ระหว่างการ request
- รองรับการ crawl แบบ concurrent เพื่อประสิทธิภาพที่ดีขึ้น

## Dependencies

- `github.com/PuerkitoBio/goquery`: สำหรับ parse HTML
- `golang.org/x/net`: สำหรับ networking utilities

## Warning

โปรดใช้โปรแกรมนี้อย่างมีความรับผิดชอบ:
- อย่า crawl เว็บไซต์ที่ไม่ได้รับอนุญาต
- เคารพไฟล์ robots.txt ของเว็บไซต์
- ไม่ควรใช้ delay ที่เร็วเกินไปเพื่อไม่ให้สร้างภาระให้เซิร์ฟเวอร์
