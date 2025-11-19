package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	commands "github.com/archesai/archesai/apis/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/apis/studio/generated/core/models"
	"github.com/archesai/archesai/pkg/auth"
)

// UpdateAccountCommandHandler handles account update commands.
type UpdateAccountCommandHandler struct {
	authService *auth.Service
}

// NewUpdateAccountCommandHandler creates a new account update command handler.
func NewUpdateAccountCommandHandler(
	authService *auth.Service,
) *UpdateAccountCommandHandler {
	return &UpdateAccountCommandHandler{
		authService: authService,
	}
}

// Handle executes the account update command.
func (h *UpdateAccountCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.UpdateAccountCommand,
) (*models.User, error) {
	if cmd.SessionID == uuid.Nil {
		return nil, fmt.Errorf("session ID is required")
	}

	// Prepare update data
	updateData := make(map[string]any)

	if cmd.Provider != nil {
		updateData["provider"] = *cmd.Provider
	}
	if cmd.ProviderAccountIdentifier != nil {
		updateData["provider_account_id"] = *cmd.ProviderAccountIdentifier
	}
	if cmd.Type != nil {
		updateData["type"] = *cmd.Type
	}

	// Update the account
	account, err := h.authService.UpdateAccount(ctx, cmd.SessionID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}
