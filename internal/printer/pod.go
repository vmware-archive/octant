/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiEquality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	podColsWithLabels    = component.NewTableCols("Name", "Labels", "Ready", "Phase", "Status", "Restarts", "Node", "Age")
	podColsWithOutLabels = component.NewTableCols("Name", "Ready", "Phase", "Status", "Restarts", "Node", "Age")
	podResourceCols      = component.NewTableCols("Container", "Request: Memory", "Request: CPU", "Limit: Memory", "Limit: CPU")
)

// PodListHandler is a printFunc that prints pods
func PodListHandler(ctx context.Context, list *corev1.PodList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	cols := podColsWithLabels
	if opts.DisableLabels {
		cols = podColsWithOutLabels
	}

	ot := NewObjectTable("Pods", "We couldn't find any pods!", cols, opts.DashConfig.ObjectStore())
	ot.AddFilters(podTableFilters())
	ot.EnablePluginStatus(opts.DashConfig.PluginManager())

	for i := range list.Items {
		row := component.TableRow{}
		pod := list.Items[i]
		nameLink, err := opts.Link.ForObject(&pod, pod.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink

		if !opts.DisableLabels {
			row["Labels"] = component.NewLabels(pod.Labels)
		}

		readyCounter := 0
		for _, c := range pod.Status.ContainerStatuses {
			if c.Ready {
				readyCounter++
			}
		}
		ready := fmt.Sprintf("%d/%d", readyCounter, len(pod.Spec.Containers))
		row["Ready"] = component.NewText(ready)

		row["Phase"] = component.NewText(string(pod.Status.Phase))

		if len(pod.Status.ContainerStatuses) > 0 {
			lastStatus := pod.Status.ContainerStatuses[len(pod.Status.ContainerStatuses)-1]

			if state := lastStatus.State.Waiting; state != nil {
				row["Status"] = component.NewText(lastStatus.State.Waiting.Reason)
			} else if state := lastStatus.State.Running; state != nil {
				row["Status"] = component.NewText("Running")
			}

			if pod.DeletionTimestamp != nil && lastStatus.State.Terminated == nil {
				row["Status"] = component.NewText("Terminating")
			}
		} else {
			row["Status"] = component.NewText("")
		}

		restartCounter := 0
		for _, c := range pod.Status.ContainerStatuses {
			restartCounter += int(c.RestartCount)
		}
		restarts := fmt.Sprintf("%d", restartCounter)
		row["Restarts"] = component.NewText(restarts)

		nodeComponent, err := podNode(&pod, opts.Link)
		if err != nil {
			return nil, err
		}

		row["Node"] = nodeComponent

		ts := pod.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		if err := ot.AddRowForObject(ctx, &pod, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	ot.SetSortOrder("Name", false)

	return ot.ToComponent()
}

func podNode(pod *corev1.Pod, linkGenerator link.Interface) (component.Component, error) {
	if nodeName := pod.Spec.NodeName; nodeName != "" {
		return linkGenerator.ForGVK("", "v1", "Node", pod.Spec.NodeName, pod.Spec.NodeName)
	}

	return component.NewText("<not scheduled>"), nil
}

var podConditionColumns = [][]string{
	{"Type", "type"},
	{"Reason", "reason"},
	{"Status", "status"},
	{"Message", "message"},
	{"Last Probe", "lastProbeTime"},
	{"Last Transition", "lastTransitionTime"},
}

// PodHandler is a printFunc that prints Pods
func PodHandler(ctx context.Context, pod *corev1.Pod, options Options) (component.Component, error) {
	o := NewObject(pod)
	o.ConditionsGen = conditionsGenFactory("", podConditionColumns, nil)
	o.EnableEvents()

	ph, err := newPodHandler(pod, o)
	if err != nil {
		return nil, err
	}

	if err := ph.Config(options); err != nil {
		return nil, errors.Wrap(err, "print pod configuration")
	}
	if err := ph.Status(options); err != nil {
		return nil, errors.Wrap(err, "print pod status")
	}
	if err := ph.InitContainers(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print pod init containers")
	}
	if err := ph.Containers(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print pod containers")
	}
	if err := ph.EphemeralContainers(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print pod ephemeral containers")
	}
	if err := ph.Additional(options); err != nil {
		return nil, errors.Wrap(err, "print pod additional items")
	}

	return o.ToComponent(ctx, options)
}

func createPodSummaryStatus(pod *corev1.Pod) (*component.Summary, error) {
	if pod == nil {
		return nil, errors.New("pod is nil")
	}

	summary := component.NewSummary("Status")

	sections := component.SummarySections{}

	sections.AddText("QoS", string(pod.Status.QOSClass))

	if pod.DeletionTimestamp != nil {
		summary.SetAlert(component.NewAlert(component.AlertStatusError, component.AlertTypeDefault, "Pod is being deleted", false, nil))

		sections = append(sections, component.SummarySection{
			Header:  "Status: Terminating",
			Content: component.NewText(pod.DeletionTimestamp.String()),
		})
		if pod.DeletionGracePeriodSeconds != nil {
			sections.AddText("Termination Grace Period", fmt.Sprintf("%ds", *pod.DeletionGracePeriodSeconds))
		}
	} else {
		sections.AddText("Phase", string(pod.Status.Phase))
	}

	if pod.Status.Reason != "" {
		sections.AddText("Reason", pod.Status.Reason)
	}
	if pod.Status.Message != "" {
		sections.AddText("Message", pod.Status.Message)
	}

	sections.AddText("Pod IP", pod.Status.PodIP)
	sections.AddText("Host IP", pod.Status.HostIP)

	if pod.Status.NominatedNodeName != "" {
		sections.AddText("NominatedNodeName", pod.Status.NominatedNodeName)
	}

	summary.Add(sections...)

	return summary, nil
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
func (p *PodConfiguration) Create(options Options) (*component.Summary, error) {
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

	contentLink, err := options.Link.ForGVK(pod.Namespace, "v1", "ServiceAccount", pod.Spec.ServiceAccountName, pod.Spec.ServiceAccountName)
	if err != nil {
		return nil, err
	}

	nodeLink, err := podNode(p.pod, options.Link)
	if err != nil {
		return nil, err
	}
	sections.Add("Node", nodeLink)

	sections = append(sections, component.SummarySection{
		Header:  "Service Account",
		Content: contentLink,
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func listPods(ctx context.Context, namespace string, selector *metav1.LabelSelector, uid types.UID, o store.Store) ([]*corev1.Pod, error) {
	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	pods, err := loadPods(ctx, key, o, selector)
	if err != nil {
		return nil, errors.Wrap(err, "load pods")
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

func loadPods(ctx context.Context, key store.Key, o store.Store, labelSelector *metav1.LabelSelector) ([]*corev1.Pod, error) {
	objects, _, err := o.List(ctx, key)
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

		if selector == kLabels.Nothing() || isEqualSelector(labelSelector, podSelector) || selector.Matches(kLabels.Set(pod.Labels)) {
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
	if s1 == nil || s2 == nil {
		return false
	}

	s1Copy := s1.DeepCopy()
	s2Copy := s2.DeepCopy()

	for _, key := range extraKeys {
		delete(s1Copy.MatchLabels, key)
		delete(s2Copy.MatchLabels, key)
	}

	return apiEquality.Semantic.DeepEqual(s1Copy, s2Copy)
}

func createPodListView(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	options.DisableLabels = true

	podList := &corev1.PodList{}

	objectStore := options.DashConfig.ObjectStore()

	accessor := meta.NewAccessor()

	namespace, err := accessor.Namespace(object)
	if err != nil {
		return nil, errors.Wrap(err, "get namespace for object")
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return nil, errors.Wrap(err, "Get apiVersion for object")
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return nil, errors.Wrap(err, "get kind for object")
	}

	name, err := accessor.Name(object)
	if err != nil {
		return nil, errors.Wrap(err, "get name for object")
	}

	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	list, _, err := objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "list all objects for key %+v", key)
	}

	for i := range list.Items {
		pod := &corev1.Pod{}
		err := kubernetes.FromUnstructured(&list.Items[i], pod)
		if err != nil {
			return nil, err
		}

		for _, ownerReference := range pod.OwnerReferences {
			if ownerReference.APIVersion == apiVersion &&
				ownerReference.Kind == kind &&
				ownerReference.Name == name {
				podList.Items = append(podList.Items, *pod)
			}
		}
	}

	return PodListHandler(ctx, podList, options)
}

func createRollingPodListView(ctx context.Context, objects []runtime.Object, options Options) (component.Component, error) {
	options.DisableLabels = true

	podList := &corev1.PodList{}

	objectStore := options.DashConfig.ObjectStore()

	accessor := meta.NewAccessor()

	for _, object := range objects {
		namespace, err := accessor.Namespace(object)
		if err != nil {
			return nil, errors.Wrap(err, "get namespace for object")
		}

		apiVersion, err := accessor.APIVersion(object)
		if err != nil {
			return nil, errors.Wrap(err, "Get apiVersion for object")
		}

		kind, err := accessor.Kind(object)
		if err != nil {
			return nil, errors.Wrap(err, "get kind for object")
		}

		name, err := accessor.Name(object)
		if err != nil {
			return nil, errors.Wrap(err, "get name for object")
		}

		key := store.Key{
			Namespace:  namespace,
			APIVersion: "v1",
			Kind:       "Pod",
		}

		list, _, err := objectStore.List(ctx, key)
		if err != nil {
			return nil, errors.Wrapf(err, "list all objects for key %+v", key)
		}

		for i := range list.Items {
			pod := &corev1.Pod{}
			err := kubernetes.FromUnstructured(&list.Items[i], pod)
			if err != nil {
				return nil, err
			}

			for _, ownerReference := range pod.OwnerReferences {
				if ownerReference.APIVersion == apiVersion &&
					ownerReference.Kind == kind &&
					ownerReference.Name == name {
					podList.Items = append(podList.Items, *pod)
				}
			}
		}
	}

	return PodListHandler(ctx, podList, options)
}

func createMountedPodListView(ctx context.Context, namespace string, persistentVolumeClaimName string, options Options) (component.Component, error) {
	options.DisableLabels = true

	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	objectStore := options.DashConfig.ObjectStore()

	mountedPodList := &corev1.PodList{}

	pods, err := loadPods(ctx, key, objectStore, nil)
	if err != nil {
		return nil, err
	}

	for _, pod := range pods {
		var volumeClaims []corev1.Volume

		for _, volume := range pod.Spec.Volumes {
			if volume.VolumeSource.PersistentVolumeClaim != nil {
				volumeClaims = append(volumeClaims, volume)
			}
		}

		for _, persistentVolumeClaim := range volumeClaims {
			if persistentVolumeClaim.PersistentVolumeClaim.ClaimName == persistentVolumeClaimName {
				mountedPodList.Items = append(mountedPodList.Items, *pod)
			}
		}
	}

	return PodListHandler(ctx, mountedPodList, options)
}

func printPodResources(podSpec corev1.PodSpec) (*component.Table, error) {
	table := component.NewTable("Resources", "Pod has no resource needs", podResourceCols)

	// for each container in the spec, there will be requests and limits
	// for memory and cpu

	for _, container := range podSpec.Containers {
		memoryRequest := ""
		if q := container.Resources.Requests.Memory(); q != nil {
			memoryRequest = q.String()
		}
		cpuRequest := ""
		if q := container.Resources.Requests.Cpu(); q != nil {
			cpuRequest = q.String()
		}
		memoryLimit := ""
		if q := container.Resources.Limits.Memory(); q != nil {
			memoryLimit = q.String()
		}
		cpuLimit := ""
		if q := container.Resources.Limits.Cpu(); q != nil {
			cpuLimit = q.String()
		}

		row := component.TableRow{
			"Container":       component.NewText(container.Name),
			"Request: Memory": component.NewText(memoryRequest),
			"Request: CPU":    component.NewText(cpuRequest),
			"Limit: Memory":   component.NewText(memoryLimit),
			"Limit: CPU":      component.NewText(cpuLimit),
		}
		table.Add(row)
	}

	return table, nil
}

type podObject interface {
	Config(options Options) error
	Status(options Options) error
	InitContainers(ctx context.Context, options Options) error
	Containers(ctx context.Context, options Options) error
	Additional(options Options) error
}

type podHandler struct {
	pod             *corev1.Pod
	configFunc      func(*corev1.Pod, Options) (*component.Summary, error)
	summaryFunc     func(*corev1.Pod, Options) (*component.Summary, error)
	containerFunc   func(ctx context.Context, pod *corev1.Pod, container *corev1.Container, isInit bool, options Options) (*component.Summary, error)
	additionalFuncs []func(*corev1.Pod, Options) ObjectPrinterFunc
	object          *Object
}

var _ podObject = (*podHandler)(nil)

var defaultPodHandlerAdditionalItems = []func(*corev1.Pod, Options) ObjectPrinterFunc{
	func(pod *corev1.Pod, options Options) ObjectPrinterFunc {
		return func() (component.Component, error) {
			return printPodResources(pod.Spec)
		}
	},
	func(pod *corev1.Pod, options Options) ObjectPrinterFunc {
		return func() (component.Component, error) {
			return printVolumes(pod.Spec.Volumes)
		}
	},
	func(pod *corev1.Pod, options Options) ObjectPrinterFunc {
		return func() (component.Component, error) {
			return printTolerations(pod.Spec)
		}
	},
	func(pod *corev1.Pod, options Options) ObjectPrinterFunc {
		return func() (component.Component, error) {
			return printAffinity(pod.Spec)
		}
	},
}

func newPodHandler(pod *corev1.Pod, object *Object) (*podHandler, error) {
	if pod == nil {
		return nil, errors.New("can't print a nil pod")
	}

	if object == nil {
		return nil, errors.New("can't print pod using a nil object printer")
	}

	ph := &podHandler{
		pod:             pod,
		configFunc:      defaultPodConfig,
		summaryFunc:     defaultPodSummary,
		containerFunc:   defaultPodContainers,
		additionalFuncs: defaultPodHandlerAdditionalItems,
		object:          object,
	}

	return ph, nil
}

func (p *podHandler) Config(options Options) error {
	out, err := p.configFunc(p.pod, options)
	if err != nil {
		return err
	}
	p.object.RegisterConfig(out)
	return nil
}

func defaultPodConfig(pod *corev1.Pod, options Options) (*component.Summary, error) {
	creator := NewPodConfiguration(pod)
	return creator.Create(options)
}

func (p *podHandler) Status(options Options) error {
	out, err := p.summaryFunc(p.pod, options)
	if err != nil {
		return err
	}

	p.object.RegisterSummary(out)
	return nil
}

func defaultPodSummary(pod *corev1.Pod, options Options) (*component.Summary, error) {
	return createPodSummaryStatus(pod)
}

func (p *podHandler) InitContainers(ctx context.Context, options Options) error {
	return p.containers(ctx, p.pod.Spec.InitContainers, true, options)
}

func (p *podHandler) EphemeralContainers(ctx context.Context, options Options) error {
	var containers []corev1.Container
	for _, container := range p.pod.Spec.EphemeralContainers {
		containers = append(containers, corev1.Container(container.EphemeralContainerCommon))
	}
	return p.containers(ctx, containers, false, options)
}

func (p *podHandler) containers(ctx context.Context, containers []corev1.Container, isInit bool, options Options) error {
	var itemDescriptors []ItemDescriptor

	for i := range containers {
		container := containers[i]

		itemDescriptors = append(itemDescriptors, ItemDescriptor{
			Width: component.WidthHalf,
			Func: func() (component.Component, error) {
				return p.containerFunc(ctx, p.pod, &container, isInit, options)
			},
		})
	}

	p.object.RegisterItems(itemDescriptors...)

	return nil
}

func (p *podHandler) Containers(ctx context.Context, options Options) error {
	return p.containers(ctx, p.pod.Spec.Containers, false, options)
}

func defaultPodContainers(ctx context.Context, pod *corev1.Pod, container *corev1.Container, isInit bool, options Options) (*component.Summary, error) {
	portForwarder := options.DashConfig.PortForwarder()
	creator := NewContainerConfiguration(ctx, pod, container, portForwarder, IsInit(isInit), WithPrintOptions(options))
	return creator.Create()
}

func (p *podHandler) Additional(options Options) error {
	var itemDescriptors []ItemDescriptor

	for i := range p.additionalFuncs {
		itemDescriptors = append(itemDescriptors, ItemDescriptor{
			Width: component.WidthHalf,
			Func:  p.additionalFuncs[i](p.pod, options),
		})
	}

	p.object.RegisterItems(itemDescriptors...)

	return nil
}

func podTableFilters() map[string]component.TableFilter {
	return map[string]component.TableFilter{
		"Phase": {
			Values:   []string{"Pending", "Running", "Succeeded", "Failed", "Unknown"},
			Selected: []string{"Pending", "Running"},
		},
	}
}

func addPodTableFilters(table *component.Table) {
	for k, v := range podTableFilters() {
		table.AddFilter(k, v)
	}
}
