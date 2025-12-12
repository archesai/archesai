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
	"gopkg.in/yaml.v3"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/internal/dev"
	"github.com/archesai/archesai/pkg/logger"
)

// devConfig is a minimal config struct for dev command.
// This avoids importing the generated models package which creates a circular dependency.
type devConfig struct {
	API      *devAPIConfig      `yaml:"api"`
	Platform *devPlatformConfig `yaml:"platform"`
	Auth     *devAuthConfig     `yaml:"auth"`
}

type devAPIConfig struct {
	URL string `yaml:"url"`
}

type devPlatformConfig struct {
	URL string `yaml:"url"`
}

type devAuthConfig struct {
	Enabled bool `yaml:"enabled"`
}

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

	// Load config to get working directory
	opts := codegen.LoadOptionsFromConfig(flags.Root.ConfigFile)
	rootDir := opts.WorkDir
	if rootDir == "" {
		var err error
		rootDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	// Load config for API and platform URLs
	// Uses a minimal local struct to avoid circular dependency on generated models package
	var cfg devConfig
	configData, err := os.ReadFile(flags.Root.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
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
	if cfg.API != nil && cfg.API.URL != "" {
		frontendEnv = append(
			frontendEnv,
			fmt.Sprintf("VITE_ARCHES_API_HOST=%s", cfg.API.URL),
		)
	}
	if cfg.Platform != nil && cfg.Platform.URL != "" {
		frontendEnv = append(
			frontendEnv,
			fmt.Sprintf("VITE_ARCHES_PLATFORM_URL=%s", cfg.Platform.URL),
		)
	}
	if cfg.Auth != nil {
		frontendEnv = append(
			frontendEnv,
			fmt.Sprintf("VITE_ARCHES_AUTH_ENABLED=%t", cfg.Auth.Enabled),
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
