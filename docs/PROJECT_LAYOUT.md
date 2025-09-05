# Project Layout & Domain Mapping

This document provides a complete map of the ArchesAI project structure, including which database tables belong to which domains and the standard files for each domain.

## Complete Project Structure

```
archesai/
├── api/                              # OpenAPI specifications
│   ├── components/                   # Reusable OpenAPI components
│   │   ├── parameters/              # Common parameters
│   │   ├── responses/               # Common responses
│   │   └── schemas/                 # Data models
│   ├── paths/                       # API endpoints by resource
│   ├── openapi.yaml                 # Main OpenAPI file
│   ├── openapi.bundled.yaml         # Generated bundled spec
│   └── redocly.yaml                 # Redocly configuration
│
├── cmd/                              # Application entry points
│   ├── api/                         # Main API server
│   │   └── main.go
│   ├── generate-converters/         # Converter generator
│   │   └── main.go
│   └── generate-defaults/           # Config defaults generator
│       └── main.go
│
├── internal/                         # Private application code
│   ├── app/                         # Application assembly
│   │   └── deps.go                  # Dependency injection container
│   │
│   ├── domains/                     # Business domains (flat structure)
│   │   ├── auth/                    # Authentication & user management
│   │   │   ├── auth.go              # Package docs & constants
│   │   │   ├── entities.go          # Domain models
│   │   │   ├── service.go           # Business logic & Repository interface
│   │   │   ├── repository.go        # PostgreSQL implementation
│   │   │   ├── handler.go           # HTTP handlers
│   │   │   ├── middleware.go        # Auth middleware
│   │   │   └── converters/          # Generated converters
│   │   │       └── converters.gen.go
│   │   │
│   │   ├── organizations/           # Organization & membership management
│   │   │   ├── organizations.go     # Package docs & constants
│   │   │   ├── entities.go          # Domain models
│   │   │   ├── service.go           # Business logic & Repository interface
│   │   │   ├── repository.go        # PostgreSQL implementation
│   │   │   ├── handler.go           # HTTP handlers
│   │   │   └── converters/          # Generated converters
│   │   │       └── converters.gen.go
│   │   │
│   │   ├── workflows/               # Pipeline workflows, runs, and tools
│   │   │   ├── workflows.go         # Package docs & constants
│   │   │   ├── entities.go          # Domain models
│   │   │   ├── service.go           # Business logic & Repository interface
│   │   │   ├── repository.go        # PostgreSQL implementation
│   │   │   ├── handler.go           # HTTP handlers
│   │   │   └── converters/          # Generated converters
│   │   │       └── converters.gen.go
│   │   │
│   │   ├── content/                 # Content artifacts and labels
│   │   │   ├── content.go           # Package docs & constants
│   │   │   ├── entities.go          # Domain models
│   │   │   ├── service.go           # Business logic & Repository interface
│   │   │   ├── repository.go        # PostgreSQL implementation
│   │   │   ├── handler.go           # HTTP handlers
│   │   │   └── converters/          # Generated converters
│   │   │       └── converters.gen.go
│   │   │
│   │   └── converters.yaml          # Converter configuration
│   │
│   ├── generated/                   # Generated code (DO NOT EDIT)
│   │   ├── api/                     # OpenAPI generated
│   │   │   └── api.gen.go
│   │   └── database/                # sqlc generated
│   │       └── postgresql/
│   │           ├── models.go
│   │           ├── querier.go
│   │           └── *.sql.go
│   │
│   └── infrastructure/              # Technical infrastructure
│       ├── config/                  # Configuration management
│       │   ├── config.go
│       │   └── defaults.gen.go
│       ├── database/                # Database setup
│       │   ├── migrations/          # Schema migrations
│       │   │   ├── postgresql/
│       │   │   └── sqlite/
│       │   ├── queries/             # SQL queries for sqlc
│       │   └── sqlc.yaml            # sqlc configuration
│       └── server/                  # HTTP server setup
│           └── server.go
│
├── web/                              # Frontend code
│   ├── platform/                    # Main web application
│   ├── client/                      # Generated API client
│   ├── ui/                          # Shared UI components
│
├── docs/                             # Documentation
│   ├── DEVELOPMENT.md
│   ├── GENERATORS.md
│   └── PROJECT_LAYOUT.md
│
├── .claude/                          # Claude Code configuration
│   └── CLAUDE.md
│
├── Makefile                          # Build automation
├── go.mod                            # Go module definition
├── go.sum                            # Go dependencies
├── package.json                      # Node.js scripts
├── pnpm-workspace.yaml              # pnpm monorepo config
└── README.md                         # Project README
```

## Database Tables to Domain Mapping

### Auth Domain

Handles authentication, authorization, and user management.

**Tables:**

- `user` - User accounts
- `account` - OAuth/social login accounts
- `session` - Active user sessions
- `verification_token` - Email verification tokens
- `api_token` - API authentication tokens

**API Endpoints:**

- `/auth/sign-in`
- `/auth/sign-up`
- `/auth/sign-out`
- `/auth/sessions`
- `/auth/email-verification/*`
- `/auth/password-reset/*`

### Organizations Domain

Manages organizations, members, and invitations.

**Tables:**

- `organization` - Organization entities
- `member` - Organization members
- `invitation` - Pending invitations

**API Endpoints:**

- `/organizations`
- `/organizations/{id}`
- `/organizations/{id}/members`
- `/organizations/{id}/members/{memberId}`
- `/organizations/{id}/invitations`
- `/organizations/{id}/invitations/{invitationId}`

### Workflows Domain (formerly Pipelines)

Manages data processing pipelines, execution runs, and tools.

**Tables:**

- `pipeline` - Pipeline definitions
- `tool` - Processing tools
- `pipeline_step` - Pipeline execution steps
- `pipeline_step_to_dependency` - Step dependencies
- `run` - Pipeline execution runs

**API Endpoints:**

- `/workflows/pipelines`
- `/workflows/pipelines/{id}`
- `/workflows/runs`
- `/workflows/runs/{id}`
- `/workflows/tools`
- `/workflows/tools/{id}`

### Content Domain (formerly Knowledge)

Manages content artifacts and labels.

**Tables:**

- `artifact` - Generated artifacts
- `run_to_artifact` - Run-artifact relationships
- `label` - Artifact labels
- `label_to_artifact` - Label associations

**API Endpoints:**

- `/content/artifacts`
- `/content/artifacts/{id}`
- `/content/labels`
- `/content/labels/{id}`

### System Endpoints

Health and configuration endpoints (not domain-specific).

**API Endpoints:**

- `/health`
- `/config`

## Standard Domain Files

Each domain follows this standard structure:

### Required Files

1. **{domain}.go** - Package documentation and shared constants

   ```go
   // Package {domain} provides {description}
   package {domain}

   type ContextKey string
   const (
       // Domain-specific constants
   )
   ```

2. **entities.go** - Domain models

   ```go
   type User struct {
       api.UserEntity
       // Additional domain fields
   }

   // Domain errors
   var (
       ErrNotFound = errors.New("not found")
   )
   ```

3. **service.go** - Business logic

   ```go
   type Repository interface {
       // Data access methods
   }

   type Service struct {
       repo Repository
       // Other dependencies
   }
   ```

4. **repository.go** - Database implementation

   ```go
   var _ Repository = (*PostgresRepository)(nil)

   type PostgresRepository struct {
       q postgresql.Querier
   }
   ```

5. **handler.go** - HTTP handlers

   ```go
   type Handler struct {
       service *Service
   }

   // Implements OpenAPI ServerInterface
   ```

### Optional Files

6. **middleware.go** - Domain-specific middleware
7. **validators.go** - Input validation logic
8. **events.go** - Domain events (if using event-driven)

### Generated Files

9. **converters/converters.gen.go** - Type converters (DO NOT EDIT)

## SQL Query Files

Located in `internal/infrastructure/database/queries/`:

**Auth Domain:**

- `users.sql` - User queries
- `accounts.sql` - OAuth account queries
- `sessions.sql` - Session management
- `api-tokens.sql` - API token queries
- `verification-tokens.sql` - Verification token queries

**Organizations Domain:**

- `organizations.sql` - Organization queries
- `members.sql` - Member queries
- `invitations.sql` - Invitation queries

**Workflows Domain:**

- `pipelines.sql` - Pipeline queries
- `tools.sql` - Tool queries
- `pipeline-steps.sql` - Pipeline step queries
- `runs.sql` - Execution run queries

**Content Domain:**

- `artifacts.sql` - Artifact queries
- `labels.sql` - Label queries

## Configuration Files

### Go Configuration

- `go.mod` - Module definition
- `go.sum` - Dependency checksums
- `.golangci.yml` - Linter configuration

### Generator Configuration

- `internal/infrastructure/database/sqlc.yaml` - sqlc config
- `internal/generated/api/oapi-codegen.yaml` - OpenAPI generator config
- `internal/domains/converters.yaml` - Converter definitions

### Frontend Configuration

- `package.json` - Root package scripts
- `pnpm-workspace.yaml` - Workspace configuration
- `nx.json` - Nx configuration
- `tsconfig.json` - TypeScript configuration

### API Configuration

- `api/openapi.yaml` - Main OpenAPI spec
- `api/redocly.yaml` - Redocly linting rules

## Environment Files

- `.env` - Local development environment
- `.env.example` - Example environment variables
- `config.yaml` - Default configuration (optional)

## Build & Deployment

- `Makefile` - Build automation
- `Dockerfile` - Container image (if exists)
- `docker-compose.yml` - Local development stack (if exists)
- `helm/` - Kubernetes Helm charts (if exists)

## Development Workflow

1. **Define API** in `api/paths/` and `api/components/schemas/`
2. **Create Migration** with `make migrate-create name=feature`
3. **Write SQL Queries** in `internal/infrastructure/database/queries/`
4. **Run Generators** with `make generate`
5. **Implement Domain Logic** in appropriate domain directory
6. **Add Converter Config** to `internal/domains/converters.yaml`
7. **Wire Dependencies** in `internal/app/deps.go`
8. **Test Implementation** with `make test`

## Import Paths

Standard import organization:

```go
import (
    // Standard library
    "context"
    "errors"

    // External packages
    "github.com/labstack/echo/v4"

    // Internal packages
    "github.com/archesai/archesai/internal/generated/api"
    "github.com/archesai/archesai/internal/generated/database/postgresql"
    "github.com/archesai/archesai/internal/domains/auth"
)
```

## Naming Conventions

- **Packages**: Lowercase, singular (e.g., `auth`, not `auths`)
- **Files**: Lowercase with underscores for multi-word (e.g., `api_tokens.sql`)
- **Types**: PascalCase (e.g., `UserEntity`)
- **Functions**: PascalCase for exported, camelCase for private
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE for groups
- **Interfaces**: Usually end with `-er` (e.g., `Querier`, `Repository`)
