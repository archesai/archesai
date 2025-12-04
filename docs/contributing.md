# Contributing

Guide for contributing to the Arches codebase.

## Prerequisites

- Go 1.22+
- Node.js 20+ with pnpm
- PostgreSQL 15+
- Docker (optional)
- Make

## Setup

```bash
# Clone the repo
git clone https://github.com/archesai/archesai.git
cd archesai

# Install dependencies
make deps

# Setup database
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=archesdb \
  -p 5432:5432 \
  postgres:15

# Generate code
make generate

# Run tests
make test
```

## Development Commands

```bash
# Start dev server with hot reload
make dev-all

# Run linters
make lint

# Format code
make format

# Run tests
make test
make test-short    # Skip integration tests
make test-verbose  # Verbose output

# Generate code
make generate
make clean-generated  # Clean and regenerate
```

## Project Structure

```text
archesai/
├── cmd/archesai/      # CLI entry point
├── internal/          # Internal packages
│   ├── core/          # Domain models and events
│   ├── codegen/       # Code generation
│   └── ...
├── api/               # OpenAPI specifications
├── web/platform/      # Frontend (React + TypeScript)
└── docs/              # Documentation
```

## Workflow

### Making Changes

1. Create a branch:

   ```bash
   git checkout -b feature/your-feature
   ```

2. Make changes

3. Generate code if you modified OpenAPI or SQL:

   ```bash
   make generate
   ```

4. Run tests and lint:

   ```bash
   make test
   make lint
   ```

5. Commit with conventional commit format:

   ```bash
   git commit -m "feat: add new feature"
   ```

### Commit Format

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code refactoring
- `test:` - Tests
- `chore:` - Maintenance

Examples:

```text
feat(auth): add refresh token rotation
fix(codegen): resolve type mapping for arrays
docs: update getting started guide
```

### Pull Requests

1. Push your branch:

   ```bash
   git push origin feature/your-feature
   ```

2. Open a PR on GitHub

3. Fill out the PR template:
   - Description of changes
   - Related issues
   - Testing performed

## Code Style

### Go

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Handle errors explicitly
- Write table-driven tests

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "foo", "bar", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Something(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### TypeScript

- Use TypeScript for all frontend code
- Use functional components with hooks
- Follow existing patterns in the codebase

## Adding Features

1. Define the API in `api/openapi.yaml` with x-codegen annotations
2. Generate code: `make generate`
3. Implement business logic in handler files
4. Write tests
5. Update documentation if needed

## Getting Help

- Check existing issues: [GitHub Issues](https://github.com/archesai/archesai/issues)
- Open a new issue for bugs or feature requests
- Email: <support@archesai.com>
