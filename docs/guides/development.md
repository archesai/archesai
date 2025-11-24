# Development Guide

Welcome to the Arches development guide! This document covers everything you need to know
to set up your development environment and contribute to the project.

## Prerequisites

- Go 1.22+
- Node.js 20+ with pnpm
- PostgreSQL 15+ (for development)
- Redis 7+ (optional, for caching and sessions)
- Docker and Docker Compose (optional, for containerized development)

## Quick Setup

1. **Clone the repository**:

   ```bash
   git clone https://github.com/archesai/archesai.git
   cd archesai
   ```

2. **Install dependencies**:

   ```bash
   make deps        # Install all dependencies
   # Or separately:
   make deps-go     # Go dependencies
   make deps-ts     # Node.js dependencies
   ```

3. **Build the CLI**:

   ```bash
   make build       # Creates ./bin/archesai
   ```

4. **Start development environment**:

   ```bash
   # Option 1: Using archesai CLI
   archesai dev     # Start with hot reload
   archesai dev --tui # With interactive TUI

   # Option 2: Using Make
   make dev-all     # Start all services
   make dev-api     # Start API only
   make dev-platform # Start platform UI only
   ```

## Essential Commands

```bash
# Code Generation
make generate        # Generate all code from OpenAPI/SQL
archesai generate --spec api/openapi.yaml --output ./generated

# Development
make dev-all         # Start everything with hot reload
make dev-api         # Start API server with hot reload
make dev-platform    # Start platform UI with hot reload

# Code Quality
make lint            # Run all linters
make format          # Format all code
make test            # Run all tests

# Building
make build           # Build all binaries
make build-api       # Build API server
make build-platform  # Build platform UI
```

For a complete list of commands, see [Makefile Commands](makefile-commands.md).

## Project Structure

For a detailed overview of the project organization, see
[Project Layout](../architecture/project-layout.md).

## Development Workflow

### For Arches Platform Development

1. **Modify OpenAPI Specification**:
   - Edit files in `api/` directory
   - Add x-codegen annotations for custom behavior

2. **Generate Code**:

   ```bash
   make generate    # Regenerate all code
   ```

3. **Implement Custom Logic**:
   - Add business logic in handler files
   - Write tests using generated mocks

4. **Test Your Changes**:

   ```bash
   make test        # Run all tests
   make test-short  # Quick tests only
   ```

### For Using Arches to Build Apps

1. **Create Your OpenAPI Spec**:

   ```yaml
   # myapp.yaml
   openapi: 3.1.0
   info:
     title: My App API
     version: 1.0.0
   # ... your API definition
   ```

2. **Generate Your Application**:

   ```bash
   archesai generate --spec myapp.yaml --output ./myapp
   ```

3. **Start Development**:

   ```bash
   cd myapp
   archesai dev
   ```

4. **Customize Generated Code**:
   - Add business logic in `handlers/`
   - Modify templates if needed
   - Extend generated models

For detailed code generation instructions, see the [Code Generation Guide](code-generation.md).

## Contributing

See our [Contributing Guide](../contributing.md) for detailed information about:

- Code style and conventions
- Pull request process
- Issue reporting
- Development best practices

## Architecture

Learn about the system design and patterns in our
[Architecture Documentation](../architecture/system-design.md).

## Need Help?

- Check the [Troubleshooting Guide](../troubleshooting/common-issues.md)
- Open an issue on GitHub
