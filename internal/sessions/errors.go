package sessions

import "errors"

var (
	// ErrSessionNotFound is returned when a session is not found.
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionExpired is returned when a session has expired.
	ErrSessionExpired = errors.New("session has expired")
	// ErrInvalidToken is returned when a session token is invalid.
	ErrInvalidToken = errors.New("invalid session token")
)
