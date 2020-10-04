/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package describer

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestObjectTabsGenerator_Generate(t *testing.T) {
	g := NewObjectTabsGenerator()

	c := component.NewText("text")
	c.SetAccessor("accessor")

	object := testutil.CreatePod("pod")
	tabsFactory := func() ([]Tab, error) {
		return []Tab{
			{
				Name: "tab",
				Factory: func(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
					return c, nil
				},
			},
		}, nil
	}

	config := TabsGeneratorConfig{
		Object:      object,
		TabsFactory: tabsFactory,
		Options:     Options{},
	}

	ctx := context.Background()
	actual, err := g.Generate(ctx, config)
	require.NoError(t, err)

	wanted := []component.Component{c}
	testutil.AssertJSONEqual(t, wanted, actual)
}

func TestCreateErrorTab(t *testing.T) {
	actual := CreateErrorTab("Name", fmt.Errorf("error"))
	wanted := component.NewError(component.TitleFromString("Name"), fmt.Errorf("error"))
	wanted.SetAccessor("Name")

	testutil.AssertJSONEqual(t, wanted, actual)
}

func Test_pluginTabFactory(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pod := testutil.CreatePod("pod")
	g := NewObjectTabsGenerator()

	ctx := context.Background()
	dashConfig := configFake.NewMockDash(controller)
	pluginManager := pluginFake.NewMockManagerInterface(controller)
	dashConfig.EXPECT().PluginManager().Return(pluginManager).AnyTimes()
	options := Options{
		Dash: dashConfig,
	}

	tabs := []component.Tab{
		{
			Name:     "foo",
			Contents: *component.NewFlexLayout("foo"),
		},
		{
			Name:     "bar",
			Contents: *component.NewFlexLayout("bar"),
		},
		{
			Name:     "baz",
			Contents: *component.NewFlexLayout("baz"),
		},
	}

	pluginManager.EXPECT().Tabs(ctx, pod).Return(tabs, nil)

	tabsFactory, err := pluginTabsFactory(ctx, pod, options)
	require.NoError(t, err)

	config := TabsGeneratorConfig{
		Object:      pod,
		TabsFactory: func() ([]Tab, error) { return tabsFactory, nil },
		Options:     options,
	}

	actual, err := g.Generate(ctx, config)
	require.NoError(t, err)

	test := []component.Component{
		component.NewFlexLayout("foo"),
		component.NewFlexLayout("bar"),
		component.NewFlexLayout("baz"),
	}
	testutil.AssertJSONEqual(t, test, actual)
}
