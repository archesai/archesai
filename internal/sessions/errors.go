// Package sessions provides session management functionality
package sessions

import "errors"

// Common session errors
var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
	ErrInvalidToken    = errors.New("invalid session token")
)
