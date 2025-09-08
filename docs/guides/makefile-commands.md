# Makefile Commands Reference

This document provides a comprehensive reference for all available Makefile targets in the ArchesAI project.

## Usage

```bash
make [target]
```

To see this help in the terminal:

```bash
make help
```

## Available Targets

### Core Development Commands

| Command         | Description                              |
| --------------- | ---------------------------------------- |
| `make all`      | Default: generate, lint, and format code |
| `make dev`      | Run all services in development mode     |
| `make generate` | Generate all code                        |
| `make lint`     | Run all linters                          |
| `make format`   | Format all code                          |
| `make test`     | Run all tests                            |

### Build Commands

| Command                | Description                  |
| ---------------------- | ---------------------------- |
| `make build`           | Build all binaries           |
| `make build-archesai`  | Build archesai server binary |
| `make build-web`       | Build web assets             |
| `make clean`           | Clean build artifacts        |
| `make clean-generated` | Clean all generated code     |

### Code Generation

| Command                       | Description                              |
| ----------------------------- | ---------------------------------------- |
| `make generate`               | Generate all code                        |
| `make generate-oapi`          | Generate OpenAPI server code             |
| `make generate-sqlc`          | Generate database code with sqlc         |
| `make generate-codegen`       | Generate codegen                         |
| `make generate-codegen-types` | Generate types for codegen configuration |
| `make generate-mocks`         | Generate test mocks using mockery        |

### Dependency Management

| Command                 | Description                       |
| ----------------------- | --------------------------------- |
| `make deps`             | Install all dependencies          |
| `make deps-go`          | Install Go dependencies and tools |
| `make deps-node`        | Install Node.js dependencies      |
| `make deps-update`      | Update all dependencies           |
| `make deps-update-go`   | Update Go dependencies            |
| `make deps-update-node` | Update Node.js dependencies       |

### Code Quality

| Command               | Description                             |
| --------------------- | --------------------------------------- |
| `make lint`           | Run all linters                         |
| `make lint-go`        | Run Go linter                           |
| `make lint-node`      | Run Node.js linter (includes typecheck) |
| `make lint-openapi`   | Lint OpenAPI specification              |
| `make format`         | Format all code                         |
| `make format-go`      | Format Go code                          |
| `make format-node`    | Format Node.js/TypeScript code          |
| `make typecheck-node` | Run TypeScript type checking            |

### Testing

| Command                        | Description                                   |
| ------------------------------ | --------------------------------------------- |
| `make test`                    | Run all tests                                 |
| `make test-short`              | Run short tests only (skip integration tests) |
| `make test-verbose`            | Run all tests with verbose output             |
| `make test-bench`              | Run benchmark tests                           |
| `make test-coverage`           | Generate test coverage report                 |
| `make test-coverage-html`      | Generate HTML coverage report                 |
| `make test-clean`              | Clean test cache and coverage files           |
| `make test-watch`              | Run tests in watch mode (requires fswatch)    |
| `make test-domain DOMAIN=auth` | Test specific domain                          |

### Database Operations

| Command                              | Description                         |
| ------------------------------------ | ----------------------------------- |
| `make migrate`                       | Alias for migrate-up                |
| `make migrate-up`                    | Apply database migrations           |
| `make migrate-down`                  | Rollback database migrations        |
| `make migrate-reset`                 | Reset database to initial state     |
| `make migrate-status`                | Show migration status               |
| `make migrate-create name=add_users` | Create new migration                |
| `make convert-schema`                | Convert PostgreSQL schema to SQLite |

### OpenAPI Operations

| Command               | Description                           |
| --------------------- | ------------------------------------- |
| `make openapi-bundle` | Bundle OpenAPI into single file       |
| `make openapi-split`  | Split OpenAPI into multiple files     |
| `make openapi-stats`  | Show OpenAPI specification statistics |

### Runtime Commands

| Command           | Description                        |
| ----------------- | ---------------------------------- |
| `make run`        | Alias for run-api                  |
| `make run-api`    | Run the API server                 |
| `make run-web`    | Run the web UI server              |
| `make run-worker` | Run the background worker          |
| `make tui`        | Launch the TUI interface           |
| `make watch`      | Run with hot reload (requires air) |

### Docker Operations

| Command            | Description                       |
| ------------------ | --------------------------------- |
| `make docker-run`  | Build and run with Docker Compose |
| `make docker-stop` | Stop Docker Compose services      |

### Kubernetes/Skaffold

| Command                | Description                   |
| ---------------------- | ----------------------------- |
| `make skaffold-dev`    | Run with Skaffold in dev mode |
| `make skaffold-run`    | Deploy with Skaffold          |
| `make skaffold-delete` | Delete Skaffold deployment    |

### Development Tools

| Command                    | Description                     |
| -------------------------- | ------------------------------- |
| `make install-tools`       | Install development tools       |
| `make install-completions` | Install shell completions guide |
| `make help`                | Show this help message          |

## Common Workflows

### Starting Development

```bash
make deps           # Install dependencies
make generate       # Generate code
make dev           # Start all services
```

### After API Changes

```bash
make generate-oapi  # Regenerate from OpenAPI
make lint          # Check for issues
make test          # Verify tests pass
```

### After Database Changes

```bash
make generate-sqlc  # Regenerate database code
make migrate-up     # Apply migrations
make test          # Verify tests pass
```

### Before Committing

```bash
make all           # Generate, lint, and format
make test          # Run all tests
```

### Production Build

```bash
make clean         # Clean previous builds
make generate      # Ensure code is current
make build         # Build all binaries
```

## Tool Requirements

Some commands require additional tools to be installed:

- **fswatch**: Required for `make test-watch`
- **air**: Required for `make watch` (hot reload)
- **Docker**: Required for Docker-related commands
- **Skaffold**: Required for Kubernetes deployment commands
- **kubectl**: Required for Kubernetes operations

Install tools with:

```bash
make install-tools
```

## Environment Variables

Some commands may require environment variables to be set:

- Database connection strings for migration commands
- API keys for integration tests
- Docker registry credentials for deployment

See [Development Guide](development.md) for complete environment setup.
