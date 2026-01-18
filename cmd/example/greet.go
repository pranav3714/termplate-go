package example

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/blacksilver/termplate-go/internal/handler"
)

var (
	name      string
	uppercase bool
)

var greetCmd = &cobra.Command{
	Use:   "greet",
	Short: "Greet a user",
	Long: `Greet a user with a personalized message.

Examples:
  ever-so-powerful-go example greet --name John
  ever-so-powerful-go example greet --name Jane --uppercase`,

	Args: cobra.NoArgs,

	PreRunE: func(_ *cobra.Command, _ []string) error {
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, _ []string) error {
		return runGreet(cmd.Context())
	},
}

func init() {
	greetCmd.Flags().StringVarP(&name, "name", "n", "", "name to greet (required)")
	greetCmd.Flags().BoolVarP(&uppercase, "uppercase", "u", false, "convert message to uppercase")

	_ = greetCmd.MarkFlagRequired("name")
}

func runGreet(ctx context.Context) error {
	slog.Debug("greeting user",
		"name", name,
		"uppercase", uppercase,
	)

	h := handler.NewGreetHandler()
	result, err := h.Greet(ctx, handler.GreetInput{
		Name:      name,
		Uppercase: uppercase,
	})
	if err != nil {
		return fmt.Errorf("greeting user: %w", err)
	}

	fmt.Println(result.Message)
	return nil
}
