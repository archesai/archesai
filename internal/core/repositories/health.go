// Package repositories provides repository interfaces and implementations
package repositories

import (
	"context"
)

// HealthCheckRepository defines the interface for health check persistence and infrastructure checks.
// This interface is implemented by the infrastructure layer to provide health status of external dependencies.
type HealthCheckRepository interface {
	// CheckDatabase performs a health check on the database connection.
	// Returns nil if the database is healthy, error otherwise.
	CheckDatabase(ctx context.Context) error

	// CheckRedis performs a health check on Redis connection.
	// Returns nil if Redis is healthy, error otherwise.
	CheckRedis(ctx context.Context) error

	// CheckEmail performs a health check on the email service.
	// Returns nil if the email service is healthy, error otherwise.
	CheckEmail(ctx context.Context) error
}
