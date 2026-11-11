package orchestration

import (
	"context"
	"time"
)

// TargetService identifies specific containers or infrastructure components
type TargetService string

const (
	TargetAll      TargetService = "all"
	TargetPostgres TargetService = "postgres"
	TargetOrderApp TargetService = "order-service"
	TargetOrderWeb TargetService = "order-web-client"
)

// ComposeOptions configures Docker Compose operations
type ComposeOptions struct {
	RootDir  string
	Services []TargetService
	Detach   bool
	Clean    bool // down -v (purge volumes)
	Verbose  bool
	Timeout  time.Duration
}

// K8sOptions configures Kustomize and kubectl overlay operations
type K8sOptions struct {
	RootDir   string
	Overlay   string // "dev" or "prod"
	Namespace string
	Verbose   bool
	Timeout   time.Duration
}

// ServiceStatus captures operational health metrics for running workloads
type ServiceStatus struct {
	Name    string
	State   string // running, stopped, healthy, unhealthy
	Ports   string
	Details string
}

// ComposeManager orchestrates Docker Compose lifecycle operations
type ComposeManager interface {
	Up(ctx context.Context, opts ComposeOptions) error
	Down(ctx context.Context, opts ComposeOptions) error
	Status(ctx context.Context, opts ComposeOptions) ([]ServiceStatus, error)
}

// K8sManager orchestrates declarative Kubernetes overlays
type K8sManager interface {
	Deploy(ctx context.Context, opts K8sOptions) error
	Down(ctx context.Context, opts K8sOptions) error
	Status(ctx context.Context, opts K8sOptions) ([]ServiceStatus, error)
}
