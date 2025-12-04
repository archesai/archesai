package generate

import (
	"fmt"

	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/pkg/storage"
)

// Run executes code generation with standard output.
func Run(opts Options) error {
	if err := BundleSpec(&opts); err != nil {
		return err
	}

	cg := codegen.NewCodegen(opts.OutputPath)

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

	if err := cg.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize code generator: %w", err)
	}

	if err := cg.GenerateAPI(opts.SpecPath); err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}

	if opts.DryRun {
		return printDryRunResults(cg, opts.OutputPath)
	}

	return nil
}
