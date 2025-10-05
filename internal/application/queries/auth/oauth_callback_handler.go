package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/core/events"
	"github.com/archesai/archesai/internal/core/services"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// OAuthCallbackQueryHandler handles the OAuth callback query.
type OAuthCallbackQueryHandler struct {
	authService services.AuthService
	publisher   events.Publisher
}

// NewOAuthCallbackQueryHandler creates a new OAuth callback query handler.
func NewOAuthCallbackQueryHandler(
	authService services.AuthService,
	publisher events.Publisher,
) *OAuthCallbackQueryHandler {
	return &OAuthCallbackQueryHandler{
		authService: authService,
		publisher:   publisher,
	}
}

// Handle executes the OAuth callback query.
func (h *OAuthCallbackQueryHandler) Handle(
	_ context.Context,
	query *OAuthCallbackQuery,
) (*valueobjects.AuthTokens, error) {
	if query.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}

	// TODO: Get code and state from query params when controller is fixed to pass them
	// For now, the query only has Provider from the controller

	// This will need to be fixed when the controller properly passes query params
	return nil, fmt.Errorf(
		"OAuth callback not fully implemented - controller needs to pass code and state",
	)
}
