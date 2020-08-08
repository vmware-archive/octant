/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package describer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestStoreResourceLoader_Load(t *testing.T) {
	myCR := testutil.CreateCustomResource("my-cr", func(u *unstructured.Unstructured) {
		u.SetNamespace("test")
	})
	myCRD := testutil.ToUnstructured(t, testutil.CreateCRD("my-crd", func(definition *apiextv1.CustomResourceDefinition) {
		definition.Spec.Group = myCR.GetObjectKind().GroupVersionKind().Group
		definition.Spec.Names.Kind = myCR.GetKind()
	}))
	myCRD.SetAPIVersion("apiextensions.k8s.io/v1")

	type ctorArgs struct {
		store func(t *testing.T, ctrl *gomock.Controller) store.Store
	}
	type args struct {
		descriptor ResourceDescriptor
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		args     args
		want     *ResourceLoadResponse
		wantErr  bool
	}{
		{
			name: "in general",
			ctorArgs: ctorArgs{
				store: func(t *testing.T, ctrl *gomock.Controller) store.Store {
					s := fake.NewMockStore(ctrl)

					crdKey, err := store.KeyFromObject(myCRD)
					require.NoError(t, err)
					s.EXPECT().
						Get(gomock.Any(), crdKey).
						Return(myCRD, nil)

					crKey, err := store.KeyFromObject(myCR)
					require.NoError(t, err)
					s.EXPECT().
						Get(gomock.Any(), crKey).
						Return(myCR, nil)

					return s
				},
			},
			args: args{
				descriptor: ResourceDescriptor{
					CustomResourceDefinitionName: "my-crd",
					Namespace:                    "test",
					CustomResourceVersion:        "v1",
					CustomResourceName:           "my-cr",
				},
			},
			want: &ResourceLoadResponse{
				CustomResource:           myCR,
				CustomResourceDefinition: myCRD,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rl := NewStoreResourceLoader(tt.ctorArgs.store(t, ctrl))

			got, err := rl.Load(context.Background(), tt.args.descriptor)
			testutil.RequireErrorOrNot(t, tt.wantErr, err)

			require.Equal(t, tt.want, got)
		})
	}
}
