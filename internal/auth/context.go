package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

// Context keys for authentication data
const (
	// AuthClaimsContextKey is the context key for storing JWT claims
	AuthClaimsContextKey ContextKey = "auth.claims"

	// AuthUserContextKey is the context key for storing user ID
	AuthUserContextKey ContextKey = "auth.user_id"

	// UserContextKey is the context key for the authenticated user
	UserContextKey ContextKey = "user"

	// SessionTokenContextKey is the context key for session token
	SessionTokenContextKey ContextKey = "session_token"

	// OrganizationContextKey is the context key for storing organization data
	OrganizationContextKey ContextKey = "auth.organization"

	// AuthStrategyContextKey is the context key for storing the authentication strategy used
	AuthStrategyContextKey ContextKey = "auth.strategy"
)

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

// GetAuthStrategy retrieves the authentication strategy from context
func GetAuthStrategy(c echo.Context) (Strategy, bool) {
	strategy, ok := c.Get(string(AuthStrategyContextKey)).(Strategy)
	return strategy, ok
}

// GetOrganizationFromContext retrieves the organization from the context
func GetOrganizationFromContext(c echo.Context) (interface{}, bool) {
	org := c.Get(string(OrganizationContextKey))
	return org, org != nil
}

// GetSessionTokenFromContext retrieves the session token from the context
func GetSessionTokenFromContext(c echo.Context) (string, bool) {
	token, ok := c.Get(string(SessionTokenContextKey)).(string)
	return token, ok
}

// GetUserEntityFromContext retrieves the full user entity from the context
func GetUserEntityFromContext(c echo.Context) (*User, bool) {
	user, ok := c.Get(string(UserContextKey)).(*User)
	return user, ok
}

// SetAuthContext sets all authentication-related values in both Echo and Go contexts
func SetAuthContext(c echo.Context, claims *EnhancedClaims, user *User, token string) {
	// Set in Echo context
	c.Set(string(AuthClaimsContextKey), claims)
	c.Set(string(AuthUserContextKey), claims.UserID)
	c.Set(string(UserContextKey), user)
	c.Set(string(SessionTokenContextKey), token)

	// Set in Go context for downstream use
	ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
	ctx = context.WithValue(ctx, AuthUserContextKey, claims.UserID)
	ctx = context.WithValue(ctx, UserContextKey, user)
	ctx = context.WithValue(ctx, SessionTokenContextKey, token)
	c.SetRequest(c.Request().WithContext(ctx))
}

// SetOrganizationContext sets organization data in both Echo and Go contexts
func SetOrganizationContext(c echo.Context, org interface{}) {
	c.Set(string(OrganizationContextKey), org)
	ctx := context.WithValue(c.Request().Context(), OrganizationContextKey, org)
	c.SetRequest(c.Request().WithContext(ctx))
}

// SetAuthStrategyContext sets the authentication strategy in both Echo and Go contexts
func SetAuthStrategyContext(c echo.Context, strategy Strategy) {
	c.Set(string(AuthStrategyContextKey), strategy)
	ctx := context.WithValue(c.Request().Context(), AuthStrategyContextKey, strategy)
	c.SetRequest(c.Request().WithContext(ctx))
}

// GetAuthContextFromGoContext retrieves authentication data from standard Go context
// This is useful for service and repository layers that don't have access to Echo context
func GetAuthContextFromGoContext(ctx context.Context) (*EnhancedClaims, uuid.UUID, bool) {
	claims, claimsOk := ctx.Value(AuthClaimsContextKey).(*EnhancedClaims)
	userID, userOk := ctx.Value(AuthUserContextKey).(uuid.UUID)
	return claims, userID, claimsOk && userOk
}
