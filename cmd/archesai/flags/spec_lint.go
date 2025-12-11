package flags

import (
	"github.com/spf13/cobra"
)

// SpecLintFlags holds the spec lint command flag values.
type SpecLintFlags struct {
	SpecPath string
}

// SpecLint is the global instance of spec lint flags.
var SpecLint SpecLintFlags

// SetSpecLintFlags configures flags on the spec lint command.
func SetSpecLintFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVar(&SpecLint.SpecPath, "spec", "", "Path to OpenAPI specification file (required)")
	_ = cmd.MarkFlagRequired("spec")
}
