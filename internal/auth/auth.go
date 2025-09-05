// Package auth provides authentication and authorization functionality.
// It includes user management, session handling, JWT token generation,
// and middleware for protecting routes.
package auth

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the context key for the authenticated user
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the context key for JWT claims
	ClaimsContextKey ContextKey = "claims"
)
