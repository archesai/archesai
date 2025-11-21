package auth

import (
	"context"
	"fmt"

	commands "github.com/archesai/archesai/apps/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/pkg/auth"
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

// Handle executes the login command and returns authentication tokens.
// This is pure business logic with no knowledge of HTTP or cookies.
func (h *LoginCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.LoginCommand,
) (*auth.Tokens, error) {
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
