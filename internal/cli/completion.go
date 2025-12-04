// Package cli provides command-line interface commands for the Arches platform
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command.
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `To load completions:

Bash:
  $ source <(archesai completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ archesai completion bash > /etc/bash_completion.d/archesai
  # macOS:
  $ archesai completion bash > $(brew --prefix)/etc/bash_completion.d/archesai

Zsh:
  $ source <(archesai completion zsh)
  # To load completions for each session, execute once:
  $ archesai completion zsh > "${fpath[1]}/_archesai"

Fish:
  $ archesai completion fish | source
  # To load completions for each session, execute once:
  $ archesai completion fish > ~/.config/fish/completions/archesai.fish

PowerShell:
  PS> archesai completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> archesai completion powershell > archesai.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:                  runCompletion,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	}
	return nil
}
