package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/domains/auth/entities"
	"github.com/archesai/archesai/internal/domains/auth/services"
	"github.com/labstack/echo/v4"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the context key for the authenticated user
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the context key for JWT claims
	ClaimsContextKey ContextKey = "claims"
)

// Middleware creates an authentication middleware
func Middleware(authService *services.Service, logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			token := extractToken(c)
			if token == "" {
				// Check for session cookie
				cookie, err := c.Cookie("session_token")
				if err == nil {
					token = cookie.Value
				}
			}

			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication token")
			}

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				logger.Warn("invalid token", "error", err)
				if err == entities.ErrTokenExpired {
					return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Set claims in context
			c.Set(string(ClaimsContextKey), claims)
			c.Set(string(UserContextKey), claims.UserID)

			// Add user info to request context for downstream use
			ctx := context.WithValue(c.Request().Context(), ClaimsContextKey, claims)
			ctx = context.WithValue(ctx, UserContextKey, claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
// It validates the token if present but doesn't require it
func OptionalAuthMiddleware(authService *services.Service, logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			token := extractToken(c)
			if token == "" {
				// Check for session cookie
				cookie, err := c.Cookie("session_token")
				if err == nil {
					token = cookie.Value
				}
			}

			if token != "" {
				// Validate token if present
				claims, err := authService.ValidateToken(token)
				if err == nil {
					// Set claims in context if valid
					c.Set(string(ClaimsContextKey), claims)
					c.Set(string(UserContextKey), claims.UserID)

					// Add user info to request context
					ctx := context.WithValue(c.Request().Context(), ClaimsContextKey, claims)
					ctx = context.WithValue(ctx, UserContextKey, claims.UserID)
					c.SetRequest(c.Request().WithContext(ctx))
				} else {
					logger.Debug("invalid optional token", "error", err)
				}
			}

			return next(c)
		}
	}
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(ClaimsContextKey)).(*entities.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			return next(c)
		}
	}
}

// RequireOrganization creates a middleware that requires organization membership
func RequireOrganization() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(ClaimsContextKey)).(*entities.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			if claims.OrganizationID == nil {
				return echo.NewHTTPError(http.StatusForbidden, "organization membership required")
			}

			return next(c)
		}
	}
}

// extractToken extracts the token from the Authorization header
func extractToken(c echo.Context) string {
	// Get token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		// Check for Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Get token from query parameter (useful for WebSocket connections)
	if token := c.QueryParam("token"); token != "" {
		return token
	}

	return ""
}

// GetUserFromContext retrieves the user ID from the context
func GetUserFromContext(c echo.Context) (string, bool) {
	userID, ok := c.Get(string(UserContextKey)).(string)
	return userID, ok
}

// GetClaimsFromContext retrieves the claims from the context
func GetClaimsFromContext(c echo.Context) (*entities.Claims, bool) {
	claims, ok := c.Get(string(ClaimsContextKey)).(*entities.Claims)
	return claims, ok
}

// RateLimitMiddleware creates a rate limiting middleware for authentication endpoints
func RateLimitMiddleware(maxAttempts int, windowMinutes int) echo.MiddlewareFunc {
	// This is a simplified version. In production, use a Redis-based solution
	attempts := make(map[string]int)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			// Check attempts
			if attempts[ip] >= maxAttempts {
				return echo.NewHTTPError(http.StatusTooManyRequests, "too many authentication attempts")
			}

			// Increment attempts
			attempts[ip]++

			// Continue with request
			err := next(c)

			// Reset on successful authentication
			if err == nil && c.Response().Status == http.StatusOK {
				delete(attempts, ip)
			}

			return err
		}
	}
}

// SetRequestContextWithTimeout will set the request context with timeout for every incoming HTTP Request
func SetRequestContextWithTimeout(d time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), d)
			defer cancel()

			newRequest := c.Request().WithContext(ctx)
			c.SetRequest(newRequest)
			return next(c)
		}
	}
}
