# Quickstart – Movie Streaming Portal

## Prerequisites
- Node.js 20+, pnpm 9+
- Go 1.23+
- Docker / Compose (MySQL, Redis, CDN stub)
- OpenSSL (signed URL key generation)

## 1. Bootstrap infrastructure
```bash
docker compose up -d mysql redis nginx-cdn
```
- MySQL: apply migrations via `backend/internal/persistence/migrations`.
- Redis: reserved for rate-limit counters and signed token cache.
- Nginx CDN stub: ensure `X-Accel-Buffering: no` header is present (per Context7 Next.js streaming guidance).

## 2. Install dependencies & generate types
```bash
pnpm install
pnpm contracts:generate
go mod tidy
go generate ./...
```
- Contracts generator emits TypeScript clients in `frontend/lib/contracts` and Go clients in `backend/internal/contracts`.

## 3. Seed development data
```bash
pnpm db:seed            # adds sample movies, captions, streams
go run ./backend/cmd/admin bootstrap --force
```

## 4. Run services
```bash
pnpm dev                # Next.js 15 App Router dev server (uses Suspense/loading states)
go run ./backend/cmd/api
```
- Ensure `.env.local` contains CDN signing keys and Redis DSN.

## 5. Validate workflows
1. Visit `http://localhost:3000/movies/{sample-slug}` – confirm loading UI streams immediately, then playback starts ≤4 s.
2. Toggle captions on/off; verify WebVTT loads.
3. Log in as admin, create a movie with new caption asset, ensure duplicate title check blocks collisions.
4. Observe Redis rate-limit keys and OpenTelemetry traces for playback + admin actions.

## 6. Test suites
```bash
pnpm test               # Jest unit
pnpm test:e2e           # Playwright (viewer + admin journeys)
go test ./...           # Backend unit/integration
k6 run tests/load.js    # Optional smoke for 1k concurrent viewers
```
- Ensure CI pipeline runs lint, tests, build, docker image publish, and uploads Playwright artifacts for regression debugging.

## 7. Deployment readiness
- Review monitoring dashboards: playback failure <5%, rate-limit alerts <50 events/5 min.
- Confirm feature flag `ADMIN_CATALOG_ENABLED` defaults off until launch sign-off.
- Prepare rollback by snapshotting MySQL schema and caching config.
