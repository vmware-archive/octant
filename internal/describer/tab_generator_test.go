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

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

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
