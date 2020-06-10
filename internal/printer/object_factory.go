/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/gvk"
)

type ObjectFactory interface {
	Factory(object runtime.Object, options Options) (*Object, error)
}

type DefaultObjectFactory struct {
	typeFactory map[schema.GroupVersionKind]func(runtime.Object, Options) (*Object, error)
}

var _ ObjectFactory = DefaultObjectFactory{}

func NewDefaultObjectFactory() *DefaultObjectFactory {
	f := DefaultObjectFactory{
		typeFactory: map[schema.GroupVersionKind]func(runtime.Object, Options) (*Object, error){
			gvk.CustomResourceDefinition: CustomResourceDefinitionObjectFactory,
		},
	}

	return &f
}

func (d DefaultObjectFactory) Factory(object runtime.Object, options Options) (*Object, error) {
	if object == nil {
		return nil, fmt.Errorf("unable to create a print object for a nil runtime object")
	}

	groupVersionKind := object.GetObjectKind().GroupVersionKind()
	fn, ok := d.typeFactory[groupVersionKind]
	if !ok {
		o := NewObject(object)
		o.EnableEvents()

		return o, nil
	}

	o, err := fn(object, options)
	if err != nil {
		return nil, fmt.Errorf("create object factory for %s: %w", groupVersionKind, err)
	}

	return o, nil
}
