package main

import (
	"fmt"
	"os"

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
	specPath := flags.SpecLint.SpecPath
	if specPath == "" {
		return fmt.Errorf("--spec flag is required")
	}

	// Create parser and parse spec (handles x-include processing)
	parser := openapi.NewParser()
	if _, err := parser.Parse(specPath); err != nil {
		return fmt.Errorf("failed to parse spec: %w", err)
	}

	// Read the spec bytes for linting
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("failed to read spec file: %w", err)
	}

	// Lint (parser.basePath is already set from Parse call)
	return parser.Lint(specBytes)
}
