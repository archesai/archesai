# ArchesAI

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Proprietary-red?style=flat-square)](LICENSE)
[![API Documentation](https://img.shields.io/badge/API-OpenAPI%203.0-green?style=flat-square)](http://localhost:8080/docs)

A high-performance data processing platform for managing, analyzing, and transforming diverse data assets. Built with Go using hexagonal architecture and domain-driven design principles.

## Features

### 🚀 Core Capabilities

- **Multi-format content processing** - Files, audio, text, images, and web content
- **Vector embeddings & semantic search** - Advanced similarity matching and retrieval
- **DAG-based workflows** - Build complex data processing pipelines
- **Transformation tools** - Text-to-speech, text-to-image, transcription, and more
- **Multi-tenant architecture** - Organization-based data isolation

### 🔒 Security & Auth

- JWT-based authentication with refresh tokens
- Database-backed session management
- Multi-organization support with data isolation
- Role-based access control (coming soon)

### 🛠 Developer Experience

- OpenAPI 3.0 specification with auto-generated types
- TypeScript SDK for frontend integration
- Comprehensive code generation from specs
- Hot-reload development environment

## Quick Start

### Prerequisites

- **Go 1.21+**
- **PostgreSQL 15+** with pgvector extension
- **Node.js 20+** and pnpm 8+
- **Make** for build commands

### Installation

```bash
# Clone repository
git clone https://github.com/archesai/archesai.git
cd archesai

# Setup environment
cp .env.example .env
# Edit .env with your configuration

# Install dependencies
make tools              # Install Go tools
go mod download         # Backend dependencies
pnpm install           # Frontend dependencies

# Setup database
createdb archesai
psql archesai -c "CREATE EXTENSION IF NOT EXISTS vector;"
make migrate-up

# Generate code
make generate          # Generate all code (SQLC, OpenAPI, codegen)

# Start development
make dev               # Backend with hot-reload (port 8080)
pnpm dev:platform      # Frontend in another terminal (port 5173)
```

### Access Points

- **API**: http://localhost:8080
- **API Docs**: http://localhost:8080/docs
- **Web UI**: http://localhost:5173

## Project Structure

```
archesai/
├── api/                    # OpenAPI specifications
├── cmd/                    # Application entry points
│   ├── archesai/          # Main server
│   └── codegen/           # Code generation tool
├── internal/
│   ├── app/               # Application wiring
│   ├── auth/              # Authentication domain
│   ├── organizations/     # Organization management
│   ├── workflows/         # Pipeline workflows
│   ├── content/           # Content management
│   ├── database/          # Database layer
│   │   ├── postgresql/    # SQLC generated code
│   │   ├── queries/       # SQL queries
│   │   └── migrations/    # Database migrations
│   └── config/            # Configuration
├── web/                   # Frontend monorepo
│   ├── platform/          # React application
│   ├── client/            # TypeScript SDK
│   └── ui/                # Component library
└── docs/                  # Documentation
```

## Development

### Essential Commands

```bash
# Development
make dev              # Start with hot reload
make build            # Build all binaries
make test             # Run tests
make lint             # Run linters

# Code Generation
make generate         # Run all generators
make generate-sqlc    # Database code
make generate-oapi    # OpenAPI types
make generate-codegen # Domain code

# Database
make migrate-up       # Apply migrations
make migrate-down     # Rollback
make migrate-create name=feature_name
```

### Adding Features

1. Define API in `api/openapi.yaml`
2. Add SQL queries in `internal/database/queries/`
3. Run `make generate`
4. Implement business logic in domain service
5. Wire dependencies in `internal/app/app.go`
6. Run `make lint && make test`

## API Overview

### Authentication

- `POST /api/auth/register` - Register user
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh token
- `GET /api/auth/me` - Current user

### Organizations

- `GET /api/organizations` - List organizations
- `POST /api/organizations` - Create organization
- `GET /api/organizations/{id}` - Get details
- `PUT /api/organizations/{id}` - Update
- `DELETE /api/organizations/{id}` - Delete

### Workflows

- `GET /api/workflows` - List workflows
- `POST /api/workflows` - Create workflow
- `POST /api/workflows/{id}/runs` - Execute

### Content

- `GET /api/content/artifacts` - List artifacts
- `POST /api/content/artifacts` - Upload
- `POST /api/content/artifacts/{id}/process` - Process

## Configuration

Environment variables use `ARCHESAI_` prefix:

```bash
# Database
ARCHESAI_DATABASE_URL=postgres://user:pass@localhost/archesai?sslmode=disable

# Server
ARCHESAI_SERVER_PORT=8080
ARCHESAI_SERVER_HOST=0.0.0.0

# Authentication
ARCHESAI_JWT_SECRET=your-secret-key
ARCHESAI_JWT_ACCESS_TOKEN_DURATION=15m
ARCHESAI_JWT_REFRESH_TOKEN_DURATION=7d

# Logging
ARCHESAI_LOGGING_LEVEL=info
ARCHESAI_LOGGING_FORMAT=json
```

## Documentation

- [Development Guide](docs/DEVELOPMENT.md) - Technical development details
- [Architecture](docs/ARCHITECTURE.md) - System design and patterns
- [API Reference](http://localhost:8080/docs) - Interactive API documentation

## License

Proprietary - All rights reserved

## Support

- Email: support@archesai.com
- Issues: https://github.com/archesai/archesai/issues
