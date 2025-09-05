package cli

import (
	"log"

	"github.com/archesai/archesai/internal/codegen/converters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// convertersCmd represents the converters command
var convertersCmd = &cobra.Command{
	Use:     "converters",
	Aliases: []string{"conv", "adapters"},
	Short:   "Generate type converters between database and API models",
	Long: `Generate type converter functions (adapters) for converting between
database models and API models.

This command reads converter configuration from a YAML file and generates
type-safe converter functions for each domain.`,
	Example: `  codegen converters
  codegen converters --config=internal/domains/adapters.yaml
  codegen conv  # Using alias`,
	RunE: runConverters,
}

var (
	convertersConfig  string
	convertersDomains []string
)

func init() {
	rootCmd.AddCommand(convertersCmd)

	// Local flags
	convertersCmd.Flags().StringVar(&convertersConfig, "config", "internal/domains/adapters.yaml", "Converters configuration file")
	convertersCmd.Flags().StringSliceVar(&convertersDomains, "domains", []string{"auth", "organizations", "workflows", "content"}, "Domains to generate converters for")

	// Bind to viper
	if err := viper.BindPFlag("converters.config", convertersCmd.Flags().Lookup("config")); err != nil {
		log.Fatalf("Failed to bind config flag: %v", err)
	}
	if err := viper.BindPFlag("converters.domains", convertersCmd.Flags().Lookup("domains")); err != nil {
		log.Fatalf("Failed to bind domains flag: %v", err)
	}
}

func runConverters(_ *cobra.Command, _ []string) error {
	generator := converters.NewGenerator()
	if err := generator.Generate(); err != nil {
		return err
	}

	log.Println("✅ Converters generated successfully")
	return nil
}
