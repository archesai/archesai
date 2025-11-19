// Package config provides configuration query handlers
package config

import (
	"context"

	"github.com/archesai/archesai/apis/studio/generated/core/models"
)

// GetConfigQueryHandler handles the get config query.
type GetConfigQueryHandler struct {
	config *models.Config
}

// NewGetConfigQueryHandler creates a new get config query handler.
func NewGetConfigQueryHandler(config *models.Config) *GetConfigQueryHandler {
	return &GetConfigQueryHandler{
		config: config,
	}
}

// Handle executes the get config query.
func (h *GetConfigQueryHandler) Handle(
	_ context.Context,
	_ *any,
) (*models.Config, error) {
	// Return the injected configuration
	return h.config, nil
}
