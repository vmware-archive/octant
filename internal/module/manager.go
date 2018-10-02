package module

import (
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/overview"
)

// ManagerInterface is an interface for managing module lifecycle.
type ManagerInterface interface {
	Load() ([]Module, error)
}

// Manager manages module lifecycle.
type Manager struct {
	clusterClient *cluster.Cluster

	loadedModules []Module
}

var _ ManagerInterface = (*Manager)(nil)

// NewManager creates an instance of Manager.
func NewManager(clusterClient *cluster.Cluster) *Manager {
	return &Manager{
		clusterClient: clusterClient,
	}
}

// Load loads modules.
func (m *Manager) Load() ([]Module, error) {
	o := overview.NewClusterOverview(m.clusterClient)

	modules := []Module{
		o,
	}

	for _, module := range modules {
		if err := module.Start(); err != nil {
			return nil, err
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
