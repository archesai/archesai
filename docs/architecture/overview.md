# Architecture Documentation

This section covers the system architecture, design patterns, and structural documentation for
ArchesAI - a high-performance data processing platform with AI-powered capabilities.

## Overview

ArchesAI is built on modern cloud-native principles with a focus on scalability, maintainability,
and developer experience. The platform combines workflow automation, AI integration, and data
processing in a cohesive architecture.

## Documentation Structure

- **[System Architecture](system-design.md)** - Complete system design, patterns, and technical decisions
- **[Authentication](authentication.md)** - JWT-based authentication and authorization architecture
- **[Project Layout](project-layout.md)** - Directory structure, code organization, and conventions

## Core Architecture Principles

### 1. Hexagonal Architecture (Ports & Adapters)

ArchesAI implements **Hexagonal Architecture** to achieve true separation of concerns:

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

#### **Workflows Domain**

- DAG-based workflow definitions
- Pipeline execution engine
- Tool registry and integration
- Run history and monitoring

#### **Content Domain**

- Artifact storage and retrieval
- Vector embeddings for semantic search
- Content processing pipelines
- Metadata management

Each domain maintains:

- **Isolated data models** - No shared database tables
- **Independent APIs** - Domain-specific endpoints
- **Clear boundaries** - No cross-domain imports
- **Event communication** - Domains interact through events

### 3. Code Generation Strategy

ArchesAI leverages extensive code generation to maintain consistency and reduce boilerplate:

#### **OpenAPI-Driven Development**

```yaml
# Define in api/openapi.yaml
paths:
  /organizations:
    post:
      operationId: createOrganization
```

Generates:

- Type definitions (`types.gen.go`)
- HTTP handler interfaces (`http.gen.go`)
- Client SDKs (`web/client/`)

#### **SQL-First Database Layer**

```sql
-- Define in internal/database/queries/
-- name: GetOrganization :one
SELECT * FROM organizations WHERE id = $1;
```

Generates:

- Type-safe database queries
- Repository interfaces
- Transaction helpers

#### **Custom Code Generation**

```yaml
# x-codegen annotations in OpenAPI
x-codegen:
  repository: true
  cache: true
  events: true
```

Generates:

- Repository implementations
- Cache layers with TTL
- Event publishers

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

- Type safety across the stack
- Reduced boilerplate code
- Consistent patterns
- Self-documenting APIs

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

When contributing to ArchesAI:

1. **Respect domain boundaries** - Don't import across domains
2. **Define first, generate second** - Use OpenAPI/SQL before coding
3. **Test at the right level** - Unit tests for logic, integration for APIs
4. **Document decisions** - Add ADRs for significant changes
5. **Optimize thoughtfully** - Measure before optimizing
