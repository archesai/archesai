package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Strategy represents the authentication method used
type Strategy string

const (
	// StrategyJWT indicates JWT token authentication
	StrategyJWT Strategy = "jwt"
	// StrategySession indicates session-based authentication
	StrategySession Strategy = "session"
	// StrategyAPIKey indicates API key authentication
	StrategyAPIKey Strategy = "api_key"
	// StrategyOAuth indicates OAuth authentication
	StrategyOAuth Strategy = "oauth"
)

// AuthStrategyContextKey is the context key for storing the authentication strategy used
const AuthStrategyContextKey contextKey = "auth.strategy"

// ComposeAuthStrategies creates a middleware that tries multiple authentication strategies
// in order, succeeding if any one of them succeeds
func ComposeAuthStrategies(strategies ...echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var lastError error

			// Try each strategy in order
			for _, strategy := range strategies {
				// Create a test handler to check if auth succeeds
				testHandler := func(c echo.Context) error {
					// Auth succeeded, continue with the real handler
					return next(c)
				}

				// Try the strategy
				err := strategy(testHandler)(c)
				if err == nil {
					// Success - one strategy worked
					return nil
				}

				// Check if it's an HTTP error
				if httpErr, ok := err.(*echo.HTTPError); ok {
					// If it's not an auth error, return it immediately
					if httpErr.Code != http.StatusUnauthorized {
						return err
					}
				}

				lastError = err
			}

			// All strategies failed
			if lastError != nil {
				return lastError
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}
	}
}

// RequireAuthConfig defines configuration for authentication middleware
type RequireAuthConfig struct {
	// Strategies specifies which authentication methods are allowed
	Strategies []Strategy
	// Optional makes authentication optional (doesn't fail if no auth provided)
	Optional bool
	// RequireVerifiedEmail requires the user to have a verified email
	RequireVerifiedEmail bool
	// RequireOrganization requires the user to belong to an organization
	RequireOrganization bool
	// RequiredScopes specifies the scopes required for this endpoint
	RequiredScopes []string
	// RequiredPermissions specifies the permissions required
	RequiredPermissions []string
	// RateLimit specifies custom rate limiting for this endpoint
	RateLimit *RateLimitConfig
}

// RateLimitConfig specifies rate limiting configuration
type RateLimitConfig struct {
	// RequestsPerMinute specifies the maximum requests per minute
	RequestsPerMinute int
	// BurstSize specifies the burst size for rate limiting
	BurstSize int
	// ByStrategy allows different rate limits per auth strategy
	ByStrategy map[Strategy]int
}

// RequireAuth creates a configurable authentication middleware
func RequireAuth(authService *Service, config RequireAuthConfig, logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var authenticated bool
			var authStrategy Strategy

			// Build list of middleware to try based on allowed strategies
			var middlewares []echo.MiddlewareFunc

			for _, strategy := range config.Strategies {
				switch strategy {
				case StrategyJWT:
					middlewares = append(middlewares, wrapWithStrategy(
						Middleware(authService, logger),
						StrategyJWT,
					))
				case StrategyAPIKey:
					middlewares = append(middlewares, wrapWithStrategy(
						APIKeyMiddleware(authService, logger),
						StrategyAPIKey,
					))
				case StrategySession:
					// Session is usually part of JWT validation
					middlewares = append(middlewares, wrapWithStrategy(
						Middleware(authService, logger),
						StrategySession,
					))
				}
			}

			// If no strategies specified, use default (JWT)
			if len(middlewares) == 0 {
				middlewares = append(middlewares, wrapWithStrategy(
					Middleware(authService, logger),
					StrategyJWT,
				))
			}

			// Try authentication with composed strategies
			if config.Optional {
				// Use optional auth middleware
				err := ComposeOptionalStrategies(middlewares...)(next)(c)
				if err != nil {
					return err
				}
			} else {
				// Use required auth middleware
				err := ComposeAuthStrategies(middlewares...)(next)(c)
				if err != nil {
					return err
				}
				authenticated = true
			}

			// Check additional requirements if authenticated
			if authenticated || c.Get(string(AuthClaimsContextKey)) != nil {
				claims, ok := GetClaimsFromContext(c)
				if ok {
					// Check verified email requirement
					if config.RequireVerifiedEmail && !claims.EmailVerified {
						return echo.NewHTTPError(http.StatusForbidden, "email verification required")
					}

					// Check organization requirement
					if config.RequireOrganization && claims.OrganizationID.String() == "" {
						return echo.NewHTTPError(http.StatusForbidden, "organization membership required")
					}

					// Check required scopes
					for _, scope := range config.RequiredScopes {
						if !claims.HasScope(scope) {
							return echo.NewHTTPError(http.StatusForbidden, "insufficient scope")
						}
					}

					// Check required permissions
					for _, perm := range config.RequiredPermissions {
						if !claims.HasPermission(perm) {
							return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
						}
					}
				}

				// Get auth strategy from context
				if strategy, ok := c.Get(string(AuthStrategyContextKey)).(Strategy); ok {
					authStrategy = strategy
				}

				// Apply rate limiting based on auth strategy
				if config.RateLimit != nil {
					limit := config.RateLimit.RequestsPerMinute
					if strategyLimit, ok := config.RateLimit.ByStrategy[authStrategy]; ok {
						limit = strategyLimit
					}

					// Apply rate limiting (simplified - in production use Redis)
					if limit > 0 {
						// This is a placeholder - actual implementation would use Redis
						logger.Debug("applying rate limit",
							"strategy", authStrategy,
							"limit", limit,
						)
					}
				}
			}

			return nil
		}
	}
}

// wrapWithStrategy wraps a middleware to set the authentication strategy in context
func wrapWithStrategy(middleware echo.MiddlewareFunc, strategy Strategy) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Run the middleware
			err := middleware(func(c echo.Context) error {
				// Set the strategy in context if auth succeeds
				c.Set(string(AuthStrategyContextKey), strategy)
				ctx := context.WithValue(c.Request().Context(), AuthStrategyContextKey, strategy)
				c.SetRequest(c.Request().WithContext(ctx))
				return next(c)
			})(c)

			return err
		}
	}
}

// ComposeOptionalStrategies creates a middleware that tries multiple authentication strategies
// but doesn't fail if none succeed (optional auth)
func ComposeOptionalStrategies(strategies ...echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try each strategy, but don't fail if none work
			for _, strategy := range strategies {
				// Create a test handler
				testHandler := func(_ echo.Context) error {
					return nil
				}

				// Try the strategy
				err := strategy(testHandler)(c)
				if err == nil {
					// Success - auth worked
					break
				}
			}

			// Always continue to the next handler (auth is optional)
			return next(c)
		}
	}
}

// GetAuthStrategy retrieves the authentication strategy from context
func GetAuthStrategy(c echo.Context) (Strategy, bool) {
	strategy, ok := c.Get(string(AuthStrategyContextKey)).(Strategy)
	return strategy, ok
}

// MiddlewarePresets provides common middleware configurations
var MiddlewarePresets = struct {
	// Public allows access without authentication
	Public RequireAuthConfig
	// Authenticated requires any valid authentication
	Authenticated RequireAuthConfig
	// APIOnly requires API key authentication
	APIOnly RequireAuthConfig
	// AdminOnly requires admin role
	AdminOnly RequireAuthConfig
	// OrganizationMember requires organization membership
	OrganizationMember RequireAuthConfig
}{
	Public: RequireAuthConfig{
		Optional: true,
		Strategies: []Strategy{
			StrategyJWT,
			StrategyAPIKey,
		},
	},
	Authenticated: RequireAuthConfig{
		Optional: false,
		Strategies: []Strategy{
			StrategyJWT,
			StrategySession,
			StrategyAPIKey,
		},
	},
	APIOnly: RequireAuthConfig{
		Optional: false,
		Strategies: []Strategy{
			StrategyAPIKey,
		},
		RateLimit: &RateLimitConfig{
			RequestsPerMinute: 100,
			BurstSize:         10,
		},
	},
	AdminOnly: RequireAuthConfig{
		Optional: false,
		Strategies: []Strategy{
			StrategyJWT,
			StrategySession,
		},
		RequiredPermissions: []string{
			"system:manage",
		},
	},
	OrganizationMember: RequireAuthConfig{
		Optional:            false,
		RequireOrganization: true,
		Strategies: []Strategy{
			StrategyJWT,
			StrategySession,
		},
		RequiredPermissions: []string{
			"org:read",
		},
	},
}

// ChainMiddleware chains multiple middleware functions together
func ChainMiddleware(middlewares ...echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Build the chain in reverse order
		handler := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

// ConditionalMiddleware applies middleware based on a condition
func ConditionalMiddleware(condition func(c echo.Context) bool, middleware echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if condition(c) {
				return middleware(next)(c)
			}
			return next(c)
		}
	}
}

// PathAuthConfig applies different authentication based on path patterns
type PathAuthConfig struct {
	Pattern string
	Config  RequireAuthConfig
}

// matchPath performs simple path pattern matching
// Supports wildcards: * for single segment, ** for multiple segments
func matchPath(pattern, path string) bool {
	// Exact match
	if pattern == path {
		return true
	}

	// Check for wildcards
	if strings.Contains(pattern, "*") {
		// Simple wildcard matching
		if strings.HasSuffix(pattern, "/*") {
			// Match any path under this prefix
			prefix := strings.TrimSuffix(pattern, "/*")
			return strings.HasPrefix(path, prefix+"/")
		}
		if pattern == "/*" {
			// Match any single-segment path
			return strings.Count(path, "/") == 1
		}
	}

	// Check if pattern is a prefix
	if strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, pattern) {
		return true
	}

	return false
}

// PathBasedAuthMiddleware applies different auth configurations based on path
func PathBasedAuthMiddleware(authService *Service, configs []PathAuthConfig, logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()

			// Find matching config for this path
			for _, config := range configs {
				// Simple path matching (can be enhanced with proper pattern matching)
				if matchPath(config.Pattern, path) {
					// Apply the auth config for this path
					return RequireAuth(authService, config.Config, logger)(next)(c)
				}
			}

			// No specific config found, use default authenticated requirement
			return RequireAuth(authService, MiddlewarePresets.Authenticated, logger)(next)(c)
		}
	}
}

// ErrorHandlingMiddleware wraps auth errors with more context
func ErrorHandlingMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				// Log authentication errors with context
				if httpErr, ok := err.(*echo.HTTPError); ok {
					if httpErr.Code == http.StatusUnauthorized || httpErr.Code == http.StatusForbidden {
						strategy, _ := GetAuthStrategy(c)
						logger.Warn("authentication failed",
							"path", c.Path(),
							"method", c.Request().Method,
							"strategy", strategy,
							"error", err,
							"ip", c.RealIP(),
						)
					}
				}
			}
			return err
		}
	}
}
