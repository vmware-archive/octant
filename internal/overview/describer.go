package overview

import (
	"context"
	"reflect"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/overview/logviewer"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/overview/resourceviewer"
	"github.com/heptio/developer-dash/internal/overview/yamlviewer"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

const (
	resourceNameRegex = "(?P<name>.*?)"
)

// LoaderFunc loads an object from the objectstore.
type LoaderFunc func(ctx context.Context, o objectstore.ObjectStore, namespace string, fields map[string]string) (*unstructured.Unstructured, error)

// DefaultLoader returns a loader that loads a single object from the objectstore
var DefaultLoader = func(objectStoreKey objectstoreutil.Key) LoaderFunc {
	return func(ctx context.Context, o objectstore.ObjectStore, namespace string, fields map[string]string) (*unstructured.Unstructured, error) {
		return loadObject(ctx, o, namespace, fields, objectStoreKey)
	}
}

// DescriberOptions provides options to describers
type DescriberOptions struct {
	Queryer        queryer.Queryer
	ObjectStore    objectstore.ObjectStore
	Fields         map[string]string
	Printer        printer.Printer
	LabelSet       *kLabels.Set
	PortForwardSvc portforward.PortForwarder
	PluginManager  printer.PluginPrinter
}

// Describer creates content.
type Describer interface {
	Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error)
	PathFilters() []pathFilter
}

type baseDescriber struct{}

func newBaseDescriber() *baseDescriber {
	return &baseDescriber{}
}

// ListDescriber describes a list of objects.
type ListDescriber struct {
	*baseDescriber

	path           string
	title          string
	listType       func() interface{}
	objectType     func() interface{}
	objectStoreKey objectstoreutil.Key
	isClusterWide  bool
}

// NewListDescriber creates an instance of ListDescriber.
func NewListDescriber(p, title string, objectStoreKey objectstoreutil.Key, listType, objectType func() interface{}, isClusterWide bool) *ListDescriber {
	return &ListDescriber{
		path:           p,
		title:          title,
		baseDescriber:  newBaseDescriber(),
		objectStoreKey: objectStoreKey,
		listType:       listType,
		objectType:     objectType,
		isClusterWide:  isClusterWide,
	}
}

// Describe creates content.
func (d *ListDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	if options.Printer == nil {
		return emptyContentResponse, errors.New("object list describer requires a printer")
	}

	// Pass through selector if provided to filter objects
	var key = d.objectStoreKey // copy
	key.Selector = options.LabelSet

	if d.isClusterWide {
		namespace = ""
	}

	objects, err := loadObjects(ctx, options.ObjectStore, namespace, options.Fields, []objectstoreutil.Key{key})
	if err != nil {
		return emptyContentResponse, err
	}

	list := component.NewList(d.title, nil)

	listType := d.listType()

	v := reflect.ValueOf(listType)
	f := reflect.Indirect(v).FieldByName("Items")

	// Convert unstructured objects to typed runtime objects
	for _, object := range objects {
		item := d.objectType()
		if err := scheme.Scheme.Convert(object, item, nil); err != nil {
			return emptyContentResponse, err
		}

		if err := copyObjectMeta(item, object); err != nil {
			return emptyContentResponse, err
		}

		newSlice := reflect.Append(f, reflect.ValueOf(item).Elem())
		f.Set(newSlice)
	}

	listObject, ok := listType.(runtime.Object)
	if !ok {
		return emptyContentResponse, errors.Errorf("expected list to be a runtime object. It was a %T",
			listType)
	}

	viewComponent, err := options.Printer.Print(ctx, listObject, options.PluginManager)
	if err != nil {
		return emptyContentResponse, err
	}

	if viewComponent != nil {

		if table, ok := viewComponent.(*component.Table); ok {
			if err := table.Config.Rows.Sort("Name"); err != nil {
				return emptyContentResponse, errors.Wrap(err, "sorting list by Name column")
			}
			list.Add(table)
		} else {
			list.Add(viewComponent)
		}
	}

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

// PathFilters returns path filters for this describer.
func (d *ListDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

// ObjectDescriber describes an object.
type ObjectDescriber struct {
	*baseDescriber

	path                  string
	baseTitle             string
	objectType            func() interface{}
	loaderFunc            LoaderFunc
	disableResourceViewer bool
}

// NewObjectDescriber creates an instance of ObjectDescriber.
func NewObjectDescriber(p, baseTitle string, loaderFunc LoaderFunc, objectType func() interface{}, disableResourceViewer bool) *ObjectDescriber {
	return &ObjectDescriber{
		path:                  p,
		baseTitle:             baseTitle,
		baseDescriber:         newBaseDescriber(),
		loaderFunc:            loaderFunc,
		objectType:            objectType,
		disableResourceViewer: disableResourceViewer,
	}
}

type tabFunc func(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options DescriberOptions) error

// Describe describes an object.
func (d *ObjectDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	logger := log.From(ctx)

	if options.Printer == nil {
		return emptyContentResponse, errors.New("object describer requires a printer")
	}

	if options.PluginManager == nil {
		return emptyContentResponse, errors.New("plugin manager is nil")
	}

	newObject, err := d.currentObject(ctx, namespace, options)
	if err != nil {
		return emptyContentResponse, err
	}

	accessor := meta.NewAccessor()
	objectName, _ := accessor.Name(newObject)

	title := append([]component.TitleComponent{}, component.NewText(d.baseTitle))
	if objectName != "" {
		title = append(title, component.NewText(objectName))
	}

	cr := component.NewContentResponse(title)

	tabFuncs := map[string]tabFunc{
		"summary":         d.addSummaryTab,
		"resource viewer": d.addResourceViewerTab,
		"yaml":            d.addYAMLViewerTab,
		"logs":            d.addLogsTab,
	}

	hasTabError := false
	for name, fn := range tabFuncs {
		if err := fn(ctx, newObject, cr, options); err != nil {
			hasTabError = true
			logger.With(
				"err", err,
				"tab-name", name).Errorf("generating object describer tab")
		}
	}

	if hasTabError {
		logger.With("tab-object", newObject).Errorf("unable to generate all tabs for object")
	}

	tabs, err := options.PluginManager.Tabs(newObject)
	if err != nil {
		return emptyContentResponse, errors.Wrap(err, "getting tabs from plugins")
	}

	for _, tab := range tabs {
		tab.Contents.SetAccessor(tab.Name)
		cr.Add(&tab.Contents)
	}

	return *cr, nil
}

func (d *ObjectDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

func (d *ObjectDescriber) currentObject(ctx context.Context, namespace string, options DescriberOptions) (runtime.Object, error) {
	object, err := d.loaderFunc(ctx, options.ObjectStore, namespace, options.Fields)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, api.NewNotFoundError(d.path)
		}
		return nil, err
	}

	if object == nil {
		return nil, api.NewNotFoundError(d.path)
	}

	item := d.objectType()

	if err := scheme.Scheme.Convert(object, item, nil); err != nil {
		return nil, err
	}

	if err := copyObjectMeta(item, object); err != nil {
		return nil, errors.Wrapf(err, "copying object metadata")
	}

	newObject, ok := item.(runtime.Object)
	if !ok {
		return nil, errors.Errorf("expected item to be a runtime object. It was a %T",
			item)
	}

	return newObject, nil
}

func (d *ObjectDescriber) addSummaryTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options DescriberOptions) error {
	vc, err := options.Printer.Print(ctx, object, options.PluginManager)
	if err != nil {
		return err
	}

	if vc == nil {
		return errors.Wrap(err, "unable to print a nil object")
	}

	vc.SetAccessor("summary")
	cr.Add(vc)

	return nil
}

func (d *ObjectDescriber) addResourceViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options DescriberOptions) error {
	logger := log.From(ctx)

	if !d.disableResourceViewer {
		rv, err := resourceviewer.New(logger, options.ObjectStore, resourceviewer.WithDefaultQueryer(options.Queryer))
		if err != nil {
			return err
		}

		resourceViewerComponent, err := rv.Visit(ctx, object)
		if err != nil {
			return err
		}

		resourceViewerComponent.SetAccessor("resourceViewer")
		cr.Add(resourceViewerComponent)
	}

	return nil
}

func (d *ObjectDescriber) addYAMLViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options DescriberOptions) error {
	yvComponent, err := yamlviewer.ToComponent(object)
	if err != nil {
		return err
	}

	yvComponent.SetAccessor("yaml")
	cr.Add(yvComponent)

	return nil
}

func (d *ObjectDescriber) addLogsTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options DescriberOptions) error {
	if isPod(object) {
		logsComponent, err := logviewer.ToComponent(object)
		if err != nil {
			return err
		}

		logsComponent.SetAccessor("logs")
		cr.Add(logsComponent)
	}

	return nil
}

func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return errors.Errorf("%T is not an object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return errors.Wrapf(err, "accessing type meta")
	}
	t.SetAPIVersion(from.GetAPIVersion())
	t.SetKind(from.GetObjectKind().GroupVersionKind().Kind)

	object.SetNamespace(from.GetNamespace())
	object.SetName(from.GetName())
	object.SetGenerateName(from.GetGenerateName())
	object.SetUID(from.GetUID())
	object.SetResourceVersion(from.GetResourceVersion())
	object.SetGeneration(from.GetGeneration())
	object.SetSelfLink(from.GetSelfLink())
	object.SetCreationTimestamp(from.GetCreationTimestamp())
	object.SetDeletionTimestamp(from.GetDeletionTimestamp())
	object.SetDeletionGracePeriodSeconds(from.GetDeletionGracePeriodSeconds())
	object.SetLabels(from.GetLabels())
	object.SetAnnotations(from.GetAnnotations())
	object.SetInitializers(from.GetInitializers())
	object.SetOwnerReferences(from.GetOwnerReferences())
	object.SetClusterName(from.GetClusterName())
	object.SetFinalizers(from.GetFinalizers())

	return nil
}

// SectionDescriber is a wrapper to combine content from multiple describers.
type SectionDescriber struct {
	path       string
	title      string
	describers []Describer
}

// NewSectionDescriber creates a SectionDescriber.
func NewSectionDescriber(p, title string, describers ...Describer) *SectionDescriber {
	return &SectionDescriber{
		path:       p,
		title:      title,
		describers: describers,
	}
}

// Describe generates content.
func (d *SectionDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	list := component.NewList(d.title, nil)

	for _, child := range d.describers {
		cResponse, err := child.Describe(ctx, prefix, namespace, clusterClient, options)
		if err != nil {
			return emptyContentResponse, err
		}

		for _, vc := range cResponse.Components {
			if nestedList, ok := vc.(*component.List); ok {
				for i := range nestedList.Config.Items {
					item := nestedList.Config.Items[i]
					if !item.IsEmpty() {
						list.Add(item)
					}
				}
			}
		}
	}

	cr := component.ContentResponse{
		Components: []component.Component{list},
		Title:      component.Title(component.NewText(d.title)),
	}

	return cr, nil
}

func (d *SectionDescriber) PathFilters() []pathFilter {
	pathFilters := []pathFilter{
		*newPathFilter(d.path, d),
	}

	for _, child := range d.describers {
		pathFilters = append(pathFilters, child.PathFilters()...)
	}

	return pathFilters
}

func isPod(object runtime.Object) bool {
	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return apiVersion == "v1" && kind == "Pod"
}
