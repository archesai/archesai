# Authentication

Arches generates applications with built-in JWT authentication.

## Enabling Authentication

Add security schemes to your OpenAPI spec:

```yaml
openapi: 3.1.0
info:
  title: My API
  version: 1.0.0

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - bearerAuth: []

paths:
  /todos:
    get:
      security:
        - bearerAuth: []
      # ... rest of endpoint
```

## Configuration

Configure authentication in `.archesai.yaml`:

```yaml
auth:
  enabled: true
  jwt_secret: ${JWT_SECRET} # Use environment variable
  access_token_expiry: 15m
  refresh_token_expiry: 7d

  password:
    min_length: 8
    bcrypt_cost: 10
```

Environment variables:

```bash
export JWT_SECRET=your-secret-key-at-least-32-chars
export ARCHES_AUTH_ACCESS_TOKEN_EXPIRY=15m
export ARCHES_AUTH_REFRESH_TOKEN_EXPIRY=7d
```

## API Endpoints

Generated auth endpoints:

| Endpoint                | Method | Description               |
| ----------------------- | ------ | ------------------------- |
| `/auth/register`        | POST   | Register new user         |
| `/auth/login`           | POST   | Login with email/password |
| `/auth/logout`          | POST   | Logout current session    |
| `/auth/refresh`         | POST   | Refresh access token      |
| `/auth/verify-email`    | POST   | Verify email address      |
| `/auth/forgot-password` | POST   | Request password reset    |
| `/auth/reset-password`  | POST   | Reset password with token |

## Usage Examples

### Register

```bash
curl -X POST http://localhost:3001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "password": "secure-password-123"
  }'
```

### Login

```bash
curl -X POST http://localhost:3001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure-password-123"
  }'
```

Response:

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "refresh_token_here",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

### Using the Token

```bash
curl http://localhost:3001/todos \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### Refresh Token

```bash
curl -X POST http://localhost:3001/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken": "refresh_token_here"}'
```

## Token System

### Access Token

- Short-lived (default: 15 minutes)
- Used for API authentication
- Include in `Authorization: Bearer <token>` header

### Refresh Token

- Long-lived (default: 7 days)
- Used to obtain new access tokens
- Store securely (HttpOnly cookie recommended)

## Middleware

Apply authentication middleware to protect routes:

```go
// Protected routes
api := e.Group("/api")
api.Use(middleware.Auth())

// Public routes (no middleware)
e.POST("/auth/login", authHandler.Login)
e.POST("/auth/register", authHandler.Register)
```

## Getting Current User

In handlers, access the authenticated user:

```go
func (h *TodoHandler) Create(c echo.Context) error {
    userID := c.Get("user_id").(uuid.UUID)
    // Use userID for authorization
}
```

## OAuth Integration (Optional)

Configure OAuth providers:

```yaml
auth:
  oauth:
    google:
      enabled: true
      client_id: ${GOOGLE_CLIENT_ID}
      client_secret: ${GOOGLE_CLIENT_SECRET}
    github:
      enabled: true
      client_id: ${GITHUB_CLIENT_ID}
      client_secret: ${GITHUB_CLIENT_SECRET}
```

OAuth endpoints:

```bash
# Initiate OAuth flow
GET /auth/oauth/google/authorize

# Callback (handled automatically)
GET /auth/oauth/google/callback
```

## Database Schema

Authentication requires these tables (auto-generated):

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255),
  password_hash VARCHAR(255),
  email_verified BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(255) UNIQUE NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

## Security Best Practices

1. **Use environment variables** for JWT_SECRET
2. **Use HTTPS** in production
3. **Set secure cookie options** (HttpOnly, Secure, SameSite)
4. **Use short access token expiry** (15 minutes recommended)
5. **Implement rate limiting** on auth endpoints
6. **Validate password strength** on registration
