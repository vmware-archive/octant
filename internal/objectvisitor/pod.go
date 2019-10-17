package objectvisitor

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
)

// Pod is a typed visitor for pods.
type Pod struct {
	queryer queryer.Queryer
}

var _ TypedVisitor = (*Pod)(nil)

// NewPod creates an instance of Pod.
func NewPod(q queryer.Queryer) *Pod {
	return &Pod{
		queryer: q,
	}
}

// Support returns the gvk this typed visitor supports.
func (p *Pod) Supports() schema.GroupVersionKind {
	return gvk.Pod
}

// Visit visits a pod. It looks for service accounts and services.
func (p *Pod) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool) error {
	ctx, span := trace.StartSpan(ctx, "visitPod")
	defer span.End()

	if p.queryer == nil {
		return errors.New("queryer is nil")
	}

	pod := &corev1.Pod{}
	if err := convertToType(object, pod); err != nil {
		return err
	}

	var g errgroup.Group

	g.Go(func() error {
		services, err := p.queryer.ServicesForPod(ctx, pod)
		if err != nil {
			return err
		}

		for i := range services {
			service := services[i]
			g.Go(func() error {
				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(service)
				if err != nil {
					return err
				}
				u := &unstructured.Unstructured{Object: m}
				if err := visitor.Visit(ctx, u, handler, true); err != nil {
					return errors.Wrapf(err, "pod %s visit service %s",
						kubernetes.PrintObject(pod), kubernetes.PrintObject(service))
				}

				return handler.AddEdge(ctx, object, u)
			})
		}

		return nil
	})
	g.Go(func() error {
		if pod.Spec.ServiceAccountName != "" {
			serviceAccount, err := p.queryer.ServiceAccountForPod(ctx, pod)
			if err != nil {
				return err
			}

			m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(serviceAccount)
			if err != nil {
				return err
			}
			u := &unstructured.Unstructured{Object: m}

			if serviceAccount != nil {
				if err := visitor.Visit(ctx, u, handler, true); err != nil {
					return errors.Wrapf(err, "pod %s visit service account %s",
						kubernetes.PrintObject(pod), kubernetes.PrintObject(serviceAccount))
				}
				return handler.AddEdge(ctx, object, u)
			}
		}

		return nil
	})

	return g.Wait()
}
