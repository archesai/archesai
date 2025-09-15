# Project Layout

## Directory Structure

```text
.
├── api/                          # OpenAPI specifications
│   ├── components/
│   │   ├── parameters/           # Reusable query/path parameters
│   │   ├── responses/            # Standard error/success responses
│   │   └── schemas/              # Data models (User, Organization, etc.)
│   ├── paths/                    # Endpoint definitions by domain
│   │   ├── auth.yaml            # /auth/* endpoints
│   │   ├── organizations.yaml   # /organizations/* endpoints
│   │   └── workflows.yaml       # /workflows/* endpoints
│   ├── openapi.yaml             # Main spec with $refs to components
│   └── openapi.bundled.yaml     # Single-file bundle (generated)
│
├── cmd/
│   └── archesai/
│       └── main.go              # Entry point, CLI commands
│
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile           # Multi-stage build
│   │   └── docker-compose.yml   # Local dev environment
│   ├── helm/                    # Kubernetes Helm charts
│   └── k3d/                     # Local k8s cluster config
│
├── docs/                        # Markdown documentation
│   ├── api-reference/
│   ├── architecture/
│   ├── deployment/
│   ├── features/
│   ├── guides/
│   ├── monitoring/
│   ├── performance/
│   ├── security/
│   └── troubleshooting/
│
├── internal/                    # Private Go packages
│   ├── app/                     # Application bootstrap, DI container
│   │   ├── app.go              # Main app struct
│   │   └── providers.go        # Service providers
│   │
│   ├── auth/                    # Authentication domain
│   │   ├── domain.go           # User, Session, Token entities
│   │   ├── service.go          # Login, Register, Verify logic
│   │   ├── handler.go          # HTTP endpoints
│   │   ├── middleware.go       # JWT validation middleware
│   │   ├── repository.gen.go   # Generated DB queries
│   │   ├── cache.gen.go        # Generated Redis cache
│   │   ├── http.gen.go         # Generated HTTP interfaces
│   │   ├── types.gen.go        # Generated OpenAPI types
│   │   ├── mappers.go          # Entity <-> DTO conversions
│   │   ├── service_test.go     # Unit tests with mocks
│   │   └── postgres_test.go    # Integration tests
│   │
│   ├── accounts/                # Account management domain
│   ├── artifacts/               # Artifact storage and management
│   ├── invitations/             # Invitation system
│   ├── labels/                  # Label management
│   ├── members/                 # Organization members
│   ├── organizations/           # Multi-tenancy domain
│   ├── pipelines/               # Pipeline execution engine
│   ├── runs/                    # Pipeline run management
│   ├── sessions/                # Session management
│   ├── tools/                   # Tool registry and integration
│   ├── users/                   # User profiles domain
│   │
│   ├── database/
│   │   ├── migrations/          # SQL migration files
│   │   │   ├── 001_initial.sql
│   │   │   └── 002_add_workflows.sql
│   │   ├── queries/             # SQLC query definitions
│   │   │   ├── auth.sql        # Auth-related queries
│   │   │   └── organizations.sql
│   │   ├── schema.sql          # Complete schema
│   │   └── db.go               # Database connection
│   │
│   ├── codegen/                # Custom code generator
│   │   ├── templates/          # Go templates
│   │   └── generator.go        # Generation logic
│   │
│   ├── cache/                   # Cache utilities
│   │
│   ├── cli/                     # CLI commands
│   │   └── commands.go         # Command definitions
│   │
│   ├── config/                  # Configuration
│   │   ├── config.go           # Env vars, validation
│   │   └── loader.go           # Config loading utilities
│   │
│   ├── health/                 # Health checks
│   │   └── handler.go          # /health/live, /health/ready
│   │
│   ├── email/                  # Email service
│   │   └── service.go          # Email sending logic
│   │
│   ├── events/                 # Event system
│   │   └── publisher.go        # Event publishing
│   │
│   ├── llm/                    # LLM integrations
│   │   ├── client.go           # Provider abstraction
│   │   ├── openai.go           # OpenAI implementation
│   │   └── anthropic.go        # Anthropic implementation
│   │
│   ├── logger/                 # Structured logging
│   ├── redis/                  # Redis client and utilities
│   ├── server/                 # HTTP server setup
│   ├── storage/                # File storage (S3, local)
│   ├── testutil/               # Test helpers
│   │   ├── containers.go       # Testcontainers setup
│   │   └── fixtures.go         # Common test data
│   └── tui/                    # Terminal UI
│       ├── app.go              # Bubble Tea app
│       └── views/              # UI components
│
├── scripts/
│   ├── generate-api-docs.sh
│   ├── setup-dev.sh
│   └── migrate.sh
│
├── test/
│   └── data/                   # Test fixtures
│       ├── sample.pdf
│       └── users.json
│
├── tools/
│   ├── codegen/                # Code generation CLI
│   └── pg-to-sqlite/           # Schema converter
│
├── web/                        # Frontend monorepo
│   ├── client/                 # Generated TypeScript client
│   │   ├── src/
│   │   │   └── generated/     # OpenAPI client code
│   │   └── package.json
│   │
│   ├── docs/                   # Zudoku documentation site
│   │   ├── apis/              # OpenAPI spec copy
│   │   ├── pages/             # MDX documentation
│   │   ├── public/            # Static assets
│   │   ├── zudoku.config.tsx  # Site configuration
│   │   └── package.json
│   │
│   ├── platform/              # Main React application
│   │   ├── app/
│   │   │   ├── routes/        # File-based routing
│   │   │   ├── components/    # React components
│   │   │   └── hooks/         # Custom hooks
│   │   ├── public/
│   │   └── package.json
│   │
│   ├── ui/                    # Shared component library
│   │   ├── src/
│   │   │   ├── Button/
│   │   │   ├── Card/
│   │   │   └── index.ts
│   │   └── package.json
│   │
│   ├── eslint/                # Shared ESLint config
│   └── typescript/            # Shared TypeScript config
│
├── .taskmaster/               # Task Master AI
│   ├── CLAUDE.md             # AI context
│   ├── config.json           # Model configuration
│   ├── docs/
│   │   └── prd.txt          # Product requirements
│   ├── reports/              # Complexity analysis
│   └── tasks/
│       ├── tasks.json        # Task database
│       └── *.md              # Individual task files
│
├── .air.toml                 # Hot reload config
├── .env.example              # Environment template
├── .gitignore
├── .golangci.yaml           # Go linter rules
├── .markdownlint.json       # Markdown linter rules
├── .mockery.yaml            # Mock generation config
├── .redocly.yaml            # OpenAPI linter rules
├── go.mod                   # Go dependencies
├── go.sum
├── LICENSE
├── Makefile                 # Build automation
├── package.json             # Root package.json
├── pnpm-lock.yaml          # pnpm lockfile
├── pnpm-workspace.yaml     # Monorepo config
├── README.md
├── sqlc.yaml               # Database code generation
└── tsconfig.json           # Root TypeScript config
```

## Domain Package Structure

Each domain in `/internal` follows this pattern:

```text
domain/
├── generate.go        # Code generation annotations
├── service.go         # Business logic, use cases
├── handler.go         # HTTP request/response handling
├── middleware.go      # Domain-specific middleware (optional)
├── repository.gen.go  # Generated database interface
├── cache.gen.go       # Generated cache layer (optional)
├── events.gen.go      # Generated event publisher (optional)
├── http.gen.go        # Generated HTTP interface from OpenAPI
├── types.gen.go       # Generated types from OpenAPI
├── mappers.go         # Convert between layers (optional)
├── service_test.go    # Unit tests with mocked dependencies
├── handler_test.go    # HTTP handler tests (optional)
├── mocks_test.go      # Generated test mocks
└── postgres_test.go   # Integration tests (optional)
```

## Generated Files

### Go Generated Files

- `*.gen.go` - Do not edit manually
- `types.gen.go` - OpenAPI struct definitions
- `http.gen.go` - HTTP handler interfaces
- `repository.gen.go` - Database query methods
- `cache.gen.go` - Redis caching layer
- `events.gen.go` - Event publishing

### TypeScript Generated Files

- `web/client/src/generated/` - Complete API client
- Generated from `api/openapi.bundled.yaml`

### SQL Generated Files

- Database queries in `internal/database/queries/*.sql`
- Generate Go code with `sqlc generate`

## File Naming Conventions

### Go Files

- `generate.go` - Code generation annotations
- `service.go` - Business logic
- `handler.go` - HTTP handlers
- `middleware.go` - Middleware functions
- `mappers.go` - Type conversions (optional)
- `*_test.go` - Test files
- `*.gen.go` - Generated (don't edit)

### Config Files

- `.*.yaml` - YAML configs (golangci, mockery, sqlc)
- `.*.toml` - TOML configs (air)
- `.*.json` - JSON configs (markdownlint, tsconfig)
- `.*rc` - RC files (prettierrc)

### Documentation

- `*.md` - Markdown docs
- `*.mdx` - MDX with components (web/docs)
- `README.md` - Package/directory docs
