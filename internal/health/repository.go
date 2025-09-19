package health

import (
	"context"
)

// Repository defines the interface for health check operations
type Repository interface {
	// CheckDatabase performs a health check on the database connection
	CheckDatabase(ctx context.Context) error

	// CheckRedis performs a health check on Redis connection (if configured)
	CheckRedis(ctx context.Context) error

	// CheckEmail performs a health check on email service (if configured)
	CheckEmail(ctx context.Context) error
}
