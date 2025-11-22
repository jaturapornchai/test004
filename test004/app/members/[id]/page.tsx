"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent, StatsCard } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input, Select, Textarea } from "@/components/ui/input";
import { Modal, ModalHeader, ModalContent, ModalFooter, ConfirmDialog } from "@/components/ui/modal";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell, Badge } from "@/components/ui/table";
import { formatCurrency, formatNumber, formatDate, formatIdCard, formatPhone, getMemberStatusLabel, getMemberStatusColor } from "@/lib/utils";

// Mock member data
const getMemberData = (id: string) => ({
  _id: id,
  memberId: id,
  coopId: "COOP001",
  prefix: "นาย",
  firstName: "สมชาย",
  lastName: "ใจดี",
  idCard: "1234567890123",
  birthDate: "1980-05-15",
  gender: "male",
  phone: "0812345678",
  email: "somchai@email.com",
  occupation: "เกษตรกร",
  status: "active" as const,
  joinDate: "2020-01-15",
  memberType: "ordinary" as const,
  address: {
    houseNo: "123",
    moo: "4",
    village: "บ้านใหม่",
    road: "เพชรเกษม",
    subDistrict: "บางกระทึก",
    district: "สามพราน",
    province: "นครปฐม",
    postalCode: "73210",
  },
  farmingInfo: {
    farmType: ["rice", "vegetable"],
    landSize: 15.5,
    landUnit: "rai" as const,
    mainCrops: ["ข้าว", "ผัก"],
    irrigationType: "canal",
    farmAddress: "ต.บางกระทึก อ.สามพราน จ.นครปฐม",
  },
  beneficiaries: [
    { name: "นางสมศรี ใจดี", relationship: "ภรรยา", percentage: 50, phone: "0898765432" },
    { name: "นายสมหมาย ใจดี", relationship: "บุตร", percentage: 50, phone: "0887654321" },
  ],
  pdpaConsent: {
    consented: true,
    consentDate: "2020-01-15",
    version: "1.0",
  },
  createdAt: "2020-01-15T10:30:00Z",
  updatedAt: "2024-01-15T14:20:00Z",
});

// Mock financial summary
const getFinancialSummary = () => ({
  shares: {
    totalShares: 500,
    totalValue: 50000,
    monthlyPayment: 100,
    lastPaymentDate: "2024-01-05",
  },
  deposits: [
    { accountId: "DEP001", type: "savings", balance: 125000, interestRate: 1.5 },
    { accountId: "DEP002", type: "special_savings", balance: 50000, interestRate: 2.5 },
  ],
  loans: [
    { contractId: "LN001", type: "ordinary", principal: 200000, balance: 150000, status: "active", monthlyPayment: 8500 },
  ],
  totalDeposits: 175000,
  totalLoanBalance: 150000,
  netWorth: 75000, // shares + deposits - loans
});

// Mock transactions
const getRecentTransactions = () => [
  { id: "TXN001", date: "2024-01-15", type: "share_payment", description: "ชำระค่าหุ้นรายเดือน", amount: 100 },
  { id: "TXN002", date: "2024-01-10", type: "deposit", description: "ฝากเงินออมทรัพย์", amount: 5000 },
  { id: "TXN003", date: "2024-01-05", type: "loan_payment", description: "ชำระหนี้เงินกู้สามัญ", amount: -8500 },
  { id: "TXN004", date: "2024-01-01", type: "interest", description: "ดอกเบี้ยเงินฝาก", amount: 156.25 },
  { id: "TXN005", date: "2023-12-25", type: "withdrawal", description: "ถอนเงินออมทรัพย์", amount: -3000 },
];

const depositTypeLabels: Record<string, string> = {
  savings: "ออมทรัพย์",
  special_savings: "ออมทรัพย์พิเศษ",
  fixed: "ประจำ",
};

const loanTypeLabels: Record<string, string> = {
  emergency: "ฉุกเฉิน",
  ordinary: "สามัญ",
  special: "พิเศษ",
};

type TabType = "overview" | "financial" | "transactions" | "documents";

export default function MemberDetailPage() {
  const params = useParams();
  const memberId = params.id as string;

  const [activeTab, setActiveTab] = useState<TabType>("overview");
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDepositModal, setShowDepositModal] = useState(false);
  const [showWithdrawModal, setShowWithdrawModal] = useState(false);
  const [showBuyShareModal, setShowBuyShareModal] = useState(false);

  const member = getMemberData(memberId);
  const financial = getFinancialSummary();
  const transactions = getRecentTransactions();

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case "deposit":
      case "interest":
        return (
          <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
          </div>
        );
      case "withdrawal":
      case "loan_payment":
        return (
          <div className="w-8 h-8 rounded-full bg-red-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
            </svg>
          </div>
        );
      case "share_payment":
        return (
          <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
        );
      default:
        return (
          <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        );
    }
  };

  const renderOverviewTab = () => (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* Personal Information */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>ข้อมูลส่วนตัว</CardTitle>
          <Button variant="outline" size="sm" onClick={() => setShowEditModal(true)}>
            <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
            แก้ไข
          </Button>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center gap-4">
              <div className="w-20 h-20 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white text-2xl font-bold">
                {member.firstName.charAt(0)}
              </div>
              <div>
                <h3 className="text-xl font-semibold text-slate-800">
                  {member.prefix}{member.firstName} {member.lastName}
                </h3>
                <p className="text-slate-500">รหัสสมาชิก: {member.memberId}</p>
                <Badge variant={getMemberStatusColor(member.status) as "success" | "warning" | "danger"}>
                  {getMemberStatusLabel(member.status)}
                </Badge>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4 pt-4 border-t">
              <div>
                <p className="text-sm text-slate-500">เลขบัตรประชาชน</p>
                <p className="font-medium">{formatIdCard(member.idCard)}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">วันเกิด</p>
                <p className="font-medium">{formatDate(member.birthDate)}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">โทรศัพท์</p>
                <p className="font-medium">{formatPhone(member.phone)}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">อีเมล</p>
                <p className="font-medium">{member.email || "-"}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">อาชีพ</p>
                <p className="font-medium">{member.occupation}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">วันที่สมัคร</p>
                <p className="font-medium">{formatDate(member.joinDate)}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Address */}
      <Card>
        <CardHeader>
          <CardTitle>ที่อยู่</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <p className="text-slate-700">
              {member.address.houseNo} หมู่ {member.address.moo} {member.address.village && `บ้าน${member.address.village}`}
            </p>
            {member.address.road && <p className="text-slate-700">ถนน{member.address.road}</p>}
            <p className="text-slate-700">
              ตำบล{member.address.subDistrict} อำเภอ{member.address.district}
            </p>
            <p className="text-slate-700">
              จังหวัด{member.address.province} {member.address.postalCode}
            </p>
          </div>

          <div className="mt-6 pt-4 border-t">
            <h4 className="font-medium text-slate-800 mb-3">ข้อมูลการเกษตร</h4>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-slate-500">ขนาดพื้นที่</p>
                <p className="font-medium">{member.farmingInfo.landSize} ไร่</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">ระบบชลประทาน</p>
                <p className="font-medium">{member.farmingInfo.irrigationType === "canal" ? "คลองชลประทาน" : member.farmingInfo.irrigationType}</p>
              </div>
              <div className="col-span-2">
                <p className="text-sm text-slate-500">พืชที่ปลูก</p>
                <div className="flex flex-wrap gap-2 mt-1">
                  {member.farmingInfo.mainCrops.map((crop, index) => (
                    <span key={index} className="px-2 py-1 bg-green-100 text-green-700 rounded-full text-sm">
                      {crop}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Beneficiaries */}
      <Card>
        <CardHeader>
          <CardTitle>ผู้รับผลประโยชน์</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {member.beneficiaries.map((beneficiary, index) => (
              <div key={index} className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
                <div>
                  <p className="font-medium text-slate-800">{beneficiary.name}</p>
                  <p className="text-sm text-slate-500">{beneficiary.relationship} | {formatPhone(beneficiary.phone)}</p>
                </div>
                <div className="text-right">
                  <p className="text-lg font-semibold text-blue-600">{beneficiary.percentage}%</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* PDPA Consent */}
      <Card>
        <CardHeader>
          <CardTitle>ความยินยอม PDPA</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-3">
            {member.pdpaConsent.consented ? (
              <>
                <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
                  <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div>
                  <p className="font-medium text-green-700">ยินยอมแล้ว</p>
                  <p className="text-sm text-slate-500">
                    วันที่ยินยอม: {formatDate(member.pdpaConsent.consentDate)} | เวอร์ชัน: {member.pdpaConsent.version}
                  </p>
                </div>
              </>
            ) : (
              <>
                <div className="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center">
                  <svg className="w-5 h-5 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
                <div>
                  <p className="font-medium text-red-700">ยังไม่ยินยอม</p>
                  <p className="text-sm text-slate-500">กรุณาขอความยินยอมจากสมาชิก</p>
                </div>
              </>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderFinancialTab = () => (
    <div className="space-y-6">
      {/* Financial Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <StatsCard
          title="ทุนเรือนหุ้น"
          value={formatCurrency(financial.shares.totalValue)}
          subtitle={`${formatNumber(financial.shares.totalShares)} หุ้น`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          }
        />
        <StatsCard
          title="เงินฝากรวม"
          value={formatCurrency(financial.totalDeposits)}
          subtitle={`${financial.deposits.length} บัญชี`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
          }
        />
        <StatsCard
          title="หนี้คงค้าง"
          value={formatCurrency(financial.totalLoanBalance)}
          subtitle={`${financial.loans.length} สัญญา`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
        />
        <StatsCard
          title="ทรัพย์สินสุทธิ"
          value={formatCurrency(financial.netWorth)}
          subtitle="หุ้น + เงินฝาก - หนี้"
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          }
        />
      </div>

      {/* Shares */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>ทุนเรือนหุ้น</CardTitle>
          <Button size="sm" onClick={() => setShowBuyShareModal(true)}>
            <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            ซื้อหุ้นเพิ่ม
          </Button>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <p className="text-sm text-blue-600">จำนวนหุ้น</p>
              <p className="text-2xl font-bold text-blue-700">{formatNumber(financial.shares.totalShares)}</p>
            </div>
            <div className="p-4 bg-blue-50 rounded-lg">
              <p className="text-sm text-blue-600">มูลค่า</p>
              <p className="text-2xl font-bold text-blue-700">{formatCurrency(financial.shares.totalValue)}</p>
            </div>
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-600">ชำระรายเดือน</p>
              <p className="text-2xl font-bold text-slate-700">{formatCurrency(financial.shares.monthlyPayment)}</p>
            </div>
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-600">ชำระล่าสุด</p>
              <p className="text-2xl font-bold text-slate-700">{formatDate(financial.shares.lastPaymentDate)}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Deposits */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>บัญชีเงินฝาก</CardTitle>
          <div className="flex gap-2">
            <Button size="sm" onClick={() => setShowDepositModal(true)}>
              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              ฝากเงิน
            </Button>
            <Button variant="outline" size="sm" onClick={() => setShowWithdrawModal(true)}>
              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
              </svg>
              ถอนเงิน
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>เลขที่บัญชี</TableHead>
                <TableHead>ประเภท</TableHead>
                <TableHead className="text-right">ยอดคงเหลือ</TableHead>
                <TableHead className="text-right">ดอกเบี้ย (%/ปี)</TableHead>
                <TableHead className="text-right">จัดการ</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {financial.deposits.map((deposit) => (
                <TableRow key={deposit.accountId}>
                  <TableCell className="font-medium">{deposit.accountId}</TableCell>
                  <TableCell>{depositTypeLabels[deposit.type]}</TableCell>
                  <TableCell className="text-right font-semibold text-green-600">
                    {formatCurrency(deposit.balance)}
                  </TableCell>
                  <TableCell className="text-right">{deposit.interestRate}</TableCell>
                  <TableCell className="text-right">
                    <Link href={`/deposits/${deposit.accountId}`}>
                      <Button variant="ghost" size="sm">ดูรายละเอียด</Button>
                    </Link>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Loans */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>สัญญาสินเชื่อ</CardTitle>
          <Link href={`/loans?member=${memberId}`}>
            <Button variant="outline" size="sm">ดูทั้งหมด</Button>
          </Link>
        </CardHeader>
        <CardContent>
          {financial.loans.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>เลขที่สัญญา</TableHead>
                  <TableHead>ประเภท</TableHead>
                  <TableHead className="text-right">วงเงินกู้</TableHead>
                  <TableHead className="text-right">ยอดคงค้าง</TableHead>
                  <TableHead className="text-right">ชำระ/เดือน</TableHead>
                  <TableHead>สถานะ</TableHead>
                  <TableHead className="text-right">จัดการ</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {financial.loans.map((loan) => (
                  <TableRow key={loan.contractId}>
                    <TableCell className="font-medium">{loan.contractId}</TableCell>
                    <TableCell>{loanTypeLabels[loan.type]}</TableCell>
                    <TableCell className="text-right">{formatCurrency(loan.principal)}</TableCell>
                    <TableCell className="text-right font-semibold text-red-600">
                      {formatCurrency(loan.balance)}
                    </TableCell>
                    <TableCell className="text-right">{formatCurrency(loan.monthlyPayment)}</TableCell>
                    <TableCell>
                      <Badge variant={loan.status === "active" ? "success" : "danger"}>
                        {loan.status === "active" ? "ปกติ" : "ค้างชำระ"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <Link href={`/loans/${loan.contractId}`}>
                        <Button variant="ghost" size="sm">ดูรายละเอียด</Button>
                      </Link>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <div className="text-center py-8 text-slate-500">
              <svg className="w-12 h-12 mx-auto text-slate-300 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p>ไม่มีสัญญาสินเชื่อ</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );

  const renderTransactionsTab = () => (
    <Card>
      <CardHeader>
        <CardTitle>ประวัติรายการ</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {transactions.map((txn) => (
            <div key={txn.id} className="flex items-center gap-4 p-4 bg-slate-50 rounded-lg">
              {getTransactionIcon(txn.type)}
              <div className="flex-1">
                <p className="font-medium text-slate-800">{txn.description}</p>
                <p className="text-sm text-slate-500">{formatDate(txn.date)}</p>
              </div>
              <div className={`text-lg font-semibold ${txn.amount >= 0 ? "text-green-600" : "text-red-600"}`}>
                {txn.amount >= 0 ? "+" : ""}{formatCurrency(txn.amount)}
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );

  const renderDocumentsTab = () => (
    <Card>
      <CardHeader>
        <CardTitle>เอกสาร</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[
            { name: "สำเนาบัตรประชาชน", date: "2020-01-15", type: "id_card" },
            { name: "สำเนาทะเบียนบ้าน", date: "2020-01-15", type: "house_reg" },
            { name: "ใบสมัครสมาชิก", date: "2020-01-15", type: "application" },
            { name: "หนังสือยินยอม PDPA", date: "2020-01-15", type: "pdpa" },
          ].map((doc, index) => (
            <div key={index} className="flex items-center gap-3 p-4 border rounded-lg hover:border-blue-300 transition-colors cursor-pointer">
              <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                <svg className="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div className="flex-1">
                <p className="font-medium text-slate-800">{doc.name}</p>
                <p className="text-sm text-slate-500">{formatDate(doc.date)}</p>
              </div>
              <svg className="w-5 h-5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
            </div>
          ))}
        </div>

        <div className="mt-6 pt-6 border-t">
          <Button variant="outline">
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
            </svg>
            อัพโหลดเอกสาร
          </Button>
        </div>
      </CardContent>
    </Card>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Link href="/members">
            <Button variant="ghost" size="sm">
              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              กลับ
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-slate-800">
              {member.prefix}{member.firstName} {member.lastName}
            </h1>
            <p className="text-slate-500">รหัสสมาชิก: {member.memberId}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
            </svg>
            พิมพ์
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 border-b border-slate-200 pb-2">
        {([
          { id: "overview", label: "ข้อมูลทั่วไป" },
          { id: "financial", label: "ข้อมูลการเงิน" },
          { id: "transactions", label: "ประวัติรายการ" },
          { id: "documents", label: "เอกสาร" },
        ] as { id: TabType; label: string }[]).map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              activeTab === tab.id
                ? "bg-blue-50 text-blue-700"
                : "text-slate-600 hover:bg-slate-100"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      {activeTab === "overview" && renderOverviewTab()}
      {activeTab === "financial" && renderFinancialTab()}
      {activeTab === "transactions" && renderTransactionsTab()}
      {activeTab === "documents" && renderDocumentsTab()}

      {/* Deposit Modal */}
      <Modal isOpen={showDepositModal} onClose={() => setShowDepositModal(false)} size="md">
        <ModalHeader onClose={() => setShowDepositModal(false)}>ฝากเงิน</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Select label="บัญชี" required>
              <option value="">เลือกบัญชี</option>
              {financial.deposits.map((dep) => (
                <option key={dep.accountId} value={dep.accountId}>
                  {dep.accountId} - {depositTypeLabels[dep.type]} (ยอด: {formatCurrency(dep.balance)})
                </option>
              ))}
            </Select>
            <Input label="จำนวนเงิน (บาท)" type="number" required />
            <Textarea label="หมายเหตุ" rows={2} />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowDepositModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowDepositModal(false)}>ฝากเงิน</Button>
        </ModalFooter>
      </Modal>

      {/* Withdraw Modal */}
      <Modal isOpen={showWithdrawModal} onClose={() => setShowWithdrawModal(false)} size="md">
        <ModalHeader onClose={() => setShowWithdrawModal(false)}>ถอนเงิน</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Select label="บัญชี" required>
              <option value="">เลือกบัญชี</option>
              {financial.deposits.map((dep) => (
                <option key={dep.accountId} value={dep.accountId}>
                  {dep.accountId} - {depositTypeLabels[dep.type]} (ยอด: {formatCurrency(dep.balance)})
                </option>
              ))}
            </Select>
            <Input label="จำนวนเงิน (บาท)" type="number" required />
            <Textarea label="หมายเหตุ" rows={2} />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowWithdrawModal(false)}>ยกเลิก</Button>
          <Button variant="danger" onClick={() => setShowWithdrawModal(false)}>ถอนเงิน</Button>
        </ModalFooter>
      </Modal>

      {/* Buy Share Modal */}
      <Modal isOpen={showBuyShareModal} onClose={() => setShowBuyShareModal(false)} size="md">
        <ModalHeader onClose={() => setShowBuyShareModal(false)}>ซื้อหุ้นเพิ่ม</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <p className="text-sm text-blue-600">หุ้นปัจจุบัน</p>
              <p className="text-2xl font-bold text-blue-700">{formatNumber(financial.shares.totalShares)} หุ้น</p>
              <p className="text-sm text-blue-600">มูลค่า {formatCurrency(financial.shares.totalValue)}</p>
            </div>
            <Input label="จำนวนหุ้นที่ต้องการซื้อ" type="number" required hint="100 บาท/หุ้น" />
            <div className="p-3 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-600">ยอดชำระ</p>
              <p className="text-xl font-bold text-slate-800">฿0.00</p>
            </div>
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowBuyShareModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowBuyShareModal(false)}>ซื้อหุ้น</Button>
        </ModalFooter>
      </Modal>

      {/* Edit Modal */}
      <Modal isOpen={showEditModal} onClose={() => setShowEditModal(false)} size="lg">
        <ModalHeader onClose={() => setShowEditModal(false)}>แก้ไขข้อมูลสมาชิก</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <Select label="คำนำหน้า" defaultValue={member.prefix}>
                <option value="นาย">นาย</option>
                <option value="นาง">นาง</option>
                <option value="นางสาว">นางสาว</option>
              </Select>
              <Input label="ชื่อ" defaultValue={member.firstName} required />
              <Input label="นามสกุล" defaultValue={member.lastName} required />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <Input label="โทรศัพท์" defaultValue={member.phone} />
              <Input label="อีเมล" type="email" defaultValue={member.email} />
            </div>
            <Input label="อาชีพ" defaultValue={member.occupation} />
            <Select label="สถานะ" defaultValue={member.status}>
              <option value="active">ใช้งาน</option>
              <option value="inactive">ไม่ใช้งาน</option>
              <option value="suspended">ระงับ</option>
              <option value="resigned">ลาออก</option>
            </Select>
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowEditModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowEditModal(false)}>บันทึก</Button>
        </ModalFooter>
      </Modal>
    </div>
  );
}
