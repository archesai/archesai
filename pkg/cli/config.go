package cli

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"

	"github.com/archesai/archesai/pkg/config"
)

var (
	outputFormat string
)

// configCmd represents the config command.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Arches configuration",
	Long: `Manage Arches configuration with various subcommands.

This command allows you to:
- Show current configuration
- Validate configuration files
- Initialize default configuration
- Show environment variables`,
}

// configShowCmd shows the current configuration.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration with all applied defaults, environment variables, and config file values.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(cfg.Config)
		case "yaml":
			encoder := yaml.NewEncoder(os.Stdout)
			encoder.SetIndent(2)
			defer func() {
				if err := encoder.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to close encoder: %v\n", err)
				}
			}()
			return encoder.Encode(cfg.Config)
		default:
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}
	},
}

// configValidateCmd validates the configuration.
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long:  `Validate the current configuration for errors and consistency.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration is invalid: %w", err)
		}

		// TODO: Add more comprehensive validation logic
		slog.Info("✓ Configuration is valid")
		slog.Info("✓ Loaded", "source", getConfigSource(cfg))
		return nil
	},
}

// configInitCmd initializes a default configuration file.
var configInitCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize default configuration",
	Long:  `Create a default configuration file with all available options and their default values.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		// Determine output path
		configPath := "config.yaml"
		if len(args) > 0 {
			configPath = args[0]
		}

		// Check if file exists
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("config file %s already exists", configPath)
		}

		// Get default config
		defaultConfig := config.New()

		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Write config file
		file, err := os.Create(configPath)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", err)
			}
		}()

		encoder := yaml.NewEncoder(file)
		encoder.SetIndent(2)
		defer func() {
			if err := encoder.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close encoder: %v\n", err)
			}
		}()

		if err := encoder.Encode(defaultConfig); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		slog.Info("Created default configuration file", "path", configPath)
		return nil
	},
}

// configEnvCmd shows environment variables.
var configEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Long:  `Display all Arches-related environment variables and their values.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		envVars := []string{
			"ARCHESAI_API_HOST",
			"ARCHESAI_API_PORT",
			"ARCHESAI_API_DOCS",
			"ARCHESAI_API_VALIDATE",
			"ARCHESAI_API_ENVIRONMENT",
			"ARCHESAI_API_CORS_ORIGINS",
			"ARCHESAI_AUTH_ENABLED",
			"ARCHESAI_AUTH_LOCAL_ENABLED",
			"ARCHESAI_AUTH_LOCAL_JWT_SECRET",
			"ARCHESAI_AUTH_LOCAL_ACCESS_TOKEN_TTL",
			"ARCHESAI_AUTH_LOCAL_REFRESH_TOKEN_TTL",
			"ARCHESAI_DATABASE_ENABLED",
			"ARCHESAI_DATABASE_URL",
			"ARCHESAI_DATABASE_TYPE",
			"ARCHESAI_DATABASE_MAX_CONNS",
			"ARCHESAI_DATABASE_MIN_CONNS",
			"ARCHESAI_DATABASE_CONN_MAX_LIFETIME",
			"ARCHESAI_DATABASE_CONN_MAX_IDLE_TIME",
			"ARCHESAI_DATABASE_HEALTH_CHECK_PERIOD",
			"ARCHESAI_DATABASE_RUN_MIGRATIONS",
			"ARCHESAI_LOGGING_LEVEL",
			"ARCHESAI_LOGGING_PRETTY",
		}

		slog.Info("Arches Environment Variables:")
		slog.Info("================================")

		found := false
		for _, envVar := range envVars {
			if value := os.Getenv(envVar); value != "" {
				slog.Info("ENV", "key", envVar, "value", value)
				found = true
			}
		}

		if !found {
			slog.Info("No Arches environment variables are currently set.")
		}

		return nil
	},
}

// getConfigSource returns a description of where the config was loaded from.
func getConfigSource(cfg *config.Configuration) string {
	viper := cfg.GetViperInstance()
	if viper.ConfigFileUsed() != "" {
		return fmt.Sprintf("config file: %s", viper.ConfigFileUsed())
	}
	return "defaults and environment variables"
}

func init() {
	// Add config command to root
	rootCmd.AddCommand(configCmd)

	// Add subcommands
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configEnvCmd)

	// Add flags
	configShowCmd.Flags().
		StringVarP(&outputFormat, "output", "o", "yaml", "Output format (yaml, json)")
}
