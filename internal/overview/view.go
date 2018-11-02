package overview

import (
	"context"
	"time"

	"github.com/heptio/developer-dash/internal/content"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func tableCols(names ...string) []content.TableColumn {
	columns := []content.TableColumn{}
	for _, name := range names {
		columns = append(columns, content.TableColumn{Name: name, Accessor: name})
	}

	return columns
}

func formatTime(t *metav1.Time) string {
	if t == nil {
		return "-"
	}

	return t.UTC().Format(time.RFC1123Z)
}

// LookupFunc is a function for looking up the contents of a cell.
type LookupFunc func(namespace, prefix string, cell interface{}) content.Text
