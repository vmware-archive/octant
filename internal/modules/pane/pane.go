package pane

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/module"
)

// Options are options for configuring Module.
type Options struct {
	DashConfig config.Dash
}

// Module is an applications module.
type Module struct {
	Options
	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*Module)(nil)

// New creates an instance of Module.
func New(ctx context.Context, options Options) *Module {
	pm := describer.NewPathMatcher("pane")
	// for _, pf := range rootDescriber.PathFilters() {
	// 	pm.Register(ctx, pf)
	// }

	appDescriber := NewPaneDescriber(options.DashConfig)
	for _, pf := range appDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	return &Module{
		Options:     options,
		pathMatcher: pm,
	}
}
