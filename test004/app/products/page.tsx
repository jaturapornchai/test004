'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input, Select } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, StatsCard } from '@/components/ui/card';
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell, Badge, Pagination } from '@/components/ui/table';
import { Modal, ModalHeader, ModalContent, ModalFooter } from '@/components/ui/modal';
import { formatCurrency } from '@/lib/utils';

// Mock data - สินค้าคงคลัง
const mockProducts = [
  {
    productId: 'P001',
    name: 'ปุ๋ยเคมี สูตร 15-15-15',
    category: 'fertilizer',
    unit: 'กระสอบ',
    costPrice: 800,
    sellingPrice: 900,
    memberPrice: 850,
    currentStock: 250,
    minStock: 50,
    status: 'active',
  },
  {
    productId: 'P002',
    name: 'ปุ๋ยเคมี สูตร 46-0-0 (ยูเรีย)',
    category: 'fertilizer',
    unit: 'กระสอบ',
    costPrice: 750,
    sellingPrice: 850,
    memberPrice: 800,
    currentStock: 180,
    minStock: 50,
    status: 'active',
  },
  {
    productId: 'P003',
    name: 'ยาฆ่าหญ้า พาราควอต',
    category: 'pesticide',
    unit: 'ขวด',
    costPrice: 120,
    sellingPrice: 150,
    memberPrice: 140,
    currentStock: 85,
    minStock: 30,
    status: 'active',
  },
  {
    productId: 'P004',
    name: 'เมล็ดพันธุ์ข้าว กข43',
    category: 'seed',
    unit: 'กก.',
    costPrice: 25,
    sellingPrice: 35,
    memberPrice: 30,
    currentStock: 15,
    minStock: 50,
    status: 'low_stock',
  },
  {
    productId: 'P005',
    name: 'ปุ๋ยอินทรีย์ชีวภาพ',
    category: 'fertilizer',
    unit: 'กระสอบ',
    costPrice: 200,
    sellingPrice: 280,
    memberPrice: 250,
    currentStock: 320,
    minStock: 100,
    status: 'active',
  },
];

const categoryLabels: Record<string, string> = {
  fertilizer: 'ปุ๋ย',
  pesticide: 'ยาปราบศัตรูพืช',
  seed: 'เมล็ดพันธุ์',
  equipment: 'อุปกรณ์',
  other: 'อื่นๆ',
};

export default function ProductsPage() {
  const [search, setSearch] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showSaleModal, setShowSaleModal] = useState(false);
  const [showStockModal, setShowStockModal] = useState(false);

  const filteredProducts = mockProducts.filter((product) => {
    const matchSearch = product.name.toLowerCase().includes(search.toLowerCase()) ||
      product.productId.toLowerCase().includes(search.toLowerCase());
    const matchCategory = !categoryFilter || product.category === categoryFilter;
    return matchSearch && matchCategory;
  });

  const totalProducts = mockProducts.length;
  const lowStockProducts = mockProducts.filter(p => p.currentStock <= p.minStock).length;
  const totalValue = mockProducts.reduce((sum, p) => sum + (p.costPrice * p.currentStock), 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">ธุรกิจจัดหาสินค้า</h1>
          <p className="text-slate-600 mt-1">จัดการคลังสินค้า ปุ๋ย ยา เมล็ดพันธุ์ และขายสินค้า</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" onClick={() => setShowStockModal(true)}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
            รับสินค้าเข้า
          </Button>
          <Button variant="secondary" onClick={() => setShowSaleModal(true)}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            ขายสินค้า
          </Button>
          <Button onClick={() => setShowCreateModal(true)}>
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            เพิ่มสินค้า
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="สินค้าทั้งหมด"
          value={`${totalProducts} รายการ`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          }
          color="blue"
        />
        <StatsCard
          title="มูลค่าคงคลัง"
          value={formatCurrency(totalValue)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
          color="green"
        />
        <StatsCard
          title="สินค้าใกล้หมด"
          value={`${lowStockProducts} รายการ`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          }
          color="yellow"
        />
        <StatsCard
          title="ยอดขายเดือนนี้"
          value={formatCurrency(125000)}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
          }
          color="purple"
        />
      </div>

      {/* Search & Filter */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col lg:flex-row gap-4">
            <div className="flex-1">
              <Input
                placeholder="ค้นหาชื่อสินค้า หรือ รหัสสินค้า..."
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
                  { value: '', label: 'หมวดหมู่ทั้งหมด' },
                  { value: 'fertilizer', label: 'ปุ๋ย' },
                  { value: 'pesticide', label: 'ยาปราบศัตรูพืช' },
                  { value: 'seed', label: 'เมล็ดพันธุ์' },
                  { value: 'equipment', label: 'อุปกรณ์' },
                  { value: 'other', label: 'อื่นๆ' },
                ]}
                value={categoryFilter}
                onChange={setCategoryFilter}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardHeader>
          <CardTitle subtitle={`แสดง ${filteredProducts.length} รายการ`}>
            รายการสินค้า
          </CardTitle>
        </CardHeader>
        <Table>
          <TableHead>
            <TableRow>
              <TableHeader>รหัสสินค้า</TableHeader>
              <TableHeader>ชื่อสินค้า</TableHeader>
              <TableHeader>หมวดหมู่</TableHeader>
              <TableHeader className="text-right">ราคาทุน</TableHeader>
              <TableHeader className="text-right">ราคาขาย</TableHeader>
              <TableHeader className="text-right">ราคาสมาชิก</TableHeader>
              <TableHeader className="text-right">คงเหลือ</TableHeader>
              <TableHeader>สถานะ</TableHeader>
              <TableHeader className="text-center">จัดการ</TableHeader>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredProducts.map((product) => (
              <TableRow key={product.productId}>
                <TableCell className="font-medium">{product.productId}</TableCell>
                <TableCell>{product.name}</TableCell>
                <TableCell>
                  <Badge variant="info">{categoryLabels[product.category]}</Badge>
                </TableCell>
                <TableCell className="text-right">{formatCurrency(product.costPrice)}</TableCell>
                <TableCell className="text-right">{formatCurrency(product.sellingPrice)}</TableCell>
                <TableCell className="text-right">{formatCurrency(product.memberPrice)}</TableCell>
                <TableCell className="text-right">
                  <span className={product.currentStock <= product.minStock ? 'text-red-600 font-semibold' : ''}>
                    {product.currentStock} {product.unit}
                  </span>
                </TableCell>
                <TableCell>
                  <Badge variant={product.currentStock <= product.minStock ? 'danger' : 'success'}>
                    {product.currentStock <= product.minStock ? 'ใกล้หมด' : 'ปกติ'}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div className="flex items-center justify-center gap-2">
                    <Button variant="ghost" size="sm">
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => setShowSaleModal(true)}>
                      <svg className="w-4 h-4 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
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
          totalItems={filteredProducts.length}
          onPageChange={setCurrentPage}
        />
      </Card>

      {/* Create Product Modal */}
      <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} size="lg">
        <ModalHeader onClose={() => setShowCreateModal(false)}>
          เพิ่มสินค้าใหม่
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <Input label="รหัสสินค้า" placeholder="P001" required />
              <div className="col-span-2">
                <Input label="ชื่อสินค้า" placeholder="ชื่อสินค้า" required />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <Select
                label="หมวดหมู่"
                options={[
                  { value: 'fertilizer', label: 'ปุ๋ย' },
                  { value: 'pesticide', label: 'ยาปราบศัตรูพืช' },
                  { value: 'seed', label: 'เมล็ดพันธุ์' },
                  { value: 'equipment', label: 'อุปกรณ์' },
                  { value: 'other', label: 'อื่นๆ' },
                ]}
                placeholder="เลือกหมวดหมู่"
                required
              />
              <Input label="หน่วยนับ" placeholder="เช่น กระสอบ, ขวด, กก." required />
            </div>
            <div className="grid grid-cols-3 gap-4">
              <Input label="ราคาทุน" type="number" placeholder="0.00" required />
              <Input label="ราคาขายทั่วไป" type="number" placeholder="0.00" required />
              <Input label="ราคาสมาชิก" type="number" placeholder="0.00" required />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <Input label="จำนวนเริ่มต้น" type="number" placeholder="0" />
              <Input label="จุดสั่งซื้อขั้นต่ำ" type="number" placeholder="50" hint="แจ้งเตือนเมื่อสต็อกต่ำกว่านี้" />
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

      {/* Sale Modal */}
      <Modal isOpen={showSaleModal} onClose={() => setShowSaleModal(false)} size="xl">
        <ModalHeader onClose={() => setShowSaleModal(false)}>
          ขายสินค้า
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="เลขทะเบียนสมาชิก (ถ้ามี)" placeholder="ค้นหาสมาชิก..." />
              <Input label="ชื่อลูกค้า (กรณีไม่ใช่สมาชิก)" placeholder="ชื่อ-นามสกุล" />
            </div>

            <div className="border rounded-lg p-4">
              <h3 className="font-medium text-slate-900 mb-3">รายการสินค้า</h3>
              <div className="grid grid-cols-12 gap-2 mb-2">
                <div className="col-span-5">
                  <Select
                    options={mockProducts.map(p => ({ value: p.productId, label: p.name }))}
                    placeholder="เลือกสินค้า"
                  />
                </div>
                <div className="col-span-2">
                  <Input type="number" placeholder="จำนวน" />
                </div>
                <div className="col-span-2">
                  <Input type="number" placeholder="ราคา" />
                </div>
                <div className="col-span-2">
                  <Input type="number" placeholder="รวม" disabled />
                </div>
                <div className="col-span-1 flex items-center">
                  <Button variant="ghost" size="sm">
                    <svg className="w-4 h-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </Button>
                </div>
              </div>
              <Button variant="ghost" size="sm">
                <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
                </svg>
                เพิ่มรายการ
              </Button>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <Select
                label="วิธีชำระเงิน"
                options={[
                  { value: 'cash', label: 'เงินสด' },
                  { value: 'credit', label: 'ขายเชื่อ (ตัดจากเงินค่าผลผลิต)' },
                  { value: 'transfer', label: 'โอนเงิน' },
                ]}
                required
              />
              <Input label="ส่วนลด" type="number" placeholder="0.00" />
            </div>

            <div className="p-4 bg-emerald-50 rounded-lg border border-emerald-200">
              <div className="flex justify-between text-lg">
                <span className="text-emerald-700 font-medium">ยอดรวมทั้งสิ้น</span>
                <span className="font-bold text-emerald-800">฿0.00</span>
              </div>
            </div>
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowSaleModal(false)}>
            ยกเลิก
          </Button>
          <Button variant="secondary" onClick={() => setShowSaleModal(false)}>
            ยืนยันการขาย
          </Button>
        </ModalFooter>
      </Modal>

      {/* Stock In Modal */}
      <Modal isOpen={showStockModal} onClose={() => setShowStockModal(false)} size="lg">
        <ModalHeader onClose={() => setShowStockModal(false)}>
          รับสินค้าเข้าคลัง
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Select
              label="สินค้า"
              options={mockProducts.map(p => ({ value: p.productId, label: `${p.productId} - ${p.name}` }))}
              placeholder="เลือกสินค้า"
              required
            />
            <div className="grid grid-cols-2 gap-4">
              <Input label="จำนวนที่รับเข้า" type="number" placeholder="0" required />
              <Input label="ราคาทุน/หน่วย" type="number" placeholder="0.00" required />
            </div>
            <Input label="เลขที่ใบส่งสินค้า" placeholder="เลขที่เอกสารอ้างอิง" />
            <Input label="หมายเหตุ" placeholder="ระบุหมายเหตุ (ถ้ามี)" />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="ghost" onClick={() => setShowStockModal(false)}>
            ยกเลิก
          </Button>
          <Button onClick={() => setShowStockModal(false)}>
            ยืนยันรับสินค้า
          </Button>
        </ModalFooter>
      </Modal>
    </div>
  );
}
