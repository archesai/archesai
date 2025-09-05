# Development Guide: Hexagonal Architecture with Domain-Driven Design

## Project Overview

ArchesAI is a Go-based API platform that follows hexagonal architecture (ports and adapters) with Domain-Driven Design principles. This guide explains the architectural patterns, development workflows, and best practices.

## Architecture Principles

### 1. Hexagonal Architecture (Ports & Adapters)

The hexagonal architecture isolates the core business logic from external concerns:

- **Core (Domain)**: Contains business logic, entities, and port interfaces
- **Ports**: Interfaces that define how the core interacts with the outside world
- **Adapters**: Implementations of ports (database, HTTP handlers, external services)

```
         ┌──────────────────────────────┐
         │      HTTP Handlers           │
         │        (Adapters)            │
         └─────────────┬────────────────┘
                       │
         ┌─────────────▼────────────────┐
         │         Core Domain          │
         │   ┌───────────────────┐      │
         │   │   Use Cases       │      │
         │   ├───────────────────┤      │
         │   │   Entities        │      │
         │   ├───────────────────┤      │
         │   │   Ports           │      │
         │   └───────────────────┘      │
         └─────────────┬────────────────┘
                       │
         ┌─────────────▼────────────────┐
         │     Infrastructure           │
         │   (Database Adapters)        │
         └──────────────────────────────┘
```

### 2. Domain-Driven Design

- **Domains** represent bounded contexts (auth, organizations, workflows, content)
- Each domain is self-contained with its own entities, use cases, and adapters
- Business logic stays within domain boundaries
- Domains communicate through well-defined interfaces

### 3. Dependency Inversion

- Dependencies point inward (toward business logic)
- Core domain doesn't depend on infrastructure or presentation layers
- Infrastructure and handlers depend on core domain interfaces

## Domain Structure (Hexagonal Pattern)

Each domain follows a consistent hexagonal structure:

```
domains/auth/
├── auth.go                    # Package documentation and constants
├── core/                      # Core business logic (hexagon center)
│   ├── entities.go           # Domain models and value objects
│   ├── ports.go              # Interface definitions (Repository, Services)
│   └── usecase.go            # Business use cases and orchestration
├── infrastructure/            # Infrastructure adapters
│   └── postgres.go           # PostgreSQL repository implementation
├── handlers/                  # Presentation layer adapters
│   └── http/
│       ├── handler.go        # HTTP handlers implementing OpenAPI
│       └── middleware.go     # HTTP middleware (auth only)
├── adapters/                  # Type converters (generated)
│   └── adapters.gen.go       # Auto-generated DB<->API converters
└── generated/                 # Domain-specific generated code
    └── api/
        ├── types.gen.go      # OpenAPI types for this domain
        └── server.gen.go     # OpenAPI server interfaces
```

### File Responsibilities

**auth.go** - Package entry point:

```go
// Package auth provides authentication and authorization functionality.
package auth

// Shared constants and configuration
type ContextKey string

const (
    UserContextKey   ContextKey = "user"
    ClaimsContextKey ContextKey = "claims"
)
```

**core/entities.go** - Domain models:

```go
// Domain entities - pure business objects
type User struct {
    ID           uuid.UUID
    Email        string
    Name         string
    PasswordHash string // Not exposed in API
    CreatedAt    time.Time
}

// Domain-specific errors
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserNotFound      = errors.New("user not found")
)
```

**core/ports.go** - Interface definitions (Dependency Inversion):

```go
// Repository interface - defined by the domain, implemented by infrastructure
type Repository interface {
    GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
    UpdateUser(ctx context.Context, user *User) error
}

// External service interfaces
type EmailService interface {
    SendVerificationEmail(ctx context.Context, email string, token string) error
}
```

**core/usecase.go** - Business logic:

```go
// Service contains business logic - the hexagon core
type Service struct {
    repo         Repository
    emailService EmailService
    jwtConfig    *config.JWTConfig
}

// NewService creates a new auth service
func NewService(repo Repository, email EmailService, jwt *config.JWTConfig) *Service {
    return &Service{
        repo:         repo,
        emailService: email,
        jwtConfig:    jwt,
    }
}

// Business use cases
func (s *Service) SignIn(ctx context.Context, email, password string) (*User, string, error) {
    // Core business logic here
}
```

**infrastructure/postgres.go** - Database adapter:

```go
// PostgresRepository implements the Repository port
type PostgresRepository struct {
    queries *postgresql.Queries
}

// Compile-time interface check
var _ authcore.Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(queries *postgresql.Queries) *PostgresRepository {
    return &PostgresRepository{queries: queries}
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*authcore.User, error) {
    // Implementation using generated SQLC queries
}
```

**handlers/http/handler.go** - HTTP adapter:

```go
// Handler adapts HTTP requests to domain use cases
type Handler struct {
    service *authcore.Service
}

// NewHandler creates a new HTTP handler
func NewHandler(service *authcore.Service) *Handler {
    return &Handler{service: service}
}

// Implements generated OpenAPI interface
func (h *Handler) PostAuthSignIn(ctx echo.Context) error {
    // Adapt HTTP request to domain use case
}
```

## Development Workflow

### 1. API-First Development

1. **Define API contract** in `api/openapi.yaml` or component files
2. **Generate code** with `make generate`
3. **Implement core domain logic** in `core/usecase.go`
4. **Create infrastructure adapters** in `infrastructure/`
5. **Implement HTTP handlers** in `handlers/http/`
6. **Wire dependencies** in `internal/app/deps.go`

### 2. Adding New Features

**Step 1: Define OpenAPI Contract**

```yaml
# api/paths/users.yaml
/api/v1/users/{id}:
  get:
    summary: Get user by ID
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
          format: uuid
```

**Step 2: Generate Code**

```bash
make generate
```

**Step 3: Define Domain Entity**

```go
// internal/domains/auth/core/entities.go
type User struct {
    ID        uuid.UUID
    Email     string
    Name      string
    Role      UserRole
    CreatedAt time.Time
}
```

**Step 4: Define Port Interface**

```go
// internal/domains/auth/core/ports.go
type Repository interface {
    GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
}
```

**Step 5: Implement Use Case**

```go
// internal/domains/auth/core/usecase.go
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    return s.repo.GetUserByID(ctx, id)
}
```

**Step 6: Implement Infrastructure Adapter**

```go
// internal/domains/auth/infrastructure/postgres.go
func (r *PostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*authcore.User, error) {
    dbUser, err := r.queries.GetUser(ctx, id.String())
    if err != nil {
        return nil, err
    }
    return r.toDomainUser(dbUser), nil
}
```

**Step 7: Implement HTTP Handler**

```go
// internal/domains/auth/handlers/http/handler.go
func (h *Handler) GetUsersId(ctx echo.Context, id openapi_types.UUID) error {
    user, err := h.service.GetUser(ctx.Request().Context(), uuid.UUID(id))
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound)
    }

    apiUser := h.toAPIUser(user)
    return ctx.JSON(http.StatusOK, apiUser)
}
```

## Database Integration

### Type-Safe Queries with SQLC

```sql
-- internal/infrastructure/database/queries/users.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;
```

### Generated Type Converters

The `adapters.yaml` file configures automatic type conversion between layers:

```yaml
# internal/domains/adapters.yaml
converters:
  - name: AuthUserDBToAPI
    from: postgresql.User
    to: api.UserEntity
    automap: true
    overrides:
      Id: "uuid.MustParse(from.Id)"
      Email: "openapi_types.Email(from.Email)"
```

Generated converters handle:

- Nullable field conversions
- Type transformations (string ↔ UUID)
- Timestamp formatting
- Custom type mappings

## Testing Strategy

### 1. Unit Tests - Core Domain Logic

Test business logic in isolation using mocks:

```go
// internal/domains/auth/core/usecase_test.go
func TestService_SignIn(t *testing.T) {
    // Create mock repository
    mockRepo := &mocks.Repository{}
    mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(&User{Email: "test@example.com"}, nil)

    // Test service with mock
    service := NewService(mockRepo, nil, testJWTConfig)
    user, token, err := service.SignIn(ctx, "test@example.com", "password")

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.NotEmpty(t, token)
}
```

### 2. Integration Tests - Infrastructure Layer

Test database interactions with real database:

```go
// internal/domains/auth/infrastructure/postgres_test.go
func TestPostgresRepository_GetUserByEmail(t *testing.T) {
    // Use test database
    db := setupTestDB(t)
    repo := NewPostgresRepository(db)

    // Test real database operations
    user, err := repo.GetUserByEmail(ctx, "test@example.com")

    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### 3. End-to-End Tests - Full Request Flow

Test complete HTTP request handling:

```go
// test/e2e/auth_test.go
func TestAuthAPI_SignIn(t *testing.T) {
    // Setup test server with real dependencies
    app := setupTestApp(t)

    // Make HTTP request
    req := httptest.NewRequest("POST", "/api/v1/auth/signin", body)
    rec := httptest.NewRecorder()

    app.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)
}
```

## Configuration Management

### Environment-Based Configuration

```go
// internal/infrastructure/config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

// Load from environment with defaults
func Load() (*Config, error) {
    viper.SetEnvPrefix("ARCHESAI")
    viper.AutomaticEnv()

    // Set defaults
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("database.pool_size", 10)

    var config Config
    return &config, viper.Unmarshal(&config)
}
```

## Best Practices

### 1. Domain Isolation

- Keep business logic in `core/` package
- Don't import infrastructure or handlers in core
- Use interfaces to define external dependencies

### 2. Error Handling

```go
// Domain errors (in core/entities.go)
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)

// Infrastructure error wrapping
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// HTTP error mapping (in handlers)
switch {
case errors.Is(err, authcore.ErrUserNotFound):
    return echo.NewHTTPError(http.StatusNotFound)
case errors.Is(err, authcore.ErrInvalidInput):
    return echo.NewHTTPError(http.StatusBadRequest)
default:
    return echo.NewHTTPError(http.StatusInternalServerError)
}
```

### 3. Dependency Injection

All dependencies are wired in `internal/app/deps.go`:

```go
// Create repositories
authRepo := authinfra.NewPostgresRepository(queries)

// Create services with dependencies
authService := authcore.NewService(authRepo, emailService, cfg.JWT)

// Create handlers
authHandler := authhandlers.NewHandler(authService)
```

### 4. Code Generation

The project uses multiple code generators:

1. **SQLC** - Type-safe database queries
2. **oapi-codegen** - OpenAPI server interfaces
3. **generate-adapters** - Type converters between layers
4. **generate-domain** - Scaffold new domains

Always run `make generate` after:

- Modifying OpenAPI specifications
- Adding SQL queries
- Updating adapters.yaml
- Creating new domains

## Creating a New Domain

Use the domain generator to scaffold a new bounded context:

```bash
make generate-domain name=billing tables=subscription,invoice
```

This creates:

- Domain structure with core/infrastructure/handlers
- Database migrations
- Basic CRUD operations
- OpenAPI specifications
- Type converters configuration

Then:

1. Define your business logic in `core/usecase.go`
2. Add domain-specific methods to the repository
3. Implement custom handlers as needed
4. Wire the domain in `internal/app/deps.go`

## Command Reference

For a complete list of development commands, see [.claude/CLAUDE.md](../.claude/CLAUDE.md).

## Resources

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [OpenAPI Specification](https://swagger.io/specification/)
- [SQLC Documentation](https://docs.sqlc.dev/)
