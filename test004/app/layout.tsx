import type { Metadata } from "next";
import { Sarabun } from "next/font/google";
import type { ReactNode } from "react";
import { NavBar } from "@/components/nav-bar";
import "./globals.css";

const sarabun = Sarabun({
  weight: ['300', '400', '500', '600', '700'],
  subsets: ["thai", "latin"],
  variable: "--font-sarabun",
});

export const metadata: Metadata = {
  title: "ระบบบริหารจัดการสหกรณ์การเกษตร",
  description: "ระบบสารสนเทศเพื่อการบริหารจัดการสหกรณ์การเกษตรแห่งชาติ รองรับหลายสหกรณ์",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  return (
    <html lang="th">
      <body
        className={`${sarabun.variable} font-sans bg-gradient-to-br from-slate-50 to-blue-50/30 text-slate-900 antialiased flex min-h-screen flex-col`}
      >
        <NavBar />
        <main className="mx-auto flex w-full max-w-7xl flex-1 flex-col gap-6 px-6 py-8">
          {children}
        </main>
        <footer className="border-t border-slate-200 bg-white py-6">
          <div className="mx-auto max-w-7xl px-6 text-center text-sm text-slate-500">
            <p>ระบบบริหารจัดการสหกรณ์การเกษตรแห่งชาติ | National Agricultural Cooperative Management System</p>
            <p className="mt-1">พัฒนาตามมาตรฐานกรมตรวจบัญชีสหกรณ์และกรมส่งเสริมสหกรณ์</p>
          </div>
        </footer>
      </body>
    </html>
  );
}
