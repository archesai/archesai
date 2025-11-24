# Troubleshooting Guide

This guide helps you diagnose and resolve common issues when using Arches or developing applications with it.

## Installation Issues

### Go Installation Problems

**Problem**: `command not found: go` or wrong Go version

**Solution**:

```bash
# Check Go version
go version

# Should be 1.22 or higher
# If not installed or wrong version:
# Visit https://go.dev/doc/install
```

### Node.js/pnpm Issues

**Problem**: `command not found: pnpm`

**Solution**:

```bash
# Install pnpm
npm install -g pnpm

# Or with corepack (Node.js 16.13+)
corepack enable
corepack prepare pnpm@latest --activate
```

### Binary Not Found

**Problem**: `command not found: archesai`

**Solution**:

```bash
# Option 1: Build from source
make build
export PATH=$PATH:$(pwd)/bin

# Option 2: Install with go
go install github.com/archesai/archesai/cmd/archesai@latest

# Option 3: Add to PATH permanently
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc
```

## Build & Generation Issues

### Code Generation Fails

**Problem**: `make generate` fails with errors

**Solution**:

```bash
# Clean and regenerate
make clean-generated
make deps
make generate

# If OpenAPI-related:
# Check your OpenAPI spec is valid
archesai generate --spec api/openapi.yaml --bundle --output test.yaml
# This will show validation errors
```

### Type Errors After Generation

**Problem**: Generated code has type mismatches

**Solution**:

```bash
# Ensure all generated files are up to date
make clean-generated
make generate
make lint

# Check for x-codegen annotations in OpenAPI spec
# They may be causing conflicts
```

### Mock Generation Issues

**Problem**: Mocks not generating or outdated

**Solution**:

```bash
# Regenerate mocks specifically
make generate-mocks

# Or clean and regenerate all
make clean
make generate
```

## Development Server Issues

### Port Already in Use

**Problem**: `bind: address already in use`

**Solution**:

```bash
# Find what's using the port
lsof -i :3000  # For platform UI
lsof -i :3001  # For API server

# Kill the process
kill -9 <PID>

# Or use different ports
export ARCHES_SERVER_PORT=8001
export ARCHES_PLATFORM_PORT=8000
archesai dev
```

### Hot Reload Not Working

**Problem**: Changes not reflected when saving files

**Solution**:

```bash
# Ensure air is installed
go install github.com/air-verse/air@latest

# Check .air.toml configuration
cat .air.toml

# Run with explicit hot reload
make dev-api  # Instead of make run-api
```

### Platform UI Won't Start

**Problem**: Frontend fails to compile or start

**Solution**:

```bash
# Clean and reinstall dependencies
cd web/platform
rm -rf node_modules pnpm-lock.yaml
pnpm install
pnpm dev

# Or from root
make clean-ts-deps
make deps-ts
make dev-platform
```

## Database Issues

### Connection Refused

**Problem**: `dial tcp 127.0.0.1:5432: connect: connection refused`

**Solution**:

```bash
# Start PostgreSQL with Docker
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=archesdb \
  -p 5432:5432 \
  postgres:15

# Or check if it's running
docker ps | grep postgres
docker start postgres  # If stopped
```

### Migration Errors

**Problem**: Database migrations fail

**Solution**:

```bash
# For Arches platform development:
# Migrations are handled by the generated code
# Check the generated migration files

# For generated apps:
# Look in the generated app's database/ directory
# Run migrations manually if needed:
psql -U postgres -d yourdb < migrations/*.sql
```

### Redis Connection Issues

**Problem**: `dial tcp 127.0.0.1:6379: connect: connection refused`

**Solution**:

```bash
# Redis is optional, but if you need it:
docker run -d --name redis -p 6379:6379 redis:7

# Or disable Redis in config:
# Edit .archesai.yaml
redis:
  enabled: false
```

## OpenAPI Generation Issues

### Invalid Schema Errors

**Problem**: Generation fails with schema validation errors

**Solution**:

```bash
# Validate your OpenAPI spec
npx @apidevtools/swagger-cli validate your-spec.yaml

# Common issues:
# - Missing required fields
# - Invalid $ref paths
# - Duplicate operation IDs
# - Missing component definitions
```

### Generated Code Missing Features

**Problem**: Expected code not generated

**Solution**:

```yaml
# Check x-codegen annotations in your spec:
paths:
  /users:
    get:
      x-codegen-custom-handler: true # For custom logic
      x-public-endpoint: true # Skip auth
```

### Bundle Mode Issues

**Problem**: Multi-file OpenAPI spec not working

**Solution**:

```bash
# Use bundle mode to combine files
archesai generate --spec api/openapi.yaml \
  --bundle \
  --output bundled.yaml

# Then generate from bundled file
archesai generate --spec bundled.yaml \
  --output ./myapp
```

## Testing Issues

### Tests Failing

**Problem**: `make test` fails

**Solution**:

```bash
# Run tests with verbose output
make test-verbose

# Run only unit tests (skip integration)
make test-short

# Clear test cache
make clean-test
make test
```

### Coverage Reports

**Problem**: Can't generate coverage reports

**Solution**:

```bash
# Generate coverage
make test-coverage

# View HTML report
make test-coverage-html
open coverage.html
```

## Docker & Kubernetes Issues

### Docker Build Fails

**Problem**: Docker image won't build

**Solution**:

```bash
# Build with correct Dockerfile
docker build -f deployments/docker/Dockerfile .

# Or use make
make docker-run

# Clean Docker cache
docker system prune -a
```

### Kubernetes Deployment Issues

**Problem**: Pods not starting

**Solution**:

```bash
# Check pod status
kubectl get pods
kubectl describe pod <pod-name>

# Check logs
kubectl logs <pod-name>

# Common issues:
# - ImagePullBackOff: Check image name/registry
# - CrashLoopBackOff: Check logs for errors
# - Pending: Check resource requirements
```

## Configuration Issues

### Config Not Loading

**Problem**: Configuration values not applied

**Solution**:

```bash
# Check config file location
archesai config show

# Validate configuration
archesai config validate

# Use environment variables
export ARCHES_DATABASE_HOST=localhost
export ARCHES_DATABASE_PORT=5432
```

### Environment Variables

**Problem**: Environment variables not working

**Solution**:

```bash
# View all Arches environment variables
archesai config env

# Format: ARCHES_<SECTION>_<KEY>
# Example:
export ARCHES_SERVER_PORT=8080
export ARCHES_DATABASE_HOST=localhost
```

## Performance Issues

### Slow Generation

**Problem**: Code generation takes too long

**Solution**:

```bash
# Use parallel generation (default)
make generate

# For large specs, split into smaller files
# Use x-codegen annotations selectively
```

### High Memory Usage

**Problem**: Process using too much memory

**Solution**:

```bash
# Limit concurrent operations
export GOMAXPROCS=2

# For Docker, set memory limits:
docker run -m 512m archesai/api
```

## Common Error Messages

### "missing required field"

**Cause**: OpenAPI spec validation error
**Fix**: Check required fields in your OpenAPI spec

### "no such file or directory"

**Cause**: Generated files missing
**Fix**: Run `make generate`

### "syntax error near unexpected token"

**Cause**: Shell script issues
**Fix**: Check shell compatibility, use bash not sh

### "permission denied"

**Cause**: File permission issues
**Fix**:

```bash
chmod +x bin/archesai
# Or run with proper permissions
sudo make install
```

## Getting Help

If you're still stuck:

1. **Check logs carefully** - Error messages usually indicate the problem
2. **Search existing issues**: [GitHub Issues](https://github.com/archesai/archesai/issues)
3. **Ask for help**:
   - Create a [new issue](https://github.com/archesai/archesai/issues/new)
   - Include: Error messages, Steps to reproduce, Environment details
4. **Email support**: <support@archesai.com>

## Debug Mode

Enable debug logging for more information:

```bash
# Set debug mode
export ARCHES_LOG_LEVEL=debug

# Run with verbose output
archesai dev --verbose

# Check all logs
journalctl -f | grep archesai
```

## See Also

- [Development Guide](../guides/development.md) - Setup and development
- [CLI Reference](../cli-reference.md) - Complete command reference
- [Configuration Guide](../guides/configuration.md) - Config options
- [Deployment Guide](../deployment/overview.md) - Production issues
