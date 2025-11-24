package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	commands "github.com/archesai/archesai/apps/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/apps/studio/generated/core"
	"github.com/archesai/archesai/apps/studio/generated/core/repositories"
)

// UpdateSessionCommandHandler handles session update commands.
type UpdateSessionCommandHandler struct {
	sessionRepo repositories.SessionRepository
}

// NewUpdateSessionCommandHandler creates a new session update command handler.
func NewUpdateSessionCommandHandler(
	sessionRepo repositories.SessionRepository,
) *UpdateSessionCommandHandler {
	return &UpdateSessionCommandHandler{
		sessionRepo: sessionRepo,
	}
}

// Handle executes the update session command.
func (h *UpdateSessionCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.UpdateSessionCommand,
) (*core.Session, error) {
	if cmd.ID == uuid.Nil {
		return nil, fmt.Errorf("session ID is required")
	}

	// Get existing session
	session, err := h.sessionRepo.Get(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Update session fields if provided
	if cmd.OrganizationID != uuid.Nil {
		session.OrganizationID = &cmd.OrganizationID
	}

	// Save updated session
	if session, err = h.sessionRepo.Update(ctx, session.ID, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}
