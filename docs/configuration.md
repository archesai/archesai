# Configuration

Arches uses a YAML configuration file (`.archesai.yaml` by default) with environment variable overrides.

## Configuration File

Create `.archesai.yaml` in your project root:

```yaml
app:
  name: my-app
  version: 1.0.0

server:
  host: localhost
  port: 3001
  cors:
    enabled: true
    origins:
      - http://localhost:3000

database:
  driver: postgres
  host: localhost
  port: 5432
  name: archesdb
  user: postgres
  password: postgres

redis:
  enabled: false
  host: localhost
  port: 6379

auth:
  enabled: true
  jwt_secret: change-me-in-production
  access_token_expiry: 15m
  refresh_token_expiry: 7d
```

## Environment Variables

All config values can be overridden with environment variables using the format `ARCHES_<SECTION>_<KEY>`:

```bash
# Server
export ARCHES_SERVER_PORT=8080
export ARCHES_SERVER_HOST=0.0.0.0

# Database
export ARCHES_DATABASE_HOST=db.example.com
export ARCHES_DATABASE_PORT=5432
export ARCHES_DATABASE_NAME=mydb
export ARCHES_DATABASE_USER=myuser
export ARCHES_DATABASE_PASSWORD=mypassword

# Redis
export ARCHES_REDIS_HOST=redis.example.com
export ARCHES_REDIS_PORT=6379

# Auth
export ARCHES_AUTH_JWT_SECRET=your-secret-key
export ARCHES_AUTH_ACCESS_TOKEN_EXPIRY=15m
```

## CLI Config Commands

```bash
# Show current configuration
archesai config show

# Show as JSON
archesai config show -o json

# Validate configuration
archesai config validate

# Create default config file
archesai config init

# View environment variables
archesai config env
```

## Database Configuration

### PostgreSQL (Default)

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  name: archesdb
  user: postgres
  password: postgres
  ssl_mode: disable # disable, require, verify-ca, verify-full
```

### Connection String

Alternatively, use a connection string:

```bash
export DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=disable"
```

## Redis Configuration

Redis is optional, used for caching and sessions:

```yaml
redis:
  enabled: true
  host: localhost
  port: 6379
  password: ""
  db: 0
```

## Auth Configuration

```yaml
auth:
  enabled: true
  jwt_secret: ${JWT_SECRET} # Use environment variable
  access_token_expiry: 15m
  refresh_token_expiry: 7d

  password:
    min_length: 8
    bcrypt_cost: 10

  oauth:
    google:
      enabled: false
      client_id: ${GOOGLE_CLIENT_ID}
      client_secret: ${GOOGLE_CLIENT_SECRET}
    github:
      enabled: false
      client_id: ${GITHUB_CLIENT_ID}
      client_secret: ${GITHUB_CLIENT_SECRET}
```

## Logging

```yaml
logging:
  level: info # debug, info, warn, error
  format: json # json, text
  output: stdout # stdout, stderr, file
```

Or via environment:

```bash
export ARCHES_LOG_LEVEL=debug
```

## Production Recommendations

1. **Use environment variables** for secrets (JWT_SECRET, database passwords)
2. **Enable SSL** for database connections
3. **Set secure CORS origins** (not `*`)
4. **Use strong JWT secrets** (32+ random bytes)
5. **Enable Redis** for session management at scale

Example production config:

```yaml
server:
  host: 0.0.0.0
  port: 3001
  cors:
    enabled: true
    origins:
      - https://myapp.com

database:
  driver: postgres
  host: ${DATABASE_HOST}
  port: 5432
  name: ${DATABASE_NAME}
  user: ${DATABASE_USER}
  password: ${DATABASE_PASSWORD}
  ssl_mode: require

auth:
  enabled: true
  jwt_secret: ${JWT_SECRET}
```
