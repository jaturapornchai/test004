'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input, Select } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, StatsCard } from '@/components/ui/card';
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell, Badge, Pagination } from '@/components/ui/table';
import { Modal, ModalHeader, ModalContent, ModalFooter } from '@/components/ui/modal';
import { formatCurrency, formatDate } from '@/lib/utils';

// Mock data - ข้อมูลทุนเรือนหุ้นสมาชิก
const mockShares = [
  {
    memberId: 'M001',
    memberNo: '0001',
    memberName: 'นายสมชาย ใจดี',
    totalShares: 500,
    shareValue: 100,
    totalAmount: 50000,
    monthlyContribution: 500,
    joinDate: '2020-03-15',
    lastPaymentDate: '2024-11-01',
    status: 'active',
  },
  {
    memberId: 'M002',
    memberNo: '0002',
    memberName: 'นางสมหญิง รักษ์ธรรม',
    totalShares: 800,
    shareValue: 100,
    totalAmount: 80000,
    monthlyContribution: 1000,
    joinDate: '2019-07-22',
    lastPaymentDate: '2024-11-01',
    status: 'active',
  },
  {
    memberId: 'M003',
    memberNo: '0003',
    memberName: 'นายวิชัย พัฒนา',
    totalShares: 1200,
    shareValue: 100,
    totalAmount: 120000,
    monthlyContribution: 1500,
    joinDate: '2018-01-10',
    lastPaymentDate: '2024-11-01',
    status: 'active',
  },
  {
    memberId: 'M004',
    memberNo: '0004',
    memberName: 'นางสาวจันทร์ แสงงาม',
    totalShares: 300,
    shareValue: 100,
    totalAmount: 30000,
    monthlyContribution: 0,
    joinDate: '2021-05-20',
    lastPaymentDate: '2024-08-01',
    status: 'resigned',
  },
];

// ประวัติการทำรายการหุ้น
const mockTransactions = [
  { date: '2024-11-01', type: 'monthly', memberId: 'M001', memberName: 'นายสมชาย ใจดี', shares: 5, amount: 500, description: 'ค่าหุ้นประจำเดือน พ.ย. 2567' },
  { date: '2024-11-01', type: 'monthly', memberId: 'M002', memberName: 'นางสมหญิง รักษ์ธรรม', shares: 10, amount: 1000, description: 'ค่าหุ้นประจำเดือน พ.ย. 2567' },
  { date: '2024-11-01', type: 'monthly', memberId: 'M003', memberName: 'นายวิชัย พัฒนา', shares: 15, amount: 1500, description: 'ค่าหุ้นประจำเดือน พ.ย. 2567' },
  { date: '2024-10-15', type: 'buy', memberId: 'M001', memberName: 'นายสมชาย ใจดี', shares: 50, amount: 5000, description: 'ซื้อหุ้นเพิ่ม' },
  { date: '2024-10-01', type: 'monthly', memberId: 'M001', memberName: 'นายสมชาย ใจดี', shares: 5, amount: 500, description: 'ค่าหุ้นประจำเดือน ต.ค. 2567' },
];

export default function SharesPage() {
  const [search, setSearch] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [showBuyModal, setShowBuyModal] = useState(false);
  const [showMonthlyModal, setShowMonthlyModal] = useState(false);
  const [activeTab, setActiveTab] = useState<'accounts' | 'transactions'>('accounts');

  const filteredShares = mockShares.filter((share) =>
    share.memberName.toLowerCase().includes(search.toLowerCase()) ||
    share.memberNo.includes(search)
  );

  const totalShareCapital = mockShares.filter(s => s.status === 'active').reduce((sum, s) => sum + s.totalAmount, 0);
  const totalMembers = mockShares.filter(s => s.status === 'active').length;
  const totalShares = mockShares.filter(s => s.status === 'active').reduce((sum, s) => sum + s.totalShares, 0);
  const avgSharePerMember = totalMembers > 0 ? totalShareCapital / totalMembers : 0;

  const transactionTypeLabels: Record<string, string> = {
    buy: 'ซื้อหุ้นเพิ่ม',
    monthly: 'ค่าหุ้นรายเดือน',
    withdrawal: 'ถอนหุ้น',
    dividend: 'เงินปันผล',
    offset: 'หักลบหนี้',
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">ทุนเรือนหุ้น</h1>
          <p className="text-slate-600 mt-1">จัดการหุ้นสมาชิก ซื้อหุ้นเพิ่ม และส่งค่าหุ้นรายเดือน</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" onClick={() => setShowMonthlyModal(true)}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            ค่าหุ้นรายเดือน
          </Button>
          <Button onClick={() => setShowBuyModal(true)}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            ซื้อหุ้นเพิ่ม
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="ทุนเรือนหุ้นรวม"
          value={formatCurrency(totalShareCapital)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2z" />
            </svg>
          }
          color="green"
        />
        <StatsCard
          title="จำนวนหุ้นรวม"
          value={totalShares.toLocaleString()}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" />
            </svg>
          }
          color="blue"
        />
        <StatsCard
          title="สมาชิกถือหุ้น"
          value={`${totalMembers} คน`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          }
          color="purple"
        />
        <StatsCard
          title="เฉลี่ยต่อคน"
          value={formatCurrency(avgSharePerMember)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
          }
          color="yellow"
        />
      </div>

      {/* Info Banner */}
      <div className="bg-gradient-to-r from-emerald-50 to-teal-50 border border-emerald-200 rounded-xl p-4">
        <div className="flex items-start gap-3">
          <div className="p-2 bg-emerald-100 rounded-lg">
            <svg className="w-5 h-5 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h3 className="font-medium text-emerald-900">ข้อบังคับสหกรณ์</h3>
            <p className="text-sm text-emerald-700 mt-1">
              มูลค่าหุ้นละ <strong>100 บาท</strong> |
              เพดานการถือหุ้น <strong>ไม่เกิน 1/5 ของทุนทั้งหมด หรือ 2,500,000 บาท</strong> |
              ค่าหุ้นรายเดือนขั้นต่ำ <strong>100 บาท</strong>
            </p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-slate-200">
        <div className="flex gap-4">
          <button
            onClick={() => setActiveTab('accounts')}
            className={`pb-3 px-1 font-medium text-sm border-b-2 transition-colors ${
              activeTab === 'accounts'
                ? 'border-emerald-600 text-emerald-600'
                : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            บัญชีหุ้นสมาชิก
          </button>
          <button
            onClick={() => setActiveTab('transactions')}
            className={`pb-3 px-1 font-medium text-sm border-b-2 transition-colors ${
              activeTab === 'transactions'
                ? 'border-emerald-600 text-emerald-600'
                : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            ประวัติการทำรายการ
          </button>
        </div>
      </div>

      {activeTab === 'accounts' ? (
        <>
          {/* Search */}
          <Card>
            <CardContent className="p-4">
              <Input
                placeholder="ค้นหาเลขสมาชิก หรือ ชื่อสมาชิก..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                leftIcon={
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                }
              />
            </CardContent>
          </Card>

          {/* Table */}
          <Card>
            <CardHeader>
              <CardTitle subtitle={`แสดง ${filteredShares.length} รายการ`}>
                บัญชีหุ้นสมาชิก
              </CardTitle>
            </CardHeader>
            <Table>
              <TableHead>
                <TableRow>
                  <TableHeader>เลขสมาชิก</TableHeader>
                  <TableHeader>ชื่อสมาชิก</TableHeader>
                  <TableHeader className="text-right">จำนวนหุ้น</TableHeader>
                  <TableHeader className="text-right">มูลค่ารวม</TableHeader>
                  <TableHeader className="text-right">ค่าหุ้น/เดือน</TableHeader>
                  <TableHeader>ชำระล่าสุด</TableHeader>
                  <TableHeader>สถานะ</TableHeader>
                  <TableHeader className="text-center">ทำรายการ</TableHeader>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredShares.map((share) => (
                  <TableRow key={share.memberId}>
                    <TableCell className="font-medium">{share.memberNo}</TableCell>
                    <TableCell>
                      <Link href={`/members/${share.memberId}`} className="text-blue-600 hover:underline">
                        {share.memberName}
                      </Link>
                    </TableCell>
                    <TableCell className="text-right">{share.totalShares.toLocaleString()} หุ้น</TableCell>
                    <TableCell className="text-right font-semibold">{formatCurrency(share.totalAmount)}</TableCell>
                    <TableCell className="text-right">
                      {share.monthlyContribution > 0 ? formatCurrency(share.monthlyContribution) : '-'}
                    </TableCell>
                    <TableCell>{formatDate(share.lastPaymentDate)}</TableCell>
                    <TableCell>
                      <Badge variant={share.status === 'active' ? 'success' : 'default'}>
                        {share.status === 'active' ? 'ปกติ' : 'ลาออก'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-center gap-2">
                        {share.status === 'active' && (
                          <Button variant="ghost" size="sm" onClick={() => setShowBuyModal(true)}>
                            <svg className="w-4 h-4 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
                            </svg>
                          </Button>
                        )}
                        <Link href={`/shares/${share.memberId}`}>
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
              totalItems={filteredShares.length}
              onPageChange={setCurrentPage}
            />
          </Card>
        </>
      ) : (
        /* Transactions Tab */
        <Card>
          <CardHeader>
            <CardTitle>ประวัติการทำรายการหุ้น</CardTitle>
          </CardHeader>
          <Table>
            <TableHead>
              <TableRow>
                <TableHeader>วันที่</TableHeader>
                <TableHeader>ประเภท</TableHeader>
                <TableHeader>สมาชิก</TableHeader>
                <TableHeader className="text-right">จำนวนหุ้น</TableHeader>
                <TableHeader className="text-right">จำนวนเงิน</TableHeader>
                <TableHeader>รายละเอียด</TableHeader>
              </TableRow>
            </TableHead>
            <TableBody>
              {mockTransactions.map((tx, idx) => (
                <TableRow key={idx}>
                  <TableCell>{formatDate(tx.date)}</TableCell>
                  <TableCell>
                    <Badge variant={tx.type === 'buy' ? 'info' : tx.type === 'monthly' ? 'success' : 'default'}>
                      {transactionTypeLabels[tx.type]}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Link href={`/members/${tx.memberId}`} className="text-blue-600 hover:underline">
                      {tx.memberName}
                    </Link>
                  </TableCell>
                  <TableCell className="text-right">{tx.shares} หุ้น</TableCell>
                  <TableCell className="text-right font-medium">{formatCurrency(tx.amount)}</TableCell>
                  <TableCell className="text-slate-500">{tx.description}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}

      {/* Buy Share Modal */}
      <Modal isOpen={showBuyModal} onClose={() => setShowBuyModal(false)} size="md">
        <ModalHeader onClose={() => setShowBuyModal(false)}>
          ซื้อหุ้นเพิ่ม
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="เลขทะเบียนสมาชิก" placeholder="ค้นหาสมาชิก..." required />
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-500">ชื่อสมาชิก</p>
              <p className="font-medium text-slate-900">-</p>
              <div className="grid grid-cols-2 gap-4 mt-2">
                <div>
                  <p className="text-sm text-slate-500">หุ้นปัจจุบัน</p>
                  <p className="font-semibold text-slate-900">0 หุ้น</p>
                </div>
                <div>
                  <p className="text-sm text-slate-500">มูลค่ารวม</p>
                  <p className="font-semibold text-slate-900">฿0.00</p>
                </div>
              </div>
            </div>
            <Input label="จำนวนหุ้นที่ต้องการซื้อ" type="number" placeholder="0" required hint="มูลค่าหุ้นละ 100 บาท" />
            <div className="p-4 bg-emerald-50 rounded-lg border border-emerald-200">
              <div className="flex justify-between">
                <span className="text-emerald-700">ยอดเงินที่ต้องชำระ</span>
                <span className="font-bold text-emerald-800 text-lg">฿0.00</span>
              </div>
            </div>
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowBuyModal(false)}>
            ยกเลิก
          </Button>
          <Button variant="secondary" onClick={() => setShowBuyModal(false)}>
            ยืนยันซื้อหุ้น
          </Button>
        </ModalFooter>
      </Modal>

      {/* Monthly Payment Modal */}
      <Modal isOpen={showMonthlyModal} onClose={() => setShowMonthlyModal(false)} size="lg">
        <ModalHeader onClose={() => setShowMonthlyModal(false)}>
          บันทึกค่าหุ้นรายเดือน
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Select
                label="เดือน"
                options={[
                  { value: '1', label: 'มกราคม' },
                  { value: '2', label: 'กุมภาพันธ์' },
                  { value: '3', label: 'มีนาคม' },
                  { value: '4', label: 'เมษายน' },
                  { value: '5', label: 'พฤษภาคม' },
                  { value: '6', label: 'มิถุนายน' },
                  { value: '7', label: 'กรกฎาคม' },
                  { value: '8', label: 'สิงหาคม' },
                  { value: '9', label: 'กันยายน' },
                  { value: '10', label: 'ตุลาคม' },
                  { value: '11', label: 'พฤศจิกายน' },
                  { value: '12', label: 'ธันวาคม' },
                ]}
                required
              />
              <Select
                label="ปี พ.ศ."
                options={[
                  { value: '2567', label: '2567' },
                  { value: '2568', label: '2568' },
                ]}
                required
              />
            </div>
            <Input label="เลขทะเบียนสมาชิก (เว้นว่างเพื่อตัดยอดทั้งหมด)" placeholder="ค้นหาสมาชิก..." />
            <div className="p-4 bg-amber-50 rounded-lg border border-amber-200">
              <p className="text-sm text-amber-800">
                <strong>หมายเหตุ:</strong> ระบบจะตัดยอดค่าหุ้นรายเดือนจากบัญชีเงินฝากของสมาชิกโดยอัตโนมัติ
                หากยอดเงินฝากไม่เพียงพอจะแจ้งเตือนให้สมาชิกมาชำระด้วยตนเอง
              </p>
            </div>
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowMonthlyModal(false)}>
            ยกเลิก
          </Button>
          <Button onClick={() => setShowMonthlyModal(false)}>
            ดำเนินการ
          </Button>
        </ModalFooter>
      </Modal>
    </div>
  );
}
