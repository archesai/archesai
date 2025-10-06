package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// LogoutAllCommandHandler handles logout all sessions commands.
type LogoutAllCommandHandler struct {
	authService *auth.Service
}

// NewLogoutAllCommandHandler creates a new logout all command handler.
func NewLogoutAllCommandHandler(authService *auth.Service) *LogoutAllCommandHandler {
	return &LogoutAllCommandHandler{
		authService: authService,
	}
}

// Handle executes the logout all command.
func (h *LogoutAllCommandHandler) Handle(
	ctx context.Context,
	cmd *LogoutAllCommand,
) (any, error) {
	if cmd.SessionID == uuid.Nil {
		return nil, fmt.Errorf("session ID is required")
	}

	// Delete all sessions for the user (gets user from session)
	err := h.authService.DeleteAllUserSessions(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to logout all sessions: %w", err)
	}

	return nil, nil
}
