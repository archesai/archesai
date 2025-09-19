# ArchesAI – AI coding agent instructions

This repo is a Go monorepo with a TypeScript frontend. It uses OpenAPI-first development, code generation, and a hexagonal architecture. Follow these conventions to be productive and avoid breaking patterns.

## Big picture

- Stack: Go (API, CLI, TUI), Node/PNPM (platform/docs), PostgreSQL/SQLite via sqlc, Redis, Echo HTTP server.
- Architecture: Hexagonal + DDD. Domains live under `internal/<domain>/` with a flat structure. Handlers → Services → Repositories. Generated code wires interfaces; you implement business logic in `service.go` only.
- API contracts: Defined in `api/openapi.yaml` (split under `api/components` and `api/paths`), bundled to `api/openapi.bundled.yaml`.
- App composition: `internal/app/app.go` builds infra, repositories, services, handlers in parallel and registers routes in `internal/app/routes.go` (public vs protected groups, org-scoped routes).

## Core workflows (use Makefile)

- Generate after changing OpenAPI/SQL/x-codegen: `make generate`
- Lint and format: `make lint` and `make format`
- Run dev backends (hot reload via air): `make dev-api` or all services: `make dev-all`
- API server: `make run-api` or `archesai api`; TUI: `archesai tui`; Chat: `archesai tui --chat`
- Tests: `make test` (see also `test-short`, `test-coverage`, `test-coverage-html`)
- Frontend: `pnpm -F @archesai/platform dev`; Docs: `pnpm -F @archesai/docs dev`

## Generated code and where to edit

- Do not edit generated files: `*.gen.go`, `api/openapi.bundled.yaml`, JS client under `web/client/src/generated`, Helm schema.
- Services: Implement business logic in `internal/<domain>/service.go`. Generated service/repo interfaces live alongside (e.g., `repository.gen.go`, `service.gen.go`).
- Repositories: Prefer generated repository interfaces; SQL is generated via `internal/database` (sqlc). If you need custom repo methods, define them through codegen templates/x-codegen and re-run `make generate`.
- Handlers: Use generated Strict types. Instantiate with `NewStrictServer(...)` and expose via `NewStrictHandler` in routes. See `internal/app/app.go` and `internal/app/routes.go` for canonical wiring.

## Routing, auth, and middleware

- Base group: `/api/v1`. Public routes: Accounts, Sessions, Health, OAuth (no `/v1` prefix for OAuth per spec). Protected routes use `AuthMiddleware.RequireAuth()`; org-scoped routes add `RequireOrganizationMember()`.
- Register with generated helpers: `accounts.RegisterHandlers`, etc. See `internal/app/routes.go` for examples.

## Data and domain patterns

- Pagination: List methods follow `List<Entity>s(ctx, params List<Entity>sParams) ([]*Entity, total int64, error)`. Use `PageQuery` inside params.
- IDs: Use `uuid.UUID` (github.com/google/uuid). Always thread `context.Context`.
- Logging: Pass `*slog.Logger` into services; prefer structured logs.

## OAuth and sessions

- Providers configured in `internal/oauth/service.go` from `config.Auth.*`. Add providers by wiring constructors and honoring `Enabled` flags.
- Tokens issued via `internal/sessions.Service` (`GenerateAccessToken`, `GenerateRefreshToken`).

## Testing and mocks

- Run tests with `make test`.
- Use Mockery-generated mocks only (v3). Generate via `make generate-mocks` or `go tool mockery`. Mocks live in `mocks_test.go`. Don’t hand-roll mocks.

## Frontend and clients

- TypeScript client is generated from OpenAPI via Orval: `make generate-js-client` (lives under `web/client/src/generated`).
- Docs site consumes `api/openapi.bundled.yaml` via `make prepare-docs`/`build-docs`.

## Common tasks cheat sheet

- After editing `api/**` or SQL: run `make generate` then `make lint`.
- Adding a domain: create `internal/<domain>/service.go`, rely on generated interfaces, wire in `internal/app/app.go` and register in `routes.go` using Strict handlers.
- Route gating: put general resources under protected `v1` group; org resources use the org group with membership middleware.

Key references: `README.md`, `Makefile`, `docs/architecture/system-design.md`, `docs/architecture/project-layout.md`, `internal/app/{app.go,routes.go}`, `internal/database`, `.claude/CLAUDE.md`.
