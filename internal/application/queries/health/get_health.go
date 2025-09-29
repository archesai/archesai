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
	ctx context.Context,
	query *GetHealthQuery,
) (*HealthStatusResponse, error) {
	response := NewHealthStatusResponse()

	// For now, return a simple healthy status
	// In a real implementation, you would check various components
	healthCheck := &valueobjects.HealthCheckResult{
		Status: "healthy",
	}

	response.WithComponentResult(healthCheck)

	return response, nil
}

// HealthStatusResponse represents the response to a health status query.
type HealthStatusResponse struct {
	// For single component check
	ComponentResult *valueobjects.HealthCheckResult

	// For full system check
	AggregatedResult *valueobjects.AggregatedHealthCheckResult
}

// NewHealthStatusResponse creates a new health status response.
func NewHealthStatusResponse() *HealthStatusResponse {
	return &HealthStatusResponse{}
}

// WithComponentResult sets a single component result.
func (r *HealthStatusResponse) WithComponentResult(
	result *valueobjects.HealthCheckResult,
) *HealthStatusResponse {
	r.ComponentResult = result
	return r
}

// WithAggregatedResult sets the aggregated result.
func (r *HealthStatusResponse) WithAggregatedResult(
	result *valueobjects.AggregatedHealthCheckResult,
) *HealthStatusResponse {
	r.AggregatedResult = result
	return r
}

// IsHealthy returns true if the system/component is healthy.
func (r *HealthStatusResponse) IsHealthy() bool {
	if r.ComponentResult != nil {
		return r.ComponentResult.IsHealthy()
	}
	if r.AggregatedResult != nil {
		return r.AggregatedResult.IsHealthy()
	}
	return false
}

// IsOperational returns true if the system/component is operational.
func (r *HealthStatusResponse) IsOperational() bool {
	if r.ComponentResult != nil {
		return r.ComponentResult.IsOperational()
	}
	if r.AggregatedResult != nil {
		return r.AggregatedResult.IsOperational()
	}
	return false
}
