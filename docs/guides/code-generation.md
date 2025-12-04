# Code Generation

Arches generates complete applications from OpenAPI specifications. Use `x-codegen` annotations to control what gets generated.

## Basic Usage

```bash
# Generate from OpenAPI spec
archesai generate --spec api.yaml --output ./myapp

# Preview without writing files
archesai generate --spec api.yaml --output ./myapp --dry-run

# Bundle multi-file spec into single file
archesai generate --spec api.yaml --bundle --output bundled.yaml
```

## Generated Files

From your OpenAPI spec, Arches generates:

```text
myapp/
├── main.gen.go                    # Entry point
├── spec/                          # OpenAPI spec (bundled)
├── models/                        # Go structs from schemas
├── controllers/                   # HTTP request handlers
├── application/                   # Application layer (use cases)
├── repositories/                  # Repository interfaces
├── bootstrap/                     # App initialization, routes, DI container
└── infrastructure/
    ├── postgres/
    │   ├── migrations/            # Auto-generated SQL migrations
    │   ├── queries/               # SQLC queries
    │   └── repositories/          # PostgreSQL implementations
    └── sqlite/
        ├── migrations/            # SQLite migrations
        ├── queries/               # SQLC queries
        └── repositories/          # SQLite implementations
```

## x-codegen Annotations

Add `x-codegen` to schemas to control generation:

### Repository Operations

```yaml
components:
  schemas:
    User:
      type: object
      x-codegen:
        repository:
          operations:
            - create # CreateUser(ctx, user)
            - get # GetUser(ctx, id)
            - update # UpdateUser(ctx, user)
            - delete # DeleteUser(ctx, id)
            - list # ListUsers(ctx, params)
            - getByEmail # GetUserByEmail(ctx, email)
          database:
            table: users
            postgres: true
            sqlite: true
      properties:
        id:
          type: string
          format: uuid
        email:
          type: string
          format: email
```

### Service Generation

```yaml
x-codegen:
  service:
    enabled: true # Generate service interface
```

### Event Publishing

```yaml
x-codegen:
  events:
    - created # PublishUserCreated(ctx, user)
    - updated # PublishUserUpdated(ctx, user)
    - deleted # PublishUserDeleted(ctx, user)
```

### Field Configuration

```yaml
properties:
  email:
    type: string
    x-codegen:
      database:
        column: email_address # Custom column name
        unique: true
        index: true
  createdAt:
    type: string
    format: date-time
    x-codegen:
      database:
        column: created_at
        default: CURRENT_TIMESTAMP
```

## Complete Example

```yaml
components:
  schemas:
    Organization:
      type: object
      x-codegen:
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
        name:
          type: string
          maxLength: 255
          x-codegen:
            database:
              unique: true
              index: true
        createdAt:
          type: string
          format: date-time
          x-codegen:
            database:
              column: created_at
              default: CURRENT_TIMESTAMP
```

This generates:

- Repository interface with CRUD + `GetOrganizationByName`
- PostgreSQL implementation
- Service interface
- Event publishers
- Test mocks

## Type Mappings

| OpenAPI                        | Go Type                |
| ------------------------------ | ---------------------- |
| `string`                       | `string`               |
| `string` + `format: uuid`      | `uuid.UUID`            |
| `string` + `format: date-time` | `time.Time`            |
| `string` + `format: email`     | `openapi_types.Email`  |
| `integer`                      | `int` or `int64`       |
| `number`                       | `float32` or `float64` |
| `boolean`                      | `bool`                 |
| `array`                        | `[]T`                  |
| `object`                       | struct                 |

## Makefile Commands

For Arches platform development:

```bash
make generate           # Generate all code
make generate-codegen   # Generate from OpenAPI/x-codegen
make generate-sqlc      # Generate SQLC database code
make generate-oapi      # Generate OpenAPI types
make generate-mocks     # Generate test mocks
make clean-generated    # Remove generated files
```

## Troubleshooting

### Generated code not updating

```bash
make clean-generated
make generate
```

### Import errors

Bundle your spec first:

```bash
archesai generate --spec api.yaml --bundle --output bundled.yaml
archesai generate --spec bundled.yaml --output ./myapp
```

### Missing operations

Check `x-codegen` format - operations must be an array:

```yaml
x-codegen:
  repository:
    operations: # Must be array
      - create
      - get
```

## Best Practices

1. **Define API first** - Write OpenAPI spec before implementing
2. **Use standard operations** - Prefer `create`, `get`, `update`, `delete`, `list`
3. **One entity per schema** - Keep schemas focused
4. **Use format hints** - `format: uuid`, `format: date-time`, etc.
5. **Document schemas** - Add descriptions for clarity
