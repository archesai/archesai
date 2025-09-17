package accounts

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordValidator validates password strength.
type PasswordValidator struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireDigit   bool
	RequireSpecial bool
}

// NewPasswordValidator creates a new PasswordValidator with given settings.
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:      8,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: false, // Keep this optional for better UX
	}
}

// DefaultPasswordValidator returns a password validator with default settings.
func DefaultPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:      8,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: false, // Keep this optional for better UX
	}
}

// Validate checks if a password meets the strength requirements.
func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.MinLength {
		return fmt.Errorf(
			"%w: password must be at least %d characters",
			ErrWeakPassword,
			v.MinLength,
		)
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if v.RequireUpper && !hasUpper {
		return fmt.Errorf(
			"%w: password must contain at least one uppercase letter",
			ErrWeakPassword,
		)
	}
	if v.RequireLower && !hasLower {
		return fmt.Errorf(
			"%w: password must contain at least one lowercase letter",
			ErrWeakPassword,
		)
	}
	if v.RequireDigit && !hasDigit {
		return fmt.Errorf("%w: password must contain at least one digit", ErrWeakPassword)
	}
	if v.RequireSpecial && !hasSpecial {
		return fmt.Errorf(
			"%w: password must contain at least one special character",
			ErrWeakPassword,
		)
	}

	return nil
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	// Use bcrypt default cost (currently 10)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against its hash.
func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid password")
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}
