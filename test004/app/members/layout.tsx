import type { ReactNode } from "react";

export default function MembersLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <div className="space-y-4">
      <p className="text-xs uppercase tracking-wide text-slate-500">
        ระบบสหกรณ์ · จัดการสมาชิก
      </p>
      {children}
    </div>
  );
}
