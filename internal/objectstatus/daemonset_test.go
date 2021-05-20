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

func Test_daemonSet(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storefake.MockStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				objectFile := "daemonset_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusOK,
				Details:    []component.Component{component.NewText("Daemon Set is OK")},
				Properties: []component.Property{{Label: "Selectors", Value: component.NewSelectors([]component.Selector{component.NewLabelSelector("name", "fluentd")})}},
			},
		},
		{
			name: "misscheduled",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				objectFile := "daemonset_misscheduled.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("Daemon Set pods are running on nodes that aren't supposed to run Daemon Set pods")},
				Properties: []component.Property{{Label: "Selectors", Value: component.NewSelectors([]component.Selector{component.NewLabelSelector("name", "fluentd")})}},
			},
		},
		{
			name: "not ready",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				objectFile := "daemonset_not_ready.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("Daemon Set pods are not ready")},
				Properties: []component.Property{{Label: "Selectors", Value: component.NewSelectors([]component.Selector{component.NewLabelSelector("name", "fluentd")})}},
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
			status, err := daemonSet(ctx, object, o, linkInterface)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
