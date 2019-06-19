package configuration

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/octant"
	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/event"
	"github.com/heptio/developer-dash/internal/kubeconfig"
	"github.com/heptio/developer-dash/internal/log"
)

// kubeContextsResponse is a response for current kube contexts.
type kubeContextsResponse struct {
	Contexts       []kubeconfig.Context `json:"contexts"`
	CurrentContext string               `json:"currentContext"`
}

// updateCurrentContextRequest is a request to update the current context.
type updateCurrentContextRequest struct {
	RequestedContext string `json:"requestedContext"`
}

// updateCurrentContextHandler updates the current context.
type updateCurrentContextHandler struct {
	logger            log.Logger
	contextUpdateFunc func(name string) error
}

var _ http.Handler = (*updateCurrentContextHandler)(nil)

func (h *updateCurrentContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	defer func() {
		cErr := r.Body.Close()
		if cErr != nil {
			h.logger.WithErr(cErr).Errorf("unable to close request body")
		}
	}()

	var req updateCurrentContextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithErr(err).Errorf("decoding update context request")
		api.RespondWithError(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := h.contextUpdateFunc(req.RequestedContext); err != nil {
		h.logger.WithErr(err).Errorf("unable to update context")
		api.RespondWithError(w, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

const (
	// eventTypeKubeConfig is an event for updating kube contexts on the front end.
	eventTypeKubeConfig octant.EventType = "kubeConfig"
)

type kubeContextGenerationOption func(generator *kubeContextGenerator)

// kubeContextGenerator generates kube contexts for the front end.
type kubeContextGenerator struct {
	ConfigLoader kubeconfig.Loader
	DashConfig   config.Dash
}

var _ octant.Generator = (*kubeContextGenerator)(nil)

func newKubeContextGenerator(dashConfig config.Dash, options ...kubeContextGenerationOption) *kubeContextGenerator {
	kcg := &kubeContextGenerator{
		ConfigLoader: kubeconfig.NewFSLoader(),
		DashConfig:   dashConfig,
	}

	for _, option := range options {
		option(kcg)
	}

	return kcg
}

func (g *kubeContextGenerator) Event(ctx context.Context) (octant.Event, error) {
	configPath := g.DashConfig.KubeConfigPath()

	kubeConfig, err := g.ConfigLoader.Load(configPath)
	if err != nil {
		return octant.Event{}, errors.Wrap(err, "unable to load kube config")
	}

	currentContext := g.DashConfig.ContextName()
	if currentContext == "" {
		currentContext = kubeConfig.CurrentContext
	}

	resp := kubeContextsResponse{
		CurrentContext: currentContext,
		Contexts:       kubeConfig.Contexts,
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		return octant.Event{}, errors.Wrap(err, "encoding kube config data")
	}

	e := octant.Event{
		Type: eventTypeKubeConfig,
		Data: data,
	}

	return e, nil
}

func (kubeContextGenerator) ScheduleDelay() time.Duration {
	return event.DefaultScheduleDelay
}

func (kubeContextGenerator) Name() string {
	return "kubeConfig"
}
