package overview

import (
	"context"
	"path"
	"reflect"

	cacheutil "github.com/heptio/developer-dash/internal/cache/util"
	"github.com/heptio/developer-dash/pkg/view/component"

	"github.com/heptio/developer-dash/internal/cluster"
)

type ResourceTitle struct {
	List   string
	Object string
}

type ResourceOptions struct {
	Path                  string
	CacheKey              cacheutil.Key
	ListType              interface{}
	ObjectType            interface{}
	Titles                ResourceTitle
	DisableResourceViewer bool
	ClusterWide           bool
}

type Resource struct {
	ResourceOptions
}

func NewResource(options ResourceOptions) *Resource {
	return &Resource{
		ResourceOptions: options,
	}
}

func (r *Resource) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	return r.List().Describe(ctx, prefix, namespace, clusterClient, options)
}

func (r *Resource) List() *ListDescriber {
	return NewListDescriber(
		r.Path,
		r.Titles.List,
		r.CacheKey,
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ListType).Elem().Type()).Interface()
		},
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
		},
		r.ClusterWide,
	)
}

func (r *Resource) Object() *ObjectDescriber {
	return NewObjectDescriber(
		path.Join(r.Path, resourceNameRegex),
		r.Titles.Object,
		DefaultLoader(r.CacheKey),
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
		},
		r.DisableResourceViewer,
	)
}

func (r *Resource) PathFilters() []pathFilter {
	filters := []pathFilter{
		*newPathFilter(r.Path, r.List()),
		*newPathFilter(path.Join(r.Path, resourceNameRegex), r.Object()),
	}

	return filters
}
