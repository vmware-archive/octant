/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/dashboard"
	"github.com/vmware/octant/pkg/plugin/fake"
	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
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

	fn(&mocks)
}

type grpcServerMocks struct {
	service       *fake.MockService
	moduleService *fake.MockModuleService
}

func (m *grpcServerMocks) genServer() *plugin.GRPCServer {
	return &plugin.GRPCServer{
		Impl: m.service,
	}
}

func (m *grpcServerMocks) genModuleServer() *plugin.GRPCServer {
	return &plugin.GRPCServer{
		Impl: m.moduleService,
	}
}

func testWithGRPCServer(t *testing.T, fn func(*grpcServerMocks)) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mocks := grpcServerMocks{
		service:       fake.NewMockService(controller),
		moduleService: fake.NewMockModuleService(controller),
	}

	fn(&mocks)
}

func Test_GRPCClient_Content(t *testing.T) {
	testWithGRPCClient(t, func(mocks *grpcClientMocks) {
		req := &dashboard.ContentRequest{
			Path: "/",
		}

		contentResponse := component.NewContentResponse(component.TitleFromString("title"))
		contentResponseBytes, err := json.Marshal(&contentResponse)
		require.NoError(t, err)

		resp := &dashboard.ContentResponse{
			ContentResponse: contentResponseBytes,
		}

		mocks.protoClient.EXPECT().
			Content(gomock.Any(), req).
			Return(resp, nil)

		client := mocks.genClient()
		ctx := context.Background()
		got, err := client.Content(ctx, "/")
		require.NoError(t, err)

		assert.Equal(t, *contentResponse, got)
	})
}

func Test_GRPCClient_Navigation(t *testing.T) {
	testWithGRPCClient(t, func(mocks *grpcClientMocks) {
		req := &dashboard.NavigationRequest{}

		resp := &dashboard.NavigationResponse{
			Navigation: &dashboard.NavigationResponse_Navigation{
				Title: "title",
			},
		}

		mocks.protoClient.EXPECT().
			Navigation(gomock.Any(), req).
			Return(resp, nil)

		client := mocks.genClient()
		ctx := context.Background()
		got, err := client.Navigation(ctx)
		require.NoError(t, err)

		expected := navigation.Navigation{
			Title: "title",
		}

		assert.Equal(t, expected, got)
	})
}

func Test_GRPCClient_Register(t *testing.T) {
	testWithGRPCClient(t, func(mocks *grpcClientMocks) {
		inGVKs := []*dashboard.RegisterResponse_GroupVersionKind{{Version: "v1", Kind: "Pod"}}

		resp := &dashboard.RegisterResponse{
			PluginName:  "my-plugin",
			Description: "description",
			Capabilities: &dashboard.RegisterResponse_Capabilities{
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
		ctx := context.Background()
		got, err := client.Register(ctx, apiAddress)
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
		objectRequest := &dashboard.ObjectRequest{
			Object: objectData,
		}

		config1 := component.NewText("config1 value")
		status1 := component.NewText("status1 value")

		printResponse := &dashboard.PrintResponse{
			Config: []*dashboard.PrintResponse_SummaryItem{
				{Header: "config1", Component: encodeComponent(t, config1)},
			},
			Status: []*dashboard.PrintResponse_SummaryItem{
				{Header: "status1", Component: encodeComponent(t, status1)},
			},
			Items: itemsData,
		}
		mocks.protoClient.EXPECT().Print(gomock.Any(), gomock.Eq(objectRequest)).Return(printResponse, nil)

		client := mocks.genClient()
		ctx := context.Background()
		got, err := client.Print(ctx, object)
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
		objectRequest := &dashboard.ObjectRequest{
			Object: objectData,
		}

		layout := flexlayout.New()
		section := layout.AddSection()
		err = section.Add(component.NewText("text"), component.WidthFull)
		require.NoError(t, err)

		tabResponse := &dashboard.PrintTabResponse{
			Name:   "tab name",
			Layout: encodeComponent(t, layout.ToComponent("component title")),
		}

		mocks.protoClient.EXPECT().
			PrintTab(gomock.Any(), gomock.Eq(objectRequest)).
			Return(tabResponse, nil)

		client := mocks.genClient()
		ctx := context.Background()
		got, err := client.PrintTab(ctx, object)
		require.NoError(t, err)

		expectedLayout := component.NewFlexLayout("component title")
		expectedLayout.AddSections(
			component.FlexLayoutSection{
				{
					Width: component.WidthFull,
					View:  component.NewText("text")},
			},
		)
		expected := plugin.TabResponse{
			Tab: &component.Tab{
				Name:     "tab name",
				Contents: *expectedLayout,
			},
		}

		testutil.AssertJSONEqual(t, expected, got)
	})
}

func Test_GRPCClient_ObjectStatus(t *testing.T) {
	testWithGRPCClient(t, func(mocks *grpcClientMocks) {
		object := testutil.CreatePod("pod")

		objectData, err := json.Marshal(object)
		require.NoError(t, err)
		objectRequest := &dashboard.ObjectRequest{
			Object: objectData,
		}

		gvk := object.GroupVersionKind()
		apiVersion, kind := gvk.ToAPIVersionAndKind()

		status := component.PodSummary{
			Status:  component.NodeStatusOK,
			Details: []component.Component{component.NewText(fmt.Sprintf("%s %s is OK", apiVersion, kind))},
		}

		statusData, err := json.Marshal(status)
		require.NoError(t, err)

		objectStatusResponse := &dashboard.ObjectStatusResponse{
			ObjectStatus: statusData,
		}

		mocks.protoClient.EXPECT().ObjectStatus(gomock.Any(), gomock.Eq(objectRequest)).Return(objectStatusResponse, nil)

		client := mocks.genClient()
		ctx := context.Background()
		got, err := client.ObjectStatus(ctx, object)
		require.NoError(t, err)

		expected := plugin.ObjectStatusResponse{
			ObjectStatus: status,
		}

		assert.Equal(t, expected, got)
	})
}

func Test_GRPCServer_Content(t *testing.T) {
	testWithGRPCServer(t, func(mocks *grpcServerMocks) {
		server := mocks.genModuleServer()

		contentResponse := component.NewContentResponse(component.TitleFromString("title"))

		mocks.moduleService.EXPECT().
			Content(gomock.Any(), "/").
			Return(*contentResponse, nil)

		ctx := context.Background()
		got, err := server.Content(ctx, &dashboard.ContentRequest{Path: "/"})
		require.NoError(t, err)

		contentResponseBytes, err := json.Marshal(contentResponse)

		expected := &dashboard.ContentResponse{
			ContentResponse: contentResponseBytes,
		}

		assert.Equal(t, expected, got)
	})
}

func Test_GRPCServer_Navigation(t *testing.T) {
	testWithGRPCServer(t, func(mocks *grpcServerMocks) {
		server := mocks.genModuleServer()

		mocks.moduleService.EXPECT().
			Navigation(gomock.Any()).
			Return(navigation.Navigation{Title: "title"}, nil)

		ctx := context.Background()
		got, err := server.Navigation(ctx, &dashboard.NavigationRequest{})
		require.NoError(t, err)

		expected := &dashboard.NavigationResponse{
			Navigation: &dashboard.NavigationResponse_Navigation{
				Title: "title",
			},
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

		mocks.service.EXPECT().Register(gomock.Any(), gomock.Eq(apiAddress)).Return(metadata, nil)

		server := mocks.genServer()

		ctx := context.Background()
		got, err := server.Register(ctx, &dashboard.RegisterRequest{
			DashboardAPIAddress: apiAddress,
		})
		require.NoError(t, err)

		outGVKs := []*dashboard.RegisterResponse_GroupVersionKind{{Version: "v1", Kind: "Pod"}}
		expected := &dashboard.RegisterResponse{
			PluginName:  "my-plugin",
			Description: "description",
			Capabilities: &dashboard.RegisterResponse_Capabilities{
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

		mocks.service.EXPECT().Print(gomock.Any(), gomock.Eq(u)).Return(pr, nil)

		objectData, err := json.Marshal(object)
		require.NoError(t, err)

		ctx := context.Background()
		objectRequest := &dashboard.ObjectRequest{
			Object: objectData,
		}

		server := mocks.genServer()
		got, err := server.Print(ctx, objectRequest)
		require.NoError(t, err)

		expectedItems, err := json.Marshal(pr.Items)
		require.NoError(t, err)

		expected := &dashboard.PrintResponse{
			Config: []*dashboard.PrintResponse_SummaryItem{
				{Header: "extra config", Component: encodeComponent(t, config)},
			},
			Status: []*dashboard.PrintResponse_SummaryItem{
				{Header: "extra status", Component: encodeComponent(t, status)},
			},
			Items: expectedItems,
		}
		assert.Equal(t, expected, got)

	})
}

func Test_GRPCServer_ObjectStatus(t *testing.T) {
	testWithGRPCServer(t, func(mocks *grpcServerMocks) {
		object := testutil.CreatePod("pod")
		gvk := object.GroupVersionKind()
		apiVersion, kind := gvk.ToAPIVersionAndKind()

		status := component.PodSummary{
			Status:  component.NodeStatusOK,
			Details: []component.Component{component.NewText(fmt.Sprintf("%s %s is OK", apiVersion, kind))},
		}

		osr := plugin.ObjectStatusResponse{
			ObjectStatus: status,
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		require.NoError(t, err)
		u := &unstructured.Unstructured{Object: m}

		mocks.service.EXPECT().ObjectStatus(gomock.Any(), gomock.Eq(u)).Return(osr, nil)

		objectData, err := json.Marshal(object)
		require.NoError(t, err)
		objectRequest := &dashboard.ObjectRequest{
			Object: objectData,
		}

		ctx := context.Background()

		server := mocks.genServer()
		got, err := server.ObjectStatus(ctx, objectRequest)
		require.NoError(t, err)

		encodedStatus, err := json.Marshal(status)
		require.NoError(t, err)

		expected := &dashboard.ObjectStatusResponse{
			ObjectStatus: encodedStatus,
		}

		assert.Equal(t, expected, got)
	})
}

func encodeComponent(t *testing.T, view component.Component) []byte {
	data, err := json.Marshal(view)
	require.NoError(t, err)
	return data
}
