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
	"github.com/archesai/archesai/internal/spec"
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

Use --only to generate specific components (comma-separated):
  go.mod, models, repositories, postgres, sqlite, application, controllers,
  hcl, sqlc, client, app, container, routes, wire, bootstrap (alias for app,container,routes,wire)

By default (no --only flag), all components are generated.

The --lint flag enables strict OpenAPI linting. If ANY violations are found,
code generation will be blocked.

Configuration can be provided via arches.yaml:
  generation:
    spec: ./api/openapi.yaml
    includes:
      - server
      - auth
    only: []
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
	// Load options from arches.yaml, CLI flags override
	opts := codegen.LoadOptionsFromConfig(flags.Root.ConfigFile)

	// CLI flags take precedence
	if flags.Generate.SpecPath != "" {
		opts.SpecPath = flags.Generate.SpecPath
	}
	if flags.Generate.OutputPath != "" {
		opts.OutputPath = flags.Generate.OutputPath
	}
	if flags.Generate.Only != "" {
		opts.Only = strings.Split(flags.Generate.Only, ",")
	}
	if flags.Generate.Lint {
		opts.Lint = true
	}

	if opts.SpecPath == "" {
		return fmt.Errorf(
			"spec path is required: use --spec flag or set generation.spec in arches.yaml",
		)
	}

	// Resolve spec path relative to working directory if not absolute
	if opts.WorkDir != "" && !filepath.IsAbs(opts.SpecPath) {
		opts.SpecPath = filepath.Join(opts.WorkDir, opts.SpecPath)
	}

	// Default output path to working directory
	if opts.OutputPath == "" {
		opts.OutputPath = opts.WorkDir
	} else if opts.WorkDir != "" && !filepath.IsAbs(opts.OutputPath) {
		opts.OutputPath = filepath.Join(opts.WorkDir, opts.OutputPath)
	}

	// Load OpenAPI document for bundling
	doc, err := spec.NewOpenAPIDocumentWithIncludes(opts.SpecPath, opts.Includes)
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
	specDir := filepath.Dir(opts.SpecPath)
	bundledPath := filepath.Join(specDir, "openapi.bundled.yaml")
	if err := os.WriteFile(bundledPath, openapiData, 0644); err != nil {
		return fmt.Errorf("failed to write bundled OpenAPI spec: %w", err)
	}

	// Lint bundled spec with vacuum
	if err := runLintStep(bundledPath); err != nil {
		return err
	}

	// Parse OpenAPI spec for code generation
	parser := spec.NewParser()
	if len(opts.Includes) > 0 {
		parser = parser.WithIncludes(opts.Includes)
	}

	s, err := parser.Parse(opts.SpecPath)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	result, err := codegen.Run(opts, s)
	if err != nil {
		return err
	}

	// Run frontend build steps if package.json was generated
	packageJSONPath := filepath.Join("web", "package.json")
	if result != nil && result.WasFileWritten(packageJSONPath) {
		webDir := result.FullPath("web")

		if err := runFrontendStep(webDir, "pnpm", "pnpm", "install", "--silent"); err != nil {
			return err
		}

		if err := runFrontendStep(webDir, "orval", "pnpm", "generate"); err != nil {
			return err
		}

		if err := runFrontendStep(webDir, "biome", "pnpm", "biome", "check", "--write", "--unsafe", "--diagnostic-level=error", "."); err != nil {
			return err
		}

		// Only run vite build if routeTree.gen.ts doesn't exist
		routeTreePath := filepath.Join(webDir, "src", "routeTree.gen.ts")
		if _, err := os.Stat(routeTreePath); os.IsNotExist(err) {
			if err := runFrontendStep(webDir, "vite", "pnpm", "build"); err != nil {
				return err
			}
		}

		if err := runFrontendStep(webDir, "tsc", "pnpm", "typecheck"); err != nil {
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

// runFrontendStep runs a command in the frontend directory with prefixed output.
func runFrontendStep(dir, name, command string, args ...string) error {
	cmd := exec.Command(command, args...)
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
