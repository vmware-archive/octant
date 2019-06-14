package clustereye

import (
	"context"
	"path"
	"sort"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/pkg/store"
)

// Navigation is a set of navigation entries.
type Navigation struct {
	Title    string       `json:"title,omitempty"`
	Path     string       `json:"path,omitempty"`
	Children []Navigation `json:"children,omitempty"`
}

// NewNavigation creates a Navigation.
func NewNavigation(title, path string) *Navigation {
	return &Navigation{Title: title, Path: path}
}

// CRDEntries generates navigation entries for crds.
func CRDEntries(ctx context.Context, prefix, namespace string, objectStore store.Store) ([]Navigation, error) {
	var list []Navigation

	crdNames, err := CustomResourceDefinitionNames(ctx, objectStore)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving CRD names")
	}

	sort.Strings(crdNames)

	for _, name := range crdNames {
		crd, err := CustomResourceDefinition(ctx, name, objectStore)
		if err != nil {
			return nil, errors.Wrapf(err, "load %q custom resource definition", name)
		}

		objects, err := ListCustomResources(ctx, crd, namespace, objectStore, nil)
		if err != nil {
			return nil, err
		}

		if len(objects) > 0 {
			list = append(list, *NewNavigation(name, path.Join(prefix, name)))
		}
	}

	return list, nil
}

// CustomResourceDefinitionNames returns the available custom resource definition names.
func CustomResourceDefinitionNames(ctx context.Context, o store.Store) ([]string, error) {
	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	if err := o.HasAccess(key, "list"); err != nil {
		return []string{}, nil
	}

	rawList, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "listing CRDs")
	}

	var list []string

	for _, object := range rawList {
		crd := &apiextv1beta1.CustomResourceDefinition{}

		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, crd); err != nil {
			return nil, errors.Wrap(err, "crd conversion failed")
		}

		list = append(list, crd.Name)
	}

	return list, nil
}

// CustomResourceDefinition retrieves a CRD.
func CustomResourceDefinition(ctx context.Context, name string, o store.Store) (*apiextv1beta1.CustomResourceDefinition, error) {
	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       name,
	}

	crd := &apiextv1beta1.CustomResourceDefinition{}
	if err := store.GetObjectAs(ctx, o, key, crd); err != nil {
		return nil, errors.Wrap(err, "get CRD from object store")
	}

	return crd, nil
}

// ListCustomResources lists all custom resources given a CRD.
func ListCustomResources(
	ctx context.Context,
	crd *apiextv1beta1.CustomResourceDefinition,
	namespace string,
	o store.Store,
	selector *labels.Set) ([]*unstructured.Unstructured, error) {
	if crd == nil {
		return nil, errors.New("crd is nil")
	}
	gvk := schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: crd.Spec.Version,
		Kind:    crd.Spec.Names.Kind,
	}

	apiVersion, kind := gvk.ToAPIVersionAndKind()

	key := store.Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Selector:   selector,
	}

	if err := o.HasAccess(key, "list"); err != nil {
		return []*unstructured.Unstructured{}, nil
	}

	objects, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "listing custom resources for %q", crd.Name)
	}

	return objects, nil
}
