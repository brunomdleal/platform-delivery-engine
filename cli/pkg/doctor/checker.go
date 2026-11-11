package doctor

import (
	"time"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
)

// RunAll executes the entire workstation diagnostic suite
func RunAll() []CheckResult {
	return runCheckers([]Checker{
		&DockerSocketChecker{},
		&DockerComposeChecker{},
		&JavaChecker{MinMajor: 21},
		&MavenChecker{},
		&NodeChecker{MinMajor: 20},
		&KubectlChecker{},
		&MinikubeChecker{Timeout: 3 * time.Second},
		&KubernetesClusterChecker{Timeout: 3 * time.Second},
		&TerraformChecker{},
		&GitChecker{},
		&TrivyChecker{},
	})
}

// RunForContext runs only the prerequisite checks relevant to a given repository
func RunForContext(repoType context.RepoType) []CheckResult {
	var checkers []Checker

	// Baseline common dependencies
	checkers = append(checkers, &DockerSocketChecker{}, &GitChecker{})

	switch repoType {
	case context.RepoOrderService:
		checkers = append(checkers,
			&JavaChecker{MinMajor: 21},
			&MavenChecker{},
			&TrivyChecker{},
		)

	case context.RepoOrderWebClient:
		checkers = append(checkers,
			&NodeChecker{MinMajor: 20},
			&TrivyChecker{},
		)

	case context.RepoPlatformDeliveryEngine:
		checkers = append(checkers,
			&DockerComposeChecker{},
			&KubectlChecker{},
			&MinikubeChecker{Timeout: 3 * time.Second},
			&KubernetesClusterChecker{Timeout: 3 * time.Second},
			&TerraformChecker{},
		)

	case context.RepoQAAutomation:
		checkers = append(checkers,
			&JavaChecker{MinMajor: 21},
			&MavenChecker{},
		)

	default:
		return RunAll()
	}

	return runCheckers(checkers)
}

func runCheckers(checkers []Checker) []CheckResult {
	results := make([]CheckResult, 0, len(checkers))

	for _, c := range checkers {
		res := c.Run()

		// If failed or warned, attempt 1-time self-healing if implemented
		if res.Status != StatusOK {
			if healer, ok := c.(Healer); ok {
				output.Info("Self-healing [%s]: attempting 1-time automated remediation...", res.Name)
				if err := healer.Heal(); err == nil {
					// Settle time for daemon/process initialization
					time.Sleep(1500 * time.Millisecond)

					// Re-evaluate check
					retryRes := c.Run()
					if retryRes.Status == StatusOK {
						retryRes.AutoHealed = true
						retryRes.Details = retryRes.Details + " (auto-healed)"
						results = append(results, retryRes)
						continue
					}
				}
			}
		}

		results = append(results, res)
	}

	return results
}
