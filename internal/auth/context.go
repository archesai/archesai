package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// Context keys for storing authentication information.
const (
	// AuthUserContextKey is the context key for user authentication.
	AuthUserContextKey contextKey = "auth.user"
	// AuthClaimsContextKey is the context key for auth claims.
	AuthClaimsContextKey contextKey = "auth.claims"
	// AuthUserIDContextKey is the context key for user ID only.
	AuthUserIDContextKey contextKey = "auth.user_id"
)

// GetAuthContextFromGoContext retrieves authentication data from standard Go context
// This is useful for service and repository layers that don't have access to Echo context.
func GetAuthContextFromGoContext(ctx context.Context) (*EnhancedClaims, uuid.UUID, bool) {
	// Try to get enhanced claims first
	if claims, ok := ctx.Value(AuthClaimsContextKey).(*EnhancedClaims); ok {
		return claims, claims.UserID, true
	}

	// Fallback to getting user ID directly
	if userID, ok := ctx.Value(AuthUserIDContextKey).(uuid.UUID); ok {
		return nil, userID, true
	}

	return nil, uuid.Nil, false
}

// GetClaimsFromContext retrieves enhanced claims from context.
func GetClaimsFromContext(ctx context.Context) (*EnhancedClaims, bool) {
	claims, ok := ctx.Value(AuthClaimsContextKey).(*EnhancedClaims)
	return claims, ok
}

// GetUserIDFromContext retrieves user ID from context.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	// Try claims first
	if claims, ok := ctx.Value(AuthClaimsContextKey).(*EnhancedClaims); ok {
		return claims.UserID, true
	}

	// Try direct user ID
	if userID, ok := ctx.Value(AuthUserIDContextKey).(uuid.UUID); ok {
		return userID, true
	}

	return uuid.Nil, false
}

// SetClaimsInContext adds enhanced claims to context.
func SetClaimsInContext(ctx context.Context, claims *EnhancedClaims) context.Context {
	ctx = context.WithValue(ctx, AuthClaimsContextKey, claims)
	if claims != nil {
		ctx = context.WithValue(ctx, AuthUserIDContextKey, claims.UserID)
	}
	return ctx
}

// GetClaimsFromEchoContext retrieves enhanced claims from Echo context.
func GetClaimsFromEchoContext(c echo.Context) (*EnhancedClaims, bool) {
	claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
	return claims, ok
}

// GetUserIDFromEchoContext retrieves user ID from Echo context.
func GetUserIDFromEchoContext(c echo.Context) (uuid.UUID, bool) {
	// Try claims first
	if claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
		return claims.UserID, true
	}

	// Try direct user ID
	if userID, ok := c.Get(string(AuthUserIDContextKey)).(uuid.UUID); ok {
		return userID, true
	}

	return uuid.Nil, false
}

// SetClaimsInEchoContext adds enhanced claims to Echo context.
func SetClaimsInEchoContext(c echo.Context, claims *EnhancedClaims) {
	c.Set(string(AuthClaimsContextKey), claims)
	if claims != nil {
		c.Set(string(AuthUserIDContextKey), claims.UserID)
	}
	// Also set in request context for downstream services
	ctx := SetClaimsInContext(c.Request().Context(), claims)
	c.SetRequest(c.Request().WithContext(ctx))
}
