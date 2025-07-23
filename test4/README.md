# Product API

REST API สำหรับการจัดการข้อมูลสินค้า เขียนด้วยภาษา Go และ Gin framework

## คุณสมบัติ

- ✅ CRUD operations สำหรับสินค้า
- ✅ การกรองสินค้าตามหมวดหมู่และราคา
- ✅ จัดการข้อมูลสินค้าในหน่วยความจำ
- ✅ Mock data พร้อมใช้งาน
- ✅ CORS support
- ✅ JSON response format

## การติดตั้ง

1. ติดตั้ง dependencies:
```bash
go mod tidy
```

2. รันเซิร์ฟเวอร์:
```bash
go run main.go
```

เซิร์ฟเวอร์จะทำงานที่ `http://localhost:8080`

## API Endpoints

### สินค้า

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/products` | ดึงรายการสินค้าทั้งหมด |
| GET | `/api/products/:id` | ดึงข้อมูลสินค้าตาม ID |
| POST | `/api/products` | สร้างสินค้าใหม่ |
| PUT | `/api/products/:id` | แก้ไขข้อมูลสินค้า |
| DELETE | `/api/products/:id` | ลบสินค้า |

### อื่นๆ

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/categories` | ดึงรายการหมวดหมู่สินค้า |
| GET | `/api/health` | ตรวจสอบสถานะเซิร์ฟเวอร์ |
| GET | `/` | หน้าแรก API |

## ตัวอย่างการใช้งาน

### 1. ดึงสินค้าทั้งหมด
```bash
curl http://localhost:8080/api/products
```

### 2. ดึงสินค้าตามหมวดหมู่
```bash
curl "http://localhost:8080/api/products?category=Electronics"
```

### 3. ดึงสินค้าตามช่วงราคา
```bash
curl "http://localhost:8080/api/products?min_price=1000&max_price=10000"
```

### 4. สร้างสินค้าใหม่
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "สินค้าใหม่",
    "description": "รายละเอียดสินค้า",
    "price": 999.99,
    "category": "Test",
    "stock": 10,
    "image_url": "https://example.com/image.jpg"
  }'
```

### 5. แก้ไขสินค้า
```bash
curl -X PUT http://localhost:8080/api/products/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "สินค้าที่แก้ไขแล้ว",
    "description": "รายละเอียดใหม่",
    "price": 1299.99,
    "category": "Updated",
    "stock": 15,
    "image_url": "https://example.com/updated-image.jpg"
  }'
```

### 6. ลบสินค้า
```bash
curl -X DELETE http://localhost:8080/api/products/{id}
```

## โครงสร้างข้อมูล

### Product Object
```json
{
  "id": "uuid-string",
  "name": "ชื่อสินค้า",
  "description": "รายละเอียดสินค้า",
  "price": 999.99,
  "category": "หมวดหมู่",
  "stock": 10,
  "image_url": "https://example.com/image.jpg",
  "created_at": "2025-06-24T10:00:00Z",
  "updated_at": "2025-06-24T10:00:00Z"
}
```

### Response Format
```json
{
  "success": true,
  "message": "ข้อความตอบกลับ",
  "data": {}
}
```

## Mock Data

API มาพร้อมกับข้อมูลตัวอย่าง:
- iPhone 15 Pro
- MacBook Air M2
- AirPods Pro 2
- Nike Air Max 270
- เสื้อยืดผ้าฝ้าย 100%
- กาแฟอาราบิก้าคั่วกลาง

## การพัฒนาต่อ

- เชื่อมต่อกับฐานข้อมูล (PostgreSQL, MySQL)
- เพิ่มระบบ Authentication/Authorization
- เพิ่มการอัพโหลดรูปภาพ
- เพิ่มระบบค้นหาแบบ full-text search
- เพิ่ม pagination สำหรับรายการสินค้า
- เพิ่ม validation และ error handling ที่ดีขึ้น
