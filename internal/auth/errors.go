package auth

import "errors"

// Core domain errors
var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserExists is returned when a user already exists
	ErrUserExists = errors.New("user already exists")
	// ErrAccountNotFound is returned when an account is not found
	ErrAccountNotFound = errors.New("account not found")
	// ErrUnauthorized is returned for unauthorized access
	ErrUnauthorized = errors.New("unauthorized")
	// ErrNotImplemented is returned when a feature is not yet implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrNotFound is a general not found error
	ErrNotFound = errors.New("not found")
)

// Authentication errors
var (
	// ErrTooManyAttempts is returned when too many failed login attempts
	ErrTooManyAttempts = errors.New("too many failed login attempts, try again later")
)

// Token errors
var (
	// ErrInvalidToken is returned for invalid tokens
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token expired")
)

// Session errors
var (
	// ErrSessionNotFound is returned when a session is not found
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionExpired is returned when a session has expired
	ErrSessionExpired = errors.New("session expired")
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

// OAuth errors
var (
	// ErrProviderNotFound is returned when an OAuth provider is not found
	ErrProviderNotFound = errors.New("provider not found")
	// ErrInvalidState is returned when OAuth state parameter is invalid
	ErrInvalidState = errors.New("invalid state parameter")
	// ErrStateNotFound is returned when OAuth state is not found or expired
	ErrStateNotFound = errors.New("state not found or expired")
	// ErrNoRefreshToken is returned when no refresh token is available
	ErrNoRefreshToken = errors.New("no refresh token available")
	// ErrOAuthAccountNotFound is returned when OAuth account is not found
	ErrOAuthAccountNotFound = errors.New("OAuth account not found")
	// ErrNoVerifiedEmail is returned when no verified email is found
	ErrNoVerifiedEmail = errors.New("no verified email found")
)

// API Key errors
var (
	// ErrInvalidAPIKeyFormat is returned when API key format is invalid
	ErrInvalidAPIKeyFormat = errors.New("invalid api key format")
	// ErrAPIKeyExpired is returned when an API key has expired
	ErrAPIKeyExpired = errors.New("api key expired")
	// ErrAPIKeyServiceNotConfigured is returned when API key service is not configured
	ErrAPIKeyServiceNotConfigured = errors.New("API key service not configured")
)

// Verification errors
var (
	// ErrEmailVerificationNotImplemented is returned when email verification is not yet implemented
	ErrEmailVerificationNotImplemented = errors.New("email verification not yet implemented")
	// ErrEmailVerificationResendNotImplemented is returned when email verification resend is not yet implemented
	ErrEmailVerificationResendNotImplemented = errors.New("email verification resend not yet implemented")
)
