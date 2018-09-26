package api

import (
	"encoding/json"
	"net/http"

	"github.com/heptio/developer-dash/internal/overview"
)

type namespacesResponse struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

type namespaces struct {
	overview overview.Interface
}

var _ http.Handler = (*namespaces)(nil)

func newNamespaces(o overview.Interface) *namespaces {
	return &namespaces{
		overview: o,
	}
}

func (n *namespaces) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	names, err := n.overview.Namespaces()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	nr := &namespacesResponse{
		Namespaces: names,
	}

	json.NewEncoder(w).Encode(nr)
}
