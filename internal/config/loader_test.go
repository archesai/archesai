package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		cleanup func()
		wantErr bool
		check   func(*Config) error
	}{
		{
			name: "load with defaults",
			setup: func() {
				// Clear any existing env vars
				_ = os.Unsetenv("ARCHESAI_API_HOST")
				_ = os.Unsetenv("ARCHESAI_API_PORT")
			},
			cleanup: func() {},
			wantErr: false,
			check: func(c *Config) error {
				if c == nil {
					return errorf("expected non-nil config")
				}
				if c.API.Host != "0.0.0.0" {
					return errorf("expected default host 0.0.0.0, got %s", c.API.Host)
				}
				return nil
			},
		},
		{
			name: "load with environment variables",
			setup: func() {
				_ = os.Setenv("ARCHESAI_API_HOST", "127.0.0.1")
				_ = os.Setenv("ARCHESAI_API_PORT", "8080")
				_ = os.Setenv("ARCHESAI_DATABASE_URL", "postgres://test")
				_ = os.Setenv("ARCHESAI_AUTH_ENABLED", "false")
			},
			cleanup: func() {
				_ = os.Unsetenv("ARCHESAI_API_HOST")
				_ = os.Unsetenv("ARCHESAI_API_PORT")
				_ = os.Unsetenv("ARCHESAI_DATABASE_URL")
				_ = os.Unsetenv("ARCHESAI_AUTH_ENABLED")
			},
			wantErr: false,
			check: func(c *Config) error {
				if c.API.Host != "127.0.0.1" {
					return errorf("expected host 127.0.0.1, got %s", c.API.Host)
				}
				if c.API.Port != 8080 {
					return errorf("expected port 8080, got %f", c.API.Port)
				}
				if c.Database.URL != "postgres://test" {
					return errorf("expected database URL postgres://test, got %s", c.Database.URL)
				}
				if c.Auth.Enabled != false {
					return errorf("expected auth disabled")
				}
				return nil
			},
		},
		{
			name: "load with config file",
			setup: func() {
				// Create a temporary config file
				configData := `
api:
  host: "192.168.1.1"
  port: 9090
  environment: "production"
database:
  type: "postgres"
  max_conns: 50
logging:
  level: "debug"
  pretty: true
`
				err := os.WriteFile("config.yaml", []byte(configData), 0644)
				if err != nil {
					t.Fatalf("Failed to create test config file: %v", err)
				}
			},
			cleanup: func() {
				_ = os.Remove("config.yaml")
			},
			wantErr: false,
			check: func(c *Config) error {
				if c.API.Host != "192.168.1.1" {
					return errorf("expected host 192.168.1.1, got %s", c.API.Host)
				}
				if c.API.Port != 9090 {
					return errorf("expected port 9090, got %f", c.API.Port)
				}
				if c.API.Environment != "production" {
					return errorf("expected environment production, got %s", c.API.Environment)
				}
				if c.Database.MaxConns != 50 {
					return errorf("expected max_conns 50, got %d", c.Database.MaxConns)
				}
				if c.Logging.Level != "debug" {
					return errorf("expected log level debug, got %s", c.Logging.Level)
				}
				if !c.Logging.Pretty {
					return errorf("expected pretty logging enabled")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			config, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.check != nil {
				if err := tt.check(config); err != nil {
					t.Errorf("Load() check failed: %v", err)
				}
			}
		})
	}
}

func TestSetupViper(t *testing.T) {
	v := viper.New()
	setupViper(v)

	// Check that config paths are set
	if v.ConfigFileUsed() == "" {
		// This is expected as no config file exists yet
		t.Log("No config file found (expected)")
	}

	// Check environment prefix
	_ = os.Setenv("ARCHESAI_TEST_VAR", "test_value")
	defer func() { _ = os.Unsetenv("ARCHESAI_TEST_VAR") }()

	if v.GetString("test.var") != "test_value" {
		t.Error("Viper not reading environment variables correctly")
	}
}

func TestApplyAPIOverrides(t *testing.T) {
	config := &ArchesConfig{
		API: APIConfig{
			Host: "0.0.0.0",
			Port: 3001,
			Cors: CORSConfig{
				Origins: "*",
			},
		},
	}

	v := viper.New()
	v.Set("api.host", "localhost")
	v.Set("api.port", 8080)
	v.Set("api.docs", true)
	v.Set("api.validate", true)
	v.Set("api.environment", "production")
	v.Set("api.cors.origins", "https://example.com")

	applyAPIOverrides(config, v)

	if config.API.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", config.API.Host)
	}
	if config.API.Port != 8080 {
		t.Errorf("Expected port 8080, got %f", config.API.Port)
	}
	if !config.API.Docs {
		t.Error("Expected docs enabled")
	}
	if !config.API.Validate {
		t.Error("Expected validate enabled")
	}
	if config.API.Environment != "production" {
		t.Errorf("Expected environment production, got %s", config.API.Environment)
	}
	if config.API.Cors.Origins != "https://example.com" {
		t.Errorf("Expected CORS origins https://example.com, got %s", config.API.Cors.Origins)
	}
}

func TestApplyAuthOverrides(t *testing.T) {
	config := &ArchesConfig{
		Auth: AuthConfig{
			Enabled: true,
			Local: LocalAuth{
				Enabled:         true,
				JwtSecret:       "default-secret",
				AccessTokenTTL:  "15m",
				RefreshTokenTTL: "7d",
			},
		},
	}

	v := viper.New()
	v.Set("auth.enabled", false)
	v.Set("auth.local.enabled", false)
	v.Set("auth.local.jwt_secret", "new-secret")
	v.Set("auth.local.access_token_ttl", "30m")
	v.Set("auth.local.refresh_token_ttl", "14d")

	applyAuthOverrides(config, v)

	if config.Auth.Enabled != false {
		t.Error("Expected auth disabled")
	}
	if config.Auth.Local.Enabled != false {
		t.Error("Expected local auth disabled")
	}
	if config.Auth.Local.JwtSecret != "new-secret" {
		t.Errorf("Expected JWT secret new-secret, got %s", config.Auth.Local.JwtSecret)
	}
	if config.Auth.Local.AccessTokenTTL != "30m" {
		t.Errorf("Expected access token TTL 30m, got %s", config.Auth.Local.AccessTokenTTL)
	}
	if config.Auth.Local.RefreshTokenTTL != "14d" {
		t.Errorf("Expected refresh token TTL 14d, got %s", config.Auth.Local.RefreshTokenTTL)
	}
}

func TestApplyDatabaseOverrides(t *testing.T) {
	config := &ArchesConfig{
		Database: DatabaseConfig{
			Enabled:           true,
			URL:               "postgres://localhost/test",
			Type:              "postgres",
			MaxConns:          25,
			MinConns:          5,
			ConnMaxLifetime:   "1h",
			ConnMaxIdleTime:   "10m",
			HealthCheckPeriod: "30s",
			RunMigrations:     false,
		},
	}

	v := viper.New()
	v.Set("database.enabled", false)
	v.Set("database.url", "postgres://remote/prod")
	v.Set("database.type", "sqlite")
	v.Set("database.max_conns", 100)
	v.Set("database.min_conns", 10)
	v.Set("database.conn_max_lifetime", "2h")
	v.Set("database.conn_max_idle_time", "20m")
	v.Set("database.health_check_period", "60s")
	v.Set("database.run_migrations", true)

	applyDatabaseOverrides(config, v)

	if config.Database.Enabled != false {
		t.Error("Expected database disabled")
	}
	if config.Database.URL != "postgres://remote/prod" {
		t.Errorf("Expected URL postgres://remote/prod, got %s", config.Database.URL)
	}
	if config.Database.Type != "sqlite" {
		t.Errorf("Expected type sqlite, got %s", config.Database.Type)
	}
	if config.Database.MaxConns != 100 {
		t.Errorf("Expected max conns 100, got %d", config.Database.MaxConns)
	}
	if config.Database.MinConns != 10 {
		t.Errorf("Expected min conns 10, got %d", config.Database.MinConns)
	}
	if config.Database.ConnMaxLifetime != "2h" {
		t.Errorf("Expected conn max lifetime 2h, got %s", config.Database.ConnMaxLifetime)
	}
	if config.Database.ConnMaxIdleTime != "20m" {
		t.Errorf("Expected conn max idle time 20m, got %s", config.Database.ConnMaxIdleTime)
	}
	if config.Database.HealthCheckPeriod != "60s" {
		t.Errorf("Expected health check period 60s, got %s", config.Database.HealthCheckPeriod)
	}
	if !config.Database.RunMigrations {
		t.Error("Expected run migrations enabled")
	}
}

func TestApplyLoggingOverrides(t *testing.T) {
	config := &ArchesConfig{
		Logging: LoggingConfig{
			Level:  "info",
			Pretty: false,
		},
	}

	v := viper.New()
	v.Set("logging.level", "debug")
	v.Set("logging.pretty", true)

	applyLoggingOverrides(config, v)

	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level debug, got %s", config.Logging.Level)
	}
	if !config.Logging.Pretty {
		t.Error("Expected pretty logging enabled")
	}
}

func TestConfigConstants(t *testing.T) {
	// Test that constants have expected values
	if DefaultConfigName != "config" {
		t.Errorf("Expected DefaultConfigName to be 'config', got %s", DefaultConfigName)
	}
	if DefaultConfigType != "yaml" {
		t.Errorf("Expected DefaultConfigType to be 'yaml', got %s", DefaultConfigType)
	}
	if EnvPrefix != "ARCHESAI" {
		t.Errorf("Expected EnvPrefix to be 'ARCHESAI', got %s", EnvPrefix)
	}
}

func TestConfigPaths(t *testing.T) {
	// Test that config paths are set correctly
	expectedPaths := []string{
		".",
		"/etc/archesai/",
		"$HOME/.config/archesai",
	}

	if len(ConfigPaths) != len(expectedPaths) {
		t.Errorf("Expected %d config paths, got %d", len(expectedPaths), len(ConfigPaths))
	}

	for i, path := range ConfigPaths {
		if path != expectedPaths[i] {
			t.Errorf("Expected config path %d to be %s, got %s", i, expectedPaths[i], path)
		}
	}
}

// Helper function to create formatted errors
func errorf(format string, args ...interface{}) error {
	if len(args) == 0 {
		return &testError{msg: format}
	}
	return &testError{msg: format, args: args}
}

type testError struct {
	msg  string
	args []interface{}
}

func (e *testError) Error() string {
	if len(e.args) == 0 {
		return e.msg
	}
	// Simple sprintf implementation for testing
	result := e.msg
	for _, arg := range e.args {
		result = replaceFirst(result, "%", fmt.Sprintf("%v", arg))
	}
	return result
}

// Helper function for string replacement
func replaceFirst(s, old, replacement string) string {
	if idx := strings.Index(s, old); idx != -1 {
		return s[:idx] + replacement + s[idx+len(old):]
	}
	return s
}
