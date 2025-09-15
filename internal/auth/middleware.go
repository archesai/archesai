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
			ctx = context.WithValue(ctx, ContextKey("api_key"), apiKeyData)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

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
