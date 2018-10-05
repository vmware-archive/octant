package fake

import "github.com/heptio/developer-dash/internal/cluster"

// NamespaceClient is a fake that implements cluster.NamespaceInterface.
type NamespaceClient struct{}

var _ cluster.NamespaceInterface = (*NamespaceClient)(nil)

// NewNamespaceClient creates an instance of NamespaceClient.
func NewNamespaceClient() *NamespaceClient {
	return &NamespaceClient{}
}

// Names returns ["default"].
func (nc *NamespaceClient) Names() ([]string, error) {
	names := []string{"default"}
	return names, nil
}
