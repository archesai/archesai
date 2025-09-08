# ArchesAI Assistant Guide

## Essential Commands

```bash
make generate         # Run after API/SQL changes
make lint            # Check code quality
make dev             # Start backend server
pnpm dev:platform    # Start frontend
make format          # Format code
```

## Project Conventions

- **Flat package structure** in domains (no subdirectories)
- **Generate first, code second** - Define in OpenAPI/SQL before implementing
- **No cross-domain imports** - Keep domains isolated
- **Use generated types** - Don't create manual type definitions

## Code Generation

After modifying:

- `api/openapi.yaml` → Run `make generate-oapi` (generates types.gen.go, http.gen.go)
- `internal/database/queries/*.sql` → Run `make generate-sqlc` (generates database code)
- Any x-codegen annotations → Run `make generate-codegen` (generates repository.gen.go)

### Currently Generated:

- ✅ OpenAPI types for all domains
- ✅ HTTP interfaces (ServerInterface) for all domains
- ✅ Repository interfaces for: Auth, Organizations, Workflows, Content
- ⚠️ Tool entity has x-codegen but not generating yet

## Domain Structure

Each domain follows this pattern:

```
internal/{domain}/
├── {domain}.go            # Package constants, errors, and go:generate directives
├── service.go             # Business logic (manual)
├── handler.go             # HTTP handler implementation (manual)
├── postgres.go            # PostgreSQL repository (manual)
├── sqlite.go              # SQLite repository (manual)
├── types.gen.go           # Generated OpenAPI types
├── http.gen.go            # Generated HTTP interface (ServerInterface)
└── repository.gen.go      # Generated repository interface
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
**Repository not generating**: Check inferDomain() in internal/codegen/parser.go
**Method signature mismatch**: Use GetXByID not GetX, Update(ctx, id, entity) not Update(ctx, entity)
**Directory moving**:Try not to have to cd into other directories all the time. You can pretty much do everything from the makefile, which is the preferable way of doing anything.
