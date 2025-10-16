# Feature Specification: Movie Streaming Portal

**Feature Branch**: `001-specify-scripts-bash`  
**Created**: 2025-10-17  
**Status**: Draft  
**Input**: User description: "เว็บที่สามารถดูหนังได้และสามารถเพิ่มหนังที่จะฉายได้ โดยเก็บ ลิงค์ .m3u8 เอาไว้ใน mysqli"  
**Tech Stack**: Next.js 15 App Router (TypeScript) frontend, Go 1.23+ services backend, shared API contracts

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Stream a Published Movie (Priority: P1)

As a visitor, I want to open a movie detail page and start watching the stream immediately so that I can enjoy available titles without friction.

**Why this priority**: Streaming playback is the core value proposition; without it, the site fails to deliver user value.

**Independent Test**: A tester can publish a single movie, navigate to its watch page, and verify playback, controls, and error handling without any other stories implemented.

**Acceptance Scenarios**:

1. **Given** a movie with an active availability window and a valid `.m3u8` link, **When** the visitor opens `app/movies/[movieId]/page.tsx`, **Then** the page loads metadata and the video begins streaming within 4 seconds.
2. **Given** a streaming error occurs while loading the `.m3u8`, **When** the visitor attempts playback, **Then** the UI presents a descriptive error state and offers a retry action without crashing the session.

---

### User Story 2 - Add a New Movie to the Catalog (Priority: P2)

As a content manager, I want to register new movies with scheduling details and a streaming link so that upcoming titles are ready for viewers.

**Why this priority**: Keeping the catalog fresh enables the business to showcase new releases and maintain engagement.

**Independent Test**: A tester can log in as a manager, fill out the movie creation form once backend contracts exist, and verify the movie appears in listings with correct schedule details.

**Acceptance Scenarios**:

1. **Given** a manager is on `app/admin/movies/new/page.tsx` with required permissions, **When** they submit title, synopsis, poster, availability window, and `.m3u8` link, **Then** the movie is persisted, validated, and appears in the upcoming list with a success confirmation.
2. **Given** the submitted `.m3u8` link fails validation, **When** the manager attempts to save the form, **Then** the system blocks the submission, explains the error, and preserves entered data.

---

### User Story 3 - Manage Upcoming Titles (Priority: P3)

As a content manager, I want to review and edit scheduled movies so that incorrect data or expired links can be corrected without developer support.

**Why this priority**: Self-service maintenance avoids downtime caused by stale assets and reduces operational workload.

**Independent Test**: A tester can edit an existing movie’s metadata, reschedule its availability, and deactivate a title while confirming audience-facing pages update accordingly.

**Acceptance Scenarios**:

1. **Given** a scheduled movie exists, **When** the manager adjusts start or end times from `app/admin/movies/[movieId]/page.tsx`, **Then** the new window is saved and reflected immediately on viewer listings.
2. **Given** a movie should be temporarily unavailable, **When** the manager toggles its visibility, **Then** viewers no longer see it in browse lists while the record remains intact for future activation.

---

### Edge Cases

- `.m3u8` link returns 404, 403, or a non-HLS payload during playback.
- Scheduled availability window is set in the past or overlaps with another event requiring exclusivity.
- Duplicate title submissions MUST trigger a blocking validation error that guides managers to select a distinct name.
- Viewer attempts to access a movie outside of its availability window.
- MySQL connectivity or transaction failure occurs while saving a movie record.
- Caption track referenced by the movie is missing, corrupted, or mismatched with the declared language when a viewer enables captions.
- Rate limits block high-frequency requests from a single IP or admin account; the UI must surface clear retry guidance.

## Requirements *(mandatory)*

### Functional Requirements

#### Frontend (Next.js 15 App Router)
- **FR-001**: The movie detail page MUST display metadata (title, synopsis, availability window) and start playback automatically when a valid `.m3u8` link is available.
- **FR-002**: Playback controls MUST cover play/pause, seek, quality selection (if provided), and fullscreen, while keeping client-side state minimal.
- **FR-003**: Actions triggered from the watch page MUST verify viewer eligibility and log watch-start events before requesting backend telemetry updates.
- **FR-004**: Movie and admin catalog pages MUST present dedicated loading, error, and maintenance states to keep visitors informed during data fetches or outages.
- **FR-005**: The movie creation form MUST enforce validation for required fields, highlight errors inline, and preserve user input after server-side validation failures.
- **FR-006**: The player MUST provide caption controls, default to an available caption track in at least one primary language, and allow viewers to toggle captions on or off during playback.

#### Backend (Go Services)
- **FR-101**: The movie service MUST offer endpoints to list, retrieve, create, update, and deactivate movies, returning structured JSON aligned with shared contracts.
- **FR-102**: Business rules MUST validate scheduling windows (start before end, exclusivity conflicts) and honor request cancellation via propagated contexts.
- **FR-103**: Movie metadata, including `.m3u8` URLs, MUST be written to MySQL within a single transaction to avoid partial state.
- **FR-104**: Streaming metadata retrievals MUST enforce timeouts and retries, surfacing structured errors when upstream content delivery networks are unreachable.
- **FR-105**: Audit logging MUST capture create, update, and visibility changes with actor identifiers and timestamps.
- **FR-106**: Playback requests MUST obtain time-limited signed `.m3u8` URLs per session, expiring within minutes and scoped to the requesting viewer context.
- **FR-107**: Title collisions MUST be rejected so each movie title in the catalog remains globally unique; duplicate submissions return actionable validation errors.
- **FR-108**: Movie creation and update endpoints MUST require at least one caption track reference and validate that caption files match declared language metadata.
- **FR-109**: Viewer read endpoints MUST enforce per-IP or per-session rate limits (target 120 requests per minute) and respond with descriptive throttle errors when exceeded.
- **FR-110**: Admin mutation endpoints MUST enforce stricter per-account limits (target 20 writes per minute) and log throttle events for audit review.

#### Shared Contracts & Integration
- **FR-201**: The movie contract specification MUST define endpoints for catalog listing, detail retrieval, creation, updates, and visibility toggling with clear backward compatibility notes.
- **FR-202**: Shared type generation MUST be refreshed before release so frontend and backend agree on request/response schemas.
- **FR-203**: Logs, traces, and metrics MUST include correlation IDs, movie identifiers, and actor identifiers for all catalog interactions.
- **FR-204**: Contract schemas MUST constrain `.m3u8` fields to HTTPS URLs, specify validation errors, and document expected retry guidance.

*Example of marking unclear requirements:*

- **FR-006**: System MUST authenticate users via [NEEDS CLARIFICATION: auth method not specified - email/password, SSO, OAuth?]
- **FR-007**: System MUST retain user data for [NEEDS CLARIFICATION: retention period not specified]

### Key Entities *(include if feature involves data)*

- **Movie**: Represents a title with fields for title (globally unique), synopsis, genres, poster asset, duration, availability window, visibility flag, and associated stream sources.
- **StreamSource**: Stores the canonical `.m3u8` URL, quality tags, DRM flags, and last validation timestamp linked to a movie.
- **ViewingSession**: Captures viewer interactions (start time, completion percentage, device info) to support analytics and error diagnostics.
- **SignedStreamToken**: Ephemeral artifact containing the time-limited playback URL, expiry timestamp, and viewer/session identifiers used to authorize stream access.
- **CaptionTrack**: Stores caption metadata (language code, format, CDN location, validation status) linked to each movie to satisfy accessibility requirements.
- **UserRole**: Defines permissions (Visitor vs Content Manager) governing access to admin routes and endpoints.

### Operational Requirements

- **OP-001**: Deployment strategy MUST include canary rollout verifying streaming playback, catalog creation, and rollback readiness with automated smoke tests.
- **OP-002**: Feature toggles MUST control exposure of the admin catalog interface, with documented owner, rollout plan, and removal timeline.
- **OP-003**: Security scans (SAST, DAST, dependency audits) MUST pass before enabling public streaming, and `.m3u8` URLs MUST be stored encrypted at rest.
- **OP-004**: Monitoring dashboards MUST alert on playback failure rate above 5% or catalog creation errors above 2% in a 10-minute window.
- **OP-005**: Rate limiting dashboards MUST track per-IP and per-account throttle events, alerting operators if limits trigger more than 50 times within 5 minutes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 95% of viewers reach playback within 4 seconds of opening a movie detail page during peak hours.
- **SC-002**: 100% of catalog submissions with valid data are available to viewers within 1 minute of manager confirmation.
- **SC-003**: Playback error rate remains below 5% over a rolling 24-hour period after launch.
- **SC-004**: Content managers report the ability to add or update a movie in under 2 minutes during user acceptance testing.
- **SC-005**: Incident alerts for catalog creation failures trigger within 2 minutes and include actionable context for on-call responders.
- **SC-006**: Infrastructure sustains at least 1,000 simultaneous viewers with no more than 5% degradation in playback quality or startup time.

## Assumptions

- Visitors can stream without authentication, while catalog management requires authenticated content manager accounts.
- One canonical `.m3u8` URL per movie is sufficient for the initial release; multi-bitrate or multi-CDN support will be considered later.
- Existing CDN infrastructure will host the HLS streams, and this feature only references those assets.
- CDN and edge delivery tooling support issuing and validating time-limited signed playback URLs.
- Caption assets for at least one primary language are available and managed alongside movie metadata.
- Launch-ready infrastructure (CDN + backend) is provisioned to handle roughly 1,000 concurrent viewers.
- API gateway or edge middleware supports configurable per-IP and per-account rate limiting policies.

## Dependencies

- Identity and role management system capable of distinguishing content managers from viewers.
- Observability stack (OpenTelemetry, logging pipeline, dashboards) ready to capture new metrics and alerts.
- MySQL database cluster provisioned with capacity for catalog growth and encrypted storage for stream URLs.

## Clarifications

### Session 2025-10-17
- Q: How should viewer access to the `.m3u8` stream be protected to prevent unauthorized sharing? → A: Time-limited signed `.m3u8` URLs generated per session
- Q: What unique identifier should the catalog use so admins avoid duplicate or conflicting entries? → A: Movie title must be unique
- Q: Which accessibility features must the movie player support at launch? → A: Captions required for at least one primary language
- Q: What peak number of simultaneous viewers should the system support at launch? → A: ~1,000 concurrent viewers
- Q: What rate limiting approach should the system enforce at launch? → A: IP/session limits for viewer endpoints and stricter quotas for admin actions
