package generate

import (
	"fmt"
	"sync"

	"github.com/archesai/archesai/internal/codegen"
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
		{ID: "bundle", Title: "Bundling OpenAPI specification"},
		{ID: "init", Title: "Initializing code generator"},
		{ID: "generate", Title: "Generating code"},
	}

	if opts.DryRun {
		steps = append(steps, tui.StepDef{ID: "summary", Title: "Preparing summary"})
	}

	var cg *codegen.Codegen

	err := runner.Steps("Code Generation", steps, func(stepID string) (string, error) {
		switch stepID {
		case "bundle":
			if err := BundleSpec(&opts); err != nil {
				return "", err
			}
			return "Bundled successfully", nil

		case "init":
			cg = codegen.NewCodegen(opts.OutputPath)

			if opts.Lint {
				cg = cg.WithLinting()
			}

			if opts.Only != "" {
				cg = cg.WithOnly(opts.Only)
			}

			if opts.DryRun {
				memStorage := storage.NewMemoryStorage()
				cg = cg.WithStorage(memStorage)
			}

			// Set up progress callback
			cg = cg.WithProgress(func(event codegen.ProgressEvent) {
				switch event.Type {
				case codegen.ProgressEventGeneratorDone:
					mu.Lock()
					completedGenerators = append(completedGenerators, event.GeneratorName)
					mu.Unlock()
				case codegen.ProgressEventError:
					generationErr = event.Error
				}
			})

			if err := cg.Initialize(); err != nil {
				return "", fmt.Errorf("failed to initialize: %w", err)
			}
			return "Initialized", nil

		case "generate":
			if err := cg.GenerateAPI(opts.SpecPath); err != nil {
				return "", fmt.Errorf("generation failed: %w", err)
			}

			mu.Lock()
			count := len(completedGenerators)
			mu.Unlock()

			return fmt.Sprintf("Generated %d components", count), nil

		case "summary":
			if opts.DryRun && cg != nil {
				memStorage := cg.GetStorage().(*storage.MemoryStorage)
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

	if opts.DryRun && cg != nil {
		return printDryRunResultsTUI(runner, cg, opts.OutputPath)
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
