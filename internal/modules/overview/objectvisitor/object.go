package objectvisitor

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

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
func (o *Object) Visit(ctx context.Context, object runtime.Object, handler ObjectHandler, visitor Visitor) error {
	ctx, span := trace.StartSpan(ctx, "handleObject")
	defer span.End()

	if object == nil {
		return errors.New("trying to visit a nil object")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return err
	}

	unstructuredObject := &unstructured.Unstructured{Object: m}

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("apiVersion", unstructuredObject.GetAPIVersion()),
		trace.StringAttribute("kind", unstructuredObject.GetKind()),
		trace.StringAttribute("name", unstructuredObject.GetName()),
		trace.StringAttribute("namespace", unstructuredObject.GetNamespace()),
	}, "handling object")

	var g errgroup.Group

	references := unstructuredObject.GetOwnerReferences()
	for i := range references {
		ownerReference := references[i]
		g.Go(func() error {
			clusterClient := o.dashConfig.ClusterClient()
			discoveryClient, err := clusterClient.DiscoveryClient()
			if err != nil {
				return err
			}

			resourceList, err := discoveryClient.ServerResourcesForGroupVersion(ownerReference.APIVersion)
			if err != nil {
				return errors.Wrapf(err, "get resource list for %s", ownerReference.APIVersion)
			}

			if resourceList == nil {
				return errors.Errorf("did not expected resource list for %s to be nil", ownerReference.APIVersion)
			}

			found := false
			isNamespaced := false
			for _, apiResource := range resourceList.APIResources {
				if apiResource.Kind == ownerReference.Kind {
					isNamespaced = apiResource.Namespaced
					found = true
				}
			}

			if !found {
				return errors.Errorf("unable to find owner references %v", ownerReference)
			}

			namespace := ""
			if isNamespaced {
				namespace = unstructuredObject.GetNamespace()
			}

			owner, err := o.queryer.OwnerReference(ctx, namespace, ownerReference)
			if err != nil {
				return errors.Wrapf(err, "unable to check owner reference for %s", kubernetes.PrintObject(unstructuredObject))
			}

			if owner == nil {
				return errors.Errorf("unable to find owner for %s", unstructuredObject)
			}

			if err := visitor.Visit(ctx, owner, handler); err != nil {
				return errors.Wrapf(err, "visit ancestor %s for %s",
					kubernetes.PrintObject(owner),
					kubernetes.PrintObject(unstructuredObject))
			}

			return handler.AddEdge(unstructuredObject, owner)
		})
	}

	children, err := o.queryer.Children(ctx, unstructuredObject)
	if err != nil {
		return err
	}

	for i := range children {
		child := children[i]
		g.Go(func() error {
			if err := visitor.Visit(ctx, child, handler); err != nil {
				return errors.Wrapf(err, "visit child %s for %s",
					kubernetes.PrintObject(child),
					kubernetes.PrintObject(unstructuredObject))
			}

			return handler.AddEdge(object, child)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return handler.Process(ctx, object)
}
