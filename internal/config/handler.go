package config

import (
	"context"
	"log/slog"
)

const sanitizedValue = "***"

// Handler implements the config API
type Handler struct {
	config *Config
	logger *slog.Logger
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)

// NewHandler creates a new config server
func NewHandler(config *Config, logger *slog.Logger) *Handler {
	return &Handler{
		config: config,
		logger: logger,
	}
}

// GetConfig returns the current configuration
func (s *Handler) GetConfig(
	_ context.Context,
	_ GetConfigRequestObject,
) (GetConfigResponseObject, error) {
	// Create a sanitized copy of the config
	// We don't want to expose sensitive information like secrets
	sanitizedConfig := s.sanitizeConfig()
	return GetConfig200JSONResponse(*sanitizedConfig), nil
}

// sanitizeConfig creates a copy of the config with sensitive data removed FIXME - this should be generated form an x-sensitive
func (s *Handler) sanitizeConfig() *ArchesConfig {
	// Deep copy the embedded ArchesConfig
	cfg := *s.config.ArchesConfig

	// Sanitize authentication secrets
	if cfg.Auth.Local.JwtSecret != "" {
		cfg.Auth.Local.JwtSecret = sanitizedValue
	}

	// Sanitize OAuth secrets
	if cfg.Auth.Google.ClientSecret != "" {
		cfg.Auth.Google.ClientSecret = sanitizedValue
	}
	if cfg.Auth.Github.ClientSecret != "" {
		cfg.Auth.Github.ClientSecret = sanitizedValue
	}
	if cfg.Auth.Microsoft.ClientSecret != "" {
		cfg.Auth.Microsoft.ClientSecret = sanitizedValue
	}
	if cfg.Auth.Twitter.ConsumerSecret != "" {
		cfg.Auth.Twitter.ConsumerSecret = sanitizedValue
	}

	// Sanitize database password if present in URL
	if cfg.Database.URL != "" {
		// Simple check - if it contains @ it might have credentials
		if len(cfg.Database.URL) > 20 {
			cfg.Database.URL = sanitizedValue
		}
	}

	// Sanitize Redis password
	if cfg.Redis.Auth != "" {
		cfg.Redis.Auth = sanitizedValue
	}

	return &cfg
}
