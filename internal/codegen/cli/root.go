// Package cli provides the CLI interface for the codegen tool.
package cli

import (
	"fmt"
	"os"

	"github.com/archesai/archesai/internal/codegen"
	"github.com/spf13/cobra"
)

// Execute runs the codegen CLI.
func Execute() {
	var configPath string
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   "codegen",
		Short: "Generate code from OpenAPI specs",
		Long: `A unified code generation tool for ArchesAI that generates
repository interfaces, cache implementations, event publishers,
and other boilerplate code from OpenAPI specifications.`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := codegen.Run(configPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "codegen.yaml", "Config file path")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
