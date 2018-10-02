package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/heptio/developer-dash/internal/hcli"
)

type navigationResponse struct {
	Sections []*hcli.Navigation `json:"sections,omitempty"`
}

type navigation struct {
	sections []*hcli.Navigation
}

var _ http.Handler = (*navigation)(nil)

func newNavigation(sections []*hcli.Navigation) *navigation {
	return &navigation{
		sections: sections,
	}
}

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nr := navigationResponse{
		Sections: n.sections,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		log.Printf("encoding navigation error: %v", err)
	}
}
