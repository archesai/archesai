package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/cli/flags"
	"github.com/archesai/archesai/internal/generate"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from OpenAPI specification",
	Long: `Generate Go code from an OpenAPI specification.

This command generates:
- Models
- Repositories
- HTTP Handlers
- Application Handlers
- React frontend and TypeScript client
- Database schema (HCL and SQLC)
- Bootstrap code (app, container, routes, wire)

Use --only to generate specific components (comma-separated):
  go.mod, models, repositories, postgres, sqlite, application, controllers,
  hcl, sqlc, client, app, container, routes, wire, bootstrap (alias for app,container,routes,wire)

By default (no --only flag), all components are generated.

The --lint flag enables strict OpenAPI linting. If ANY violations are found,
code generation will be blocked.

The --dry-run flag will show what files would be generated without actually
writing them to disk.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	flags.SetGenerateFlags(generateCmd)
}

func runGenerate(_ *cobra.Command, _ []string) error {
	// Output is always required
	if flags.Generate.OutputPath == "" {
		return fmt.Errorf("--output flag is required")
	}

	// Spec path is required
	if flags.Generate.SpecPath == "" {
		return fmt.Errorf("--spec flag is required")
	}

	opts := generate.Options{
		OutputPath: flags.Generate.OutputPath,
		SpecPath:   flags.Generate.SpecPath,
		OrvalFix:   flags.Generate.OrvalFix,
		DryRun:     flags.Generate.DryRun,
		Lint:       flags.Generate.Lint,
		Only:       flags.Generate.Only,
	}

	if flags.Generate.TUI {
		return generate.RunTUI(opts)
	}
	return generate.Run(opts)
}
