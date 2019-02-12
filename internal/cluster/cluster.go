package cluster

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// ClientInterface is a client for cluster operations.
type ClientInterface interface {
	KubernetesClient() (kubernetes.Interface, error)
	DynamicClient() (dynamic.Interface, error)
	DiscoveryClient() (discovery.DiscoveryInterface, error)
	NamespaceClient() (NamespaceInterface, error)
	InfoClient() (InfoInterface, error)
}

// Cluster is a client for cluster operations
type Cluster struct {
	clientConfig clientcmd.ClientConfig
	restClient   *rest.Config
}

var _ ClientInterface = (*Cluster)(nil)

// KubernetesClient returns a Kubernetes client.
func (c *Cluster) KubernetesClient() (kubernetes.Interface, error) {
	return kubernetes.NewForConfig(c.restClient)
}

// NamespaceClient returns a namespace client.
func (c *Cluster) NamespaceClient() (NamespaceInterface, error) {
	dc, err := c.DynamicClient()
	if err != nil {
		return nil, err
	}

	ns, _, err := c.clientConfig.Namespace()
	if err != nil {
		return nil, errors.Wrap(err, "resolving initial namespace")
	}
	return newNamespaceClient(dc, ns), nil
}

// DynamicClient returns a dynamic client.
func (c *Cluster) DynamicClient() (dynamic.Interface, error) {
	return dynamic.NewForConfig(c.restClient)
}

// DiscoveryClient returns a DiscoveryClient for the cluster.
func (c *Cluster) DiscoveryClient() (discovery.DiscoveryInterface, error) {
	return discovery.NewDiscoveryClientForConfig(c.restClient)
}

// InfoClient returns an InfoClient for the cluster.
func (c *Cluster) InfoClient() (InfoInterface, error) {
	return newClusterInfo(c.clientConfig), nil
}

// Version returns a ServerVersion for the cluster
func (c *Cluster) Version() (string, error) {
	dc, err := c.DiscoveryClient()
	if err != nil {
		return "", err
	}
	serverVersion, err := dc.ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(serverVersion), nil
}

// FromKubeconfig creates a Cluster from a kubeconfig.
func FromKubeconfig(kubeconfig string) (*Cluster, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		rules.ExplicitPath = kubeconfig
	}
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := cc.ClientConfig()
	if err != nil {
		return nil, err
	}

	return &Cluster{
		clientConfig: cc,
		restClient:   config,
	}, nil
}
