// Package config provides configuration query handlers
package config

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetConfigQueryHandler handles the get config query.
type GetConfigQueryHandler struct {
	config *valueobjects.Config
}

// NewGetConfigQueryHandler creates a new get config query handler.
func NewGetConfigQueryHandler(config *valueobjects.Config) *GetConfigQueryHandler {
	return &GetConfigQueryHandler{
		config: config,
	}
}

// Handle executes the get config query.
func (h *GetConfigQueryHandler) Handle(
	_ context.Context,
	_ *GetConfigQuery,
) (*valueobjects.Config, error) {
	// Return the injected configuration
	return h.config, nil
}
