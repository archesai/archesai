package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/auth"
	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	// AuthClaimsContextKey is the context key for storing JWT claims
	AuthClaimsContextKey contextKey = "auth.claims"
	// AuthUserContextKey is the context key for storing user ID
	AuthUserContextKey contextKey = "auth.user_id"
)

// Middleware creates an authentication middleware
func Middleware(authService *auth.Service, logger *slog.Logger) echo.MiddlewareFunc {
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
				if err == auth.ErrTokenExpired {
					return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Set claims in context
			c.Set(string(AuthClaimsContextKey), claims)
			c.Set(string(AuthUserContextKey), claims.UserID)

			// Add user info to request context for downstream use
			ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
			ctx = context.WithValue(ctx, AuthUserContextKey, claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
// It validates the token if present but doesn't require it
func OptionalAuthMiddleware(authService *auth.Service, logger *slog.Logger) echo.MiddlewareFunc {
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
					c.Set(string(AuthClaimsContextKey), claims)
					c.Set(string(AuthUserContextKey), claims.UserID)

					// Add user info to request context
					ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
					ctx = context.WithValue(ctx, AuthUserContextKey, claims.UserID)
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
			claims, ok := c.Get(string(AuthClaimsContextKey)).(*auth.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				// FIX ME
				if claims.Email == role {
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
			claims, ok := c.Get(string(AuthClaimsContextKey)).(*auth.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// FIXME
			if claims.Audience == nil {
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
	userID, ok := c.Get(string(AuthUserContextKey)).(string)
	return userID, ok
}

// GetClaimsFromContext retrieves the claims from the context
func GetClaimsFromContext(c echo.Context) (*auth.Claims, bool) {
	claims, ok := c.Get(string(AuthClaimsContextKey)).(*auth.Claims)
	return claims, ok
}

// RateLimitMiddleware creates a rate limiting middleware for authentication endpoints
func RateLimitMiddleware(maxAttempts int, _ int) echo.MiddlewareFunc {
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
