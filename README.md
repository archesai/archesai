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

AI-powered full-stack app builder that generates production-ready applications from OpenAPI schemas.

<a href="#-quick-start"><strong>Quick Start</strong></a> ·
<a href="#documentation"><strong>Documentation</strong></a> ·
<a href="#-features"><strong>Features</strong></a> ·
<a href="#support"><strong>Support</strong></a>

## Introduction

**Arches** is an open-source app builder that transforms OpenAPI schemas into complete, production-ready applications. Using AI and powerful code generation, Arches creates full-stack apps with authentication, CRUD operations, real-time features, and deployment configurations - all from a single API specification.

Built for developers who want:

🎯 **OpenAPI → Full App** - Generate complete applications from schemas<br />
🤖 **AI-Powered Development** - Natural language to OpenAPI, automatic handler implementation<br />
🔥 **Hot Reload Everything** - Instant regeneration as you modify schemas<br />
🎨 **Visual Schema Designer** - Build APIs visually or with AI assistance<br />
🚀 **Multi-Language Support** - Python, JavaScript, or Go for custom logic<br />
📦 **Deploy Anywhere** - Docker, Kubernetes, or single binary<br />

## 🚀 Quick Start

```bash
# Create a new app
arches new my-app

# Start the development server
cd my-app
arches dev

# Open the studio
open http://localhost:3000
```

For detailed installation and setup instructions, see
[Development Guide](docs/guides/development.md).

## Documentation

### Core Documentation

- [Codebase Analysis](docs/codebase-analysis.md) - Detailed analysis of current architecture
- [Development Guide](docs/guides/development.md) - Setup, build, and contribution guide
- [Architecture](docs/architecture/system-design.md) - System design and patterns
- [Code Generation](docs/guides/codegen.md) - How the code generation works
- [API Reference](api/openapi.yaml) - OpenAPI specification format
- [Makefile Commands](docs/guides/makefile-commands.md) - Complete command reference

## ✨ Features

### What Gets Generated

From a single OpenAPI specification, Arches generates:

- ✅ **Backend API** - Complete REST API with CRUD operations
- ✅ **Database Layer** - Migrations, models, and type-safe queries
- ✅ **Frontend App** - React application with routing and components
- ✅ **Authentication** - JWT, OAuth, magic links, and sessions
- ✅ **WebSocket Support** - Real-time features out of the box
- ✅ **RBAC** - Role-based access control from security schemes
- ✅ **API Client** - Type-safe SDK with validation
- ✅ **Docker Setup** - Containerization and orchestration
- ✅ **Tests** - Unit and integration tests
- ✅ **Documentation** - Auto-generated API docs

### Development Experience

- 🔥 **Hot Reload** - Instant regeneration on schema changes
- 🎨 **Visual Designer** - Build schemas with drag-and-drop
- 🤖 **AI Assistant** - Generate schemas from natural language
- 🔧 **Custom Logic** - Write handlers in Python, JavaScript, or Go
- 📊 **Live Preview** - See your app as you build it
- 🚀 **One-Click Deploy** - Deploy to cloud platforms instantly

### Technology Stack

- **Code Generation**: Template-based with parallel processing
- **Backend Languages**: Go (default), Python, Node.js (coming soon)
- **Frontend**: React 19 + TypeScript + Vite
- **Database**: PostgreSQL, SQLite, MySQL (coming soon)
- **Deployment**: Docker, Kubernetes, Single Binary

## License

See [LICENSE](LICENSE) file for details.

## Support

- Email: support@archesai.com
- Issues: https://github.com/archesai/archesai/issues
