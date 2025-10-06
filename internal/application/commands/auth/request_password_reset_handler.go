package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// RequestPasswordResetCommandHandler handles password reset request commands.
type RequestPasswordResetCommandHandler struct {
	authService *auth.Service
}

// NewRequestPasswordResetCommandHandler creates a new password reset request command handler.
func NewRequestPasswordResetCommandHandler(
	authService *auth.Service,
) *RequestPasswordResetCommandHandler {
	return &RequestPasswordResetCommandHandler{
		authService: authService,
	}
}

// Handle executes the password reset request command.
func (h *RequestPasswordResetCommandHandler) Handle(
	ctx context.Context,
	cmd *RequestPasswordResetCommand,
) (any, error) {
	if cmd.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Request password reset
	err := h.authService.RequestPasswordReset(ctx, cmd.Email)
	if err != nil {
		// Don't reveal if email exists or not
		// Log the actual error but return generic message
		return nil, nil
	}

	return nil, nil
}
