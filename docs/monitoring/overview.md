# Monitoring and Observability

Arches deploys a comprehensive monitoring stack with Grafana and Loki, providing built-in
dashboards for immediate visibility into system health and performance.

## Stack Overview

### Core Components

- **Grafana**: Visualization and dashboarding platform with pre-configured dashboards
- **Loki**: Log aggregation system for centralized logging
- **Prometheus**: Metrics collection and alerting
- **Promtail**: Log shipping agent for Loki

## Built-in Dashboards

Arches comes with pre-configured Grafana dashboards that provide immediate insights:

### 1. Application Overview Dashboard

- Request rate and response times
- Error rates by endpoint
- Active connections and goroutines
- Memory and CPU usage
- Database connection pool metrics

### 2. Infrastructure Dashboard

- Node resource utilization
- Pod status and restarts
- Network traffic patterns
- Disk I/O and usage
- Container resource limits

### 3. Business Metrics Dashboard

- User registrations and logins
- API usage by endpoint
- Organization activity
- Workflow execution metrics
- Content processing statistics

### 4. Logs Dashboard (Loki)

- Real-time log streaming
- Log level distribution
- Error log aggregation
- Request tracing
- Structured query capabilities

## Quick Start

### Deploy Monitoring Stack

```bash
# Deploy using Helm
helm install monitoring archesai/monitoring-stack

# Or using docker-compose
docker-compose -f docker-compose.monitoring.yml up -d
```

### Access Dashboards

```bash
# Default URLs
Grafana: http://localhost:3000
Loki: http://localhost:3100
Prometheus: http://localhost:9090

# Default credentials
Username: admin
Password: admin (change on first login)
```

## Configuration

### Grafana Data Sources

Pre-configured data sources include:

```yaml
datasources:
  - name: Prometheus
    type: prometheus
    url: http://prometheus:9090
    isDefault: true

  - name: Loki
    type: loki
    url: http://loki:3100
    jsonData:
      maxLines: 1000
```

### Loki Configuration

```yaml
# loki-config.yaml
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    address: 127.0.0.1
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/boltdb-shipper-active
    cache_location: /loki/boltdb-shipper-cache
    shared_store: filesystem
  filesystem:
    directory: /loki/chunks

limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h
```

## Dashboard Features

### Real-time Metrics

- Live updating graphs with 5-second refresh
- Customizable time ranges
- Drill-down capabilities
- Correlation between metrics

### Log Analysis

- Full-text search across all logs
- LogQL query language support
- Context viewing for log entries
- Export capabilities

### Alerting

- Pre-configured alert rules
- Multiple notification channels
- Alert history and silencing
- SLA tracking

## Integration

### Application Instrumentation

Arches automatically exports metrics and logs:

```go
// Metrics are automatically exposed at /metrics
// Logs are shipped to Loki via Promtail
// No additional configuration required
```

### Custom Metrics

Add custom metrics easily:

```go
// Custom counter example
requestCounter := prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "archesai_custom_requests_total",
        Help: "Total number of custom requests",
    },
    []string{"endpoint", "status"},
)
```

## Alerts

Pre-configured alerts include:

- High error rate (greater than 5% of requests)
- High latency (p95 greater than 1s)
- Low disk space (less than 10% free)
- Pod restarts (more than 3 in 5 minutes)
- Database connection issues
- Memory pressure (greater than 80% usage)

## Best Practices

1. **Retention**: Logs retained for 30 days, metrics for 90 days
2. **Sampling**: Automatic sampling for high-volume endpoints
3. **Cardinality**: Labels kept minimal to prevent metric explosion
4. **Security**: TLS enabled, authentication required
5. **Backup**: Daily backups of Grafana dashboards and configurations

## Troubleshooting

### Common Issues

1. **No data in dashboards**

   ```bash
   # Check Prometheus targets
   curl http://localhost:9090/api/v1/targets

   # Verify Loki is receiving logs
   curl http://localhost:3100/ready
   ```

2. **High memory usage**
   - Adjust retention policies
   - Increase resource limits
   - Enable log sampling

3. **Slow queries**
   - Add appropriate indexes
   - Optimize LogQL queries
   - Use time range filters
