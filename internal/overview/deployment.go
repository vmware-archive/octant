package overview

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

type DeploymentSummary struct{}

var _ View = (*DeploymentSummary)(nil)

func NewDeploymentSummary() *DeploymentSummary {
	return &DeploymentSummary{}
}

func (ds *DeploymentSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	deployment, err := retrieveDeployment(object)
	if err != nil {
		return nil, err
	}

	return ds.summary(deployment)
}

func (ds *DeploymentSummary) summary(deployment *extensions.Deployment) ([]content.Content, error) {
	section, err := printDeploymentSummary(deployment)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{section})
	return []content.Content{
		&summary,
	}, nil
}

type DeploymentReplicaSets struct{}

var _ View = (*DeploymentReplicaSets)(nil)

func NewDeploymentReplicaSets() *DeploymentReplicaSets {
	return &DeploymentReplicaSets{}
}

func (drs *DeploymentReplicaSets) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	var contents []content.Content

	deployment, err := retrieveDeployment(object)
	if err != nil {
		log.Printf("wtf: %v", err)
		return nil, err
	}

	replicaSetContent, err := drs.replicaSets(deployment, c)
	if err != nil {
		log.Printf("wtf2: %v", err)
		return nil, err
	}
	contents = append(contents, replicaSetContent...)

	return contents, nil
}

func (drs *DeploymentReplicaSets) replicaSets(deployment *extensions.Deployment, c Cache) ([]content.Content, error) {
	contents := []content.Content{}

	replicaSets, err := listReplicaSets(deployment, c)
	if err != nil {
		return nil, err
	}

	newReplicaSet := findNewReplicaSet(deployment, replicaSets)

	err = printContentObject(
		"New Replica Set",
		"ns",
		"prefix",
		replicaSetTransforms,
		newReplicaSet,
		&contents,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to print new replica set")
	}

	oldList := &extensions.ReplicaSetList{}
	for _, rs := range findOldReplicaSets(deployment, replicaSets) {
		oldList.Items = append(oldList.Items, *rs)
	}

	err = printContentObject(
		"Old Replica Sets",
		"ns",
		"prefix",
		replicaSetTransforms,
		oldList,
		&contents,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to print old replica sets")
	}

	return contents, nil
}

func printContentObject(title, namespace, prefix string, transforms map[string]lookupFunc, object runtime.Object, contents *[]content.Content) error {
	if reflect.ValueOf(object).IsNil() {
		return errors.New("unable to print a nil object")
	}

	otf := summaryFunc(title, transforms)
	transformed := otf(namespace, prefix, contents)
	return printObject(object, transformed)
}

func retrieveDeployment(object runtime.Object) (*extensions.Deployment, error) {
	deployment, ok := object.(*extensions.Deployment)
	if !ok {
		return nil, errors.Errorf("expected object to be a Deployment, it was %T", object)
	}

	return deployment, nil
}

func listReplicaSets(deployment *extensions.Deployment, c Cache) ([]*extensions.ReplicaSet, error) {
	key := CacheKey{
		Namespace:  deployment.GetNamespace(),
		APIVersion: deployment.APIVersion,
		Kind:       "ReplicaSet",
	}

	replicaSets, err := loadReplicaSets(key, c, deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	var owned []*extensions.ReplicaSet
	for _, rs := range replicaSets {
		if metav1.IsControlledBy(rs, deployment) {
			owned = append(owned, rs)
		}
	}

	return owned, nil
}

func loadReplicaSets(key CacheKey, c Cache, selector *metav1.LabelSelector) ([]*extensions.ReplicaSet, error) {
	objects, err := c.Retrieve(key)
	if err != nil {
		return nil, err
	}

	var list []*extensions.ReplicaSet

	for _, object := range objects {
		rs := &extensions.ReplicaSet{}
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

func isEqualSelector(s1, s2 *metav1.LabelSelector) bool {
	s1Copy := s1.DeepCopy()
	s2Copy := s2.DeepCopy()

	delete(s1Copy.MatchLabels, extensions.DefaultDeploymentUniqueLabelKey)
	delete(s2Copy.MatchLabels, extensions.DefaultDeploymentUniqueLabelKey)

	return apiequality.Semantic.DeepEqual(s1Copy, s2Copy)
}

func equalIgnoreHash(template1, template2 *core.PodTemplateSpec) bool {
	t1Copy := template1.DeepCopy()
	t2Copy := template2.DeepCopy()
	// Remove hash labels from template.Labels before comparing
	delete(t1Copy.Labels, extensions.DefaultDeploymentUniqueLabelKey)
	delete(t2Copy.Labels, extensions.DefaultDeploymentUniqueLabelKey)

	return apiequality.Semantic.DeepEqual(*t1Copy, *t2Copy)
}

func findNewReplicaSet(deployment *extensions.Deployment, rsList []*extensions.ReplicaSet) *extensions.ReplicaSet {
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
func findOldReplicaSets(deployment *extensions.Deployment, rsList []*extensions.ReplicaSet) []*extensions.ReplicaSet {
	var requiredRSs []*extensions.ReplicaSet
	newRS := findNewReplicaSet(deployment, rsList)
	for _, rs := range rsList {
		// Filter out new replica set
		if newRS != nil && rs.UID == newRS.UID {
			continue
		}
		if rs.Spec.Replicas != 0 {
			requiredRSs = append(requiredRSs, rs)
		}
	}
	return requiredRSs
}

// replicaSetsByCreationTimestamp sorts a list of ReplicaSet by creation timestamp, using their names as a tie breaker.
type replicaSetsByCreationTimestamp []*extensions.ReplicaSet

func (o replicaSetsByCreationTimestamp) Len() int      { return len(o) }
func (o replicaSetsByCreationTimestamp) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o replicaSetsByCreationTimestamp) Less(i, j int) bool {
	if o[i].CreationTimestamp.Equal(&o[j].CreationTimestamp) {
		return o[i].Name < o[j].Name
	}
	return o[i].CreationTimestamp.Before(&o[j].CreationTimestamp)
}

func printDeploymentSummary(deployment *extensions.Deployment) (content.Section, error) {
	minReadySeconds := fmt.Sprintf("%d", deployment.Spec.MinReadySeconds)

	var revisionHistoryLimit string
	if rhl := deployment.Spec.RevisionHistoryLimit; rhl != nil {
		revisionHistoryLimit = fmt.Sprintf("%d", *rhl)
	}

	var rollingUpdateStrategy string
	if rus := deployment.Spec.Strategy.RollingUpdate; rus != nil {
		rollingUpdateStrategy = fmt.Sprintf("Max Surge: %s, Max unavailable: %s",
			rus.MaxSurge.String(), rus.MaxUnavailable.String())
	}

	status := fmt.Sprintf("%d updated, %d total, %d available, %d unavailable",
		deployment.Status.UpdatedReplicas,
		deployment.Status.Replicas,
		deployment.Status.AvailableReplicas,
		deployment.Status.UnavailableReplicas,
	)

	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}

	section := content.Section{
		Items: []content.Item{
			content.TextItem("Name", deployment.GetName()),
			content.TextItem("Namespace", deployment.GetNamespace()),
			content.LabelsItem("Labels", deployment.GetLabels()),
			content.LabelsItem("Annotations", deployment.GetAnnotations()),
			content.TextItem("Creation Time", deployment.CreationTimestamp.Time.UTC().Format(time.RFC1123Z)),
			content.TextItem("Selector", selector.String()),
			content.TextItem("Strategy", string(deployment.Spec.Strategy.Type)),
			content.TextItem("Min Ready Seconds", minReadySeconds),
			content.TextItem("Revision History Limit", revisionHistoryLimit),
			content.TextItem("Rolling Update Strategy", rollingUpdateStrategy),
			content.TextItem("Status", status),
		},
	}

	return section, nil
}
