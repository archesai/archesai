# ArchesAI Documentation

## Quick Links

### Core Documentation

- [Terminal UI (TUI)](TUI.md) - Configuration viewer and AI chat interface
- [API Reference](../api/openapi.yaml) - OpenAPI specification
- [Contributing](../CONTRIBUTING.md) - How to contribute

### Examples

- [TUI Demo](examples/tui.go) - Complete TUI usage examples with AI agents

### Packages

#### AI & Chat

- [LLM Package](../internal/llm/) - Multi-provider LLM interface with chat clients

#### Core Domains

- [Auth](../internal/auth/) - Authentication & authorization
- [Organizations](../internal/organizations/) - Organization management
- [Workflows](../internal/workflows/) - Workflow automation
- [Content](../internal/content/) - Content management

#### Infrastructure

- [Database](../internal/database/) - Database layer
- [Config](../internal/config/) - Configuration management
- [CLI](../internal/cli/) - Command-line interface
- [TUI](../internal/tui/) - Terminal user interface

## Getting Started

### 1. Configuration Viewer (No API Key Needed)

```bash
archesai tui
```

### 2. API Server

```bash
archesai api
```

### 3. AI Chat Interface

```bash
export OPENAI_API_KEY=your-key
archesai tui --chat

# Or with Ollama (local)
archesai tui --chat --provider=ollama
```

## Architecture

ArchesAI follows a domain-driven design with:

- Flat package structure in domains
- Generated code from OpenAPI specs
- Repository pattern for data access
- Clean separation of concerns
- Direct LLM client interfaces for simplicity

## Development

### Code Generation

```bash
make generate      # Generate all code
make generate-oapi # Generate OpenAPI types
make generate-sqlc # Generate database code
```

### Testing

```bash
make test          # Run all tests
make lint          # Run linters
```

### Building

```bash
make build         # Build binary
make dev           # Run in development mode
```

## Features

- **Multi-Provider AI**: Support for OpenAI, Claude, Gemini, Ollama
- **Chat Interface**: Simple persona-based chat system with session management
- **Beautiful TUI**: Terminal interface for configuration and chat
- **Code Generation**: OpenAPI and SQL-driven development
- **Modern Stack**: Go, PostgreSQL/SQLite, Redis
- **Clean Architecture**: Direct LLM usage without complex abstractions

## Chat Interface Features

- **Multiple Personas**: Switch between different AI personalities
- **Session Management**: Automatic conversation history tracking
- **Provider Support**: OpenAI (full), Ollama (local), others ready
- **Simple API**: Easy-to-use chat client interfaces

## License

See [LICENSE](../LICENSE) file for details.
