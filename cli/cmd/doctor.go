package cmd

import (
	"fmt"
	"os"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/doctor"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

// RunDoctorDiagnostics executes checks and reports status.
func RunDoctorDiagnostics(targetType context.RepoType) bool {
	if targetType != "" && targetType != context.RepoUnknown {
		output.Header(fmt.Sprintf("Workstation Diagnostics for Context: %s", targetType))
	} else {
		output.Header("Running Complete Workstation Toolchain Diagnostics")
	}

	var results []doctor.CheckResult
	if targetType != "" && targetType != context.RepoUnknown {
		results = doctor.RunForContext(targetType)
	} else {
		results = doctor.RunAll()
	}

	hasFailures := false
	for _, r := range results {
		switch r.Status {
		case doctor.StatusOK:
			if r.AutoHealed {
				output.Success("%-28s: %s [RECOVERED]", r.Name, r.Details)
			} else {
				output.Success("%-28s: %s", r.Name, r.Details)
			}
		case doctor.StatusWarn:
			output.Warn("%-28s: %s", r.Name, r.Details)
			if r.Remediation != "" {
				fmt.Printf("   -> %s\n", r.Remediation)
			}
		case doctor.StatusFail:
			hasFailures = true
			output.Error("%-28s: %s", r.Name, r.Details)
			if r.Remediation != "" {
				fmt.Printf("   -> Hint: %s\n", r.Remediation)
			}
		}
	}

	fmt.Println()
	if hasFailures {
		output.Error("Diagnostic audit failed. Resolve missing required toolchains above.")
		return false
	}

	output.Success("Workstation prerequisites satisfied!")
	return true
}

var doctorCmd = &cobra.Command{
	Use:   "doctor [repoType]",
	Short: "Audit host prerequisites globally or for a specific repository context",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var target context.RepoType
		if len(args) > 0 {
			target = context.RepoType(args[0])
		}

		if !RunDoctorDiagnostics(target) {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
