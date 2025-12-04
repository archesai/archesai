# Database

Arches generates database implementations for both PostgreSQL and SQLite, with automatic schema generation and migrations.

## Multi-Database Support

Arches generates repository implementations for two databases:

- **PostgreSQL** - Production database
- **SQLite** - Development and testing

Both implementations share the same interface, making it easy to switch between them:

```go
// Use PostgreSQL in production
repo := repositories.NewPostgresUserRepository(db)

// Use SQLite for testing
repo := repositories.NewSQLiteUserRepository(db)
```

## Automatic Schema Generation

Database schemas are automatically generated from your OpenAPI spec. Define entities with `x-codegen` annotations:

```yaml
components:
  schemas:
    User:
      type: object
      x-codegen:
        repository:
          database:
            table: users
            postgres: true
            sqlite: true
      properties:
        id:
          type: string
          format: uuid
          x-codegen:
            database:
              primaryKey: true
        email:
          type: string
          x-codegen:
            database:
              unique: true
              index: true
        name:
          type: string
          maxLength: 255
        createdAt:
          type: string
          format: date-time
          x-codegen:
            database:
              column: created_at
              default: CURRENT_TIMESTAMP
```

Arches generates:

- Table definitions with proper types
- Primary keys and foreign keys
- Unique constraints
- Indexes
- Default values

## Automatic Migrations

SQL migrations are generated automatically when you define or change your OpenAPI spec:

```bash
# Generate/regenerate after spec changes
archesai generate --spec api.yaml --output ./myapp

# Migrations appear in database/migrations/
```

Generated migrations include:

- `CREATE TABLE` statements
- `ALTER TABLE` for schema changes
- Index creation
- Constraint definitions

## Field Configuration

Configure database columns with `x-codegen.database`:

```yaml
properties:
  email:
    type: string
    x-codegen:
      database:
        column: email_address # Custom column name
        unique: true # UNIQUE constraint
        index: true # Create index
        nullable: false # NOT NULL (default)

  status:
    type: string
    x-codegen:
      database:
        default: "'active'" # Default value

  createdAt:
    type: string
    format: date-time
    x-codegen:
      database:
        column: created_at
        default: CURRENT_TIMESTAMP
```

## Type Mappings

| OpenAPI Type                   | PostgreSQL         | SQLite    |
| ------------------------------ | ------------------ | --------- |
| `string`                       | `VARCHAR(255)`     | `TEXT`    |
| `string` + `format: uuid`      | `UUID`             | `TEXT`    |
| `string` + `format: date-time` | `TIMESTAMPTZ`      | `TEXT`    |
| `integer`                      | `INTEGER`          | `INTEGER` |
| `number`                       | `DOUBLE PRECISION` | `REAL`    |
| `boolean`                      | `BOOLEAN`          | `INTEGER` |

## Configuration

Configure database connection in `.archesai.yaml`:

```yaml
database:
  driver: postgres # or sqlite
  host: localhost
  port: 5432
  name: mydb
  user: postgres
  password: postgres
```

Or via environment variables:

```bash
export ARCHES_DATABASE_DRIVER=postgres
export ARCHES_DATABASE_HOST=localhost
export ARCHES_DATABASE_PORT=5432
export ARCHES_DATABASE_NAME=mydb
export ARCHES_DATABASE_USER=postgres
export ARCHES_DATABASE_PASSWORD=postgres
```

## Generated Repository Interface

For each entity, Arches generates a repository interface:

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Get(ctx context.Context, id uuid.UUID) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, params *ListParams) ([]*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
}
```

With implementations in the infrastructure layer:

```text
infrastructure/
├── postgres/
│   ├── migrations/           # SQL migration files
│   ├── queries/              # SQLC query files
│   └── repositories/         # Generated repository implementations
└── sqlite/
    ├── migrations/
    ├── queries/
    └── repositories/
```

## Running Migrations

```bash
# Migrations run automatically on startup
# Or run manually:
psql -U postgres -d mydb < infrastructure/postgres/migrations/*.sql
```
