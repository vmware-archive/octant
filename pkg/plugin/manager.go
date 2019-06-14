package plugin

//go:generate mockgen -destination=./fake/mock_manager.go -package=fake github.com/heptio/developer-dash/pkg/plugin ManagerInterface

import (
	"context"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/pkg/plugin/api"
	"github.com/heptio/developer-dash/pkg/view/component"
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
	Store(name string, client Client, metadata Metadata) error
	GetMetadata(name string) (Metadata, error)
	GetService(name string) (Service, error)
	Clients() map[string]Client
	ClientNames() []string
}

// DefaultStore is the default implement of ManagerStore.
type DefaultStore struct {
	clients  map[string]Client
	metadata map[string]Metadata
}

var _ ManagerStore = (*DefaultStore)(nil)

// NewDefaultStore creates an instance of DefaultStore.
func NewDefaultStore() *DefaultStore {
	return &DefaultStore{
		clients:  make(map[string]Client),
		metadata: make(map[string]Metadata),
	}
}

// Store stores information for a plugin.
func (s *DefaultStore) Store(name string, client Client, metadata Metadata) error {
	if _, ok := s.clients[name]; ok {
		return errors.Errorf("plugin %q is already stored", name)
	}

	s.clients[name] = client
	s.metadata[name] = metadata

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
func (s *DefaultStore) GetMetadata(name string) (Metadata, error) {
	metadata, ok := s.metadata[name]
	if !ok {
		return Metadata{}, errors.Errorf("plugin %q doesn't have metadata", name)
	}

	return metadata, nil
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
	Print(object runtime.Object) (*PrintResponse, error)

	// Tabs retrieves tabs for an object.
	Tabs(object runtime.Object) ([]component.Tab, error)

	// Store returns the manager's storage.
	Store() ManagerStore

	// ObjectStatus returns the object status
	ObjectStatus(object runtime.Object) (*ObjectStatusResponse, error)
}

// ManagerOption is an option for configuring Manager.
type ManagerOption func(*Manager)

// Manager manages plugins
type Manager struct {
	PortForwarder portforward.PortForwarder
	API           api.API
	ClientFactory ClientFactory

	Runners Runners

	configs []config
	store   ManagerStore

	lock sync.Mutex
}

var _ ManagerInterface = (*Manager)(nil)

// NewManager creates an instance of Manager.
func NewManager(apiService api.API, options ...ManagerOption) *Manager {
	m := &Manager{
		store:         NewDefaultStore(),
		ClientFactory: NewDefaultClientFactory(),
		Runners:       newDefaultRunners(),
		API:           apiService,
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

		pluginLogger := logger.With("plugin-name", c.name)

		client := m.ClientFactory.Init(ctx, c.cmd)

		rpcClient, err := client.Client()
		if err != nil {
			return errors.Wrapf(err, "get rpc client for %q", c.name)
		}

		raw, err := rpcClient.Dispense("plugin")
		if err != nil {
			return errors.Wrapf(err, "dispensing plugin for %q", c.name)
		}

		p, ok := raw.(Service)
		if !ok {
			return errors.Errorf("unknown type for plugin %q: %T", c.name, raw)
		}

		metadata, err := p.Register(m.API.Addr())
		if err != nil {
			return errors.Wrapf(err, "register plugin %q", c.name)
		}

		if err := m.store.Store(c.name, client, metadata); err != nil {
			return errors.Wrapf(err, "storing plugin")
		}

		pluginLogger.With(
			"cmd", c.cmd,
			"metadata", metadata,
		).Infof("registered plugin %q", metadata.Name)
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
func (m *Manager) Print(object runtime.Object) (*PrintResponse, error) {
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

	if err := runner.Run(object, m.store.ClientNames()); err != nil {
		return nil, err
	}
	close(ch)

	<-done
	return &pr, nil
}

// Tabs queries plugins for tabs for an object.
func (m *Manager) Tabs(object runtime.Object) ([]component.Tab, error) {
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

	if err := runner.Run(object, m.store.ClientNames()); err != nil {
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
func (m *Manager) ObjectStatus(object runtime.Object) (*ObjectStatusResponse, error) {
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

	if err := runner.Run(object, m.store.ClientNames()); err != nil {
		return nil, err
	}
	close(ch)

	<-done
	return &osr, nil
}
