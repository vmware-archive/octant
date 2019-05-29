package overview

import (
	"context"
	"fmt"
	"regexp"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/overview/resourceviewer"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/discovery"
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

func (pf *pathFilter) String() string {
	return pf.path
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

type realGenerator struct {
	objectStore        objectstore.ObjectStore
	visitCache         resourceviewer.VisitCache
	pathMatcher        *pathMatcher
	clusterClient      cluster.ClientInterface
	printer            printer.Printer
	portForwardSvc     portforward.PortForwarder
	discoveryInterface discovery.DiscoveryInterface
	pluginStore        plugin.ManagerStore
}

// GeneratorOptions are additional options to pass a generator
type GeneratorOptions struct {
	LabelSet       *kLabels.Set
	PortForwardSvc portforward.PortForwarder
	PluginManager  *plugin.Manager
}

func newGenerator(objectStore objectstore.ObjectStore, di discovery.DiscoveryInterface, pm *pathMatcher, clusterClient cluster.ClientInterface, portForwardSvc portforward.PortForwarder, pluginStore plugin.ManagerStore) (*realGenerator, error) {
	p := printer.NewResource(objectStore, portForwardSvc, pluginStore)

	if err := AddPrintHandlers(p); err != nil {
		return nil, errors.Wrap(err, "add print handlers")
	}

	if pm == nil {
		return nil, errors.New("path matcher is nil")
	}

	visitCache, err := resourceviewer.NewVisitCache(100)
	if err != nil {
		return nil, errors.Wrap(err, "new visit cache")
	}

	return &realGenerator{
		objectStore:        objectStore,
		visitCache:         visitCache,
		discoveryInterface: di,
		pathMatcher:        pm,
		clusterClient:      clusterClient,
		portForwardSvc:     portForwardSvc,
		printer:            p,
		pluginStore:        pluginStore,
	}, nil
}

func (g *realGenerator) Generate(ctx context.Context, path, prefix, namespace string, opts GeneratorOptions) (component.ContentResponse, error) {
	ctx, span := trace.StartSpan(ctx, "Generate")
	defer span.End()

	pf, err := g.pathMatcher.Find(path)
	if err != nil {
		if err == errPathNotFound {
			return emptyContentResponse, api.NewNotFoundError(path)
		}
		return emptyContentResponse, err
	}

	q := queryer.New(g.objectStore, g.discoveryInterface)

	fields := pf.Fields(path)
	options := DescriberOptions{
		ObjectStore:        g.objectStore,
		VCache:             g.visitCache,
		Queryer:            q,
		Fields:             fields,
		Printer:            g.printer,
		LabelSet:           opts.LabelSet,
		PortForwardSvc:     opts.PortForwardSvc,
		PluginPrinter:      opts.PluginManager,
		PluginManagerStore: g.pluginStore,
	}

	cResponse, err := pf.describer.Describe(ctx, prefix, namespace, g.clusterClient, options)
	if err != nil {
		return emptyContentResponse, err
	}

	return cResponse, nil
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
		printer.ClusterRoleBindingListHandler,
		printer.ClusterRoleBindingHandler,
		printer.ConfigMapListHandler,
		printer.ConfigMapHandler,
		printer.CronJobListHandler,
		printer.CronJobHandler,
		printer.ClusterRoleListHandler,
		printer.ClusterRoleHandler,
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
		printer.PersistentVolumeClaimHandler,
		printer.PersistentVolumeClaimListHandler,
		printer.ServiceAccountListHandler,
		printer.ServiceAccountHandler,
		printer.ServiceHandler,
		printer.ServiceListHandler,
		printer.SecretHandler,
		printer.SecretListHandler,
		printer.StatefulSetHandler,
		printer.StatefulSetListHandler,
		printer.RoleBindingListHandler,
		printer.RoleBindingHandler,
		printer.RoleListHandler,
		printer.RoleHandler,
	}

	for _, handler := range handlers {
		if err := p.Handler(handler); err != nil {
			return err
		}
	}

	return nil
}
