/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/store"
	storefake "github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

func mockObjectsEvents(t *testing.T, appObjectStore *storefake.MockStore, namespace string, events ...corev1.Event) {
	require.NotNil(t, appObjectStore)

	objects := &unstructured.UnstructuredList{}

	for _, event := range events {
		objects.Items = append(objects.Items, *testutil.ToUnstructured(t, &event))
	}

	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Event",
	}

	appObjectStore.EXPECT().
		List(gomock.Any(), key).
		Return(objects, false, nil).
		AnyTimes()
}

func stubMetadataForObject(t *testing.T, object runtime.Object, fl *flexlayout.FlexLayout) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	metadata, err := NewMetadata(object, tpo.link)
	require.NoError(t, err)
	require.NoError(t, metadata.AddToFlexLayout(fl))
}
