// Package auth provides authentication and authorization services
package auth

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Additional context keys for middleware
const (
	// AuthAPITokenContextKey is the context key for API token
	AuthAPITokenContextKey = "auth_api_token"
	// AuthMethodContextKey is the context key for the auth method used
	AuthMethodContextKey = "auth_method"
)

// Middleware creates a unified authentication middleware.
func (s *Service) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try JWT authentication first
			if authHeader := c.Request().Header.Get("Authorization"); authHeader != "" {
				if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
					token := strings.TrimSpace(authHeader[7:])

					// Check if it's an API key format
					if strings.HasPrefix(token, "sk_") {
						return s.authenticateAPIKey(c, token, next)
					}

					// Try JWT authentication
					return s.authenticateJWT(c, token, next)
				}
			}

			// Try API key authentication from other headers
			if apiKey := s.extractAPIKey(c); apiKey != "" {
				return s.authenticateAPIKey(c, apiKey, next)
			}

			// No authentication provided
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}
	}
}

// OptionalMiddleware creates middleware that allows unauthenticated requests.
func (s *Service) OptionalMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try authentication but don't fail if not provided
			if authHeader := c.Request().Header.Get("Authorization"); authHeader != "" {
				if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
					token := strings.TrimSpace(authHeader[7:])

					if strings.HasPrefix(token, "sk_") {
						_ = s.authenticateAPIKey(c, token, func(c echo.Context) error {
							return next(c)
						})
						return nil
					}

					_ = s.authenticateJWT(c, token, func(c echo.Context) error {
						return next(c)
					})
					return nil
				}
			}

			if apiKey := s.extractAPIKey(c); apiKey != "" {
				_ = s.authenticateAPIKey(c, apiKey, func(c echo.Context) error {
					return next(c)
				})
				return nil
			}

			return next(c)
		}
	}
}

// RequireEmailVerified ensures the user's email is verified.
func RequireEmailVerified() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check JWT claims
			if claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
				if !claims.EmailVerified {
					return echo.NewHTTPError(http.StatusForbidden, "email verification required")
				}
			}

			return next(c)
		}
	}
}

// RequireOrganizationMember ensures the user is a member of the specified organization.
func RequireOrganizationMember(orgID uuid.UUID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check JWT claims
			if claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
				if !claims.IsOrgMember(orgID) {
					return echo.NewHTTPError(
						http.StatusForbidden,
						"organization membership required",
					)
				}
			}

			// Check API token
			if token, ok := c.Get(AuthAPITokenContextKey).(*APIToken); ok {
				if token.OrganizationID != orgID {
					return echo.NewHTTPError(http.StatusForbidden, "organization access required")
				}
			}

			return next(c)
		}
	}
}

// RequireRole ensures the user has the specified role.
func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check JWT claims
			if claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
				if !claims.HasRole(role) {
					return echo.NewHTTPError(http.StatusForbidden, "insufficient role")
				}
			}

			return next(c)
		}
	}
}

// RequireScopes ensures the API token has the required scopes.
func (s *Service) RequireScopes(scopes ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check API token scopes
			if token, ok := c.Get(AuthAPITokenContextKey).(*APIToken); ok {
				if err := s.apiTokenValidator.ValidateScopes(token.Scopes, scopes); err != nil {
					return echo.NewHTTPError(http.StatusForbidden, "insufficient scopes")
				}
			}

			// Check JWT scopes
			if claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
				for _, scope := range scopes {
					if !claims.HasScope(scope) {
						return echo.NewHTTPError(http.StatusForbidden, "insufficient scopes")
					}
				}
			}

			return next(c)
		}
	}
}

// authenticateJWT handles JWT token authentication.
func (s *Service) authenticateJWT(c echo.Context, token string, next echo.HandlerFunc) error {
	// Validate JWT token
	claims, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	// Get user information (you may want to cache this)
	user, err := s.usersService.GetByID(c.Request().Context(), claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	// Set context values
	c.Set(string(AuthUserContextKey), user)
	c.Set(string(AuthClaimsContextKey), claims)
	c.Set(AuthMethodContextKey, "jwt")

	return next(c)
}

// authenticateAPIKey handles API key authentication.
func (s *Service) authenticateAPIKey(c echo.Context, apiKey string, next echo.HandlerFunc) error {
	// Validate API key
	token, err := s.apiTokenStore.ValidateToken(c.Request().Context(), apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid api key")
	}

	// Get user information
	user, err := s.usersService.GetByID(c.Request().Context(), token.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	// Set context values
	c.Set(string(AuthUserContextKey), user)
	c.Set(AuthAPITokenContextKey, token)
	c.Set(AuthMethodContextKey, "api_key")

	return next(c)
}

// extractAPIKey extracts API key from various headers.
func (s *Service) extractAPIKey(c echo.Context) string {
	headers := map[string]string{
		"Authorization": c.Request().Header.Get("Authorization"),
		"X-API-Key":     c.Request().Header.Get("X-API-Key"),
		"Api-Key":       c.Request().Header.Get("Api-Key"),
	}

	return s.apiTokenValidator.ExtractAPIKeyFromHeaders(headers)
}

// GetAuthenticatedUser returns the authenticated user from context.
func GetAuthenticatedUser(c echo.Context) interface{} {
	return c.Get(string(AuthUserContextKey))
}

// GetAuthenticatedUserID returns the authenticated user ID from context.
func GetAuthenticatedUserID(c echo.Context) (uuid.UUID, bool) {
	// Try JWT claims first
	if claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
		return claims.UserID, true
	}

	// Try API token
	if token, ok := c.Get(AuthAPITokenContextKey).(*APIToken); ok {
		return token.UserID, true
	}

	return uuid.Nil, false
}

// GetJWTClaims returns JWT claims from context.
func GetJWTClaims(c echo.Context) (*EnhancedClaims, bool) {
	claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims)
	return claims, ok
}

// GetAPIToken returns API token from context.
func GetAPIToken(c echo.Context) (*APIToken, bool) {
	token, ok := c.Get(AuthAPITokenContextKey).(*APIToken)
	return token, ok
}

// GetAuthMethod returns the authentication method used.
func GetAuthMethod(c echo.Context) string {
	if method, ok := c.Get(AuthMethodContextKey).(string); ok {
		return method
	}
	return ""
}

// IsAuthenticated checks if the request is authenticated.
func IsAuthenticated(c echo.Context) bool {
	_, hasUser := GetAuthenticatedUserID(c)
	return hasUser
}

// GetOrganizationID returns the organization ID from the authenticated context.
func GetOrganizationID(c echo.Context) (uuid.UUID, bool) {
	// Try JWT claims first
	if claims, ok := c.Get(string(AuthClaimsContextKey)).(*EnhancedClaims); ok {
		if claims.OrganizationID != uuid.Nil {
			return claims.OrganizationID, true
		}
	}

	// Try API token
	if token, ok := c.Get(AuthAPITokenContextKey).(*APIToken); ok {
		return token.OrganizationID, true
	}

	return uuid.Nil, false
}
