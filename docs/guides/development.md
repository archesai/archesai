# Development Guide

Welcome to the Arches development guide! This document covers everything you need to know to set
up your development environment and contribute to the project.

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- PostgreSQL 15 or later (for development)
- Redis (for caching and sessions)
- Docker and Docker Compose (for containerized development)

## Quick Setup

1. **Clone the repository**:

   ```bash
   git clone https://github.com/archesai/archesai.git
   cd archesai
   ```

2. **Install dependencies**:

   ```bash
   make deps
   ```

3. **Generate code**:

   ```bash
   make generate
   ```

4. **Set up the database**:

   ```bash
   make db-migrate-up
   ```

5. **Start development servers**:

   ```bash
   make dev          # Start backend
   pnpm dev:platform # Start frontend (in another terminal)
   ```

## Essential Commands

```bash
make generate     # Run after API/SQL/x-codegen changes
make lint         # Check code quality
make dev          # Start backend server
pnpm dev:platform # Start frontend
make format       # Format code
make test         # Run tests
```

For a complete list of commands, see [Makefile Commands](makefile-commands.md).

## Project Structure

For a detailed overview of the project organization, see
[Project Layout](../architecture/project-layout.md).

## Development Workflow

1. **Define First, Generate Second**:
   - Add x-codegen annotations to OpenAPI schemas
   - Define database queries in SQL files
   - Run `make generate` to create all boilerplate code
2. **Implement Business Logic**: Write your custom logic in `service.go` files
3. **Test with Generated Mocks**: Use the generated `mocks_test.gen.go` files for testing
4. **Follow the testing strategy** outlined in [Testing Documentation](testing.md)
5. **Use the TUI for configuration** as described in the [TUI Guide](../features/tui.md)

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
