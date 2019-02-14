package overview

import (
	"path"
	"sort"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

// TODO: these are here to support all the items that have been deprecated by the printer move

func getSelector(object runtime.Object) (*metav1.LabelSelector, error) {
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

func listPods(namespace string, selector *metav1.LabelSelector, uid types.UID, c cache.Cache) ([]*corev1.Pod, error) {
	key := cache.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	pods, err := loadPods(key, c, selector)
	if err != nil {
		return nil, err
	}

	var owned []*corev1.Pod
	for _, pod := range pods {
		controllerRef := metav1.GetControllerOf(pod)
		if controllerRef == nil || controllerRef.UID != uid {
			continue
		}

		owned = append(owned, pod)
	}

	return owned, nil
}

func loadPods(key cache.Key, c cache.Cache, selector *metav1.LabelSelector) ([]*corev1.Pod, error) {
	objects, err := c.List(key)
	if err != nil {
		return nil, err
	}

	var list []*corev1.Pod

	for _, object := range objects {
		pod := &corev1.Pod{}
		if err := scheme.Scheme.Convert(object, pod, runtime.InternalGroupVersioner); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(pod, object); err != nil {
			return nil, err
		}

		podSelector := &metav1.LabelSelector{
			MatchLabels: pod.GetLabels(),
		}

		if selector == nil || isEqualSelector(selector, podSelector) {
			list = append(list, pod)
		}
	}

	return list, nil
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

func gvkPath(apiVersion, kind, name string) string {
	var p string

	switch {
	case apiVersion == "apps/v1" && kind == "DaemonSet":
		p = "/content/overview/workloads/daemon-sets"
	case apiVersion == "extensions/v1beta1" && kind == "ReplicaSet":
		p = "/content/overview/workloads/replica-sets"
	case apiVersion == "apps/v1" && kind == "StatefulSet":
		p = "/content/overview/workloads/stateful-sets"
	case apiVersion == "extensions/v1beta1" && kind == "Deployment":
		p = "/content/overview/workloads/deployments"
	case apiVersion == "apps/v1" && kind == "Deployment":
		p = "/content/overview/workloads/deployments"
	case apiVersion == "batch/v1beta1" && kind == "CronJob":
		p = "/content/overview/workloads/cron-jobs"
	case (apiVersion == "batch/v1beta1" || apiVersion == "batch/v1") && kind == "Job":
		p = "/content/overview/workloads/jobs"
	case apiVersion == "v1" && kind == "ReplicationController":
		p = "/content/overview/workloads/replication-controllers"
	case apiVersion == "v1" && kind == "Secret":
		p = "/content/overview/config-and-storage/secrets"
	case apiVersion == "v1" && kind == "ConfigMap":
		p = "/content/overview/config-and-storage/configmaps"
	case apiVersion == "v1" && kind == "PersistentVolumeClaim":
		p = "/content/overview/config-and-storage/persistent-volume-claims"
	case apiVersion == "v1" && kind == "ServiceAccount":
		p = "/content/overview/config-and-storage/service-accounts"
	case apiVersion == "v1" && kind == "Service":
		p = "/content/overview/discovery-and-load-balancing/services"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "Role":
		p = "/content/overview/rbac/roles"
	case apiVersion == "v1" && kind == "Event":
		p = "/content/overview/events"
	default:
		return "/content/overview"
	}

	return path.Join(p, name)
}

func listReplicaSets(deployment *appsv1.Deployment, c cache.Cache) ([]*appsv1.ReplicaSet, error) {
	key := cache.Key{
		Namespace:  deployment.GetNamespace(),
		APIVersion: deployment.APIVersion,
		Kind:       "ReplicaSet",
	}

	replicaSets, err := loadReplicaSets(key, c, deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	var owned []*appsv1.ReplicaSet
	for _, rs := range replicaSets {
		if metav1.IsControlledBy(rs, deployment) {
			owned = append(owned, rs)
		}
	}

	return owned, nil
}

func loadReplicaSets(key cache.Key, c cache.Cache, selector *metav1.LabelSelector) ([]*appsv1.ReplicaSet, error) {
	objects, err := c.List(key)
	if err != nil {
		return nil, err
	}

	var list []*appsv1.ReplicaSet

	for _, object := range objects {
		rs := &appsv1.ReplicaSet{}
		if err := scheme.Scheme.Convert(object, rs, 0); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(rs, object); err != nil {
			return nil, err
		}

		if selector == nil {
			list = append(list, rs)
		} else if isEqualSelector(selector, rs.Spec.Selector) {
			list = append(list, rs)
		}
	}

	return list, nil
}

func findNewReplicaSet(deployment *appsv1.Deployment, rsList []*appsv1.ReplicaSet) *appsv1.ReplicaSet {
	sort.Sort(replicaSetsByCreationTimestamp(rsList))
	for i := range rsList {
		if equalIgnoreHash(&rsList[i].Spec.Template, &deployment.Spec.Template) {
			// In rare cases, such as after cluster upgrades, Deployment may end up with
			// having more than one new ReplicaSets that have the same template as its template,
			// see https://github.com/kubernetes/kubernetes/issues/40415
			// We deterministically choose the oldest new ReplicaSet.
			return rsList[i]
		}
	}

	// new ReplicaSet does not exist.
	return nil
}

// findOldReplicaSets returns the old replica sets targeted by the given Deployment, with the given slice of replica sets.
// Note that the first set of old replica sets doesn't include the ones with no pods, and the second set of old replica sets include all old replica sets.
func findOldReplicaSets(deployment *appsv1.Deployment, rsList []*appsv1.ReplicaSet) []*appsv1.ReplicaSet {
	var requiredRSs []*appsv1.ReplicaSet
	newRS := findNewReplicaSet(deployment, rsList)
	for _, rs := range rsList {
		// Filter out new replica set
		if newRS != nil && rs.UID == newRS.UID {
			continue
		}
		if rs.Spec.Replicas != nil && *rs.Spec.Replicas != 0 {
			requiredRSs = append(requiredRSs, rs)
		}
	}
	return requiredRSs
}

// replicaSetsByCreationTimestamp sorts a list of ReplicaSet by creation timestamp, using their names as a tie breaker.
type replicaSetsByCreationTimestamp []*appsv1.ReplicaSet

func (o replicaSetsByCreationTimestamp) Len() int      { return len(o) }
func (o replicaSetsByCreationTimestamp) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o replicaSetsByCreationTimestamp) Less(i, j int) bool {
	if o[i].CreationTimestamp.Equal(&o[j].CreationTimestamp) {
		return o[i].Name < o[j].Name
	}
	return o[i].CreationTimestamp.Before(&o[j].CreationTimestamp)
}

// EqualIgnoreHash returns true if two given podTemplateSpec are equal, ignoring the diff in value of Labels[pod-template-hash]
// We ignore pod-template-hash because:
// 1. The hash result would be different upon podTemplateSpec API changes
//    (e.g. the addition of a new field will cause the hash code to change)
// 2. The deployment template won't have hash labels
func equalIgnoreHash(template1, template2 *corev1.PodTemplateSpec) bool {
	t1Copy := template1.DeepCopy()
	t2Copy := template2.DeepCopy()

	// Remove hash labels from template.Labels before comparing
	for _, key := range extraKeys {
		delete(t1Copy.Labels, key)
		delete(t2Copy.Labels, key)
	}

	return apiequality.Semantic.DeepEqual(t1Copy, t2Copy)
}

func listSecrets(namespace string, c cache.Cache) ([]*corev1.Secret, error) {
	key := cache.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Secret",
	}

	return loadSecrets(key, c)
}

func loadSecrets(key cache.Key, c cache.Cache) ([]*corev1.Secret, error) {
	objects, err := c.List(key)
	if err != nil {
		return nil, err
	}

	var list []*corev1.Secret

	for _, object := range objects {
		e := &corev1.Secret{}
		if err := scheme.Scheme.Convert(object, e, runtime.InternalGroupVersioner); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(e, object); err != nil {
			return nil, err
		}

		list = append(list, e)
	}

	return list, nil
}

// loadSecret loads a single secret from the cache.
// Note if the secret is not found, a nil error and secret is returned,
// i.e. it is not an error.
func loadSecret(key cache.Key, c cache.Cache) (*corev1.Secret, error) {
	object, err := c.Get(key)
	if err != nil {
		return nil, err
	}

	// Object not found is not an error
	if object == nil {
		return nil, nil
	}

	secret := &corev1.Secret{}
	if err := scheme.Scheme.Convert(object, secret, runtime.InternalGroupVersioner); err != nil {
		return nil, err
	}

	if err := copyObjectMeta(secret, object); err != nil {
		return nil, err
	}

	return secret, nil
}
