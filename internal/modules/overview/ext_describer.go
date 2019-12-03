/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type ExtDescriber struct {
	extFactory func(ctx context.Context, namespace string, options describer.Options) ([]component.Component, error)
}

var _ describer.Describer = (*ExtDescriber)(nil)

func NewExtDescriber() *ExtDescriber {
	d := &ExtDescriber{
		extFactory: extFactory,
	}

	return d
}

func (e *ExtDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	// ext, err := extFactory(ctx, namespace, options)
	// if err != nil {
	// 	return component.EmptyContentResponse, err
	// }

	resp := component.ContentResponse{
		Title: nil,
		Components: []component.Component{
			component.NewFlexLayout("testing"),
		},
	}

	return resp, nil
}

func extFactory(ctx context.Context, namespace string, options describer.Options) ([]component.Component, error) {
	terminalManager := options.Dash.TerminalManager()
	terminals := terminalManager.List()

	var extTabs []component.Component
	for _, terminal := range terminals {
		fl := component.NewFlexLayout(terminal.Command())

		details := component.TerminalDetails{
			Container: terminal.Container(),
			Command:   terminal.Command(),
			UUID:      terminal.ID(),
			CreatedAt: terminal.CreatedAt(),
		}

		fl.AddSections([]component.FlexLayoutItem{
			{
				Width: component.WidthFull,
				View:  component.NewTerminal(namespace, terminal.Command(), details),
			},
		})

		extTabs = append(extTabs, fl)
	}
	return extTabs, nil
}

func (e *ExtDescriber) PathFilters() []describer.PathFilter {
	PathFilters := []describer.PathFilter{
		*describer.NewPathFilter("/extension", e),
	}

	return PathFilters
}

func (e *ExtDescriber) Reset(ctx context.Context) error {
	return nil
}
