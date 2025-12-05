package main

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/tui"
	"github.com/archesai/archesai/pkg/config"
)

// configShowCmd shows the current configuration.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: `Display the current configuration with all applied defaults,
environment variables, and config file values.

Examples:
  archesai config show
  archesai config show --output json
  archesai config show -o tui`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runConfigShow,
}

func init() {
	configCmd.AddCommand(configShowCmd)
	flags.SetConfigShowFlags(configShowCmd)
}

func runConfigShow(_ *cobra.Command, _ []string) error {
	cfg, err := config.NewParser[any]().Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch flags.ConfigShow.OutputFormat {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(cfg.Config)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		defer func() {
			if err := encoder.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close encoder: %v\n", err)
			}
		}()
		return encoder.Encode(cfg.Config)
	case "tui":
		model := tui.NewConfigModel()
		program := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := program.Run(); err != nil {
			return fmt.Errorf("error running config TUI: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", flags.ConfigShow.OutputFormat)
	}
}
