// Package main provides the entry point for the codegen tool.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/codegen"
)

var (
	verbose    bool
	outputPath string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "codegen",
	Short: "Arches code generation tool",
	Long: `Generate Go code from OpenAPI specifications or JSON Schema files.

This tool supports two modes:
1. OpenAPI mode: Generates multiple files based on x-codegen extensions
2. JSON Schema mode: Generates a single struct file from a JSON Schema`,
}

// openapiCmd represents the openapi subcommand
var openapiCmd = &cobra.Command{
	Use:   "openapi [path]",
	Short: "Generate code from OpenAPI specification",
	Long: `Generate multiple Go files based on x-codegen extensions in the OpenAPI spec.

If no path is provided, defaults to api/openapi.bundled.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		path := "api/openapi.bundled.yaml"
		if len(args) > 0 {
			path = args[0]
		}

		if verbose {
			if err := os.Setenv("ARCHESAI_LOGGING_LEVEL", "debug"); err != nil {
				return fmt.Errorf("failed to set logging level: %w", err)
			}
		} else {
			if err := os.Setenv("ARCHESAI_LOGGING_LEVEL", "error"); err != nil {
				return fmt.Errorf("failed to set logging level: %w", err)
			}
		}

		generator := codegen.NewGenerator()
		if err := generator.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize code generator: %w", err)
		}

		if _, err := generator.GenerateAPI(path); err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		return nil
	},
}

// jsonschemaCmd represents the jsonschema subcommand
var jsonschemaCmd = &cobra.Command{
	Use:   "jsonschema [path]",
	Short: "Generate code from JSON Schema",
	Long: `Generate a single Go struct file from a JSON Schema.

The path argument is the JSON Schema file to process.
Requires --output flag for the output directory.
The package name is automatically inferred from the output directory.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		schemaPath := args[0]

		if outputPath == "" {
			return fmt.Errorf("--output flag is required")
		}

		if verbose {
			if err := os.Setenv("ARCHESAI_LOGGING_LEVEL", "debug"); err != nil {
				return fmt.Errorf("failed to set logging level: %w", err)
			}
		} else {
			if err := os.Setenv("ARCHESAI_LOGGING_LEVEL", "error"); err != nil {
				return fmt.Errorf("failed to set logging level: %w", err)
			}
		}

		generator := codegen.NewGenerator()
		if err := generator.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize code generator: %w", err)
		}

		if _, err := generator.GenerateJSONSchema(schemaPath, outputPath); err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		return nil
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add flags to jsonschema command
	jsonschemaCmd.Flags().StringVar(&outputPath, "output", "", "Output file path")

	// Add subcommands
	rootCmd.AddCommand(openapiCmd)
	rootCmd.AddCommand(jsonschemaCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
