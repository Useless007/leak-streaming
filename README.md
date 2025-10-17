# leak-streaming (Educational Project)

โครงการตัวอย่างสำหรับสาธิตระบบสตรีมมิ่งวิดีโอแบบ HLS ที่มีทั้ง Backend (Go) และ Frontend (Next.js 15 App Router) เน้นสถาปัตยกรรมที่อ่านง่าย ทดสอบได้ และใช้งานจริงได้ในสภาพแวดล้อมพัฒนา เหมาะสำหรับการศึกษาและทดลองแนวทางการทำงานแบบ end-to-end

> หมายเหตุ: โปรเจกต์นี้มีข้อมูลและลิงก์ตัวอย่างเพื่อการศึกษาเท่านั้น ไม่ควรใช้กับสื่อที่มีลิขสิทธิ์จริง

## คุณสมบัติเด่น
- แคตตาล็อกภาพยนตร์ + หน้าแสดงรายละเอียดและเล่นวิดีโอ (HLS)
- ระบบขอ playback token จาก backend แล้วปรับแต่ง manifest + proxy segment พร้อมตรวจสอบ allowed hosts
- คำบรรยาย (WebVTT) ตัวอย่าง พร้อม UI สลับเปิด/ปิด
- Admin UI (แบบง่าย) สำหรับเพิ่มภาพยนตร์ใหม่ พร้อม validation ทั้งฝั่ง client และ server
- มี integration rate limit, Redis token cache, และ OpenTelemetry hook ไว้ต่อยอดการสังเกตการณ์

## โครงสร้างโปรเจกต์ (ย่อ)
```
backend/        # Go API, migrations, services
frontend/       # Next.js 15 App Router UI
infrastructure/ # manifests/pipelines ตัวอย่าง
specs/          # บันทึกแผน/สเปก ใช้ประกอบการศึกษา
```

## เทคสแตก
- Frontend: Next.js 15 (App Router, React 18), Tailwind CSS, Radix UI, Vitest, Playwright (E2E)
- Backend: Go 1.23+, chi, goose (migrations), pgx (PostgreSQL), Redis (token cache), OpenTelemetry
- Database: PostgreSQL 16 (Docker)
- Cache: Redis 7 (Docker)

## ข้อกำหนดเบื้องต้น
- Docker และ Docker Compose
- Go 1.23+
- Bun 1.1+ หรือ Node.js 20+ (ใช้รันสคริปต์ฝั่ง frontend)

## การติดตั้งและรัน (เริ่มตั้งแต่ clone จนใช้งานได้)

### 1) Clone โค้ดและติดตั้ง dependencies
```bash
# Clone
git clone <your-fork-or-repo-url> leak-streaming
cd leak-streaming

# ติดตั้ง dependencies ฝั่ง frontend (ใช้ Bun ตามค่าเริ่มต้นของ workspace)
bun install

# จัดการโมดูล Go
cd backend && go mod tidy && cd -
```

### 2) บูทฐานข้อมูลและแคช
```bash
docker compose up -d postgres redis
```
- PostgreSQL: เชื่อมต่อที่ `postgres://leakstream:leakstream@127.0.0.1:5432/leakstream`
- Redis: เปิดที่ `127.0.0.1:6379`

### 3) ตั้งค่าตัวแปรสภาพแวดล้อม
- Backend ใช้ค่าเริ่มต้นจากโค้ดโดยตรง (ไม่จำเป็นต้องมีไฟล์ `.env` ในเบื้องต้น) โดยชี้ไปที่ DB/Redis บนเครื่อง
- Frontend ต้องกำหนดปลายทาง API ให้ถูกต้อง:

สร้างไฟล์ `frontend/.env.local` (ถ้ายังไม่มี)
```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### 4) รัน migrations และ seed ข้อมูลตัวอย่าง
คำสั่งนี้จะสร้างตารางและใส่ข้อมูล demo movie 2 เรื่อง พร้อมสตรีม+คำบรรยายตัวอย่าง
```bash
# รันจากโฟลเดอร์รากของโปรเจกต์
go run ./backend/cmd/migrate up
```
สิ่งที่จะถูกสร้างและ seed:
- ตาราง `movies`, `movie_streams`, `movie_captions`, `playback_tokens`
- ตัวอย่างภาพยนตร์:
  - slug: `sample-movie`
  - slug: `demo-movie-2`
- คำบรรยายภาษาอังกฤษไฟล์ตัวอย่าง: `frontend/public/captions/sample-en.vtt`

### 5) รัน Backend และ Frontend
เปิด 2 เทอร์มินัลแยกกัน:
```bash
# เทอร์มินัลที่ 1 – Backend API (พอร์ต 8080)
go run ./backend/cmd/api

# เทอร์มินัลที่ 2 – Frontend (Next.js dev server, พอร์ต 3000)
cd frontend
bun run dev
```

จากนั้นเปิดเบราว์เซอร์ไปที่:
- หน้าแรก/รายการหนัง: `http://localhost:3000/movies`
- หน้ารายละเอียดและเล่น: `http://localhost:3000/movies/sample-movie`
- Admin สร้างหนังใหม่: `http://localhost:3000/admin/movies/new`

หมายเหตุ: เมื่อสร้างหนังใหม่ ฟิลด์ที่ต้องกรอกให้ครบ ได้แก่ title, synopsis, posterUrl, streamUrl, allowedHosts, availabilityStart, availabilityEnd (slug ถูกสร้างอัตโนมัติจาก title)

## Running Tests
- Frontend unit tests (Vitest)
```bash
bun run test
```
- Backend unit/integration
```bash
go test ./...
```
- Playwright E2E (ต้องรัน backend+frontend ก่อน)
```bash
PLAYWRIGHT_BASE_URL=http://localhost:3000 npx playwright test
```

## ทิปส์และการแก้ปัญหา
- ถ้ารัน backend แล้วเชื่อม DB ไม่ได้ ให้ตรวจสอบ Docker Compose ว่า postgres ทำงานและพอร์ต 5432 ไม่ชน
- ห้องสมุด `NEXT_PUBLIC_API_BASE_URL` ต้องชี้ไปที่ backend API จริง มิฉะนั้น UI จะเรียก API ไม่เจอ
- หากต้องการรีเซ็ตฐานข้อมูล ทดลองลบ volume ที่ชื่อ `postgres-data` แล้วรัน migrations ใหม่

## ใบอนุญาตใช้งาน
โครงการนี้จัดทำเพื่อการศึกษาเท่านั้น โปรดตรวจสอบสิทธิ์การใช้งานของสื่อ/ลิงก์ที่คุณเพิ่มเข้าไปในระบบของคุณเอง
