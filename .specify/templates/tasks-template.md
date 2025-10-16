---
description: "Task list template for feature implementation"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`  
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/  
**Tests**: Contract, integration, and end-to-end tests are mandatory for every story that touches
shared API contracts or production-critical flows; unit tests are required for Go services and
client components.  
**Organization**: Tasks are grouped by user story so each slice can be delivered, tested, and
released independently.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions (frontend, backend, contracts, infrastructure)

## Path Conventions
- Frontend: `frontend/app/`, `frontend/components/`, `frontend/lib/`, `frontend/tests/`
- Backend: `backend/cmd/`, `backend/internal/{api,service,domain,persistence}`, `backend/tests/`
- Contracts: `contracts/{openapi,proto}/`, generators in `contracts/generators/`
- Infrastructure: `infrastructure/{k8s,terraform,pipelines}/`
- Scripts/tooling: `scripts/`, `.github/workflows/`

<!--
  ============================================================================
  IMPORTANT: The tasks below are SAMPLE TASKS for illustration.

  The /speckit.tasks command MUST replace these with actual tasks based on:
  - User stories from spec.md (priorities P1, P2, P3…)
  - Requirements from plan.md
  - Contracts and data models
  - Constitution Principles (I–V)

  DO NOT keep these sample tasks in the generated tasks.md file.
  ============================================================================
-->

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Ensure tooling, contracts, and environments are ready.

- [ ] T001 Align repo structure with plan (`frontend/`, `backend/`, `contracts/`, `infrastructure/`)
- [ ] T002 Install/update dependencies (`pnpm install`, `go mod tidy`, contracts codegen)
- [ ] T003 [P] Configure linting/formatting (`pnpm lint`, `golangci-lint`, `prettier`, `gofmt`)
- [ ] T004 [P] Verify CI pipelines run lint, test, build, and preview deployments

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared work required before any user story can begin.

**⚠️ CRITICAL**: No user story work starts until this phase is complete.

- [ ] T010 Update `contracts/openapi|proto/[domain]` and regenerate Go/TS types
- [ ] T011 Create database migrations in `backend/internal/persistence/migrations`
- [ ] T012 [P] Extend Go middleware (auth, observability) in `backend/internal/api/middleware`
- [ ] T013 [P] Prepare App Router layouts/loading/error in `frontend/app/[segment]/`
- [ ] T014 Configure feature flags/config entries with documentation in `specs/.../data-model.md`
- [ ] T015 Hook up OpenTelemetry logging/tracing/metrics across both stacks

**Checkpoint**: Foundation complete – user story development can proceed.

---

## Phase 3: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [Brief description of what this story delivers]  
**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 1 (must be written first)

- [ ] T020 [P] [US1] Contract test in `backend/tests/contract/[name]_test.go`
- [ ] T021 [P] [US1] Integration test in `backend/tests/integration/[name]_test.go`
- [ ] T022 [P] [US1] Playwright/Cypress E2E in `frontend/tests/e2e/[name].spec.ts`

### Implementation for User Story 1

- [ ] T023 [P] [US1] Implement domain logic in `backend/internal/domain/[entity].go`
- [ ] T024 [P] [US1] Add service orchestration in `backend/internal/service/[service].go`
- [ ] T025 [US1] Expose handler in `backend/internal/api/[segment]/handler.go`
- [ ] T026 [US1] Implement App Router page in `frontend/app/[segment]/page.tsx`
- [ ] T027 [US1] Create client component (if needed) in `frontend/components/[component].tsx` with `"use client"`
- [ ] T028 [US1] Wire observability (logs/traces/metrics) end-to-end
- [ ] T029 [US1] Update documentation (`specs/.../quickstart.md`, API docs, changelog)

**Checkpoint**: User Story 1 is independently testable and deployable.

---

## Phase 4: User Story 2 - [Title] (Priority: P2)

**Goal**: [Brief description of what this story delivers]  
**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 2

- [ ] T030 [P] [US2] Extend contract coverage for new/updated endpoints
- [ ] T031 [P] [US2] Visual regression or accessibility check in `frontend/tests`

### Implementation for User Story 2

- [ ] T032 [P] [US2] Extend repository in `backend/internal/persistence/[repo].go`
- [ ] T033 [US2] Add background worker/cron in `backend/cmd/[worker]/main.go`
- [ ] T034 [US2] Implement edge UI states in `frontend/app/[segment]/loading.tsx`/`error.tsx`
- [ ] T035 [US2] Document contract changes and notify consumers

**Checkpoint**: User Stories 1 and 2 both operate independently.

---

## Phase 5: User Story 3 - [Title] (Priority: P3)

**Goal**: [Brief description of what this story delivers]  
**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 3

- [ ] T040 [P] [US3] Synthetic monitoring or canary scenario scripts
- [ ] T041 [P] [US3] Load/performance test for Go endpoints (k6, vegeta)

### Implementation for User Story 3

- [ ] T042 [P] [US3] Evolve shared contract and regenerate artifacts
- [ ] T043 [US3] Implement streaming/background updates in App Router route
- [ ] T044 [US3] Harden Go concurrency (context, worker pools, retries)

**Checkpoint**: All stories function independently and are production ready.

---

[Add more user story phases as needed, following the same pattern]

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Stabilize, document, and operationalize the release.

- [ ] T050 [P] Update documentation in `docs/` and `specs/.../quickstart.md`
- [ ] T051 Remove temporary flags and perform code cleanup
- [ ] T052 Run Go profiling (`pprof`, `benchstat`) and frontend performance audits (Lighthouse)
- [ ] T053 [P] Accessibility and internationalization review of App Router surfaces
- [ ] T054 Security hardening (dependency bumps, threat model updates)
- [ ] T055 Validate quickstart.md steps against real environment and update runbooks

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately
- **Foundational (Phase 2)**: Blocks all user stories
- **User Stories (Phase 3+)**: Unblocked once Phase 2 is complete; can run in parallel per capacity
- **Polish (Final Phase)**: Begins after targeted user stories are complete

### User Story Dependencies

- **User Story 1 (P1)**: No dependency on other stories once foundation is ready
- **User Story 2 (P2)**: May integrate with US1 but must remain independently deployable/testable
- **User Story 3 (P3)**: May integrate with earlier stories but must stay independently testable

### Within Each User Story

- Contracts/migrations before service logic
- Domain/service layers before transport handlers
- Transport handlers before frontend integration
- Tests written first, failing, then implementation to make them pass
- Story complete (including docs, observability, flags) before moving on

### Parallel Opportunities

- All tasks marked `[P]` can proceed concurrently
- Different user stories may run in parallel after foundational tasks finish
- Frontend and backend work for the same story often run in parallel once contracts are final

---

## Parallel Example: User Story 1

```bash
# Launch key automated checks for User Story 1 in parallel:
Task: "Contract test in backend/tests/contract/[name]_test.go"
Task: "Playwright E2E in frontend/tests/e2e/[name].spec.ts"
Task: "Integration test in backend/tests/integration/[name]_test.go"
```
