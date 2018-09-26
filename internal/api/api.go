package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type notFoundResponse struct {
	Error errorResponse `json:"error,omitempty"`
}

// API is the API for the dashboard client
type API struct {
	mux *mux.Router
}

var _ http.Handler = (*API)(nil)

// New creates an instance of API.
func New(prefix string) *API {
	router := mux.NewRouter()
	s := router.PathPrefix(prefix).Subrouter()

	namespacesService := &namespaces{}
	s.Handle("/namespaces", namespacesService)

	navigationService := &navigation{}
	s.Handle("/navigation", navigationService)

	contentService := &content{}
	s.Handle("/content/{path:.*}", contentService)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		resp := &notFoundResponse{
			Error: errorResponse{
				Code:    http.StatusNotFound,
				Message: "not found",
			},
		}

		json.NewEncoder(w).Encode(resp)
	})

	return &API{
		mux: router,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
