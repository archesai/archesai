# Clean Architecture Implementation Guide for Arches

## Overview

This guide documents the complete restructuring of the Arches project to follow clean/hexagonal architecture principles with `core` as the domain layer.

## Complete Clean Architecture Structure for Arches

```
archesai/
├── api/                         # OpenAPI specs (unchanged)
│   ├── openapi.yaml
│   ├── components/
│   └── paths/
│
├── cmd/                         # Application entry points (unchanged)
│   └── archesai/
│       └── main.go
│
├── deployments/                 # Deployment configs (unchanged)
├── docs/                        # Documentation (unchanged)
├── scripts/                     # Utility scripts (unchanged)
├── tools/                       # Development tools (unchanged)
├── web/                         # Frontend applications (unchanged)
│
└── internal/                    # THE MAIN REFACTORING HAPPENS HERE
    │
    ├── core/                    # 🎯 PURE DOMAIN LAYER (New)
    │   ├── entities/            # Domain entities
    │   │   ├── user.go
    │   │   ├── organization.go
    │   │   ├── pipeline.go
    │   │   ├── artifact.go
    │   │   ├── run.go
    │   │   ├── tool.go
    │   │   ├── label.go
    │   │   ├── member.go
    │   │   ├── invitation.go
    │   │   └── account.go
    │   │
    │   ├── valueobjects/        # Value objects
    │   │   ├── email.go
    │   │   ├── ids.go          # UserID, OrgID, PipelineID, etc.
    │   │   ├── role.go
    │   │   ├── status.go        # PipelineStatus, RunStatus
    │   │   ├── provider.go      # OAuth Provider
    │   │   ├── permissions.go
    │   │   └── credits.go
    │   │
    │   ├── aggregates/          # Aggregate roots
    │   │   ├── workspace.go     # Organization + Members + Settings
    │   │   ├── project.go       # Pipeline + Runs + Artifacts
    │   │   └── session.go       # User + Auth + Tokens
    │   │
    │   ├── services/            # Domain services
    │   │   ├── auth_service.go       # Authentication logic
    │   │   ├── permission_service.go # Authorization rules
    │   │   ├── pipeline_executor.go  # Pipeline execution logic
    │   │   ├── credit_calculator.go  # Billing logic
    │   │   └── invitation_service.go # Invitation workflow
    │   │
    │   ├── events/              # Domain events
    │   │   ├── events.go        # Base event interface
    │   │   ├── user_events.go   # UserCreated, EmailChanged, etc.
    │   │   ├── org_events.go    # OrgCreated, MemberAdded, etc.
    │   │   └── pipeline_events.go # PipelineStarted, RunCompleted
    │   │
    │   ├── errors/              # Domain errors
    │   │   ├── errors.go        # Base error types
    │   │   ├── validation.go    # Validation errors
    │   │   └── business.go      # Business rule violations
    │   │
    │   └── ports/               # 🔌 INTERFACES DEFINED BY DOMAIN
    │       ├── repositories/    # Repository interfaces
    │       │   ├── user_repository.go
    │       │   ├── organization_repository.go
    │       │   ├── pipeline_repository.go
    │       │   ├── artifact_repository.go
    │       │   ├── run_repository.go
    │       │   └── session_repository.go
    │       │
    │       ├── services/        # External service interfaces
    │       │   ├── email_service.go
    │       │   ├── storage_service.go
    │       │   ├── llm_service.go
    │       │   ├── oauth_provider.go
    │       │   └── payment_service.go
    │       │
    │       └── events/          # Event bus interfaces
    │           ├── publisher.go
    │           └── subscriber.go
    │
    ├── application/             # 📋 USE CASES LAYER (New)
    │   ├── commands/            # Write operations
    │   │   ├── users/
    │   │   │   ├── create_user.go
    │   │   │   ├── update_user.go
    │   │   │   ├── verify_email.go
    │   │   │   └── delete_user.go
    │   │   ├── organizations/
    │   │   │   ├── create_organization.go
    │   │   │   ├── invite_member.go
    │   │   │   └── update_settings.go
    │   │   ├── pipelines/
    │   │   │   ├── create_pipeline.go
    │   │   │   ├── execute_pipeline.go
    │   │   │   └── update_pipeline.go
    │   │   └── auth/
    │   │       ├── login.go
    │   │       ├── logout.go
    │   │       ├── refresh_token.go
    │   │       └── request_magic_link.go
    │   │
    │   ├── queries/             # Read operations
    │   │   ├── users/
    │   │   │   ├── get_user.go
    │   │   │   ├── list_users.go
    │   │   │   └── search_users.go
    │   │   ├── pipelines/
    │   │   │   ├── get_pipeline.go
    │   │   │   ├── list_pipelines.go
    │   │   │   └── get_execution_plan.go
    │   │   └── analytics/
    │   │       ├── usage_stats.go
    │   │       └── billing_summary.go
    │   │
    │   ├── dto/                 # Data transfer objects
    │   │   ├── user_dto.go
    │   │   ├── organization_dto.go
    │   │   ├── pipeline_dto.go
    │   │   └── common_dto.go    # PageInfo, SortParams
    │   │
    │   └── mappers/             # Entity <-> DTO mappers
    │       ├── user_mapper.go
    │       ├── organization_mapper.go
    │       └── pipeline_mapper.go
    │
    ├── adapters/                # 🔄 INTERFACE ADAPTERS (Refactored from existing)
    │   ├── http/                # HTTP layer (from current handlers)
    │   │   ├── server/          # Server setup (from internal/server)
    │   │   │   └── server.go
    │   │   ├── handlers/        # HTTP handlers
    │   │   │   ├── user_handler.go      # Refactored from users/handler.gen.go
    │   │   │   ├── org_handler.go       # Refactored from organizations/
    │   │   │   ├── pipeline_handler.go  # Refactored from pipelines/
    │   │   │   ├── auth_handler.go      # Refactored from auth/
    │   │   │   └── health_handler.go    # From health/
    │   │   ├── middleware/      # HTTP middleware (from internal/middleware)
    │   │   │   ├── auth.go
    │   │   │   ├── cors.go
    │   │   │   └── logging.go
    │   │   ├── dto/             # HTTP-specific DTOs
    │   │   │   ├── requests.go  # Generated from OpenAPI
    │   │   │   └── responses.go # Generated from OpenAPI
    │   │   └── mappers/         # HTTP DTO <-> Application DTO
    │   │       └── http_mapper.go
    │   │
    │   ├── cli/                 # CLI adapter (from internal/cli)
    │   │   ├── commands/
    │   │   │   ├── api.go
    │   │   │   ├── tui.go
    │   │   │   └── config.go
    │   │   └── app.go
    │   │
    │   └── tui/                 # TUI adapter (from internal/tui)
    │       ├── screens/
    │       └── app.go
    │
    ├── infrastructure/          # 🔧 EXTERNAL IMPLEMENTATIONS (Refactored)
    │   ├── persistence/         # Database layer
    │   │   ├── postgres/        # PostgreSQL implementations
    │   │   │   ├── user_repository.go      # Implements core/ports/repositories
    │   │   │   ├── org_repository.go
    │   │   │   ├── pipeline_repository.go
    │   │   │   ├── queries/                # From internal/database/queries
    │   │   │   ├── migrations/             # From internal/migrations/postgresql
    │   │   │   └── connection.go
    │   │   │
    │   │   ├── sqlite/          # SQLite implementations
    │   │   │   ├── repositories/
    │   │   │   └── migrations/             # From internal/migrations/sqlite
    │   │   │
    │   │   ├── memory/          # In-memory implementations (for testing)
    │   │   │   └── repositories/
    │   │   │
    │   │   └── sqlc/            # SQLC generated code
    │   │       ├── queries.sql.go
    │   │       └── models.sql.go
    │   │
    │   ├── cache/               # Caching implementations (from internal/cache & redis)
    │   │   ├── redis/
    │   │   │   ├── user_cache.go
    │   │   │   └── connection.go
    │   │   └── memory/
    │   │       └── cache.go
    │   │
    │   ├── events/              # Event bus implementations (from internal/events)
    │   │   ├── nats/
    │   │   │   └── publisher.go
    │   │   └── redis/
    │   │       └── publisher.go
    │   │
    │   ├── storage/             # File storage (from internal/storage)
    │   │   ├── s3/
    │   │   │   └── storage.go
    │   │   └── local/
    │   │       └── storage.go
    │   │
    │   ├── external/            # Third-party service implementations
    │   │   ├── openai/          # LLM implementation (from internal/llm)
    │   │   │   └── client.go
    │   │   ├── stripe/          # Payment service
    │   │   │   └── client.go
    │   │   ├── oauth/           # OAuth providers (from internal/auth/providers)
    │   │   │   ├── google.go
    │   │   │   ├── github.go
    │   │   │   └── microsoft.go
    │   │   └── email/           # Email service (from internal/auth/deliverers)
    │   │       ├── smtp.go
    │   │       └── sendgrid.go
    │   │
    │   └── config/              # Configuration (from internal/config)
    │       ├── loader.go
    │       └── types.go
    │
    ├── shared/                  # 🛠️ SHARED UTILITIES (Refactored)
    │   ├── logger/              # From internal/logger
    │   ├── testutil/            # From internal/testutil
    │   └── app/                 # From internal/app
    │
    └── generated/               # 🤖 GENERATED CODE (New location)
        ├── openapi/             # Generated from OpenAPI
        │   ├── types.gen.go     # Request/Response types
        │   └── server.gen.go    # Echo server interfaces
        ├── sqlc/                # Generated from SQL
        │   └── *.sql.go
        └── mocks/               # Generated mocks for testing
            └── *.mock.go
```

## How Current Packages Map to Clean Architecture

| Current Package            | New Location                                                | Purpose               |
| -------------------------- | ----------------------------------------------------------- | --------------------- |
| **internal/users**         | Split across multiple layers                                |                       |
| → service logic            | `core/entities/user.go`                                     | Domain entity         |
| → handler                  | `adapters/http/handlers/user_handler.go`                    | HTTP adapter          |
| → repository               | `infrastructure/persistence/postgres/user_repository.go`    | DB implementation     |
| → types.gen.go             | `generated/openapi/types.gen.go`                            | Generated types       |
| **internal/auth**          | Split across layers                                         |                       |
| → business logic           | `core/services/auth_service.go`                             | Domain service        |
| → providers                | `infrastructure/external/oauth/`                            | OAuth implementations |
| → stores                   | `infrastructure/persistence/postgres/session_repository.go` | Session storage       |
| → handler                  | `adapters/http/handlers/auth_handler.go`                    | HTTP adapter          |
| **internal/pipelines**     | Split across layers                                         |                       |
| → core logic               | `core/entities/pipeline.go`                                 | Domain entity         |
| → execution                | `core/services/pipeline_executor.go`                        | Domain service        |
| → handler                  | `adapters/http/handlers/pipeline_handler.go`                | HTTP adapter          |
| → queue                    | `infrastructure/events/redis/queue.go`                      | Queue implementation  |
| **internal/organizations** | Split across layers                                         |                       |
| → entity                   | `core/entities/organization.go`                             | Domain entity         |
| → aggregate                | `core/aggregates/workspace.go`                              | Workspace aggregate   |
| → handler                  | `adapters/http/handlers/org_handler.go`                     | HTTP adapter          |
| **internal/artifacts**     | Split across layers                                         |                       |
| → entity                   | `core/entities/artifact.go`                                 | Domain entity         |
| → storage logic            | `infrastructure/storage/`                                   | File storage          |
| **internal/database**      | `infrastructure/persistence/`                               | Database layer        |
| **internal/redis**         | `infrastructure/cache/redis/`                               | Cache implementation  |
| **internal/llm**           | `infrastructure/external/openai/`                           | LLM service           |
| **internal/server**        | `adapters/http/server/`                                     | HTTP server setup     |
| **internal/cli**           | `adapters/cli/`                                             | CLI adapter           |
| **internal/tui**           | `adapters/tui/`                                             | TUI adapter           |
| **internal/config**        | `infrastructure/config/`                                    | Configuration         |
| **internal/codegen**       | `tools/codegen/`                                            | Code generation tool  |
| **internal/templates**     | `tools/codegen/templates/`                                  | Generation templates  |
| **internal/migrations**    | `infrastructure/persistence/*/migrations/`                  | DB migrations         |
| **internal/events**        | `infrastructure/events/`                                    | Event bus             |
| **internal/cache**         | `infrastructure/cache/`                                     | Caching layer         |
| **internal/storage**       | `infrastructure/storage/`                                   | File storage          |
| **internal/middleware**    | `adapters/http/middleware/`                                 | HTTP middleware       |
| **internal/logger**        | `shared/logger/`                                            | Logging utility       |
| **internal/testutil**      | `shared/testutil/`                                          | Test utilities        |
| **internal/app**           | `shared/app/`                                               | App utilities         |

## Code Examples

### Core Layer - Domain Entity

```go
// internal/core/entities/user.go
package entities

import (
    "time"
    "github.com/archesai/archesai/internal/core/valueobjects"
    "github.com/archesai/archesai/internal/core/errors"
    "github.com/archesai/archesai/internal/core/events"
)

type User struct {
    id        valueobjects.UserID
    email     valueobjects.Email
    name      valueobjects.PersonName
    role      valueobjects.Role
    verified  bool
    createdAt time.Time
    updatedAt time.Time
    events    []events.DomainEvent
}

// NewUser creates a new user with validation
func NewUser(email valueobjects.Email, name valueobjects.PersonName) (*User, error) {
    user := &User{
        id:        valueobjects.NewUserID(),
        email:     email,
        name:      name,
        role:      valueobjects.RoleMember,
        verified:  false,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }

    user.addEvent(events.UserCreated{
        UserID: user.id,
        Email:  email,
        At:     user.createdAt,
    })

    return user, nil
}

// Business methods with rules
func (u *User) ChangeEmail(newEmail valueobjects.Email) error {
    if !u.verified {
        return errors.ErrUnverifiedUser
    }

    if u.email.Equals(newEmail) {
        return errors.ErrNoChange
    }

    oldEmail := u.email
    u.email = newEmail
    u.verified = false // Require re-verification
    u.updatedAt = time.Now()

    u.addEvent(events.EmailChanged{
        UserID:   u.id,
        OldEmail: oldEmail,
        NewEmail: newEmail,
        At:       u.updatedAt,
    })

    return nil
}

func (u *User) PromoteToAdmin() error {
    if u.role == valueobjects.RoleAdmin {
        return errors.ErrAlreadyAdmin
    }

    if !u.verified {
        return errors.ErrUnverifiedUser
    }

    u.role = valueobjects.RoleAdmin
    u.updatedAt = time.Now()

    u.addEvent(events.UserPromoted{
        UserID: u.id,
        Role:   u.role,
        At:     u.updatedAt,
    })

    return nil
}

// Getters (no setters - maintain encapsulation)
func (u *User) ID() valueobjects.UserID { return u.id }
func (u *User) Email() valueobjects.Email { return u.email }
func (u *User) Name() valueobjects.PersonName { return u.name }
func (u *User) Role() valueobjects.Role { return u.role }
func (u *User) IsVerified() bool { return u.verified }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }
func (u *User) Events() []events.DomainEvent { return u.events }

func (u *User) addEvent(event events.DomainEvent) {
    u.events = append(u.events, event)
}

// ReconstructUser is used by repository to reconstruct from DB
func ReconstructUser(
    id valueobjects.UserID,
    email valueobjects.Email,
    name valueobjects.PersonName,
    role valueobjects.Role,
    verified bool,
    createdAt, updatedAt time.Time,
) *User {
    return &User{
        id:        id,
        email:     email,
        name:      name,
        role:      role,
        verified:  verified,
        createdAt: createdAt,
        updatedAt: updatedAt,
        events:    []events.DomainEvent{},
    }
}
```

### Value Objects

```go
// internal/core/valueobjects/email.go
package valueobjects

import (
    "regexp"
    "strings"
    "github.com/archesai/archesai/internal/core/errors"
)

type Email struct {
    value string
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func NewEmail(s string) (Email, error) {
    s = strings.ToLower(strings.TrimSpace(s))

    if s == "" {
        return Email{}, errors.ErrEmptyEmail
    }

    if !emailRegex.MatchString(s) {
        return Email{}, errors.ErrInvalidEmail
    }

    return Email{value: s}, nil
}

func MustParseEmail(s string) Email {
    email, err := NewEmail(s)
    if err != nil {
        panic(err)
    }
    return email
}

func (e Email) String() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

func (e Email) Domain() string {
    parts := strings.Split(e.value, "@")
    if len(parts) == 2 {
        return parts[1]
    }
    return ""
}
```

### Ports (Interfaces defined by domain)

```go
// internal/core/ports/repositories/user_repository.go
package repositories

import (
    "context"
    "github.com/archesai/archesai/internal/core/entities"
    "github.com/archesai/archesai/internal/core/valueobjects"
)

type UserRepository interface {
    // Save creates or updates a user
    Save(ctx context.Context, user *entities.User) error

    // FindByID retrieves a user by ID
    FindByID(ctx context.Context, id valueobjects.UserID) (*entities.User, error)

    // FindByEmail retrieves a user by email
    FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error)

    // List retrieves users with pagination
    List(ctx context.Context, offset, limit int) ([]*entities.User, int64, error)

    // Delete removes a user
    Delete(ctx context.Context, id valueobjects.UserID) error
}
```

### Application Layer - Use Case

```go
// internal/application/commands/users/create_user.go
package users

import (
    "context"
    "github.com/archesai/archesai/internal/core/entities"
    "github.com/archesai/archesai/internal/core/valueobjects"
    "github.com/archesai/archesai/internal/core/errors"
    "github.com/archesai/archesai/internal/core/ports/repositories"
    "github.com/archesai/archesai/internal/core/ports/events"
    "github.com/archesai/archesai/internal/application/dto"
)

type CreateUserCommand struct {
    Email string
    Name  string
}

type CreateUserHandler struct {
    userRepo repositories.UserRepository
    events   events.EventPublisher
}

func NewCreateUserHandler(
    userRepo repositories.UserRepository,
    events events.EventPublisher,
) *CreateUserHandler {
    return &CreateUserHandler{
        userRepo: userRepo,
        events:   events,
    }
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*dto.UserDTO, error) {
    // Validate and create value objects
    email, err := valueobjects.NewEmail(cmd.Email)
    if err != nil {
        return nil, err
    }

    name, err := valueobjects.NewPersonName(cmd.Name)
    if err != nil {
        return nil, err
    }

    // Check if user exists
    existing, _ := h.userRepo.FindByEmail(ctx, email)
    if existing != nil {
        return nil, errors.ErrUserAlreadyExists
    }

    // Create domain entity
    user, err := entities.NewUser(email, name)
    if err != nil {
        return nil, err
    }

    // Persist
    if err := h.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }

    // Publish domain events
    for _, event := range user.Events() {
        if err := h.events.Publish(ctx, event); err != nil {
            // Log but don't fail
        }
    }

    // Return DTO
    return dto.FromUser(user), nil
}
```

### Infrastructure - Repository Implementation

```go
// internal/infrastructure/persistence/postgres/user_repository.go
package postgres

import (
    "context"
    "database/sql"
    "github.com/google/uuid"
    "github.com/archesai/archesai/internal/core/entities"
    "github.com/archesai/archesai/internal/core/valueobjects"
    "github.com/archesai/archesai/internal/core/ports/repositories"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *entities.User) error {
    query := `
        INSERT INTO users (id, email, name, role, verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            name = EXCLUDED.name,
            role = EXCLUDED.role,
            verified = EXCLUDED.verified,
            updated_at = EXCLUDED.updated_at
    `

    _, err := r.db.ExecContext(ctx, query,
        user.ID().String(),
        user.Email().String(),
        user.Name().String(),
        user.Role().String(),
        user.IsVerified(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )

    return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error) {
    query := `
        SELECT id, email, name, role, verified, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    var id string
    var emailStr string
    var name string
    var role string
    var verified bool
    var createdAt, updatedAt time.Time

    err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
        &id,
        &emailStr,
        &name,
        &role,
        &verified,
        &createdAt,
        &updatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    // Reconstruct domain entity
    userID, _ := valueobjects.ParseUserID(id)
    userEmail := valueobjects.MustParseEmail(emailStr)
    userName := valueobjects.MustParsePersonName(name)
    userRole := valueobjects.ParseRole(role)

    return entities.ReconstructUser(
        userID,
        userEmail,
        userName,
        userRole,
        verified,
        createdAt,
        updatedAt,
    ), nil
}

func (r *UserRepository) FindByID(ctx context.Context, id valueobjects.UserID) (*entities.User, error) {
    query := `
        SELECT id, email, name, role, verified, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var idStr string
    var email string
    var name string
    var role string
    var verified bool
    var createdAt, updatedAt time.Time

    err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
        &idStr,
        &email,
        &name,
        &role,
        &verified,
        &createdAt,
        &updatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    // Reconstruct domain entity
    userEmail := valueobjects.MustParseEmail(email)
    userName := valueobjects.MustParsePersonName(name)
    userRole := valueobjects.ParseRole(role)

    return entities.ReconstructUser(
        id,
        userEmail,
        userName,
        userRole,
        verified,
        createdAt,
        updatedAt,
    ), nil
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entities.User, int64, error) {
    // Count total
    var total int64
    countQuery := `SELECT COUNT(*) FROM users`
    if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
        return nil, 0, err
    }

    // Get page
    query := `
        SELECT id, email, name, role, verified, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var users []*entities.User
    for rows.Next() {
        var id string
        var email string
        var name string
        var role string
        var verified bool
        var createdAt, updatedAt time.Time

        if err := rows.Scan(&id, &email, &name, &role, &verified, &createdAt, &updatedAt); err != nil {
            return nil, 0, err
        }

        userID, _ := valueobjects.ParseUserID(id)
        userEmail := valueobjects.MustParseEmail(email)
        userName := valueobjects.MustParsePersonName(name)
        userRole := valueobjects.ParseRole(role)

        user := entities.ReconstructUser(
            userID,
            userEmail,
            userName,
            userRole,
            verified,
            createdAt,
            updatedAt,
        )

        users = append(users, user)
    }

    return users, total, nil
}

func (r *UserRepository) Delete(ctx context.Context, id valueobjects.UserID) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id.String())
    return err
}
```

### Adapter - HTTP Handler

```go
// internal/adapters/http/handlers/user_handler.go
package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/archesai/archesai/internal/application/commands/users"
    "github.com/archesai/archesai/internal/application/queries/users"
    "github.com/archesai/archesai/internal/core/errors"
)

type UserHandler struct {
    createUser *users.CreateUserHandler
    getUser    *users.GetUserHandler
    listUsers  *users.ListUsersHandler
}

func NewUserHandler(
    createUser *users.CreateUserHandler,
    getUser *users.GetUserHandler,
    listUsers *users.ListUsersHandler,
) *UserHandler {
    return &UserHandler{
        createUser: createUser,
        getUser:    getUser,
        listUsers:  listUsers,
    }
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    // Validate request
    if err := req.Validate(); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    // Execute use case
    user, err := h.createUser.Handle(c.Request().Context(), users.CreateUserCommand{
        Email: req.Email,
        Name:  req.Name,
    })

    if err != nil {
        return h.handleError(err)
    }

    // Map to response
    return c.JSON(http.StatusCreated, CreateUserResponse{
        Data: UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            Name:      user.Name,
            Role:      user.Role,
            Verified:  user.Verified,
            CreatedAt: user.CreatedAt,
        },
    })
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c echo.Context) error {
    id := c.Param("id")

    user, err := h.getUser.Handle(c.Request().Context(), users.GetUserQuery{
        UserID: id,
    })

    if err != nil {
        return h.handleError(err)
    }

    return c.JSON(http.StatusOK, GetUserResponse{
        Data: UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            Name:      user.Name,
            Role:      user.Role,
            Verified:  user.Verified,
            CreatedAt: user.CreatedAt,
        },
    })
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c echo.Context) error {
    page := c.QueryParam("page")
    size := c.QueryParam("size")

    users, total, err := h.listUsers.Handle(c.Request().Context(), users.ListUsersQuery{
        Page: page,
        Size: size,
    })

    if err != nil {
        return h.handleError(err)
    }

    response := ListUsersResponse{
        Data: make([]UserResponse, len(users)),
        Meta: MetaResponse{
            Total: total,
        },
    }

    for i, user := range users {
        response.Data[i] = UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            Name:      user.Name,
            Role:      user.Role,
            Verified:  user.Verified,
            CreatedAt: user.CreatedAt,
        }
    }

    return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) handleError(err error) error {
    switch err {
    case errors.ErrUserNotFound:
        return echo.NewHTTPError(http.StatusNotFound, "User not found")
    case errors.ErrUserAlreadyExists:
        return echo.NewHTTPError(http.StatusConflict, "User already exists")
    case errors.ErrInvalidEmail:
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid email format")
    case errors.ErrUnverifiedUser:
        return echo.NewHTTPError(http.StatusForbidden, "User not verified")
    default:
        return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
    }
}
```

## Migration Strategy

### Phase 1: Foundation (Week 1)

1. Create directory structure
2. Define core entities and value objects for one domain (start with `labels` as it's simple)
3. Create port interfaces
4. Set up dependency injection container

### Phase 2: Pilot Migration - Labels Domain (Week 2)

1. Create `core/entities/label.go`
2. Create `core/valueobjects/label_name.go`
3. Create `core/ports/repositories/label_repository.go`
4. Create `application/commands/labels/create_label.go`
5. Create `infrastructure/persistence/postgres/label_repository.go`
6. Create `adapters/http/handlers/label_handler.go`
7. Update tests to use new structure

### Phase 3: Core Domains Migration (Weeks 3-4)

1. **Users Domain**:
   - Entity: User with business rules
   - Value Objects: Email, UserID, PersonName
   - Use Cases: CreateUser, UpdateUser, VerifyEmail

2. **Organizations Domain**:
   - Entity: Organization
   - Aggregate: Workspace (Org + Members)
   - Use Cases: CreateOrg, InviteMember, UpdateSettings

3. **Pipelines Domain**:
   - Entity: Pipeline
   - Value Objects: PipelineStatus, StepType
   - Domain Service: PipelineExecutor
   - Use Cases: CreatePipeline, ExecutePipeline

### Phase 4: Auth Refactoring (Week 5)

1. Split monolithic auth service:
   - `core/services/auth_service.go` - Core authentication logic
   - `core/services/permission_service.go` - Authorization rules
   - `application/commands/auth/login.go` - Login use case
   - `infrastructure/external/oauth/` - OAuth providers

2. Create proper session aggregate:
   - `core/aggregates/session.go` - User session with tokens

### Phase 5: Infrastructure Layer (Week 6)

1. Move database code to `infrastructure/persistence/`
2. Move cache to `infrastructure/cache/`
3. Move external services to `infrastructure/external/`
4. Update dependency injection

### Phase 6: Cleanup (Week 7)

1. Remove old package structure
2. Update code generation templates
3. Update documentation
4. Performance testing

## Code Generation Strategy

### What Gets Generated

```yaml
# .archesai.codegen.yaml
generators:
  # Generate HTTP DTOs and interfaces from OpenAPI
  openapi:
    input: api/openapi.yaml
    outputs:
      types: internal/generated/openapi/types.gen.go
      server: internal/generated/openapi/server.gen.go

  # Generate SQL models from schema
  sqlc:
    input: internal/infrastructure/persistence/postgres/queries/
    output: internal/generated/sqlc/

  # Generate mocks for testing
  mockery:
    interfaces:
      - internal/core/ports/**/*.go
    output: internal/generated/mocks/
```

### What Remains Hand-Written

- **`core/`** - All domain logic, entities, value objects, domain services
- **`application/`** - All use cases and application logic
- **`adapters/`** - Handler logic (uses generated interfaces)
- **`infrastructure/`** - Repository implementations (uses generated SQL)

## Testing Strategy

### Core Layer Testing (Pure Unit Tests - No Mocks)

```go
// internal/core/entities/user_test.go
func TestUser_ChangeEmail(t *testing.T) {
    // Create user with valid email
    email := valueobjects.MustParseEmail("old@example.com")
    name := valueobjects.MustParsePersonName("John Doe")
    user, _ := entities.NewUser(email, name)

    // Verify the user
    user.Verify()

    // Change email
    newEmail := valueobjects.MustParseEmail("new@example.com")
    err := user.ChangeEmail(newEmail)

    assert.NoError(t, err)
    assert.Equal(t, newEmail, user.Email())
    assert.False(t, user.IsVerified()) // Should require re-verification
    assert.Len(t, user.Events(), 2) // UserCreated + EmailChanged
}

func TestUser_CannotChangeEmailWhenUnverified(t *testing.T) {
    email := valueobjects.MustParseEmail("test@example.com")
    name := valueobjects.MustParsePersonName("John Doe")
    user, _ := entities.NewUser(email, name)

    // Try to change email without verification
    newEmail := valueobjects.MustParseEmail("new@example.com")
    err := user.ChangeEmail(newEmail)

    assert.Equal(t, errors.ErrUnverifiedUser, err)
}
```

### Application Layer Testing (Mock Only Ports)

```go
// internal/application/commands/users/create_user_test.go
func TestCreateUserHandler(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    mockEvents := mocks.NewMockEventPublisher()

    handler := NewCreateUserHandler(mockRepo, mockEvents)

    // Setup expectations
    mockRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, nil)
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    mockEvents.On("Publish", mock.Anything, mock.Anything).Return(nil)

    // Execute
    result, err := handler.Handle(context.Background(), CreateUserCommand{
        Email: "test@example.com",
        Name:  "Test User",
    })

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test@example.com", result.Email)
    mockRepo.AssertExpectations(t)
    mockEvents.AssertExpectations(t)
}
```

### Infrastructure Layer Testing (Integration Tests)

```go
// internal/infrastructure/persistence/postgres/user_repository_test.go
func TestUserRepository_Save(t *testing.T) {
    // Use test database
    db := setupTestDB(t)
    defer db.Close()

    repo := NewUserRepository(db)

    // Create domain entity
    email := valueobjects.MustParseEmail("test@example.com")
    name := valueobjects.MustParsePersonName("Test User")
    user, _ := entities.NewUser(email, name)

    // Save
    err := repo.Save(context.Background(), user)

    // Assert
    assert.NoError(t, err)

    // Verify by finding
    found, err := repo.FindByID(context.Background(), user.ID())
    assert.NoError(t, err)
    assert.Equal(t, user.Email(), found.Email())
}
```

## Benefits of This Structure

1. **True Independence**: Core business logic has zero dependencies on frameworks, databases, or external services
2. **Testability**: Core can be unit tested without any mocks - it's pure business logic
3. **Flexibility**: Easy to swap implementations (PostgreSQL → MongoDB, Echo → Gin, etc.)
4. **Clear Boundaries**: Each layer has a specific responsibility
5. **Code Generation Compatibility**: You can still use OpenAPI generation for the adapters layer while keeping core clean
6. **Gradual Migration**: Can be done incrementally, one domain at a time
7. **Maintainability**: Business rules are centralized in the core layer
8. **Documentation**: The code structure itself documents the business domain

## Key Principles to Follow

1. **Dependency Rule**: Dependencies only point inward. Core knows nothing about outer layers.
2. **Interface Segregation**: Ports are small, focused interfaces defined by the domain.
3. **Single Responsibility**: Each layer has one reason to change.
4. **Don't Generate Core**: Business logic should always be hand-written and carefully crafted.
5. **Test the Right Thing**: Unit test the core, integration test the infrastructure.

## Next Steps

1. Start with the `labels` domain as a pilot (it's the simplest)
2. Create the directory structure
3. Implement core entities and value objects
4. Create application use cases
5. Adapt existing handlers to use the new structure
6. Gradually migrate other domains
7. Update code generation to work with new structure

This architecture will give you a maintainable, testable, and flexible codebase that can grow with your needs while still leveraging code generation where it makes sense.
