package fake

import "github.com/heptio/developer-dash/internal/overview"

// SimpleClusterOverview is a fake that implements overview.Interface.
type SimpleClusterOverview struct{}

// NewSimpleClusterOverview creates an instance of SimpleClusterOverview.
func NewSimpleClusterOverview() *SimpleClusterOverview {
	return &SimpleClusterOverview{}
}

var _ overview.Interface = (*SimpleClusterOverview)(nil)

// Namespaces returns a list of namespaces. ["default"].
func (sco *SimpleClusterOverview) Namespaces() ([]string, error) {
	names := []string{"default"}
	return names, nil
}

// Navigation is a no-op.
func (sco *SimpleClusterOverview) Navigation() error {
	return nil
}

// Content is a no-op.
func (sco *SimpleClusterOverview) Content(path string) error {
	return nil
}
