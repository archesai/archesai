package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// LoginCommandHandler handles login commands.
type LoginCommandHandler struct {
	authService *auth.Service
}

// NewLoginCommandHandler creates a new login command handler.
func NewLoginCommandHandler(authService *auth.Service) *LoginCommandHandler {
	return &LoginCommandHandler{
		authService: authService,
	}
}

// Handle executes the login command.
func (h *LoginCommandHandler) Handle(
	ctx context.Context,
	cmd *LoginCommand,
) (*entities.Session, error) {
	if cmd.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if cmd.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Authenticate user and create session
	tokens, err := h.authService.AuthenticateWithPassword(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Get the session from the auth service
	session, err := h.authService.GetSessionByToken(ctx, tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	return session, nil
}
