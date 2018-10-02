package overview

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
	mux *mux.Router
}

var _ http.Handler = (*handler)(nil)

func newHandler(prefix string) *handler {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix(prefix).Subrouter()

	paths := []string{
		"/workloads/cron-jobs",
		"/workloads/daemon-sets",
		"/workloads/deployments",
		"/workloads/jobs",
		"/workloads/pods",
		"/workloads/replica-sets",
		"/workloads/replication-controllers",
		"/workloads/stateful-sets",
		"/workloads",

		"/discovery-and-load-balancing/ingresses",
		"/discovery-and-load-balancing/services",
		"/discovery-and-load-balancing",

		"/config-and-storage/config-maps",
		"/config-and-storage/persistent-volume-claims",
		"/config-and-storage/secrets",
		"/config-and-storage",

		"/custom-resources",

		"/rbac/roles",
		"/rbac/role-bindings",
		"/rbac",

		"/events",

		"/",
	}

	for _, p := range paths {
		s.HandleFunc(p, stubHandler(p)).Methods("GET")

	}

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handler not found: %s", r.URL.String())
	})

	return &handler{
		mux: router,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	r := &notFoundResponse{
		Error: errorResponse{
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

type notFoundResponse struct {
	Error errorResponse `json:"error,omitempty"`
}

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func stubHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := newTable(name)
		t.Columns = []tableColumn{
			{Name: "foo", Accessor: "foo"},
			{Name: "bar", Accessor: "bar"},
			{Name: "baz", Accessor: "baz"},
		}

		t.Rows = []tableRow{
			{
				"foo": "r1c1",
				"bar": "r1c2",
				"baz": "r1c3",
			},
			{
				"foo": "r2c1",
				"bar": "r2c2",
				"baz": "r2c3",
			},
			{
				"foo": "r3c1",
				"bar": "r3c2",
				"baz": "r3c3",
			},
		}

		cr := &contentResponse{
			Contents: []content{
				t,
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := json.NewEncoder(w).Encode(cr); err != nil {
			log.Printf("encoding response: %v", err)
		}
	}
}
