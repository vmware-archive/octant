package service

import (
	"github.com/gobwas/glob"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// HandleFunc is a function that generates a content response.
type HandleFunc func(request *Request) (component.ContentResponse, error)

// Request represents a path request from Octant. It will always be a
// GET style request with a path.
type Request struct {
	baseRequest

	dashboardClient Dashboard

	// Path is path that Octant is requesting. It is scoped to the plugin.
	// i.e. If Octant wants to render /content/plugin/foo, Path will be
	// `/foo`.
	Path string
}

// DashboardClient returns a dashboard client for the request.
func (r *Request) DashboardClient() Dashboard {
	return r.dashboardClient
}

type route struct {
	path       string
	glob       glob.Glob
	handleFunc HandleFunc
}

// Router is a router for the plugin. A plugin can register a HandleFuncs to
// a path.
type Router struct {
	routes []route
}

// NewRouter creates a Router.
func NewRouter() *Router {
	return &Router{}
}

// HandleFunc registers a HandleFunc to a path. Paths can contain globs.
// e.g `/*` will match `/foo` if an explicit `/foo` path (or glob) has
// not already been registered. Routes are evaluated in the order they
// were added.
func (r *Router) HandleFunc(routePath string, handleFunc HandleFunc) {
	pathGlob, err := glob.Compile(routePath)
	if err != nil {
		return
	}

	r.routes = append(r.routes, route{
		path:       routePath,
		glob:       pathGlob,
		handleFunc: handleFunc,
	})
}

// Match matches a path against a exiting route. If no route is found,
// it will return false in the second return value.
func (r *Router) Match(contentPath string) (HandleFunc, bool) {
	for _, route := range r.routes {
		if route.glob.Match(contentPath) {
			return route.handleFunc, true
		}
	}

	return nil, false
}
