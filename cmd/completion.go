package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `To load completions:

Bash:
  $ source <(mycli completion bash)
  $ mycli completion bash > /etc/bash_completion.d/mycli

Zsh:
  $ mycli completion zsh > "${fpath[1]}/_mycli"
  $ source ~/.zshrc

Fish:
  $ mycli completion fish | source
  $ mycli completion fish > ~/.config/fish/completions/mycli.fish

PowerShell:
  PS> mycli completion powershell | Out-String | Invoke-Expression
`,
	Args:                  cobra.ExactArgs(1),
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	DisableFlagsInUseLine: true,
	RunE: func(_ *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletion(os.Stdout)
		}
		return nil
	},
}
