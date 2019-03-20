package fake

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/sugarloaf"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
)

// Module is a fake module.
type Module struct {
	name   string
	logger log.Logger

	ObservedContentPath string
	ObservedNamespace   string
}

var _ module.Module = (*Module)(nil)

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
func (m *Module) Navigation(ctx context.Context, namespace, prefix string) (*sugarloaf.Navigation, error) {
	nav := &sugarloaf.Navigation{
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

func (m *Module) Content(ctx context.Context, contentPath, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
	m.ObservedContentPath = contentPath
	m.ObservedNamespace = namespace

	switch contentPath {
	case "/":
		return component.ContentResponse{
			Title: component.Title(component.NewText("/")),
		}, nil
	case "/nested":
		return component.ContentResponse{
			Title: component.Title(component.NewText("/nested")),
		}, nil
	default:
		return component.ContentResponse{}, errors.New("not found")
	}
}

func (m *Module) Handlers(ctx context.Context) map[string]http.Handler {
	return make(map[string]http.Handler)
}
