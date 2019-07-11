/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/vmware/octant/internal/config/fake"
	printerFake "github.com/vmware/octant/internal/modules/overview/printer/fake"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/plugin"
	pluginFake "github.com/vmware/octant/pkg/plugin/fake"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

func TestObjectDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	thePath := "/"

	pod := testutil.CreatePod("pod")
	pod.CreationTimestamp = metav1.Time{
		Time: time.Unix(1547472896, 0),
	}

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

	cResponse, err := d.Describe(ctx, "/path", pod.Namespace, options)
	require.NoError(t, err)

	summary := component.NewText("summary")
	summary.SetAccessor("summary")

	expected := component.ContentResponse{
		Title:      component.Title(component.NewText("object"), component.NewText("pod")),
		IconName:   "icon-name",
		IconSource: "icon-source",
		Components: []component.Component{
			summary,
		},
	}
	assert.Equal(t, expected, cResponse)

}
