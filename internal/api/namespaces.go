package api

import (
	"encoding/json"
	"net/http"
)

type namespacesResponse struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

type namespaces struct {
}

var _ http.Handler = (*namespaces)(nil)

func (n *namespaces) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nr := &namespacesResponse{
		Namespaces: []string{
			"default",
			"app-1",
			"app-2",
		},
	}

	json.NewEncoder(w).Encode(nr)
}
