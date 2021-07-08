/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package queryer

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	dashstrings "github.com/vmware-tanzu/octant/internal/util/strings"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
)

//go:generate mockgen -destination=./fake/mock_queryer.go -package=fake github.com/vmware-tanzu/octant/internal/queryer Queryer
//go:generate mockgen -destination=./fake/mock_discovery.go -package=fake k8s.io/client-go/discovery DiscoveryInterface

type Queryer interface {
	Children(ctx context.Context, object *unstructured.Unstructured) (*unstructured.UnstructuredList, error)
	Events(ctx context.Context, object metav1.Object) ([]*corev1.Event, error)
	IngressesForService(ctx context.Context, service *corev1.Service) ([]*networkingv1.Ingress, error)
	APIServicesForService(ctx context.Context, service *corev1.Service) ([]*apiregistrationv1.APIService, error)
	MutatingWebhookConfigurationsForService(ctx context.Context, service *corev1.Service) ([]*admissionregistrationv1.MutatingWebhookConfiguration, error)
	ValidatingWebhookConfigurationsForService(ctx context.Context, service *corev1.Service) ([]*admissionregistrationv1.ValidatingWebhookConfiguration, error)
	OwnerReference(ctx context.Context, object *unstructured.Unstructured) (bool, []*unstructured.Unstructured, error)
	ScaleTarget(ctx context.Context, hpa *autoscalingv1.HorizontalPodAutoscaler) (map[string]interface{}, error)
	PodsForService(ctx context.Context, service *corev1.Service) ([]*corev1.Pod, error)
	ServicesForIngress(ctx context.Context, ingress *networkingv1.Ingress) (*unstructured.UnstructuredList, error)
	ServicesForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.Service, error)
	ServiceAccountForPod(ctx context.Context, pod *corev1.Pod) (*corev1.ServiceAccount, error)
	ConfigMapsForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.ConfigMap, error)
	SecretsForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.Secret, error)
	PersistentVolumeClaimsForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.PersistentVolumeClaim, error)
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
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return nil, fmt.Errorf("objectStoreQueryer children: %w", err)
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
	gvk.AppReplicaSet,
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
	gvk.NetworkPolicy,
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
		err := kubernetes.FromUnstructured(&unstructuredEvent, event)
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

func (osq *ObjectStoreQueryer) IngressesForService(ctx context.Context, service *corev1.Service) ([]*networkingv1.Ingress, error) {
	if service == nil {
		return nil, errors.New("nil service")
	}

	key := store.Key{
		Namespace:  service.Namespace,
		APIVersion: "networking.k8s.io/v1",
		Kind:       "Ingress",
	}
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving ingresses")
	}

	var results []*networkingv1.Ingress

	for i := range ul.Items {
		ingress := &networkingv1.Ingress{}
		err := kubernetes.FromUnstructured(&ul.Items[i], ingress)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured ingress")
		}
		backends := osq.listIngressBackends(*ingress)
		if !containsBackend(backends, service.Name) {
			continue
		}

		results = append(results, ingress)
	}
	return results, nil
}

func (osq *ObjectStoreQueryer) listIngressBackends(ingress networkingv1.Ingress) []networkingv1.IngressBackend {
	var backends []networkingv1.IngressBackend

	if ingress.Spec.DefaultBackend != nil && ingress.Spec.DefaultBackend.Service.Name != "" {
		backends = append(backends, *ingress.Spec.DefaultBackend)
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}
		for _, p := range rule.IngressRuleValue.HTTP.Paths {
			if p.Backend.Service.Name == "" {
				continue
			}
			backends = append(backends, p.Backend)
		}
	}

	return backends
}

func (osq *ObjectStoreQueryer) APIServicesForService(ctx context.Context, service *corev1.Service) ([]*apiregistrationv1.APIService, error) {
	if service == nil {
		return nil, errors.New("nil service")
	}

	key := store.KeyFromGroupVersionKind(gvk.APIService)
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving apiservices")
	}

	var results []*apiregistrationv1.APIService

	for i := range ul.Items {
		apiservice := &apiregistrationv1.APIService{}
		err := kubernetes.FromUnstructured(&ul.Items[i], apiservice)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured apiservice")
		}
		if apiservice.Spec.Service != nil &&
			apiservice.Spec.Service.Namespace == service.Namespace &&
			apiservice.Spec.Service.Name == service.Name {
			results = append(results, apiservice)
		}
	}

	return results, nil
}

func (osq *ObjectStoreQueryer) MutatingWebhookConfigurationsForService(ctx context.Context, service *corev1.Service) ([]*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	if service == nil {
		return nil, errors.New("nil service")
	}

	key := store.KeyFromGroupVersionKind(gvk.MutatingWebhookConfiguration)
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving mutatingwebhookconfigurations")
	}

	var results []*admissionregistrationv1.MutatingWebhookConfiguration

	for i := range ul.Items {
		mutatingwebhookconfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{}
		err := kubernetes.FromUnstructured(&ul.Items[i], mutatingwebhookconfiguration)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured mutatingwebhookconfiguration")
		}
		for _, mutatingwebhook := range mutatingwebhookconfiguration.Webhooks {
			if mutatingwebhook.ClientConfig.Service != nil &&
				mutatingwebhook.ClientConfig.Service.Namespace == service.Namespace &&
				mutatingwebhook.ClientConfig.Service.Name == service.Name {
				results = append(results, mutatingwebhookconfiguration)
				break
			}
		}
	}

	return results, nil
}

func (osq *ObjectStoreQueryer) ValidatingWebhookConfigurationsForService(ctx context.Context, service *corev1.Service) ([]*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	if service == nil {
		return nil, errors.New("nil service")
	}

	key := store.KeyFromGroupVersionKind(gvk.ValidatingWebhookConfiguration)
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving validatingwebhookconfigurations")
	}

	var results []*admissionregistrationv1.ValidatingWebhookConfiguration

	for i := range ul.Items {
		validatingwebhookconfiguration := &admissionregistrationv1.ValidatingWebhookConfiguration{}
		err := kubernetes.FromUnstructured(&ul.Items[i], validatingwebhookconfiguration)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured validatingwebhookconfiguration")
		}
		for _, validatingwebhook := range validatingwebhookconfiguration.Webhooks {
			if validatingwebhook.ClientConfig.Service != nil &&
				validatingwebhook.ClientConfig.Service.Namespace == service.Namespace &&
				validatingwebhook.ClientConfig.Service.Name == service.Name {
				results = append(results, validatingwebhookconfiguration)
				break
			}
		}
	}

	return results, nil
}

func (osq *ObjectStoreQueryer) OwnerReference(ctx context.Context, object *unstructured.Unstructured) (bool, []*unstructured.Unstructured, error) {
	if object == nil {
		return false, nil, errors.New("can't find owner for nil object")
	}

	ownerReferences := object.GetOwnerReferences()
	switch len(ownerReferences) {
	case 0:
		return false, nil, nil
	default:
		var list []*unstructured.Unstructured

		found := false

		for _, ownerReference := range ownerReferences {
			objectFound, object, err := osq.handle(ctx, object, ownerReference)
			if err != nil {
				return false, nil, err
			}

			list = append(list, object)
			if objectFound {
				found = true
			}
		}

		return found, list, nil
	}
}

func (osq *ObjectStoreQueryer) handle(
	ctx context.Context,
	object *unstructured.Unstructured,
	ownerReference metav1.OwnerReference) (bool, *unstructured.Unstructured, error) {
	resourceList, err := osq.discoveryClient.ServerResourcesForGroupVersion(ownerReference.APIVersion)
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return false, nil, fmt.Errorf("objectStoreQueryer handle: %w", err)
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

	owner, err := osq.objectStore.Get(ctx, key)
	if err != nil {
		return false, nil, errors.Wrap(err, "get owner from store")
	}

	if owner == nil {
		return false, nil, errors.Errorf("owner %s not found", key)
	}

	osq.owner.set(key, owner)

	return true, owner, nil
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

	u, err := osq.objectStore.Get(ctx, key)
	if err != nil {
		return nil, errors.WithMessagef(err, "retrieve scale target %q from namespace %q", key.Name, key.Namespace)
	}

	if u != nil {
		switch key.Kind {
		case "Deployment":
			deployment := &appsv1.Deployment{}
			if err := kubernetes.FromUnstructured(u, deployment); err != nil {
				return nil, errors.WithMessage(err, "converting unstructured object to deployment")
			}

			object, err := runtime.DefaultUnstructuredConverter.ToUnstructured(deployment)
			if err != nil {
				return nil, err
			}
			return object, nil
		case "ReplicaSet":
			replicaSet := &appsv1.ReplicaSet{}
			if err := kubernetes.FromUnstructured(u, replicaSet); err != nil {
				return nil, errors.WithMessage(err, "converting unstructured object to replica set")
			}

			object, err := runtime.DefaultUnstructuredConverter.ToUnstructured(replicaSet)
			if err != nil {
				return nil, err
			}
			return object, nil
		case "ReplicationController":
			replicationController := &corev1.ReplicationController{}
			if err := kubernetes.FromUnstructured(u, replicationController); err != nil {
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
		if err := kubernetes.FromUnstructured(&objects.Items[i], pod); err != nil {
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

func (osq *ObjectStoreQueryer) ServicesForIngress(ctx context.Context, ingress *networkingv1.Ingress) (*unstructured.UnstructuredList, error) {
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
			Name:       backend.Service.Name,
		}
		u, err := osq.objectStore.Get(ctx, key)
		if err != nil && !kerrors.IsNotFound(err) {
			return nil, errors.Wrapf(err, "retrieving service backend: %v", backend)
		}

		if u == nil {
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
		err := kubernetes.FromUnstructured(&ul.Items[i], svc)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured service")
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

	u, err := osq.objectStore.Get(ctx, key)
	if err != nil {
		return nil, errors.WithMessagef(err, "retrieve service account %q from namespace %q",
			key.Name, key.Namespace)
	}

	if u == nil {
		return nil, errors.Errorf("service account %q from namespace %q does not exist",
			key.Name, key.Namespace)
	}

	serviceAccount := &corev1.ServiceAccount{}
	if err := kubernetes.FromUnstructured(u, serviceAccount); err != nil {
		return nil, errors.WithMessage(err, "converting unstructured object to service account")
	}

	return serviceAccount, nil

}

func (osq *ObjectStoreQueryer) ConfigMapsForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.ConfigMap, error) {
	if pod == nil {
		return nil, errors.New("pod is nil")
	}

	var configMaps []*corev1.ConfigMap
	key := store.Key{
		Namespace:  pod.Namespace,
		APIVersion: "v1",
		Kind:       "ConfigMap",
	}
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving configmaps")
	}

	for i := range ul.Items {
		configMap := &corev1.ConfigMap{}
		err := kubernetes.FromUnstructured(&ul.Items[i], configMap)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured configmap")
		}

		for _, v := range pod.Spec.Volumes {
			if v.ConfigMap != nil && v.ConfigMap.Name == configMap.Name {
				configMaps = append(configMaps, configMap)
			}
		}

		for ci := range pod.Spec.Containers {
			c := &pod.Spec.Containers[ci]
			for _, e := range c.Env {
				if e.ValueFrom != nil && e.ValueFrom.ConfigMapKeyRef != nil {
					ref := e.ValueFrom.ConfigMapKeyRef
					if ref.Name == configMap.Name {
						configMaps = append(configMaps, configMap)
					}
				}
			}

			for _, e := range c.EnvFrom {
				if e.ConfigMapRef != nil {
					ref := e.ConfigMapRef
					if ref.Name == configMap.Name {
						configMaps = append(configMaps, configMap)
					}
				}
			}
		}
	}

	return configMaps, nil
}

func (osq *ObjectStoreQueryer) SecretsForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.Secret, error) {
	if pod == nil {
		return nil, errors.New("pod is nil")
	}

	var secrets []*corev1.Secret
	key := store.Key{
		Namespace:  pod.Namespace,
		APIVersion: "v1",
		Kind:       "Secret",
	}
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving secrets")
	}

	for i := range ul.Items {
		secret := &corev1.Secret{}
		err := kubernetes.FromUnstructured(&ul.Items[i], secret)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured secret")
		}

		for vi := range pod.Spec.Volumes {
			v := &pod.Spec.Volumes[vi]
			if v.Secret != nil && v.Secret.SecretName == secret.Name {
				secrets = append(secrets, secret)
			}
		}
		for ci := range pod.Spec.Containers {
			c := &pod.Spec.Containers[ci]
			for _, e := range c.Env {
				if e.ValueFrom != nil && e.ValueFrom.SecretKeyRef != nil {
					ref := e.ValueFrom.SecretKeyRef
					if ref.Name == secret.Name {
						secrets = append(secrets, secret)
					}
				}
			}

			for _, e := range c.EnvFrom {
				if e.SecretRef != nil {
					ref := e.SecretRef
					if ref.Name == secret.Name {
						secrets = append(secrets, secret)
					}
				}
			}
		}
	}

	return secrets, nil
}

func (osq *ObjectStoreQueryer) PersistentVolumeClaimsForPod(ctx context.Context, pod *corev1.Pod) ([]*corev1.PersistentVolumeClaim, error) {
	if pod == nil {
		return nil, errors.New("pod is nil")
	}

	var persistentVolumeClaims []*corev1.PersistentVolumeClaim
	key := store.Key{
		Namespace:  pod.Namespace,
		APIVersion: "v1",
		Kind:       "PersistentVolumeClaim",
	}
	ul, _, err := osq.objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving persistentVolumeClaims")
	}

	for i := range ul.Items {
		pvc := &corev1.PersistentVolumeClaim{}
		err := kubernetes.FromUnstructured(&ul.Items[i], pvc)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured persistentVolumeClaim")
		}

		for _, v := range pod.Spec.Volumes {
			if v.PersistentVolumeClaim != nil && v.PersistentVolumeClaim.ClaimName == pvc.Name {
				persistentVolumeClaims = append(persistentVolumeClaims, pvc)
			}
		}
	}

	return persistentVolumeClaims, nil
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
	default:
		return nil, errors.Errorf("unable to retrieve selector for type %T", object)
	}
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

func containsBackend(lst []networkingv1.IngressBackend, s string) bool {
	for _, item := range lst {
		if item.Service.Name == s {
			return true
		}
	}
	return false
}
