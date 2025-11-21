// Package config provides handlers for configuration-related operations.
package config

import (
	"context"

	queries "github.com/archesai/archesai/apps/studio/generated/application/queries/config"
	"github.com/archesai/archesai/apps/studio/generated/core/models"
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
	_ *queries.GetConfigQuery,
) (*models.Config, error) {
	// Return the injected configuration
	return h.config, nil
}
