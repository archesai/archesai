package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// RequestMagicLinkCommandHandler handles magic link request commands.
type RequestMagicLinkCommandHandler struct {
	authService *auth.Service
	// TODO: Add notification service for sending emails
}

// NewRequestMagicLinkCommandHandler creates a new magic link request command handler.
func NewRequestMagicLinkCommandHandler(authService *auth.Service) *RequestMagicLinkCommandHandler {
	return &RequestMagicLinkCommandHandler{
		authService: authService,
	}
}

// Handle executes the magic link request command.
func (h *RequestMagicLinkCommandHandler) Handle(
	_ context.Context,
	cmd *RequestMagicLinkCommand,
) (string, error) {
	if cmd.Identifier == "" {
		return "", fmt.Errorf("identifier (email) is required")
	}

	// Get redirect URL if provided
	redirectURL := ""
	if cmd.RedirectUrl != nil {
		redirectURL = *cmd.RedirectUrl
	}

	// Generate magic link
	link, err := h.authService.GenerateMagicLink(cmd.Identifier, redirectURL)
	if err != nil {
		return "", fmt.Errorf("failed to generate magic link: %w", err)
	}

	// TODO: Send magic link via email using notification service
	// For now, just return the link (in production, this would be sent via email)

	return link, nil
}
