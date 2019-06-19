package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/trace"

	"github.com/heptio/developer-dash/internal/octant"
	"github.com/heptio/developer-dash/internal/log"
)

type navSections interface {
	Sections(ctx context.Context, namespace string) ([]octant.Navigation, error)
}

type navigationResponse struct {
	Sections []octant.Navigation `json:"sections,omitempty"`
}

type navigation struct {
	navSections navSections
	logger      log.Logger
}

var _ http.Handler = (*navigation)(nil)

func newNavigation(ns navSections, logger log.Logger) *navigation {
	if logger == nil {
		logger = log.NopLogger()
	}

	return &navigation{
		navSections: ns,
		logger:      logger,
	}
}

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "api:navigation")
	defer span.End()

	if n.navSections == nil {
		RespondWithError(w, http.StatusInternalServerError,
			"unable to generate navigation sections", n.logger)
		return
	}

	vars := mux.Vars(r)
	namespace := vars["namespace"] // optional
	if namespace == "" {
		// Fallback to legacy query parameter
		namespace = "default"
	}

	n.logger.Debugf("navigation for namespace %s", namespace)

	ns, err := n.navSections.Sections(ctx, namespace)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"unable to generate navigation sections", n.logger)
		return
	}

	nr := navigationResponse{
		Sections: ns,
	}

	serveAsJSON(w, &nr, n.logger)
}
