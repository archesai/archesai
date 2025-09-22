// Package tokens provides token management and validation services
package tokens

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/auth"
)

// APIValidator handles API token validation with rate limiting and scope checking.
type APIValidator struct {
	store auth.APITokenStore
}

// NewAPIValidator creates a new API token validator.
func NewAPIValidator(store auth.APITokenStore) *APIValidator {
	return &APIValidator{
		store: store,
	}
}

// NewAPITokenValidator creates a new API token validator without store (dummy implementation).
func NewAPITokenValidator() auth.APITokenValidator {
	return &APITokenValidatorImpl{}
}

// ValidateAPIKey validates an API key and returns the token data.
func (v *APIValidator) ValidateAPIKey(ctx context.Context, key string) (*auth.APIToken, error) {
	if key == "" {
		return nil, auth.ErrInvalidAPIKey
	}

	// Validate and retrieve token data
	token, err := v.store.ValidateToken(ctx, key)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ValidateAPIKeyWithScopes validates an API key and checks required scopes.
func (v *APIValidator) ValidateAPIKeyWithScopes(
	ctx context.Context,
	key string,
	requiredScopes []string,
) (*auth.APIToken, error) {
	token, err := v.ValidateAPIKey(ctx, key)
	if err != nil {
		return nil, err
	}

	// Check required scopes
	if !v.hasRequiredScopes(token.Scopes, requiredScopes) {
		return nil, auth.ErrInsufficientScopes
	}

	return token, nil
}

// ValidateAPIKeyForOrganization validates an API key for a specific organization.
func (v *APIValidator) ValidateAPIKeyForOrganization(
	ctx context.Context,
	key string,
	organizationID uuid.UUID,
) (*auth.APIToken, error) {
	token, err := v.ValidateAPIKey(ctx, key)
	if err != nil {
		return nil, err
	}

	// Check organization access
	if token.OrganizationID != organizationID {
		return nil, auth.ErrUnauthorizedOrganization
	}

	return token, nil
}

// ExtractAPIKeyFromHeaders extracts API key from various header formats.
func (v *APIValidator) ExtractAPIKeyFromHeaders(headers map[string]string) string {
	// Check Authorization header with various schemes
	if authHeader, ok := headers["Authorization"]; ok {
		if key := v.store.ParseAPIKey(authHeader); key != "" {
			return key
		}
	}

	// Check X-API-Key header
	if apiKey, ok := headers["X-API-Key"]; ok {
		if key := v.store.ParseAPIKey(apiKey); key != "" {
			return key
		}
	}

	// Check Api-Key header (alternative)
	if apiKey, ok := headers["Api-Key"]; ok {
		if key := v.store.ParseAPIKey(apiKey); key != "" {
			return key
		}
	}

	return ""
}

// CheckRateLimit checks if the API key has exceeded its rate limit.
func (v *APIValidator) CheckRateLimit(
	_ context.Context,
	token *auth.APIToken,
	window time.Duration,
) error {
	// This is a simplified rate limit check
	// In production, you would use Redis or similar for distributed rate limiting

	if token.LastUsedAt == nil {
		return nil // First use
	}

	// Check if last use was within the rate limit window
	if time.Since(*token.LastUsedAt) < window {
		// For now, just check if token was used recently
		// This should be replaced with proper rate limiting logic
		return nil
	}

	return nil
}

// ValidateScopes checks if a token has all required scopes.
func (v *APIValidator) ValidateScopes(tokenScopes, requiredScopes []string) error {
	if !v.hasRequiredScopes(tokenScopes, requiredScopes) {
		return auth.ErrInsufficientScopes
	}
	return nil
}

// GetScopeDescription returns human-readable descriptions for scopes.
func (v *APIValidator) GetScopeDescription(scope string) string {
	descriptions := map[string]string{
		"read":              "Read access to resources",
		"write":             "Write access to resources",
		"admin":             "Administrative access",
		"pipelines:read":    "Read access to pipelines",
		"pipelines:write":   "Write access to pipelines",
		"pipelines:execute": "Execute pipelines",
		"artifacts:read":    "Read access to artifacts",
		"artifacts:write":   "Write access to artifacts",
		"tools:read":        "Read access to tools",
		"tools:write":       "Write access to tools",
		"tools:execute":     "Execute tools",
		"organizations":     "Access to organization resources",
		"users:read":        "Read user information",
		"users:write":       "Modify user information",
	}

	if desc, ok := descriptions[scope]; ok {
		return desc
	}
	return scope
}

// ListScopesForToken returns all scopes for a given token.
func (v *APIValidator) ListScopesForToken(
	_ context.Context,
	_ uuid.UUID,
) ([]string, error) {
	// This would typically get the token from the store and return its scopes
	// For now, returning empty as we need to implement the store method
	return []string{}, nil
}

// hasRequiredScopes checks if token scopes contain all required scopes.
func (v *APIValidator) hasRequiredScopes(tokenScopes, requiredScopes []string) bool {
	if len(requiredScopes) == 0 {
		return true
	}

	// Create a map for quick lookup
	scopeMap := make(map[string]bool)
	for _, scope := range tokenScopes {
		scopeMap[scope] = true

		// Also check for wildcard scopes
		if scope == "*" || scope == "admin" {
			return true // Admin or wildcard scope grants all access
		}
	}

	// Check each required scope
	for _, required := range requiredScopes {
		// Check exact match
		if scopeMap[required] {
			continue
		}

		// Check parent scope (e.g., "pipelines" covers "pipelines:read")
		if strings.Contains(required, ":") {
			parent := strings.Split(required, ":")[0]
			if scopeMap[parent] {
				continue
			}
		}

		// Required scope not found
		return false
	}

	return true
}

// APITokenValidatorImpl is a simple implementation of APITokenValidator
type APITokenValidatorImpl struct{}

// ValidateAPIKey validates an API key
func (v *APITokenValidatorImpl) ValidateAPIKey(
	_ context.Context,
	_ string,
) (*auth.APIToken, error) {
	// Dummy implementation - would need actual store
	return nil, auth.ErrInvalidAPIKey
}

// ValidateAPIKeyWithScopes validates an API key with required scopes
func (v *APITokenValidatorImpl) ValidateAPIKeyWithScopes(
	_ context.Context,
	_ string,
	_ []string,
) (*auth.APIToken, error) {
	// Dummy implementation - would need actual store
	return nil, auth.ErrInvalidAPIKey
}

// ValidateAPIKeyForOrganization validates an API key for an organization
func (v *APITokenValidatorImpl) ValidateAPIKeyForOrganization(
	_ context.Context,
	_ string,
	_ uuid.UUID,
) (*auth.APIToken, error) {
	// Dummy implementation - would need actual store
	return nil, auth.ErrInvalidAPIKey
}

// CheckRateLimit checks rate limit for a token
func (v *APITokenValidatorImpl) CheckRateLimit(
	_ context.Context,
	_ *auth.APIToken,
	_ time.Duration,
) error {
	// Simple implementation - no rate limiting
	return nil
}

// ValidateScopes checks if provided scopes contain all required scopes
func (v *APITokenValidatorImpl) ValidateScopes(tokenScopes, requiredScopes []string) error {
	scopeMap := make(map[string]bool)
	for _, s := range tokenScopes {
		scopeMap[s] = true
	}

	for _, required := range requiredScopes {
		if !scopeMap[required] {
			return auth.ErrInsufficientScopes
		}
	}
	return nil
}

// ExtractAPIKeyFromHeaders extracts API key from various header formats
func (v *APITokenValidatorImpl) ExtractAPIKeyFromHeaders(headers map[string]string) string {
	// Check Authorization header
	if auth := headers["Authorization"]; auth != "" {
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			token := strings.TrimSpace(auth[7:])
			if strings.HasPrefix(token, "sk_") {
				return token
			}
		}
	}

	// Check X-API-Key header
	if apiKey := headers["X-API-Key"]; apiKey != "" {
		return apiKey
	}

	// Check Api-Key header
	if apiKey := headers["Api-Key"]; apiKey != "" {
		return apiKey
	}

	return ""
}
