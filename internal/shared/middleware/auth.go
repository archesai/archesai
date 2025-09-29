// Package middleware provides authentication and authorization middleware functions.
// It includes JWT validation, context management, and role-based access control.
package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/infrastructure/config"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// AuthUserContextKey is the context key for user authentication.
const AuthUserContextKey contextKey = "auth.user"

// AuthClaimsContextKey is the context key for auth claims.
const AuthClaimsContextKey contextKey = "auth.claims"

// Claims represents JWT claims with user information.
type Claims struct {
	UserID             uuid.UUID `json:"user_id"`
	Email              string    `json:"email"`
	OrganizationID     uuid.UUID `json:"organization_id,omitempty"`
	Role               string    `json:"role,omitempty"`
	SessionID          uuid.UUID `json:"session_id,omitempty"`
	IsEmailVerified    bool      `json:"email_verified"`
	RequiresOnboarding bool      `json:"requires_onboarding"`
	jwt.RegisteredClaims
}

// User represents the authenticated user context.
type User struct {
	ID                 uuid.UUID `json:"id"`
	Email              string    `json:"email"`
	OrganizationID     uuid.UUID `json:"organization_id,omitempty"`
	Role               string    `json:"role,omitempty"`
	SessionID          uuid.UUID `json:"session_id,omitempty"`
	IsEmailVerified    bool      `json:"email_verified"`
	RequiresOnboarding bool      `json:"requires_onboarding"`
}

// AuthMiddleware provides authentication middleware functionality.
// This is a bridge to the unified auth service middleware.
type AuthMiddleware struct {
	jwtSecret []byte
	logger    *slog.Logger
	config    *config.Config
	authSvc   interface{ Middleware() echo.MiddlewareFunc } // Auth service with middleware method
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(jwtSecret string, cfg *config.Config, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
		config:    cfg,
	}
}

// SetAuthService sets the auth service for unified authentication.
func (am *AuthMiddleware) SetAuthService(authSvc interface{ Middleware() echo.MiddlewareFunc }) {
	am.authSvc = authSvc
}

// RequireAuth creates middleware that validates JWT tokens.
func (am *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	// Use unified auth service if available
	if am.authSvc != nil {
		return am.authSvc.Middleware()
	}

	// Fallback to legacy JWT validation
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

			// Parse and validate JWT token
			token, err := jwt.ParseWithClaims(
				tokenString,
				&Claims{},
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return am.jwtSecret, nil
				},
			)

			if err != nil {
				am.logger.Warn("JWT validation failed", "error", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			// Check token expiration
			if time.Now().After(claims.ExpiresAt.Time) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
			}

			// Create user context
			user := &User{
				ID:                 claims.UserID,
				Email:              claims.Email,
				OrganizationID:     claims.OrganizationID,
				Role:               claims.Role,
				SessionID:          claims.SessionID,
				IsEmailVerified:    claims.IsEmailVerified,
				RequiresOnboarding: claims.RequiresOnboarding,
			}

			// Add user and claims to context
			ctx := context.WithValue(c.Request().Context(), AuthUserContextKey, user)
			ctx = context.WithValue(ctx, AuthClaimsContextKey, claims)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// GetUserFromContext retrieves the authenticated user from the request context.
func GetUserFromContext(c echo.Context) (*User, bool) {
	user, ok := c.Request().Context().Value(AuthUserContextKey).(*User)
	return user, ok
}

// GetClaimsFromContext retrieves the JWT claims from the request context.
func GetClaimsFromContext(c echo.Context) (*Claims, bool) {
	claims, ok := c.Request().Context().Value(AuthClaimsContextKey).(*Claims)
	return claims, ok
}

// RequireEmailVerified creates middleware that requires email verification.
func (am *AuthMiddleware) RequireEmailVerified() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := GetUserFromContext(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}

			if !user.IsEmailVerified {
				return echo.NewHTTPError(http.StatusForbidden, "email verification required")
			}

			return next(c)
		}
	}
}

// RequireOrganizationMember creates middleware that requires organization membership.
func (am *AuthMiddleware) RequireOrganizationMember() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := GetUserFromContext(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}

			if user.OrganizationID == uuid.Nil {
				return echo.NewHTTPError(http.StatusForbidden, "organization membership required")
			}

			return next(c)
		}
	}
}

// RequireRole creates middleware that requires a specific role.
func (am *AuthMiddleware) RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := GetUserFromContext(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}

			if user.Role != role {
				return echo.NewHTTPError(
					http.StatusForbidden,
					fmt.Sprintf("role '%s' required", role),
				)
			}

			return next(c)
		}
	}
}
