/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

//go:generate mockgen -destination=./fake/mock_manager.go -package=fake github.com/vmware-tanzu/octant/pkg/plugin ManagerInterface
//go:generate mockgen -destination=./fake/mock_module_registrar.go -package=fake github.com/vmware-tanzu/octant/pkg/plugin ModuleRegistrar
//go:generate mockgen -destination=./fake/mock_action_registrar.go -package=fake github.com/vmware-tanzu/octant/pkg/plugin ActionRegistrar

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ClientFactory is a factory for creating clients.
type ClientFactory interface {
	// Init initializes a client.
	Init(ctx context.Context, cmd string) Client
}

// DefaultClientFactory is the default client factory
type DefaultClientFactory struct{}

var _ ClientFactory = (*DefaultClientFactory)(nil)

// NewDefaultClientFactory creates an instance of DefaultClientFactory.
func NewDefaultClientFactory() *DefaultClientFactory {
	return &DefaultClientFactory{}
}

// Init creates a new client.
func (f *DefaultClientFactory) Init(ctx context.Context, cmd string) Client {
	loggerAdapter := &zapAdapter{
		dashLogger: log.From(ctx),
	}

	return plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         pluginMap,
		Cmd:             exec.Command(cmd),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		Logger: loggerAdapter,
	})
}

// Client is an interface that describes a plugin client.
type Client interface {
	Client() (plugin.ClientProtocol, error)
	Kill()
}

// ManagerStore is the data store for Manager.
type ManagerStore interface {
	Store(name string, client Client, metadata *Metadata, cmd string) error
	GetMetadata(name string) (*Metadata, error)
	GetService(name string) (Service, error)
	GetCommand(name string) (string, error)
	Clients() map[string]Client
	ClientNames() []string
}

// DefaultStore is the default implement of ManagerStore.
type DefaultStore struct {
	clients  map[string]Client
	metadata map[string]Metadata
	commands map[string]string
}

var _ ManagerStore = (*DefaultStore)(nil)

// NewDefaultStore creates an instance of DefaultStore.
func NewDefaultStore() *DefaultStore {
	return &DefaultStore{
		clients:  make(map[string]Client),
		metadata: make(map[string]Metadata),
		commands: make(map[string]string),
	}
}

// Store stores information for a plugin.
func (s *DefaultStore) Store(name string, client Client, metadata *Metadata, cmd string) error {
	if metadata == nil {
		return errors.New("metadata is nil")
	}

	s.clients[name] = client
	s.metadata[name] = *metadata
	s.commands[name] = cmd

	return nil
}

// GetService gets the service for a plugin.
func (s *DefaultStore) GetService(name string) (Service, error) {
	client, ok := s.clients[name]
	if !ok {
		return nil, errors.Errorf("plugin %q doesn't have a client", name)
	}

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense("plugin")
	if err != nil {
		return nil, errors.Wrapf(err, "dispensing plugin for %q", name)
	}

	service, ok := raw.(Service)
	if !ok {
		return nil, errors.Errorf("unknown type for plugin %q: %T", name, raw)
	}

	return service, nil
}

// GetMetadata gets the metadata for a plugin.
func (s *DefaultStore) GetMetadata(name string) (*Metadata, error) {
	metadata, ok := s.metadata[name]
	if !ok {
		return nil, errors.Errorf("plugin %q doesn't have metadata", name)
	}

	return &metadata, nil
}

// GetCommand gets the command for a plugin.
func (s *DefaultStore) GetCommand(name string) (string, error) {
	cmd, ok := s.commands[name]
	if !ok {
		return "", errors.Errorf("plugin %q doesn't have command", name)
	}

	return cmd, nil
}

// Clients returns all the clients in the store.
func (s *DefaultStore) Clients() map[string]Client {
	return s.clients
}

// ClientNames returns the client names in the store.
func (s *DefaultStore) ClientNames() []string {
	var list []string
	for name := range s.Clients() {
		list = append(list, name)
	}
	return list
}

type config struct {
	cmd  string
	name string
}

// ManagerInterface is an interface which represent a plugin manager.
type ManagerInterface interface {
	// Print prints an object.
	Print(ctx context.Context, object runtime.Object) (*PrintResponse, error)

	// Tabs retrieves tabs for an object.
	Tabs(ctx context.Context, object runtime.Object) ([]component.Tab, error)

	// Store returns the manager's storage.
	Store() ManagerStore

	// ObjectStatus returns the object status
	ObjectStatus(ctx context.Context, object runtime.Object) (*ObjectStatusResponse, error)
}

// ModuleRegistrar is a module registrar.
type ModuleRegistrar interface {
	// Register registers a module.
	Register(mod module.Module) error
}

// ActionRegistrar is an action registrar.
type ActionRegistrar interface {
	// Register registers an action.
	Register(actionPath string, actionFunc action.DispatcherFunc) error
}

// ManagerOption is an option for configuring Manager.
type ManagerOption func(*Manager)

// Manager manages plugins
type Manager struct {
	PortForwarder   portforward.PortForwarder
	API             api.API
	ClientFactory   ClientFactory
	ModuleRegistrar ModuleRegistrar
	ActionRegistrar ActionRegistrar

	Runners Runners

	configs []config
	store   ManagerStore

	lock sync.Mutex
}

var _ ManagerInterface = (*Manager)(nil)

// NewManager creates an instance of Manager.
func NewManager(apiService api.API, moduleRegistrar ModuleRegistrar, actionRegistrar ActionRegistrar, options ...ManagerOption) *Manager {
	m := &Manager{
		store:           NewDefaultStore(),
		ClientFactory:   NewDefaultClientFactory(),
		Runners:         newDefaultRunners(),
		API:             apiService,
		ModuleRegistrar: moduleRegistrar,
		ActionRegistrar: actionRegistrar,
	}

	for _, option := range options {
		option(m)
	}

	return m
}

// Store returns the store for the manager.
func (m *Manager) Store() ManagerStore {
	return m.store
}

// SetStore sets the store for the manager.
func (m *Manager) SetStore(store ManagerStore) {
	m.store = store
}

// Load loads a plugin.
func (m *Manager) Load(cmd string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	name := filepath.Base(cmd)

	for _, c := range m.configs {
		if name == c.name {
			return errors.Errorf("tried to load plugin %q more than once", name)
		}
	}

	c := config{
		name: name,
		cmd:  cmd,
	}

	m.configs = append(m.configs, c)

	return nil
}

// Start starts all plugins.
func (m *Manager) Start(ctx context.Context) error {
	if m.store == nil {
		return errors.New("manager store is nil")
	}

	if m.ClientFactory == nil {
		return errors.New("manager client factory is nil")
	}

	if err := m.API.Start(ctx); err != nil {
		return errors.Wrap(err, "start api service")
	}

	logger := log.From(ctx)
	logger.With("addr", m.API.Addr()).Debugf("starting plugin api service")

	m.lock.Lock()
	defer m.lock.Unlock()

	for i := range m.configs {
		c := m.configs[i]

		if err := m.start(ctx, c); err != nil {
			return err
		}
	}

	go m.watchPlugins(ctx)

	return nil
}

func (m *Manager) watchPlugins(ctx context.Context) {
	logger := log.From(ctx)

	timer := time.NewTimer(5 * time.Second)
	running := true

	for running {
		select {
		case <-ctx.Done():
			logger.Infof("shutting down plugin watcher")
			running = false
			break
		case <-timer.C:
			for clientName, client := range m.store.Clients() {
				rpcClient, err := client.Client()
				if err != nil {
					logger.WithErr(err).Errorf("retrieve plugin client for ping")
				}

				if err := rpcClient.Ping(); err != nil {
					logger.With("plugin-name", clientName).Infof("restarting plugin")

					cmd, err := m.store.GetCommand(clientName)
					if err != nil {
						logger.WithErr(err).Errorf("unable to find command for plugin")
						continue
					}

					c := config{
						name: clientName,
						cmd:  cmd,
					}

					if err := m.start(ctx, c); err != nil {
						logger.WithErr(err).Errorf("unable to restart plugin")
						continue
					}
				}
			}

			timer.Reset(5 * time.Second)
		}
	}

}

func (m *Manager) start(ctx context.Context, c config) error {
	client := m.ClientFactory.Init(ctx, c.cmd)

	rpcClient, err := client.Client()
	if err != nil {
		return errors.Wrapf(err, "get rpc client for %q", c.name)
	}

	pluginLogger := log.From(ctx).With("plugin-name", c.name)

	raw, err := rpcClient.Dispense("plugin")
	if err != nil {
		return errors.Wrapf(err, "dispensing plugin for %q", c.name)
	}

	service, ok := raw.(Service)
	if !ok {
		return errors.Errorf("unknown type for plugin %q: %T", c.name, raw)
	}

	metadata, err := service.Register(ctx, m.API.Addr())
	if err != nil {
		return errors.Wrapf(err, "register plugin %q", c.name)
	}

	if err := m.store.Store(c.name, client, &metadata, c.cmd); err != nil {
		return errors.Wrapf(err, "storing plugin")
	}

	for _, actionName := range metadata.Capabilities.ActionNames {
		pluginLogger.With("action-path", actionName).Infof("registering plugin action")
		err := m.ActionRegistrar.Register(actionName, func(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
			return service.HandleAction(ctx, payload)
		})

		if err != nil {
			return errors.Wrap(err, "configuring plugin action")
		}
	}

	pluginLogger.With(
		"cmd", c.cmd,
		"metadata", metadata,
	).Infof("registered plugin %q", metadata.Name)

	if metadata.Capabilities.IsModule {
		service, ok := raw.(ModuleService)
		if !ok {
			return errors.Errorf("plugin type %T is a not a module", raw)
		}

		pluginLogger.Infof("plugin supports navigation")

		mp, err := NewModuleProxy(c.name, &metadata, service)
		if err != nil {
			return errors.Wrap(err, "creating module proxy")
		}

		if err := m.ModuleRegistrar.Register(mp); err != nil {
			return errors.Wrapf(err, "register module %s", metadata.Name)
		}
	}

	return nil
}

// Stop stops all plugins.
func (m *Manager) Stop(ctx context.Context) {
	logger := log.From(ctx)

	m.lock.Lock()
	defer m.lock.Unlock()

	for name, client := range m.store.Clients() {
		logger.With("plugin-name", name).Debugf("stopping plugin")
		client.Kill()
	}
}

// Print prints an object with plugins which are configured to print the objects's
// GVK.
func (m *Manager) Print(ctx context.Context, object runtime.Object) (*PrintResponse, error) {
	if m.Runners == nil {
		return nil, errors.New("runners is nil")
	}

	runner, ch := m.Runners.Print(m.store)
	done := make(chan bool)

	var pr PrintResponse

	go func() {
		for resp := range ch {
			pr.Config = append(pr.Config, resp.Config...)
			pr.Status = append(pr.Status, resp.Status...)
			pr.Items = append(pr.Items, resp.Items...)
		}

		done <- true
	}()

	if err := runner.Run(ctx, object, m.store.ClientNames()); err != nil {
		return nil, fmt.Errorf("print runner failed: %w", err)
	}
	close(ch)

	<-done

	// Attempt to eliminate whitespace before fallback
	sort.Slice(pr.Items, func(i, j int) bool {
		if a, b := pr.Items[i].Width, pr.Items[j].Width; a != b {
			return a < b
		}

		a, _ := component.TitleFromTitleComponent(pr.Items[i].View.GetMetadata().Title)
		b, _ := component.TitleFromTitleComponent(pr.Items[j].View.GetMetadata().Title)

		return a < b
	})

	return &pr, nil
}

// Tabs queries plugins for tabs for an object.
func (m *Manager) Tabs(ctx context.Context, object runtime.Object) ([]component.Tab, error) {
	if m.Runners == nil {
		return nil, errors.New("runners is nil")
	}

	runner, ch := m.Runners.Tab(m.store)
	done := make(chan bool)

	var tabs []component.Tab

	go func() {
		for tab := range ch {
			tabs = append(tabs, tab)
		}

		done <- true
	}()

	if err := runner.Run(ctx, object, m.store.ClientNames()); err != nil {
		return nil, err
	}

	close(ch)
	<-done

	sort.Slice(tabs, func(i, j int) bool {
		return tabs[i].Name < tabs[j].Name
	})

	return tabs, nil
}

// ObjectStatus updates the object status of an object configured from a plugin
func (m *Manager) ObjectStatus(ctx context.Context, object runtime.Object) (*ObjectStatusResponse, error) {
	if m.Runners == nil {
		return nil, errors.New("runners is nil")
	}

	runner, ch := m.Runners.ObjectStatus(m.store)
	done := make(chan bool)

	var osr ObjectStatusResponse

	go func() {
		for resp := range ch {
			osr.ObjectStatus = resp.ObjectStatus
		}

		done <- true
	}()

	if err := runner.Run(ctx, object, m.store.ClientNames()); err != nil {
		return nil, err
	}
	close(ch)

	<-done
	return &osr, nil
}
