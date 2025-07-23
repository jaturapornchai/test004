"""
สร้างคู่มือโครงสร้างข้อมูลเป็น PDF พร้อมรองรับภาษาไทย
"""

from reportlab.lib.pagesizes import A4
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle, PageBreak
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.enums import TA_LEFT, TA_CENTER
from reportlab.lib import colors
from reportlab.lib.units import inch
from reportlab.pdfbase import pdfutils
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.pdfbase import pdfmetrics
import os
from datetime import datetime
import clickhouse_connect

# ลงทะเบียนฟอนต์ภาษาไทย
try:
    # ใช้ฟอนต์ระบบ Windows ที่รองรับไทย
    pdfmetrics.registerFont(TTFont('THSarabun', 'C:/Windows/Fonts/thsarabunnew.ttf'))
    pdfmetrics.registerFont(TTFont('THSarabun-Bold', 'C:/Windows/Fonts/thsarabunnew-bold.ttf'))
    THAI_FONT = 'THSarabun'
    THAI_FONT_BOLD = 'THSarabun-Bold'
    print("Using Thai fonts: THSarabun")
except Exception as e:
    # ถ้าไม่พบฟอนต์ไทย ใช้ฟอนต์เดิม
    THAI_FONT = 'Helvetica'
    THAI_FONT_BOLD = 'Helvetica-Bold'
    print(f"Thai fonts not found ({e}), using Helvetica")

def get_clickhouse_client():
    """เชื่อมต่อกับ ClickHouse"""
    return clickhouse_connect.get_client(
        host='xkruwixlx4.europe-west4.gcp.clickhouse.cloud',
        port=8443,
        username='default',
        password='9a.dPQwLg~Kj0',
        database='default',
        secure=True
    )

def get_table_schema_info():
    """ดึงข้อมูลโครงสร้าง table ทั้งหมด"""
    client = get_clickhouse_client()
    
    # รายชื่อ table ทั้งหมด
    tables = [
        'ภาพรวม', 'bookings', 'services', 'payments', 
        'staff_assignments', 'staff_wages', 'sales_by_items', 
        'product_sales', 'cashier_data', 'stock_data', 'customer_rejections'
    ]
    
    table_info = {}
    
    for table in tables:
        try:
            # ดึงโครงสร้าง table
            result = client.query(f"DESCRIBE TABLE {table}")
            columns = []
            for row in result.result_rows:
                columns.append({
                    'name': row[0],
                    'type': row[1],
                    'default': row[2] if len(row) > 2 else '',
                    'comment': row[3] if len(row) > 3 else ''
                })
            
            # ดึงจำนวน record
            count_result = client.query(f"SELECT COUNT(*) FROM {table}")
            record_count = count_result.result_rows[0][0]
            
            table_info[table] = {
                'columns': columns,
                'record_count': record_count
            }
            
        except Exception as e:
            print(f"Error getting info for table {table}: {e}")
            table_info[table] = {
                'columns': [],
                'record_count': 0
            }
    
    return table_info

def create_title_page(story):
    """สร้างหน้าปก"""
    # สร้าง styles สำหรับภาษาไทย
    styles = getSampleStyleSheet()
    title_style = ParagraphStyle(
        'CustomTitle',
        parent=styles['Title'],
        fontName=THAI_FONT_BOLD,
        fontSize=24,
        alignment=TA_CENTER,
        spaceAfter=30
    )
    
    subtitle_style = ParagraphStyle(
        'CustomSubtitle',
        parent=styles['Normal'],
        fontName=THAI_FONT,
        fontSize=16,
        alignment=TA_CENTER,
        spaceAfter=20
    )
    
    # เนื้อหาหน้าปก
    story.append(Spacer(1, 2*inch))
    story.append(Paragraph("คู่มือโครงสร้างข้อมูลระบบ Spa & Massage", title_style))
    story.append(Paragraph("Database Structure Manual", subtitle_style))
    story.append(Spacer(1, inch))
    
    # ข้อมูลเพิ่มเติม
    info_data = [
        ["ประเภทฐานข้อมูล:", "ClickHouse Cloud"],
        ["จำนวนตาราง:", "11 tables"],
        ["แหล่งข้อมูล:", "Excel (8 sheets)"],
        ["วันที่สร้าง:", datetime.now().strftime("%d/%m/%Y")],
        ["เวลา:", datetime.now().strftime("%H:%M:%S")]
    ]
    
    intro_table = Table(info_data, colWidths=[2.5*inch, 2.5*inch])
    intro_table.setStyle(TableStyle([
        ('FONTNAME', (0, 0), (-1, -1), THAI_FONT),
        ('FONTSIZE', (0, 0), (-1, -1), 12),
        ('ALIGN', (0, 0), (0, -1), 'RIGHT'),
        ('ALIGN', (1, 0), (1, -1), 'LEFT'),
        ('VALIGN', (0, 0), (-1, -1), 'MIDDLE'),
        ('BOTTOMPADDING', (0, 0), (-1, -1), 8),
    ]))
    
    story.append(intro_table)
    story.append(PageBreak())

def create_table_of_contents(story, table_info):
    """สร้างสารบัญ"""
    styles = getSampleStyleSheet()
    heading_style = ParagraphStyle(
        'CustomHeading',
        parent=styles['Heading1'],
        fontName=THAI_FONT_BOLD,
        fontSize=18,
        alignment=TA_CENTER,
        spaceAfter=20
    )
    
    story.append(Paragraph("สารบัญ (Table of Contents)", heading_style))
    story.append(Spacer(1, 0.2*inch))
    
    toc_data = [["หมายเลข", "ชื่อตาราง", "จำนวน Records", "หน้า"]]
    
    page_num = 4  # เริ่มต้นหน้าที่ 4
    for i, (table_name, info) in enumerate(table_info.items(), 1):
        thai_name = get_thai_table_name(table_name)
        toc_data.append([
            str(i),
            f"{thai_name}\n({table_name})",
            f"{info['record_count']:,}",
            str(page_num)
        ])
        page_num += 1
    
    toc_table = Table(toc_data, colWidths=[0.8*inch, 3*inch, 1.2*inch, 0.8*inch])
    toc_table.setStyle(TableStyle([
        ('BACKGROUND', (0, 0), (-1, 0), colors.darkblue),
        ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
        ('ALIGN', (0, 0), (-1, -1), 'CENTER'),
        ('FONTNAME', (0, 0), (-1, 0), THAI_FONT_BOLD),
        ('FONTNAME', (0, 1), (-1, -1), THAI_FONT),
        ('FONTSIZE', (0, 0), (-1, -1), 10),
        ('GRID', (0, 0), (-1, -1), 1, colors.black),
        ('VALIGN', (0, 0), (-1, -1), 'TOP')
    ]))
    
    story.append(toc_table)
    story.append(PageBreak())

def get_thai_table_name(table_name):
    """แปลงชื่อตารางเป็นภาษาไทย"""
    name_mapping = {
        'ภาพรวม': 'ตารางภาพรวม',
        'bookings': 'ตารางการจอง',
        'services': 'ตารางบริการ',
        'payments': 'ตารางการชำระเงิน',
        'staff_assignments': 'ตารางมอบหมายพนักงาน',
        'staff_wages': 'ตารางค่าแรงพนักงาน',
        'sales_by_items': 'ตารางยอดขายรายการ',
        'product_sales': 'ตารางยอดขายสินค้า',
        'cashier_data': 'ตารางข้อมูลแคชเชียร์',
        'stock_data': 'ตารางข้อมูลสต๊อก',
        'customer_rejections': 'ตารางปฏิเสธลูกค้า'
    }
    return name_mapping.get(table_name, table_name)

def get_table_description(table_name):
    """คำอธิบายตาราง"""
    descriptions = {
        'ภาพรวม': 'ข้อมูลภาพรวมทั้งหมดของระบบ รวมสถิติและสรุปผล',
        'bookings': 'ข้อมูลการจองของลูกค้า เก็บรายละเอียดการนัดหมายและสถานะ',
        'services': 'รายละเอียดบริการที่ให้กับลูกค้า ประเภทบริการ ราคา และระยะเวลา',
        'payments': 'ข้อมูลการชำระเงิน ยอดเงิน วิธีการชำระ และสถานะการชำระ',
        'staff_assignments': 'การมอบหมายงานให้กับพนักงาน ระบุผู้รับผิดชอบในแต่ละบริการ',
        'staff_wages': 'ข้อมูลค่าแรงและค่าคอมมิชชั่นของพนักงาน',
        'sales_by_items': 'สถิติยอดขายแยกตามรายการบริการ สำหรับการวิเคราะห์',
        'product_sales': 'ข้อมูลยอดขายสินค้า เก็บปริมาณและมูลค่าการขาย',
        'cashier_data': 'ข้อมูลการทำงานของแคชเชียร์ สรุปยอดขายรายวัน',
        'stock_data': 'ข้อมูลสินค้าคงคลัง ปริมาณคงเหลือและการเคลื่อนไหว',
        'customer_rejections': 'บันทึกกรณีที่ปฏิเสธลูกค้า เหตุผลและรายละเอียด'
    }
    return descriptions.get(table_name, f'ตาราง {table_name}')

def create_overview_page(story, table_info):
    """สร้างหน้าภาพรวม"""
    styles = getSampleStyleSheet()
    heading_style = ParagraphStyle(
        'CustomHeading',
        parent=styles['Heading1'],
        fontName=THAI_FONT_BOLD,
        fontSize=16,
        spaceAfter=15
    )
    
    normal_style = ParagraphStyle(
        'NormalThai',
        parent=styles['Normal'],
        fontName=THAI_FONT,
        fontSize=11,
        spaceAfter=10
    )
    
    story.append(Paragraph("ภาพรวมโครงสร้างฐานข้อมูล", heading_style))
    
    overview_text = """
    ระบบฐานข้อมูลนี้ถูกออกแบบสำหรับธุรกิจสปาและนวด โดยแยกข้อมูลออกเป็น 11 ตาราง 
    เพื่อให้การจัดเก็บข้อมูลมีประสิทธิภาพและลดความซ้ำซ้อน (Database Normalization)
    
    ข้อมูลถูกแปลงจากไฟล์ Excel 8 แผ่นงาน ให้เป็นโครงสร้างฐานข้อมูลเชิงสัมพันธ์
    ที่สามารถสืบค้นและวิเคราะห์ได้อย่างมีประสิทธิภาพ
    """
    
    story.append(Paragraph(overview_text, normal_style))
    story.append(Spacer(1, 0.3*inch))
    
    # สรุปสถิติ
    total_records = sum(info['record_count'] for info in table_info.values())
    
    overview_data = [
        ["สถิติ", "จำนวน"],
        ["จำนวนตารางทั้งหมด", f"{len(table_info)} ตาราง"],
        ["จำนวนข้อมูลรวม", f"{total_records:,} records"],
        ["ฐานข้อมูลที่ใช้", "ClickHouse Cloud"],
        ["Engine ที่ใช้", "MergeTree"]
    ]
    
    overview_table = Table(overview_data, colWidths=[2.5*inch, 2*inch])
    overview_table.setStyle(TableStyle([
        ('BACKGROUND', (0, 0), (-1, 0), colors.navy),
        ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
        ('ALIGN', (0, 0), (-1, -1), 'CENTER'),
        ('FONTNAME', (0, 0), (-1, 0), THAI_FONT_BOLD),
        ('FONTNAME', (0, 1), (-1, -1), THAI_FONT),
        ('FONTSIZE', (0, 0), (-1, -1), 11),
        ('GRID', (0, 0), (-1, -1), 1, colors.black),
        ('VALIGN', (0, 0), (-1, -1), 'MIDDLE'),
    ]))
    
    story.append(overview_table)
    story.append(PageBreak())

def create_table_detail_page(story, table_name, table_info):
    """สร้างหน้ารายละเอียดของแต่ละตาราง"""
    styles = getSampleStyleSheet()
    
    heading_style = ParagraphStyle(
        'TableHeading',
        parent=styles['Heading1'],
        fontName=THAI_FONT_BOLD,
        fontSize=16,
        spaceAfter=15
    )
    
    desc_style = ParagraphStyle(
        'Description',
        parent=styles['Normal'],
        fontName=THAI_FONT,
        fontSize=11,
        spaceAfter=15
    )
    
    # หัวข้อ
    thai_name = get_thai_table_name(table_name)
    story.append(Paragraph(f"{thai_name} ({table_name})", heading_style))
    
    # คำอธิบาย
    description = get_table_description(table_name)
    story.append(Paragraph(f"<b>คำอธิบาย:</b> {description}", desc_style))
    story.append(Paragraph(f"<b>จำนวนข้อมูล:</b> {table_info['record_count']:,} records", desc_style))
    
    # ตารางโครงสร้าง
    if table_info['columns']:
        column_data = [["ลำดับ", "ชื่อ Column", "ชนิดข้อมูล", "หมายเหตุ"]]
        
        for i, col in enumerate(table_info['columns'], 1):
            column_data.append([
                str(i),
                col['name'],
                col['type'],
                col.get('comment', '-') or '-'
            ])
        
        if table_name == 'bookings':
            table_style_color = colors.darkgreen
        elif table_name in ['services', 'payments']:
            table_style_color = colors.darkblue  
        elif table_name in ['staff_assignments', 'staff_wages']:
            table_style_color = colors.purple
        else:
            table_style_color = colors.darkred
            
        detail_table = Table(column_data, colWidths=[0.5*inch, 2*inch, 1.5*inch, 1.8*inch])
        detail_table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), table_style_color),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
            ('ALIGN', (0, 0), (0, -1), 'CENTER'),  # ลำดับให้อยู่กลาง
            ('FONTNAME', (0, 0), (-1, 0), THAI_FONT_BOLD),
            ('FONTNAME', (0, 1), (-1, -1), THAI_FONT),
            ('FONTSIZE', (0, 0), (-1, -1), 9),
            ('GRID', (0, 0), (-1, -1), 1, colors.black),
            ('VALIGN', (0, 0), (-1, -1), 'TOP')
        ]))
        
        story.append(detail_table)
    
    story.append(PageBreak())

def create_pdf_manual():
    """สร้างไฟล์ PDF คู่มือ"""
    filename = "database_structure_manual_thai.pdf"
    doc = SimpleDocTemplate(filename, pagesize=A4, topMargin=0.8*inch)
    
    story = []
    
    print("กำลังดึงข้อมูลโครงสร้างตาราง...")
    table_info = get_table_schema_info()
    
    print("สร้างหน้าปก...")
    create_title_page(story)
    
    print("สร้างสารบัญ...")
    create_table_of_contents(story, table_info)
    
    print("สร้างหน้าภาพรวม...")
    create_overview_page(story, table_info)
    
    print("สร้างหน้ารายละเอียดแต่ละตาราง...")
    for table_name, info in table_info.items():
        print(f"  - สร้างหน้า {table_name}...")
        create_table_detail_page(story, table_name, info)
    
    print("กำลังสร้างไฟล์ PDF...")
    doc.build(story)
    
    print(f"สร้างไฟล์ {filename} เรียบร้อยแล้ว!")
    return filename

if __name__ == "__main__":
    create_pdf_manual()
