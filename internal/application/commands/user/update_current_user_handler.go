package user

import (
	"context"
	"fmt"
	"time"

	"github.com/archesai/archesai/internal/core/events"
	"github.com/archesai/archesai/internal/core/models"
	"github.com/archesai/archesai/internal/core/repositories"
)

// UpdateCurrentUserCommandHandler handles the update current user command.
type UpdateCurrentUserCommandHandler struct {
	userRepo  repositories.UserRepository
	publisher events.Publisher
}

// NewUpdateCurrentUserCommandHandler creates a new update current user command handler.
func NewUpdateCurrentUserCommandHandler(
	userRepo repositories.UserRepository,
	publisher events.Publisher,
) *UpdateCurrentUserCommandHandler {
	return &UpdateCurrentUserCommandHandler{
		userRepo:  userRepo,
		publisher: publisher,
	}
}

// Handle executes the update current user command.
func (h *UpdateCurrentUserCommandHandler) Handle(
	ctx context.Context,
	cmd *UpdateCurrentUserCommand,
) (*models.User, error) {
	// Get existing user
	user, err := h.userRepo.GetUserBySessionID(ctx, cmd.SessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if cmd.Name != nil {
		user.Name = *cmd.Name
	}
	if cmd.Image != nil {
		user.Image = cmd.Image
	}
	user.UpdatedAt = time.Now().UTC()

	// Save updated user
	updated, err := h.userRepo.Update(ctx, user.ID, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Publish domain event
	event := events.NewUserUpdatedEvent(updated.ID)
	if err := h.publisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		_ = err
	}

	return updated, nil
}
