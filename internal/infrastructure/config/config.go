// Package config provides application configuration management using Viper
// for environment variables, config files, and default values.
//
// Configuration is loaded from multiple sources in this order of precedence:
// 1. Environment variables (prefixed with ARCHES)
// 2. Configuration files (config.yaml)
// 3. Default values from OpenAPI specification
package config

import (
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// Type aliases for convenience
type (
	Config          = valueobjects.Config
	ConfigAPI       = valueobjects.ConfigAPI
	ConfigAuth      = valueobjects.ConfigAuth
	ConfigAuthLocal = valueobjects.ConfigAuthLocal
	ConfigDatabase  = valueobjects.ConfigDatabase
	ConfigRedis     = valueobjects.ConfigRedis
	ConfigLogging   = valueobjects.ConfigLogging
)

// Configuration constants.
const (
	// DefaultConfigName is the default config file name.
	DefaultConfigName = "config"

	// DefaultConfigType is the default config file type.
	DefaultConfigType = "yaml"

	// EnvPrefix is the environment variable prefix.
	EnvPrefix = "ARCHES"
)

// ConfigPaths defines the search paths for configuration files.
var ConfigPaths = []string{
	".",
	"/etc/archesai/",
	"$HOME/.config/archesai",
}

// New returns the default configuration.
func New() *Config {
	return &Config{
		API: &ConfigAPI{
			Host:        "0.0.0.0",
			Port:        8080,
			Cors:        "*",
			Docs:        true,
			Environment: "development",
			Validate:    true,
		},
		Database: &ConfigDatabase{
			Enabled:       true,
			URL:           "postgres://localhost/archesai",
			Type:          "postgresql",
			MaxConns:      25,
			MinConns:      5,
			RunMigrations: true,
		},
		Redis: &ConfigRedis{
			Enabled: true,
			Host:    "localhost",
			Port:    6379,
			Auth:    "password",
		},
		Auth: &ConfigAuth{
			Enabled: true,
			Local: &ConfigAuthLocal{
				Enabled:         true,
				JWTSecret:       "change-me-in-production",
				AccessTokenTTL:  "15m",
				RefreshTokenTTL: "7d",
			},
		},
		Logging: &ConfigLogging{
			Level:  "info",
			Pretty: true,
		},
	}
}
