/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestDiscoveryResourceInfo_PreferredVersion(t *testing.T) {
	resourceLists := []*metav1.APIResourceList{
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Kind: "Pod"},
			},
		},
		{
			GroupVersion: "apps/v1",
			APIResources: []metav1.APIResource{
				{Kind: "ReplicaSet"},
				{Kind: "Deployment"},
			},
		},
	}

	type ctorArgs struct {
		resourceLists []*metav1.APIResourceList
		options       func() []Option
	}
	type args struct {
		groupKind schema.GroupKind
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		args     args
		wantErr  bool
		want     string
	}{
		{
			name: "version exists",
			ctorArgs: ctorArgs{
				resourceLists: resourceLists,
				options: func() []Option {
					return nil
				},
			},
			args: args{
				groupKind: schema.GroupKind{Group: "apps", Kind: "Deployment"},
			},
			want: "v1",
		},
		{
			name: "invalid group version",
			ctorArgs: ctorArgs{
				resourceLists: resourceLists,
				options: func() []Option {
					return []Option{
						GroupVersionParser(func(groupVersion string) (schema.GroupVersion, error) {
							return schema.GroupVersion{}, fmt.Errorf("error")
						}),
					}
				},
			},
			args: args{
				groupKind: schema.GroupKind{Group: "apps", Kind: "Deployment"},
			},
			wantErr: true,
		},
		{
			name: "not found",
			ctorArgs: ctorArgs{
				resourceLists: resourceLists,
				options: func() []Option {
					return nil
				},
			},
			args: args{
				groupKind: schema.GroupKind{Group: "apps", Kind: "Invalid"},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dri := NewDiscoveryResourceInfo(test.ctorArgs.resourceLists, test.ctorArgs.options()...)
			got, err := dri.PreferredVersion(test.args.groupKind)
			testutil.RequireErrorOrNot(t, test.wantErr, err, func() {
				require.Equal(t, test.want, got)
			})

		})
	}
}
