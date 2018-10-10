package module

import (
	"log"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/overview"
)

// ManagerInterface is an interface for managing module lifecycle.
type ManagerInterface interface {
	Modules() []Module
	SetNamespace(namespace string)
	GetNamespace() string
}

// Manager manages module lifecycle.
type Manager struct {
	clusterClient cluster.ClientInterface
	namespace     string

	loadedModules []Module
}

var _ ManagerInterface = (*Manager)(nil)

// NewManager creates an instance of Manager.
func NewManager(clusterClient cluster.ClientInterface, namespace string) (*Manager, error) {
	manager := &Manager{
		clusterClient: clusterClient,
		namespace:     namespace,
	}

	if err := manager.Load(); err != nil {
		return nil, err
	}

	return manager, nil
}

// Load loads modules.
func (m *Manager) Load() error {
	modules := []Module{
		overview.NewClusterOverview(m.clusterClient, m.namespace),
	}

	for _, module := range modules {
		if err := module.Start(); err != nil {
			return errors.Wrapf(err, "%s module failed to start", module.Name())
		}
	}

	m.loadedModules = modules

	return nil
}

// Modules returns a list of modules.
func (m *Manager) Modules() []Module {
	return m.loadedModules
}

// Unload unloads modules.
func (m *Manager) Unload() {
	for _, module := range m.loadedModules {
		module.Stop()
	}
}

// SetNamespace sets the current namespace.
func (m *Manager) SetNamespace(namespace string) {
	m.namespace = namespace
	log.Printf("Setting namespace to %s", namespace)
	for _, module := range m.loadedModules {
		if err := module.SetNamespace(namespace); err != nil {
			log.Printf("ERROR: setting namespace for module %q: %v",
				module.Name(), err)
		}
	}
}

// GetNamespace gets the current namespace.
func (m *Manager) GetNamespace() string {
	return m.namespace
}
