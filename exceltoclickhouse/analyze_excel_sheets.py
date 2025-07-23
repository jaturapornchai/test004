import pandas as pd
import clickhouse_connect
from datetime import datetime

def analyze_all_sheets(excel_file):
    """
    วิเคราะห์ข้อมูลทุก sheet ใน Excel file
    """
    print(f"🔍 กำลังวิเคราะห์ไฟล์: {excel_file}")
    
    try:
        xl = pd.ExcelFile(excel_file)
        sheets_info = {}
        
        for sheet_name in xl.sheet_names:
            print(f"\n📋 กำลังวิเคราะห์ Sheet: '{sheet_name}'")
            try:
                df = pd.read_excel(excel_file, sheet_name=sheet_name)
                print(f"   ✅ จำนวนแถว: {len(df)}, จำนวนคอลัมน์: {len(df.columns)}")
                print(f"   📝 คอลัมน์: {list(df.columns[:10])}{'...' if len(df.columns) > 10 else ''}")
                
                # เก็บข้อมูลสำหรับแต่ละ sheet
                sheets_info[sheet_name] = {
                    'dataframe': df,
                    'rows': len(df),
                    'columns': len(df.columns),
                    'column_names': list(df.columns)
                }
                
                # แสดงตัวอย่างข้อมูล 3 แถวแรก
                if len(df) > 0:
                    print(f"   🔍 ตัวอย่างข้อมูล 3 แถวแรก:")
                    print(df.head(3).to_string())
                
            except Exception as e:
                print(f"   ❌ ไม่สามารถอ่าน sheet '{sheet_name}' ได้: {e}")
                sheets_info[sheet_name] = None
        
        return sheets_info
        
    except Exception as e:
        print(f"❌ ไม่สามารถอ่านไฟล์ได้: {e}")
        return None

if __name__ == '__main__':
    # วิเคราะห์ไฟล์ Excel ทุก sheet
    excel_file = "sample.xlsx"
    sheets_info = analyze_all_sheets(excel_file)
    
    if sheets_info:
        print(f"\n🎯 สรุปข้อมูลทั้งหมด:")
        print(f"📊 จำนวน Sheet ทั้งหมด: {len(sheets_info)}")
        
        for sheet_name, info in sheets_info.items():
            if info:
                print(f"   - {sheet_name}: {info['rows']} แถว, {info['columns']} คอลัมน์")
            else:
                print(f"   - {sheet_name}: ไม่สามารถอ่านได้")
