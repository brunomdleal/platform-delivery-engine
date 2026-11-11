package orchestration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
)

type DockerComposeManager struct{}

func NewComposeManager() *DockerComposeManager {
	return &DockerComposeManager{}
}

func (m *DockerComposeManager) Up(ctx context.Context, opts ComposeOptions) error {
	args := []string{"compose", "up"}
	if opts.Detach {
		args = append(args, "-d")
	}

	for _, svc := range opts.Services {
		if svc != TargetAll {
			args = append(args, string(svc))
		}
	}

	output.Info("Running: docker %s", strings.Join(args, " "))
	return runStreamCmd(ctx, opts.RootDir, opts.Verbose, "docker", args...)
}

func (m *DockerComposeManager) Down(ctx context.Context, opts ComposeOptions) error {
	args := []string{"compose", "down"}
	if opts.Clean {
		args = append(args, "-v")
	}

	output.Info("Running: docker %s", strings.Join(args, " "))
	return runStreamCmd(ctx, opts.RootDir, opts.Verbose, "docker", args...)
}

func (m *DockerComposeManager) Status(ctx context.Context, opts ComposeOptions) ([]ServiceStatus, error) {
	cmd := exec.CommandContext(ctx, "docker", "compose", "ps", "--format", "{{.Name}}|{{.State}}|{{.Health}}|{{.Ports}}")
	cmd.Dir = opts.RootDir

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to inspect compose services: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var statuses []ServiceStatus
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "|")
		status := ServiceStatus{Name: parts[0]}
		if len(parts) > 1 {
			status.State = parts[1]
		}
		if len(parts) > 2 && parts[2] != "" {
			status.State = fmt.Sprintf("%s (%s)", status.State, parts[2])
		}
		if len(parts) > 3 {
			status.Ports = parts[3]
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func runStreamCmd(ctx context.Context, dir string, verbose bool, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = dir

	var errBuf bytes.Buffer
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = &errBuf
	}

	if err := cmd.Run(); err != nil {
		if !verbose && errBuf.Len() > 0 {
			fmt.Fprintf(os.Stderr, "\n%s\n", errBuf.String())
		}
		return fmt.Errorf("command %s %v failed: %w", command, args, err)
	}
	return nil
}
