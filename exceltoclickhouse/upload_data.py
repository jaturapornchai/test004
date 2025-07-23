import pandas as pd
import clickhouse_connect
from datetime import datetime
import numpy as np

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

def get_thai_to_english_mapping():
    """
    ส่งคืน mapping ชื่อคอลัมน์ไทย-อังกฤษ
    """
    return {
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

def upload_bookings_data(client, df_overview):
    """
    อัพโหลดข้อมูลการจองหลัก
    """
    print("📤 กำลังอัพโหลดข้อมูลการจอง...")
    
    # เลือกเฉพาะคอลัมน์ที่เกี่ยวข้องกับการจอง
    booking_columns = [
        'รหัสการจอง', 'วันที่', 'เวลา', 'วัน-เวลา', 'เลขใบเสร็จ', 
        'ช่องทางการจอง', 'รหัสเอเจนซี่', 'ประเภท', 'id ของลูกค้า', 
        'ลูกค้า', 'เบอร์โทร', 'หมายเหตุการจอง', 'ผู้ทำรายการ'
    ]
    
    # สร้าง DataFrame สำหรับการจอง
    df_bookings = df_overview[booking_columns].copy()
    
    # ลบข้อมูลซ้ำ (เพราะ 1 การจองอาจมีหลายบริการ)
    df_bookings = df_bookings.drop_duplicates(subset=['รหัสการจอง'])
    
    # เปลี่ยนชื่อคอลัมน์
    mapping = get_thai_to_english_mapping()
    df_bookings = df_bookings.rename(columns=mapping)
    
    # ทำความสะอาดข้อมูล
    for col in df_bookings.columns:
        if df_bookings[col].dtype == 'object':
            df_bookings[col] = df_bookings[col].fillna('')
        else:
            df_bookings[col] = df_bookings[col].fillna(0)
    
    try:
        # อัพโหลดข้อมูล
        client.insert_df('bookings', df_bookings)
        print(f"✅ อัพโหลดข้อมูลการจองสำเร็จ: {len(df_bookings)} รายการ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถอัพโหลดข้อมูลการจองได้: {e}")
        return False

def upload_services_data(client, df_overview):
    """
    อัพโหลดข้อมูลบริการ
    """
    print("📤 กำลังอัพโหลดข้อมูลบริการ...")
    
    # เลือกเฉพาะคอลัมน์ที่เกี่ยวข้องกับบริการ
    service_columns = [
        'รหัสการจอง', 'รายการ', 'เวลาบริการ', 'ราคาต่อบริการ', 
        'ส่วนลดต่อบริการ', 'จำนวนรายการทั้งหมด', 'ต้นทุนบริการ'
    ]
    
    df_services = df_overview[service_columns].copy()
    
    # เปลี่ยนชื่อคอลัมน์
    mapping = get_thai_to_english_mapping()
    df_services = df_services.rename(columns=mapping)
    
    # กรองเฉพาะแถวที่มีข้อมูลบริการ
    df_services = df_services[df_services['service_item'].notna()]
    
    # ทำความสะอาดข้อมูล
    df_services = df_services.fillna(0)
    
    try:
        client.insert_df('services', df_services)
        print(f"✅ อัพโหลดข้อมูลบริการสำเร็จ: {len(df_services)} รายการ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถอัพโหลดข้อมูลบริการได้: {e}")
        return False

def upload_payments_data(client, df_overview):
    """
    อัพโหลดข้อมูลการชำระเงิน
    """
    print("📤 กำลังอัพโหลดข้อมูลการชำระเงิน...")
    
    # เลือกเฉพาะคอลัมน์ที่เกี่ยวข้องกับการชำระเงิน
    payment_columns = [
        'รหัสการจอง', 'ราคารวม', 'ส่วนลดรวม', 'ราคาหลังหักส่วนลด', 'เซอร์วิสชาร์จ',
        'ชาร์จ', 'ราคาก่อน VAT', 'VAT (7%)', 'ทิปรวม', 'ค่าคอมมิชชันเอเจนซี',
        'ยอดชำระสุทธิ', 'เงินสด', 'เงินโอน', 'บัตรเครดิต', 'เอเจนซี', 
        'แพ็กเกจ', 'E-Wallet', 'ประเภท E-Wallet', 'หมายเหตุการชำระเงิน'
    ]
    
    df_payments = df_overview[payment_columns].copy()
    
    # ลบข้อมูลซ้ำ
    df_payments = df_payments.drop_duplicates(subset=['รหัสการจอง'])
    
    # เปลี่ยนชื่อคอลัมน์
    mapping = get_thai_to_english_mapping()
    df_payments = df_payments.rename(columns=mapping)
    
    # กรองเฉพาะแถวที่มีข้อมูลการชำระเงิน
    df_payments = df_payments[df_payments['net_payment'].notna()]
    
    # ทำความสะอาดข้อมูล
    df_payments = df_payments.fillna(0)
    
    try:
        client.insert_df('payments', df_payments)
        print(f"✅ อัพโหลดข้อมูลการชำระเงินสำเร็จ: {len(df_payments)} รายการ")
        return True
    except Exception as e:
        print(f"❌ ไม่สามารถอัพโหลดข้อมูลการชำระเงินได้: {e}")
        return False

def upload_staff_assignments_data(client, df_overview):
    """
    อัพโหลดข้อมูลการจัดพนักงาน
    """
    print("📤 กำลังอัพโหลดข้อมูลการจัดพนักงาน...")
    
    staff_assignments = []
    
    for _, row in df_overview.iterrows():
        booking_id = row['รหัสการจอง']
        
        # พนักงานคนที่ 1
        if pd.notna(row.get('พนักงาน (1)')):
            staff_assignments.append({
                'booking_id': booking_id,
                'staff_sequence': 1,
                'staff_name': row.get('พนักงาน (1)', ''),
                'staff_fee': row.get('ค่ามือหมอนวด (1)', 0),
                'staff_tip': row.get('ทิปหมอนวด (1)', 0),
                'request_fee': row.get('ค่ารีเควส (1)', 0)
            })
        
        # พนักงานคนที่ 2
        if pd.notna(row.get('พนักงาน (2)')):
            staff_assignments.append({
                'booking_id': booking_id,
                'staff_sequence': 2,
                'staff_name': row.get('พนักงาน (2)', ''),
                'staff_fee': row.get('ค่ามือหมอนวด (2)', 0),
                'staff_tip': row.get('ทิปหมอนวด (2)', 0),
                'request_fee': row.get('ค่ารีเควส (2)', 0)
            })
    
    if staff_assignments:
        df_staff = pd.DataFrame(staff_assignments)
        
        try:
            client.insert_df('staff_assignments', df_staff)
            print(f"✅ อัพโหลดข้อมูลการจัดพนักงานสำเร็จ: {len(df_staff)} รายการ")
            return True
        except Exception as e:
            print(f"❌ ไม่สามารถอัพโหลดข้อมูลการจัดพนักงานได้: {e}")
            return False
    else:
        print("⚠️  ไม่พบข้อมูลการจัดพนักงาน")
        return True

def upload_staff_wages_data(client, df_wages):
    """
    อัพโหลดข้อมูลค่ามือพนักงาน (จาก sheet ค่ามือ)
    """
    print("📤 กำลังอัพโหลดข้อมูลค่ามือพนักงาน...")
    
    staff_wages = []
    
    for _, row in df_wages.iterrows():
        date_record = row.get('วันที่', '')
        record_type = row.get('Unnamed: 1', '')
        
        # ข้ามแถวที่ไม่มีข้อมูลสำคัญ
        if pd.isna(date_record) and pd.isna(record_type):
            continue
            
        # วนลูปผ่านคอลัมน์ชื่อพนักงาน (เริ่มจากคอลัมน์ที่ 2)
        for col in df_wages.columns[2:]:
            amount = row.get(col)
            if pd.notna(amount) and amount != 0:
                staff_wages.append({
                    'date_record': date_record if pd.notna(date_record) else '',
                    'record_type': record_type if pd.notna(record_type) else '',
                    'staff_name': col,
                    'amount': float(amount)
                })
    
    if staff_wages:
        df_staff_wages = pd.DataFrame(staff_wages)
        
        try:
            client.insert_df('staff_wages', df_staff_wages)
            print(f"✅ อัพโหลดข้อมูลค่ามือพนักงานสำเร็จ: {len(df_staff_wages)} รายการ")
            return True
        except Exception as e:
            print(f"❌ ไม่สามารถอัพโหลดข้อมูลค่ามือพนักงานได้: {e}")
            return False
    else:
        print("⚠️  ไม่พบข้อมูลค่ามือพนักงาน")
        return True

def upload_other_sheets_data(client, excel_file):
    """
    อัพโหลดข้อมูลจาก sheet อื่นๆ
    """
    try:
        # ยอดเงิน
        df_money = pd.read_excel(excel_file, sheet_name='ยอดเงิน')
        if not df_money.empty:
            # ปรับโครงสร้างข้อมูล - ข้ามแถวแรกที่เป็นหัวคอลัมน์
            money_records = []
            for i, row in df_money.iterrows():
                if i == 0:  # ข้ามแถวแรก
                    continue
                if pd.notna(row.iloc[0]) and str(row.iloc[0]) not in ['เงินสด', 'บัตรเครดิต', 'เงินโอน', 'e-Wallet']:
                    continue
                if pd.notna(row.iloc[0]):  # ถ้าคอลัมน์แรกมีค่า
                    money_records.append({
                        'summary_date': '2025-01-01',  # ใส่วันที่ default
                        'payment_type': str(row.iloc[0]),
                        'total_amount': float(row.iloc[1]) if pd.notna(row.iloc[1]) and str(row.iloc[1]).replace('.','').replace(',','').isdigit() else 0
                    })
            
            if money_records:
                df_daily = pd.DataFrame(money_records)
                client.insert_df('daily_summary', df_daily)
                print(f"✅ อัพโหลดยอดเงินสำเร็จ: {len(df_daily)} รายการ")
        
        # ยอดขายตามรายการ
        df_service_sales = pd.read_excel(excel_file, sheet_name='ยอดขายตามรายการ')
        if not df_service_sales.empty:
            # เปลี่ยนชื่อคอลัมน์
            df_service_sales = df_service_sales.rename(columns={
                'กลุ่มบริการ': 'service_group',
                'บริการ': 'service_name', 
                'จำนวน': 'quantity',
                'ราคาขาย': 'sales_amount'
            })
            client.insert_df('service_sales', df_service_sales)
            print(f"✅ อัพโหลดยอดขายตามรายการสำเร็จ: {len(df_service_sales)} รายการ")
        
        # ยอดขายสินค้า
        df_product_sales = pd.read_excel(excel_file, sheet_name='ยอดขายสินค้า')
        if not df_product_sales.empty:
            df_product_sales = df_product_sales.rename(columns={
                'กลุ่มสินค้า': 'product_group',
                'สินค้า': 'product_name',
                'จำนวน': 'quantity', 
                'ราคาขาย': 'sales_amount'
            })
            client.insert_df('product_sales', df_product_sales)
            print(f"✅ อัพโหลดยอดขายสินค้าสำเร็จ: {len(df_product_sales)} รายการ")
        
        # แคชเชียร์
        df_cashier = pd.read_excel(excel_file, sheet_name='แคชเชียร์')
        if not df_cashier.empty:
            df_cashier = df_cashier.rename(columns={
                'วันที่': 'transaction_date',
                'รายละเอียด': 'description',
                'ประเภท': 'transaction_type',
                'ราคา': 'amount'
            })
            client.insert_df('cashier_transactions', df_cashier)
            print(f"✅ อัพโหลดรายการแคชเชียร์สำเร็จ: {len(df_cashier)} รายการ")
        
        # สต๊อกสินค้า
        df_inventory = pd.read_excel(excel_file, sheet_name='สต๊อกสินค้า')
        if not df_inventory.empty:
            df_inventory = df_inventory.rename(columns={
                'วันที่': 'inventory_date',
                'กลุ่มสินค้า': 'product_group',
                'สินค้า': 'product_name',
                'ขาย': 'sold_qty',
                'เติม': 'added_qty',
                'ลด': 'reduced_qty',
                'ดึงกลับ': 'returned_qty',
                'คงเหลือ': 'remaining_qty'
            })
            client.insert_df('inventory', df_inventory)
            print(f"✅ อัพโหลดสต๊อกสินค้าสำเร็จ: {len(df_inventory)} รายการ")
        
        # ปฏิเสธลูกค้า
        df_rejected = pd.read_excel(excel_file, sheet_name='ปฎิเสธลูกค้า')
        if not df_rejected.empty:
            df_rejected = df_rejected.rename(columns={
                'วัน-เวลา': 'datetime_rejected',
                'วันที่จอง': 'booking_date',
                'เวลาที่จอง': 'booking_time',
                'ช่องทางการจอง': 'booking_channel',
                'จำนวนลูกค้า': 'customer_count',
                'เหตุผลการไม่รับลูกค้า': 'rejection_reason',
                'ลูกค้า': 'customer_name',
                'เบอร์โทร': 'phone_number',
                'หมายเหตุการจอง': 'booking_note'
            })
            client.insert_df('rejected_customers', df_rejected)
            print(f"✅ อัพโหลดลูกค้าที่ปฏิเสธสำเร็จ: {len(df_rejected)} รายการ")
        
        return True
        
    except Exception as e:
        print(f"❌ เกิดข้อผิดพลาดในการอัพโหลดข้อมูล sheet อื่นๆ: {e}")
        return False

def main():
    """
    ฟังก์ชันหลักสำหรับอัพโหลดข้อมูล
    """
    excel_file = "sample.xlsx"
    
    print("🔗 กำลังเชื่อมต่อ ClickHouse...")
    try:
        client = get_clickhouse_client()
        print("✅ เชื่อมต่อ ClickHouse สำเร็จ!")
    except Exception as e:
        print(f"❌ ไม่สามารถเชื่อมต่อ ClickHouse ได้: {e}")
        return
    
    print("📖 กำลังอ่านไฟล์ Excel...")
    try:
        # อ่านข้อมูลจาก sheet หลัก
        df_overview = pd.read_excel(excel_file, sheet_name='ภาพรวม')
        df_wages = pd.read_excel(excel_file, sheet_name='ค่ามือ')
        print(f"✅ อ่านไฟล์ Excel สำเร็จ!")
    except Exception as e:
        print(f"❌ ไม่สามารถอ่านไฟล์ Excel ได้: {e}")
        return
    
    print("\n📤 กำลังอัพโหลดข้อมูลทั้งหมด...")
    print("=" * 60)
    
    # อัพโหลดข้อมูลทีละตาราง
    success_count = 0
    total_uploads = 6
    
    if upload_bookings_data(client, df_overview):
        success_count += 1
    if upload_services_data(client, df_overview):
        success_count += 1
    if upload_payments_data(client, df_overview):
        success_count += 1
    if upload_staff_assignments_data(client, df_overview):
        success_count += 1
    if upload_staff_wages_data(client, df_wages):
        success_count += 1
    if upload_other_sheets_data(client, excel_file):
        success_count += 1
    
    print("=" * 60)
    print(f"🎉 อัพโหลดข้อมูลเสร็จสิ้น! ({success_count}/{total_uploads} งาน)")
    
    if success_count == total_uploads:
        print("✅ อัพโหลดข้อมูลทั้งหมดสำเร็จ!")
        
        # แสดงสรุปข้อมูลในตาราง
        print(f"\n📊 สรุปข้อมูลในตาราง:")
        tables_to_check = [
            'bookings', 'services', 'payments', 'staff_assignments', 
            'staff_wages', 'daily_summary', 'service_sales', 'product_sales',
            'cashier_transactions', 'inventory', 'rejected_customers'
        ]
        
        for table in tables_to_check:
            try:
                result = client.query(f"SELECT COUNT(*) FROM {table}")
                count = result.result_set[0][0]
                print(f"  📋 {table}: {count:,} รายการ")
            except:
                print(f"  ❌ {table}: ไม่สามารถตรวจสอบได้")
    else:
        print(f"⚠️  อัพโหลดข้อมูลได้ {success_count} จาก {total_uploads} งาน")

if __name__ == '__main__':
    main()
