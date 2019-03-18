package plugin_test

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/gvk"
	"github.com/heptio/developer-dash/internal/view/component"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/testutil"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-plugin"
	dashplugin "github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/plugin/fake"
	"github.com/stretchr/testify/require"
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

	name := "plugin1"

	clientFactory := fake.NewMockClientFactory(controller)
	client := newFakePluginClient(name, controller)
	clientFactory.EXPECT().Init(gomock.Eq(name)).Return(client)

	metadata := dashplugin.Metadata{
		Name: name,
	}
	store.EXPECT().Store(gomock.Eq(name), gomock.Eq(client), gomock.Eq(metadata))
	store.EXPECT().Clients().Return(map[string]dashplugin.Client{name: client})

	options = append(options, func(m *dashplugin.Manager) {
		m.Store = store
		m.ClientFactory = clientFactory
	})

	manager := dashplugin.NewManager(options...)

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

	name1 := "plugin1"
	name2 := "plugin2"
	pod := testutil.CreatePod("pod")

	var options []dashplugin.ManagerOption

	store := fake.NewMockManagerStore(controller)
	clientFactory := fake.NewMockClientFactory(controller)

	client1 := newFakePluginClient(name1, controller)
	client2 := newFakePluginClient(name2, controller)
	clients := map[string]dashplugin.Client{
		name1: client1,
		name2: client2,
	}
	store.EXPECT().Clients().Return(clients)

	setupClient := func(client *fakePluginClient) {
		pr := dashplugin.PrintResponse{
			Config: []component.SummarySection{
				{Header: client.name},
			},
		}

		client.service.EXPECT().Print(gomock.Eq(pod)).
			Return(pr, nil).AnyTimes()
		store.EXPECT().GetMetadata(gomock.Eq(client.name)).
			Return(client.metadata(), nil).AnyTimes()
		store.EXPECT().GetService(gomock.Eq(client.name)).
			Return(client.service, nil).AnyTimes()
		client.service.EXPECT().Print(pod).Return(pr, nil).AnyTimes()
	}

	setupClient(client1)
	setupClient(client2)

	options = append(options, func(m *dashplugin.Manager) {
		m.Store = store
		m.ClientFactory = clientFactory
	})

	manager := dashplugin.NewManager(options...)

	_, err := manager.Print(pod)
	require.NoError(t, err)

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
	service.EXPECT().Register().Return(metadata, nil).AnyTimes()

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

func (c *fakePluginClient) metadata() dashplugin.Metadata {
	return dashplugin.Metadata{
		Name: c.name,
		Capabilities: dashplugin.Capabilities{
			SupportsPrinterConfig: []schema.GroupVersionKind{gvk.PodGVK},
		},
	}
}
