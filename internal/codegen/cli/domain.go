package cli

import (
	"fmt"
	"log"

	"github.com/archesai/archesai/internal/codegen/domain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"d", "scaffold"},
	Short:   "Generate domain scaffolding for a new domain",
	Long: `Generate a complete domain structure including entities, service,
repository, handlers, and optionally middleware and events.

This command creates the boilerplate code needed to implement a new domain
following the project's Domain-Driven Design patterns.`,
	Example: `  codegen domain --name=billing --tables=subscription,invoice
  codegen domain -n payment -t order,transaction --auth --events
  codegen d -n users  # Using alias with short flags`,
	RunE: runDomain,
}

var (
	domainName        string
	domainDescription string
	domainTables      string
	domainAuth        bool
	domainEvents      bool
)

func init() {
	rootCmd.AddCommand(domainCmd)

	// Local flags
	domainCmd.Flags().StringVarP(&domainName, "name", "n", "", "Domain name (required)")
	domainCmd.Flags().StringVarP(&domainDescription, "desc", "d", "", "Domain description")
	domainCmd.Flags().StringVarP(&domainTables, "tables", "t", "", "Comma-separated list of database tables")
	domainCmd.Flags().BoolVar(&domainAuth, "auth", false, "Include auth middleware")
	domainCmd.Flags().BoolVar(&domainEvents, "events", false, "Include domain events")

	// Mark required
	if err := domainCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("Failed to mark name flag as required: %v", err)
	}

	// Bind to viper
	if err := viper.BindPFlag("domain.name", domainCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Failed to bind name flag: %v", err)
	}
	if err := viper.BindPFlag("domain.desc", domainCmd.Flags().Lookup("desc")); err != nil {
		log.Fatalf("Failed to bind desc flag: %v", err)
	}
	if err := viper.BindPFlag("domain.tables", domainCmd.Flags().Lookup("tables")); err != nil {
		log.Fatalf("Failed to bind tables flag: %v", err)
	}
	if err := viper.BindPFlag("domain.auth", domainCmd.Flags().Lookup("auth")); err != nil {
		log.Fatalf("Failed to bind auth flag: %v", err)
	}
	if err := viper.BindPFlag("domain.events", domainCmd.Flags().Lookup("events")); err != nil {
		log.Fatalf("Failed to bind events flag: %v", err)
	}
}

func runDomain(_ *cobra.Command, _ []string) error {
	config := domain.Config{
		Name:        viper.GetString("domain.name"),
		Description: viper.GetString("domain.desc"),
		Tables:      viper.GetString("domain.tables"),
		HasAuth:     viper.GetBool("domain.auth"),
		HasEvents:   viper.GetBool("domain.events"),
	}

	generator := domain.NewGenerator()
	if err := generator.Generate(config); err != nil {
		return err
	}

	fmt.Printf("✅ Domain '%s' generated successfully!\n", config.Name)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Add adapter configuration to internal/adapters.yaml")
	fmt.Println("2. Wire dependencies in internal/app/deps.go")
	fmt.Println("3. Run 'make generate' to generate adapters")
	fmt.Println("4. Implement business logic in service.go")

	return nil
}
