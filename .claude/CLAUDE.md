# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ArchesAI is a comprehensive data processing platform with a hybrid architecture:

- **Backend**: Go API server using Echo framework with Hexagonal Architecture (Ports & Adapters)
- **Frontend**: TypeScript/React with TanStack Router, built with Vite
- **Database**: PostgreSQL with vector extensions for embeddings
- **Monorepo**: pnpm workspaces for TypeScript, standalone Go module

## Essential Development Commands

### Backend (Go)

```bash
# Development with hot reload
make dev  # or make watch

# Build and run
make build            # Build all binaries
make build-archesai   # Build server binary only
make build-codegen    # Build codegen tool only
make run             # Run the API server
make run-api         # Run API server (alias)
make run-web         # Run web UI server
make run-worker      # Run background worker

# Code generation (required after API/DB schema changes)
make generate              # Runs all generators
make generate-sqlc         # Generate database code from SQL queries
make generate-oapi         # Generate server code from OpenAPI spec
make generate-defaults     # Generate config defaults from OpenAPI
make generate-adapters     # Generate type converters between layers
make generate-domain       # Scaffold new domain (usage: make generate-domain name=billing tables=subscription,invoice)

# Testing and quality
make test             # Run all tests
make test-unit        # Run unit tests only
make test-integration # Run integration tests only
make test-coverage    # Generate coverage report
make lint            # Run all linters (Go + OpenAPI)
make lint-go        # Run Go linter only
make lint-oapi      # Run OpenAPI linter only
make format         # Format all code
make format-go      # Format Go code only
make format-sql     # Format SQL files

# Database migrations
make migrate-up      # Apply all pending migrations
make migrate-down    # Rollback last migration
make migrate-create name=<migration_name>  # Create new migration
make migrate-status  # Show migration status
make migrate-force version=<version>  # Force set migration version

# Docker operations
make docker-build    # Build Docker image
make docker-run      # Run in Docker container
make docker-down     # Stop Docker containers
make docker-logs     # View Docker logs

# Utilities
make clean          # Clean build artifacts
make deps           # Download dependencies
make tools          # Install development tools
make help           # Show all available commands
```

### Key Patterns and Conventions

#### Hexagonal Architecture Flow

1. **HTTP Request** → Handler (Adapter)
2. **Handler** → Use Case (Core)
3. **Use Case** → Repository Interface (Port)
4. **Repository** → Database (Infrastructure Adapter)
5. **Response** flows back through the same layers

#### API Development Flow

1. Define endpoints in `api/openapi.yaml` or component files
2. Run `make generate-oapi` to generate server interfaces
3. Implement handlers in `internal/domains/*/handlers/http/`
4. Handlers must satisfy generated interfaces from `internal/domains/*/generated/api/`

#### Database Development Flow

1. Create migration: `make migrate-create name=add_users_table`
2. Write SQL queries in `internal/infrastructure/database/queries/`
3. Run `make generate-sqlc` to generate type-safe query functions
4. Use generated queries in repository implementations

#### Type Converter Flow

1. Define converters in `internal/domains/adapters.yaml`
2. Run `make generate-adapters` to generate converter functions
3. Use converters in handlers and repositories for type mapping

#### Frontend Development Flow

1. API client is auto-generated from OpenAPI spec
2. Use `@archesai/client` package to make API calls
3. Components go in `@archesai/ui` for reusability
4. Routes are file-based in `web/platform/src/routes/`

#### Dependency Injection Pattern

- All dependencies are wired in `internal/app/deps.go`
- Container pattern provides dependencies to all layers
- Services receive repositories via interfaces (for testability)
- Handlers receive services for business logic

#### Authentication Flow

- JWT-based authentication with refresh tokens
- Auth middleware validates tokens and adds claims to context
- Protected routes use `middleware.RequireAuth()`
- Session management with database-backed sessions

## Environment Configuration

Backend configuration uses Viper and reads from:

- Environment variables (prefix: `ARCHESAI_`)
- `.env` file (local development)
- `config.yaml` (defaults)

Key environment variables:

```bash
# Database
ARCHESAI_DATABASE_URL=postgres://user:pass@localhost/archesai?sslmode=disable
ARCHESAI_DATABASE_POOL_SIZE=10
ARCHESAI_DATABASE_MAX_IDLE_TIME=30m

# Server
ARCHESAI_SERVER_PORT=8080
ARCHESAI_SERVER_HOST=0.0.0.0
ARCHESAI_SERVER_READ_TIMEOUT=30s
ARCHESAI_SERVER_WRITE_TIMEOUT=30s

# Authentication
ARCHESAI_JWT_SECRET=your-secret-key
ARCHESAI_JWT_ACCESS_TOKEN_DURATION=15m
ARCHESAI_JWT_REFRESH_TOKEN_DURATION=7d

# Logging
ARCHESAI_LOGGING_LEVEL=info
ARCHESAI_LOGGING_FORMAT=json
```

## Common Tasks

### Adding a New API Endpoint

1. Define in `api/openapi.yaml` or create new path file in `api/paths/`
2. Run `make generate` to regenerate all code
3. Implement use case in `internal/domains/{domain}/core/usecase.go`
4. Implement handler in `internal/domains/{domain}/handlers/http/handler.go`
5. Wire handler in `internal/app/routes.go`

### Adding a New Database Table

1. Create migration: `make migrate-create name=create_table_name`
2. Write SQL queries in `internal/infrastructure/database/queries/table_name.sql`
3. Run `make generate-sqlc` to generate query functions
4. Define domain entity in `core/entities.go`
5. Add repository methods to `core/ports.go`
6. Implement repository methods in `infrastructure/postgres.go`

### Creating a New Domain

1. Use the generator: `make generate-domain name=billing tables=subscription,invoice`
2. Define business logic in `core/usecase.go`
3. Add domain-specific repository methods
4. Implement custom handlers as needed
5. Wire domain in `internal/app/deps.go`
6. Add routes in `internal/app/routes.go`

### Creating a New React Component

1. Add component to `web/ui/src/components/`
2. Export from `web/ui/src/index.ts`
3. Import in platform app as `@archesai/ui`
4. Add stories for Storybook if applicable

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
- **Output**: Type-safe query functions in `internal/infrastructure/database/generated/`
- **Usage**: Access via repository implementations

### 2. oapi-codegen (OpenAPI → Go)

- **Config**: Per-domain generation in `internal/domains/*/generated/api/generate.go`
- **Input**: `api/openapi.bundled.yaml`
- **Output**: Server interfaces and types in `internal/domains/*/generated/api/`
- **Usage**: Implement interfaces in domain handlers

### 3. generate-defaults (OpenAPI → Go Config)

- **Source**: Custom generator
- **Input**: OpenAPI schema definitions
- **Output**: `internal/infrastructure/config/defaults.gen.go`
- **Purpose**: Generate config struct with default values from OpenAPI

### 4. generate-adapters (YAML → Go Converters)

- **Config**: `internal/domains/adapters.yaml`
- **Output**: `internal/domains/*/adapters/adapters.gen.go`
- **Features**:
  - Automap: Automatically maps fields with matching names
  - Type-aware conversions (nullable handling, UUID parsing)
  - Deterministic output (alphabetically sorted fields)
  - Custom field mappings via overrides

## Important Notes

- **Generated Code**: Never edit files in `generated/` directories or `*.gen.go` files - they are overwritten
- **OpenAPI First**: API changes start in OpenAPI spec, not code
- **Type Safety**: Both Go (sqlc) and TypeScript (orval) use code generation for type safety
- **Monorepo**: Use pnpm workspaces - dependencies are shared via catalog
- **Domain Boundaries**: Keep domains isolated - communicate through interfaces
- **Migration Safety**: Always review migrations before applying to production
- **Hexagonal Pattern**: Core business logic should never depend on infrastructure
- **Interface Segregation**: Define small, focused interfaces (ports) for external dependencies
- **Always remember to run `make generate` and `make lint` before you sign off on a task**
- **For code generation tools**: Use make commands (`make generate-sqlc`, `make generate-oapi`) for standard operations. The tools are installed via go.mod and can be run directly with `go tool` if needed (e.g., `go tool sqlc`)
- **Error Handling**: Implement proper error handling with descriptive messages. Use domain-specific errors defined in `domain/entities.go`
- **Testing**: Write tests for all new functionality. Use table-driven tests for Go code
- **Documentation**: Update relevant documentation when adding new features or changing existing behavior
