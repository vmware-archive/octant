package overview

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/cluster"
)

type notImplemented struct {
	name string
}

func (e *notImplemented) Error() string {
	return fmt.Sprintf("%s not implemented", e.name)
}

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client *cluster.Cluster
}

var _ Interface = (*ClusterOverview)(nil)

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client *cluster.Cluster) *ClusterOverview {
	return &ClusterOverview{
		client: client,
	}
}

// Namespaces returns a list of namespace names for a cluster.
func (co *ClusterOverview) Namespaces() ([]string, error) {
	return co.client.Namespace.Names()
}

// Navigation returns navigation entries for overview.
func (co *ClusterOverview) Navigation() (*Navigation, error) {
	return navigationEntries()
}

// Content returns content for a path.
func (co *ClusterOverview) Content(path string) error {
	return &notImplemented{name: "Content"}
}
