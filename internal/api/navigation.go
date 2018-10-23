package api

import (
	"encoding/json"
	"net/http"

	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
)

type navigationResponse struct {
	Sections []*hcli.Navigation `json:"sections,omitempty"`
}

type navigation struct {
	sections []*hcli.Navigation
	logger   log.Logger
}

var _ http.Handler = (*navigation)(nil)

func newNavigation(sections []*hcli.Navigation, logger log.Logger) *navigation {
	return &navigation{
		sections: sections,
		logger:   logger,
	}
}

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nr := navigationResponse{
		Sections: n.sections,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		n.logger.Errorf("encoding navigation error: %v", err)
	}
}
