// Package health provides health check and readiness probe functionality
// for monitoring service status and dependencies.
//
// The package includes:
// - Liveness checks for the application
// - Readiness checks for dependencies (database, redis, storage)
// - Health check HTTP endpoints
// - Dependency status aggregation
// - Custom health check registration
package health

import "time"

// Health check status constants
const (
	// StatusHealthy indicates the service is healthy
	StatusHealthy = "healthy"

	// StatusUnhealthy indicates the service is unhealthy
	StatusUnhealthy = "unhealthy"

	// StatusDegraded indicates the service is degraded but operational
	StatusDegraded = "degraded"
)

// Health check configuration
const (
	// DefaultTimeout is the default health check timeout
	DefaultTimeout = 5 * time.Second

	// DefaultInterval is the default health check interval
	DefaultInterval = 30 * time.Second

	// DefaultRetries is the default number of retries
	DefaultRetries = 3

	// DefaultRetryDelay is the delay between retries
	DefaultRetryDelay = 1 * time.Second
)

// Health check endpoints
const (
	// LivenessPath is the liveness probe endpoint
	LivenessPath = "/health/live"

	// ReadinessPath is the readiness probe endpoint
	ReadinessPath = "/health/ready"

	// HealthPath is the combined health status endpoint
	HealthPath = "/health"
)

// Component names for health checks
const (
	// ComponentAPI is the API server component
	ComponentAPI = "api"

	// ComponentDatabase is the database component
	ComponentDatabase = "database"

	// ComponentRedis is the Redis component
	ComponentRedis = "redis"

	// ComponentStorage is the storage component
	ComponentStorage = "storage"

	// ComponentWorker is the background worker component
	ComponentWorker = "worker"
)
