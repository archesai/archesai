package flags

import (
	"github.com/spf13/cobra"
)

// SpecShowFlags holds the spec show command flag values.
type SpecShowFlags struct {
	JSON bool
}

// SpecShow is the global instance of spec show flags.
var SpecShow SpecShowFlags

// SetSpecShowFlags configures flags on the spec show command.
func SetSpecShowFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&SpecShow.JSON, "json", false, "Output as JSON instead of YAML")
}
