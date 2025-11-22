"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input, Textarea } from "@/components/ui/input";
import { Modal, ModalHeader, ModalContent, ModalFooter } from "@/components/ui/modal";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell, Badge } from "@/components/ui/table";
import { formatCurrency, formatNumber, formatDate } from "@/lib/utils";

// Mock loan data
const getLoanData = (id: string) => ({
  _id: id,
  contractId: id,
  coopId: "COOP001",
  memberId: "M001",
  memberName: "นายสมชาย ใจดี",
  loanType: "ordinary" as const,
  purpose: "ซื้อเมล็ดพันธุ์และปุ๋ยสำหรับฤดูกาลเพาะปลูก",
  principal: 200000,
  interestRate: 9,
  term: 24,
  startDate: "2023-06-01",
  endDate: "2025-05-31",
  monthlyPayment: 9200,
  totalInterest: 20800,
  totalPayment: 220800,
  paidPrincipal: 50000,
  paidInterest: 12000,
  balance: 150000,
  status: "active" as const,
  guarantors: [
    { name: "นายวิชัย พัฒนา", memberId: "M002", amount: 100000 },
    { name: "นางสมศรี รักดี", memberId: "M003", amount: 100000 },
  ],
  collateral: "โฉนดที่ดิน เลขที่ 12345",
  approvedBy: "นายประสิทธิ์ มั่นคง",
  approvedDate: "2023-05-28",
  createdAt: "2023-05-25T10:30:00Z",
  updatedAt: "2024-01-15T14:20:00Z",
});

// Mock payment schedule
const getPaymentSchedule = () => [
  { no: 1, dueDate: "2023-07-01", principal: 7200, interest: 1500, total: 8700, paid: true, paidDate: "2023-07-01", balance: 192800 },
  { no: 2, dueDate: "2023-08-01", principal: 7254, interest: 1446, total: 8700, paid: true, paidDate: "2023-08-01", balance: 185546 },
  { no: 3, dueDate: "2023-09-01", principal: 7309, interest: 1391, total: 8700, paid: true, paidDate: "2023-09-05", balance: 178237 },
  { no: 4, dueDate: "2023-10-01", principal: 7364, interest: 1336, total: 8700, paid: true, paidDate: "2023-10-02", balance: 170873 },
  { no: 5, dueDate: "2023-11-01", principal: 7419, interest: 1281, total: 8700, paid: true, paidDate: "2023-11-01", balance: 163454 },
  { no: 6, dueDate: "2023-12-01", principal: 7475, interest: 1225, total: 8700, paid: true, paidDate: "2023-12-04", balance: 155979 },
  { no: 7, dueDate: "2024-01-01", principal: 7531, interest: 1169, total: 8700, paid: true, paidDate: "2024-01-05", balance: 148448 },
  { no: 8, dueDate: "2024-02-01", principal: 7588, interest: 1112, total: 8700, paid: false, paidDate: null, balance: 150000 },
  { no: 9, dueDate: "2024-03-01", principal: 7645, interest: 1055, total: 8700, paid: false, paidDate: null, balance: 150000 },
];

// Mock payment history
const getPaymentHistory = () => [
  { id: "PAY007", date: "2024-01-05", amount: 8700, principal: 7531, interest: 1169, method: "เงินสด", receiptNo: "RCP-2024-0105" },
  { id: "PAY006", date: "2023-12-04", amount: 8700, principal: 7475, interest: 1225, method: "โอนเงิน", receiptNo: "RCP-2023-1204" },
  { id: "PAY005", date: "2023-11-01", amount: 8700, principal: 7419, interest: 1281, method: "หักบัญชี", receiptNo: "RCP-2023-1101" },
  { id: "PAY004", date: "2023-10-02", amount: 8700, principal: 7364, interest: 1336, method: "เงินสด", receiptNo: "RCP-2023-1002" },
  { id: "PAY003", date: "2023-09-05", amount: 8700, principal: 7309, interest: 1391, method: "เงินสด", receiptNo: "RCP-2023-0905" },
];

const loanTypeLabels: Record<string, string> = {
  emergency: "เงินกู้ฉุกเฉิน",
  ordinary: "เงินกู้สามัญ",
  special: "เงินกู้พิเศษ",
};

const loanStatusLabels: Record<string, string> = {
  active: "ปกติ",
  overdue: "ค้างชำระ",
  closed: "ปิดบัญชี",
  default: "ผิดนัด",
};

type TabType = "overview" | "schedule" | "history";

export default function LoanDetailPage() {
  const params = useParams();
  const loanId = params.id as string;

  const [activeTab, setActiveTab] = useState<TabType>("overview");
  const [showPaymentModal, setShowPaymentModal] = useState(false);

  const loan = getLoanData(loanId);
  const schedule = getPaymentSchedule();
  const history = getPaymentHistory();

  const progressPercent = ((loan.principal - loan.balance) / loan.principal) * 100;

  const renderOverviewTab = () => (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* Loan Summary */}
      <Card>
        <CardHeader>
          <CardTitle>ข้อมูลสัญญา</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <Badge variant={loan.status === "active" ? "success" : loan.status === "overdue" ? "danger" : "default"}>
                {loanStatusLabels[loan.status]}
              </Badge>
              <span className="text-sm text-slate-500">{loanTypeLabels[loan.loanType]}</span>
            </div>

            {/* Progress Bar */}
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-slate-600">ความคืบหน้าการชำระ</span>
                <span className="font-medium">{progressPercent.toFixed(1)}%</span>
              </div>
              <div className="h-3 bg-slate-100 rounded-full overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-blue-500 to-blue-600 rounded-full transition-all"
                  style={{ width: `${progressPercent}%` }}
                />
              </div>
              <div className="flex justify-between text-xs text-slate-500">
                <span>ชำระแล้ว {formatCurrency(loan.paidPrincipal)}</span>
                <span>คงเหลือ {formatCurrency(loan.balance)}</span>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4 pt-4 border-t">
              <div>
                <p className="text-sm text-slate-500">วงเงินกู้</p>
                <p className="text-xl font-bold text-slate-800">{formatCurrency(loan.principal)}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">ยอดคงค้าง</p>
                <p className="text-xl font-bold text-red-600">{formatCurrency(loan.balance)}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">อัตราดอกเบี้ย</p>
                <p className="font-medium">{loan.interestRate}% ต่อปี</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">ระยะเวลา</p>
                <p className="font-medium">{loan.term} เดือน</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">วันที่เริ่มสัญญา</p>
                <p className="font-medium">{formatDate(loan.startDate)}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">วันครบกำหนด</p>
                <p className="font-medium">{formatDate(loan.endDate)}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Member & Payment Info */}
      <Card>
        <CardHeader>
          <CardTitle>ข้อมูลผู้กู้</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white text-xl font-bold">
                {loan.memberName.charAt(0)}
              </div>
              <div>
                <Link href={`/members/${loan.memberId}`} className="font-medium text-blue-600 hover:underline">
                  {loan.memberName}
                </Link>
                <p className="text-sm text-slate-500">รหัสสมาชิก: {loan.memberId}</p>
              </div>
            </div>

            <div className="pt-4 border-t">
              <p className="text-sm text-slate-500 mb-1">วัตถุประสงค์</p>
              <p className="text-slate-700">{loan.purpose}</p>
            </div>

            <div className="pt-4 border-t">
              <p className="text-sm text-slate-500 mb-2">หลักประกัน</p>
              <p className="text-slate-700">{loan.collateral || "-"}</p>
            </div>

            <div className="pt-4 border-t grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-slate-500">ผู้อนุมัติ</p>
                <p className="font-medium">{loan.approvedBy}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">วันที่อนุมัติ</p>
                <p className="font-medium">{formatDate(loan.approvedDate)}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Payment Summary */}
      <Card>
        <CardHeader>
          <CardTitle>สรุปการชำระ</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">ค่างวดต่อเดือน</span>
              <span className="font-semibold text-lg">{formatCurrency(loan.monthlyPayment)}</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-green-50 rounded-lg">
              <span className="text-green-600">เงินต้นชำระแล้ว</span>
              <span className="font-semibold text-green-700">{formatCurrency(loan.paidPrincipal)}</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-amber-50 rounded-lg">
              <span className="text-amber-600">ดอกเบี้ยชำระแล้ว</span>
              <span className="font-semibold text-amber-700">{formatCurrency(loan.paidInterest)}</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-blue-50 rounded-lg">
              <span className="text-blue-600">ยอดรวมที่ต้องชำระ</span>
              <span className="font-semibold text-blue-700">{formatCurrency(loan.totalPayment)}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Guarantors */}
      <Card>
        <CardHeader>
          <CardTitle>ผู้ค้ำประกัน</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {loan.guarantors.map((guarantor, index) => (
              <div key={index} className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-slate-300 flex items-center justify-center text-white font-bold">
                    {guarantor.name.charAt(0)}
                  </div>
                  <div>
                    <Link href={`/members/${guarantor.memberId}`} className="font-medium text-blue-600 hover:underline">
                      {guarantor.name}
                    </Link>
                    <p className="text-sm text-slate-500">รหัส: {guarantor.memberId}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-semibold">{formatCurrency(guarantor.amount)}</p>
                  <p className="text-xs text-slate-500">วงเงินค้ำ</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderScheduleTab = () => (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>ตารางผ่อนชำระ</CardTitle>
        <Button variant="outline" size="sm">
          <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
          </svg>
          พิมพ์ตาราง
        </Button>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="text-center">งวดที่</TableHead>
              <TableHead>วันครบกำหนด</TableHead>
              <TableHead className="text-right">เงินต้น</TableHead>
              <TableHead className="text-right">ดอกเบี้ย</TableHead>
              <TableHead className="text-right">รวม</TableHead>
              <TableHead className="text-right">คงเหลือ</TableHead>
              <TableHead>สถานะ</TableHead>
              <TableHead>วันที่ชำระ</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {schedule.map((item) => (
              <TableRow key={item.no} className={!item.paid && new Date(item.dueDate) < new Date() ? "bg-red-50" : ""}>
                <TableCell className="text-center font-medium">{item.no}</TableCell>
                <TableCell>{formatDate(item.dueDate)}</TableCell>
                <TableCell className="text-right">{formatCurrency(item.principal)}</TableCell>
                <TableCell className="text-right">{formatCurrency(item.interest)}</TableCell>
                <TableCell className="text-right font-medium">{formatCurrency(item.total)}</TableCell>
                <TableCell className="text-right">{formatCurrency(item.balance)}</TableCell>
                <TableCell>
                  {item.paid ? (
                    <Badge variant="success">ชำระแล้ว</Badge>
                  ) : new Date(item.dueDate) < new Date() ? (
                    <Badge variant="danger">ค้างชำระ</Badge>
                  ) : (
                    <Badge variant="warning">รอชำระ</Badge>
                  )}
                </TableCell>
                <TableCell>{item.paidDate ? formatDate(item.paidDate) : "-"}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );

  const renderHistoryTab = () => (
    <Card>
      <CardHeader>
        <CardTitle>ประวัติการชำระ</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>วันที่</TableHead>
              <TableHead className="text-right">เงินต้น</TableHead>
              <TableHead className="text-right">ดอกเบี้ย</TableHead>
              <TableHead className="text-right">รวม</TableHead>
              <TableHead>วิธีชำระ</TableHead>
              <TableHead>เลขที่ใบเสร็จ</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {history.map((item) => (
              <TableRow key={item.id}>
                <TableCell>{formatDate(item.date)}</TableCell>
                <TableCell className="text-right">{formatCurrency(item.principal)}</TableCell>
                <TableCell className="text-right">{formatCurrency(item.interest)}</TableCell>
                <TableCell className="text-right font-semibold">{formatCurrency(item.amount)}</TableCell>
                <TableCell>{item.method}</TableCell>
                <TableCell className="text-blue-600">{item.receiptNo}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Link href="/loans">
            <Button variant="ghost" size="sm">
              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              กลับ
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-slate-800">สัญญาเงินกู้ {loan.contractId}</h1>
            <p className="text-slate-500">{loanTypeLabels[loan.loanType]} - {loan.memberName}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => setShowPaymentModal(true)}>
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            รับชำระ
          </Button>
          <Button variant="outline">
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
            </svg>
            พิมพ์สัญญา
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 border-b border-slate-200 pb-2">
        {([
          { id: "overview", label: "ภาพรวม" },
          { id: "schedule", label: "ตารางผ่อนชำระ" },
          { id: "history", label: "ประวัติการชำระ" },
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
      {activeTab === "schedule" && renderScheduleTab()}
      {activeTab === "history" && renderHistoryTab()}

      {/* Payment Modal */}
      <Modal isOpen={showPaymentModal} onClose={() => setShowPaymentModal(false)} size="md">
        <ModalHeader onClose={() => setShowPaymentModal(false)}>รับชำระเงินกู้</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <div className="flex justify-between mb-2">
                <span className="text-blue-600">ยอดคงค้าง</span>
                <span className="font-bold text-blue-700">{formatCurrency(loan.balance)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-blue-600">ค่างวดต่อเดือน</span>
                <span className="font-bold text-blue-700">{formatCurrency(loan.monthlyPayment)}</span>
              </div>
            </div>
            <Input label="จำนวนเงินที่ชำระ (บาท)" type="number" defaultValue={loan.monthlyPayment} required />
            <Input label="วันที่ชำระ" type="date" defaultValue={new Date().toISOString().split("T")[0]} required />
            <Textarea label="หมายเหตุ" rows={2} />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowPaymentModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowPaymentModal(false)}>บันทึกการชำระ</Button>
        </ModalFooter>
      </Modal>
    </div>
  );
}
