package cmd

import (
	"os"
	"path/filepath"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/install"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	targetHookPath string
	installAll     bool
)

var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Manage Git automation hooks across repositories",
}

var hooksInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Inject pre-push verification hook running 'pde pre-flight'",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Header("Installing Git Pre-Push Hook")

		if installAll {
			return installAcrossSiblings()
		}

		target := targetHookPath
		if target == "" {
			wd, err := os.Getwd()
			if err != nil {
				output.Error("Failed to resolve current directory: %v", err)
				return err
			}
			target = wd
		}

		hookPath, err := install.InstallHook(target)
		if err != nil {
			output.Error("Failed to install hook: %v", err)
			return err
		}

		output.Success("Pre-push hook successfully installed at: %s", hookPath)
		return nil
	},
}

func installAcrossSiblings() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Assuming standard portfolio layout where repos share a common parent
	parentDir := filepath.Dir(wd)
	if filepath.Base(wd) == "cli" {
		parentDir = filepath.Dir(filepath.Dir(wd))
	}

	siblingRepos := []string{
		"order-service",
		"order-web-client",
		"qa-automation",
		"platform-delivery-engine",
	}

	for _, repo := range siblingRepos {
		repoPath := filepath.Join(parentDir, repo)
		hookPath, err := install.InstallHook(repoPath)
		if err != nil {
			output.Warn("Skipping %s: %v", repo, err)
			continue
		}
		output.Success("Installed hook in %s -> %s", repo, hookPath)
	}

	return nil
}

func init() {
	hooksInstallCmd.Flags().StringVarP(&targetHookPath, "path", "p", "", "Target repository path (defaults to current directory)")
	hooksInstallCmd.Flags().BoolVar(&installAll, "all", false, "Install pre-push hook across all sibling portfolio repositories")

	hooksCmd.AddCommand(hooksInstallCmd)
	rootCmd.AddCommand(hooksCmd)
}
