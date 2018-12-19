package overview

import (
	"context"
	"fmt"
	"reflect"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/heptio/developer-dash/internal/printers"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
	kprinters "k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
)

type LoaderFunc func(ctx context.Context, c Cache, namespace string, fields map[string]string) ([]*unstructured.Unstructured, error)

var DefaultLoader = func(cacheKey CacheKey) LoaderFunc {
	return func(ctx context.Context, c Cache, namespace string, fields map[string]string) ([]*unstructured.Unstructured, error) {
		cacheKeys := []CacheKey{cacheKey}
		return loadObjects(ctx, c, namespace, fields, cacheKeys)
	}
}

type ObjectTransformFunc func(namespace, prefix string, contents *[]content.Content) func(*metav1beta1.Table) error

type DescriberOptions struct {
	Cache  Cache
	Fields map[string]string
}

// Describer creates content.
type Describer interface {
	Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (ContentResponse, error)
	PathFilters(namespace string) []pathFilter
}

type baseDescriber struct{}

func newBaseDescriber() *baseDescriber {
	return &baseDescriber{}
}

func (d *baseDescriber) clock() clock.Clock {
	return &clock.RealClock{}
}

type ListDescriber struct {
	*baseDescriber

	path                string
	title               string
	listType            func() interface{}
	objectType          func() interface{}
	cacheKey            CacheKey
	objectTransformFunc ObjectTransformFunc
}

func NewListDescriber(p, title string, cacheKey CacheKey, listType, objectType func() interface{}, otf ObjectTransformFunc) *ListDescriber {
	return &ListDescriber{
		path:                p,
		title:               title,
		baseDescriber:       newBaseDescriber(),
		cacheKey:            cacheKey,
		listType:            listType,
		objectType:          objectType,
		objectTransformFunc: otf,
	}
}

// Describe creates content.
func (d *ListDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (ContentResponse, error) {
	var contents []content.Content

	objects, err := loadObjects(ctx, options.Cache, namespace, options.Fields, []CacheKey{d.cacheKey})
	if err != nil {
		return emptyContentResponse, err
	}

	list := d.listType()

	v := reflect.ValueOf(list)
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

	listObject, ok := list.(runtime.Object)
	if !ok {
		return emptyContentResponse, errors.Errorf("expected list to be a runtime object. It was a %T",
			list)
	}

	otf := d.objectTransformFunc(namespace, prefix, &contents)
	if err := printObject(listObject, otf); err != nil {
		return emptyContentResponse, err
	}

	return ContentResponse{
		Views: []Content{
			{
				Contents: contents,
				Title:    "",
			},
		},
	}, nil
}

func (d *ListDescriber) PathFilters(namespace string) []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

type ObjectDescriber struct {
	*baseDescriber

	path       string
	baseTitle  string
	objectType func() interface{}
	loaderFunc LoaderFunc
	sections   []ContentSection
}

func NewObjectDescriber(p, baseTitle string, loaderFunc LoaderFunc, objectType func() interface{}, sections []ContentSection) *ObjectDescriber {
	return &ObjectDescriber{
		path:          p,
		baseTitle:     baseTitle,
		baseDescriber: newBaseDescriber(),
		loaderFunc:    loaderFunc,
		objectType:    objectType,
		sections:      sections,
	}
}

func (d *ObjectDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (ContentResponse, error) {
	objects, err := d.loaderFunc(ctx, options.Cache, namespace, options.Fields)
	if err != nil {
		return emptyContentResponse, err
	}

	if len(objects) != 1 {
		return emptyContentResponse, errors.Errorf("expected exactly one object")
	}

	object := objects[0]

	item := d.objectType()

	if err := scheme.Scheme.Convert(object, item, nil); err != nil {
		return emptyContentResponse, err
	}

	if err := copyObjectMeta(item, object); err != nil {
		return emptyContentResponse, errors.Wrapf(err, "copying object metadata")
	}

	objectName := object.GetName()

	var title string

	if objectName == "" {
		title = d.baseTitle
	} else {
		title = fmt.Sprintf("%s: %s", d.baseTitle, objectName)
	}

	newObject, ok := item.(runtime.Object)
	if !ok {
		return emptyContentResponse, errors.Errorf("expected item to be a runtime object. It was a %T",
			item)
	}

	cr := ContentResponse{
		Title: title,
	}

	cl := &clock.RealClock{}

	for _, section := range d.sections {
		var contents []content.Content
		for _, viewFactory := range section.Views {
			view := viewFactory(prefix, namespace, cl)
			viewContent, err := view.Content(ctx, newObject, options.Cache)
			if err != nil {
				return emptyContentResponse, err
			}

			contents = append(contents, viewContent...)
		}

		cr.Views = append(cr.Views, Content{
			Contents: contents,
			Title:    section.Title,
		})
	}

	return cr, nil
}

func (d *ObjectDescriber) PathFilters(namespace string) []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
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

func convertToInternal(in runtime.Object) (runtime.Object, error) {
	return scheme.Scheme.ConvertToVersion(in, runtime.InternalGroupVersioner)
}

func printObject(object runtime.Object, transformFunc func(*metav1beta1.Table) error) error {
	options := kprinters.PrintOptions{
		Wide:       true,
		ShowLabels: true,
		WithKind:   true,
	}

	decoder := scheme.Codecs.UniversalDecoder()
	p := printers.NewHumanReadablePrinter(decoder, options)

	printersinternal.AddHandlers(p)

	internal, err := convertToInternal(object)
	if err != nil {
		return errors.Wrapf(err, "converting to internal: %T", object)
	}

	tbl, err := p.PrintTable(internal, options)
	if err != nil {
		return err
	}

	if transformFunc != nil {
		return transformFunc(tbl)
	}

	return nil
}

func printContentTable(title, namespace, prefix, emptyMessage string, tbl *metav1beta1.Table, m map[string]lookupFunc) (*content.Table, error) {
	contentTable := content.NewTable(title, emptyMessage)

	headers := make(map[int]string)

	for i, column := range tbl.ColumnDefinitions {

		headers[i] = column.Name

		contentTable.Columns = append(contentTable.Columns, content.TableColumn{
			Name:     column.Name,
			Accessor: column.Name,
		})
	}

	transforms := buildTransforms(m)

	for _, row := range tbl.Rows {
		contentRow := content.TableRow{}

		for pos, header := range headers {
			cell := row.Cells[pos]

			c, ok := transforms[header]
			if !ok {
				contentRow[header] = content.NewStringText(fmt.Sprintf("%v", cell))
			} else {
				contentRow[header] = c(namespace, prefix, cell)
			}
		}

		contentTable.AddRow(contentRow)
	}

	return &contentTable, nil
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
func (d *SectionDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (ContentResponse, error) {
	var contents []content.Content

	for _, child := range d.describers {
		cResponse, err := child.Describe(ctx, prefix, namespace, clusterClient, options)
		if err != nil {
			return emptyContentResponse, err
		}

		for _, views := range cResponse.Views {
			for _, childContent := range views.Contents {
				if !childContent.IsEmpty() {
					contents = append(contents, childContent)
				}
			}
		}
	}

	if len(contents) == 0 {
		tbl := content.NewTable(d.title, fmt.Sprintf("Namespace %s does not have any resources of this type", namespace))
		contents = append(contents, &tbl)
	}

	cr := ContentResponse{
		Views: []Content{
			{Contents: contents, Title: d.title},
		},
	}

	return cr, nil
}

func (d *SectionDescriber) PathFilters(namespace string) []pathFilter {
	pathFilters := []pathFilter{
		*newPathFilter(d.path, d),
	}

	for _, child := range d.describers {
		pathFilters = append(pathFilters, child.PathFilters(namespace)...)
	}

	return pathFilters
}
