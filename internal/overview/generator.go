package overview

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/heptio/developer-dash/internal/queryer"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/pkg/errors"
)

type pathFilter struct {
	path      string
	describer Describer

	re *regexp.Regexp
}

func newPathFilter(path string, describer Describer) *pathFilter {
	re := regexp.MustCompile(fmt.Sprintf("^%s/?$", path))

	return &pathFilter{
		re:        re,
		path:      path,
		describer: describer,
	}
}

func (pf *pathFilter) Match(path string) bool {
	return pf.re.MatchString(path)
}

func (pf *pathFilter) Fields(path string) map[string]string {
	out := make(map[string]string)

	match := pf.re.FindStringSubmatch(path)
	for i, name := range pf.re.SubexpNames() {
		if i != 0 && name != "" {
			out[name] = match[i]
		}
	}

	return out
}

var contentNotFound = errors.Errorf("content not found")

type realGenerator struct {
	cache         cache.Cache
	queryer       queryer.Queryer
	pathFilters   []pathFilter
	clusterClient cluster.ClientInterface
	printer       printer.Printer

	mu sync.Mutex
}

func newGenerator(cache cache.Cache, q queryer.Queryer, pathFilters []pathFilter, clusterClient cluster.ClientInterface) (*realGenerator, error) {
	p := printer.NewResource(cache)

	if err := AddPrintHandlers(p); err != nil {
		return nil, errors.Wrap(err, "add print handlers")
	}

	return &realGenerator{
		cache:         cache,
		queryer:       q,
		pathFilters:   pathFilters,
		clusterClient: clusterClient,
		printer:       p,
	}, nil
}

func (g *realGenerator) Generate(ctx context.Context, path, prefix, namespace string) (component.ContentResponse, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, pf := range g.pathFilters {
		if !pf.Match(path) {
			continue
		}

		fields := pf.Fields(path)
		options := DescriberOptions{
			Cache:   g.cache,
			Queryer: g.queryer,
			Fields:  fields,
			Printer: g.printer,
		}

		cResponse, err := pf.describer.Describe(ctx, prefix, namespace, g.clusterClient, options)
		if err != nil {
			return emptyContentResponse, err
		}

		return cResponse, nil
	}

	fmt.Println("content not found for", path)
	return emptyContentResponse, contentNotFound
}

// PrinterHandler configures handlers for a printer.
type PrinterHandler interface {
	Handler(printFunc interface{}) error
}

// AddPrintHandlers adds print handlers to a printer.
func AddPrintHandlers(p PrinterHandler) error {
	handlers := []interface{}{
		printer.EventListHandler,
		printer.EventHandler,
		printer.ConfigMapListHandler,
		printer.ConfigMapHandler,
		printer.DeploymentHandler,
		printer.DeploymentListHandler,
		printer.ReplicaSetHandler,
		printer.ReplicaSetListHandler,
		printer.PodHandler,
		printer.PodListHandler,
		printer.ServiceHandler,
		printer.ServiceListHandler,
		printer.SecretHandler,
		printer.SecretListHandler,
	}

	for _, handler := range handlers {
		if err := p.Handler(handler); err != nil {
			return err
		}
	}

	return nil
}
