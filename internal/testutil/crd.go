package testutil

import (
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func CreateCR(group, version, kind, name string) *unstructured.Unstructured {
	m := make(map[string]interface{})
	u := &unstructured.Unstructured{Object: m}

	u.SetName(name)
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})

	return u
}

func CreateCRDWithKind(name, kind string, isClusterScoped bool) *apiextv1.CustomResourceDefinition {
	scope := apiextv1.ClusterScoped
	if !isClusterScoped {
		scope = apiextv1.NamespaceScoped
	}

	crd := CreateCRD(name)
	crd.Spec.Scope = scope
	crd.Spec.Group = "testing"
	crd.Spec.Names = apiextv1.CustomResourceDefinitionNames{
		Kind: kind,
	}
	crd.Spec.Versions = []apiextv1.CustomResourceDefinitionVersion{
		{
			Name:    "v1",
			Served:  true,
			Storage: true,
		},
	}

	return crd
}
