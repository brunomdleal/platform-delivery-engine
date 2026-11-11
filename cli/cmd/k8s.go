package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/orchestration"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	k8sNamespace string
	k8sRootDir   string
)

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Manage declarative Kubernetes overlays (Minikube)",
	Long:  `Deploy, inspect, and tear down dev and prod Kustomize overlays against the local Minikube cluster.`,
}

var k8sDeployCmd = &cobra.Command{
	Use:   "deploy [dev|prod]",
	Short: "Deploy an environment overlay using Kustomize",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overlay := args[0]
		output.Header(fmt.Sprintf("Deploying Kubernetes Overlay: %s", overlay))

		mgr := orchestration.NewK8sManager()
		verbose, _ := cmd.Flags().GetBool("verbose")

		opts := orchestration.K8sOptions{
			RootDir: resolveK8sRootDir(),
			Overlay: overlay,
			Verbose: verbose,
		}

		if err := mgr.Deploy(context.Background(), opts); err != nil {
			output.Error("Deployment failed: %v", err)
			return err
		}

		output.Success("Overlay '%s' successfully applied.", overlay)
		printK8sStatus(mgr, overlay, overlay)
		return nil
	},
}

var k8sStatusCmd = &cobra.Command{
	Use:   "status [dev|prod]",
	Short: "Check running pod status across overlay namespaces",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetEnv := "dev"
		if len(args) == 1 {
			targetEnv = args[0]
		}
		if k8sNamespace != "" {
			targetEnv = k8sNamespace
		}

		mgr := orchestration.NewK8sManager()
		printK8sStatus(mgr, targetEnv, targetEnv)
		return nil
	},
}

var k8sDownCmd = &cobra.Command{
	Use:   "down [dev|prod]",
	Short: "Delete overlay resources from the cluster",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overlay := args[0]
		output.Header(fmt.Sprintf("Tearing Down Kubernetes Overlay: %s", overlay))

		mgr := orchestration.NewK8sManager()
		verbose, _ := cmd.Flags().GetBool("verbose")

		opts := orchestration.K8sOptions{
			RootDir: resolveK8sRootDir(),
			Overlay: overlay,
			Verbose: verbose,
		}

		if err := mgr.Down(context.Background(), opts); err != nil {
			output.Error("Failed to delete overlay: %v", err)
			return err
		}

		output.Success("Overlay '%s' resources deleted.", overlay)
		return nil
	},
}

func resolveK8sRootDir() string {
	if k8sRootDir != "" {
		return k8sRootDir
	}
	return resolveRootDir()
}

func printK8sStatus(mgr orchestration.K8sManager, overlay, namespace string) {
	statuses, err := mgr.Status(context.Background(), orchestration.K8sOptions{
		Overlay:   overlay,
		Namespace: namespace,
	})
	if err != nil {
		output.Warn("Could not query pods in namespace '%s': %v", namespace, err)
		return
	}

	if len(statuses) == 0 {
		output.Info("No active pods found in namespace '%s'.", namespace)
		return
	}

	output.Header(fmt.Sprintf("Workloads in Namespace: %s", namespace))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "POD\tSTATUS\tDETAILS")
	for _, s := range statuses {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.Name, s.State, s.Details)
	}
	w.Flush()
}

func init() {
	k8sCmd.PersistentFlags().StringVarP(&k8sRootDir, "root", "r", "", "Root directory of platform-delivery-engine")
	k8sStatusCmd.Flags().StringVarP(&k8sNamespace, "namespace", "n", "", "Explicit namespace to query")

	k8sCmd.AddCommand(k8sDeployCmd)
	k8sCmd.AddCommand(k8sStatusCmd)
	k8sCmd.AddCommand(k8sDownCmd)

	rootCmd.AddCommand(k8sCmd)
}
