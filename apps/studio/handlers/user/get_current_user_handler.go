// Package user provides query handlers for user operations.
package user

import (
	"context"
	"fmt"

	queries "github.com/archesai/archesai/apps/studio/generated/application/queries/user"
	"github.com/archesai/archesai/apps/studio/generated/core/models"
	"github.com/archesai/archesai/apps/studio/generated/core/repositories"
)

// GetCurrentUserQueryHandler handles the get current user query.
type GetCurrentUserQueryHandler struct {
	userRepo repositories.UserRepository
}

// NewGetCurrentUserQueryHandler creates a new get current user query handler.
func NewGetCurrentUserQueryHandler(
	userRepo repositories.UserRepository,
) *GetCurrentUserQueryHandler {
	return &GetCurrentUserQueryHandler{
		userRepo: userRepo,
	}
}

// Handle executes the get current user query.
func (h *GetCurrentUserQueryHandler) Handle(
	ctx context.Context,
	query *queries.GetCurrentUserQuery,
) (*models.User, error) {
	// Get user by ID from session
	user, err := h.userRepo.GetUserBySessionID(ctx, query.SessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
