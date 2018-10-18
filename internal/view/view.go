package view

import (
	"context"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
)

type View interface {
	Content(ctx context.Context, object runtime.Object, clusterClient cluster.ClientInterface) ([]content.Content, error)
}

func tableCol(name string) content.TableColumn {
	return content.TableColumn{
		Name:     name,
		Accessor: name,
	}
}
