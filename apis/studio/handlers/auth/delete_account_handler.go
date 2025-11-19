package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	commands "github.com/archesai/archesai/apis/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/pkg/auth"
)

// DeleteAccountCommandHandler handles account deletion commands.
type DeleteAccountCommandHandler struct {
	authService *auth.Service
}

// NewDeleteAccountCommandHandler creates a new account deletion command handler.
func NewDeleteAccountCommandHandler(authService *auth.Service) *DeleteAccountCommandHandler {
	return &DeleteAccountCommandHandler{
		authService: authService,
	}
}

// Handle executes the account deletion command.
func (h *DeleteAccountCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.DeleteAccountCommand,
) error {
	if cmd.SessionID == uuid.Nil {
		return fmt.Errorf("session ID is required")
	}

	// Delete the account
	_, err := h.authService.DeleteAccount(ctx, cmd.SessionID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}
