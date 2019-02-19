package overview

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/heptio/developer-dash/internal/queryer"
	"k8s.io/apimachinery/pkg/labels"

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

// Fields extracts parameters from the request path.
// In practice, this finds the field "name" for an object request.
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
	queryer       queryer.Queryer // Queryer is used by the ResourceViewer and should not be filtered
	pathFilters   []pathFilter
	clusterClient cluster.ClientInterface
	printer       printer.Printer

	mu sync.Mutex
}

// GeneratorOptions are additional options to pass a generator
type GeneratorOptions struct {
	Selector labels.Selector
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

func (g *realGenerator) Generate(ctx context.Context, path, prefix, namespace string, opts GeneratorOptions) (component.ContentResponse, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, pf := range g.pathFilters {
		if !pf.Match(path) {
			continue
		}

		fields := pf.Fields(path)
		options := DescriberOptions{
			Cache:    g.cache,
			Queryer:  g.queryer,
			Fields:   fields,
			Printer:  g.printer,
			Selector: opts.Selector,
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
		printer.CronJobListHandler,
		printer.CronJobHandler,
		printer.DaemonSetListHandler,
		printer.DaemonSetHandler,
		printer.DeploymentHandler,
		printer.DeploymentListHandler,
		printer.IngressListHandler,
		printer.IngressHandler,
		printer.JobListHandler,
		printer.JobHandler,
		printer.ReplicaSetHandler,
		printer.ReplicaSetListHandler,
		printer.ReplicationControllerHandler,
		printer.ReplicationControllerListHandler,
		printer.PodHandler,
		printer.PodListHandler,
		printer.ServiceHandler,
		printer.ServiceListHandler,
		printer.SecretHandler,
		printer.SecretListHandler,
		printer.StatefulSetHandler,
		printer.StatefulSetListHandler,
	}

	for _, handler := range handlers {
		if err := p.Handler(handler); err != nil {
			return err
		}
	}

	return nil
}
