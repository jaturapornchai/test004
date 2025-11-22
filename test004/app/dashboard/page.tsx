"use client";

import { useState } from "react";
import Link from "next/link";
import { Card, CardHeader, CardTitle, CardContent, StatsCard } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell, Badge } from "@/components/ui/table";
import { formatCurrency, formatNumber, formatDate } from "@/lib/utils";

// Mock data for dashboard
const summaryData = {
  totalMembers: 1247,
  activeMembers: 1189,
  newMembersThisMonth: 23,
  totalShares: 45678900,
  totalDeposits: 128456780,
  totalLoans: 89234560,
  overdueLoans: 3,
  monthlyRevenue: 2345670,
  monthlyExpense: 1234560,
  netProfit: 1111110,
};

const recentMembers = [
  { id: "MEM001", name: "นายสมชาย ใจดี", joinDate: "2024-01-15", shares: 50000, status: "active" },
  { id: "MEM002", name: "นางสาวสมหญิง รักษ์ดี", joinDate: "2024-01-14", shares: 30000, status: "active" },
  { id: "MEM003", name: "นายวิชัย พัฒนา", joinDate: "2024-01-12", shares: 100000, status: "active" },
];

const recentLoans = [
  { id: "LN001", memberName: "นายประสิทธิ์ มั่นคง", amount: 150000, type: "สามัญ", status: "active", dueDate: "2024-06-15" },
  { id: "LN002", memberName: "นางมาลี สุขสันต์", amount: 50000, type: "ฉุกเฉิน", status: "active", dueDate: "2024-02-28" },
  { id: "LN003", memberName: "นายอนันต์ รุ่งเรือง", amount: 200000, type: "พิเศษ", status: "overdue", dueDate: "2024-01-10" },
];

const recentTransactions = [
  { id: "TXN001", type: "deposit", description: "ฝากเงินออมทรัพย์ - นายสมชาย ใจดี", amount: 5000, date: "2024-01-15 10:30" },
  { id: "TXN002", type: "share", description: "ซื้อหุ้นเพิ่ม - นางสาวสมหญิง รักษ์ดี", amount: 10000, date: "2024-01-15 09:45" },
  { id: "TXN003", type: "loan_payment", description: "ชำระหนี้เงินกู้ - นายวิชัย พัฒนา", amount: 8500, date: "2024-01-15 09:15" },
  { id: "TXN004", type: "sale", description: "ขายปุ๋ยเคมี 10 กระสอบ", amount: 12000, date: "2024-01-14 16:20" },
  { id: "TXN005", type: "withdrawal", description: "ถอนเงินออมทรัพย์ - นายประสิทธิ์ มั่นคง", amount: -3000, date: "2024-01-14 14:00" },
];

const alerts = [
  { id: 1, type: "warning", message: "มีสินเชื่อค้างชำระ 3 รายการ", link: "/loans?status=overdue" },
  { id: 2, type: "info", message: "สินค้าใกล้หมด: ปุ๋ยยูเรีย (เหลือ 15 กระสอบ)", link: "/products?filter=low_stock" },
  { id: 3, type: "success", message: "รายได้เดือนนี้เพิ่มขึ้น 12% จากเดือนก่อน", link: "/reports/revenue" },
];

// Simple bar chart component
function SimpleBarChart({ data, title }: { data: { label: string; value: number; color: string }[]; title: string }) {
  const maxValue = Math.max(...data.map(d => d.value));

  return (
    <div className="space-y-3">
      {title && <h4 className="font-medium text-slate-700">{title}</h4>}
      <div className="space-y-2">
        {data.map((item, index) => (
          <div key={index} className="space-y-1">
            <div className="flex justify-between text-sm">
              <span className="text-slate-600">{item.label}</span>
              <span className="font-medium">{formatCurrency(item.value)}</span>
            </div>
            <div className="h-3 bg-slate-100 rounded-full overflow-hidden">
              <div
                className={`h-full rounded-full transition-all duration-500 ${item.color}`}
                style={{ width: `${(item.value / maxValue) * 100}%` }}
              />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// Donut chart component (CSS-based)
function DonutChart({ data, total, centerLabel }: {
  data: { label: string; value: number; color: string }[];
  total: number;
  centerLabel: string;
}) {
  let cumulativePercent = 0;

  const segments = data.map(item => {
    const percent = (item.value / total) * 100;
    const segment = {
      ...item,
      percent,
      offset: cumulativePercent,
    };
    cumulativePercent += percent;
    return segment;
  });

  return (
    <div className="flex items-center gap-6">
      <div className="relative w-32 h-32">
        <svg viewBox="0 0 36 36" className="w-full h-full -rotate-90">
          {segments.map((segment, index) => (
            <circle
              key={index}
              cx="18"
              cy="18"
              r="15.9155"
              fill="transparent"
              stroke={segment.color}
              strokeWidth="3"
              strokeDasharray={`${segment.percent} ${100 - segment.percent}`}
              strokeDashoffset={`${-segment.offset}`}
              className="transition-all duration-500"
            />
          ))}
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className="text-lg font-bold text-slate-800">{formatNumber(total)}</span>
          <span className="text-xs text-slate-500">{centerLabel}</span>
        </div>
      </div>
      <div className="space-y-2">
        {data.map((item, index) => (
          <div key={index} className="flex items-center gap-2 text-sm">
            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: item.color }} />
            <span className="text-slate-600">{item.label}</span>
            <span className="font-medium">{formatNumber(item.value)}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export default function DashboardPage() {
  const [selectedPeriod, setSelectedPeriod] = useState<"day" | "week" | "month" | "year">("month");

  const loanStatusData = [
    { label: "ปกติ", value: 145, color: "#22c55e" },
    { label: "ค้างชำระ", value: 3, color: "#ef4444" },
    { label: "ปิดบัญชี", value: 52, color: "#94a3b8" },
  ];

  const revenueData = [
    { label: "ดอกเบี้ยเงินกู้", value: 1234560, color: "bg-blue-500" },
    { label: "ค่าธรรมเนียม", value: 234560, color: "bg-green-500" },
    { label: "ขายสินค้า", value: 876550, color: "bg-amber-500" },
  ];

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
      case "share":
        return (
          <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
        );
      case "loan_payment":
        return (
          <div className="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        );
      case "sale":
        return (
          <div className="w-8 h-8 rounded-full bg-amber-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
          </div>
        );
      default:
        return (
          <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center">
            <svg className="w-4 h-4 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        );
    }
  };

  const getAlertIcon = (type: string) => {
    switch (type) {
      case "warning":
        return (
          <div className="w-8 h-8 rounded-full bg-amber-100 flex items-center justify-center flex-shrink-0">
            <svg className="w-4 h-4 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
        );
      case "info":
        return (
          <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
            <svg className="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        );
      case "success":
        return (
          <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center flex-shrink-0">
            <svg className="w-4 h-4 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
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
        <div>
          <h1 className="text-2xl font-bold text-slate-800">แดชบอร์ด</h1>
          <p className="text-slate-500 mt-1">ภาพรวมการดำเนินงานสหกรณ์</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex bg-slate-100 rounded-lg p-1">
            {(["day", "week", "month", "year"] as const).map((period) => (
              <button
                key={period}
                onClick={() => setSelectedPeriod(period)}
                className={`px-3 py-1.5 text-sm font-medium rounded-md transition-all ${
                  selectedPeriod === period
                    ? "bg-white text-blue-600 shadow-sm"
                    : "text-slate-600 hover:text-slate-900"
                }`}
              >
                {period === "day" && "วัน"}
                {period === "week" && "สัปดาห์"}
                {period === "month" && "เดือน"}
                {period === "year" && "ปี"}
              </button>
            ))}
          </div>
          <Button variant="outline" size="sm">
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            ดาวน์โหลดรายงาน
          </Button>
        </div>
      </div>

      {/* Alerts */}
      {alerts.length > 0 && (
        <div className="space-y-2">
          {alerts.map((alert) => (
            <Link
              key={alert.id}
              href={alert.link}
              className={`flex items-center gap-3 p-3 rounded-lg border transition-all hover:shadow-sm ${
                alert.type === "warning"
                  ? "bg-amber-50 border-amber-200 hover:border-amber-300"
                  : alert.type === "info"
                  ? "bg-blue-50 border-blue-200 hover:border-blue-300"
                  : "bg-green-50 border-green-200 hover:border-green-300"
              }`}
            >
              {getAlertIcon(alert.type)}
              <span className="text-sm font-medium text-slate-700 flex-1">{alert.message}</span>
              <svg className="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>
          ))}
        </div>
      )}

      {/* Main Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="สมาชิกทั้งหมด"
          value={formatNumber(summaryData.totalMembers)}
          subtitle={`ใช้งาน ${formatNumber(summaryData.activeMembers)} ราย`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          }
          trend={{ value: 1.9, isPositive: true }}
        />
        <StatsCard
          title="ทุนเรือนหุ้น"
          value={formatCurrency(summaryData.totalShares)}
          subtitle="มูลค่าหุ้นทั้งหมด"
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          }
          trend={{ value: 3.2, isPositive: true }}
        />
        <StatsCard
          title="เงินฝากรวม"
          value={formatCurrency(summaryData.totalDeposits)}
          subtitle="ยอดเงินฝากทั้งหมด"
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
          }
          trend={{ value: 5.1, isPositive: true }}
        />
        <StatsCard
          title="สินเชื่อคงค้าง"
          value={formatCurrency(summaryData.totalLoans)}
          subtitle={`ค้างชำระ ${summaryData.overdueLoans} รายการ`}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
          trend={{ value: 2.3, isPositive: false }}
        />
      </div>

      {/* Financial Summary */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <Card className="bg-gradient-to-br from-green-500 to-green-600 text-white border-0">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-green-100 text-sm">รายได้เดือนนี้</p>
                <p className="text-2xl font-bold mt-1">{formatCurrency(summaryData.monthlyRevenue)}</p>
              </div>
              <div className="w-12 h-12 rounded-full bg-white/20 flex items-center justify-center">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-red-500 to-red-600 text-white border-0">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-red-100 text-sm">ค่าใช้จ่ายเดือนนี้</p>
                <p className="text-2xl font-bold mt-1">{formatCurrency(summaryData.monthlyExpense)}</p>
              </div>
              <div className="w-12 h-12 rounded-full bg-white/20 flex items-center justify-center">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6" />
                </svg>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-blue-500 to-blue-600 text-white border-0">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-blue-100 text-sm">กำไรสุทธิ</p>
                <p className="text-2xl font-bold mt-1">{formatCurrency(summaryData.netProfit)}</p>
              </div>
              <div className="w-12 h-12 rounded-full bg-white/20 flex items-center justify-center">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Revenue Chart */}
        <Card>
          <CardHeader>
            <CardTitle>รายได้ตามประเภท</CardTitle>
          </CardHeader>
          <CardContent>
            <SimpleBarChart data={revenueData} title="" />
          </CardContent>
        </Card>

        {/* Loan Status */}
        <Card>
          <CardHeader>
            <CardTitle>สถานะสินเชื่อ</CardTitle>
          </CardHeader>
          <CardContent>
            <DonutChart data={loanStatusData} total={200} centerLabel="สัญญา" />
          </CardContent>
        </Card>
      </div>

      {/* Recent Activities Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Transactions */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>รายการล่าสุด</CardTitle>
            <Link href="/transactions" className="text-sm text-blue-600 hover:text-blue-700 font-medium">
              ดูทั้งหมด
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentTransactions.map((txn) => (
                <div key={txn.id} className="flex items-center gap-3">
                  {getTransactionIcon(txn.type)}
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-slate-800 truncate">{txn.description}</p>
                    <p className="text-xs text-slate-500">{txn.date}</p>
                  </div>
                  <span className={`text-sm font-semibold ${txn.amount >= 0 ? "text-green-600" : "text-red-600"}`}>
                    {txn.amount >= 0 ? "+" : ""}{formatCurrency(txn.amount)}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Loans */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>สินเชื่อล่าสุด</CardTitle>
            <Link href="/loans" className="text-sm text-blue-600 hover:text-blue-700 font-medium">
              ดูทั้งหมด
            </Link>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>สมาชิก</TableHead>
                  <TableHead>ประเภท</TableHead>
                  <TableHead className="text-right">จำนวน</TableHead>
                  <TableHead>สถานะ</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recentLoans.map((loan) => (
                  <TableRow key={loan.id}>
                    <TableCell className="font-medium">{loan.memberName}</TableCell>
                    <TableCell>{loan.type}</TableCell>
                    <TableCell className="text-right">{formatCurrency(loan.amount)}</TableCell>
                    <TableCell>
                      <Badge variant={loan.status === "overdue" ? "error" : "success"}>
                        {loan.status === "active" ? "ปกติ" : loan.status === "overdue" ? "ค้างชำระ" : loan.status}
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      {/* New Members */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>สมาชิกใหม่</CardTitle>
          <Link href="/members" className="text-sm text-blue-600 hover:text-blue-700 font-medium">
            ดูทั้งหมด
          </Link>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>รหัสสมาชิก</TableHead>
                <TableHead>ชื่อ-นามสกุล</TableHead>
                <TableHead>วันที่สมัคร</TableHead>
                <TableHead className="text-right">ทุนเรือนหุ้น</TableHead>
                <TableHead>สถานะ</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {recentMembers.map((member) => (
                <TableRow key={member.id}>
                  <TableCell className="font-medium">{member.id}</TableCell>
                  <TableCell>{member.name}</TableCell>
                  <TableCell>{formatDate(member.joinDate)}</TableCell>
                  <TableCell className="text-right">{formatCurrency(member.shares)}</TableCell>
                  <TableCell>
                    <Badge variant="success">ใช้งาน</Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>ทางลัดการทำงาน</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            <Link
              href="/members?action=new"
              className="flex flex-col items-center gap-2 p-4 rounded-xl bg-blue-50 hover:bg-blue-100 transition-colors"
            >
              <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center">
                <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                </svg>
              </div>
              <span className="text-sm font-medium text-slate-700">เพิ่มสมาชิก</span>
            </Link>
            <Link
              href="/deposits?action=deposit"
              className="flex flex-col items-center gap-2 p-4 rounded-xl bg-green-50 hover:bg-green-100 transition-colors"
            >
              <div className="w-10 h-10 rounded-full bg-green-500 flex items-center justify-center">
                <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
              </div>
              <span className="text-sm font-medium text-slate-700">รับฝากเงิน</span>
            </Link>
            <Link
              href="/loans?action=new"
              className="flex flex-col items-center gap-2 p-4 rounded-xl bg-purple-50 hover:bg-purple-100 transition-colors"
            >
              <div className="w-10 h-10 rounded-full bg-purple-500 flex items-center justify-center">
                <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <span className="text-sm font-medium text-slate-700">สร้างสินเชื่อ</span>
            </Link>
            <Link
              href="/shares?action=buy"
              className="flex flex-col items-center gap-2 p-4 rounded-xl bg-amber-50 hover:bg-amber-100 transition-colors"
            >
              <div className="w-10 h-10 rounded-full bg-amber-500 flex items-center justify-center">
                <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <span className="text-sm font-medium text-slate-700">ซื้อหุ้น</span>
            </Link>
            <Link
              href="/products?action=sale"
              className="flex flex-col items-center gap-2 p-4 rounded-xl bg-rose-50 hover:bg-rose-100 transition-colors"
            >
              <div className="w-10 h-10 rounded-full bg-rose-500 flex items-center justify-center">
                <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
                </svg>
              </div>
              <span className="text-sm font-medium text-slate-700">ขายสินค้า</span>
            </Link>
            <Link
              href="/settings"
              className="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 hover:bg-slate-100 transition-colors"
            >
              <div className="w-10 h-10 rounded-full bg-slate-500 flex items-center justify-center">
                <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
              <span className="text-sm font-medium text-slate-700">ตั้งค่า</span>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
