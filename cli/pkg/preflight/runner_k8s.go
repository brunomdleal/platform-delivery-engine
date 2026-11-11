package preflight

import (
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
)

type K8sPlatformRunner struct{}

func (r *K8sPlatformRunner) Run(ctx *context.RepoContext, opts Options) ([]Result, error) {
	var results []Result

	// 1. Kustomize overlays dev rendering
	res := execute("Kustomize Template (dev)", ctx.Path, opts.Verbose, "kubectl", "kustomize", "k8s/overlays/dev")
	results = append(results, res)
	if res.Status == StatusFailed && opts.FailFast {
		return results, res.Err
	}

	// 2. Kustomize overlays prod rendering
	res = execute("Kustomize Template (prod)", ctx.Path, opts.Verbose, "kubectl", "kustomize", "k8s/overlays/prod")
	results = append(results, res)
	if res.Status == StatusFailed && opts.FailFast {
		return results, res.Err
	}

	// 3. Client dry-run dev
	res = execute("Kubectl Client Dry-Run (dev)", ctx.Path, opts.Verbose, "kubectl", "apply", "--dry-run=client", "-k", "k8s/overlays/dev")
	results = append(results, res)
	if res.Status == StatusFailed && opts.FailFast {
		return results, res.Err
	}

	return results, nil
}
