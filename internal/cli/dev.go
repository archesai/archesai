package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/dev"
	"github.com/archesai/archesai/internal/tui"
	"github.com/archesai/archesai/pkg/logger"
)

var (
	devTUI bool
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
	devCmd.Flags().BoolVar(&devTUI, "tui", false, "Enable TUI mode for interactive log viewing")
}

func runDev(_ *cobra.Command, _ []string) error {
	// Load configuration
	// cfg, err := config.Load()
	// if err != nil {
	// 	return fmt.Errorf("failed to load configuration: %w", err)
	// }

	// Create logger configuration
	logCfg := logger.Config{
		Level:  "info",
		Pretty: !devTUI,
	}
	if devTUI {
		logCfg.Level = "silent"
	}
	baseLogger := logger.New(logCfg)

	// Create process manager
	manager := dev.NewManager(baseLogger)

	// Get project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Configure API process with custom hot reload
	apiConfig := dev.ProcessConfig{
		Name:       "api",
		Command:    "./bin/studio",
		Dir:        rootDir,
		Env:        []string{},
		HotReload:  true,
		BuildCmd:   "go",
		BuildArgs:  []string{"build", "-o", "./bin/studio", "./apps/studio/main.gen.go"},
		WatchPaths: []string{"."},
		WatchExts:  []string{".go", ".mod", ".sum"},
	}

	if err := manager.AddProcess(apiConfig); err != nil {
		return fmt.Errorf("failed to add API process: %w", err)
	}

	// Configure Platform process
	platformConfig := dev.ProcessConfig{
		Name:    "platform",
		Command: "pnpm",
		Args:    []string{"-F", "@archesai/studio", "dev"},
		Dir:     rootDir,
		// Env: []string{
		// 	fmt.Sprintf("VITE_API_URL=http://%s:%d", cfg.API.Host, cfg.API.Port),
		// },
	}

	if err := manager.AddProcess(platformConfig); err != nil {
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
	if devTUI {
		// Run TUI - this blocks until user quits
		if err := tui.RunDevTUI(manager); err != nil {
			baseLogger.Error("TUI error", "error", err)
		}
		// TUI exited, now shutdown
		fmt.Println("Shutting down development server...")
	} else {
		// Just wait for interrupt in non-TUI mode
		<-quit
		baseLogger.Info("Shutting down development server...")
	}

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
		baseLogger.Warn("Shutdown timeout exceeded, forcing exit")
		fmt.Println("Shutdown timeout exceeded, forcing exit")
	}

	baseLogger.Info("Development server stopped")
	fmt.Println("Development server stopped")
	return nil
}
