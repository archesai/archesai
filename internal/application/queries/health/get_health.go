// Package health provides health check query handlers
package health

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetHealthQuery represents a query to get the health status of the system.
type GetHealthQuery struct {
	// ComponentName is optional - if provided, returns health for specific component
	ComponentName *valueobjects.ComponentName
}

// NewGetHealthQuery creates a new query to get health status.
func NewGetHealthQuery() *GetHealthQuery {
	return &GetHealthQuery{}
}

// GetHealthQueryHandler handles the get health status query.
type GetHealthQueryHandler struct {
	// Add any infrastructure dependencies needed for health checks
	// For now, we'll return a simple healthy status
}

// NewGetHealthQueryHandler creates a new get health query handler.
func NewGetHealthQueryHandler() *GetHealthQueryHandler {
	return &GetHealthQueryHandler{}
}

// Handle executes the get health query.
func (h *GetHealthQueryHandler) Handle(
	_ context.Context,
	_ *GetHealthQuery,
) (*StatusResponse, error) {
	response := NewStatusResponse()

	// For now, return a simple healthy status
	// In a real implementation, you would check various components
	healthCheck := &valueobjects.HealthCheckResult{
		Status: "healthy",
	}

	response.WithComponentResult(healthCheck)

	return response, nil
}

// StatusResponse represents the response to a health status query
type StatusResponse struct {
	// For single component check
	ComponentResult *valueobjects.HealthCheckResult

	// For full system check
	AggregatedResult *valueobjects.AggregatedHealthCheckResult
}

// NewStatusResponse creates a new health status response.
func NewStatusResponse() *StatusResponse {
	return &StatusResponse{}
}

// WithComponentResult sets a single component result.
func (r *StatusResponse) WithComponentResult(
	result *valueobjects.HealthCheckResult,
) *StatusResponse {
	r.ComponentResult = result
	return r
}

// WithAggregatedResult sets the aggregated result.
func (r *StatusResponse) WithAggregatedResult(
	result *valueobjects.AggregatedHealthCheckResult,
) *StatusResponse {
	r.AggregatedResult = result
	return r
}

// IsHealthy returns true if the system/component is healthy.
func (r *StatusResponse) IsHealthy() bool {
	if r.ComponentResult != nil {
		return r.ComponentResult.IsHealthy()
	}
	if r.AggregatedResult != nil {
		return r.AggregatedResult.IsHealthy()
	}
	return false
}

// IsOperational returns true if the system/component is operational.
func (r *StatusResponse) IsOperational() bool {
	if r.ComponentResult != nil {
		return r.ComponentResult.IsOperational()
	}
	if r.AggregatedResult != nil {
		return r.AggregatedResult.IsOperational()
	}
	return false
}
