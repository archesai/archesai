# Deployment Documentation

This section provides comprehensive deployment strategies, infrastructure guidance, and production
configuration for Arches across various environments and scales.

## Deployment Philosophy

Arches is designed for **cloud-native deployment** with support for containerization,
orchestration, and horizontal scaling. The platform can be deployed from a single Docker container
for development to a full Kubernetes cluster for enterprise production.

## Documentation Structure

- **[Deployment Overview](overview.md)** - Deployment strategies and options (this document)
- **[Docker](docker.md)** - Container-based deployment with Docker Compose
- **[Kubernetes](kubernetes.md)** - Production orchestration with Helm charts
- **[Production](production.md)** - Production deployment best practices

## Deployment Options

### 1. Local Development

**Best for**: Individual developers, testing, and development

```bash
# Using Make commands
make dev        # Start all services locally
make run-server # Run API server only
make run-web    # Run frontend only

# Using Docker Compose
docker-compose -f deployments/docker/docker-compose.dev.yml up
```

**Features:**

- Hot reload for rapid development
- Local PostgreSQL and Redis
- Mock external services
- Debug-friendly configuration

### 2. Docker Deployment

**Best for**: Small teams, staging environments, simple production setups

```yaml
# docker-compose.yml
version: "3.8"
services:
  api:
    image: archesai/api:latest
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://...
      REDIS_URL: redis://...
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
```

**Features:**

- Single-host deployment
- Built-in service discovery
- Volume persistence
- Easy backup and restore

### 3. Kubernetes Deployment

**Best for**: Production workloads, high availability, auto-scaling

```yaml
# kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: archesai-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: archesai-api
  template:
    metadata:
      labels:
        app: archesai-api
    spec:
      containers:
        - name: api
          image: archesai/api:v1.0.0
          ports:
            - containerPort: 8080
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
```

**Features:**

- Horizontal pod autoscaling
- Rolling updates with zero downtime
- Service mesh integration
- Multi-region deployment

### 4. Serverless Deployment

**Best for**: Variable workloads, cost optimization, event-driven processing

```yaml
# serverless.yml for AWS Lambda
service: archesai-api

provider:
  name: aws
  runtime: go1.x
  region: us-east-1

functions:
  api:
    handler: bin/api
    events:
      - httpApi:
          path: /{proxy+}
          method: ANY
```

**Features:**

- Pay-per-request pricing
- Automatic scaling
- No infrastructure management
- Integration with cloud services

## Infrastructure Components

### Core Services

#### **API Server**

- **Technology**: Go with Echo framework
- **Port**: 8080 (configurable)
- **Scaling**: Horizontal, stateless
- **Health checks**: `/health/live` and `/health/ready`

#### **PostgreSQL Database**

- **Version**: 15+ with pgvector extension
- **Connection pooling**: PgBouncer recommended
- **Replication**: Primary-replica for HA
- **Backup**: Daily automated backups

#### **Redis Cache**

- **Version**: 7+
- **Usage**: Sessions, cache, rate limiting
- **Persistence**: AOF for durability
- **Clustering**: Redis Cluster for HA

#### **Storage Layer**

- **Local**: Development only
- **S3/MinIO**: Production file storage
- **CDN**: CloudFront/Cloudflare for assets

### Supporting Services

#### **Load Balancer**

```nginx
upstream archesai_api {
    least_conn;
    server api1:8080 weight=1;
    server api2:8080 weight=1;
    server api3:8080 weight=1;
}

server {
    listen 443 ssl http2;
    server_name api.archesai.com;

    location / {
        proxy_pass http://archesai_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

#### **Monitoring Stack**

- **Metrics**: Prometheus + Grafana
- **Logging**: Loki + Promtail
- **Tracing**: Jaeger/Tempo
- **Alerting**: AlertManager

## Deployment Configurations

### Environment Variables

```bash
# Core Configuration
API_PORT=8080
API_HOST=0.0.0.0
ENV=production

# Database
DATABASE_URL=postgres://user:pass@host:5432/archesai
DATABASE_MAX_CONNECTIONS=100
DATABASE_CONNECTION_TIMEOUT=30s

# Redis
REDIS_URL=redis://host:6379/0
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=10

# Security
JWT_SECRET=your-secret-key
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d

# Storage
STORAGE_TYPE=s3
AWS_BUCKET=archesai-storage
AWS_REGION=us-east-1

# Monitoring
METRICS_ENABLED=true
METRICS_PORT=9090
TRACING_ENABLED=true
TRACING_ENDPOINT=http://jaeger:14268
```

### Resource Requirements

| Component      | Development         | Staging             | Production       |
| -------------- | ------------------- | ------------------- | ---------------- |
| **API Server** | 256MB RAM, 0.25 CPU | 512MB RAM, 0.5 CPU  | 1GB RAM, 1 CPU   |
| **PostgreSQL** | 512MB RAM, 0.5 CPU  | 2GB RAM, 1 CPU      | 8GB RAM, 4 CPU   |
| **Redis**      | 128MB RAM, 0.1 CPU  | 256MB RAM, 0.25 CPU | 1GB RAM, 0.5 CPU |
| **Storage**    | 10GB                | 100GB               | 1TB+             |

## Deployment Process

### 1. Pre-Deployment Checklist

- [ ] Code review completed
- [ ] Tests passing (unit, integration, e2e)
- [ ] Security scan completed
- [ ] Documentation updated
- [ ] Database migrations prepared
- [ ] Environment variables configured
- [ ] Monitoring alerts configured
- [ ] Rollback plan prepared

### 2. Deployment Steps

```bash
# 1. Build and tag Docker image
docker build -t archesai/api:v1.0.0 .
docker push archesai/api:v1.0.0

# 2. Run database migrations
make db-migrate-up

# 3. Deploy to staging
kubectl apply -f deployments/kubernetes/staging/

# 4. Run smoke tests
make test-staging

# 5. Deploy to production
kubectl apply -f deployments/kubernetes/production/

# 6. Verify deployment
kubectl rollout status deployment/archesai-api
```

### 3. Post-Deployment Verification

- [ ] Health checks passing
- [ ] Metrics being collected
- [ ] Logs aggregating properly
- [ ] API endpoints responding
- [ ] Database connections stable
- [ ] Cache hit rates normal
- [ ] No error spike in monitoring

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
name: Deploy to Production

on:
  push:
    tags:
      - "v*"

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5

      - name: Build and push Docker image
        run: |
          docker build -t archesai/api:${{ github.ref_name }} .
          docker push archesai/api:${{ github.ref_name }}

      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/archesai-api \
            api=archesai/api:${{ github.ref_name }}

      - name: Wait for rollout
        run: kubectl rollout status deployment/archesai-api
```

## Scaling Strategies

### Horizontal Scaling

```yaml
# Horizontal Pod Autoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: archesai-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: archesai-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

### Database Scaling

- **Read replicas**: Distribute read queries
- **Connection pooling**: PgBouncer for connection management
- **Partitioning**: Time-based partitioning for large tables
- **Caching**: Redis for frequently accessed data

## Disaster Recovery

### Backup Strategy

```bash
# Automated daily backups
0 2 * * * pg_dump $DATABASE_URL | gzip > backup_$(date +%Y%m%d).sql.gz

# Upload to S3
aws s3 cp backup_*.sql.gz s3://archesai-backups/
```

### Recovery Procedures

1. **Database recovery**: Restore from latest backup
2. **File storage recovery**: Sync from S3 backup
3. **Configuration recovery**: Restore from git repository
4. **State recovery**: Rebuild from event log

## Cost Optimization

### Recommendations

- Use spot instances for non-critical workloads
- Implement auto-scaling to match demand
- Use reserved instances for baseline capacity
- Optimize container images (multi-stage builds)
- Enable CDN for static assets
- Implement caching aggressively
- Use managed services where appropriate

## Quick Start Guides

- **[Local Development](docker.md)** - Get started in 5 minutes
- **[Production Deployment](production.md)** - Complete production guide
- **[Kubernetes Setup](kubernetes.md)** - Deploy with Helm
- **[Monitoring Setup](../monitoring/overview.md)** - Observability configuration

## Getting Help

For deployment assistance:

- **Documentation**: Full guides in this section
- **Issues**: [GitHub Issues](https://github.com/archesai/archesai/issues)
- **Support**: <support@archesai.com>
- **Community**: [Discord Server](https://discord.gg/archesai)
