package cmd

import (
	"os"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var targetPath string

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Detect and display the active repository context",
	Run: func(cmd *cobra.Command, args []string) {
		path := targetPath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				output.Error("Failed to get current directory: %v", err)
				os.Exit(1)
			}
		}

		ctx, err := context.Detect(path)
		if err != nil {
			output.Error("Error detecting context: %v", err)
			os.Exit(1)
		}

		output.Header("Repository Context Diagnostic")
		output.Info("Evaluated Path: %s", ctx.Type)

		if ctx.Type == context.RepoUnknown {
			output.Warn("Detected Type: %s", ctx.Type)
			output.Warn("No known microservice or platform markers found in this directory.")
			return
		}

		output.Success("Detected Type : %s", ctx.Type)
		output.Success("Description   : %s", ctx.Description)
		output.Success("Toolchain     : %s", ctx.Toolchain)
	},
}

func init() {
	contextCmd.Flags().StringVarP(
		&targetPath,
		"path",
		"p",
		"",
		"Path to evaluate (defaults to current working directory)",
	)
	rootCmd.AddCommand(contextCmd)
}
