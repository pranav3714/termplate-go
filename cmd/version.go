package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/blacksilver/termplate-go/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version, commit, build date, and Go version.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		info := version.Get()

		switch output {
		case "json":
			data, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling to JSON: %w", err)
			}
			fmt.Println(string(data))
		case "yaml":
			data, err := yaml.Marshal(info)
			if err != nil {
				return fmt.Errorf("marshaling to YAML: %w", err)
			}
			fmt.Print(string(data))
		default:
			fmt.Printf("Termplate Go %s\n", info.String())
		}

		return nil
	},
}
