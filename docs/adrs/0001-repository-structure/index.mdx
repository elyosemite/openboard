# 0001. Repository structure and HTTP stack

- Status: Accepted
- Date: 2026-11-15

## Context

OpenBoard has a Go backend, a React frontend and a PostgreSQL schema managed by Liquibase. The three evolve together: an API change usually touches the contract, the backend, the frontend and sometimes a migration. The project values simplicity before abstraction and small, reviewable increments.

## Decision

1. **Single monorepo** with one top-level directory per deliverable:
   - `backend/`: Go modular monolith (`cmd/api`, `internal/platform`, `internal/<module>`).
   - `frontend/`: React + TypeScript app organized by feature (`src/features/<feature>`), with a shared design system and API client under `src/shared`.
   - `database/`: Liquibase migrations organized by object type, plus the runner image used by local compose and CD.
   - `docs/`: ADRs and the OpenAPI contract (`docs/api/openapi.yaml`).
2. **Standard library HTTP routing** (`net/http` `ServeMux` with method patterns, Go 1.22+). No router dependency until a concrete need appears (for example, middleware groups the stdlib cannot express cleanly).
3. **RFC 9457 Problem Details** for every API error, produced centrally by `internal/platform/httpx`. Handlers return errors; they never write error bodies themselves.
4. **Operational probes outside `/api/v1`**: `GET /healthz` (liveness, no dependency checks) and `GET /readyz` (readiness, runs registered checkers).
5. **The frontend reaches the API on the same origin**: the Vite dev server and the nginx image proxy `/api`, `/healthz` and `/readyz` to the backend, so no CORS configuration is needed.

## Consequences

- One pull request can change contract, backend, frontend and migrations atomically.
- CI must run backend and frontend checks independently (separate toolchains in one repo).
- Adopting a router library later means replacing `httpx.Route` and the mux wiring in `cmd/api/routes.go`. Handlers stay unchanged because they use the plain `http.ResponseWriter`/`*http.Request` signature.
