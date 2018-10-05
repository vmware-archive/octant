package overview

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Error errorMessage `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Printf("encoding response error: %v", err)
	}
}

type handler struct {
	mux       *mux.Router
	generator generator
}

var _ http.Handler = (*handler)(nil)

func newHandler(prefix string, g generator) *handler {
	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, prefix)
		namespace := r.URL.Query().Get("namespace")

		contents, err := g.Generate(path, prefix, namespace)
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		cr := &contentResponse{
			Contents: contents,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := json.NewEncoder(w).Encode(cr); err != nil {
			log.Printf("encoding response: %v", err)
		}
	})

	return &handler{
		mux:       router,
		generator: g,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
