package config

import (
	"fmt"
	"os"
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
				_ = os.Unsetenv("ARCHES_API_HOST")
				_ = os.Unsetenv("ARCHES_API_PORT")
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
				_ = os.Setenv("ARCHES_API_HOST", "127.0.0.1")
				_ = os.Setenv("ARCHES_API_PORT", "8080")
				_ = os.Setenv("ARCHES_DATABASE_URL", "postgres://test")
				_ = os.Setenv("ARCHES_AUTH_ENABLED", "false")
			},
			cleanup: func() {
				_ = os.Unsetenv("ARCHES_API_HOST")
				_ = os.Unsetenv("ARCHES_API_PORT")
				_ = os.Unsetenv("ARCHES_DATABASE_URL")
				_ = os.Unsetenv("ARCHES_AUTH_ENABLED")
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
  maxConns: 50
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
	_ = os.Setenv("ARCHES_TEST_VAR", "test_value")
	defer func() { _ = os.Unsetenv("ARCHES_TEST_VAR") }()

	if v.GetString("test.var") != "test_value" {
		t.Error("Viper not reading environment variables correctly")
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
	if EnvPrefix != "ARCHES" {
		t.Errorf("Expected EnvPrefix to be 'ARCHES', got %s", EnvPrefix)
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

// Helper function to create formatted errors.
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
	// Use fmt.Sprintf for proper formatting
	return fmt.Sprintf(e.msg, e.args...)
}
