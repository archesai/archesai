package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// RequestEmailVerificationCommandHandler handles email verification request commands.
type RequestEmailVerificationCommandHandler struct {
	authService *auth.Service
}

// NewRequestEmailVerificationCommandHandler creates a new email verification request command handler.
func NewRequestEmailVerificationCommandHandler(
	authService *auth.Service,
) *RequestEmailVerificationCommandHandler {
	return &RequestEmailVerificationCommandHandler{
		authService: authService,
	}
}

// Handle executes the email verification request command.
func (h *RequestEmailVerificationCommandHandler) Handle(
	ctx context.Context,
	cmd *RequestEmailVerificationCommand,
) error {
	if cmd.SessionID == uuid.Nil {
		return fmt.Errorf("session ID is required")
	}

	// Request email verification
	err := h.authService.RequestEmailVerification(ctx, cmd.SessionID)
	if err != nil {
		return fmt.Errorf("failed to request email verification: %w", err)
	}

	return nil
}
