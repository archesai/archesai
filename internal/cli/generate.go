package cli

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/storage"
)

var (
	outputPath string
	specPath   string
	bundleFlag bool
	orvalFix   bool
	dryRun     bool
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
- Events
- JavaScript/TypeScript client
- Database schema (HCL and SQLC)
- Bootstrap code

The --bundle flag will output a bundled version of the OpenAPI specification
instead of generating code.

The --dry-run flag will show what files would be generated without actually
writing them to disk.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(_ *cobra.Command, _ []string) error {
		// If bundle flag is set, bundle the OpenAPI spec
		if bundleFlag {
			if specPath == "" {
				return fmt.Errorf("--spec flag is required when using --bundle")
			}
			if outputPath == "" {
				return fmt.Errorf("--output flag is required when using --bundle")
			}

			parser := parsers.NewOpenAPIParser()
			if err := parser.Bundle(specPath, outputPath, orvalFix); err != nil {
				return fmt.Errorf("bundling failed: %w", err)
			}

			slog.Info("Bundled OpenAPI specification", slog.String("output", outputPath))
			return nil
		}

		// Regular code generation
		if outputPath == "" {
			return fmt.Errorf("--output flag is required")
		}

		// Use default spec path if not provided
		if specPath == "" {
			specPath = "api/openapi.bundled.yaml"
		}

		generator := codegen.NewGenerator(outputPath)

		// Use memory storage for dry-run
		if dryRun {
			memStorage := storage.NewMemoryStorage()
			generator = generator.WithStorage(memStorage)
		}

		if err := generator.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize code generator: %w", err)
		}

		if err := generator.GenerateAPI(specPath); err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		// If dry-run, print what would be generated
		if dryRun {
			memStorage := generator.GetStorage().(*storage.MemoryStorage)
			files := memStorage.GetFiles()

			// Sort file paths for consistent output
			var paths []string
			for path := range files {
				paths = append(paths, path)
			}
			sort.Strings(paths)

			slog.Info("🔍 Dry-run mode - Files that would be generated")

			// Group files by directory
			dirFiles := make(map[string][]string)
			totalSize := 0
			for _, path := range paths {
				dir := filepath.Dir(path)
				dirFiles[dir] = append(dirFiles[dir], filepath.Base(path))
				totalSize += len(files[path])
			}

			// Sort directories
			var dirs []string
			for dir := range dirFiles {
				dirs = append(dirs, dir)
			}
			sort.Strings(dirs)

			// Log each directory and its files
			for _, dir := range dirs {
				relDir := strings.TrimPrefix(dir, outputPath+"/")
				if relDir == dir {
					relDir = dir
				}

				// Log files in this directory
				for _, file := range dirFiles[dir] {
					fileInfo := files[filepath.Join(dir, file)]
					slog.Info("  File would be generated",
						slog.String("dir", relDir),
						slog.String("file", file),
						slog.Int("size", len(fileInfo)),
					)
				}
			}

			slog.Info("✨ Dry-run complete",
				slog.Int("total_files", len(files)),
				slog.Int("total_bytes", totalSize),
			)
			return nil
		}

		slog.Info("Code generation completed successfully")
		return nil
	},
}

func init() {
	// Add generate command to root
	rootCmd.AddCommand(generateCmd)

	// Generate command flags
	generateCmd.Flags().
		StringVar(&outputPath, "output", "", "Output directory for generated code (required)")
	generateCmd.Flags().
		StringVar(&specPath, "spec", "", "Path to OpenAPI specification file (default: api/openapi.bundled.yaml)")
	generateCmd.Flags().
		BoolVar(&bundleFlag, "bundle", false, "Bundle the OpenAPI spec into a single file instead of generating code")
	generateCmd.Flags().
		BoolVar(&orvalFix, "orval-fix", false, "Apply fixes for Orval compatibility (only used with --bundle)")
	generateCmd.Flags().
		BoolVar(&dryRun, "dry-run", false, "Show what would be generated without writing files")
	_ = generateCmd.MarkFlagRequired("output")
}
