/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/octant"
	printerFake "github.com/vmware-tanzu/octant/internal/printer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestObjectDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	thePath := "/"

	pod := testutil.CreatePod("pod")
	pod.CreationTimestamp = *testutil.CreateTimestamp()

	key, err := store.KeyFromObject(pod)
	require.NoError(t, err)

	dashConfig := configFake.NewMockDash(controller)
	moduleRegistrar := pluginFake.NewMockModuleRegistrar(controller)
	actionRegistrar := pluginFake.NewMockActionRegistrar(controller)

	pluginManager := plugin.NewManager(nil, moduleRegistrar, actionRegistrar)
	dashConfig.EXPECT().PluginManager().Return(pluginManager).AnyTimes()

	objectPrinter := printerFake.NewMockPrinter(controller)

	podSummary := component.NewText("summary")
	objectPrinter.EXPECT().Print(gomock.Any(), pod, pluginManager).Return(podSummary, nil)

	options := Options{
		Dash:    dashConfig,
		Printer: objectPrinter,
		LoadObject: func(ctx context.Context, namespace string, fields map[string]string, objectStoreKey store.Key) (*unstructured.Unstructured, error) {
			return testutil.ToUnstructured(t, pod), nil
		},
	}

	objectConfig := ObjectConfig{
		Path:                  thePath,
		BaseTitle:             "object",
		StoreKey:              key,
		ObjectType:            podObjectType,
		DisableResourceViewer: true,
		IconName:              "icon-name",
		IconSource:            "icon-source",
	}
	d := NewObject(objectConfig)

	d.tabFuncDescriptors = []tabFuncDescriptor{
		{name: "summary", tabFunc: d.addSummaryTab},
	}

	cResponse, err := d.Describe(ctx, pod.Namespace, options)
	require.NoError(t, err)

	summary := component.NewText("summary")
	summary.SetAccessor("summary")

	buttonGroup := component.NewButtonGroup()

	buttonGroup.AddButton(
		component.NewButton("Delete",
			action.CreatePayload(octant.ActionDeleteObject, key.ToActionPayload()),
			component.WithButtonConfirmation(
				"Delete Pod",
				"Are you sure you want to delete *Pod* **pod**? This action is permanent and cannot be recovered.",
			)))

	expected := component.ContentResponse{
		Title:      component.Title(component.NewText("object"), component.NewText("pod")),
		IconName:   "icon-name",
		IconSource: "icon-source",
		Components: []component.Component{
			summary,
		},
		ButtonGroup: buttonGroup,
	}
	assert.Equal(t, expected, cResponse)

}

func Test_deleteObjectConfirmation(t *testing.T) {
	pod := testutil.CreatePod("pod")
	option, err := deleteObjectConfirmation(pod)
	require.NoError(t, err)

	button := component.Button{}
	option(&button)

	expected := component.Button{
		Confirmation: &component.Confirmation{
			Title: "Delete Pod",
			Body:  "Are you sure you want to delete *Pod* **pod**? This action is permanent and cannot be recovered.",
		},
	}

	assert.Equal(t, expected, button)
}
