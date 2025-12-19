package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/pkg/config"
	configschemas "github.com/archesai/archesai/pkg/config/schemas"
	"github.com/archesai/archesai/pkg/logger"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "archesai",
	Short: "Arches server and utilities",
	Long: `Arches is a comprehensive data processing platform.

This command provides various modes to run the Arches server:
- Development server with hot reload (dev command)
- Worker for background job processing
- Configuration viewer and TUI interface`,
}

// Config is the loaded configuration (nil if not loaded).
var Config *config.Configuration[configschemas.Config]

func init() {
	// Initialize configuration
	cobra.OnInitialize(loadConfig)

	// Global flags
	flags.SetRootFlags(rootCmd)
}

// loadConfig reads in config file and ENV variables if set.
func loadConfig() {
	// Load typed config
	parser := config.NewParser[configschemas.Config]()
	cfg, err := parser.LoadFrom(flags.Root.ConfigFile)
	if err != nil {
		// Config loading is optional for some commands (e.g., version, completion)
		if flags.Root.Verbose {
			fmt.Fprintln(os.Stderr, "Warning: failed to load config:", err)
		}
	} else {
		Config = cfg
		if flags.Root.Verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", Config.ConfigPath)
		}
	}

	// Setup viper for env var support
	viper.SetEnvPrefix("ARCHESAI")
	viper.AutomaticEnv()

	// Configure logger based on flags and set as default
	logLevel := "info"
	if flags.Root.Verbose {
		logLevel = "debug"
	}
	logCfg := logger.Config{
		Level:  logLevel,
		Pretty: flags.Root.Pretty,
	}
	slog.SetDefault(logger.New(logCfg))
}
