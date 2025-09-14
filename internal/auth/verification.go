package auth

import (
	"context"
	"fmt"
)

// VerifyEmail verifies a user's email address using a verification token
func (s *Service) VerifyEmail(_ context.Context, _ string) error {
	// TODO: Implement email verification once database migration is added
	// For now, return not implemented error
	return fmt.Errorf("email verification not yet implemented")
}

// ResendVerificationEmail resends the email verification email
func (s *Service) ResendVerificationEmail(_ context.Context, _ string) error {
	// TODO: Implement email verification resend once database migration is added
	// For now, return not implemented error
	return fmt.Errorf("email verification resend not yet implemented")
}
