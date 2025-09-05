package cli

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the API server",
	Long: `Start the ArchesAI API server which provides REST endpoints
for all platform functionality.

The API server handles authentication, organizations, workflows,
content management, and more.`,
	Example: `  archesai api --port=8080
  archesai api --host=0.0.0.0 --port=3000
  archesai api --config=production.yaml`,
	RunE: runAPI,
}

var (
	apiHost string
	apiPort int
)

func init() {
	rootCmd.AddCommand(apiCmd)

	// Local flags for API server
	apiCmd.Flags().StringVar(&apiHost, "host", "0.0.0.0", "Host to bind the server to")
	apiCmd.Flags().IntVar(&apiPort, "port", 8080, "Port to bind the server to")

	// Bind to viper
	if err := viper.BindPFlag("server.host", apiCmd.Flags().Lookup("host")); err != nil {
		log.Fatalf("Failed to bind host flag: %v", err)
	}
	if err := viper.BindPFlag("server.port", apiCmd.Flags().Lookup("port")); err != nil {
		log.Fatalf("Failed to bind port flag: %v", err)
	}
}

func runAPI(cmd *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override with command line flags
	if cmd.Flags().Changed("host") {
		cfg.Api.Host = viper.GetString("server.host")
	}
	if cmd.Flags().Changed("port") {
		cfg.Api.Port = float32(viper.GetInt("server.port"))
	}

	// Create application container
	appContainer, err := container.NewContainer(cfg)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	defer func() {
		if err := appContainer.Close(); err != nil {
			slog.Error("Failed to close container: %v", "error", err)
		}
	}()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run server in goroutine
	go func() {
		log.Printf("🚀 API server starting on %s:%d", cfg.Api.Host, int(cfg.Api.Port))
		if err := appContainer.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := appContainer.Server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
