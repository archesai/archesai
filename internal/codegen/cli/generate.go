// Package cli provides the code generation CLI commands.
package cli

import (
	"fmt"
	"log"

	"github.com/archesai/archesai/internal/codegen/cache"
	"github.com/archesai/archesai/internal/codegen/events"
	"github.com/archesai/archesai/internal/codegen/repository"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from OpenAPI specifications",
	Long: `Generate repository interfaces, cache implementations, and event publishers
from OpenAPI specifications and configuration.`,
	Run: runGenerate,
}

// generateFlags holds the flags for the generate command
type generateFlags struct {
	config     string
	repository bool
	cache      bool
	events     bool
	all        bool
}

var genFlags generateFlags

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&genFlags.config, "config", "c", "codegen.yaml", "Configuration file path")
	generateCmd.Flags().BoolVar(&genFlags.repository, "repository", false, "Generate repository code")
	generateCmd.Flags().BoolVar(&genFlags.cache, "cache", false, "Generate cache code")
	generateCmd.Flags().BoolVar(&genFlags.events, "events", false, "Generate events code")
	generateCmd.Flags().BoolVar(&genFlags.all, "all", false, "Generate all code")
}

func runGenerate(cmd *cobra.Command, args []string) {
	// If --all is specified, enable all generators
	if genFlags.all {
		genFlags.repository = true
		genFlags.cache = true
		genFlags.events = true
	}

	// If no specific generators are specified, generate all
	if !genFlags.repository && !genFlags.cache && !genFlags.events {
		genFlags.repository = true
		genFlags.cache = true
		genFlags.events = true
	}

	// Run repository generator
	if genFlags.repository {
		log.Println("Generating repository code...")
		gen := repository.NewGenerator()
		if err := gen.Generate(genFlags.config); err != nil {
			log.Fatalf("Failed to generate repository code: %v", err)
		}
		log.Println("Repository code generated successfully")
	}

	// Run cache generator
	if genFlags.cache {
		log.Println("Generating cache code...")
		gen := cache.NewGenerator()
		if err := gen.Generate(genFlags.config); err != nil {
			log.Fatalf("Failed to generate cache code: %v", err)
		}
		log.Println("Cache code generated successfully")
	}

	// Run events generator
	if genFlags.events {
		log.Println("Generating events code...")
		gen := events.NewGenerator()
		if err := gen.Generate(genFlags.config); err != nil {
			log.Fatalf("Failed to generate events code: %v", err)
		}
		log.Println("Events code generated successfully")
	}

	fmt.Println("✓ Code generation complete!")
}

// repositoryCmd represents the repository subcommand
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Generate repository interfaces and implementations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating repository code...")
		gen := repository.NewGenerator()
		if err := gen.Generate(genFlags.config); err != nil {
			log.Fatalf("Failed to generate repository code: %v", err)
		}
		log.Println("✓ Repository code generated successfully")
	},
}

// cacheCmd represents the cache subcommand
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Generate cache interfaces and implementations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating cache code...")
		gen := cache.NewGenerator()
		if err := gen.Generate(genFlags.config); err != nil {
			log.Fatalf("Failed to generate cache code: %v", err)
		}
		log.Println("✓ Cache code generated successfully")
	},
}

// eventsCmd represents the events subcommand
var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Generate event publisher interfaces and implementations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating events code...")
		gen := events.NewGenerator()
		if err := gen.Generate(genFlags.config); err != nil {
			log.Fatalf("Failed to generate events code: %v", err)
		}
		log.Println("✓ Events code generated successfully")
	},
}

func init() {
	generateCmd.AddCommand(repositoryCmd)
	generateCmd.AddCommand(cacheCmd)
	generateCmd.AddCommand(eventsCmd)

	// Add config flag to all subcommands
	repositoryCmd.Flags().StringVarP(&genFlags.config, "config", "c", "codegen.yaml", "Configuration file path")
	cacheCmd.Flags().StringVarP(&genFlags.config, "config", "c", "codegen.yaml", "Configuration file path")
	eventsCmd.Flags().StringVarP(&genFlags.config, "config", "c", "codegen.yaml", "Configuration file path")
}
