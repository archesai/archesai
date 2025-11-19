package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	commands "github.com/archesai/archesai/apis/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/pkg/auth"
)

// LogoutCommandHandler handles logout commands.
type LogoutCommandHandler struct {
	authService *auth.Service
}

// NewLogoutCommandHandler creates a new logout command handler.
func NewLogoutCommandHandler(authService *auth.Service) *LogoutCommandHandler {
	return &LogoutCommandHandler{
		authService: authService,
	}
}

// Handle executes the logout command.
func (h *LogoutCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.LogoutCommand,
) (any, error) {
	if cmd.SessionID == uuid.Nil {
		return nil, fmt.Errorf("session ID is required")
	}

	// Delete the current session by ID
	err := h.authService.DeleteSessionByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to logout: %w", err)
	}

	return nil, nil
}
