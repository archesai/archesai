package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/archesai/archesai/pkg/config/models"
)

func TestParserLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		cleanup func()
		wantErr bool
		check   func(*testing.T, *Configuration[models.Config])
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
			check: func(t *testing.T, c *Configuration[models.Config]) {
				require.NotNil(t, c)
				require.NotNil(t, c.Config)
				// With no config file and no env vars, we get zero values
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
			check: func(t *testing.T, c *Configuration[models.Config]) {
				require.NotNil(t, c)
				require.NotNil(t, c.Config)
				if c.Config.API != nil {
					assert.Equal(t, "127.0.0.1", c.Config.API.Host)
					assert.Equal(t, int32(8080), c.Config.API.Port)
				}
				if c.Config.Database != nil {
					assert.Equal(t, "postgres://test", c.Config.Database.URL)
				}
				if c.Config.Auth != nil {
					assert.False(t, c.Config.Auth.Enabled)
				}
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
  type: "postgresql"
  maxConns: 50
  url: "postgres://localhost/test"
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
			check: func(t *testing.T, c *Configuration[models.Config]) {
				require.NotNil(t, c)
				require.NotNil(t, c.Config)
				if c.Config.API != nil {
					assert.Equal(t, "192.168.1.1", c.Config.API.Host)
					assert.Equal(t, int32(9090), c.Config.API.Port)
					assert.Equal(t, models.APIConfigEnvironmentProduction, c.Config.API.Environment)
				}
				if c.Config.Database != nil {
					assert.Equal(t, int32(50), c.Config.Database.MaxConns)
				}
				if c.Config.Logging != nil {
					assert.Equal(t, models.LoggingConfigLevelDebug, c.Config.Logging.Level)
					assert.True(t, c.Config.Logging.Pretty)
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

			parser := NewParser[models.Config]()
			config, err := parser.Load()

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.check != nil {
				tt.check(t, config)
			}
		})
	}
}

func TestConfigConstants(t *testing.T) {
	assert.Equal(t, "config", DefaultConfigName)
	assert.Equal(t, "yaml", DefaultConfigType)
	assert.Equal(t, "ARCHES", EnvPrefix)
}

func TestConfigPaths(t *testing.T) {
	expectedPaths := []string{
		".",
		"/etc/archesai/",
		"$HOME/.config/archesai",
	}

	assert.Equal(t, len(expectedPaths), len(ConfigPaths))
	for i, path := range ConfigPaths {
		assert.Equal(t, expectedPaths[i], path)
	}
}

func TestConfigFileNames(t *testing.T) {
	expectedNames := []string{
		"arches",
		"config",
		".archesai",
	}

	assert.Equal(t, len(expectedNames), len(ConfigFileNames))
	for i, name := range ConfigFileNames {
		assert.Equal(t, expectedNames[i], name)
	}
}

func TestGetViperInstance(t *testing.T) {
	parser := NewParser[models.Config]()
	config, err := parser.Load()
	require.NoError(t, err)

	v := config.GetViperInstance()
	assert.NotNil(t, v)
}
