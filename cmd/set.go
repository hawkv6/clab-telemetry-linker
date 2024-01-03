package cmd

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/spf13/cobra"
)

var (
	Node      string
	Interface string
	Delay     uint64
	Jitter    uint64
	Loss      float64
	Rate      uint64
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set impairments on a containerlab interface",
	Run: func(cmd *cobra.Command, args []string) {
		command.CreateBasicCommand(Node, Interface)

		config.SetPrefix(Node, Interface)
		config.SetDelay(Delay)
		config.SetJitter(Jitter)
		config.SetLoss(Loss)
		config.SetRate(Rate)

		command.AddDelay(Delay)
		command.AddJitter(Jitter)
		command.AddLoss(Loss)
		command.AddRate(Rate)
		command.ExecuteCommand()

		if err := config.WriteConfig(); err != nil {
			command.CreateBasicCommand(Node, Interface)
			command.ExecuteCommand()
			log.Fatalf("All settings reverted due to error: %v\n", err)
		}
	},
}

func markRequiredFlags(flags []string) {
	for _, flag := range flags {
		if err := setCmd.MarkFlagRequired(flag); err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&Node, "node", "n", "", "node to apply the impairment to ")
	setCmd.Flags().StringVarP(&Interface, "interface", "i", "", "interface to apply the impairment to")
	setCmd.Flags().Uint64VarP(&Delay, "delay", "d", 0, "outgoing delay in ms")
	setCmd.Flags().Uint64VarP(&Jitter, "jitter", "j", 0, "outgoing delay variation (jitter) in ms")
	setCmd.Flags().Float64VarP(&Loss, "loss", "l", 0, "packet loss in %")
	setCmd.Flags().Uint64VarP(&Rate, "rate", "r", 0, "link rate / bandwidth in kbit/s")

	markRequiredFlags([]string{"node", "interface"})
}
