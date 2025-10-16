# Data Model – Movie Streaming Portal

## Movie
- **Primary Key**: `movie_id` (UUID)
- **Natural Key**: `title` (globally unique, case-insensitive)
- **Fields**:
  - `title` (string, unique, required)
  - `slug` (string, generated from title, immutable)
  - `synopsis` (text, required)
  - `genres` (string[], minimum 1)
  - `duration_minutes` (int, >0)
  - `availability_start` (timestamp with timezone, required)
  - `availability_end` (timestamp with timezone, required, > start)
  - `visibility` (enum: `VISIBLE`, `HIDDEN`)
  - `poster_url` (string, HTTPS URL)
  - `created_at` / `updated_at` (timestamps)
  - `created_by` / `updated_by` (UUID referencing User)
- **Relationships**:
  - 1:N with `StreamSource`
  - 1:N with `CaptionTrack`
- **Validation Rules**:
  - Availability window must be in the future for new movies.
  - Duration must match associated stream metadata (validated on ingest).

## StreamSource
- **Primary Key**: `stream_source_id` (UUID)
- **Fields**:
  - `movie_id` (FK → Movie)
  - `stream_url` (string, HTTPS `.m3u8`, encrypted at rest)
  - `quality_tag` (string, e.g., `1080p`, nullable for single quality)
  - `drm_required` (bool)
  - `last_validated_at` (timestamp)
  - `validation_status` (enum: `VALID`, `INVALID`, `UNKNOWN`)
- **Validation Rules**:
  - URL must end with `.m3u8`.
  - Only one canonical source marked `primary=true`.

## CaptionTrack
- **Primary Key**: `caption_track_id` (UUID)
- **Fields**:
  - `movie_id` (FK → Movie)
  - `language_code` (ISO 639-1, required, unique per movie)
  - `caption_url` (string, HTTPS VTT or TTML)
  - `format` (enum: `WEBVTT`, `TTML`)
  - `is_default` (bool, at most one true per movie)
  - `last_validated_at` (timestamp)
  - `validation_status` (enum)
- **Validation Rules**:
  - At least one caption per movie; one default.
  - Caption URL must be reachable during ingest.

## ViewingSession
- **Primary Key**: `session_id` (UUID)
- **Fields**:
  - `movie_id` (FK → Movie)
  - `viewer_id` or `anonymous_id` (string)
  - `started_at` (timestamp)
  - `ended_at` (timestamp, nullable)
  - `completion_percent` (0–100)
  - `device_type` (enum)
  - `player_errors` (jsonb array of codes)
- **Usage**: Analytics and incident triage (observability Principle IV).

## SignedStreamToken
- **Primary Key**: composite (`token_id` UUID, `movie_id`)
- **Fields**:
  - `movie_id` (FK → Movie)
  - `viewer_context` (string hash of viewer/session)
  - `signed_url` (string, HTTPS, expires quickly)
  - `expires_at` (timestamp)
  - `issued_at` (timestamp)
- **Validation Rules**:
  - TTL ≤ 5 minutes.
  - Unique per viewer/session pair for a sliding window.
- **Storage**: Redis (ephemeral) + audit trail in MySQL for compliance.

## RateLimitBucket
- **Primary Key**: composite (`scope`, `bucket_key`)
- **Fields**:
  - `scope` (enum: `VIEWER_READ`, `ADMIN_MUTATION`)
  - `bucket_key` (string: IP, session, or user ID)
  - `hits` (int)
  - `window_start` (timestamp)
- **Storage**: Redis for counters, with Go service emitting logs when thresholds exceeded.
