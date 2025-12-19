package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"

	"github.com/archesai/archesai/cmd/archesai/flags"
)

// configShowCmd shows the current configuration.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: `Display the current configuration with all applied defaults,
environment variables, and config file values.

Examples:
  archesai config show
  archesai config show --output json`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runConfigShow,
}

func init() {
	configCmd.AddCommand(configShowCmd)
	flags.SetConfigShowFlags(configShowCmd)
}

func runConfigShow(_ *cobra.Command, _ []string) error {
	switch flags.ConfigShow.OutputFormat {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(Config.Config)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		defer func() {
			if err := encoder.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close encoder: %v\n", err)
			}
		}()
		return encoder.Encode(Config.Config)
	default:
		return fmt.Errorf(
			"unsupported output format: %s (supported: yaml, json)",
			flags.ConfigShow.OutputFormat,
		)
	}
}
