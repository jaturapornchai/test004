// ========================================
// Utility Functions
// ========================================

/**
 * รวม class names
 */
export function cn(...classes: (string | undefined | false | null)[]): string {
  return classes.filter(Boolean).join(' ');
}

/**
 * จัดรูปแบบตัวเลขเป็นเงินบาท
 */
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('th-TH', {
    style: 'currency',
    currency: 'THB',
    minimumFractionDigits: 2,
  }).format(amount);
}

/**
 * จัดรูปแบบตัวเลข
 */
export function formatNumber(num: number, decimals: number = 0): string {
  return new Intl.NumberFormat('th-TH', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(num);
}

/**
 * จัดรูปแบบวันที่แบบไทย
 */
export function formatDate(date: string | Date): string {
  const d = new Date(date);
  return d.toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

/**
 * จัดรูปแบบวันที่แบบสั้น
 */
export function formatDateShort(date: string | Date): string {
  const d = new Date(date);
  return d.toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

/**
 * จัดรูปแบบวันที่และเวลา
 */
export function formatDateTime(date: string | Date): string {
  const d = new Date(date);
  return d.toLocaleString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * แปลงวันที่เป็น พ.ศ.
 */
export function toBuddhistYear(date: string | Date): number {
  const d = new Date(date);
  return d.getFullYear() + 543;
}

/**
 * จัดรูปแบบเลขบัตรประชาชน
 */
export function formatIdCard(idCard: string): string {
  if (!idCard || idCard.length !== 13) return idCard;
  return `${idCard.slice(0, 1)}-${idCard.slice(1, 5)}-${idCard.slice(5, 10)}-${idCard.slice(10, 12)}-${idCard.slice(12)}`;
}

/**
 * จัดรูปแบบเบอร์โทรศัพท์
 */
export function formatPhone(phone: string): string {
  if (!phone) return phone;
  if (phone.length === 10) {
    return `${phone.slice(0, 3)}-${phone.slice(3, 6)}-${phone.slice(6)}`;
  }
  if (phone.length === 9) {
    return `${phone.slice(0, 2)}-${phone.slice(2, 5)}-${phone.slice(5)}`;
  }
  return phone;
}

/**
 * สถานะสมาชิก
 */
export const memberStatusLabels: Record<string, string> = {
  active: 'ปกติ',
  inactive: 'ไม่ใช้งาน',
  suspended: 'ระงับ',
  resigned: 'ลาออก',
  deceased: 'ถึงแก่กรรม',
  expelled: 'ให้ออก',
};

export const memberStatusColors: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  inactive: 'bg-slate-100 text-slate-800',
  suspended: 'bg-amber-100 text-amber-800',
  resigned: 'bg-gray-100 text-gray-800',
  deceased: 'bg-gray-100 text-gray-800',
  expelled: 'bg-red-100 text-red-800',
};

const memberStatusVariants: Record<string, 'default' | 'success' | 'warning' | 'danger' | 'info'> = {
  active: 'success',
  inactive: 'default',
  suspended: 'warning',
  resigned: 'warning',
  deceased: 'warning',
  expelled: 'danger',
};

/**
 * แปลงสถานะสมาชิกเป็นข้อความที่อ่านง่าย
 */
export function getMemberStatusLabel(status: string): string {
  return memberStatusLabels[status] ?? 'ไม่ระบุ';
}

/**
 * แปลงสถานะสมาชิกเป็นรูปแบบ Badge variant
 */
export function getMemberStatusColor(status: string): 'default' | 'success' | 'warning' | 'danger' | 'info' {
  return memberStatusVariants[status] ?? 'default';
}

/**
 * สถานะสินเชื่อ
 */
export const loanStatusLabels: Record<string, string> = {
  draft: 'ร่าง',
  pending: 'รอพิจารณา',
  approved: 'อนุมัติ',
  rejected: 'ไม่อนุมัติ',
  cancelled: 'ยกเลิก',
  active: 'ปกติ',
  closed: 'ปิดบัญชี',
  default: 'ผิดนัด',
  restructured: 'ปรับโครงสร้าง',
};

export const loanStatusColors: Record<string, string> = {
  draft: 'bg-gray-100 text-gray-800',
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
  cancelled: 'bg-gray-100 text-gray-800',
  active: 'bg-blue-100 text-blue-800',
  closed: 'bg-gray-100 text-gray-800',
  default: 'bg-red-100 text-red-800',
  restructured: 'bg-orange-100 text-orange-800',
};

/**
 * ประเภทสินเชื่อ
 */
export const loanTypeLabels: Record<string, string> = {
  emergency: 'เงินกู้ฉุกเฉิน',
  ordinary: 'เงินกู้สามัญ',
  special: 'เงินกู้พิเศษ',
};

/**
 * ประเภทเงินฝาก
 */
export const depositTypeLabels: Record<string, string> = {
  savings: 'เงินฝากออมทรัพย์',
  special_savings: 'เงินฝากออมทรัพย์พิเศษ',
  fixed: 'เงินฝากประจำ',
};

/**
 * คำนำหน้าชื่อ
 */
export const prefixes = [
  { value: 'นาย', label: 'นาย' },
  { value: 'นาง', label: 'นาง' },
  { value: 'นางสาว', label: 'นางสาว' },
  { value: 'ดร.', label: 'ดร.' },
];

/**
 * จังหวัดในประเทศไทย
 */
export const provinces = [
  'กรุงเทพมหานคร', 'กระบี่', 'กาญจนบุรี', 'กาฬสินธุ์', 'กำแพงเพชร',
  'ขอนแก่น', 'จันทบุรี', 'ฉะเชิงเทรา', 'ชลบุรี', 'ชัยนาท',
  'ชัยภูมิ', 'ชุมพร', 'เชียงราย', 'เชียงใหม่', 'ตรัง',
  'ตราด', 'ตาก', 'นครนายก', 'นครปฐม', 'นครพนม',
  'นครราชสีมา', 'นครศรีธรรมราช', 'นครสวรรค์', 'นนทบุรี', 'นราธิวาส',
  'น่าน', 'บึงกาฬ', 'บุรีรัมย์', 'ปทุมธานี', 'ประจวบคีรีขันธ์',
  'ปราจีนบุรี', 'ปัตตานี', 'พระนครศรีอยุธยา', 'พะเยา', 'พังงา',
  'พัทลุง', 'พิจิตร', 'พิษณุโลก', 'เพชรบุรี', 'เพชรบูรณ์',
  'แพร่', 'ภูเก็ต', 'มหาสารคาม', 'มุกดาหาร', 'แม่ฮ่องสอน',
  'ยโสธร', 'ยะลา', 'ร้อยเอ็ด', 'ระนอง', 'ระยอง',
  'ราชบุรี', 'ลพบุรี', 'ลำปาง', 'ลำพูน', 'เลย',
  'ศรีสะเกษ', 'สกลนคร', 'สงขลา', 'สตูล', 'สมุทรปราการ',
  'สมุทรสงคราม', 'สมุทรสาคร', 'สระแก้ว', 'สระบุรี', 'สิงห์บุรี',
  'สุโขทัย', 'สุพรรณบุรี', 'สุราษฎร์ธานี', 'สุรินทร์', 'หนองคาย',
  'หนองบัวลำภู', 'อ่างทอง', 'อำนาจเจริญ', 'อุดรธานี', 'อุตรดิตถ์',
  'อุทัยธานี', 'อุบลราชธานี',
];

/**
 * รายชื่อจังหวัดในประเทศไทย (alias)
 */
export const thaiProvinces = provinces;

/**
 * ประเภทพืช
 */
export const cropTypes = [
  'ข้าว', 'ข้าวโพด', 'มันสำปะหลัง', 'อ้อย', 'ยางพารา',
  'ปาล์มน้ำมัน', 'ลำไย', 'ทุเรียน', 'มังคุด', 'เงาะ',
  'ส้ม', 'มะม่วง', 'มะพร้าว', 'กาแฟ', 'ลิ้นจี่',
  'สับปะรด', 'พืชผัก', 'ไม้ดอก', 'อื่นๆ',
];

/**
 * ตรวจสอบ validation
 */
export function validateIdCard(idCard: string): boolean {
  if (!idCard || idCard.length !== 13 || !/^\d+$/.test(idCard)) {
    return false;
  }

  let sum = 0;
  for (let i = 0; i < 12; i++) {
    sum += parseInt(idCard[i]) * (13 - i);
  }

  const checkDigit = (11 - (sum % 11)) % 10;
  return checkDigit === parseInt(idCard[12]);
}

export function validatePhone(phone: string): boolean {
  return /^0\d{8,9}$/.test(phone);
}

export function validateEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}
