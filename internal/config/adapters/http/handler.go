// Package http provides HTTP handlers for configuration operations
package http

import (
	"context"
	"log/slog"

	"github.com/archesai/archesai/internal/config"
)

// ConfigHandler handles configuration operations
type ConfigHandler struct {
	config *config.Config
	logger *slog.Logger
}

// NewConfigHandler creates a new config handler
func NewConfigHandler(cfg *config.Config, logger *slog.Logger) *ConfigHandler {
	return &ConfigHandler{
		config: cfg,
		logger: logger,
	}
}

// GetConfig implements the get configuration endpoint
func (h *ConfigHandler) GetConfig(_ context.Context, _ GetConfigRequestObject) (GetConfigResponseObject, error) {
	h.logger.Debug("config requested")

	// Return the current configuration
	return GetConfig200JSONResponse(*h.config.ArchesConfig), nil
}
