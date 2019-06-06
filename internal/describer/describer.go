package describer

import (
	"context"
	"reflect"
	"sort"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"

	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/link"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/modules/overview/logviewer"
	"github.com/heptio/developer-dash/internal/modules/overview/printer"
	"github.com/heptio/developer-dash/internal/modules/overview/resourceviewer"
	"github.com/heptio/developer-dash/internal/modules/overview/yamlviewer"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

// EmptyContentResponse is an empty content response.
var EmptyContentResponse = component.ContentResponse{}

type ObjectLoaderFactory struct {
	dashConfig config.Dash
}

func NewObjectLoaderFactory(dashConfig config.Dash) *ObjectLoaderFactory {
	return &ObjectLoaderFactory{
		dashConfig: dashConfig,
	}
}

func (f *ObjectLoaderFactory) LoadObject(ctx context.Context, namespace string, fields map[string]string, objectStoreKey objectstoreutil.Key) (*unstructured.Unstructured, error) {
	return LoadObject(ctx, f.dashConfig.ObjectStore(), namespace, fields, objectStoreKey)
}

func (f *ObjectLoaderFactory) LoadObjects(ctx context.Context, namespace string, fields map[string]string, objectStoreKeys []objectstoreutil.Key) ([]*unstructured.Unstructured, error) {
	return LoadObjects(ctx, f.dashConfig.ObjectStore(), namespace, fields, objectStoreKeys)
}

// loadObject loads a single object from the object store.
func LoadObject(ctx context.Context, objectStore objectstore.ObjectStore, namespace string, fields map[string]string, objectStoreKey objectstoreutil.Key) (*unstructured.Unstructured, error) {
	objectStoreKey.Namespace = namespace

	if name, ok := fields["name"]; ok && name != "" {
		objectStoreKey.Name = name
	}

	object, err := objectStore.Get(ctx, objectStoreKey)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// loadObjects loads objects from the object store sorted by their name.
func LoadObjects(ctx context.Context, objectStore objectstore.ObjectStore, namespace string, fields map[string]string, objectStoreKeys []objectstoreutil.Key) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	for _, objectStoreKey := range objectStoreKeys {
		objectStoreKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			objectStoreKey.Name = name
		}

		storedObjects, err := objectStore.List(ctx, objectStoreKey)
		if err != nil {
			return nil, err
		}

		objects = append(objects, storedObjects...)
	}

	sort.SliceStable(objects, func(i, j int) bool {
		a, b := objects[i], objects[j]
		return a.GetName() < b.GetName()
	})

	return objects, nil
}

// LoaderFunc loads an object from the object store.
type LoaderFunc func(ctx context.Context, o objectstore.ObjectStore, namespace string, fields map[string]string) (*unstructured.Unstructured, error)

// Options provides options to describers
type Options struct {
	config.Dash

	Queryer  queryer.Queryer
	Fields   map[string]string
	Printer  printer.Printer
	LabelSet *kLabels.Set
	Link     link.Interface
	ComponentCache resourceviewer.ComponentCache

	LoadObjects func(ctx context.Context, namespace string, fields map[string]string, objectStoreKeys []objectstoreutil.Key) ([]*unstructured.Unstructured, error)
	LoadObject  func(ctx context.Context, namespace string, fields map[string]string, objectStoreKey objectstoreutil.Key) (*unstructured.Unstructured, error)
}

// Describer creates content.
type Describer interface {
	Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error)
	PathFilters() []PathFilter
}

type base struct{}

func newBaseDescriber() *base {
	return &base{}
}

// List describes a list of objects.
type List struct {
	*base

	path           string
	title          string
	listType       func() interface{}
	objectType     func() interface{}
	objectStoreKey objectstoreutil.Key
	isClusterWide  bool
}

// NewList creates an instance of List.
func NewList(p, title string, objectStoreKey objectstoreutil.Key, listType, objectType func() interface{}, isClusterWide bool) *List {
	return &List{
		path:           p,
		title:          title,
		base:           newBaseDescriber(),
		objectStoreKey: objectStoreKey,
		listType:       listType,
		objectType:     objectType,
		isClusterWide:  isClusterWide,
	}
}

// Describe creates content.
func (d *List) Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error) {
	if options.Printer == nil {
		return EmptyContentResponse, errors.New("object list Describer requires a printer")
	}

	// Pass through selector if provided to filter objects
	var key = d.objectStoreKey // copy
	key.Selector = options.LabelSet

	if d.isClusterWide {
		namespace = ""
	}

	objects, err := options.LoadObjects(ctx, namespace, options.Fields, []objectstoreutil.Key{key})
	if err != nil {
		return EmptyContentResponse, err
	}

	list := component.NewList(d.title, nil)

	listType := d.listType()

	v := reflect.ValueOf(listType)
	f := reflect.Indirect(v).FieldByName("Items")

	// Convert unstructured objects to typed runtime objects
	for _, object := range objects {
		item := d.objectType()
		if err := scheme.Scheme.Convert(object, item, nil); err != nil {
			return EmptyContentResponse, err
		}

		if err := copyObjectMeta(item, object); err != nil {
			return EmptyContentResponse, err
		}

		newSlice := reflect.Append(f, reflect.ValueOf(item).Elem())
		f.Set(newSlice)
	}

	listObject, ok := listType.(runtime.Object)
	if !ok {
		return EmptyContentResponse, errors.Errorf("expected list to be a runtime object. It was a %T",
			listType)
	}

	viewComponent, err := options.Printer.Print(ctx, listObject, options.PluginManager())
	if err != nil {
		return EmptyContentResponse, err
	}

	if viewComponent != nil {
		if table, ok := viewComponent.(*component.Table); ok {
			list.Add(table)
		} else {
			list.Add(viewComponent)
		}
	}

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

// PathFilters returns path filters for this Describer.
func (d *List) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}

// Object describes an object.
type Object struct {
	*base

	path                  string
	baseTitle             string
	objectType            func() interface{}
	objectStoreKey        objectstoreutil.Key
	disableResourceViewer bool
	tabFuncDescriptors    []tabFuncDescriptor
}

// NewObjectDescriber creates an instance of Object.
func NewObjectDescriber(p, baseTitle string, objectStoreKey objectstoreutil.Key, objectType func() interface{}, disableResourceViewer bool) *Object {
	o := &Object{
		path:                  p,
		baseTitle:             baseTitle,
		base:                  newBaseDescriber(),
		objectStoreKey:        objectStoreKey,
		objectType:            objectType,
		disableResourceViewer: disableResourceViewer,
	}

	o.tabFuncDescriptors = []tabFuncDescriptor{
		{name: "summary", tabFunc: o.addSummaryTab},
		{name: "resource viewer", tabFunc: o.addResourceViewerTab},
		{name: "yaml", tabFunc: o.addYAMLViewerTab},
		{name: "logs", tabFunc: o.addLogsTab},
	}

	return o
}

type tabFunc func(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error

type tabFuncDescriptor struct {
	name    string
	tabFunc tabFunc
}

// Describe describes an object.
func (d *Object) Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error) {
	logger := log.From(ctx)

	object, err := options.LoadObject(ctx, namespace, options.Fields, d.objectStoreKey)
	if err != nil {
		return EmptyContentResponse, errors.Wrapf(err, "loading object with %s", d.objectStoreKey.String())
	}

	item := d.objectType()

	if err := scheme.Scheme.Convert(object, item, nil); err != nil {
		return EmptyContentResponse, errors.Wrapf(err, "converting dynamic object to a type")
	}

	if err := copyObjectMeta(item, object); err != nil {
		return EmptyContentResponse, errors.Wrap(err, "copying object metadata")
	}

	accessor := meta.NewAccessor()
	objectName, _ := accessor.Name(object)

	title := append([]component.TitleComponent{}, component.NewText(d.baseTitle))
	if objectName != "" {
		title = append(title, component.NewText(objectName))
	}

	cr := component.NewContentResponse(title)

	currentObject, ok := item.(runtime.Object)
	if !ok {
		return EmptyContentResponse, errors.Errorf("expected item to be a runtime object. It was a %T",
			item)
	}

	hasTabError := false
	for _, tfd := range d.tabFuncDescriptors {
		if err := tfd.tabFunc(ctx, currentObject, cr, options); err != nil {
			hasTabError = true
			logger.With(
				"err", err,
				"tab-name", tfd.name,
			).Errorf("generating object Describer tab")
		}
	}

	if hasTabError {
		logger.With("tab-object", object).Errorf("unable to generate all tabs for object")
	}

	tabs, err := options.PluginManager().Tabs(object)
	if err != nil {
		return EmptyContentResponse, errors.Wrap(err, "getting tabs from plugins")
	}

	for _, tab := range tabs {
		tab.Contents.SetAccessor(tab.Name)
		cr.Add(&tab.Contents)
	}

	return *cr, nil
}

func (d *Object) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}

func (d *Object) addSummaryTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	vc, err := options.Printer.Print(ctx, object, options.PluginManager())
	if err != nil {
		return errors.Wrap(err, "printing object")
	}

	if vc == nil {
		return errors.Wrap(err, "unable to print a nil object")
	}

	vc.SetAccessor("summary")
	cr.Add(vc)

	return nil
}

func (d *Object) addResourceViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	if !d.disableResourceViewer {

		resourceViewComponent, err := options.ComponentCache.Get(ctx, object)
		if err != nil {
			return err
		}

		resourceViewComponent.SetAccessor("resourceViewer")
		cr.Add(resourceViewComponent)
	}

	return nil
}

func (d *Object) addYAMLViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	yvComponent, err := yamlviewer.ToComponent(object)
	if err != nil {
		return errors.Wrap(err, "converting object to YAML")
	}
	yvComponent.SetAccessor("yaml")
	cr.Add(yvComponent)

	return nil
}

func (d *Object) addLogsTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	if isPod(object) {
		logsComponent, err := logviewer.ToComponent(object)
		if err != nil {
			return errors.Wrap(err, "retrieving logs for pod")
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

// Section is a wrapper to combine content from multiple describers.
type Section struct {
	path       string
	title      string
	describers []Describer
}

// NewSectionDescriber creates a Section.
func NewSectionDescriber(p, title string, describers ...Describer) *Section {
	return &Section{
		path:       p,
		title:      title,
		describers: describers,
	}
}

// Describe generates content.
func (d *Section) Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error) {
	list := component.NewList(d.title, nil)

	for _, child := range d.describers {
		cResponse, err := child.Describe(ctx, prefix, namespace, options)
		if err != nil {
			return EmptyContentResponse, err
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

func (d *Section) PathFilters() []PathFilter {
	PathFilters := []PathFilter{
		*NewPathFilter(d.path, d),
	}

	for _, child := range d.describers {
		PathFilters = append(PathFilters, child.PathFilters()...)
	}

	return PathFilters
}

func isPod(object runtime.Object) bool {
	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return apiVersion == "v1" && kind == "Pod"
}
