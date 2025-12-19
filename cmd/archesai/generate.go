package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/internal/located"
	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/pkg/storage"
)

// genPrefixColors maps tool names to lipgloss colors for log output
var genPrefixColors = map[string]lipgloss.Color{
	"vacuum": lipgloss.Color("135"), // purple
	"pnpm":   lipgloss.Color("208"), // orange
	"orval":  lipgloss.Color("213"), // pink
	"biome":  lipgloss.Color("39"),  // cyan
	"vite":   lipgloss.Color("226"), // yellow
	"tsc":    lipgloss.Color("33"),  // blue
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from OpenAPI specification",
	Long: `Generate Go code from an OpenAPI specification.

This command generates:
- Models
- Repositories
- HTTP Handlers
- Application Handlers
- React frontend and TypeScript client
- Database schema (HCL and SQLC)
- Bootstrap code (app, container, routes, wire)

The --lint flag enables strict OpenAPI linting. If ANY violations are found,
code generation will be blocked.

Configuration can be provided via arches.yaml:
  generation:
    spec: ./spec/openapi.yaml
    includes:
      - server
      - auth
    groups:
      module: true
      schemas: true
      http: true
      wire: true
      postgres: true
      sqlite: true
      web: true
    lint: false`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	flags.SetGenerateFlags(generateCmd)
}

func runGenerate(_ *cobra.Command, _ []string) error {
	if Config == nil {
		return fmt.Errorf("config not loaded: ensure arches.yaml exists")
	}

	// Initialize generation config if nil
	gen := Config.Config.Generation

	// CLI flags take precedence
	specPath := ""
	if gen.Spec != nil {
		specPath = *gen.Spec
	}
	if flags.Generate.SpecPath != "" {
		specPath = flags.Generate.SpecPath
	}

	outputPath := ""
	if gen.Output != nil {
		outputPath = *gen.Output
	}
	if flags.Generate.OutputPath != "" {
		outputPath = flags.Generate.OutputPath
	}

	if flags.Generate.Lint {
		lint := true
		gen.Lint = &lint
	}

	if specPath == "" {
		return fmt.Errorf(
			"spec path is required: use --spec flag or set generation.spec in arches.yaml",
		)
	}

	workDir := Config.WorkDir()

	// Resolve spec path relative to working directory if not absolute
	if workDir != "" && !filepath.IsAbs(specPath) {
		specPath = filepath.Join(workDir, specPath)
	}

	// Default output path to working directory
	if outputPath == "" {
		outputPath = workDir
	} else if workDir != "" && !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(workDir, outputPath)
	}

	// Build composite filesystem with includes
	baseFS := os.DirFS(filepath.Dir(specPath))
	compositeFS := spec.BuildIncludeFS(baseFS, gen.Includes)

	// Load OpenAPI document
	doc, err := spec.NewOpenAPIDocumentFromFS(compositeFS, filepath.Base(specPath))
	if err != nil {
		return fmt.Errorf("failed to load OpenAPI document: %w", err)
	}

	// Generate bundled OpenAPI spec
	bundler := spec.NewBundler(doc)
	openapiData, err := bundler.BundleToYAML()
	if err != nil {
		return fmt.Errorf("failed to generate bundled OpenAPI spec: %w", err)
	}

	// Write bundled spec
	specDir := filepath.Dir(specPath)
	bundledPath := filepath.Join(specDir, "openapi.bundled.yaml")
	if err := os.WriteFile(bundledPath, openapiData, 0644); err != nil {
		return fmt.Errorf("failed to write bundled OpenAPI spec: %w", err)
	}

	// Lint bundled spec with vacuum
	if err := runLintStep(bundledPath); err != nil {
		return err
	}

	// Parse OpenAPI spec for code generation
	specParser := spec.NewParser(doc).WithIncludes(gen.Includes)
	if gen.Output != nil {
		specParser = specParser.WithCodegenOutput(*gen.Output)
	}
	s, err := specParser.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Setup and run code generator
	cg, err := codegen.New(outputPath)
	if err != nil {
		return fmt.Errorf("failed to initialize code generator: %w", err)
	}

	locatedSpec := &located.Located[spec.Spec]{Value: s, Path: specPath}
	if err := cg.Generate(locatedSpec, *gen); err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}

	// Get storage for file tracking
	trackedStorage, _ := cg.GetStorage().(*storage.TrackedStorage)

	// Run frontend build steps if package.json was generated
	packageJSONPath := filepath.Join("web", "package.json")
	if trackedStorage != nil && trackedStorage.WasFileWritten(packageJSONPath) {
		webDir := filepath.Join(trackedStorage.BaseDir(), "web")

		if err := runFrontendStep(webDir, "pnpm", "install", "--silent"); err != nil {
			return err
		}

		if err := runFrontendStep(webDir, "orval", "generate"); err != nil {
			return err
		}

		if err := runFrontendStep(webDir, "biome", "biome", "check", "--write", "--unsafe", "--diagnostic-level=error", "."); err != nil {
			return err
		}

		// Only run vite build if routeTree.gen.ts doesn't exist
		routeTreePath := filepath.Join(webDir, "src", "routeTree.gen.ts")
		if _, err := os.Stat(routeTreePath); os.IsNotExist(err) {
			if err := runFrontendStep(webDir, "vite", "build"); err != nil {
				return err
			}
		}

		if err := runFrontendStep(webDir, "tsc", "typecheck"); err != nil {
			return err
		}
	}

	return nil
}

// genStepDescriptions maps tool names to human-readable descriptions
var genStepDescriptions = map[string]string{
	"vacuum": "Linting OpenAPI spec",
	"pnpm":   "Installing dependencies",
	"orval":  "Generating API client",
	"biome":  "Linting and formatting",
	"vite":   "Building frontend",
	"tsc":    "Type checking",
}

// runLintStep runs vacuum lint on the bundled OpenAPI spec.
func runLintStep(specPath string) error {
	cmd := exec.Command("vacuum", "lint", specPath,
		"--details",
		"--no-banner",
		"--hard-mode",
		"--no-clip",
		"--all-results",
		"--pipeline-output",
		"--no-style",
		"--errors",
	)

	name := "vacuum"
	color := genPrefixColors[name]
	prefixStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
	prefix := prefixStyle.Render(fmt.Sprintf("[%s]", name))
	errorPrefix := makeErrorPrefix(name)

	if flags.Generate.Verbose {
		fmt.Printf("%s %s...\n", prefix, genStepDescriptions[name])
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s failed to start: %w", name, err)
	}

	var hasOutput bool
	var outputMu sync.Mutex

	done := make(chan struct{})
	go func() {
		if flags.Generate.Verbose {
			if streamWithPrefix(stdout, prefix, errorPrefix) {
				outputMu.Lock()
				hasOutput = true
				outputMu.Unlock()
			}
		} else {
			if streamErrorsOnly(stdout, errorPrefix) {
				outputMu.Lock()
				hasOutput = true
				outputMu.Unlock()
			}
		}
		done <- struct{}{}
	}()
	go func() {
		if streamWithPrefix(stderr, errorPrefix, errorPrefix) {
			outputMu.Lock()
			hasOutput = true
			outputMu.Unlock()
		}
		done <- struct{}{}
	}()

	<-done
	<-done

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}

	if flags.Generate.Verbose && !hasOutput {
		fmt.Printf("%s Done\n", prefix)
	}

	return nil
}

// streamErrorsOnly reads from r and prints only error lines with the given prefix.
// Returns true if any output was printed.
func streamErrorsOnly(r io.Reader, errorPrefix string) bool {
	scanner := bufio.NewScanner(r)
	hasOutput := false
	inErrorBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			// Track error detail blocks (vacuum output)
			if strings.Contains(line, "<details><summary>🔴") {
				inErrorBlock = true
			}

			// Only print if it's an error line or in error block
			isError := inErrorBlock ||
				strings.Contains(line, "error TS") ||
				strings.Contains(line, ": error") ||
				strings.Contains(line, "Error:") ||
				strings.Contains(line, "errors detected") ||
				strings.Contains(line, "🔴")

			if isError {
				fmt.Printf("%s %s\n", errorPrefix, line)
				hasOutput = true
			}

			if strings.Contains(line, "</details>") && inErrorBlock {
				inErrorBlock = false
			}
		}
	}
	return hasOutput
}

// runFrontendStep runs a pnpm command in the frontend directory with prefixed output.
func runFrontendStep(dir, name string, args ...string) error {
	cmd := exec.Command("pnpm", args...)
	cmd.Dir = dir

	// Get color for this tool
	color, ok := genPrefixColors[name]
	if !ok {
		color = lipgloss.Color("7") // default gray
	}
	prefixStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
	prefix := prefixStyle.Render(fmt.Sprintf("[%s]", name))

	// Error prefix in red with ERR marker
	errorPrefix := makeErrorPrefix(name)

	// Show step description if verbose
	if flags.Generate.Verbose {
		desc := genStepDescriptions[name]
		if desc != "" {
			fmt.Printf("%s %s...\n", prefix, desc)
		}
	}

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s failed to start: %w", name, err)
	}

	// Track if any output was shown
	var hasOutput bool
	var outputMu sync.Mutex

	// Stream stdout and stderr with prefix
	done := make(chan struct{})
	go func() {
		if streamWithPrefix(stdout, prefix, errorPrefix) {
			outputMu.Lock()
			hasOutput = true
			outputMu.Unlock()
		}
		done <- struct{}{}
	}()
	go func() {
		if streamWithPrefix(stderr, errorPrefix, errorPrefix) {
			outputMu.Lock()
			hasOutput = true
			outputMu.Unlock()
		}
		done <- struct{}{}
	}()

	// Wait for streaming to complete
	<-done
	<-done

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}

	// Show completion in verbose mode for steps with no output
	if flags.Generate.Verbose && !hasOutput {
		fmt.Printf("%s Done\n", prefix)
	}

	return nil
}

// streamWithPrefix reads from r and prints each line with the given prefix.
// Returns true if any output was printed.
func streamWithPrefix(r io.Reader, prefix, errorPrefix string) bool {
	scanner := bufio.NewScanner(r)
	hasOutput := false
	inErrorBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			// Track error detail blocks (vacuum output)
			if strings.Contains(line, "<details><summary>🔴") {
				inErrorBlock = true
			}
			if strings.Contains(line, "</details>") && inErrorBlock {
				// Print this line as error, then exit block
				fmt.Printf("%s %s\n", errorPrefix, line)
				hasOutput = true
				inErrorBlock = false
				continue
			}

			// Use error prefix for lines containing error patterns or inside error block
			p := prefix
			if inErrorBlock ||
				strings.Contains(line, "error TS") ||
				strings.Contains(line, ": error") ||
				strings.Contains(line, "Error:") ||
				strings.Contains(line, "errors detected") ||
				strings.Contains(line, "🔴") {
				p = errorPrefix
			}
			fmt.Printf("%s %s\n", p, line)
			hasOutput = true
		}
	}
	return hasOutput
}

// makeErrorPrefix creates a styled error prefix with ERR marker.
func makeErrorPrefix(name string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	return style.Render(fmt.Sprintf("[%s ERR]", name))
}
