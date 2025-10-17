# Quickstart – Movie Streaming Portal

## Prerequisites
- Bun 1.1+ (หรือ Node.js 20+ หากจำเป็นต้องใช้สคริปต์ pnpm เดิม)
- Go 1.23+
- Docker / Compose (PostgreSQL, Redis, CDN stub)
- OpenSSL (signed URL key generation)

## 1. Bootstrap infrastructure
```bash
docker compose up -d postgres redis
```
- PostgreSQL 16: ฐานข้อมูลหลัก (user/password/db = `leakstream`) เปิดพอร์ต 5432
- Redis: reserved for rate-limit counters and signed token cache (service ใช้ภาพ `redis:7.4-alpine` เปิดพอร์ต 6379)

## 2. Configure environment variables

1. คัดลอก `.env.example` ไปเป็น `.env` สำหรับ backend และปรับค่าเมื่อจำเป็น
2. คัดลอก `frontend/.env.local.example` ไปเป็น `frontend/.env.local` เพื่อให้ Next.js รู้ปลายทาง API (`NEXT_PUBLIC_API_BASE_URL`)

## 3. Install dependencies & generate types
```bash
bun install
# TODO: contracts:generate / go generate จะผูกกับ contract tooling ใน phase ถัดไป
go mod tidy
```

## 4. Seed development data
```bash
go run ./backend/cmd/migrate up
```
- Migration `0002_seed_sample.sql` จะเติมข้อมูลหนังตัวอย่างลง PostgreSQL ให้พร้อมทดสอบทันที

## 5. Run services
```bash
bun dev                 # Next.js 15 App Router dev server
go run ./backend/cmd/api
```
- Ensure `.env.local` contains CDN signing keys and Redis DSN.

## 6. Validate workflows
1. Visit `http://localhost:3000/movies/{sample-slug}` – confirm loading UI streams immediately, then playback starts ≤4 s.
2. Toggle captions on/off; verify WebVTT loads.
3. Log in as admin, create a movie with new caption asset, ensure duplicate title check blocks collisions.
4. Observe Redis rate-limit keys and OpenTelemetry traces for playback + admin actions.

## 7. Test suites
```bash
bun run test                                        # Vitest unit (frontend)
PLAYWRIGHT_BASE_URL=http://localhost:3000 npx playwright test  # Viewer E2E
go test ./...                                      # Backend unit/integration
# k6 run tests/load/viewer-stream.js               # Optional load (เตรียมไว้ทีหลัง)
```
- Ensure CI pipeline runs lint, tests, build, docker image publish, and uploads Playwright artifacts for regression debugging.

## 8. Deployment readiness
- Review monitoring dashboards: playback failure <5%, rate-limit alerts <50 events/5 min.
- Confirm feature flag `ADMIN_CATALOG_ENABLED` defaults off until launch sign-off.
- Prepare rollback by snapshotting PostgreSQL schema and caching config.
