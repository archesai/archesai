# Features Documentation

Arches is a powerful code generation platform that transforms OpenAPI specifications
into production-ready applications. This section documents the key features and capabilities.

## Core Features

### [Code Generation](code-generation.md)

The heart of Arches - transforming OpenAPI specs into complete applications:

- Generate models, controllers, and handlers from schemas
- Create database migrations and repositories
- Build TypeScript/JavaScript client SDKs
- Generate test files and mocks

### [Authentication System](auth.md)

Comprehensive authentication features generated for your apps:

- JWT-based authentication with refresh tokens
- Email/password registration and login
- OAuth provider integration
- Magic link (passwordless) authentication
- Session management across devices

### [CLI Tools](../cli-reference.md)

Powerful command-line interface for development:

- `archesai generate` - Generate code from specifications
- `archesai dev` - Development server with hot reload
- `archesai config` - Configuration management
- Shell completions for all major shells

### [Development Experience](development.md)

Tools and features for efficient development:

- Hot reload for instant feedback
- TUI mode for interactive development
- Integrated testing framework
- Mock generation for unit tests

## Platform Capabilities

### OpenAPI-First Development

- **Single Source of Truth**: Your OpenAPI spec drives everything
- **x-codegen Extensions**: Custom annotations for advanced generation
- **Multi-file Support**: Split specs across multiple files
- **Validation**: Built-in schema validation

### Database Integration

- **PostgreSQL Support**: First-class PostgreSQL integration
- **Migrations**: Auto-generated SQL migrations
- **Type Safety**: Generated models with proper types
- **GORM Integration**: ORM support out of the box

### Deployment Ready

- **Docker Support**: Containerization configs included
- **Kubernetes Manifests**: Production-ready K8s deployment
- **Health Checks**: Built-in liveness and readiness probes
- **Configuration Management**: Environment-based configs

### Testing Infrastructure

- **Unit Tests**: Generated test scaffolding
- **Integration Tests**: API endpoint testing
- **Mock Generation**: Automatic mock creation with mockery
- **Coverage Reports**: Built-in coverage tracking

## Generated Components

When you use Arches to generate an application, you get:

### Backend (Go)

- **Models**: Data structures from OpenAPI schemas
- **Controllers**: HTTP request handlers
- **Handlers**: Business logic layer
- **Repositories**: Database access layer
- **Migrations**: SQL migration scripts
- **Bootstrap**: Application initialization
- **Middleware**: Auth, CORS, logging

### Frontend Support

- **TypeScript SDK**: Type-safe client library
- **API Client**: Axios-based HTTP client
- **Type Definitions**: Full TypeScript types
- **Validation**: Runtime validation with zod

### Infrastructure

- **Dockerfile**: Multi-stage production build
- **docker-compose.yml**: Local development setup
- **Kubernetes YAML**: Deployment manifests
- **Helm Charts**: For advanced deployments

## Configuration Features

### Application Configuration

- YAML-based configuration
- Environment variable overrides
- Multi-environment support
- Secret management

### Generation Options

- Custom output paths
- Template overrides
- Selective generation
- Bundle mode for specs

## Security Features

### Built-in Security

- **Password Hashing**: bcrypt by default
- **JWT Signing**: Secure token generation
- **CORS Support**: Configurable CORS policies
- **Rate Limiting**: Built-in rate limiters
- **Input Validation**: Schema-based validation

### Authentication Methods

- Traditional email/password
- OAuth 2.0 providers
- Magic links
- API keys (coming soon)

## Developer Tools

### CLI Features

- Interactive TUI mode
- Shell completions
- Config validation
- Version management

### Development Server

- Hot reload on file changes
- Concurrent service running
- Error reporting
- Request logging

## Extensibility

### Custom Logic

- Add business logic to generated handlers
- Override generated methods
- Custom middleware
- Plugin system (planned)

### Template Customization

- Override default templates
- Create custom templates
- Template variables
- Conditional generation

## Integration Capabilities

### Database Support

- PostgreSQL (primary)
- SQLite (development)
- MySQL (coming soon)
- MongoDB (planned)

### Cache & Session

- Redis integration
- In-memory caching
- Session management
- Distributed caching

### External Services

- S3/MinIO for storage
- Email providers
- SMS providers (planned)
- Payment gateways (planned)

## See Also

- [Getting Started](../getting-started.md) - Quick start guide
- [CLI Reference](../cli-reference.md) - Complete CLI documentation
- [Code Generation Guide](../guides/code-generation.md) - How generation works
- [Development Guide](../guides/development.md) - Development workflow
- [ROADMAP](../ROADMAP.md) - Planned features and enhancements
