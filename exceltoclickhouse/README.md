# Excel to ClickHouse Data Manager

🎯 **ระบบจัดการข้อมูลจาก Excel ไปยัง ClickHouse** สำหรับธุรกิจสปา/นวด

## 📋 โครงสร้างโปรเจค

```
exceltoclickhouse/
├── main.py                     # สคริปต์หลัก
├── create_tables.py            # สร้างตารางใน ClickHouse
├── upload_data.py              # อัพโหลดข้อมูลจาก Excel
├── clickhouse_connection.py    # ทดสอบการเชื่อมต่อ
├── sample.xlsx                 # ไฟล์ Excel ตัวอย่าง
└── README.md                   # คู่มือใช้งาน
```

## 🏗️ โครงสร้างตาราง (11 ตาราง)

### 📊 ตารางหลัก
1. **bookings** - ข้อมูลการจองหลัก (customer, date, channel)
2. **services** - ข้อมูลบริการ (service items, duration, price)
3. **payments** - ข้อมูลการชำระเงิน (payment methods, amounts)
4. **staff_assignments** - การจัดพนักงาน (staff, fees, tips)

### 📈 ตารางรายงาน
5. **staff_wages** - ค่ามือพนักงาน (wages by staff and date)
6. **daily_summary** - สรุปยอดขายรายวัน (payment summary)
7. **service_sales** - ยอดขายตามรายการบริการ
8. **product_sales** - ยอดขายสินค้า

### 🏪 ตารางจัดการ
9. **cashier_transactions** - รายการแคชเชียร์ (income/expense)
10. **inventory** - สต๊อกสินค้า (stock movements)
11. **rejected_customers** - ลูกค้าที่ปฏิเสธ (rejection reasons)

## 🚀 การใช้งาน

### 1. สร้างตารางใหม่
```bash
python main.py --create-tables
```

### 2. อัพโหลดข้อมูลจาก sample.xlsx
```bash
python main.py --upload-data
```

### 3. อัพโหลดข้อมูลจากไฟล์อื่น
```bash
python main.py --upload-only new_data.xlsx
```

### 4. รันทั้งสองขั้นตอน (สร้าง + อัพโหลด)
```bash
python main.py --all
```

### 5. ดูข้อมูลตาราง
```bash
python main.py --info
```

## 🔧 คุณสมบัติพิเศษ

✅ **ชื่อฟิลด์ภาษาอังกฤษ** - ใช้งานง่าย สำหรับ SQL Query
✅ **Comment ภาษาไทย** - เข้าใจง่าย อธิบายข้อมูล
✅ **แยกตารางตามหน้าที่** - ออกแบบตาม Normalization
✅ **รองรับข้อมูลใหม่** - อัพโหลดไฟล์ Excel ใหม่ได้ตลอด
✅ **Data Type Optimization** - ใช้ประเภทข้อมูลที่เหมาะสม

## 💡 ตัวอย่าง SQL Query

### จำนวนการจองทั้งหมด
```sql
SELECT COUNT(*) as total_bookings FROM bookings;
```

### ยอดขายรายวัน
```sql
SELECT 
    booking_date,
    COUNT(*) as bookings_count,
    SUM(net_payment) as total_revenue
FROM bookings b
JOIN payments p ON b.booking_id = p.booking_id
GROUP BY booking_date
ORDER BY booking_date;
```

### บริการยอดนิยม
```sql
SELECT 
    service_item,
    COUNT(*) as usage_count,
    AVG(service_price) as avg_price
FROM services 
WHERE service_item IS NOT NULL
GROUP BY service_item
ORDER BY usage_count DESC
LIMIT 10;
```

### ผลงานพนักงาน
```sql
SELECT 
    staff_name,
    COUNT(*) as services_count,
    SUM(staff_fee) as total_fee,
    SUM(staff_tip) as total_tip
FROM staff_assignments
WHERE staff_name IS NOT NULL
GROUP BY staff_name
ORDER BY total_fee DESC;
```

### สรุปช่องทางการจอง
```sql
SELECT 
    booking_channel,
    COUNT(*) as bookings,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
FROM bookings
WHERE booking_channel IS NOT NULL
GROUP BY booking_channel
ORDER BY bookings DESC;
```

## 🔗 การตั้งค่าการเชื่อมต่อ

แก้ไขข้อมูลการเชื่อมต่อใน `create_tables.py` และ `upload_data.py`:

```python
def get_clickhouse_client():
    return clickhouse_connect.get_client(
        host='your-clickhouse-host',
        user='your-username', 
        password='your-password',
        secure=True
    )
```

## 📦 Dependencies

```bash
pip install pandas openpyxl clickhouse-connect
```

## 🎯 วิธีเพิ่มข้อมูลใหม่

1. วาง Excel file ใหม่ในโฟลเดอร์เดียวกัน
2. รัน: `python main.py --upload-only your_new_file.xlsx`
3. ข้อมูลจะถูกเพิ่มเข้าตารางที่มีอยู่

## ⚠️ หมายเหตุ

- โครงสร้างตารางจะถูกสร้างใหม่ทุกครั้งที่รัน `--create-tables`
- การอัพโหลดข้อมูลจะเพิ่มข้อมูลใหม่เข้าตารางที่มีอยู่
- รองรับไฟล์ Excel ที่มี 8 sheet ตามโครงสร้างของ sample.xlsx

---

📊 **สร้างโดย**: AI Assistant  
🔧 **เวอร์ชัน**: 1.0  
📅 **วันที่**: July 2025
