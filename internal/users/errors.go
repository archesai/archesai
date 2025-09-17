package users

import "errors"

// Domain errors.
var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserExists is returned when a user already exists.
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidUserData is returned when user data is invalid.
	ErrInvalidUserData = errors.New("invalid user data")
)
