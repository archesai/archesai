# Security Best Practices

## Overview

This guide outlines security best practices for ArchesAI, covering authentication, authorization,
data protection, infrastructure security, and compliance requirements.

## Authentication Security

### Password Requirements

```go
type PasswordPolicy struct {
    MinLength            int  // Minimum 12 characters
    RequireUppercase     bool // At least one uppercase
    RequireLowercase     bool // At least one lowercase
    RequireNumbers       bool // At least one number
    RequireSpecialChars  bool // At least one special character
    PreventCommon        bool // Block common passwords
    PreventUserInfo      bool // Block passwords containing user info
    PasswordHistory      int  // Prevent reuse of last N passwords
}
```

### Multi-Factor Authentication

```go
// TOTP Implementation
import "github.com/pquerna/otp/totp"

func GenerateTOTP(user *User) (string, error) {
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "ArchesAI",
        AccountName: user.Email,
        Algorithm:   otp.AlgorithmSHA256,
    })
    return key.Secret(), err
}

func ValidateTOTP(secret, code string) bool {
    return totp.Validate(code, secret)
}
```

### Session Management

```go
// Secure session configuration
type SessionConfig struct {
    Secure     bool          // HTTPS only
    HttpOnly   bool          // No JavaScript access
    SameSite   http.SameSite // CSRF protection
    MaxAge     time.Duration // Session timeout
    IdleTimeout time.Duration // Idle timeout
}

var secureSession = SessionConfig{
    Secure:      true,
    HttpOnly:    true,
    SameSite:    http.SameSiteStrictMode,
    MaxAge:      24 * time.Hour,
    IdleTimeout: 30 * time.Minute,
}
```

## Authorization

### Role-Based Access Control (RBAC)

```go
// Permission checking
func CheckPermission(user *User, resource string, action string) bool {
    // Get user's roles
    roles := getUserRoles(user)

    // Check each role's permissions
    for _, role := range roles {
        if hasPermission(role, resource, action) {
            return true
        }
    }

    return false
}

// Middleware for authorization
func RequirePermission(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := getCurrentUser(c)
        if !CheckPermission(user, resource, action) {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        c.Next()
    }
}
```

### Resource-Level Security

```sql
-- Row-level security in PostgreSQL
CREATE POLICY organization_isolation ON content FOR ALL TO application_user USING (
  organization_id = current_setting('app.current_org_id')::uuid
);

ALTER TABLE content ENABLE ROW LEVEL SECURITY;
```

## Data Protection

### Encryption at Rest

```go
// AES-256 encryption for sensitive data
import "crypto/aes"
import "crypto/cipher"

func EncryptData(plaintext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}
```

### Encryption in Transit

```nginx
# Nginx SSL configuration
server {
    listen 443 ssl http2;

    # Modern SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;

    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;

    # SSL session caching
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # OCSP stapling
    ssl_stapling on;
    ssl_stapling_verify on;
}
```

### Secrets Management

```yaml
# Kubernetes secrets
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  database-url: <base64-encoded>
  jwt-secret: <base64-encoded>
  api-keys: <base64-encoded>
```

```go
// Environment variable validation
func LoadSecrets() error {
    required := []string{
        "DATABASE_URL",
        "JWT_SECRET",
        "ENCRYPTION_KEY",
        "REDIS_URL",
    }

    for _, key := range required {
        if os.Getenv(key) == "" {
            return fmt.Errorf("missing required secret: %s", key)
        }
    }

    return nil
}
```

## Input Validation

### SQL Injection Prevention

```go
// Always use parameterized queries
func GetUser(email string) (*User, error) {
    // Safe - uses parameterized query
    query := "SELECT * FROM users WHERE email = $1"
    return db.QueryRow(query, email).Scan(&user)

    // NEVER do this:
    // query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
}
```

### XSS Prevention

```go
// Sanitize user input
import "html/template"

func RenderHTML(userContent string) template.HTML {
    // Automatically escapes HTML
    return template.HTMLEscapeString(userContent)
}

// Content Security Policy
func SetSecurityHeaders(c *gin.Context) {
    c.Header("Content-Security-Policy",
        "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
    c.Header("X-Content-Type-Options", "nosniff")
    c.Header("X-Frame-Options", "DENY")
    c.Header("X-XSS-Protection", "1; mode=block")
}
```

### Request Validation

```go
// Input validation using struct tags
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=12"`
    Name     string `json:"name" binding:"required,min=2,max=100"`
}

// Rate limiting
import "golang.org/x/time/rate"

var limiter = rate.NewLimiter(rate.Every(time.Second), 10)

func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.AbortWithStatus(http.StatusTooManyRequests)
            return
        }
        c.Next()
    }
}
```

## API Security

### API Key Management

```go
// API key generation and validation
func GenerateAPIKey() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func ValidateAPIKey(key string) (*Application, error) {
    // Hash the key before database lookup
    hashedKey := hashAPIKey(key)

    var app Application
    err := db.Get(&app, "SELECT * FROM applications WHERE api_key_hash = $1", hashedKey)
    if err != nil {
        return nil, ErrInvalidAPIKey
    }

    // Check if key is expired or revoked
    if app.KeyExpiredAt.Before(time.Now()) || app.KeyRevoked {
        return nil, ErrAPIKeyExpired
    }

    return &app, nil
}
```

### CORS Configuration

```go
// Strict CORS policy
func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"https://app.archesai.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:          12 * time.Hour,
    })
}
```

## Infrastructure Security

### Network Security

```yaml
# Kubernetes NetworkPolicy
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-netpol
spec:
  podSelector:
    matchLabels:
      app: api
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: nginx
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: database
      ports:
        - protocol: TCP
          port: 5432
```

### Container Security

```dockerfile
# Secure Docker image
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache ca-certificates

# Build with security flags
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags='-w -s -extldflags "-static"' \
  -a -installsuffix cgo -o app .

# Minimal runtime image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

# Run as non-root user
USER 1000:1000
ENTRYPOINT ["/app"]
```

### Database Security

```sql
-- Least privilege principle
CREATE USER app_user
WITH
  PASSWORD 'secure_password';

GRANT CONNECT ON DATABASE archesai TO app_user;

GRANT USAGE ON SCHEMA public TO app_user;

GRANT
SELECT
,
  INSERT,
UPDATE,
DELETE ON ALL TABLES IN SCHEMA public TO app_user;

-- Audit logging
CREATE TABLE audit_log (
  id SERIAL PRIMARY KEY,
  user_id UUID,
  action VARCHAR(50),
  resource VARCHAR(100),
  resource_id UUID,
  ip_address INET,
  user_agent TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);
```

## Logging and Monitoring

### Security Logging

```go
// Structured security logging
type SecurityEvent struct {
    Type      string    `json:"type"`
    UserID    string    `json:"user_id"`
    IP        string    `json:"ip"`
    Action    string    `json:"action"`
    Resource  string    `json:"resource"`
    Result    string    `json:"result"`
    Timestamp time.Time `json:"timestamp"`
}

func LogSecurityEvent(event SecurityEvent) {
    logger.Info("security_event",
        zap.String("type", event.Type),
        zap.String("user_id", event.UserID),
        zap.String("ip", event.IP),
        zap.String("action", event.Action),
        zap.String("resource", event.Resource),
        zap.String("result", event.Result),
    )
}
```

### Intrusion Detection

```go
// Detect suspicious activities
func DetectAnomalies(userID string) bool {
    // Check for rapid failed login attempts
    failedAttempts := getRecentFailedLogins(userID, 5*time.Minute)
    if failedAttempts > 5 {
        triggerAlert("Multiple failed login attempts", userID)
        return true
    }

    // Check for unusual access patterns
    locations := getRecentAccessLocations(userID, 24*time.Hour)
    if hasGeographicalAnomaly(locations) {
        triggerAlert("Geographical anomaly detected", userID)
        return true
    }

    return false
}
```

## Incident Response

### Response Plan

1. **Detection**
   - Monitor security alerts
   - Review audit logs
   - User reports

2. **Containment**
   - Isolate affected systems
   - Revoke compromised credentials
   - Block malicious IPs

3. **Eradication**
   - Remove malicious code
   - Patch vulnerabilities
   - Update security rules

4. **Recovery**
   - Restore from backups
   - Verify system integrity
   - Resume normal operations

5. **Lessons Learned**
   - Document incident
   - Update security policies
   - Improve detection mechanisms

### Security Checklist

#### Development

- [ ] Code review for security issues
- [ ] Dependency vulnerability scanning
- [ ] Static code analysis (SAST)
- [ ] Dynamic testing (DAST)
- [ ] Secrets scanning in code

#### Deployment

- [ ] SSL/TLS certificates valid
- [ ] Security headers configured
- [ ] CORS properly configured
- [ ] Rate limiting enabled
- [ ] WAF rules updated

#### Operations

- [ ] Regular security patches
- [ ] Audit log review
- [ ] Access control review
- [ ] Backup verification
- [ ] Incident response drills

## Compliance

### GDPR Requirements

```go
// Data deletion for GDPR
func DeleteUserData(userID string) error {
    // Start transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Delete or anonymize user data
    queries := []string{
        "DELETE FROM user_sessions WHERE user_id = $1",
        "DELETE FROM user_content WHERE user_id = $1",
        "UPDATE users SET email = 'deleted@user.com', name = 'Deleted User' WHERE id = $1",
    }

    for _, query := range queries {
        if _, err := tx.Exec(query, userID); err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

### SOC 2 Requirements

- Encryption of data at rest and in transit
- Access controls and authentication
- Availability monitoring
- Change management procedures
- Incident response procedures

### HIPAA Compliance

- PHI encryption (AES-256)
- Access logging and auditing
- Business Associate Agreements (BAAs)
- Data backup and recovery
- Security risk assessments

## Security Tools

### Dependency Scanning

```bash
# Go vulnerability check
go install github.com/sonatype-nexus-community/nancy@latest
go list -json -deps | nancy sleuth

# Node.js vulnerability check
npm audit
yarn audit
```

### Security Testing

```bash
# OWASP ZAP scanning
docker run -t owasp/zap2docker-stable zap-baseline.py \
  -t https://app.archesai.com

# SQLMap for SQL injection testing
sqlmap -u "https://api.archesai.com/users?id=1" \
  --batch --random-agent
```

## Related Documentation

- [Authentication](../features/auth.md) - Authentication implementation
- [Architecture](../architecture/system-design.md) - Security architecture
- [Deployment](../deployment/production.md) - Production security
- [Monitoring](../monitoring/overview.md) - Security monitoring
