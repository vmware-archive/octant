package overview

import (
	"context"
	"reflect"
	"sort"

	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

type DeploymentSummary struct{}

var _ View = (*DeploymentSummary)(nil)

func NewDeploymentSummary(prefix, namespace string, c clock.Clock) View {
	return &DeploymentSummary{}
}

func (ds *DeploymentSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	deployment, err := retrieveDeployment(object)
	if err != nil {
		return nil, err
	}

	return ds.summary(deployment)
}

func (ds *DeploymentSummary) summary(deployment *appsv1.Deployment) ([]content.Content, error) {
	section, err := printDeploymentSummary(deployment)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{section})
	contents := []content.Content{
		&summary,
	}

	return contents, nil
}

type DeploymentReplicaSets struct{}

var _ View = (*DeploymentReplicaSets)(nil)

func NewDeploymentReplicaSets(prefix, namespace string, c clock.Clock) View {
	return &DeploymentReplicaSets{}
}

func (drs *DeploymentReplicaSets) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	var contents []content.Content

	deployment, err := retrieveDeployment(object)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving deployment")
	}

	replicaSetContent, err := drs.replicaSets(deployment, c)
	if err != nil {
		return nil, errors.Wrap(err, "rendering replicasets")
	}
	contents = append(contents, replicaSetContent...)

	return contents, nil
}

func (drs *DeploymentReplicaSets) replicaSets(deployment *appsv1.Deployment, c cache.Cache) ([]content.Content, error) {
	contents := []content.Content{}

	replicaSets, err := listReplicaSets(deployment, c)
	if err != nil {
		return nil, err
	}

	newReplicaSet := findNewReplicaSet(deployment, replicaSets)

	if newReplicaSet != nil {
		err = printContentObject(
			"New Replica Set",
			"",
			"",
			"This Deployment does not have a current Replica",
			replicaSetTransforms,
			newReplicaSet,
			&contents,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to print new replica set")
		}
	}

	oldList := &appsv1.ReplicaSetList{}
	for _, rs := range findOldReplicaSets(deployment, replicaSets) {
		oldList.Items = append(oldList.Items, *rs)
	}

	err = printContentObject(
		"Old Replica Sets",
		"",
		"",
		"This Deployment does not have any old Replicas",
		replicaSetTransforms,
		oldList,
		&contents,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to print old replica sets")
	}

	return contents, nil
}

func printContentObject(title, namespace, prefix, emptyMessage string, transforms map[string]lookupFunc, object runtime.Object, contents *[]content.Content) error {
	if reflect.ValueOf(object).IsNil() {
		return errors.New("unable to print a nil object")
	}

	otf := summaryFunc(title, emptyMessage, transforms)
	transformed := otf(namespace, prefix, contents)
	return printObject(object, transformed)
}

func retrieveDeployment(object runtime.Object) (*appsv1.Deployment, error) {
	deployment, ok := object.(*appsv1.Deployment)
	if !ok {
		return nil, errors.Errorf("expected object to be a Deployment, it was %T", object)
	}

	return deployment, nil
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
	objects, err := c.Retrieve(key)
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

// findOldReplicaSets returns the old replica sets targeted by the given Deployment, with the given slice of RSes.
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
