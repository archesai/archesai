// Package main implements the entry point for the Arches Studio API server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/archesai/archesai/apps/studio/generated/infrastructure/bootstrap"
	"github.com/archesai/archesai/pkg/config"
	"github.com/archesai/archesai/pkg/logger"
)

// apiCmd represents the api command.
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the API server",
	Long: `Start the Arches API server which provides REST endpoints
for all platform functionality.

The API server handles authentication, organizations, workflows,
content management, and more.`,
	Example: `  archesai api --port=8080
  archesai api --host=0.0.0.0 --port=3000
  archesai api --config=production.yaml`,
	RunE: runAPI,
}

func main() {
	_ = apiCmd.Execute()
}

func runAPI(cmd *cobra.Command, _ []string) error {

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logCfg := logger.Config{
		Level:  cfg.Logging.Level.String(),
		Pretty: cfg.Logging.Pretty,
	}

	// Set as default logger
	slog.SetDefault(logger.New(logCfg))

	// Override with command line flags
	if cmd.Flags().Changed("host") {
		cfg.API.Host = viper.GetString("server.host")
	}
	if cmd.Flags().Changed("port") {
		cfg.API.Port = int32(viper.GetInt("server.port"))
	}

	// Create application container
	appContainer, err := bootstrap.NewApp(cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	defer func() {
		if err := appContainer.Close(); err != nil {
			slog.Error("failed to close container: %v", "error", err)
		}
	}()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run server in goroutine
	go func() {
		slog.Info("api server starting", "host", cfg.API.Host, "port", int(cfg.API.Port))
		if err := appContainer.APIServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			slog.Error("failed to start server", "err", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	slog.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := appContainer.APIServer.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
		return err
	}

	slog.Info("server gracefully stopped")
	return nil
}
