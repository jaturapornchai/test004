'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell, Badge, EmptyState, Pagination } from '@/components/ui/table';
import { Modal, ModalHeader, ModalContent, ModalFooter, ConfirmDialog } from '@/components/ui/modal';

// Mock data
const mockCooperatives = [
  {
    coopId: 'COOP001',
    name: 'สหกรณ์การเกษตรเมืองเชียงใหม่ จำกัด',
    registrationNo: 'ชม.1234/2540',
    type: 'agriculture',
    province: 'เชียงใหม่',
    memberCount: 2450,
    totalAssets: 125000000,
    status: 'active',
  },
  {
    coopId: 'COOP002',
    name: 'สหกรณ์การเกษตรดอยสะเก็ด จำกัด',
    registrationNo: 'ชม.5678/2545',
    type: 'agriculture',
    province: 'เชียงใหม่',
    memberCount: 1830,
    totalAssets: 89000000,
    status: 'active',
  },
  {
    coopId: 'COOP003',
    name: 'สหกรณ์การเกษตรสันกำแพง จำกัด',
    registrationNo: 'ชม.9012/2548',
    type: 'agriculture',
    province: 'เชียงใหม่',
    memberCount: 980,
    totalAssets: 45000000,
    status: 'active',
  },
];

const typeLabels: Record<string, string> = {
  agriculture: 'สหกรณ์การเกษตร',
  savings: 'สหกรณ์ออมทรัพย์',
  consumer: 'สหกรณ์ร้านค้า',
  service: 'สหกรณ์บริการ',
  credit_union: 'สหกรณ์เครดิตยูเนี่ยน',
};

export default function CooperativesPage() {
  const [search, setSearch] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [selectedCoop, setSelectedCoop] = useState<typeof mockCooperatives[0] | null>(null);
  const [currentPage, setCurrentPage] = useState(1);

  const filteredCoops = mockCooperatives.filter(
    (coop) =>
      coop.name.toLowerCase().includes(search.toLowerCase()) ||
      coop.registrationNo.toLowerCase().includes(search.toLowerCase())
  );

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('th-TH', {
      style: 'currency',
      currency: 'THB',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">จัดการสหกรณ์</h1>
          <p className="text-slate-600 mt-1">รายการสหกรณ์ในระบบทั้งหมด</p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          เพิ่มสหกรณ์ใหม่
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-500">สหกรณ์ทั้งหมด</p>
                <p className="text-2xl font-bold text-slate-900">{mockCooperatives.length}</p>
              </div>
              <div className="p-3 bg-blue-100 rounded-lg">
                <svg className="w-6 h-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5" />
                </svg>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-500">สมาชิกรวม</p>
                <p className="text-2xl font-bold text-slate-900">
                  {mockCooperatives.reduce((sum, c) => sum + c.memberCount, 0).toLocaleString()}
                </p>
              </div>
              <div className="p-3 bg-emerald-100 rounded-lg">
                <svg className="w-6 h-6 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-500">สินทรัพย์รวม</p>
                <p className="text-2xl font-bold text-slate-900">
                  {formatCurrency(mockCooperatives.reduce((sum, c) => sum + c.totalAssets, 0))}
                </p>
              </div>
              <div className="p-3 bg-amber-100 rounded-lg">
                <svg className="w-6 h-6 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-500">ใช้งานอยู่</p>
                <p className="text-2xl font-bold text-emerald-600">
                  {mockCooperatives.filter((c) => c.status === 'active').length}
                </p>
              </div>
              <div className="p-3 bg-emerald-100 rounded-lg">
                <svg className="w-6 h-6 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Search & Filter */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col lg:flex-row gap-4">
            <div className="flex-1">
              <Input
                placeholder="ค้นหาชื่อสหกรณ์ หรือ เลขทะเบียน..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                leftIcon={
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                }
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardHeader>
          <CardTitle>รายการสหกรณ์</CardTitle>
        </CardHeader>
        {filteredCoops.length > 0 ? (
          <>
            <Table>
              <TableHead>
                <TableRow>
                  <TableHeader>รหัส</TableHeader>
                  <TableHeader>ชื่อสหกรณ์</TableHeader>
                  <TableHeader>ประเภท</TableHeader>
                  <TableHeader>จังหวัด</TableHeader>
                  <TableHeader className="text-right">จำนวนสมาชิก</TableHeader>
                  <TableHeader className="text-right">สินทรัพย์</TableHeader>
                  <TableHeader>สถานะ</TableHeader>
                  <TableHeader className="text-center">จัดการ</TableHeader>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredCoops.map((coop) => (
                  <TableRow key={coop.coopId}>
                    <TableCell className="font-medium">{coop.coopId}</TableCell>
                    <TableCell>
                      <Link href={`/cooperatives/${coop.coopId}`} className="text-blue-600 hover:underline">
                        {coop.name}
                      </Link>
                      <p className="text-xs text-slate-500">{coop.registrationNo}</p>
                    </TableCell>
                    <TableCell>{typeLabels[coop.type]}</TableCell>
                    <TableCell>{coop.province}</TableCell>
                    <TableCell className="text-right font-medium">
                      {coop.memberCount.toLocaleString()} คน
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {formatCurrency(coop.totalAssets)}
                    </TableCell>
                    <TableCell>
                      <Badge variant={coop.status === 'active' ? 'success' : 'default'}>
                        {coop.status === 'active' ? 'ใช้งาน' : 'ปิดใช้งาน'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-center gap-2">
                        <Link href={`/cooperatives/${coop.coopId}`}>
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
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => {
                            setSelectedCoop(coop);
                            setShowDeleteDialog(true);
                          }}
                        >
                          <svg className="w-4 h-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
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
              totalItems={filteredCoops.length}
              onPageChange={setCurrentPage}
            />
          </>
        ) : (
          <EmptyState
            title="ไม่พบข้อมูลสหกรณ์"
            description="ยังไม่มีสหกรณ์ในระบบ หรือไม่พบข้อมูลที่ค้นหา"
            action={
              <Button onClick={() => setShowCreateModal(true)}>
                เพิ่มสหกรณ์ใหม่
              </Button>
            }
          />
        )}
      </Card>

      {/* Create Modal */}
      <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} size="lg">
        <ModalHeader onClose={() => setShowCreateModal(false)}>
          เพิ่มสหกรณ์ใหม่
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="ชื่อสหกรณ์" placeholder="เช่น สหกรณ์การเกษตรเมืองเชียงใหม่ จำกัด" required />
            <div className="grid grid-cols-2 gap-4">
              <Input label="เลขทะเบียนสหกรณ์" placeholder="เช่น ชม.1234/2540" required />
              <Input label="เลขประจำตัวผู้เสียภาษี" placeholder="13 หลัก" />
            </div>
            <Input label="ที่อยู่" placeholder="เลขที่ ถนน ตำบล" required />
            <div className="grid grid-cols-3 gap-4">
              <Input label="อำเภอ" placeholder="อำเภอ" required />
              <Input label="จังหวัด" placeholder="จังหวัด" required />
              <Input label="รหัสไปรษณีย์" placeholder="5 หลัก" required />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <Input label="โทรศัพท์" placeholder="0XX-XXX-XXXX" required />
              <Input label="อีเมล" type="email" placeholder="example@coop.or.th" />
            </div>
            <Input label="วันที่จดทะเบียน" type="date" required />
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
        title="ยืนยันการลบสหกรณ์"
        message={`คุณต้องการลบ "${selectedCoop?.name}" ออกจากระบบใช่หรือไม่? การดำเนินการนี้ไม่สามารถยกเลิกได้`}
        confirmText="ลบสหกรณ์"
        variant="danger"
      />
    </div>
  );
}
