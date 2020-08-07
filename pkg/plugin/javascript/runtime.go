package javascript

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/dop251/goja_nodejs/require"

	olog "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/log"
)

func CreateRuntimeLoop(ctx context.Context, logName string) (*eventloop.EventLoop, error) {
	loop := eventloop.NewEventLoop()
	loop.Start()

	errCh := make(chan error)

	loop.RunOnLoop(func(vm *goja.Runtime) {
		vm.Set("global", vm.GlobalObject())
		vm.Set("self", vm.GlobalObject())

		_, err := vm.RunString(`
var module = { exports: {} };
var exports = module.exports;
`)
		if err != nil {
			errCh <- fmt.Errorf("runtime global values: %w", err)
			return
		}

		registry := new(require.Registry)
		registry.Enable(vm)

		logger := olog.From(ctx).With("plugin", logName)
		printer := logPrinter{logger: logger}
		registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))
		console.Enable(vm)

		// This maps Caps fields to lower fields from struct to Object based on the JSON annotations.
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
		errCh <- nil
	})

	err := <-errCh
	if err != nil {
		return nil, err
	}

	return loop, nil
}

func ExtractDefaultClass(vm *goja.Runtime) (*goja.Object, error) {
	// This is the location of a export default class that implements the Octant
	// TypeScript module definition.
	instantiateClass := "var _concretePlugin = new module.exports.default; _concretePlugin"
	// This is the library name the Octant webpack configuration uses.
	if vm.Get("_octantPlugin") != nil {
		instantiateClass = "var _concretePlugin = new _octantPlugin(dashboardClient, httpClient); _concretePlugin"
	}

	v, err := vm.RunString(instantiateClass)
	if err != nil {
		return nil, fmt.Errorf("unable to create new plugin: %w", err)
	}
	pluginClass := v.ToObject(vm)
	return pluginClass, nil
}

type logPrinter struct {
	logger log.Logger
}

func (l logPrinter) Log(msg string) {
	l.logger.Infof(msg)
}

func (l logPrinter) Warn(msg string) {
	l.logger.Warnf(msg)
}

func (l logPrinter) Error(msg string) {
	l.logger.Errorf(msg)
}
