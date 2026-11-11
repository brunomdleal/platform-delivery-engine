package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/orchestration"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	cleanVolumes bool
	composePath  string
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the full local Docker Compose stack in detached mode",
	Long:  `Boots order-service, order-web-client, and postgres containers with health checks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Header("Starting Local Delivery Stack (Docker Compose)")

		mgr := orchestration.NewComposeManager()
		verbose, _ := cmd.Flags().GetBool("verbose")

		targetDir := resolveRootDir()

		opts := orchestration.ComposeOptions{
			RootDir:  targetDir,
			Services: []orchestration.TargetService{orchestration.TargetAll},
			Detach:   true,
			Verbose:  verbose,
		}

		if err := mgr.Up(context.Background(), opts); err != nil {
			output.Error("Failed to start Compose services: %v", err)
			return err
		}

		output.Success("Full stack started successfully.")
		printComposeStatus(mgr, targetDir)
		return nil
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and tear down the local Docker Compose stack",
	Long:  `Stops running application and database containers. Use --clean to purge persistent database volumes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Header("Tearing Down Local Delivery Stack")

		mgr := orchestration.NewComposeManager()
		verbose, _ := cmd.Flags().GetBool("verbose")

		targetDir := resolveRootDir()

		opts := orchestration.ComposeOptions{
			RootDir:  targetDir,
			Services: []orchestration.TargetService{orchestration.TargetAll},
			Clean:    cleanVolumes,
			Verbose:  verbose,
		}

		if err := mgr.Down(context.Background(), opts); err != nil {
			output.Error("Failed to tear down Compose services: %v", err)
			return err
		}

		if cleanVolumes {
			output.Success("Compose stack and persistent volumes removed.")
		} else {
			output.Success("Compose stack stopped. Volumes preserved.")
		}
		return nil
	},
}

func resolveRootDir() string {
	if composePath != "" {
		return filepath.Clean(composePath)
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func printComposeStatus(mgr orchestration.ComposeManager, rootDir string) {
	statuses, err := mgr.Status(context.Background(), orchestration.ComposeOptions{RootDir: rootDir})
	if err != nil || len(statuses) == 0 {
		return
	}

	output.Header("Active Service Status")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tSTATUS\tPORTS")
	for _, s := range statuses {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.Name, s.State, s.Ports)
	}
	w.Flush()
}

func init() {
	upCmd.Flags().StringVarP(&composePath, "root", "r", "", "Root directory of platform-delivery-engine (defaults to CWD)")
	downCmd.Flags().StringVarP(&composePath, "root", "r", "", "Root directory of platform-delivery-engine (defaults to CWD)")
	downCmd.Flags().BoolVarP(&cleanVolumes, "clean", "c", false, "Purge persistent volumes (docker compose down -v)")

	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}
