/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"testing"

	linkFake "github.com/vmware-tanzu/octant/internal/link/fake"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/testutil"
	storefake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_replicaSetAppsV1(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storefake.MockStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				objectFile := "replicaset_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusOK,
				Details:    []component.Component{component.NewText("Replica Set is OK")},
				Properties: []component.Property{{Label: "Current Replicas", Value: component.NewText("1")},
					{Label: "Desired Replicas", Value: component.NewText("1")}},
			},
		},
		{
			name: "no replicas",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				objectFile := "replicaset_no_replicas.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusError,
				Details:    []component.Component{component.NewText("Replica Set has no replicas available")},
				Properties: []component.Property{{Label: "Current Replicas", Value: component.NewText("1")},
					{Label: "Desired Replicas", Value: component.NewText("1")}},
			},
		},
		{
			name: "zero replicas",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				objectFile := "replicaset_zero_replicas.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusOK,
				Details:    []component.Component{component.NewText("Replica Set is OK")},
				Properties: []component.Property{{Label: "Current Replicas", Value: component.NewText("0")},
					{Label: "Desired Replicas", Value: component.NewText("0")}},
			},
		},
		{
			name: "not available",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				objectFile := "replicaset_not_available.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("Expected 1 replicas, but 0 are available")},
				Properties: []component.Property{{Label: "Current Replicas", Value: component.NewText("1")},
					{Label: "Desired Replicas", Value: component.NewText("1")}},
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not a daemon set",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			linkInterface := linkFake.NewMockInterface(controller)
			defer controller.Finish()

			o := storefake.NewMockStore(controller)

			object := tc.init(t, o)

			ctx := context.Background()
			status, err := replicaSetAppsV1(ctx, object, o, linkInterface)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
