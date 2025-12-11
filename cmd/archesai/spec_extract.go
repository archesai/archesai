package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/spec"
)

// specExtractCmd represents the spec extract command
var specExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract inline OpenAPI definitions into separate files",
	Long: `Extract inline OpenAPI definitions into separate component files.

This command analyzes your OpenAPI specification, finds inline definitions
(schemas, responses, parameters, headers, request bodies), and extracts them
into separate files in the standard directory structure.

Examples:
  archesai spec extract --spec api.yaml
  archesai spec extract --spec api.yaml --dry-run
  archesai spec extract --spec api.yaml --output ./extracted
  archesai spec extract --spec api.yaml --force --verbose`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSpecExtract,
}

func init() {
	specCmd.AddCommand(specExtractCmd)
	flags.SetSpecExtractFlags(specExtractCmd)
}

func runSpecExtract(_ *cobra.Command, _ []string) error {
	extractor := spec.NewExtractor(spec.ExtractorConfig{
		SpecPath: flags.SpecExtract.SpecPath,
		Output:   flags.SpecExtract.Output,
		DryRun:   flags.SpecExtract.DryRun,
		Force:    flags.SpecExtract.Force,
		Verbose:  flags.SpecExtract.Verbose,
	})

	result, err := extractor.Extract()
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	// Print summary
	fmt.Printf("Extraction complete:\n")
	fmt.Printf("  Schemas extracted: %d\n", result.SchemasExtracted)
	fmt.Printf("  Request bodies extracted: %d\n", result.RequestBodiesExtracted)
	fmt.Printf("  Responses extracted: %d\n", result.ResponsesExtracted)
	fmt.Printf("  Parameters extracted: %d\n", result.ParametersExtracted)
	fmt.Printf("  Headers extracted: %d\n", result.HeadersExtracted)
	fmt.Printf("  Files modified: %d\n", result.FilesModified)
	if result.Skipped > 0 {
		fmt.Printf("  Skipped (conflicts): %d\n", result.Skipped)
	}

	return nil
}
