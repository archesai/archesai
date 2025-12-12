package flags

import (
	"github.com/spf13/cobra"
)

// SpecExtractFlags holds the spec extract command flag values.
type SpecExtractFlags struct {
	SpecPath string
	Output   string
	DryRun   bool
	Force    bool
	Verbose  bool
}

// SpecExtract is the global instance of spec extract flags.
var SpecExtract SpecExtractFlags

// SetSpecExtractFlags configures flags on the spec extract command.
func SetSpecExtractFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVar(&SpecExtract.SpecPath, "spec", "", "Path to OpenAPI specification file (required)")
	cmd.Flags().
		StringVar(&SpecExtract.Output, "output", "", "Output directory (default: same as spec directory)")
	cmd.Flags().
		BoolVar(&SpecExtract.DryRun, "dry-run", false, "Print what would be extracted without making changes")
	cmd.Flags().
		BoolVar(&SpecExtract.Force, "force", false, "Overwrite existing files")
	cmd.Flags().
		BoolVar(&SpecExtract.Verbose, "verbose", false, "Print detailed output")
	_ = cmd.MarkFlagRequired("spec")
}
