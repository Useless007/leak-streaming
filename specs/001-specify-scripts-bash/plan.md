# Implementation Plan: Movie Streaming Portal

**Branch**: `001-specify-scripts-bash` | **Date**: 2025-10-17 | **Spec**: [specs/001-specify-scripts-bash/spec.md](specs/001-specify-scripts-bash/spec.md)
**Input**: Feature specification from `/specs/001-specify-scripts-bash/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Deliver a streaming portal that lets visitors watch HLS (`.m3u8`) movies with signed session URLs while content managers curate the catalog (create, schedule, update, toggle visibility) via an admin interface. We will use Next.js 15 App Router for the viewer/admin experiences, Go services backed by MySQL for catalog and token issuance, and enforce accessibility (captions) plus rate limiting and observability guardrails to meet the constitution’s reliability and security expectations.

## Technical Context

**Language/Version**: TypeScript (Next.js 15 App Router), Go 1.23+  
**Primary Dependencies**: Next.js 15 (App Router, React Server Components, Suspense streaming), Tailwind CSS, Radix UI, Go chi + sqlc, MySQL driver, OpenTelemetry, Redis (token cache)  
**Storage**: MySQL (movie catalog, caption metadata) + Redis (signed URL cache / rate limit counters)  
**Testing**: Playwright (viewer/admin E2E), Jest/Testing Library (frontend unit), `go test` + testify + httptest (backend unit/integration), k6 smoke tests for streaming concurrency  
**Target Platform**: Web (SSR + streaming) deployed on Kubernetes with CDN-backed HLS delivery  
**Project Type**: B2C streaming web application with shared contracts across frontend and backend  
**Performance Goals**: 95% viewers start playback ≤4 s, sustain ≥1,000 concurrent viewers with <5% quality degradation, backend movie endpoints p95 ≤200 ms  
**Constraints**: Enforce signed `.m3u8` URLs, caption availability, App Router streaming (loading UI + Suspense per Context7 Next.js guidance), zero-trust admin APIs, rate limiting (120 rpm viewers / 20 rpm admins)  
**Scale/Scope**: Initial catalog of hundreds of titles, concurrency target 1k, single region deployment with CDN fan-out

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Principle I (Next.js App Router Discipline): Viewer and admin routes live in `frontend/app/`, default to server components, use Suspense/loading states to stream UI (per Context7 Next.js production checklist), and isolate `"use client"` playback controls with strict props typing.
- Principle II (Go Backend Reliability): Go services remain layered (`internal/api` → `internal/service` → `internal/persistence`), propagate `context.Context`, enforce request cancellation/timeouts, and guard signed URL issuance with concurrency-safe Redis cache.
- Principle III (Shared Contracts & Type Safety): Define `contracts/openapi/movies.yaml`, generate Go + TS types (sqlc + openapi generator), validate payloads via zod + Go structs, and version the contract (v1.0.0) with backward-compatible warnings.
- Principle IV (Observability & Security Assurance): Emit structured logs/traces (correlation IDs, movieId, actorId), monitor rate-limit alerts (OP-005), store `.m3u8` URLs encrypted, run SAST/DAST in CI, and gate releases on green telemetry dashboards.
- Principle V (Continuous Delivery Excellence): Enforce CI pipelines (lint, tests, build, docker image, Playwright), create preview envs for PRs, document rollback playbook, and keep feature flags for admin UI rollout.

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
```
backend/
├── cmd/                 # Service entrypoints
├── internal/
│   ├── api/             # Handlers, transport adapters
│   ├── domain/          # Business logic
│   ├── persistence/     # Repositories
│   └── platform/        # Observability, config, clients
├── pkg/                 # Shared libraries (exported)
└── tests/
    ├── integration/
    └── contract/

frontend/
├── app/                 # Next.js App Router entrypoints
├── components/          # Server/client components (clearly labeled)
├── lib/                 # Shared utilities (server-first)
├── styles/
└── tests/               # Jest + Playwright specs

contracts/
├── openapi/             # REST contracts
├── proto/               # gRPC contracts
└── generators/          # Scripts to emit TS + Go types

infrastructure/
├── k8s/
├── terraform/
└── pipelines/

scripts/
├── dev.sh
└── ci/

```

**Structure Decision**: Adopt the documented dual-front/back structure with `frontend/app` (App Router) and `backend/internal/...` service layering, plus `contracts/openapi` for shared schemas; extend `infrastructure/k8s` and `scripts/ci` to cover preview deployments and rate-limit observability pipelines.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |

## Constitution Check Status
- **Principle I**: PASS – Viewer/admin surfaces reside in App Router server components with route-level Suspense + loading UIs, aligning with Context7 Next.js streaming checklist; client components limited to playback controls.
- **Principle II**: PASS – Go handlers maintain layered architecture, enforce `context.Context` deadlines, and guard Redis-backed signed tokens to avoid resource leaks.
- **Principle III**: PASS – OpenAPI contract (`contracts/movies.yaml`) defined; generators kept in quickstart; zod + Go validation mirror contract schema.
- **Principle IV**: PASS – Structured logging/tracing, rate-limit alert dashboards (OP-005), encrypted storage of `.m3u8` URLs, and CI security scans codified.
- **Principle V**: PASS – CI pipeline covers lint/unit/e2e/load + preview deploys; feature flag controls admin rollout with rollback plan documented.

## Phase 0 – Research Actions
1. **Streaming UX** – Implement Suspense/loading states and disable proxy buffering (`X-Accel-Buffering: no`) so pages stream immediately (Context7 Next.js production checklist).
2. **Parallel Fetching & Caching** – Fetch movie metadata and signed tokens in parallel server actions and apply explicit cache directives (`fetch` caching or `unstable_cache`) to avoid waterfalls.
3. **Go Context Discipline** – Use `context.WithTimeout` and respect cancellation down to persistence for token issuance (Go stdlib context/http best practices).
4. **Rate Limiting Strategy** – Store viewer (IP/session) and admin (account) counters in Redis, throttle at 120/20 rpm respectively, surface telemetry to dashboards.
5. **Caption Compliance** – Require at least one caption track per movie, validate on ingest, and expose caption toggles in player UI to meet accessibility promise.

Deliverable: [research.md](specs/001-specify-scripts-bash/research.md).

## Phase 1 – Design & Contracts
### Data Model Highlights
- Entities: `Movie`, `StreamSource`, `CaptionTrack`, `ViewingSession`, `SignedStreamToken`, `RateLimitBucket`.
- Key rules: unique title per movie, mandatory captions, `.m3u8` validation, availability windows, signed token TTL ≤5 min.

### API Endpoints (OpenAPI v1.0.0)
- `GET /api/movies` (filters for visibility/schedule).
- `POST /api/movies`, `PATCH /api/movies/{movieId}`, `POST /api/movies/{movieId}/visibility`.
- `POST /api/movies/{movieId}/stream-token` for per-session signed URLs.
- Standardized responses for validation, rate-limit, caption-missing errors.

### Frontend Plan
- Server components for viewer/admin pages with `generateMetadata` to ensure head tags before streaming.
- Client `MoviePlayer` handles playback controls, caption toggles, error retries.
- Admin forms use server actions + zod validation, showing inline errors, preserving input.

### Backend Plan
- chi router with handlers per endpoint, services enforcing business logic and rate limits.
- sqlc-generated repositories for MySQL; `.m3u8` encryption helper.
- Redis token cache with request-scoped contexts; audit logging for mutations.

### Observability & Security
- OpenTelemetry middleware attaches correlation IDs, movieId, viewer/admin identifiers (where allowed).
- Rate-limit middleware logs throttle events; Grafana dashboards alert >50 events/5 min.
- Secrets via Vault/SSM; automated SAST/DAST gates prior to deployment.

Deliverables: [data-model.md](specs/001-specify-scripts-bash/data-model.md), [contracts/movies.yaml](specs/001-specify-scripts-bash/contracts/movies.yaml), [quickstart.md](specs/001-specify-scripts-bash/quickstart.md).

## Phase 2 – Engineering Plan
### Frontend
1. Scaffold viewer route with Suspense + streaming metadata.
2. Build client-side player controls with caption toggle and retry messaging.
3. Implement admin create/update pages with server actions and validation feedback.
4. Integrate generated TypeScript API client, enforce type-safe responses.
5. Add Playwright flows: viewer playback, admin CRUD, rate-limit error surfacing.

### Backend
1. Implement chi handlers + services honoring contexts/timeouts.
2. Create sqlc queries and migrations for movies, streams, captions.
3. Add Redis rate limiting and signed token issuance with audit logging.
4. Instrument OpenTelemetry traces/logs/metrics; cover rate-limit + playback KPIs.
5. Write unit/integration tests (httptest) for validation, rate-limit, caption enforcement.

### Infrastructure & CI
1. Extend docker-compose (MySQL, Redis, nginx-cdn with buffering disabled).
2. Configure GitHub Actions: lint → unit → integration → Playwright → k6 smoke → docker build.
3. Provision preview namespaces with feature flag gating admin UI.
4. Build Grafana dashboards + alerts for playback error rate and throttle events.

### Risks & Mitigations
- **Streaming latency spikes** – Use parallel fetch/caching (Context7), CDN tuning, preloading posters.
- **Caption ingestion failures** – Validate assets on upload, block publish if missing, provide admin retries.
- **Aggressive throttling** – Start with conservative thresholds, monitor dashboards, allow QA bypass header.

### Outstanding Follow-Ups
- Document explicit out-of-scope features (multi-CDN switching, DRM) in kickoff notes.
- Confirm CDN contract for signed URL validation behavior.
