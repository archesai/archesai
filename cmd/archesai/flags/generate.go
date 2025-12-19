package flags

import (
	"github.com/spf13/cobra"
)

// GenerateFlags holds the CLI flags for the generate command.
type GenerateFlags struct {
	OutputPath string
	SpecPath   string
	Lint       bool
	Verbose    bool
}

// Generate is the global instance of generate flags.
var Generate GenerateFlags

// SetGenerateFlags configures flags on the generate command.
func SetGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVar(&Generate.OutputPath, "output", "", "Output directory for generated code (defaults to generation.output in arches.yaml)")
	cmd.Flags().
		StringVar(&Generate.SpecPath, "spec", "", "Path to OpenAPI specification file (defaults to generation.spec in arches.yaml)")
	cmd.Flags().
		BoolVar(&Generate.Lint, "lint", false, "Enable strict OpenAPI linting (blocks generation on ANY violations)")
	cmd.Flags().
		BoolVarP(&Generate.Verbose, "verbose", "v", false, "Show verbose output with step descriptions")
}
