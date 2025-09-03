# Development Guide: Domain-Driven Go Architecture

## Project Overview

ArchesAI is a Go-based API platform that follows domain-driven design principles with clean architecture patterns. This guide explains the project structure, development workflows, and best practices.

## Project Structure

```
archesai/
├── api/                        # OpenAPI specifications
│   ├── components/            # Shared OpenAPI components
│   ├── paths/                 # API endpoints organized by resource
│   └── openapi.yaml           # Main OpenAPI specification
├── internal/                   # Private application code
│   ├── domains/               # Business domains (4 domains)
│   │   ├── auth/             # Authentication & user management
│   │   ├── organizations/    # Organization & membership management
│   │   ├── workflows/        # Pipeline workflows, runs, and tools
│   │   └── content/          # Content artifacts and labels
│   ├── infrastructure/       # Shared infrastructure
│   │   ├── database/         # DB connection & migrations
│   │   ├── server/           # HTTP server setup
│   │   ├── config/           # Configuration loading
│   │   └── logging/          # Structured logging
│   ├── generated/            # Generated code
│   │   ├── api/              # OpenAPI-generated types & handlers
│   │   └── database/         # SQLC-generated queries
│   └── app/                  # Application assembly
├── cmd/                      # Application entry points
└── docs/                     # Documentation
```

## Architecture Principles

### 1. Domain-Driven Design

- **Domains** represent business areas (auth, intelligence, admin)
- Each domain is self-contained with its own entities, services, and adapters
- Business logic stays within domain boundaries

### 2. Clean Architecture Layers

```
┌─────────────────┐
│   HTTP API      │ ← Generated OpenAPI handlers
├─────────────────┤
│   Services      │ ← Business logic & orchestration
├─────────────────┤
│   Repositories  │ ← Data access interfaces
├─────────────────┤
│   Adapters      │ ← Database, external APIs
└─────────────────┘
```

### 3. Dependency Direction

- Dependencies point inward (toward business logic)
- External concerns (HTTP, DB) depend on business logic
- Business logic doesn't depend on external systems

## Domain Structure (Flat Go-Centric Pattern)

Each domain follows a flat, Go-centric structure:

```
domains/auth/
├── auth.go            # Package documentation and shared constants
├── entities.go        # Domain models (often embedding API types)
├── service.go         # Business logic & Repository interface definition
├── repository.go      # PostgreSQL implementation of Repository
├── handler.go         # HTTP handlers implementing OpenAPI interfaces
├── middleware.go      # Domain-specific middleware (optional)
└── converters/        # Generated type converters (DO NOT EDIT)
    └── converters.gen.go
```

### File Responsibilities

**auth.go** - Package entry point:

```go
// Package auth provides authentication and authorization functionality.
package auth

// Shared constants
type ContextKey string

const (
    UserContextKey   ContextKey = "user"
    ClaimsContextKey ContextKey = "claims"
)
```

**entities.go** - Domain models:

```go
// User extends the API type with domain-specific fields
type User struct {
    api.UserEntity
    PasswordHash string `json:"-"`
}

// Domain-specific errors
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserNotFound      = errors.New("user not found")
)
```

**service.go** - Business logic (consumer defines interface):

```go
// Repository interface defined by the service (consumer)
type Repository interface {
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
}

// Service contains business logic
type Service struct {
    repo Repository
    jwt  *config.JWTConfig
}
```

**repository.go** - Database implementation:

```go
// Compile-time interface check
var _ Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
    q postgresql.Querier
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    // Implementation using sqlc-generated queries
}
```

**handler.go** - HTTP handlers:

```go
type Handler struct {
    service *Service
}

// Implements generated OpenAPI interface
func (h *Handler) PostAuthSignIn(ctx echo.Context) error {
    // Handle sign-in request
}
```

## Development Workflow

### 1. API-First Development

1. **Define API specs** in `api/specifications/`
2. **Generate code** with `make generate`
3. **Implement business logic** in domain services
4. **Create adapters** for external systems
5. **Wire dependencies** in `app/container.go`

### 2. Code Generation

The project uses two code generators:

**OpenAPI → Go Types & Handlers**

```bash
# Generate API types and server interfaces
make generate-api
```

**SQL → Go Database Code**

```bash
# Generate type-safe database queries
make generate-db
```

### 3. Adding New Features

**Step 1: Define API Contract**

```yaml
# api/specifications/intelligence/models/paths.yaml
paths:
  /api/v1/intelligence/models:
    get:
      summary: List AI models
      responses:
        '200':
          content:
            application/json:
              schema:
                $ref: './schemas.yaml#/components/schemas/ModelList'
```

**Step 2: Generate Code**

```bash
make generate
```

**Step 3: Implement Domain Logic**

```go
// domains/intelligence/entities/model.go
type Model struct {
    ID       string
    Name     string
    Provider string
    Status   ModelStatus
}

// domains/intelligence/services/model.go
type ModelService struct {
    repo repositories.ModelRepository
}

func (s *ModelService) ListModels(ctx context.Context) ([]*entities.Model, error) {
    return s.repo.FindAll(ctx)
}
```

**Step 4: Create HTTP Handler**

```go
// domains/intelligence/handlers/model.go
type ModelHandler struct {
    service *services.ModelService
}

func (h *ModelHandler) GetModels(ctx echo.Context) error {
    models, err := h.service.ListModels(ctx.Request().Context())
    if err != nil {
        return err
    }

    return ctx.JSON(http.StatusOK, convertModels(models))
}
```

**Step 5: Wire Dependencies**

```go
// app/container.go
type Container struct {
    // ... other dependencies
    ModelService *intelligence.ModelService
    ModelHandler *intelligence.ModelHandler
}

func NewContainer(cfg *config.Config) *Container {
    // ... initialization
    modelService := intelligence.NewModelService(modelRepo)
    modelHandler := intelligence.NewModelHandler(modelService)

    return &Container{
        ModelService: modelService,
        ModelHandler: modelHandler,
    }
}
```

## Database Integration

### Schema Migrations

```sql
-- internal/infrastructure/database/migrations/002_models.up.sql
CREATE TABLE models (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    provider text NOT NULL,
    status text NOT NULL DEFAULT 'active',
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
);
```

### Type-Safe Queries

```sql
-- internal/infrastructure/database/queries/models.sql
-- name: FindAllModels :many
SELECT id, name, provider, status, created_at, updated_at
FROM models
WHERE status = $1;

-- name: CreateModel :one
INSERT INTO models (name, provider, status)
VALUES ($1, $2, $3)
RETURNING *;
```

### Repository Implementation

```go
// domains/intelligence/adapters/postgres/model.go
type ModelRepository struct {
    db *database.Queries
}

func (r *ModelRepository) FindAll(ctx context.Context) ([]*entities.Model, error) {
    models, err := r.db.FindAllModels(ctx, "active")
    if err != nil {
        return nil, err
    }

    return convertDBModels(models), nil
}
```

## Configuration Management

### Go-Friendly Schema Design

**Instead of discriminated unions:**

```yaml
# ❌ Problematic - discriminated union
auth:
  anyOf:
    - properties: { mode: { const: disabled } }
    - properties: { mode: { const: enabled }, firebase: { ... } }
```

**Use embedded structs:**

```yaml
# ✅ Go-friendly
auth:
  type: object
  properties:
    enabled: { type: boolean, default: true }
    firebase:
      type: object
      properties:
        enabled: { type: boolean, default: false }
        projectId: { type: string }
        privateKey: { type: string }
    local:
      type: object
      properties:
        enabled: { type: boolean, default: true }
```

**Generated Go struct:**

```go
type AuthConfig struct {
    Enabled  bool          `json:"enabled"`
    Firebase *FirebaseAuth `json:"firebase,omitempty"`
    Local    *LocalAuth    `json:"local,omitempty"`
}
```

## Testing Strategy

### 1. Unit Tests - Domain Logic

```go
// domains/auth/services/auth_test.go
func TestAuthService_Login(t *testing.T) {
    // Test business logic in isolation
    mockRepo := &mocks.UserRepository{}
    service := NewAuthService(mockRepo)

    user, err := service.Login(ctx, "user@example.com", "password")

    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### 2. Integration Tests - Database Layer

```go
// domains/auth/adapters/postgres/user_test.go
func TestUserRepository_Create(t *testing.T) {
    // Test with real database
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    user := &entities.User{Email: "test@example.com"}
    err := repo.Create(ctx, user)

    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

### 3. End-to-End Tests - HTTP API

```go
// test/e2e/auth_test.go
func TestAuthAPI_Login(t *testing.T) {
    // Test complete request flow
    app := setupTestApp(t)

    resp := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/api/v1/auth/login", body)

    app.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
}
```

## Development Commands

```bash
# Code generation
make generate     # Generate all code (API + DB)
make generate-api # Generate OpenAPI code only
make generate-db  # Generate database code only

# Database operations
make migrate-up   # Apply database migrations
make migrate-down # Rollback migrations
make db-reset     # Reset database (dev only)

# Development server
make dev              # Run with hot reload
make build            # Build production binary
make test             # Run all tests
make test-unit        # Run unit tests only
make test-integration # Run integration tests only

# Code quality
make lint     # Run linter
make format   # Format code
make security # Security scan
```

## Best Practices

### 1. Domain Boundaries

- Keep business logic within domain services
- Use repository interfaces for data access
- External dependencies stay in adapters

### 2. Error Handling

```go
// Domain errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
)

// HTTP error mapping
func mapDomainError(err error) error {
    switch {
    case errors.Is(err, ErrUserNotFound):
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    case errors.Is(err, ErrInvalidCredentials):
        return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
    default:
        return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
    }
}
```

### 3. Configuration

- Use environment variables for deployment config
- Keep sensitive data out of code
- Provide sensible defaults

### 4. Logging & Monitoring

```go
// Structured logging
log.Info("user login attempt",
    zap.String("user_id", userID),
    zap.String("ip", clientIP),
    zap.Duration("duration", time.Since(start)))

// Metrics
metrics.Counter("auth.login.attempts").Inc()
metrics.Histogram("auth.login.duration").Observe(duration.Seconds())
```

## Code Generation Strategy

ArchesAI heavily leverages code generation to reduce boilerplate:

1. **sqlc** - Generates type-safe database queries from SQL
2. **oapi-codegen** - Generates server interfaces from OpenAPI spec
3. **generate-defaults** - Generates config with defaults from OpenAPI
4. **generate-converters** - Generates type converters between layers

### Converter Configuration

The `internal/domains/converters.yaml` file configures type conversions:

```yaml
converters:
  - name: PipelineDBToAPI
    from: postgresql.Pipeline
    to: api.PipelineEntity
    automap: true # Automatically map matching fields
    fields:
      # Only specify fields that need custom conversion
      OrganizationId: 'openapi_types.UUID(uuid.MustParse(from.OrganizationId))'
```

### Benefits of Code Generation

- **Type Safety**: Compile-time checking for database queries and API contracts
- **Reduced Boilerplate**: Auto-generate repetitive conversion code
- **Consistency**: Ensure uniform patterns across domains
- **Maintainability**: Changes to schemas automatically propagate

## Contributing

1. **Follow the architecture** - respect domain boundaries
2. **Generate code first** - define APIs before implementation
3. **Write tests** - unit tests for business logic, integration tests for adapters
4. **Use meaningful names** - functions and variables should be self-documenting
5. **Keep functions small** - single responsibility principle
6. **Document complex logic** - add comments for non-obvious business rules

## Resources

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [OpenAPI Specification](https://swagger.io/specification/)
- [SQLC Documentation](https://docs.sqlc.dev/)
