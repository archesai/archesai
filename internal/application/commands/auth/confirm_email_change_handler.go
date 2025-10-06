// Package auth provides command handlers for authentication operations.
package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/infrastructure/auth"
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
	cmd *ConfirmEmailChangeCommand,
) (any, error) {
	if cmd.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if cmd.NewEmail == "" {
		return nil, fmt.Errorf("new email is required")
	}
	if cmd.UserID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	// Confirm email change
	err := h.authService.ConfirmEmailChange(ctx, cmd.Token, cmd.NewEmail, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm email change: %w", err)
	}

	return nil, nil
}
