# Code Generators Documentation

ArchesAI uses a comprehensive code generation strategy to maintain type safety, reduce boilerplate, and ensure consistency across the codebase.

## Overview

The project employs four primary code generators:

1. **sqlc** - Database queries to Go code
2. **oapi-codegen** - OpenAPI spec to Go server code
3. **generate-defaults** - OpenAPI schema to Go config defaults
4. **generate-converters** - YAML config to type converters

## Generator Details

### 1. sqlc - Database Query Generation

**Purpose**: Generate type-safe Go code from SQL queries

**Configuration**: `internal/infrastructure/database/sqlc.yaml`

```yaml
version: '2'
sql:
  - engine: 'postgresql'
    queries: './queries'
    schema: './migrations/postgresql'
    gen:
      go:
        package: 'postgresql'
        out: '../generated/database/postgresql'
        emit_json_tags: true
        emit_empty_slices: true
```

**Input Files**: `internal/infrastructure/database/queries/*.sql`

```sql
-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (email, name, email_verified, image)
VALUES ($1, $2, $3, $4)
RETURNING *;
```

**Generated Output**: `internal/generated/database/postgresql/`

- `models.go` - Struct definitions for tables
- `queries.sql.go` - Type-safe query functions
- `querier.go` - Interface for all queries

**Usage in Code**:

```go
user, err := container.Queries.GetUserByEmail(ctx, email)
```

### 2. oapi-codegen - OpenAPI to Go Server

**Purpose**: Generate Go server interfaces and types from OpenAPI specification

**Configuration**: `internal/generated/api/generate.go`

```go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
//  --config=oapi-codegen.yaml \
//  --package=api \
//  ../../../api/openapi.bundled.yaml
```

**Input**: `api/openapi.bundled.yaml` (bundled from component files)

**Generated Output**: `internal/generated/api/`

- `api.gen.go` - Server interfaces and types

**Implementation Pattern**:

```go
// Generated interface (DO NOT EDIT)
type ServerInterface interface {
    PostAuthSignIn(ctx echo.Context) error
    GetAuthUsers(ctx echo.Context, params GetAuthUsersParams) error
}

// Domain implementation
type Handler struct {
    service *Service
}

func (h *Handler) PostAuthSignIn(ctx echo.Context) error {
    // Implementation
}
```

### 3. generate-defaults - Config Defaults from OpenAPI

**Purpose**: Generate Go config struct with default values from OpenAPI schema

**Source Code**: `cmd/generate-defaults/main.go`

**Input**: OpenAPI schema components

```yaml
components:
  schemas:
    ArchesConfig:
      properties:
        api:
          properties:
            port:
              type: integer
              default: 3001
            host:
              type: string
              default: '0.0.0.0'
```

**Output**: `internal/infrastructure/config/defaults.gen.go`

```go
func GetDefaultConfig() *api.ArchesConfig {
    config := &api.ArchesConfig{
        Api: api.APIConfig{
            Port: 3001,
            Host: "0.0.0.0",
        },
    }
    return config
}
```

### 4. generate-converters - Type Converter Generation

**Purpose**: Generate converters between database and API types

**Source Code**: `cmd/generate-converters/main.go`

**Configuration**: `internal/domains/converters.yaml`

```yaml
converters:
  - name: PipelineDBToAPI
    from: postgresql.Pipeline
    to: api.PipelineEntity
    domain: workflows
    automap: true # Automatically map fields with same names
    fields:
      # Only specify custom conversions
      OrganizationId: 'openapi_types.UUID(uuid.MustParse(from.OrganizationId))'

  - name: AuthUserDBToAPI
    from: postgresql.User
    to: api.UserEntity
    domain: auth
    automap: true
    fields:
      Email: 'openapi_types.Email(from.Email)'
      Id: 'openapi_types.UUID(uuid.MustParse(from.Id))'
```

**Features**:

- **Automap**: Automatically maps fields with matching names and compatible types
- **Type Awareness**: Handles nullable fields, UUIDs, timestamps
- **Deterministic Output**: Fields are alphabetically sorted for consistent generation
- **Helper Functions**: Generates reusable helpers for common conversions

**Generated Output**: `internal/domains/{domain}/converters/converters.gen.go`

```go
// Helper functions
func handleNullableString(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

// PipelineDBToAPI converts postgresql.Pipeline to api.PipelineEntity
func PipelineDBToAPI(from *postgresql.Pipeline) api.PipelineEntity {
    return api.PipelineEntity{
        CreatedAt:      from.CreatedAt,
        Description:    handleNullableString(from.Description),
        Id:             openapi_types.UUID(uuid.MustParse(from.Id)),
        Name:           handleNullableString(from.Name),
        OrganizationId: openapi_types.UUID(uuid.MustParse(from.OrganizationId)),
        UpdatedAt:      from.UpdatedAt,
    }
}
```

## Domain-Specific Generator Usage

Each domain has its own converter configuration:

### Auth Domain

- User entities with email verification and sessions
- Account linking for OAuth providers
- API token management

### Organizations Domain

- Organization membership management
- Member role conversions
- Invitation status handling

### Workflows Domain

- Pipeline and run status enums
- Tool configuration objects
- Step dependency relationships

### Content Domain

- Artifact type classifications
- Label relationship mappings
- Content metadata handling

## Running Generators

### Individual Generators

```bash
make sqlc                # Generate database code
make oapi                # Generate OpenAPI server code
make generate-defaults   # Generate config defaults
make generate-converters # Generate type converters
```

### All Generators

```bash
make generate # Runs all generators in correct order
```

## Adding New Converters

1. Add converter spec to `internal/domains/converters.yaml`
2. Set `automap: true` to automatically map matching fields
3. Only specify fields that need custom conversion logic
4. Run `make generate-converters`

## Best Practices

### 1. Never Edit Generated Files

Files in these directories are regenerated:

- `internal/generated/`
- `*/converters/converters.gen.go`
- `*.gen.go` files

### 2. Use Automap Feature

When adding converters, use `automap: true` to reduce configuration:

```yaml
converters:
  - name: MyConverter
    from: postgresql.MyTable
    to: api.MyEntity
    automap: true # Handles 90% of fields automatically
    fields:
      # Only special cases here
      SpecialField: 'customConversion(from.SpecialField)'
```

### 3. Keep Generators Deterministic

- Sort fields alphabetically
- Use consistent naming patterns
- Avoid random or time-based values

### 4. Type Safety First

- Let generators handle type conversions
- Use compile-time checks where possible
- Avoid `interface{}` in generated code

## Extending the Generator System

### Creating a New Generator

Example structure for a domain generator:

```go
// cmd/generate-domain/main.go
package main

import (
    "flag"
    "text/template"
)

type DomainSpec struct {
    Name       string
    Tables     []string
    Endpoints  []string
}

func main() {
    var spec DomainSpec
    flag.StringVar(&spec.Name, "name", "", "Domain name")
    flag.Parse()

    // Generate domain files
    generateDomainFiles(spec)
}

func generateDomainFiles(spec DomainSpec) {
    files := map[string]string{
        "auth.go":       domainTemplate,
        "entities.go":   entitiesTemplate,
        "service.go":    serviceTemplate,
        "repository.go": repositoryTemplate,
        "handler.go":    handlerTemplate,
    }

    for filename, tmpl := range files {
        // Generate each file
    }
}
```

### Generator Pipeline

The current generator pipeline for ArchesAI's four-domain architecture:

1. **OpenAPI Spec** → API types & interfaces (auth, organizations, workflows, content)
2. **SQL Migrations** → Database models (PostgreSQL and SQLite support)
3. **Converter Config** → Domain-specific type converters
4. **Default Config** → Configuration with sensible defaults

**Domain Order**: Auth → Organizations → Workflows → Content

## Troubleshooting

### Common Issues

1. **Import Cycles**: Keep generated code in separate packages
2. **Type Mismatches**: Check nullable fields in database vs API
3. **Missing Fields**: Verify automap compatibility or add explicit mapping
4. **Compilation Errors**: Run `make generate` after schema changes

### Debugging Tips

- Check generator logs for warnings
- Verify input files are valid (YAML, SQL, OpenAPI)
- Use `go generate -x` to see actual commands
- Compare generated output with previous versions

## Future Enhancements

Planned improvements to the generator system:

1. **Domain Generator**: Scaffold complete domains from configuration
2. **Migration Generator**: Generate migrations from struct definitions
3. **Test Generator**: Generate test cases from OpenAPI examples
4. **Client Generator**: Generate TypeScript client from OpenAPI
5. **Documentation Generator**: Generate API docs from OpenAPI + code comments
