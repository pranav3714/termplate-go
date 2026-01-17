package example

import "github.com/spf13/cobra"

// Cmd is the parent command for example operations
var Cmd = &cobra.Command{
	Use:   "example",
	Short: "Example command demonstrating CLI structure",
	Long:  `Example command showing how to implement commands, handlers, and services.`,
}

func init() {
	Cmd.AddCommand(greetCmd)
}
