# leak-streaming Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-10-17

## Active Technologies
- TypeScript (Next.js 15 App Router), Go 1.23+ + Next.js 15 (App Router, React Server Components, Suspense streaming), Tailwind CSS, Radix UI, Go chi + sqlc, MySQL driver, OpenTelemetry, Redis (token cache) (001-specify-scripts-bash)

## Project Structure
```
backend/
frontend/
tests/
```

## Commands
npm test && npm run lint

## Code Style
TypeScript (Next.js 15 App Router), Go 1.23+: Follow standard conventions

## Recent Changes
- 001-specify-scripts-bash: Added TypeScript (Next.js 15 App Router), Go 1.23+ + Next.js 15 (App Router, React Server Components, Suspense streaming), Tailwind CSS, Radix UI, Go chi + sqlc, MySQL driver, OpenTelemetry, Redis (token cache)

<!-- MANUAL ADDITIONS START -->
- สรุปสถานะ (อัปเดตโดย Codex ณ 2025-10-17):
  - ติดตั้งและตั้งค่า Next.js 15 (App Router) + Bun + Vitest/Playwright ฝั่ง frontend และ Go + chi + sqlc + Redis token cache ฝั่ง backend ครบ พร้อม lint/test script (`bun run lint`, `bun run test`, `GOCACHE=... go test ./...`).
  - ระบบดูหนังทำงานผ่าน flow เต็ม: frontend เรียก backend เพื่อขอ `playback-token`, ดึง manifest ที่ปรับปรุงแล้ว และ proxy segment พร้อมตรวจ allowed host; มีคำบรรยายตัวอย่างแบบไฟล์ภายใน (`frontend/public/captions/sample-en.vtt`).
  - Backend เชื่อมต่อ PostgreSQL (docker-compose) พร้อม goose migrations/seed สำหรับ demo movies 2 เรื่อง และรีโปซิทอรีอ่านเขียนข้อมูลจริงได้; Redis สำหรับ cache token พร้อมใช้งานใน docker.
  - ฝั่ง frontend มีหน้า `/movies` แสดงรายการจากฐานข้อมูล, หน้า detail สตรีมวิดีโอ, และ admin UI (`/admin/movies/new`) สำหรับเพิ่มหนังใหม่ (รวม server action + validation + textarea component).
  - เพิ่ม Playwright specs (viewer streaming, movie catalogue, admin create) แม้ยังต้องรันพร้อม backend/frontend/ฐานข้อมูล; lint, unit และ go test ผ่านในสภาพแวดล้อม dev.
  - งานที่เหลือ: เสริม integration tests ฝั่ง backend, เปิดใช้ Playwright ใน CI เมื่อสแตกพร้อม, ปรับปรุง security ของ allowed hosts (ยังไม่รองรับ wildcard `*` โดยตั้งใจ), จัดการเครื่องมือ admin เพิ่มเติม (เช่น list/delete) หากจำเป็น.
<!-- MANUAL ADDITIONS END -->
