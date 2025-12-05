package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/openapi"
)

// specShowCmd represents the spec show command
var specShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show parsed OpenAPI specification",
	Long: `Show a parsed and rendered OpenAPI specification.

This command parses an OpenAPI specification (processing any x-include-*
extensions) and outputs the rendered result. By default, output is in YAML
format. Use --json to output as JSON.

Examples:
  archesai spec show --spec api.yaml
  archesai spec show --spec api.yaml --json`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSpecShow,
}

func init() {
	specCmd.AddCommand(specShowCmd)
	flags.SetSpecShowFlags(specShowCmd)
}

func runSpecShow(_ *cobra.Command, _ []string) error {
	specPath := flags.SpecShow.SpecPath
	if specPath == "" {
		return fmt.Errorf("--spec flag is required")
	}

	// Determine output format
	format := openapi.RenderFormatYAML
	if flags.SpecShow.JSON {
		format = openapi.RenderFormatJSON
	}

	// Parse and render spec (handles x-include processing)
	parser := openapi.NewParser()
	output, err := parser.ParseAndRender(specPath, format)
	if err != nil {
		return fmt.Errorf("failed to render spec: %w", err)
	}

	// Write to stdout
	_, err = os.Stdout.Write(output)
	return err
}
