package cmd

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/impairments"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show impairments on a containerlab node",
	Run: func(cmd *cobra.Command, args []string) {
		defaultConfig, err := config.NewDefaultConfig()
		if err != nil {
			log.Fatalf("Error reading/creating config: %v\n", err)
		}
		helper := helpers.NewDefaultHelper()
		command := command.NewDefaultShowCommand(Node, defaultConfig.GetValue(helper.GetDefaultClabNameKey()))
		manager := impairments.NewDefaultViewer(Node, command)
		if err := manager.ShowImpairments(); err != nil {
			log.Fatalf("Error showing impairments: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().StringVarP(&Node, "node", "n", "", "node to apply the impairment to ")
	markRequiredFlags(showCmd, []string{"node"})
}
