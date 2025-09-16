# Development Setup

This section contains detailed guides for setting up your development environment and working with
the Arches codebase.

## Getting Started

- **[Development Guide](development.md)** - Main development setup and workflow
- **[Code Generation Guide](code-generation.md)** - Unified code generation system
- **[Testing Strategy](testing.md)** - Comprehensive testing documentation
- **[Contributing Guidelines](../contributing.md)** - How to contribute to the project
- **[Makefile Commands](makefile-commands.md)** - Complete command reference
- **[TUI Guide](../features/tui.md)** - Terminal user interface documentation

## Reports and Coverage

- **[Test Coverage Report](test-coverage-report.md)** - Current test coverage status

## Development Topics

### Code Generation

Learn about our unified code generation system:

- **[Code Generation Guide](code-generation.md)** - Complete documentation
- OpenAPI-driven development with x-codegen annotations
- Automatic repository and service generation
- Multi-database support (PostgreSQL/SQLite)
- Test mock generation with Mockery

### Testing

- Unit testing with mocks
- Integration testing with testcontainers
- End-to-end testing strategies
- Coverage reporting

### Debugging

- Using the TUI for debugging
- Log analysis
- Performance profiling
- Common debugging techniques

## Tools and Setup

### Required Tools

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis
- Docker

### Optional Tools

- k3d (local Kubernetes)
- Air (hot reload)
- Make (build automation)

## Quick Commands

```bash
# Start development
make dev

# Generate all code
make generate

# Run tests
make test

# Lint and format
make lint format
```
