package cli

import (
	"log"

	"github.com/archesai/archesai/internal/codegen/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Generate configuration defaults from OpenAPI schemas",
	Long: `Generate Go code with default configuration values by parsing
OpenAPI schema files.

This command reads schema files from api/components/schemas and generates
a Go file with default values extracted from the schemas.`,
	Example: `  codegen defaults
  codegen defaults --output=internal/config/defaults.gen.go`,
	RunE: runDefaults,
}

var (
	defaultsOutput string
	schemasPath    string
)

func init() {
	rootCmd.AddCommand(defaultsCmd)

	// Local flags
	defaultsCmd.Flags().StringVar(&defaultsOutput, "output", "internal/config/defaults.gen.go", "Output file path")
	defaultsCmd.Flags().StringVar(&schemasPath, "schemas", "api/components/schemas", "Path to OpenAPI schemas directory")

	// Bind to viper
	if err := viper.BindPFlag("defaults.output", defaultsCmd.Flags().Lookup("output")); err != nil {
		log.Fatalf("Failed to bind output flag: %v", err)
	}
	if err := viper.BindPFlag("defaults.schemas", defaultsCmd.Flags().Lookup("schemas")); err != nil {
		log.Fatalf("Failed to bind schemas flag: %v", err)
	}
}

func runDefaults(_ *cobra.Command, _ []string) error {
	generator := defaults.NewGenerator()
	if err := generator.Generate(); err != nil {
		return err
	}

	log.Println("✅ Config defaults generated successfully")
	return nil
}
