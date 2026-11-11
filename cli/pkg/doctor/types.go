package doctor

// Status represents the health check result state
type Status string

const (
	StatusOK   Status = "OK"
	StatusWarn Status = "WARN"
	StatusFail Status = "FAIL"
)

// CheckResult holds the outcome of a single toolchain check
type CheckResult struct {
	Name        string
	Status      Status
	Details     string
	Remediation string
	AutoHealed  bool
}

// Checker defines the interface for running a diagnostic check
type Checker interface {
	Run() CheckResult
}

// Healer defines an optional self-healing contract for a diagnostic check
type Healer interface {
	Checker
	// Heal attempts a one-time automated repair (start daemon, pull image, start cluster, etc.)
	Heal() error
}
