package cmd

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/impairments"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete impairments on a containerlab interface",
	Run: func(cmd *cobra.Command, args []string) {
		err, defaultConfig := config.NewDefaultConfig()
		if err != nil {
			log.Fatalf("Error creating config: %v\n", err)
		}
		helper := helpers.NewDefaultHelper()
		command := command.NewDefaultSetCommand(Node, Interface, defaultConfig.GetValue(helper.GetDefaultClabNameKey()))
		manager := impairments.NewDefaultSetter(defaultConfig, Node, Interface, helper, command)
		// Delete is setting all values to 0
		handleError(manager.SetDelay(0), manager, "Error setting delay")
		handleError(manager.SetJitter(0), manager, "Error setting jitter")
		handleError(manager.SetLoss(0), manager, "Error setting loss")
		handleError(manager.SetRate(0), manager, "Error setting rate")
		handleError(manager.ApplyImpairments(), manager, "Error applying impairments")
		handleError(manager.WriteConfig(), manager, "Error writing config")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&Node, "node", "n", "", "node to delete the impairments from ")
	deleteCmd.Flags().StringVarP(&Interface, "interface", "i", "", "interface to delete the impairments from")
	markRequiredFlags(deleteCmd, []string{"node", "interface"})
}
