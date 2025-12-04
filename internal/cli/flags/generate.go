package flags

import (
	"github.com/spf13/cobra"
)

// GenerateFlags holds the generate command flag values.
type GenerateFlags struct {
	OutputPath string
	SpecPath   string
	OrvalFix   bool
	DryRun     bool
	Lint       bool
	Only       string
	TUI        bool
}

// Generate is the global instance of generate flags.
var Generate GenerateFlags

// SetGenerateFlags configures flags on the generate command.
func SetGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().
		StringVar(&Generate.OutputPath, "output", "", "Output directory for generated code (required)")
	cmd.Flags().
		StringVar(&Generate.SpecPath, "spec", "", "Path to OpenAPI specification file (required)")
	cmd.Flags().
		BoolVar(&Generate.OrvalFix, "orval-fix", false, "Apply fixes for Orval compatibility during bundling")
	cmd.Flags().
		BoolVar(&Generate.Lint, "lint", false, "Enable strict OpenAPI linting (blocks generation on ANY violations)")
	cmd.Flags().
		BoolVar(&Generate.DryRun, "dry-run", false, "Show what would be generated without writing files")
	cmd.Flags().
		StringVar(&Generate.Only, "only", "", "Only generate specific components (comma-separated: models,repositories,postgres,sqlite,application,controllers,hcl,sqlc,client,bootstrap)")
	cmd.Flags().BoolVarP(&Generate.TUI, "tui", "t", false, "Enable TUI mode with progress display")

	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagRequired("spec")
}
