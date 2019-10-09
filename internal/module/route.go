/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package module

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router allows registering handlers for a path pattern.
// Routes form a tree and subroutes can be registered.
// Route is a subset of mux.Router.
type Router interface {
	Handle(path string, handler http.Handler) Route
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) Route
	PathPrefix(path string) Route
}

// Route allows further tuning the matching of a route.
// Route is a subset of mux.Route.
type Route interface {
	Handler(handler http.Handler) Route
	HandlerFunc(f func(http.ResponseWriter, *http.Request)) Route
	Methods(methods ...string) Route
	Subrouter() Router
}

type MuxRouter struct {
	*mux.Router
}

type MuxRoute struct {
	*mux.Route
}

func (m MuxRoute) Handler(handler http.Handler) Route {
	return MuxRoute{m.Route.Handler(handler)}
}

func (m MuxRoute) HandlerFunc(f func(http.ResponseWriter, *http.Request)) Route {
	return MuxRoute{m.Route.HandlerFunc(f)}
}

func (m MuxRoute) Methods(methods ...string) Route {
	return MuxRoute{m.Route.Methods(methods...)}
}

func (m MuxRoute) Subrouter() Router {
	return MuxRouter{m.Route.Subrouter()}
}

func (m MuxRouter) PathPrefix(path string) Route {
	r := m.Router.PathPrefix(path)
	return MuxRoute{r}
}

func (m MuxRouter) Handle(path string, handler http.Handler) Route {
	r := m.Router.Handle(path, handler)
	return MuxRoute{r}
}

func (m MuxRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) Route {
	return MuxRoute{m.Router.HandleFunc(path, f)}

}
