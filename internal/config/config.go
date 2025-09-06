// Package config provides application configuration management using Viper
// for environment variables, config files, and default values.
//
// Configuration is loaded from multiple sources in this order of precedence:
// 1. Environment variables (prefixed with ARCHESAI_)
// 2. Configuration files (config.yaml)
// 3. Default values from OpenAPI specification
package config

// Generate config types from OpenAPI specification
//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package config --include-tags Config ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package config --include-tags Config ../../api/openapi.bundled.yaml

// Configuration constants
const (
	// DefaultConfigName is the default config file name
	DefaultConfigName = "config"

	// DefaultConfigType is the default config file type
	DefaultConfigType = "yaml"

	// EnvPrefix is the environment variable prefix
	EnvPrefix = "ARCHESAI"

	// DefaultPort is the default API server port
	DefaultPort = 3001

	// DefaultHost is the default API server host
	DefaultHost = "0.0.0.0"
)

// ConfigPaths defines the search paths for configuration files
var ConfigPaths = []string{
	".",
	"/etc/archesai/",
	"$HOME/.config/archesai",
}
