// Package auth provides authentication and authorization functionality.
// It includes user management, session handling, JWT token generation,
// and middleware for protecting routes.
package auth

//go:generate go tool oapi-codegen --config=models.cfg.yaml ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=server.cfg.yaml ../../api/openapi.bundled.yaml

// User represents a user with authentication data
type User struct {
	UserEntity
	PasswordHash string // Not exposed in API
}

// Session represents an authenticated session
type Session struct {
	SessionEntity
	Token     string // Session token (not exposed in API)
	IpAddress string // Client IP
	UserAgent string // Client user agent
}

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the context key for the authenticated user
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the context key for JWT claims
	ClaimsContextKey ContextKey = "claims"
)
