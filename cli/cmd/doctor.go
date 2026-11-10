package cmd

import (
	"fmt"
	"os"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/doctor"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Audit host prerequisites, sockets, and required developer toolchains",
	Run: func(cmd *cobra.Command, args []string) {
		output.Header("Running PDE Environment & Toolchain Diagnostics")

		results := doctor.RunAll()
		hasFailures := false

		for _, r := range results {
			switch r.Status {
			case doctor.StatusOK:
				output.Success(fmt.Sprintf("%-26s: %s", r.Name, r.Details))
			case doctor.StatusWarn:
				output.Warn(fmt.Sprintf("%-26s: %s", r.Name, r.Details))
				if r.Remediation != "" {
					fmt.Printf("   -> %s\n", r.Remediation)
				}
			case doctor.StatusFail:
				hasFailures = true
				output.Error(fmt.Sprintf("%-26s: %s", r.Name, r.Details))
				if r.Remediation != "" {
					fmt.Printf("   -> Hint: %s\n", r.Remediation)
				}
			}
		}

		fmt.Println()
		if hasFailures {
			output.Error("Diagnostic checks failed. Resolve missing required toolchains above.")
			os.Exit(1)
		}

		output.Success("All required workstation prerequisites satisfied!")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
