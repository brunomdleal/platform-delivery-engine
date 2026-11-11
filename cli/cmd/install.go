package cmd

import (
	"path/filepath"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/install"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install pde binary to ~/.local/bin for global system access",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Header("Installing PDE CLI Globally")

		destPath, err := install.InstallBinary()
		if err != nil {
			output.Error("Installation failed: %v", err)
			return err
		}

		output.Success("Binary installed to: %s", destPath)

		installDir := filepath.Dir(destPath)
		if !install.CheckPathInEnvironment(installDir) {
			output.Warn("Directory %s is not in your $PATH!", installDir)
			output.Info("Add it to your shell configuration (e.g. fish: 'fish_add_path %s' or bash: 'export PATH=\"$PATH:%s\"')", installDir, installDir)
		} else {
			output.Success("Directory %s is active in $PATH.", installDir)
		}

		output.Success("You can now run 'pde' directly from any directory!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
