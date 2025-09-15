// Package users provides user profile management functionality.
// It handles user creation, updates, deletion, and profile information
// management separate from authentication concerns.
package users

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package users --include-tags Users ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package users --include-tags Users ../../api/openapi.bundled.yaml

import (
	"errors"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the context key for the current user
	UserContextKey ContextKey = "user"
)

// Domain errors
var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrUserExists is returned when a user already exists
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidUserData is returned when user data is invalid
	ErrInvalidUserData = errors.New("invalid user data")
)
