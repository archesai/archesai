package auth

import (
	"context"
	"fmt"

	commands "github.com/archesai/archesai/apps/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/pkg/auth"
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
	cmd *commands.ConfirmPasswordResetCommand,
) error {
	if cmd.Token == "" {
		return fmt.Errorf("token is required")
	}
	if cmd.NewPassword == "" {
		return fmt.Errorf("new password is required")
	}

	// Confirm password reset
	err := h.authService.ConfirmPasswordReset(ctx, cmd.Token, cmd.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}
