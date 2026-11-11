package preflight

import (
	"fmt"

	"github.com/brunomdleal/platform-delivery-engine/cli/pkg/context"
)

// GetRunner returns the appropriate runner matching the detected repository type
func GetRunner(repoType context.RepoType) (Runner, error) {
	switch repoType {
	case context.RepoOrderService:
		return &JavaRunner{}, nil
	case context.RepoOrderWebClient:
		return &NodeRunner{}, nil
	case context.RepoPlatformDeliveryEngine:
		return &K8sPlatformRunner{}, nil
	case context.RepoQAAutomation:
		return &QARunner{}, nil
	default:
		return nil, fmt.Errorf("unsupported repository context: %s", repoType)
	}
}
