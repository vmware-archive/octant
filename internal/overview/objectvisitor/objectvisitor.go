package objectvisitor

import (
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/heptio/developer-dash/internal/gvk"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// ClusterObject is a cluster object.
// NOTE: this might not be the most succinct description.
type ClusterObject interface {
	runtime.Object
}

// ObjectHandler performs actions on an object. Can be used to augment
// visitor actions with extra functionality.
type ObjectHandler interface {
	AddChild(parent ClusterObject, children ...ClusterObject) error
	Process(object ClusterObject) error
}

// Visitor is a visitor for cluster objects. It will visit an object and all of
// its ancestors and descendants.
type Visitor interface {
	Visit(object ClusterObject) error
}

// ObjectHandlerFactory creates ObjectHandler given a ClusterObject.
type ObjectHandlerFactory func(ClusterObject) (ObjectHandler, error)

// DefaultFactoryGenerator generates ObjectHandlerFactory based on GVK.
type DefaultFactoryGenerator struct {
	m map[schema.GroupVersionKind]ObjectHandlerFactory
}

// NewDefaultFactoryGenerator creates an instance of NewDefaultFactoryGenerator.
func NewDefaultFactoryGenerator() *DefaultFactoryGenerator {
	return &DefaultFactoryGenerator{
		m: make(map[schema.GroupVersionKind]ObjectHandlerFactory),
	}
}

// Register registers an ObjectHandlerFactory for a GVK.
func (dfg *DefaultFactoryGenerator) Register(gvk schema.GroupVersionKind, fn ObjectHandlerFactory) error {
	if _, ok := dfg.m[gvk]; ok {
		return errors.Errorf("%s has already been registered", gvk)
	}

	dfg.m[gvk] = fn

	return nil
}

// FactoryFunc creates an ObjectHandlerFactory for a GVK.
func (dfg *DefaultFactoryGenerator) FactoryFunc() ObjectHandlerFactory {
	return func(object ClusterObject) (ObjectHandler, error) {
		if object == nil {
			return nil, errors.New("unable to find factory for nil object")
		}

		objectKind := object.GetObjectKind()
		if objectKind == nil {
			return nil, errors.New("object kind is nil")
		}

		gvk := objectKind.GroupVersionKind()
		factory, ok := dfg.m[gvk]
		if !ok {
			return nil, errors.Errorf("%s was not registered",
				gvk)
		}

		return factory(object)
	}
}

// DefaultVisitor is the default implementation of Visitor.
type DefaultVisitor struct {
	queryer        queryer.Queryer
	handlerFactory ObjectHandlerFactory
	visited        map[types.UID]bool
	visitedMu      sync.Mutex
}

var _ Visitor = (*DefaultVisitor)(nil)

// NewDefaultVisitor creates an instance of DefaultVisitor.
func NewDefaultVisitor(q queryer.Queryer, factory ObjectHandlerFactory) (*DefaultVisitor, error) {
	if factory == nil {
		return nil, errors.Errorf("factory was nil")
	}

	return &DefaultVisitor{
		queryer:        q,
		handlerFactory: factory,
		visited:        make(map[types.UID]bool),
	}, nil
}

// hasVisited returns true if this object has already been visited. If the
// object has not been visited, it returns false, and sets the object
// visit status to true.
func (dv *DefaultVisitor) hasVisited(object runtime.Object) (bool, error) {
	dv.visitedMu.Lock()
	defer dv.visitedMu.Unlock()

	accessor := meta.NewAccessor()
	uid, err := accessor.UID(object)
	if err != nil {
		return false, errors.Wrap(err, "get uid from object")
	}

	if _, ok := dv.visited[uid]; ok {
		return true, nil
	}

	dv.visited[uid] = true

	return false, nil
}

// Visit visits a ClusterObject.
func (dv *DefaultVisitor) Visit(object ClusterObject) error {
	if object == nil {
		return errors.New("trying to visit a nil object")
	}

	hasVisited, err := dv.hasVisited(object)
	if err != nil {
		return errors.Wrap(err, "check for visit object")
	}
	if hasVisited {
		return nil
	}

	// Create a handler factory for this object. This allows the visitor's caller to
	// interact with the ancestors and descendants of the object.
	o, err := dv.handlerFactory(object)
	if err != nil {
		return err
	}

	return dv.visitObject(object, o)
}

// visitIngress visits an ingress' service backends.
func (dv *DefaultVisitor) visitIngress(ingress *extv1beta1.Ingress) ([]ClusterObject, error) {
	services, err := dv.queryer.ServicesForIngress(ingress)
	if err != nil {
		return nil, err
	}

	var children []ClusterObject

	for _, service := range services {
		if err := dv.Visit(service); err != nil {
			return nil, err
		}

		children = append(children, service)
	}

	return children, nil
}

// visitPod visits a pod's services.
func (dv *DefaultVisitor) visitPod(pod *corev1.Pod) error {
	services, err := dv.queryer.ServicesForPod(pod)
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := dv.Visit(service); err != nil {
			return err
		}
	}

	return nil
}

// visitService visits a service's ingresses and pods.
func (dv *DefaultVisitor) visitService(service *corev1.Service) ([]ClusterObject, error) {
	pods, err := dv.queryer.PodsForService(service)
	if err != nil {
		return nil, err
	}

	var children []ClusterObject

	for _, pod := range pods {
		if err := dv.Visit(pod); err != nil {
			return nil, err
		}

		children = append(children, pod)
	}

	ingresses, err := dv.queryer.IngressesForService(service)
	if err != nil {
		return nil, err
	}

	for _, ingress := range ingresses {
		if err := dv.Visit(ingress); err != nil {
			return nil, err
		}
	}

	return children, nil
}

// handleObject attempts to visit parents and children of the object.
func (dv *DefaultVisitor) handleObject(object ClusterObject, visitorObject ObjectHandler) error {
	if object == nil {
		return errors.New("trying to visit a nil object")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return err
	}

	u := &unstructured.Unstructured{Object: m}

	var g errgroup.Group

	for _, ownerReference := range u.GetOwnerReferences() {
		g.Go(func() error {
			o, err := dv.queryer.OwnerReference(u.GetNamespace(), ownerReference)
			if err != nil {
				return err
			}

			owner := o.(ClusterObject)

			if object == nil {
				return nil
			}

			if err := dv.Visit(owner); err != nil {
				return err
			}

			return nil
		})
	}

	children, err := dv.queryer.Children(u)
	if err != nil {
		return err
	}

	for i := range children {
		index := i
		g.Go(func() error {
			o := children[index].(ClusterObject)
			if err := dv.Visit(o); err != nil {
				return err
			}

			if err := visitorObject.AddChild(object, o); err != nil {
				return errors.Wrap(err, "add child")
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return visitorObject.Process(object)
}

// visitObject visits an object. If the object is a service, ingress, or pod, it
// also runs custom visitor code for them.
func (dv *DefaultVisitor) visitObject(object ClusterObject, visitorObject ObjectHandler) error {

	if object == nil {
		return errors.New("can't visit a nil object")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return err
	}

	u := &unstructured.Unstructured{Object: m}

	apiVersion := u.GetAPIVersion()
	kind := u.GetKind()

	objectGVK := schema.FromAPIVersionAndKind(apiVersion, kind)

	switch objectGVK {
	case gvk.IngressGVK:
		ingress := &extv1beta1.Ingress{}
		if err := dv.convertToType(u, ingress); err != nil {
			return err
		}
		children, err := dv.visitIngress(ingress)
		if err != nil {
			return err
		}

		if err := visitorObject.AddChild(object, children...); err != nil {
			return err
		}
	case gvk.PodGVK:
		pod := &corev1.Pod{}
		if err := dv.convertToType(u, pod); err != nil {
			return err
		}
		if err := dv.visitPod(pod); err != nil {
			return err
		}
	case gvk.ServiceGVK:
		service := &corev1.Service{}
		if err := dv.convertToType(u, service); err != nil {
			return err
		}
		children, err := dv.visitService(service)
		if err != nil {
			return err
		}

		if err := visitorObject.AddChild(object, children...); err != nil {
			return err
		}
	}

	return dv.handleObject(object, visitorObject)
}

func (dv *DefaultVisitor) convertToType(object runtime.Object, objectType interface{}) error {
	u, ok := object.(*unstructured.Unstructured)
	if !ok {
		return errors.Errorf("object is not an unstructured (%T)", object)
	}

	return runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, objectType)
}
