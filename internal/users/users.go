// Package users provides user profile management functionality.
// It handles user creation, updates, deletion, and profile information
// management separate from authentication concerns.
package users

// ContextKey is a type for context keys.
type ContextKey string

const (
	// UserContextKey is the context key for the current user.
	UserContextKey ContextKey = "user"
)
