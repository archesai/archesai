# Development Guide

## Architecture Overview

ArchesAI follows **Hexagonal Architecture** (Ports & Adapters) with **Domain-Driven Design** principles. The core business logic is isolated from external concerns through well-defined interfaces.

```
┌─────────────────────────────────┐
│     HTTP Handlers (Adapters)     │
└────────────┬─────────────────────┘
             │
┌────────────▼─────────────────────┐
│       Domain Service (Core)      │
│   ┌───────────────────────┐      │
│   │   Business Logic      │      │
│   ├───────────────────────┤      │
│   │   Repository Interface│      │
│   └───────────────────────┘      │
└────────────┬─────────────────────┘
             │
┌────────────▼─────────────────────┐
│   PostgreSQL/SQLite (Adapters)   │
└──────────────────────────────────┘
```

## Domain Structure

Each domain uses a flat package structure with clear naming conventions:

```
internal/auth/
├── auth.go                    # Package constants and errors
├── service.go                 # Business logic (manual)
├── handler.go                 # HTTP handler implementation (manual)
├── http.gen.go                # Generated HTTP interface (ServerInterface)
├── middleware_http.go         # Authentication middleware
├── repository_postgres.go     # PostgreSQL repository (manual)
├── repository_sqlite.go       # SQLite repository (manual)
├── repository.gen.go          # Generated repository interface
├── types.gen.go               # Generated OpenAPI types
├── cache.gen.go               # Generated cache interface (future)
├── cache_memory.gen.go        # Memory cache implementation (future)
├── cache_redis.gen.go         # Redis cache implementation (future)
├── events.gen.go              # Generated event types (future)
├── events_redis.gen.go        # Redis event publisher (future)
└── events_nats.gen.go         # NATS event publisher (future)
```

## Code Generation Pipeline

### 1. OpenAPI → Go Types & Interfaces

Define schemas in `api/components/schemas/`:

```yaml
# api/components/schemas/User.yaml
User:
  type: object
  x-codegen:
    repository:
      operations: [create, read, update, delete, list]
      indices: [email]
      additional_methods:
        - name: GetUserByEmail
          params: [email]
    cache:
      enabled: true
      ttl: 300
    events:
      - created
      - updated
      - deleted
  properties:
    id:
      type: string
      format: uuid
      x-codegen:
        primary-key: true
    email:
      type: string
      format: email
      x-codegen:
        unique: true
        index: true
```

This generates:

- `types.gen.go` - OpenAPI types (User struct)
- `http.gen.go` - ServerInterface with HTTP handler methods
- `repository.gen.go` - Repository interface with CRUD operations

### 2. SQL → Database Code

Define queries in `internal/database/queries/`:

```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, name, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $2, name = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
```

### 3. Run Generators

```bash
make generate         # Run all generators
make generate-sqlc    # Database queries only
make generate-oapi    # OpenAPI types only
make generate-codegen # Repository interfaces from x-codegen
```

### Current Generation Status

✅ **Working:**

- OpenAPI types generation (`types.gen.go`)
- HTTP handler interfaces (`http.gen.go`)
- Repository interfaces for schemas with x-codegen:
  - Auth: User, Session, Account
  - Organizations: Organization, Member, Invitation
  - Workflows: Pipeline, Run
  - Content: Artifact, Label

⚠️ **In Progress:**

- Tool entity (has x-codegen but not generating)
- Cache interfaces (x-codegen defined but generator not active)
- Event publishers (x-codegen defined but generator not active)

📝 **Manual Implementation Required:**

- Repository implementations (PostgreSQL/SQLite)
- Service business logic
- HTTP handler implementations

## Adding a New Feature

### Step 1: Define API Contract

Add endpoint to `api/paths/`:

```yaml
# api/paths/users.yaml
/users/{id}/profile:
  get:
    operationId: getUserProfile
    tags: [Users]
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      "200":
        description: User profile
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/UserProfile"
```

### Step 2: Add Database Queries

Create `internal/database/queries/user_profile.sql`:

```sql
-- name: GetUserProfile :one
SELECT u.*, COUNT(o.id) as organization_count
FROM users u
LEFT JOIN members m ON m.user_id = u.id
LEFT JOIN organizations o ON o.id = m.organization_id
WHERE u.id = $1
GROUP BY u.id;
```

### Step 3: Generate Code

```bash
make generate
```

### Step 4: Implement Business Logic

Update `internal/auth/service.go`:

```go
func (s *Service) GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
    // Check cache first
    if cached, err := s.cache.GetUserProfile(ctx, userID); err == nil {
        return cached, nil
    }

    // Get from database
    profile, err := s.repo.GetUserProfile(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("get user profile: %w", err)
    }

    // Update cache
    _ = s.cache.SetUserProfile(ctx, profile, 5*time.Minute)

    // Publish event
    _ = s.events.PublishProfileViewed(ctx, userID)

    return profile, nil
}
```

### Step 5: Implement HTTP Handler

Update `internal/auth/handler.go`:

```go
func (h *Handler) GetUserProfile(ctx echo.Context, id openapi_types.UUID) error {
    profile, err := h.service.GetUserProfile(ctx.Request().Context(), uuid.UUID(id))
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return echo.NewHTTPError(http.StatusNotFound, "User not found")
        }
        return fmt.Errorf("get user profile: %w", err)
    }

    return ctx.JSON(http.StatusOK, profile)
}
```

### Step 6: Wire Dependencies

Update `internal/app/app.go`:

```go
func New(cfg *config.Config) (*App, error) {
    // ... existing code ...

    // Wire auth domain
    authRepo := auth.NewPostgresRepository(queries)
    authCache := auth.NewRedisCache(redisClient)
    authEvents := auth.NewRedisEventPublisher(redisClient)
    authService := auth.NewService(authRepo, authCache, authEvents, cfg.JWT)
    authHandler := auth.NewHandler(authService)

    // ... register routes ...
}
```

## Testing

### Unit Tests

Test business logic in isolation:

```go
// internal/auth/service_test.go
func TestService_CreateUser(t *testing.T) {
    mockRepo := &MockRepository{}
    mockCache := &MockCache{}
    mockEvents := &MockEventPublisher{}

    service := auth.NewService(mockRepo, mockCache, mockEvents, testConfig)

    user, err := service.CreateUser(ctx, CreateUserRequest{
        Email: "test@example.com",
        Name:  "Test User",
    })

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### Integration Tests

Test with real database:

```go
// internal/auth/repository_postgres_test.go
func TestPostgresRepository_CreateUser(t *testing.T) {
    db := testutil.SetupTestDB(t)
    defer db.Close()

    repo := auth.NewPostgresRepository(postgresql.New(db))

    user := &auth.User{
        ID:    uuid.New(),
        Email: "test@example.com",
        Name:  "Test User",
    }

    err := repo.CreateUser(context.Background(), user)
    assert.NoError(t, err)

    retrieved, err := repo.GetUserByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Email, retrieved.Email)
}
```

### API Tests

Test HTTP endpoints:

```go
// internal/auth/handler_test.go
func TestHandler_Login(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
        strings.NewReader(`{"email":"test@example.com","password":"secret"}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    handler := setupTestHandler(t)
    err := handler.PostAuthLogin(c)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)

    var response LoginResponse
    json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NotEmpty(t, response.AccessToken)
}
```

## Database Migrations

### Create Migration

```bash
make migrate-create name=add_user_profile
```

### Migration File Structure

```sql
-- internal/migrations/postgresql/TIMESTAMP_add_user_profile.sql
-- +goose Up
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
CREATE INDEX idx_users_avatar_url ON users(avatar_url);

-- +goose Down
DROP INDEX IF EXISTS idx_users_avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS bio;
```

### Apply Migrations

```bash
make migrate-up        # Apply all pending
make migrate-down      # Rollback last
make migrate-status    # Check status
```

## Configuration Management

### Environment Variables

All config uses `ARCHESAI_` prefix and maps to structs:

```go
// internal/config/config.go
type Config struct {
    Database DatabaseConfig `envconfig:"DATABASE"`
    Server   ServerConfig   `envconfig:"SERVER"`
    JWT      JWTConfig      `envconfig:"JWT"`
    Redis    RedisConfig    `envconfig:"REDIS"`
    Logging  LoggingConfig  `envconfig:"LOGGING"`
}

type DatabaseConfig struct {
    URL         string        `envconfig:"URL" required:"true"`
    PoolSize    int           `envconfig:"POOL_SIZE" default:"10"`
    MaxIdleTime time.Duration `envconfig:"MAX_IDLE_TIME" default:"30m"`
}
```

### Loading Configuration

```go
cfg, err := config.Load()
if err != nil {
    log.Fatal("Failed to load config:", err)
}
```

## Error Handling

### Domain Errors

Define in domain package:

```go
// internal/auth/auth.go
var (
    ErrUserNotFound       = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserAlreadyExists  = errors.New("user already exists")
    ErrSessionExpired     = errors.New("session expired")
)
```

### Error Wrapping

Always wrap errors with context:

```go
user, err := repo.GetUserByEmail(ctx, email)
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrUserNotFound
    }
    return nil, fmt.Errorf("get user by email: %w", err)
}
```

### HTTP Error Responses

Convert to appropriate HTTP status:

```go
if err != nil {
    switch {
    case errors.Is(err, auth.ErrUserNotFound):
        return echo.NewHTTPError(http.StatusNotFound, "User not found")
    case errors.Is(err, auth.ErrInvalidCredentials):
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
    default:
        log.Error("Unexpected error:", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
    }
}
```

## Logging

### Structured Logging

Use structured logging with context:

```go
log := logger.FromContext(ctx)
log.Info("Creating user",
    zap.String("email", email),
    zap.String("organization_id", orgID.String()),
)
```

### Log Levels

- **Debug**: Detailed debugging information
- **Info**: General informational messages
- **Warn**: Warning messages for recoverable issues
- **Error**: Error messages for failures
- **Fatal**: Fatal errors that cause shutdown

## Performance Optimization

### Database Queries

1. **Use indexes** for frequently queried columns
2. **Batch operations** when possible
3. **Use prepared statements** (SQLC does this automatically)
4. **Limit result sets** with pagination

### Caching Strategy

1. **Cache frequently accessed data** (users, sessions)
2. **Use appropriate TTLs** based on data volatility
3. **Invalidate on updates** to maintain consistency
4. **Use cache-aside pattern** for simplicity

### Concurrent Processing

```go
// Process items concurrently with bounded parallelism
sem := make(chan struct{}, 10) // Max 10 concurrent
var wg sync.WaitGroup

for _, item := range items {
    wg.Add(1)
    sem <- struct{}{}

    go func(item Item) {
        defer wg.Done()
        defer func() { <-sem }()

        processItem(item)
    }(item)
}

wg.Wait()
```

## Debugging

### Enable Debug Logging

```bash
ARCHESAI_LOGGING_LEVEL=debug make dev
```

### Database Query Logging

```bash
ARCHESAI_DATABASE_LOG_QUERIES=true make dev
```

### Performance Profiling

```go
import _ "net/http/pprof"

// In main.go
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

Access profiles at http://localhost:6060/debug/pprof/

## Troubleshooting

### Repository Not Generating

If a domain's repository isn't being generated despite having x-codegen:

1. **Check schema name detection** in `internal/codegen/parser.go`:

```go
// The inferDomain function must detect your schema
case strings.Contains(schemaLower, "artifact") ||
     strings.Contains(schemaLower, "label"):
    return "content"
```

2. **Verify x-codegen annotation**:

```yaml
x-codegen:
  repository:
    operations: [create, read, update, delete, list]
```

3. **Rebundle OpenAPI spec**:

```bash
make bundle
make generate-codegen
```

### Method Signature Mismatches

If you get errors like "GetArtifact undefined":

1. Generated repository uses standardized names:
   - `GetArtifactByID` (not `GetArtifact`)
   - `UpdateArtifact(ctx, id, entity)` (not `UpdateArtifact(ctx, entity)`)

2. Update service calls to match:

```go
// Before
artifact, err := s.repo.GetArtifact(ctx, id)

// After
artifact, err := s.repo.GetArtifactByID(ctx, id)
```

### Package Conflicts

If you get "found packages X and Y in same directory":

1. Ensure all files in a directory use the same package name
2. Check generated files match manual files:
   - `internal/server/http/` → `package http`
   - `internal/auth/` → `package auth`

### Unused Parameter Warnings

For stub implementations, use underscore:

```go
// Instead of:
func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    panic("unimplemented")
}

// Use:
func (r *Repository) GetUser(_ context.Context, _ uuid.UUID) (*User, error) {
    panic("unimplemented")
}
```

## Best Practices

### 1. **Keep domains isolated** - No cross-domain imports

### 2. **Use dependency injection** - Pass interfaces, not implementations

### 3. **Generate what you can** - Reduce manual boilerplate

### 4. **Test business logic** - Focus tests on service layer

### 5. **Handle errors explicitly** - No silent failures

### 6. **Use context everywhere** - For cancellation and tracing

### 7. **Validate at boundaries** - Input validation in handlers

### 8. **Log actions, not state** - Log what happened, not data dumps

### 9. **Cache judiciously** - Only cache what's expensive to compute

### 10. **Document decisions** - Use ADRs for architectural choices
