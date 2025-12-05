package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/openapi"
)

// specStatsCmd represents the spec stats command
var specStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics of an OpenAPI specification",
	Long: `Show statistics of an OpenAPI specification.

This command validates your OpenAPI specification against OpenAPI recommended
rules and OWASP security rules. Any violations will be reported with their
location and severity.

Examples:
  archesai spec stats --spec api.yaml
  archesai spec stats --spec ./spec/openapi.yaml`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSpecStats,
}

func init() {
	specCmd.AddCommand(specStatsCmd)
	flags.SetSpecLintFlags(specStatsCmd)
}

func runSpecStats(_ *cobra.Command, _ []string) error {
	parser := openapi.NewParser()
	if _, err := parser.Parse(flags.SpecLint.SpecPath); err != nil {
		return fmt.Errorf("failed to parse spec: %w", err)
	}

	// Lint (parser.basePath is already set from Parse call)
	return parser.GetStats()
}
