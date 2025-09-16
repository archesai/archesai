<div align=center>

<a href="https://archesai.com" alt="ArchesAI">
  <img src="./assets/github-hero.png" width=630 alt="ArchesAI Platform">
</a>

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&labelColor=000000)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Proprietary-red?style=flat&labelColor=000000)](LICENSE)
[![API Documentation](https://img.shields.io/badge/API-OpenAPI%203.0-green?style=flat&labelColor=000000)](http://localhost:8080/docs)
[![GitHub Stars](https://img.shields.io/github/stars/archesai/archesai?style=flat&labelColor=000000)](https://github.com/archesai/archesai)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8.svg?style=flat&logo=go&labelColor=000)](https://go.dev)

</div>

# ArchesAI

AI-powered data processing platform with workflow automation and beautiful terminal interfaces.

<a href="#-quick-start"><strong>Quick Start</strong></a> ·
<a href="#documentation"><strong>Documentation</strong></a> ·
<a href="#-features"><strong>Features</strong></a> ·
<a href="#development"><strong>Development</strong></a> ·
<a href="#support"><strong>Support</strong></a>

## Introduction

**ArchesAI** is a high-performance data processing platform that combines AI-powered chat
interfaces, workflow automation, and beautiful terminal UIs to create powerful developer
experiences.

Built for developers who need:

🚀 Fast & efficient data processing<br /> 🤖 Multi-provider AI integration<br /> 💬 Beautiful
terminal chat interface<br /> ⚡ Workflow automation with DAG support<br /> 🔧 Code-first
development with OpenAPI<br />

## 🚀 Quick Start

Get started with ArchesAI in seconds:

### Configuration Viewer (No API Key Needed)

```bash
archesai tui
```

### API Server

```bash
archesai api
```

### AI Chat Interface

```bash
export OPENAI_API_KEY=your-key
archesai tui --chat

# Or with Ollama (local)
archesai tui --chat --provider=ollama
```

For detailed installation and setup instructions, see
[Development Guide](docs/guides/development.md).

## Documentation

### Core Documentation

- [Terminal UI (TUI)](docs/features/tui.md) - Configuration viewer and AI chat interface
- [API Reference](api/openapi.yaml) - OpenAPI specification
- [Development Guide](docs/guides/development.md) - Setup, build, and contribution guide
- [Architecture](docs/architecture/system-design.md) - System design and patterns
- [Contributing](docs/contributing.md) - How to contribute
- [Project Layout](docs/architecture/project-layout.md) - Directory structure and organization
- [Makefile Commands](docs/guides/makefile-commands.md) - Complete command reference

### Packages

#### AI & Chat

- [LLM Package](internal/llm/) - Multi-provider LLM interface with chat clients

#### Core Domains

- [OAuth](internal/oauth) - Authentication & authorization
- [Organizations](internal/organizations/) - Organization management
- [Pipelines](internal/pipelines/) - Pipeline automation
- [Runs](internal/runs/) - Run management

#### Infrastructure

- [Database](internal/database/) - Database layer
- [Config](internal/config/) - Configuration management
- [CLI](internal/cli/) - Command-line interface
- [TUI](internal/tui/) - Terminal user interface

## ✨ Features

- 🤖 **Multi-Provider AI**: Support for OpenAI, Claude, Gemini, Ollama
- 💬 **Chat Interface**: Simple persona-based chat system with session management
- 🎨 **Beautiful TUI**: Terminal interface for configuration and chat
- ⚙️ **Workflow Automation**: DAG-based data processing pipelines
- 🔩 **Code Generation**: OpenAPI and SQL-driven development
- 🚀 **Modern Stack**: Go, PostgreSQL/SQLite, Redis

## Development

### Essential Commands

```bash
make generate     # Run after API/SQL changes
make lint         # Check code quality
make dev          # Start backend server
pnpm dev:platform # Start frontend
```

For detailed development instructions, see [Development Guide](docs/guides/development.md).

## License

See [LICENSE](LICENSE) file for details.

## Support

- Email: support@archesai.com
- Issues: https://github.com/archesai/archesai/issues
