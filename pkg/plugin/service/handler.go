package service

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/view/component"
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
	router           *Router
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
func (p *Handler) Register(ctx context.Context, dashboardAPIAddress string) (plugin.Metadata, error) {
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
func (p *Handler) Print(ctx context.Context, object runtime.Object) (plugin.PrintResponse, error) {
	if p.HandlerFuncs.Print == nil {
		return plugin.PrintResponse{}, nil
	}

	request := &PrintRequest{
		baseRequest:     newBaseRequest(ctx, p.name),
		DashboardClient: p.dashboardClient,
		Object:          object,
	}

	return p.HandlerFuncs.Print(request)
}

// PrintTab prints a tab for an object.
func (p *Handler) PrintTab(ctx context.Context, object runtime.Object) (plugin.TabResponse, error) {
	if p.HandlerFuncs.PrintTab == nil {
		return plugin.TabResponse{}, nil
	}

	request := &PrintRequest{
		baseRequest:     newBaseRequest(ctx, p.name),
		DashboardClient: p.dashboardClient,
		Object:          object,
	}

	return p.HandlerFuncs.PrintTab(request)
}

// ObjectStatus creates status for an object.
func (p *Handler) ObjectStatus(ctx context.Context, object runtime.Object) (plugin.ObjectStatusResponse, error) {
	if p.HandlerFuncs.ObjectStatus == nil {
		return plugin.ObjectStatusResponse{}, nil
	}

	request := &PrintRequest{
		baseRequest:     newBaseRequest(ctx, p.name),
		DashboardClient: p.dashboardClient,
		Object:          object,
	}

	return p.HandlerFuncs.ObjectStatus(request)
}

// HandleAction handles actions given a payload.
func (p *Handler) HandleAction(ctx context.Context, payload action.Payload) error {
	if p.HandlerFuncs.HandleAction == nil {
		return nil
	}

	request := &ActionRequest{
		baseRequest:     newBaseRequest(ctx, p.name),
		DashboardClient: p.dashboardClient,
		Payload:         payload,
	}

	return p.HandlerFuncs.HandleAction(request)
}

// Navigation creates navigation.
func (p *Handler) Navigation(ctx context.Context) (navigation.Navigation, error) {
	if p.HandlerFuncs.Navigation == nil {
		return navigation.Navigation{}, nil
	}

	request := &NavigationRequest{
		baseRequest:     newBaseRequest(ctx, p.name),
		DashboardClient: p.dashboardClient,
	}

	return p.HandlerFuncs.Navigation(request)
}

// Content creates content for a given content path.
func (p *Handler) Content(ctx context.Context, contentPath string) (component.ContentResponse, error) {
	handlerFunc, ok := p.router.Match(contentPath)
	if !ok {
		return component.ContentResponse{}, nil
	}

	request := &Request{
		baseRequest:     newBaseRequest(ctx, p.name),
		dashboardClient: p.dashboardClient,
		Path:            contentPath,
	}

	return handlerFunc(request)
}
