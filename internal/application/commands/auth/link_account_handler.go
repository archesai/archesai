package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/core/models"
	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// LinkAccountCommandHandler handles account linking commands.
type LinkAccountCommandHandler struct {
	authService *auth.Service
}

// NewLinkAccountCommandHandler creates a new account linking command handler.
func NewLinkAccountCommandHandler(authService *auth.Service) *LinkAccountCommandHandler {
	return &LinkAccountCommandHandler{
		authService: authService,
	}
}

// Handle executes the account linking command.
func (h *LinkAccountCommandHandler) Handle(
	ctx context.Context,
	cmd *LinkAccountCommand,
) (*models.Account, error) {
	if cmd.SessionID == uuid.Nil {
		return nil, fmt.Errorf("session ID is required")
	}
	if cmd.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}

	// Link the account
	account, err := h.authService.LinkAccount(ctx, cmd.SessionID, cmd.Provider, cmd.RedirectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to link account: %w", err)
	}

	return account, nil
}
