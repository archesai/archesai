# Project Layout

This document provides an overview of the ArchesAI project structure and organization.

## Directory Structure

```text
.
├── api/                          # OpenAPI specifications
│   ├── components/               # Reusable OpenAPI components
│   │   ├── parameters/           # Common parameters
│   │   ├── responses/            # Common responses
│   │   └── schemas/              # Data schemas
│   ├── paths/                    # API endpoint definitions
│   ├── openapi.bundled.yaml     # Bundled specification
│   └── openapi.yaml             # Main OpenAPI spec
├── cmd/
│   └── archesai/                 # Main application entry point
│       └── main.go
├── deployments/                  # Deployment configurations
│   ├── development/              # Development environment
│   ├── docker/                   # Docker configurations
│   ├── gcp/                      # Google Cloud Platform
│   ├── helm/                     # Kubernetes Helm charts
│   └── k3d/                      # Local k3d setup
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md           # System architecture
│   ├── AUTHENTICATION.md         # Authentication documentation
│   ├── CONTRIBUTING.md           # Contribution guidelines
│   ├── TEST_COVERAGE_REPORT.md   # Test coverage reports
│   ├── TESTING.md               # Testing strategy
│   └── TUI.md                   # Terminal UI documentation
├── internal/                     # Private application code
│   ├── app/                      # Application setup
│   ├── auth/                     # Authentication domain
│   ├── cli/                      # Command-line interface
│   ├── codegen/                  # Code generation system
│   ├── config/                   # Configuration management
│   ├── content/                  # Content management domain
│   ├── database/                 # Database layer
│   ├── health/                   # Health check endpoints
│   ├── llm/                      # Large Language Model integration
│   ├── logger/                   # Logging utilities
│   ├── migrations/               # Database migrations
│   ├── organizations/            # Organization management domain
│   ├── redis/                    # Redis integration
│   ├── server/                   # HTTP server setup
│   ├── storage/                  # File storage
│   ├── testutil/                 # Testing utilities
│   ├── tui/                      # Terminal user interface
│   ├── users/                    # User management domain
│   └── workflows/                # Workflow automation domain
├── scripts/                      # Build and utility scripts
├── .taskmaster/                  # Task Master AI configuration
│   ├── docs/                     # Task Master documentation
│   ├── reports/                  # Analysis reports
│   ├── tasks/                    # Task definitions
│   └── templates/                # Task templates
├── test/
│   └── data/                     # Test data files
├── tools/                        # Development tools
│   ├── codegen/                  # Code generation tool
│   └── pg-to-sqlite/             # Database conversion tool
└── web/                          # Frontend applications
    ├── client/                   # API client library
    ├── eslint/                   # ESLint configuration
    ├── platform/                 # Main platform SPA
    ├── prettier/                 # Prettier configuration
    ├── typescript/               # TypeScript configuration
    └── ui/                       # Shared UI components
```

## Key Directories Explained

### `/api`

Contains the OpenAPI 3.0 specification split into logical components:

- **components/**: Reusable schemas, parameters, and responses
- **paths/**: Individual endpoint definitions organized by domain
- **openapi.yaml**: Main specification file that references all components

### `/internal`

Private Go packages following domain-driven design:

- **Domain packages** (auth, organizations, workflows, content): Core business logic
- **Infrastructure packages** (database, redis, server): Technical implementations
- **Shared packages** (config, logger, testutil): Common utilities

### `/web`

Frontend applications and shared packages:

- **platform/**: Main React SPA using TanStack Router and Start
- **client/**: Generated TypeScript API client
- **ui/**: Shared component library
- **Config packages**: ESLint, Prettier, TypeScript configurations

### `/deployments`

All deployment-related configurations:

- **docker/**: Docker Compose for local development
- **helm/**: Kubernetes Helm charts for production
- **k3d/**: Local Kubernetes cluster setup

### Generated Files

The project uses extensive code generation:

- `*.gen.go`: Generated Go code (types, handlers, repositories)
- `types.gen.go`: OpenAPI-generated type definitions
- `http.gen.go`: HTTP handler interfaces
- `repository.gen.go`: Database repository interfaces

## Architecture Patterns

### Hexagonal Architecture

Each domain follows hexagonal architecture:

```text
domain/
├── domain.go          # Domain entities and business rules
├── service.go         # Use cases and business logic
├── handler.go         # HTTP adapter (inbound port)
├── repository.gen.go  # Database adapter (outbound port)
├── cache.gen.go       # Cache adapter (outbound port)
└── events.gen.go      # Event publisher (outbound port)
```

### Code Generation Flow

1. Define in OpenAPI (`api/`) and SQL (`internal/database/queries/`)
2. Run `make generate` to create Go types and interfaces
3. Implement business logic in service layer
4. Generated adapters handle HTTP, database, and caching

### Testing Structure

Each domain includes comprehensive tests:

- `service_test.go`: Business logic unit tests with mocks
- `handler_test.go`: HTTP handler tests
- `*_postgres_test.go`: Integration tests with real database

## File Naming Conventions

- `*.gen.go`: Generated files (do not edit manually)
- `*_test.go`: Test files
- `mocks_test.go`: Test mocks and fixtures
- `mappers.go`: Data transformation between layers
- `adapters/`: External service integrations

## Configuration Files

- `.air.toml`: Hot reload configuration
- `.golangci.yaml`: Go linter settings
- `.mockery.yaml`: Mock generation settings
- `sqlc.yaml`: Database code generation
- `*.codegen.yaml`: Custom code generation templates
