package console

import (
	"github.com/vmware-tanzu/octant/pkg/log"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	_ "github.com/dop251/goja_nodejs/util"
)

type Console struct {
	runtime *goja.Runtime
	util    *goja.Object
	logger  log.Logger
}

func (c *Console) log(fn func(string, ...interface{})) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if format, ok := goja.AssertFunction(c.util.Get("format")); ok {
			ret, err := format(c.util, call.Arguments...)
			if err != nil {
				panic(err)
			}

			fn(ret.String())
		} else {
			panic(c.runtime.NewTypeError("util.format is not a function"))
		}

		return nil
	}
}

func RequireFactory(logger log.Logger) require.ModuleLoader {
	return func(runtime *goja.Runtime, module *goja.Object) {
		c := &Console{
			runtime: runtime,
			logger:  logger,
		}

		c.util = require.Require(runtime, "util").(*goja.Object)

		o := module.Get("exports").(*goja.Object)
		o.Set("log", c.log(c.logger.Infof))
		o.Set("error", c.log(c.logger.Errorf))
		o.Set("warn", c.log(c.logger.Warnf))
	}

}

func Enable(runtime *goja.Runtime) {
	runtime.Set("console", require.Require(runtime, "console"))
}

func CustomInit(logger log.Logger) {
	require.RegisterNativeModule("console", RequireFactory(logger))
}
