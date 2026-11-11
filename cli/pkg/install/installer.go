package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const PrePushHookScript = `#!/usr/bin/env bash
# Ensure common user-space and system bin paths are present for GUI git clients
export PATH="$HOME/.local/bin:/usr/local/bin:$PATH"

echo "🛡️  Running pde pre-flight verification gate..."
pde pre-flight
if [ $? -ne 0 ]; then
  echo "❌ Pre-flight check failed! Push aborted."
  exit 1
fi
`

// InstallBinary copies the currently running or built binary into user's PATH
func InstallBinary() (string, error) {
	destDir, err := resolveInstallDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to resolve current executable: %w", err)
	}

	destPath := filepath.Join(destDir, "pde")

	// Open source
	srcFile, err := os.Open(execPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source executable: %w", err)
	}
	defer srcFile.Close()

	// Write destination
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to write to %s: %w", destPath, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", fmt.Errorf("failed to copy binary: %w", err)
	}

	return destPath, nil
}

// CheckPathInEnvironment verifies whether target directory is present in $PATH
func CheckPathInEnvironment(dir string) bool {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, string(os.PathListSeparator))
	absTarget, _ := filepath.Abs(dir)

	for _, d := range dirs {
		absD, _ := filepath.Abs(d)
		if absD == absTarget {
			return true
		}
	}
	return false
}

// InstallHook installs the pre-push script into targetRepo/.git/hooks/pre-push
func InstallHook(targetRepoPath string) (string, error) {
	absPath, err := filepath.Abs(targetRepoPath)
	if err != nil {
		return "", err
	}

	gitDir := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory %s is not a git repository (missing .git/)", absPath)
	}

	hooksDir := filepath.Join(gitDir, "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create hooks directory: %w", err)
	}

	hookFile := filepath.Join(hooksDir, "pre-push")
	if err := os.WriteFile(hookFile, []byte(PrePushHookScript), 0755); err != nil {
		return "", fmt.Errorf("failed to write pre-push hook: %w", err)
	}

	return hookFile, nil
}

func resolveInstallDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/usr/local/bin", nil
	}

	localBin := filepath.Join(home, ".local", "bin")
	return localBin, nil
}
