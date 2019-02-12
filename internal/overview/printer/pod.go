package printer

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/view/flexlayout"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
)

// PodListHandler is a printFunc that prints pods
func PodListHandler(list *corev1.PodList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Ready", "Status", "Restarts", "Age")
	tbl := component.NewTable("Pods", cols)

	for _, p := range list.Items {
		row := component.TableRow{}
<<<<<<< HEAD
		podPath := gvkPath(p.TypeMeta.APIVersion, p.TypeMeta.Kind, p.Name)
		row["Name"] = component.NewLink("", p.Name, podPath)
		row["Labels"] = component.NewLabels(p.Labels)
=======
		podPath, err := gvkPathFromObject(&d)
		if err != nil {
			return nil, errors.Wrap(err, "get path for pod")
		}

		row["Name"] = component.NewLink("", d.Name, podPath)
		row["Labels"] = component.NewLabels(d.Labels)
>>>>>>> Create log viewer

		readyCounter := 0
		for _, c := range p.Status.ContainerStatuses {
			if c.Ready {
				readyCounter++
			}
		}
		ready := fmt.Sprintf("%d/%d", readyCounter, len(p.Spec.Containers))
		row["Ready"] = component.NewText(ready)

<<<<<<< HEAD
		status := fmt.Sprintf("%s", p.Status.Phase)
		row["Status"] = component.NewText(status)
=======
		row["Status"] = component.NewText(string(d.Status.Phase))
>>>>>>> Create log viewer

		restartCounter := 0
		for _, c := range p.Status.ContainerStatuses {
			restartCounter += int(c.RestartCount)
		}
		restarts := fmt.Sprintf("%d", restartCounter)
		row["Restarts"] = component.NewText(restarts)

		ts := p.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// PodHandler is a printFunc that prints Pods
func PodHandler(p *corev1.Pod, opts Options) (component.ViewComponent, error) {
	fl := flexlayout.New()

	configSection := fl.AddSection()
	podConfigGen := NewPodConfiguration(p)
	configView, err := podConfigGen.Create()
	if err != nil {
		return nil, err
	}

	if err := configSection.Add(configView, 16); err != nil {
		return nil, errors.Wrap(err, "add pod config to layout")
	}

	view := fl.ToComponent("Summary")
	return view, nil

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

// PodConfiguration generates pod configuration.
type PodConfiguration struct {
	pod *corev1.Pod
}

// NewPodConfiguration creates an instance of PodConfiguration.
func NewPodConfiguration(p *corev1.Pod) *PodConfiguration {
	return &PodConfiguration{
		pod: p,
	}
}

// Create creates a pod configuration summary.
func (p *PodConfiguration) Create() (*component.Summary, error) {
	if p.pod == nil {
		return nil, errors.New("pod is nil")
	}
	pod := p.pod

	sections := component.SummarySections{}

	if pod.Spec.Priority != nil {
		sections.AddText("Priority", fmt.Sprintf("%d", *pod.Spec.Priority))
	}
	if pod.Spec.PriorityClassName != "" {
		sections.AddText("PriorityClassName", pod.Spec.PriorityClassName)
	}

	if pod.Status.StartTime != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Start Time",
			Content: component.NewTimestamp(pod.Status.StartTime.Time),
		})
	}

	if pod.DeletionTimestamp != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Status: Terminating",
			Content: component.NewTimestamp(pod.DeletionTimestamp.Time),
		})
		if pod.DeletionGracePeriodSeconds != nil {
			sections.AddText("Termination Grace Period", fmt.Sprintf("%ds", *pod.DeletionGracePeriodSeconds))
		}
	} else {
		sections.AddText("Status", string(pod.Status.Phase))
	}

	if pod.Status.Reason != "" {
		sections.AddText("Reason", pod.Status.Reason)
	}
	if pod.Status.Message != "" {
		sections.AddText("Message", pod.Status.Message)
	}

	if controllerRef := metav1.GetControllerOf(pod); controllerRef != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Controlled By",
			Content: linkForOwner(pod, controllerRef),
		})
	}

	if pod.Status.NominatedNodeName != "" {
		sections.AddText("NominatedNodeName", pod.Status.NominatedNodeName)
	}

	if pod.Status.QOSClass != "" {
		sections.AddText("QoS Class", string(pod.Status.QOSClass))
	}

	sections = append(sections, component.SummarySection{
		Header: "Service Account",
		Content: linkForGVK(pod.Namespace, "v1", "ServiceAccount", pod.Spec.ServiceAccountName,
			pod.Spec.ServiceAccountName),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
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
	object.SetInitializers(from.GetInitializers())
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
