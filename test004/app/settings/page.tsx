"use client";

import { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input, Select, Textarea } from "@/components/ui/input";
import { Modal, ModalHeader, ModalContent, ModalFooter, ConfirmDialog } from "@/components/ui/modal";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell, Badge } from "@/components/ui/table";
import { formatCurrency, formatDate, thaiProvinces } from "@/lib/utils";

type SettingsTab = "cooperative" | "shares" | "loans" | "deposits" | "users" | "notifications" | "system";

// Mock data
const mockCooperative = {
  name: "สหกรณ์การเกษตรตำบลบางกระทึก จำกัด",
  registrationNumber: "กบส.123/2540",
  taxId: "1234567890123",
  address: "123 หมู่ 4 ตำบลบางกระทึก อำเภอสามพราน จังหวัดนครปฐม 73210",
  phone: "034-123456",
  email: "contact@bangkrateuk-coop.or.th",
  website: "www.bangkrateuk-coop.or.th",
  foundedDate: "1997-05-15",
  type: "agricultural",
  province: "นครปฐม",
};

const mockBylaws = {
  shareValue: 100,
  minSharesRequired: 10,
  maxSharesAllowed: 20000,
  monthlySharePayment: 100,
  shareTransferFee: 50,
  membershipFee: 100,
  withdrawalPenaltyPercent: 2,
};

const mockLoanProducts = [
  { id: "LP001", name: "เงินกู้ฉุกเฉิน", type: "emergency", interestRate: 12, maxAmount: 50000, maxTerm: 12, status: "active" },
  { id: "LP002", name: "เงินกู้สามัญ", type: "ordinary", interestRate: 9, maxAmount: 500000, maxTerm: 60, status: "active" },
  { id: "LP003", name: "เงินกู้พิเศษ", type: "special", interestRate: 7.5, maxAmount: 2000000, maxTerm: 120, status: "active" },
];

const mockDepositProducts = [
  { id: "DP001", name: "ออมทรัพย์", type: "savings", interestRate: 1.5, minBalance: 0, status: "active" },
  { id: "DP002", name: "ออมทรัพย์พิเศษ", type: "special_savings", interestRate: 2.5, minBalance: 10000, status: "active" },
  { id: "DP003", name: "ประจำ 12 เดือน", type: "fixed", interestRate: 3.0, minBalance: 5000, status: "active" },
];

const mockUsers = [
  { id: "USR001", name: "นายสมชาย ใจดี", email: "somchai@coop.or.th", role: "admin", status: "active", lastLogin: "2024-01-15 10:30" },
  { id: "USR002", name: "นางสาวสมหญิง รักษ์ดี", email: "somying@coop.or.th", role: "accountant", status: "active", lastLogin: "2024-01-15 09:15" },
  { id: "USR003", name: "นายวิชัย พัฒนา", email: "wichai@coop.or.th", role: "staff", status: "active", lastLogin: "2024-01-14 16:45" },
];

const roleLabels: Record<string, string> = {
  admin: "ผู้ดูแลระบบ",
  manager: "ผู้จัดการ",
  accountant: "เจ้าหน้าที่การเงิน",
  staff: "เจ้าหน้าที่ทั่วไป",
  auditor: "ผู้ตรวจสอบ",
};

const loanTypeLabels: Record<string, string> = {
  emergency: "ฉุกเฉิน",
  ordinary: "สามัญ",
  special: "พิเศษ",
};

const depositTypeLabels: Record<string, string> = {
  savings: "ออมทรัพย์",
  special_savings: "ออมทรัพย์พิเศษ",
  fixed: "ประจำ",
};

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<SettingsTab>("cooperative");
  const [showUserModal, setShowUserModal] = useState(false);
  const [showLoanProductModal, setShowLoanProductModal] = useState(false);
  const [showDepositProductModal, setShowDepositProductModal] = useState(false);
  const [showSaveConfirm, setShowSaveConfirm] = useState(false);

  const tabs: { id: SettingsTab; label: string; icon: string }[] = [
    { id: "cooperative", label: "ข้อมูลสหกรณ์", icon: "M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" },
    { id: "shares", label: "ทุนเรือนหุ้น", icon: "M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" },
    { id: "loans", label: "ผลิตภัณฑ์สินเชื่อ", icon: "M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" },
    { id: "deposits", label: "ผลิตภัณฑ์เงินฝาก", icon: "M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" },
    { id: "users", label: "ผู้ใช้งาน", icon: "M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" },
    { id: "notifications", label: "การแจ้งเตือน", icon: "M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" },
    { id: "system", label: "ระบบ", icon: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" },
  ];

  const renderCooperativeSettings = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>ข้อมูลทั่วไป</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input label="ชื่อสหกรณ์" defaultValue={mockCooperative.name} />
            <Input label="เลขทะเบียนสหกรณ์" defaultValue={mockCooperative.registrationNumber} />
            <Input label="เลขประจำตัวผู้เสียภาษี" defaultValue={mockCooperative.taxId} />
            <Select label="ประเภทสหกรณ์" defaultValue={mockCooperative.type}>
              <option value="agricultural">สหกรณ์การเกษตร</option>
              <option value="savings">สหกรณ์ออมทรัพย์</option>
              <option value="consumer">สหกรณ์ร้านค้า</option>
              <option value="service">สหกรณ์บริการ</option>
              <option value="credit_union">สหกรณ์เครดิตยูเนี่ยน</option>
            </Select>
            <Input label="วันที่จดทะเบียน" type="date" defaultValue={mockCooperative.foundedDate} />
            <Select label="จังหวัด" defaultValue={mockCooperative.province}>
              {thaiProvinces.map((province) => (
                <option key={province} value={province}>{province}</option>
              ))}
            </Select>
            <div className="md:col-span-2">
              <Textarea label="ที่อยู่" defaultValue={mockCooperative.address} rows={2} />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>ช่องทางติดต่อ</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input label="หมายเลขโทรศัพท์" defaultValue={mockCooperative.phone} />
            <Input label="อีเมล" type="email" defaultValue={mockCooperative.email} />
            <Input label="เว็บไซต์" defaultValue={mockCooperative.website} />
            <Input label="LINE Official Account" placeholder="@example" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>โลโก้และตราประทับ</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">โลโก้สหกรณ์</label>
              <div className="border-2 border-dashed border-slate-300 rounded-lg p-6 text-center hover:border-blue-400 transition-colors cursor-pointer">
                <svg className="w-12 h-12 mx-auto text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <p className="mt-2 text-sm text-slate-500">คลิกเพื่ออัพโหลดรูปภาพ</p>
                <p className="text-xs text-slate-400">PNG, JPG ขนาดไม่เกิน 2MB</p>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">ตราประทับ</label>
              <div className="border-2 border-dashed border-slate-300 rounded-lg p-6 text-center hover:border-blue-400 transition-colors cursor-pointer">
                <svg className="w-12 h-12 mx-auto text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
                <p className="mt-2 text-sm text-slate-500">คลิกเพื่ออัพโหลดรูปภาพ</p>
                <p className="text-xs text-slate-400">PNG ขนาดไม่เกิน 1MB</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderSharesSettings = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>ข้อบังคับทุนเรือนหุ้น</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="มูลค่าหุ้น (บาท/หุ้น)"
              type="number"
              defaultValue={mockBylaws.shareValue}
              hint="ตามข้อบังคับสหกรณ์"
            />
            <Input
              label="หุ้นขั้นต่ำที่ต้องถือ (หุ้น)"
              type="number"
              defaultValue={mockBylaws.minSharesRequired}
              hint="จำนวนหุ้นขั้นต่ำเมื่อสมัครสมาชิก"
            />
            <Input
              label="หุ้นสูงสุดที่ถือได้ (หุ้น)"
              type="number"
              defaultValue={mockBylaws.maxSharesAllowed}
              hint="ไม่เกิน 1/5 ของทุนเรือนหุ้นทั้งหมด"
            />
            <Input
              label="ค่าหุ้นรายเดือนขั้นต่ำ (บาท)"
              type="number"
              defaultValue={mockBylaws.monthlySharePayment}
              hint="ยอดหักจากเงินเดือนหรือชำระด้วยตนเอง"
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>ค่าธรรมเนียม</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="ค่าธรรมเนียมแรกเข้า (บาท)"
              type="number"
              defaultValue={mockBylaws.membershipFee}
            />
            <Input
              label="ค่าธรรมเนียมโอนหุ้น (บาท)"
              type="number"
              defaultValue={mockBylaws.shareTransferFee}
            />
            <Input
              label="อัตราหักเงินถอนหุ้นก่อนกำหนด (%)"
              type="number"
              step="0.01"
              defaultValue={mockBylaws.withdrawalPenaltyPercent}
              hint="กรณีลาออกก่อนครบ 1 ปี"
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>การจ่ายเงินปันผล</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="อัตราเงินปันผลปีก่อน (%)"
              type="number"
              step="0.01"
              defaultValue={5.5}
              disabled
            />
            <Select label="วิธีการจ่ายเงินปันผล">
              <option value="deposit">โอนเข้าบัญชีเงินฝาก</option>
              <option value="share">ซื้อหุ้นเพิ่ม</option>
              <option value="cash">เงินสด</option>
            </Select>
            <Input
              label="วันที่จ่ายเงินปันผลประจำปี"
              type="date"
              defaultValue="2024-03-31"
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderLoansSettings = () => (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-lg font-semibold text-slate-800">ผลิตภัณฑ์สินเชื่อ</h2>
          <p className="text-sm text-slate-500">กำหนดประเภทและเงื่อนไขสินเชื่อของสหกรณ์</p>
        </div>
        <Button onClick={() => setShowLoanProductModal(true)}>
          <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          เพิ่มผลิตภัณฑ์
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ชื่อผลิตภัณฑ์</TableHead>
                <TableHead>ประเภท</TableHead>
                <TableHead className="text-right">ดอกเบี้ย (%/ปี)</TableHead>
                <TableHead className="text-right">วงเงินสูงสุด</TableHead>
                <TableHead className="text-right">ระยะเวลา (เดือน)</TableHead>
                <TableHead>สถานะ</TableHead>
                <TableHead className="text-right">จัดการ</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {mockLoanProducts.map((product) => (
                <TableRow key={product.id}>
                  <TableCell className="font-medium">{product.name}</TableCell>
                  <TableCell>{loanTypeLabels[product.type]}</TableCell>
                  <TableCell className="text-right">{product.interestRate}</TableCell>
                  <TableCell className="text-right">{formatCurrency(product.maxAmount)}</TableCell>
                  <TableCell className="text-right">{product.maxTerm}</TableCell>
                  <TableCell>
                    <Badge variant="success">ใช้งาน</Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button variant="ghost" size="sm">แก้ไข</Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>เงื่อนไขทั่วไป</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="อัตราค่าปรับล่าช้า (%/เดือน)"
              type="number"
              step="0.01"
              defaultValue={1.5}
            />
            <Input
              label="จำนวนวันผ่อนผันก่อนคิดค่าปรับ"
              type="number"
              defaultValue={7}
            />
            <Select label="วิธีคำนวณดอกเบี้ย">
              <option value="flat">อัตราคงที่ (Flat Rate)</option>
              <option value="effective">อัตราลดต้นลดดอก (Effective Rate)</option>
            </Select>
            <Select label="รอบตัดชำระ">
              <option value="monthly">รายเดือน</option>
              <option value="biweekly">ทุก 2 สัปดาห์</option>
              <option value="weekly">รายสัปดาห์</option>
            </Select>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderDepositsSettings = () => (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-lg font-semibold text-slate-800">ผลิตภัณฑ์เงินฝาก</h2>
          <p className="text-sm text-slate-500">กำหนดประเภทและอัตราดอกเบี้ยเงินฝาก</p>
        </div>
        <Button onClick={() => setShowDepositProductModal(true)}>
          <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          เพิ่มผลิตภัณฑ์
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ชื่อผลิตภัณฑ์</TableHead>
                <TableHead>ประเภท</TableHead>
                <TableHead className="text-right">ดอกเบี้ย (%/ปี)</TableHead>
                <TableHead className="text-right">ยอดขั้นต่ำ</TableHead>
                <TableHead>สถานะ</TableHead>
                <TableHead className="text-right">จัดการ</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {mockDepositProducts.map((product) => (
                <TableRow key={product.id}>
                  <TableCell className="font-medium">{product.name}</TableCell>
                  <TableCell>{depositTypeLabels[product.type]}</TableCell>
                  <TableCell className="text-right">{product.interestRate}</TableCell>
                  <TableCell className="text-right">{formatCurrency(product.minBalance)}</TableCell>
                  <TableCell>
                    <Badge variant="success">ใช้งาน</Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button variant="ghost" size="sm">แก้ไข</Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>การคำนวณดอกเบี้ย</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Select label="รอบคำนวณดอกเบี้ย">
              <option value="daily">รายวัน</option>
              <option value="monthly">รายเดือน</option>
            </Select>
            <Select label="รอบจ่ายดอกเบี้ย (ออมทรัพย์)">
              <option value="monthly">ทุกเดือน</option>
              <option value="quarterly">ทุกไตรมาส</option>
              <option value="yearly">ทุกปี</option>
            </Select>
            <Input
              label="ยอดเงินขั้นต่ำเพื่อรับดอกเบี้ย"
              type="number"
              defaultValue={100}
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderUsersSettings = () => (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-lg font-semibold text-slate-800">ผู้ใช้งานระบบ</h2>
          <p className="text-sm text-slate-500">จัดการบัญชีผู้ใช้และสิทธิ์การเข้าถึง</p>
        </div>
        <Button onClick={() => setShowUserModal(true)}>
          <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
          </svg>
          เพิ่มผู้ใช้
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ชื่อ</TableHead>
                <TableHead>อีเมล</TableHead>
                <TableHead>บทบาท</TableHead>
                <TableHead>เข้าใช้งานล่าสุด</TableHead>
                <TableHead>สถานะ</TableHead>
                <TableHead className="text-right">จัดการ</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {mockUsers.map((user) => (
                <TableRow key={user.id}>
                  <TableCell className="font-medium">{user.name}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>
                    <Badge variant={user.role === "admin" ? "info" : "default"}>
                      {roleLabels[user.role]}
                    </Badge>
                  </TableCell>
                  <TableCell>{user.lastLogin}</TableCell>
                  <TableCell>
                    <Badge variant="success">ใช้งาน</Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button variant="ghost" size="sm">แก้ไข</Button>
                      <Button variant="ghost" size="sm" className="text-red-600">ระงับ</Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>บทบาทและสิทธิ์</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {Object.entries(roleLabels).map(([role, label]) => (
              <div key={role} className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
                <div>
                  <p className="font-medium text-slate-800">{label}</p>
                  <p className="text-sm text-slate-500">
                    {role === "admin" && "สิทธิ์เต็ม สามารถจัดการทุกส่วนของระบบ"}
                    {role === "manager" && "จัดการสมาชิก สินเชื่อ และดูรายงาน"}
                    {role === "accountant" && "จัดการธุรกรรมทางการเงิน และบัญชี"}
                    {role === "staff" && "ดูข้อมูล และทำรายการพื้นฐาน"}
                    {role === "auditor" && "ดูรายงานและตรวจสอบเอกสาร"}
                  </p>
                </div>
                <Button variant="outline" size="sm">กำหนดสิทธิ์</Button>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderNotificationsSettings = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>การแจ้งเตือนทางอีเมล</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-slate-800">สมาชิกใหม่สมัคร</p>
                <p className="text-sm text-slate-500">แจ้งเตือนเมื่อมีสมาชิกใหม่ลงทะเบียน</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" defaultChecked />
                <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-slate-800">คำขอสินเชื่อใหม่</p>
                <p className="text-sm text-slate-500">แจ้งเตือนเมื่อมีคำขอสินเชื่อรอพิจารณา</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" defaultChecked />
                <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-slate-800">สินเชื่อค้างชำระ</p>
                <p className="text-sm text-slate-500">แจ้งเตือนเมื่อมีสินเชื่อค้างชำระเกินกำหนด</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" defaultChecked />
                <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-slate-800">สินค้าใกล้หมด</p>
                <p className="text-sm text-slate-500">แจ้งเตือนเมื่อสินค้าในคลังต่ำกว่าจุดสั่งซื้อ</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" defaultChecked />
                <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>การแจ้งเตือนทาง LINE</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <Input
              label="LINE Notify Token"
              type="password"
              placeholder="กรอก Token เพื่อเชื่อมต่อ LINE Notify"
            />
            <div className="flex items-center gap-4">
              <Button variant="outline">ทดสอบการเชื่อมต่อ</Button>
              <span className="text-sm text-slate-500">ยังไม่ได้เชื่อมต่อ</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>การแจ้งเตือนอัตโนมัติ</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-slate-800">แจ้งเตือนก่อนครบกำหนดชำระ</p>
                <p className="text-sm text-slate-500">แจ้งเตือนสมาชิกก่อนครบกำหนดชำระหนี้</p>
              </div>
              <Select className="w-32">
                <option value="3">3 วัน</option>
                <option value="5">5 วัน</option>
                <option value="7">7 วัน</option>
              </Select>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-slate-800">แจ้งเตือนเงินฝากครบกำหนด</p>
                <p className="text-sm text-slate-500">แจ้งเตือนสมาชิกก่อนเงินฝากประจำครบกำหนด</p>
              </div>
              <Select className="w-32">
                <option value="7">7 วัน</option>
                <option value="14">14 วัน</option>
                <option value="30">30 วัน</option>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderSystemSettings = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>ข้อมูลระบบ</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-500">เวอร์ชัน</p>
              <p className="text-lg font-semibold text-slate-800">1.0.0</p>
            </div>
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-500">อัพเดทล่าสุด</p>
              <p className="text-lg font-semibold text-slate-800">{formatDate("2024-01-15")}</p>
            </div>
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-500">ฐานข้อมูล</p>
              <p className="text-lg font-semibold text-slate-800">MongoDB Atlas</p>
            </div>
            <div className="p-4 bg-slate-50 rounded-lg">
              <p className="text-sm text-slate-500">สถานะการเชื่อมต่อ</p>
              <p className="text-lg font-semibold text-green-600">ปกติ</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>การตั้งค่าทั่วไป</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Select label="ภาษา">
              <option value="th">ไทย</option>
              <option value="en">English</option>
            </Select>
            <Select label="รูปแบบวันที่">
              <option value="th">วัน/เดือน/ปี (พ.ศ.)</option>
              <option value="en">DD/MM/YYYY (A.D.)</option>
            </Select>
            <Select label="เขตเวลา">
              <option value="Asia/Bangkok">Asia/Bangkok (GMT+7)</option>
            </Select>
            <Select label="สกุลเงิน">
              <option value="THB">บาท (THB)</option>
            </Select>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>สำรองข้อมูล</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
              <div>
                <p className="font-medium text-slate-800">สำรองข้อมูลอัตโนมัติ</p>
                <p className="text-sm text-slate-500">ระบบจะสำรองข้อมูลทุกวันเวลา 02:00 น.</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" defaultChecked />
                <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div className="flex items-center gap-4">
              <Button variant="outline">
                <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                </svg>
                ส่งออกข้อมูล
              </Button>
              <Button variant="outline">
                <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                นำเข้าข้อมูล
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="border-red-200 bg-red-50">
        <CardHeader>
          <CardTitle className="text-red-700">โซนอันตราย</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-red-800">รีเซ็ตข้อมูลทดสอบ</p>
                <p className="text-sm text-red-600">ลบข้อมูลทดสอบทั้งหมด (ใช้เฉพาะในโหมด Development)</p>
              </div>
              <Button variant="danger" size="sm">รีเซ็ต</Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-800">ตั้งค่า</h1>
          <p className="text-slate-500 mt-1">จัดการการตั้งค่าระบบและข้อมูลสหกรณ์</p>
        </div>
        <Button onClick={() => setShowSaveConfirm(true)}>
          <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
          บันทึกการตั้งค่า
        </Button>
      </div>

      {/* Tabs */}
      <div className="flex flex-wrap gap-2 border-b border-slate-200 pb-2">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              activeTab === tab.id
                ? "bg-blue-50 text-blue-700"
                : "text-slate-600 hover:bg-slate-100"
            }`}
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={tab.icon} />
            </svg>
            <span className="hidden sm:inline">{tab.label}</span>
          </button>
        ))}
      </div>

      {/* Content */}
      {activeTab === "cooperative" && renderCooperativeSettings()}
      {activeTab === "shares" && renderSharesSettings()}
      {activeTab === "loans" && renderLoansSettings()}
      {activeTab === "deposits" && renderDepositsSettings()}
      {activeTab === "users" && renderUsersSettings()}
      {activeTab === "notifications" && renderNotificationsSettings()}
      {activeTab === "system" && renderSystemSettings()}

      {/* User Modal */}
      <Modal isOpen={showUserModal} onClose={() => setShowUserModal(false)} size="md">
        <ModalHeader onClose={() => setShowUserModal(false)}>
          เพิ่มผู้ใช้งาน
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="ชื่อ-นามสกุล" required />
            <Input label="อีเมล" type="email" required />
            <Input label="รหัสผ่านชั่วคราว" type="password" required />
            <Select label="บทบาท" required>
              <option value="">เลือกบทบาท</option>
              {Object.entries(roleLabels).map(([value, label]) => (
                <option key={value} value={value}>{label}</option>
              ))}
            </Select>
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowUserModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowUserModal(false)}>เพิ่มผู้ใช้</Button>
        </ModalFooter>
      </Modal>

      {/* Loan Product Modal */}
      <Modal isOpen={showLoanProductModal} onClose={() => setShowLoanProductModal(false)} size="md">
        <ModalHeader onClose={() => setShowLoanProductModal(false)}>
          เพิ่มผลิตภัณฑ์สินเชื่อ
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="ชื่อผลิตภัณฑ์" required />
            <Select label="ประเภท" required>
              <option value="">เลือกประเภท</option>
              {Object.entries(loanTypeLabels).map(([value, label]) => (
                <option key={value} value={value}>{label}</option>
              ))}
            </Select>
            <Input label="อัตราดอกเบี้ย (%/ปี)" type="number" step="0.01" required />
            <Input label="วงเงินสูงสุด (บาท)" type="number" required />
            <Input label="ระยะเวลาสูงสุด (เดือน)" type="number" required />
            <Textarea label="เงื่อนไขพิเศษ" rows={3} />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowLoanProductModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowLoanProductModal(false)}>บันทึก</Button>
        </ModalFooter>
      </Modal>

      {/* Deposit Product Modal */}
      <Modal isOpen={showDepositProductModal} onClose={() => setShowDepositProductModal(false)} size="md">
        <ModalHeader onClose={() => setShowDepositProductModal(false)}>
          เพิ่มผลิตภัณฑ์เงินฝาก
        </ModalHeader>
        <ModalContent>
          <div className="space-y-4">
            <Input label="ชื่อผลิตภัณฑ์" required />
            <Select label="ประเภท" required>
              <option value="">เลือกประเภท</option>
              {Object.entries(depositTypeLabels).map(([value, label]) => (
                <option key={value} value={value}>{label}</option>
              ))}
            </Select>
            <Input label="อัตราดอกเบี้ย (%/ปี)" type="number" step="0.01" required />
            <Input label="ยอดขั้นต่ำ (บาท)" type="number" required />
            <Input label="ระยะเวลาฝาก (เดือน)" type="number" hint="เฉพาะเงินฝากประจำ" />
          </div>
        </ModalContent>
        <ModalFooter>
          <Button variant="outline" onClick={() => setShowDepositProductModal(false)}>ยกเลิก</Button>
          <Button onClick={() => setShowDepositProductModal(false)}>บันทึก</Button>
        </ModalFooter>
      </Modal>

      {/* Save Confirm Dialog */}
      <ConfirmDialog
        isOpen={showSaveConfirm}
        onClose={() => setShowSaveConfirm(false)}
        onConfirm={() => {
          setShowSaveConfirm(false);
          // Save logic here
        }}
        title="บันทึกการตั้งค่า"
        message="คุณต้องการบันทึกการเปลี่ยนแปลงทั้งหมดหรือไม่?"
        confirmText="บันทึก"
        cancelText="ยกเลิก"
      />
    </div>
  );
}
