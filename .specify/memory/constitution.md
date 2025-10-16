<!--
Sync Impact Report
Version change: 0.0.0 → 1.0.0
Modified principles:
- I. Next.js 15 App Router Discipline
- II. Go Backend Reliability
- III. Shared Contracts & Type Safety
- IV. Observability & Security Assurance
- V. Continuous Delivery Excellence
Added sections:
- Technology Standards
- Delivery Workflow
Removed sections:
- None
Templates requiring updates:
- ✅ .specify/templates/plan-template.md
- ✅ .specify/templates/spec-template.md
- ✅ .specify/templates/tasks-template.md
Follow-up TODOs:
- None
-->
# Leak Streaming Constitution

## Core Principles

### I. Next.js 15 App Router Discipline
- MUST implement all frontend surfaces in the `app/` directory with server components by default and
  compose shared layouts, loading states, and metadata through the App Router primitives.
- MUST enforce fully typed data flows (TypeScript, Zod/Valibot schemas, and generated client types
  from shared contracts) and prefer server actions for mutations with explicit input validation.
- MUST optimize rendering with streaming, caching, and incremental revalidation while keeping
  client components minimal, isolated, and stateful only when unavoidable.
Rationale: Aligns the frontend with Next.js 15 best practices so multi-route features stay
predictable, performant, and testable.

### II. Go Backend Reliability
- MUST organize Go services using Go modules with `internal/` and `pkg/` boundaries, clean layering
  (transport → service → domain → persistence), and dependency injection for test seams.
- MUST expose APIs via HTTP handlers or gRPC endpoints that honor shared contracts, propagate
  `context.Context`, and return structured errors with traceable codes.
- MUST adopt Go concurrency carefully (timeouts, cancellation, worker pools) and guard all external
  calls with circuit breakers, retries, or idempotency keys as appropriate.
Rationale: Enforces idiomatic, reliable Go services that integrate cleanly with the Next.js frontend
and scale safely.

### III. Shared Contracts & Type Safety
- MUST define API and event contracts once (OpenAPI/Buf/Protobuf) and generate both Go server and
  TypeScript client types from the same source of truth.
- MUST version contracts semantically, document breaking changes before release, and ship backward
  compatible transitions whenever possible.
- MUST validate all external inputs at the boundary (Go request decoders, server actions, middleware)
  and log rejected payloads without leaking secrets.
Rationale: Guarantees that frontend and backend evolve together, eliminating drift and runtime type
errors.

### IV. Observability & Security Assurance
- MUST provide structured logging, distributed tracing, and metrics (OpenTelemetry) across Go
  services and Next.js route handlers, with correlation IDs passed end-to-end.
- MUST enforce security baselines: 12-factor config, secret rotation, HTTPS/TLS everywhere, OWASP
  protections, and zero-trust defaults on internal APIs.
- MUST gate deployments with automated security scans (SAST, dependency checks) and runtime health
  probes, raising incidents when guardrails fail.
Rationale: Maintains trust and rapid debugging for a streaming product handling sensitive data.

### V. Continuous Delivery Excellence
- MUST run automated unit, integration, and end-to-end tests (Playwright/Cypress + Go test suites)
  on every merge, blocking releases on red pipelines.
- MUST publish preview environments for frontend and backend changes, capturing observability and
  performance baselines before production rollout.
- MUST maintain clear rollback plans, versioned infrastructure manifests, and post-release audits to
  continually improve the deployment process.
Rationale: Sustains fast iteration without compromising stability or user experience.

## Technology Standards
- Frontend stack is Next.js 15 with the App Router, React Server Components, TypeScript, Tailwind or
  CSS Modules, and Vite-powered tooling where applicable.
- Backend stack is Go 1.23+, using chi/fiber or gRPC for transport, sqlc or ent for data access, and
  Dockerized services orchestrated via Terraform/Kubernetes manifests.
- Shared tooling includes pnpm (workspace) for frontend packages, Go toolchain for backend, Turbo
  or Nx for task orchestration, and Renovate/Dependabot for dependency hygiene.
- All secrets are managed through Vault or SSM, infrastructure as code is mandatory, and CI/CD runs
  via GitHub Actions with required status checks.
- Best practice definitions follow Context7 MCP guidance; any deviation requires an approved
  exception documented in the relevant spec or plan.

## Delivery Workflow
- Discovery → Plan → Spec → Tasks flow is mandatory; each phase must explicitly check Principles I–V
  and document compliance decisions.
- Code reviews require at least two approvals: one for frontend (App Router expert) and one for
  backend (Go maintainer), each confirming contract adherence and observability hooks.
- Feature branches must include migration scripts, contract changes, and documentation updates in
  the same PR to keep history atomic.
- Production releases follow progressive deployment (canary → staged rollout), with automated
  verification from observability dashboards before marking complete.

## Governance
- This constitution supersedes conflicting project documents; deviations require a written waiver
  approved by the tech lead and recorded in the repo.
- Amendments follow an RFC cycle: proposal → impact analysis (including version delta) → maintainer
  vote; approved changes update this file and all dependent templates in the same commit.
- Constitution versions follow semver: MAJOR for principle changes/removals, MINOR for new guidance,
  PATCH for clarifications. Each amendment logs rationale in commit history.
- Compliance reviews occur quarterly and before major releases; unresolved violations block release
  candidates until addressed or explicitly waived.

**Version**: 1.0.0 | **Ratified**: 2025-10-17 | **Last Amended**: 2025-10-17
