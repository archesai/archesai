// Package config provides application configuration management.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config wraps the generated ArchesConfig with simplified access methods
type Config struct {
	v *viper.Viper

	// Simplified fields for easy access
	Server   ServerConfig
	Auth     AuthConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

// ServerConfig holds HTTP server configuration options.
type ServerConfig struct {
	Host        string
	Port        int
	DocsEnabled bool
}

// AuthConfig holds authentication and JWT configuration.
type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	URL string
}

// LoggingConfig holds logging configuration options.
type LoggingConfig struct {
	Level  string
	Pretty bool
}

// Load reads configuration from environment variables and returns a Config.
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/archesai/")
	v.AddConfigPath("$HOME/.archesai")

	v.SetEnvPrefix("ARCHESAI")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	return &Config{
		v: v,
		Server: ServerConfig{
			Host:        v.GetString("api.host"),
			Port:        v.GetInt("api.port"),
			DocsEnabled: v.GetBool("api.docs"),
		},
		Auth: AuthConfig{
			JWTSecret:       v.GetString("auth.jwt_secret"),
			AccessTokenTTL:  v.GetDuration("auth.access_token_ttl"),
			RefreshTokenTTL: v.GetDuration("auth.refresh_token_ttl"),
		},
		Database: DatabaseConfig{
			URL: v.GetString("database.url"),
		},
		Logging: LoggingConfig{
			Level:  v.GetString("logging.level"),
			Pretty: v.GetBool("logging.pretty"),
		},
	}, nil
}

func setDefaults(v *viper.Viper) {
	// API defaults
	v.SetDefault("api.host", "0.0.0.0")
	v.SetDefault("api.port", 8080)
	v.SetDefault("api.docs", true)
	v.SetDefault("api.validate", true)
	v.SetDefault("api.cors.origins", "http://localhost:3000")

	// Auth defaults
	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.access_token_ttl", "15m")
	v.SetDefault("auth.refresh_token_ttl", "168h") // 7 days

	// Database defaults
	v.SetDefault("database.url", "postgres://postgres:postgres@localhost:5432/archesai?sslmode=disable")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.pretty", false)
}

// GetAllowedOrigins returns CORS allowed origins as a slice
func (c *Config) GetAllowedOrigins() []string {
	origins := c.v.GetString("api.cors.origins")
	if origins == "" {
		return []string{}
	}
	return strings.Split(origins, ",")
}

// GetServerAddress returns the formatted server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
