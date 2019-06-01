package config

import (
	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/pkg/plugin"
)

//go:generate mockgen -source=dash.go -destination=./fake/mock_dash.go -package=fake github.com/heptio/developer-dash/internal/config Dash

// Config is configuration for dash.
type Dash interface {
	ObjectPath(namespace, apiVersion, kind, name string) (string, error)

	ClusterClient() cluster.ClientInterface

	ObjectStore() objectstore.ObjectStore

	Logger() log.Logger

	PluginManager() plugin.ManagerInterface

	PortForwarder() portforward.PortForwarder

	Validate() error
}

type Live struct {
	clusterClient cluster.ClientInterface
	logger        log.Logger
	moduleManager module.ManagerInterface
	objectStore   objectstore.ObjectStore
	pluginManager plugin.ManagerInterface
	portForwarder portforward.PortForwarder
}

var _ Dash = (*Live)(nil)

func NewLiveConfig(
	clusterClient cluster.ClientInterface,
	logger log.Logger,
	moduleManager module.ManagerInterface,
	objectStore objectstore.ObjectStore,
	pluginManager *plugin.Manager,
	portForwarder portforward.PortForwarder,
) *Live {
	return &Live{
		clusterClient: clusterClient,
		logger:        logger,
		moduleManager: moduleManager,
		objectStore:   objectStore,
		pluginManager: pluginManager,
		portForwarder: portForwarder,
	}
}

func (l *Live) ObjectPath(namespace, apiVersion, kind, name string) (string, error) {
	return l.moduleManager.ObjectPath(namespace, apiVersion, kind, name)
}

func (l *Live) ClusterClient() cluster.ClientInterface {
	return l.clusterClient
}

func (l *Live) ObjectStore() objectstore.ObjectStore {
	return l.objectStore
}

func (l *Live) Logger() log.Logger {
	return l.logger
}

func (l *Live) PluginManager() plugin.ManagerInterface {
	return l.pluginManager
}

func (l *Live) PortForwarder() portforward.PortForwarder {
	return l.portForwarder
}

func (l *Live) Validate() error {
	if l.clusterClient == nil {
		return errors.New("cluster client is nil")
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
