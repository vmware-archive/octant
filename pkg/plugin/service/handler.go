package service

import (
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/view/component"
)

// Handler is the plugin service helper handler. Functions on this struct are called from Octant.
type Handler struct {
	HandlerFuncs

	mu sync.Mutex

	name         string
	description  string
	capabilities *plugin.Capabilities

	dashboardFactory func(dashboardAPIAddress string) (Dashboard, error)
	dashboardClient  Dashboard
}

var _ plugin.Service = (*Handler)(nil)

// Validate validates Handler.
func (p *Handler) Validate() error {
	if p.dashboardFactory == nil {
		return errors.New("plugin handler doesn't know how to create a dashboard client")
	}

	return nil
}

// Register registers a plugin with Octant.
func (p *Handler) Register(dashboardAPIAddress string) (plugin.Metadata, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	client, err := p.dashboardFactory(dashboardAPIAddress)
	if err != nil {
		return plugin.Metadata{}, errors.Wrap(err, "create api client")
	}

	p.dashboardClient = client

	return plugin.Metadata{
		Name:         p.name,
		Description:  p.description,
		Capabilities: *p.capabilities,
	}, nil
}

// Print prints components for an object.
func (p *Handler) Print(object runtime.Object) (plugin.PrintResponse, error) {
	if p.HandlerFuncs.Print == nil {
		return plugin.PrintResponse{}, nil
	}

	return p.HandlerFuncs.Print(p.dashboardClient, object)
}

// PrintTab prints a tab for an object.
func (p *Handler) PrintTab(object runtime.Object) (*component.Tab, error) {
	if p.HandlerFuncs.PrintTab == nil {
		return &component.Tab{}, nil
	}

	return p.HandlerFuncs.PrintTab(p.dashboardClient, object)
}

// ObjectStatus creates status for an object.
func (p *Handler) ObjectStatus(object runtime.Object) (plugin.ObjectStatusResponse, error) {
	if p.HandlerFuncs.ObjectStatus == nil {
		return plugin.ObjectStatusResponse{}, nil
	}

	return p.HandlerFuncs.ObjectStatus(p.dashboardClient, object)
}

// HandleAction handles actions given a payload.
func (p *Handler) HandleAction(payload action.Payload) error {
	if p.HandlerFuncs.HandleAction == nil {
		return nil
	}

	return p.HandlerFuncs.HandleAction(p.dashboardClient, payload)
}
