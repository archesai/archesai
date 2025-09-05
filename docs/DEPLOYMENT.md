# Deployment Guide

This guide covers deploying ArchesAI to various environments, from local development to production cloud infrastructure.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development](#local-development)
3. [Docker Deployment](#docker-deployment)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [Cloud Deployments](#cloud-deployments)
6. [Production Configuration](#production-configuration)
7. [Monitoring & Maintenance](#monitoring--maintenance)
8. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Tools

- Docker 24.0+
- Docker Compose 2.20+
- Kubernetes 1.28+ (for K8s deployment)
- Helm 3.12+ (for K8s deployment)
- PostgreSQL client tools
- Make

### Required Services

- PostgreSQL 15+ with pgvector extension
- Redis 7+ (for caching and queues)
- S3-compatible object storage (MinIO for local)
- SMTP server (for email notifications)

## Local Development

### Quick Start

```bash
# Clone repository
git clone https://github.com/archesai/archesai.git
cd archesai

# Set up environment
cp .env.example .env
# Edit .env with your configuration

# Start dependencies with Docker Compose
docker-compose -f docker-compose.dev.yml up -d postgres redis minio

# Run database migrations
make migrate-up

# Start the application
make dev
```

### Docker Compose Development

```yaml
# docker-compose.dev.yml
version: "3.8"

services:
  postgres:
    image: pgvector/pgvector:pg15
    environment:
      POSTGRES_USER: archesai
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: archesai
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data

volumes:
  postgres_data:
  redis_data:
  minio_data:
```

## Docker Deployment

### Building Docker Images

```bash
# Build API server image
docker build -f deployments/docker/Dockerfile.api -t archesai/api:latest .

# Build worker image
docker build -f deployments/docker/Dockerfile.worker -t archesai/worker:latest .

# Build frontend image
docker build -f deployments/docker/Dockerfile.web -t archesai/web:latest .
```

### Production Dockerfile

```dockerfile
# deployments/docker/Dockerfile.api
# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build-archesai

# Runtime stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/archesai /app/
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/archesai"]
CMD ["serve"]
```

### Docker Compose Production

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  api:
    image: archesai/api:latest
    environment:
      ARCHESAI_DATABASE_URL: postgres://archesai:${DB_PASSWORD}@postgres:5432/archesai?sslmode=require
      ARCHESAI_REDIS_URL: redis://redis:6379
      ARCHESAI_JWT_SECRET: ${JWT_SECRET}
      ARCHESAI_SERVER_HOST: 0.0.0.0
      ARCHESAI_SERVER_PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  worker:
    image: archesai/worker:latest
    environment:
      ARCHESAI_DATABASE_URL: postgres://archesai:${DB_PASSWORD}@postgres:5432/archesai?sslmode=require
      ARCHESAI_REDIS_URL: redis://redis:6379
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    deploy:
      replicas: 3

  web:
    image: archesai/web:latest
    ports:
      - "80:80"
    depends_on:
      - api
    restart: unless-stopped

  postgres:
    image: pgvector/pgvector:pg15
    environment:
      POSTGRES_USER: archesai
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: archesai
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

### Deployment Commands

```bash
# Start production stack
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f api

# Scale workers
docker-compose -f docker-compose.prod.yml up -d --scale worker=5

# Stop services
docker-compose -f docker-compose.prod.yml down
```

## Kubernetes Deployment

### Namespace and ConfigMap

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: archesai

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: archesai-config
  namespace: archesai
data:
  SERVER_HOST: "0.0.0.0"
  SERVER_PORT: "8080"
  LOGGING_LEVEL: "info"
  LOGGING_FORMAT: "json"
```

### Secret Management

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: archesai-secrets
  namespace: archesai
type: Opaque
stringData:
  database-url: "postgres://user:pass@postgres:5432/archesai?sslmode=require"
  jwt-secret: "your-super-secret-jwt-key-minimum-32-characters"
  redis-url: "redis://redis:6379"
```

### API Deployment

```yaml
# k8s/deployment-api.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: archesai-api
  namespace: archesai
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
          image: archesai/api:latest
          ports:
            - containerPort: 8080
          env:
            - name: ARCHESAI_DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: archesai-secrets
                  key: database-url
            - name: ARCHESAI_JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: archesai-secrets
                  key: jwt-secret
            - name: ARCHESAI_REDIS_URL
              valueFrom:
                secretKeyRef:
                  name: archesai-secrets
                  key: redis-url
          envFrom:
            - configMapRef:
                name: archesai-config
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "1000m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

### Service Configuration

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: archesai-api
  namespace: archesai
spec:
  selector:
    app: archesai-api
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP

---
apiVersion: v1
kind: Service
metadata:
  name: archesai-api-nodeport
  namespace: archesai
spec:
  type: NodePort
  selector:
    app: archesai-api
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
      nodePort: 30080
```

### Ingress Configuration

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: archesai-ingress
  namespace: archesai
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
    - hosts:
        - api.archesai.com
      secretName: archesai-tls
  rules:
    - host: api.archesai.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: archesai-api
                port:
                  number: 80
```

### Database StatefulSet

```yaml
# k8s/postgres.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: archesai
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: pgvector/pgvector:pg15
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              value: archesai
            - name: POSTGRES_USER
              value: archesai
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: password
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
          resources:
            requests:
              memory: "1Gi"
              cpu: "500m"
            limits:
              memory: "2Gi"
              cpu: "1000m"
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 20Gi
```

### Helm Chart

```yaml
# helm/archesai/values.yaml
replicaCount: 3

image:
  repository: archesai/api
  pullPolicy: IfNotPresent
  tag: "latest"

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: api.archesai.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: archesai-tls
      hosts:
        - api.archesai.com

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 250m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80

postgresql:
  enabled: true
  auth:
    database: archesai
    username: archesai
    password: changeme
  primary:
    persistence:
      size: 20Gi

redis:
  enabled: true
  auth:
    enabled: false
  master:
    persistence:
      size: 8Gi
```

### Helm Deployment

```bash
# Add Helm repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install PostgreSQL
helm install postgres bitnami/postgresql \
  --namespace archesai \
  --set auth.database=archesai \
  --set auth.username=archesai \
  --set auth.password=$DB_PASSWORD \
  --set image.tag=15

# Install Redis
helm install redis bitnami/redis \
  --namespace archesai \
  --set auth.enabled=false

# Install ArchesAI
helm install archesai ./helm/archesai \
  --namespace archesai \
  --values ./helm/archesai/values.yaml
```

## Cloud Deployments

### AWS Deployment

#### ECS with Fargate

```json
// ecs-task-definition.json
{
  "family": "archesai-api",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "1024",
  "memory": "2048",
  "containerDefinitions": [
    {
      "name": "api",
      "image": "archesai/api:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "ARCHESAI_SERVER_PORT",
          "value": "8080"
        }
      ],
      "secrets": [
        {
          "name": "ARCHESAI_DATABASE_URL",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:archesai/db"
        },
        {
          "name": "ARCHESAI_JWT_SECRET",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:archesai/jwt"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/archesai",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "api"
        }
      }
    }
  ]
}
```

#### RDS Configuration

```bash
# Create RDS PostgreSQL instance with pgvector
aws rds create-db-instance \
  --db-instance-identifier archesai-db \
  --db-instance-class db.t3.medium \
  --engine postgres \
  --engine-version 15.4 \
  --master-username archesai \
  --master-user-password $DB_PASSWORD \
  --allocated-storage 100 \
  --vpc-security-group-ids sg-xxxxx \
  --backup-retention-period 7 \
  --multi-az
```

### Google Cloud Platform

#### Cloud Run Deployment

```yaml
# cloudrun-service.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: archesai-api
  annotations:
    run.googleapis.com/launch-stage: GA
spec:
  template:
    metadata:
      annotations:
        run.googleapis.com/cloudsql-instances: project:region:instance
        run.googleapis.com/cpu-throttling: "false"
    spec:
      containerConcurrency: 100
      containers:
        - image: gcr.io/project/archesai-api:latest
          ports:
            - containerPort: 8080
          env:
            - name: ARCHESAI_DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: database-url
                  key: latest
          resources:
            limits:
              cpu: "2"
              memory: "2Gi"
```

#### Deployment Script

```bash
# Deploy to Cloud Run
gcloud run deploy archesai-api \
  --image gcr.io/project/archesai-api:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="ARCHESAI_SERVER_PORT=8080" \
  --set-secrets="ARCHESAI_DATABASE_URL=database-url:latest" \
  --set-secrets="ARCHESAI_JWT_SECRET=jwt-secret:latest" \
  --min-instances=1 \
  --max-instances=10 \
  --cpu=2 \
  --memory=2Gi
```

### Azure Deployment

#### Azure Container Instances

```yaml
# aci-deployment.yaml
apiVersion: 2019-12-01
location: eastus
name: archesai-containers
properties:
  containers:
    - name: api
      properties:
        image: archesai/api:latest
        resources:
          requests:
            cpu: 1.0
            memoryInGb: 1.5
        ports:
          - port: 8080
        environmentVariables:
          - name: ARCHESAI_SERVER_PORT
            value: 8080
          - name: ARCHESAI_DATABASE_URL
            secureValue: postgres://...
  osType: Linux
  ipAddress:
    type: Public
    ports:
      - protocol: tcp
        port: 8080
```

## Production Configuration

### Environment Variables

```bash
# .env.production
# Database
ARCHESAI_DATABASE_URL=postgres://user:pass@host:5432/archesai?sslmode=require
ARCHESAI_DATABASE_POOL_SIZE=25
ARCHESAI_DATABASE_MAX_IDLE_TIME=15m

# Redis
ARCHESAI_REDIS_URL=redis://host:6379
ARCHESAI_REDIS_POOL_SIZE=10

# Server
ARCHESAI_SERVER_HOST=0.0.0.0
ARCHESAI_SERVER_PORT=8080
ARCHESAI_SERVER_READ_TIMEOUT=30s
ARCHESAI_SERVER_WRITE_TIMEOUT=30s
ARCHESAI_SERVER_SHUTDOWN_TIMEOUT=10s

# Authentication
ARCHESAI_JWT_SECRET=your-super-secret-key-minimum-32-characters
ARCHESAI_JWT_ACCESS_TOKEN_DURATION=15m
ARCHESAI_JWT_REFRESH_TOKEN_DURATION=7d

# Storage
ARCHESAI_STORAGE_TYPE=s3
ARCHESAI_STORAGE_S3_BUCKET=archesai-artifacts
ARCHESAI_STORAGE_S3_REGION=us-east-1
ARCHESAI_STORAGE_S3_ENDPOINT=https://s3.amazonaws.com

# Email
ARCHESAI_EMAIL_SMTP_HOST=smtp.sendgrid.net
ARCHESAI_EMAIL_SMTP_PORT=587
ARCHESAI_EMAIL_SMTP_USER=apikey
ARCHESAI_EMAIL_SMTP_PASSWORD=SG.xxxxx
ARCHESAI_EMAIL_FROM=noreply@archesai.com

# Logging
ARCHESAI_LOGGING_LEVEL=info
ARCHESAI_LOGGING_FORMAT=json

# Monitoring
ARCHESAI_METRICS_ENABLED=true
ARCHESAI_METRICS_PORT=9090
ARCHESAI_TRACING_ENABLED=true
ARCHESAI_TRACING_ENDPOINT=http://tempo:14268/api/traces
```

### SSL/TLS Configuration

```nginx
# nginx.conf
server {
    listen 80;
    server_name api.archesai.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.archesai.com;

    ssl_certificate /etc/letsencrypt/live/api.archesai.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.archesai.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    location / {
        proxy_pass http://archesai-api:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Database Migrations

```bash
# Run migrations before deployment
docker run --rm \
  -e ARCHESAI_DATABASE_URL=$DATABASE_URL \
  archesai/api:latest \
  migrate up

# Rollback if needed
docker run --rm \
  -e ARCHESAI_DATABASE_URL=$DATABASE_URL \
  archesai/api:latest \
  migrate down 1
```

## Monitoring & Maintenance

### Health Checks

```go
// Health check endpoints
GET /health     # Basic health check
GET /ready      # Readiness check (includes DB connection)
GET /metrics    # Prometheus metrics
```

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "archesai-api"
    static_configs:
      - targets: ["archesai-api:9090"]

  - job_name: "postgres"
    static_configs:
      - targets: ["postgres-exporter:9187"]

  - job_name: "redis"
    static_configs:
      - targets: ["redis-exporter:9121"]
```

### Grafana Dashboards

```json
// grafana-dashboard.json
{
  "dashboard": {
    "title": "ArchesAI Monitoring",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])"
          }
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])"
          }
        ]
      },
      {
        "title": "Response Time",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, http_request_duration_seconds_bucket)"
          }
        ]
      }
    ]
  }
}
```

### Backup Strategy

```bash
# PostgreSQL Backup
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"

# Database backup
pg_dump $DATABASE_URL > $BACKUP_DIR/archesai_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/archesai_$DATE.sql

# Upload to S3
aws s3 cp $BACKUP_DIR/archesai_$DATE.sql.gz s3://archesai-backups/

# Clean old backups (keep 30 days)
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
```

### Log Aggregation

```yaml
# fluentd.conf
<source>
@type forward
port 24224
</source>

<filter archesai.**>
@type parser
key_name log
format json
</filter>

<match archesai.**>
@type elasticsearch
host elasticsearch
port 9200
logstash_format true
logstash_prefix archesai
</match>
```

## Troubleshooting

### Common Issues

#### Database Connection Issues

```bash
# Test database connection
psql $DATABASE_URL -c "SELECT 1"

# Check connection pool
curl http://localhost:8080/metrics | grep db_connections

# Reset connections
docker-compose restart api
```

#### High Memory Usage

```bash
# Check memory usage
docker stats archesai-api

# Analyze memory profile
curl http://localhost:8080/debug/pprof/heap > heap.pprof
go tool pprof heap.pprof

# Adjust memory limits
docker update --memory=2g archesai-api
```

#### Slow Queries

```sql
-- Find slow queries
SELECT query, calls, mean_exec_time
FROM pg_stat_statements
WHERE mean_exec_time > 1000
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Analyze query plan
EXPLAIN ANALYZE SELECT ...;
```

### Debugging

```bash
# Enable debug logging
export ARCHESAI_LOGGING_LEVEL=debug

# View container logs
docker logs -f archesai-api --tail=100

# Access container shell
docker exec -it archesai-api /bin/sh

# Check environment variables
docker exec archesai-api env | grep ARCHESAI

# Test API endpoint
curl -X GET http://localhost:8080/health
```

### Performance Tuning

```bash
# Database tuning
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET random_page_cost = 1.1;

# Redis tuning
CONFIG SET maxmemory 2gb
CONFIG SET maxmemory-policy allkeys-lru
CONFIG SET tcp-keepalive 60
```

## Security Hardening

### Production Checklist

- [ ] Enable SSL/TLS for all connections
- [ ] Use strong passwords and rotate regularly
- [ ] Enable audit logging
- [ ] Configure firewall rules
- [ ] Disable unnecessary services
- [ ] Keep dependencies updated
- [ ] Implement rate limiting
- [ ] Enable CORS properly
- [ ] Use secrets management service
- [ ] Regular security scans
- [ ] Backup encryption
- [ ] Network segmentation
- [ ] Intrusion detection system
- [ ] DDoS protection
- [ ] Container image scanning

### Security Headers

```go
// middleware/security.go
func SecurityHeaders() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            c.Response().Header().Set("X-Content-Type-Options", "nosniff")
            c.Response().Header().Set("X-Frame-Options", "DENY")
            c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
            c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
            return next(c)
        }
    }
}
```

## Disaster Recovery

### Backup and Restore

```bash
# Full backup script
#!/bin/bash
# backup.sh

# Database backup
pg_dumpall -h $DB_HOST -U $DB_USER > backup.sql

# Application data backup
tar -czf artifacts.tar.gz /data/artifacts

# Configuration backup
tar -czf config.tar.gz /etc/archesai

# Upload to disaster recovery site
rsync -avz backup.sql artifacts.tar.gz config.tar.gz dr-site:/backups/
```

### Restore Procedure

```bash
# Restore database
psql -h $DB_HOST -U $DB_USER < backup.sql

# Restore artifacts
tar -xzf artifacts.tar.gz -C /data

# Restore configuration
tar -xzf config.tar.gz -C /etc

# Restart services
docker-compose up -d
```

## Maintenance Windows

### Rolling Updates

```bash
# Kubernetes rolling update
kubectl set image deployment/archesai-api \
  api=archesai/api:v2.0.0 \
  --namespace=archesai

# Monitor rollout
kubectl rollout status deployment/archesai-api \
  --namespace=archesai

# Rollback if needed
kubectl rollout undo deployment/archesai-api \
  --namespace=archesai
```

### Zero-Downtime Deployment

```bash
# Blue-Green deployment
# 1. Deploy to green environment
docker-compose -f docker-compose.green.yml up -d

# 2. Run smoke tests
./scripts/smoke-tests.sh green

# 3. Switch traffic
./scripts/switch-traffic.sh green

# 4. Monitor for issues
./scripts/monitor.sh --duration=5m

# 5. Remove blue environment
docker-compose -f docker-compose.blue.yml down
```

## Support

For deployment support:

- Documentation: https://docs.archesai.com/deployment
- Issues: https://github.com/archesai/archesai/issues
- Email: devops@archesai.com
