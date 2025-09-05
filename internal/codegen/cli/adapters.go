package cli

import (
	"log"

	"github.com/archesai/archesai/internal/codegen/adapters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// adaptersCmd represents the adapters command
var adaptersCmd = &cobra.Command{
	Use:     "adapters",
	Aliases: []string{"adapters-gen", "converters"},
	Short:   "Generate type adapters between database and API models",
	Long: `Generate type adapter functions for converting between
database models and API models.

This command reads adapter configuration from a YAML file and generates
type-safe adapter functions for each domain.`,
	Example: `  codegen adapters
  codegen adapters --config=internal/adapters.yaml
  codegen converters  # Using alias for backwards compatibility`,
	RunE: runAdapters,
}

var (
	adaptersConfig  string
	adaptersDomains []string
)

func init() {
	rootCmd.AddCommand(adaptersCmd)

	// Local flags
	adaptersCmd.Flags().StringVar(&adaptersConfig, "config", "internal/adapters.yaml", "Adapters configuration file")
	adaptersCmd.Flags().StringSliceVar(&adaptersDomains, "domains", []string{"auth", "organizations", "workflows", "content"}, "Domains to generate adapters for")

	// Bind to viper
	if err := viper.BindPFlag("adapters.config", adaptersCmd.Flags().Lookup("config")); err != nil {
		log.Fatalf("Failed to bind config flag: %v", err)
	}
	if err := viper.BindPFlag("adapters.domains", adaptersCmd.Flags().Lookup("domains")); err != nil {
		log.Fatalf("Failed to bind domains flag: %v", err)
	}
}

func runAdapters(_ *cobra.Command, _ []string) error {
	generator := adapters.NewGenerator()
	if err := generator.Generate(); err != nil {
		return err
	}

	log.Println("✅ Adapters generated successfully")
	return nil
}
