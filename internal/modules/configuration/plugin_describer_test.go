/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/gvk"
	dashPlugin "github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/fake"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestPluginDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	name := "plugin-test"
	namespace := "default"
	metadata := &dashPlugin.Metadata{
		Name:        name,
		Description: "this is a test",
		Capabilities: dashPlugin.Capabilities{
			SupportsPrinterConfig: []schema.GroupVersionKind{gvk.Pod},
			SupportsPrinterStatus: []schema.GroupVersionKind{gvk.Pod},
			SupportsPrinterItems:  []schema.GroupVersionKind{gvk.Pod},
			SupportsObjectStatus:  []schema.GroupVersionKind{gvk.Pod},
			SupportsTab:           []schema.GroupVersionKind{gvk.Pod},
			IsModule:              true,
			ActionNames:           []string{"action"},
		},
	}

	store := dashPlugin.NewDefaultStore()
	client := newFakePluginClient(name, controller)
	require.NoError(t, store.Store(name, client, metadata, "cmd"))

	pluginManager := pluginFake.NewMockManagerInterface(controller)
	pluginManager.EXPECT().Store().Return(store).AnyTimes()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().PluginManager().Return(pluginManager)

	p := NewPluginListDescriber()

	options := describer.Options{
		Dash: dashConfig,
	}

	ctx := context.Background()
	cResponse, err := p.Describe(ctx, namespace, options)
	require.NoError(t, err)

	capabilitiesData := "[Module], [Actions: action], [Object Status: v1 Pod], [Printer Config: v1 Pod], [Printer Items: v1 Pod], [Printer Status: v1 Pod], [Tab: v1 Pod]"

	list := component.NewList("Plugins", nil)
	tableCols := component.NewTableCols("Name", "Description", "Capabilities")
	table := component.NewTable("Plugins", "There are no plugins!", tableCols)
	table.Add(component.TableRow{
		"Name":         component.NewText(name),
		"Description":  component.NewText("this is a test"),
		"Capabilities": component.NewText(capabilitiesData),
	})

	list.Add(table)

	require.Len(t, cResponse.Components, 1)
	component.AssertEqual(t, list, cResponse.Components[0])
}

func newFakePluginClient(name string, controller *gomock.Controller) *fakePluginClient {
	service := fake.NewMockService(controller)
	metadata := dashPlugin.Metadata{
		Name: name,
	}
	service.EXPECT().Register(gomock.Any(), gomock.Eq("localhost:54321")).Return(metadata, nil).AnyTimes()

	clientProtocol := fake.NewMockClientProtocol(controller)
	clientProtocol.EXPECT().Dispense("plugin").Return(service, nil).AnyTimes()

	return &fakePluginClient{
		service:        service,
		clientProtocol: clientProtocol,
		name:           name,
	}
}

type fakePluginClient struct {
	clientProtocol *fake.MockClientProtocol
	service        *fake.MockService
	name           string
}

var _ dashPlugin.Client = (*fakePluginClient)(nil)

func (c *fakePluginClient) Client() (plugin.ClientProtocol, error) {
	return c.clientProtocol, nil
}

func (c *fakePluginClient) Kill() {}
