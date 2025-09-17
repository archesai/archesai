package accounts

import "errors"

// Account errors.
var (
	ErrAccountNotFound = errors.New("account not found")
	// ErrInvalidCredentials is returned when credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidProvider is returned when an invalid provider is specified.
	ErrInvalidProvider = errors.New("invalid provider")
	// ErrDuplicateAccount is returned when an account already exists.
	ErrDuplicateAccount = errors.New("account already exists")
	// ErrWeakPassword is returned when a password doesn't meet strength requirements.
	ErrWeakPassword = errors.New("password does not meet strength requirements")
)
