package cluster

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// Cluster is a client cluster operations
type Cluster struct {
	Namespace NamespaceInterface
}

// New creates an instance of Cluster.
func New(dynamicClient dynamic.Interface) *Cluster {
	return &Cluster{
		Namespace: newNamespaceClient(dynamicClient),
	}
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

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	clusterClient := New(dc)

	return clusterClient, nil
}
