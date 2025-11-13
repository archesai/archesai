// Package config provides application configuration management using Viper
// for environment variables, config files, and default values.
//
// Configuration is loaded from multiple sources in this order of precedence:
// 1. Environment variables (prefixed with ARCHES)
// 2. Configuration files (config.yaml)
// 3. Default values from OpenAPI specification
package config

import (
	"github.com/archesai/archesai/internal/core/models"
)

// Config is the main application configuration structure
type Config = models.Config

// Type aliases for convenience - avoiding stuttering names
type (
	// API configuration
	API = models.APIConfig
	// Auth configuration
	Auth = models.AuthConfig
	// AuthLocal configuration
	AuthLocal = models.LocalAuthConfig
	// Database configuration
	Database = models.DatabaseConfig
	// Redis configuration
	Redis = models.RedisConfig
	// Logging configuration
	Logging = models.LoggingConfig
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

// ConfigFileNames defines the configuration file names to search for (in order of priority).
var ConfigFileNames = []string{
	"arches",    // arches.yaml
	"config",    // config.yaml
	".archesai", // .archesai.yaml
}

// New returns the default configuration.
func New() *Config {
	return &Config{
		API: &API{
			Host:        "0.0.0.0",
			Port:        8080,
			Cors:        "*",
			Docs:        true,
			Environment: "development",
			Validation:  true,
		},
		Database: &Database{
			Enabled:       true,
			URL:           "postgres://localhost/archesai",
			Type:          "postgresql",
			MaxConns:      25,
			MinConns:      5,
			RunMigrations: true,
		},
		Redis: &Redis{
			Enabled: true,
			Host:    "localhost",
			Port:    6379,
			Auth:    "password",
		},
		Auth: &Auth{
			Enabled: true,
			Local: &AuthLocal{
				Enabled:         true,
				JWTSecret:       "change-me-in-production",
				AccessTokenTTL:  "15m",
				RefreshTokenTTL: "7d",
			},
		},
		Logging: &Logging{
			Level:  "info",
			Pretty: true,
		},
	}
}
