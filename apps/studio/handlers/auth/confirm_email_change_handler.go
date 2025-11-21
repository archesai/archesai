// Package auth provides command handlers for authentication operations.
package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	commands "github.com/archesai/archesai/apps/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/pkg/auth"
)

// ConfirmEmailChangeCommandHandler handles email change confirmation commands.
type ConfirmEmailChangeCommandHandler struct {
	authService *auth.Service
}

// NewConfirmEmailChangeCommandHandler creates a new email change confirmation command handler.
func NewConfirmEmailChangeCommandHandler(
	authService *auth.Service,
) *ConfirmEmailChangeCommandHandler {
	return &ConfirmEmailChangeCommandHandler{
		authService: authService,
	}
}

// Handle executes the email change confirmation command.
func (h *ConfirmEmailChangeCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.ConfirmEmailChangeCommand,
) error {
	if cmd.Token == "" {
		return fmt.Errorf("token is required")
	}
	if cmd.NewEmail == "" {
		return fmt.Errorf("new email is required")
	}
	if cmd.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	// Confirm email change
	err := h.authService.ConfirmEmailChange(ctx, cmd.Token, cmd.NewEmail, cmd.UserID)
	if err != nil {
		return fmt.Errorf("failed to confirm email change: %w", err)
	}

	return nil
}
