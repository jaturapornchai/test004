'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input, Select } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, StatsCard } from '@/components/ui/card';
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell, Badge, EmptyState, Pagination } from '@/components/ui/table';
import { Modal, ModalHeader, ModalContent, ModalFooter } from '@/components/ui/modal';
import { formatCurrency, formatDate, loanStatusLabels, loanTypeLabels } from '@/lib/utils';

// Mock data
const mockLoans = [
  {
    contractNo: 'L2567-0001',
    memberId: 'M001',
    memberName: 'นายสมชาย ใจดี',
    productName: 'เงินกู้สามัญ',
    loanType: 'ordinary',
    principalAmount: 200000,
    interestRate: 6.5,
    termMonths: 36,
    monthlyPayment: 6150,
    principalBalance: 150000,
    disbursementDate: '2024-01-15',
    nextPaymentDate: '2024-12-15',
    daysOverdue: 0,
    status: 'active',
  },
  {
    contractNo: 'L2567-0002',
    memberId: 'M003',
    memberName: 'นายวิชัย พัฒนา',
    productName: 'เงินกู้พิเศษ',
    loanType: 'special',
    principalAmount: 500000,
    interestRate: 5.5,
    termMonths: 120,
    monthlyPayment: 5500,
    principalBalance: 450000,
    disbursementDate: '2023-06-20',
    nextPaymentDate: '2024-12-20',
    daysOverdue: 0,
    status: 'active',
  },
  {
    contractNo: 'L2567-0003',
    memberId: 'M002',
    memberName: 'นางสมหญิง รักษ์ธรรม',
    productName: 'เงินกู้ฉุกเฉิน',
    loanType: 'emergency',
    principalAmount: 30000,
    interestRate: 6.0,
    termMonths: 12,
    monthlyPayment: 2650,
    principalBalance: 0,
    disbursementDate: '2024-03-01',
    nextPaymentDate: '-',
    daysOverdue: 0,
    status: 'closed',
  },
  {
    contractNo: 'L2566-0045',
    memberId: 'M005',
    memberName: 'นายประสิทธิ์ การดี',
    productName: 'เงินกู้สามัญ',
    loanType: 'ordinary',
    principalAmount: 150000,
    interestRate: 6.5,
    termMonths: 24,
    monthlyPayment: 6700,
    principalBalance: 45000,
    disbursementDate: '2023-01-10',
    nextPaymentDate: '2024-11-10',
    daysOverdue: 35,
    status: 'default',
  },
];

const mockApplications = [
  {
    applicationId: 'APP2567-0012',
    memberId: 'M006',
    memberName: 'นางมาลี สุขใจ',
    productName: 'เงินกู้สามัญ',
    requestedAmount: 100000,
    purpose: 'ซื้อปุ๋ยและเมล็ดพันธุ์',
    applicationDate: '2024-11-18',
    status: 'pending',
  },
  {
    applicationId: 'APP2567-0011',
    memberId: 'M007',
    memberName: 'นายสมศักดิ์ ดีมาก',
    productName: 'เงินกู้พิเศษ',
    requestedAmount: 300000,
    purpose: 'ซื้อที่ดินทำกิน',
    applicationDate: '2024-11-15',
    status: 'pending',
  },
];

export default function LoansPage() {
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [typeFilter, setTypeFilter] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [activeTab, setActiveTab] = useState<'contracts' | 'applications'>('contracts');

  const filteredLoans = mockLoans.filter((loan) => {
    const matchSearch =
      loan.contractNo.toLowerCase().includes(search.toLowerCase()) ||
      loan.memberName.toLowerCase().includes(search.toLowerCase());
    const matchStatus = !statusFilter || loan.status === statusFilter;
    const matchType = !typeFilter || loan.loanType === typeFilter;
    return matchSearch && matchStatus && matchType;
  });

  const activeLoans = mockLoans.filter((l) => l.status === 'active');
  const totalPrincipal = activeLoans.reduce((sum, l) => sum + l.principalBalance, 0);
  const overdueLoans = mockLoans.filter((l) => l.status === 'default');
  const overdueAmount = overdueLoans.reduce((sum, l) => sum + l.principalBalance, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">ระบบสินเชื่อ</h1>
          <p className="text-slate-600 mt-1">จัดการคำขอกู้ สัญญาเงินกู้ และติดตามการชำระหนี้</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            รายงาน
          </Button>
          <Button onClick={() => setShowCreateModal(true)}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            คำขอสินเชื่อใหม่
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="สัญญาที่ใช้งาน"
          value={activeLoans.length.toLocaleString()}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          }
          color="blue"
        />
        <StatsCard
          title="เงินกู้คงค้างรวม"
          value={formatCurrency(totalPrincipal)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
          color="green"
        />
        <StatsCard
          title="คำขอรออนุมัติ"
          value={mockApplications.length.toLocaleString()}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
          color="yellow"
        />
        <StatsCard
          title="หนี้ค้างชำระ"
          value={formatCurrency(overdueAmount)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          }
          color="red"
        />
      </div>

      {/* Tabs */}
      <div className="border-b border-slate-200">
        <div className="flex gap-4">
          <button
            onClick={() => setActiveTab('contracts')}
            className={`pb-3 px-1 font-medium text-sm border-b-2 transition-colors ${
              activeTab === 'contracts'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            สัญญาเงินกู้
          </button>
          <button
            onClick={() => setActiveTab('applications')}
            className={`pb-3 px-1 font-medium text-sm border-b-2 transition-colors ${
              activeTab === 'applications'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            คำขอสินเชื่อ
            {mockApplications.length > 0 && (
              <span className="ml-2 px-2 py-0.5 bg-amber-100 text-amber-700 text-xs rounded-full">
                {mockApplications.length}
              </span>
            )}
          </button>
        </div>
      </div>

      {activeTab === 'contracts' ? (
        <>
          {/* Search & Filter */}
          <Card>
            <CardContent className="p-4">
              <div className="flex flex-col lg:flex-row gap-4">
                <div className="flex-1">
                  <Input
                    placeholder="ค้นหาเลขที่สัญญา หรือ ชื่อสมาชิก..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    leftIcon={
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                      </svg>
                    }
                  />
                </div>
                <div className="w-full lg:w-40">
                  <Select
                    options={[
                      { value: '', label: 'ประเภททั้งหมด' },
                      { value: 'emergency', label: 'ฉุกเฉิน' },
                      { value: 'ordinary', label: 'สามัญ' },
                      { value: 'special', label: 'พิเศษ' },
                    ]}
                    value={typeFilter}
                    onChange={setTypeFilter}
                  />
                </div>
                <div className="w-full lg:w-40">
                  <Select
                    options={[
                      { value: '', label: 'สถานะทั้งหมด' },
                      { value: 'active', label: 'ปกติ' },
                      { value: 'closed', label: 'ปิดบัญชี' },
                      { value: 'default', label: 'ค้างชำระ' },
                    ]}
                    value={statusFilter}
                    onChange={setStatusFilter}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Table */}
          <Card>
            <CardHeader>
              <CardTitle subtitle={`แสดง ${filteredLoans.length} รายการ`}>
                รายการสัญญาเงินกู้
              </CardTitle>
            </CardHeader>
            {filteredLoans.length > 0 ? (
              <>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableHeader>เลขที่สัญญา</TableHeader>
                      <TableHeader>สมาชิก</TableHeader>
                      <TableHeader>ประเภท</TableHeader>
                      <TableHeader className="text-right">วงเงิน</TableHeader>
                      <TableHeader className="text-right">ยอดคงค้าง</TableHeader>
                      <TableHeader className="text-right">ดอกเบี้ย</TableHeader>
                      <TableHeader>ชำระงวดถัดไป</TableHeader>
                      <TableHeader>สถานะ</TableHeader>
                      <TableHeader className="text-center">จัดการ</TableHeader>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredLoans.map((loan) => (
                      <TableRow key={loan.contractNo}>
                        <TableCell className="font-medium">{loan.contractNo}</TableCell>
                        <TableCell>
                          <Link href={`/members/${loan.memberId}`} className="text-blue-600 hover:underline">
                            {loan.memberName}
                          </Link>
                        </TableCell>
                        <TableCell>
                          <Badge variant="info">{loanTypeLabels[loan.loanType]}</Badge>
                        </TableCell>
                        <TableCell className="text-right">{formatCurrency(loan.principalAmount)}</TableCell>
                        <TableCell className="text-right font-medium">
                          {loan.principalBalance > 0 ? (
                            <span className={loan.status === 'default' ? 'text-red-600' : ''}>
                              {formatCurrency(loan.principalBalance)}
                            </span>
                          ) : (
                            <span className="text-slate-400">-</span>
                          )}
                        </TableCell>
                        <TableCell className="text-right">{loan.interestRate}%</TableCell>
                        <TableCell>
                          {loan.nextPaymentDate !== '-' ? (
                            <div>
                              <span>{formatDate(loan.nextPaymentDate)}</span>
                              {loan.daysOverdue > 0 && (
                                <p className="text-xs text-red-600">เกินกำหนด {loan.daysOverdue} วัน</p>
                              )}
                            </div>
                          ) : (
                            <span className="text-slate-400">-</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={
                              loan.status === 'active' ? 'success' :
                              loan.status === 'closed' ? 'default' :
                              loan.status === 'default' ? 'danger' : 'warning'
                            }
                          >
                            {loanStatusLabels[loan.status]}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center justify-center gap-2">
                            <Link href={`/loans/${loan.contractNo}`}>
                              <Button variant="ghost" size="sm">
                                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                  <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                  <path strokeLinecap="round" strokeLinejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                </svg>
                              </Button>
                            </Link>
                            {loan.status === 'active' && (
                              <Button variant="ghost" size="sm">
                                <svg className="w-4 h-4 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                  <path strokeLinecap="round" strokeLinejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2z" />
                                </svg>
                              </Button>
                            )}
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
                <Pagination
                  currentPage={currentPage}
                  totalPages={1}
                  totalItems={filteredLoans.length}
                  onPageChange={setCurrentPage}
                />
              </>
            ) : (
              <EmptyState
                title="ไม่พบข้อมูลสัญญา"
                description="ยังไม่มีสัญญาเงินกู้ในระบบ หรือไม่พบข้อมูลที่ค้นหา"
              />
            )}
          </Card>
        </>
      ) : (
        /* Applications Tab */
        <Card>
          <CardHeader
            action={
              <Button onClick={() => setShowCreateModal(true)} size="sm">
                คำขอใหม่
              </Button>
            }
          >
            <CardTitle subtitle={`${mockApplications.length} รายการรออนุมัติ`}>
              คำขอสินเชื่อ
            </CardTitle>
          </CardHeader>
          {mockApplications.length > 0 ? (
            <Table>
              <TableHead>
                <TableRow>
                  <TableHeader>เลขที่คำขอ</TableHeader>
                  <TableHeader>สมาชิก</TableHeader>
                  <TableHeader>ประเภท</TableHeader>
                  <TableHeader className="text-right">จำนวนเงิน</TableHeader>
                  <TableHeader>วัตถุประสงค์</TableHeader>
                  <TableHeader>วันที่ยื่น</TableHeader>
                  <TableHeader>สถานะ</TableHeader>
                  <TableHeader className="text-center">จัดการ</TableHeader>
                </TableRow>
              </TableHead>
              <TableBody>
                {mockApplications.map((app) => (
                  <TableRow key={app.applicationId}>
                    <TableCell className="font-medium">{app.applicationId}</TableCell>
                    <TableCell>{app.memberName}</TableCell>
                    <TableCell>{app.productName}</TableCell>
                    <TableCell className="text-right font-medium">{formatCurrency(app.requestedAmount)}</TableCell>
                    <TableCell className="max-w-[200px] truncate">{app.purpose}</TableCell>
                    <TableCell>{formatDate(app.applicationDate)}</TableCell>
                    <TableCell>
                      <Badge variant="warning">รอพิจารณา</Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-center gap-2">
                        <Button variant="ghost" size="sm">
                          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path strokeLinecap="round" strokeLinejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                          </svg>
                        </Button>
                        <Button size="sm" variant="secondary">
                          อนุมัติ
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <EmptyState
              title="ไม่มีคำขอสินเชื่อ"
              description="ไม่มีคำขอสินเชื่อที่รอพิจารณา"
            />
          )}
        </Card>
      )}

      {/* Create Application Modal */}
      <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} size="lg">
        <ModalHeader onClose={() => setShowCreateModal(false)}>
          ยื่นคำขอสินเชื่อใหม่
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="เลขทะเบียนสมาชิก" placeholder="ค้นหาสมาชิก..." required />
            <Select
              label="ประเภทสินเชื่อ"
              options={[
                { value: 'emergency', label: 'เงินกู้ฉุกเฉิน (สูงสุด 50,000 บาท)' },
                { value: 'ordinary', label: 'เงินกู้สามัญ (สูงสุด 500,000 บาท)' },
                { value: 'special', label: 'เงินกู้พิเศษ (ต้องมีหลักประกัน)' },
              ]}
              placeholder="เลือกประเภทสินเชื่อ"
              required
            />
            <div className="grid grid-cols-2 gap-4">
              <Input label="จำนวนเงินที่ขอกู้" type="number" placeholder="0.00" required />
              <Input label="ระยะเวลา (เดือน)" type="number" placeholder="12" required />
            </div>
            <Input label="วัตถุประสงค์การกู้" placeholder="ระบุวัตถุประสงค์..." required />
            <Select
              label="ผู้ค้ำประกัน 1"
              options={[]}
              placeholder="เลือกผู้ค้ำประกัน"
            />
            <Select
              label="ผู้ค้ำประกัน 2"
              options={[]}
              placeholder="เลือกผู้ค้ำประกัน"
            />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowCreateModal(false)}>
            ยกเลิก
          </Button>
          <Button onClick={() => setShowCreateModal(false)}>
            ยื่นคำขอ
          </Button>
        </ModalFooter>
      </Modal>
    </div>
  );
}
