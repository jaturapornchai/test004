import pandas as pd
import clickhouse_connect
from datetime import datetime
import numpy as np

def analyze_excel_structure(file_path):
    """
    วิเคราะห์โครงสร้างข้อมูลใน Excel file
    และกำหนด data types ที่เหมาะสมสำหรับ ClickHouse
    """
    print(f"กำลังวิเคราะห์ไฟล์: {file_path}")
    
    # อ่านไฟล์ Excel
    try:
        df = pd.read_excel(file_path)
        print(f"✅ อ่านไฟล์สำเร็จ! จำนวนแถว: {len(df)}, จำนวนคอลัมน์: {len(df.columns)}")
    except Exception as e:
        print(f"❌ ไม่สามารถอ่านไฟล์ได้: {e}")
        return None, None, None
    
    # แสดงตัวอย่างข้อมูล 5 แถวแรก
    print("\n📊 ตัวอย่างข้อมูล 5 แถวแรก:")
    print(df.head())
    
    # แสดงข้อมูลเกี่ยวกับ data types
    print("\n🔍 ข้อมูลเกี่ยวกับ Data Types:")
    print(df.info())
    
    # วิเคราะห์และสร้าง ClickHouse schema
    clickhouse_schema = {}
    column_comments = {}
    
    print("\n🎯 การแปลง Data Types สำหรับ ClickHouse:")
    print("-" * 60)
    
    # สร้าง mapping สำหรับแปลงชื่อคอลัมน์เป็นภาษาอังกฤษ
    thai_to_english = {
        'รหัสการจอง': 'booking_id',
        'วันที่': 'booking_date', 
        'เวลา': 'booking_time',
        'วัน-เวลา': 'datetime_display',
        'เลขใบเสร็จ': 'receipt_number',
        'ช่องทางการจอง': 'booking_channel',
        'รหัสเอเจนซี่': 'agency_id',
        'ประเภท': 'booking_type',
        'id ของลูกค้า': 'customer_id',
        'ลูกค้า': 'customer_name',
        'เบอร์โทร': 'phone_number',
        'รายการ': 'service_item',
        'เวลาบริการ': 'service_duration',
        'ราคาต่อบริการ': 'service_price',
        'ส่วนลดต่อบริการ': 'service_discount',
        'จำนวนรายการทั้งหมด': 'total_items',
        'ราคารวม': 'total_price',
        'ส่วนลดรวม': 'total_discount',
        'ราคาหลังหักส่วนลด': 'price_after_discount',
        'เซอร์วิสชาร์จ': 'service_charge',
        'ชาร์จ': 'additional_charge',
        'ราคาก่อน VAT': 'price_before_vat',
        'VAT (7%)': 'vat_amount',
        'ทิปรวม': 'total_tip',
        'ค่าคอมมิชชันเอเจนซี': 'agency_commission',
        'ยอดชำระสุทธิ': 'net_payment',
        'เงินสด': 'cash_payment',
        'เงินโอน': 'transfer_payment',
        'บัตรเครดิต': 'credit_card_payment',
        'เอเจนซี': 'agency_payment',
        'แพ็กเกจ': 'package_payment',
        'E-Wallet': 'ewallet_payment',
        'ประเภท E-Wallet': 'ewallet_type',
        'พนักงาน (1)': 'staff_1',
        'ค่ามือหมอนวด (1)': 'staff_fee_1',
        'ทิปหมอนวด (1)': 'staff_tip_1',
        'ค่ารีเควส (1)': 'request_fee_1',
        'พนักงาน (2)': 'staff_2',
        'ค่ามือหมอนวด (2)': 'staff_fee_2',
        'ทิปหมอนวด (2)': 'staff_tip_2',
        'ค่ารีเควส (2)': 'request_fee_2',
        'ต้นทุนบริการ': 'service_cost',
        'หมายเหตุการจอง': 'booking_note',
        'หมายเหตุการชำระเงิน': 'payment_note',
        'ผู้ทำรายการ': 'created_by'
    }
    
    for col in df.columns:
        # ใช้ mapping หรือทำความสะอาดชื่อคอลัมน์
        if col in thai_to_english:
            clean_col_name = thai_to_english[col]
        else:
            # ถ้าไม่มีใน mapping ให้ทำความสะอาดตามปกติ
            clean_col_name = col.strip().replace(' ', '_').replace('-', '_').replace('.', '_').replace('(', '').replace(')', '').lower()
        
        # วิเคราะห์ข้อมูลในคอลัมน์
        col_data = df[col].dropna()  # ลบค่า null ออกก่อนวิเคราะห์
        
        if len(col_data) == 0:
            # ถ้าคอลัมน์ว่างทั้งหมด
            clickhouse_type = "Nullable(String)"
            comment = "คอลัมน์ว่าง - ใช้ String เป็น default"
        elif pd.api.types.is_integer_dtype(col_data):
            # ตัวเลขจำนวนเต็ม
            max_val = col_data.max()
            min_val = col_data.min()
            
            if min_val >= 0:
                # จำนวนเต็มบวก
                if max_val <= 255:
                    clickhouse_type = "UInt8" if col == 'รหัสการจอง' else "Nullable(UInt8)"
                elif max_val <= 65535:
                    clickhouse_type = "UInt16" if col == 'รหัสการจอง' else "Nullable(UInt16)"
                elif max_val <= 4294967295:
                    clickhouse_type = "UInt32" if col == 'รหัสการจอง' else "Nullable(UInt32)"
                else:
                    clickhouse_type = "UInt64" if col == 'รหัสการจอง' else "Nullable(UInt64)"
            else:
                # จำนวนเต็มที่อาจเป็นลบได้
                if min_val >= -128 and max_val <= 127:
                    clickhouse_type = "Int8" if col == 'รหัสการจอง' else "Nullable(Int8)"
                elif min_val >= -32768 and max_val <= 32767:
                    clickhouse_type = "Int16" if col == 'รหัสการจอง' else "Nullable(Int16)"
                elif min_val >= -2147483648 and max_val <= 2147483647:
                    clickhouse_type = "Int32" if col == 'รหัสการจอง' else "Nullable(Int32)"
                else:
                    clickhouse_type = "Int64" if col == 'รหัสการจอง' else "Nullable(Int64)"
            
            comment = f"ตัวเลขจำนวนเต็ม (ช่วง: {min_val} ถึง {max_val})"
            
        elif pd.api.types.is_float_dtype(col_data):
            # ตัวเลขทศนิยม
            clickhouse_type = "Nullable(Float64)"
            comment = f"ตัวเลขทศนิยม (ค่าเฉลี่ย: {col_data.mean():.2f})"
            
        elif pd.api.types.is_datetime64_any_dtype(col_data):
            # วันที่และเวลา
            clickhouse_type = "Nullable(DateTime)"
            comment = "วันที่และเวลา"
            
        elif pd.api.types.is_bool_dtype(col_data):
            # Boolean
            clickhouse_type = "Nullable(Bool)"
            comment = "ค่าจริง/เท็จ"
            
        else:
            # String (default)
            max_length = col_data.astype(str).str.len().max() if len(col_data) > 0 else 0
            clickhouse_type = "Nullable(String)"
            comment = f"ข้อความ (ความยาวสูงสุด: {max_length} ตัวอักษร)"
        
        # เก็บข้อมูล schema
        clickhouse_schema[clean_col_name] = clickhouse_type
        column_comments[clean_col_name] = f"{col} - {comment}"  # ใส่ชื่อภาษาไทยใน comment
        
        print(f"{col:30} -> {clean_col_name:30} -> {clickhouse_type:20} | {col} - {comment}")
    
    return df, clickhouse_schema, column_comments

def create_clickhouse_table(client, table_name, schema, comments, engine="MergeTree() ORDER BY booking_id"):
    """
    สร้าง table บน ClickHouse ตาม schema ที่กำหนด
    """
    print(f"\n🏗️  กำลังสร้างตาราง '{table_name}' บน ClickHouse...")
    
    # สร้าง SQL statement สำหรับการสร้างตาราง
    columns_def = []
    for col_name, col_type in schema.items():
        comment = comments.get(col_name, "")
        columns_def.append(f"    {col_name} {col_type} COMMENT '{comment}'")
    
    create_table_sql = f"""
CREATE TABLE IF NOT EXISTS {table_name} (
{',\n'.join(columns_def)}
) ENGINE = {engine}
COMMENT 'ตารางที่สร้างจากไฟล์ Excel - {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}'
"""
    
    print("📝 SQL Statement ที่จะใช้สร้างตาราง:")
    print("-" * 60)
    print(create_table_sql)
    print("-" * 60)
    
    try:
        # ลบตารางเก่าถ้ามี (เพื่อการทดสอบ)
        client.command(f"DROP TABLE IF EXISTS {table_name}")
        print(f"🗑️  ลบตารางเก่า '{table_name}' (ถ้ามี)")
        
        # สร้างตารางใหม่
        client.command(create_table_sql)
        print(f"✅ สร้างตาราง '{table_name}' สำเร็จ!")
        
        # แสดงข้อมูลตาราง
        result = client.query(f"DESCRIBE TABLE {table_name}")
        print(f"\n📋 โครงสร้างตาราง '{table_name}':")
        for row in result.result_set:
            print(f"  {row[0]:25} | {row[1]:20} | {row[6]}")  # name, type, comment
        
        return True
        
    except Exception as e:
        print(f"❌ เกิดข้อผิดพลาดในการสร้างตาราง: {e}")
        return False

def insert_data_to_clickhouse(client, table_name, df, schema):
    """
    นำเข้าข้อมูลจาก DataFrame ไปยัง ClickHouse
    """
    print(f"\n📤 กำลังนำเข้าข้อมูล {len(df)} แถว ไปยังตาราง '{table_name}'...")
    
    try:
        # ทำความสะอาดชื่อคอลัมน์ใน DataFrame ให้ตรงกับ schema
        df_clean = df.copy()
        
        # สร้าง mapping จากชื่อเก่าไปชื่อใหม่
        thai_to_english = {
            'รหัสการจอง': 'booking_id',
            'วันที่': 'booking_date', 
            'เวลา': 'booking_time',
            'วัน-เวลา': 'datetime_display',
            'เลขใบเสร็จ': 'receipt_number',
            'ช่องทางการจอง': 'booking_channel',
            'รหัสเอเจนซี่': 'agency_id',
            'ประเภท': 'booking_type',
            'id ของลูกค้า': 'customer_id',
            'ลูกค้า': 'customer_name',
            'เบอร์โทร': 'phone_number',
            'รายการ': 'service_item',
            'เวลาบริการ': 'service_duration',
            'ราคาต่อบริการ': 'service_price',
            'ส่วนลดต่อบริการ': 'service_discount',
            'จำนวนรายการทั้งหมด': 'total_items',
            'ราคารวม': 'total_price',
            'ส่วนลดรวม': 'total_discount',
            'ราคาหลังหักส่วนลด': 'price_after_discount',
            'เซอร์วิสชาร์จ': 'service_charge',
            'ชาร์จ': 'additional_charge',
            'ราคาก่อน VAT': 'price_before_vat',
            'VAT (7%)': 'vat_amount',
            'ทิปรวม': 'total_tip',
            'ค่าคอมมิชชันเอเจนซี': 'agency_commission',
            'ยอดชำระสุทธิ': 'net_payment',
            'เงินสด': 'cash_payment',
            'เงินโอน': 'transfer_payment',
            'บัตรเครดิต': 'credit_card_payment',
            'เอเจนซี': 'agency_payment',
            'แพ็กเกจ': 'package_payment',
            'E-Wallet': 'ewallet_payment',
            'ประเภท E-Wallet': 'ewallet_type',
            'พนักงาน (1)': 'staff_1',
            'ค่ามือหมอนวด (1)': 'staff_fee_1',
            'ทิปหมอนวด (1)': 'staff_tip_1',
            'ค่ารีเควส (1)': 'request_fee_1',
            'พนักงาน (2)': 'staff_2',
            'ค่ามือหมอนวด (2)': 'staff_fee_2',
            'ทิปหมอนวด (2)': 'staff_tip_2',
            'ค่ารีเควส (2)': 'request_fee_2',
            'ต้นทุนบริการ': 'service_cost',
            'หมายเหตุการจอง': 'booking_note',
            'หมายเหตุการชำระเงิน': 'payment_note',
            'ผู้ทำรายการ': 'created_by'
        }
        
        column_mapping = {}
        for original_col in df.columns:
            if original_col in thai_to_english:
                column_mapping[original_col] = thai_to_english[original_col]
            else:
                # ถ้าไม่มีใน mapping ให้ทำความสะอาดตามปกติ
                clean_name = original_col.strip().replace(' ', '_').replace('-', '_').replace('.', '_').replace('(', '').replace(')', '').lower()
                column_mapping[original_col] = clean_name
        
        df_clean = df_clean.rename(columns=column_mapping)
        
        # แปลงข้อมูลให้เหมาะสมกับ ClickHouse
        for col_name, col_type in schema.items():
            if col_name in df_clean.columns:
                if 'DateTime' in col_type:
                    df_clean[col_name] = pd.to_datetime(df_clean[col_name], errors='coerce')
                elif 'Int' in col_type or 'UInt' in col_type:
                    df_clean[col_name] = pd.to_numeric(df_clean[col_name], errors='coerce').astype('Int64')
                elif 'Float' in col_type:
                    df_clean[col_name] = pd.to_numeric(df_clean[col_name], errors='coerce')
                elif 'Bool' in col_type:
                    df_clean[col_name] = df_clean[col_name].astype('boolean')
        
        # นำเข้าข้อมูล
        client.insert_df(table_name, df_clean)
        print(f"✅ นำเข้าข้อมูลสำเร็จ! {len(df_clean)} แถว")
        
        # ตรวจสอบข้อมูลที่นำเข้า
        count_result = client.query(f"SELECT COUNT(*) FROM {table_name}")
        record_count = count_result.result_set[0][0]
        print(f"📊 จำนวนข้อมูลในตาราง: {record_count} แถว")
        
        # แสดงตัวอย่างข้อมูล 5 แถวแรก
        sample_result = client.query(f"SELECT * FROM {table_name} LIMIT 5")
        print(f"\n🔍 ตัวอย่างข้อมูลในตาราง '{table_name}' (5 แถวแรก):")
        for i, row in enumerate(sample_result.result_set, 1):
            print(f"  แถว {i}: {row}")
        
        return True
        
    except Exception as e:
        print(f"❌ เกิดข้อผิดพลาดในการนำเข้าข้อมูล: {e}")
        return False

if __name__ == '__main__':
    # เชื่อมต่อ ClickHouse
    print("🔗 กำลังเชื่อมต่อ ClickHouse...")
    client = clickhouse_connect.get_client(
        host='npomobbg93.germanywestcentral.azure.clickhouse.cloud',
        user='default',
        password='1S.6V_z9Lr9Wc',
        secure=True
    )
    print("✅ เชื่อมต่อ ClickHouse สำเร็จ!")
    
    # วิเคราะห์ไฟล์ Excel
    excel_file = "sample.xlsx"
    table_name = "sample_data"  # ชื่อตารางที่จะสร้าง
    
    df, schema, comments = analyze_excel_structure(excel_file)
    
    if df is not None and schema:
        # สร้างตาราง
        if create_clickhouse_table(client, table_name, schema, comments):
            # นำเข้าข้อมูล
            insert_data_to_clickhouse(client, table_name, df, schema)
        
        print("\n🎉 กระบวนการเสร็จสิ้น!")
        print(f"📝 สรุปผลลัพธ์:")
        print(f"   - ไฟล์ Excel: {excel_file}")
        print(f"   - ตาราง ClickHouse: {table_name}")
        print(f"   - จำนวนคอลัมน์: {len(schema)}")
        print(f"   - จำนวนแถวข้อมูล: {len(df)}")
    else:
        print("❌ ไม่สามารถวิเคราะห์ไฟล์ Excel ได้")
