## Excel to ClickHouse Cooperative

Blank Next.js 14 (App Router) workspace forสร้างระบบสหกรณ์ที่เชื่อมข้อมูลจาก Excel ไปยัง ClickHouse. รวม Tailwind CSS v4, TypeScript, ESLint, Prettier, และหน้าพื้นฐานสำหรับแดชบอร์ด สมาชิก สินเชื่อ และการตั้งค่า

### Scripts

```bash
npm run dev      # start the turbopack dev server
npm run build    # create a production build
npm run start    # run the production build locally
npm run lint     # run ESLint across the project
npm run format   # run Prettier in check mode
```

### Directory Highlights

- `app/` – App Router pages and layouts (dashboard, members, loans, settings)
- `components/nav-bar.tsx` – ท็อปบาร์นำทางพร้อมไฮไลต์เมนูปัจจุบัน
- `.prettierrc` – Formatting rules shared between Prettier and editors

### Next Steps

1. เชื่อมต่อ ClickHouse client และ service layer สำหรับ ingest/extract ข้อมูล
2. สร้างตาราง UI/ฟอร์มจริงสำหรับสมาชิกและสินเชื่อ พร้อม state management
3. เพิ่มฟีเจอร์การนำเข้าไฟล์ Excel และ mapping เข้าสู่ ClickHouse ใน background job

### Development Notes

- Tailwind 4 อยู่ในสถานะ alpha: ตรวจสอบ release notes หากมี breaking changes
- แก้ไขสไตล์กลางใน `app/globals.css` หรือใช้ utility class เพิ่มเติมตามต้องการ
- หน้าใหม่สามารถเพิ่มภายใต้ `app/<route>/page.tsx` และจะถูกเพิ่มใน NavBar ได้โดยแก้ `components/nav-bar.tsx`
