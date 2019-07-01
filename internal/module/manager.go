/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package module

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/action"
	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
)

//go:generate mockgen -destination=./fake/mock_manager.go -package=fake github.com/vmware/octant/internal/module ManagerInterface
//go:generate mockgen -destination=./fake/mock_action_registrar.go -package=fake github.com/vmware/octant/internal/module ActionRegistrar

type ActionReceiver interface {
	ActionPaths() map[string]action.DispatcherFunc
}

type ActionRegistrar interface {
	Register(actionPath string, actionFunc action.DispatcherFunc) error
}

// ManagerInterface is an interface for managing module lifecycle.
type ManagerInterface interface {
	Modules() []Module
	SetNamespace(namespace string)
	GetNamespace() string
	UpdateContext(ctx context.Context, contextName string) error

	ObjectPath(namespace, apiVersion, kind, name string) (string, error)
	RegisterObjectPath(Module, schema.GroupVersionKind)
	DeregisterObjectPath(schema.GroupVersionKind)
}

// Manager manages module lifecycle.
type Manager struct {
	clusterClient   cluster.ClientInterface
	namespace       string
	actionRegistrar ActionRegistrar
	logger          log.Logger

	registeredModules []Module

	loadedModules []Module
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
			if err := m.actionRegistrar.Register(actionPath, actionFunc); err != nil {
				return err
			}
		}
	}

	return nil
}

// Load loads modules.
func (m *Manager) Load() error {
	for _, module := range m.registeredModules {
		if err := module.Start(); err != nil {
			return errors.Wrapf(err, "%s module failed to start", module.Name())
		}
	}

	m.loadedModules = m.registeredModules

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
		return "", errors.Errorf("no module claimed ownership of %s", gvk.String())
	}

	return owner.GroupVersionKindPath(namespace, apiVersion, kind, name)
}

func (m *Manager) RegisterObjectPath(mod Module, gvk schema.GroupVersionKind) {
	//m.objectPaths[gvk] = mod
}

func (m *Manager) DeregisterObjectPath(gvk schema.GroupVersionKind) {
	//delete(m.objectPaths, gvk)
}
