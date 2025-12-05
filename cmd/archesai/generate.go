package main

import (
	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/codegen"
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
code generation will be blocked.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	flags.SetGenerateFlags(generateCmd)
}

func runGenerate(_ *cobra.Command, _ []string) error {
	opts := codegen.Options{
		OutputPath: flags.Generate.OutputPath,
		SpecPath:   flags.Generate.SpecPath,
		Lint:       flags.Generate.Lint,
		Only:       flags.Generate.Only,
	}

	if flags.Generate.TUI {
		return codegen.RunTUI(opts)
	}
	return codegen.Run(opts)
}
