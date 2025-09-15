// Package auth provides enhanced JWT claims structures for authentication
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	// AccessTokenType represents an access token
	AccessTokenType TokenType = "access"
	// RefreshTokenType represents a refresh token
	RefreshTokenType TokenType = "refresh"
	// APIKeyTokenType represents an API key token
	APIKeyTokenType TokenType = "api_key"
	// SessionTokenType represents a session token
	SessionTokenType TokenType = "session"
)

// Method represents the authentication method used
type Method string

const (
	// AuthMethodPassword represents password authentication
	AuthMethodPassword Method = "password"
	// AuthMethodOAuth represents OAuth authentication
	AuthMethodOAuth Method = "oauth"
	// AuthMethodAPIKey represents API key authentication
	AuthMethodAPIKey Method = "api_key"
	// AuthMethodMFA represents multi-factor authentication
	AuthMethodMFA Method = "mfa"
)

// EnhancedClaims represents comprehensive JWT claims with rich context
type EnhancedClaims struct {
	// Standard JWT claims
	jwt.RegisteredClaims

	// User Identity
	UserID    uuid.UUID `json:"uid"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`

	// Organization Context
	OrganizationID   uuid.UUID           `json:"org_id,omitempty"`
	OrganizationName string              `json:"org_name,omitempty"`
	OrganizationRole string              `json:"org_role,omitempty"`
	Organizations    []OrganizationClaim `json:"orgs,omitempty"`

	// Permissions and Roles
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`

	// Security Metadata
	TokenType     TokenType `json:"token_type"`
	AuthMethod    Method    `json:"auth_method"`
	SessionID     string    `json:"sid,omitempty"`
	IPAddress     string    `json:"ip,omitempty"`
	UserAgent     string    `json:"ua,omitempty"`
	EmailVerified bool      `json:"email_verified"`
	MFAEnabled    bool      `json:"mfa_enabled"`
	MFAVerified   bool      `json:"mfa_verified,omitempty"`

	// Provider Information (for OAuth)
	Provider         string `json:"provider,omitempty"`
	ProviderID       string `json:"provider_id,omitempty"`
	ProviderTokenExp *int64 `json:"provider_token_exp,omitempty"`

	// Feature Flags
	Features map[string]bool `json:"features,omitempty"`

	// Custom Claims
	CustomClaims map[string]interface{} `json:"custom,omitempty"`
}

// OrganizationClaim represents organization membership in claims
type OrganizationClaim struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions,omitempty"`
}

// RefreshClaims represents minimal claims for refresh tokens
type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID     uuid.UUID `json:"uid"`
	TokenType  TokenType `json:"token_type"`
	SessionID  string    `json:"sid"`
	AuthMethod Method    `json:"auth_method"`
}

// SessionClaims represents minimal claims for session tokens
type SessionClaims struct {
	jwt.RegisteredClaims
	UserID    uuid.UUID `json:"uid"`
	TokenType TokenType `json:"token_type"`
}

// APIKeyClaims represents claims for API key tokens
type APIKeyClaims struct {
	jwt.RegisteredClaims
	KeyID          string    `json:"kid"`
	UserID         uuid.UUID `json:"uid"`
	OrganizationID uuid.UUID `json:"org_id"`
	Name           string    `json:"name"`
	Scopes         []string  `json:"scopes"`
	RateLimit      int       `json:"rate_limit,omitempty"`
}

// HasPermission checks if the claims contain a specific permission
func (c *EnhancedClaims) HasPermission(permission string) bool {
	for _, p := range c.Permissions {
		if p == permission {
			return true
		}
	}
	// Check organization-specific permissions
	if c.OrganizationID != uuid.Nil {
		for _, org := range c.Organizations {
			if org.ID == c.OrganizationID {
				for _, p := range org.Permissions {
					if p == permission {
						return true
					}
				}
			}
		}
	}
	return false
}

// HasRole checks if the claims contain a specific role
func (c *EnhancedClaims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return c.OrganizationRole == role
}

// HasScope checks if the claims contain a specific scope
func (c *EnhancedClaims) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// IsOrgMember checks if the user is a member of a specific organization
func (c *EnhancedClaims) IsOrgMember(orgID uuid.UUID) bool {
	if c.OrganizationID == orgID {
		return true
	}
	for _, org := range c.Organizations {
		if org.ID == orgID {
			return true
		}
	}
	return false
}

// GetOrgRole returns the user's role in a specific organization
func (c *EnhancedClaims) GetOrgRole(orgID uuid.UUID) string {
	if c.OrganizationID == orgID {
		return c.OrganizationRole
	}
	for _, org := range c.Organizations {
		if org.ID == orgID {
			return org.Role
		}
	}
	return ""
}

// IsValid checks if the claims are valid
func (c *EnhancedClaims) IsValid() bool {
	now := time.Now()

	// Check expiration
	if c.ExpiresAt != nil && now.After(c.ExpiresAt.Time) {
		return false
	}

	// Check not before
	if c.NotBefore != nil && now.Before(c.NotBefore.Time) {
		return false
	}

	// Check required fields
	if c.UserID == uuid.Nil || c.Email == "" {
		return false
	}

	return true
}

// ValidateForEndpoint checks if claims are valid for a specific endpoint
func (c *EnhancedClaims) ValidateForEndpoint(requiredScopes []string, requiredPermissions []string) bool {
	// Check if claims are valid
	if !c.IsValid() {
		return false
	}

	// Check required scopes
	for _, scope := range requiredScopes {
		if !c.HasScope(scope) {
			return false
		}
	}

	// Check required permissions
	for _, perm := range requiredPermissions {
		if !c.HasPermission(perm) {
			return false
		}
	}

	return true
}
