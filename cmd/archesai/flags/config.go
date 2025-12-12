package flags

import (
	"github.com/spf13/cobra"
)

// ConfigShowFlags holds the config show command flag values.
type ConfigShowFlags struct {
	OutputFormat string
}

// ConfigShow is the global instance of config show flags.
var ConfigShow ConfigShowFlags

// SetConfigShowFlags configures flags on the config show command.
func SetConfigShowFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVarP(&ConfigShow.OutputFormat, "output", "o", "yaml", "Output format (yaml, json)")
}
