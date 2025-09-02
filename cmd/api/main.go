// Package main provides the entry point for the ArchesAI API server.
package main

import (
	"log/slog"
	"os"

	"github.com/archesai/archesai/internal/app"
	"github.com/archesai/archesai/internal/infrastructure/config"
)

func main() {
	// Load configuration using Viper
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration: %v", "error", err)
		os.Exit(1)
	}

	// Initialize application container (includes server)
	container, err := app.NewContainer(cfg)
	if err != nil {
		slog.Error("Failed to initialize application: %v", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := container.Close(); err != nil {
			slog.Error("Failed to close container: %v", "error", err)
		}
	}()

	// Start the server (container owns it now)
	if err := container.Server.Start(); err != nil {
		slog.Error("Failed to start server: %v", "error", err)
	}
}
