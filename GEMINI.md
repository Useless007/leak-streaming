# leak-streaming Development Guidelines

## Project Overview

This is a monorepo for a video streaming application. The project consists of a Go backend and a Next.js frontend.

**Backend:**

*   Written in Go (1.23+)
*   Uses PostgreSQL for data storage (with `sqlc` for query generation)
*   Uses Redis for caching and token signing
*   Uses `chi` for routing
*   Provides a RESTful API for the frontend
*   Includes OpenTelemetry for observability

**Frontend:**

*   Built with Next.js 15 (App Router, React Server Components, Suspense streaming) and TypeScript
*   Uses Tailwind CSS and Radix UI for styling
*   Provides a user interface for browsing and watching movies
*   Includes an admin section for managing content

## Active Technologies

- TypeScript (Next.js 15 App Router)
- Go 1.23+
- Next.js 15 (App Router, React Server Components, Suspense streaming)
- Tailwind CSS
- Radix UI
- Go chi + sqlc
- MySQL driver
- OpenTelemetry
- Redis (token cache)

## Building and Running

### Docker

The easiest way to get started is to use Docker.

1.  Make sure you have Docker and Docker Compose installed.
2.  Run `docker-compose up -d` to start the Postgres and Redis containers.
3.  Run the backend: `cd backend && go run ./cmd/api`
4.  Run the frontend: `cd frontend && bun dev`

### Manual

**Backend:**

1.  Install Go 1.23 or later.
2.  Install dependencies: `cd backend && go mod tidy`
3.  Run the application: `cd backend && go run ./cmd/api`

**Frontend:**

1.  Install Bun.
2.  Install dependencies: `cd frontend && bun install`
3.  Run the application: `cd frontend && bun dev`

## Development Conventions

*   **Code Style:** Follow standard conventions for TypeScript (Next.js 15 App Router) and Go 1.23+. Adhere to context7 mcp best practices.
*   **Testing:**
    *   Frontend: `cd frontend && bun run test`
    *   Backend: `GOCACHE=... go test ./...`
*   **Linting:** `cd frontend && bun run lint`
*   **Project Structure:**
    ```
    backend/
    frontend/
    tests/
    ```

## Current Status (as of 2025-10-17)

*   **Setup:**
    *   Frontend: Next.js 15 (App Router) + Bun + Vitest/Playwright are set up.
    *   Backend: Go + chi + sqlc + Redis token cache are set up.
    *   Linting and testing scripts are available: `bun run lint`, `bun run test`, `GOCACHE=... go test ./...`.
*   **Features:**
    *   The movie watching flow is fully functional: the frontend requests a `playback-token` from the backend, fetches the modified manifest, and proxies segments with host validation.
    *   Sample subtitles are available at `frontend/public/captions/sample-en.vtt`.
    *   The backend is connected to a PostgreSQL database (managed via docker-compose) with goose migrations and seed data for two demo movies. The repository for reading and writing data is functional.
    *   Redis is used for token caching and is available in docker.
    *   The frontend has a `/movies` page that lists movies from the database, a detail page for streaming videos, and an admin UI at `/admin/movies/new` for adding new movies (including server action, validation, and a textarea component).
*   **Testing:**
    *   Playwright specs for viewer streaming, movie catalogue, and admin creation are available, but require the backend, frontend, and database to be running.
    *   Linting, unit tests, and Go tests pass in the development environment.
*   **Next Steps:**
    *   Enhance integration tests for the backend.
    *   Enable Playwright in the CI/CD pipeline once the stack is ready.
    *   Improve the security of allowed hosts (currently does not support wildcards `*` by design).
    *   Add more admin tools (e.g., list/delete) as needed.