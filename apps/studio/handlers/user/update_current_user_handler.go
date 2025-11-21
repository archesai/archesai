package user

import (
	"context"
	"fmt"
	"time"

	queries "github.com/archesai/archesai/apps/studio/generated/application/commands/user"
	domainevents "github.com/archesai/archesai/apps/studio/generated/core/events"
	"github.com/archesai/archesai/apps/studio/generated/core/models"
	"github.com/archesai/archesai/apps/studio/generated/core/repositories"
	"github.com/archesai/archesai/pkg/events"
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
	cmd *queries.UpdateCurrentUserCommand,
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
	event := domainevents.NewUserUpdatedEvent(updated.ID)
	if err := events.PublishDomainEvent(ctx, h.publisher, event); err != nil {
		// Log error but don't fail the operation
		_ = err
	}

	return updated, nil
}
