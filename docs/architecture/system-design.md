# Architecture Documentation

## Overview

Arches implements **Hexagonal Architecture** (Ports & Adapters) with **Domain-Driven Design**
principles, ensuring separation of concerns, testability, and business logic independence from
infrastructure.

## Core Principles

### Hexagonal Architecture

- **Core Domain**: Business logic and rules
- **Ports**: Interfaces defining external interactions
- **Adapters**: Implementations connecting to external systems

### Dependency Rule

Dependencies flow inward toward the domain core:

```text
External World → Adapters → Ports → Domain Core
```

### Domain Isolation

Each bounded context (auth, organizations, workflows, content) operates independently with:

- Own entities and business rules
- Dedicated database tables
- Separate API endpoints
- No cross-domain imports

## System Architecture

### High-Level Components

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[React SPA]
        CLI[CLI Tools]
        SDK[TypeScript SDK]
    end

    subgraph "API Layer"
        GATEWAY[API Gateway/Load Balancer]
        API[REST API Server]
        DOCS[OpenAPI Docs]
    end

    subgraph "Business Layer"
        AUTH[Auth Domain]
        ORG[Organizations Domain]
        WORK[Workflows Domain]
        CONT[Content Domain]
    end

    subgraph "Data Layer"
        PG[(PostgreSQL + pgvector)]
        REDIS[(Redis Cache)]
        S3[Object Storage]
    end

    subgraph "External Services"
        OAUTH[OAuth Providers]
        AI[AI/ML Services]
        EMAIL[Email Service]
    end

    WEB --> GATEWAY
    CLI --> GATEWAY
    SDK --> GATEWAY

    GATEWAY --> API
    API --> DOCS

    API --> AUTH
    API --> ORG
    API --> WORK
    API --> CONT

    AUTH --> PG
    AUTH --> REDIS
    ORG --> PG
    WORK --> PG
    WORK --> S3
    CONT --> PG
    CONT --> S3

    AUTH --> OAUTH
    AUTH --> EMAIL
    WORK --> AI
```

### Request Flow

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Middleware
    participant Service
    participant Repository
    participant Database

    Client->>Handler: HTTP Request
    Handler->>Middleware: Authentication
    Middleware->>Handler: User Context
    Handler->>Service: Business Operation
    Service->>Repository: Data Operation
    Repository->>Database: SQL Query
    Database->>Repository: Result
    Repository->>Service: Domain Entity
    Service->>Handler: Response
    Handler->>Client: HTTP Response
```

## Domain Architecture

### Flat Package Structure

Each domain follows a flat package structure for simplicity:

```text
internal/auth/
├── auth.go                    # Package documentation, constants, errors
├── service.go                 # Business logic implementation
├── handler.go                 # HTTP request handlers
├── middleware_http.go         # HTTP middleware (auth domain only)
├── repository.gen.go          # Generated repository interface
├── postgres.gen.go            # Generated PostgreSQL implementation
├── sqlite.gen.go              # Generated SQLite implementation
├── service.gen.go             # Generated service interface
├── server.gen.go              # Generated HTTP server implementation
├── types.gen.go               # Generated OpenAPI types
├── handler.gen.go                 # Generated API client interface
└── mocks_test.gen.go          # Generated test mocks
```

### Service Layer Pattern

The service layer orchestrates business operations with generated interfaces:

```go
// Generated service interface (service.gen.go)
type Service interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id uuid.UUID) (*User, error)
    UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*User, error)
    DeleteUser(ctx context.Context, id uuid.UUID) error
    ListUsers(ctx context.Context, params ListUsersParams) ([]*User, error)
}

// Manual implementation (service.go)
type ServiceImpl struct {
    repo Repository // Generated repository interface
    log  *slog.Logger
}

// Business operation example
func (s *ServiceImpl) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // 1. Validate business rules
    if err := s.validateUserCreation(req); err != nil {
        return nil, err
    }

    // 2. Execute business logic
    user := &User{
        ID:        uuid.New(),
        Email:     req.Email,
        Name:      req.Name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // 3. Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    user.PasswordHash = string(hash)

    // 4. Persist to repository (generated interface)
    if err := s.repo.CreateUser(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}
```

### Repository Pattern

Repositories abstract data persistence with generated implementations:

```go
// Generated repository interface (repository.gen.go)
type Repository interface {
    CreateUser(ctx context.Context, user *User) error
    GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id uuid.UUID) error
}

// Generated PostgreSQL implementation (postgres.gen.go)
type PostgresRepository struct {
    db *sql.DB
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *User) error {
    // Generated SQL execution
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO users (id, email, name, password_hash) VALUES ($1, $2, $3, $4)`,
        user.ID, user.Email, user.Name, user.PasswordHash)
    return err
}

// Generated SQLite implementation (sqlite.gen.go)
type SQLiteRepository struct {
    db *sql.DB
}

func (r *SQLiteRepository) CreateUser(ctx context.Context, user *User) error {
    // Generated SQL execution with SQLite syntax
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO users (id, email, name, password_hash) VALUES (?, ?, ?, ?)`,
        user.ID, user.Email, user.Name, user.PasswordHash)
    return err
}
```

## Data Flow Architecture

### Code Generation Pipeline

```mermaid
graph LR
    OPENAPI[OpenAPI Spec] --> UNIFIED[Unified Code Generator]
    XCODEGEN[x-codegen Annotations] --> UNIFIED

    UNIFIED --> REPO[Repository Interface]
    UNIFIED --> POSTGRES[PostgreSQL Implementation]
    UNIFIED --> SQLITE[SQLite Implementation]
    UNIFIED --> SERVICE[Service Interface]
    UNIFIED --> SERVER[HTTP Server]
    UNIFIED --> TYPES[Go Types]
    UNIFIED --> MOCKS[Test Mocks]

    SQL[SQL Queries] --> SQLC[SQLC Generator]
    SQLC --> QUERIES[Database Queries]
    SQLC --> MODELS[Database Models]

    SERVICE --> REPO
    REPO --> POSTGRES
    REPO --> SQLITE
    SERVER --> SERVICE
```

### Authentication Flow

```mermaid
sequenceDiagram
    participant User
    participant API
    participant AuthMiddleware
    participant AuthService
    participant Database
    participant JWT

    User->>API: POST /auth/login
    API->>AuthService: Authenticate(email, password)
    AuthService->>Database: GetUserByEmail
    Database->>AuthService: User
    AuthService->>AuthService: VerifyPassword
    AuthService->>JWT: GenerateTokens
    JWT->>AuthService: AccessToken, RefreshToken
    AuthService->>Database: CreateSession
    AuthService->>API: Tokens
    API->>User: 200 OK + Tokens

    User->>API: GET /api/resource + Bearer Token
    API->>AuthMiddleware: ValidateToken
    AuthMiddleware->>JWT: VerifyToken
    JWT->>AuthMiddleware: Claims
    AuthMiddleware->>Database: GetSession
    Database->>AuthMiddleware: Session
    AuthMiddleware->>API: User Context
    API->>User: 200 OK + Resource
```

## Database Architecture

### Schema Design

```sql
-- Multi-tenant foundation
CREATE TABLE organizations (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  plan VARCHAR(50) DEFAULT 'free',
  credits INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User management
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Organization membership
CREATE TABLE members (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users (id),
  organization_id UUID REFERENCES organizations (id),
  role VARCHAR(50) NOT NULL,
  UNIQUE (user_id, organization_id)
);

-- Session management
CREATE TABLE sessions (
  token VARCHAR(255) PRIMARY KEY,
  user_id UUID REFERENCES users (id),
  expires_at TIMESTAMP NOT NULL,
  active_organization_id UUID REFERENCES organizations (id)
);

-- Content storage with vectors
CREATE TABLE artifacts (
  id UUID PRIMARY KEY,
  organization_id UUID REFERENCES organizations (id),
  name VARCHAR(255) NOT NULL,
  content TEXT,
  embedding vector (1536), -- pgvector for similarity search
  metadata JSONB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workflow definitions
CREATE TABLE pipelines (
  id UUID PRIMARY KEY,
  organization_id UUID REFERENCES organizations (id),
  name VARCHAR(255) NOT NULL,
  definition JSONB NOT NULL, -- DAG structure
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workflow executions
CREATE TABLE runs (
  id UUID PRIMARY KEY,
  pipeline_id UUID REFERENCES pipelines (id),
  status VARCHAR(50) NOT NULL,
  progress DECIMAL(5, 2) DEFAULT 0,
  started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP
);
```

### Indexing Strategy

```sql
-- Performance indexes
CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_sessions_token ON sessions (token);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);

CREATE INDEX idx_members_user_org ON members (user_id, organization_id);

CREATE INDEX idx_artifacts_org ON artifacts (organization_id);

CREATE INDEX idx_pipelines_org ON pipelines (organization_id);

CREATE INDEX idx_runs_pipeline ON runs (pipeline_id);

CREATE INDEX idx_runs_status ON runs (status);

-- Vector similarity search
CREATE INDEX idx_artifacts_embedding ON artifacts USING ivfflat (embedding vector_cosine_ops)
WITH
  (lists = 100);
```

## Caching Architecture

### Cache Layers

1. **Application Cache** (Redis)
   - Session data
   - User profiles
   - Organization metadata
   - Temporary computation results

2. **Database Cache** (PostgreSQL)
   - Query result caching
   - Prepared statement caching
   - Connection pooling

3. **CDN Cache** (CloudFlare/CloudFront)
   - Static assets
   - API responses for public data

### Cache Patterns

```go
// Cache-Aside Pattern
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    // Try cache first
    user, err := s.cache.GetUser(ctx, id)
    if err == nil {
        return user, nil
    }

    // Cache miss - get from database
    user, err = s.repo.GetUserByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Update cache for next time
    _ = s.cache.SetUser(ctx, user, 5*time.Minute)

    return user, nil
}

// Write-Through Pattern
func (s *Service) UpdateUser(ctx context.Context, user *User) error {
    // Update database
    if err := s.repo.UpdateUser(ctx, user); err != nil {
        return err
    }

    // Update cache
    _ = s.cache.SetUser(ctx, user, 5*time.Minute)

    // Publish update event
    _ = s.events.PublishUserUpdated(ctx, user)

    return nil
}
```

## Security Architecture

### Defense in Depth

1. **Network Layer**
   - TLS/HTTPS encryption
   - Rate limiting
   - DDoS protection

2. **Application Layer**
   - JWT authentication
   - CORS configuration
   - Input validation
   - SQL injection prevention (via SQLC)

3. **Data Layer**
   - Encryption at rest
   - Row-level security
   - Audit logging

### Authentication & Authorization

```go
// JWT Claims Structure
type Claims struct {
    UserID         uuid.UUID `json:"user_id"`
    Email          string    `json:"email"`
    OrganizationID uuid.UUID `json:"organization_id"`
    Role           string    `json:"role"`
    jwt.StandardClaims
}

// Middleware for protected routes
func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Extract token
            token := extractToken(c.Request())

            // Validate token
            claims, err := validateToken(token, jwtSecret)
            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized)
            }

            // Set user context
            c.Set("user", claims)

            return next(c)
        }
    }
}

// Role-based access control
func RequireRole(roles ...string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            claims := c.Get("user").(*Claims)

            for _, role := range roles {
                if claims.Role == role {
                    return next(c)
                }
            }

            return echo.NewHTTPError(http.StatusForbidden)
        }
    }
}
```

## Scalability Considerations

### Horizontal Scaling

- **Stateless API servers** - Can scale horizontally behind load balancer
- **Database read replicas** - Distribute read load
- **Redis clustering** - Distributed caching
- **Queue-based processing** - Async job processing

### Performance Optimizations

1. **Connection Pooling**

   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

2. **Batch Processing**

   ```go
   // Process in batches to avoid memory issues
   const batchSize = 100
   for i := 0; i < len(items); i += batchSize {
       end := i + batchSize
       if end > len(items) {
           end = len(items)
       }
       processBatch(items[i:end])
   }
   ```

3. **Concurrent Processing**

   ```go
   // Use worker pool pattern
   jobs := make(chan Job, 100)
   results := make(chan Result, 100)

   // Start workers
   for w := 1; w <= numWorkers; w++ {
       go worker(jobs, results)
   }

   // Send jobs
   for _, job := range allJobs {
       jobs <- job
   }
   close(jobs)
   ```

## Monitoring & Observability

### Metrics Collection

<!-- - **Application Metrics** (Prometheus)
  - Request latency
  - Error rates
  - Business metrics

- **Infrastructure Metrics** (Grafana)
  - CPU/Memory usage
  - Database connections
  - Cache hit rates -->

### Logging Strategy

```go
// Structured logging with context
slog.Info("Processing request",
    zap.String("request_id", requestID),
    zap.String("user_id", userID),
    zap.String("action", "create_artifact"),
    zap.Duration("duration", duration),
)
```

### Distributed Tracing

- OpenTelemetry integration
- Request correlation IDs
- Cross-service tracing

## Deployment Architecture

### Hybrid Kustomize + Helm Strategy

Arches uses a **hybrid deployment approach** that combines the best of both Kustomize and Helm:

```text
┌─────────────────────────────────────────────────┐
│            Helm Templates Only                  │
│        kustomization.yaml File                  │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│              Kustomize                          │
│         Components + Base                       │
│    (Plain YAML Kubernetes Resources)           │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│          Final Kubernetes                       │
│            Manifests                            │
└─────────────────────────────────────────────────┘
```

**Benefits:**

- **Plain YAML Resources**: All Kubernetes manifests are valid YAML (no templating)
- **Component-based**: Modular architecture with reusable components
- **Environment Control**: Helm values control which components are enabled
- **GitOps Friendly**: Generated manifests can be committed and tracked
- **Consistent Labels**: Automatic label injection via Kustomize `commonLabels`

### Component Structure

```text
deployments/
├── kustomize/
│   ├── base/
│   │   ├── namespace.yaml          # Base namespace
│   │   └── kustomization.yaml      # Base configuration
│   └── components/
│       ├── api/                    # API service component
│       │   ├── deployment.yaml     # Plain YAML (no templates)
│       │   ├── service.yaml
│       │   └── kustomization.yaml  # commonLabels applied
│       ├── database/               # PostgreSQL component
│       │   ├── statefulset.yaml
│       │   ├── service.yaml
│       │   └── kustomization.yaml
│       └── ...                     # Other components
└── helm-minimal/
    ├── templates/
    │   └── kustomization.yaml      # ONLY file templated by Helm
    ├── values-dev.yaml             # Dev environment config
    └── values-prod.yaml            # Prod environment config
```

### Deployment Process

```bash
# 1. Helm templates the kustomization.yaml file
helm template archesai deployments/helm-minimal \
  -f deployments/helm-minimal/values-prod.yaml \
  --set namespace=production > /tmp/kustomization.yaml

# 2. Kustomize builds final manifests with labels
kustomize build /tmp | kubectl apply -f -
```

### Generated Manifest Example

```yaml
# Generated by Kustomize with labels automatically injected
apiVersion: apps/v1
kind: Deployment
metadata:
  name: archesai-api
  labels:
    app.kubernetes.io/name: archesai # From commonLabels
    app.kubernetes.io/component: api # From commonLabels
spec:
  replicas: 3 # From Helm values
  selector:
    matchLabels:
      app.kubernetes.io/name: archesai # Auto-injected
      app.kubernetes.io/component: api # Auto-injected
  template:
    metadata:
      labels:
        app.kubernetes.io/name: archesai # Auto-injected
        app.kubernetes.io/component: api # Auto-injected
    spec:
      serviceAccountName: archesai
      containers:
        - name: api
          image: archesai/api:v1.0.0 # From Helm image tags
          ports:
            - name: http
              containerPort: 3001
```

## Future Considerations

### Planned Enhancements

1. **GraphQL API** - Alternative query interface
2. **gRPC Services** - Internal service communication
3. **Event Sourcing** - Audit trail and replay capability
4. **CQRS Pattern** - Separate read/write models
5. **Service Mesh** - Istio/Linkerd for microservices
6. **Multi-region** - Geographic distribution

### Technology Evaluations

- **Message Queue**: Kafka vs NATS vs RabbitMQ
- **Search Engine**: Elasticsearch vs Meilisearch
- **Time Series DB**: InfluxDB vs TimescaleDB
- **Graph Database**: Neo4j vs DGraph
