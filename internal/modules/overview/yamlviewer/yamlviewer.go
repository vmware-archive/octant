/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package yamlviewer

import (
	"github.com/vmware/octant/pkg/view/component"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ToComponent converts an object into a YAML component.
func ToComponent(object runtime.Object) (*component.YAML, error) {
	yv, err := new(object)
	if err != nil {
		return nil, errors.Wrap(err, "create YAML viewer")
	}

	return yv.ToComponent()
}

// YAMLViewer is a YAML viewer for objects.
type yamlViewer struct {
	object runtime.Object
}

// New creates an instance of YAMLViewer.
func new(object runtime.Object) (*yamlViewer, error) {
	if object == nil {
		return nil, errors.New("can't create YAML view for nil object")
	}

	return &yamlViewer{
		object: object,
	}, nil
}

// ToComponent converts the YAMLViewer to a component.
func (yv *yamlViewer) ToComponent() (*component.YAML, error) {
	y := component.NewYAML(component.TitleFromString("YAML"), "")
	if err := y.Data(yv.object); err != nil {
		return nil, errors.Wrap(err, "add YAML data")
	}

	return y, nil
}
