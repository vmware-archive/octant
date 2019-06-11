package config

import (
	"context"

	"github.com/heptio/developer-dash/internal/componentcache"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/objectstore"
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

	ObjectStore() objectstore.ObjectStore

	ComponentCache() componentcache.ComponentCache

	Logger() log.Logger

	PluginManager() plugin.ManagerInterface

	PortForwarder() portforward.PortForwarder

	KubeConfigsPaths() []string

	Validate() error
}

// Live is a live version of dash config.
type Live struct {
	clusterClient   cluster.ClientInterface
	crdWatcher      CRDWatcher
	logger          log.Logger
	moduleManager   module.ManagerInterface
	objectStore     objectstore.ObjectStore
	componentCache  componentcache.ComponentCache
	pluginManager   plugin.ManagerInterface
	portForwarder   portforward.PortForwarder
	kubeConfigPaths []string
}



var _ Dash = (*Live)(nil)

// NewLiveConfig creates an instance of Live.
func NewLiveConfig(
	clusterClient cluster.ClientInterface,
	crdWatcher CRDWatcher,
	kubeConfigPaths []string,
	logger log.Logger,
	moduleManager module.ManagerInterface,
	objectStore objectstore.ObjectStore,
	componentCache componentcache.ComponentCache,
	pluginManager plugin.ManagerInterface,
	portForwarder portforward.PortForwarder,
) *Live {
	return &Live{
		clusterClient:  clusterClient,
		crdWatcher:     crdWatcher,
		kubeConfigPaths: kubeConfigPaths,
		logger:         logger,
		moduleManager:  moduleManager,
		objectStore:    objectStore,
		componentCache: componentCache,
		pluginManager:  pluginManager,
		portForwarder:  portForwarder,
	}
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

// ObjectStore returns an object store.
func (l *Live) ObjectStore() objectstore.ObjectStore {
	return l.objectStore
}

// ComponentCache returns an component cache.
func (l *Live) ComponentCache() componentcache.ComponentCache {
	return l.componentCache
}

// KubeConfigsPaths returns a slice of kube config paths.
func (l *Live) KubeConfigsPaths() []string {
	return l.kubeConfigPaths
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
