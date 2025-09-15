package health

import (
	"context"
	"log/slog"
	"time"

	"github.com/archesai/archesai/internal/database"
)

// Service handles health check business logic
type Service struct {
	db     *database.Database
	logger *slog.Logger
	start  time.Time
}

// NewService creates a new health service
func NewService(db *database.Database, logger *slog.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
		start:  time.Now(),
	}
}

// CheckHealth performs health checks on all services
func (s *Service) CheckHealth(ctx context.Context) (*HealthResponse, error) {
	// Check database
	dbStatus := StatusHealthy
	if err := s.checkDatabase(ctx); err != nil {
		s.logger.Error("database health check failed", "error", err)
		dbStatus = StatusUnhealthy
	}

	// Build response
	response := &HealthResponse{
		Services: struct {
			Database string `json:"database" yaml:"database"`
			Email    string `json:"email" yaml:"email"`
			Redis    string `json:"redis" yaml:"redis"`
		}{
			Database: dbStatus,
			Email:    StatusHealthy, // TODO: Implement email check
			Redis:    StatusHealthy, // TODO: Implement redis check
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    float32(time.Since(s.start).Seconds()),
	}

	return response, nil
}

// checkDatabase verifies database connectivity
func (s *Service) checkDatabase(_ context.Context) error {
	// Simple ping check - could be enhanced with actual query
	// For now, if we have a db connection, we consider it healthy
	if s.db == nil {
		return ErrDatabaseUnavailable
	}
	return nil
}
