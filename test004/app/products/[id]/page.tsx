"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent, StatsCard } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input, Textarea } from "@/components/ui/input";
import { Modal, ModalHeader, ModalContent, ModalFooter } from "@/components/ui/modal";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell, Badge } from "@/components/ui/table";
import { formatCurrency, formatNumber, formatDate } from "@/lib/utils";

// Mock product data
const getProductData = (id: string) => ({
  _id: id,
  productId: id,
  coopId: "COOP001",
  name: "ปุ๋ยยูเรีย 46-0-0",
  description: "ปุ๋ยยูเรียสูตร 46-0-0 บรรจุกระสอบ 50 กก. เหมาะสำหรับข้าว พืชผัก และไม้ผล",
  category: "fertilizer" as const,
  sku: "FER-UREA-001",
  barcode: "8850000000001",
  unit: "กระสอบ",
  unitSize: "50 กก.",
  costPrice: 850,
  sellingPrice: 950,
  memberPrice: 900,
  stock: 45,
  minStock: 20,
  maxStock: 200,
  supplier: "บริษัท ปุ๋ยไทย จำกัด",
  status: "active" as const,
  createdAt: "2023-01-15T10:30:00Z",
  updatedAt: "2024-01-15T14:20:00Z",
});

// Mock stock movements
const getStockMovements = () => [
  { id: "STK001", date: "2024-01-15", type: "sale", description: "ขายให้สมาชิก M001", quantity: -5, balance: 45, reference: "INV-2024-0115" },
  { id: "STK002", date: "2024-01-10", type: "stock_in", description: "รับสินค้าเข้าคลัง", quantity: 30, balance: 50, reference: "PO-2024-0110" },
  { id: "STK003", date: "2024-01-05", type: "sale", description: "ขายให้สมาชิก M003", quantity: -10, balance: 20, reference: "INV-2024-0105" },
  { id: "STK004", date: "2023-12-28", type: "sale", description: "ขายให้บุคคลภายนอก", quantity: -8, balance: 30, reference: "INV-2023-1228" },
  { id: "STK005", date: "2023-12-20", type: "stock_in", description: "รับสินค้าเข้าคลัง", quantity: 25, balance: 38, reference: "PO-2023-1220" },
  { id: "STK006", date: "2023-12-15", type: "adjustment", description: "ปรับปรุงสต็อก (ตรวจนับ)", quantity: -2, balance: 13, reference: "ADJ-2023-1215" },
];

// Mock sales summary
const getSalesSummary = () => ({
  totalSold: 156,
  totalRevenue: 148200,
  totalProfit: 15600,
  avgMonthlySales: 13,
  lastSaleDate: "2024-01-15",
  topBuyers: [
    { memberId: "M001", name: "นายสมชาย ใจดี", quantity: 35, amount: 31500 },
    { memberId: "M003", name: "นายวิชัย พัฒนา", quantity: 28, amount: 25200 },
    { memberId: "M005", name: "นางมาลี สุขสันต์", quantity: 22, amount: 19800 },
  ],
});

const categoryLabels: Record<string, string> = {
  fertilizer: "ปุ๋ย",
  pesticide: "ยาปราบศัตรูพืช",
  seed: "เมล็ดพันธุ์",
  equipment: "อุปกรณ์การเกษตร",
  feed: "อาหารสัตว์",
  other: "อื่นๆ",
};

type TabType = "overview" | "stock" | "sales";

export default function ProductDetailPage() {
  const params = useParams();
  const productId = params.id as string;

  const [activeTab, setActiveTab] = useState<TabType>("overview");
  const [showStockInModal, setShowStockInModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);

  const product = getProductData(productId);
  const stockMovements = getStockMovements();
  const salesSummary = getSalesSummary();

  const profit = product.sellingPrice - product.costPrice;
  const profitPercent = (profit / product.costPrice) * 100;
  const isLowStock = product.stock <= product.minStock;

  const getMovementIcon = (type: string) => {
    switch (type) {
      case "stock_in":
        return (
          <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
          </div>
        );
      case "sale":
        return (
          <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
          </div>
        );
      case "adjustment":
        return (
          <div className="w-8 h-8 rounded-full bg-amber-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
        );
      default:
        return null;
    }
  };

  const renderOverviewTab = () => (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Product Info */}
      <Card className="lg:col-span-2">
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>ข้อมูลสินค้า</CardTitle>
          <Button variant="outline" size="sm" onClick={() => setShowEditModal(true)}>
            <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
            แก้ไข
          </Button>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <p className="text-sm text-slate-500">ชื่อสินค้า</p>
                <p className="text-lg font-semibold text-slate-800">{product.name}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">รายละเอียด</p>
                <p className="text-slate-700">{product.description}</p>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-slate-500">หมวดหมู่</p>
                  <Badge variant="info">{categoryLabels[product.category]}</Badge>
                </div>
                <div>
                  <p className="text-sm text-slate-500">หน่วย</p>
                  <p className="font-medium">{product.unit} ({product.unitSize})</p>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-slate-500">รหัสสินค้า (SKU)</p>
                  <p className="font-medium font-mono">{product.sku}</p>
                </div>
                <div>
                  <p className="text-sm text-slate-500">บาร์โค้ด</p>
                  <p className="font-medium font-mono">{product.barcode}</p>
                </div>
              </div>
              <div>
                <p className="text-sm text-slate-500">ผู้จัดจำหน่าย</p>
                <p className="font-medium">{product.supplier}</p>
              </div>
            </div>

            <div className="space-y-4">
              {/* Pricing */}
              <div className="p-4 bg-slate-50 rounded-lg space-y-3">
                <h4 className="font-medium text-slate-800">ราคา</h4>
                <div className="flex justify-between">
                  <span className="text-slate-500">ราคาต้นทุน</span>
                  <span className="font-medium">{formatCurrency(product.costPrice)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">ราคาขาย (ทั่วไป)</span>
                  <span className="font-semibold text-lg">{formatCurrency(product.sellingPrice)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">ราคาสมาชิก</span>
                  <span className="font-semibold text-green-600">{formatCurrency(product.memberPrice)}</span>
                </div>
                <div className="pt-2 border-t flex justify-between">
                  <span className="text-slate-500">กำไร/หน่วย</span>
                  <span className="font-semibold text-blue-600">
                    {formatCurrency(profit)} ({profitPercent.toFixed(1)}%)
                  </span>
                </div>
              </div>

              {/* Stock Levels */}
              <div className="p-4 bg-slate-50 rounded-lg space-y-3">
                <h4 className="font-medium text-slate-800">ระดับสต็อก</h4>
                <div className="flex justify-between">
                  <span className="text-slate-500">สต็อกต่ำสุด</span>
                  <span className="font-medium">{formatNumber(product.minStock)} {product.unit}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">สต็อกสูงสุด</span>
                  <span className="font-medium">{formatNumber(product.maxStock)} {product.unit}</span>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Stock Status */}
      <Card>
        <CardHeader>
          <CardTitle>สถานะคลัง</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            {/* Current Stock */}
            <div className={`p-6 rounded-xl text-center ${isLowStock ? "bg-red-50" : "bg-green-50"}`}>
              <p className={`text-sm ${isLowStock ? "text-red-600" : "text-green-600"}`}>คงเหลือ</p>
              <p className={`text-4xl font-bold mt-1 ${isLowStock ? "text-red-700" : "text-green-700"}`}>
                {formatNumber(product.stock)}
              </p>
              <p className={`text-sm ${isLowStock ? "text-red-500" : "text-green-500"}`}>{product.unit}</p>
              {isLowStock && (
                <div className="mt-3 flex items-center justify-center gap-2 text-red-600">
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  <span className="text-sm font-medium">สินค้าใกล้หมด</span>
                </div>
              )}
            </div>

            {/* Stock Progress */}
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-slate-500">ระดับสต็อก</span>
                <span className="font-medium">{((product.stock / product.maxStock) * 100).toFixed(0)}%</span>
              </div>
              <div className="h-3 bg-slate-100 rounded-full overflow-hidden">
                <div
                  className={`h-full rounded-full transition-all ${
                    isLowStock ? "bg-red-500" : product.stock > product.maxStock * 0.5 ? "bg-green-500" : "bg-amber-500"
                  }`}
                  style={{ width: `${Math.min((product.stock / product.maxStock) * 100, 100)}%` }}
                />
              </div>
              <div className="flex justify-between text-xs text-slate-400">
                <span>0</span>
                <span>{formatNumber(product.maxStock)}</span>
              </div>
            </div>

            {/* Quick Actions */}
            <div className="space-y-2">
              <Button className="w-full" onClick={() => setShowStockInModal(true)}>
                <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                นำเข้าสินค้า
              </Button>
              <Button variant="outline" className="w-full">
                <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
                ตรวจนับสต็อก
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderStockTab = () => (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>ประวัติเคลื่อนไหวสต็อก</CardTitle>
        <Button onClick={() => setShowStockInModal(true)}>
          <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          นำเข้าสินค้า
        </Button>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>วันที่</TableHead>
              <TableHead>รายการ</TableHead>
              <TableHead className="text-right">เปลี่ยนแปลง</TableHead>
              <TableHead className="text-right">คงเหลือ</TableHead>
              <TableHead>อ้างอิง</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {stockMovements.map((movement) => (
              <TableRow key={movement.id}>
                <TableCell>{formatDate(movement.date)}</TableCell>
                <TableCell>
                  <div className="flex items-center gap-3">
                    {getMovementIcon(movement.type)}
                    <span>{movement.description}</span>
                  </div>
                </TableCell>
                <TableCell className={`text-right font-medium ${movement.quantity >= 0 ? "text-green-600" : "text-red-600"}`}>
                  {movement.quantity >= 0 ? "+" : ""}{formatNumber(movement.quantity)}
                </TableCell>
                <TableCell className="text-right">{formatNumber(movement.balance)}</TableCell>
                <TableCell className="text-blue-600 font-mono text-sm">{movement.reference}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );

  const renderSalesTab = () => (
    <div className="space-y-6">
      {/* Sales Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <StatsCard
          title="ยอดขายรวม"
          value={formatNumber(salesSummary.totalSold)}
          subtitle={product.unit}
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" /></svg>}
        />
        <StatsCard
          title="รายได้รวม"
          value={formatCurrency(salesSummary.totalRevenue)}
          subtitle="บาท"
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>}
        />
        <StatsCard
          title="กำไรรวม"
          value={formatCurrency(salesSummary.totalProfit)}
          subtitle="บาท"
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" /></svg>}
        />
        <StatsCard
          title="ขายเฉลี่ย/เดือน"
          value={formatNumber(salesSummary.avgMonthlySales)}
          subtitle={product.unit}
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>}
        />
      </div>

      {/* Top Buyers */}
      <Card>
        <CardHeader>
          <CardTitle>ลูกค้าซื้อมากที่สุด</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>สมาชิก</TableHead>
                <TableHead className="text-right">จำนวน ({product.unit})</TableHead>
                <TableHead className="text-right">มูลค่า (บาท)</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {salesSummary.topBuyers.map((buyer, index) => (
                <TableRow key={buyer.memberId}>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <div className={`w-8 h-8 rounded-full flex items-center justify-center text-white font-bold ${
                        index === 0 ? "bg-amber-500" : index === 1 ? "bg-slate-400" : "bg-amber-700"
                      }`}>
                        {index + 1}
                      </div>
                      <div>
                        <Link href={`/members/${buyer.memberId}`} className="font-medium text-blue-600 hover:underline">
                          {buyer.name}
                        </Link>
                        <p className="text-sm text-slate-500">รหัส: {buyer.memberId}</p>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell className="text-right font-medium">{formatNumber(buyer.quantity)}</TableCell>
                  <TableCell className="text-right font-semibold text-green-600">{formatCurrency(buyer.amount)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Link href="/products">
            <Button variant="ghost" size="sm">
              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              กลับ
            </Button>
          </Link>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-bold text-slate-800">{product.name}</h1>
              {isLowStock && <Badge variant="danger">สินค้าใกล้หมด</Badge>}
            </div>
            <p className="text-slate-500">รหัส: {product.sku}</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 border-b border-slate-200 pb-2">
        {([
          { id: "overview", label: "ข้อมูลสินค้า" },
          { id: "stock", label: "ประวัติสต็อก" },
          { id: "sales", label: "ยอดขาย" },
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
      {activeTab === "stock" && renderStockTab()}
      {activeTab === "sales" && renderSalesTab()}

      {/* Stock In Modal */}
      <Modal isOpen={showStockInModal} onClose={() => setShowStockInModal(false)} size="md">
        <ModalHeader onClose={() => setShowStockInModal(false)}>นำเข้าสินค้า</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <p className="text-sm text-blue-600">สต็อกปัจจุบัน</p>
              <p className="text-2xl font-bold text-blue-700">{formatNumber(product.stock)} {product.unit}</p>
            </div>
            <Input label={`จำนวนที่นำเข้า (${product.unit})`} type="number" required />
            <Input label="ราคาต้นทุน/หน่วย (บาท)" type="number" defaultValue={product.costPrice} required />
            <Input label="เลขที่ใบสั่งซื้อ" placeholder="PO-XXXX-XXXX" />
            <Textarea label="หมายเหตุ" rows={2} />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowStockInModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowStockInModal(false)}>บันทึก</Button>
        </ModalFooter>
      </Modal>

      {/* Edit Modal */}
      <Modal isOpen={showEditModal} onClose={() => setShowEditModal(false)} size="lg">
        <ModalHeader onClose={() => setShowEditModal(false)}>แก้ไขข้อมูลสินค้า</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="ชื่อสินค้า" defaultValue={product.name} required />
            <Textarea label="รายละเอียด" defaultValue={product.description} rows={2} />
            <div className="grid grid-cols-2 gap-4">
              <Input label="ราคาต้นทุน (บาท)" type="number" defaultValue={product.costPrice} required />
              <Input label="ราคาขาย (บาท)" type="number" defaultValue={product.sellingPrice} required />
            </div>
            <Input label="ราคาสมาชิก (บาท)" type="number" defaultValue={product.memberPrice} required />
            <div className="grid grid-cols-2 gap-4">
              <Input label="สต็อกต่ำสุด" type="number" defaultValue={product.minStock} />
              <Input label="สต็อกสูงสุด" type="number" defaultValue={product.maxStock} />
            </div>
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
