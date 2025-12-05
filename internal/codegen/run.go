package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/archesai/archesai/internal/openapi"
	"github.com/archesai/archesai/pkg/storage"
)

// PreparedGeneration contains the bundled spec path and configured orchestrator
// ready for code generation.
type PreparedGeneration struct {
	Orchestrator *Orchestrator
	BundledPath  string
}

// prepareGeneration bundles the OpenAPI spec and configures the orchestrator.
// This shared function is used by both Run() and RunTUI() to avoid duplication.
func prepareGeneration(opts Options) (*PreparedGeneration, error) {
	// Parse spec
	parser := openapi.NewParser()
	_, err := parser.Parse(opts.SpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Bundle spec
	dir := filepath.Dir(opts.SpecPath)
	bundledPath := filepath.Join(dir, "openapi.bundled.yaml")
	_, err = parser.Bundle(bundledPath, opts.OrvalFix)
	if err != nil {
		return nil, err
	}

	// Setup orchestrator
	orch := NewOrchestrator(opts.OutputPath)
	if opts.Only != "" {
		orch = orch.WithOnly(opts.Only)
	}
	if opts.DryRun {
		orch = orch.WithStorage(storage.NewMemoryStorage())
	}

	return &PreparedGeneration{
		Orchestrator: orch,
		BundledPath:  bundledPath,
	}, nil
}

// Run executes code generation with standard output.
func Run(opts Options) error {
	prep, err := prepareGeneration(opts)
	if err != nil {
		return err
	}

	if err := prep.Orchestrator.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize code generator: %w", err)
	}

	if err := prep.Orchestrator.Generate(prep.BundledPath); err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}

	if opts.DryRun {
		return printDryRunResults(prep.Orchestrator, opts.OutputPath)
	}

	return nil
}
