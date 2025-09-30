// Package auth provides query handlers for authentication operations.
package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/infrastructure/auth"
)

// OAuthAuthorizeQueryHandler handles the OAuth authorize query.
type OAuthAuthorizeQueryHandler struct {
	authService *auth.Service
}

// NewOAuthAuthorizeQueryHandler creates a new OAuth authorize query handler.
func NewOAuthAuthorizeQueryHandler(authService *auth.Service) *OAuthAuthorizeQueryHandler {
	return &OAuthAuthorizeQueryHandler{
		authService: authService,
	}
}

// Handle executes the OAuth authorize query.
func (h *OAuthAuthorizeQueryHandler) Handle(
	_ context.Context,
	query *OAuthAuthorizeQuery,
) (string, error) {
	if query.Provider == "" {
		return "", fmt.Errorf("provider is required")
	}

	// Generate state token for CSRF protection
	state := uuid.New().String()

	// Get authorization URL
	authURL, err := h.authService.GetOAuthAuthorizationURL(query.Provider, state)
	if err != nil {
		return "", fmt.Errorf("failed to get authorization URL: %w", err)
	}

	return authURL, nil
}
