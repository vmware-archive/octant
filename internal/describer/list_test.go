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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	printerFake "github.com/vmware-tanzu/octant/internal/printer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestListDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	thePath := "/"

	pod := testutil.CreatePod("pod")
	pod.CreationTimestamp = *testutil.CreateTimestamp()

	key, err := store.KeyFromObject(pod)
	require.NoError(t, err)

	ctx := context.Background()
	namespace := "default"

	dashConfig := configFake.NewMockDash(controller)

	podListTable := createPodTable(*pod)

	objectPrinter := printerFake.NewMockPrinter(controller)
	podList := &corev1.PodList{Items: []corev1.Pod{*pod}}
	objectPrinter.EXPECT().Print(gomock.Any(), podList).Return(podListTable, nil)

	options := Options{
		Dash:    dashConfig,
		Printer: objectPrinter,
		LoadObjects: func(ctx context.Context, namespace string, fields map[string]string, objectStoreKeys []store.Key) (*unstructured.UnstructuredList, error) {
			return testutil.ToUnstructuredList(t, pod), nil
		},
	}

	listConfig := ListConfig{
		Path:          thePath,
		Title:         "list",
		StoreKey:      key,
		ListType:      PodListType,
		ObjectType:    PodObjectType,
		IsClusterWide: false,
	}
	d := NewList(listConfig)
	cResponse, err := d.Describe(ctx, namespace, options)
	require.NoError(t, err)

	list := component.NewList(append([]component.TitleComponent{}, component.NewText("list")), nil)
	list.Add(podListTable)
	expected := component.ContentResponse{
		Components: []component.Component{list},
	}

	assert.Equal(t, expected.Title, cResponse.Title)
}
