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

	storefake "github.com/heptio/developer-dash/pkg/store/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/store"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_service(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storefake.MockStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Endpoints",
					Name:       "stateful",
				}

				endpoints := testutil.LoadObjectFromFile(t, "endpoints_ok.yaml")

				o.EXPECT().Get(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructured(t, endpoints), nil)

				objectFile := "service_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusOK,
				Details:    []component.Component{component.NewText("Service is OK")},
			},
		},
		{
			name: "no endpoint subsets",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Endpoints",
					Name:       "stateful",
				}

				endpoints := testutil.LoadObjectFromFile(t, "endpoints_no_subsets.yaml")

				o.EXPECT().Get(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructured(t, endpoints), nil)

				objectFile := "service_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("Service has no endpoints")},
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
			status, err := service(ctx, object, o)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
