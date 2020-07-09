/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package module

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
)

//go:generate mockgen -destination=./fake/mock_manager.go -package=fake github.com/vmware-tanzu/octant/internal/module ManagerInterface
//go:generate mockgen -destination=./fake/mock_action_registrar.go -package=fake github.com/vmware-tanzu/octant/internal/module ActionRegistrar

type ActionReceiver interface {
	ActionPaths() map[string]action.DispatcherFunc
}

type ActionRegistrar interface {
	Register(actionPath, pluginPath string, actionFunc action.DispatcherFunc) error
	Unregister(actionPath, pluginPath string)
}

// ManagerInterface is an interface for managing module lifecycle.
type ManagerInterface interface {
	Modules() []Module
	Register(mod Module) error
	Unregister(mod Module)
	SetNamespace(namespace string)
	GetNamespace() string
	UpdateContext(ctx context.Context, contextName string) error

	ModuleForContentPath(contentPath string) (Module, bool)

	ClientRequestHandlers() []octant.ClientRequestHandler

	ObjectPath(namespace, apiVersion, kind, name string) (string, error)
}

// Manager manages module lifecycle.
type Manager struct {
	clusterClient   cluster.ClientInterface
	namespace       string
	actionRegistrar ActionRegistrar
	logger          log.Logger

	registeredModules []Module

	loadedModules []Module
	mu            sync.Mutex
}

var _ ManagerInterface = (*Manager)(nil)

// NewManager creates an instance of Manager.
func NewManager(clusterClient cluster.ClientInterface, namespace string, actionRegistrar ActionRegistrar, logger log.Logger) (*Manager, error) {
	manager := &Manager{
		clusterClient:   clusterClient,
		namespace:       namespace,
		actionRegistrar: actionRegistrar,
		logger:          logger.With("component", "module-manager"),
	}

	return manager, nil
}

// Register register a module with the manager.
func (m *Manager) Register(mod Module) error {
	m.registeredModules = append(m.registeredModules, mod)

	if receiver, ok := mod.(ActionReceiver); ok {
		for actionPath, actionFunc := range receiver.ActionPaths() {
			m.logger.With("actionPath", actionPath, "module-name", mod.Name()).Infof("registering action")
			if err := m.actionRegistrar.Register(actionPath, mod.Name(), actionFunc); err != nil {
				return err
			}
		}
	}

	if err := mod.Start(); err != nil {
		return errors.Wrapf(err, "%s module failed to start", mod.Name())
	}

	m.loadedModules = append(m.loadedModules, mod)

	return nil
}

// Register register a module with the manager.
func (m *Manager) Unregister(mod Module) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var delIdx int
	var found bool

	for i := 0; i < len(m.loadedModules); i++ {
		if m.loadedModules[i].Name() == mod.Name() {
			delIdx = i
			found = true
			break
		}
	}

	if found {
		m.loadedModules[delIdx].Stop()
		m.loadedModules = append(m.loadedModules[:delIdx], m.loadedModules[delIdx+1:]...)
		for i := 0; i < len(m.registeredModules); i++ {
			if m.registeredModules[i].Name() == mod.Name() {
				delIdx = i
			}
		}
		m.registeredModules = append(m.registeredModules[:delIdx], m.registeredModules[delIdx+1:]...)
	}
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
	for _, module := range m.loadedModules {
		if err := module.SetNamespace(namespace); err != nil {
			m.logger.Errorf("setting namespace for module %q: %v",
				module.Name(), err)
		}
	}
}

// GetNamespace gets the current namespace.
func (m *Manager) GetNamespace() string {
	return m.namespace
}

func (m *Manager) UpdateContext(ctx context.Context, contextName string) error {
	for _, module := range m.loadedModules {
		if err := module.SetContext(ctx, contextName); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) ObjectPath(namespace, apiVersion, kind, name string) (string, error) {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return "", err
	}

	gvk := schema.GroupVersionKind{
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    kind,
	}

	objectPaths := make(map[schema.GroupVersionKind]Module)
	for _, registered := range m.registeredModules {
		for _, supported := range registered.SupportedGroupVersionKind() {
			objectPaths[supported] = registered
		}
	}

	owner, ok := objectPaths[gvk]
	if !ok {
		return "", nil
	}

	return owner.GroupVersionKindPath(namespace, apiVersion, kind, name)
}

func (m *Manager) ModuleForContentPath(contentPath string) (Module, bool) {
	for _, m := range m.Modules() {
		if strings.HasPrefix(contentPath, m.ContentPath()) {
			return m, true
		}
	}

	return nil, false
}

// ClientRequestHandlers returns client request handlers for all modules.
func (m *Manager) ClientRequestHandlers() []octant.ClientRequestHandler {
	var list []octant.ClientRequestHandler

	for _, m := range m.loadedModules {
		list = append(list, m.ClientRequestHandlers()...)
	}

	return list
}
