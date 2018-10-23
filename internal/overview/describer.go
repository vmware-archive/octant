package overview

import (
	"fmt"
	"reflect"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/heptio/developer-dash/internal/printers"
	"github.com/heptio/developer-dash/internal/view"
	"github.com/pkg/errors"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/kubernetes/scheme"
	kprinters "k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
)

type ObjectTransformFunc func(namespace, prefix string, contents *[]content.Content) func(*metav1beta1.Table) error

type DescriberOptions struct {
	Cache  Cache
	Fields map[string]string
}

// Describer creates content.
type Describer interface {
	Describe(prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) ([]content.Content, string, error)
	PathFilters() []pathFilter
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
func (d *ListDescriber) Describe(prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) ([]content.Content, string, error) {
	var contents []content.Content

	objects, err := loadObjects(options.Cache, namespace, options.Fields, []CacheKey{d.cacheKey})
	if err != nil {
		return nil, "", err
	}

	list := d.listType()

	v := reflect.ValueOf(list)
	f := reflect.Indirect(v).FieldByName("Items")

	for _, object := range objects {
		item := d.objectType()
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, item)
		if err != nil {
			return nil, "", err
		}

		setItemName(item, object.GetName())

		newSlice := reflect.Append(f, reflect.ValueOf(item).Elem())
		f.Set(newSlice)
	}

	listObject, ok := list.(runtime.Object)
	if !ok {
		return nil, "", errors.Errorf("expected list to be a runtime object. It was a %T",
			list)
	}

	otf := d.objectTransformFunc(namespace, prefix, &contents)
	if err := printObject(listObject, otf); err != nil {
		return nil, "", err
	}

	return contents, d.title, nil
}

func (d *ListDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

type ObjectDescriber struct {
	*baseDescriber

	path                string
	baseTitle           string
	objectType          func() interface{}
	cacheKey            CacheKey
	objectTransformFunc ObjectTransformFunc
	views               []view.View
}

func NewObjectDescriber(p, baseTitle string, cacheKey CacheKey, objectType func() interface{}, otf ObjectTransformFunc, views []view.View) *ObjectDescriber {
	return &ObjectDescriber{
		path:                p,
		baseTitle:           baseTitle,
		baseDescriber:       newBaseDescriber(),
		cacheKey:            cacheKey,
		objectType:          objectType,
		objectTransformFunc: otf,
		views:               views,
	}
}

func (d *ObjectDescriber) Describe(prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) ([]content.Content, string, error) {
	objects, err := loadObjects(options.Cache, namespace, options.Fields, []CacheKey{d.cacheKey})
	if err != nil {
		return nil, "", err
	}

	var contents []content.Content

	if len(objects) != 1 {
		return nil, "", errors.Errorf("expected exactly one object")
	}

	object := objects[0]

	item := d.objectType()
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, item)
	if err != nil {
		return nil, "", err
	}

	objectName := object.GetName()
	setItemName(item, objectName)

	var title string

	if objectName == "" {
		title = d.baseTitle
	} else {
		title = fmt.Sprintf("%s: %s", d.baseTitle, objectName)
	}

	newObject, ok := item.(runtime.Object)
	if !ok {
		return nil, "", errors.Errorf("expected item to be a runtime object. It was a %T",
			item)
	}

	otf := d.objectTransformFunc(namespace, prefix, &contents)
	if err := printObject(newObject, otf); err != nil {
		return nil, "", err
	}

	// TODO should show parents here
	// TODO will need to register a map of object transformers?

	for _, v := range d.views {
		viewContent, err := v.Content(nil, newObject, nil)
		if err != nil {
			return nil, "", err
		}

		contents = append(contents, viewContent...)
	}

	eventsTable, err := eventsForObject(object, options.Cache, prefix, namespace, d.clock())
	if err != nil {
		return nil, "", err
	}

	contents = append(contents, eventsTable)

	return contents, title, nil
}

func (d *ObjectDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

func setItemName(item interface{}, name string) {
	setNameVal := reflect.ValueOf(item).MethodByName("SetName")
	setNameIface := setNameVal.Interface()
	setName := setNameIface.(func(string))
	setName(name)
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

	tbl, err := p.PrintTable(object, options)
	if err != nil {
		return err
	}

	if transformFunc != nil {
		return transformFunc(tbl)
	}

	return nil
}

func printContentTable(title, namespace, prefix string, tbl *metav1beta1.Table, m map[string]lookupFunc) (*content.Table, error) {
	contentTable := content.NewTable(title)

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
func (d *SectionDescriber) Describe(prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) ([]content.Content, string, error) {
	var contents []content.Content

	for _, child := range d.describers {
		childContents, _, err := child.Describe(prefix, namespace, clusterClient, options)
		if err != nil {
			return nil, "", err
		}

		for _, childContent := range childContents {
			if !childContent.IsEmpty() {
				contents = append(contents, childContent)
			}
		}
	}

	return contents, d.title, nil
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
