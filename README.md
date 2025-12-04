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

Open-source code generation platform that transforms OpenAPI specifications into production-ready applications.

<a href="#-quick-start"><strong>Quick Start</strong></a> ·
<a href="#documentation"><strong>Documentation</strong></a> ·
<a href="#-features"><strong>Features</strong></a> ·
<a href="#support"><strong>Support</strong></a>

## Introduction

**Arches** is a powerful code generation platform that creates complete, production-ready applications from OpenAPI specifications. Through sophisticated templating and code generation, Arches produces full-stack applications with authentication, database layers, API endpoints, and deployment configurations - all from your OpenAPI schema.

## 🚀 Quick Start

```bash
# Install Arches
go install github.com/archesai/archesai/cmd/archesai@v0.20.0-rc.1

# Create your OpenAPI specification
cat > api.yaml << EOF
openapi: 3.1.0
info:
  title: My API
  version: 1.0.0
tags:
  - name: User
    description: User profile management
components:
  schemas:
    User:
      type: object
      x-codegen-schema-type: entity
      required:
        - id
        - name
        - email
        - createdAt
        - updatedAt
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        email:
          type: string
          format: email
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time
    UserResponse:
      type: object
      properties:
        data:
          $ref: '#/components/schemas/User'
      required:
        - data
      additionalProperties: false
paths:
  /users:
    get:
      operationId: GetUser
      tags:
        - User
      summary: Get user
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserResponse'
        '500':
          description: Internal server error
          content:
            application/problem+json:
              schema:
                type: object
                properties:
                  type:
                    type: string
                  title:
                    type: string
                  status:
                    type: integer
                  detail:
                    type: string
EOF

# Generate application code
archesai generate --spec api.yaml --output ./myapp

# Start development server
archesai dev

```

To see what is generaetd, a full example project is available in the [examples](./examples/basic/) directory.

For detailed setup instructions, see [Getting Started](docs/getting-started.md).

## Documentation

- [Getting Started](docs/getting-started.md) - Installation and first app
- [Quickstart](docs/guides/quickstart.md) - 5-minute tutorial
- [CLI Reference](docs/cli-reference.md) - Command documentation
- [Configuration](docs/configuration.md) - Config options
- [Code Generation](docs/guides/code-generation.md) - x-codegen extensions
- [Custom Handlers](docs/guides/custom-handlers.md) - Adding business logic
- [Database](docs/features/database.md) - Multi-database support and migrations
- [Authentication](docs/features/authentication.md) - Auth system
- [Troubleshooting](docs/troubleshooting/common-issues.md) - Common issues
- [Contributing](docs/contributing.md) - Development guide
- [Roadmap](docs/ROADMAP.md) - Planned features

## ✨ Features

### What Gets Generated

From a single OpenAPI specification, Arches generates:

- ✅ **Backend API** - Complete REST API with CRUD operations
- ✅ **Multi-Database Support** - PostgreSQL and SQLite implementations
- ✅ **Automatic Migrations** - SQL migrations generated from your spec
- ✅ **Schema Generation** - Database schemas from OpenAPI definitions
- ✅ **Authentication** - JWT-based auth with role-based access control
- ✅ **API Client** - Type-safe TypeScript/JavaScript SDK
- ✅ **Docker Setup** - Containerization and orchestration configs
- ✅ **Kubernetes Manifests** - Production-ready K8s deployments
- ✅ **Bootstrap Code** - Application initialization and configuration

### Current Capabilities

- 🔥 **Hot Reload Development** - Instant feedback during development
- 🔧 **Extensible Handlers** - Add custom business logic in Go
- 📝 **OpenAPI-First** - Your API specification drives everything
- 🏗️ **Production Ready** - Generated code follows best practices
- 🔒 **Secure by Default** - Built-in security patterns

### Technology Stack

- **Code Generation**: Go templates with parallel processing
- **Backend**: Go with Echo framework
- **Frontend**: Platform UI with React + TypeScript + Vite
- **Database**: PostgreSQL with GORM
- **Deployment**: Docker, Kubernetes, Binary

### Roadmap

See [ROADMAP.md](docs/ROADMAP.md) for planned features including:

- Visual schema designer
- AI-powered features
- Multi-language support (Python, Node.js)
- Cloud deployment integrations

## License

See [LICENSE](LICENSE) file for details.

## Support

- Email: support@archesai.com
- Issues: https://github.com/archesai/archesai/issues
