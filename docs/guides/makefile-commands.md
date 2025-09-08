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

## Command Naming Convention

All commands follow a consistent `<action>-<target>` pattern:

- `build-*` - Build various components
- `run-*` - Run services or tools
- `test-*` - Run different test suites
- `lint-*` - Run linters for different languages
- `format-*` - Format code
- `generate-*` - Generate code from specifications
- `clean-*` - Clean various artifacts
- `db-*` - Database operations
- `api-*` - OpenAPI/API operations
- `deps-*` - Dependency management

## Available Targets

### Primary Commands

| Command     | Description                              |
| ----------- | ---------------------------------------- |
| `make all`  | Default: generate, lint, and format code |
| `make help` | Show help message with all targets       |
| `make dev`  | Run all services in development mode     |

### Build Commands

| Command             | Description                     |
| ------------------- | ------------------------------- |
| `make build`        | Build all binaries              |
| `make build-server` | Build archesai server binary    |
| `make build-web`    | Build web assets                |
| `make build-docs`   | Build documentation static site |

### Run Commands

| Command           | Description                             |
| ----------------- | --------------------------------------- |
| `make run`        | Alias for run-server                    |
| `make run-server` | Run the API server                      |
| `make run-web`    | Run the web UI server                   |
| `make run-worker` | Run the background worker               |
| `make run-watch`  | Run with hot reload (requires air)      |
| `make run-tui`    | Launch the TUI interface                |
| `make run-docs`   | Serve documentation locally with Docker |

### Generate Commands

| Command                       | Description                              |
| ----------------------------- | ---------------------------------------- |
| `make generate`               | Generate all code                        |
| `make generate-sqlc`          | Generate database code with sqlc         |
| `make generate-schema-sqlite` | Convert PostgreSQL schema to SQLite      |
| `make generate-oapi`          | Generate OpenAPI server code             |
| `make generate-codegen`       | Generate codegen                         |
| `make generate-codegen-types` | Generate types for codegen configuration |
| `make generate-mocks`         | Generate test mocks using mockery        |

### Test Commands

| Command                   | Description                                   |
| ------------------------- | --------------------------------------------- |
| `make test`               | Run all tests                                 |
| `make test-verbose`       | Run all tests with verbose output             |
| `make test-short`         | Run short tests only (skip integration tests) |
| `make test-coverage`      | Generate test coverage report                 |
| `make test-coverage-html` | Generate HTML coverage report                 |
| `make test-bench`         | Run benchmark tests                           |
| `make test-watch`         | Run tests in watch mode (requires fswatch)    |

### Lint Commands

| Command               | Description                             |
| --------------------- | --------------------------------------- |
| `make lint`           | Run all linters                         |
| `make lint-go`        | Run Go linter                           |
| `make lint-node`      | Run Node.js linter (includes typecheck) |
| `make lint-openapi`   | Lint OpenAPI specification              |
| `make lint-typecheck` | Run TypeScript type checking            |
| `make lint-docs`      | Lint documentation with markdownlint    |

### Format Commands

| Command            | Description                    |
| ------------------ | ------------------------------ |
| `make format`      | Format all code                |
| `make format-go`   | Format Go code                 |
| `make format-node` | Format Node.js/TypeScript code |

### Clean Commands

| Command                | Description                         |
| ---------------------- | ----------------------------------- |
| `make clean`           | Clean build artifacts               |
| `make clean-generated` | Clean all generated code            |
| `make clean-test`      | Clean test cache and coverage files |
| `make clean-docs`      | Clean documentation build           |

### Database Commands

| Command                                 | Description                     |
| --------------------------------------- | ------------------------------- |
| `make db-migrate`                       | Alias for db-migrate-up         |
| `make db-migrate-up`                    | Apply database migrations       |
| `make db-migrate-down`                  | Rollback database migrations    |
| `make db-migrate-create name=add_users` | Create new migration            |
| `make db-migrate-status`                | Show migration status           |
| `make db-migrate-reset`                 | Reset database to initial state |

### API/OpenAPI Commands

| Command           | Description                           |
| ----------------- | ------------------------------------- |
| `make api-bundle` | Bundle OpenAPI into single file       |
| `make api-split`  | Split OpenAPI into multiple files     |
| `make api-stats`  | Show OpenAPI specification statistics |

### Dependency Commands

| Command                 | Description                       |
| ----------------------- | --------------------------------- |
| `make deps`             | Install all dependencies          |
| `make deps-go`          | Install Go dependencies and tools |
| `make deps-node`        | Install Node.js dependencies      |
| `make deps-update`      | Update all dependencies           |
| `make deps-update-go`   | Update Go dependencies            |
| `make deps-update-node` | Update Node.js dependencies       |

### Install Commands

| Command                    | Description                     |
| -------------------------- | ------------------------------- |
| `make install-tools`       | Install development tools       |
| `make install-completions` | Install shell completions guide |

### Docker Commands

| Command            | Description                       |
| ------------------ | --------------------------------- |
| `make docker-run`  | Build and run with Docker Compose |
| `make docker-stop` | Stop Docker Compose services      |

### Kubernetes Commands

| Command                  | Description             |
| ------------------------ | ----------------------- |
| `make k8s-cluster-start` | Start k3d cluster       |
| `make k8s-cluster-stop`  | Stop k3d cluster        |
| `make k8s-deploy`        | Deploy with Helm        |
| `make k8s-upgrade`       | Upgrade Helm deployment |

### Skaffold Commands

| Command                | Description                   |
| ---------------------- | ----------------------------- |
| `make skaffold-dev`    | Run with Skaffold in dev mode |
| `make skaffold-run`    | Deploy with Skaffold          |
| `make skaffold-delete` | Delete Skaffold deployment    |

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
make generate-sqlc     # Regenerate database code
make db-migrate-up     # Apply migrations
make test             # Verify tests pass
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

### Documentation Development

```bash
make lint-docs     # Check markdown files
make build-docs    # Build static site
make run-docs      # Serve documentation locally
```

## Tool Requirements

Some commands require additional tools to be installed:

- **fswatch**: Required for `make test-watch`
- **air**: Required for `make run-watch` (hot reload)
- **Docker**: Required for Docker-related commands and documentation
- **Skaffold**: Required for Kubernetes deployment commands
- **kubectl**: Required for Kubernetes operations
- **pnpm**: Required for Node.js-related commands

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
