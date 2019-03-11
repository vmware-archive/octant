package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	defaultEventTimeout = 5 * time.Second
)

type notFound interface {
	NotFound() bool
	Path() string
}

type eventGenerator interface {
	// Generate generates an event or an error.
	Generate(context.Context) (event, error)
	// RunEvery schedules the event to run every x.
	RunEvery() time.Duration
	// Name is the generator name.
	Name() string
}

// contentEventGenerator generates events that contain components.
type contentEventGenerator struct {
	// generatorFn is a function that generates the component.
	generatorFn func(ctx context.Context, path, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error)
	// path is the path to the content.
	path string
	// prefix the API path prefix. It could be prepended to the path to create
	// a resolvable path.
	prefix string
	// namespace is the current namespace.
	namespace string
	// selector is a label selector to filter any content.
	selector labels.Selector
	// runEvery is how often the event generator should be run.
	runEvery time.Duration
}

// Generate generates an event from a component using `generatorFn` and wraps it in a
// `dashResponse`.
func (g *contentEventGenerator) Generate(ctx context.Context) (event, error) {
	resp, err := g.generatorFn(ctx, g.path, g.prefix, g.namespace, module.ContentOptions{Selector: g.selector})
	if err != nil {
		return event{}, err
	}

	dr := dashResponse{
		Content: resp,
	}

	data, err := json.Marshal(dr)
	if err != nil {
		return event{}, err
	}

	return event{data: data}, nil
}

func (g *contentEventGenerator) RunEvery() time.Duration {
	return g.runEvery
}

func (g *contentEventGenerator) Name() string {
	return "content"
}

// navigationEventGenerator generates events to update navigation.
type navigationEventGenerator struct {
	// modules is a list of modules to query for events.
	modules   []module.Module
	namespace string
}

func (g *navigationEventGenerator) Generate(ctx context.Context) (event, error) {
	ans := newAPINavSections(g.modules)

	ns, err := ans.Sections(ctx, g.namespace)
	if err != nil {
		return event{}, err
	}

	nr := navigationResponse{
		Sections: ns,
	}

	data, err := json.Marshal(nr)
	if err != nil {
		return event{}, err
	}

	return event{name: "navigation", data: data}, nil
}

func (g *navigationEventGenerator) RunEvery() time.Duration {
	return 5 * time.Second
}

func (g *navigationEventGenerator) Name() string {
	return "navigation"
}

type namespaceEventGenerator struct {
	nsClient cluster.NamespaceInterface
}

func (g *namespaceEventGenerator) Generate(ctx context.Context) (event, error) {
	if g.nsClient == nil {
		return event{}, errors.New("unable to query namespaces, client is nil")
	}

	names, err := g.nsClient.Names()
	if err != nil {
		initialNamespace := g.nsClient.InitialNamespace()
		names = []string{initialNamespace}
	}

	nr := &namespacesResponse{Namespaces: names}
	data, err := json.Marshal(nr)
	if err != nil {
		return event{}, errors.New("unable to marshal namespaces")
	}

	return event{name: "namespaces", data: data}, nil
}

func (g *namespaceEventGenerator) RunEvery() time.Duration {
	return 5 * time.Second
}

func (g *namespaceEventGenerator) Name() string {
	return "namespace"
}

type dashResponse struct {
	Content component.ContentResponse `json:"content,omitempty"`
}

type streamFn func(ctx context.Context, w http.ResponseWriter, ch chan event)

type contentStreamer struct {
	eventGenerators []eventGenerator
	w               http.ResponseWriter
	streamFn        streamFn
	logger          log.Logger
	contentPath     string
}

func (cs *contentStreamer) content(ctx context.Context) error {
	ch := make(chan event, 1)

	if cs.eventGenerators == nil {
		return errors.Errorf("event generators are not configured")
	}

	if cs.streamFn == nil {
		return errors.Errorf("stream function is not configured")
	}

	if cs.logger == nil {
		return errors.Errorf("logger is not configured")
	}

	var wg sync.WaitGroup
	wg.Add(len(cs.eventGenerators))

	for _, eg := range cs.eventGenerators {
		go func(ctx context.Context, eg eventGenerator, ch chan<- event) {
			defer wg.Done()
			timer := time.NewTimer(0)
			isRunning := true

			for isRunning {
				select {
				case <-ctx.Done():
					isRunning = false
				case <-timer.C:
					now := time.Now()
					e, err := eg.Generate(ctx)
					if err != nil {
						if nfe, ok := err.(notFound); ok && nfe.NotFound() {
							cs.logger.With(
								"path", nfe.Path(),
							).Infof("content not found")
							break
						}

						cs.logger.Errorf("event generator error: %v", err)

						// This could be one time error, or it could be a huge failure.
						// Either way, log, and move on. If this becomes a problem,
						// a circuit breaker or some other pattern could be employed here.
						break
					}

					cs.logger.With(
						"elapsed", time.Since(now),
						"generator", eg.Name(),
						"contentPath", cs.contentPath,
					).Debugf("generate complete")
					ch <- e

					nextTick := eg.RunEvery()
					if nextTick == 0 {
						isRunning = false
					} else {
						timer.Reset(nextTick)
					}
				}

			}
		}(ctx, eg, ch)
	}

	cs.streamFn(ctx, cs.w, ch)

	wg.Wait()
	close(ch)

	return nil
}

type event struct {
	name string
	data []byte
}

func stream(ctx context.Context, w http.ResponseWriter, ch chan event) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server sent events are unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	isStreaming := true

	for isStreaming {
		select {
		case <-ctx.Done():
			isStreaming = false
		case e := <-ch:
			if e.name != "" {
				fmt.Fprintf(w, "event: %s\n", e.name)
			}
			fmt.Fprintf(w, "data: %s\n\n", string(e.data))
			flusher.Flush()
		}
	}
}
