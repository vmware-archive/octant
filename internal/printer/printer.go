/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/pkg/config"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

//go:generate mockgen -destination=./fake/mock_printer.go -package=fake github.com/vmware-tanzu/octant/internal/printer Printer

// Options provides options to a print handler
type Options struct {
	DisableLabels bool
	DashConfig    config.Dash
	Link          link.Interface
	ObjectFactory ObjectFactory
}

// Printer is an interface for printing runtime objects.
type Printer interface {
	// Print prints a runtime object.
	Print(ctx context.Context, object runtime.Object) (component.Component, error)
}

// Resource prints runtime objects.
type Resource struct {
	handlerMap map[reflect.Type]reflect.Value
	dashConfig config.Dash
}

var _ Printer = (*Resource)(nil)

// NewResource creates an instance of ResourcePrinter.
func NewResource(dashConfig config.Dash) *Resource {
	return &Resource{
		handlerMap: make(map[reflect.Type]reflect.Value),
		dashConfig: dashConfig,
	}
}

// Print prints a runtime object. If not handler can be found for the type,
// it will print using `DefaultPrintFunc`.
func (p *Resource) Print(ctx context.Context, object runtime.Object) (component.Component, error) {
	l, err := link.NewFromDashConfig(p.dashConfig)
	if err != nil {
		return nil, err
	}

	printOptions := Options{
		DashConfig:    p.dashConfig,
		Link:          l,
		ObjectFactory: NewDefaultObjectFactory(),
	}

	t := reflect.TypeOf(object)
	printFunc, ok := p.handlerMap[t]
	if ok {
		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(object),
			reflect.ValueOf(printOptions)}
		results := printFunc.Call(args)
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}

		viewComponent := results[0].Interface().(component.Component)
		return viewComponent, nil
	}

	return DefaultPrintFunc(ctx, object, printOptions)
}

// Handler adds a printer handler.
// See ValidatePrintHandlerFunc for required method signature.
func (p *Resource) Handler(printFunc interface{}) error {
	printFuncValue := reflect.ValueOf(printFunc)
	if err := ValidatePrintHandlerFunc(printFuncValue); err != nil {
		return err
	}

	objType := printFuncValue.Type().In(1)
	if _, ok := p.handlerMap[objType]; ok {
		return errors.Errorf("registered duplicate printer for %v", objType)
	}

	p.handlerMap[objType] = printFuncValue

	return nil
}

// ValidatePrintHandlerFunc validates print handler signature.
// printFunc is the function that will be called to print an object.
// printFunc must be of the following type:
//   func printFunc(ctx context.Context, object ObjectType, options Options) (component.Component, error)
// where:
//   ObjectType is the type of object that will be printed
func ValidatePrintHandlerFunc(printFunc reflect.Value) error {
	if printFunc.Kind() != reflect.Func {
		return errors.Errorf("invalid print handler. %#v is not a function", printFunc)
	}

	funcType := printFunc.Type()
	if numIn, numOut := funcType.NumIn(), funcType.NumOut(); numIn != 3 || numOut != 2 {
		return errors.Errorf("invalid printer handler. "+
			"Must accept 3 parameters and 2 return values. "+
			"It accepted %d parameters and returned %d values",
			numIn, numOut,
		)
	}

	if funcType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() ||
		funcType.In(2) != reflect.TypeOf((*Options)(nil)).Elem() ||
		funcType.Out(0) != reflect.TypeOf((*component.Component)(nil)).Elem() ||
		funcType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.Errorf("invalid print handler. The expected signature is: "+
			"func handler(ctx context.Context, obj %v, options *PrintOptions) (component.Component, error)",
			funcType.In(0))
	}

	return nil
}

// DefaultPrintFunc is a default object printer. It prints Kubernetes resource
// lists with three columns: name, labels, age. Returns nil if the object
// should not be printed.
func DefaultPrintFunc(_ context.Context, object runtime.Object, _ Options) (component.Component, error) {
	if object == nil {
		return nil, errors.New("unable to print nil objects")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, err
	}

	if _, ok := m["items"]; !ok {
		// It's not a list, so return empty content
		return nil, nil
	}

	if m["items"] == nil {
		// List is empty, so return empty content
		return nil, nil
	}

	cols := component.NewTableCols("Name", "Labels", "Age")

	title := strings.TrimPrefix(fmt.Sprintf("%T", object), "*")
	desc := strings.Split(title, ".")
	gvk := schema.FromAPIVersionAndKind(desc[0], desc[1])
	title = gvk.String()

	table := component.NewTable(title, "We couldn't find any objects!", cols)

	items := m["items"].([]interface{})

	for _, item := range items {
		r, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("item was not a resource")
		}

		u := unstructured.Unstructured{Object: r}

		name := component.NewText(u.GetName())
		labels := component.NewLabels(u.GetLabels())
		age := component.NewTimestamp(u.GetCreationTimestamp().Time)

		row := component.TableRow{
			"Name":   name,
			"Labels": labels,
			"Age":    age,
		}

		table.Add(row)
	}

	return table, nil
}
