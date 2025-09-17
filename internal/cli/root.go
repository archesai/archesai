package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "archesai",
	Short: "Arches server and utilities",
	Long: `Arches is a comprehensive data processing platform.

This command provides various modes to run the Arches server:
- API server for REST endpoints
- Web server for the UI (future)
- Worker for background job processing
- All mode for development`,
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

	// Bind flags to viper
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		log.Fatalf("Failed to bind verbose flag: %v", err)
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
}
