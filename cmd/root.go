package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/blacksilver/termplate-go/cmd/example"
	"github.com/blacksilver/termplate-go/internal/config"
	"github.com/blacksilver/termplate-go/internal/logger"
)

var (
	cfgFile string
	verbose bool
	output  string
)

var rootCmd = &cobra.Command{
	Use:   "termplate",
	Short: "Termplate Go - A powerful CLI template for developers",
	Long: `Termplate Go is a production-ready CLI tool template built with Go.

It demonstrates best practices for building CLI applications with:
- Cobra for command structure
- Viper for configuration management
- Structured logging with slog
- Clean architecture patterns

Examples:
  termplate --help
  termplate version
  termplate example greet --name World`,

	// Runs before any subcommand
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Skip for completion and help
		if cmd.Name() == "completion" || cmd.Name() == "help" {
			return nil
		}

		// Initialize logger
		level := slog.LevelInfo
		if verbose {
			level = slog.LevelDebug
		}
		logger.Init(level, os.Getenv("ENV") == "production")

		// Bind flags to viper
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return fmt.Errorf("binding flags: %w", err)
		}

		return nil
	},

	SilenceUsage:  true, // Don't show usage on error
	SilenceErrors: true, // We handle errors ourselves
}

// Execute is the entry point called from main
func Execute() error {
	// Set up context with signal handling
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("executing command: %w", err)
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags (available to all subcommands)
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config", "c",
		"",
		"config file (default: $HOME/.termplate.yaml)",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose", "v",
		false,
		"enable verbose output",
	)
	rootCmd.PersistentFlags().StringVarP(
		&output,
		"output", "o",
		"text",
		"output format (text, json, yaml)",
	)

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(example.Cmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("failed to get home directory", "error", err)
			return
		}

		// Search paths
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ever-so-powerful-go")
	}

	// Environment variables
	viper.SetEnvPrefix("TERMPLATE")
	viper.AutomaticEnv()

	// Set defaults
	config.SetDefaults()

	// Read config (ignore if not found)
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("using config file", "file", viper.ConfigFileUsed())
	}
}
