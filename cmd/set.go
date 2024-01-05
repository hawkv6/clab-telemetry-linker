package cmd

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/impairments"
	"github.com/spf13/cobra"
)

func handleError(err error, manager *impairments.DefaultSetter, message string) {
	if err != nil {
		log.Errorf("%s: %v\n", message, err)
		if err := manager.DeleteImpairments(); err != nil {
			log.Fatalf("Error reverting impairments: %v\n", err)
		}
		log.Fatalf("All settings reverted due to error: %v\n", err)
	}
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set impairments on a containerlab interface",
	Run: func(cmd *cobra.Command, args []string) {
		err, defaultConfig := config.NewDefaultConfig()
		if err != nil {
			log.Fatalf("Error creating config: %v\n", err)
		}
		helper := helpers.NewDefaultHelper()
		command := command.NewDefaultSetCommand(Node, Interface, defaultConfig.GetValue(helper.GetDefaultClabNameKey()))
		manager := impairments.NewDefaultSetter(defaultConfig, Node, Interface, helper, command)
		handleError(manager.SetDelay(Delay), manager, "Error setting delay")
		handleError(manager.SetJitter(Jitter), manager, "Error setting jitter")
		handleError(manager.SetLoss(Loss), manager, "Error setting loss")
		handleError(manager.SetRate(Rate), manager, "Error setting rate")
		handleError(manager.ApplyImpairments(), manager, "Error applying impairments")
		handleError(manager.WriteConfig(), manager, "Error writing config")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&Node, "node", "n", "", "node to apply the impairment to ")
	setCmd.Flags().StringVarP(&Interface, "interface", "i", "", "interface to apply the impairment to")
	setCmd.Flags().Uint64VarP(&Delay, "delay", "d", 0, "outgoing delay in ms")
	setCmd.Flags().Uint64VarP(&Jitter, "jitter", "j", 0, "outgoing delay variation (jitter) in ms")
	setCmd.Flags().Float64VarP(&Loss, "loss", "l", 0, "packet loss in %")
	setCmd.Flags().Uint64VarP(&Rate, "rate", "r", 0, "link rate / bandwidth in kbit/s")

	markRequiredFlags(setCmd, []string{"node", "interface"})
}
