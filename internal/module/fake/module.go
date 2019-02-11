package fake

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
)

// Module is a fake module.
type Module struct {
	name   string
	logger log.Logger

	ObservedContentPath string
	ObservedNamespace   string
}

// NewModule creates an instance of Module.
func NewModule(name string, logger log.Logger) *Module {
	return &Module{
		name:   name,
		logger: logger,
	}
}

// Name is the name of the module.
func (m *Module) Name() string {
	return m.name
}

// ContentPath is the path to the module's content.
func (m *Module) ContentPath() string {
	return fmt.Sprintf("/%s", m.name)
}

// Navigation returns navigation entries for the module.
func (m *Module) Navigation(prefix string) (*hcli.Navigation, error) {
	nav := &hcli.Navigation{
		Path:  prefix,
		Title: m.name,
	}

	return nav, nil
}

// SetNamespace sets the current namespace.
func (m *Module) SetNamespace(namespace string) error {
	return nil
}

// Start doesn't do anything.
func (m *Module) Start() error {
	return nil
}

// Stop doesn't do anything.
func (m *Module) Stop() {
}

func (m *Module) Content(ctx context.Context, contentPath, prefix, namespace string) (component.ContentResponse, error) {
	m.ObservedContentPath = contentPath
	m.ObservedNamespace = namespace

	switch contentPath {
	case "/":
		return component.ContentResponse{
			Title: []component.TitleViewComponent{
				component.NewText("/"),
			},
		}, nil
	case "/nested":
		return component.ContentResponse{
			Title: []component.TitleViewComponent{
				component.NewText("/nested"),
			},
		}, nil
	default:
		return component.ContentResponse{}, errors.New("not found")
	}
}
