package kubernetes

import (
	"errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// CRDResources returns a list of resources identified by group/version/kind a CRD supports.
func CRDResources(crd *unstructured.Unstructured) ([]schema.GroupVersionKind, error) {
	if crd == nil {
		return nil, errors.New("crd is nil")
	}

	var list []schema.GroupVersionKind

	group, ok, err := unstructured.NestedString(crd.Object, "spec", "group")
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("crd did not have a spec.group")
	}

	kind, ok, err := unstructured.NestedString(crd.Object, "spec", "names", "kind")
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("crd did not have a spec.names.kind")
	}

	versionsRaw, ok, err := unstructured.NestedSlice(crd.Object, "spec", "versions")
	if err != nil {
		return nil, err
	}

	if ok {
		for _, versionRaw := range versionsRaw {
			version, ok := versionRaw.(map[string]interface{})
			if !ok {
				return nil, errors.New("version was of an unknown type")
			}

			isServed, ok, err := unstructured.NestedBool(version, "served")
			if err != nil {
				return nil, err
			}

			if !ok {
				return nil, errors.New("version doesn't have served entry")
			}

			name, ok, err := unstructured.NestedString(version, "name")
			if err != nil {
				return nil, err
			}

			if !ok {
				return nil, errors.New("version doesn't have name")
			}

			if isServed {
				g := schema.GroupVersionKind{
					Group:   group,
					Version: name,
					Kind:    kind,
				}
				if !groupVersionKindsContains(g, list) {
					list = append(list, g)
				}

			}
		}
	}

	version, ok, err := unstructured.NestedString(crd.Object, "spec", "version")
	if err != nil {
		return nil, err
	}

	if ok {
		g := schema.GroupVersionKind{
			Group:   group,
			Version: version,
			Kind:    kind,
		}
		if !groupVersionKindsContains(g, list) {
			list = append(list, g)
		}

	}

	return list, nil
}

// CRDContainsResource returns true if a CRD contains a resource.
func CRDContainsResource(crd *unstructured.Unstructured, g schema.GroupVersionKind) (bool, error) {
	list, err := CRDResources(crd)
	if err != nil {
		return false, err
	}

	return groupVersionKindsContains(g, list), nil
}

func groupVersionKindsContains(g schema.GroupVersionKind, list []schema.GroupVersionKind) bool {
	for i := range list {
		if g.String() == list[i].String() {
			return true
		}
	}
	return false
}
