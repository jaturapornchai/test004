'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input, Select } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, StatsCard } from '@/components/ui/card';
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell, Badge, EmptyState, Pagination } from '@/components/ui/table';
import { Modal, ModalHeader, ModalContent, ModalFooter, ConfirmDialog } from '@/components/ui/modal';
import { formatCurrency, formatDate, formatIdCard, formatPhone, memberStatusLabels, prefixes, provinces } from '@/lib/utils';

// Mock data
const mockMembers = [
  {
    memberId: 'M001',
    memberNo: '0001',
    prefix: 'นาย',
    firstName: 'สมชาย',
    lastName: 'ใจดี',
    idCardNo: '5101234567890',
    phone: '0812345678',
    address: { subDistrict: 'ช้างคลาน', district: 'เมือง', province: 'เชียงใหม่' },
    joinDate: '2020-03-15',
    status: 'active',
    shares: 500,
    shareAmount: 50000,
    savings: 125000,
    loanBalance: 150000,
  },
  {
    memberId: 'M002',
    memberNo: '0002',
    prefix: 'นาง',
    firstName: 'สมหญิง',
    lastName: 'รักษ์ธรรม',
    idCardNo: '5109876543210',
    phone: '0898765432',
    address: { subDistrict: 'สันทราย', district: 'สันทราย', province: 'เชียงใหม่' },
    joinDate: '2019-07-22',
    status: 'active',
    shares: 800,
    shareAmount: 80000,
    savings: 250000,
    loanBalance: 0,
  },
  {
    memberId: 'M003',
    memberNo: '0003',
    prefix: 'นาย',
    firstName: 'วิชัย',
    lastName: 'พัฒนา',
    idCardNo: '5105555555555',
    phone: '0865554444',
    address: { subDistrict: 'ดอยสะเก็ด', district: 'ดอยสะเก็ด', province: 'เชียงใหม่' },
    joinDate: '2018-01-10',
    status: 'active',
    shares: 1200,
    shareAmount: 120000,
    savings: 380000,
    loanBalance: 200000,
  },
  {
    memberId: 'M004',
    memberNo: '0004',
    prefix: 'นางสาว',
    firstName: 'จันทร์',
    lastName: 'แสงงาม',
    idCardNo: '5107777777777',
    phone: '0877778888',
    address: { subDistrict: 'หางดง', district: 'หางดง', province: 'เชียงใหม่' },
    joinDate: '2021-05-20',
    status: 'resigned',
    shares: 0,
    shareAmount: 0,
    savings: 0,
    loanBalance: 0,
  },
];

export default function MembersPage() {
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [selectedMember, setSelectedMember] = useState<typeof mockMembers[0] | null>(null);
  const [currentPage, setCurrentPage] = useState(1);

  const filteredMembers = mockMembers.filter((member) => {
    const matchSearch =
      member.firstName.toLowerCase().includes(search.toLowerCase()) ||
      member.lastName.toLowerCase().includes(search.toLowerCase()) ||
      member.memberNo.includes(search) ||
      member.idCardNo.includes(search);
    const matchStatus = !statusFilter || member.status === statusFilter;
    return matchSearch && matchStatus;
  });

  const activeMembers = mockMembers.filter((m) => m.status === 'active');
  const totalShares = activeMembers.reduce((sum, m) => sum + m.shareAmount, 0);
  const totalSavings = activeMembers.reduce((sum, m) => sum + m.savings, 0);
  const totalLoans = activeMembers.reduce((sum, m) => sum + m.loanBalance, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">ทะเบียนสมาชิก</h1>
          <p className="text-slate-600 mt-1">จัดการข้อมูลสมาชิกสหกรณ์</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
            </svg>
            นำเข้าข้อมูล
          </Button>
          <Button onClick={() => setShowCreateModal(true)}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
            </svg>
            เพิ่มสมาชิกใหม่
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="สมาชิกทั้งหมด"
          value={activeMembers.length.toLocaleString()}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          }
          color="blue"
          trend={{ value: 5.2, label: 'จากเดือนก่อน', isPositive: true }}
        />
        <StatsCard
          title="ทุนเรือนหุ้นรวม"
          value={formatCurrency(totalShares)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2z" />
            </svg>
          }
          color="green"
        />
        <StatsCard
          title="เงินฝากรวม"
          value={formatCurrency(totalSavings)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
          }
          color="purple"
        />
        <StatsCard
          title="หนี้สินคงค้าง"
          value={formatCurrency(totalLoans)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
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
                placeholder="ค้นหาชื่อ นามสกุล เลขสมาชิก หรือ เลขบัตรประชาชน..."
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
                  { value: '', label: 'สถานะทั้งหมด' },
                  { value: 'active', label: 'ปกติ' },
                  { value: 'resigned', label: 'ลาออก' },
                  { value: 'deceased', label: 'ถึงแก่กรรม' },
                  { value: 'expelled', label: 'ให้ออก' },
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
          <CardTitle subtitle={`แสดง ${filteredMembers.length} รายการ`}>
            รายชื่อสมาชิก
          </CardTitle>
        </CardHeader>
        {filteredMembers.length > 0 ? (
          <>
            <Table>
              <TableHead>
                <TableRow>
                  <TableHeader>เลขสมาชิก</TableHeader>
                  <TableHeader>ชื่อ-นามสกุล</TableHeader>
                  <TableHeader>เลขบัตรประชาชน</TableHeader>
                  <TableHeader>โทรศัพท์</TableHeader>
                  <TableHeader className="text-right">ทุนเรือนหุ้น</TableHeader>
                  <TableHeader className="text-right">เงินฝาก</TableHeader>
                  <TableHeader className="text-right">หนี้คงค้าง</TableHeader>
                  <TableHeader>สถานะ</TableHeader>
                  <TableHeader className="text-center">จัดการ</TableHeader>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredMembers.map((member) => (
                  <TableRow key={member.memberId}>
                    <TableCell className="font-medium">{member.memberNo}</TableCell>
                    <TableCell>
                      <Link href={`/members/${member.memberId}`} className="text-blue-600 hover:underline">
                        {member.prefix}{member.firstName} {member.lastName}
                      </Link>
                      <p className="text-xs text-slate-500">สมาชิกเมื่อ {formatDate(member.joinDate)}</p>
                    </TableCell>
                    <TableCell className="font-mono text-sm">{formatIdCard(member.idCardNo)}</TableCell>
                    <TableCell>{formatPhone(member.phone)}</TableCell>
                    <TableCell className="text-right font-medium">{formatCurrency(member.shareAmount)}</TableCell>
                    <TableCell className="text-right font-medium">{formatCurrency(member.savings)}</TableCell>
                    <TableCell className="text-right font-medium">
                      {member.loanBalance > 0 ? (
                        <span className="text-amber-600">{formatCurrency(member.loanBalance)}</span>
                      ) : (
                        <span className="text-slate-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={member.status === 'active' ? 'success' : member.status === 'resigned' ? 'default' : 'danger'}
                      >
                        {memberStatusLabels[member.status]}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-center gap-2">
                        <Link href={`/members/${member.memberId}`}>
                          <Button variant="ghost" size="sm">
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                              <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                              <path strokeLinecap="round" strokeLinejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                            </svg>
                          </Button>
                        </Link>
                        <Button variant="ghost" size="sm">
                          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                          </svg>
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
            <Pagination
              currentPage={currentPage}
              totalPages={1}
              totalItems={filteredMembers.length}
              onPageChange={setCurrentPage}
            />
          </>
        ) : (
          <EmptyState
            title="ไม่พบข้อมูลสมาชิก"
            description="ยังไม่มีสมาชิกในระบบ หรือไม่พบข้อมูลที่ค้นหา"
            action={
              <Button onClick={() => setShowCreateModal(true)}>
                เพิ่มสมาชิกใหม่
              </Button>
            }
          />
        )}
      </Card>

      {/* Create Modal */}
      <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} size="xl">
        <ModalHeader onClose={() => setShowCreateModal(false)}>
          เพิ่มสมาชิกใหม่
        </ModalHeader>
        <ModalContent>
          <div className="space-y-6">
            <div>
              <h3 className="text-sm font-semibold text-slate-900 mb-3">ข้อมูลส่วนบุคคล</h3>
              <div className="grid grid-cols-12 gap-4">
                <div className="col-span-2">
                  <Select
                    label="คำนำหน้า"
                    options={prefixes}
                    required
                  />
                </div>
                <div className="col-span-5">
                  <Input label="ชื่อ" placeholder="ชื่อจริง" required />
                </div>
                <div className="col-span-5">
                  <Input label="นามสกุล" placeholder="นามสกุล" required />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4 mt-4">
                <Input label="เลขบัตรประชาชน" placeholder="เลขบัตรประชาชน 13 หลัก" required />
                <Input label="วันเดือนปีเกิด" type="date" required />
              </div>
              <div className="grid grid-cols-2 gap-4 mt-4">
                <Input label="โทรศัพท์" placeholder="0XX-XXX-XXXX" required />
                <Input label="อีเมล" type="email" placeholder="example@email.com" />
              </div>
            </div>

            <div>
              <h3 className="text-sm font-semibold text-slate-900 mb-3">ที่อยู่</h3>
              <Input label="ที่อยู่" placeholder="บ้านเลขที่ หมู่ ถนน" required />
              <div className="grid grid-cols-4 gap-4 mt-4">
                <Input label="ตำบล" placeholder="ตำบล" required />
                <Input label="อำเภอ" placeholder="อำเภอ" required />
                <Select
                  label="จังหวัด"
                  options={provinces.map((p) => ({ value: p, label: p }))}
                  placeholder="เลือกจังหวัด"
                  required
                />
                <Input label="รหัสไปรษณีย์" placeholder="5 หลัก" required />
              </div>
            </div>

            <div>
              <h3 className="text-sm font-semibold text-slate-900 mb-3">ข้อมูลการเกษตร</h3>
              <div className="grid grid-cols-3 gap-4">
                <Input label="อาชีพ" placeholder="อาชีพหลัก" required />
                <Input label="พื้นที่เพาะปลูก (ไร่)" type="number" placeholder="0" />
                <Select
                  label="สิทธิ์ในที่ดิน"
                  options={[
                    { value: 'owned', label: 'เจ้าของ' },
                    { value: 'rented', label: 'เช่า' },
                    { value: 'mixed', label: 'เจ้าของและเช่า' },
                  ]}
                />
              </div>
            </div>

            <div>
              <h3 className="text-sm font-semibold text-slate-900 mb-3">ทุนเรือนหุ้นแรกเข้า</h3>
              <div className="grid grid-cols-2 gap-4">
                <Input label="จำนวนหุ้น" type="number" placeholder="จำนวนหุ้นที่ต้องการซื้อ" required />
                <Input label="ค่าหุ้นรายเดือน" type="number" placeholder="ค่าหุ้นรายเดือน (บาท)" required />
              </div>
            </div>
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowCreateModal(false)}>
            ยกเลิก
          </Button>
          <Button onClick={() => setShowCreateModal(false)}>
            บันทึก
          </Button>
        </ModalFooter>
      </Modal>

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
        onConfirm={() => setShowDeleteDialog(false)}
        title="ยืนยันการลบสมาชิก"
        message={`คุณต้องการลบ "${selectedMember?.prefix}${selectedMember?.firstName} ${selectedMember?.lastName}" ออกจากระบบใช่หรือไม่?`}
        confirmText="ลบสมาชิก"
        variant="danger"
      />
    </div>
  );
}
