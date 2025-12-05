package flags

import (
	"github.com/spf13/cobra"
)

// CapabilitiesFlags holds the capabilities command flag values.
type CapabilitiesFlags struct {
	JSON bool
	All  bool
}

// Capabilities is the global instance of capabilities flags.
var Capabilities CapabilitiesFlags

// SetCapabilitiesFlags configures flags on the capabilities command.
func SetCapabilitiesFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&Capabilities.JSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&Capabilities.All, "all", false, "Show all capabilities including optional")
}
