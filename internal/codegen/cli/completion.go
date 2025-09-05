package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `To load completions:

Bash:
  $ source <(codegen completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ codegen completion bash > /etc/bash_completion.d/codegen
  # macOS:
  $ codegen completion bash > $(brew --prefix)/etc/bash_completion.d/codegen

Zsh:
  $ source <(codegen completion zsh)
  # To load completions for each session, execute once:
  $ codegen completion zsh > "${fpath[1]}/_codegen"

Fish:
  $ codegen completion fish | source
  # To load completions for each session, execute once:
  $ codegen completion fish > ~/.config/fish/completions/codegen.fish

PowerShell:
  PS> codegen completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> codegen completion powershell > codegen.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
