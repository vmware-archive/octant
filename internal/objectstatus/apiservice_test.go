/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
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
	"k8s.io/client-go/kubernetes/scheme"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_apiService(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storeFake.MockStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				objectFile := "apiservice_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusOK,
				Details:    []component.Component{component.NewText("API Service is OK")},
			},
		},
		{
			name: "unavailable",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				objectFile := "apiservice_unavailable.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusError,
				Details:    []component.Component{component.NewText("Not available: (ServiceNotFound) service/metrics-server in \"kube-system\" is not present")},
			},
		},
		{
			name: "unknown",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				objectFile := "apiservice_unknown.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("No available condition for this apiservice")},
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
			name: "object is not an apiservice",
			init: func(t *testing.T, o *storeFake.MockStore) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)

			object := tc.init(t, o)

			ctx := context.Background()
			status, err := apiService(ctx, object, o)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}

func init() {
	apiregistrationv1.AddToScheme(scheme.Scheme)
}
