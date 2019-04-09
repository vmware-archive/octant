package objectstatus

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_replicaSetAppsV1(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storefake.MockObjectStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				objectFile := "replicaset_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusOK,
				Details:    component.TitleFromString("Replica Set is OK"),
			},
		},
		{
			name: "no replicas",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				objectFile := "replicaset_no_replicas.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusError,
				Details:    component.TitleFromString("Replica Set has no replicas available"),
			},
		},
		{
			name: "not available",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				objectFile := "replicaset_not_available.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    component.TitleFromString("Expected 1 replicas, but 0 are available"),
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not a daemon set",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)

			object := tc.init(t, o)

			ctx := context.Background()
			status, err := replicaSetAppsV1(ctx, object, o)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}

func Test_replicaSetExtV1Beta1(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storefake.MockObjectStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				objectFile := "replicaset_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusOK,
				Details:    component.TitleFromString("Replica Set is OK"),
			},
		},
		{
			name: "no replicas",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				objectFile := "replicaset_no_replicas.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusError,
				Details:    component.TitleFromString("Replica Set has no replicas available"),
			},
		},
		{
			name: "not available",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				objectFile := "replicaset_not_available.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    component.TitleFromString("Expected 1 replicas, but 0 are available"),
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not a daemon set",
			init: func(t *testing.T, o *storefake.MockObjectStore) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)

			object := tc.init(t, o)

			ctx := context.Background()
			status, err := replicaSetExtV1Beta1(ctx, object, o)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
