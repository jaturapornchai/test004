"""
🎯 สคริปต์หลักสำหรับจัดการข้อมูล Excel ไปยัง ClickHouse

การใช้งาน:
1. รันเพื่อสร้างตารางใหม่: python main.py --create-tables
2. รันเพื่ออัพโหลดข้อมูล: python main.py --upload-data  
3. รันทั้งสองอย่าง: python main.py --all
4. รันเฉพาะไฟล์ Excel ใหม่: python main.py --upload-only new_file.xlsx
"""

import sys
import argparse
import subprocess
import os

def run_create_tables():
    """
    เรียกใช้สคริปต์สร้างตาราง
    """
    print("🏗️  เริ่มสร้างตารางใน ClickHouse...")
    result = subprocess.run([sys.executable, "create_tables.py"], capture_output=True, text=True)
    
    if result.returncode == 0:
        print("✅ สร้างตารางสำเร็จ!")
        print(result.stdout)
    else:
        print("❌ เกิดข้อผิดพลาดในการสร้างตาราง:")
        print(result.stderr)
    
    return result.returncode == 0

def run_upload_data(excel_file="sample.xlsx"):
    """
    เรียกใช้สคริปต์อัพโหลดข้อมูล
    """
    if not os.path.exists(excel_file):
        print(f"❌ ไม่พบไฟล์ {excel_file}")
        return False
        
    print(f"📤 เริ่มอัพโหลดข้อมูลจาก {excel_file}...")
    
    # ถ้าเป็นไฟล์ที่ไม่ใช่ sample.xlsx ให้แก้ไขชื่อในสคริปต์ชั่วคราว
    if excel_file != "sample.xlsx":
        # สร้างสคริปต์ชั่วคราว
        with open("upload_data.py", "r", encoding="utf-8") as f:
            content = f.read()
        
        temp_content = content.replace('excel_file = "sample.xlsx"', f'excel_file = "{excel_file}"')
        
        with open("upload_data_temp.py", "w", encoding="utf-8") as f:
            f.write(temp_content)
        
        result = subprocess.run([sys.executable, "upload_data_temp.py"], capture_output=True, text=True)
        
        # ลบไฟล์ชั่วคราว
        if os.path.exists("upload_data_temp.py"):
            os.remove("upload_data_temp.py")
    else:
        result = subprocess.run([sys.executable, "upload_data.py"], capture_output=True, text=True)
    
    if result.returncode == 0:
        print("✅ อัพโหลดข้อมูลสำเร็จ!")
        print(result.stdout)
    else:
        print("❌ เกิดข้อผิดพลาดในการอัพโหลดข้อมูล:")
        print(result.stderr)
    
    return result.returncode == 0

def show_table_info():
    """
    แสดงข้อมูลตารางที่มีใน ClickHouse
    """
    print("\n📊 ข้อมูลตารางใน ClickHouse:")
    print("=" * 60)
    
    tables = [
        "bookings - ข้อมูลการจองหลัก",
        "services - ข้อมูลบริการ", 
        "payments - ข้อมูลการชำระเงิน",
        "staff_assignments - ข้อมูลการจัดพนักงาน",
        "staff_wages - ข้อมูลค่ามือพนักงาน",
        "daily_summary - สรุปยอดขายรายวัน",
        "service_sales - ยอดขายตามรายการบริการ",
        "product_sales - ยอดขายสินค้า",
        "cashier_transactions - รายการแคชเชียร์",
        "inventory - สต๊อกสินค้า",
        "rejected_customers - ลูกค้าที่ปฏิเสธ"
    ]
    
    for i, table in enumerate(tables, 1):
        print(f"  {i:2d}. {table}")

def main():
    parser = argparse.ArgumentParser(description="จัดการข้อมูล Excel ไปยัง ClickHouse")
    
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--create-tables", action="store_true", help="สร้างตารางใหม่ใน ClickHouse")
    group.add_argument("--upload-data", action="store_true", help="อัพโหลดข้อมูลจาก sample.xlsx")
    group.add_argument("--upload-only", metavar="FILE", help="อัพโหลดข้อมูลจากไฟล์ Excel ที่กำหนด")
    group.add_argument("--all", action="store_true", help="สร้างตารางและอัพโหลดข้อมูล")
    group.add_argument("--info", action="store_true", help="แสดงข้อมูลตาราง")
    
    args = parser.parse_args()
    
    print("🚀 Excel to ClickHouse Data Manager")
    print("=" * 60)
    
    success = True
    
    if args.info:
        show_table_info()
    elif args.create_tables:
        success = run_create_tables()
    elif args.upload_data:
        success = run_upload_data()
    elif args.upload_only:
        success = run_upload_data(args.upload_only)
    elif args.all:
        print("🔄 รันทั้งสองขั้นตอน: สร้างตาราง + อัพโหลดข้อมูล")
        success = run_create_tables()
        if success:
            print("\n" + "="*60)
            success = run_upload_data()
    
    print("=" * 60)
    if success:
        print("🎉 ดำเนินการเสร็จสิ้น!")
        
        if args.all or args.upload_data or args.upload_only:
            print("\n💡 ตัวอย่างการใช้งาน SQL:")
            print("SELECT COUNT(*) FROM bookings;")
            print("SELECT booking_date, COUNT(*) FROM bookings GROUP BY booking_date;")
            print("SELECT service_item, COUNT(*) FROM services GROUP BY service_item LIMIT 10;")
    else:
        print("❌ ดำเนินการไม่สำเร็จ!")
        sys.exit(1)

if __name__ == "__main__":
    main()
