// Package main provides the unified code generation tool for ArchesAI.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/archesai/archesai/internal/codegen/converters"
	"github.com/archesai/archesai/internal/codegen/defaults"
	"github.com/archesai/archesai/internal/codegen/domain"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Remove the command from args for flag parsing
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)

	switch command {
	case "defaults":
		runDefaults()
	case "converters":
		runConverters()
	case "domain":
		runDomain()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("ArchesAI Code Generator")
	fmt.Println()
	fmt.Println("Usage: codegen <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  defaults      Generate configuration defaults from OpenAPI schemas")
	fmt.Println("  converters    Generate type converters between database and API models")
	fmt.Println("  domain        Generate domain scaffolding for a new domain")
	fmt.Println()
	fmt.Println("Use 'codegen <command> -h' for more information about a command.")
}

func runDefaults() {
	var help bool
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	if help {
		fmt.Println("Generate configuration defaults from OpenAPI schemas")
		fmt.Println()
		fmt.Println("Usage: codegen defaults")
		fmt.Println()
		fmt.Println("This command parses OpenAPI schema files from api/components/schemas")
		fmt.Println("and generates Go code with default configuration values.")
		fmt.Println()
		fmt.Println("Output: internal/infrastructure/config/defaults.gen.go")
		return
	}

	generator := defaults.NewGenerator()
	if err := generator.Generate(); err != nil {
		log.Fatalf("Failed to generate defaults: %v", err)
	}

	log.Println("✅ Config defaults generated successfully")
}

func runConverters() {
	var help bool
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	if help {
		fmt.Println("Generate type converters between database and API models")
		fmt.Println()
		fmt.Println("Usage: codegen converters")
		fmt.Println()
		fmt.Println("This command reads converter configuration from internal/domains/adapters.yaml")
		fmt.Println("and generates type converter functions for each domain.")
		fmt.Println()
		fmt.Println("Output: internal/domains/*/adapters/adapters.gen.go")
		return
	}

	generator := converters.NewGenerator()
	if err := generator.Generate(); err != nil {
		log.Fatalf("Failed to generate converters: %v", err)
	}

	log.Println("✅ Converters generated successfully")
}

func runDomain() {
	var (
		name        string
		description string
		tables      string
		hasAuth     bool
		hasEvents   bool
		help        bool
	)

	flag.StringVar(&name, "name", "", "Domain name (e.g., billing)")
	flag.StringVar(&description, "desc", "", "Domain description")
	flag.StringVar(&tables, "tables", "", "Comma-separated list of database tables")
	flag.BoolVar(&hasAuth, "auth", false, "Include auth middleware")
	flag.BoolVar(&hasEvents, "events", false, "Include domain events")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	if help || name == "" {
		fmt.Println("Generate domain scaffolding for a new domain")
		fmt.Println()
		fmt.Println("Usage: codegen domain -name=<name> [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -name string    Domain name (required, e.g., billing)")
		fmt.Println("  -desc string    Domain description")
		fmt.Println("  -tables string  Comma-separated list of database tables")
		fmt.Println("  -auth           Include auth middleware")
		fmt.Println("  -events         Include domain events")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  codegen domain -name=billing -tables=subscription,invoice -auth")
		if name == "" && !help {
			os.Exit(1)
		}
		return
	}

	config := domain.Config{
		Name:        name,
		Description: description,
		Tables:      tables,
		HasAuth:     hasAuth,
		HasEvents:   hasEvents,
	}

	generator := domain.NewGenerator()
	if err := generator.Generate(config); err != nil {
		log.Fatalf("Failed to generate domain: %v", err)
	}

	fmt.Printf("✅ Domain '%s' generated successfully!\n", name)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Add converter configuration to internal/domains/adapters.yaml")
	fmt.Println("2. Wire dependencies in internal/app/deps.go")
	fmt.Println("3. Run 'make generate' to generate converters")
	fmt.Println("4. Implement business logic in service.go")
}
