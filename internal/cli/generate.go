package cli

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/internal/tui"
	"github.com/archesai/archesai/pkg/storage"
)

var (
	outputPath  string
	specPath    string
	orvalFix    bool
	dryRun      bool
	lintFlag    bool
	onlyFlag    string
	generateTUI bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from OpenAPI specification",
	Long: `Generate Go code from an OpenAPI specification.

This command generates:
- Models
- Repositories
- Controllers
- Command/Query handlers
- JavaScript/TypeScript client
- Database schema (HCL and SQLC)
- Bootstrap code (app, container, routes, wire)

Use --only to generate specific components (comma-separated):
  go.mod, models, repositories, postgres, sqlite, application, controllers,
  hcl, sqlc, client, app, container, routes, wire, bootstrap (alias for app,container,routes,wire)

By default (no --only flag), all components are generated.

The --lint flag enables strict OpenAPI linting. If ANY violations are found,
code generation will be blocked.

The --dry-run flag will show what files would be generated without actually
writing them to disk.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(_ *cobra.Command, _ []string) error {
		// Output is always required
		if outputPath == "" {
			return fmt.Errorf("--output flag is required")
		}

		// Spec path is required
		if specPath == "" {
			return fmt.Errorf("--spec flag is required")
		}

		if generateTUI {
			return runGenerateTUI()
		}
		return runGenerateStandard()
	},
}

// bundleSpec handles bundling the OpenAPI spec if needed.
func bundleSpec() error {
	shouldBundle := true
	if onlyFlag != "" {
		shouldBundle = false
		for _, component := range strings.Split(onlyFlag, ",") {
			if strings.TrimSpace(strings.ToLower(component)) == "bundle" {
				shouldBundle = true
				break
			}
		}
	}

	if shouldBundle {
		dir := filepath.Dir(specPath)
		bundledPath := filepath.Join(dir, "openapi.bundled.yaml")

		parser := parsers.NewOpenAPIParser()
		if err := parser.Bundle(specPath, bundledPath, orvalFix); err != nil {
			return fmt.Errorf("bundling failed: %w", err)
		}

		slog.Debug("Bundled OpenAPI specification", slog.String("output", bundledPath))
		specPath = bundledPath
	}
	return nil
}

// runGenerateStandard runs generation with standard slog output.
func runGenerateStandard() error {
	if err := bundleSpec(); err != nil {
		return err
	}

	cg := codegen.NewCodegen(outputPath)

	if lintFlag {
		cg = cg.WithLinting()
	}

	if onlyFlag != "" {
		cg = cg.WithOnly(onlyFlag)
	}

	if dryRun {
		memStorage := storage.NewMemoryStorage()
		cg = cg.WithStorage(memStorage)
	}

	if err := cg.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize code generator: %w", err)
	}

	if err := cg.GenerateAPI(specPath); err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}

	if dryRun {
		return printDryRunResults(cg)
	}

	return nil
}

// printDryRunResults displays what would be generated in dry-run mode.
func printDryRunResults(cg *codegen.Codegen) error {
	memStorage := cg.GetStorage().(*storage.MemoryStorage)
	files := memStorage.GetFiles()

	var paths []string
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	slog.Info("Dry-run mode - Files that would be generated")

	dirFiles := make(map[string][]string)
	totalSize := 0
	for _, path := range paths {
		dir := filepath.Dir(path)
		dirFiles[dir] = append(dirFiles[dir], filepath.Base(path))
		totalSize += len(files[path])
	}

	var dirs []string
	for dir := range dirFiles {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		relDir := strings.TrimPrefix(dir, outputPath+"/")
		if relDir == dir {
			relDir = dir
		}

		for _, file := range dirFiles[dir] {
			fileInfo := files[filepath.Join(dir, file)]
			slog.Info("  File would be generated",
				slog.String("dir", relDir),
				slog.String("file", file),
				slog.Int("size", len(fileInfo)),
			)
		}
	}

	slog.Info("Dry-run complete",
		slog.Int("total_files", len(files)),
		slog.Int("total_bytes", totalSize),
	)
	return nil
}

// runGenerateTUI runs generation with TUI progress display.
func runGenerateTUI() error {
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

	if dryRun {
		steps = append(steps, tui.StepDef{ID: "summary", Title: "Preparing summary"})
	}

	var cg *codegen.Codegen

	err := runner.Steps("Code Generation", steps, func(stepID string) (string, error) {
		switch stepID {
		case "bundle":
			if err := bundleSpec(); err != nil {
				return "", err
			}
			return "Bundled successfully", nil

		case "init":
			cg = codegen.NewCodegen(outputPath)

			if lintFlag {
				cg = cg.WithLinting()
			}

			if onlyFlag != "" {
				cg = cg.WithOnly(onlyFlag)
			}

			if dryRun {
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
			if err := cg.GenerateAPI(specPath); err != nil {
				return "", fmt.Errorf("generation failed: %w", err)
			}

			mu.Lock()
			count := len(completedGenerators)
			mu.Unlock()

			return fmt.Sprintf("Generated %d components", count), nil

		case "summary":
			if dryRun && cg != nil {
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

	if dryRun && cg != nil {
		return printDryRunResultsTUI(runner, cg)
	}

	// Show what was generated
	mu.Lock()
	generators := completedGenerators
	mu.Unlock()

	if len(generators) > 0 {
		summary := tui.NewSummary("Generation Complete")
		summary.AddCount("Components", len(generators), "success")
		summary.AddMessage(fmt.Sprintf("Output: %s", outputPath), "info")
		runner.PrintSummary(summary)
	}

	return generationErr
}

// printDryRunResultsTUI displays dry-run results using TUI components.
func printDryRunResultsTUI(runner *tui.Runner, cg *codegen.Codegen) error {
	memStorage := cg.GetStorage().(*storage.MemoryStorage)
	files := memStorage.GetFiles()

	var paths []string
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Group by directory
	dirFiles := make(map[string][]string)
	totalSize := 0
	for _, path := range paths {
		dir := filepath.Dir(path)
		dirFiles[dir] = append(dirFiles[dir], filepath.Base(path))
		totalSize += len(files[path])
	}

	// Create summary
	summary := tui.NewSummary("Dry Run Results")
	summary.AddCount("Files", len(files), "info")
	summary.AddCount("Directories", len(dirFiles), "info")
	summary.AddMessage(fmt.Sprintf("Total size: %d bytes", totalSize), "info")
	runner.PrintSummary(summary)

	runner.PrintNewline()

	// Create table of files
	table := tui.NewTable("Files that would be generated", "Directory", "File", "Size")
	for _, path := range paths {
		dir := filepath.Dir(path)
		relDir := strings.TrimPrefix(dir, outputPath+"/")
		if relDir == dir {
			relDir = dir
		}
		file := filepath.Base(path)
		size := fmt.Sprintf("%d B", len(files[path]))
		table.AddRow(relDir, file, size)
	}
	runner.PrintTable(table)

	return nil
}

func init() {
	// Add generate command to root
	rootCmd.AddCommand(generateCmd)

	// Generate command flags
	generateCmd.Flags().
		StringVar(&outputPath, "output", "", "Output directory for generated code (required)")
	generateCmd.Flags().
		StringVar(&specPath, "spec", "", "Path to OpenAPI specification file (required)")
	generateCmd.Flags().
		BoolVar(&orvalFix, "orval-fix", false, "Apply fixes for Orval compatibility during bundling")
	generateCmd.Flags().
		BoolVar(&lintFlag, "lint", false, "Enable strict OpenAPI linting (blocks generation on ANY violations)")
	generateCmd.Flags().
		BoolVar(&dryRun, "dry-run", false, "Show what would be generated without writing files")
	generateCmd.Flags().
		StringVar(&onlyFlag, "only", "", "Only generate specific components (comma-separated: models,repositories,postgres,sqlite,application,controllers,hcl,sqlc,client,bootstrap)")
	generateCmd.Flags().
		BoolVarP(&generateTUI, "tui", "t", false, "Enable TUI mode with progress display")
	_ = generateCmd.MarkFlagRequired("output")
	_ = generateCmd.MarkFlagRequired("spec")
}
