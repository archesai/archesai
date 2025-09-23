// Package config provides application configuration management using Viper
// for environment variables, config files, and default values.
//
// Configuration is loaded from multiple sources in this order of precedence:
// 1. Environment variables (prefixed with ARCHES)
// 2. Configuration files (config.yaml)
// 3. Default values from OpenAPI specification
package config

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
	host := "0.0.0.0"
	port := float64(8080)
	dbURL := "postgres://localhost/archesai"
	redisHost := "localhost"
	redisAuth := "password"
	authEnabled := true

	return &Config{
		API: &ConfigAPI{
			Host: host,
			Port: port,
		},
		Database: &ConfigDatabase{
			URL: dbURL,
		},
		Redis: &ConfigRedis{
			Host: redisHost,
			Auth: redisAuth,
		},
		Auth: &ConfigAuth{
			Enabled: authEnabled,
		},
	}
}
