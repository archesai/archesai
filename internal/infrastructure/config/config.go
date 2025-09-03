// Package config provides application configuration management.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/generated/api"
	"github.com/spf13/viper"
)

// Config wraps the generated ArchesConfig with simplified access methods
type Config struct {
	v *viper.Viper

	// Simplified fields for easy access
	Server   ServerConfig
	Auth     AuthConfig
	Database api.DatabaseConfig
	Logging  LoggingConfig
}

// ServerConfig holds HTTP server configuration options.
type ServerConfig struct {
	Host        string
	Port        int
	DocsEnabled bool
	Environment string
}

// AuthConfig holds authentication and JWT configuration.
type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
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

	// Build database config
	dbType := api.DatabaseConfigType(v.GetString("database.type"))
	maxConns := v.GetInt("database.max_conns")
	minConns := v.GetInt("database.min_conns")
	runMigrations := v.GetBool("database.run_migrations")

	databaseConfig := api.DatabaseConfig{
		Enabled:       true,
		Url:           v.GetString("database.url"),
		Type:          &dbType,
		MaxConns:      &maxConns,
		MinConns:      &minConns,
		RunMigrations: &runMigrations,
	}

	if lifetime := v.GetDuration("database.conn_max_lifetime"); lifetime > 0 {
		lifetimeStr := lifetime.String()
		databaseConfig.ConnMaxLifetime = &lifetimeStr
	}

	if idleTime := v.GetDuration("database.conn_max_idle_time"); idleTime > 0 {
		idleTimeStr := idleTime.String()
		databaseConfig.ConnMaxIdleTime = &idleTimeStr
	}

	if period := v.GetDuration("database.health_check_period"); period > 0 {
		periodStr := period.String()
		databaseConfig.HealthCheckPeriod = &periodStr
	}

	return &Config{
		v: v,
		Server: ServerConfig{
			Host:        v.GetString("api.host"),
			Port:        v.GetInt("api.port"),
			DocsEnabled: v.GetBool("api.docs"),
			Environment: v.GetString("api.environment"),
		},
		Auth: AuthConfig{
			JWTSecret:       v.GetString("auth.jwt_secret"),
			AccessTokenTTL:  v.GetDuration("auth.access_token_ttl"),
			RefreshTokenTTL: v.GetDuration("auth.refresh_token_ttl"),
		},
		Database: databaseConfig,
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
	v.SetDefault("api.environment", "development")

	// Auth defaults
	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.access_token_ttl", "15m")
	v.SetDefault("auth.refresh_token_ttl", "168h") // 7 days

	// Database defaults
	v.SetDefault("database.url", "postgres://postgres:postgres@localhost:5432/archesai?sslmode=disable")
	v.SetDefault("database.type", "postgresql")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.conn_max_idle_time", "30m")
	v.SetDefault("database.health_check_period", "30s")
	v.SetDefault("database.run_migrations", false)

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
