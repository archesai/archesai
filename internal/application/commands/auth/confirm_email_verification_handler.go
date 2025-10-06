package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// ConfirmEmailVerificationCommandHandler handles email verification confirmation commands.
type ConfirmEmailVerificationCommandHandler struct {
	authService *auth.Service
}

// NewConfirmEmailVerificationCommandHandler creates a new email verification confirmation command handler.
func NewConfirmEmailVerificationCommandHandler(
	authService *auth.Service,
) *ConfirmEmailVerificationCommandHandler {
	return &ConfirmEmailVerificationCommandHandler{
		authService: authService,
	}
}

// Handle executes the email verification confirmation command.
func (h *ConfirmEmailVerificationCommandHandler) Handle(
	ctx context.Context,
	cmd *ConfirmEmailVerificationCommand,
) (any, error) {
	if cmd.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Confirm email verification
	err := h.authService.ConfirmEmailVerification(ctx, cmd.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	return nil, nil
}
