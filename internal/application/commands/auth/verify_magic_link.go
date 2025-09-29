package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// VerifyMagicLinkCommand represents the magic link verification command.
type VerifyMagicLinkCommand struct {
	Token string
}

// VerifyMagicLinkCommandHandler handles magic link verification commands.
type VerifyMagicLinkCommandHandler struct {
	authService *auth.Service
}

// NewVerifyMagicLinkCommandHandler creates a new magic link verification command handler.
func NewVerifyMagicLinkCommandHandler(authService *auth.Service) *VerifyMagicLinkCommandHandler {
	return &VerifyMagicLinkCommandHandler{
		authService: authService,
	}
}

// Handle executes the magic link verification command.
func (h *VerifyMagicLinkCommandHandler) Handle(
	ctx context.Context,
	cmd *VerifyMagicLinkCommand,
) (*auth.AuthTokens, error) {
	if cmd.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Verify magic link and create session
	return h.authService.VerifyMagicLink(ctx, cmd.Token)
}
