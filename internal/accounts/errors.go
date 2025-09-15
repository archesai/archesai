package accounts

import "errors"

// Account errors
var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrUserExists is returned when a user already exists
	ErrUserExists = errors.New("user already exists")
	// ErrAccountNotFound is returned when an account is not found
	ErrAccountNotFound = errors.New("account not found")
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidProvider is returned when an invalid provider is specified
	ErrInvalidProvider = errors.New("invalid provider")
	// ErrDuplicateAccount is returned when an account already exists
	ErrDuplicateAccount = errors.New("account already exists")
)

// Password errors
var (
	// ErrInvalidPassword is returned for invalid passwords
	ErrInvalidPassword = errors.New("invalid password")
	// ErrPasswordTooShort is returned when password is less than 8 characters
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	// ErrPasswordTooLong is returned when password exceeds 128 characters
	ErrPasswordTooLong = errors.New("password must not exceed 128 characters")
	// ErrPasswordResetServiceNotConfigured is returned when password reset service is not configured
	ErrPasswordResetServiceNotConfigured = errors.New("password reset service not configured")
	// ErrPasswordResetConfirmationNotImplemented is returned when password reset confirmation is not yet implemented
	ErrPasswordResetConfirmationNotImplemented = errors.New("password reset confirmation not yet implemented")
)

// Verification errors
var (
	// ErrEmailVerificationNotImplemented is returned when email verification is not yet implemented
	ErrEmailVerificationNotImplemented = errors.New("email verification not yet implemented")
	// ErrEmailVerificationResendNotImplemented is returned when email verification resend is not yet implemented
	ErrEmailVerificationResendNotImplemented = errors.New("email verification resend not yet implemented")
)
