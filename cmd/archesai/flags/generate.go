package flags

import (
	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/codegen"
)

// Generate is the global instance of generate flags.
var Generate codegen.Options

// SetGenerateFlags configures flags on the generate command.
func SetGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVar(&Generate.OutputPath, "output", "", "Output directory for generated code (required)")
	cmd.Flags().
		StringVar(&Generate.SpecPath, "spec", "", "Path to OpenAPI specification file (required)")
	cmd.Flags().
		BoolVar(&Generate.Lint, "lint", false, "Enable strict OpenAPI linting (blocks generation on ANY violations)")
	cmd.Flags().
		StringVar(&Generate.Only, "only", "", "Only generate specific components (comma-separated: models,repositories,postgres,sqlite,application,controllers,hcl,sqlc,client,bootstrap)")
	cmd.Flags().BoolVarP(&Generate.TUI, "tui", "t", false, "Enable TUI mode with progress display")

	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagRequired("spec")
}
