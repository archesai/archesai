// Package user provides command handlers for user operations.
package user

import (
	"context"
	"fmt"
	"log/slog"

	commands "github.com/archesai/archesai/apps/studio/generated/application/commands/user"
	"github.com/archesai/archesai/apps/studio/generated/core"
	"github.com/archesai/archesai/apps/studio/generated/core/repositories"
	"github.com/archesai/archesai/pkg/events"
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
	cmd *commands.DeleteCurrentUserCommand,
) error {
	// Validate confirmation string
	if cmd.XConfirm != "DELETE_MY_ACCOUNT" {
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
	event := core.NewUserDeletedEvent(user.ID)
	if err := h.publisher.Publish(ctx, event); err != nil {
		slog.Error("failed to publish event", "error", err)
	}

	return nil
}
