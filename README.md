# Arches AI

> A comprehensive data processing platform for managing, analyzing, and transforming diverse data assets

## Introduction

**Arches AI** is a comprehensive data processing platform designed to empower businesses to efficiently manage, analyze, and transform their diverse data assets. Similar to Palantir Foundry, Arches AI enables organizations to upload various types of content—including files, audio, text, images, and websites—and index them for seamless parsing, querying, and transformation. Leveraging advanced embedding models and a suite of transformation tools, Arches AI provides flexible and powerful data processing capabilities tailored to meet the unique needs of different industries.

## Core Features

### Data Upload and Indexing

- **Multi-Format Support:** Seamlessly upload and manage files, audio, text, images, and websites.
- **Automated Indexing:** Efficiently index all uploaded content for quick retrieval and management.
- **Vector Embeddings:** Generate and store embeddings for semantic search and similarity matching.

### Transformation Tools

- **Text-to-Speech:** Convert textual data into natural-sounding audio.
- **Text-to-Image:** Generate high-quality images based on textual descriptions.
- **Text-to-Text:** Advanced text manipulation, generation, and transformation capabilities.
- **Document Processing:** Extract and convert content from various file types (PDF, DOCX, etc.) into structured text.
- **Audio Transcription:** Convert audio content to text with high accuracy.

### Embedding Models

- **Advanced Embeddings:** Utilize state-of-the-art models to embed text content into vector representations.
- **Semantic Search:** Enable sophisticated querying and semantic search for enhanced data accessibility.
- **Similarity Matching:** Find related content based on vector similarity.

### Data Querying and Transformation

- **Intuitive Query Interface:** User-friendly tools for querying indexed data with ease.
- **Advanced Filtering:** Filter data by metadata, labels, and custom attributes.
- **Data Transformation Tools:** Flexible tools to transform data to meet specific business requirements.
- **Batch Processing:** Process large volumes of data efficiently.

### Workflow Building

- **Custom Workflows:** Design and implement data processing workflows using individual tools through the workflows domain.
- **Automation:** Automate complex data workflows tailored to organizational needs.
- **Directed Acyclical Graph (DAG):** Workflows are DAGs, representing all possible processing chains.
- **Pipeline Runs:** Track and monitor workflow execution with detailed run history and status.
- **Tool Orchestration:** Chain multiple transformation tools together for complex processing.

### Authentication & Security

- **JWT-based Authentication:** Secure token-based authentication with refresh tokens.
- **Session Management:** Database-backed session tracking and management.
- **Multi-Organization Support:** Isolate data and permissions across different organizations.
- **Role-Based Access Control:** Fine-grained permissions system (coming soon).

### Support and Consulting

- **Integration Support:** Expert assistance in integrating Arches AI with existing systems.
- **Data Strategy Consulting:** Help businesses optimize their data strategies for maximum impact.
- **Custom Tool Development:** Build specialized transformation tools for unique requirements.

## Design Concepts

### Scalability

- **Modular Architecture:** Easily add or remove components to scale with business growth.
- **Cloud-Native Infrastructure:** Built on scalable cloud platforms to handle increasing data volumes.
- **Horizontal Scaling:** Support for distributed processing across multiple nodes.

### Usability

- **Intuitive Interface:** User-friendly dashboards and interfaces to lower the barrier to entry.
- **Customizable Workflows:** Flexible pipeline creation to suit various business processes.
- **Real-time Feedback:** Immediate processing status and results visualization.

### Security

- **Data Encryption:** Ensure data is securely stored and transmitted using advanced encryption standards.
- **Access Controls:** Robust authentication and authorization mechanisms to protect sensitive data.
- **Audit Logging:** Comprehensive logging of all data access and modifications.

### Integration

- **RESTful APIs:** Well-documented REST APIs for seamless integration.
- **TypeScript SDK:** Auto-generated TypeScript client for frontend applications.
- **Webhook Support:** Event-driven integrations with external systems (coming soon).
- **Third-Party Integrations:** Support for popular services and data sources.

## Use Cases by Industry

### Finance

- **Fraud Detection:** Analyze transaction data to identify and prevent fraudulent activities.
- **Risk Management:** Assess and manage financial risks through comprehensive data analysis.
- **Customer Insights:** Gain deeper understanding of customer behaviors and preferences.
- **Regulatory Compliance:** Automate document processing for compliance reporting.

### Healthcare

- **Medical Records Management:** Organize and analyze patient data for improved healthcare delivery.
- **Research and Development:** Facilitate medical research by managing and processing large datasets.
- **Clinical Trial Analysis:** Process and analyze clinical trial documentation.
- **Medical Image Analysis:** Extract insights from medical imaging data.

### Legal

- **Document Management:** Organize and search through large volumes of legal documents.
- **Case Analysis:** Analyze case data to identify patterns and support legal strategies.
- **Contract Review:** Automated extraction and analysis of contract terms.
- **E-Discovery:** Efficient document discovery for litigation support.

### Technology

- **Code Documentation:** Generate and maintain technical documentation from codebases.
- **Log Analysis:** Process and analyze system logs for insights.
- **User Behavior Analytics:** Understand user interactions and improve products.
- **API Documentation:** Auto-generate and maintain API documentation.

### Manufacturing

- **Quality Control:** Analyze production data to identify quality issues.
- **Supply Chain Optimization:** Process supplier and logistics data.
- **Predictive Maintenance:** Analyze sensor data to predict equipment failures.
- **Production Planning:** Optimize production schedules based on data insights.

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+ with pgvector extension
- Node.js 20+ and pnpm 8+
- Docker & Docker Compose (optional, for containerized development)
- Make (for running build commands)

### Quick Start

1. **Clone the repository**

   ```bash
   git clone https://github.com/archesai/archesai.git
   cd archesai
   ```

2. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

   Key environment variables:

   ```bash
   # Database
   ARCHESAI_DATABASE_URL=postgres://user:pass@localhost/archesai?sslmode=disable

   # Server
   ARCHESAI_SERVER_PORT=8080
   ARCHESAI_SERVER_HOST=0.0.0.0

   # Authentication
   ARCHESAI_JWT_SECRET=your-secret-key-change-in-production
   ARCHESAI_JWT_ACCESS_TOKEN_DURATION=15m
   ARCHESAI_JWT_REFRESH_TOKEN_DURATION=7d

   # Logging
   ARCHESAI_LOGGING_LEVEL=info
   ARCHESAI_LOGGING_FORMAT=json
   ```

3. **Install dependencies**

   ```bash
   # Install Go tools
   make tools

   # Backend dependencies
   go mod download

   # Frontend dependencies
   pnpm install
   ```

4. **Set up the database**

   ```bash
   # Create database (if not exists)
   createdb archesai

   # Enable pgvector extension
   psql archesai -c "CREATE EXTENSION IF NOT EXISTS vector;"

   # Run migrations
   make migrate-up
   ```

5. **Generate code**

   ```bash
   # Generate all code (SQLC, OpenAPI, adapters)
   make generate
   ```

6. **Start development servers**

   ```bash
   # Backend API (with hot reload)
   make dev

   # Frontend (in another terminal)
   pnpm dev:platform
   ```

7. **Access the application**
   - API: http://localhost:8080
   - API Documentation: http://localhost:8080/docs
   - Web UI: http://localhost:5173

### Docker Development

For a containerized development environment:

```bash
# Build and start all services
docker-compose up

# Run migrations
docker-compose exec api make migrate-up

# Access services
# - API: http://localhost:8080
# - Web UI: http://localhost:5173
# - PostgreSQL: localhost:5432
```

## Project Structure

ArchesAI uses a hexagonal (ports and adapters) architecture with Domain-Driven Design principles:

```
archesai/
├── api/                     # OpenAPI specifications
│   ├── openapi.yaml        # Main OpenAPI spec
│   ├── paths/              # Path definitions
│   └── components/         # Reusable components
├── cmd/                    # Application entry points
│   ├── archesai/          # Main server binary
│   ├── codegen/           # Code generation tool
│   └── worker/            # Background worker
├── internal/
│   ├── app/               # Application assembly
│   │   ├── deps.go        # Dependency injection
│   │   ├── routes.go      # HTTP route registration
│   │   └── server.go      # Server configuration
│   ├── auth/              # Authentication domain
│   │   ├── domain/        # Business logic & entities
│   │   ├── adapters/      # Type converters & repos
│   │   ├── handlers/      # HTTP handlers
│   │   └── generated/     # Generated code
│   ├── organizations/     # Organization management
│   ├── workflows/         # Pipeline workflows
│   ├── content/          # Content & artifacts
│   ├── database/         # Database layer
│   │   ├── postgresql/   # Generated SQLC code
│   │   ├── queries/      # SQL queries
│   │   └── migrations/   # Database migrations
│   └── config/           # Configuration management
├── web/                  # Frontend monorepo
│   ├── platform/         # Main React application
│   │   ├── src/
│   │   │   ├── routes/   # File-based routing
│   │   │   ├── components/
│   │   │   └── hooks/
│   │   └── package.json
│   ├── client/          # Generated TypeScript client
│   └── ui/              # Shared component library
├── deployments/         # Deployment configurations
│   ├── kubernetes/      # K8s manifests
│   └── docker/          # Dockerfiles
├── docs/                # Documentation
│   ├── DEVELOPMENT.md   # Development guide
│   ├── ARCHITECTURE.md  # Architecture details
│   └── API.md          # API documentation
└── scripts/            # Utility scripts
```

### Domain Structure

Each domain follows hexagonal architecture:

```
domain/
├── domain/              # Core business logic
│   ├── entities.go     # Domain entities
│   ├── ports.go        # Repository interfaces
│   ├── usecase.go      # Business use cases
│   └── types.gen.go    # Generated types
├── adapters/           # Adapters layer
│   ├── postgres/       # PostgreSQL implementation
│   └── adapters.gen.go # Generated converters
├── handlers/           # HTTP handlers
│   └── http/
│       └── handler.go
└── generated/          # Generated code
    └── api/           # OpenAPI generated types
```

## Development

### Essential Commands

```bash
# Development
make dev              # Start with hot reload
make build            # Build all binaries
make test             # Run all tests
make lint             # Run all linters

# Code Generation
make generate         # Run all generators
make generate-sqlc    # Generate database code
make generate-oapi    # Generate OpenAPI code
make generate-adapters # Generate type converters

# Database
make migrate-up       # Apply migrations
make migrate-down     # Rollback last migration
make migrate-create name=<name> # Create new migration

# Docker
make docker-build     # Build Docker image
make docker-run       # Run in Docker
```

### Adding New Features

1. **Define API** in `api/openapi.yaml`
2. **Run generators**: `make generate`
3. **Create migration**: `make migrate-create name=feature_name`
4. **Write SQL queries** in `internal/database/queries/`
5. **Implement business logic** in domain's `usecase.go`
6. **Implement handlers** in domain's `handlers/http/`
7. **Wire dependencies** in `internal/app/deps.go`
8. **Add routes** in `internal/app/routes.go`
9. **Run tests**: `make test`
10. **Lint code**: `make lint`

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific domain tests
go test ./internal/auth/...

# Frontend tests
pnpm test
pnpm test:e2e
```

## API Documentation

The API follows RESTful principles and is fully documented using OpenAPI 3.0.

### Authentication Endpoints

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout user
- `GET /api/auth/me` - Get current user

### Organization Endpoints

- `GET /api/organizations` - List organizations
- `POST /api/organizations` - Create organization
- `GET /api/organizations/{id}` - Get organization
- `PUT /api/organizations/{id}` - Update organization
- `DELETE /api/organizations/{id}` - Delete organization

### Workflow Endpoints

- `GET /api/workflows` - List workflows
- `POST /api/workflows` - Create workflow
- `GET /api/workflows/{id}` - Get workflow
- `PUT /api/workflows/{id}` - Update workflow
- `DELETE /api/workflows/{id}` - Delete workflow
- `POST /api/workflows/{id}/runs` - Execute workflow

### Content Endpoints

- `GET /api/content/artifacts` - List artifacts
- `POST /api/content/artifacts` - Upload artifact
- `GET /api/content/artifacts/{id}` - Get artifact
- `DELETE /api/content/artifacts/{id}` - Delete artifact
- `POST /api/content/artifacts/{id}/process` - Process artifact

For detailed API documentation, run the server and visit http://localhost:8080/docs

## Configuration

Configuration is managed through environment variables and config files:

### Environment Variables

All environment variables use the `ARCHESAI_` prefix:

```bash
# Database Configuration
ARCHESAI_DATABASE_URL           # PostgreSQL connection string
ARCHESAI_DATABASE_POOL_SIZE     # Connection pool size (default: 10)
ARCHESAI_DATABASE_MAX_IDLE_TIME # Max idle time (default: 30m)

# Server Configuration
ARCHESAI_SERVER_PORT            # Server port (default: 8080)
ARCHESAI_SERVER_HOST            # Server host (default: 0.0.0.0)
ARCHESAI_SERVER_READ_TIMEOUT    # Read timeout (default: 30s)
ARCHESAI_SERVER_WRITE_TIMEOUT   # Write timeout (default: 30s)

# Authentication
ARCHESAI_JWT_SECRET             # JWT signing secret (required)
ARCHESAI_JWT_ACCESS_TOKEN_DURATION  # Access token duration
ARCHESAI_JWT_REFRESH_TOKEN_DURATION # Refresh token duration

# Logging
ARCHESAI_LOGGING_LEVEL          # Log level (debug, info, warn, error)
ARCHESAI_LOGGING_FORMAT         # Log format (json, text)

# Feature Flags
ARCHESAI_FEATURES_ENABLE_WEBHOOKS   # Enable webhook support
ARCHESAI_FEATURES_ENABLE_ANALYTICS  # Enable analytics
```

### Configuration Files

- `.env` - Local environment variables
- `config.yaml` - Default configuration values
- `config.local.yaml` - Local overrides (gitignored)

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on:

- Code of Conduct
- Development workflow
- Coding standards
- Testing requirements
- Pull request process

## Documentation

- [Development Guide](docs/DEVELOPMENT.md) - Detailed development instructions
- [Architecture](docs/ARCHITECTURE.md) - System architecture and design decisions
- [API Documentation](docs/API.md) - Complete API reference
- [Deployment Guide](docs/DEPLOYMENT.md) - Production deployment instructions
- [Claude.md](.claude/CLAUDE.md) - AI assistant instructions

## License

Proprietary - All rights reserved

## Support

For support, consulting, or enterprise inquiries:

- Email: support@archesai.com
- Documentation: https://docs.archesai.com
- Issues: https://github.com/archesai/archesai/issues

## Acknowledgments

Built with:

- [Echo](https://echo.labstack.com/) - High performance Go web framework
- [SQLC](https://sqlc.dev/) - Type-safe SQL for Go
- [pgvector](https://github.com/pgvector/pgvector) - Vector similarity search for PostgreSQL
- [React](https://react.dev/) - UI library
- [TanStack Router](https://tanstack.com/router) - Type-safe routing
- [Vite](https://vitejs.dev/) - Fast build tool
