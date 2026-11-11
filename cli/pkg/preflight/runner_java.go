package preflight

import (
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
)

type JavaRunner struct{}

func (r *JavaRunner) Run(ctx *context.RepoContext, opts Options) ([]Result, error) {
	var results []Result
	imageTag := "order-service:preflight-local"

	// 1. Maven clean verify
	res := execute("Maven Verification (Unit + Integration Tests)", ctx.Path, opts.Verbose, "./mvnw", "clean", "verify")
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
