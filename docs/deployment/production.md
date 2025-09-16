# Production Deployment Guide

This guide covers best practices, security hardening, and optimization for deploying Arches in
production environments.

## Pre-Deployment Checklist

### Infrastructure Requirements

- [ ] **Compute**: Minimum 3 nodes with 4 vCPUs, 8GB RAM each
- [ ] **Database**: PostgreSQL 14+ with pgvector extension
- [ ] **Cache**: Redis 7+ with persistence enabled
- [ ] **Storage**: 100GB+ SSD for database, 20GB for application logs
- [ ] **Network**: Load balancer with SSL termination
- [ ] **Monitoring**: Prometheus, Grafana, and alerting configured
- [ ] **Backup**: Automated backup strategy implemented

### Security Checklist

- [ ] SSL/TLS certificates configured
- [ ] Secrets management system in place
- [ ] Network policies configured
- [ ] RBAC and authentication configured
- [ ] Security scanning integrated in CI/CD
- [ ] WAF (Web Application Firewall) configured
- [ ] DDoS protection enabled

## Environment Configuration

### Production Environment Variables

```bash
# Application
NODE_ENV=production
API_PORT=8080
API_HOST=0.0.0.0
LOG_LEVEL=info
LOG_FORMAT=json
CORS_ORIGINS=https://app.archesai.com

# Database
DATABASE_URL=postgresql://user:pass@db.archesai.internal:5432/archesai?sslmode=require
DB_POOL_MIN=10
DB_POOL_MAX=50
DB_CONNECTION_TIMEOUT=5000
DB_IDLE_TIMEOUT=10000

# Redis
REDIS_URL=redis://:password@redis.archesai.internal:6379
REDIS_MAX_RETRIES=3
REDIS_RETRY_DELAY=1000

# Security
JWT_SECRET=${JWT_SECRET}
JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
JWT_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
SESSION_SECRET=${SESSION_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}

# Rate Limiting
RATE_LIMIT_WINDOW=60000
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS=false

# OpenAI (if using chat features)
OPENAI_API_KEY=${OPENAI_API_KEY}
OPENAI_MAX_TOKENS=4096
OPENAI_TIMEOUT=30000

# Monitoring
METRICS_ENABLED=true
METRICS_PORT=9090
TRACING_ENABLED=true
TRACING_ENDPOINT=http://jaeger:14268/api/traces
```

## Database Configuration

### PostgreSQL Production Settings

```sql
-- postgresql.conf
max_connections = 200
shared_buffers = 2GB
effective_cache_size = 6GB
maintenance_work_mem = 512MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 10485kB
min_wal_size = 1GB
max_wal_size = 4GB

-- Enable pgvector
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Create indexes
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY idx_organizations_name ON organizations(name);
CREATE INDEX CONCURRENTLY idx_sessions_token ON sessions(token);
CREATE INDEX CONCURRENTLY idx_sessions_expires ON sessions(expires_at);
```

### Connection Pooling with PgBouncer

```ini
[databases]
archesai = host=postgres port=5432 dbname=archesai

[pgbouncer]
listen_port = 6432
listen_addr = *
auth_type = md5
auth_file = /etc/pgbouncer/userlist.txt
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 25
min_pool_size = 10
reserve_pool_size = 5
reserve_pool_timeout = 3
server_lifetime = 3600
server_idle_timeout = 600
```

## Redis Configuration

### Redis Production Settings

```conf
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
appendonly yes
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
requirepass ${REDIS_PASSWORD}
```

### Redis Sentinel for HA

```conf
# sentinel.conf
port 26379
sentinel monitor archesai-master redis-master 6379 2
sentinel auth-pass archesai-master ${REDIS_PASSWORD}
sentinel down-after-milliseconds archesai-master 5000
sentinel parallel-syncs archesai-master 1
sentinel failover-timeout archesai-master 10000
```

## Load Balancing

### Nginx Configuration

```nginx
upstream archesai_backend {
    least_conn;
    server api1.archesai.internal:8080 max_fails=3 fail_timeout=30s;
    server api2.archesai.internal:8080 max_fails=3 fail_timeout=30s;
    server api3.archesai.internal:8080 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

server {
    listen 443 ssl http2;
    server_name api.archesai.com;

    # SSL Configuration
    ssl_certificate /etc/nginx/ssl/archesai.crt;
    ssl_certificate_key /etc/nginx/ssl/archesai.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header Content-Security-Policy "default-src 'self'" always;

    # Rate Limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;

    # Compression
    gzip on;
    gzip_types text/plain application/json application/javascript;
    gzip_min_length 1000;

    location / {
        proxy_pass http://archesai_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts
        proxy_connect_timeout 10s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;

        # Buffering
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }

    location /health {
        access_log off;
        proxy_pass http://archesai_backend;
    }
}
```

## Security Hardening

### Application Security

```go
// main.go security middleware
app.Use(helmet.New())
app.Use(cors.New(cors.Config{
    AllowOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Authorization", "Content-Type"},
    MaxAge:       86400,
}))

app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()
    },
}))
```

### Secrets Management

```yaml
# Using Kubernetes Secrets with Sealed Secrets
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: archesai-secrets
  namespace: production
spec:
  encryptedData:
    jwt-secret: AgA1B2C3D4E5F6...
    db-password: AgF6E5D4C3B2A1...
```

### Network Security

```yaml
# Network Policy with correct labels
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: archesai-network-policy
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: archesai
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: nginx-ingress
      ports:
        - protocol: TCP
          port: 3001
  egress:
    - to:
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: archesai
              app.kubernetes.io/component: postgres
    - to:
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: archesai
              app.kubernetes.io/component: redis
```

## Monitoring and Observability

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: "archesai"
    static_configs:
      - targets: ["api1:9090", "api2:9090", "api3:9090"]
    metrics_path: "/metrics"
```

### Key Metrics to Monitor

```yaml
# Alert Rules
groups:
  - name: archesai
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "High error rate detected"

      - alert: HighLatency
        expr: histogram_quantile(0.95, http_request_duration_seconds_bucket) > 1
        for: 5m
        annotations:
          summary: "95th percentile latency above 1s"

      - alert: HighMemoryUsage
        expr: process_resident_memory_bytes / 1024 / 1024 > 1000
        for: 5m
        annotations:
          summary: "Memory usage above 1GB"
```

### Logging Configuration

```json
{
  "level": "info",
  "format": "json",
  "output": "stdout",
  "fields": {
    "app": "archesai",
    "env": "production",
    "version": "${APP_VERSION}"
  },
  "hooks": [
    {
      "type": "sentry",
      "dsn": "${SENTRY_DSN}",
      "level": "error"
    }
  ]
}
```

## Backup and Disaster Recovery

### Database Backup Strategy

```bash
#!/bin/bash
# backup.sh

# Daily backup
pg_dump $DATABASE_URL | gzip > backup-$(date +%Y%m%d).sql.gz

# Upload to S3
aws s3 cp backup-$(date +%Y%m%d).sql.gz \
  s3://archesai-backups/postgres/

# Cleanup old backups (keep 30 days)
find /backups -name "backup-*.sql.gz" -mtime +30 -delete
```

### Disaster Recovery Plan

1. **RTO (Recovery Time Objective)**: < 1 hour
2. **RPO (Recovery Point Objective)**: < 1 hour

```yaml
# Velero backup schedule
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: archesai-backup
spec:
  schedule: "0 */1 * * *" # Hourly
  template:
    ttl: 720h # 30 days
    includedNamespaces:
      - production
    storageLocation: s3-backup
```

## Performance Optimization

### Application Optimization

```go
// Connection pooling
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)

// Caching strategy
cache := redis.NewClient(&redis.Options{
    Addr:         os.Getenv("REDIS_URL"),
    PoolSize:     100,
    MinIdleConns: 10,
})

// Query optimization
db.Raw(`
    SELECT * FROM users
    WHERE email = ?
    AND deleted_at IS NULL
    LIMIT 1
`, email).Scan(&user)
```

### CDN Configuration

```nginx
# CloudFlare configuration
location /static {
    add_header Cache-Control "public, max-age=31536000, immutable";
    add_header X-Cache-Status $upstream_cache_status;
}
```

## CI/CD Pipeline

### GitHub Actions Production Deployment

```yaml
name: Deploy to Production

on:
  push:
    tags:
      - "v*"

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v5

      - name: Run tests
        run: |
          go test -v -race ./...

      - name: Security scan
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: archesai/api:${{ github.sha }}

      - name: Build and push images
        run: |
          # Build and tag images
          docker build -t archesai/api:${{ github.sha }} .
          docker build -f web/platform/Dockerfile -t archesai/platform:${{ github.sha }} web/platform/

          # Push to registry
          docker push archesai/api:${{ github.sha }}
          docker push archesai/platform:${{ github.sha }}

      - name: Deploy with Kustomize + Helm
        run: |
          # Template kustomization with new image tags
          helm template archesai deployments/helm-minimal \
            -f deployments/helm-minimal/values-prod.yaml \
            --set api.image.tag=${{ github.sha }} \
            --set platform.image.tag=${{ github.sha }} \
            --set namespace=production > /tmp/kustomization.yaml

          # Apply with Kustomize
          kustomize build /tmp | kubectl apply -f -

      - name: Wait for rollout
        run: |
          kubectl rollout status deployment/archesai-api -n production
          kubectl rollout status deployment/archesai-platform -n production

      - name: Smoke test
        run: |
          curl -f https://api.archesai.com/health || exit 1
```

## Maintenance

### Zero-Downtime Deployments

```yaml
# Rolling update strategy
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
```

### Database Migrations

```bash
#!/bin/bash
# migrate.sh

# Run migrations with zero downtime
goose -dir migrations postgres "$DATABASE_URL" up

# Verify migration
goose -dir migrations postgres "$DATABASE_URL" status
```

## Compliance and Auditing

### Audit Logging

```go
// Audit middleware
func AuditLog() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        err := c.Next()

        log.WithFields(log.Fields{
            "method":     c.Method(),
            "path":       c.Path(),
            "ip":         c.IP(),
            "user_id":    c.Locals("user_id"),
            "duration":   time.Since(start),
            "status":     c.Response().StatusCode(),
        }).Info("API Request")

        return err
    }
}
```

### GDPR Compliance

```sql
-- Data retention policy
DELETE FROM user_sessions
WHERE
  created_at < NOW() - INTERVAL '90 days';

-- Data anonymization
UPDATE users
SET
  email = CONCAT('deleted-', id, '@example.com'),
  name = 'Deleted User',
  phone = NULL
WHERE
  deleted_at IS NOT NULL
  AND deleted_at < NOW() - INTERVAL '30 days';
```

## Troubleshooting Production Issues

### Debug Commands

```bash
# Check application logs
kubectl logs -f deployment/archesai-api -n production

# Database connection issues
psql $DATABASE_URL -c "SELECT count(*) FROM pg_stat_activity;"

# Redis connection issues
redis-cli -u $REDIS_URL ping

# Memory profiling
go tool pprof http://api.archesai.com/debug/pprof/heap

# CPU profiling
go tool pprof http://api.archesai.com/debug/pprof/profile?seconds=30
```

### Common Issues and Solutions

1. **High Memory Usage**
   - Check for memory leaks in goroutines
   - Review cache expiration policies
   - Optimize database queries

2. **Slow Response Times**
   - Enable query logging
   - Add database indexes
   - Implement caching layer

3. **Connection Pool Exhaustion**
   - Increase pool size
   - Reduce connection timeout
   - Implement circuit breaker

## Next Steps

- [Monitoring Setup](../monitoring/overview.md) for detailed observability
- [Security Guide](../security/overview.md) for additional hardening
- [Performance Tuning](../performance/overview.md) for optimization

## Billing
