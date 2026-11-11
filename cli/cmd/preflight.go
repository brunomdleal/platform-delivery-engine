package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/preflight"
	"github.com/spf13/cobra"
)

var (
	preflightOpts preflight.Options
	skipDoctor    bool
	preflightPath string
)

var preflightCmd = &cobra.Command{
	Use:   "pre-flight",
	Short: "Execute local pre-flight quality and security gates mirroring CI pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		// 1. Resolve Target Evaluation Path
		evalPath := preflightPath
		if evalPath == "" {
			wd, err := os.Getwd()
			if err != nil {
				output.Error("Failed to resolve current directory: %v", err)
				os.Exit(1)
			}
			evalPath = wd
		}

		// 2. Detect Context (same engine used by 'pde context')
		ctx, err := context.Detect(evalPath)

		if err != nil || ctx.Type == context.RepoUnknown {
			output.Error("Pre-flight cannot run outside a known repository context (evaluated path: %s).", evalPath)
			os.Exit(1)
		}
		output.Header(fmt.Sprintf("Pre-Flight Conveyor Belt: %s", ctx.Description))
		output.Info("Repository Path : %s", ctx.Path)
		output.Info("Detected Context: %s", ctx.Type)

		// 3. Run Context-Targeted Doctor (equivalent to 'pde doctor <repoType>')
		if !skipDoctor {
			if !RunDoctorDiagnostics(ctx.Type) {
				output.Error("Pre-flight halted: Prerequisite checks failed for %s.", ctx.Type)
				os.Exit(1)
			}
		}

		// 4. Run Pipeline Stages
		runner, err := preflight.GetRunner(ctx.Type)
		if err != nil {
			output.Error("%v", err)
			os.Exit(1)
		}

		results, _ := runner.Run(ctx, preflightOpts)

		hasFailures := false
		output.Header("Pre-Flight Verification Summary")

		for _, r := range results {
			switch r.Status {
			case preflight.StatusPassed:
				output.Success("%-52s [%s] (%s)", r.Name, r.Status, r.Duration.Round(time.Millisecond))
			case preflight.StatusSkipped:
				output.Warn("%-52s [%s]", r.Name, r.Status)
			case preflight.StatusFailed:
				hasFailures = true
				output.Error("%-52s [%s] (%s)", r.Name, r.Status, r.Duration.Round(time.Millisecond))
				if r.Err != nil {
					fmt.Printf("   -> Reason: %v\n", r.Err)
				}
			}
		}

		totalElapsed := time.Since(start).Round(time.Millisecond)
		fmt.Println()
		if hasFailures {
			output.Error("Pre-flight verification failed after %s. Resolve errors before pushing.", totalElapsed)
			os.Exit(1)
		}

		output.Success("All local CI mirror gates passed in %s! Ready for git push.", totalElapsed)
	},
}

func init() {
	preflightCmd.Flags().StringVarP(&preflightPath, "path", "p", "", "Target repository path to evaluate (defaults to current working directory)")
	preflightCmd.Flags().BoolVar(&skipDoctor, "skip-doctor", false, "Bypass doctor prerequisite checks")
	preflightCmd.Flags().BoolVar(&preflightOpts.SkipScan, "skip-scan", false, "Skip Trivy vulnerability container scans")
	preflightCmd.Flags().BoolVarP(&preflightOpts.Verbose, "verbose", "v", false, "Stream all stdout/stderr from underlying tools")
	preflightCmd.Flags().BoolVar(&preflightOpts.FailFast, "fail-fast", true, "Halt execution on the first stage failure")
	rootCmd.AddCommand(preflightCmd)
}
