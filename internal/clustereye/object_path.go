package clustereye

import (
	"context"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/log"
	dashStrings "github.com/heptio/developer-dash/internal/util/strings"
)

// CRDPathGenFunc is a function that generates a custom resource path.
type CRDPathGenFunc func(namespace, crdName, name string) (string, error)

// PathLookupFunc looks up paths for an object.
type PathLookupFunc func(namespace, apiVersion, kind, name string) (string, error)

// ObjectPathConfig is configuration for ObjectPath.
type ObjectPathConfig struct {
	ModuleName     string
	SupportedGVKs  []schema.GroupVersionKind
	PathLookupFunc PathLookupFunc
	CRDPathGenFunc CRDPathGenFunc
}

// Validate returns an error if the configuration is invalid.
func (opc *ObjectPathConfig) Validate() error {
	var errorStrings []string

	if opc.ModuleName == "" {
		errorStrings = append(errorStrings, "module name is blank")
	}

	if opc.PathLookupFunc == nil {
		errorStrings = append(errorStrings, "object path lookup func is nil")
	}

	if opc.CRDPathGenFunc == nil {
		errorStrings = append(errorStrings, "object path gen func is nil")
	}

	if len(errorStrings) > 0 {
		return errors.New(strings.Join(errorStrings, ", "))
	}

	return nil
}

// ObjectPath contains functions for generating paths for an object. Typically this is a
// helper which can be embedded in modules.
type ObjectPath struct {
	crds           map[string]*unstructured.Unstructured
	moduleName     string
	supportedGVKs  []schema.GroupVersionKind
	lookupFunc     PathLookupFunc
	crdPathGenFunc CRDPathGenFunc

	mu sync.Mutex
}

// NewObjectPath creates ObjectPath.
func NewObjectPath(config ObjectPathConfig) (*ObjectPath, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Wrap(err, "object path config is invalid")
	}

	return &ObjectPath{
		moduleName:     config.ModuleName,
		supportedGVKs:  config.SupportedGVKs,
		lookupFunc:     config.PathLookupFunc,
		crdPathGenFunc: config.CRDPathGenFunc,
	}, nil
}

// AddCRD adds support for a CRD to the ObjectPath.
func (op *ObjectPath) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	op.mu.Lock()
	defer op.mu.Unlock()

	if crd == nil {
		return errors.New("unable to add nil crd")
	}

	if op.crds == nil {
		op.crds = make(map[string]*unstructured.Unstructured)
	}
	op.crds[crd.GetName()] = crd

	logger := log.From(ctx)
	logger.
		With("module", op.moduleName, "crd", crd.GetName()).
		Debugf("adding CRD from module")
	return nil
}

// RemoveCRD removes support for a CRD from the ObjectPath.
func (op *ObjectPath) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	op.mu.Lock()
	defer op.mu.Unlock()

	if crd == nil {
		return errors.New("unable to remove nil crd")
	}

	delete(op.crds, crd.GetName())

	logger := log.From(ctx)
	logger.
		With("module", op.moduleName, "crd", crd.GetName()).
		Debugf("removing CRD from module")
	return nil
}

// SupportedGroupVersionKind returns a slice of GVKs this object path can handle.
func (op *ObjectPath) SupportedGroupVersionKind() []schema.GroupVersionKind {
	op.mu.Lock()
	defer op.mu.Unlock()

	gvks := make([]schema.GroupVersionKind, len(op.supportedGVKs))
	copy(gvks, op.supportedGVKs)

	for _, crd := range op.crds {
		r, err := CRDResourceGVKs(crd)
		if err != nil {
			continue
		}

		gvks = append(gvks, r...)
	}

	return gvks
}

// GroupVersionKind returns a path for an object.
func (op *ObjectPath) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	op.mu.Lock()
	defer op.mu.Unlock()

	gvk := schema.FromAPIVersionAndKind(apiVersion, kind)

	// if apiVersion matches a crd, build up path dynamically
	for i := range op.crds {
		crd := op.crds[i]

		supports, err := crdSupportsGVK(crd, gvk)
		if err != nil {
			return "", err
		}

		if !supports {
			continue
		}

		list, err := CRDAPIVersions(crd)
		if err != nil {
			return "", errors.WithMessagef(err, "unable to find api versions for %s", crd.GetName())
		}

		var apiVersions []string
		for _, gv := range list {
			apiVersions = append(apiVersions, path.Join(gv.Group, gv.Version))
		}

		if dashStrings.Contains(apiVersion, apiVersions) {
			return op.crdPathGenFunc(namespace, crd.GetName(), name)
		}
	}

	return op.lookupFunc(namespace, apiVersion, kind, name)
}

// CRDResourceGVKs returns the GVKs contained within a CRD.
func CRDResourceGVKs(crd *unstructured.Unstructured) ([]schema.GroupVersionKind, error) {
	apiVersions, err := CRDAPIVersions(crd)
	if err != nil {
		return nil, err
	}

	spec, ok := crd.Object["spec"].(map[string]interface{})
	if !ok {
		return nil, errors.New("crd did not have spec")
	}

	names, ok := spec["names"].(map[string]interface{})
	if !ok {
		return nil, errors.New("crd spec did not have names")
	}

	kind, ok := names["kind"].(string)
	if !ok {
		return nil, errors.New("crd spec names did not have kind")
	}

	var list []schema.GroupVersionKind

	for _, apiVersion := range apiVersions {
		list = append(list, schema.GroupVersionKind{Group: apiVersion.Group, Version: apiVersion.Version, Kind: kind})
	}

	return list, nil
}

// CRDAPIVersions returns the group versions that are contained within a CRD.
func CRDAPIVersions(crd *unstructured.Unstructured) ([]schema.GroupVersion, error) {
	if crd == nil {
		return nil, errors.New("crd is nil")
	}

	var list []schema.GroupVersion

	spec, ok := crd.Object["spec"].(map[string]interface{})
	if !ok {
		return nil, errors.New("crd did not have spec")
	}

	crdName, ok := spec["group"].(string)
	if !ok {
		return nil, errors.New("crd spec did not have group")
	}

	versions, ok := spec["versions"].([]interface{})
	if !ok {
		return nil, errors.New("crd spec did not have versions")
	}

	for _, rawVersion := range versions {
		version, ok := rawVersion.(map[string]interface{})
		if !ok {
			return nil, errors.New("crd version was not an object")
		}

		versionName, ok := version["name"].(string)
		if !ok {
			return nil, errors.New("crd version did not have a name")
		}

		list = append(list, schema.GroupVersion{Group: crdName, Version: versionName})
	}

	return list, nil
}

func crdSupportsGVK(crd *unstructured.Unstructured, gvk schema.GroupVersionKind) (bool, error) {
	if crd == nil {
		return false, errors.New("crd is nil")
	}

	spec, ok := crd.Object["spec"].(map[string]interface{})
	if !ok {
		return false, errors.New("crd did not have spec")
	}

	group, ok := spec["group"].(string)
	if !ok {
		return false, errors.New("crd spec did not have group")
	}

	names, ok := spec["names"].(map[string]interface{})
	if !ok {
		return false, errors.New("crd spec did not have names")
	}

	kind, ok := names["kind"].(string)
	if !ok {
		return false, errors.New("crd spec names did not have kind")
	}

	rawVersions, ok := spec["versions"].([]interface{})
	if !ok {
		return false, errors.New("crd spec did not have versions")
	}

	for _, rawVersion := range rawVersions {
		version, ok := rawVersion.(map[string]interface{})
		if !ok {
			return false, errors.New("crd version was not an object")
		}

		versionName, ok := version["name"].(string)
		if !ok {
			return false, errors.New("crd version did not have a name")
		}

		current := schema.GroupVersionKind{
			Group:   group,
			Kind:    kind,
			Version: versionName,
		}

		if current.String() == gvk.String() {
			return true, nil
		}
	}

	return false, nil
}
