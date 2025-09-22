# Architecture Documentation

This section covers the system architecture, design patterns, and structural documentation for
Arches - a high-performance data processing platform with AI-powered capabilities.

## Overview

Arches is built on modern cloud-native principles with a focus on scalability, maintainability,
and developer experience. The platform combines workflow automation, AI integration, and data
processing in a cohesive architecture.

## Documentation Structure

- **[System Architecture](system-design.md)** - Complete system design, patterns, and technical decisions
- **[Authentication](authentication.md)** - JWT-based authentication and authorization architecture
- **[Project Layout](project-layout.md)** - Directory structure, code organization, and conventions

## Core Architecture Principles

### 1. Hexagonal Architecture (Ports & Adapters)

Arches implements **Hexagonal Architecture** to achieve true separation of concerns:

```text
┌─────────────────────────────────────────────┐
│             Presentation Layer              │
│         (HTTP Handlers, GraphQL, CLI)       │
└─────────────┬───────────────────┬───────────┘
              │                   │
┌─────────────▼───────────────────▼───────────┐
│            Application Layer                │
│        (Use Cases, Business Logic)          │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │          Domain Core                │   │
│  │   (Entities, Value Objects, Rules)  │   │
│  └─────────────────────────────────────┘   │
│                                             │
└─────────────┬───────────────────┬───────────┘
              │                   │
┌─────────────▼───────────────────▼───────────┐
│          Infrastructure Layer               │
│    (Database, Cache, External Services)     │
└─────────────────────────────────────────────┘
```

**Benefits:**

- Business logic remains independent of frameworks
- Easy to test with mock implementations
- Flexible infrastructure changes without affecting core logic
- Clear separation between technical and business concerns

### 2. Domain-Driven Design (DDD)

The system is organized into distinct bounded contexts, each representing a core business domain:

#### **Auth Domain**

- User authentication and session management
- JWT token generation and validation
- OAuth provider integration
- Password reset and email verification

#### **Organizations Domain**

- Multi-tenant organization management
- Member invitations and role management
- Billing and subscription handling
- Organization-level settings

#### **Pipelines Domain**

- DAG-based pipeline definitions
- Pipeline execution engine
- Task orchestration and scheduling
- Execution monitoring

#### **Artifacts Domain**

- Artifact storage and retrieval
- File management and versioning
- Content processing
- Metadata management

#### **Tools Domain**

- Tool registry and integration
- Tool capability management
- External service integration
- Tool execution framework

#### **Runs Domain**

- Pipeline run management
- Run history and monitoring
- Execution state tracking
- Result aggregation

#### **Sessions Domain**

- User session management
- Session persistence and caching
- Session-based authentication
- Activity tracking

Each domain maintains:

- **Isolated data models** - No shared database tables
- **Independent APIs** - Domain-specific endpoints
- **Clear boundaries** - No cross-domain imports
- **Event communication** - Domains interact through events

### 3. Code Generation Strategy

Arches leverages a unified code generation system to maintain consistency and reduce boilerplate:

#### **OpenAPI-Driven Development**

```yaml
# Define in api/openapi.yaml with x-codegen annotations
components:
  schemas:
    Organization:
      x-codegen:
        repository:
          operations:
            - create
            - get
            - update
            - delete
            - list
        service:
          enabled: true
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
```

Generates:

- Type definitions (`types.gen.go`)
- Repository interface (`repository.gen.go`)
- PostgreSQL implementation (`postgres.gen.go`)
- SQLite implementation (`sqlite.gen.go`)
- Service interface (`service.gen.go`)
- HTTP server implementation (`server.gen.go`)
- API client interface (`handler.gen.go`)
- Test mocks (`mocks_test.gen.go`)

#### **SQL-First Database Layer**

```sql
-- Define in internal/database/queries/
-- name: GetOrganization :one
SELECT * FROM organizations WHERE id = $1;
```

Generates:

- Type-safe database queries via SQLC
- Database models and interfaces

#### **Unified Code Generation**

The unified generator reads x-codegen annotations to produce:

- **Repository Layer**: Interface and database implementations (PostgreSQL/SQLite)
- **Service Layer**: Business logic interfaces and HTTP servers
- **Events Layer**: Optional event publishing with NATS/Redis
- **Configuration**: Default values and environment mappings

### 4. Technology Stack

#### **Backend**

- **Language**: Go 1.21+ for performance and simplicity
- **Framework**: Echo for HTTP routing
- **Database**: PostgreSQL with pgvector extension
- **Cache**: Redis for session and data caching
- **Message Queue**: NATS for event streaming

#### **Frontend**

- **Framework**: React with TanStack Router
- **UI Library**: Custom component library
- **State Management**: TanStack Query
- **Build Tool**: Vite for fast development

#### **Infrastructure**

- **Container**: Docker with multi-stage builds
- **Orchestration**: Kubernetes with Helm charts
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logging with Loki

### 5. Scalability Patterns

#### **Horizontal Scaling**

- Stateless services enable easy horizontal scaling
- Load balancing across multiple instances
- Database connection pooling
- Redis cluster for distributed caching

#### **Async Processing**

- Background job queues for heavy operations
- Event-driven architecture for decoupling
- Workflow execution in separate workers
- Rate limiting and backpressure handling

#### **Performance Optimization**

- Query optimization with indexes
- Caching at multiple layers
- CDN for static assets
- Database read replicas

## Key Design Decisions

### Why Go?

- Excellent performance for I/O-bound operations
- Simple deployment (single binary)
- Strong standard library
- Great tooling and testing support

### Why PostgreSQL?

- ACID compliance for critical data
- pgvector for AI/ML embeddings
- Rich feature set (JSONB, full-text search)
- Mature ecosystem

### Why Code Generation?

- **Type safety across the stack** - Generated types ensure compile-time safety
- **Reduced boilerplate code** - Automatic generation of CRUD operations
- **Consistent patterns** - Unified structure across all domains
- **Self-documenting APIs** - OpenAPI serves as single source of truth
- **Automatic mock generation** - Mockery generates test mocks from interfaces
- **Multi-database support** - Same interface, multiple implementations

## Architecture Evolution

The architecture is designed to evolve:

1. **Current State**: Monolithic with domain separation
2. **Next Phase**: Service extraction for high-traffic domains
3. **Future State**: Full microservices with service mesh

## Getting Started

To understand the architecture in practice:

1. Review the [System Design](system-design.md) for detailed patterns
2. Explore the [Project Layout](project-layout.md) to understand code organization
3. Check [Authentication](authentication.md) for security architecture
4. See working examples in the `internal/` directory

## Architecture Guidelines

When contributing to Arches:

1. **Respect domain boundaries** - Don't import across domains
2. **Define first, generate second** - Use OpenAPI/SQL before coding
3. **Test at the right level** - Unit tests for logic, integration for APIs
4. **Document decisions** - Add ADRs for significant changes
5. **Optimize thoughtfully** - Measure before optimizing
