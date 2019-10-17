/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	storefake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_status(t *testing.T) {
	deployObjectStatus := ObjectStatus{
		nodeStatus: component.NodeStatusOK,
		Details:    []component.Component{component.NewText("apps/v1 Deployment is OK")},
	}

	lookup := statusLookup{
		{apiVersion: "v1", kind: "Object"}: func(context.Context, runtime.Object, store.Store) (ObjectStatus, error) {
			return deployObjectStatus, nil
		},
	}

	cases := []struct {
		name     string
		object   runtime.Object
		lookup   statusLookup
		expected ObjectStatus
		isErr    bool
	}{
		{
			name:     "in general",
			object:   testutil.CreateDeployment("deployment"),
			lookup:   lookup,
			expected: deployObjectStatus,
		},
		{
			name:   "nil object",
			object: nil,
			lookup: lookup,
			isErr:  true,
		},
		{
			name:   "nil lookup",
			object: testutil.CreateDeployment("deployment"),
			lookup: nil,
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockStore(controller)

			ctx := context.Background()
			got, err := status(ctx, tc.object, o, tc.lookup)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}

func Test_ObjectStatus_AddDetail(t *testing.T) {
	os := ObjectStatus{}
	os.AddDetail("detail")

	expected := []component.Component{component.NewText("detail")}
	assert.Equal(t, expected, os.Details)
}

func Test_ObjectStatus_AddDetailf(t *testing.T) {
	os := ObjectStatus{}
	os.AddDetailf("detail %d", 1)

	expected := []component.Component{component.NewText("detail 1")}
	assert.Equal(t, expected, os.Details)
}

func Test_ObjectStatus_SetError(t *testing.T) {
	os := ObjectStatus{}
	os.SetError()
	assert.Equal(t, component.NodeStatusError, os.Status())
}

func Test_ObjectStatus_SetWarning(t *testing.T) {
	os := ObjectStatus{}
	os.SetWarning()
	assert.Equal(t, component.NodeStatusWarning, os.Status())

	os.SetError()
	os.SetWarning()
	assert.Equal(t, component.NodeStatusError, os.Status())
}

func Test_ObjectStatus_Default(t *testing.T) {
	os := ObjectStatus{}

	expected := component.NodeStatusOK
	assert.Equal(t, expected, os.Status())
}
