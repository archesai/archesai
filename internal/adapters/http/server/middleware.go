package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/archesai/archesai/internal/infrastructure/auth"
	"github.com/archesai/archesai/internal/infrastructure/config"
)

// SetupMiddleware configures all middleware for the server.
func (s *Server) SetupMiddleware() {
	// Request ID middleware
	s.echo.Use(middleware.RequestID())

	// Logger middleware
	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			s.logger.Info("request",
				"id", v.RequestID,
				"method", c.Request().Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
				"remote_ip", c.RealIP(),
				"error", v.Error,
			)
			return nil
		},
	}))

	// TODO: Update to use domain-scoped validation middleware
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	// swagger, err := api.GetSwagger()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
	// 	os.Exit(1)
	// }
	// s.echo.Use(echomiddleware.OapiRequestValidator(swagger))

	// Recover middleware
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			s.logger.Error("panic recovered",
				"id", c.Response().Header().Get(echo.HeaderXRequestID),
				"error", err,
				"stack", string(stack),
			)
			return nil
		},
	}))

	// CORS middleware
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(s.config.Cors, ","),
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Request-ID",
			"X-Requested-With",
		},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Request-ID"},
		MaxAge:           86400,
	}))

	// Compression middleware
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health"
		},
	}))

	// Security middleware
	s.echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		ContentSecurityPolicy: "default-src 'self'",
		// FIXME - adjust ContentSecurityPolicy as needed
		//     directives: {
		//       defaultSrc: [`'self'`],
		//       fontSrc: [`'self'`, 'fonts.scalar.com', 'data:'],
		//       imgSrc: [`'self'`, 'data:'],
		//       scriptSrc: [`'self'`, `https: 'unsafe-inline'`, `'unsafe-eval'`],
		//       styleSrc: [`'self'`, `'unsafe-inline'`, 'fonts.scalar.com']
		//     }
		//   }
		ReferrerPolicy: "strict-origin-when-cross-origin",
	}))

	// Rate limiting (basic example - consider using a Redis-based solution in production)
	s.echo.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      10,
				Burst:     30,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// Use IP address as identifier
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			s.logger.Warn("rate limiter error", "error", err, "ip", c.RealIP())
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests",
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			s.logger.Info(
				"rate limit exceeded",
				"identifier",
				identifier,
				"path",
				c.Request().URL.Path,
				"error",
				err,
			)
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Rate limit exceeded",
			})
		},
	}))

	// Body limit middleware
	s.echo.Use(middleware.BodyLimit("10M"))

	// Timeout middleware
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      30 * time.Second,
		ErrorMessage: "Request timeout",
	}))
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// Context keys for authentication
const (
	// AuthUserContextKey is the context key for user authentication
	AuthUserContextKey contextKey = "auth.user"

	// AuthClaimsContextKey is the context key for auth claims
	AuthClaimsContextKey contextKey = "auth.claims"

	// SessionIDContextKey is the context key for session ID
	SessionIDContextKey contextKey = "auth.sessionID"

	// BearerAuthScopes is used by generated handlers for bearer token authentication
	BearerAuthScopes = "bearerAuth.Scopes"

	// SessionCookieScopes is used by generated handlers for session cookie authentication
	SessionCookieScopes = "sessionCookie.Scopes"

	// AuthAPIKeyContextKey is the context key for API token
	AuthAPIKeyContextKey = "auth_api_token"

	// AuthMethodContextKey is the context key for the auth method used
	AuthMethodContextKey = "auth_method"

	// BearerPrefix is the prefix for Bearer token authentication
	BearerPrefix = "Bearer"
)

// Claims represents JWT claims with user information
type Claims struct {
	UserID         uuid.UUID `json:"user_id"`
	SessionID      uuid.UUID `json:"session_id"`
	Email          string    `json:"email"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	Roles          []string  `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware provides authentication middleware using the auth service.
type AuthMiddleware struct {
	authService *auth.Service
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(authService *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth creates middleware that validates JWT tokens.
func (am *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Check for Bearer token format
			tokenString := strings.TrimPrefix(authHeader, BearerPrefix+" ")
			if tokenString == authHeader {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"invalid authorization header format",
				)
			}

			// Validate access token using auth service
			claims, err := am.authService.ValidateAccessToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Add claims to context
			ctx := context.WithValue(c.Request().Context(), AuthUserContextKey, claims.UserID)
			ctx = context.WithValue(ctx, AuthClaimsContextKey, claims)
			ctx = context.WithValue(ctx, SessionIDContextKey, claims.SessionID)
			c.SetRequest(c.Request().WithContext(ctx))

			// Also set in Echo context for controllers
			c.Set("sessionID", claims.SessionID)

			return next(c)
		}
	}
}

// OptionalAuth creates middleware for optional authentication.
func (am *AuthMiddleware) OptionalAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// No auth header, continue without authentication
				return next(c)
			}

			// Check for Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != BearerPrefix {
				// Invalid format, continue without authentication
				return next(c)
			}

			tokenString := parts[1]

			// Try to validate token, but don't fail if invalid
			claims, err := am.authService.ValidateAccessToken(tokenString)
			if err == nil && claims != nil {
				// Add claims to context
				ctx := context.WithValue(c.Request().Context(), AuthUserContextKey, claims.UserID)
				ctx = context.WithValue(ctx, AuthClaimsContextKey, claims)
				ctx = context.WithValue(ctx, SessionIDContextKey, claims.SessionID)
				c.SetRequest(c.Request().WithContext(ctx))

				// Also set in Echo context for controllers
				c.Set("sessionID", claims.SessionID)
			}

			return next(c)
		}
	}
}

// ConfigAuthMiddleware creates JWT authentication middleware using config (legacy)
func ConfigAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			// First, try to get token from cookie
			if cookie, err := c.Cookie("access_token"); err == nil && cookie.Value != "" {
				tokenString = cookie.Value
			} else {
				// Fall back to Authorization header
				authHeader := c.Request().Header.Get("Authorization")
				if authHeader == "" {
					return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization")
				}

				// Check Bearer prefix
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					return echo.NewHTTPError(
						http.StatusUnauthorized,
						"invalid authorization header format",
					)
				}

				tokenString = parts[1]
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(
				tokenString,
				&Claims{},
				func(token *jwt.Token) (any, error) {
					// Validate signing method
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					// Return the secret key for validation
					return []byte(cfg.Auth.Local.JWTSecret), nil
				},
			)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if claims, ok := token.Claims.(*Claims); ok && token.Valid {
				// Set claims in context
				ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
				ctx = context.WithValue(ctx, AuthUserContextKey, claims.UserID)
				ctx = context.WithValue(ctx, SessionIDContextKey, claims.SessionID)
				c.SetRequest(c.Request().WithContext(ctx))

				// Also set in Echo context for controllers
				c.Set("sessionID", claims.SessionID)
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}
	}
}

// OptionalConfigAuthMiddleware creates optional JWT authentication middleware using config (legacy)
func OptionalConfigAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			// First, try to get token from cookie
			if cookie, err := c.Cookie("access_token"); err == nil && cookie.Value != "" {
				tokenString = cookie.Value
			} else {
				// Fall back to Authorization header
				authHeader := c.Request().Header.Get("Authorization")
				if authHeader == "" {
					// No auth header or cookie, continue without authentication
					return next(c)
				}

				// Check Bearer prefix
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					// Invalid format, continue without authentication
					return next(c)
				}

				tokenString = parts[1]
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(
				tokenString,
				&Claims{},
				func(token *jwt.Token) (any, error) {
					// Validate signing method
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					// Return the secret key for validation
					return []byte(cfg.Auth.Local.JWTSecret), nil
				},
			)

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*Claims); ok {
					// Set claims in context
					ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
					ctx = context.WithValue(ctx, AuthUserContextKey, claims.UserID)
					ctx = context.WithValue(ctx, SessionIDContextKey, claims.SessionID)
					c.SetRequest(c.Request().WithContext(ctx))

					// Also set in Echo context for controllers
					c.Set("sessionID", claims.SessionID)
				}
			}

			return next(c)
		}
	}
}

// GetUserFromContext retrieves the user ID from the context
func GetUserFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(AuthUserContextKey).(uuid.UUID)
	return userID, ok
}

// GetUserFromEcho retrieves the user ID from the Echo context
func GetUserFromEcho(c echo.Context) (uuid.UUID, bool) {
	return GetUserFromContext(c.Request().Context())
}

// GetSessionIDFromContext retrieves the session ID from the context
func GetSessionIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	sessionID, ok := ctx.Value(SessionIDContextKey).(uuid.UUID)
	return sessionID, ok
}

// GetClaimsFromContext retrieves the claims from the context
func GetClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(AuthClaimsContextKey).(*Claims)
	return claims, ok
}

// GetClaimsFromEcho retrieves the claims from the Echo context
func GetClaimsFromEcho(c echo.Context) (*Claims, bool) {
	return GetClaimsFromContext(c.Request().Context())
}

// GetTokenClaimsFromContext retrieves the auth service token claims from the context
func GetTokenClaimsFromContext(ctx context.Context) (*auth.TokenClaims, bool) {
	claims, ok := ctx.Value(AuthClaimsContextKey).(*auth.TokenClaims)
	return claims, ok
}

// GetTokenClaimsFromEcho retrieves the auth service token claims from the Echo context
func GetTokenClaimsFromEcho(c echo.Context) (*auth.TokenClaims, bool) {
	return GetTokenClaimsFromContext(c.Request().Context())
}

// RequireRole creates middleware that requires specific roles
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := GetClaimsFromContext(c.Request().Context())
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "no authentication claims")
			}

			// Check if user has any of the required roles
			for _, requiredRole := range roles {
				for _, userRole := range claims.Roles {
					if userRole == requiredRole {
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}

// CustomRateLimiter holds the rate limiters for each client
type CustomRateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// visitor holds the rate limiter and last seen time for a client
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewCustomRateLimiter creates a new rate limiter
func NewCustomRateLimiter(r rate.Limit, b int) *CustomRateLimiter {
	rl := &CustomRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    b,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// cleanupVisitors removes old entries from the visitors map
func (rl *CustomRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// getVisitor returns the rate limiter for the given IP
func (rl *CustomRateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// CustomRateLimitMiddleware creates rate limiting middleware
func CustomRateLimitMiddleware(requestsPerSecond float64, burst int) echo.MiddlewareFunc {
	rl := NewCustomRateLimiter(rate.Limit(requestsPerSecond), burst)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := rl.getVisitor(ip)

			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}

			return next(c)
		}
	}
}

// APIKeyRateLimitMiddleware creates rate limiting middleware based on API key
func APIKeyRateLimitMiddleware(requestsPerSecond float64, burst int) echo.MiddlewareFunc {
	rl := NewCustomRateLimiter(rate.Limit(requestsPerSecond), burst)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get API key from header or context
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				// Fall back to IP-based rate limiting
				apiKey = c.RealIP()
			}

			limiter := rl.getVisitor(apiKey)

			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}

			return next(c)
		}
	}
}
