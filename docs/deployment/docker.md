# Docker Deployment

This guide covers deploying ArchesAI using Docker and Docker Compose.

## Quick Start

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f
```

## Docker Images

### Production Image

Multi-stage build optimized for production:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o archesai cmd/archesai/main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/archesai .
EXPOSE 8080
CMD ["./archesai", "api"]
```

### Development Image

With hot reload support:

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY . .
EXPOSE 8080
CMD ["air"]
```

## Docker Compose Configuration

### Development Setup

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: password
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

  api:
    build:
      context: .
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgresql://admin:password@postgres:5432/archesai
      REDIS_URL: redis://redis:6379
      JWT_SECRET: your-secret-key
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    volumes:
      - .:/app
      - /app/bin

volumes:
  postgres_data:
  redis_data:
```

### Production Setup

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - archesai

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    restart: unless-stopped
    networks:
      - archesai

  api:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      DATABASE_URL: ${DATABASE_URL}
      REDIS_URL: ${REDIS_URL}
      JWT_SECRET: ${JWT_SECRET}
      LOG_LEVEL: info
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    networks:
      - archesai
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  archesai:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
```

## Environment Variables

Create a `.env` file for production:

```bash
# Database
DB_USER=archesai
DB_PASSWORD=secure-password
DB_NAME=archesai_prod
DATABASE_URL=postgresql://archesai:secure-password@postgres:5432/archesai_prod

# Redis
REDIS_PASSWORD=redis-secure-password
REDIS_URL=redis://:redis-secure-password@redis:6379

# Security
JWT_SECRET=your-jwt-secret-key
JWT_REFRESH_SECRET=your-refresh-secret-key

# API Configuration
API_PORT=8080
API_HOST=0.0.0.0
LOG_LEVEL=info

# OpenAI (for chat features)
OPENAI_API_KEY=your-openai-key
```

## Building Images

### Build for Production

```bash
# Build the image
docker build -t archesai:latest .

# Tag for registry
docker tag archesai:latest registry.example.com/archesai:latest

# Push to registry
docker push registry.example.com/archesai:latest
```

### Build for Multiple Architectures

```bash
# Create builder
docker buildx create --name archesai-builder --use

# Build and push multi-arch image
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag registry.example.com/archesai:latest \
  --push .
```

## Volume Management

### Backup Volumes

```bash
# Backup PostgreSQL
docker run --rm \
  -v archesai_postgres_data:/data \
  -v $(pwd)/backups:/backup \
  alpine tar czf /backup/postgres-$(date +%Y%m%d).tar.gz -C /data .

# Backup Redis
docker run --rm \
  -v archesai_redis_data:/data \
  -v $(pwd)/backups:/backup \
  alpine tar czf /backup/redis-$(date +%Y%m%d).tar.gz -C /data .
```

### Restore Volumes

```bash
# Restore PostgreSQL
docker run --rm \
  -v archesai_postgres_data:/data \
  -v $(pwd)/backups:/backup \
  alpine tar xzf /backup/postgres-20240101.tar.gz -C /data

# Restore Redis
docker run --rm \
  -v archesai_redis_data:/data \
  -v $(pwd)/backups:/backup \
  alpine tar xzf /backup/redis-20240101.tar.gz -C /data
```

## Networking

### Custom Network Configuration

```yaml
networks:
  frontend:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
  backend:
    driver: bridge
    internal: true
```

### Reverse Proxy with Traefik

```yaml
services:
  traefik:
    image: traefik:v3.0
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  api:
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`api.archesai.com`)"
      - "traefik.http.routers.api.entrypoints=websecure"
      - "traefik.http.routers.api.tls=true"
```

## Health Checks

### API Health Check

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### Database Health Check

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U $$POSTGRES_USER"]
  interval: 10s
  timeout: 5s
  retries: 5
```

## Monitoring

### With Prometheus and Grafana

```yaml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"
```

## Security Best Practices

1. **Use Secrets Management**

   ```yaml
   secrets:
     db_password:
       external: true
     jwt_secret:
       external: true
   ```

2. **Run as Non-Root User**

   ```dockerfile
   USER 1000:1000
   ```

3. **Limit Resources**

   ```yaml
   deploy:
     resources:
       limits:
         cpus: "2"
         memory: 2G
       reservations:
         cpus: "1"
         memory: 1G
   ```

4. **Network Isolation**

   ```yaml
   networks:
     backend:
       internal: true
   ```

## Troubleshooting

### Common Issues

1. **Container fails to start**

   ```bash
   docker logs archesai_api_1
   docker inspect archesai_api_1
   ```

2. **Database connection issues**

   ```bash
   docker exec -it archesai_postgres_1 psql -U admin -d archesai
   ```

3. **Permission issues**

   ```bash
   docker exec -it archesai_api_1 ls -la /app
   chmod -R 755 ./volumes
   ```

### Debug Mode

```yaml
services:
  api:
    environment:
      LOG_LEVEL: debug
      DEBUG: "true"
    command: ["./archesai", "api", "--debug"]
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Docker Build
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          push: true
          tags: ${{ secrets.REGISTRY }}/archesai:${{ github.sha }}
```

## Next Steps

- [Kubernetes Deployment](kubernetes.md) for orchestration at scale
- [Production Guide](production.md) for hardening and optimization
- [Monitoring Setup](../monitoring/overview.md) for observability
