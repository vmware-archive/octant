package overview

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/log"
)

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Error errorMessage `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string, logger log.Logger) {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		logger.Errorf("encoding response error: %v", err)
	}
}

type handler struct {
	mux       *mux.Router
	generator generator
	streamFn  streamFn
}

var _ http.Handler = (*handler)(nil)

func newHandler(prefix string, g generator, sfn streamFn, logger log.Logger) *handler {
	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := log.WithLoggerContext(r.Context(), logger)
		path := strings.TrimPrefix(r.URL.Path, prefix)
		namespace := r.URL.Query().Get("namespace")
		poll := r.URL.Query().Get("poll")

		if poll != "" {
			var eventTimeout time.Duration
			timeout, err := strconv.Atoi(poll)
			if err != nil {
				eventTimeout = defaultEventTimeout
			} else {
				eventTimeout = time.Duration(timeout) * time.Second
			}

			cs := contentStreamer{
				generator:    g,
				w:            w,
				path:         path,
				prefix:       prefix,
				namespace:    namespace,
				streamFn:     stream,
				eventTimeout: eventTimeout,
				logger:       logger,
			}

			cs.content(ctx)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		cResponse, err := g.Generate(ctx, path, prefix, namespace)
		if err != nil {
			switch {
			case err == contentNotFound:
				respondWithError(w, http.StatusNotFound, err.Error(), logger)
			default:
				respondWithError(w, http.StatusInternalServerError, err.Error(), logger)
			}
			return
		}

		cr := &cResponse

		if err := json.NewEncoder(w).Encode(cr); err != nil {
			logger.Errorf("encoding response: %v", err)
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
