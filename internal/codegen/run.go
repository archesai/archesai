package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/archesai/archesai/internal/codegen/generators"
	"github.com/archesai/archesai/internal/spec"
)

// GenerationResult contains the outcome of a code generation run.
type GenerationResult struct {
	storage *generators.LocalStorage
}

// WasFileWritten returns true if the specified file was written during generation.
func (r *GenerationResult) WasFileWritten(path string) bool {
	if r.storage == nil {
		return false
	}
	return r.storage.WasFileWritten(path)
}

// WrittenFiles returns all files that were written during generation.
func (r *GenerationResult) WrittenFiles() []string {
	if r.storage == nil {
		return nil
	}
	return r.storage.WrittenFiles()
}

// FullPath returns the full path to a file in the output directory.
func (r *GenerationResult) FullPath(path string) string {
	if r.storage == nil {
		return path
	}
	return filepath.Join(r.storage.BaseDir(), path)
}

// PreparedGeneration contains the spec and configured orchestrator
// ready for code generation.
type PreparedGeneration struct {
	Orchestrator *Orchestrator
	Spec         *spec.Spec
	SpecPath     string
}

// prepareGeneration configures the orchestrator with the provided spec.
// This shared function is used by both Run() and RunTUI() to avoid duplication.
func prepareGeneration(opts Options, s *spec.Spec) (*PreparedGeneration, error) {
	// Merge spec options with CLI options (CLI takes precedence)
	mergedOpts := mergeOptions(opts, s)

	// Setup orchestrator
	orch := NewOrchestrator(mergedOpts.OutputPath)
	if len(mergedOpts.Only) > 0 {
		orch = orch.WithOnly(mergedOpts.Only)
	}

	return &PreparedGeneration{
		Orchestrator: orch,
		Spec:         s,
		SpecPath:     opts.SpecPath,
	}, nil
}

// mergeOptions combines CLI options with spec options.
// CLI options take precedence over spec options.
func mergeOptions(opts Options, s *spec.Spec) Options {
	result := opts

	// Use spec values as defaults if CLI didn't specify
	if len(result.Only) == 0 && len(s.CodegenOnly) > 0 {
		result.Only = s.CodegenOnly
	}
	if !result.Lint && s.CodegenLint {
		result.Lint = true
	}

	return result
}

// Run executes code generation with the provided spec.
func Run(opts Options, s *spec.Spec) (*GenerationResult, error) {
	prep, err := prepareGeneration(opts, s)
	if err != nil {
		return nil, err
	}

	if err := prep.Orchestrator.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize code generator: %w", err)
	}

	if err := prep.Orchestrator.Generate(prep.Spec, prep.SpecPath); err != nil {
		return nil, fmt.Errorf("code generation failed: %w", err)
	}

	// Return result with storage for file tracking
	storage, _ := prep.Orchestrator.GetStorage().(*generators.LocalStorage)
	return &GenerationResult{storage: storage}, nil
}
