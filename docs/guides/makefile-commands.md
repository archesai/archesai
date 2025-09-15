# Makefile Commands

Run `make help` to see all available commands.

## Essential Commands

```bash
make all            # Generate, lint, format
make dev            # Start development
make build          # Build binaries
make test           # Run tests
make clean          # Clean artifacts
```

## Code Generation

```bash
make generate              # Run all generators
make generate-oapi         # Generate from OpenAPI
make generate-sqlc         # Generate database code
make generate-codegen      # Generate custom code
make generate-mocks        # Generate test mocks
make generate-js-client    # Generate TypeScript client
```

## Development

```bash
make run-server     # Run API server
make run-web        # Run frontend
make run-watch      # Hot reload (needs air)
make run-tui        # Terminal UI
make run-docs       # Documentation site
```

## Testing

```bash
make test                  # All tests
make test-short            # Skip integration tests
make test-coverage         # Coverage report
make test-coverage-html    # HTML coverage
make test-bench            # Benchmarks
make test-watch            # Watch mode (needs fswatch)
```

## Linting & Formatting

```bash
make lint           # Run all linters
make lint-go        # Go linter
make lint-ts      # Node/TypeScript
make lint-openapi   # OpenAPI spec
make lint-docs      # Markdown

make format         # Format all code
make format-go      # Format Go
make format-node    # Format JS/TS
```

## Database

```bash
make db-migrate-up         # Apply migrations
make db-migrate-down       # Rollback
make db-migrate-create name=feature  # New migration
make db-migrate-status     # Check status
make db-migrate-reset      # Reset database
```

## Dependencies

```bash
make deps           # Install all
make deps-go        # Go dependencies
make deps-node      # Node dependencies
make deps-update    # Update all
make install-tools  # Dev tools
```

## Build & Deploy

```bash
make build-server   # Build server binary
make build-web      # Build frontend
make build-docs     # Build documentation

make docker-run     # Docker Compose up
make docker-stop    # Docker Compose down
```

## API Tools

```bash
make api-bundle     # Bundle OpenAPI spec
make api-split      # Split OpenAPI spec
make api-stats      # Show API statistics
```

## Cleanup

```bash
make clean-generated  # Remove generated code
make clean-test       # Remove test cache
make clean-deps       # Remove dependencies
make clean-docs       # Remove docs build
```
