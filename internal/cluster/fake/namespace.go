package fake

import "github.com/heptio/developer-dash/internal/cluster"

// NamespaceClient is a fake that implements cluster.NamespaceInterface.
type NamespaceClient struct {
	namespaces       []string
	errOnList        error
	initialNamespace string
}

var _ cluster.NamespaceInterface = (*NamespaceClient)(nil)

// NewNamespaceClient creates an instance of NamespaceClient.
func NewNamespaceClient(namespaces []string, errOnList error, initialNamespace string) *NamespaceClient {
	return &NamespaceClient{
		namespaces:       namespaces,
		errOnList:        errOnList,
		initialNamespace: initialNamespace,
	}
}

// Names returns ["default"].
func (nc *NamespaceClient) Names() ([]string, error) {
	return nc.namespaces, nc.errOnList
}

// InitialNamespace returns an initial namespace
func (nc *NamespaceClient) InitialNamespace() string {
	return nc.initialNamespace
}
