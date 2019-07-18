package kubernetes

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

// PrintObject prints an object's name and group/version/kind
func PrintObject(object runtime.Object) (s string) {

	if object == nil {
		return "<nil kubernetes object>"
	}

	accessor, err := meta.Accessor(object)
	if err != nil {
		return fmt.Sprintf("<invalid object: %s", err)
	}

	defer func() {
		if r := recover(); r != nil {
			s = fmt.Sprintf("print object paniced: %v; print accessor: %s", r, spew.Sdump(accessor, object))
		}

	}()

	return fmt.Sprintf("<%s %s>", object.GetObjectKind().GroupVersionKind(), accessor.GetName())
}
