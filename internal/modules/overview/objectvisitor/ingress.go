package objectvisitor

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/internal/util/kubernetes"
)

// Ingress is a typed visitor for ingress objects.
type Ingress struct {
	queryer queryer.Queryer
}

var _ TypedVisitor = (*Ingress)(nil)

// NewIngress creates Ingress.
func NewIngress(q queryer.Queryer) *Ingress {
	return &Ingress{queryer: q}
}

// Supports returns the gvk this typed visitor supports.
func (i *Ingress) Supports() schema.GroupVersionKind {
	return gvk.Ingress
}

// Visit visits an ingress. It looks for associated ingresses.
func (i *Ingress) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool) error {
	ctx, span := trace.StartSpan(ctx, "visitIngress")
	defer span.End()

	ingress := &extv1beta1.Ingress{}
	if err := convertToType(object, ingress); err != nil {
		return err
	}

	services, err := i.queryer.ServicesForIngress(ctx, ingress)
	if err != nil {
		return err
	}

	var g errgroup.Group

	for i := range services {
		service := services[i]
		g.Go(func() error {
			m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(service)
			if err != nil {
				return err
			}
			u := &unstructured.Unstructured{Object: m}

			if err := visitor.Visit(ctx, u, handler, true); err != nil {
				return errors.Wrapf(err, "ingress %s visit service %s",
					kubernetes.PrintObject(ingress), kubernetes.PrintObject(u))
			}
			return handler.AddEdge(ctx, object, u)
		})

	}

	return g.Wait()
}
