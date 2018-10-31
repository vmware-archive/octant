package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
)

// View is a view that can be embbeded in the resource overview.
type View interface {
	Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error)
}

func tableCol(name string) content.TableColumn {
	return content.TableColumn{
		Name:     name,
		Accessor: name,
	}
}

// LookupFunc is a function for looking up the contents of a cell.
type LookupFunc func(namespace, prefix string, cell interface{}) content.Text
