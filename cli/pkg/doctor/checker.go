package doctor

import "time"

// RunAll executes all standard prerequisite checks in sequence
func RunAll() []CheckResult {
	checkers := []Checker{
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
	}

	results := make([]CheckResult, 0, len(checkers))
	for _, c := range checkers {
		results = append(results, c.Run())
	}
	return results
}
