# ArchesAI Assistant Guide

## Essential Commands

```bash
make generate         # Run after API/SQL changes
make lint            # Check code quality
make dev             # Start backend server
pnpm dev:platform    # Start frontend
```

## Project Conventions

- **Flat package structure** in domains (no subdirectories)
- **Generate first, code second** - Define in OpenAPI/SQL before implementing
- **No cross-domain imports** - Keep domains isolated
- **Use generated types** - Don't create manual type definitions

## Code Generation

After modifying:

- `api/openapi.yaml` → Run `make generate-oapi`
- `internal/database/queries/*.sql` → Run `make generate-sqlc`
- Any x-codegen annotations → Run `make generate-codegen`

## Domain Structure

Each domain follows this pattern:

```
internal/{domain}/
├── {domain}.go            # Package with go:generate directives
├── service.go              # Business logic (manual)
├── handler_http.go         # HTTP handlers (manual)
├── repository_postgres.go  # Database impl (manual)
└── *.gen.go               # Generated files (don't edit)
```

## Testing

```bash
make test              # Run all tests
go test ./internal/auth/...  # Test specific domain
```

## Database

```bash
make migrate-up        # Apply migrations
make migrate-create name=feature  # New migration
```

## Important Files

- `internal/app/app.go` - Dependency wiring
- `internal/config/config.go` - Configuration
- `codegen.yaml` - Code generation config

## Quick Fixes

**Build fails**: `make generate && make lint`
**Type errors**: Check generated files are up to date
**Import errors**: No cross-domain imports allowed

## Environment

All config uses `ARCHESAI_` prefix:

```bash
ARCHESAI_DATABASE_URL=postgres://...
ARCHESAI_JWT_SECRET=secret-key
ARCHESAI_SERVER_PORT=8080
```
