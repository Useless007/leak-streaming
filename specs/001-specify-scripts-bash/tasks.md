---
description: "Task list template for feature implementation"
---

# Tasks: Movie Streaming Portal

**Input**: Design documents from `/specs/001-specify-scripts-bash/`  
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/  
**Tests**: Contract, integration, end-to-end, และ load tests เป็นข้อบังคับสำหรับทุกส่วนที่แตะ shared API contracts หรือ performance criteria; unit tests ต้องครอบคลุมบริการ Go และ client component เสมอ  
**Organization**: แบ่งงานตาม User Story เพื่อให้ส่งมอบเป็น slice ที่ทดสอบและปล่อยได้อย่างอิสระ

## Format: `[ID] [P?] [Story] Description`
- **[P]**: ทำงานขนานได้ (ไฟล์/ดีเพนเดนซีไม่ชนกัน)
- **[Story]**: ระบุ US ที่งานนั้นรองรับ เช่น `[US1]`
- คำอธิบายต้องมี path ไฟล์ชัดเจน

## Path Conventions
- Frontend: `frontend/app/`, `frontend/components/`, `frontend/lib/`, `frontend/styles/`, `frontend/tests/`
- Backend: `backend/cmd/`, `backend/internal/{api,service,domain,persistence,platform}`, `backend/tests/`
- Contracts: `specs/001-specify-scripts-bash/contracts/`
- Infrastructure: `infrastructure/{k8s,terraform,pipelines}/`
- Scripts/tooling: `scripts/`, `.github/workflows/`

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 ติดตั้ง shadcn/ui CLI และ initialize registry (`npx shadcn@latest init`) ใน `frontend/`
- [ ] T002 เพิ่ม base components (button, input, form, card, dialog, navigation) ผ่าน shadcn CLI
- [ ] T003 ติดตั้ง dependencies (`pnpm install`, `go mod tidy`) ตาม quickstart
- [ ] T004 สร้าง shared contract clients (`pnpm contracts:generate`, `go generate ./...`)
- [ ] T005 เปิด docker services (MySQL, Redis, nginx-cdn) ด้วย `docker compose` ที่ repo root
- [ ] T006 seed ข้อมูลตัวอย่างและบัญชี admin (`pnpm db:seed`, `go run ./backend/cmd/admin bootstrap`)

## Phase 2: Foundational (Blocking Prerequisites)

- [ ] T007 ใช้งาน database migrations สำหรับ movies/streams/captions/tokens ใน `backend/internal/persistence/migrations`
- [ ] T008 สร้าง Redis client/config ใน `backend/internal/platform/cache/redis.go`
- [ ] T009 สร้าง rate-limit middleware skeleton ใน `backend/internal/api/middleware/ratelimit.go`
- [ ] T010 จัดทำ Zod schema + API wrapper ที่ `frontend/lib/api/` จากโค้ดที่ generate
- [ ] T011 ตั้งค่า ThemeProvider + Radix primitives ใน `frontend/app/layout.tsx` (รองรับ dark mode)
- [ ] T012 ติดตั้ง OpenTelemetry exporters และ correlation ID middleware ใน `backend/internal/platform/telemetry`
- [ ] T013 ขยาย workflow CI (`.github/workflows/ci.yml`) ให้ครอบคลุม lint/test/build/e2e/load
- [ ] T014 เพิ่ม Grafana/Prometheus dashboard definition ใน `infrastructure/pipelines/observability/`
- [ ] T015 กำหนด canary deployment pipeline และ progressive rollout script ใน `.github/workflows/deploy.yml`
- [ ] T016 เพิ่ม health/readiness probes และ runtime checks ใน `infrastructure/k8s/` สำหรับ backend/frontend

## Phase 3: User Story 1 - Stream a Published Movie (Priority: P1) 🎯 MVP

**Goal**: ผู้ชมเปิดหน้า `app/movies/[movieId]/page.tsx` แล้วเริ่มเล่นภายใน 4 วินาที  
**Independent Test**: Playwright เปิดหน้าหนัง → เริ่มเล่น, toggle คำบรรยาย, จำลอง error retry สำเร็จ

### Tests
- [ ] T017 [P] [US1] เพิ่ม Playwright spec `frontend/tests/e2e/viewer-stream.spec.ts` ครอบคลุม playback + retry
- [ ] T018 [P] [US1] เพิ่ม Go integration test `backend/tests/integration/movies_stream_test.go` (signed token + rate limit)
- [ ] T019 [P] [US1] สร้าง k6 load test `tests/load/viewer-stream.js` ทดสอบ ≥1,000 concurrent viewers (SC-006)

### Implementation
- [ ] T020 [US1] พัฒนา server component `frontend/app/movies/[movieId]/page.tsx` พร้อม Suspense/loading/error
- [ ] T021 [P] [US1] สร้าง `frontend/app/movies/[movieId]/metadata.ts` ไล่เติมเมตาดาทาแบบ server-side
- [ ] T022 [P] [US1] ทำ `frontend/components/movie/player.tsx` (caption toggle, retry แสดง error)
- [ ] T023 [P] [US1] เพิ่ม hook `frontend/lib/hooks/usePlayback.ts` สำหรับสถานะผู้เล่นและ telemetry
- [ ] T024 [US1] สร้าง handler `backend/internal/api/movies/stream_token.go` คืน signed URL พร้อม context deadline
- [ ] T025 [P] [US1] พัฒนาบริการ `backend/internal/service/movies/token_service.go` (Redis cache + TTL)
- [ ] T026 [US1] ปรับ repository `backend/internal/persistence/movies/repository.go` ดึง stream source + captions
- [ ] T027 [P] [US1] ผูก rate-limit middleware กับ playback route ใน `backend/internal/api/router.go`
- [ ] T028 [US1] เชื่อม observability (log/trace/metric) สำหรับ playback path ใน `backend/internal/platform/telemetry`
- [ ] T029 [P] [US1] ตกแต่งหน้า viewer ด้วย shadcn components ใน `frontend/app/movies/[movieId]/page.tsx`
- [ ] T030 [US1] อัปเดต quickstart ส่วนผู้ชมใน `specs/001-specify-scripts-bash/quickstart.md`

## Phase 4: User Story 2 - Add a New Movie to the Catalog (Priority: P2)

**Goal**: ผู้จัดการเนื้อหาเพิ่มภาพยนตร์ใหม่ (schedule, stream, caption) ผ่าน admin UI  
**Independent Test**: Playwright ล็อกอิน → กรอกฟอร์ม → บันทึกสำเร็จ → เห็นแสดงในรายการ

### Tests
- [ ] T031 [P] [US2] สร้าง Playwright spec `frontend/tests/e2e/admin-create.spec.ts` สำหรับการเพิ่มภาพยนตร์
- [ ] T032 [P] [US2] เพิ่ม Go integration test `backend/tests/integration/movies_create_test.go` (validation + duplicate title)
- [ ] T033 [P] [US2] เพิ่มการทดสอบ rate limit สำหรับการสร้าง (`backend/tests/integration/movies_create_ratelimit_test.go`)

### Implementation
- [ ] T034 [US2] สร้าง layout ฝั่ง admin พร้อม Sidebar/Card ใน `frontend/app/admin/layout.tsx`
- [ ] T035 [P] [US2] ทำฟอร์ม `frontend/app/admin/movies/new/page.tsx` ด้วย shadcn form + zod resolver
- [ ] T036 [P] [US2] สร้าง server action `frontend/app/admin/movies/new/actions.ts` (optimistic validation)
- [ ] T037 [US2] เพิ่ม mutation client `frontend/lib/api/movies/createMovie.ts`
- [ ] T038 [US2] เขียน handler สร้างภาพยนตร์ที่ `backend/internal/api/movies/create.go`
- [ ] T039 [P] [US2] ขยาย service `backend/internal/service/movies/mutation_service.go` (availability window rules)
- [ ] T040 [P] [US2] เพิ่ม sqlc statement `backend/internal/persistence/movies/create_movie.sql`
- [ ] T041 [US2] สร้างโมดูล `backend/internal/service/movies/caption_validator.go`
- [ ] T042 [P] [US2] ผูก rate-limit middleware กับเส้นทางสร้างภาพยนตร์ใน `backend/internal/api/router.go`
- [ ] T043 [P] [US2] ใส่ shadcn toast แสดง success/error ใน `frontend/app/admin/movies/new/page.tsx`
- [ ] T044 [US2] อัปเดต quickstart ส่วน admin create ใน `specs/001-specify-scripts-bash/quickstart.md`

## Phase 5: User Story 3 - Manage Upcoming Titles (Priority: P3)

**Goal**: ผู้จัดการปรับแก้ metadata, schedule และ visibility ได้เอง  
**Independent Test**: Playwright แก้ schedule → toggle visibility → ตรวจว่าหน้า viewer อัปเดตทันที

### Tests
- [ ] T045 [P] [US3] เพิ่ม Playwright spec `frontend/tests/e2e/admin-manage.spec.ts`
- [ ] T046 [P] [US3] เพิ่ม Go integration test `backend/tests/integration/movies_update_test.go`
- [ ] T047 [P] [US3] เพิ่มการทดสอบ rate limit สำหรับ update/visibility (`backend/tests/integration/movies_manage_ratelimit_test.go`)

### Implementation
- [ ] T048 [US3] สร้างหน้า `frontend/app/admin/movies/[movieId]/page.tsx` (shadcn tabs)
- [ ] T049 [P] [US3] ทำ component ฟอร์มแก้ไขใน `frontend/app/admin/movies/[movieId]/_components/edit-forms.tsx`
- [ ] T050 [P] [US3] สร้าง server actions `frontend/app/admin/movies/[movieId]/actions.ts` (update + visibility toggle)
- [ ] T051 [US3] เพิ่ม client `frontend/lib/api/movies/updateMovie.ts` และ `toggleVisibility.ts`
- [ ] T052 [US3] ปรับ handler อัปเดตที่ `backend/internal/api/movies/update.go`
- [ ] T053 [P] [US3] ทำ handler visibility `backend/internal/api/movies/visibility.go` พร้อม audit log
- [ ] T054 [P] [US3] เพิ่ม sqlc statement สำหรับ update/visibility ใน `backend/internal/persistence/movies/update_movie.sql`
- [ ] T055 [US3] สร้าง audit logger ใน `backend/internal/service/movies/audit_logger.go`
- [ ] T056 [P] [US3] รีเฟรช list หน้า viewer ใน `frontend/app/movies/page.tsx` หลัง visibility เปลี่ยน
- [ ] T057 [US3] อัปเดต quickstart ส่วน maintenance ใน `specs/001-specify-scripts-bash/quickstart.md`
- [ ] T058 [P] [US3] ผูก rate-limit middleware กับเส้นทาง update/visibility ใน `backend/internal/api/router.go`

## Phase N: Polish & Cross-Cutting Concerns

- [ ] T059 รัน Lighthouse/Performance audit ของหน้า viewer และเก็บรายงานใน `docs/perf/`
- [ ] T060 เสริม error boundary/fallback UI ด้วย shadcn alerts ทั้งระบบ
- [ ] T061 ตรวจ accessibility (caption default, focus state) ด้วย axe (`frontend/tests/accessibility.spec.ts`)
- [ ] T062 ปรับ threshold rate limit และบันทึกขั้นตอน override ใน `docs/operations/rate-limits.md`
- [ ] T063 อัปเดต changelog ใน `docs/changelog.md` สรุปการเปิดตัว portal
- [ ] T064 ตรวจสอบว่า alert (Grafana) ทำงานเมื่อจำลองเหตุการณ์ล้มเหลว (`scripts/ci/test-alerts.sh`)

## Dependencies & Execution Order

- Phase 1 ต้องเสร็จจึงเริ่ม Phase 2 ได้
- Phase 2 ต้องครบก่อนเข้าสู่ User Stories
- US1 (P1) คือ MVP และต้องเสร็จก่อน US2/US3 เพื่อพิสูจน์ end-to-end streaming
- US2 และ US3 ทำขนานกันได้หลัง US1 หากทีมพร้อม
- Phase N ทำหลังปิดทุก story ที่ต้องการ

## Parallel Execution Examples

- US1: หลัง T020 เสร็จ สามารถรัน T017, T018, T019, T022 ขนานกันได้
- US2: T035 (ฟอร์ม) และ T038 (handler) ทำขนานหลัง foundational พร้อม, T043 (UI toast) ทำคู่กับ T036
- US3: T049 (form) และ T053 (visibility endpoint) รันขนานหลัง API client พร้อม; T047 ทดสอบ rate limit ควบคู่กับ T058

## Implementation Strategy

1. MVP = Phase 1–2 + US1 (T001–T030) เพื่อปล่อยฟังก์ชันสตรีมพื้นฐาน
2. Increment ถัดไป: US2 (T031–T044) เปิดให้ทีมคอนเทนต์เพิ่มภาพยนตร์
3. Increment ต่อไป: US3 (T045–T058) สำหรับการดูแลและแก้ไขข้อมูล
4. ปิดท้ายด้วย Phase N (T059–T064) ทำ performance, accessibility, alert และการปรับแต่งขั้นสุดท้าย
