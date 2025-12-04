package flags

import (
	"github.com/spf13/cobra"
)

// DevFlags holds the dev command flag values.
type DevFlags struct {
	DisableTUI bool
}

// Dev is the global instance of dev flags.
var Dev DevFlags

// SetDevFlags configures flags on the dev command.
func SetDevFlags(cmd *cobra.Command) {
	cmd.Flags().
		BoolVar(&Dev.DisableTUI, "disable-tui", false, "Disable TUI mode for interactive log viewing")
}
