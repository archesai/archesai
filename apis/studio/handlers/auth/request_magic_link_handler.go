package auth

import (
	"context"
	"fmt"

	commands "github.com/archesai/archesai/apis/studio/generated/application/commands/auth"
	"github.com/archesai/archesai/pkg/auth"
)

// RequestMagicLinkCommandHandler handles magic link request commands.
type RequestMagicLinkCommandHandler struct {
	authService *auth.Service
	// TODO: Add notification service for sending emails
}

// NewRequestMagicLinkCommandHandler creates a new magic link request command handler.
func NewRequestMagicLinkCommandHandler(
	authService *auth.Service,
) *RequestMagicLinkCommandHandler {
	return &RequestMagicLinkCommandHandler{
		authService: authService,
	}
}

// Handle executes the magic link request command.
func (h *RequestMagicLinkCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.RequestMagicLinkCommand,
) (string, error) {
	if cmd.Identifier == "" {
		return "", fmt.Errorf("identifier (email) is required")
	}

	// Get redirect URL if provided
	redirectURL := ""
	if cmd.RedirectURL != nil {
		redirectURL = *cmd.RedirectURL
	}

	// Generate and send magic link
	link, err := h.authService.GenerateMagicLink(ctx, cmd.Identifier, redirectURL)
	if err != nil {
		return "", fmt.Errorf("failed to generate magic link: %w", err)
	}

	return link, nil
}
