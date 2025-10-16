# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

**Language/Version**: TypeScript (Next.js 15 App Router), Go 1.23+  
**Primary Dependencies**: Next.js 15, React Server Components, Tailwind CSS or CSS Modules, Go chi/fiber or gRPC, sqlc/ent, OpenTelemetry  
**Storage**: PostgreSQL (primary), Redis (caching/queues) unless feature specifies otherwise  
**Testing**: Playwright or Cypress (E2E), Jest/Testing Library (frontend unit), `go test` + testify (backend unit/integration)  
**Target Platform**: Web (SSR + streaming) deployed on Linux containers via Kubernetes  
**Project Type**: Web application with contracts shared between frontend and backend  
**Performance Goals**: ≤200 ms p95 for critical route handlers, sustain 1k req/s per service, 95% CLS < 0.1  
**Constraints**: Must honor App Router streaming, graceful degradation offline, zero-trust API access  
**Scale/Scope**: Multi-tenant streaming platform with progressive feature rollout

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Principle I (Next.js App Router Discipline): Document how the feature stays within the `app/`
  router, enforces TypeScript types, and scopes client components.
- Principle II (Go Backend Reliability): Describe service layering, `context.Context` handling,
  error strategy, and concurrency controls for new or changed Go endpoints.
- Principle III (Shared Contracts & Type Safety): Identify contract artifacts (OpenAPI/Buf/etc.),
  required version bumps, and generated types for both stacks.
- Principle IV (Observability & Security Assurance): Plan logging, tracing, metrics, and security
  gates (auth scopes, secret usage, scanners).
- Principle V (Continuous Delivery Excellence): Explain CI coverage, preview environment needs, and
  rollback/feature flag strategy.

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

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
