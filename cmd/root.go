package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
)

var (
	log       = logging.DefaultLogger.WithField("subsystem", "cmd")
	Node      string
	Interface string
	Delay     uint64
	Jitter    uint64
	Loss      float64
	Rate      uint64
)

func markRequiredFlags(cmd *cobra.Command, flags []string) {
	for _, flag := range flags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			log.Fatal(err)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "clab-telemetry-linker",
	Short: "clab-telemetry-linker is a tool to enrich telemetry data with the underlying containerlab impairments",
	Long: `clab-telemetry-linker is a tool to enrich telemetry data with the underlying containerlab impairments
More detailed info: https://github.com/hawkv6/clab-telemetry-linker
Example usage:
	sudo clab-telemetry-linker start -b 172.16.19.77:9094 -r hawkv6.telemetry.unprocessed -p hawkv6.telemetry.processed
	sudo clab-telemetry-linker set -n XR-1 -i Gi0-0-0-0 --delay 1ms --jitter 1ms --loss 5 --rate 100000 	
	sudo clab-telemetry-linker show -n XR-1 	
	sudo clab-telemetry-linker delete -n XR-1 -i Gi0-0-0-0
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !helpers.NewDefaultHelper().IsRoot() {
			log.Fatalln("You must be root to run this command")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
