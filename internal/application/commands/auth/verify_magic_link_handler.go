package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/core/services"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// VerifyMagicLinkCommandHandler handles magic link verification commands.
type VerifyMagicLinkCommandHandler struct {
	authService services.AuthService
}

// NewVerifyMagicLinkCommandHandler creates a new magic link verification command handler.
func NewVerifyMagicLinkCommandHandler(
	authService services.AuthService,
) *VerifyMagicLinkCommandHandler {
	return &VerifyMagicLinkCommandHandler{
		authService: authService,
	}
}

// Handle executes the magic link verification command.
func (h *VerifyMagicLinkCommandHandler) Handle(
	ctx context.Context,
	cmd *VerifyMagicLinkCommand,
) (*valueobjects.AuthTokens, error) {
	// Check for token
	token := ""
	if cmd.Token != nil {
		token = *cmd.Token
	} else if cmd.Code != nil {
		token = *cmd.Code
	}

	if token == "" {
		return nil, fmt.Errorf("token or code is required")
	}

	// Verify magic link and create session
	return h.authService.VerifyMagicLink(ctx, token)
}
