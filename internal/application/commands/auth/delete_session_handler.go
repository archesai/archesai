package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/core/repositories"
)

// DeleteSessionCommandHandler handles session deletion commands.
type DeleteSessionCommandHandler struct {
	sessionRepo repositories.SessionRepository
}

// NewDeleteSessionCommandHandler creates a new session deletion command handler.
func NewDeleteSessionCommandHandler(
	sessionRepo repositories.SessionRepository,
) *DeleteSessionCommandHandler {
	return &DeleteSessionCommandHandler{
		sessionRepo: sessionRepo,
	}
}

// Handle executes the delete session command.
func (h *DeleteSessionCommandHandler) Handle(
	ctx context.Context,
	cmd *DeleteSessionCommand,
) error {
	if cmd.ID == uuid.Nil {
		return fmt.Errorf("session ID is required")
	}

	// Delete the session
	err := h.sessionRepo.Delete(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
