package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// configEnvCmd shows environment variables.
var configEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Long: `Display all Arches-related environment variables and their values.

Examples:
  archesai config env`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runConfigEnv,
}

func init() {
	configCmd.AddCommand(configEnvCmd)
}

func runConfigEnv(_ *cobra.Command, _ []string) error {
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
}
