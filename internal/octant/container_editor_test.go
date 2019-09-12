/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/store/fake"
)

func TestNewContainerEditor(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	objectStore := fake.NewMockStore(controller)

	key := store.Key{
		Namespace:  "default",
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "deployment",
	}

	objectStore.EXPECT().
		Update(gomock.Any(), key, gomock.Any()).
		Return(nil)

	editor := NewContainerEditor(objectStore)

	ctx := context.Background()

	payload := action.CreatePayload("overview/containerEditor", map[string]interface{}{
		"apiVersion":     "apps/v1",
		"kind":           "Deployment",
		"namespace":      "default",
		"name":           "deployment",
		"containersPath": `["foo", "bar"]`,
		"containerName":  "nginx",
		"containerImage": "nginx:stable",
	})

	require.NoError(t, editor.Handle(ctx, payload))
}

func Test_updateContainer(t *testing.T) {
	deployment := testutil.ToUnstructured(
		t,
		testutil.CreateDeployment("deployment", testutil.WithGenericDeployment()))

	containersPath := []string{"spec", "template", "spec", "containers"}
	containerName := "container-name"
	containerImage := "new-image"

	fn := updateContainer(containersPath, log.NopLogger(), containerName, containerImage)

	require.NoError(t, fn(deployment))

	got, found, err := unstructured.NestedSlice(deployment.Object, containersPath...)

	require.NoError(t, err)
	require.True(t, found)

	expected := []interface{}{
		map[string]interface{}{
			"image":     "new-image",
			"name":      "container-name",
			"resources": map[string]interface{}{},
		},
	}
	require.Equal(t, expected, got)

}
