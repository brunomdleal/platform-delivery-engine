package doctor

import (
	"testing"
	"time"
)

func TestRunAllDoesNotPanic(t *testing.T) {
	results := RunAll()
	if len(results) == 0 {
		t.Fatalf("expected non-empty diagnostic results")
	}

	for _, r := range results {
		if r.Name == "" {
			t.Errorf("encountered check result with empty Name")
		}
		if r.Status != StatusOK && r.Status != StatusWarn && r.Status != StatusFail {
			t.Errorf("unexpected status %v for check %s", r.Status, r.Name)
		}
	}
}

func TestKubernetesClusterCheckerTimeout(t *testing.T) {
	checker := &KubernetesClusterChecker{Timeout: 100 * time.Millisecond}
	res := checker.Run()
	if res.Name != "Kubernetes Cluster" {
		t.Errorf("expected name 'Kubernetes Cluster', got %s", res.Name)
	}
}
