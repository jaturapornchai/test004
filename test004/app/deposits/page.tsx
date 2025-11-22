'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input, Select } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, StatsCard } from '@/components/ui/card';
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell, Badge, Pagination } from '@/components/ui/table';
import { Modal, ModalHeader, ModalContent, ModalFooter } from '@/components/ui/modal';
import { formatCurrency, formatDate, depositTypeLabels } from '@/lib/utils';

// Mock data
const mockDeposits = [
  {
    accountNo: 'SAV001-0001',
    memberId: 'M001',
    memberName: 'นายสมชาย ใจดี',
    accountType: 'savings',
    balance: 125000,
    interestRate: 1.5,
    lastTransaction: '2024-11-15',
    status: 'active',
  },
  {
    accountNo: 'SAV001-0002',
    memberId: 'M002',
    memberName: 'นางสมหญิง รักษ์ธรรม',
    accountType: 'savings',
    balance: 250000,
    interestRate: 1.5,
    lastTransaction: '2024-11-18',
    status: 'active',
  },
  {
    accountNo: 'SPS001-0001',
    memberId: 'M003',
    memberName: 'นายวิชัย พัฒนา',
    accountType: 'special_savings',
    balance: 380000,
    interestRate: 2.5,
    lastTransaction: '2024-11-10',
    status: 'active',
  },
  {
    accountNo: 'FIX001-0001',
    memberId: 'M001',
    memberName: 'นายสมชาย ใจดี',
    accountType: 'fixed',
    balance: 100000,
    interestRate: 3.0,
    maturityDate: '2025-06-15',
    lastTransaction: '2024-06-15',
    status: 'active',
  },
];

export default function DepositsPage() {
  const [search, setSearch] = useState('');
  const [typeFilter, setTypeFilter] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [showDepositModal, setShowDepositModal] = useState(false);
  const [showWithdrawModal, setShowWithdrawModal] = useState(false);
  const [transactionType, setTransactionType] = useState<'deposit' | 'withdraw'>('deposit');

  const filteredDeposits = mockDeposits.filter((deposit) => {
    const matchSearch =
      deposit.accountNo.toLowerCase().includes(search.toLowerCase()) ||
      deposit.memberName.toLowerCase().includes(search.toLowerCase());
    const matchType = !typeFilter || deposit.accountType === typeFilter;
    return matchSearch && matchType;
  });

  const totalDeposits = mockDeposits.reduce((sum, d) => sum + d.balance, 0);
  const savingsTotal = mockDeposits.filter(d => d.accountType === 'savings').reduce((sum, d) => sum + d.balance, 0);
  const specialTotal = mockDeposits.filter(d => d.accountType === 'special_savings').reduce((sum, d) => sum + d.balance, 0);
  const fixedTotal = mockDeposits.filter(d => d.accountType === 'fixed').reduce((sum, d) => sum + d.balance, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">ระบบเงินฝาก</h1>
          <p className="text-slate-600 mt-1">จัดการบัญชีเงินฝากและทำธุรกรรม ฝาก-ถอน</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" onClick={() => { setTransactionType('withdraw'); setShowWithdrawModal(true); }}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 19l-7-7 7-7m8 14l-7-7 7-7" />
            </svg>
            ถอนเงิน
          </Button>
          <Button onClick={() => { setTransactionType('deposit'); setShowDepositModal(true); }}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            รับฝากเงิน
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="เงินฝากรวมทั้งหมด"
          value={formatCurrency(totalDeposits)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
          }
          color="blue"
        />
        <StatsCard
          title="ออมทรัพย์"
          value={formatCurrency(savingsTotal)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2z" />
            </svg>
          }
          color="green"
        />
        <StatsCard
          title="ออมทรัพย์พิเศษ"
          value={formatCurrency(specialTotal)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          }
          color="purple"
        />
        <StatsCard
          title="เงินฝากประจำ"
          value={formatCurrency(fixedTotal)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
          color="yellow"
        />
      </div>

      {/* Search & Filter */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col lg:flex-row gap-4">
            <div className="flex-1">
              <Input
                placeholder="ค้นหาเลขบัญชี หรือ ชื่อสมาชิก..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                leftIcon={
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                }
              />
            </div>
            <div className="w-full lg:w-48">
              <Select
                options={[
                  { value: '', label: 'ประเภททั้งหมด' },
                  { value: 'savings', label: 'ออมทรัพย์' },
                  { value: 'special_savings', label: 'ออมทรัพย์พิเศษ' },
                  { value: 'fixed', label: 'เงินฝากประจำ' },
                ]}
                value={typeFilter}
                onChange={setTypeFilter}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardHeader>
          <CardTitle subtitle={`แสดง ${filteredDeposits.length} บัญชี`}>
            รายการบัญชีเงินฝาก
          </CardTitle>
        </CardHeader>
        <Table>
          <TableHead>
            <TableRow>
              <TableHeader>เลขบัญชี</TableHeader>
              <TableHeader>ชื่อเจ้าของบัญชี</TableHeader>
              <TableHeader>ประเภท</TableHeader>
              <TableHeader className="text-right">ยอดคงเหลือ</TableHeader>
              <TableHeader className="text-right">อัตราดอกเบี้ย</TableHeader>
              <TableHeader>ธุรกรรมล่าสุด</TableHeader>
              <TableHeader>สถานะ</TableHeader>
              <TableHeader className="text-center">ทำรายการ</TableHeader>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredDeposits.map((deposit) => (
              <TableRow key={deposit.accountNo}>
                <TableCell className="font-medium font-mono">{deposit.accountNo}</TableCell>
                <TableCell>
                  <Link href={`/members/${deposit.memberId}`} className="text-blue-600 hover:underline">
                    {deposit.memberName}
                  </Link>
                </TableCell>
                <TableCell>
                  <Badge variant={
                    deposit.accountType === 'savings' ? 'success' :
                    deposit.accountType === 'special_savings' ? 'info' : 'warning'
                  }>
                    {depositTypeLabels[deposit.accountType]}
                  </Badge>
                </TableCell>
                <TableCell className="text-right font-semibold">{formatCurrency(deposit.balance)}</TableCell>
                <TableCell className="text-right">{deposit.interestRate}% ต่อปี</TableCell>
                <TableCell>{formatDate(deposit.lastTransaction)}</TableCell>
                <TableCell>
                  <Badge variant="success">ใช้งาน</Badge>
                </TableCell>
                <TableCell>
                  <div className="flex items-center justify-center gap-2">
                    <Button variant="ghost" size="sm" onClick={() => { setTransactionType('deposit'); setShowDepositModal(true); }}>
                      <svg className="w-4 h-4 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
                      </svg>
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => { setTransactionType('withdraw'); setShowWithdrawModal(true); }}>
                      <svg className="w-4 h-4 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M20 12H4" />
                      </svg>
                    </Button>
                    <Link href={`/deposits/${deposit.accountNo}`}>
                      <Button variant="ghost" size="sm">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                      </Button>
                    </Link>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <Pagination
          currentPage={currentPage}
          totalPages={1}
          totalItems={filteredDeposits.length}
          onPageChange={setCurrentPage}
        />
      </Card>

      {/* Deposit Modal */}
      <Modal isOpen={showDepositModal} onClose={() => setShowDepositModal(false)} size="md">
        <ModalHeader onClose={() => setShowDepositModal(false)}>
          รับฝากเงิน
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="เลขบัญชี" placeholder="ค้นหาเลขบัญชี..." required />
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-500">ชื่อเจ้าของบัญชี</p>
              <p className="font-medium text-slate-900">-</p>
              <p className="text-sm text-slate-500 mt-2">ยอดคงเหลือปัจจุบัน</p>
              <p className="font-semibold text-lg text-slate-900">฿0.00</p>
            </div>
            <Input label="จำนวนเงินที่ฝาก" type="number" placeholder="0.00" required />
            <Input label="หมายเหตุ" placeholder="ระบุหมายเหตุ (ถ้ามี)" />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowDepositModal(false)}>
            ยกเลิก
          </Button>
          <Button variant="secondary" onClick={() => setShowDepositModal(false)}>
            ยืนยันการฝาก
          </Button>
        </ModalFooter>
      </Modal>

      {/* Withdraw Modal */}
      <Modal isOpen={showWithdrawModal} onClose={() => setShowWithdrawModal(false)} size="md">
        <ModalHeader onClose={() => setShowWithdrawModal(false)}>
          ถอนเงิน
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="เลขบัญชี" placeholder="ค้นหาเลขบัญชี..." required />
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-500">ชื่อเจ้าของบัญชี</p>
              <p className="font-medium text-slate-900">-</p>
              <p className="text-sm text-slate-500 mt-2">ยอดคงเหลือปัจจุบัน</p>
              <p className="font-semibold text-lg text-slate-900">฿0.00</p>
            </div>
            <Input label="จำนวนเงินที่ถอน" type="number" placeholder="0.00" required />
            <Input label="หมายเหตุ" placeholder="ระบุหมายเหตุ (ถ้ามี)" />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowWithdrawModal(false)}>
            ยกเลิก
          </Button>
          <Button variant="danger" onClick={() => setShowWithdrawModal(false)}>
            ยืนยันการถอน
          </Button>
        </ModalFooter>
      </Modal>
    </div>
  );
}
