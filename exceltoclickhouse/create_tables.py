import clickhouse_connect
from datetime import datetime

def get_clickhouse_client():
    """
    เชื่อมต่อ ClickHouse
    """
    return clickhouse_connect.get_client(
        host='npomobbg93.germanywestcentral.azure.clickhouse.cloud',
        user='default',
        password='1S.6V_z9Lr9Wc',
        secure=True
    )

def create_bookings_table(client):
    """
    สร้างตาราง bookings - ข้อมูลการจองหลัก (จาก sheet ภาพรวม)
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS bookings (
        booking_id UInt32 COMMENT 'รหัสการจอง',
        booking_date Nullable(String) COMMENT 'วันที่จอง',
        booking_time Nullable(String) COMMENT 'เวลาจอง', 
        datetime_display Nullable(String) COMMENT 'วัน-เวลาแสดงผล',
        receipt_number Nullable(String) COMMENT 'เลขใบเสร็จ',
        booking_channel Nullable(String) COMMENT 'ช่องทางการจอง',
        agency_id Nullable(Float64) COMMENT 'รหัสเอเจนซี่',
        booking_type Nullable(String) COMMENT 'ประเภทการจอง (treatment/package)',
        customer_id Nullable(Float64) COMMENT 'รหัสลูกค้า', 
        customer_name Nullable(String) COMMENT 'ชื่อลูกค้า',
        phone_number Nullable(String) COMMENT 'เบอร์โทรศัพท์',
        booking_note Nullable(String) COMMENT 'หมายเหตุการจอง',
        created_by Nullable(String) COMMENT 'ผู้ทำรายการ',
        created_at DateTime DEFAULT now() COMMENT 'วันที่สร้างข้อมูล'
    ) ENGINE = MergeTree() 
    ORDER BY booking_id
    COMMENT 'ตารางข้อมูลการจองหลัก';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS bookings")
        client.command(create_sql)
        print("✅ สร้างตาราง bookings สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง bookings ได้: {e}")
        return False

def create_services_table(client):
    """
    สร้างตาราง services - ข้อมูลบริการ
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS services (
        booking_id UInt32 COMMENT 'รหัสการจอง',
        service_item Nullable(String) COMMENT 'รายการบริการ',
        service_duration Nullable(Float64) COMMENT 'ระยะเวลาบริการ (นาที)',
        service_price Nullable(UInt16) COMMENT 'ราคาต่อบริการ',
        service_discount Nullable(Float64) COMMENT 'ส่วนลดต่อบริการ',
        total_items Nullable(Float64) COMMENT 'จำนวนรายการทั้งหมด',
        service_cost Nullable(Float64) COMMENT 'ต้นทุนบริการ'
    ) ENGINE = MergeTree()
    ORDER BY booking_id
    COMMENT 'ตารางข้อมูลบริการ';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS services")
        client.command(create_sql)
        print("✅ สร้างตาราง services สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง services ได้: {e}")
        return False

def create_payments_table(client):
    """
    สร้างตาราง payments - ข้อมูลการชำระเงิน
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS payments (
        booking_id UInt32 COMMENT 'รหัสการจอง',
        total_price Nullable(Float64) COMMENT 'ราคารวม',
        total_discount Nullable(Float64) COMMENT 'ส่วนลดรวม', 
        price_after_discount Nullable(Float64) COMMENT 'ราคาหลังหักส่วนลด',
        service_charge Nullable(Float64) COMMENT 'ค่าบริการ',
        additional_charge Nullable(Float64) COMMENT 'ค่าธรรมเนียมเพิ่มเติม',
        price_before_vat Nullable(String) COMMENT 'ราคาก่อน VAT',
        vat_amount Nullable(String) COMMENT 'จำนวน VAT 7%',
        total_tip Nullable(Float64) COMMENT 'ทิปรวม',
        agency_commission Nullable(UInt8) COMMENT 'ค่าคอมมิชชันเอเจนซี',
        net_payment Nullable(Float64) COMMENT 'ยอดชำระสุทธิ',
        cash_payment Nullable(Float64) COMMENT 'ชำระด้วยเงินสด',
        transfer_payment Nullable(Float64) COMMENT 'ชำระด้วยการโอน',
        credit_card_payment Nullable(Float64) COMMENT 'ชำระด้วยบัตรเครดิต',
        agency_payment Nullable(Float64) COMMENT 'ชำระผ่านเอเจนซี',
        package_payment Nullable(Float64) COMMENT 'ชำระด้วยแพ็กเกจ',
        ewallet_payment Nullable(Float64) COMMENT 'ชำระด้วย E-Wallet',
        ewallet_type Nullable(String) COMMENT 'ประเภท E-Wallet',
        payment_note Nullable(String) COMMENT 'หมายเหตุการชำระเงิน'
    ) ENGINE = MergeTree()
    ORDER BY booking_id
    COMMENT 'ตารางข้อมูลการชำระเงิน';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS payments")
        client.command(create_sql)
        print("✅ สร้างตาราง payments สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง payments ได้: {e}")
        return False

def create_staff_assignments_table(client):
    """
    สร้างตาราง staff_assignments - ข้อมูลการจัดพนักงาน
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS staff_assignments (
        booking_id UInt32 COMMENT 'รหัสการจอง',
        staff_sequence UInt8 COMMENT 'ลำดับพนักงาน (1, 2)',
        staff_name Nullable(String) COMMENT 'ชื่อพนักงาน',
        staff_fee Nullable(Float64) COMMENT 'ค่ามือหมอนวด',
        staff_tip Nullable(Float64) COMMENT 'ทิปหมอนวด',
        request_fee Nullable(Float64) COMMENT 'ค่ารีเควส'
    ) ENGINE = MergeTree()
    ORDER BY (booking_id, staff_sequence)
    COMMENT 'ตารางข้อมูลการจัดพนักงาน';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS staff_assignments")
        client.command(create_sql)
        print("✅ สร้างตาราง staff_assignments สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง staff_assignments ได้: {e}")
        return False

def create_staff_wages_table(client):
    """
    สร้างตาราง staff_wages - ข้อมูลค่ามือพนักงาน (จาก sheet ค่ามือ)
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS staff_wages (
        date_record String DEFAULT '' COMMENT 'วันที่บันทึก',
        record_type String DEFAULT '' COMMENT 'ประเภทรายการ (ค่ามือ/ประกันมือ, ค่ารีเควส)',
        staff_name String COMMENT 'ชื่อพนักงาน',
        amount Nullable(Float64) COMMENT 'จำนวนเงิน'
    ) ENGINE = MergeTree()
    ORDER BY (date_record, staff_name)
    COMMENT 'ตารางข้อมูลค่ามือพนักงาน';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS staff_wages")
        client.command(create_sql)
        print("✅ สร้างตาราง staff_wages สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง staff_wages ได้: {e}")
        return False

def create_daily_summary_table(client):
    """
    สร้างตาราง daily_summary - สรุปยอดขายรายวัน (จาก sheet ยอดเงิน)
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS daily_summary (
        summary_date String DEFAULT '' COMMENT 'วันที่สรุป',
        payment_type String COMMENT 'ประเภทการชำระเงิน',
        total_amount Nullable(Float64) COMMENT 'ยอดรวม'
    ) ENGINE = MergeTree()
    ORDER BY payment_type
    COMMENT 'ตารางสรุปยอดขายรายวัน';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS daily_summary")
        client.command(create_sql)
        print("✅ สร้างตาราง daily_summary สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง daily_summary ได้: {e}")
        return False

def create_service_sales_table(client):
    """
    สร้างตาราง service_sales - ยอดขายตามรายการบริการ
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS service_sales (
        service_group String DEFAULT '' COMMENT 'กลุ่มบริการ',
        service_name String COMMENT 'ชื่อบริการ',
        quantity Nullable(UInt32) COMMENT 'จำนวนที่ขาย',
        sales_amount Nullable(Float64) COMMENT 'ยอดขาย'
    ) ENGINE = MergeTree()
    ORDER BY service_name
    COMMENT 'ตารางยอดขายตามรายการบริการ';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS service_sales")
        client.command(create_sql)
        print("✅ สร้างตาราง service_sales สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง service_sales ได้: {e}")
        return False

def create_product_sales_table(client):
    """
    สร้างตาราง product_sales - ยอดขายสินค้า
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS product_sales (
        product_group String DEFAULT '' COMMENT 'กลุ่มสินค้า',
        product_name String COMMENT 'ชื่อสินค้า',
        quantity Nullable(UInt32) COMMENT 'จำนวนที่ขาย',
        sales_amount Nullable(Float64) COMMENT 'ยอดขาย'
    ) ENGINE = MergeTree()
    ORDER BY product_name
    COMMENT 'ตารางยอดขายสินค้า';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS product_sales")
        client.command(create_sql)
        print("✅ สร้างตาราง product_sales สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง product_sales ได้: {e}")
        return False

def create_cashier_transactions_table(client):
    """
    สร้างตาราง cashier_transactions - รายการแคชเชียร์
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS cashier_transactions (
        transaction_date String DEFAULT '' COMMENT 'วันที่ทำรายการ',
        description Nullable(String) COMMENT 'รายละเอียด',
        transaction_type Nullable(String) COMMENT 'ประเภทรายการ (รายรับ/รายจ่าย)',
        amount Nullable(Float64) COMMENT 'จำนวนเงิน'
    ) ENGINE = MergeTree()
    ORDER BY transaction_date
    COMMENT 'ตารางรายการแคชเชียร์';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS cashier_transactions")
        client.command(create_sql)
        print("✅ สร้างตาราง cashier_transactions สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง cashier_transactions ได้: {e}")
        return False

def create_inventory_table(client):
    """
    สร้างตาราง inventory - สต๊อกสินค้า
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS inventory (
        inventory_date String DEFAULT '' COMMENT 'วันที่บันทึก',
        product_group Nullable(String) COMMENT 'กลุ่มสินค้า',
        product_name Nullable(String) COMMENT 'ชื่อสินค้า',
        sold_qty Nullable(Float64) COMMENT 'จำนวนที่ขาย',
        added_qty Nullable(Float64) COMMENT 'จำนวนที่เติม',
        reduced_qty Nullable(Float64) COMMENT 'จำนวนที่ลด',
        returned_qty Nullable(Float64) COMMENT 'จำนวนที่ดึงกลับ',
        remaining_qty Nullable(Float64) COMMENT 'จำนวนคงเหลือ'
    ) ENGINE = MergeTree()
    ORDER BY inventory_date
    COMMENT 'ตารางสต๊อกสินค้า';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS inventory")
        client.command(create_sql)
        print("✅ สร้างตาราง inventory สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง inventory ได้: {e}")
        return False

def create_rejected_customers_table(client):
    """
    สร้างตาราง rejected_customers - ลูกค้าที่ปฏิเสธ
    """
    create_sql = """
    CREATE TABLE IF NOT EXISTS rejected_customers (
        datetime_rejected String DEFAULT '' COMMENT 'วัน-เวลาที่ปฏิเสธ',
        booking_date Nullable(String) COMMENT 'วันที่จอง',
        booking_time Nullable(String) COMMENT 'เวลาที่จอง',
        booking_channel Nullable(String) COMMENT 'ช่องทางการจอง',
        customer_count Nullable(UInt8) COMMENT 'จำนวนลูกค้า',
        rejection_reason Nullable(String) COMMENT 'เหตุผลการปฏิเสธ',
        customer_name Nullable(String) COMMENT 'ชื่อลูกค้า',
        phone_number Nullable(String) COMMENT 'เบอร์โทรศัพท์',
        booking_note Nullable(String) COMMENT 'หมายเหตุการจอง'
    ) ENGINE = MergeTree()
    ORDER BY datetime_rejected
    COMMENT 'ตารางลูกค้าที่ปฏิเสธ';
    """
    
    try:
        client.command("DROP TABLE IF EXISTS rejected_customers")
        client.command(create_sql)
        print("✅ สร้างตาราง rejected_customers สำเร็จ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถสร้างตาราง rejected_customers ได้: {e}")
        return False

def main():
    """
    ฟังก์ชันหลักสำหรับสร้างตารางทั้งหมด
    """
    print("🔗 กำลังเชื่อมต่อ ClickHouse...")
    
    try:
        client = get_clickhouse_client()
        print("✅ เชื่อมต่อ ClickHouse สำเร็จ!")
    except Exception as e:
        print(f"❌ ไม่สามารถเชื่อมต่อ ClickHouse ได้: {e}")
        return
    
    print("\n🏗️ กำลังสร้างตารางทั้งหมด...")
    print("=" * 60)
    
    # สร้างตารางทั้งหมด
    tables_created = 0
    total_tables = 10
    
    if create_bookings_table(client):
        tables_created += 1
    if create_services_table(client):
        tables_created += 1
    if create_payments_table(client):
        tables_created += 1
    if create_staff_assignments_table(client):
        tables_created += 1
    if create_staff_wages_table(client):
        tables_created += 1
    if create_daily_summary_table(client):
        tables_created += 1
    if create_service_sales_table(client):
        tables_created += 1
    if create_product_sales_table(client):
        tables_created += 1
    if create_cashier_transactions_table(client):
        tables_created += 1
    if create_inventory_table(client):
        tables_created += 1
    if create_rejected_customers_table(client):
        tables_created += 1
        total_tables += 1
    
    print("=" * 60)
    print(f"🎉 สร้างตารางเสร็จสิ้น! ({tables_created}/{total_tables} ตาราง)")
    
    if tables_created == total_tables:
        print("✅ สร้างตารางทั้งหมดสำเร็จ!")
        
        # แสดงรายการตารางที่สร้าง
        print("\n📋 รายการตารางที่สร้าง:")
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
            
    else:
        print(f"⚠️  สร้างตารางได้ {tables_created} จาก {total_tables} ตาราง")

if __name__ == '__main__':
    main()
