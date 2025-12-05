package codegen

import (
	"fmt"
	"sync"

	"github.com/archesai/archesai/internal/tui"
	"github.com/archesai/archesai/pkg/storage"
)

// RunTUI executes code generation with TUI progress display.
func RunTUI(opts Options) error {
	runner := tui.NewRunner()

	// Track completed generators
	var completedGenerators []string
	var mu sync.Mutex
	var generationErr error

	// Define steps for the TUI
	steps := []tui.StepDef{
		{ID: "prepare", Title: "Preparing code generation"},
		{ID: "init", Title: "Initializing code generator"},
		{ID: "generate", Title: "Generating code"},
	}

	if opts.DryRun {
		steps = append(steps, tui.StepDef{ID: "summary", Title: "Preparing summary"})
	}

	var prep *PreparedGeneration

	err := runner.Steps("Code Generation", steps, func(stepID string) (string, error) {
		switch stepID {
		case "prepare":
			var err error
			prep, err = prepareGeneration(opts)
			if err != nil {
				return "", err
			}
			return "Bundled successfully", nil

		case "init":
			// Set up progress callback
			prep.Orchestrator = prep.Orchestrator.WithProgress(func(event ProgressEvent) {
				switch event.Type {
				case ProgressEventGeneratorDone:
					mu.Lock()
					completedGenerators = append(completedGenerators, event.GeneratorName)
					mu.Unlock()
				case ProgressEventError:
					generationErr = event.Error
				}
			})

			if err := prep.Orchestrator.Initialize(); err != nil {
				return "", fmt.Errorf("failed to initialize: %w", err)
			}
			return "Initialized", nil

		case "generate":
			if err := prep.Orchestrator.Generate(prep.BundledPath); err != nil {
				return "", fmt.Errorf("generation failed: %w", err)
			}

			mu.Lock()
			count := len(completedGenerators)
			mu.Unlock()

			return fmt.Sprintf("Generated %d components", count), nil

		case "summary":
			if opts.DryRun && prep != nil {
				memStorage := prep.Orchestrator.GetStorage().(*storage.MemoryStorage)
				files := memStorage.GetFiles()
				return fmt.Sprintf("%d files would be created", len(files)), nil
			}
			return "", nil
		}

		return "", nil
	})

	if err != nil {
		return err
	}

	// Show final summary
	runner.PrintNewline()

	if opts.DryRun && prep != nil {
		return printDryRunResultsTUI(runner, prep.Orchestrator, opts.OutputPath)
	}

	// Show what was generated
	mu.Lock()
	generators := completedGenerators
	mu.Unlock()

	if len(generators) > 0 {
		summary := tui.NewSummary("Generation Complete")
		summary.AddCount("Components", len(generators), "success")
		summary.AddMessage(fmt.Sprintf("Output: %s", opts.OutputPath), "info")
		runner.PrintSummary(summary)
	}

	return generationErr
}
