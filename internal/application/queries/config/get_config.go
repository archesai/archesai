package config

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetConfigQuery represents a query to get the configuration.
type GetConfigQuery struct{}

// NewGetConfigQuery creates a new get config query.
func NewGetConfigQuery() *GetConfigQuery {
	return &GetConfigQuery{}
}

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
	ctx context.Context,
	query *GetConfigQuery,
) (*valueobjects.Config, error) {
	// Return the injected configuration
	return h.config, nil
}
