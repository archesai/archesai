package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/core/services"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// LoginCommandHandler handles login commands.
type LoginCommandHandler struct {
	authService services.AuthService
}

// NewLoginCommandHandler creates a new login command handler.
func NewLoginCommandHandler(authService services.AuthService) *LoginCommandHandler {
	return &LoginCommandHandler{
		authService: authService,
	}
}

// Handle executes the login command and returns authentication tokens.
// This is pure business logic with no knowledge of HTTP or cookies.
func (h *LoginCommandHandler) Handle(
	ctx context.Context,
	cmd *LoginCommand,
) (*valueobjects.AuthTokens, error) {
	if cmd.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if cmd.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Authenticate user and return tokens
	tokens, err := h.authService.AuthenticateWithPassword(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return tokens, nil
}
