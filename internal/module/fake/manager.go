package fake

import (
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/module"
)

// StubManager is a stub for module.Module.
type StubManager struct {
	modules   []module.Module
	namespace string
}

// NewStubManager creates an instance of StubManager.
func NewStubManager(namespace string, modules []module.Module) *StubManager {
	return &StubManager{
		modules:   modules,
		namespace: namespace,
	}
}

var _ module.ManagerInterface = (*StubManager)(nil)

// Modules returns the modules stored in the stub.
func (m *StubManager) Modules() []module.Module {
	return m.modules
}

// SetNamespace sets the namespace
func (m *StubManager) SetNamespace(namespace string) {
	m.namespace = namespace
}

// GetNamespace returns the namespace
func (m *StubManager) GetNamespace() string {
	return m.namespace
}

func (m *StubManager) ObjectPath(namespace, apiVersion, kind, name string) (string, error) {
	return "/pod", nil
}

func (m *StubManager) RegisterObjectPath(module.Module, schema.GroupVersionKind) {
	panic("implement me")
}

func (m *StubManager) DeregisterObjectPath(schema.GroupVersionKind) {
	panic("implement me")
}
