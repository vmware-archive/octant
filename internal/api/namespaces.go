package api

import (
	"encoding/json"
	"net/http"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
)

type namespacesResponse struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

type namespaces struct {
	nsClient cluster.NamespaceInterface
	logger   log.Logger
}

var _ http.Handler = (*namespaces)(nil)

func newNamespaces(nsClient cluster.NamespaceInterface, logger log.Logger) *namespaces {
	return &namespaces{
		nsClient: nsClient,
		logger:   logger,
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
		n.logger.Errorf("encoding namespaces: %v", err)
	}
}
