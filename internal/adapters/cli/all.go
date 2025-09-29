// Package cli provides command-line interface functionality
package cli

import (
	"context"
	"fmt"
	"log"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/archesai/archesai/internal/application/app"
	"github.com/archesai/archesai/internal/infrastructure/config"
)

// allCmd represents the all command.
var allCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"dev"},
	Short:   "Run all services (API, web, worker) for development",
	Long: `Start all Arches services in a single process for development.

This command runs the API server, web server, and background worker
together, making it convenient for local development.

Note: This mode is NOT recommended for production. In production,
run each service separately for better scaling and isolation.`,
	Example: `  archesai all
  archesai dev  # Using alias`,
	RunE: runAll,
}

func init() {
	rootCmd.AddCommand(allCmd)

	// Reuse flags from individual commands with dev defaults
	allCmd.Flags().String("api-host", "localhost", "API server host")
	allCmd.Flags().Int("api-port", 8080, "API server port")
	allCmd.Flags().Int("web-port", 3000, "Web server port")

	// Bind to viper with prefixes
	if err := viper.BindPFlag("server.host", allCmd.Flags().Lookup("api-host")); err != nil {
		log.Fatalf("Failed to bind api-host flag: %v", err)
	}
	if err := viper.BindPFlag("server.port", allCmd.Flags().Lookup("api-port")); err != nil {
		log.Fatalf("Failed to bind api-port flag: %v", err)
	}
	if err := viper.BindPFlag("web.port", allCmd.Flags().Lookup("web-port")); err != nil {
		log.Fatalf("Failed to bind web-port flag: %v", err)
	}
}

func runAll(_ *cobra.Command, _ []string) error {
	log.Println("🚀 Starting all services in development mode...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override with viper values if set
	if viper.IsSet("server.host") {
		cfg.API.Host = viper.GetString("server.host")
	}
	if viper.IsSet("server.port") {
		cfg.API.Port = float64(viper.GetInt("server.port"))
	}

	// Create application container
	appContainer, err := app.NewApp(cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	defer func() {
		if err := appContainer.Close(); err != nil {
			slog.Error("Failed to close container: %v", "error", err)
		}
	}()

	// Create wait group for services
	var wg sync.WaitGroup

	// Channel to collect errors
	errChan := make(chan error, 3)

	// Start API server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("🚀 API server starting on %s:%d", &cfg.API.Host, int(cfg.API.Port))
		if err := appContainer.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("api server: %w", err)
		}
	}()

	// Give API server time to start
	time.Sleep(1 * time.Second)

	// Start web server (when implemented)
	wg.Add(1)
	go func() {
		defer wg.Done()
		webPort := viper.GetInt("web.port")
		log.Printf("🌐 Web server would start on port %d", webPort)
		// TODO: Implement web server
	}()

	// Start worker (when implemented)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("⚙️  Worker would start processing jobs")
		// TODO: Implement worker
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for either error or interrupt
	select {
	case err := <-errChan:
		log.Printf("Error starting services: %v", err)
		return err
	case <-quit:
		log.Println("\nShutting down all services...")
		// Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := appContainer.Server.Shutdown(shutdownCtx); err != nil {
			log.Printf("API server forced to shutdown: %v", err)
		}
	}

	// Wait for all services to stop
	wg.Wait()
	log.Println("All services stopped")
	return nil
}
