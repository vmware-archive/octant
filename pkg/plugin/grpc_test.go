package plugin_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/plugin/fake"
	"github.com/heptio/developer-dash/pkg/plugin/proto"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type grpcClientMocks struct {
	broker      *fake.MockBroker
	protoClient *fake.MockPluginClient
}

func (m *grpcClientMocks) genClient() *plugin.GRPCClient {
	return plugin.NewGRPCClient(m.broker, m.protoClient)
}

func testWithGRPCClient(t *testing.T, fn func(*grpcClientMocks)) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mocks := grpcClientMocks{
		broker:      fake.NewMockBroker(controller),
		protoClient: fake.NewMockPluginClient(controller),
	}

	mocks.broker.EXPECT().NextId().Return(uint32(1))
	mocks.broker.EXPECT().AcceptAndServe(gomock.Eq(uint32(1)), gomock.Any()).AnyTimes()

	fn(&mocks)
}

type grpcServerMocks struct {
	service *fake.MockService
}

func (m *grpcServerMocks) genServer() *plugin.GRPCServer {
	return &plugin.GRPCServer{
		Impl: m.service,
	}
}

func testWithGRPCServer(t *testing.T, fn func(*grpcServerMocks)) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mocks := grpcServerMocks{
		service: fake.NewMockService(controller),
	}

	fn(&mocks)
}

func Test_GRPCClient_Register(t *testing.T) {
	testWithGRPCClient(t, func(mocks *grpcClientMocks) {
		inGVKs := []*proto.RegisterResponse_GroupVersionKind{{Version: "v1", Kind: "Pod"}}

		resp := &proto.RegisterResponse{
			PluginName:  "my-plugin",
			Description: "description",
			Capabilities: &proto.RegisterResponse_Capabilities{
				SupportsPrinterConfig: inGVKs,
				SupportsPrinterStatus: inGVKs,
				SupportsPrinterItems:  inGVKs,
				SupportsObjectStatus:  inGVKs,
				SupportsTab:           inGVKs,
			},
		}

		mocks.protoClient.EXPECT().Register(gomock.Any(), gomock.Any()).Return(resp, nil)

		client := mocks.genClient()
		apiAddress := "localhost:54321"
		got, err := client.Register(apiAddress)
		require.NoError(t, err)

		outGVKs := []schema.GroupVersionKind{{Version: "v1", Kind: "Pod"}}

		expected := plugin.Metadata{
			Name:        "my-plugin",
			Description: "description",
			Capabilities: plugin.Capabilities{
				SupportsPrinterConfig: outGVKs,
				SupportsPrinterStatus: outGVKs,
				SupportsPrinterItems:  outGVKs,
				SupportsObjectStatus:  outGVKs,
				SupportsTab:           outGVKs,
			},
		}
		assert.Equal(t, expected, got)
	})
}

func Test_GRPCClient_Print(t *testing.T) {
	testWithGRPCClient(t, func(mocks *grpcClientMocks) {
		object := testutil.CreateDeployment("deployment")

		items := component.FlexLayoutSection{
			{Width: component.WidthFull, View: component.NewText("section 1")},
		}
		itemsData, err := json.Marshal(items)
		require.NoError(t, err)

		objectData, err := json.Marshal(object)
		require.NoError(t, err)
		objectRequest := &proto.ObjectRequest{
			Object: objectData,
		}

		config1 := component.NewText("config1 value")
		status1 := component.NewText("status1 value")

		printResponse := &proto.PrintResponse{
			Config: []*proto.PrintResponse_SummaryItem{
				{Header: "config1", Component: encodeComponent(t, config1)},
			},
			Status: []*proto.PrintResponse_SummaryItem{
				{Header: "status1", Component: encodeComponent(t, status1)},
			},
			Items: itemsData,
		}
		mocks.protoClient.EXPECT().Print(gomock.Any(), gomock.Eq(objectRequest)).Return(printResponse, nil)

		client := mocks.genClient()
		got, err := client.Print(object)
		require.NoError(t, err)

		expected := plugin.PrintResponse{
			Config: []component.SummarySection{
				{Header: "config1", Content: config1},
			},
			Status: []component.SummarySection{
				{Header: "status1", Content: status1},
			},
			Items: component.FlexLayoutSection{
				{Width: component.WidthFull, View: component.NewText("section 1")},
			},
		}

		assert.Equal(t, expected, got)
	})
}

func Test_GRPCClient_PrintTab(t *testing.T) {
	testWithGRPCClient(t, func(mocks *grpcClientMocks) {
		object := testutil.CreateDeployment("deployment")

		objectData, err := json.Marshal(object)
		require.NoError(t, err)
		objectRequest := &proto.ObjectRequest{
			Object: objectData,
		}

		layout := flexlayout.New()
		section := layout.AddSection()
		err = section.Add(component.NewText("text"), component.WidthFull)
		require.NoError(t, err)

		tabResponse := &proto.PrintTabResponse{
			Name:   "tab name",
			Layout: encodeComponent(t, layout.ToComponent("component title")),
		}

		mocks.protoClient.EXPECT().
			PrintTab(gomock.Any(), gomock.Eq(objectRequest)).
			Return(tabResponse, nil)

		client := mocks.genClient()
		got, err := client.PrintTab(object)
		require.NoError(t, err)

		expectedLayout := component.NewFlexLayout("component title")
		expectedLayout.AddSections(
			component.FlexLayoutSection{
				{
					Width: component.WidthFull,
					View:  component.NewText("text")},
			},
		)
		expected := &component.Tab{
			Name:     "tab name",
			Contents: *expectedLayout,
		}

		assert.Equal(t, expected, got)
	})
}

func Test_GRPCServer_Register(t *testing.T) {
	inGVKs := []schema.GroupVersionKind{{Version: "v1", Kind: "Pod"}}

	testWithGRPCServer(t, func(mocks *grpcServerMocks) {
		metadata := plugin.Metadata{
			Name:        "my-plugin",
			Description: "description",
			Capabilities: plugin.Capabilities{
				SupportsPrinterConfig: inGVKs,
				SupportsPrinterStatus: inGVKs,
				SupportsPrinterItems:  inGVKs,
				SupportsObjectStatus:  inGVKs,
				SupportsTab:           inGVKs,
			},
		}

		apiAddress := "localhost:54321"

		mocks.service.EXPECT().Register(gomock.Eq(apiAddress)).Return(metadata, nil)

		server := mocks.genServer()

		ctx := context.Background()
		got, err := server.Register(ctx, &proto.RegisterRequest{
			DashboardAPIAddress: apiAddress,
		})
		require.NoError(t, err)

		outGVKs := []*proto.RegisterResponse_GroupVersionKind{{Version: "v1", Kind: "Pod"}}
		expected := &proto.RegisterResponse{
			PluginName:  "my-plugin",
			Description: "description",
			Capabilities: &proto.RegisterResponse_Capabilities{
				SupportsPrinterConfig: outGVKs,
				SupportsPrinterStatus: outGVKs,
				SupportsPrinterItems:  outGVKs,
				SupportsObjectStatus:  outGVKs,
				SupportsTab:           outGVKs,
			},
		}

		assert.Equal(t, expected, got)
	})
}

func Test_GRPCServer_Print(t *testing.T) {
	testWithGRPCServer(t, func(mocks *grpcServerMocks) {
		object := testutil.CreateDeployment("deployment")

		config := component.NewText("config")
		status := component.NewText("config")

		pr := plugin.PrintResponse{
			Config: []component.SummarySection{
				{Header: "extra config", Content: config},
			},
			Status: []component.SummarySection{
				{Header: "extra status", Content: status},
			},
			Items: []component.FlexLayoutItem{
				{Width: 24, View: component.NewText("item1")},
			},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		require.NoError(t, err)
		u := &unstructured.Unstructured{Object: m}

		mocks.service.EXPECT().Print(gomock.Eq(u)).Return(pr, nil)

		objectData, err := json.Marshal(object)
		require.NoError(t, err)

		ctx := context.Background()
		objectRequest := &proto.ObjectRequest{
			Object: objectData,
		}

		server := mocks.genServer()
		got, err := server.Print(ctx, objectRequest)
		require.NoError(t, err)

		expectedItems, err := json.Marshal(pr.Items)
		require.NoError(t, err)

		expected := &proto.PrintResponse{
			Config: []*proto.PrintResponse_SummaryItem{
				{Header: "extra config", Component: encodeComponent(t, config)},
			},
			Status: []*proto.PrintResponse_SummaryItem{
				{Header: "extra status", Component: encodeComponent(t, status)},
			},
			Items: expectedItems,
		}
		assert.Equal(t, expected, got)

	})
}

func encodeComponent(t *testing.T, view component.Component) []byte {
	data, err := json.Marshal(view)
	require.NoError(t, err)
	return data
}
