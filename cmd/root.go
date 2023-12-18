package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "clab-mdt-linker",
	Short: "clab-mdt-linker is a tool to enricht mdt with containerlab impairments",
	Long: `
	clab-mdt-linker is a tool to enricht mdt with containerlab impairments

	Start the tool to listen for mdt data, enrich it with the containerlab impairements and send it to kafka
	clab-mdt-linker start --input-topic mdt --output-topic mdt-enriched --kafka-broker localhost:9092

	Set the impairments on the containerlab interface
	clab-mdt-linker set -n clab-hawkv6-XR-1 -i Gi0-0-0-0 --delay 1ms --jitter 1ms --loss 5 --rate 100
	-----------+-------+--------+-------------+-------------+
	| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
	+-----------+-------+--------+-------------+-------------+
	| Gi0-0-0-0 | 1ms   | 1ms    | 5.00%       |         100 |
	+-----------+-------+--------+-------------+-------------+

	clab-mdt-linker show -n clab-hawkv6-XR-1 -i Gi0-0-0-1 
	+-----------+-------+--------+-------------+-------------+
	| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
	+-----------+-------+--------+-------------+-------------+
	| Gi0-0-0-1 | 100ms | 0s     | 10.00%      |           0 |
	+-----------+-------+--------+-------------+-------------+

	clab-mdt-linker delete -n clab-hawkv6-XR-1 -i Gi0-0-0-0
	+-----------+-------+--------+-------------+-------------+
	| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
	+-----------+-------+--------+-------------+-------------+
	| Gi0-0-0-0 | 0ms   | 0s     | 0.00%       |           0 |
	+-----------+-------+--------+-------------+-------------+
	
	clab-mdt-linker will forward the impairments to the specific containerlab command
	More information: https://containerlab.dev/cmd/tools/netem/set/
	`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	checkIsRoot()
}

func checkIsRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("Hawkwing must be run as root")
		os.Exit(1)
	}
}
