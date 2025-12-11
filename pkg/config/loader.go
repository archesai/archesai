package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// Configuration wraps the Config for easy access.
type Configuration[C any] struct {
	Config     *C
	ConfigPath string // Path to the loaded config file
}

// WorkDir returns the directory containing the config file.
// If no config file was loaded, returns the current working directory.
func (c *Configuration[C]) WorkDir() string {
	if c.ConfigPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "."
		}
		return cwd
	}
	return filepath.Dir(c.ConfigPath)
}

// Load reads configuration from config files and returns a Configuration.
func (p *Parser[C]) Load() (*Configuration[C], error) {
	return p.LoadFrom("")
}

// LoadFrom reads configuration from the specified file path (or searches for config files if empty).
func (p *Parser[C]) LoadFrom(configFile string) (*Configuration[C], error) {
	v := viper.New()
	v.SetConfigType(DefaultConfigType)

	for _, path := range ConfigPaths {
		v.AddConfigPath(path)
	}

	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
		}
	} else {
		// Try to read config file with multiple possible names
		var configFound bool
		for _, name := range ConfigFileNames {
			v.SetConfigName(name)
			if err := v.ReadInConfig(); err == nil {
				configFound = true
				break
			} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config file %s: %w", name, err)
			}
		}
		if !configFound {
			return nil, fmt.Errorf("no config file found (searched for %v in %v)", ConfigFileNames, ConfigPaths)
		}
	}

	// Unmarshal using yaml tags since generated models don't have mapstructure tags
	config := new(C)
	if err := v.Unmarshal(config, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &Configuration[C]{
		Config:     config,
		ConfigPath: v.ConfigFileUsed(),
	}, nil
}

// EnvLoader loads config with environment variable support.
type EnvLoader[C any] struct{}

// Load reads configuration with environment variable overrides.
func (l *EnvLoader[C]) Load() (*Configuration[C], error) {
	return l.LoadFrom("")
}

// LoadFrom reads configuration from the specified file with env var overrides.
func (l *EnvLoader[C]) LoadFrom(configFile string) (*Configuration[C], error) {
	v := viper.New()
	v.SetConfigType(DefaultConfigType)

	for _, path := range ConfigPaths {
		v.AddConfigPath(path)
	}

	// Enable environment variables
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
		}
	} else {
		// Try to read config file with multiple possible names
		for _, name := range ConfigFileNames {
			v.SetConfigName(name)
			if err := v.ReadInConfig(); err == nil {
				break
			} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config file %s: %w", name, err)
			}
		}
	}

	// Unmarshal using yaml tags
	config := new(C)
	if err := v.Unmarshal(config, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &Configuration[C]{
		Config:     config,
		ConfigPath: v.ConfigFileUsed(),
	}, nil
}
