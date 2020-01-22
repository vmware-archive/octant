/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"context"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	dashStrings "github.com/vmware-tanzu/octant/internal/util/strings"
)

// CRDPathGenFunc is a function that generates a custom resource path.
type CRDPathGenFunc func(namespace, crdName, version, name string) (string, error)

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
		// nothing to do if crd is nil
		return nil
	}

	delete(op.crds, crd.GetName())

	logger := log.From(ctx)
	logger.
		With("module", op.moduleName, "crd", crd.GetName()).
		Debugf("removing CRD from module")
	return nil
}

// ResetCRDs deletes all the CRD paths ObjectPath is tracking.
func (op *ObjectPath) ResetCRDs(ctx context.Context) error {
	op.mu.Lock()
	defer op.mu.Unlock()

	for k := range op.crds {
		delete(op.crds, k)
	}

	return nil
}

// SupportedGroupVersionKind returns a slice of GVKs this object path can handle.
func (op *ObjectPath) SupportedGroupVersionKind() []schema.GroupVersionKind {
	op.mu.Lock()
	defer op.mu.Unlock()

	list := make([]schema.GroupVersionKind, len(op.supportedGVKs))
	copy(list, op.supportedGVKs)

	for _, crd := range op.crds {
		r, err := kubernetes.CRDResources(crd)
		if err != nil {
			continue
		}

		list = append(list, r...)
	}

	return list
}

// GroupVersionKind returns a path for an object.
func (op *ObjectPath) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	op.mu.Lock()
	defer op.mu.Unlock()

	g := schema.FromAPIVersionAndKind(apiVersion, kind)

	// if apiVersion matches a crd, build up path dynamically
	for i := range op.crds {
		crd := op.crds[i]

		supports, err := kubernetes.CRDContainsResource(crd, g)
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
			return op.crdPathGenFunc(namespace, crd.GetName(), g.Version, name)
		}

	}

	return op.lookupFunc(namespace, apiVersion, kind, name)
}

// CRDAPIVersions returns the group versions that are contained within a CRD.
func CRDAPIVersions(crd *unstructured.Unstructured) ([]schema.GroupVersion, error) {

	resources, err := kubernetes.CRDResources(crd)
	if err != nil {
		return nil, err
	}

	list := make([]schema.GroupVersion, len(resources))
	for i := range resources {
		list[i] = resources[i].GroupVersion()
	}

	return list, nil
}
