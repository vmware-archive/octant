/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"context"

	"github.com/heptio/developer-dash/internal/componentcache"
	"github.com/heptio/developer-dash/pkg/store"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/pkg/plugin"
)

//go:generate mockgen -source=dash.go -destination=./fake/mock_dash.go -package=fake github.com/heptio/developer-dash/internal/config Dash

// CRDWatcher watches for CRDs.
type CRDWatcher interface {
	Watch(ctx context.Context, config *CRDWatchConfig) error
}

// ObjectHandler is a function that is run when a new object is available.
type ObjectHandler func(ctx context.Context, object *unstructured.Unstructured)

// CRDWatchConfig is configuration for CRDWatcher.
type CRDWatchConfig struct {
	Add          ObjectHandler
	Delete       ObjectHandler
	IsNamespaced bool
}

// CanPerform returns true if config can perform actions on an object.
func (c *CRDWatchConfig) CanPerform(u *unstructured.Unstructured) bool {
	spec, ok := u.Object["spec"].(map[string]interface{})
	if !ok {
		return false
	}

	scope, ok := spec["scope"].(string)
	if !ok {
		return false
	}

	if c.IsNamespaced && scope != string(apiextv1beta1.NamespaceScoped) {
		return false
	}

	if !c.IsNamespaced && scope != string(apiextv1beta1.ClusterScoped) {
		return false
	}

	return true
}

// Config is configuration for dash. It has knowledge of the all the major sections of
// dash.
type Dash interface {
	ObjectPath(namespace, apiVersion, kind, name string) (string, error)

	ClusterClient() cluster.ClientInterface

	CRDWatcher() CRDWatcher

	ObjectStore() store.Store

	ComponentCache() componentcache.ComponentCache

	Logger() log.Logger

	PluginManager() plugin.ManagerInterface

	PortForwarder() portforward.PortForwarder

	KubeConfigPath() string

	UseContext(ctx context.Context, contextName string) error

	ContextName() string

	Validate() error
}

// Live is a live version of dash config.
type Live struct {
	clusterClient      cluster.ClientInterface
	crdWatcher         CRDWatcher
	logger             log.Logger
	moduleManager      module.ManagerInterface
	objectStore        store.Store
	componentCache     componentcache.ComponentCache
	pluginManager      plugin.ManagerInterface
	portForwarder      portforward.PortForwarder
	kubeConfigPath     string
	currentContextName string
}

var _ Dash = (*Live)(nil)

// NewLiveConfig creates an instance of Live.
func NewLiveConfig(
	clusterClient cluster.ClientInterface,
	crdWatcher CRDWatcher,
	kubeConfigPath string,
	logger log.Logger,
	moduleManager module.ManagerInterface,
	objectStore store.Store,
	componentCache componentcache.ComponentCache,
	pluginManager plugin.ManagerInterface,
	portForwarder portforward.PortForwarder,
	currentContextName string,
) *Live {
	l := &Live{
		clusterClient:      clusterClient,
		crdWatcher:         crdWatcher,
		kubeConfigPath:     kubeConfigPath,
		logger:             logger,
		moduleManager:      moduleManager,
		objectStore:        objectStore,
		componentCache:     componentCache,
		pluginManager:      pluginManager,
		portForwarder:      portForwarder,
		currentContextName: currentContextName,
	}
	objectStore.RegisterOnUpdate(func(store store.Store) {
		l.objectStore = store
	})

	return l
}

// ObjectPath returns the path given an object description.
func (l *Live) ObjectPath(namespace, apiVersion, kind, name string) (string, error) {
	return l.moduleManager.ObjectPath(namespace, apiVersion, kind, name)
}

// ClusterClient returns a cluster client.
func (l *Live) ClusterClient() cluster.ClientInterface {
	return l.clusterClient
}

// CRDWatcher returns a CRD watcher.
func (l *Live) CRDWatcher() CRDWatcher {
	return l.crdWatcher
}

// Store returns an object store.
func (l *Live) ObjectStore() store.Store {
	return l.objectStore
}

// ComponentCache returns an component cache.
func (l *Live) ComponentCache() componentcache.ComponentCache {
	return l.componentCache
}

// KubeConfigPath returns the kube config path.
func (l *Live) KubeConfigPath() string {
	return l.kubeConfigPath
}

// Logger returns a logger.
func (l *Live) Logger() log.Logger {
	return l.logger
}

// PluginManager returns a plugin manager.
func (l *Live) PluginManager() plugin.ManagerInterface {
	return l.pluginManager
}

// PortForwarder returns a port forwarder.
func (l *Live) PortForwarder() portforward.PortForwarder {
	return l.portForwarder
}

// UseContext switches context name.
func (l *Live) UseContext(ctx context.Context, contextName string) error {
	client, err := cluster.FromKubeConfig(ctx, l.kubeConfigPath, contextName)
	if err != nil {
		return err
	}

	l.ClusterClient().Close()
	l.clusterClient = client

	if err := l.objectStore.UpdateClusterClient(ctx, client); err != nil {
		return err
	}

	if err := l.moduleManager.UpdateContext(ctx, contextName); err != nil {
		return err
	}

	l.currentContextName = contextName
	l.Logger().With("new-kube-context", contextName).Infof("updated kube config context")

	return nil
}

// ContextName returns the current context name
func (l *Live) ContextName() string {
	return l.currentContextName
}

// Validate validates the configuration and returns an error if there is an issue.
func (l *Live) Validate() error {
	if l.clusterClient == nil {
		return errors.New("cluster client is nil")
	}

	if l.crdWatcher == nil {
		return errors.New("crd watcher is nil")
	}

	if l.logger == nil {
		return errors.New("logger is nil")
	}

	if l.moduleManager == nil {
		return errors.New("module manager is nil")
	}

	if l.objectStore == nil {
		return errors.New("object store is nil")
	}

	if l.pluginManager == nil {
		return errors.New("plugin manager is nil")
	}

	if l.portForwarder == nil {
		return errors.New("port forwarder is nil")
	}

	return nil
}
