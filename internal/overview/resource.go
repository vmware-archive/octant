package overview

import (
	"context"
	"fmt"
	"path"
	"reflect"

	cacheutil "github.com/heptio/developer-dash/internal/cache/util"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/content"
)

func resourceLink(sectionType, resourceType string) lookupFunc {
	return func(namespace, prefix string, cell interface{}) content.Text {
		name := fmt.Sprintf("%v", cell)
		resourcePath := path.Join("/content", "overview", sectionType, resourceType, name)
		return content.NewLinkText(name, resourcePath)
	}
}

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

var defaultTransforms = map[string]lookupFunc{
	"Labels": func(namespace, prefix string, cell interface{}) content.Text {
		text := fmt.Sprintf("%v", cell)
		return content.NewStringText(text)
	},
}

func buildTransforms(transforms map[string]lookupFunc) map[string]lookupFunc {
	m := make(map[string]lookupFunc)
	for k, v := range defaultTransforms {
		m[k] = v
	}
	for k, v := range transforms {
		m[k] = v
	}

	return m
}
