package main

import (
	"log"

	"github.com/archesai/archesai/internal/app"
	"github.com/archesai/archesai/internal/infrastructure/config"
)

func main() {
	// Load configuration using Viper
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize application container (includes server)
	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	// Start the server (container owns it now)
	if err := container.Server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
