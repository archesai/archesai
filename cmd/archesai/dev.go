package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/dev"
	"github.com/archesai/archesai/internal/tui"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run development server with hot reload",
	Long: `Start the Arches development environment with both the API server and platform UI.

The dev command runs:
- API server with hot reload (using air)
- Platform UI with Vite dev server

Both services run concurrently and logs are combined for easy monitoring.`,
	Example: `  archesai dev
  archesai dev --tui`,
	RunE: runDev,
}

func init() {
	rootCmd.AddCommand(devCmd)
	flags.SetDevFlags(devCmd)
}

func runDev(_ *cobra.Command, _ []string) error {
	// Get base logger
	logger := slog.Default()

	// Create process manager
	manager := dev.NewManager(logger)

	// Get project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Add API process with custom hot reload
	if err := manager.AddProcess(dev.ProcessConfig{
		Name:       "api",
		Command:    "./bin/studio",
		Dir:        rootDir,
		HotReload:  true,
		BuildCmd:   "go",
		BuildArgs:  []string{"build", "-o", "./bin/studio", "./apps/studio/main.gen.go"},
		WatchPaths: []string{"."},
		WatchExts:  []string{".go", ".mod", ".sum"},
	}); err != nil {
		return fmt.Errorf("failed to add API process: %w", err)
	}

	// Add frontend process
	if err := manager.AddProcess(dev.ProcessConfig{
		Name:    "frontend",
		Command: "pnpm",
		Args:    []string{"run", "dev"},
		Dir:     rootDir,
		// Env: []string{
		// 	fmt.Sprintf("VITE_API_URL=http://%s:%d", cfg.API.Host, cfg.API.Port),
		// },
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

	// Run TUI or wait for interrupt
	if flags.Dev.DisableTUI {
		// Just wait for interrupt in non-TUI mode
		<-quit
	} else {
		// Run TUI - this blocks until user quits
		if err := tui.RunDevTUI(manager); err != nil {
			logger.Error("TUI error", "error", err)
		}
	}
	logger.Info("Shutting down development server...")

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
		logger.Warn("Shutdown timeout exceeded, forcing exit")
	}
	logger.Info("Development server stopped")

	return nil
}
