package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
func Middleware(authService *Service, logger *slog.Logger) echo.MiddlewareFunc {
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
				if err == ErrTokenExpired {
					return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Verify user exists
			user, err := authService.GetUserByID(c.Request().Context(), claims.UserID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			// Set claims, user, and session token in context
			c.Set(string(AuthClaimsContextKey), claims)
			c.Set(string(AuthUserContextKey), claims.UserID)
			c.Set(string(UserContextKey), user)
			c.Set(string(SessionTokenContextKey), token) // Add session token to context

			// Add user info to request context for downstream use
			ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
			ctx = context.WithValue(ctx, AuthUserContextKey, claims.UserID)
			ctx = context.WithValue(ctx, SessionTokenContextKey, token) // Add session token to request context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
// It validates the token if present but doesn't require it
func OptionalAuthMiddleware(authService *Service, logger *slog.Logger) echo.MiddlewareFunc {
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
					// Get user if token is valid
					user, userErr := authService.GetUserByID(c.Request().Context(), claims.UserID)
					if userErr == nil {
						// Set claims, user, and session token in context if valid
						c.Set(string(AuthClaimsContextKey), claims)
						c.Set(string(AuthUserContextKey), claims.UserID)
						c.Set(string(UserContextKey), user)
						c.Set(string(SessionTokenContextKey), token) // Add session token to context

						// Add user info to request context
						ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
						ctx = context.WithValue(ctx, AuthUserContextKey, claims.UserID)
						ctx = context.WithValue(ctx, UserContextKey, user)
						ctx = context.WithValue(ctx, SessionTokenContextKey, token) // Add session token to request context
						c.SetRequest(c.Request().WithContext(ctx))
					} else {
						logger.Debug("user not found for valid token", "error", userErr)
					}
				} else {
					logger.Debug("invalid optional token", "error", err)
				}
			}

			return next(c)
		}
	}
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(roles ...Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if claims.HasRole(string(role)) {
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

// RequirePermission creates a middleware that requires specific permissions
func RequirePermission(permissions ...Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// Check if user has all required permissions
			for _, permission := range permissions {
				if !claims.HasPermission(string(permission)) {
					return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
				}
			}

			return next(c)
		}
	}
}

// RequireScope creates a middleware that requires specific API scopes
func RequireScope(scopes ...Scope) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// Check if user has all required scopes
			for _, scope := range scopes {
				if !claims.HasScope(string(scope)) {
					return echo.NewHTTPError(http.StatusForbidden, "insufficient scope")
				}
			}

			return next(c)
		}
	}
}

// RequireOrganization creates a middleware that requires organization membership
func RequireOrganization() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// Check if user has an active organization
			if claims.OrganizationID == uuid.Nil && len(claims.Organizations) == 0 {
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
func GetUserFromContext(c echo.Context) (uuid.UUID, bool) {
	userID, ok := c.Get(string(AuthUserContextKey)).(uuid.UUID)
	return userID, ok
}

// GetClaimsFromContext retrieves the enhanced claims from the context
func GetClaimsFromContext(c echo.Context) (*EnhancedClaims, bool) {
	claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
	return claims, ok
}

// GetLegacyClaimsFromContext retrieves legacy claims from the context (for backward compatibility)
func GetLegacyClaimsFromContext(c echo.Context) (*Claims, bool) {
	// Try to get enhanced claims first and convert
	if enhanced, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
		legacy := &Claims{
			UserID:           enhanced.UserID,
			Email:            enhanced.Email,
			RegisteredClaims: enhanced.RegisteredClaims,
		}
		return legacy, true
	}
	// Fall back to direct legacy claims
	claims, ok := c.Get(string(AuthClaimsContextKey)).(*Claims)
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
