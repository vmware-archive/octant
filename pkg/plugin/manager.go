package plugin

import (
	"context"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
)

// ClientFactory is a factory for creating clients.
type ClientFactory interface {
	// Init initializes a client.
	Init(cmd string) Client
}

// DefaultClientFactory is the default client factory
type DefaultClientFactory struct{}

var _ ClientFactory = (*DefaultClientFactory)(nil)

// NewDefaultClientFactory creates an instance of DefaultClientFactory.
func NewDefaultClientFactory() *DefaultClientFactory {
	return &DefaultClientFactory{}
}

// Init creates a new client.
func (f *DefaultClientFactory) Init(cmd string) Client {
	return plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         pluginMap,
		Cmd:             exec.Command("sh", "-c", cmd),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})
}

// Client is an interface that describes a plugin client.
type Client interface {
	Client() (plugin.ClientProtocol, error)
	Kill()
}

// ManagerStore is the data store for Manager.
type ManagerStore interface {
	Store(name string, client Client, metadata Metadata)
	GetMetadata(name string) (Metadata, error)
	GetService(name string) (Service, error)
	Clients() map[string]Client
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
func (s *DefaultStore) Store(name string, client Client, metadata Metadata) {
	s.clients[name] = client
	s.metadata[name] = metadata
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

type config struct {
	cmd  string
	name string
}

// ManagerOption is an option for configuring Manager.
type ManagerOption func(*Manager)

// Manager manages plugins
type Manager struct {
	Store         ManagerStore
	ClientFactory ClientFactory

	configs []config

	lock sync.Mutex
}

// NewManager creates an instance of Manager.
func NewManager(options ...ManagerOption) *Manager {
	m := &Manager{
		Store:         NewDefaultStore(),
		ClientFactory: NewDefaultClientFactory(),
	}

	for _, option := range options {
		option(m)
	}

	return m
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

// Start stars all plugins.
func (m *Manager) Start(ctx context.Context) error {
	if m.Store == nil {
		return errors.New("manager store is nil")
	}

	if m.ClientFactory == nil {
		return errors.New("manager client factory is nil")
	}

	logger := log.From(ctx)

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, c := range m.configs {
		logger = logger.With("plugin-name", c.name)

		client := m.ClientFactory.Init(c.cmd)

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

		metadata, err := p.Register()
		if err != nil {
			return errors.Wrapf(err, "register plugin %q", c.name)
		}

		m.Store.Store(c.name, client, metadata)

		logger.With(
			"cmd", c.cmd,
			"metadata", metadata,
		).Debugf("registered plugin")
	}

	return nil
}

// Stop stops all plugins.
func (m *Manager) Stop(ctx context.Context) {
	logger := log.From(ctx)

	m.lock.Lock()
	defer m.lock.Unlock()

	for name, client := range m.Store.Clients() {
		logger.With("plugin-name", name).Debugf("stopping plugin")
		client.Kill()
	}
}

// Print prints an object with plugins which are configured to print the objects's
// GVK.
func (m *Manager) Print(object runtime.Object) (*PrintResponse, error) {
	if object == nil {
		return nil, errors.New("can't print a nil object")
	}

	var mu sync.Mutex
	var g errgroup.Group

	gvk := object.GetObjectKind().GroupVersionKind()

	var pr PrintResponse

	for name := range m.Store.Clients() {
		metadata, err := m.Store.GetMetadata(name)
		if err != nil {
			return nil, err
		}

		fn := func(name string) func() error {
			return func() error {
				resp, err := printObject(m.Store, name, object)
				if err != nil {
					return err
				}

				mu.Lock()

				defer mu.Unlock()

				pr.Config = append(pr.Config, resp.Config...)
				pr.Status = append(pr.Status, resp.Config...)
				pr.Items = append(pr.Items, resp.Items...)

				return nil
			}
		}

		if metadata.Capabilities.SupportsPrinter(gvk) {
			g.Go(fn(name))
		}
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Wrap(err, "print object")
	}

	return &pr, nil
}

func printObject(store ManagerStore, pluginName string, object runtime.Object) (PrintResponse, error) {
	service, err := store.GetService(pluginName)
	if err != nil {
		return PrintResponse{}, err
	}

	resp, err := service.Print(object)
	if err != nil {
		return PrintResponse{}, errors.Wrapf(err, "print object with plugin %q", pluginName)
	}

	return resp, nil
}
