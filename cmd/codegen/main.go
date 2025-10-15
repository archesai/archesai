// Package main provides the entry point for the codegen tool.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/internal/shared/logger"
)

var (
	verbose    bool
	pretty     bool
	outputPath string
	orvalFix   bool
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
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true, // Don't show usage on execution errors, only on arg errors
	RunE: func(_ *cobra.Command, args []string) error {
		path := "api/openapi.bundled.yaml"
		if len(args) > 0 {
			path = args[0]
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
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true, // Don't show usage on execution errors, only on arg errors
	RunE: func(_ *cobra.Command, args []string) error {
		schemaPath := args[0]

		if outputPath == "" {
			return fmt.Errorf("--output flag is required")
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

// bundleCmd represents the bundle subcommand
var bundleCmd = &cobra.Command{
	Use:   "bundle [input] [output]",
	Short: "Bundle OpenAPI specification into a single file",
	Long: `Bundle an OpenAPI specification with external references into a single document.
All references will be resolved and the resulting document will be a valid OpenAPI
specification, containing no external references.

Example:
  codegen bundle api/openapi.yaml api/openapi.bundled.yaml
  codegen bundle api/openapi.yaml api/openapi.bundled.yaml --orval-fix`,
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true, // Don't show usage on execution errors, only on arg errors
	RunE: func(_ *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		parser := parsers.NewOpenAPIParser()
		if err := parser.Bundle(inputPath, outputPath, orvalFix); err != nil {
			return fmt.Errorf("bundling failed: %w", err)
		}

		return nil
	},
}

func init() {
	cobra.OnInitialize(initLogger)

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "enable pretty logging output")

	// Add flags to jsonschema command
	jsonschemaCmd.Flags().StringVar(&outputPath, "output", "", "Output file path")

	// Add flags to bundle command
	bundleCmd.Flags().
		BoolVar(&orvalFix, "orval-fix", false, "Resolve pathItems references for Orval compatibility")

	// Add subcommands
	rootCmd.AddCommand(openapiCmd)
	rootCmd.AddCommand(jsonschemaCmd)
	rootCmd.AddCommand(bundleCmd)
}

// initLogger initializes the logger based on flags
func initLogger() {
	logLevel := "info"
	if verbose {
		logLevel = "debug"
	}

	logCfg := logger.Config{
		Level:  logLevel,
		Pretty: pretty,
	}

	slog.SetDefault(logger.New(logCfg))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
