/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/link"
	linkFake "github.com/vmware-tanzu/octant/internal/link/fake"

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
	deployment := testutil.CreateDeployment("deployment")
	deployObjectStatus := ObjectStatus{
		NodeStatus: component.NodeStatusOK,
		Details:    []component.Component{component.NewText("apps/v1 Deployment is OK")},
		Properties: []component.Property{{Label: "Namespace", Value: component.NewText("namespace")},
			{Label: "Created", Value: component.NewTimestamp(deployment.CreationTimestamp.Time)}},
	}

	lookup := statusLookup{
		{apiVersion: "v1", kind: "Object"}: func(context.Context, runtime.Object, store.Store, link.Interface) (ObjectStatus, error) {
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
			object:   deployment,
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
			linkInterface := linkFake.NewMockInterface(controller)
			defer controller.Finish()

			o := storefake.NewMockStore(controller)
			o.EXPECT().List(gomock.Any(), gomock.Any()).Return(&unstructured.UnstructuredList{}, false, nil).AnyTimes()

			ctx := context.Background()
			got, err := status(ctx, tc.object, o, tc.lookup, linkInterface, 5)
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

func Test_getDeletedObjectStatus_StatusError(t *testing.T) {
	now := time.Now()
	timeDeletion := &metav1.Time{Time: now.Add(-6 * time.Minute)}

	given, _ := getDeletedObjectStatus(timeDeletion, int64(5), "service")

	e := ObjectStatus{
		NodeStatus: component.NodeStatusError,
		Details: []component.Component{
			component.NewTextf("%s %s", "service", "was been deleted for longer than the given finalize threshold"),
		},
		Properties: []component.Property{
			{
				Label: "Deleted Date",
				Value: component.NewTimestamp(timeDeletion.Time),
			},
		},
	}

	assert.Equal(t, e.NodeStatus, given.NodeStatus)

	ej, _ := e.Details[0].MarshalJSON()
	gj, _ := given.Details[0].MarshalJSON()

	assert.Equal(t, ej, gj)
	assert.Equal(t, e.NodeStatus, given.NodeStatus)
}

func Test_getDeletedObjectStatus_StatusWarning(t *testing.T) {
	timeDeletion := &metav1.Time{Time: time.Now()}

	given, _ := getDeletedObjectStatus(timeDeletion, int64(5), "service")

	e := ObjectStatus{
		NodeStatus: component.NodeStatusWarning,
		Details: []component.Component{
			component.NewTextf("%s %s", "service", "is being deleted"),
		},
		Properties: []component.Property{
			{
				Label: "Deleted Date",
				Value: component.NewTimestamp(timeDeletion.Time),
			},
		},
	}

	assert.Equal(t, e.NodeStatus, given.NodeStatus)

	ej, _ := e.Details[0].MarshalJSON()
	gj, _ := given.Details[0].MarshalJSON()

	assert.Equal(t, ej, gj)
	assert.Equal(t, e.NodeStatus, given.NodeStatus)
}
