package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/archesai/archesai/pkg/tui"
)

// tuiCmd represents the tui command.
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI for configuration and AI agents",
	Long: `Launch an interactive terminal user interface (TUI) for viewing
configuration and optionally interacting with AI agents.

The TUI provides:
- Configuration viewer for all settings
- Database, server, and auth configuration display
- AI provider status and configuration
- Optional chat interface with AI agents`,
	Example: `  archesai tui                  # Launch config viewer
  archesai tui --chat            # Launch AI chat interface
  archesai tui --chat --provider=openai --model=gpt-4`,
	RunE: runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(_ *cobra.Command, _ []string) error {
	model := tui.NewConfigModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("error running config TUI: %w", err)
	}
	return nil
}
