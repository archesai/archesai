package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/archesai/archesai/pkg/logger"
)

var (
	cfgFile string
	verbose bool
	pretty  bool
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
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is .archesai.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "enable pretty logging output")

	// Bind flags to viper
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		slog.Error("Failed to bind verbose flag", "err", err)
	}
	if err := viper.BindPFlag("pretty", rootCmd.PersistentFlags().Lookup("pretty")); err != nil {
		slog.Error("Failed to bind pretty flag", "err", err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Configure logger based on flags
	logLevel := "info"
	if verbose {
		logLevel = "debug"
	}

	logCfg := logger.Config{
		Level:  logLevel,
		Pretty: pretty,
	}

	// Set as default logger
	slog.SetDefault(logger.New(logCfg))
}
