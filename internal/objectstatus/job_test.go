/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_runJobStatus(t *testing.T) {
	cases := []struct {
		name     string
		init     func(t *testing.T) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "job succeeded",
			init: func(t *testing.T) runtime.Object {
				objectFile := "job_success.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				Details: []component.Component{
					component.NewText("Job has succeeded 1 time"),
					component.NewText("Job completed in 11s"),
				},
			},
		},
		{
			name: "job in progress",
			init: func(t *testing.T) runtime.Object {
				objectFile := "job_in_progress.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details: []component.Component{
					component.NewText("Job has failed 2 times"),
					component.NewText("Job is in progress"),
				},
			},
		},
		{
			name: "job failed",
			init: func(t *testing.T) runtime.Object {
				objectFile := "job_failed.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusError,
				Details: []component.Component{
					component.NewText("Job has failed 5 times"),
					component.NewText("Job has reached the specified backoff limit"),
				},
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
			name: "object is not an ingress",
			init: func(t *testing.T) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := tc.init(t)

			ctx := context.Background()
			status, err := runJobStatus(ctx, object, nil)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
