package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/testutil"
)

func Test_CRDResources(t *testing.T) {
	crd1 := testutil.CreateCRD("test")
	crd1.Spec.Group = "group"
	crd1.Spec.Names.Kind = "kind"
	crd1.Spec.Version = ""
	crd1.Spec.Versions = []apiextv1beta1.CustomResourceDefinitionVersion{
		{
			Name:   "v1",
			Served: true,
		},
	}

	crd2 := testutil.CreateCRD("test")
	crd2.Spec.Group = "group"
	crd2.Spec.Names.Kind = "kind"
	crd2.Spec.Version = "v1"

	tests := []struct {
		name     string
		crd      *unstructured.Unstructured
		expected []schema.GroupVersionKind
		isErr    bool
	}{
		{
			name: "with versions",
			crd:  testutil.ToUnstructured(t, crd1),
			expected: []schema.GroupVersionKind{
				{
					Group:   "group",
					Version: "v1",
					Kind:    "kind",
				},
			},
		},
		{
			name: "with version (deprecated)",
			crd:  testutil.ToUnstructured(t, crd1),
			expected: []schema.GroupVersionKind{
				{
					Group:   "group",
					Version: "v1",
					Kind:    "kind",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CRDResources(test.crd)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}

}
