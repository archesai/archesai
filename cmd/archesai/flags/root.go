// Package flags provides flag definitions for CLI commands.
package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootFlags holds the root command flag values.
type RootFlags struct {
	ConfigFile string
	Verbose    bool
	Pretty     bool
}

// Root is the global instance of root flags.
var Root RootFlags

// SetRootFlags configures persistent flags on the root command.
func SetRootFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().
		StringVar(&Root.ConfigFile, "config", "", "config file (default is .archesai.yaml)")
	cmd.PersistentFlags().BoolVarP(&Root.Verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolVar(&Root.Pretty, "pretty", false, "enable pretty logging output")

	// Bind flags to viper for configuration file support
	_ = viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("pretty", cmd.PersistentFlags().Lookup("pretty"))
}
