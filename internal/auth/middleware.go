package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

			// If session ID is present in claims, validate session with Redis
			if claims.SessionID != "" && authService.sessionManager != nil {
				// ValidateSession actually validates by token, not session ID
				// So we use the token that was already extracted
				session, err := authService.sessionManager.ValidateSession(c.Request().Context(), token)
				if err != nil {
					logger.Warn("invalid session", "session_id", claims.SessionID, "error", err)
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
				}
				// Session activity is automatically updated by ValidateSession
				c.Set("session", session)
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

// OrganizationContextKey is the context key for storing organization data
const OrganizationContextKey contextKey = "auth.organization"

// EnrichOrganizationContext enriches the context with full organization data
func EnrichOrganizationContext(getOrgFunc func(ctx context.Context, orgID uuid.UUID) (interface{}, error), logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get claims from context
			claims, ok := GetClaimsFromContext(c)
			if !ok || claims.OrganizationID == uuid.Nil {
				// No organization to enrich
				return next(c)
			}

			// Check if organization is already in context (from cache)
			if org := c.Get(string(OrganizationContextKey)); org != nil {
				return next(c)
			}

			// Load organization data
			org, err := getOrgFunc(c.Request().Context(), claims.OrganizationID)
			if err != nil {
				logger.Warn("failed to load organization",
					"org_id", claims.OrganizationID,
					"error", err,
				)
				// Continue without organization data
				return next(c)
			}

			// Set organization in context
			c.Set(string(OrganizationContextKey), org)
			ctx := context.WithValue(c.Request().Context(), OrganizationContextKey, org)
			c.SetRequest(c.Request().WithContext(ctx))

			// Enrich claims with organization-specific permissions if applicable
			if claims.OrganizationRole != "" {
				// Add role-based permissions for the organization
				rolePerms := GetRolePermissions(Role(claims.OrganizationRole))
				for _, perm := range rolePerms {
					if !claims.HasPermission(string(perm)) {
						claims.Permissions = append(claims.Permissions, string(perm))
					}
				}
				// Update claims in context
				c.Set(string(AuthClaimsContextKey), claims)
			}

			return next(c)
		}
	}
}

// RequireOrganizationRole creates a middleware that requires a specific organization role
func RequireOrganizationRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication")
			}

			// Check if user has required organization role
			hasRole := false
			for _, role := range roles {
				if claims.OrganizationRole == role {
					hasRole = true
					break
				}
				// Also check in Organizations list
				for _, org := range claims.Organizations {
					if org.Role == role {
						hasRole = true
						break
					}
				}
			}

			if !hasRole {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient organization role")
			}

			return next(c)
		}
	}
}

// SwitchOrganization allows switching the active organization via header
func SwitchOrganization(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for X-Organization-ID header
			orgIDStr := c.Request().Header.Get("X-Organization-ID")
			if orgIDStr == "" {
				return next(c)
			}

			// Parse organization ID
			orgID, err := uuid.Parse(orgIDStr)
			if err != nil {
				logger.Warn("invalid organization ID in header",
					"org_id", orgIDStr,
					"error", err,
				)
				return echo.NewHTTPError(http.StatusBadRequest, "invalid organization ID")
			}

			// Get claims
			claims, ok := GetClaimsFromContext(c)
			if !ok {
				return next(c)
			}

			// Verify user is a member of this organization
			if !claims.IsOrgMember(orgID) {
				return echo.NewHTTPError(http.StatusForbidden, "not a member of this organization")
			}

			// Switch active organization in claims
			originalOrgID := claims.OrganizationID
			claims.OrganizationID = orgID

			// Find and set the role for this organization
			for _, org := range claims.Organizations {
				if org.ID == orgID {
					claims.OrganizationName = org.Name
					claims.OrganizationRole = org.Role
					break
				}
			}

			// Update claims in context
			c.Set(string(AuthClaimsContextKey), claims)
			ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
			c.SetRequest(c.Request().WithContext(ctx))

			logger.Debug("switched organization context",
				"user_id", claims.UserID,
				"from_org", originalOrgID,
				"to_org", orgID,
			)

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

// APIKeyMiddleware creates an API key authentication middleware
func APIKeyMiddleware(authService *Service, logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for API key in X-API-Key header
			apiKey := c.Request().Header.Get("X-API-Key")

			// If not found, check Authorization header
			if apiKey == "" {
				authHeader := c.Request().Header.Get("Authorization")
				apiKey = ParseAPIKey(authHeader)
			}

			if apiKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing API key")
			}

			// Validate API key through service
			apiKeyData, err := authService.ValidateAPIKey(c.Request().Context(), apiKey)
			if err != nil {
				logger.Warn("invalid API key", "error", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
			}

			// Create claims from API key
			claims := &EnhancedClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   apiKeyData.UserID.String(),
					ExpiresAt: jwt.NewNumericDate(apiKeyData.ExpiresAt),
					IssuedAt:  jwt.NewNumericDate(apiKeyData.CreatedAt),
					NotBefore: jwt.NewNumericDate(apiKeyData.CreatedAt),
					Issuer:    "archesai",
					ID:        apiKeyData.ID.String(),
				},
				UserID:         apiKeyData.UserID,
				OrganizationID: apiKeyData.OrganizationID,
				TokenType:      APIKeyTokenType,
				AuthMethod:     AuthMethodAPIKey,
				Scopes:         apiKeyData.Scopes,
				CustomClaims: map[string]interface{}{
					"api_key_id":   apiKeyData.ID.String(),
					"api_key_name": apiKeyData.Name,
					"rate_limit":   apiKeyData.RateLimit,
				},
			}

			// Set claims and user in context
			c.Set(string(AuthClaimsContextKey), claims)
			c.Set(string(AuthUserContextKey), apiKeyData.UserID)
			c.Set("api_key", apiKeyData)

			// Add to request context for downstream use
			ctx := context.WithValue(c.Request().Context(), AuthClaimsContextKey, claims)
			ctx = context.WithValue(ctx, AuthUserContextKey, apiKeyData.UserID)
			ctx = context.WithValue(ctx, contextKey("api_key"), apiKeyData)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
