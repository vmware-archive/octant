/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ApplyYamlDescriber describes an apply
type ApplyYamlDescriber struct {
}

var _ describer.Describer = (*ApplyYamlDescriber)(nil)

// Describe describes the apply yaml interface
func (d *ApplyYamlDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	title := append([]component.TitleComponent{}, component.NewText("Apply YAML"))
	editor := component.NewEditor(component.TitleFromString("YAML"), "", false)
	editor.Config.SubmitLabel = "Apply"
	editor.Config.SubmitAction = octant.ActionApplyYaml
	list := component.NewList(title, []component.Component{editor})

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func (d *ApplyYamlDescriber) PathFilters() []describer.PathFilter {
	filter := describer.NewPathFilter("/apply", d)
	return []describer.PathFilter{*filter}
}

func (d *ApplyYamlDescriber) Reset(ctx context.Context) error {
	return nil
}

func NewApplyYamlDescriber() *ApplyYamlDescriber {
	return &ApplyYamlDescriber{}
}
