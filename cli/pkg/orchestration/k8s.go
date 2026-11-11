package orchestration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
)

type KubectlManager struct{}

func NewK8sManager() *KubectlManager {
	return &KubectlManager{}
}

func (k *KubectlManager) Deploy(ctx context.Context, opts K8sOptions) error {
	overlayPath, err := k.resolveOverlayPath(opts.RootDir, opts.Overlay)
	if err != nil {
		return err
	}

	output.Info("Applying Kustomize overlay: %s", overlayPath)
	return runStreamCmd(ctx, opts.RootDir, opts.Verbose, "kubectl", "apply", "-k", overlayPath)
}

func (k *KubectlManager) Down(ctx context.Context, opts K8sOptions) error {
	overlayPath, err := k.resolveOverlayPath(opts.RootDir, opts.Overlay)
	if err != nil {
		return err
	}

	output.Info("Deleting Kustomize overlay resources: %s", overlayPath)
	return runStreamCmd(ctx, opts.RootDir, opts.Verbose, "kubectl", "delete", "-k", overlayPath, "--ignore-not-found=true")
}

func (k *KubectlManager) Status(ctx context.Context, opts K8sOptions) ([]ServiceStatus, error) {
	ns := opts.Namespace
	if ns == "" {
		ns = opts.Overlay
	}
	if ns == "" {
		ns = "dev"
	}

	cmd := exec.CommandContext(ctx, "kubectl", "get", "pods", "-n", ns, "-o", "custom-columns=NAME:.metadata.name,READY:.status.containerStatuses[0].ready,STATUS:.status.phase,RESTARTS:.status.containerStatuses[0].restartCount", "--no-headers")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query pod status in namespace '%s': %w", ns, err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var statuses []ServiceStatus
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			readyStr := "not ready"
			if fields[1] == "true" {
				readyStr = "ready"
			}
			statuses = append(statuses, ServiceStatus{
				Name:    fields[0],
				State:   fmt.Sprintf("%s (%s)", fields[2], readyStr),
				Details: fmt.Sprintf("Restarts: %s", fields[3]),
			})
		}
	}
	return statuses, nil
}

func (k *KubectlManager) resolveOverlayPath(rootDir, overlay string) (string, error) {
	if overlay != "dev" && overlay != "prod" {
		return "", fmt.Errorf("invalid overlay target '%s': must be either 'dev' or 'prod'", overlay)
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		absRoot = rootDir
	}

	target := filepath.Join(absRoot, "k8s", "overlays", overlay)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return "", fmt.Errorf("overlay path does not exist: %s", target)
	}
	return target, nil
}
