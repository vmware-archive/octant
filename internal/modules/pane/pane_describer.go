package pane

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type PaneDescriber struct {
	printer     *printer.Resource
	paneFactory func(ctx context.Context, namespace string, options describer.Options) ([]component.Component, error)
}

var _ describer.Describer = (*PaneDescriber)(nil)

func NewPaneDescriber(dashConfig config.Dash) *PaneDescriber {
	p := printer.NewResource(dashConfig)
	d := &PaneDescriber{
		printer:     p,
		paneFactory: paneFactory,
	}

	return d
}

func (p *PaneDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	pane, err := p.paneFactory(ctx, namespace, options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	resp := component.ContentResponse{
		Title:      nil,
		Components: pane,
	}

	return resp, nil
}

func paneFactory(ctx context.Context, namespace string, options describer.Options) ([]component.Component, error) {
	terminalManager := options.Dash.TerminalManager()
	terminals := terminalManager.List()

	var paneTabs []component.Component
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

		paneTabs = append(paneTabs, fl)
	}
	return paneTabs, nil
}

func (p *PaneDescriber) PathFilters() []describer.PathFilter {
	return nil
}

func (p *PaneDescriber) Reset(ctx context.Context) error {
	return nil
}
