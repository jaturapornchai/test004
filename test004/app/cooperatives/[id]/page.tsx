"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent, StatsCard } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input, Select, Textarea } from "@/components/ui/input";
import { Modal, ModalHeader, ModalContent, ModalFooter } from "@/components/ui/modal";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell, Badge } from "@/components/ui/table";
import { formatCurrency, formatNumber, formatDate, formatPhone } from "@/lib/utils";

// Mock cooperative data
const getCooperativeData = (id: string) => ({
  _id: id,
  coopId: id,
  name: "สหกรณ์การเกษตรตำบลบางกระทึก จำกัด",
  nameEn: "Bangkrateuk Agricultural Cooperative Limited",
  registrationNumber: "กบส.123/2540",
  taxId: "1234567890123",
  foundedDate: "1997-05-15",
  type: "agricultural" as const,
  status: "active" as const,
  address: {
    houseNo: "123",
    moo: "4",
    road: "เพชรเกษม",
    subDistrict: "บางกระทึก",
    district: "สามพราน",
    province: "นครปฐม",
    postalCode: "73210",
  },
  contact: {
    phone: "034-123456",
    fax: "034-123457",
    email: "contact@bangkrateuk-coop.or.th",
    website: "www.bangkrateuk-coop.or.th",
    lineId: "@bangkrateuk-coop",
  },
  bylaws: {
    shareValue: 100,
    minSharesRequired: 10,
    maxSharesAllowed: 20000,
    membershipFee: 100,
    dividendPaymentMonth: 3,
    fiscalYearEnd: "12-31",
  },
  boardMembers: [
    { name: "นายประสิทธิ์ มั่นคง", position: "ประธานกรรมการ", phone: "081-234-5678", startDate: "2023-01-01" },
    { name: "นางสมศรี รักดี", position: "รองประธานกรรมการ", phone: "082-345-6789", startDate: "2023-01-01" },
    { name: "นายวิชัย พัฒนา", position: "เลขานุการ", phone: "083-456-7890", startDate: "2023-01-01" },
    { name: "นางมาลี สุขสันต์", position: "เหรัญญิก", phone: "084-567-8901", startDate: "2023-01-01" },
    { name: "นายสมชาย ใจดี", position: "กรรมการ", phone: "085-678-9012", startDate: "2023-01-01" },
  ],
  createdAt: "2020-01-01T00:00:00Z",
  updatedAt: "2024-01-15T14:20:00Z",
});

// Mock statistics
const getCooperativeStats = () => ({
  totalMembers: 1247,
  activeMembers: 1189,
  totalShares: 456789,
  shareCapital: 45678900,
  totalDeposits: 128456780,
  totalLoans: 89234560,
  overdueLoans: 3,
  monthlyRevenue: 2345670,
  monthlyExpense: 1234560,
  netProfit: 1111110,
  totalProducts: 45,
  lowStockProducts: 3,
});

const coopTypeLabels: Record<string, string> = {
  agricultural: "สหกรณ์การเกษตร",
  savings: "สหกรณ์ออมทรัพย์",
  consumer: "สหกรณ์ร้านค้า",
  service: "สหกรณ์บริการ",
  credit_union: "สหกรณ์เครดิตยูเนี่ยน",
};

type TabType = "overview" | "board" | "bylaws" | "statistics";

export default function CooperativeDetailPage() {
  const params = useParams();
  const coopId = params.id as string;

  const [activeTab, setActiveTab] = useState<TabType>("overview");
  const [showEditModal, setShowEditModal] = useState(false);

  const cooperative = getCooperativeData(coopId);
  const stats = getCooperativeStats();

  const renderOverviewTab = () => (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* Basic Info */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>ข้อมูลทั่วไป</CardTitle>
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
              <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                <svg className="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
              </div>
              <div>
                <h3 className="text-xl font-semibold text-slate-800">{cooperative.name}</h3>
                <p className="text-sm text-slate-500">{cooperative.nameEn}</p>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4 pt-4 border-t">
              <div>
                <p className="text-sm text-slate-500">เลขทะเบียน</p>
                <p className="font-medium">{cooperative.registrationNumber}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">เลขประจำตัวผู้เสียภาษี</p>
                <p className="font-medium">{cooperative.taxId}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">ประเภท</p>
                <p className="font-medium">{coopTypeLabels[cooperative.type]}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">วันที่จดทะเบียน</p>
                <p className="font-medium">{formatDate(cooperative.foundedDate)}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">สถานะ</p>
                <Badge variant="success">ใช้งาน</Badge>
              </div>
              <div>
                <p className="text-sm text-slate-500">สิ้นปีบัญชี</p>
                <p className="font-medium">31 ธันวาคม</p>
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
              {cooperative.address.houseNo} หมู่ {cooperative.address.moo}
            </p>
            {cooperative.address.road && <p className="text-slate-700">ถนน{cooperative.address.road}</p>}
            <p className="text-slate-700">
              ตำบล{cooperative.address.subDistrict} อำเภอ{cooperative.address.district}
            </p>
            <p className="text-slate-700">
              จังหวัด{cooperative.address.province} {cooperative.address.postalCode}
            </p>
          </div>

          <div className="mt-6 pt-4 border-t">
            <h4 className="font-medium text-slate-800 mb-3">ช่องทางติดต่อ</h4>
            <div className="space-y-2">
              <div className="flex items-center gap-3">
                <svg className="w-5 h-5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                </svg>
                <span>{cooperative.contact.phone}</span>
              </div>
              <div className="flex items-center gap-3">
                <svg className="w-5 h-5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                <span>{cooperative.contact.email}</span>
              </div>
              {cooperative.contact.website && (
                <div className="flex items-center gap-3">
                  <svg className="w-5 h-5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                  </svg>
                  <span>{cooperative.contact.website}</span>
                </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Quick Stats */}
      <Card className="lg:col-span-2">
        <CardHeader>
          <CardTitle>สรุปข้อมูล</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="p-4 bg-blue-50 rounded-lg text-center">
              <p className="text-sm text-blue-600">สมาชิก</p>
              <p className="text-2xl font-bold text-blue-700">{formatNumber(stats.totalMembers)}</p>
              <p className="text-xs text-blue-500">ราย</p>
            </div>
            <div className="p-4 bg-green-50 rounded-lg text-center">
              <p className="text-sm text-green-600">ทุนเรือนหุ้น</p>
              <p className="text-2xl font-bold text-green-700">{formatCurrency(stats.shareCapital)}</p>
              <p className="text-xs text-green-500">{formatNumber(stats.totalShares)} หุ้น</p>
            </div>
            <div className="p-4 bg-amber-50 rounded-lg text-center">
              <p className="text-sm text-amber-600">เงินฝาก</p>
              <p className="text-2xl font-bold text-amber-700">{formatCurrency(stats.totalDeposits)}</p>
              <p className="text-xs text-amber-500">ยอดรวม</p>
            </div>
            <div className="p-4 bg-purple-50 rounded-lg text-center">
              <p className="text-sm text-purple-600">สินเชื่อคงค้าง</p>
              <p className="text-2xl font-bold text-purple-700">{formatCurrency(stats.totalLoans)}</p>
              <p className="text-xs text-purple-500">ยอดรวม</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderBoardTab = () => (
    <Card>
      <CardHeader>
        <CardTitle>คณะกรรมการดำเนินการ</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {cooperative.boardMembers.map((member, index) => (
            <div key={index} className="p-4 border rounded-lg">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-slate-400 to-slate-500 flex items-center justify-center text-white font-bold">
                  {member.name.charAt(0)}
                </div>
                <div>
                  <p className="font-medium text-slate-800">{member.name}</p>
                  <p className="text-sm text-blue-600">{member.position}</p>
                </div>
              </div>
              <div className="mt-3 pt-3 border-t text-sm text-slate-500">
                <p>โทร: {member.phone}</p>
                <p>เริ่มดำรงตำแหน่ง: {formatDate(member.startDate)}</p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );

  const renderBylawsTab = () => (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <Card>
        <CardHeader>
          <CardTitle>ข้อบังคับทุนเรือนหุ้น</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">มูลค่าหุ้น</span>
              <span className="font-semibold">{formatCurrency(cooperative.bylaws.shareValue)}/หุ้น</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">หุ้นขั้นต่ำที่ต้องถือ</span>
              <span className="font-semibold">{formatNumber(cooperative.bylaws.minSharesRequired)} หุ้น</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">หุ้นสูงสุดที่ถือได้</span>
              <span className="font-semibold">{formatNumber(cooperative.bylaws.maxSharesAllowed)} หุ้น</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">ค่าธรรมเนียมแรกเข้า</span>
              <span className="font-semibold">{formatCurrency(cooperative.bylaws.membershipFee)}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>ปฏิทินกิจกรรม</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">สิ้นปีบัญชี</span>
              <span className="font-semibold">31 ธันวาคม</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">ประชุมใหญ่สามัญประจำปี</span>
              <span className="font-semibold">มีนาคม</span>
            </div>
            <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
              <span className="text-slate-600">จ่ายเงินปันผล</span>
              <span className="font-semibold">มีนาคม</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderStatisticsTab = () => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="สมาชิกทั้งหมด"
          value={formatNumber(stats.totalMembers)}
          subtitle={`ใช้งาน ${formatNumber(stats.activeMembers)} ราย`}
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" /></svg>}
          trend={{ value: 1.9, isPositive: true }}
        />
        <StatsCard
          title="ทุนเรือนหุ้น"
          value={formatCurrency(stats.shareCapital)}
          subtitle={`${formatNumber(stats.totalShares)} หุ้น`}
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" /></svg>}
          trend={{ value: 3.2, isPositive: true }}
        />
        <StatsCard
          title="เงินฝากรวม"
          value={formatCurrency(stats.totalDeposits)}
          subtitle="ยอดเงินฝากทั้งหมด"
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" /></svg>}
          trend={{ value: 5.1, isPositive: true }}
        />
        <StatsCard
          title="สินเชื่อคงค้าง"
          value={formatCurrency(stats.totalLoans)}
          subtitle={`ค้างชำระ ${stats.overdueLoans} รายการ`}
          icon={<svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>}
          trend={{ value: 2.3, isPositive: false }}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <Card className="bg-gradient-to-br from-green-500 to-green-600 text-white border-0">
          <CardContent className="pt-6">
            <p className="text-green-100 text-sm">รายได้เดือนนี้</p>
            <p className="text-2xl font-bold mt-1">{formatCurrency(stats.monthlyRevenue)}</p>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-red-500 to-red-600 text-white border-0">
          <CardContent className="pt-6">
            <p className="text-red-100 text-sm">ค่าใช้จ่ายเดือนนี้</p>
            <p className="text-2xl font-bold mt-1">{formatCurrency(stats.monthlyExpense)}</p>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-blue-500 to-blue-600 text-white border-0">
          <CardContent className="pt-6">
            <p className="text-blue-100 text-sm">กำไรสุทธิ</p>
            <p className="text-2xl font-bold mt-1">{formatCurrency(stats.netProfit)}</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Link href="/cooperatives">
            <Button variant="ghost" size="sm">
              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              กลับ
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-slate-800">{cooperative.name}</h1>
            <p className="text-slate-500">รหัส: {cooperative.coopId}</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 border-b border-slate-200 pb-2">
        {([
          { id: "overview", label: "ข้อมูลทั่วไป" },
          { id: "board", label: "คณะกรรมการ" },
          { id: "bylaws", label: "ข้อบังคับ" },
          { id: "statistics", label: "สถิติ" },
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
      {activeTab === "board" && renderBoardTab()}
      {activeTab === "bylaws" && renderBylawsTab()}
      {activeTab === "statistics" && renderStatisticsTab()}

      {/* Edit Modal */}
      <Modal isOpen={showEditModal} onClose={() => setShowEditModal(false)} size="lg">
        <ModalHeader onClose={() => setShowEditModal(false)}>แก้ไขข้อมูลสหกรณ์</ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="ชื่อสหกรณ์" defaultValue={cooperative.name} required />
            <Input label="ชื่อภาษาอังกฤษ" defaultValue={cooperative.nameEn} />
            <div className="grid grid-cols-2 gap-4">
              <Input label="โทรศัพท์" defaultValue={cooperative.contact.phone} />
              <Input label="อีเมล" type="email" defaultValue={cooperative.contact.email} />
            </div>
            <Input label="เว็บไซต์" defaultValue={cooperative.contact.website} />
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
