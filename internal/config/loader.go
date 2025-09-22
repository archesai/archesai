package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Config wraps the generated ArchesConfig for easy access.
type Config struct {
	*ArchesConfig
	v *viper.Viper
}

// Load reads configuration from environment variables and returns a Config.
func Load() (*Config, error) {
	// Setup Viper and populate with defaults
	v := viper.New()
	setupViper(v)

	// Read config file if exists
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Start with default config
	config := GetDefaultConfig()

	// Unmarshal from viper, which will override defaults with any configured values
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &Config{
		ArchesConfig: config,
		v:            v,
	}, nil
}

// setupViper configures viper for reading config.
func setupViper(v *viper.Viper) {
	v.SetConfigName(DefaultConfigName)
	v.SetConfigType(DefaultConfigType)
	for _, path := range ConfigPaths {
		v.AddConfigPath(path)
	}

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
	defaults := GetDefaultConfig()
	setStructDefaults(v, defaults, "")
}

// setStructDefaults recursively sets defaults from a struct using reflection.
func setStructDefaults(v *viper.Viper, data interface{}, prefix string) {
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
func (c *Config) GetViperInstance() *viper.Viper {
	return c.v
}
