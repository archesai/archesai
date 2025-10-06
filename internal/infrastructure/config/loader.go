package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Configuration wraps the generated Config for easy access.
type Configuration struct {
	*Config
	v *viper.Viper
}

// Load reads configuration from environment variables and returns a Configuration.
func Load() (*Configuration, error) {
	// Setup Viper and populate with defaults
	v := viper.New()
	setupViper(v)

	// Try to read config file with multiple possible names
	var configFound bool
	for _, name := range ConfigFileNames {
		v.SetConfigName(name)
		if err := v.ReadInConfig(); err == nil {
			configFound = true
			break
		} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// If it's not a "file not found" error, it's a real error
			return nil, fmt.Errorf("failed to read config file %s: %w", name, err)
		}
	}

	// It's OK if no config file is found, we'll use defaults and env vars
	_ = configFound

	// Start with default config
	config := New()

	// Unmarshal from viper, which will override defaults with any configured values
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &Configuration{
		Config: config,
		v:      v,
	}, nil
}

// setupViper configures viper for reading config.
func setupViper(v *viper.Viper) {
	// Try multiple config file names in order of priority
	v.SetConfigType(DefaultConfigType)
	for _, path := range ConfigPaths {
		v.AddConfigPath(path)
	}

	// Set the first config name, we'll try others if this fails
	v.SetConfigName(ConfigFileNames[0])

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Set defaults from struct
	setDefaultsFromStruct(v)

	// Explicitly bind environment variables for nested config
	// This is needed because AutomaticEnv doesn't work with nested structs
	_ = v.BindEnv("logging.pretty")
	_ = v.BindEnv("logging.level")
	_ = v.BindEnv("database.url")
	_ = v.BindEnv("database.runmigrations")
	_ = v.BindEnv("api.host")
	_ = v.BindEnv("api.port")
	_ = v.BindEnv("auth.enabled")
}

// setDefaultsFromStruct uses reflection to set viper defaults from a struct with default values.
func setDefaultsFromStruct(v *viper.Viper) {
	defaults := New()
	setStructDefaults(v, defaults, "")
}

// setStructDefaults recursively sets defaults from a struct using reflection.
func setStructDefaults(v *viper.Viper, data any, prefix string) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Get the mapstructure tag
		tag := fieldType.Tag.Get("mapstructure")
		if tag == "" || tag == "-" {
			continue
		}

		// Build the full key path
		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		// Handle different field types
		switch field.Kind() {
		case reflect.Struct:
			// Recursively handle nested structs
			setStructDefaults(v, field.Interface(), key)
		case reflect.Slice, reflect.Array:
			// Set slice/array defaults if not empty
			if field.Len() > 0 {
				v.SetDefault(key, field.Interface())
			}
		default:
			// Set scalar defaults if not zero value
			if !field.IsZero() {
				v.SetDefault(key, field.Interface())
			}
		}
	}
}

// GetViperInstance returns the underlying viper instance for advanced usage.
func (c *Configuration) GetViperInstance() *viper.Viper {
	return c.v
}
