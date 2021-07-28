/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/store"

	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/pkg/cluster"
	oerrors "github.com/vmware-tanzu/octant/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/plugin"
)

//go:generate mockgen -destination=./fake/mock_dash.go -package=fake github.com/vmware-tanzu/octant/pkg/config Dash
//go:generate mockgen -destination=./fake/mock_kubecontextdecorator.go -package=fake github.com/vmware-tanzu/octant/internal/config KubeContextDecorator

// KubeContextDecorator handles context changes
type KubeContextDecorator interface {
	SwitchContext(context.Context, string) error
	ClusterClient() cluster.ClientInterface
	CurrentContext() string
	Contexts() []kubeconfig.Context
}

func StaticClusterClient(client cluster.ClientInterface) *staticClusterClient {
	return &staticClusterClient{client}
}

type staticClusterClient struct {
	cluster.ClientInterface
}

func (scc *staticClusterClient) SwitchContext(
	ctx context.Context,
	contextName string,
) error {
	return nil
}
func (scc *staticClusterClient) ClusterClient() cluster.ClientInterface {
	return scc.ClientInterface
}
func (scc *staticClusterClient) CurrentContext() string {
	return ""
}
func (scc *staticClusterClient) Contexts() []kubeconfig.Context {
	return nil
}

// UseFSContext is used to indicate a context switch to the file system Kubeconfig context
const UseFSContext = ""

// Live is a live version of dash config.
type Live struct {
	kubeContextDecorator KubeContextDecorator
	crdWatcher           config.CRDWatcher
	logger               log.Logger
	moduleManager        module.ManagerInterface
	objectStore          store.Store
	errorStore           oerrors.ErrorStore
	pluginManager        plugin.ManagerInterface
	portForwarder        portforward.PortForwarder
	restConfigOptions    cluster.RESTConfigOptions
	buildInfo            config.BuildInfo
	kubeConfigPath       string
	contextChosenInUI    bool
}

var _ config.Dash = (*Live)(nil)

// NewLiveConfig creates an instance of Live.
func NewLiveConfig(
	kubeContextDecorator KubeContextDecorator,
	crdWatcher config.CRDWatcher,
	logger log.Logger,
	moduleManager module.ManagerInterface,
	objectStore store.Store,
	errorStore oerrors.ErrorStore,
	pluginManager plugin.ManagerInterface,
	portForwarder portforward.PortForwarder,
	restConfigOptions cluster.RESTConfigOptions,
	buildInfo config.BuildInfo,
	kubeConfigPath string,
	contextChosenInUI bool,
) *Live {
	l := &Live{
		kubeContextDecorator: kubeContextDecorator,
		crdWatcher:           crdWatcher,
		logger:               logger,
		moduleManager:        moduleManager,
		objectStore:          objectStore,
		errorStore:           errorStore,
		pluginManager:        pluginManager,
		portForwarder:        portForwarder,
		restConfigOptions:    restConfigOptions,
		buildInfo:            buildInfo,
		kubeConfigPath:       kubeConfigPath,
		contextChosenInUI:    contextChosenInUI,
	}

	return l
}

// ObjectPath returns the path given an object description.
func (l *Live) ObjectPath(namespace, apiVersion, kind, name string) (string, error) {
	return l.moduleManager.ObjectPath(namespace, apiVersion, kind, name)
}

// ClusterClient returns a cluster client.
func (l *Live) ClusterClient() cluster.ClientInterface {
	return l.kubeContextDecorator.ClusterClient()
}

// CRDWatcher returns a CRD watcher.
func (l *Live) CRDWatcher() config.CRDWatcher {
	return l.crdWatcher
}

// ObjectStore returns an object store.
func (l *Live) ObjectStore() store.Store {
	return l.objectStore
}

// ErrorStore returns an error store.
func (l *Live) ErrorStore() oerrors.ErrorStore {
	return l.errorStore
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

func (l *Live) SetContextChosenInUI(contextChosen bool) {
	l.contextChosenInUI = contextChosen
}

func (l *Live) UseFSContext(ctx context.Context) error {
	return l.UseContext(ctx, UseFSContext)
}

// UseContext switches context name. This process should have synchronously.
func (l *Live) UseContext(ctx context.Context, contextName string) error {
	if l.contextChosenInUI && contextName == UseFSContext {
		contextName = l.CurrentContext()
	}

	err := l.kubeContextDecorator.SwitchContext(ctx, contextName)
	if err != nil {
		return err
	}

	client := l.kubeContextDecorator.ClusterClient()
	if err := l.objectStore.UpdateClusterClient(ctx, client); err != nil {
		return err
	}

	if err := l.moduleManager.UpdateContext(ctx, contextName); err != nil {
		return err
	}

	l.Logger().With("new-kube-context", contextName).Infof("updated kube config context")

	for _, m := range l.moduleManager.Modules() {
		if err := m.ResetCRDs(ctx); err != nil {
			return fmt.Errorf("unable to reset CRDs for module %s, %w", m.Name(), err)
		}
	}

	l.pluginManager.SetOctantClient(l)

	return nil
}

// CurrentContext returns the current context name
func (l *Live) CurrentContext() string {
	return l.kubeContextDecorator.CurrentContext()
}

// Contexts returns the set of all contexts
func (l *Live) Contexts() []kubeconfig.Context {
	return l.kubeContextDecorator.Contexts()
}

// DefaultNamespace returns the default namespace for the current cluster..
func (l *Live) DefaultNamespace() string {
	return l.ClusterClient().DefaultNamespace()
}

// Validate validates the configuration and returns an error if there is an issue.
func (l *Live) Validate() error {
	if l.ClusterClient() == nil {
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

func (l *Live) ModuleManager() module.ManagerInterface {
	return l.moduleManager
}

// BuildInfo returns build ldflag strings for version, commit hash, and build time
func (l *Live) BuildInfo() (string, string, string) {
	return l.buildInfo.Version, l.buildInfo.Commit, l.buildInfo.Time
}

func (l *Live) KubeConfigPath() string {
	return l.kubeConfigPath
}
