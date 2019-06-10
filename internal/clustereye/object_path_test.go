package clustereye

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/testutil"
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
	crd := testutil.CreateCRD("my-crd")
	crd.Spec.Group = "group"

	crd.Spec.Versions = []apiextv1beta1.CustomResourceDefinitionVersion{
		{
			Name: "v1",
		},
	}

	crd.Spec.Names = apiextv1beta1.CustomResourceDefinitionNames{
		Kind: "kind",
	}

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

func TestCRDResourceGVKs(t *testing.T) {
	crd := testutil.CreateCRD("my-crd")
	crd.Spec.Group = "group"

	crd.Spec.Versions = []apiextv1beta1.CustomResourceDefinitionVersion{
		{
			Name: "v1",
		},
	}

	crd.Spec.Names = apiextv1beta1.CustomResourceDefinitionNames{
		Kind: "kind",
	}

	got, err := CRDResourceGVKs(testutil.ToUnstructured(t, crd))
	require.NoError(t, err)

	expected := []schema.GroupVersionKind{
		{Group: "group", Version: "v1", Kind: "kind"},
	}

	assert.Equal(t, expected, got)
}

func TestCRDAPIVersions(t *testing.T) {
	crd := testutil.CreateCRD("my-crd")
	crd.Spec.Group = "group"

	crd.Spec.Versions = []apiextv1beta1.CustomResourceDefinitionVersion{
		{
			Name: "v1",
		},
	}

	got, err := CRDAPIVersions(testutil.ToUnstructured(t, crd))
	require.NoError(t, err)

	expected := []schema.GroupVersion{
		{Group: "group", Version: "v1"},
	}

	assert.Equal(t, expected, got)
}
