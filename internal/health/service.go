package health

import (
	"context"
	"log/slog"
	"time"
)

// Service handles health check business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
	start  time.Time
}

// NewService creates a new health service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
		start:  time.Now(),
	}
}

// CheckHealth performs health checks on all services.
func (s *Service) CheckHealth(ctx context.Context) (*HealthResponse, error) {
	// Check database
	dbStatus := StatusHealthy
	if err := s.repo.CheckDatabase(ctx); err != nil {
		s.logger.Error("database health check failed", "error", err)
		dbStatus = StatusUnhealthy
	}

	// Check Redis
	redisStatus := StatusHealthy
	if err := s.repo.CheckRedis(ctx); err != nil {
		s.logger.Error("redis health check failed", "error", err)
		redisStatus = StatusUnhealthy
	}

	// Check Email
	emailStatus := StatusHealthy
	if err := s.repo.CheckEmail(ctx); err != nil {
		s.logger.Error("email health check failed", "error", err)
		emailStatus = StatusUnhealthy
	}

	// Build response
	response := &HealthResponse{
		Services: struct {
			Database string `json:"database" yaml:"database"`
			Email    string `json:"email" yaml:"email"`
			Redis    string `json:"redis" yaml:"redis"`
		}{
			Database: dbStatus,
			Email:    emailStatus,
			Redis:    redisStatus,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    float32(time.Since(s.start).Seconds()),
	}

	return response, nil
}
