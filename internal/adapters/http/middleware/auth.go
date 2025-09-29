// Package middleware provides HTTP-specific middleware for the server.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// Context keys for authentication.
const (
	AuthUserContextKey   contextKey = "auth.user"
	AuthClaimsContextKey contextKey = "auth.claims"
	AuthMethodContextKey contextKey = "auth.method"
)

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
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
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
			c.SetRequest(c.Request().WithContext(ctx))

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
			if len(parts) != 2 || parts[0] != "Bearer" {
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
				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}

// GetUserFromContext retrieves the user ID from the context.
func GetUserFromContext(c echo.Context) (uuid.UUID, bool) {
	userID, ok := c.Request().Context().Value(AuthUserContextKey).(uuid.UUID)
	return userID, ok
}

// GetClaimsFromContext retrieves the token claims from the context.
func GetClaimsFromContext(c echo.Context) (*auth.TokenClaims, bool) {
	claims, ok := c.Request().Context().Value(AuthClaimsContextKey).(*auth.TokenClaims)
	return claims, ok
}
