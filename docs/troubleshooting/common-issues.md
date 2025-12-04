# Troubleshooting

Common issues and solutions.

## Installation

### `command not found: archesai`

```bash
# Install with go
go install github.com/archesai/archesai/cmd/archesai@latest

# Or build from source
make build
export PATH=$PATH:$(pwd)/bin
```

### `command not found: go`

Install Go 1.22+ from [go.dev/doc/install](https://go.dev/doc/install)

## Code Generation

### Generation fails with errors

```bash
# Clean and regenerate
make clean-generated
make generate

# Validate your OpenAPI spec
archesai generate --spec api.yaml --bundle --output test.yaml
```

### Type errors after generation

Ensure generated files are up to date:

```bash
make clean-generated
make generate
```

### Mocks not generating

```bash
make generate-mocks
```

## Development Server

### Port already in use

```bash
# Find what's using the port
lsof -i :3000
lsof -i :3001

# Kill the process
kill -9 <PID>

# Or use different ports
export ARCHES_SERVER_PORT=8001
export ARCHES_PLATFORM_PORT=8000
archesai dev
```

### Hot reload not working

```bash
# Install air
go install github.com/air-verse/air@latest

# Use make dev-api instead of make run-api
make dev-api
```

## Database

### Connection refused

```bash
# Start PostgreSQL with Docker
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=archesdb \
  -p 5432:5432 \
  postgres:15

# Check if running
docker ps | grep postgres
docker start postgres  # If stopped
```

### Redis connection refused

Redis is optional. If not needed:

```yaml
# .archesai.yaml
redis:
  enabled: false
```

Or start Redis:

```bash
docker run -d --name redis -p 6379:6379 redis:7
```

## OpenAPI

### Invalid schema errors

```bash
# Validate your spec
npx @apidevtools/swagger-cli validate your-spec.yaml
```

Common issues:

- Missing required fields
- Invalid `$ref` paths
- Duplicate operation IDs

### Multi-file spec not working

Bundle first:

```bash
archesai generate --spec api.yaml --bundle --output bundled.yaml
archesai generate --spec bundled.yaml --output ./myapp
```

## Configuration

### Config not loading

```bash
# Check current config
archesai config show

# Validate
archesai config validate

# View environment variables
archesai config env
```

### Environment variables not working

Format: `ARCHES_<SECTION>_<KEY>`

```bash
export ARCHES_SERVER_PORT=8080
export ARCHES_DATABASE_HOST=localhost
```

## Tests

### Tests failing

```bash
# Verbose output
make test-verbose

# Skip integration tests
make test-short

# Clear cache
make clean-test
make test
```

## Common Error Messages

| Error                       | Solution                           |
| --------------------------- | ---------------------------------- |
| `missing required field`    | Check OpenAPI spec required fields |
| `no such file or directory` | Run `make generate`                |
| `permission denied`         | `chmod +x bin/archesai`            |

## Debug Mode

```bash
export ARCHES_LOG_LEVEL=debug
archesai dev --verbose
```

## Getting Help

- [GitHub Issues](https://github.com/archesai/archesai/issues)
- Email: <support@archesai.com>
