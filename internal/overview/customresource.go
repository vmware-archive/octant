package overview

import (
	"context"
	"fmt"
	"path"
	"sort"
	"sync"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/overview/resourceviewer"
	"github.com/heptio/developer-dash/internal/overview/yamlviewer"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kcache "k8s.io/client-go/tools/cache"
)

func customResourceDefinitionNames(ctx context.Context, o objectstore.ObjectStore) ([]string, error) {
	key := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	rawList, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "listing CRDs")
	}

	var list []string

	for _, object := range rawList {
		crd := &apiextv1beta1.CustomResourceDefinition{}

		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, crd); err != nil {
			return nil, errors.Wrap(err, "crd conversion failed")
		}

		list = append(list, crd.Name)
	}

	return list, nil
}

func customResourceDefinition(ctx context.Context, name string, o objectstore.ObjectStore) (*apiextv1beta1.CustomResourceDefinition, error) {
	key := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       name,
	}

	crd := &apiextv1beta1.CustomResourceDefinition{}
	if err := objectstore.GetAs(ctx, o, key, crd); err != nil {
		return nil, errors.Wrap(err, "get CRD from objecstore")
	}

	return crd, nil
}

type crdSectionDescriber struct {
	describers map[string]Describer
	path       string
	title      string

	mu sync.Mutex
}

var _ (Describer) = (*crdSectionDescriber)(nil)

func newCRDSectionDescriber(p, title string) *crdSectionDescriber {
	return &crdSectionDescriber{
		describers: make(map[string]Describer),
		path:       p,
		title:      title,
	}
}

func (csd *crdSectionDescriber) Add(name string, describer Describer) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	csd.describers[name] = describer
}

func (csd *crdSectionDescriber) Remove(name string) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	delete(csd.describers, name)
}

func (csd *crdSectionDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	var names []string
	for name := range csd.describers {
		names = append(names, name)
	}

	sort.Strings(names)

	list := component.NewList("Custom Resources", nil)

	for _, name := range names {
		resp, err := csd.describers[name].Describe(ctx, prefix, namespace, clusterClient, options)
		if err != nil {
			return emptyContentResponse, err
		}

		for i := range resp.Components {
			c := resp.Components[i]
			if !c.IsEmpty() {
				list.Add(c)
			}
		}
	}

	cr := component.ContentResponse{
		Components: []component.Component{list},
		Title:      component.TitleFromString(csd.title),
	}

	return cr, nil
}

func (csd *crdSectionDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(csd.path, csd),
	}
}

type crdListPrinter func(ctx context.Context, name, namespace string, crd *apiextv1beta1.CustomResourceDefinition, objects []*unstructured.Unstructured) (component.Component, error)

type crdListDescriptionOption func(*crdListDescriber)

type crdListDescriber struct {
	name    string
	path    string
	printer crdListPrinter
}

var _ (Describer) = (*crdListDescriber)(nil)

func newCRDListDescriber(name, path string, options ...crdListDescriptionOption) *crdListDescriber {
	d := &crdListDescriber{
		name:    name,
		path:    path,
		printer: printer.CustomResourceListHandler,
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (cld *crdListDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	crd, err := customResourceDefinition(ctx, cld.name, options.ObjectStore)
	if err != nil {
		return emptyContentResponse, err
	}

	objects, err := listCustomResources(ctx, crd, namespace, options.ObjectStore)
	if err != nil {
		return emptyContentResponse, err
	}

	table, err := cld.printer(ctx, cld.name, namespace, crd, objects)
	if err != nil {
		return emptyContentResponse, err
	}

	list := component.NewList(fmt.Sprintf("Custom Resources / %s", cld.name), []component.Component{
		table,
	})

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func listCustomResources(
	ctx context.Context,
	crd *apiextv1beta1.CustomResourceDefinition,
	namespace string,
	o objectstore.ObjectStore) ([]*unstructured.Unstructured, error) {
	if crd == nil {
		return nil, errors.New("crd is nil")
	}
	gvk := schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: crd.Spec.Version,
		Kind:    crd.Spec.Names.Kind,
	}

	apiVersion, kind := gvk.ToAPIVersionAndKind()

	key := objectstoreutil.Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
	}

	objects, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "listing custom resources for %q", crd.Name)
	}

	return objects, nil
}

func (cld *crdListDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(cld.path, cld),
	}
}

type crdPrinter func(ctx context.Context, crd *apiextv1beta1.CustomResourceDefinition, object *unstructured.Unstructured, options printer.Options) (component.Component, error)
type resourceViewerPrinter func(ctx context.Context, object *unstructured.Unstructured, o objectstore.ObjectStore, q queryer.Queryer) (component.Component, error)
type yamlPrinter func(runtime.Object) (*component.YAML, error)

type crdDescriberOption func(*crdDescriber)

type crdDescriber struct {
	path                  string
	name                  string
	summaryPrinter        crdPrinter
	resourceViewerPrinter resourceViewerPrinter
	yamlPrinter           yamlPrinter
}

var _ (Describer) = (*crdDescriber)(nil)

func newCRDDescriber(name, path string, options ...crdDescriberOption) *crdDescriber {
	d := &crdDescriber{
		path:                  path,
		name:                  name,
		summaryPrinter:        printer.CustomResourceHandler,
		resourceViewerPrinter: createCRDResourceViewer,
		yamlPrinter:           yamlviewer.ToComponent,
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (cd *crdDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	crd, err := customResourceDefinition(ctx, cd.name, options.ObjectStore)
	if err != nil {
		return emptyContentResponse, err
	}

	gvk := schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: crd.Spec.Version,
		Kind:    crd.Spec.Names.Kind,
	}

	apiVersion, kind := gvk.ToAPIVersionAndKind()

	key := objectstoreutil.Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       options.Fields["name"],
	}

	object, err := options.ObjectStore.Get(ctx, key)
	if err != nil {
		return emptyContentResponse, err
	}

	// TODO: shouldn't use the nil, should use the error
	if object == nil {
		return emptyContentResponse, err
	}

	title := component.Title(
		link.ForCustomResourceDefinition(cd.name, namespace),
		component.NewText(object.GetName()))

	cr := component.NewContentResponse(title)

	printOptions := printer.Options{
		ObjectStore:   options.ObjectStore,
		PluginPrinter: options.PluginManager,
	}

	summary, err := cd.summaryPrinter(ctx, crd, object, printOptions)
	if err != nil {
		return emptyContentResponse, err
	}
	summary.SetAccessor("summary")

	cr.Add(summary)

	resourceViewerComponent, err := cd.resourceViewerPrinter(ctx, object, options.ObjectStore, options.Queryer)
	if err != nil {
		return emptyContentResponse, err
	}

	resourceViewerComponent.SetAccessor("resourceViewer")
	cr.Add(resourceViewerComponent)

	yvComponent, err := cd.yamlPrinter(object)
	if err != nil {
		return emptyContentResponse, err
	}

	yvComponent.SetAccessor("yaml")
	cr.Add(yvComponent)

	tabs, err := options.PluginManager.Tabs(object)
	if err != nil {
		return emptyContentResponse, errors.Wrap(err, "getting tabs from plugins")
	}

	for _, tab := range tabs {
		tab.Contents.SetAccessor(tab.Name)
		cr.Add(&tab.Contents)
	}

	return *cr, nil
}

func (cd *crdDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(cd.path, cd),
	}
}

func createCRDResourceViewer(ctx context.Context, object *unstructured.Unstructured, o objectstore.ObjectStore, q queryer.Queryer) (component.Component, error) {
	logger := log.From(ctx)

	rv, err := resourceviewer.New(logger, o, resourceviewer.WithDefaultQueryer(q))
	if err != nil {
		return nil, err
	}

	return rv.Visit(ctx, object)
}

type objectHandler func(ctx context.Context, object *unstructured.Unstructured)

func watchCRDs(ctx context.Context, o objectstore.ObjectStore, crdAddFunc, crdDeleteFunc objectHandler) {
	handler := &kcache.ResourceEventHandlerFuncs{}

	if crdAddFunc != nil {
		handler.AddFunc = func(object interface{}) {
			u, ok := object.(*unstructured.Unstructured)
			if ok {
				crdAddFunc(ctx, u)
			}
		}
	}

	if crdDeleteFunc != nil {
		handler.DeleteFunc = func(object interface{}) {
			u, ok := object.(*unstructured.Unstructured)
			if ok {
				crdDeleteFunc(ctx, u)
			}
		}
	}

	key := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	logger := log.From(ctx)

	if err := o.Watch(key, handler); err != nil {
		logger.Errorf("crd watcher has failed: %v", err)
	}
}

func addCRD(ctx context.Context, name string, pm *pathMatcher, sectionDescriber *crdSectionDescriber) {
	logger := log.From(ctx)
	logger.Debugf("adding CRD %s", name)

	cld := newCRDListDescriber(name, crdListPath(name))

	sectionDescriber.Add(name, cld)

	for _, pf := range cld.PathFilters() {
		pm.Register(ctx, pf)
	}

	cd := newCRDDescriber(name, crdObjectPath(name))
	for _, pf := range cd.PathFilters() {
		pm.Register(ctx, pf)
	}
}

func deleteCRD(ctx context.Context, name string, pm *pathMatcher, sectionDescriber *crdSectionDescriber) {
	logger := log.From(ctx)
	logger.Debugf("deleting CRD %s", name)

	pm.Deregister(ctx, crdListPath(name))
	pm.Deregister(ctx, crdObjectPath(name))

	sectionDescriber.Remove(name)

}

func crdListPath(name string) string {
	return path.Join("/custom-resources", name)
}

func crdObjectPath(name string) string {
	return path.Join(crdListPath(name), resourceNameRegex)
}
