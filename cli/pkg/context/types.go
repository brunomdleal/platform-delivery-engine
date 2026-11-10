package context

// RepoType represents the identified classification of a repository
type RepoType string

const (
	RepoOrderService           RepoType = "order-service"
	RepoOrderWebClient         RepoType = "order-web-client"
	RepoPlatformDeliveryEngine RepoType = "platform-delivery-engine"
	RepoQAAutomation           RepoType = "qa-automation"
	RepoUnknown                RepoType = "unknown"
)

// RepoContext encapsulates the detected repository attributes
type RepoContext struct {
	Type        RepoType
	Path        string
	Description string
	Toolchain   string
}
