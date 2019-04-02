package plugin_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-plugin"
	"github.com/heptio/developer-dash/internal/testutil"
	dashplugin "github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/plugin/api"
	"github.com/heptio/developer-dash/pkg/plugin/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestDefaultStore(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	name := "name"
	client := newFakePluginClient(name, controller)
	metadata := dashplugin.Metadata{Name: name}

	s := dashplugin.NewDefaultStore()
	s.Store(name, client, metadata)

	gotMetadata, err := s.GetMetadata(name)
	require.NoError(t, err)
	require.Equal(t, metadata, gotMetadata)

	_, err = s.GetMetadata("invalid")
	require.Error(t, err)

	_, err = s.GetService(name)
	require.NoError(t, err)

	_, err = s.GetService("invalid")
	require.Error(t, err)
}

func TestManager(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var options []dashplugin.ManagerOption

	store := fake.NewMockManagerStore(controller)
	clientFactory := fake.NewMockClientFactory(controller)

	name := "plugin1"

	client := newFakePluginClient(name, controller)
	clientFactory.EXPECT().Init(gomock.Any(), gomock.Eq(name)).Return(client)

	metadata := dashplugin.Metadata{
		Name: name,
	}
	store.EXPECT().Store(gomock.Eq(name), gomock.Eq(client), gomock.Eq(metadata))
	store.EXPECT().Clients().Return(map[string]dashplugin.Client{name: client})

	options = append(options, func(m *dashplugin.Manager) {
		m.Store = store
		m.ClientFactory = clientFactory
	})

	apiService := &stubAPIService{}
	manager := dashplugin.NewManager(apiService, options...)

	err := manager.Load(name)
	require.NoError(t, err)

	ctx := context.Background()
	err = manager.Start(ctx)
	require.NoError(t, err)

	manager.Stop(ctx)
}

func TestManager_Print(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pod := testutil.CreatePod("pod")

	var options []dashplugin.ManagerOption

	store := fake.NewMockManagerStore(controller)

	store.EXPECT().ClientNames().Return([]string{"plugin1", "plugin2"})

	ch := make(chan dashplugin.PrintResponse)
	printRunner := dashplugin.DefaultRunner{
		RunFunc: func(name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			if name == "plugin1" {
				resp1 := dashplugin.PrintResponse{
					Config: []component.SummarySection{{Header: "resp1"}},
				}
				resp2 := dashplugin.PrintResponse{
					Config: []component.SummarySection{{Header: "resp2"}},
				}
				ch <- resp1
				ch <- resp2
			}

			return nil
		},
	}

	runners := fake.NewMockRunners(controller)
	runners.EXPECT().
		Print(gomock.Eq(store)).Return(printRunner, ch)

	options = append(options, func(m *dashplugin.Manager) {
		m.Store = store
		m.Runners = runners
	})

	apiService := &stubAPIService{}
	manager := dashplugin.NewManager(apiService, options...)

	got, err := manager.Print(pod)
	require.NoError(t, err)

	expected := &dashplugin.PrintResponse{
		Config: []component.SummarySection{
			{Header: "resp1"},
			{Header: "resp2"},
		},
	}
	assert.Equal(t, expected, got)
}

func TestManager_Tabs(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pod := testutil.CreatePod("pod")

	var options []dashplugin.ManagerOption

	store := fake.NewMockManagerStore(controller)

	store.EXPECT().ClientNames().Return([]string{"plugin1", "plugin2"})

	ch := make(chan component.Tab)
	tabRunner := dashplugin.DefaultRunner{
		RunFunc: func(name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			ch <- component.Tab{Name: name}

			return nil
		},
	}

	runners := fake.NewMockRunners(controller)
	runners.EXPECT().
		Tab(gomock.Eq(store)).Return(tabRunner, ch)

	options = append(options, func(m *dashplugin.Manager) {
		m.Store = store
		m.Runners = runners
	})

	apiService := &stubAPIService{}
	manager := dashplugin.NewManager(apiService, options...)

	got, err := manager.Tabs(pod)
	require.NoError(t, err)

	expected := []component.Tab{
		{
			Name: "plugin1",
		},
		{
			Name: "plugin2",
		},
	}
	assert.Equal(t, expected, got)
}

type fakePluginClient struct {
	clientProtocol *fake.MockClientProtocol
	service        *fake.MockService
	name           string
}

var _ dashplugin.Client = (*fakePluginClient)(nil)

func newFakePluginClient(name string, controller *gomock.Controller) *fakePluginClient {
	service := fake.NewMockService(controller)
	metadata := dashplugin.Metadata{
		Name: name,
	}
	service.EXPECT().Register(gomock.Eq("localhost:54321")).Return(metadata, nil).AnyTimes()

	clientProtocol := fake.NewMockClientProtocol(controller)
	clientProtocol.EXPECT().Dispense("plugin").Return(service, nil).AnyTimes()

	return &fakePluginClient{
		service:        service,
		clientProtocol: clientProtocol,
		name:           name,
	}
}

func (c *fakePluginClient) Client() (plugin.ClientProtocol, error) {
	return c.clientProtocol, nil
}

func (c *fakePluginClient) Kill() {}

type stubAPIService struct{}

var _ api.API = (*stubAPIService)(nil)

func (f *stubAPIService) Addr() string {
	return "localhost:54321"
}

func (f *stubAPIService) Start(context.Context) error {
	return nil
}
