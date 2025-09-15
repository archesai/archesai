# Performance Optimization

## Overview

This guide covers performance optimization strategies for ArchesAI, including database tuning,
application optimization, caching strategies, and monitoring best practices.

## Database Performance

### Connection Pooling

```yaml
database:
  pool:
    max_connections: 100
    min_connections: 10
    max_lifetime: 30m
    idle_timeout: 10m
    connection_timeout: 5s
```

### Query Optimization

#### Index Strategy

```sql
-- Essential indexes for performance
CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_organizations_slug ON organizations (slug);

CREATE INDEX idx_members_org_user ON organization_members (organization_id, user_id);

CREATE INDEX idx_content_org_created ON content (organization_id, created_at DESC);

CREATE INDEX idx_workflows_org_status ON workflows (organization_id, status);
```

#### Query Analysis

```sql
-- Analyze query performance
EXPLAIN
ANALYZE
SELECT
  *
FROM
  users
WHERE
  email = 'user@example.com';

-- Find slow queries
SELECT
  query,
  mean_exec_time,
  calls
FROM
  pg_stat_statements
WHERE
  mean_exec_time > 100
ORDER BY
  mean_exec_time DESC;
```

### PostgreSQL Tuning

```ini
# postgresql.conf optimizations
shared_buffers = 256MB          # 25% of RAM
effective_cache_size = 1GB      # 50-75% of RAM
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1          # For SSD storage
effective_io_concurrency = 200  # For SSD storage
work_mem = 4MB
huge_pages = try
```

### pgvector Optimization

```sql
-- Optimize vector search performance
CREATE INDEX ON content USING ivfflat (embedding vector_cosine_ops)
WITH
  (lists = 100);

-- Tune for accuracy vs speed
SET
  ivfflat.probes = 10;

-- Increase for better accuracy
-- Parallel vector search
SET
  max_parallel_workers_per_gather = 4;

SET
  max_parallel_workers = 8;
```

## Application Performance

### Go Performance Best Practices

#### Memory Management

```go
// Use sync.Pool for frequently allocated objects
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func ProcessData(data []byte) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    // Use buffer
}
```

#### Concurrent Processing

```go
// Efficient worker pool pattern
func WorkerPool(jobs <-chan Job, results chan<- Result) {
    var wg sync.WaitGroup
    workers := runtime.NumCPU()

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                result := processJob(job)
                results <- result
            }
        }()
    }

    wg.Wait()
    close(results)
}
```

#### HTTP Client Optimization

```go
// Reuse HTTP clients with connection pooling
var httpClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  false,
        DisableKeepAlives:   false,
    },
}
```

### API Response Optimization

#### Pagination

```go
// Cursor-based pagination for large datasets
func GetContentPaginated(cursor string, limit int) ([]Content, string, error) {
    query := `
        SELECT * FROM content
        WHERE created_at < $1
        ORDER BY created_at DESC
        LIMIT $2
    `
    // Implementation
}
```

#### Field Selection

```graphql
# Allow clients to request only needed fields
query {
  users(limit: 10) {
    id
    email
    # Skip expensive fields unless needed
  }
}
```

#### Response Compression

```go
// Enable gzip compression
import "github.com/gin-contrib/gzip"

router.Use(gzip.Gzip(gzip.DefaultCompression))
```

## Caching Strategies

### Redis Configuration

```yaml
redis:
  max_memory: 1gb
  max_memory_policy: allkeys-lru
  save: "" # Disable persistence for cache
  tcp_keepalive: 60
  timeout: 0
  databases: 16
```

### Cache Patterns

#### Cache-Aside Pattern

```go
func GetUser(id string) (*User, error) {
    // Check cache first
    cached, err := redis.Get(ctx, "user:"+id)
    if err == nil {
        return unmarshalUser(cached), nil
    }

    // Load from database
    user, err := db.GetUser(id)
    if err != nil {
        return nil, err
    }

    // Cache for future requests
    redis.Set(ctx, "user:"+id, marshalUser(user), 5*time.Minute)
    return user, nil
}
```

#### Write-Through Cache

```go
func UpdateUser(user *User) error {
    // Update database
    if err := db.UpdateUser(user); err != nil {
        return err
    }

    // Update cache
    redis.Set(ctx, "user:"+user.ID, marshalUser(user), 5*time.Minute)
    return nil
}
```

### Cache Invalidation

```go
// Tag-based invalidation
func InvalidateOrganizationCache(orgID string) {
    pattern := fmt.Sprintf("org:%s:*", orgID)
    keys, _ := redis.Keys(ctx, pattern).Result()
    if len(keys) > 0 {
        redis.Del(ctx, keys...)
    }
}
```

## CDN and Static Assets

### CDN Configuration

```nginx
# Nginx configuration for static assets
location ~* \.(jpg|jpeg|png|gif|ico|css|js|woff2)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
    add_header Vary "Accept-Encoding";
    gzip_static on;
}
```

### Asset Optimization

```javascript
// Webpack configuration for production
module.exports = {
  optimization: {
    minimize: true,
    splitChunks: {
      chunks: "all",
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name: "vendors",
          priority: 10,
        },
      },
    },
  },
};
```

## Load Balancing

### HAProxy Configuration

```haproxy
global
    maxconn 4096
    tune.ssl.default-dh-param 2048

defaults
    mode http
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms
    option httplog

backend api_servers
    balance roundrobin
    option httpchk GET /health
    server api1 10.0.1.10:8080 check
    server api2 10.0.1.11:8080 check
    server api3 10.0.1.12:8080 check
```

## Monitoring and Profiling

### Application Metrics

```go
// Prometheus metrics
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request latency",
        },
        []string{"method", "endpoint", "status"},
    )

    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "db_query_duration_seconds",
            Help: "Database query latency",
        },
        []string{"query_type"},
    )
)
```

### CPU Profiling

```go
import _ "net/http/pprof"

// Enable profiling endpoint
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// Profile CPU usage
// go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

### Memory Profiling

```bash
# Capture heap profile
curl http://localhost:6060/debug/pprof/heap > heap.prof

# Analyze memory usage
go tool pprof heap.prof
```

## Optimization Checklist

### Database

- [ ] Connection pooling configured
- [ ] Appropriate indexes created
- [ ] Query performance analyzed
- [ ] Vacuum and analyze scheduled
- [ ] Slow query log enabled

### Application

- [ ] HTTP client connection pooling
- [ ] Concurrent processing implemented
- [ ] Memory pools for hot paths
- [ ] Response compression enabled
- [ ] API pagination implemented

### Caching

- [ ] Redis configured for caching
- [ ] Cache-aside pattern implemented
- [ ] Cache invalidation strategy
- [ ] CDN for static assets
- [ ] Browser caching headers

### Monitoring

- [ ] Application metrics collection
- [ ] Database monitoring
- [ ] Error tracking
- [ ] Performance alerting
- [ ] Regular profiling

## Performance Targets

### API Response Times

- p50: < 50ms
- p95: < 200ms
- p99: < 500ms

### Database Queries

- Simple queries: < 10ms
- Complex queries: < 100ms
- Batch operations: < 1s

### Throughput

- API requests: 10,000 req/s
- Concurrent users: 5,000
- Database connections: 100

## Scaling Strategies

### Horizontal Scaling

```yaml
# Kubernetes HPA
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api
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

```yaml
# Read replicas configuration
primary:
  host: db-primary.example.com

replicas:
  - host: db-replica-1.example.com
  - host: db-replica-2.example.com

routing:
  reads: replicas # Route reads to replicas
  writes: primary # Route writes to primary
```

## Troubleshooting Performance Issues

### Slow API Responses

1. Check database query performance
2. Review cache hit rates
3. Analyze CPU and memory usage
4. Check network latency
5. Review application logs

### High Memory Usage

1. Profile heap allocations
2. Check for memory leaks
3. Review cache sizes
4. Optimize data structures
5. Implement object pooling

### Database Bottlenecks

1. Analyze slow query log
2. Check missing indexes
3. Review connection pool settings
4. Consider query optimization
5. Evaluate need for read replicas

## Related Documentation

- [Monitoring](../monitoring/overview.md) - Monitoring setup
- [Architecture](../architecture/system-design.md) - System design
- [Deployment](../deployment/production.md) - Production deployment
- [Testing](../guides/testing.md) - Performance testing
