package api

import (
	"net/http"

	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/mime"
)

type navSections interface {
	Sections() ([]*hcli.Navigation, error)
}

type navigationResponse struct {
	Sections []*hcli.Navigation `json:"sections,omitempty"`
}

type navigation struct {
	navSections navSections
	logger      log.Logger
}

var _ http.Handler = (*navigation)(nil)

func newNavigation(ns navSections, logger log.Logger) *navigation {
	return &navigation{
		navSections: ns,
		logger:      logger,
	}
}

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if n.navSections == nil {
		msg := map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "unable to generate navigation sections",
		}

		w.Header().Set("Content-Type", mime.JSONContentType)
		w.WriteHeader(http.StatusInternalServerError)

		serveAsJSON(w, msg, n.logger)
		return
	}
	ns, err := n.navSections.Sections()
	if err != nil {
		msg := map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "unable to generate navigation sections",
		}

		w.Header().Set("Content-Type", mime.JSONContentType)
		w.WriteHeader(http.StatusInternalServerError)

		serveAsJSON(w, msg, n.logger)
	}

	nr := navigationResponse{
		Sections: ns,
	}

	serveAsJSON(w, &nr, n.logger)
}
