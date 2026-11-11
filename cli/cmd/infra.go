package cmd

import (
	"context"
	"path/filepath"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/orchestration"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var infraClean bool

var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Manage shared backing infrastructure dependencies (PostgreSQL)",
	Long:  `Boots or destroys backing infrastructure containers in detached mode for hybrid development.`,
}

var infraUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start containerized PostgreSQL infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Header("Starting Local Backing Infrastructure")
		mgr := orchestration.NewComposeManager()
		root, _ := cmd.Flags().GetString("root")
		verbose, _ := cmd.Flags().GetBool("verbose")

		opts := orchestration.ComposeOptions{
			RootDir:  filepath.Clean(root),
			Services: []orchestration.TargetService{orchestration.TargetPostgres},
			Detach:   true,
			Verbose:  verbose,
		}

		if err := mgr.Up(context.Background(), opts); err != nil {
			output.Error("Failed to start PostgreSQL: %v", err)
			return err
		}
		output.Success("PostgreSQL is running and accessible on port 5432.")
		return nil
	},
}

var infraDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop containerized PostgreSQL infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Header("Stopping Local Backing Infrastructure")
		mgr := orchestration.NewComposeManager()
		root, _ := cmd.Flags().GetString("root")
		verbose, _ := cmd.Flags().GetBool("verbose")

		opts := orchestration.ComposeOptions{
			RootDir:  filepath.Clean(root),
			Services: []orchestration.TargetService{orchestration.TargetPostgres},
			Clean:    infraClean,
			Verbose:  verbose,
		}

		if err := mgr.Down(context.Background(), opts); err != nil {
			output.Error("Failed to stop PostgreSQL: %v", err)
			return err
		}
		output.Success("PostgreSQL stopped successfully.")
		return nil
	},
}

func init() {
	infraDownCmd.Flags().BoolVarP(&infraClean, "clean", "c", false, "Remove persistent storage volumes")
	infraCmd.AddCommand(infraUpCmd)
	infraCmd.AddCommand(infraDownCmd)
	rootCmd.AddCommand(infraCmd)
}
