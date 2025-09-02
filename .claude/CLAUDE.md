# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ArchesAI is a comprehensive data processing platform with a hybrid architecture:

- **Backend**: Go API server using Echo framework with Domain-Driven Design
- **Frontend**: TypeScript/React with TanStack Router, built with Vite
- **Database**: PostgreSQL with vector extensions for embeddings
- **Monorepo**: pnpm workspaces + Nx for TypeScript, standalone Go module

## Essential Development Commands

### Backend (Go)

```bash
# Development with hot reload
make dev  # or make watch

# Build and run
make build
make run

# Code generation (required after API/DB schema changes)
make generate  # Runs both sqlc and oapi-codegen
make sqlc      # Generate database code from SQL queries
make oapi      # Generate server code from OpenAPI spec

# Testing and quality
make test             # Run tests
make test-coverage    # Generate coverage report
make lint            # Run all linters (Go + OpenAPI)
make lint-go        # Run Go linter only

# Database migrations
make migrate-up      # Apply migrations
make migrate-down    # Rollback migrations
make migrate-create name=<migration_name>  # Create new migration
```

### Frontend (TypeScript/React)

```bash
# Development
pnpm dev:platform # Run platform web app
nx dev platform   # Alternative using Nx

# Build
pnpm build # Build all packages
nx run-many -t build -t tsc:build

# Quality checks
pnpm lint       # Run ESLint on all packages
pnpm lint:fix   # Fix linting issues
pnpm typecheck  # TypeScript type checking
pnpm format     # Check Prettier formatting
pnpm format:fix # Fix formatting issues

# Client SDK generation
pnpm client:generate # Generate TypeScript client from OpenAPI
```

### Running a Single Test

```bash
# Go tests
go test -v -run TestFunctionName ./path/to/package

# TypeScript/React tests
pnpm test -- --run path/to/test.spec.ts
vitest run path/to/test.spec.ts
```

## Architecture & Code Structure

### Backend Architecture (Go)

The Go backend follows Domain-Driven Design with clean architecture:

```
cmd/api/          # Application entry point
  main.go         # Server initialization

internal/
  app/            # Application layer - dependency injection, route registration
    deps.go       # Container with all dependencies

  domains/        # Business logic organized by domain
    auth/         # Authentication domain
      entities/   # Domain models
      handlers/   # HTTP handlers implementing OpenAPI interfaces
      services/   # Business logic
      repositories/  # Data access interfaces
      adapters/postgres/  # PostgreSQL implementation
      middleware.go       # Auth middleware

    admin/        # Admin domain (similar structure)
    intelligence/ # AI/ML features domain

  generated/      # Generated code (DO NOT EDIT)
    api/          # OpenAPI server stubs
    database/     # sqlc generated queries

  infrastructure/ # Technical concerns
    config/       # Configuration management (Viper)
    database/     # Database connection, migrations, queries
    server/       # HTTP server setup
```

### Frontend Architecture (TypeScript/React)

Monorepo structure with shared packages:

```
web/
  platform/       # Main web application
    src/
      routes/     # TanStack Router file-based routing
      lib/        # Utilities and helpers

  client/         # Generated API client
    src/
      services/   # Auto-generated from OpenAPI

  ui/             # Shared UI components library
    src/
      components/ # Reusable React components
      hooks/      # Custom React hooks
      lib/        # UI utilities

  schemas/        # Shared Zod schemas

tools/            # Build tools and configs
  eslint/         # Shared ESLint config
  prettier/       # Shared Prettier config
  typescript/     # Shared TypeScript configs
```

### Key Patterns and Conventions

#### API Development Flow

1. Define endpoints in `api/openapi.yaml` or component files
2. Run `make oapi` to generate server interfaces
3. Implement handlers in `internal/domains/*/handlers/`
4. Handlers must satisfy generated interfaces from `internal/generated/api/`

#### Database Development Flow

1. Create migration: `make migrate-create name=add_users_table`
2. Write SQL queries in `internal/infrastructure/database/queries/`
3. Run `make sqlc` to generate type-safe query functions
4. Use generated queries via `container.Queries` in services

#### Frontend Development Flow

1. API client is auto-generated from OpenAPI spec
2. Use `@archesai/client` package to make API calls
3. Components go in `@archesai/ui` for reusability
4. Routes are file-based in `web/platform/src/routes/`

#### Dependency Injection Pattern

- All dependencies are wired in `internal/app/deps.go`
- Container pattern provides dependencies to handlers
- Services receive repositories via interfaces (for testability)

#### Authentication Flow

- JWT-based authentication with refresh tokens
- Auth middleware validates tokens and adds claims to context
- Protected routes use `middleware.RequireAuth()`

## Environment Configuration

Backend configuration uses Viper and reads from:

- Environment variables (prefix: `ARCHESAI_`)
- `.env` file (local development)
- `config.yaml` (defaults)

Key environment variables:

```bash
ARCHESAI_DATABASE_URL=postgres://user:pass@localhost/archesai?sslmode=disable
ARCHESAI_JWT_SECRET=your-secret-key
ARCHESAI_SERVER_PORT=8080
ARCHESAI_SERVER_HOST=0.0.0.0
```

## Common Tasks

### Adding a New API Endpoint

1. Define in `api/openapi.yaml` or create new path file in `api/paths/`
2. Run `make oapi` to regenerate interfaces
3. Implement handler methods to satisfy the interface
4. Register handler in `internal/app/deps.go`

### Adding a New Database Table

1. Create migration: `make migrate-create name=create_table_name`
2. Write SQL queries in `internal/infrastructure/database/queries/table_name.sql`
3. Run `make sqlc` to generate query functions
4. Create repository interface in domain
5. Implement repository using generated queries

### Creating a New React Component

1. Add component to `web/ui/src/components/`
2. Export from `web/ui/src/index.ts`
3. Import in platform app as `@archesai/ui`

## Testing Strategy

- **Go**: Table-driven tests, interfaces for mocking
- **TypeScript**: Vitest for unit tests, Playwright for E2E
- Run `make test` (Go) or `pnpm test` (TS) before commits
- Coverage reports: `make test-coverage` (Go), `pnpm test:coverage` (TS)

## Important Notes

- **Generated Code**: Never edit files in `internal/generated/` - they are overwritten
- **OpenAPI First**: API changes start in OpenAPI spec, not code
- **Type Safety**: Both Go (sqlc) and TypeScript (orval) use code generation for type safety
- **Monorepo**: Use pnpm workspaces - dependencies are shared via catalog in `pnpm-workspace.yaml`
- **Domain Boundaries**: Keep domains isolated - communicate through interfaces
- **Migration Safety**: Always review migrations before applying to production
