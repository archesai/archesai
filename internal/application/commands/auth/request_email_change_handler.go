package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// RequestEmailChangeCommandHandler handles email change request commands.
type RequestEmailChangeCommandHandler struct {
	authService *auth.Service
}

// NewRequestEmailChangeCommandHandler creates a new email change request command handler.
func NewRequestEmailChangeCommandHandler(
	authService *auth.Service,
) *RequestEmailChangeCommandHandler {
	return &RequestEmailChangeCommandHandler{
		authService: authService,
	}
}

// Handle executes the email change request command.
func (h *RequestEmailChangeCommandHandler) Handle(
	ctx context.Context,
	cmd *RequestEmailChangeCommand,
) error {
	if cmd.SessionID == uuid.Nil {
		return fmt.Errorf("session ID is required")
	}
	if cmd.NewEmail == "" {
		return fmt.Errorf("new email is required")
	}

	// Request email change
	err := h.authService.RequestEmailChange(ctx, cmd.SessionID, cmd.NewEmail)
	if err != nil {
		return fmt.Errorf("failed to request email change: %w", err)
	}

	return nil
}
