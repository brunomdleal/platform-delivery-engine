package preflight

import (
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
)

type NodeRunner struct{}

func (r *NodeRunner) Run(ctx *context.RepoContext, opts Options) ([]Result, error) {
	var results []Result
	imageTag := "order-web-client:preflight-local"

	// 1. npm ci & npm run build
	res := execute("Node Dependencies Installation (npm ci)", ctx.Path, opts.Verbose, "npm", "ci")
	results = append(results, res)
	if res.Status == StatusFailed && opts.FailFast {
		return results, res.Err
	}

	res = execute("Frontend Bundle & Typecheck (npm run build)", ctx.Path, opts.Verbose, "npm", "run", "build")
	results = append(results, res)
	if res.Status == StatusFailed && opts.FailFast {
		return results, res.Err
	}

	// 2. Docker container build
	res = execute("Container Image Build", ctx.Path, opts.Verbose, "docker", "build", "-t", imageTag, ".")
	results = append(results, res)
	if res.Status == StatusFailed && opts.FailFast {
		return results, res.Err
	}

	// 3. Security Audit
	if opts.SkipScan {
		results = append(results, Result{
			Name:   "Container Security Audit (Trivy)",
			Status: StatusSkipped,
		})
	} else {
		res = runTrivyScan(imageTag, ctx.Path, opts.Verbose)
		results = append(results, res)
		if res.Status == StatusFailed && opts.FailFast {
			return results, res.Err
		}
	}

	return results, nil
}
