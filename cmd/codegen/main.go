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
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "api/openapi.bundled.yaml"
		if len(args) > 0 {
			path = args[0]
		}

		if verbose {
			os.Setenv("ARCHESAI_LOGGING_LEVEL", "debug")
		} else {
			os.Setenv("ARCHESAI_LOGGING_LEVEL", "error")
		}

		if err := codegen.Run(path); err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		fmt.Printf("Successfully generated code from %s\n", path)
		return nil
	},
}

// jsonschemaCmd represents the jsonschema subcommand
var jsonschemaCmd = &cobra.Command{
	Use:   "jsonschema [path]",
	Short: "Generate code from JSON Schema",
	Long: `Generate a single Go struct file from a JSON Schema.

The path argument is the JSON Schema file to process.
Requires --output flag for the output file path.
The package name is automatically inferred from the output path.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		schemaPath := args[0]

		if outputPath == "" {
			return fmt.Errorf("--output flag is required")
		}

		if verbose {
			os.Setenv("ARCHESAI_LOGGING_LEVEL", "debug")
		} else {
			os.Setenv("ARCHESAI_LOGGING_LEVEL", "error")
		}

		if err := codegen.RunFromJSONSchema(schemaPath, outputPath); err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		fmt.Printf("Successfully generated %s\n", outputPath)
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