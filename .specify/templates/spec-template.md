# Feature Specification: [FEATURE NAME]

**Feature Branch**: `[###-feature-name]`  
**Created**: [DATE]  
**Status**: Draft  
**Input**: User description: "$ARGUMENTS"
**Tech Stack**: Next.js 15 App Router (TypeScript) frontend, Go 1.23+ services backend, shared API
contracts

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - [Brief Title] (Priority: P1)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently - e.g., "Can be fully tested by [specific action] and delivers [specific value]"]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]
2. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 2 - [Brief Title] (Priority: P2)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 3 - [Brief Title] (Priority: P3)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- What happens when streaming responses are delayed or cancelled mid-flight?
- How does the system handle Go backend timeouts, retries, and idempotency?
- What is the fallback when App Router data fetching fails (loading/error states)?
- How are unauthorized or unauthenticated requests surfaced to the user?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

#### Frontend (Next.js 15 App Router)
- **FR-001**: App Router route at `app/[segment]/page.tsx` MUST [deliver capability] with server
  components by default.
- **FR-002**: Client components MUST be annotated with `"use client"` and limited to [stateful UI].
- **FR-003**: Server actions MUST validate input using shared schemas before invoking backend calls.
- **FR-004**: Loading, error, and metadata files MUST be defined for the route hierarchy.

#### Backend (Go Services)
- **FR-101**: Go handler at `backend/internal/api/[domain]/handler.go` MUST [perform capability].
- **FR-102**: Service layer MUST enforce business rules and propagate `context.Context`.
- **FR-103**: Persistence layer MUST use repository interfaces defined in `internal/persistence`.
- **FR-104**: Concurrency-sensitive operations MUST use cancellation, timeouts, and retries.

#### Shared Contracts & Integration
- **FR-201**: Contract file `[openapi|proto]/[domain].(yaml|proto)` MUST define or extend endpoints
  for this feature with backward compatibility notes.
- **FR-202**: Generated Go and TypeScript types MUST be updated via `pnpm contracts:generate` and
  `go generate ./...`.
- **FR-203**: Observability metadata (logs, traces, metrics) MUST include correlation IDs and user
  context where permitted.

*Example of marking unclear requirements:*

- **FR-006**: System MUST authenticate users via [NEEDS CLARIFICATION: auth method not specified - email/password, SSO, OAuth?]
- **FR-007**: System MUST retain user data for [NEEDS CLARIFICATION: retention period not specified]

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

### Operational Requirements

- **OP-001**: Deployment strategy MUST include canary rollout steps and monitoring checkpoints.
- **OP-002**: Feature toggles MUST be documented with owner, rollout plan, and removal criteria.
- **OP-003**: Security scans (SAST/DAST/dependency) MUST be green before release.

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: [Measurable metric, e.g., "Users can complete account creation in under 2 minutes"]
- **SC-002**: [Measurable metric, e.g., "System handles 1000 concurrent users without degradation"]
- **SC-003**: [User satisfaction metric, e.g., "90% of users successfully complete primary task on first attempt"]
- **SC-004**: [Business metric, e.g., "Reduce support tickets related to [X] by 50%"]
- **SC-005**: [Operational metric, e.g., "p95 latency ≤200 ms for new endpoints during canary"]
