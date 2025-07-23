"""
สร้างคู่มือโครงสร้างข้อมูลเป็น PDF พร้อมรองรับภาษาไทย (Static Version)
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

def get_static_table_info():
    """ข้อมูลโครงสร้าง table แบบ static"""
    return {
        'ภาพรวม': {
            'columns': [
                {'name': 'date', 'type': 'Date', 'comment': 'วันที่'},
                {'name': 'total_bookings', 'type': 'UInt32', 'comment': 'จำนวนการจองทั้งหมด'},
                {'name': 'total_revenue', 'type': 'Float64', 'comment': 'รายได้รวม'},
                {'name': 'staff_count', 'type': 'UInt16', 'comment': 'จำนวนพนักงาน'}
            ],
            'record_count': 365
        },
        'bookings': {
            'columns': [
                {'name': 'booking_id', 'type': 'UInt32', 'comment': 'รหัสการจอง'},
                {'name': 'booking_date', 'type': 'Date', 'comment': 'วันที่จอง'},
                {'name': 'booking_time', 'type': 'String', 'comment': 'เวลาจอง'},
                {'name': 'customer_name', 'type': 'String', 'comment': 'ชื่อลูกค้า'},
                {'name': 'phone_number', 'type': 'Nullable(String)', 'comment': 'เบอร์โทรศัพท์'},
                {'name': 'service_type', 'type': 'String', 'comment': 'ประเภทบริการ'},
                {'name': 'status', 'type': 'String', 'comment': 'สถานะการจอง'}
            ],
            'record_count': 3113
        },
        'services': {
            'columns': [
                {'name': 'service_id', 'type': 'UInt32', 'comment': 'รหัสบริการ'},
                {'name': 'booking_id', 'type': 'UInt32', 'comment': 'รหัสการจอง'},
                {'name': 'service_name', 'type': 'String', 'comment': 'ชื่อบริการ'},
                {'name': 'duration', 'type': 'UInt16', 'comment': 'ระยะเวลา (นาที)'},
                {'name': 'price', 'type': 'Float64', 'comment': 'ราคาบริการ'},
                {'name': 'staff_assigned', 'type': 'Nullable(String)', 'comment': 'พนักงานที่ได้รับมอบหมาย'}
            ],
            'record_count': 5053
        },
        'payments': {
            'columns': [
                {'name': 'payment_id', 'type': 'UInt32', 'comment': 'รหัสการชำระเงิน'},
                {'name': 'booking_id', 'type': 'UInt32', 'comment': 'รหัสการจอง'},
                {'name': 'amount', 'type': 'Float64', 'comment': 'จำนวนเงิน'},
                {'name': 'payment_method', 'type': 'String', 'comment': 'วิธีการชำระเงิน'},
                {'name': 'payment_status', 'type': 'String', 'comment': 'สถานะการชำระ'},
                {'name': 'payment_date', 'type': 'DateTime', 'comment': 'วันเวลาชำระเงิน'}
            ],
            'record_count': 3113
        },
        'staff_assignments': {
            'columns': [
                {'name': 'assignment_id', 'type': 'UInt32', 'comment': 'รหัสการมอบหมาย'},
                {'name': 'booking_id', 'type': 'UInt32', 'comment': 'รหัสการจอง'},
                {'name': 'staff_name', 'type': 'String', 'comment': 'ชื่อพนักงาน'},
                {'name': 'service_type', 'type': 'String', 'comment': 'ประเภทบริการ'},
                {'name': 'assignment_date', 'type': 'Date', 'comment': 'วันที่มอบหมาย'},
                {'name': 'completion_status', 'type': 'String', 'comment': 'สถานะการทำงาน'}
            ],
            'record_count': 4017
        },
        'staff_wages': {
            'columns': [
                {'name': 'wage_id', 'type': 'UInt32', 'comment': 'รหัสค่าแรง'},
                {'name': 'staff_name', 'type': 'String', 'comment': 'ชื่อพนักงาน'},
                {'name': 'work_date', 'type': 'Date', 'comment': 'วันที่ทำงาน'},
                {'name': 'base_wage', 'type': 'Float64', 'comment': 'ค่าแรงพื้นฐาน'},
                {'name': 'commission', 'type': 'Float64', 'comment': 'ค่าคอมมิชชั่น'},
                {'name': 'total_wage', 'type': 'Float64', 'comment': 'ค่าแรงรวม'}
            ],
            'record_count': 1892
        },
        'sales_by_items': {
            'columns': [
                {'name': 'item_id', 'type': 'UInt32', 'comment': 'รหัสรายการ'},
                {'name': 'item_name', 'type': 'String', 'comment': 'ชื่อรายการ'},
                {'name': 'sale_date', 'type': 'Date', 'comment': 'วันที่ขาย'},
                {'name': 'quantity_sold', 'type': 'UInt32', 'comment': 'จำนวนที่ขาย'},
                {'name': 'unit_price', 'type': 'Float64', 'comment': 'ราคาต่อหน่วย'},
                {'name': 'total_amount', 'type': 'Float64', 'comment': 'ยอดรวม'}
            ],
            'record_count': 2456
        },
        'product_sales': {
            'columns': [
                {'name': 'product_id', 'type': 'UInt32', 'comment': 'รหัสสินค้า'},
                {'name': 'product_name', 'type': 'String', 'comment': 'ชื่อสินค้า'},
                {'name': 'category', 'type': 'String', 'comment': 'หมวดหมู่สินค้า'},
                {'name': 'sale_date', 'type': 'Date', 'comment': 'วันที่ขาย'},
                {'name': 'quantity', 'type': 'UInt32', 'comment': 'จำนวน'},
                {'name': 'price', 'type': 'Float64', 'comment': 'ราคา'},
                {'name': 'total_value', 'type': 'Float64', 'comment': 'มูลค่ารวม'}
            ],
            'record_count': 1834
        },
        'cashier_data': {
            'columns': [
                {'name': 'cashier_id', 'type': 'UInt32', 'comment': 'รหัสแคชเชียร์'},
                {'name': 'cashier_name', 'type': 'String', 'comment': 'ชื่อแคชเชียร์'},
                {'name': 'shift_date', 'type': 'Date', 'comment': 'วันที่เข้าเวร'},
                {'name': 'shift_start', 'type': 'DateTime', 'comment': 'เวลาเริ่มเวร'},
                {'name': 'shift_end', 'type': 'DateTime', 'comment': 'เวลาสิ้นสุดเวร'},
                {'name': 'total_sales', 'type': 'Float64', 'comment': 'ยอดขายรวม'}
            ],
            'record_count': 892
        },
        'stock_data': {
            'columns': [
                {'name': 'stock_id', 'type': 'UInt32', 'comment': 'รหัสสต๊อก'},
                {'name': 'product_name', 'type': 'String', 'comment': 'ชื่อสินค้า'},
                {'name': 'category', 'type': 'String', 'comment': 'หมวดหมู่'},
                {'name': 'current_stock', 'type': 'UInt32', 'comment': 'สต๊อกปัจจุบัน'},
                {'name': 'min_stock', 'type': 'UInt32', 'comment': 'สต๊อกขั้นต่ำ'},
                {'name': 'last_updated', 'type': 'DateTime', 'comment': 'อัปเดตล่าสุด'}
            ],
            'record_count': 567
        },
        'customer_rejections': {
            'columns': [
                {'name': 'rejection_id', 'type': 'UInt32', 'comment': 'รหัสการปฏิเสธ'},
                {'name': 'customer_name', 'type': 'String', 'comment': 'ชื่อลูกค้า'},
                {'name': 'rejection_date', 'type': 'Date', 'comment': 'วันที่ปฏิเสธ'},
                {'name': 'reason', 'type': 'String', 'comment': 'เหตุผล'},
                {'name': 'staff_involved', 'type': 'String', 'comment': 'พนักงานที่เกี่ยวข้อง'},
                {'name': 'resolution', 'type': 'Nullable(String)', 'comment': 'วิธีแก้ไข'}
            ],
            'record_count': 123
        }
    }

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
    
    print("กำลังเตรียมข้อมูลโครงสร้างตาราง...")
    table_info = get_static_table_info()
    
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
