package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
)

var rootCmd = &cobra.Command{
	Use:   "clab-telemetry-linker",
	Short: "clab-telemetry-linker is a tool to enricht telemetry with containerlab impairments",
	Long: `
clab-telemetry-linker is a tool to enricht telemetry with containerlab impairments

Start the tool to listen for telemetry data, enrich it with the containerlab impairements and send it to kafka
clab-telemetry-linker start --input-topic telemetry --output-topic telemetry-enriched --kafka-broker localhost:9092

Set the impairments on the containerlab interface
clab-telemetry-linker set -n clab-hawkv6-XR-1 -i Gi0-0-0-0 --delay 1ms --jitter 1ms --loss 5 --rate 100
-----------+-------+--------+-------------+-------------+
| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
+-----------+-------+--------+-------------+-------------+
| Gi0-0-0-0 | 1ms   | 1ms    | 5.00%       |         100 |
+-----------+-------+--------+-------------+-------------+

clab-telemetry-linker show -n clab-hawkv6-XR-1 -i Gi0-0-0-1 
+-----------+-------+--------+-------------+-------------+
| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
+-----------+-------+--------+-------------+-------------+
| Gi0-0-0-1 | 100ms | 0s     | 10.00%      |           0 |
+-----------+-------+--------+-------------+-------------+

clab-telemetry-linker delete -n clab-hawkv6-XR-1 -i Gi0-0-0-0
+-----------+-------+--------+-------------+-------------+
| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
+-----------+-------+--------+-------------+-------------+
| Gi0-0-0-0 | 0ms   | 0s     | 0.00%       |           0 |
+-----------+-------+--------+-------------+-------------+
	
clab-telemetry-linker will forward the impairments to the specific containerlab command
More information: https://containerlab.dev/cmd/tools/netem/set/
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.DefaultLogger.Fatalln(err)
	}
}

func init() {
	checkIsRoot()
}

func checkIsRoot() {
	if os.Geteuid() != 0 {
		logging.DefaultLogger.Fatalln("Hawkwing must be run as root")
	}
}
