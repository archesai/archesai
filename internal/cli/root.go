package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/archesai/archesai/internal/cli/flags"
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

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration
	cobra.OnInitialize(loadConfig)

	// Global flags
	flags.SetRootFlags(rootCmd)
}

// loadConfig reads in config file and ENV variables if set.
func loadConfig() {
	if flags.Root.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flags.Root.ConfigFile)
	} else {
		// Search for config in current directory and home directory
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".archesai")
	}

	// Read in environment variables
	viper.SetEnvPrefix("ARCHESAI")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && flags.Root.Verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

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
