# ArchesAI

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Proprietary-red?style=flat-square)](LICENSE)
[![API Documentation](https://img.shields.io/badge/API-OpenAPI%203.0-green?style=flat-square)](http://localhost:8080/docs)

A high-performance data processing platform with AI-powered chat interface and workflow automation.

## Quick Start

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

For detailed installation and setup instructions, see [Development Guide](docs/DEVELOPMENT.md).

## Documentation

### Core Documentation

- [Terminal UI (TUI)](docs/TUI.md) - Configuration viewer and AI chat interface
- [API Reference](api/openapi.yaml) - OpenAPI specification
- [Development Guide](docs/DEVELOPMENT.md) - Setup, build, and contribution guide
- [Architecture](docs/ARCHITECTURE.md) - System design and patterns
- [Contributing](docs/CONTRIBUTING.md) - How to contribute

### Packages

#### AI & Chat

- [LLM Package](internal/llm/) - Multi-provider LLM interface with chat clients

#### Core Domains

- [Auth](internal/auth/) - Authentication & authorization
- [Organizations](internal/organizations/) - Organization management
- [Workflows](internal/workflows/) - Workflow automation
- [Content](internal/content/) - Content management

#### Infrastructure

- [Database](internal/database/) - Database layer
- [Config](internal/config/) - Configuration management
- [CLI](internal/cli/) - Command-line interface
- [TUI](internal/tui/) - Terminal user interface

## Features

- **Multi-Provider AI**: Support for OpenAI, Claude, Gemini, Ollama
- **Chat Interface**: Simple persona-based chat system with session management
- **Beautiful TUI**: Terminal interface for configuration and chat
- **Workflow Automation**: DAG-based data processing pipelines
- **Code Generation**: OpenAPI and SQL-driven development
- **Modern Stack**: Go, PostgreSQL/SQLite, Redis

## Development

### Essential Commands

```bash
make generate         # Run after API/SQL changes
make lint            # Check code quality
make dev             # Start backend server
pnpm dev:platform    # Start frontend
```

For detailed development instructions, see [Development Guide](docs/DEVELOPMENT.md).

## License

See [LICENSE](LICENSE) file for details.

## Support

- Email: support@archesai.com
- Issues: https://github.com/archesai/archesai/issues
