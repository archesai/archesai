# ArchesAI Documentation

## Quick Links

### Core Documentation

- [Terminal UI (TUI)](TUI.md) - Configuration viewer and AI chat interface
- [API Reference](../api/openapi.yaml) - OpenAPI specification
- [Contributing](../CONTRIBUTING.md) - How to contribute

### Examples

- [TUI Demo](examples/tui_demo.go) - Complete TUI usage examples with AI agents

### Packages

#### AI & Agents

- [LLM Package](../internal/llm/) - Multi-provider LLM interface
- [Swarm Package](../internal/swarm/) - Multi-agent orchestration system

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
```

## Architecture

ArchesAI follows a domain-driven design with:

- Flat package structure in domains
- Generated code from OpenAPI specs
- Repository pattern for data access
- Clean separation of concerns

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
- **Multi-Agent System**: SwarmGo integration for agent orchestration
- **Beautiful TUI**: Terminal interface for configuration and chat
- **Code Generation**: OpenAPI and SQL-driven development
- **Modern Stack**: Go, PostgreSQL/SQLite, Redis

## License

See [LICENSE](../LICENSE) file for details.
