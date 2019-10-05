package objectvisitor

import (
	"context"

	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/pkg/errors"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/internal/util/kubernetes"
)

// HorizontalPodAutoscaler is a typed visitor for services.
type HorizontalPodAutoscaler struct {
	queryer queryer.Queryer
}

var _ TypedVisitor = (*HorizontalPodAutoscaler)(nil)

// NewHorizontalPodAutoscaler creates an instance of HorizontalPodAutoscaler
func NewHorizontalPodAutoscaler(q queryer.Queryer) *HorizontalPodAutoscaler {
	return &HorizontalPodAutoscaler{queryer: q}
}

// Supports returns the gvk this typed visitor supports.
func (HorizontalPodAutoscaler) Supports() schema.GroupVersionKind {
	return gvk.HorizontalPodAutoscaler
}

// Visit visits a hpa. It looks for an associated scale target (replication controllers, deployments, and replica sets)
func (s *HorizontalPodAutoscaler) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool) error {
	ctx, span := trace.StartSpan(ctx, "visitHorizontalPodAutoscaler")
	defer span.End()

	hpa := &autoscalingv1.HorizontalPodAutoscaler{}
	if err := convertToType(object, hpa); err != nil {
		return err
	}

	var g errgroup.Group

	g.Go(func() error {
		target, err := s.queryer.ScaleTarget(ctx, hpa)
		if err != nil {
			return err
		}

		if target != nil {
			g.Go(func() error {
				u := &unstructured.Unstructured{Object: target}
				if err := visitor.Visit(ctx, u, handler, true); err != nil {
					return errors.Wrapf(err, "horizontal pod scaler %s visit scale target",
						kubernetes.PrintObject(hpa))
				}

				return handler.AddEdge(ctx, object, u)
			})
		}

		return nil
	})

	return g.Wait()
}
