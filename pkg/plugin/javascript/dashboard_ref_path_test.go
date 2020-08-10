/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/octant/fake"
)

func TestDashboardRefPath_Name(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	linkGenerator := fake.NewMockLinkGenerator(ctrl)

	d := NewDashboardRefPath(linkGenerator)

	want := "RefPath"
	got := d.Name()

	require.Equal(t, want, got)
}

func TestDashboardRefPath_Call(t *testing.T) {
	type ctorArgs struct {
		linkGenerator func(ctx context.Context, ctrl *gomock.Controller) octant.LinkGenerator
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		call     string
		wantErr  bool
	}{
		{
			name: "in general",
			ctorArgs: ctorArgs{
				linkGenerator: func(ctx context.Context, ctrl *gomock.Controller) octant.LinkGenerator {
					linkGenerator := fake.NewMockLinkGenerator(ctrl)
					linkGenerator.EXPECT().
						ObjectPath("test", "v1", "ReplicaSet", "my-replica-set").
						Return("/path", nil)
					return linkGenerator
				},
			},
			call: `dashClient.RefPath({namespace:'test', apiVersion: 'v1', kind:'ReplicaSet', name: 'my-replica-set'})`,
		},
		{
			name: "ref path fails",
			ctorArgs: ctorArgs{
				linkGenerator: func(ctx context.Context, ctrl *gomock.Controller) octant.LinkGenerator {
					linkGenerator := fake.NewMockLinkGenerator(ctrl)
					linkGenerator.EXPECT().
						ObjectPath("test", "v1", "ReplicaSet", "my-replica-set").
						Return("", errors.New("error"))
					return linkGenerator
				},
			},
			call:    `dashClient.RefPath({namespace:'test', apiVersion: 'v1', kind:'ReplicaSet', name: 'my-replica-set'})`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			d := NewDashboardRefPath(tt.ctorArgs.linkGenerator(ctx, ctrl))

			runner := functionRunner{wantErr: tt.wantErr}
			runner.run(ctx, t, d, tt.call)
		})
	}
}
