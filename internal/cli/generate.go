package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/internal/parsers"
)

var (
	outputPath string
	bundleFlag bool
	orvalFix   bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from specifications",
	Long: `Generate code from OpenAPI or JSON Schema specifications.

This command provides various code generation capabilities including:
- OpenAPI to Go code generation
- JSON Schema to Go struct generation`,
}

// openapiCmd represents the openapi subcommand
var openapiCmd = &cobra.Command{
	Use:   "openapi [spec-path]",
	Short: "Generate code from an OpenAPI specification",
	Long: `Generate Go code from an OpenAPI specification.

This command generates:
- Models (entities and value objects)
- Repositories
- Controllers
- Command/Query handlers
- Events
- JavaScript/TypeScript client
- Database schema (HCL and SQLC)
- Bootstrap code

The --bundle flag will output a bundled version of the OpenAPI specification
instead of generating code.`,
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(_ *cobra.Command, args []string) error {
		// If bundle flag is set, bundle the OpenAPI spec
		if bundleFlag {
			if len(args) < 1 {
				return fmt.Errorf("input path required when using --bundle")
			}
			if outputPath == "" {
				return fmt.Errorf("--output flag is required when using --bundle")
			}

			inputPath := args[0]
			parser := parsers.NewOpenAPIParser()
			if err := parser.Bundle(inputPath, outputPath, orvalFix); err != nil {
				return fmt.Errorf("bundling failed: %w", err)
			}

			fmt.Printf("✅ Bundled OpenAPI specification written to %s\n", outputPath)
			return nil
		}

		// Regular code generation
		if outputPath == "" {
			return fmt.Errorf("--output flag is required")
		}

		path := "api/openapi.bundled.yaml"
		if len(args) > 0 {
			path = args[0]
		}

		generator := codegen.NewGenerator(outputPath)
		if err := generator.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize code generator: %w", err)
		}

		if _, err := generator.GenerateAPI(path); err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		fmt.Println("✅ Code generation completed successfully")
		return nil
	},
}

// jsonschemaCmd represents the jsonschema subcommand
var jsonschemaCmd = &cobra.Command{
	Use:   "jsonschema [spec-path]",
	Short: "Generate Go structs from a JSON Schema",
	Long: `Generate Go structs from a JSON Schema specification.

This command converts JSON Schema definitions into Go structs with
appropriate JSON tags and validation annotations.`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(_ *cobra.Command, args []string) error {
		schemaPath := args[0]

		if outputPath == "" {
			return fmt.Errorf("--output flag is required")
		}

		// For JSON schema, use the output path's directory as the base
		generator := codegen.NewGenerator(outputPath)
		if err := generator.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize code generator: %w", err)
		}

		if _, err := generator.GenerateJSONSchema(schemaPath, outputPath); err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		fmt.Println("✅ JSON Schema generation completed successfully")
		return nil
	},
}

func init() {
	// Add generate command to root
	rootCmd.AddCommand(generateCmd)

	// Add subcommands to generate
	generateCmd.AddCommand(openapiCmd)
	generateCmd.AddCommand(jsonschemaCmd)

	// OpenAPI command flags
	openapiCmd.Flags().
		StringVar(&outputPath, "output", "", "Output directory for generated code (required)")
	openapiCmd.Flags().
		BoolVar(&bundleFlag, "bundle", false, "Bundle the OpenAPI spec into a single file instead of generating code")
	openapiCmd.Flags().
		BoolVar(&orvalFix, "orval-fix", false, "Apply fixes for Orval compatibility (only used with --bundle)")
	_ = openapiCmd.MarkFlagRequired("output")

	// JSON Schema command flags
	jsonschemaCmd.Flags().StringVar(&outputPath, "output", "", "Output file path")
}
