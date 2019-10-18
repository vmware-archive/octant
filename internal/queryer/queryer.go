/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package queryer

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"

	"github.com/vmware-tanzu/octant/internal/gvk"
	dashstrings "github.com/vmware-tanzu/octant/internal/util/strings"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
)

//go:generate mockgen -destination=./fake/mock_queryer.go -package=fake github.com/vmware-tanzu/octant/internal/queryer Queryer
//go:generate mockgen -source=../../vendor/k8s.io/client-go/discovery/discovery_client.go -imports=openapi_v2=github.com/googleapis/gnostic/OpenAPIv2 -destination=./fake/mock_discovery.go -package=fake k8s.io/client-go/discovery DiscoveryInterface

type Queryer interface {
	Children(ctx context.Context, object *unstructured.Unstructured) (*unstructured.UnstructuredList, error)
	Events(ctx context.Context, object metav1.Object) ([]*corev1.Event, error)
	IngressesForService(ctx context.Context, service *corev1.Service) ([]*extv1beta1.Ingress, error)
	OwnerReference(ctx context.Context, object *unstructured.Unstructured) (bool, *unstructured.Unstructured, error)
	ScaleTarget(ctx context.Context, hpa *autoscalingv1.HorizontalPodAutoscaler) (map[string]interface{}, error)
	PodsForService(ctx context.Context, service *corev1.Service) ([]*corev1.Pod, error)
	ServicesForIngress(ctx context.Context, ingress *extv1beta1.Ingress) (*unstructured.UnstructuredList, error)
	ServicesForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.Service, error)
	ServiceAccountForPod(ctx context.Context, pod *corev1.Pod) (*corev1.ServiceAccount, error)
}

type childrenCache struct {
	children map[types.UID]*unstructured.UnstructuredList
	mu       sync.RWMutex
}

func initChildrenCache() *childrenCache {
	return &childrenCache{
		children: make(map[types.UID]*unstructured.UnstructuredList),
	}
}

func (c *childrenCache) get(key types.UID) (*unstructured.UnstructuredList, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.children[key]
	return v, ok
}

func (c *childrenCache) set(key types.UID, value *unstructured.UnstructuredList) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.children[key] = value
}

type ownerCache struct {
	owner map[store.Key]*unstructured.Unstructured
	mu    sync.Mutex
}

func initOwnerCache() *ownerCache {
	return &ownerCache{
		owner: make(map[store.Key]*unstructured.Unstructured),
	}
}

func (c *ownerCache) set(key store.Key, value *unstructured.Unstructured) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value == nil {
		return
	}

	c.owner[key] = value
}

func (c *ownerCache) get(key store.Key) (*unstructured.Unstructured, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.owner[key]
	return v, ok
}

type podsForServicesCache struct {
	podsForServices map[types.UID][]*corev1.Pod
	mu              sync.Mutex
}

func initPodsForServicesCache() *podsForServicesCache {
	return &podsForServicesCache{
		podsForServices: make(map[types.UID][]*corev1.Pod),
	}
}

func (c *podsForServicesCache) set(key types.UID, value []*corev1.Pod) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.podsForServices[key] = value
}

func (c *podsForServicesCache) get(key types.UID) ([]*corev1.Pod, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.podsForServices[key]
	return v, ok
}

type ObjectStoreQueryer struct {
	objectStore     store.Store
	discoveryClient discovery.DiscoveryInterface

	children        *childrenCache
	podsForServices *podsForServicesCache
	owner           *ownerCache

	// mu sync.Mutex
}

var _ Queryer = (*ObjectStoreQueryer)(nil)

func New(o store.Store, discoveryClient discovery.DiscoveryInterface) *ObjectStoreQueryer {
	return &ObjectStoreQueryer{
		objectStore:     o,
		discoveryClient: discoveryClient,

		children:        initChildrenCache(),
		podsForServices: initPodsForServicesCache(),
		owner:           initOwnerCache(),
	}
}

func (osq *ObjectStoreQueryer) Children(ctx context.Context, owner *unstructured.Unstructured) (*unstructured.UnstructuredList, error) {
	if owner == nil {
		return nil, errors.New("owner is nil")
	}

	ctx, span := trace.StartSpan(ctx, "queryer:Children")
	defer span.End()

	stored, ok := osq.children.get(owner.GetUID())

	if ok {
		return stored, nil
	}

	out := &unstructured.UnstructuredList{}

	ch := make(chan *unstructured.Unstructured)
	childrenProcessed := make(chan bool, 1)
	go func() {
		for child := range ch {
			if child == nil {
				continue
			}
			out.Items = append(out.Items, *child)
		}
		childrenProcessed <- true
	}()

	list := append(allowed[:0:0], allowed...)

	crds, _, err := navigation.CustomResourceDefinitions(ctx, osq.objectStore)
	if err == nil {
		for _, crd := range crds {
			for _, version := range crd.Spec.Versions {
				list = append(list, schema.GroupVersionKind{
					Group:   crd.Spec.Group,
					Version: version.Name,
					Kind:    crd.Spec.Names.Kind,
				})

			}
		}
	}

	resourceLists, err := osq.discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	var g errgroup.Group

	sem := semaphore.NewWeighted(5)

	for resourceListIndex := range resourceLists {
		resourceList := resourceLists[resourceListIndex]
		if resourceList == nil {
			continue
		}

		for i := range resourceList.APIResources {
			apiResource := resourceList.APIResources[i]
			if !apiResource.Namespaced {
				continue
			}

			gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
			if err != nil {
				return nil, err
			}

			found := false
			for i := range list {
				if list[i].Group == gv.Group &&
					list[i].Version == gv.Version &&
					list[i].Kind == apiResource.Kind {
					found = true
				}
			}

			if !found {
				continue
			}

			key := store.Key{
				Namespace:  owner.GetNamespace(),
				APIVersion: resourceList.GroupVersion,
				Kind:       apiResource.Kind,
			}

			if osq.canList(apiResource) {
				continue
			}

			g.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return err
				}
				defer sem.Release(1)
				objects, _, err := osq.objectStore.List(ctx, key)
				if err != nil {
					return errors.Wrapf(err, "unable to retrieve %+v", key)
				}

				for i := range objects.Items {
					if metav1.IsControlledBy(&objects.Items[i], owner) {
						ch <- &objects.Items[i]
					}
				}

				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		if err != context.Canceled {
			return nil, errors.Wrap(err, "find children")
		}
	}

	close(ch)
	<-childrenProcessed
	close(childrenProcessed)

	osq.children.set(owner.GetUID(), out)

	return out, nil
}

var allowed = []schema.GroupVersionKind{
	gvk.CronJob,
	gvk.DaemonSet,
	gvk.Deployment,
	gvk.Pod,
	gvk.Job,
	gvk.ExtReplicaSet,
	gvk.ReplicationController,
	gvk.StatefulSet,
	gvk.HorizontalPodAutoscaler,
	gvk.Ingress,
	gvk.Service,
	gvk.ConfigMap,
	gvk.PersistentVolumeClaim,
	gvk.Secret,
	gvk.ServiceAccount,
}

func (osq *ObjectStoreQueryer) canList(apiResource metav1.APIResource) bool {
	return !dashstrings.Contains("watch", apiResource.Verbs) ||
		!dashstrings.Contains("list", apiResource.Verbs)
}

func (osq *ObjectStoreQueryer) Events(ctx context.Context, object metav1.Object) ([]*corev1.Event, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, err
	}

	u := &unstructured.Unstructured{Object: m}

	key := store.Key{
		Namespace:  u.GetNamespace(),
		APIVersion: "v1",
		Kind:       "Event",
	}

	allEvents, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, err
	}

	var events []*corev1.Event
	for _, unstructuredEvent := range allEvents.Items {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredEvent.Object, event)
		if err != nil {
			return nil, err
		}

		involvedObject := event.InvolvedObject
		if involvedObject.Namespace == u.GetNamespace() &&
			involvedObject.APIVersion == u.GetAPIVersion() &&
			involvedObject.Kind == u.GetKind() &&
			involvedObject.Name == u.GetName() {
			events = append(events, event)
		}
	}

	return events, nil
}

func (osq *ObjectStoreQueryer) IngressesForService(ctx context.Context, service *corev1.Service) ([]*v1beta1.Ingress, error) {
	if service == nil {
		return nil, errors.New("nil service")
	}

	key := store.Key{
		Namespace:  service.Namespace,
		APIVersion: "extensions/v1beta1",
		Kind:       "Ingress",
	}
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving ingresses")
	}

	var results []*v1beta1.Ingress

	for i := range ul.Items {
		ingress := &v1beta1.Ingress{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(ul.Items[i].Object, ingress)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured ingress")
		}
		if err = copyObjectMeta(ingress, &ul.Items[i]); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		backends := osq.listIngressBackends(*ingress)
		if !containsBackend(backends, service.Name) {
			continue
		}

		results = append(results, ingress)
	}
	return results, nil
}

func (osq *ObjectStoreQueryer) listIngressBackends(ingress v1beta1.Ingress) []extv1beta1.IngressBackend {
	var backends []v1beta1.IngressBackend

	if ingress.Spec.Backend != nil && ingress.Spec.Backend.ServiceName != "" {
		backends = append(backends, *ingress.Spec.Backend)
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}
		for _, p := range rule.IngressRuleValue.HTTP.Paths {
			if p.Backend.ServiceName == "" {
				continue
			}
			backends = append(backends, p.Backend)
		}
	}

	return backends
}

func (osq *ObjectStoreQueryer) OwnerReference(ctx context.Context, object *unstructured.Unstructured) (bool, *unstructured.Unstructured, error) {
	if object == nil {
		return false, nil, errors.New("can't find owner for nil object")
	}

	ownerReferences := object.GetOwnerReferences()
	switch len(ownerReferences) {
	case 0:
		return false, nil, nil
	case 1:
		ownerReference := ownerReferences[0]

		resourceList, err := osq.discoveryClient.ServerResourcesForGroupVersion(ownerReference.APIVersion)
		if err != nil {
			return false, nil, err
		}
		if resourceList == nil {
			return false, nil, errors.Errorf("did not expect resource list for %s to be nil", ownerReference.APIVersion)
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
			return false, nil, errors.Errorf("unable to find owner references %v", ownerReference)
		}

		namespace := ""
		if isNamespaced {
			namespace = object.GetNamespace()
		}

		key := store.Key{
			Namespace:  namespace,
			APIVersion: ownerReference.APIVersion,
			Kind:       ownerReference.Kind,
			Name:       ownerReference.Name,
		}

		object, ok := osq.owner.get(key)
		if ok {
			return true, object, nil
		}

		owner, found, err := osq.objectStore.Get(ctx, key)
		if err != nil {
			return false, nil, errors.Wrap(err, "get owner from store")
		}

		if !found {
			return false, nil, errors.Errorf("owner %s not found", key)
		}

		osq.owner.set(key, owner)

		return true, owner, nil
	default:
		return false, nil, errors.New("unable to handle more than one owner reference")
	}
}

func (osq *ObjectStoreQueryer) ScaleTarget(ctx context.Context, hpa *autoscalingv1.HorizontalPodAutoscaler) (map[string]interface{}, error) {
	if hpa == nil {
		return nil, errors.New("can't find scale target for nil hpa")
	}

	key := store.Key{
		Namespace:  hpa.Namespace,
		APIVersion: hpa.Spec.ScaleTargetRef.APIVersion,
		Kind:       hpa.Spec.ScaleTargetRef.Kind,
		Name:       hpa.Spec.ScaleTargetRef.Name,
	}

	u, found, err := osq.objectStore.Get(ctx, key)
	if err != nil {
		return nil, errors.WithMessagef(err, "retrieve scale target %q from namespace %q", key.Name, key.Namespace)
	}

	if found {
		switch key.Kind {
		case "Deployment":
			deployment := &appsv1.Deployment{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, deployment); err != nil {
				return nil, errors.WithMessage(err, "converting unstructured object to deployment")
			}

			object, err := runtime.DefaultUnstructuredConverter.ToUnstructured(deployment)
			if err != nil {
				return nil, err
			}
			return object, nil
		case "ReplicaSet":
			replicaSet := &appsv1.ReplicaSet{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, replicaSet); err != nil {
				return nil, errors.WithMessage(err, "converting unstructured object to replica set")
			}

			object, err := runtime.DefaultUnstructuredConverter.ToUnstructured(replicaSet)
			if err != nil {
				return nil, err
			}
			return object, nil
		case "ReplicationController":
			replicationController := &corev1.ReplicationController{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, replicationController); err != nil {
				return nil, errors.WithMessage(err, "converting unstructured object to replication controller")
			}

			object, err := runtime.DefaultUnstructuredConverter.ToUnstructured(replicationController)
			if err != nil {
				return nil, err
			}
			return object, nil
		}
	}

	return nil, errors.Wrap(err, "invalid scale target")
}

func (osq *ObjectStoreQueryer) PodsForService(ctx context.Context, service *corev1.Service) ([]*corev1.Pod, error) {
	if service == nil {
		return nil, errors.New("nil service")
	}

	stored, ok := osq.podsForServices.get(service.UID)
	if ok {
		return stored, nil
	}

	key := store.Key{
		Namespace:  service.Namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	selector, err := osq.getSelector(service)
	if err != nil {
		return nil, errors.Wrapf(err, "creating pod selector for service: %v", service.Name)
	}
	pods, err := osq.loadPods(ctx, key, selector)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching pods for service: %v", service.Name)
	}

	osq.podsForServices.set(service.UID, pods)

	return pods, nil
}

func (osq *ObjectStoreQueryer) loadPods(ctx context.Context, key store.Key, labelSelector *metav1.LabelSelector) ([]*corev1.Pod, error) {
	objects, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, err
	}

	var list []*corev1.Pod

	for i := range objects.Items {
		pod := &corev1.Pod{}
		if err := scheme.Scheme.Convert(&objects.Items[i], pod, runtime.InternalGroupVersioner); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(pod, &objects.Items[i]); err != nil {
			return nil, err
		}

		podSelector := &metav1.LabelSelector{
			MatchLabels: pod.GetLabels(),
		}

		selector, err := metav1.LabelSelectorAsSelector(labelSelector)
		if err != nil {
			return nil, err
		}

		if selector == nil || isEqualSelector(labelSelector, podSelector) || selector.Matches(kLabels.Set(pod.Labels)) {
			list = append(list, pod)
		}
	}

	return list, nil
}

func (osq *ObjectStoreQueryer) ServicesForIngress(ctx context.Context, ingress *extv1beta1.Ingress) (*unstructured.UnstructuredList, error) {
	if ingress == nil {
		return nil, errors.New("ingress is nil")
	}

	backends := osq.listIngressBackends(*ingress)
	list := &unstructured.UnstructuredList{}
	for _, backend := range backends {
		key := store.Key{
			Namespace:  ingress.Namespace,
			APIVersion: "v1",
			Kind:       "Service",
			Name:       backend.ServiceName,
		}
		u, found, err := osq.objectStore.Get(ctx, key)
		if err != nil {
			return nil, errors.Wrapf(err, "retrieving service backend: %v", backend)
		}

		if !found {
			continue
		}

		list.Items = append(list.Items, *u)
	}
	return list, nil
}

func (osq *ObjectStoreQueryer) ServicesForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.Service, error) {
	var results []*corev1.Service
	if pod == nil {
		return nil, errors.New("nil pod")
	}

	key := store.Key{
		Namespace:  pod.Namespace,
		APIVersion: "v1",
		Kind:       "Service",
	}
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving services")
	}
	for i := range ul.Items {
		svc := &corev1.Service{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(ul.Items[i].Object, svc)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured service")
		}
		if err = copyObjectMeta(svc, &ul.Items[i]); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		labelSelector, err := osq.getSelector(svc)
		if err != nil {
			return nil, errors.Wrapf(err, "creating pod selector for service: %v", svc.Name)
		}
		selector, err := metav1.LabelSelectorAsSelector(labelSelector)
		if err != nil {
			return nil, errors.Wrap(err, "invalid selector")
		}

		if selector.Empty() || !selector.Matches(kLabels.Set(pod.Labels)) {
			continue
		}
		results = append(results, svc)
	}
	return results, nil
}

func (osq *ObjectStoreQueryer) ServiceAccountForPod(ctx context.Context, pod *corev1.Pod) (*corev1.ServiceAccount, error) {
	if pod == nil {
		return nil, errors.New("pod is nil")
	}

	if pod.Spec.ServiceAccountName == "" {
		return nil, nil
	}

	key := store.Key{
		Namespace:  pod.Namespace,
		APIVersion: "v1",
		Kind:       "ServiceAccount",
		Name:       pod.Spec.ServiceAccountName,
	}

	u, found, err := osq.objectStore.Get(ctx, key)
	if err != nil {
		return nil, errors.WithMessagef(err, "retrieve service account %q from namespace %q",
			key.Name, key.Namespace)
	}

	if !found {
		return nil, errors.Errorf("service account %q from namespace %q does not exist",
			key.Name, key.Namespace)
	}

	serviceAccount := &corev1.ServiceAccount{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, serviceAccount); err != nil {
		return nil, errors.WithMessage(err, "converting unstructured object to service account")
	}

	if err = copyObjectMeta(serviceAccount, u); err != nil {
		return nil, errors.Wrap(err, "copying object metadata")
	}

	return serviceAccount, nil

}

func (osq *ObjectStoreQueryer) getSelector(object runtime.Object) (*metav1.LabelSelector, error) {
	switch t := object.(type) {
	case *appsv1.DaemonSet:
		return t.Spec.Selector, nil
	case *appsv1.StatefulSet:
		return t.Spec.Selector, nil
	case *batchv1beta1.CronJob:
		return nil, nil
	case *corev1.ReplicationController:
		selector := &metav1.LabelSelector{
			MatchLabels: t.Spec.Selector,
		}
		return selector, nil
	case *v1beta1.ReplicaSet:
		return t.Spec.Selector, nil
	case *appsv1.ReplicaSet:
		return t.Spec.Selector, nil
	case *appsv1.Deployment:
		return t.Spec.Selector, nil
	case *corev1.Service:
		selector := &metav1.LabelSelector{
			MatchLabels: t.Spec.Selector,
		}
		return selector, nil
	case *apps.DaemonSet:
		return t.Spec.Selector, nil
	case *apps.StatefulSet:
		return t.Spec.Selector, nil
	case *batch.CronJob:
		return nil, nil
	case *core.ReplicationController:
		selector := &metav1.LabelSelector{
			MatchLabels: t.Spec.Selector,
		}
		return selector, nil
	case *apps.ReplicaSet:
		return t.Spec.Selector, nil
	case *apps.Deployment:
		return t.Spec.Selector, nil
	case *core.Service:
		selector := &metav1.LabelSelector{
			MatchLabels: t.Spec.Selector,
		}
		return selector, nil
	default:
		return nil, errors.Errorf("unable to retrieve selector for type %T", object)
	}
}

func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return errors.Errorf("%T is not an object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return errors.Wrapf(err, "accessing type meta")
	}
	t.SetAPIVersion(from.GetAPIVersion())
	t.SetKind(from.GetObjectKind().GroupVersionKind().Kind)

	object.SetNamespace(from.GetNamespace())
	object.SetName(from.GetName())
	object.SetGenerateName(from.GetGenerateName())
	object.SetUID(from.GetUID())
	object.SetResourceVersion(from.GetResourceVersion())
	object.SetGeneration(from.GetGeneration())
	object.SetSelfLink(from.GetSelfLink())
	object.SetCreationTimestamp(from.GetCreationTimestamp())
	object.SetDeletionTimestamp(from.GetDeletionTimestamp())
	object.SetDeletionGracePeriodSeconds(from.GetDeletionGracePeriodSeconds())
	object.SetLabels(from.GetLabels())
	object.SetAnnotations(from.GetAnnotations())
	object.SetOwnerReferences(from.GetOwnerReferences())
	object.SetClusterName(from.GetClusterName())
	object.SetFinalizers(from.GetFinalizers())

	return nil
}

// extraKeys are keys that should be ignored in labels. These keys are added
// by tools or by Kubernetes itself.
var extraKeys = []string{
	"statefulset.kubernetes.io/pod-name",
	appsv1.DefaultDeploymentUniqueLabelKey,
	"controller-revision-hash",
	"pod-template-generation",
}

func isEqualSelector(s1, s2 *metav1.LabelSelector) bool {
	s1Copy := s1.DeepCopy()
	s2Copy := s2.DeepCopy()

	for _, key := range extraKeys {
		delete(s1Copy.MatchLabels, key)
		delete(s2Copy.MatchLabels, key)
	}

	return apiequality.Semantic.DeepEqual(s1Copy, s2Copy)
}

func containsBackend(lst []v1beta1.IngressBackend, s string) bool {
	for _, item := range lst {
		if item.ServiceName == s {
			return true
		}
	}
	return false
}
