// Package user provides command handlers for user operations.
package user

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/core/events"
	"github.com/archesai/archesai/internal/core/repositories"
)

// DeleteCurrentUserCommandHandler handles the delete current user command.
type DeleteCurrentUserCommandHandler struct {
	userRepo  repositories.UserRepository
	publisher events.Publisher
}

// NewDeleteCurrentUserCommandHandler creates a new delete current user command handler.
func NewDeleteCurrentUserCommandHandler(
	userRepo repositories.UserRepository,
	publisher events.Publisher,
) *DeleteCurrentUserCommandHandler {
	return &DeleteCurrentUserCommandHandler{
		userRepo:  userRepo,
		publisher: publisher,
	}
}

// Handle executes the delete current user command.
func (h *DeleteCurrentUserCommandHandler) Handle(
	ctx context.Context,
	cmd *DeleteCurrentUserCommand,
) error {
	// Validate confirmation string
	if cmd.Confirmation != "DELETE_MY_ACCOUNT" {
		return fmt.Errorf("invalid confirmation string")
	}

	// Get session to find user ID
	user, err := h.userRepo.GetUserBySessionID(ctx, cmd.SessionID.String())
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Delete user (should cascade to sessions and accounts)
	if err := h.userRepo.Delete(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Publish domain event
	event := events.NewUserDeletedEvent(user.ID)
	if err := h.publisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		_ = err
	}

	return nil
}
