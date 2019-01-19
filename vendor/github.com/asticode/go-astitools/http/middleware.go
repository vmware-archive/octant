package astihttp

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// ChainMiddlewares chains middlewares
func ChainMiddlewares(h http.Handler, ms ...Middleware) http.Handler {
	return ChainMiddlewaresWithPrefix(h, []string{}, ms...)
}

// ChainMiddlewaresWithPrefix chains middlewares if one of prefixes is present
func ChainMiddlewaresWithPrefix(h http.Handler, prefixes []string, ms ...Middleware) http.Handler {
	for _, m := range ms {
		if len(prefixes) == 0 {
			h = m(h)
		} else {
			t := h
			h = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				for _, prefix := range prefixes {
					if strings.HasPrefix(r.URL.EscapedPath(), prefix) {
						m(t).ServeHTTP(rw, r)
						return
					}
				}
				t.ServeHTTP(rw, r)
			})
		}
	}
	return h
}

// ChainRouterMiddlewares chains router middlewares
func ChainRouterMiddlewares(h httprouter.Handle, ms ...RouterMiddleware) httprouter.Handle {
	return ChainRouterMiddlewaresWithPrefix(h, []string{}, ms...)
}

// ChainRouterMiddlewares chains router middlewares if one of prefixes is present
func ChainRouterMiddlewaresWithPrefix(h httprouter.Handle, prefixes []string, ms ...RouterMiddleware) httprouter.Handle {
	for _, m := range ms {
		if len(prefixes) == 0 {
			h = m(h)
		} else {
			t := h
			h = func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				for _, prefix := range prefixes {
					if strings.HasPrefix(r.URL.EscapedPath(), prefix) {
						m(t)(rw, r, p)
						return
					}
				}
				t(rw, r, p)
			}
		}
	}
	return h
}

// Middleware represents a middleware
type Middleware func(http.Handler) http.Handler

// RouterMiddleware represents a router middleware
type RouterMiddleware func(httprouter.Handle) httprouter.Handle

func handleBasicAuth(username, password string, rw http.ResponseWriter, r *http.Request) bool {
	if len(username) > 0 && len(password) > 0 {
		if u, p, ok := r.BasicAuth(); !ok || u != username || p != password {
			rw.Header().Set("WWW-Authenticate", "Basic Realm=Please enter your credentials")
			rw.WriteHeader(http.StatusUnauthorized)
			return true
		}
	}
	return false
}

// MiddlewareBasicAuth adds basic HTTP auth to a handler
func MiddlewareBasicAuth(username, password string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Basic auth
			if handleBasicAuth(username, password, rw, r) {
				return
			}

			// Next handler
			h.ServeHTTP(rw, r)
		})
	}
}

// RouterMiddlewareBasicAuth adds basic HTTP auth to a router handler
func RouterMiddlewareBasicAuth(username, password string) RouterMiddleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
			// Basic auth
			if handleBasicAuth(username, password, rw, r) {
				return
			}

			// Next handler
			h(rw, r, p)
		}
	}
}

func handleContentType(contentType string, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", contentType)
}

// MiddlewareContentType adds a content type to a handler
func MiddlewareContentType(contentType string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Content type
			handleContentType(contentType, rw)

			// Next handler
			h.ServeHTTP(rw, r)
		})
	}
}

// RouterMiddlewareContentType adds a content type to a router handler
func RouterMiddlewareContentType(contentType string) RouterMiddleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
			// Content type
			handleContentType(contentType, rw)

			// Next handler
			h(rw, r, p)
		}
	}
}

func handleHeaders(vs map[string]string, rw http.ResponseWriter) {
	for k, v := range vs {
		rw.Header().Set(k, v)
	}
}

// MiddlewareHeaders adds headers to a handler
func MiddlewareHeaders(vs map[string]string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Add headers
			handleHeaders(vs, rw)

			// Next handler
			h.ServeHTTP(rw, r)
		})
	}
}

// RouterMiddlewareHeaders adds headers to a router handler
func RouterMiddlewareHeaders(vs map[string]string) RouterMiddleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
			// Add headers
			handleHeaders(vs, rw)

			// Next handler
			h(rw, r, p)
		}
	}
}

func handleTimeout(timeout time.Duration, rw http.ResponseWriter, fn func()) {
	// Init context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Serve
	var done = make(chan bool)
	go func() {
		fn()
		done <- true
	}()

	// Wait for done or timeout
	for {
		select {
		case <-ctx.Done():
			astilog.Error(errors.Wrap(ctx.Err(), "astihttp: serving HTTP failed"))
			rw.WriteHeader(http.StatusGatewayTimeout)
			return
		case <-done:
			return
		}
	}
}

// MiddlewareTimeout adds a timeout to a handler
func MiddlewareTimeout(timeout time.Duration) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			handleTimeout(timeout, rw, func() { h.ServeHTTP(rw, r) })
		})
	}
}

// RouterMiddlewareTimeout adds a timeout to a router handler
func RouterMiddlewareTimeout(timeout time.Duration) RouterMiddleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
			handleTimeout(timeout, rw, func() { h(rw, r, p) })
		}
	}
}
