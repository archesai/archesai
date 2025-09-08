// Package health provides HTTP handlers for health checks
package health

import (
	"context"
	"log/slog"
	"time"
)

// Handler handles health check operations
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)

// NewHandler creates a new health handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// GetHealth implements the health check endpoint
func (h *Handler) GetHealth(ctx context.Context, _ GetHealthRequestObject) (GetHealthResponseObject, error) {
	h.logger.Debug("health check requested")

	status := h.service.CheckHealth(ctx)

	response := HealthResponse{
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
