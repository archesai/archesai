package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/dev"
	"github.com/archesai/archesai/pkg/logger"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run development server with hot reload",
	Long: `Start the Arches development environment with both the API server and platform UI.

The dev command runs:
- API server with hot reload
- Platform UI with Vite dev server

Both services run concurrently with colored log output.`,
	Example:       `  archesai dev`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runDev,
}

func init() {
	rootCmd.AddCommand(devCmd)
}

func runDev(_ *cobra.Command, _ []string) error {
	// Force pretty logging for dev mode
	logLevel := "info"
	if flags.Root.Verbose {
		logLevel = "debug"
	}
	devLogger := logger.NewPretty(logger.Config{
		Level:  logLevel,
		Pretty: true,
	})
	slog.SetDefault(devLogger)

	// Create process manager
	manager := dev.NewManager(devLogger)

	if Config == nil {
		return fmt.Errorf("config not loaded: ensure arches.yaml exists")
	}

	rootDir := Config.WorkDir()
	if rootDir == "" {
		var err error
		rootDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	// Add API process with custom hot reload
	if err := manager.AddProcess(dev.ProcessConfig{
		Name:       "api",
		Command:    "./bin/app",
		Dir:        rootDir,
		HotReload:  true,
		BuildCmd:   "go",
		BuildArgs:  []string{"build", "-o", "./bin/app", "./main.gen.go"},
		WatchPaths: []string{"."},
		WatchExts:  []string{".go", ".mod", ".sum"},
	}); err != nil {
		return fmt.Errorf("failed to add API process: %w", err)
	}

	// Build environment variables for frontend
	var frontendEnv []string
	if Config.Config.API != nil && Config.Config.API.URL != nil {
		frontendEnv = append(
			frontendEnv,
			fmt.Sprintf("VITE_ARCHES_API_HOST=%s", *Config.Config.API.URL),
		)
	}
	if Config.Config.Platform != nil && Config.Config.Platform.URL != nil {
		frontendEnv = append(
			frontendEnv,
			fmt.Sprintf("VITE_ARCHES_PLATFORM_URL=%s", *Config.Config.Platform.URL),
		)
	}
	if Config.Config.Auth != nil {
		frontendEnv = append(
			frontendEnv,
			fmt.Sprintf("VITE_ARCHES_AUTH_ENABLED=%t", Config.Config.Auth.Enabled),
		)
	}

	// Add frontend process (runs from web/ subdirectory)
	frontendDir := filepath.Join(rootDir, "web")
	if err := manager.AddProcess(dev.ProcessConfig{
		Name:    "frontend",
		Command: "pnpm",
		Args:    []string{"run", "dev"},
		Dir:     frontendDir,
		Env:     frontendEnv,
	}); err != nil {
		return fmt.Errorf("failed to add Platform process: %w", err)
	}

	// Start all processes
	if err := manager.StartAll(); err != nil {
		return fmt.Errorf("failed to start processes: %w", err)
	}

	// Handle interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-quit
	devLogger.Info("Shutting down development server...")

	// Shutdown manager
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- manager.Shutdown()
	}()

	select {
	case err := <-shutdownDone:
		if err != nil {
			return fmt.Errorf("failed to shutdown cleanly: %w", err)
		}
	case <-ctx.Done():
		devLogger.Warn("Shutdown timeout exceeded, forcing exit")
	}
	devLogger.Info("Development server stopped")

	return nil
}
