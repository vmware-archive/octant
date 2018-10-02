package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/heptio/developer-dash/internal/cluster"
)

type namespacesResponse struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

type namespaces struct {
	nsClient cluster.NamespaceInterface
}

var _ http.Handler = (*namespaces)(nil)

func newNamespaces(nsClient cluster.NamespaceInterface) *namespaces {
	return &namespaces{
		nsClient: nsClient,
	}
}

func (n *namespaces) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	names, err := n.nsClient.Names()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	nr := &namespacesResponse{
		Namespaces: names,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		log.Printf("encoding namespaces error: %v", err)
	}
}
