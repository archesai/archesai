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

// ConfigFileNames defines the configuration file names to search for (in order of priority).
var ConfigFileNames = []string{
	"arches",    // arches.yaml
	"config",    // config.yaml
	".archesai", // .archesai.yaml
}

// Config is a constraint interface for configuration types.
// Any struct with a Validate method can be used as a configuration.
type Config interface {
	// Validate checks if the configuration is valid.
	// Validate() error
}

// Parser handles loading and parsing of configuration.
type Parser[C Config] struct{}

// NewParser creates a new configuration parser for type C.
func NewParser[C Config]() *Parser[C] {
	return &Parser[C]{}
}
