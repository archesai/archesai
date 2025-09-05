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
make generate  # Runs all generators (sqlc, oapi, defaults, converters)
make sqlc      # Generate database code from SQL queries
make oapi      # Generate server code from OpenAPI spec
make generate-defaults   # Generate config defaults from OpenAPI
make generate-converters # Generate type converters

# Testing and quality
make test             # Run tests
make test-coverage    # Generate coverage report
make lint            # Run all linters (Go + OpenAPI)
make lint-go        # Run Go linter only
make format          # Format all code

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

The Go backend follows Domain-Driven Design with a flat, Go-centric structure:

```
cmd/api/          # Application entry point
  main.go         # Server initialization

internal/
  app/            # Application layer - dependency injection, route registration
    deps.go       # Container with all dependencies

  domains/        # Business logic organized by domain (flat structure)
    auth/         # Authentication and user management domain
      auth.go     # Package docs and shared constants
      entities.go # Domain models (extending API types)
      service.go  # Business logic with Repository interface
      repository.go # PostgreSQL implementation
      handler.go  # HTTP handlers implementing OpenAPI interfaces
      middleware.go # Auth middleware
      converters/ # Generated type converters (DO NOT EDIT)

    organizations/ # Organization, membership, and invitation management
      organizations.go # Package docs and shared constants
      entities.go     # Domain models (extending API types)
      service.go      # Business logic with Repository interface
      repository.go   # PostgreSQL implementation
      handler.go      # HTTP handlers implementing OpenAPI interfaces
      converters/     # Generated type converters (DO NOT EDIT)

    workflows/    # Pipeline workflows, runs, and tools domain
      workflows.go # Package docs and shared constants
      entities.go  # Domain models (extending API types)
      service.go   # Business logic with Repository interface
      repository.go # PostgreSQL implementation
      handler.go   # HTTP handlers implementing OpenAPI interfaces
      converters/  # Generated type converters (DO NOT EDIT)

    content/      # Content artifacts and labels domain
      content.go  # Package docs and shared constants
      entities.go # Domain models (extending API types)
      service.go  # Business logic with Repository interface
      repository.go # PostgreSQL implementation
      handler.go  # HTTP handlers implementing OpenAPI interfaces
      converters/ # Generated type converters (DO NOT EDIT)

  generated/      # Generated code (DO NOT EDIT)
    api/          # OpenAPI server stubs
    database/     # sqlc generated queries

  infrastructure/ # Technical concerns
    config/       # Configuration management (Viper)
    database/     # Database connection, migrations, queries
    server/       # HTTP server setup
```

### Domain Pattern

Each domain follows a consistent flat structure where:

- **{domain}.go**: Package documentation and shared constants
- **entities.go**: Domain models, often embedding API types with additional fields
- **service.go**: Business logic, defines Repository interface (consumer defines interface pattern)
- **repository.go**: Database implementation of Repository interface
- **handler.go**: HTTP handlers satisfying generated OpenAPI interfaces
- **middleware.go**: Domain-specific middleware (optional - auth only)
- **converters/**: Generated converters between DB and API types

### Current Domains

1. **auth**: Authentication and user management (users, sessions, accounts)
2. **organizations**: Organization, membership, and invitation management
3. **workflows**: Pipeline workflows, runs, and tools (formerly pipelines)
4. **content**: Content artifacts and labels (formerly knowledge)

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

## Code Generators

ArchesAI uses multiple code generators to reduce boilerplate and ensure type safety:

### 1. sqlc (Database → Go)

- **Config**: `internal/infrastructure/database/sqlc.yaml`
- **Input**: SQL queries in `internal/infrastructure/database/queries/*.sql`
- **Output**: Type-safe query functions in `internal/generated/database/postgresql/`
- **Usage**: Access via `container.Queries` in services

### 2. oapi-codegen (OpenAPI → Go)

- **Config**: `internal/generated/api/generate.go`
- **Input**: `api/openapi.bundled.yaml`
- **Output**: Server interfaces and types in `internal/generated/api/`
- **Usage**: Implement interfaces in domain handlers

### 3. generate-defaults (OpenAPI → Go Config)

- **Source**: `cmd/generate-defaults/main.go`
- **Input**: OpenAPI schema definitions
- **Output**: `internal/infrastructure/config/defaults.gen.go`
- **Purpose**: Generate config struct with default values from OpenAPI

### 4. generate-converters (YAML → Go Converters)

- **Source**: `cmd/generate-converters/main.go`
- **Config**: `internal/domains/converters.yaml`
- **Output**: `internal/domains/*/converters/converters.gen.go`
- **Features**:
  - Automap: Automatically maps fields with matching names
  - Type-aware conversions (nullable handling, UUID parsing)
  - Deterministic output (alphabetically sorted fields)

## Important Notes

- **Generated Code**: Never edit files in `internal/generated/` or `*/converters.gen.go` - they are overwritten
- **OpenAPI First**: API changes start in OpenAPI spec, not code
- **Type Safety**: Both Go (sqlc) and TypeScript (orval) use code generation for type safety
- **Monorepo**: Use pnpm workspaces - dependencies are shared via catalog in `pnpm-workspace.yaml`
- **Domain Boundaries**: Keep domains isolated - communicate through interfaces
- **Migration Safety**: Always review migrations before applying to production
- **Interface Pattern**: Consumer defines interface (e.g., Service defines Repository interface)
- Always remember to run make generate and make lint before you sign off on a task
