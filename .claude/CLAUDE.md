# ArchesAI Assistant Guide

## Essential Commands

```bash
make generate     # Run after API/SQL changes
make lint         # Check code quality
make dev          # Start backend server
pnpm dev:platform # Start frontend
make format       # Format code
make help         # See all make commands
```

## Project Conventions

- **Generate first, code second** - Define in OpenAPI/SQL before implementing
- **Use generated types** - Don't create manual type definitions

## Code Generation

After modifying:

- `api/**/*.yaml` → Run `make generate-oapi` (generates types.gen.go, http.gen.go)
- `internal/database/**` or `internal/migrations/**` → Run `make generate-sqlc` (generates database
  code)
- Any x-codegen annotations → Run `make generate-codegen` (generates repository.gen.go)

## Testing

```bash
make test                   # Run all tests
go test ./internal/auth/... # Test specific domain
```

## Database

```bash
make db-migrate-up                  # Apply migrations
make db-migrate-create name=feature # New migration
```

## Quick Fixes

**Build fails**: `make generate && make lint` **Type errors**: Check generated files are up to date
**Directory moving**: Try not to have to cd into other directories all the time. You can pretty much
do everything from the makefile, which is the preferable way of doing anything.

## Docs - MAKE SURE TO ALWAYS UPDATE THESE FILES AFTER MAKING A CHANGE

@../docs/architecture/project-layout.md
@../docs/architecture/overview.md
@../docs/guides/makefile-commands.md
@../docs/guides/testing.md
@../README.md

## Task Master AI Instructions

**Import Task Master's development workflow commands and guidelines, treat as if import is in the
main CLAUDE.md file.**

@../.taskmaster/CLAUDE.md

TIPS:

- DO NOT SWITCH DIRECTORIES, STAY IN THE ROOT AT ALL TIMES
- Do not create your own mocks. Always try to use mockery and generate from an interface. We have done this many times in this project.

ALWAYS USER MOCKERY FOR GETTING MOCKS, NEVER CREATE MOCKED SERVICES OR REPOSITORIES OR ANYTHING MANUALLY.
RUN go tool mockery
WE ARE RUNNING MOCKERY v3
MOCKERY CONFIG IS .mockery.yaml

DO NOT KEEP DEPRECATED OR LEGACY CODE, ALWAYS MAKE SURE YOU JUST IMPLEMENT LATEST

i want you to improve test coverage as much as possible.
you should only ever use mocks from mockery that will be found in mocks_test.go.
if you need to get a mocked interface from another package, alias the interface
in your local package and add it to .mockery.yaml
