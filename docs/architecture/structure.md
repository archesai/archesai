# Application Structure

This document describes the target architecture for Arches-generated applications.

## Overview

Arches generates applications organized by **OpenAPI tag** (feature/domain), not by layer. This keeps
related code together and follows Go idioms.

```text
app/
  pipelines/                   # Tag: Pipelines
    types.gen.go               # Entities, interfaces, DTOs
    http.gen.go                # HTTP handlers
    create_pipeline.go         # User implementation
    get_pipeline.go            # User implementation

  users/                       # Tag: Users
    types.gen.go
    http.gen.go
    create_user.go

  adapters/
    postgres/
      migrations/
        001_create_pipelines.sql
        002_create_users.sql
      pipelines.gen.go         # PipelineRepository implementation
      users.gen.go             # UserRepository implementation
    sqlite/
      migrations/
      pipelines.gen.go
      users.gen.go

  main.gen.go                  # Bootstrap and wiring
```

## Package Structure

### Feature Packages (`{tag}/`)

Each OpenAPI tag becomes a Go package containing everything for that feature.

#### `types.gen.go` - Generated types

```go
package pipelines

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/archesai/archesai/pkg/executor"
)

// ============================================================================
// Entities
// ============================================================================

type Pipeline struct {
    ID          uuid.UUID
    Name        string
    Description *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type PipelineStep struct {
    ID         uuid.UUID
    PipelineID uuid.UUID
    Name       string
    ToolID     uuid.UUID
    Position   int32
}

// ============================================================================
// Repositories
// ============================================================================

type PipelineRepository interface {
    Save(ctx context.Context, p *Pipeline) error
    FindByID(ctx context.Context, id uuid.UUID) (*Pipeline, error)
    List(ctx context.Context, filter PipelineFilter) ([]Pipeline, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type PipelineStepRepository interface {
    Save(ctx context.Context, s *PipelineStep) error
    FindByPipelineID(ctx context.Context, pipelineID uuid.UUID) ([]PipelineStep, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

// ============================================================================
// CreatePipeline
// ============================================================================

type CreatePipelineInput struct {
    Name        string
    Description *string
}

type CreatePipelineOutput struct {
    Data Pipeline
}

type CreatePipeline = executor.Executor[CreatePipelineInput, CreatePipelineOutput]

// ============================================================================
// GetPipeline
// ============================================================================

type GetPipelineInput struct {
    ID uuid.UUID
}

type GetPipelineOutput struct {
    Data Pipeline
}

type GetPipeline = executor.Executor[GetPipelineInput, GetPipelineOutput]

// ============================================================================
// ListPipelines
// ============================================================================

type ListPipelinesInput struct {
    Filter *PipelineFilter
    Page   *int
    Limit  *int
}

type ListPipelinesOutput struct {
    Data []Pipeline
    Meta PaginationMeta
}

type ListPipelines = executor.Executor[ListPipelinesInput, ListPipelinesOutput]
```

#### `http.gen.go` - HTTP handlers

```go
package pipelines

import (
    "errors"
    "log/slog"
    "net/http"

    "github.com/archesai/archesai/pkg/httputil"
    "github.com/archesai/archesai/pkg/validation"
)

// ============================================================================
// CreatePipeline - POST /pipelines
// ============================================================================

type CreatePipelineHTTPHandler struct {
    executor CreatePipeline
}

func NewCreatePipelineHTTPHandler(executor CreatePipeline) *CreatePipelineHTTPHandler {
    return &CreatePipelineHTTPHandler{executor: executor}
}

func RegisterCreatePipelineRoute(mux *http.ServeMux, handler *CreatePipelineHTTPHandler) {
    mux.HandleFunc("POST /pipelines", handler.ServeHTTP)
}

// ----------------------------------------------------------------------------
// Request Body
// ----------------------------------------------------------------------------

type CreatePipelineRequestBody struct {
    Name        string  `json:"name"`
    Description *string `json:"description,omitempty"`
}

func (b *CreatePipelineRequestBody) Validate() validation.Errors {
    var errs validation.Errors
    validation.RequiredString(b.Name, "name", &errs)
    validation.MinLengthString(b.Name, 1, "name", &errs)
    validation.MaxLengthString(b.Name, 255, "name", &errs)
    return errs
}

// ----------------------------------------------------------------------------
// Response Types
// ----------------------------------------------------------------------------

type CreatePipeline201Response struct {
    Data Pipeline `json:"data"`
}

func (CreatePipeline201Response) StatusCode() int { return 201 }

type CreatePipeline400Response struct {
    httputil.ProblemDetails
}

func (CreatePipeline400Response) StatusCode() int { return 400 }

type CreatePipeline500Response struct {
    httputil.ProblemDetails
}

func (CreatePipeline500Response) StatusCode() int { return 500 }

// ----------------------------------------------------------------------------
// Handler
// ----------------------------------------------------------------------------

func (h *CreatePipelineHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var body CreatePipelineRequestBody
    if err := httputil.DecodeAndValidate(r, &body); err != nil {
        var validErrs validation.Errors
        if errors.As(err, &validErrs) {
            _ = httputil.WriteValidationErrors(w, validErrs, r.URL.Path)
            return
        }
        _ = httputil.WriteResponse(w, CreatePipeline400Response{
            ProblemDetails: httputil.NewBadRequestResponse(err.Error(), r.URL.Path),
        })
        return
    }

    input := &CreatePipelineInput{
        Name:        body.Name,
        Description: body.Description,
    }

    result, err := h.executor.Execute(ctx, input)
    if err != nil {
        var httpErr httputil.HTTPError
        if errors.As(err, &httpErr) {
            _ = httputil.WriteHTTPError(w, httpErr, r.URL.Path)
            return
        }
        slog.Error("handler error", "operation", "CreatePipeline", "error", err)
        _ = httputil.WriteResponse(w, CreatePipeline500Response{
            ProblemDetails: httputil.NewInternalServerErrorResponse(err.Error(), r.URL.Path),
        })
        return
    }

    _ = httputil.WriteResponse(w, CreatePipeline201Response{
        Data: result.Data,
    })
}
```

#### User implementation files

Users create implementation files alongside the generated code:

```go
// pipelines/create_pipeline.go

package pipelines

import (
    "context"
    "time"

    "github.com/google/uuid"
)

type createPipeline struct {
    repo PipelineRepository
}

func NewCreatePipeline(repo PipelineRepository) *createPipeline {
    return &createPipeline{repo: repo}
}

func (c *createPipeline) Execute(ctx context.Context, input *CreatePipelineInput) (*CreatePipelineOutput, error) {
    pipeline := &Pipeline{
        ID:          uuid.New(),
        Name:        input.Name,
        Description: input.Description,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := c.repo.Save(ctx, pipeline); err != nil {
        return nil, err
    }

    return &CreatePipelineOutput{Data: *pipeline}, nil
}
```

### Adapters (`adapters/`)

Database adapters are centralized because:

- Migrations must be sequential across all entities
- SQLC generates from all queries together
- Connection pooling is shared

#### `adapters/postgres/`

```text
adapters/postgres/
  migrations/
    001_create_pipelines.sql
    002_create_pipeline_steps.sql
    003_create_users.sql
  db.go                    # Connection setup
  pipelines.gen.go         # PipelineRepository, PipelineStepRepository
  users.gen.go             # UserRepository
```

```go
// adapters/postgres/pipelines.gen.go

package postgres

import (
    "context"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"

    "myapp/pipelines"
)

type PipelineRepository struct {
    db *pgxpool.Pool
}

func NewPipelineRepository(db *pgxpool.Pool) *PipelineRepository {
    return &PipelineRepository{db: db}
}

func (r *PipelineRepository) Save(ctx context.Context, p *pipelines.Pipeline) error {
    _, err := r.db.Exec(ctx, `
        INSERT INTO pipelines (id, name, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            description = EXCLUDED.description,
            updated_at = EXCLUDED.updated_at
    `, p.ID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt)
    return err
}

func (r *PipelineRepository) FindByID(ctx context.Context, id uuid.UUID) (*pipelines.Pipeline, error) {
    var p pipelines.Pipeline
    err := r.db.QueryRow(ctx, `
        SELECT id, name, description, created_at, updated_at
        FROM pipelines WHERE id = $1
    `, id).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return &p, nil
}

// ... List, Delete, etc.
```

#### `adapters/sqlite/`

Same structure, different SQL dialect:

```text
adapters/sqlite/
  migrations/
  db.go
  pipelines.gen.go
  users.gen.go
```

### Bootstrap (`main.gen.go`)

Wires everything together:

```go
// main.gen.go

package main

import (
    "log"
    "net/http"

    "myapp/adapters/postgres"
    "myapp/pipelines"
    "myapp/users"
)

func main() {
    // Database
    db, err := postgres.NewDB("postgres://...")
    if err != nil {
        log.Fatal(err)
    }

    // Repositories
    pipelineRepo := postgres.NewPipelineRepository(db)
    userRepo := postgres.NewUserRepository(db)

    // Executors
    createPipeline := pipelines.NewCreatePipeline(pipelineRepo)
    getPipeline := pipelines.NewGetPipeline(pipelineRepo)
    createUser := users.NewCreateUser(userRepo)

    // HTTP Handlers
    createPipelineHTTP := pipelines.NewCreatePipelineHTTPHandler(createPipeline)
    getPipelineHTTP := pipelines.NewGetPipelineHTTPHandler(getPipeline)
    createUserHTTP := users.NewCreateUserHTTPHandler(createUser)

    // Routes
    mux := http.NewServeMux()
    pipelines.RegisterCreatePipelineRoute(mux, createPipelineHTTP)
    pipelines.RegisterGetPipelineRoute(mux, getPipelineHTTP)
    users.RegisterCreateUserRoute(mux, createUserHTTP)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## File Naming Conventions

| File        | Description                       |
| ----------- | --------------------------------- |
| `*.gen.go`  | Generated by Arches - do not edit |
| `*.go`      | User-written code                 |
| `*_test.go` | Tests                             |

## What Gets Generated

| Source                    | Generated                                     |
| ------------------------- | --------------------------------------------- |
| OpenAPI tag               | `{tag}/` package                              |
| Schema with `x-entity`    | Entity struct in `types.gen.go`               |
| Schema (response/request) | DTO struct in `types.gen.go`                  |
| Operation                 | Executor type alias + HTTP handler            |
| Entity with DB            | Repository interface + adapter implementation |
| All entities              | Migrations in `adapters/{db}/migrations/`     |

## Comparison with Current Structure

### Current (layer-based)

```text
app/
  handlers/        # Interfaces + DTOs
  routes/          # HTTP handlers
  models/          # Entities
  database/        # Migrations, queries
```

### New (feature-based)

```text
app/
  pipelines/       # Everything for pipelines
  users/           # Everything for users
  adapters/        # Database implementations
```

### Benefits

1. **Cohesion** - Related code lives together
2. **Go-idiomatic** - Flat packages, clear boundaries
3. **Easier navigation** - Find all pipeline code in `pipelines/`
4. **Clear ownership** - Each tag is self-contained
5. **Simpler imports** - `import "myapp/pipelines"` gets everything

## Migration Path

1. Generate new structure alongside old
2. Move user implementations to new packages
3. Update imports
4. Delete old packages

Or regenerate fresh and re-implement business logic.
