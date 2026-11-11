package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "pde",
	Short:   "pde - Developer Experience (DevEx) & Platform Delivery CLI",
	Long:    "A unified platform engineering CLI tool to inspect repository context, validate toolchains, and mirror CI quality gates locally.",
	Version: "0.0.1",
}

// Execute runs the root command entrypoint
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
