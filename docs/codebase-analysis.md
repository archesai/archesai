# Arches Codebase Analysis Report

Based on my thorough exploration of the Arches codebase, here is a comprehensive analysis covering all your requested areas:

---

## 1. CURRENT ARCHITECTURE ANALYSIS

### Purpose & Functionality

Arches is an **AI-powered data processing platform** with workflow automation and beautiful terminal interfaces. It provides:

- REST API for multi-provider AI integration (OpenAI, Claude, Gemini, Ollama)
- Terminal UI for configuration viewing and chat interfaces
- Workflow automation with DAG-based pipeline support
- Multi-tenant organization management with auth
- Executor system for custom code execution in isolated containers

### Key Components & Responsibilities

**Backend Architecture (Hexagonal/Clean Architecture):**

- **Adapters Layer** (`/internal/adapters/`):
  - HTTP controllers for REST endpoints
  - CLI interface (Cobra-based)
  - TUI (Terminal UI) using Bubble Tea
  - WebSocket support for real-time communication

- **Application Layer** (`/internal/application/`):
  - CQRS pattern: `/commands` (write operations) and `/queries` (read operations)
  - Commands organized by domain: auth, organization, pipeline, run, user, executor, etc.
  - Queries for fetching data with filtering/pagination

- **Core Domain Layer** (`/internal/core/`):
  - Entities and Value Objects (combined in `/models`)
  - Repositories (interfaces for data access)
  - Events (domain events for entities)
  - Services (domain services)

- **Infrastructure Layer** (`/internal/infrastructure/`):
  - **Persistence**: PostgreSQL and SQLite support with SQLC-generated queries
  - **Auth**: OAuth, JWT, magic links, email verification
  - **Executor**: Docker-based code execution in containers (Python, Go, Node.js)
  - **Cache**: Redis and in-memory caching
  - **LLM**: Multi-provider LLM interface
  - **Notifications**: Email/OTP delivery
  - **Bootstrap**: Dependency injection wiring

### API Generation Approach

The system uses a **code-first with OpenAPI as source of truth** pattern:

1. OpenAPI specification defined in `/api/` with custom x-codegen extensions
2. Generate all Go code from OpenAPI specs
3. Controllers, handlers, models all generated from single source

---

## 2. CODE GENERATION CAPABILITIES

### How Codegen Works

**Tool Entry Points:**

- `cmd/codegen/main.go` - CLI tool with three subcommands:
  1. `openapi` - Generate Go code from OpenAPI specs
  2. `jsonschema` - Generate Go structs from JSON Schema
  3. `bundle` - Bundle OpenAPI specs with external references into single file

**Generator Architecture** (`/internal/codegen/`):

- `generate.go` - Main orchestrator, runs 8+ generators in parallel:
  - GenerateSchemas
  - GenerateRepositories
  - GenerateEvents
  - GenerateCommandQueryHandlers
  - GenerateControllers
  - GenerateHCL (database migrations)
  - GenerateSQLC (SQL query files)
  - GenerateJSClient (TypeScript client)
  - GenerateBootstrap (dependency injection - sequential)

**Key Files:**

- `generate_controllers.go` - HTTP handler generation
- `generate_schemas.go` - Model/DTO generation
- `generate_cqrs.go` - Command/Query handler generation
- `generate_repositories.go` - Repository interface generation
- `filewriter.go` - Writes generated files with overwrite protection

### Templates (in `/internal/codegen/tmpl/`)

| Template                   | Purpose                                  | Size      |
| -------------------------- | ---------------------------------------- | --------- |
| `schema.tmpl`              | Entity/Value Object models               | 352 lines |
| `controller.tmpl`          | HTTP endpoint handlers                   | 503 lines |
| `command_handler.tmpl`     | Command handlers (POST/PUT/PATCH/DELETE) | 198 lines |
| `query_handler.tmpl`       | Query handlers (GET)                     | 96 lines  |
| `repository.tmpl`          | Repository interface                     | 33 lines  |
| `repository_postgres.tmpl` | PostgreSQL implementation                | 230 lines |
| `repository_sqlite.tmpl`   | SQLite implementation                    | 80 lines  |
| `bootstrap.tmpl`           | Dependency injection wiring              | 341 lines |
| `infrastructure.tmpl`      | Infrastructure setup                     | 196 lines |
| `events.tmpl`              | Domain events                            | 64 lines  |
| `schema_hcl.tmpl`          | HCL for Atlas migrations                 | 109 lines |

### What Gets Generated

**From OpenAPI Specs:**

1. **Models** - Entities and Value Objects from component schemas
   - Full struct definitions with JSON/YAML tags
   - Enum types with Parse/IsValid methods
   - Constructor functions (New\*)

2. **HTTP Controllers** - Grouped by domain tag
   - Request/Response type definitions
   - HTTP route registration
   - Parameter binding (path, query, headers, body)

3. **CQRS Handlers** - Automatically routed based on HTTP method
   - Commands for POST/PUT/PATCH/DELETE
   - Queries for GET
   - Built with repository injection

4. **Repositories** - Data access layer
   - Interface definitions
   - PostgreSQL implementations with SQLC
   - SQLite implementations
   - Support for indices and relations via x-codegen extensions

5. **Database Migrations**
   - HCL files for Atlas schema management
   - SQL files for migrations

6. **TypeScript Client** - Using Orval
   - Generated in `web/client/src/generated/`
   - Includes Zod schemas for validation

7. **Bootstrap File** - Dependency injection container
   - Sets up all repositories
   - Wires controllers and handlers
   - Initializes auth, cache, notifications

### Controller Generation Pattern

Controllers are **REST ordered**: POST → GET (singular) → GET (list) → PATCH/PUT → DELETE

**Type System:**

- `ControllersTemplateData` contains Tag and Operations
- Operations sorted by method and ID pattern
- Generates handler functions with request/response interfaces
- Response types support visitor pattern

---

## 3. RUNTIME CAPABILITIES

### Executor System (`/internal/infrastructure/executor/`)

**What It Does:**

- Executes user-provided code in isolated Docker containers
- Supports 3 languages: Python, Go, Node.js
- Input/Output schema validation
- Resource limits (CPU, memory)
- Async execution with timeouts

**Components:**

- `executor.go` - Interface definition for generic typed execution
- `executor_service.go` - Main service (ExecutorService[A,B])
- `container.go` - Docker container management
- `builder.go` - Image building
- `local.go` - Local execution (for testing)

**Execution Flow:**

1. Get executor config from database
2. Create temp directory for code
3. Write execute code + extra files (mounts)
4. Build/get Docker image with dependencies
5. Run container with mounted code
6. Return output + execution time

**Database Model:**

- `Executor` entity stores:
  - ExecuteCode (source code)
  - Language (python/go/node)
  - SchemaIn/SchemaOut (JSON Schema validation)
  - ExtraFiles (mounted as read-only)
  - Dependencies (package.json/requirements.txt)
  - Timeout, MemoryMB, CPUShares
  - Environment variables

### Runner Containers

**Base images** in `/deployments/containers/runners/`:

- **Node.js** - archesai/runner-node:latest
- **Python** - archesai/runner-python:latest
- **Go** - archesai/runner-go:latest

Each runner has:

- Dockerfile with base runtime
- Entry point that executes user code
- Standard input/output handling

### Pipeline & Workflow System

**Core Models:**

- `Pipeline` - Workflow definition with steps
- `PipelineStep` - Individual steps with dependencies
- `Run` - Execution instance of a pipeline
- `Artifact` - Content/file storage
- `Tool` - Reusable tools in workflows

**Key Features:**

- DAG-based step execution
- Step dependencies
- Artifact storage and retrieval
- Run tracking with status

---

## 4. DEVELOPMENT EXPERIENCE

### Hot Reload Setup (`.air.toml`)

- **Build command:** `go build -o ./tmp/archesai ./cmd/archesai/main.go`
- **Excluded dirs:** assets, api, node_modules, deployments, docs
- **Included extensions:** go, tpl, tmpl, html
- **Watch delay:** 1000ms
- **Runs with:** `CI=1 CLICOLOR_FORCE=1 ./tmp/archesai api`

### Makefile Targets (80+ targets)

**Code Generation:**

- `make generate` - Run all generators (codegen + mocks)
- `make g` - Shortcut for generate
- `make bundle-openapi` - Bundle OpenAPI specs

**Development:**

- `make dev-api` - API with hot reload (air)
- `make dev-all` - All services (API, platform, docs)
- `make w` - Shortcut for dev-all

**Building:**

- `make build` - Build all binaries
- `make build-api` - Server binary
- `make build-platform` - Frontend assets
- `make build-runners` - Container images

**Testing:**

- `make test` - Run all tests
- `make test-short` - Skip integration tests
- `make test-coverage` - Generate coverage report
- `make test-watch` - Watch mode (requires fswatch)

**Linting & Formatting:**

- `make lint` - All linters (Go, TypeScript, OpenAPI, Docs)
- `make format` - Format all code
- `make lint-go` - Go linter (golangci-lint)
- `make lint-openapi` - OpenAPI validation

**Deployment:**

- `make docker-run` - Docker Compose
- `make skaffold-dev` - Kubernetes dev mode

### Build Process

**CLI Architecture:**

- Cobra-based CLI (`/internal/adapters/cli/`)
- Main entry: `cmd/archesai/main.go` → `cli.Execute()`
- Commands: api, worker, tui, web, config, completion, version

**Subcommands:**

- `archesai api` - REST server
- `archesai tui` - Terminal UI
- `archesai worker` - Background job processor
- `archesai web` - Platform UI
- `archesai config show` - Configuration viewer

---

## 5. FRONTEND COMPONENTS

### Platform UI (`/web/platform/`)

- **Framework:** React 19 + TypeScript
- **Routing:** TanStack Router (TanStart SSR)
- **UI Components:** Custom component library (@archesai/ui)
- **State Management:** TanStack Query (React Query)
- **Build:** Vite

**Structure:**

- `/src/components/` - Reusable components
- `/src/hooks/` - Custom React hooks (5 files)
- `/src/lib/` - Utilities
- `/src/routes/` - Route definitions
- `/src/styles/` - CSS modules

**Key Dependencies:**

- @tanstack/react-router - Routing
- @tanstack/react-query - Data fetching
- @tanstack/react-table - Tables
- @xyflow/react - Workflow/DAG visualization
- motion - Animations
- socket.io-client - Real-time communication

### Client SDK (`/web/client/`)

- **Purpose:** Type-safe API client for React apps
- **Generation:** Orval (from OpenAPI specs)
- **Validation:** Zod for runtime schema validation
- **Output:** `src/generated/`

**Generated Files:**

- `orval.schemas.ts` - Type-safe models
- `zod.ts` - Zod validation schemas

**Build:** TypeScript with tsconfig alias imports

### UI Component Library (`/web/ui/`)

- Shared reusable components
- Custom hooks
- Layout utilities
- Provider components (themes, etc.)
- Type definitions

---

## 6. DATABASE & PERSISTENCE

### Approach: SQLC + Atlas Migrations

**SQL Generation:**

- **SQLC** (`/internal/infrastructure/persistence/sqlc.yaml`)
- Generates type-safe Go code from SQL
- Files in `/postgres/queries/` and `/sqlite/queries/`
- Generated code in `/postgres/repositories/` and `/sqlite/repositories/`

**Migrations:**

- **Atlas** for schema management
- HCL schema files: `schema.gen.hcl`
- SQL migration files with timestamps
- Automatic generation from codegen

**Supported Databases:**

- PostgreSQL 15+ (primary)
- SQLite (development/testing)

**Schema:**

- Tables for each domain entity
- Foreign keys and indices from x-codegen
- Timestamps (created_at, updated_at)
- UUID primary keys

**Sample Tables:**

- users, accounts (auth)
- organizations, members, invitations
- pipelines, pipeline_steps, runs
- artifacts, labels, tools
- executors, api-keys, sessions

---

## 7. AUTHENTICATION & AUTHORIZATION

### Auth Infrastructure (`/internal/infrastructure/auth/`)

**Features:**

- **JWT tokens** - Bearer token auth
- **OAuth** - Pluggable OAuth providers
- **Magic Links** - Email-based auth
- **Sessions** - Cookie-based sessions
- **Email Verification** - OTP delivery
- **Password Reset** - Reset flow

**Flows:**

- Login/Logout
- Register
- Magic link request/verify
- Email verification
- Password reset/change
- OAuth provider authorize/callback
- Account linking

**Service:**

- `Service` implements `services.AuthService`
- Token manager for JWT handling
- Magic link provider
- OAuth provider registry
- Cache integration (Redis/Memory)

**Security:**

- Bcrypt password hashing
- JWT with expiration
- CSRF protection via session tokens
- Rate limiting middleware
- CORS middleware

---

## 8. OPENAPI HANDLING

### Structure

**Root:** `/api/openapi.yaml` (with includes)

**Organization:**

- `/paths/` - 47 endpoint definitions
- `/components/schemas/` - Reusable schemas
  - 22 base schemas
  - `/config/` subdirectory (37 config-related schemas)
- `/components/parameters/` - 29 parameter definitions
  - Filter, Sort, Pagination
  - ResourceID, OrganizationID
- `/components/responses/` - 42 response definitions
- `/components/headers/` - 5 header definitions

**Bundling:**

- Uses pb33f/libopenapi for parsing
- `codegen bundle` merges all references
- Orval fix option for pathItem resolution
- Single bundled output: `api/openapi.bundled.yaml`

### X-Codegen Extensions

Custom extensions on operations and schemas:

**On Operations:**

- `x-codegen-custom-handler` - Skip auto-generation
- `x-codegen-repository` - Custom repo name

**On Schemas:**

- `x-codegen-schema-type` - "entity" or "valueobject"
- `x-codegen-repository` - Repository config
  - `indices` - Database indices
  - `relations` - Foreign key relations

**Parsing:**

- `XCodegenExtension` type
- `xcodegenextension.go` parser
- Extracted during OpenAPI parsing

---

## 9. KEY ARCHITECTURAL PATTERNS

### CQRS (Command Query Responsibility Segregation)

- **Commands** - Modify state (POST/PUT/PATCH/DELETE)
- **Queries** - Read state (GET)
- Separate handler files per domain
- Generated with full type safety

### Hexagonal Architecture

- Clear adapter/application/domain separation
- Ports (interfaces) and adapters (implementations)
- Easy to test and swap components

### Repository Pattern

- Database abstraction through repositories
- Support for multiple databases (PostgreSQL, SQLite)
- SQLC for type-safe queries

### Dependency Injection

- Bootstrap container wires all dependencies
- Configuration-driven service creation
- Testable through injection

### Domain Events

- Entities can publish events
- Event publisher interface
- Domain event types generated

---

## 10. TESTING INFRASTRUCTURE

**Approach:**

- Mockery v3 for mocks (`go tool -modfile=tools.mod mockery`)
- `.mockery.yaml` configuration
- `mocks_test.go` in packages
- TestContainers for integration tests
- Docker for database/Redis testing

**Test Utilities:**

- `/internal/shared/testutil/`
- Postgres and Redis test containers
- Test data fixtures in `/test/data/`

---

## KEY INSIGHTS

1. **Code-First Development:** OpenAPI specs are the single source of truth, all code generated from them
2. **Parallel Generation:** 8+ generators run in parallel for fast builds
3. **Type Safety:** Generated types at every layer (models, handlers, repositories)
4. **Container-Native:** Executor system uses Docker for isolated code execution
5. **Multi-Stack:** Go backend + TypeScript/React frontend + SQLite/PostgreSQL databases
6. **Developer Experience:** Hot reload, comprehensive Makefile, code generation automation
7. **Recent Refactor:** Moved to hexagonal architecture (commit 88866aab) with combined entities/valueobjects (cb5b0c9a)
8. **Enterprise-Ready:** Multi-tenant, OAuth, RBAC planning, monitoring support

---

## PROJECT FILES SUMMARY

**Total Codegen Templates:** 2,202 lines across 11 templates
**Parsers:** OpenAPI + JSONSchema + x-codegen extension parsing
**Database:** 14+ domain entities with PostgreSQL and SQLite support
**API Endpoints:** 47+ endpoints across 18 tags
**Frontend:** React 19 with TanStack ecosystem
**CLI:** 8 commands via Cobra
**Development:** Vite (frontend), Air (backend hot reload), Atlas migrations

---

## ABSOLUTE FILE PATHS FOR KEY COMPONENTS

- `/home/jonathan/apis/archesai/internal/codegen/` - Code generation engine
- `/home/jonathan/apis/archesai/internal/infrastructure/executor/` - Custom execution system
- `/home/jonathan/apis/archesai/internal/adapters/http/` - HTTP layer
- `/home/jonathan/apis/archesai/internal/infrastructure/persistence/` - Database layer
- `/home/jonathan/apis/archesai/web/platform/` - Frontend UI
- `/home/jonathan/apis/archesai/api/` - OpenAPI specification
- `/home/jonathan/apis/archesai/cmd/` - CLI entry points

---

## TRANSFORMATION ASSESSMENT FOR APP BUILDER

### What We Can Keep (70% of codebase)

#### Core Strengths

1. **Code Generation Engine** - Already generates full backend from OpenAPI
2. **Executor System** - Perfect for custom handler execution
3. **Template System** - Extensible for multiple languages/frameworks
4. **OpenAPI Infrastructure** - Parsing, bundling, x-codegen extensions
5. **Authentication System** - Complete auth flows ready to use
6. **Database Layer** - Multi-DB support with migrations
7. **Development Tools** - Hot reload, Makefile, testing infrastructure

### What Needs Enhancement (20% modification)

#### Required Additions

1. **Frontend Generation** - Add React component generation from OpenAPI
2. **WebSocket Support** - Add real-time templates
3. **RBAC Generation** - Permission system from OpenAPI security schemes
4. **Multi-Backend** - Templates for Python/Node.js backends
5. **Visual Designer** - Web UI for schema creation
6. **AI Integration** - LLM-powered generation endpoints

### What to Replace (10% removal)

#### Current Domain-Specific Code

1. **Pipeline/Run System** - Replace with app build pipeline
2. **Organization Management** - Simplify for app projects
3. **Current API Endpoints** - Replace with studio API
4. **Worker System** - Unless needed for async builds

---

## RECOMMENDED ARCHITECTURE FOR APP BUILDER

### Three-Layer Architecture

#### Layer 1: Studio (Development Environment)

- Web-based IDE for app creation
- Visual OpenAPI designer
- AI chat for schema generation
- Project management
- Deployment controls

#### Layer 2: Engine (Generation & Runtime)

- Enhanced codegen from OpenAPI
- Multi-language template system
- Hot reload development server
- Build and packaging system
- Deployment generators

#### Layer 3: Generated Apps (Output)

- Self-contained applications
- Full-stack with frontend/backend/database
- Docker/K8s ready
- Single binary option
- Cloud-deployable

### Development Workflow

1. **Design Phase**
   - User creates OpenAPI schema (visual/AI/code)
   - Defines custom handler logic
   - Configures app settings

2. **Generation Phase**
   - Engine generates full application
   - Creates CRUD operations
   - Generates frontend components
   - Sets up authentication

3. **Customization Phase**
   - User implements business logic
   - AI assists with handler code
   - Tests in development environment

4. **Deployment Phase**
   - Build for target platform
   - Generate deployment configs
   - One-click deploy

---

## CONCLUSION

Arches has an exceptional foundation for transformation into an AI-powered app builder. The existing code generation infrastructure, executor system, and development tooling provide ~70% of what's needed. The main work involves:

1. **Extending the template system** for full-stack generation
2. **Building the Studio interface** for visual development
3. **Integrating AI capabilities** for natural language → code
4. **Creating the deployment pipeline** for multiple targets

The transformation is not only feasible but can leverage most of the existing high-quality codebase, making it an evolution rather than a rewrite.
