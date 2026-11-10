package cmd

import (
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows tool version",
	Run: func(cmd *cobra.Command, args []string) {
		output.Version(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
