package javascript

import (
	"context"
	"fmt"

	"github.com/dop251/goja"

	"github.com/vmware-tanzu/octant/internal/octant"
)

type DashboardCreateOrUpdateFromYAML struct {
	storage octant.Storage
}

var _ octant.DashboardClientFunction = &DashboardCreateOrUpdateFromYAML{}

// NewDashboardGet creates an instance of DashboardGet.
func NewDashboardCreateOrUpdateFromYAML(storage octant.Storage) *DashboardCreateOrUpdateFromYAML {
	d := &DashboardCreateOrUpdateFromYAML{
		storage: storage,
	}
	return d
}

// Name returns the name of this function. It will always return "Get".
func (d *DashboardCreateOrUpdateFromYAML) Name() string {
	return "CreateOrUpdateFromYAML"
}

// Call creates a function call that gets an object by key. If the key is invalid, or if the
// get is unsuccessful, it will throw a javascript exception.
func (d *DashboardCreateOrUpdateFromYAML) Call(ctx context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value {
	return func(c goja.FunctionCall) goja.Value {
		fmt.Println("HERE in call")
		var namespace, yaml string

		obj0 := c.Argument(0).ToString()
		obj1 := c.Argument(1).ToString()

		// This will never error since &key is a pointer to a type.
		_ = vm.ExportTo(obj0, &namespace)
		_ = vm.ExportTo(obj1, &yaml)

		_, err := d.storage.ObjectStore().CreateOrUpdateFromYAML(ctx, namespace, yaml)
		if err != nil {
			panic(panicMessage(vm, err, ""))
		}

		return goja.Undefined()
	}
}
