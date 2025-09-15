// Package sessions provides session management functionality
package sessions

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package sessions --include-tags Sessions ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package sessions --include-tags Sessions ../../api/openapi.bundled.yaml

import "errors"

var (
	// ErrSessionNotFound is returned when a session is not found
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionExpired is returned when a session has expired
	ErrSessionExpired = errors.New("session has expired")
	// ErrInvalidToken is returned when a session token is invalid
	ErrInvalidToken = errors.New("invalid session token")
)
