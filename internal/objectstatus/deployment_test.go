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
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_deploymentAppsV1(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storeFake.MockStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				objectFile := "deployment_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusOK,
				Details:    []component.Component{component.NewText("Deployment is OK")},
				Properties: []component.Property{{Label: "Deployment Strategy", Value: component.NewText("RollingUpdate")},
					{Label: "Selectors", Value: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "hello-node")})}},
			},
		},
		{
			name: "no replicas",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				objectFile := "deployment_no_replicas.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusError,
				Details:    []component.Component{component.NewText("No replicas exist for this deployment")},
				Properties: []component.Property{{Label: "Deployment Strategy", Value: component.NewText("RollingUpdate")},
					{Label: "Selectors", Value: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "hello-node")})}},
			},
		},
		{
			name: "not available",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				objectFile := "deployment_not_available.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("Expected 1 replicas, but 0 are available")},
				Properties: []component.Property{{Label: "Deployment Strategy", Value: component.NewText("RollingUpdate")},
					{Label: "Selectors", Value: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "hello-node")})}},
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not a daemon set",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
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

			o := storeFake.NewMockStore(controller)

			object := tc.init(t, o)

			ctx := context.Background()
			status, err := deploymentAppsV1(ctx, object, o, linkInterface)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
