package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/openapi"
)

// specLintCmd represents the spec lint command
var specLintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint an OpenAPI specification",
	Long: `Lint an OpenAPI specification using strict validation rules.

This command validates your OpenAPI specification against OpenAPI recommended
rules and OWASP security rules. Any violations will be reported with their
location and severity.

Examples:
  archesai spec lint --spec api.yaml
  archesai spec lint --spec ./spec/openapi.yaml`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSpecLint,
}

func init() {
	specCmd.AddCommand(specLintCmd)
	flags.SetSpecLintFlags(specLintCmd)
}

func runSpecLint(_ *cobra.Command, _ []string) error {
	parser := openapi.NewParser()
	if _, err := parser.Parse(flags.SpecLint.SpecPath); err != nil {
		return fmt.Errorf("failed to parse spec: %w", err)
	}

	// Lint (parser.basePath is already set from Parse call)
	return parser.Lint()
}
