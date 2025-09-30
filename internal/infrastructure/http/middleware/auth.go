// Package middleware provides HTTP middleware for authentication and authorization
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/infrastructure/config"
)

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

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"invalid authorization header format",
				)
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.ParseWithClaims(
				tokenString,
				&Claims{},
				func(token *jwt.Token) (interface{}, error) {
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

// OptionalAuthMiddleware creates optional JWT authentication middleware
func OptionalAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// No auth header, continue without authentication
				return next(c)
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format, continue without authentication
				return next(c)
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.ParseWithClaims(
				tokenString,
				&Claims{},
				func(token *jwt.Token) (interface{}, error) {
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
