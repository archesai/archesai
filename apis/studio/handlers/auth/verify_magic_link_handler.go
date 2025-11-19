package auth

import (
	"context"
	"fmt"

	commands "github.com/archesai/archesai/apis/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/pkg/auth"
)

// VerifyMagicLinkCommandHandler handles magic link verification commands.
type VerifyMagicLinkCommandHandler struct {
	authService *auth.Service
}

// NewVerifyMagicLinkCommandHandler creates a new magic link verification command handler.
func NewVerifyMagicLinkCommandHandler(
	authService *auth.Service,
) *VerifyMagicLinkCommandHandler {
	return &VerifyMagicLinkCommandHandler{
		authService: authService,
	}
}

// Handle executes the magic link verification command.
func (h *VerifyMagicLinkCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.VerifyMagicLinkCommand,
) (*auth.AuthTokens, error) {
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
