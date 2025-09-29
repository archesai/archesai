# Code Generation Guide

## Overview

Arches uses a unified code generation system that reads OpenAPI specifications with x-codegen annotations to
automatically generate repository interfaces, database implementations, service interfaces, HTTP servers, and
test mocks. This approach ensures type safety, reduces boilerplate, and maintains consistency across the codebase.

## Quick Start

```bash
# Generate all code
make generate

# Generate only from OpenAPI/x-codegen
make generate-codegen

# Generate only SQLC database code
make generate-sqlc

# Generate only OpenAPI types
make generate-oapi

# Generate test mocks
make generate-mocks
```

## How It Works

### 1. Define OpenAPI Schema with x-codegen

Add x-codegen annotations to your OpenAPI schema:

```yaml
components:
  schemas:
    Organization:
      description: Organization entity
      x-codegen:
        # Repository generation configuration
        repository:
          operations:
            - create
            - get
            - update
            - delete
            - list
            - getByName
          database:
            table: organizations
            postgres: true
            sqlite: true

        # Service generation configuration
        service:
          enabled: true

        # Events configuration (optional)
        events:
          - created
          - updated
          - deleted

      properties:
        id:
          type: string
          format: uuid
          x-codegen:
            database:
              column: id
              primaryKey: true
        name:
          type: string
          maxLength: 255
          x-codegen:
            database:
              column: name
              unique: true
        createdAt:
          type: string
          format: date-time
          x-codegen:
            database:
              column: created_at
              default: CURRENT_TIMESTAMP
```

### 2. Run Code Generation

```bash
make generate-codegen
```

### 3. Generated Files

The generator creates the following files in each domain:

- `repository.gen.go` - Repository interface with CRUD operations
- `postgres.gen.go` - PostgreSQL implementation of the repository
- `sqlite.gen.go` - SQLite implementation of the repository
- `service.gen.go` - Service interface for business logic
- `server.gen.go` - HTTP server implementation
- `types.gen.go` - OpenAPI type definitions
- `handler.gen.go` - API client interface
- `mocks_test.gen.go` - Test mocks for all interfaces

## x-codegen Annotations

### Repository Configuration

```yaml
x-codegen:
  repository:
    # Standard CRUD operations
    operations:
      - create       # Create{Entity}(ctx, entity) error
      - get          # Get{Entity}(ctx, id) (*Entity, error)
      - update       # Update{Entity}(ctx, entity) error
      - delete       # Delete{Entity}(ctx, id) error
      - list         # List{Entity}s(ctx, params) ([]*Entity, error)

    # Custom operations
    operations:
      - getByEmail   # Get{Entity}ByEmail(ctx, email) (*Entity, error)
      - getByName    # Get{Entity}ByName(ctx, name) (*Entity, error)
      - findActive   # FindActive{Entity}s(ctx) ([]*Entity, error)

    # Database configuration
    database:
      table: table_name     # Database table name
      postgres: true        # Generate PostgreSQL implementation
      sqlite: true          # Generate SQLite implementation
```

### Service Configuration

```yaml
x-codegen:
  service:
    enabled: true # Generate service interface
    implementation: false # Generate service implementation (optional)
```

### Events Configuration

```yaml
x-codegen:
  events:
    - created # Publish{Entity}Created(ctx, entity) error
    - updated # Publish{Entity}Updated(ctx, entity) error
    - deleted # Publish{Entity}Deleted(ctx, entity) error
```

### Field Configuration

```yaml
properties:
  fieldName:
    type: string
    x-codegen:
      database:
        column: field_name # Database column name
        primaryKey: true # Mark as primary key
        unique: true # Add unique constraint
        nullable: false # Allow NULL values
        default: "value" # Default value
        index: true # Create index
```

## Configuration File

Create a `codegen.yaml` file in your project root:

```yaml
# OpenAPI specification path
openapi: api/openapi.bundled.yaml

# Output directory for generated files
output: internal

# Generator configurations
generators:
  repository:
    interface: repository.gen.go
    postgres: postgres.gen.go
    sqlite: sqlite.gen.go

  service:
    interface: service.gen.go
    implementation: server.gen.go

  events:
    interface: events.gen.go
    redis: events_redis.gen.go
    nats: events_nats.gen.go

  defaults: internal/infrastructure/config/defaults.gen.go

# Domain configurations (optional)
domains:
  auth:
    path: internal/auth
    schemas:
      - User
      - Session

  organizations:
    path: internal/organizations
    schemas:
      - Organization
      - Member
```

## Best Practices

### 1. Define First, Generate Second

Always define your API contract in OpenAPI before implementing:

```yaml
# 1. Define in api/components/schemas/user.yaml
User:
  x-codegen:
    repository:
      operations: [create, get, update, delete, list, getByEmail]
  properties:
    id:
      type: string
      format: uuid
    email:
      type: string
      format: email
```

```bash
# 2. Generate code
make generate-codegen

# 3. Implement business logic in service.go
```

### 2. Use Standard Operations

Prefer standard CRUD operations over custom ones:

```yaml
# Good - uses standard operations
x-codegen:
  repository:
    operations: [create, get, update, delete, list]

# Avoid custom operations unless necessary
x-codegen:
  repository:
    operations: [customFind, specialGet, uniqueUpdate]
```

### 3. Keep Schemas Focused

Each schema should represent a single domain entity:

```yaml
# Good - focused schema
User:
  x-codegen:
    repository:
      operations: [create, get, update]
  properties:
    id: ...
    email: ...
    name: ...

# Avoid - mixed concerns
UserWithOrganizationAndPermissions:
  properties:
    userId: ...
    organizationId: ...
    permissions: ...
```

### 4. Document Generation

Add descriptions to help future developers:

```yaml
Organization:
  description: Represents a multi-tenant organization
  x-codegen:
    repository:
      operations: [create, get, update, delete, list]
      # Generate both PostgreSQL and SQLite for testing
      database:
        postgres: true
        sqlite: true
```

## Custom Templates

The generator uses Go templates located in `internal/codegen/templates/`:

- `repository.go.tmpl` - Repository interface template
- `repository_postgres.go.tmpl` - PostgreSQL implementation
- `repository_sqlite.go.tmpl` - SQLite implementation
- `service.go.tmpl` - Service interface template
- `handler.go.tmpl` - HTTP handler template
- `events.go.tmpl` - Event publisher interface

## Troubleshooting

### Generated Code Not Updating

```bash
# Clean and regenerate
make clean-generated
make generate
```

### Import Errors

Ensure your OpenAPI spec is bundled:

```bash
make api-bundle
make generate-codegen
```

### Missing Operations

Check x-codegen annotations are properly formatted:

```yaml
x-codegen:
  repository:
    operations: # Must be an array
      - create
      - get
      - update
```

### Type Mismatches

Ensure OpenAPI types match Go types:

```yaml
# OpenAPI format -> Go type
format: uuid      -> uuid.UUID
format: date-time -> time.Time
format: email     -> openapi_types.Email
type: integer     -> int or int64
type: number      -> float32 or float64
```

## Advanced Usage

### Multi-Database Support

Generate implementations for multiple databases:

```yaml
x-codegen:
  repository:
    database:
      postgres: true
      sqlite: true
      mysql: false # Coming soon
```

### Custom Repository Methods

Add custom repository operations:

```yaml
x-codegen:
  repository:
    operations:
      - create
      - get
      - getByEmail # Generates GetUserByEmail
      - findActive # Generates FindActiveUsers
      - countByStatus # Generates CountUsersByStatus
```

### Conditional Generation

Skip generation for existing implementations:

```bash
# Generator checks for existing service.go
# If found, skips service generation
# Useful for custom business logic
```

## Examples

### Complete User Domain

```yaml
User:
  description: User account entity
  x-codegen:
    repository:
      operations:
        - create
        - get
        - update
        - delete
        - list
        - getByEmail
      database:
        table: users
        postgres: true
        sqlite: true
    service:
      enabled: true
    events:
      - created
      - updated
      - deleted
  properties:
    id:
      type: string
      format: uuid
      x-codegen:
        database:
          primaryKey: true
    email:
      type: string
      format: email
      x-codegen:
        database:
          unique: true
          index: true
    name:
      type: string
      maxLength: 255
    passwordHash:
      type: string
      x-codegen:
        database:
          column: password_hash
    emailVerified:
      type: boolean
      default: false
      x-codegen:
        database:
          column: email_verified
    createdAt:
      type: string
      format: date-time
      x-codegen:
        database:
          column: created_at
          default: CURRENT_TIMESTAMP
    updatedAt:
      type: string
      format: date-time
      x-codegen:
        database:
          column: updated_at
          default: CURRENT_TIMESTAMP
```

This generates a complete user domain with:

- Repository interface with CRUD + getByEmail
- PostgreSQL and SQLite implementations
- Service interface for business logic
- HTTP server implementation
- Event publishers for user lifecycle
- Test mocks for all interfaces

## Integration with Other Tools

### SQLC

The code generator works alongside SQLC:

```sql
-- internal/database/queries/users.sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (id, email, name) VALUES ($1, $2, $3);
```

### Mockery

Generate mocks for testing:

```bash
make generate-mocks
```

Uses interfaces generated by the code generator to create test mocks.

### API Client Generation

TypeScript/JavaScript clients are generated from the same OpenAPI spec:

```bash
make generate-js-client
```

## Related Documentation

- [Development Guide](development.md) - Overall development workflow
- [Testing Guide](testing.md) - Testing with generated mocks
- [Makefile Commands](makefile-commands.md) - All generation commands
- [Project Layout](../architecture/project-layout.md) - Where generated files go
