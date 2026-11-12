package context

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/output"
)

// Detect inspects a directory path and identifies the repository context
func Detect(dirPath string) (*RepoContext, error) {
	absPath, err := filepath.Abs(dirPath)
	output.Directory("%s", absPath)
	if err != nil {
		return nil, err
	}

	//	check for platform-delivery-engine
	output.Info("Checking for platform-delivery-engine")
	if exists(filepath.Join(absPath, "k8s")) && exists(filepath.Join(absPath, "terraform")) {
		return &RepoContext{
			Type:        RepoPlatformDeliveryEngine,
			Path:        absPath,
			Description: "Platform Orchestration & Delivery Engine",
			Toolchain:   "Terraform / Kubernetes / Docker Compose",
		}, nil
	}

	// order-web-client
	output.Info("Checking for order-web-client")
	if exists(filepath.Join(absPath, "package.json")) {
		return &RepoContext{
			Type:        RepoOrderWebClient,
			Path:        absPath,
			Description: "Order Web Client SPA",
			Toolchain:   "Node.js / React / TypeScript / Vite",
		}, nil
	}

	// order-service
	output.Info("Checking for maven projects")
	pomFilePath := filepath.Join(absPath, "pom.xml")
	if exists(pomFilePath) {
		output.Info("Checking for order-service")
		if exists(filepath.Join(absPath, "src", "main", "java")) {
			return &RepoContext{
				Type:        RepoOrderService,
				Path:        absPath,
				Description: "Order Management Core API",
				Toolchain:   "Java 21 / Maven / Spring Boot",
			}, nil
		}

		// qa-automation
		output.Info("Checking for qa-automation")
		pomContent, err := os.ReadFile(pomFilePath)
		if err == nil {
			pomContentStr := string(pomContent)
			if strings.Contains(pomContentStr, "qa-automation") || strings.Contains(pomContentStr, "atf-") {
				return &RepoContext{
					Type:        RepoQAAutomation,
					Path:        absPath,
					Description: "QA Automated Testing Framework",
					Toolchain:   "Java 21 / Maven / REST-Assured",
				}, nil
			}
		}
	}

	return &RepoContext{
		Type:        RepoUnknown,
		Path:        absPath,
		Description: "Unrecognized repository context",
		Toolchain:   "Unknown",
	}, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
