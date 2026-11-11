package preflight

import (
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
)

type QARunner struct{}

func (r *QARunner) Run(ctx *context.RepoContext, opts Options) ([]Result, error) {
	var results []Result
	res := execute("Test Suite Compilation", ctx.Path, opts.Verbose, "./mvnw", "test-compile")
	results = append(results, res)
	return results, nil
}
