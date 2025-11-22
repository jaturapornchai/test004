"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input, Textarea } from "@/components/ui/input";
import { Modal, ModalHeader, ModalContent, ModalFooter } from "@/components/ui/modal";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell, Badge } from "@/components/ui/table";
import { formatCurrency, formatDate } from "@/lib/utils";

// Mock deposit account data
const getDepositData = (id: string) => ({
  _id: id,
  accountId: id,
  coopId: "COOP001",
  memberId: "M001",
  memberName: "นายสมชาย ใจดี",
  accountType: "savings" as const,
  balance: 125000,
  interestRate: 1.5,
  accruedInterest: 156.25,
  openDate: "2020-01-15",
  lastInterestDate: "2024-01-01",
  status: "active" as const,
  createdAt: "2020-01-15T10:30:00Z",
  updatedAt: "2024-01-15T14:20:00Z",
});

// Mock transactions
const getTransactions = () => [
  { id: "TXN001", date: "2024-01-15", type: "deposit", description: "ฝากเงิน", amount: 5000, balance: 125000 },
  { id: "TXN002", date: "2024-01-01", type: "interest", description: "ดอกเบี้ยรายเดือน", amount: 156.25, balance: 120000 },
  { id: "TXN003", date: "2023-12-25", type: "withdrawal", description: "ถอนเงิน", amount: -3000, balance: 119843.75 },
  { id: "TXN004", date: "2023-12-15", type: "deposit", description: "ฝากเงิน", amount: 10000, balance: 122843.75 },
  { id: "TXN005", date: "2023-12-01", type: "interest", description: "ดอกเบี้ยรายเดือน", amount: 150.50, balance: 112843.75 },
  { id: "TXN006", date: "2023-11-20", type: "deposit", description: "ฝากเงิน", amount: 8000, balance: 112693.25 },
  { id: "TXN007", date: "2023-11-01", type: "interest", description: "ดอกเบี้ยรายเดือน", amount: 145.25, balance: 104693.25 },
];

const depositTypeLabels: Record<string, string> = {
  savings: "ออมทรัพย์",
  special_savings: "ออมทรัพย์พิเศษ",
  fixed: "ประจำ",
};

export default function DepositDetailPage() {
  const params = useParams();
  const accountId = params.id as string;

  const [showDepositModal, setShowDepositModal] = useState(false);
  const [showWithdrawModal, setShowWithdrawModal] = useState(false);

  const account = getDepositData(accountId);
  const transactions = getTransactions();

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case "deposit":
        return (
          <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
          </div>
        );
      case "withdrawal":
        return (
          <div className="w-8 h-8 rounded-full bg-red-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
            </svg>
          </div>
        );
      case "interest":
        return (
          <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Link href="/deposits">
            <Button variant="ghost" size="sm">
              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              กลับ
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-slate-800">บัญชีเงินฝาก {account.accountId}</h1>
            <p className="text-slate-500">{depositTypeLabels[account.accountType]} - {account.memberName}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => setShowDepositModal(true)}>
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            ฝากเงิน
          </Button>
          <Button variant="outline" onClick={() => setShowWithdrawModal(true)}>
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
            </svg>
            ถอนเงิน
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Account Summary */}
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>ข้อมูลบัญชี</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              {/* Balance */}
              <div className="p-6 bg-gradient-to-br from-green-500 to-green-600 rounded-xl text-white">
                <p className="text-green-100 text-sm">ยอดคงเหลือ</p>
                <p className="text-3xl font-bold mt-1">{formatCurrency(account.balance)}</p>
                <div className="mt-4 flex justify-between text-sm">
                  <span className="text-green-100">ดอกเบี้ยค้างรับ</span>
                  <span className="font-medium">{formatCurrency(account.accruedInterest)}</span>
                </div>
              </div>

              {/* Account Info */}
              <div className="space-y-4">
                <div className="flex justify-between">
                  <span className="text-slate-500">ประเภทบัญชี</span>
                  <span className="font-medium">{depositTypeLabels[account.accountType]}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">อัตราดอกเบี้ย</span>
                  <span className="font-medium text-green-600">{account.interestRate}% ต่อปี</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">วันที่เปิดบัญชี</span>
                  <span className="font-medium">{formatDate(account.openDate)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">คำนวณดอกเบี้ยล่าสุด</span>
                  <span className="font-medium">{formatDate(account.lastInterestDate)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">สถานะ</span>
                  <Badge variant="success">ใช้งาน</Badge>
                </div>
              </div>

              {/* Member Info */}
              <div className="pt-4 border-t">
                <p className="text-sm text-slate-500 mb-2">เจ้าของบัญชี</p>
                <Link href={`/members/${account.memberId}`} className="flex items-center gap-3 p-3 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white font-bold">
                    {account.memberName.charAt(0)}
                  </div>
                  <div>
                    <p className="font-medium text-blue-600">{account.memberName}</p>
                    <p className="text-sm text-slate-500">รหัส: {account.memberId}</p>
                  </div>
                </Link>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Transaction History */}
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>ประวัติรายการ</CardTitle>
            <Button variant="outline" size="sm">
              <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              ดาวน์โหลด
            </Button>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>วันที่</TableHead>
                  <TableHead>รายการ</TableHead>
                  <TableHead className="text-right">จำนวนเงิน</TableHead>
                  <TableHead className="text-right">คงเหลือ</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {transactions.map((txn) => (
                  <TableRow key={txn.id}>
                    <TableCell>{formatDate(txn.date)}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        {getTransactionIcon(txn.type)}
                        <span>{txn.description}</span>
                      </div>
                    </TableCell>
                    <TableCell className={`text-right font-medium ${txn.amount >= 0 ? "text-green-600" : "text-red-600"}`}>
                      {txn.amount >= 0 ? "+" : ""}{formatCurrency(txn.amount)}
                    </TableCell>
                    <TableCell className="text-right">{formatCurrency(txn.balance)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      {/* Deposit Modal */}
      <Modal isOpen={showDepositModal} onClose={() => setShowDepositModal(false)} size="md">
        <ModalHeader onClose={() => setShowDepositModal(false)}>ฝากเงิน</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="p-4 bg-green-50 rounded-lg">
              <p className="text-sm text-green-600">ยอดคงเหลือปัจจุบัน</p>
              <p className="text-2xl font-bold text-green-700">{formatCurrency(account.balance)}</p>
            </div>
            <Input label="จำนวนเงินที่ฝาก (บาท)" type="number" required />
            <Input label="วันที่ฝาก" type="date" defaultValue={new Date().toISOString().split("T")[0]} required />
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
            <div className="p-4 bg-blue-50 rounded-lg">
              <p className="text-sm text-blue-600">ยอดคงเหลือปัจจุบัน</p>
              <p className="text-2xl font-bold text-blue-700">{formatCurrency(account.balance)}</p>
            </div>
            <Input label="จำนวนเงินที่ถอน (บาท)" type="number" required />
            <Input label="วันที่ถอน" type="date" defaultValue={new Date().toISOString().split("T")[0]} required />
            <Textarea label="หมายเหตุ" rows={2} />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowWithdrawModal(false)}>ยกเลิก</Button>
          <Button variant="danger" onClick={() => setShowWithdrawModal(false)}>ถอนเงิน</Button>
        </ModalFooter>
      </Modal>
    </div>
  );
}
