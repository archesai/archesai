// Package auth provides authentication and authorization functionality.
// It includes user management, session handling, JWT token generation,
// and middleware for protecting routes.
package auth

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package auth --include-tags Auth,Sessions,Accounts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package auth --include-tags Auth,Sessions,Accounts ../../api/openapi.bundled.yaml

import (
	"errors"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the context key for the authenticated user
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the context key for JWT claims
	ClaimsContextKey ContextKey = "claims"
	// SessionTokenContextKey is the context key for session token
	SessionTokenContextKey ContextKey = "session_token"
)

// Domain errors
var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserExists is returned when a user already exists
	ErrUserExists = errors.New("user already exists")
)
