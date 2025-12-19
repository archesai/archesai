package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// specLintCmd represents the spec lint command
var specLintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint OpenAPI specification with vacuum",
	Long: `Lint an OpenAPI specification using vacuum.

This command runs vacuum lint on the bundled OpenAPI specification.
The bundled spec must exist (run 'archesai generate' first).

Examples:
  archesai spec lint
  archesai spec lint --config arches.yaml`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSpecLint,
}

func init() {
	specCmd.AddCommand(specLintCmd)
}

func runSpecLint(_ *cobra.Command, _ []string) error {
	if Config == nil {
		return fmt.Errorf("config not loaded: ensure arches.yaml exists")
	}

	// Get generation config
	gen := Config.Config.Generation

	specPath := ""
	if gen != nil && gen.Spec != nil {
		specPath = *gen.Spec
	}

	if specPath == "" {
		return fmt.Errorf(
			"spec path is required: set generation.spec in arches.yaml",
		)
	}

	// Resolve spec path relative to working directory if not absolute
	if Config.WorkDir() != "" && !filepath.IsAbs(specPath) {
		specPath = filepath.Join(Config.WorkDir(), specPath)
	}

	// Find bundled spec path
	specDir := filepath.Dir(specPath)
	bundledPath := filepath.Join(specDir, "openapi.bundled.yaml")

	// Check if bundled spec exists
	if _, err := os.Stat(bundledPath); os.IsNotExist(err) {
		return fmt.Errorf(
			"bundled spec not found at %s: run 'archesai generate' first",
			bundledPath,
		)
	}

	// Run vacuum lint
	cmd := exec.Command("vacuum", "lint", bundledPath,
		"--details",
		"--no-banner",
		"--hard-mode",
		"--no-clip",
		"--all-results",
		"--no-style",
		// "--errors",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("vacuum lint failed: %w", err)
	}

	return nil
}
