# Research Summary – Movie Streaming Portal

## Next.js App Router Streaming Strategy
- **Decision**: Use server components with route-level Suspense/loading states and ensure proxies disable buffering (`X-Accel-Buffering: no`) so movie pages stream metadata and player shell before playback is ready.
- **Rationale**: Context7 Next.js production guidance highlights Suspense + Loading UI to prevent blocking renders and recommends turning off proxy buffering for streaming responses, keeping first paint fast even with remote fetches.
- **Alternatives Considered**: Fully client-rendered pages (rejected: slower TTFB, larger bundles); static pre-rendering (rejected: schedule-driven data too dynamic).

## Data Fetching & Caching
- **Decision**: Fetch movie metadata and signed stream tokens in parallel server actions, using App Router caching controls and `unstable_cache` for non-`fetch` lookups where safe.
- **Rationale**: Context7 docs emphasize parallel fetches and explicit caching in App Router to avoid waterfalls while keeping dynamic data opt-in, fitting our live schedule and token lifetimes.
- **Alternatives Considered**: Sequential fetches (rejected: increases latency); client-side fetch for tokens (rejected: exposes secrets, increases bundle size).

## Go Service Concurrency & Cancellation
- **Decision**: Wrap movie APIs in contexts with per-request deadlines, leveraging `context.WithTimeout` and honoring cancellation in repository calls; propagate cancellation to stop token generation early.
- **Rationale**: Go standard library guidance (Context package, `http.Request.Context`) stresses deadline-aware handlers to prevent resource leaks, aligning with our rate limits and signed URL issuance.
- **Alternatives Considered**: Ignoring context cancellations (rejected: risks runaway goroutines); global timeouts without context (rejected: less granular control).

## Rate Limiting & Observability
- **Decision**: Implement IP/session rate limits for viewer endpoints (120 rpm) and stricter per-account limits for admin mutations (20 rpm), storing counters in Redis and logging exceed events for dashboards.
- **Rationale**: Aligns with specification guardrails and Go net/http configurability; tracking via OpenTelemetry matches constitution Principle IV requirements.
- **Alternatives Considered**: No rate limits (rejected: abuse risk); third-party gateway only (rejected: keeps core logic unaware, harder to test).

## Caption Asset Management
- **Decision**: Require at least one caption track per movie, validate metadata during admin submissions, and expose toggle controls in the player.
- **Rationale**: Meets accessibility success metrics and integrates smoothly with App Router server components that stream metadata plus available caption tracks.
- **Alternatives Considered**: Optional captions (rejected: fails accessibility target); outsourcing to external caption service (rejected: adds latency, cost for MVP).
