package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/spec"
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
  archesai spec show
  archesai spec show --config arches.yaml --json`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSpecShow,
}

func init() {
	specCmd.AddCommand(specShowCmd)
	flags.SetSpecShowFlags(specShowCmd)
}

func runSpecShow(_ *cobra.Command, _ []string) error {
	if Config == nil {
		return fmt.Errorf("config not loaded: ensure arches.yaml exists")
	}

	// Get generation config
	gen := Config.Config.Generation

	specPath := ""
	if gen != nil && gen.Spec != nil {
		specPath = *gen.Spec
	}

	if specPath == "" {
		return fmt.Errorf(
			"spec path is required: set generation.spec in arches.yaml or use --config flag",
		)
	}

	// Resolve spec path relative to working directory if not absolute
	if Config.WorkDir() != "" && !filepath.IsAbs(specPath) {
		specPath = filepath.Join(Config.WorkDir(), specPath)
	}

	var includeNames []string
	if gen != nil {
		includeNames = gen.Includes
	}

	// Build composite filesystem with includes
	baseFS := os.DirFS(filepath.Dir(specPath))
	compositeFS := spec.BuildIncludeFS(baseFS, includeNames)

	// Load OpenAPI document
	doc, err := spec.NewOpenAPIDocumentFromFS(compositeFS, filepath.Base(specPath))
	if err != nil {
		return fmt.Errorf("failed to load OpenAPI document: %w", err)
	}

	// Determine output format
	format := spec.RenderFormatYAML
	if flags.SpecShow.JSON {
		format = spec.RenderFormatJSON
	}

	// Bundle and render the document
	bundler := spec.NewBundler(doc)
	output, err := bundler.Render(format)
	if err != nil {
		return fmt.Errorf("failed to render specification: %w", err)
	}

	_, err = os.Stdout.Write(output)
	return err
}
