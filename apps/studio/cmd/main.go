// Package main implements the entry point for the Arches Studio API server.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/archesai/archesai/apps/studio/generated/infrastructure/bootstrap"
	"github.com/archesai/archesai/pkg/config"
	"github.com/archesai/archesai/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Set as default logger
	slog.SetDefault(logger.New(logger.Config{
		Level:  cfg.Logging.Level.String(),
		Pretty: cfg.Logging.Pretty,
	}))

	// Create application container
	appContainer, err := bootstrap.NewApp(cfg.Config)
	if err != nil {
		slog.Error("failed to create app container", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := appContainer.Close(); err != nil {
			slog.Error("failed to close container: %v", "error", err)
		}
	}()

	// Start API server
	if err := appContainer.APIServer.Start(); err != nil {
		slog.Error("failed to start API server", "error", err)
		os.Exit(1)
	}

	slog.Info("server gracefully stopped")
}
