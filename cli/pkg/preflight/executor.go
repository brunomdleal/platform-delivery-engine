package preflight

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
)

func execute(name string, dir string, verbose bool, command string, args ...string) Result {
	start := time.Now()
	output.Info("Running: %s (%s %v)", name, command, args)

	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	var outBuf, errBuf bytes.Buffer
	if verbose {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	} else {
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
	}

	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		if !verbose && (outBuf.Len() > 0 || errBuf.Len() > 0) {
			fmt.Println("\n--- Output Failure ---")
			if outBuf.Len() > 0 {
				fmt.Print(outBuf.String())
			}
			if errBuf.Len() > 0 {
				fmt.Print(errBuf.String())
			}
			fmt.Println("----------------------------")
		}
		return Result{
			Name:     name,
			Status:   StatusFailed,
			Duration: duration,
			Err:      fmt.Errorf("%q failed: %w", name, err),
		}
	}

	return Result{
		Name:     name,
		Status:   StatusPassed,
		Duration: duration,
	}
}

func runTrivyScan(imageTag string, dir string, verbose bool) Result {
	start := time.Now()
	output.Info("Running security scan on image: %s", imageTag)

	_, err := exec.LookPath("trivy")
	var cmd *exec.Cmd
	if err == nil {
		cmd = exec.Command("trivy", "image", "--severity", "HIGH,CRITICAL", "--exit-code", "1", imageTag)
	} else {
		output.Warn("Host 'trivy' not found; using ephemeral Docker container fallback...")
		cmd = exec.Command("docker", "run", "--rm",
			"-v", "/var/run/docker.sock:/var/run/docker.sock",
			"aquasec/trivy:latest",
			"image",
			"--severity", "HIGH,CRITICAL",
			"--exit-code", "1",
			imageTag,
		)
	}
	cmd.Dir = dir

	var outBuf, errBuf bytes.Buffer
	if verbose {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	} else {
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
	}

	runErr := cmd.Run()
	duration := time.Since(start)

	if runErr != nil {
		if !verbose {
			if outBuf.Len() > 0 {
				fmt.Print(outBuf.String())
			}
			if errBuf.Len() > 0 {
				fmt.Print(errBuf.String())
			}
		}
		return Result{
			Name:     "Container Security Audit (Trivy)",
			Status:   StatusFailed,
			Duration: duration,
			Err:      fmt.Errorf("vulnerabilities detected in %s: %w", imageTag, runErr),
		}
	}

	return Result{
		Name:     "Container Security Audit (Trivy)",
		Status:   StatusPassed,
		Duration: duration,
	}
}
