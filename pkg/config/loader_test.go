package config

import (
	"os"
	"testing"

	"github.com/archesai/archesai/pkg/config/schemas"
)

func TestParserLoad(t *testing.T) {
	tests := []struct {
		name       string
		setup      func()
		cleanup    func()
		configFile string // If set, use LoadFrom with this file path
		wantErr    bool
		check      func(*testing.T, *Configuration[schemas.Config])
	}{
		{
			name: "load with defaults",
			setup: func() {
				// Clear existing env vars
				_ = os.Unsetenv("ARCHES_API_HOST")
				_ = os.Unsetenv("ARCHES_API_PORT")
			},
			cleanup: func() {},
			wantErr: false,
			check: func(t *testing.T, c *Configuration[schemas.Config]) {
				t.Helper()
				if c == nil {
					t.Fatal("expected non-nil Configuration")
				}
				if c.Config == nil {
					t.Fatal("expected non-nil Config")
				}
				// With no config file and no env vars, we get zero values
			},
		},
		// NOTE: Parser doesn't read environment variables.
		// Use EnvLoader for environment variable tests.
		{
			name: "load with config file",
			setup: func() {
				// Create a temporary config file with unique name to avoid conflict with arches.yaml
				configData := `
api:
  host: "192.168.1.1"
  port: 9090
  environment: "production"
database:
  type: "postgresql"
  maxConns: 50
  url: "postgres://localhost/test"
logging:
  level: "debug"
  pretty: true
`
				err := os.WriteFile("test_config.yaml", []byte(configData), 0644)
				if err != nil {
					panic("Failed to create test config file: " + err.Error())
				}
			},
			cleanup: func() {
				_ = os.Remove("test_config.yaml")
			},
			configFile: "test_config.yaml",
			wantErr:    false,
			check: func(t *testing.T, c *Configuration[schemas.Config]) {
				t.Helper()
				if c == nil {
					t.Fatal("expected non-nil Configuration")
				}
				if c.Config == nil {
					t.Fatal("expected non-nil Config")
				}
				if c.Config.API.Host != "192.168.1.1" {
					t.Errorf("API.Host = %q, want %q", c.Config.API.Host, "192.168.1.1")
				}
				if c.Config.API.Port != 9090 {
					t.Errorf("API.Port = %d, want %d", c.Config.API.Port, 9090)
				}
				if c.Config.API.Environment != schemas.ConfigAPIEnvironmentProduction {
					t.Errorf(
						"API.Environment = %v, want %v",
						c.Config.API.Environment,
						schemas.ConfigAPIEnvironmentProduction,
					)
				}
				if c.Config.Database.MaxConns != 50 {
					t.Errorf("Database.MaxConns = %d, want %d", c.Config.Database.MaxConns, 50)
				}
				if c.Config.Logging.Level != schemas.ConfigLoggingLevelDebug {
					t.Errorf(
						"Logging.Level = %v, want %v",
						c.Config.Logging.Level,
						schemas.ConfigLoggingLevelDebug,
					)
				}
				if !c.Config.Logging.Pretty {
					t.Error("Logging.Pretty = false, want true")
				}
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

			parser := NewParser[schemas.Config]()
			var config *Configuration[schemas.Config]
			var err error
			if tt.configFile != "" {
				config, err = parser.LoadFrom(tt.configFile)
			} else {
				config, err = parser.Load()
			}

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.check != nil {
				tt.check(t, config)
			}
		})
	}
}

func TestConfigConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"DefaultConfigName", DefaultConfigName, "config"},
		{"DefaultConfigType", DefaultConfigType, "yaml"},
		{"EnvPrefix", EnvPrefix, "ARCHES"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestConfigPaths(t *testing.T) {
	expectedPaths := []string{
		".",
		"/etc/archesai/",
		"$HOME/.config/archesai",
	}

	if len(ConfigPaths) != len(expectedPaths) {
		t.Fatalf("len(ConfigPaths) = %d, want %d", len(ConfigPaths), len(expectedPaths))
	}
	for i, path := range ConfigPaths {
		if path != expectedPaths[i] {
			t.Errorf("ConfigPaths[%d] = %q, want %q", i, path, expectedPaths[i])
		}
	}
}

func TestConfigFileNames(t *testing.T) {
	expectedNames := []string{
		"arches",
		"config",
		".archesai",
	}

	if len(ConfigFileNames) != len(expectedNames) {
		t.Fatalf("len(ConfigFileNames) = %d, want %d", len(ConfigFileNames), len(expectedNames))
	}
	for i, name := range ConfigFileNames {
		if name != expectedNames[i] {
			t.Errorf("ConfigFileNames[%d] = %q, want %q", i, name, expectedNames[i])
		}
	}
}
