package module

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/overview"
)

// ManagerInterface is an interface for managing module lifecycle.
type ManagerInterface interface {
	Load() ([]Module, error)
}

// Manager manages module lifecycle.
type Manager struct {
	clusterClient cluster.ClientInterface
	namespace     string

	loadedModules []Module
}

var _ ManagerInterface = (*Manager)(nil)

// NewManager creates an instance of Manager.
func NewManager(clusterClient cluster.ClientInterface, namespace string) *Manager {
	return &Manager{
		clusterClient: clusterClient,
		namespace:     namespace,
	}
}

// Load loads modules.
func (m *Manager) Load() ([]Module, error) {
	modules := []Module{
		overview.NewClusterOverview(m.clusterClient, m.namespace),
	}

	for _, module := range modules {
		if err := module.Start(); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("%s module failed to start", module.Name()))
		}
	}

	m.loadedModules = modules

	return modules, nil
}

// Unload unloads modules.
func (m *Manager) Unload() {
	for _, module := range m.loadedModules {
		module.Stop()
	}
}

// SetNamespace sets the current namespace.
func (m *Manager) SetNamespace(namespace string) {
	for _, module := range m.loadedModules {
		if err := module.SetNamespace(namespace); err != nil {
			log.Printf("ERROR: setting namespace for module %q: %v",
				module.Name(), err)
		}
	}
}
