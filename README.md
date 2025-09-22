<div align=center>

<a href="https://archesai.com" alt="Arches">
  <img src="./assets/github-hero.svg" width="100%" alt="Arches Platform">
</a>

<br/>

[![License](https://img.shields.io/badge/license-AGPLv3-purple?style=for-the-badge&labelColor=000000)](LICENSE)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1.1-6BA539?style=for-the-badge&logo=openapi-initiative&labelColor=000000)](https://www.openapis.org)
<br/>
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&labelColor=000000)](https://github.com/archesai/archesai/pkgs/container/archesai)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5?style=for-the-badge&logo=kubernetes&labelColor=000000)](https://kubernetes.io)
[![Helm](https://img.shields.io/badge/Helm-Charts-0F1689?style=for-the-badge&logo=helm&labelColor=000000)](https://helm.sh)
<br/>
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8.svg?style=for-the-badge&logo=go&labelColor=000)](https://go.dev)
[![Made with TypeScript](https://img.shields.io/badge/Made%20with-TypeScript-3178C6.svg?style=for-the-badge&logo=typescript&labelColor=000)](https://www.typescriptlang.org)
<br/>
[![OpenAI](https://img.shields.io/badge/OpenAI-Compatible-ffffff?style=for-the-badge&logo=openai&labelColor=000000&logoColor=white)](https://openai.com)
[![Anthropic](https://img.shields.io/badge/Anthropic-Compatible-FF6600?style=for-the-badge&labelColor=000000)](https://anthropic.com)
[![Ollama](https://img.shields.io/badge/Ollama-Compatible-000000?style=for-the-badge&labelColor=ffffff)](https://ollama.ai)
<br/>
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql&labelColor=000000)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=for-the-badge&logo=redis&labelColor=000000)](https://redis.io)

</div>

# Arches

AI-powered data processing platform with workflow automation and beautiful terminal interfaces.

<a href="#-quick-start"><strong>Quick Start</strong></a> ·
<a href="#documentation"><strong>Documentation</strong></a> ·
<a href="#-features"><strong>Features</strong></a> ·
<a href="#support"><strong>Support</strong></a>

## Introduction

**Arches** is a high-performance data processing platform that combines AI-powered chat
interfaces, workflow automation, and beautiful terminal UIs to create powerful developer
experiences.

Built for developers who need:

🚀 Fast & efficient data processing<br /> 🤖 Multi-provider AI integration<br /> 💬 Beautiful
terminal chat interface<br /> ⚡ Workflow automation with DAG support<br /> 🔧 Code-first
development with OpenAPI<br />

## 🚀 Quick Start

Get started with Arches in seconds:

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

- [OAuth](internal/auth) - Authentication & authorization
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

## License

See [LICENSE](LICENSE) file for details.

## Support

- Email: support@archesai.com
- Issues: https://github.com/archesai/archesai/issues
