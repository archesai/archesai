// Package http provides HTTP handlers for health checks
package http

import (
	"context"
	"log/slog"
	"time"

	"github.com/archesai/archesai/internal/health"
)

// HealthHandler handles health check operations
type HealthHandler struct {
	service *health.Service
	logger  *slog.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(service *health.Service, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		service: service,
		logger:  logger,
	}
}

// GetHealth implements the health check endpoint
func (h *HealthHandler) GetHealth(ctx context.Context, _ GetHealthRequestObject) (GetHealthResponseObject, error) {
	h.logger.Debug("health check requested")

	status := h.service.CheckHealth(ctx)

	response := health.HealthResponse{
		Services: struct {
			Database string `json:"database" yaml:"database"`
			Email    string `json:"email" yaml:"email"`
			Redis    string `json:"redis" yaml:"redis"`
		}{
			Database: status.Database,
			Email:    status.Email,
			Redis:    status.Redis,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    float32(status.Uptime),
	}

	return GetHealth200JSONResponse(response), nil
}
