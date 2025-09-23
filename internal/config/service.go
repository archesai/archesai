package config

import (
	"context"
	"log/slog"
)

// Service handles configuration business logic.
type Service struct {
	config *Config
	logger *slog.Logger
}

// NewService creates a new config service.
func NewService(config *Config, logger *slog.Logger) *Service {
	return &Service{
		config: config,
		logger: logger,
	}
}

// GetConfig returns the current configuration (expected by generated handler).
func (s *Service) GetConfig(ctx context.Context) (*Config, error) {
	return s.config, nil
}
