package flags

import (
	"github.com/spf13/cobra"
)

// SpecLintFlags holds the spec lint command flag values.
type SpecLintFlags struct {
	SpecPath string
}

// SpecShowFlags holds the spec show command flag values.
type SpecShowFlags struct {
	SpecPath string
	JSON     bool
}

// SpecLint is the global instance of spec lint flags.
var SpecLint SpecLintFlags

// SpecShow is the global instance of spec show flags.
var SpecShow SpecShowFlags

// SetSpecLintFlags configures flags on the spec lint command.
func SetSpecLintFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVar(&SpecLint.SpecPath, "spec", "", "Path to OpenAPI specification file (required)")
	_ = cmd.MarkFlagRequired("spec")
}

// SetSpecShowFlags configures flags on the spec show command.
func SetSpecShowFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVar(&SpecShow.SpecPath, "spec", "", "Path to OpenAPI specification file (required)")
	cmd.Flags().BoolVar(&SpecShow.JSON, "json", false, "Output as JSON instead of YAML")
	_ = cmd.MarkFlagRequired("spec")
}
