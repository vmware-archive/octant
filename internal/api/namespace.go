package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/heptio/developer-dash/internal/module"
)

type namespace struct {
	moduleManager module.ManagerInterface
}

func newNamespace(moduleManager module.ManagerInterface) *namespace {
	return &namespace{
		moduleManager: moduleManager,
	}
}

type namespaceRequest struct {
	Namespace string `json:"namespace,omitempty"`
}

func (n *namespace) update(w http.ResponseWriter, r *http.Request) {
	var nr namespaceRequest

	err := json.NewDecoder(r.Body).Decode(&nr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode request")
		return
	}

	if nr.Namespace == "" {
		respondWithError(w, http.StatusBadRequest, "unable to decode request")
		return
	}

	n.moduleManager.SetNamespace(nr.Namespace)

	w.WriteHeader(http.StatusNoContent)
}

type namespaceResponse struct {
	Namespace string `json:"namespace,omitempty"`
}

func (n *namespace) read(w http.ResponseWriter, r *http.Request) {
	ns := n.moduleManager.GetNamespace()
	nr := &namespaceResponse{
		Namespace: ns,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		log.Printf("encoding namespace error: %v", err)
	}
}
