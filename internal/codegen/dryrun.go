package codegen

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/tui"
	"github.com/archesai/archesai/pkg/storage"
)

// printDryRunResults displays what would be generated in dry-run mode.
func printDryRunResults(orch *Orchestrator, outputPath string) error {
	memStorage := orch.GetStorage().(*storage.MemoryStorage)
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

// printDryRunResultsTUI displays dry-run results using TUI components.
func printDryRunResultsTUI(runner *tui.Runner, orch *Orchestrator, outputPath string) error {
	memStorage := orch.GetStorage().(*storage.MemoryStorage)
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
