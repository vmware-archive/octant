/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	storefake "github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
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
			defer controller.Finish()

			o := storefake.NewMockStore(controller)

			object := tc.init(t, o)

			ctx := context.Background()
			status, err := daemonSet(ctx, object, o)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
