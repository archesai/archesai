package entities

import "errors"

var (
	// ErrInvalidCredentials indicates the provided login credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserNotFound indicates the requested user was not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserExists indicates a user already exists with the given email.
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidToken indicates the provided token is malformed or invalid.
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired indicates the provided token has expired.
	ErrTokenExpired = errors.New("token expired")
	// ErrUnauthorized indicates the user lacks required permissions.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrSessionNotFound indicates the requested session was not found.
	ErrSessionNotFound = errors.New("session not found")
)
