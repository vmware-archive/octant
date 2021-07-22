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

func Test_pod(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T) runtime.Object {
				objectFile := "pod_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusOK,
				Details: []component.Component{
					component.NewText("Pod is OK"),
				},
				Properties: []component.Property{{Label: "ServiceAccount", Value: component.NewLink("", "ServiceAccount", "some-url/service-account")},
					{Label: "Node", Value: component.NewLink("", "Node", "some-url/node")},
					{Label: "Controlled By", Value: component.NewLink("", "ReplicaSet", "some-url/replica")}},
			},
		},
		{
			name: "pod is in unknown state",
			init: func(t *testing.T) runtime.Object {
				objectFile := "pod_unknown.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusError,
				Details: []component.Component{
					component.NewText("Pod is unhealthy"),
				},
				Properties: []component.Property{{Label: "ServiceAccount", Value: component.NewLink("", "ServiceAccount", "some-url/service-account")},
					{Label: "Node", Value: component.NewLink("", "Node", "some-url/node")},
					{Label: "Controlled By", Value: component.NewLink("", "ReplicaSet", "some-url/replica")}},
			},
		},
		{
			name: "pod is pending",
			init: func(t *testing.T) runtime.Object {
				objectFile := "pod_pending.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusWarning,
				Details: []component.Component{
					component.NewText("Pod may require additional action"),
				},
				Properties: []component.Property{{Label: "ServiceAccount", Value: component.NewLink("", "ServiceAccount", "some-url/service-account")},
					{Label: "Node", Value: component.NewLink("", "Node", "some-url/node")},
					{Label: "Controlled By", Value: component.NewLink("", "ReplicaSet", "some-url/replica")}},
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not a pod",
			init: func(t *testing.T) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
		{
			name: "pod has ephemeral containers",
			init: func(t *testing.T) runtime.Object {
				objectFile := "pod_ephemeral_container.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusWarning,
				Details: []component.Component{
					component.NewText("Pod is OK"),
					component.NewText("Ephemeral container is running"),
				},
				Properties: []component.Property{{Label: "ServiceAccount", Value: component.NewLink("", "ServiceAccount", "some-url/service-account")},
					{Label: "Node", Value: component.NewLink("", "Node", "some-url/node")},
					{Label: "Controlled By", Value: component.NewLink("", "ReplicaSet", "some-url/replica")}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			linkInterface := linkFake.NewMockInterface(controller)
			linkSA := component.NewLink("", "ServiceAccount", "some-url/service-account")
			linkNode := component.NewLink("", "Node", "some-url/node")
			linkReplica := component.NewLink("", "ReplicaSet", "some-url/replica")
			linkInterface.EXPECT().ForGVK(gomock.Any(), gomock.Any(), "ServiceAccount", gomock.Any(), gomock.Any()).Return(linkSA, nil).AnyTimes()
			linkInterface.EXPECT().ForGVK(gomock.Any(), gomock.Any(), "Node", gomock.Any(), gomock.Any()).Return(linkNode, nil).AnyTimes()
			linkInterface.EXPECT().ForGVK(gomock.Any(), gomock.Any(), "ReplicaSet", gomock.Any(), gomock.Any()).Return(linkReplica, nil).AnyTimes()
			defer controller.Finish()

			o := storefake.NewMockStore(controller)

			object := tc.init(t)

			ctx := context.Background()
			status, err := pod(ctx, object, o, linkInterface)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
