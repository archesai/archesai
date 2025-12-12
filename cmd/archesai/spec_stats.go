package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/spec"
)

// specStatsCmd represents the spec stats command
var specStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics of an OpenAPI specification",
	Long: `Show statistics of an OpenAPI specification.

This command analyzes your OpenAPI specification and displays statistics
about paths, operations, schemas, and other components.

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
	parser := spec.NewParser()
	if _, err := parser.Parse(flags.SpecLint.SpecPath); err != nil {
		return fmt.Errorf("failed to parse spec: %w", err)
	}

	stats, err := parser.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	fmt.Printf("OpenAPI Specification Statistics:\n")
	fmt.Printf("  Title: %s\n", stats.Title)
	fmt.Printf("  Version: %s\n", stats.Version)
	fmt.Printf("  Total Paths: %d\n", stats.TotalPaths)
	fmt.Printf("  Total Operations: %d\n", stats.TotalOperations)
	fmt.Printf("  Total Schemas: %d\n", stats.TotalSchemas)
	fmt.Printf("  Total Parameters: %d\n", stats.TotalParameters)
	fmt.Printf("  Total Responses: %d\n", stats.TotalResponses)
	fmt.Printf("  Total Security Schemes: %d\n", stats.TotalSecuritySchemes)

	return nil
}
