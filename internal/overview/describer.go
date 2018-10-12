package overview

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/kubernetes/scheme"
	printers "k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
)

type ObjecTransformFunc func(namespace, prefix string, contents *[]Content) func(*metav1beta1.Table) error

// Describer creates content.
type Describer interface {
	Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error)
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
	listType            func() interface{}
	objectType          func() interface{}
	cacheKey            CacheKey
	objectTransformFunc ObjecTransformFunc
}

func NewListDescriber(p string, cacheKey CacheKey, listType, objectType func() interface{}, otf ObjecTransformFunc) *ListDescriber {
	return &ListDescriber{
		path:                p,
		baseDescriber:       newBaseDescriber(),
		cacheKey:            cacheKey,
		listType:            listType,
		objectType:          objectType,
		objectTransformFunc: otf,
	}
}

// Describe creates content.
func (d *ListDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	var contents []Content

	objects, err := loadObjects(cache, namespace, fields, []CacheKey{d.cacheKey})
	if err != nil {
		return nil, err
	}

	list := d.listType()

	v := reflect.ValueOf(list)
	f := reflect.Indirect(v).FieldByName("Items")

	for _, object := range objects {
		item := d.objectType()
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, item)
		if err != nil {
			return nil, err
		}

		setItemName(item, object.GetName())

		newSlice := reflect.Append(f, reflect.ValueOf(item).Elem())
		f.Set(newSlice)
	}

	listObject, ok := list.(runtime.Object)
	if !ok {
		return nil, errors.Errorf("expected list to be a runtime object. It was a %T",
			list)
	}

	otf := d.objectTransformFunc(namespace, prefix, &contents)
	if err := printObject(listObject, otf); err != nil {
		return nil, err
	}

	return contents, nil
}

func (d *ListDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

type ObjectDescriber struct {
	*baseDescriber

	path                string
	objectType          func() interface{}
	cacheKey            CacheKey
	objectTransformFunc ObjecTransformFunc
}

func NewObjectDescriber(p string, cacheKey CacheKey, objectType func() interface{}, otf ObjecTransformFunc) *ObjectDescriber {
	return &ObjectDescriber{
		path:                p,
		baseDescriber:       newBaseDescriber(),
		cacheKey:            cacheKey,
		objectType:          objectType,
		objectTransformFunc: otf,
	}
}

func (d *ObjectDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	objects, err := loadObjects(cache, namespace, fields, []CacheKey{d.cacheKey})
	if err != nil {
		return nil, err
	}

	var contents []Content

	if len(objects) != 1 {
		return nil, errors.Errorf("expected exactly one object")
	}

	object := objects[0]

	item := d.objectType()
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, item)
	if err != nil {
		return nil, err
	}

	setItemName(item, object.GetName())

	newObject, ok := item.(runtime.Object)
	if !ok {
		return nil, errors.Errorf("expected item to be a runtime object. It was a %T",
			item)
	}

	otf := d.objectTransformFunc(prefix, namespace, &contents)
	if err := printObject(newObject, otf); err != nil {
		return nil, err
	}

	eventsTable, err := eventsForObject(object, cache, prefix, namespace, d.clock())
	if err != nil {
		return nil, err
	}

	contents = append(contents, eventsTable)

	return contents, nil
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
	options := printers.PrintOptions{
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

func printContentTable(title, namespace, prefix string, tbl *metav1beta1.Table, m map[string]lookupFunc) (*table, error) {
	contentTable := newTable(title)

	headers := make(map[int]string)

	for i, column := range tbl.ColumnDefinitions {

		headers[i] = column.Name

		contentTable.Columns = append(contentTable.Columns, tableColumn{
			Name:     column.Name,
			Accessor: column.Name,
		})
	}

	transforms := buildTransforms(m)

	for _, row := range tbl.Rows {
		contentRow := tableRow{}

		for pos, header := range headers {
			cell := row.Cells[pos]

			c, ok := transforms[header]
			if !ok {
				contentRow[header] = newStringText(fmt.Sprintf("%v", cell))
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
	describers []Describer
}

// NewSectionDescriber creates a SectionDescriber.
func NewSectionDescriber(p string, describers ...Describer) *SectionDescriber {
	return &SectionDescriber{
		path:       p,
		describers: describers,
	}
}

// Describe generates content.
func (d *SectionDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	var contents []Content

	for _, child := range d.describers {
		childContent, err := child.Describe(prefix, namespace, cache, fields)
		if err != nil {
			return nil, err
		}

		contents = append(contents, childContent...)
	}

	return contents, nil
}

func (d *SectionDescriber) PathFilters() []pathFilter {
	var pathFilters []pathFilter

	for _, child := range d.describers {
		pathFilters = append(pathFilters, child.PathFilters()...)
	}

	return pathFilters
}
