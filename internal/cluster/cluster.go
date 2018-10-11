package cluster

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// gcp
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// ClientInterface is a client for cluster operations.
type ClientInterface interface {
	DynamicClient() (dynamic.Interface, error)
	DiscoveryClient() (discovery.DiscoveryInterface, error)
	NamespaceClient() (NamespaceInterface, error)
}

// Cluster is a client for cluster operations
type Cluster struct {
	restClient *rest.Config
}

var _ ClientInterface = (*Cluster)(nil)

// New creates an instance of Cluster.
func New(restClient *rest.Config) *Cluster {
	return &Cluster{
		restClient: restClient,
	}
}

// NamespaceClient returns a namespace client.
func (c *Cluster) NamespaceClient() (NamespaceInterface, error) {
	dc, err := c.DynamicClient()
	if err != nil {
		return nil, err
	}

	return newNamespaceClient(dc), nil
}

// DynamicClient returns a dynamic client.
func (c *Cluster) DynamicClient() (dynamic.Interface, error) {
	return dynamic.NewForConfig(c.restClient)
}

// DiscoveryClient returns a DiscoveryClient for the cluster.
func (c *Cluster) DiscoveryClient() (discovery.DiscoveryInterface, error) {
	return discovery.NewDiscoveryClientForConfig(c.restClient)
}

// NamespaceInterface is an interface for querying namespace details.
type NamespaceInterface interface {
	Names() ([]string, error)
}

// FromKubeconfig creates a Cluster from a kubeconfig.
func FromKubeconfig(kubeconfig string) (*Cluster, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clusterClient := New(config)

	return clusterClient, nil
}
