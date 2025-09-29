// Package auth provides authentication command handlers
package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// OAuthLoginCommand represents the OAuth login command.
type OAuthLoginCommand struct {
	Provider string
	Code     string
	State    string
}

// OAuthLoginCommandHandler handles OAuth login commands.
type OAuthLoginCommandHandler struct {
	authService *auth.Service
}

// NewOAuthLoginCommandHandler creates a new OAuth login command handler.
func NewOAuthLoginCommandHandler(authService *auth.Service) *OAuthLoginCommandHandler {
	return &OAuthLoginCommandHandler{
		authService: authService,
	}
}

// Handle executes the OAuth login command.
func (h *OAuthLoginCommandHandler) Handle(
	ctx context.Context,
	cmd *OAuthLoginCommand,
) (*auth.Tokens, error) {
	if cmd.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	if cmd.Code == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

	// Use auth service to handle OAuth callback
	return h.authService.HandleOAuthCallback(ctx, cmd.Provider, cmd.Code, cmd.State)
}
