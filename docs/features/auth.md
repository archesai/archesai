# Authentication

## Overview

Arches provides a comprehensive authentication system built on modern security standards with JWT
tokens, session management, and role-based access control.

## Core Features

### JWT Authentication

- **Token Generation**: Secure JWT tokens with configurable expiration
- **Refresh Tokens**: Automatic token refresh for seamless user experience
- **Token Validation**: Middleware-based validation on all protected routes
- **Revocation**: Support for token blacklisting and immediate invalidation

### Session Management

- **Redis-Backed Sessions**: High-performance session storage
- **Session Persistence**: Configurable session timeout and persistence
- **Multi-Device Support**: Users can maintain sessions across multiple devices
- **Session Monitoring**: Track active sessions and last activity

### OAuth Integration

- **Multiple Providers**: Support for Google, GitHub, Microsoft, and custom OAuth2 providers
- **Single Sign-On**: Seamless SSO experience for enterprise deployments
- **Account Linking**: Link multiple OAuth providers to a single account
- **Custom Scopes**: Configure OAuth scopes per provider

## Implementation

### Authentication Flow

```go
// Login endpoint
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "secure_password"
}

// Response
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "refresh_token_here",
  "user": {
    "id": "user_id",
    "email": "user@example.com",
    "role": "user"
  }
}
```

### Middleware Configuration

```go
// Protected route example
router.GET("/api/v1/protected",
  middleware.RequireAuth(),
  middleware.RequireRole("admin"),
  handler.ProtectedEndpoint,
)
```

### Password Management

- **Bcrypt Hashing**: Industry-standard password hashing
- **Password Policies**: Configurable complexity requirements
- **Reset Flow**: Secure password reset via email tokens
- **Password History**: Prevent reuse of recent passwords

## Security Features

### Rate Limiting

- Login attempt throttling
- IP-based rate limiting
- Distributed rate limiting with Redis

### Two-Factor Authentication

- TOTP support (Google Authenticator compatible)
- Backup codes for recovery
- SMS authentication (optional)

### Security Headers

- CSRF protection
- XSS prevention
- Content Security Policy
- HSTS enforcement

## Role-Based Access Control

### Roles and Permissions

```yaml
roles:
  admin:
    permissions:
      - users:read
      - users:write
      - organizations:manage
      - system:configure

  member:
    permissions:
      - content:read
      - content:write
      - workflows:execute

  viewer:
    permissions:
      - content:read
      - workflows:view
```

### Organization-Level Roles

- **Owner**: Full organization control
- **Admin**: Member management and settings
- **Member**: Standard access to resources
- **Guest**: Read-only access

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/register` - New user registration
- `POST /api/v1/auth/verify-email` - Email verification
- `POST /api/v1/auth/reset-password` - Password reset request
- `POST /api/v1/auth/reset-password/confirm` - Confirm password reset

### User Management

- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/me` - Update current user
- `DELETE /api/v1/users/me` - Delete account
- `GET /api/v1/users/me/sessions` - List active sessions
- `DELETE /api/v1/users/me/sessions/:id` - Revoke session

## Configuration

### Environment Variables

```bash
# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# OAuth Providers
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# Session Configuration
SESSION_TIMEOUT=24h
SESSION_SECURE=true
SESSION_HTTPONLY=true
SESSION_SAMESITE=strict

# Security
BCRYPT_COST=12
RATE_LIMIT_LOGIN=5/min
RATE_LIMIT_API=100/min
```

## Testing

### Unit Tests

```bash
go test ./internal/auth/...
```

### Integration Tests

```bash
go test -tags=integration ./internal/auth/...
```

### Security Testing

- Automated penetration testing with OWASP ZAP
- JWT token validation tests
- Session hijacking prevention tests
- SQL injection prevention tests

## Best Practices

### Secure Defaults

- Passwords require minimum 12 characters
- Sessions expire after 24 hours of inactivity
- Failed login attempts trigger exponential backoff
- All tokens are cryptographically signed

### Monitoring

- Track failed login attempts
- Monitor unusual session patterns
- Alert on privilege escalation attempts
- Log all authentication events

### Compliance

- GDPR-compliant data handling
- SOC 2 audit logging
- HIPAA-ready encryption
- PCI DSS password standards

## Migration Guide

### From v1 to v2

1. Update JWT library to latest version
2. Migrate session storage from memory to Redis
3. Update password hashing from SHA-256 to bcrypt
4. Implement new RBAC system

### Database Schema

```sql
-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  email_verified BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- Sessions table
CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users (id),
  token_hash VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL
);

-- Roles table
CREATE TABLE user_roles (
  user_id UUID REFERENCES users (id),
  organization_id UUID REFERENCES organizations (id),
  role VARCHAR(50) NOT NULL,
  PRIMARY KEY (user_id, organization_id)
);
```

## Troubleshooting

### Common Issues

#### Invalid JWT Token

- Verify JWT_SECRET is correctly set
- Check token expiration time
- Ensure clock synchronization between servers

#### Session Not Persisting

- Verify Redis connection
- Check session cookie settings
- Ensure CORS configuration allows credentials

#### OAuth Login Failing

- Verify redirect URLs in provider configuration
- Check client ID and secret
- Ensure callback URL is whitelisted

## Related Documentation

- [Deployment Guide](../deployment/production.md) - Production deployment
