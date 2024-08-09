package cmd

import (
	"fmt"

	"github.com/hawkv6/clab-telemetry-linker/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("clab-telemetry-linker version: %s\n", version.GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
