package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// ConfirmPasswordResetCommandHandler handles password reset confirmation commands.
type ConfirmPasswordResetCommandHandler struct {
	authService *auth.Service
}

// NewConfirmPasswordResetCommandHandler creates a new password reset confirmation command handler.
func NewConfirmPasswordResetCommandHandler(
	authService *auth.Service,
) *ConfirmPasswordResetCommandHandler {
	return &ConfirmPasswordResetCommandHandler{
		authService: authService,
	}
}

// Handle executes the password reset confirmation command.
func (h *ConfirmPasswordResetCommandHandler) Handle(
	ctx context.Context,
	cmd *ConfirmPasswordResetCommand,
) (interface{}, error) {
	if cmd.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if cmd.NewPassword == "" {
		return nil, fmt.Errorf("new password is required")
	}

	// Confirm password reset
	err := h.authService.ConfirmPasswordReset(ctx, cmd.Token, cmd.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to reset password: %w", err)
	}

	return nil, nil
}
