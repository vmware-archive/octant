package api

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
)

type contentHandler struct {
	modulePaths map[string]module.Module
	logger      log.Logger
	prefix      string
}

func (h *contentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range h.modulePaths {
		p := path.Join(h.prefix, k)
		if strings.HasPrefix(r.URL.Path, p) {
			ctx := log.WithLoggerContext(r.Context(), h.logger)
			contentPath := strings.TrimPrefix(r.URL.Path, h.modulePrefix(v))
			namespace := r.URL.Query().Get("namespace")
			poll := r.URL.Query().Get("poll")

			if poll != "" {
				h.handlePoll(ctx, poll, namespace, w, r, v)
				return
			}

			resp, err := v.Content(ctx, contentPath, h.prefix, namespace)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error(), h.logger)
				return
			}

			serveAsJSON(w, resp, h.logger)
			return
		}
	}
}

func (h *contentHandler) handlePoll(ctx context.Context, poll, namespace string, w http.ResponseWriter, r *http.Request, m module.Module) {
	eventTimeout := defaultEventTimeout
	timeout, err := strconv.Atoi(poll)
	if err == nil {
		eventTimeout = time.Duration(timeout) * time.Second
	}

	cs := contentStreamer{
		generatorFn:  m.Content,
		w:            w,
		path:         h.contentPath(r, m),
		prefix:       h.prefix,
		namespace:    namespace,
		streamFn:     stream,
		eventTimeout: eventTimeout,
		logger:       h.logger,
	}

	cs.content(ctx)

}

func (h *contentHandler) modulePrefix(m module.Module) string {
	return path.Join(h.prefix, "content", m.ContentPath())
}

func (h *contentHandler) contentPath(r *http.Request, m module.Module) string {
	return strings.TrimPrefix(r.URL.Path, h.modulePrefix(m))
}
