package overview

import (
	"log"
	"net/http"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/hcli"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client *cluster.Cluster

	stopFn func()
}

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client *cluster.Cluster) *ClusterOverview {
	return &ClusterOverview{
		client: client,
	}
}

// ContentPath returns the content path for overview.
func (co *ClusterOverview) ContentPath() string {
	return "/overview"
}

// Handler returns a handler for serving overview HTTP content.
func (co *ClusterOverview) Handler(prefix string) http.Handler {
	return newHandler(prefix)
}

// Namespaces returns a list of namespace names for a cluster.
func (co *ClusterOverview) Namespaces() ([]string, error) {
	nsClient, err := co.client.NamespaceClient()
	if err != nil {
		return nil, err
	}

	return nsClient.Names()
}

// Navigation returns navigation entries for overview.
func (co *ClusterOverview) Navigation(root string) (*hcli.Navigation, error) {
	return navigationEntries(root)
}

// Content returns content for a path.
func (co *ClusterOverview) Content() error {
	return nil
}

// Start starts overview.
func (co *ClusterOverview) Start() error {
	log.Printf("Starting cluster overview")
	return nil
}

// Stop stops overview.
func (co *ClusterOverview) Stop() {
	if co.stopFn != nil {
		log.Printf("Stopping cluster overview")

		co.stopFn()
	}
}
