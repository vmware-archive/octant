package objectvisitor

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/internal/util/kubernetes"
)

// Object is the default visitor for an object.
type Object struct {
	queryer    queryer.Queryer
	dashConfig config.Dash
}

// NewObject creates Object.
func NewObject(dashConfig config.Dash, q queryer.Queryer) *Object {
	return &Object{
		dashConfig: dashConfig,
		queryer:    q,
	}
}

// Visit visits an objects. It looks at immediate ancestors and descendants.
func (o *Object) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool) error {
	if object == nil {
		return errors.New("can't visit nil object")
	}

	ctx, span := trace.StartSpan(ctx, "handleObject")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("apiVersion", object.GetAPIVersion()),
		trace.StringAttribute("kind", object.GetKind()),
		trace.StringAttribute("name", object.GetName()),
		trace.StringAttribute("namespace", object.GetNamespace()),
	}, "handling object")

	var g errgroup.Group

	object = object.DeepCopy()

	g.Go(func() error {
		found, owner, err := o.queryer.OwnerReference(ctx, object)
		if err != nil {
			return errors.Wrapf(err, "unable to check owner reference for %s", kubernetes.PrintObject(object))
		}

		if found {
			if owner == nil {
				return errors.Errorf("unable to find owner for %s", object)
			}

			if err := visitor.Visit(ctx, owner, handler, false); err != nil {
				return errors.Wrapf(err, "visit ancestor %s for %s",
					kubernetes.PrintObject(owner),
					kubernetes.PrintObject(object))
			}
			return handler.AddEdge(ctx, object, owner)
		}

		return nil
	})

	if visitDescendants {
		children, err := o.queryer.Children(ctx, object)
		if err != nil {
			return err
		}

		for i := range children.Items {
			child := &children.Items[i]
			g.Go(func() error {
				if err := visitor.Visit(ctx, child, handler, true); err != nil {
					return errors.Wrapf(err, "visit child %s for %s",
						kubernetes.PrintObject(child),
						kubernetes.PrintObject(object))
				}

				return handler.AddEdge(ctx, object, child)
			})
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return handler.Process(ctx, object)
}
