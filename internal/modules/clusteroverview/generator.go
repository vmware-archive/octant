package clusteroverview

import (
	"context"

	gen "github.com/heptio/developer-dash/internal/generator"
	"github.com/heptio/developer-dash/pkg/view/component"
)

type generator struct{}

var _ gen.Generator = (*generator)(nil)

func (generator) Generate(ctx context.Context, path, prefix, namespace string, opts gen.Options) (component.ContentResponse, error) {
	panic("implement me")
}

