package overview

import (
	"context"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

type PodList struct{}

func NewPodList(prefix, namespace string, c clock.Clock) view.View {
	return &PodList{}
}

func (pc *PodList) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	mobject, ok := object.(metav1.Object)
	if !ok {
		return nil, errors.Errorf("%T is not an object", object)
	}

	selector, err := getSelector(object)
	if err != nil {
		return nil, err
	}

	pods, err := listPods(mobject.GetNamespace(), selector, mobject.GetUID(), c)
	if err != nil {
		return nil, err
	}

	list := &corev1.PodList{}
	for _, pod := range pods {
		list.Items = append(list.Items, *pod)
	}

	var contents []content.Content

	err = printContentObject(
		"Pods",
		"ns",
		"prefix",
		"No pods were found",
		podTransforms,
		list,
		&contents,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to print pods")
	}

	return contents, nil
}

type PodSummary struct{}

var _ view.View = (*PodSummary)(nil)

func NewPodSummary(prefix, namespace string, c clock.Clock) view.View {
	return &PodSummary{}
}

func (ps *PodSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	// TODO this clock should come from somewhere else
	clk := &clock.RealClock{}

	pod, err := retrievePod(object)
	if err != nil {
		return nil, err
	}

	detail, err := printPodSummary(pod, clk)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type PodCondition struct{}

func NewPodCondition(prefix, namespace string, c clock.Clock) view.View {
	return &PodCondition{}
}

func (pc *PodCondition) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	pod, err := retrievePod(object)
	if err != nil {
		return nil, err
	}

	conditions := pod.Status.Conditions

	table := content.NewTable("Conditions", "No conditions")
	table.Columns = []content.TableColumn{
		view.TableCol("Type"),
		view.TableCol("Status"),
		view.TableCol("Last probe time"),
		view.TableCol("Last transition time"),
		view.TableCol("Reason"),
		view.TableCol("Message"),
	}

	for _, condition := range conditions {

		lastProbeTime := condition.LastProbeTime.UTC().Format(time.RFC3339)
		lastTransitionTime := condition.LastTransitionTime.UTC().Format(time.RFC3339)

		row := content.TableRow{
			"Type":                 content.NewStringText(string(condition.Type)),
			"Status":               content.NewStringText(string(condition.Status)),
			"Last probe time":      content.NewTimeText(lastProbeTime),
			"Last transition time": content.NewTimeText(lastTransitionTime),
			"Reason":               content.NewStringText(condition.Reason),
			"Message":              content.NewStringText(condition.Message),
		}

		table.AddRow(row)
	}

	return []content.Content{&table}, nil
}

type PodContainer struct{}

func NewPodContainer(prefix, namespace string, c clock.Clock) view.View {
	return &PodContainer{}
}

func (pc *PodContainer) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	pod, err := retrievePod(object)
	if err != nil {
		return nil, err
	}

	containers := pod.Spec.Containers
	statuses := pod.Status.ContainerStatuses

	return describePodContainers(containers, statuses)
}

type PodVolume struct{}

func NewPodVolume(prefix, namespace string, c clock.Clock) view.View {
	return &PodVolume{}
}

func (pc *PodVolume) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	pod, err := retrievePod(object)
	if err != nil {
		return nil, err
	}

	sections := []content.Section{}

	for _, volume := range pod.Spec.Volumes {
		sections = append(sections, summarizeVolume(volume))
	}

	volumes := content.NewSummary("Volumes", sections)

	return []content.Content{&volumes}, nil
}

func retrievePod(object runtime.Object) (*corev1.Pod, error) {
	pod, ok := object.(*corev1.Pod)
	if !ok {
		return nil, errors.Errorf("expected object to be a Pod, it was %T", object)
	}

	return pod, nil
}

type podStatus struct {
	Running   int
	Waiting   int
	Succeeded int
	Failed    int
}

func createPodStatus(pods []*corev1.Pod) podStatus {
	var ps podStatus

	for _, pod := range pods {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			ps.Running++
		case corev1.PodPending:
			ps.Waiting++
		case corev1.PodSucceeded:
			ps.Succeeded++
		case corev1.PodFailed:
			ps.Failed++
		}
	}

	return ps
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
	objects, err := c.Retrieve(key)
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
