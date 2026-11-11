package preflight

import (
	"time"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
)

// StageStatus represents the execution state of an individual gate
type StageStatus string

const (
	StatusPassed  StageStatus = "PASSED"
	StatusFailed  StageStatus = "FAILED"
	StatusSkipped StageStatus = "SKIPPED"
)

// Result captures execution metrics and diagnostics for one pipeline stage
type Result struct {
	Name     string
	Status   StageStatus
	Duration time.Duration
	Err      error
}

// Options controls pre-flight behavior via CLI flags
type Options struct {
	SkipScan bool
	Verbose  bool
	FailFast bool
}

// Runner defines the interface for repository-specific CI mirrors
type Runner interface {
	Run(ctx *context.RepoContext, opts Options) ([]Result, error)
}
