// Package accounts provides account management functionality.
package accounts

import "errors"

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package accounts --include-tags Accounts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package accounts --include-tags Accounts ../../api/openapi.bundled.yaml

// Account errors
var (
	ErrAccountNotFound = errors.New("account not found")
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidProvider is returned when an invalid provider is specified
	ErrInvalidProvider = errors.New("invalid provider")
	// ErrDuplicateAccount is returned when an account already exists
	ErrDuplicateAccount = errors.New("account already exists")
)
