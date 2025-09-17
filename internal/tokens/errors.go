package tokens

import "errors"

// Domain errors.
var (
	// ErrTokenNotFound is returned when a token is not found.
	ErrTokenNotFound = errors.New("token not found")
	// ErrTokenExists is returned when a token already exists.
	ErrTokenExists = errors.New("token already exists")
	// ErrInvalidTokenData is returned when token data is invalid.
	ErrInvalidTokenData = errors.New("invalid token data")
)
