# Architecture: Combining Generated Code with Domain Logic

## Overview

This document explains how to properly integrate OpenAPI-generated code (`@gen/`) with your domain logic (`@internal/features/`).

## Directory Structure

```
archesai/
├── gen/                          # Generated code (DO NOT EDIT)
│   ├── api/                      # OpenAPI generated
│   │   └── features/
│   │       └── auth/
│   │           ├── users/        # User endpoints
│   │           │   ├── paths.gen.go     # ServerInterface & handlers
│   │           │   └── schemas.gen.go   # Request/Response types
│   │           ├── accounts/
│   │           ├── sessions/
│   │           └── auth.gen.go  # Shared auth types
│   └── db/                       # SQLC generated
│       └── postgresql/
│           └── *.sql.go
│
└── internal/                     # Your domain logic
    └── features/
        └── auth/
            ├── domain/           # Business entities
            ├── ports/            # Interfaces
            ├── usecase/          # Business logic
            └── adapters/         # Implementations
                ├── http/         # HTTP handlers
                └── postgresql/   # Database

```

## Integration Pattern

### 1. Generated Interfaces (from OpenAPI)

Each OpenAPI path generates a `ServerInterface` that your handlers must implement:

```go
// gen/api/features/auth/users/paths.gen.go
type ServerInterface interface {
    FindManyUsers(ctx echo.Context, params FindManyUsersParams) error
    DeleteUser(ctx echo.Context, id openapi_types.UUID) error
    GetOneUser(ctx echo.Context, id openapi_types.UUID) error
    UpdateUser(ctx echo.Context, id openapi_types.UUID) error
}
```

### 2. Your Implementation

Your handler implements these interfaces:

```go
// internal/features/auth/adapters/http/handler.go
type Handler struct {
    service ports.Service
    logger  *zap.Logger
}

// Implement the generated interface
var _ users.ServerInterface = (*Handler)(nil)

func (h *Handler) GetOneUser(ctx echo.Context, id openapi_types.UUID) error {
    // Convert generated types to domain types
    userID := uuid.UUID(id)

    // Call your domain service
    user, err := h.service.GetUser(ctx.Request().Context(), userID)

    // Convert domain response to generated types
    response := convertToGeneratedUser(user)
    return ctx.JSON(http.StatusOK, response)
}
```

### 3. Route Registration

Routes are registered in two ways:

```go
// internal/app/routes.go
func RegisterRoutes(e *echo.Echo, container *Container) {
    v1 := e.Group("/api/v1")

    // 1. Custom routes (not in OpenAPI)
    container.AuthHandler.RegisterRoutes(v1)

    // 2. OpenAPI-generated routes
    userGroup := v1.Group("/auth")
    users.RegisterHandlers(userGroup, container.AuthHandler)
}
```

## Implementation Steps

### Step 1: Check Generated Interfaces

For each feature, check what interfaces are generated:

```bash
# Find all ServerInterface definitions
grep -r "type ServerInterface" ./gen/api/features/auth/
```

### Step 2: Implement Required Interfaces

For each generated interface, create implementations in your handler:

```go
// For each ServerInterface method:
// 1. Accept generated request types
// 2. Convert to domain types
// 3. Call domain service
// 4. Convert response to generated types
// 5. Return appropriate HTTP response
```

### Step 3: Type Conversions

Create conversion helpers between generated and domain types:

```go
// internal/features/auth/adapters/http/converters.go
func convertToGeneratedUser(u *domain.User) *users.User {
    return &users.User{
        Id:        openapi_types.UUID(u.ID),
        Email:     openapi_types.Email(u.Email),
        Name:      u.Name,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}

func convertFromGeneratedUser(gu *users.UpdateUserRequest) *domain.UpdateUserInput {
    return &domain.UpdateUserInput{
        Name:  gu.Name,
        Email: string(gu.Email),
    }
}
```

### Step 4: Register All Handlers

Update `routes.go` to register all implemented interfaces:

```go
// Register each generated handler group
accounts.RegisterHandlers(authGroup, container.AccountHandler)
sessions.RegisterHandlers(authGroup, container.SessionHandler)
organizations.RegisterHandlers(authGroup, container.OrgHandler)
```

## Benefits of This Approach

1. **Type Safety**: OpenAPI ensures consistent types between frontend and backend
2. **Separation of Concerns**: Generated code stays separate from business logic
3. **Clean Architecture**: Domain logic doesn't depend on HTTP concerns
4. **Easy Updates**: Regenerate from OpenAPI without losing business logic

## Common Patterns

### Pagination

Generated pagination params should be converted to domain pagination:

```go
func (h *Handler) FindManyUsers(ctx echo.Context, params users.FindManyUsersParams) error {
    // Convert generated pagination to domain pagination
    opts := domain.QueryOptions{
        Limit:  convertLimit(params.Page),
        Offset: convertOffset(params.Page),
        Filter: convertFilter(params.Filter),
    }

    users, total, err := h.service.ListUsers(ctx.Request().Context(), opts)
    // ...
}
```

### Error Handling

Map domain errors to HTTP errors:

```go
func mapDomainError(err error) error {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    case errors.Is(err, domain.ErrUnauthorized):
        return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
    default:
        return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
    }
}
```

### Authentication Middleware

Use Echo middleware with generated scopes:

```go
// Check auth scopes from generated constants
requiredScopes := []string{"users:read"}
ctx.Set(auth.BearerAuthScopes, requiredScopes)
```

## Testing

Test your implementations separately:

1. **Unit tests**: Test domain logic without HTTP
2. **Integration tests**: Test HTTP handlers with mocked services
3. **E2E tests**: Test full flow with real database

## Regeneration Workflow

When OpenAPI spec changes:

1. Regenerate code: `make generate`
2. Check for interface changes
3. Update handler implementations
4. Update type converters
5. Run tests

## Example: Complete Feature Implementation

See `internal/features/auth/` for a complete example implementing:

- User CRUD operations
- Authentication flows
- Session management
- All following the patterns described above
