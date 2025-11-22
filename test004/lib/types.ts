// ========================================
// Types สำหรับระบบบริหารจัดการสหกรณ์การเกษตร
// ========================================

// ========== สหกรณ์ (Cooperative) ==========
export interface Cooperative {
  _id?: string;
  coopId: string;           // รหัสสหกรณ์ (unique)
  name: string;             // ชื่อสหกรณ์
  nameEn?: string;          // ชื่อภาษาอังกฤษ
  registrationNo: string;   // เลขทะเบียนสหกรณ์
  type: 'agriculture' | 'savings' | 'consumer' | 'service' | 'credit_union';
  address: Address;
  phone: string;
  email?: string;
  taxId?: string;           // เลขประจำตัวผู้เสียภาษี
  establishedDate: string;  // วันที่จดทะเบียน

  // ข้อบังคับและพารามิเตอร์
  bylaws: Bylaws;

  // สถานะ
  status: 'active' | 'inactive' | 'suspended';
  createdAt: string;
  updatedAt: string;
}

export interface Address {
  addressLine1: string;
  addressLine2?: string;
  subDistrict: string;      // ตำบล
  district: string;         // อำเภอ
  province: string;         // จังหวัด
  postalCode: string;
}

export interface Bylaws {
  maxSharePerMember: number;       // เพดานการถือหุ้นสูงสุด (บาท)
  maxSharePercentage: number;      // ไม่เกิน x% ของทุนทั้งหมด (default: 20%)
  shareValue: number;              // มูลค่าหุ้นละ (บาท)
  minMonthlyShare: number;         // ค่าหุ้นรายเดือนขั้นต่ำ
  maxMonthlyShare: number;         // ค่าหุ้นรายเดือนสูงสุด
  dividendRate: number;            // อัตราเงินปันผล (%)
  patronageRefundRate: number;     // อัตราเฉลี่ยคืน (%)
  reserveFundRate: number;         // อัตรากันทุนสำรอง (%)

  // เงื่อนไขสินเชื่อ
  maxGuaranteeContracts: number;   // จำนวนสัญญาที่ค้ำประกันได้สูงสุด
  maxDSR: number;                  // Debt Service Ratio สูงสุด (%)

  // เงื่อนไขเงินฝาก
  savingsInterestRate: number;     // อัตราดอกเบี้ยเงินฝากออมทรัพย์ (%)
  specialSavingsInterestRate: number;
  fixedDepositRates: FixedDepositRate[];

  // ไฟล์ข้อบังคับ
  bylawDocuments: Document[];
}

export interface FixedDepositRate {
  months: number;           // ระยะเวลา (เดือน)
  interestRate: number;     // อัตราดอกเบี้ย (%)
}

export interface Document {
  name: string;
  url: string;
  uploadedAt: string;
  version: string;
}

// ========== สมาชิก (Member) ==========
export interface Member {
  _id?: string;
  coopId: string;           // รหัสสหกรณ์ (foreign key)
  memberId: string;         // รหัสสมาชิก (unique per coop)
  memberNo: string;         // เลขทะเบียนสมาชิก

  // ข้อมูลส่วนบุคคล
  prefix: string;           // คำนำหน้า
  firstName: string;
  lastName: string;
  idCardNo: string;         // เลขบัตรประชาชน
  birthDate: string;
  gender: 'male' | 'female' | 'other';
  phone: string;
  email?: string;
  address: Address;

  // ข้อมูลการเกษตร
  occupation: string;
  farmingInfo?: FarmingInfo;

  // ข้อมูลสมาชิกภาพ
  memberType: 'regular' | 'associate';  // สมาชิกสามัญ/สมทบ
  joinDate: string;
  resignDate?: string;
  status: 'active' | 'resigned' | 'deceased' | 'expelled';

  // ความยินยอม PDPA
  pdpaConsent: PDPAConsent;

  // ผู้รับผลประโยชน์
  beneficiaries: Beneficiary[];

  createdAt: string;
  updatedAt: string;
}

export interface FarmingInfo {
  landArea: number;         // พื้นที่ (ไร่)
  landOwnership: 'owned' | 'rented' | 'mixed';
  crops: CropInfo[];        // พืชที่ปลูก
  livestock?: string[];     // ปศุสัตว์
}

export interface CropInfo {
  cropType: string;         // ชนิดพืช
  area: number;             // พื้นที่ปลูก (ไร่)
  harvestMonth: number[];   // เดือนที่เก็บเกี่ยว
}

export interface PDPAConsent {
  consentGiven: boolean;
  consentDate: string;
  consentVersion: string;
  biometricConsent?: boolean;
  marketingConsent?: boolean;
}

export interface Beneficiary {
  name: string;
  relationship: string;
  idCardNo: string;
  phone?: string;
  percentage: number;       // สัดส่วนรับผลประโยชน์ (%)
}

// ========== ทุนเรือนหุ้น (Share Capital) ==========
export interface ShareAccount {
  _id?: string;
  coopId: string;
  memberId: string;

  totalShares: number;      // จำนวนหุ้นทั้งหมด
  totalAmount: number;      // มูลค่ารวม (บาท)

  transactions: ShareTransaction[];

  status: 'active' | 'locked' | 'withdrawn';
  createdAt: string;
  updatedAt: string;
}

export interface ShareTransaction {
  transactionId: string;
  type: 'buy' | 'monthly' | 'withdrawal' | 'dividend' | 'offset';
  shares: number;
  amount: number;
  date: string;
  description?: string;
  reference?: string;       // อ้างอิงเอกสาร
}

// ========== เงินฝาก (Deposits) ==========
export interface DepositAccount {
  _id?: string;
  coopId: string;
  memberId: string;
  accountNo: string;        // เลขที่บัญชี

  accountType: 'savings' | 'special_savings' | 'fixed';
  accountName: string;

  balance: number;          // ยอดคงเหลือ
  interestRate: number;     // อัตราดอกเบี้ย (%)
  accruedInterest: number;  // ดอกเบี้ยค้างรับ

  // สำหรับเงินฝากประจำ
  fixedDepositInfo?: FixedDepositInfo;

  // เงื่อนไขพิเศษ
  withdrawalRestriction?: WithdrawalRestriction;

  status: 'active' | 'closed' | 'dormant';
  openedDate: string;
  closedDate?: string;
  lastInterestDate?: string;

  createdAt: string;
  updatedAt: string;
}

export interface FixedDepositInfo {
  principalAmount: number;
  termMonths: number;
  maturityDate: string;
  autoRenew: boolean;
}

export interface WithdrawalRestriction {
  maxWithdrawalsPerMonth: number;
  minBalance: number;
  penaltyRate: number;      // ค่าธรรมเนียมถอนเกินกำหนด (%)
}

export interface DepositTransaction {
  _id?: string;
  coopId: string;
  accountNo: string;
  memberId: string;

  transactionId: string;
  type: 'deposit' | 'withdrawal' | 'interest' | 'transfer' | 'fee';
  amount: number;
  balanceAfter: number;

  date: string;
  description?: string;
  reference?: string;
  teller?: string;          // ผู้ทำรายการ

  createdAt: string;
}

// ========== สินเชื่อ (Loans) ==========
export interface LoanProduct {
  _id?: string;
  coopId: string;
  productId: string;
  productName: string;

  loanType: 'emergency' | 'ordinary' | 'special';
  description?: string;

  // เงื่อนไข
  minAmount: number;
  maxAmount: number;
  minTermMonths: number;
  maxTermMonths: number;
  interestRate: number;     // อัตราดอกเบี้ย (% ต่อปี)
  interestMethod: 'flat' | 'effective' | 'reducing_balance';

  // หลักประกัน
  collateralRequired: 'guarantor' | 'share' | 'property' | 'mixed';
  minGuarantors: number;

  // เงื่อนไขพิเศษ
  requiresApproval: boolean;
  approvalLevel: 'staff' | 'committee' | 'board';

  status: 'active' | 'inactive';
  createdAt: string;
  updatedAt: string;
}

export interface LoanApplication {
  _id?: string;
  coopId: string;
  applicationId: string;
  memberId: string;
  productId: string;

  // รายละเอียดคำขอ
  requestedAmount: number;
  requestedTermMonths: number;
  purpose: string;

  // การวิเคราะห์สินเชื่อ (5 C's)
  creditAnalysis: CreditAnalysis;

  // ผู้ค้ำประกัน
  guarantors: Guarantor[];

  // หลักทรัพย์จำนอง
  collaterals?: Collateral[];

  // สถานะการอนุมัติ
  status: 'draft' | 'pending' | 'approved' | 'rejected' | 'cancelled';
  approvalHistory: ApprovalHistory[];

  applicationDate: string;
  approvedDate?: string;
  approvedAmount?: number;
  approvedTermMonths?: number;
  approvedInterestRate?: number;

  createdAt: string;
  updatedAt: string;
}

export interface CreditAnalysis {
  // Character
  paymentHistory: 'excellent' | 'good' | 'fair' | 'poor';
  membershipYears: number;

  // Capacity
  monthlyIncome: number;
  monthlyExpense: number;
  existingDebt: number;
  dsr: number;              // Debt Service Ratio

  // Capital
  shareCapital: number;
  savings: number;
  otherAssets: number;

  // Collateral
  collateralValue: number;
  ltv: number;              // Loan to Value ratio

  // Conditions
  cropCycle?: string;
  expectedHarvest?: number;
  marketConditions?: string;

  creditScore?: number;
  recommendation?: string;
  analysisDate: string;
  analyzedBy: string;
}

export interface Guarantor {
  memberId: string;
  memberName: string;
  amount: number;           // จำนวนเงินที่ค้ำประกัน
  status: 'pending' | 'accepted' | 'rejected';
  acceptedDate?: string;
}

export interface Collateral {
  collateralId: string;
  type: 'land' | 'building' | 'vehicle' | 'equipment' | 'other';
  description: string;

  // สำหรับที่ดิน
  landInfo?: {
    titleDeedNo: string;
    area: number;           // ไร่
    location: string;
  };

  appraisedValue: number;
  appraisalDate: string;
  appraiser: string;

  documents: Document[];
}

export interface ApprovalHistory {
  date: string;
  action: 'submit' | 'review' | 'approve' | 'reject' | 'return';
  actor: string;
  role: string;
  comments?: string;
}

export interface LoanContract {
  _id?: string;
  coopId: string;
  contractNo: string;
  applicationId: string;
  memberId: string;

  // รายละเอียดสัญญา
  principalAmount: number;
  interestRate: number;
  termMonths: number;
  monthlyPayment: number;
  totalPayment: number;
  totalInterest: number;

  disbursementDate: string;
  firstPaymentDate: string;
  lastPaymentDate: string;

  // ยอดคงเหลือ
  principalBalance: number;
  interestAccrued: number;
  totalBalance: number;

  // สถานะ
  status: 'active' | 'closed' | 'default' | 'restructured';

  // การชำระเงิน
  paymentSchedule: PaymentSchedule[];

  createdAt: string;
  updatedAt: string;
}

export interface PaymentSchedule {
  installmentNo: number;
  dueDate: string;
  principalDue: number;
  interestDue: number;
  totalDue: number;

  // การชำระจริง
  paidAmount?: number;
  paidDate?: string;
  principalPaid?: number;
  interestPaid?: number;

  status: 'pending' | 'paid' | 'partial' | 'overdue';
  daysOverdue?: number;
}

export interface LoanPayment {
  _id?: string;
  coopId: string;
  paymentId: string;
  contractNo: string;
  memberId: string;

  paymentDate: string;
  amount: number;
  principalPaid: number;
  interestPaid: number;
  feesPaid?: number;

  paymentMethod: 'cash' | 'transfer' | 'deduct_savings' | 'deduct_produce';
  reference?: string;

  balanceAfter: number;

  createdAt: string;
}

// ========== ธุรกิจสินค้า (Business) ==========
export interface Product {
  _id?: string;
  coopId: string;
  productId: string;

  name: string;
  category: 'fertilizer' | 'pesticide' | 'seed' | 'equipment' | 'other';
  unit: string;             // หน่วยนับ

  costPrice: number;        // ราคาทุน
  sellingPrice: number;     // ราคาขาย
  memberPrice?: number;     // ราคาสมาชิก

  currentStock: number;
  minStock: number;         // จุดสั่งซื้อ

  // สำหรับสินค้าหมดอายุ
  hasExpiry: boolean;
  expiryMonths?: number;

  status: 'active' | 'inactive' | 'discontinued';
  createdAt: string;
  updatedAt: string;
}

export interface StockMovement {
  _id?: string;
  coopId: string;
  movementId: string;
  productId: string;

  type: 'in' | 'out' | 'adjust';
  reason: 'purchase' | 'sale' | 'return' | 'damage' | 'adjustment';

  quantity: number;
  unitCost?: number;
  totalValue: number;

  reference?: string;       // เลขที่เอกสาร
  batchNo?: string;
  expiryDate?: string;

  stockBefore: number;
  stockAfter: number;

  date: string;
  createdBy: string;
  createdAt: string;
}

export interface Sale {
  _id?: string;
  coopId: string;
  saleId: string;

  memberId?: string;        // ถ้าขายให้สมาชิก
  customerName?: string;    // ถ้าขายให้บุคคลทั่วไป

  items: SaleItem[];
  subtotal: number;
  discount: number;
  total: number;

  paymentType: 'cash' | 'credit' | 'transfer';

  // ถ้าขายเชื่อ
  creditInfo?: {
    dueDate: string;
    loanContractNo?: string;
  };

  status: 'completed' | 'pending' | 'cancelled';
  saleDate: string;
  createdBy: string;
  createdAt: string;
}

export interface SaleItem {
  productId: string;
  productName: string;
  quantity: number;
  unitPrice: number;
  discount: number;
  total: number;
}

// ========== การรวบรวมผลผลิต (Produce Collection) ==========
export interface ProduceCollection {
  _id?: string;
  coopId: string;
  collectionId: string;

  memberId: string;
  memberName: string;

  produceType: string;      // ชนิดผลผลิต

  // น้ำหนัก
  grossWeight: number;      // น้ำหนักรวม
  tareWeight: number;       // น้ำหนักภาชนะ
  netWeight: number;        // น้ำหนักสุทธิ

  // คุณภาพ (สำหรับยางพารา)
  qualityGrade?: string;
  drcPercent?: number;      // เปอร์เซ็นต์เนื้อยางแห้ง

  // ราคา
  pricePerUnit: number;
  totalAmount: number;

  // การชำระเงิน
  deductions: Deduction[];
  netPayment: number;
  paymentStatus: 'pending' | 'paid';
  paymentDate?: string;

  collectionDate: string;
  createdBy: string;
  createdAt: string;
}

export interface Deduction {
  type: 'share' | 'loan' | 'credit_sale' | 'fee' | 'other';
  description: string;
  amount: number;
  reference?: string;
}

// ========== บัญชี (Accounting) ==========
export interface ChartOfAccount {
  _id?: string;
  coopId: string;
  accountCode: string;
  accountName: string;
  accountType: 'asset' | 'liability' | 'equity' | 'revenue' | 'expense';
  parentCode?: string;
  level: number;
  isActive: boolean;
  createdAt: string;
}

export interface JournalEntry {
  _id?: string;
  coopId: string;
  journalId: string;

  date: string;
  description: string;
  reference?: string;
  sourceModule?: string;    // member, loan, deposit, sale, etc.
  sourceId?: string;

  entries: JournalLine[];

  status: 'draft' | 'posted' | 'reversed';
  postedDate?: string;
  postedBy?: string;

  createdAt: string;
  updatedAt: string;
}

export interface JournalLine {
  accountCode: string;
  accountName: string;
  debit: number;
  credit: number;
  description?: string;
}

// ========== การจัดสรรกำไร (Profit Allocation) ==========
export interface ProfitAllocation {
  _id?: string;
  coopId: string;
  fiscalYear: string;

  netProfit: number;

  allocations: AllocationItem[];

  // เงินปันผล
  dividendPool: number;
  dividendRate: number;

  // เฉลี่ยคืน
  patronagePool: number;
  patronageRate: number;

  status: 'draft' | 'approved' | 'distributed';
  approvedDate?: string;
  approvedBy?: string;

  createdAt: string;
  updatedAt: string;
}

export interface AllocationItem {
  type: 'reserve' | 'education' | 'federation' | 'dividend' | 'patronage' | 'other';
  name: string;
  percentage: number;
  amount: number;
}

export interface DividendDistribution {
  _id?: string;
  coopId: string;
  fiscalYear: string;
  memberId: string;

  // ปันผล
  averageShares: number;    // หุ้นเฉลี่ยถ่วงน้ำหนัก
  dividendRate: number;
  dividendAmount: number;

  // เฉลี่ยคืน
  totalInterestPaid: number;
  patronageRate: number;
  patronageAmount: number;

  totalDistribution: number;

  // การจ่าย
  paymentMethod: 'deposit' | 'share' | 'check';
  paymentDate?: string;
  status: 'calculated' | 'approved' | 'paid';

  createdAt: string;
}

// ========== ผู้ใช้งานและสิทธิ์ (Users & Permissions) ==========
export interface User {
  _id?: string;
  coopId: string;
  userId: string;
  username: string;
  passwordHash: string;

  employeeId?: string;
  fullName: string;
  email: string;
  phone?: string;

  role: UserRole;
  permissions: string[];

  status: 'active' | 'inactive' | 'locked';
  lastLogin?: string;

  createdAt: string;
  updatedAt: string;
}

export type UserRole =
  | 'admin'           // ผู้ดูแลระบบ
  | 'manager'         // ผู้จัดการ
  | 'accountant'      // พนักงานบัญชี
  | 'teller'          // พนักงานเคาน์เตอร์
  | 'credit_officer'  // เจ้าหน้าที่สินเชื่อ
  | 'auditor'         // ผู้สอบบัญชี (read-only)
  | 'member';         // สมาชิก (app)

// ========== API Response Types ==========
export interface ApiResponse<T = unknown> {
  status: 'success' | 'error';
  code: number;
  data?: T;
  count?: number;
  message?: string;
  error?: string;
}

// ========== Filter & Pagination ==========
export interface PaginationParams {
  page: number;
  limit: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface FilterParams {
  search?: string;
  status?: string;
  dateFrom?: string;
  dateTo?: string;
  [key: string]: string | number | boolean | undefined;
}
