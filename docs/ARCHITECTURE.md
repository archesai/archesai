# Architecture Documentation

## Overview

ArchesAI follows a **Hexagonal Architecture** (Ports and Adapters) pattern combined with **Domain-Driven Design (DDD)** principles. This architecture ensures separation of concerns, testability, and maintainability while keeping the business logic independent of external dependencies.

## Core Architectural Principles

### 1. Hexagonal Architecture

The hexagonal architecture isolates the core business logic from external concerns:

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│                  (HTTP Handlers, CLI)                    │
└────────────────────────┬───────────────────────────────┘
                         │
┌────────────────────────▼───────────────────────────────┐
│                    Application Layer                    │
│               (Use Cases, Orchestration)                │
└────────────────────────┬───────────────────────────────┘
                         │
┌────────────────────────▼───────────────────────────────┐
│                     Domain Layer                        │
│            (Entities, Business Rules, Ports)            │
└────────────────────────┬───────────────────────────────┘
                         │
┌────────────────────────▼───────────────────────────────┐
│                 Infrastructure Layer                    │
│        (Database, External APIs, File System)           │
└─────────────────────────────────────────────────────────┘
```

### 2. Domain-Driven Design

Each bounded context represents a distinct business domain:

- **Auth Domain**: User authentication and authorization
- **Organizations Domain**: Multi-tenant organization management
- **Workflows Domain**: DAG-based data processing pipelines
- **Content Domain**: Artifact and content management

### 3. Clean Architecture Principles

- **Dependency Rule**: Dependencies point inward toward the domain
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions, not concretions
- **Single Responsibility**: Each component has one reason to change

## System Architecture

### High-Level Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[React Web App]
        CLI[CLI Tools]
        API_CLIENT[API Clients]
    end

    subgraph "API Gateway"
        NGINX[Nginx/Load Balancer]
    end

    subgraph "Application Layer"
        API[REST API Server]
        WORKER[Background Workers]
        WEBHOOK[Webhook Service]
    end

    subgraph "Data Layer"
        PG[(PostgreSQL)]
        REDIS[(Redis Cache)]
        S3[Object Storage]
        VECTOR[(Vector DB)]
    end

    subgraph "External Services"
        AUTH_PROVIDER[OAuth Providers]
        AI_SERVICE[AI/ML Services]
        EMAIL[Email Service]
    end

    WEB --> NGINX
    CLI --> NGINX
    API_CLIENT --> NGINX

    NGINX --> API
    API --> PG
    API --> REDIS
    API --> S3
    API --> VECTOR

    WORKER --> PG
    WORKER --> S3
    WORKER --> AI_SERVICE

    API --> AUTH_PROVIDER
    API --> EMAIL
    WEBHOOK --> API
```

### Component Architecture

```mermaid
graph LR
    subgraph "API Server"
        HANDLER[HTTP Handlers]
        MIDDLEWARE[Middleware]
        ROUTER[Router]
    end

    subgraph "Domain Layer"
        USE_CASE[Use Cases]
        ENTITY[Entities]
        PORT[Port Interfaces]
    end

    subgraph "Infrastructure"
        REPO[Repositories]
        SERVICE[External Services]
        ADAPTER[Adapters]
    end

    ROUTER --> HANDLER
    HANDLER --> MIDDLEWARE
    MIDDLEWARE --> USE_CASE
    USE_CASE --> PORT
    PORT --> REPO
    PORT --> SERVICE
    REPO --> ADAPTER
```

## Project Structure

### Directory Layout

```
archesai/
├── api/                          # API Specifications
│   ├── openapi.yaml             # Main OpenAPI spec
│   ├── paths/                   # Path definitions
│   │   ├── auth.yaml
│   │   ├── organizations.yaml
│   │   └── workflows.yaml
│   └── components/              # Reusable components
│       ├── schemas/
│       ├── parameters/
│       └── responses/
│
├── cmd/                         # Application Entry Points
│   ├── archesai/               # Main server
│   │   └── main.go
│   ├── worker/                 # Background worker
│   │   └── main.go
│   └── cli/                    # CLI tool
│       └── main.go
│
├── internal/                    # Private Application Code
│   ├── app/                    # Application Assembly
│   │   ├── deps.go            # Dependency injection
│   │   ├── routes.go          # Route registration
│   │   ├── middleware.go      # Global middleware
│   │   └── server.go          # Server configuration
│   │
│   ├── auth/                   # Auth Domain
│   │   ├── domain/            # Core business logic
│   │   │   ├── entities.go   # Domain entities
│   │   │   ├── ports.go      # Repository interfaces
│   │   │   ├── usecase.go    # Business use cases
│   │   │   └── types.gen.go  # Generated types
│   │   ├── adapters/          # Infrastructure adapters
│   │   │   ├── postgres/     # PostgreSQL implementation
│   │   │   └── adapters.gen.go
│   │   ├── handlers/          # HTTP handlers
│   │   │   └── http/
│   │   └── generated/         # Generated code
│   │
│   ├── organizations/          # Organizations Domain
│   ├── workflows/              # Workflows Domain
│   ├── content/                # Content Domain
│   │
│   ├── database/               # Database Layer
│   │   ├── postgresql/        # Generated SQLC code
│   │   ├── queries/           # SQL queries
│   │   └── migrations/        # Database migrations
│   │
│   ├── config/                 # Configuration
│   │   ├── config.go
│   │   └── defaults.gen.go
│   │
│   └── middleware/             # Shared Middleware
│       ├── auth.go
│       ├── cors.go
│       └── ratelimit.go
│
├── pkg/                        # Public Packages
│   ├── errors/                # Error handling
│   ├── logger/                # Logging utilities
│   └── validator/             # Validation helpers
│
├── web/                        # Frontend Applications
│   ├── platform/              # Main React app
│   ├── client/                # TypeScript API client
│   └── ui/                    # Shared UI components
│
└── deployments/                # Deployment Configurations
    ├── docker/                # Dockerfiles
    ├── kubernetes/            # K8s manifests
    └── terraform/             # Infrastructure as Code
```

## Domain Architecture

### Domain Structure

Each domain follows this structure:

```
domain/
├── domain/                     # Core Business Logic
│   ├── entities.go            # Domain entities
│   ├── ports.go              # Repository interfaces
│   ├── usecase.go            # Business use cases
│   ├── errors.go             # Domain-specific errors
│   └── types.gen.go          # Generated types from OpenAPI
│
├── adapters/                   # Adapters Layer
│   ├── postgres/              # Database implementation
│   │   └── repository.go     # Repository implementation
│   ├── http/                  # HTTP adapters
│   │   └── dto.go           # Data transfer objects
│   └── adapters.gen.go       # Generated type converters
│
├── handlers/                   # Presentation Layer
│   └── http/
│       ├── handler.go        # HTTP request handlers
│       └── routes.go         # Route definitions
│
└── generated/                  # Generated Code
    └── api/
        ├── server.gen.go     # Generated server interface
        └── types.gen.go      # Generated request/response types
```

### Domain Interactions

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant UseCase
    participant Repository
    participant Database

    Client->>Handler: HTTP Request
    Handler->>Handler: Validate Input
    Handler->>UseCase: Execute Business Logic
    UseCase->>Repository: Query/Command
    Repository->>Database: SQL Query
    Database-->>Repository: Result
    Repository-->>UseCase: Domain Entity
    UseCase-->>Handler: Response
    Handler-->>Client: HTTP Response
```

## Data Flow Architecture

### Request Processing Pipeline

1. **Client Request** → Load Balancer
2. **Load Balancer** → API Server
3. **API Server** → Router
4. **Router** → Middleware Chain
5. **Middleware** → Handler
6. **Handler** → Use Case
7. **Use Case** → Repository/Service
8. **Repository** → Database
9. **Response** flows back through the same layers

### Workflow Execution Architecture

```mermaid
graph TD
    subgraph "Workflow Engine"
        SCHEDULER[Scheduler]
        EXECUTOR[DAG Executor]
        TASK_QUEUE[Task Queue]
    end

    subgraph "Workers"
        WORKER1[Worker 1]
        WORKER2[Worker 2]
        WORKER3[Worker N]
    end

    subgraph "Tools"
        TEXT_TOOL[Text Processing]
        IMAGE_TOOL[Image Processing]
        EMBED_TOOL[Embedding Generation]
        CUSTOM_TOOL[Custom Tools]
    end

    API --> SCHEDULER
    SCHEDULER --> TASK_QUEUE
    TASK_QUEUE --> EXECUTOR
    EXECUTOR --> WORKER1
    EXECUTOR --> WORKER2
    EXECUTOR --> WORKER3

    WORKER1 --> TEXT_TOOL
    WORKER2 --> IMAGE_TOOL
    WORKER3 --> EMBED_TOOL
    WORKER3 --> CUSTOM_TOOL
```

## Database Architecture

### Schema Design

```sql
-- Core Tables
users
organizations
organization_members

-- Auth Tables
sessions
accounts
verification_tokens

-- Workflow Tables
workflows
workflow_runs
workflow_nodes
workflow_edges

-- Content Tables
artifacts
artifact_embeddings
labels
artifact_labels

-- Audit Tables
audit_logs
activity_logs
```

### Database Patterns

1. **UUID Primary Keys**: All tables use UUIDs for global uniqueness
2. **Soft Deletes**: Critical data uses soft deletes with `deleted_at`
3. **Audit Trails**: All changes logged to audit tables
4. **Multi-tenancy**: Row-level security with `organization_id`
5. **Vector Storage**: pgvector for embedding storage and similarity search

### Query Patterns

- **SQLC**: Type-safe SQL queries
- **Transactions**: ACID compliance for critical operations
- **Prepared Statements**: Protection against SQL injection
- **Connection Pooling**: Efficient resource usage

## Security Architecture

### Authentication & Authorization

```mermaid
graph TD
    subgraph "Authentication Flow"
        LOGIN[Login Request]
        VALIDATE[Validate Credentials]
        JWT[Generate JWT]
        REFRESH[Refresh Token]
    end

    subgraph "Authorization Flow"
        TOKEN[Validate Token]
        CLAIMS[Extract Claims]
        PERMS[Check Permissions]
        ACCESS[Grant/Deny Access]
    end

    LOGIN --> VALIDATE
    VALIDATE --> JWT
    JWT --> REFRESH

    TOKEN --> CLAIMS
    CLAIMS --> PERMS
    PERMS --> ACCESS
```

### Security Layers

1. **Network Security**
   - TLS/SSL encryption
   - Rate limiting
   - DDoS protection
   - WAF rules

2. **Application Security**
   - JWT authentication
   - RBAC authorization
   - Input validation
   - Output encoding
   - CSRF protection

3. **Data Security**
   - Encryption at rest
   - Encryption in transit
   - Key management
   - Data masking
   - Audit logging

## Scalability Architecture

### Horizontal Scaling

```mermaid
graph LR
    subgraph "Load Balancer"
        LB[HAProxy/Nginx]
    end

    subgraph "API Servers"
        API1[Server 1]
        API2[Server 2]
        API3[Server N]
    end

    subgraph "Workers"
        W1[Worker 1]
        W2[Worker 2]
        W3[Worker N]
    end

    subgraph "Data Layer"
        PG_PRIMARY[(PG Primary)]
        PG_REPLICA1[(PG Replica 1)]
        PG_REPLICA2[(PG Replica 2)]
        REDIS_CLUSTER[(Redis Cluster)]
    end

    LB --> API1
    LB --> API2
    LB --> API3

    API1 --> PG_PRIMARY
    API2 --> PG_REPLICA1
    API3 --> PG_REPLICA2

    W1 --> PG_PRIMARY
    W2 --> PG_PRIMARY
    W3 --> PG_PRIMARY
```

### Caching Strategy

1. **Application Cache** (Redis)
   - Session data
   - Frequently accessed data
   - Rate limiting counters

2. **Database Cache**
   - Query result caching
   - Prepared statement caching

3. **CDN** (Static Assets)
   - Frontend assets
   - Public artifacts
   - API documentation

### Queue Architecture

```mermaid
graph LR
    subgraph "Producers"
        API[API Server]
        WEBHOOK[Webhooks]
        SCHEDULER[Scheduler]
    end

    subgraph "Message Queue"
        QUEUE[(Redis/RabbitMQ)]
    end

    subgraph "Consumers"
        WORKER1[Worker Pool 1]
        WORKER2[Worker Pool 2]
        WORKER3[Worker Pool 3]
    end

    API --> QUEUE
    WEBHOOK --> QUEUE
    SCHEDULER --> QUEUE

    QUEUE --> WORKER1
    QUEUE --> WORKER2
    QUEUE --> WORKER3
```

## Development Architecture

### Code Generation Pipeline

```mermaid
graph TD
    OPENAPI[OpenAPI Spec]
    SQL[SQL Queries]
    CONFIG[Config Schema]

    OPENAPI --> OAPI_GEN[oapi-codegen]
    SQL --> SQLC[sqlc]
    CONFIG --> DEFAULTS_GEN[defaults-gen]

    OAPI_GEN --> TYPES[Go Types]
    OAPI_GEN --> SERVER[Server Interfaces]
    OAPI_GEN --> CLIENT[TS Client]

    SQLC --> QUERIES[Query Functions]
    SQLC --> MODELS[DB Models]

    DEFAULTS_GEN --> CONFIG_GO[Config Structs]
```

### Testing Architecture

1. **Unit Tests**
   - Domain logic testing
   - Pure functions
   - Mocked dependencies

2. **Integration Tests**
   - API endpoint testing
   - Database operations
   - External service integration

3. **E2E Tests**
   - Complete user flows
   - Multi-service interactions
   - Performance testing

## Deployment Architecture

### Container Architecture

```yaml
services:
  api:
    image: archesai/api
    replicas: 3
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi

  worker:
    image: archesai/worker
    replicas: 5
    resources:
      limits:
        cpu: 2000m
        memory: 2Gi

  postgres:
    image: postgres:15-pgvector
    replicas: 1
    storage: 100Gi

  redis:
    image: redis:7-alpine
    replicas: 3
    mode: cluster
```

### Kubernetes Architecture

```mermaid
graph TD
    subgraph "Ingress"
        INGRESS[Nginx Ingress]
    end

    subgraph "Services"
        API_SVC[API Service]
        WORKER_SVC[Worker Service]
    end

    subgraph "Deployments"
        API_DEPLOY[API Deployment]
        WORKER_DEPLOY[Worker Deployment]
    end

    subgraph "StatefulSets"
        PG_STS[PostgreSQL]
        REDIS_STS[Redis]
    end

    subgraph "Jobs"
        MIGRATE_JOB[Migration Job]
        BACKUP_JOB[Backup CronJob]
    end

    INGRESS --> API_SVC
    API_SVC --> API_DEPLOY
    WORKER_SVC --> WORKER_DEPLOY

    API_DEPLOY --> PG_STS
    API_DEPLOY --> REDIS_STS
    WORKER_DEPLOY --> PG_STS
```

## Monitoring & Observability

### Metrics Collection

```mermaid
graph LR
    subgraph "Application"
        API[API Metrics]
        WORKER[Worker Metrics]
        DB[Database Metrics]
    end

    subgraph "Collection"
        PROM[Prometheus]
        LOKI[Loki]
        TEMPO[Tempo]
    end

    subgraph "Visualization"
        GRAFANA[Grafana]
        ALERTS[Alert Manager]
    end

    API --> PROM
    WORKER --> PROM
    DB --> PROM

    API --> LOKI
    WORKER --> LOKI

    API --> TEMPO

    PROM --> GRAFANA
    LOKI --> GRAFANA
    TEMPO --> GRAFANA
    PROM --> ALERTS
```

### Key Metrics

1. **Application Metrics**
   - Request rate
   - Response time
   - Error rate
   - Throughput

2. **Infrastructure Metrics**
   - CPU usage
   - Memory usage
   - Disk I/O
   - Network traffic

3. **Business Metrics**
   - User registrations
   - Workflow executions
   - Artifact processing
   - API usage

## Performance Considerations

### Optimization Strategies

1. **Database Optimization**
   - Proper indexing
   - Query optimization
   - Connection pooling
   - Read replicas

2. **Caching**
   - Redis for hot data
   - HTTP caching headers
   - Query result caching

3. **Async Processing**
   - Background jobs
   - Message queues
   - Event-driven architecture

4. **Resource Management**
   - Connection pooling
   - Goroutine management
   - Memory profiling
   - CPU profiling

## Technology Stack

### Backend

- **Language**: Go 1.21+
- **Framework**: Echo v4
- **Database**: PostgreSQL 15+ with pgvector
- **Cache**: Redis 7+
- **Queue**: Redis/RabbitMQ
- **Authentication**: JWT

### Frontend

- **Framework**: React 18+
- **Language**: TypeScript 5+
- **Routing**: TanStack Router
- **State**: Zustand/TanStack Query
- **Build**: Vite
- **UI**: Tailwind CSS

### Infrastructure

- **Container**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus/Grafana
- **Logging**: Loki
- **Tracing**: Tempo/Jaeger

### Code Generation

- **API**: OpenAPI 3.0 + oapi-codegen
- **Database**: SQLC
- **Configuration**: Custom generators

## Best Practices

### Architectural Principles

1. **SOLID Principles**
   - Single Responsibility
   - Open/Closed
   - Liskov Substitution
   - Interface Segregation
   - Dependency Inversion

2. **DRY** (Don't Repeat Yourself)
   - Code generation for repetitive tasks
   - Shared libraries for common functionality

3. **KISS** (Keep It Simple, Stupid)
   - Simple solutions over complex ones
   - Clear and readable code

4. **YAGNI** (You Aren't Gonna Need It)
   - Build only what's needed
   - Avoid premature optimization

### Code Organization

1. **Domain Isolation**
   - Each domain is self-contained
   - Clear boundaries between domains
   - Minimal cross-domain dependencies

2. **Layered Architecture**
   - Clear separation of concerns
   - Dependencies flow inward
   - Business logic in the core

3. **Interface-Based Design**
   - Program to interfaces
   - Dependency injection
   - Easy testing and mocking

## Future Architecture Considerations

### Planned Enhancements

1. **Microservices Migration**
   - Service mesh (Istio)
   - gRPC communication
   - Independent deployments

2. **Event Sourcing**
   - Event store
   - CQRS pattern
   - Event replay capability

3. **GraphQL API**
   - Federation
   - Subscriptions
   - Schema stitching

4. **Multi-Region Deployment**
   - Global load balancing
   - Data replication
   - Edge computing

## References

- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://martinfowler.com/books/ddd.html)
- [Twelve-Factor App](https://12factor.net/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
