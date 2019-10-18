/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestObjectPathConfig_Validate(t *testing.T) {
	tests := []struct {
		name           string
		moduleName     string
		pathLookupFunc PathLookupFunc
		crdPathGenFunc CRDPathGenFunc
		isErr          bool
	}{
		{
			name:       "in general",
			moduleName: "module",
			pathLookupFunc: func(string, string, string, string) (string, error) {
				return "/path", nil
			},
			crdPathGenFunc: func(string, string, string) (string, error) {
				return "/path", nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := ObjectPathConfig{
				ModuleName:     test.moduleName,
				PathLookupFunc: test.pathLookupFunc,
				CRDPathGenFunc: test.crdPathGenFunc,
			}

			err := config.Validate()
			if test.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

		})
	}
}

func TestObjectPath(t *testing.T) {
	config := ObjectPathConfig{
		ModuleName: "module",
		PathLookupFunc: func(string, string, string, string) (string, error) {
			return "/path", nil
		},
		CRDPathGenFunc: func(string, string, string) (string, error) {
			return "/crd-path", nil
		},
	}

	objectPath, err := NewObjectPath(config)
	require.NoError(t, err)

	ctx := context.Background()
	crd := testutil.CreateCRD("my-crd", testutil.WithGenericCRD())

	err = objectPath.AddCRD(ctx, testutil.ToUnstructured(t, crd))
	require.NoError(t, err)

	require.Contains(t, objectPath.crds, crd.Name)

	expectedSupported := schema.GroupVersionKind{Group: "group", Version: "v1", Kind: "kind"}
	supported := objectPath.SupportedGroupVersionKind()
	assert.Contains(t, supported, expectedSupported)

	crPath, err := objectPath.GroupVersionKindPath("namespace", "group/v1", "kind", "name")
	require.NoError(t, err)
	assert.Equal(t, "/crd-path", crPath)

	err = objectPath.RemoveCRD(ctx, testutil.ToUnstructured(t, crd))
	require.NoError(t, err)

	require.NotContains(t, objectPath.crds, crd.Name)
}

func TestCRDAPIVersions(t *testing.T) {
	crd := testutil.CreateCRD("my-crd", testutil.WithGenericCRD())
	got, err := CRDAPIVersions(testutil.ToUnstructured(t, crd))
	require.NoError(t, err)

	expected := []schema.GroupVersion{
		{Group: "group", Version: "v1"},
	}

	assert.Equal(t, expected, got)
}
