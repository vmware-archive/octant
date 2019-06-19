/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	configFake "github.com/heptio/developer-dash/internal/config/fake"
	"github.com/heptio/developer-dash/internal/describer"
	dashplugin "github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/plugin/fake"
	pluginFake "github.com/heptio/developer-dash/pkg/plugin/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func TestPluginDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	name := "plugin-test"
	namespace := "default"
	metadata := dashplugin.Metadata{
		Name:         name,
		Description:  "this is a test",
		Capabilities: dashplugin.Capabilities{},
	}

	store := dashplugin.NewDefaultStore()
	client := newFakePluginClient(name, controller)
	require.NoError(t, store.Store(name, client, metadata))

	pluginManager := pluginFake.NewMockManagerInterface(controller)
	pluginManager.EXPECT().Store().Return(store).AnyTimes()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().PluginManager().Return(pluginManager)

	p := NewPluginListDescriber()

	options := describer.Options{
		Dash: dashConfig,
	}

	ctx := context.Background()
	cResponse, err := p.Describe(ctx, "/plugins", namespace, options)
	require.NoError(t, err)

	list := component.NewList("Plugins", nil)
	tableCols := component.NewTableCols("Name", "Description", "Capabilities")
	table := component.NewTable("Plugins", tableCols)
	table.Add(component.TableRow{
		"Name":        component.NewText(name),
		"Description": component.NewText("this is a test"),
		"Capability":  component.NewText("{\"SupportsPrinterConfig\":null,\"SupportsPrinterStatus\":null,\"SupportsPrinterItems\":null,\"SupportsObjectStatus\":null,\"SupportsTab\":null}"),
	})

	list.Add(table)

	expected := component.ContentResponse{
		Components: []component.Component{list},
	}

	assert.Equal(t, expected, cResponse)
}

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

type fakePluginClient struct {
	clientProtocol *fake.MockClientProtocol
	service        *fake.MockService
	name           string
}

var _ dashplugin.Client = (*fakePluginClient)(nil)

func (c *fakePluginClient) Client() (plugin.ClientProtocol, error) {
	return c.clientProtocol, nil
}

func (c *fakePluginClient) Kill() {}
