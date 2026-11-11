package doctor

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
)

// DockerSocketChecker verifies /var/run/docker.sock connectivity
type DockerSocketChecker struct{}

func (c *DockerSocketChecker) Run() CheckResult {
	socketPath := "/var/run/docker.sock"
	conn, err := net.DialTimeout("unix", socketPath, 1*time.Second)
	if err != nil {
		output.Warn("Docker is not instantiated. Running command: sudo systemctl start docker")
		exec.Command("sudo systemctl", "start", "docker")
		conn2, err2 := net.DialTimeout("unix", socketPath, 5*time.Second)
		if err2 != nil {
			return CheckResult{
				Name:        "Docker Daemon",
				Status:      StatusFail,
				Details:     fmt.Sprintf("Cannot connect to %s", socketPath),
				Remediation: "Ensure Docker is installed and running.  (`sudo systemctl start docker`)",
			}
		}
		_ = conn2.Close()
	}
	_ = conn.Close()

	return CheckResult{
		Name:    "Docker Daemon",
		Status:  StatusOK,
		Details: fmt.Sprintf("Active & listening on %s", socketPath),
	}
}

func (c *DockerSocketChecker) Heal() error {
	cmd := exec.Command("sudo", "systemctl", "start", "docker")
	if err := cmd.Run(); err != nil {
		// Fallback for non-systemd or alternative WSL inits
		return exec.Command("sudo", "service", "docker", "start").Run()
	}
	return nil
}

// DockerComposeChecker checks docker compose plugin/binary
type DockerComposeChecker struct{}

func (c *DockerComposeChecker) Run() CheckResult {
	cmd := exec.Command("docker", "compose", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Docker Compose",
			Status:      StatusFail,
			Details:     "docker compose command failed or plugin not installed",
			Remediation: "Install docker-compose plugin (`sudo pacman -S docker-compose`)",
		}
	}

	raw := strings.TrimSpace(string(out))
	return CheckResult{
		Name:    "Docker Compose",
		Status:  StatusOK,
		Details: raw,
	}
}

// JavaChecker verifies Java version >= MinMajor
type JavaChecker struct {
	MinMajor int
}

func (c *JavaChecker) Run() CheckResult {
	cmd := exec.Command("java", "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Java Toolchain",
			Status:      StatusFail,
			Details:     "java binary not found in PATH",
			Remediation: fmt.Sprintf("Install JDK %d+ (e.g. Temurin or Arch jdk-openjdk)", c.MinMajor),
		}
	}

	outputStr := string(out)
	re := regexp.MustCompile(`version "(?P<major>\d+)`)
	matches := re.FindStringSubmatch(outputStr)
	if len(matches) < 2 {
		return CheckResult{
			Name:        "Java Toolchain",
			Status:      StatusWarn,
			Details:     "Unable to parse Java version output",
			Remediation: fmt.Sprintf("Verify installed JDK version manually satisfies >= %d", c.MinMajor),
		}
	}

	major, _ := strconv.Atoi(matches[1])
	if major < c.MinMajor {
		return CheckResult{
			Name:        "Java Toolchain",
			Status:      StatusFail,
			Details:     fmt.Sprintf("Detected Java version %d (< required %d)", major, c.MinMajor),
			Remediation: fmt.Sprintf("Upgrade to JDK %d+", c.MinMajor),
		}
	}

	firstLine := strings.Split(outputStr, "\n")[0]
	return CheckResult{
		Name:    "Java Toolchain",
		Status:  StatusOK,
		Details: strings.TrimSpace(firstLine),
	}
}

// MavenChecker verifies optional host maven installation
type MavenChecker struct{}

func (c *MavenChecker) Run() CheckResult {
	cmd := exec.Command("mvn", "-v")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Apache Maven (Optional)",
			Status:      StatusWarn,
			Details:     "Host mvn binary not found (repos use bundled ./mvnw)",
			Remediation: "Optional: Install maven (`sudo pacman -S maven`) for global usage",
		}
	}

	firstLine := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	return CheckResult{
		Name:    "Apache Maven (Optional)",
		Status:  StatusOK,
		Details: strings.TrimSpace(firstLine),
	}
}

// NodeChecker verifies Node.js version >= MinMajor
type NodeChecker struct {
	MinMajor int
}

func (c *NodeChecker) Run() CheckResult {
	cmd := exec.Command("node", "-v")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Node.js Runtime",
			Status:      StatusFail,
			Details:     "node binary not found in PATH",
			Remediation: fmt.Sprintf("Install Node.js >= %d (`sudo pacman -S nodejs npm`)", c.MinMajor),
		}
	}

	raw := strings.TrimSpace(string(out))
	clean := strings.TrimPrefix(raw, "v")
	parts := strings.Split(clean, ".")
	if len(parts) == 0 {
		return CheckResult{
			Name:        "Node.js Runtime",
			Status:      StatusWarn,
			Details:     fmt.Sprintf("Unrecognized format: %s", raw),
			Remediation: fmt.Sprintf("Verify Node.js version >= %d", c.MinMajor),
		}
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return CheckResult{
			Name:        "Node.js Runtime",
			Status:      StatusWarn,
			Details:     fmt.Sprintf("Unrecognized format: %s", raw),
			Remediation: fmt.Sprintf("Verify Node.js version >= %d", c.MinMajor),
		}
	}

	if major < c.MinMajor {
		return CheckResult{
			Name:        "Node.js Runtime",
			Status:      StatusFail,
			Details:     fmt.Sprintf("Detected Node %s (< required v%d)", raw, c.MinMajor),
			Remediation: fmt.Sprintf("Upgrade Node.js to >= %d", c.MinMajor),
		}
	}

	return CheckResult{
		Name:    "Node.js Runtime",
		Status:  StatusOK,
		Details: raw,
	}
}

// KubectlChecker verifies kubectl is installed and reports client version
type KubectlChecker struct{}

func (c *KubectlChecker) Run() CheckResult {
	cmd := exec.Command("kubectl", "version", "--client", "--output=yaml")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Kubernetes CLI (kubectl)",
			Status:      StatusFail,
			Details:     "kubectl binary not found in PATH",
			Remediation: "Install kubectl (`sudo pacman -S kubectl`)",
		}
	}

	re := regexp.MustCompile(`gitVersion: (v\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(out))
	version := "Installed"
	if len(matches) >= 2 {
		version = matches[1]
	}

	return CheckResult{
		Name:    "Kubernetes CLI (kubectl)",
		Status:  StatusOK,
		Details: version,
	}
}

// MinikubeChecker verifies minikube CLI presence and runtime profile state
type MinikubeChecker struct {
	Timeout time.Duration
}

func (c *MinikubeChecker) Run() CheckResult {
	cmdVersion := exec.Command("minikube", "version", "--short")
	outVersion, err := cmdVersion.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Minikube CLI",
			Status:      StatusFail,
			Details:     "minikube binary not found in PATH",
			Remediation: "Install minikube (`sudo pacman -S minikube`)",
		}
	}

	version := strings.TrimSpace(string(outVersion))

	timeout := c.Timeout
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmdStatus := exec.CommandContext(ctx, "minikube", "status", "--format={{.Host}}")
	outStatus, err := cmdStatus.CombinedOutput()
	statusText := strings.TrimSpace(string(outStatus))

	if err != nil || statusText == "" || strings.EqualFold(statusText, "Stopped") {
		return CheckResult{
			Name:        "Minikube CLI",
			Status:      StatusWarn,
			Details:     fmt.Sprintf("%s (profile stopped/not running)", version),
			Remediation: "Start your cluster with `minikube start --driver=docker --ports=8443:8443`",
		}
	}

	return CheckResult{
		Name:    "Minikube CLI",
		Status:  StatusOK,
		Details: fmt.Sprintf("%s (cluster host: %s)", version, statusText),
	}
}

func (c *MinikubeChecker) Heal() error {
	cmd := exec.Command("minikube", "start", "--driver=docker", "--ports=8443:8443")
	return cmd.Run()
}

// KubernetesClusterChecker verifies connectivity to the cluster control plane and node readiness
type KubernetesClusterChecker struct {
	Timeout time.Duration
}

func (c *KubernetesClusterChecker) Run() CheckResult {
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 3 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmdCluster := exec.CommandContext(ctx, "kubectl", "cluster-info")
	out, err := cmdCluster.CombinedOutput()
	outLower := strings.ToLower(string(out))
	if err != nil {
		return CheckResult{
			Name:        "Kubernetes Cluster",
			Status:      StatusWarn,
			Details:     "Control plane unreachable or no active cluster context",
			Remediation: "Start Minikube (`minikube start`) or ensure ~/.kube/config points to an active cluster",
		}
	} else if !strings.Contains(outLower, "control plane") && !strings.Contains(outLower, "master") {
		return CheckResult{
			Name:        "Kubernetes Cluster",
			Status:      StatusWarn,
			Details:     "Unexpected response from cluster-info",
			Remediation: "Check cluster status with `kubectl cluster-info` or `minikube status`",
		}
	}

	cmdNodes := exec.CommandContext(ctx, "kubectl", "get", "nodes", "--no-headers")
	outNodes, err := cmdNodes.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Kubernetes Cluster",
			Status:      StatusWarn,
			Details:     "Unable to query nodes",
			Remediation: "Verify node state with `kubectl get nodes`",
		}
	}

	lines := strings.Split(strings.TrimSpace(string(outNodes)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return CheckResult{
			Name:        "Kubernetes Cluster",
			Status:      StatusWarn,
			Details:     "Cluster reachable but 0 nodes found",
			Remediation: "Ensure Minikube has completed bootstrapping nodes",
		}
	}

	readyCount := 0
	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) >= 2 && fields[1] == "Ready" {
			readyCount++
		}
	}

	return CheckResult{
		Name:    "Kubernetes Cluster",
		Status:  StatusOK,
		Details: fmt.Sprintf("Control plane online, %d/%d nodes Ready", readyCount, len(lines)),
	}
}

// TerraformChecker verifies terraform CLI is installed
type TerraformChecker struct{}

func (c *TerraformChecker) Run() CheckResult {
	cmd := exec.Command("terraform", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Terraform CLI",
			Status:      StatusFail,
			Details:     "terraform binary not found in PATH",
			Remediation: "Install Terraform (`sudo pacman -S terraform`)",
		}
	}

	firstLine := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	return CheckResult{
		Name:    "Terraform CLI",
		Status:  StatusOK,
		Details: strings.TrimSpace(firstLine),
	}
}

// GitChecker verifies git is installed (optional check, does not fail)
type GitChecker struct{}

func (c *GitChecker) Run() CheckResult {
	cmd := exec.Command("git", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Git CLI (Optional)",
			Status:      StatusWarn,
			Details:     "git binary not found in PATH",
			Remediation: "Install Git (`sudo pacman -S git`) for repository operations and hooks",
		}
	}

	return CheckResult{
		Name:    "Git CLI (Optional)",
		Status:  StatusOK,
		Details: strings.TrimSpace(string(out)),
	}
}

// TrivyChecker verifies trivy host binary or notes Docker fallback
type TrivyChecker struct{}

func (c *TrivyChecker) Run() CheckResult {
	cmd := exec.Command("trivy", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:        "Aqua Trivy",
			Status:      StatusWarn,
			Details:     "Host binary not found (CLI will fallback to aquasec/trivy container)",
			Remediation: "Optional: Install host trivy for faster scans (`sudo pacman -S trivy`)",
		}
	}

	firstLine := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	return CheckResult{
		Name:    "Aqua Trivy",
		Status:  StatusOK,
		Details: strings.TrimSpace(firstLine),
	}
}

func (c *TrivyChecker) Heal() error {
	// Pre-cache container fallback so preflight can scan without latency
	cmd := exec.Command("docker", "pull", "aquasec/trivy:latest")
	return cmd.Run()
}
