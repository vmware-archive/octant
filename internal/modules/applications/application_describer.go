package applications

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/internal/resourceviewer"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var allowed = []schema.GroupVersionKind{
	gvk.CronJob,
	gvk.DaemonSet,
	gvk.Deployment,
	gvk.Pod,
	gvk.Job,
	gvk.ExtReplicaSet,
	gvk.ReplicationController,
	gvk.StatefulSet,
	gvk.Ingress,
	gvk.Service,
	gvk.ConfigMap,
	gvk.PersistentVolumeClaim,
	gvk.Secret,
	gvk.ServiceAccount,
}

// ApplicationDescriber describes an application.
type ApplicationDescriber struct {
	printer *printer.Resource

	overviewFactory       func(ctx context.Context, namespace string, options describer.Options) (component.Component, error)
	resourceViewerFactory func(ctx context.Context, namespace string, options describer.Options) (component.Component, error)
}

var _ describer.Describer = (*ApplicationDescriber)(nil)

// NewApplicationDescriber creates an instance of ApplicationDescriber.
func NewApplicationDescriber(dashConfig config.Dash) *ApplicationDescriber {
	p := printer.NewResource(dashConfig)
	d := &ApplicationDescriber{
		printer:               p,
		overviewFactory:       overviewFactory,
		resourceViewerFactory: resourceViewerFactory,
	}

	return d
}

// Describe creates an application content response. It includes a list overview and
// a resource viewer for all the objects in an application.
func (a *ApplicationDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	d, err := descriptorFromFields(options.Fields)
	if err != nil {
		return component.EmptyContentResponse, errors.Wrap(err, "extract descriptor from fields")
	}

	overview, err := a.overviewFactory(ctx, namespace, options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	resourceViewer, err := a.resourceViewerFactory(ctx, namespace, options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	resp := component.ContentResponse{
		Title: component.TitleFromString(d.applicationTitle()),
		Components: []component.Component{
			overview,
			resourceViewer,
		},
	}

	return resp, nil
}

// PathFilters creates PathFilters for an application. The path for an application
// is /app-name/app-instance/app-version.
func (a *ApplicationDescriber) PathFilters() []describer.PathFilter {
	return []describer.PathFilter{
		*describer.NewPathFilter("/(?P<name>[^/]*)/(?P<instance>[^/]*)/(?P<version>[^/]*)", a),
	}
}

// Reset does nothing
func (a ApplicationDescriber) Reset(ctx context.Context) error {
	return nil
}

func loadAppObjects(ctx context.Context, dashConfig config.Dash, namespace, name, instance, version string) (*unstructured.UnstructuredList, error) {
	out := &unstructured.UnstructuredList{}

	ch := make(chan *unstructured.Unstructured)
	childrenProcessed := make(chan bool, 1)
	go func() {
		for child := range ch {
			if child == nil {
				continue
			}
			out.Items = append(out.Items, *child)
		}
		childrenProcessed <- true
	}()

	discoveryClient, err := dashConfig.ClusterClient().DiscoveryClient()
	if err != nil {
		return nil, err
	}

	resourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		//TODO: determine the best way to handle these types of errors for all resources, not just metrics.
		if discovery.IsGroupDiscoveryFailedError(err) {
			logger := log.From(ctx)
			logger.Debugf("preferred resources: %w", err)
		} else {
			return nil, err
		}
	}

	var g errgroup.Group

	sem := semaphore.NewWeighted(5)

	for resourceListIndex := range resourceLists {
		resourceList := resourceLists[resourceListIndex]
		if resourceList == nil {
			continue
		}

		for i := range resourceList.APIResources {
			apiResource := resourceList.APIResources[i]
			if !apiResource.Namespaced {
				continue
			}

			gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
			if err != nil {
				return nil, err
			}

			found := false
			for j := range allowed {
				if allowed[j].Group == gv.Group &&
					allowed[j].Version == gv.Version &&
					allowed[j].Kind == apiResource.Kind {
					found = true
				}
			}

			if !found {
				continue
			}

			key := store.Key{
				Namespace:  namespace,
				APIVersion: resourceList.GroupVersion,
				Kind:       apiResource.Kind,
				Selector: &labels.Set{
					appLabelName:     name,
					appLabelInstance: instance,
					appLabelVersion:  version,
				},
			}

			g.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return err
				}
				defer sem.Release(1)
				objects, _, err := dashConfig.ObjectStore().List(ctx, key)
				if err != nil {
					return errors.Wrapf(err, "unable to retrieve %+v", key)
				}

				for objectIndex := range objects.Items {
					ch <- &objects.Items[objectIndex]
				}

				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		if err != context.Canceled {
			return nil, errors.Wrap(err, "find children")
		}
	}

	close(ch)
	<-childrenProcessed
	close(childrenProcessed)

	sort.Slice(out.Items, func(i, j int) bool {
		if out.Items[i].GetAPIVersion() < out.Items[j].GetAPIVersion() {
			return true
		}
		if out.Items[i].GetAPIVersion() > out.Items[j].GetAPIVersion() {
			return false
		}
		if out.Items[i].GetKind() < out.Items[j].GetKind() {
			return true
		}
		if out.Items[i].GetKind() > out.Items[j].GetKind() {
			return false
		}

		return out.Items[i].GetName() < out.Items[j].GetName()
	})

	return out, nil
}

func overviewFactory(ctx context.Context, namespace string, options describer.Options) (component.Component, error) {
	d, err := descriptorFromFields(options.Fields)
	if err != nil {
		return nil, errors.Wrap(err, "extract descriptor from fields")
	}

	rootDescriber := describer.NamespacedOverview()

	options.LabelSet = &labels.Set{
		appLabelName:     d.name,
		appLabelInstance: d.instance,
		appLabelVersion:  d.version,
	}

	overview, err := rootDescriber.Component(ctx, namespace, options)
	if err != nil {
		return nil, err
	}

	overview.SetTitleText("Overview")
	overview.SetAccessor("overview")

	return overview, nil
}

func resourceViewerFactory(ctx context.Context, namespace string, options describer.Options) (component.Component, error) {
	d, err := descriptorFromFields(options.Fields)
	if err != nil {
		return nil, errors.Wrap(err, "extract descriptor from fields")
	}

	objects, err := loadAppObjects(ctx, options.Dash, namespace, d.name, d.instance, d.version)
	if err != nil {
		return nil, err
	}

	handler, err := resourceviewer.NewHandler(options.Dash)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := options.Dash.ClusterClient().DiscoveryClient()
	if err != nil {
		return nil, err
	}
	q := queryer.New(options.Dash.ObjectStore(), discoveryClient)

	rv, err := resourceviewer.New(options.Dash, resourceviewer.WithDefaultQueryer(options.Dash, q))
	if err != nil {
		return nil, err
	}

	for i := range objects.Items {
		if err := rv.Visit(ctx, &objects.Items[i], handler); err != nil {
			return nil, err
		}
	}

	return resourceviewer.GenerateComponent(ctx, handler, "")
}
