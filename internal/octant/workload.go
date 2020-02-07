/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"fmt"
	"path"
	"sort"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/objectstatus"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	// WorkloadStatusColorOK is the color for ok workload status.
	WorkloadStatusColorOK = "#60b515"
	// WorkloadStatusColorWarning is the color for warning workload status.
	WorkloadStatusColorWarning = "#f57600"
	// WorkloadStatusColorError is the color for error workload status.
	WorkloadStatusColorError = "#e12200"
)

// PodWithMetric combines a pod and resource list into a single type.
type PodWithMetric struct {
	Pod          *unstructured.Unstructured
	ResourceList corev1.ResourceList
}

// Workload is a workload.
type Workload struct {
	// IconName is the name of the icon for this workload.
	IconName string
	// Name is the name of the workload
	Name string
	// Owner is the ancestor that ultimately own the workload.
	Owner *unstructured.Unstructured

	SegmentCounter map[component.NodeStatus][]PodWithMetric

	podMetricsDisabled bool

	mu *sync.Mutex
}

// NewWorkload creates a workload.
func NewWorkload(name, iconName string) *Workload {
	w := &Workload{
		IconName:       iconName,
		Name:           name,
		SegmentCounter: make(map[component.NodeStatus][]PodWithMetric),
		mu:             &sync.Mutex{},
	}

	return w
}

func (w *Workload) DonutChart(size component.DonutChartSize) (*component.DonutChart, error) {
	chart := component.NewDonutChart()

	segments := w.donutSegments()
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].Status < segments[j].Status
	})

	chart.SetSegments(segments)
	chart.SetLabels("Pods", "Pod")
	chart.SetSize(size)

	return chart, nil
}

// PodsWithMetrics returns a slice of PodWithMetric.
func (w *Workload) PodsWithMetrics() []PodWithMetric {
	var list []PodWithMetric
	for _, podWithMetric := range w.SegmentCounter {
		list = append(list, podWithMetric...)
	}
	return list
}

func (w *Workload) Pods() *unstructured.UnstructuredList {
	list := &unstructured.UnstructuredList{}

	for _, pwm := range w.PodsWithMetrics() {
		list.Items = append(list.Items, *pwm.Pod)
	}

	return list
}

// AddStatus adds a pod status to the workload.
func (w *Workload) AddPodStatus(status component.NodeStatus, object *unstructured.Unstructured, resourceList corev1.ResourceList) {
	w.mu.Lock()
	defer w.mu.Unlock()

	podWithMetrics, ok := w.SegmentCounter[status]
	if !ok {
		podWithMetrics = []PodWithMetric{}
	}
	pwm := PodWithMetric{
		Pod:          object,
		ResourceList: resourceList,
	}
	podWithMetrics = append(podWithMetrics, pwm)
	w.SegmentCounter[status] = podWithMetrics
}

func (w *Workload) SetPodMetricsDisabled() {
	w.podMetricsDisabled = true
}

func (w *Workload) PodMetricsEnabled() bool {
	return !w.podMetricsDisabled
}

func (w *Workload) donutSegments() []component.DonutSegment {
	w.mu.Lock()
	defer w.mu.Unlock()

	var segments []component.DonutSegment

	for k, v := range w.SegmentCounter {
		segments = append(segments, component.DonutSegment{
			Count:  len(v),
			Status: k,
		})
	}

	return segments
}

func sortWorkloads(workloads []Workload) {
	sort.SliceStable(workloads, func(i, j int) bool {
		return workloads[i].Name < workloads[j].Name
	})
}

func printResourceUsage(resourceType corev1.ResourceName, quantity resource.Quantity) string {
	switch resourceType {
	case corev1.ResourceCPU:
		return fmt.Sprintf("%vm", quantity.MilliValue())
	case corev1.ResourceMemory:
		return fmt.Sprintf("%vMi", quantity.Value()/(1024*1024))
	default:
		return fmt.Sprintf("%v", quantity.Value())
	}
}

func objectOwner(ctx context.Context, objectStore store.Store, object *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if object == nil {
		return nil, fmt.Errorf("can't find owner for nil object")
	}

	ownerReferences := object.GetOwnerReferences()
	if len(ownerReferences) == 0 {
		return object, nil
	} else if len(ownerReferences) > 1 {
		return object, nil
	}

	ownerReference := ownerReferences[0]
	ownerRefKey := store.Key{
		Namespace:  object.GetNamespace(),
		APIVersion: ownerReference.APIVersion,
		Kind:       ownerReference.Kind,
		Name:       ownerReference.Name,
	}

	owner, err := objectStore.Get(ctx, ownerRefKey)
	if err != nil {
		return nil, fmt.Errorf("find object %s: %w", ownerRefKey, err)
	}

	if owner == nil {
		return object, nil
	}

	return objectOwner(ctx, objectStore, owner)
}

func parseAsResourceList(object map[string]interface{}, field string) (resource.Quantity, error) {
	s, found, err := unstructured.NestedString(object, field)
	if err != nil {
		return resource.Quantity{}, fmt.Errorf("parse %s: %w", field, err)
	}

	if !found {
		return resource.Quantity{}, fmt.Errorf("unable to find field '%s'", field)
	}

	return resource.ParseQuantity(s)
}

func addResourceList(a corev1.ResourceList, q resource.Quantity, resourceName corev1.ResourceName) corev1.ResourceList {
	rl := a.DeepCopy()

	orig := a[resourceName]
	orig.Add(q)
	rl[resourceName] = orig

	return rl
}

// WorkloadLoaderInterface loads workloads from a namespace.
type WorkloadLoaderInterface interface {
	// Load loads workloads from a namespace.
	Load(ctx context.Context, namespace string) ([]Workload, error)
}

// WorkloadCardCollector creates cards for workloads in a namespace.
type WorkloadCardCollector struct {
	WorkloadLoader WorkloadLoaderInterface
}

// NewWorkloadCardCollector creates an instance of WorkloadCardCollector.
func NewWorkloadCardCollector(loader WorkloadLoaderInterface) (*WorkloadCardCollector, error) {
	if loader == nil {
		return nil, fmt.Errorf("workload loader is nil")
	}

	return &WorkloadCardCollector{
		WorkloadLoader: loader,
	}, nil
}

// Collect collects cards.
func (wc *WorkloadCardCollector) Collect(ctx context.Context, namespace string) ([]*component.Card, bool, error) {
	workloads, err := wc.WorkloadLoader.Load(ctx, namespace)
	if err != nil {
		return nil, false, fmt.Errorf("load workloads: %w", err)
	}

	var cards []*component.Card

	fullMetrics := true

	for i := range workloads {
		workload := &workloads[i]

		card, supportsMetrics, err := CreateCard(workload, namespace)
		if err != nil {
			return nil, false, fmt.Errorf("create workload card: %w", err)
		}

		if !supportsMetrics {
			fullMetrics = false
		}

		cards = append(cards, card)
	}

	return cards, fullMetrics, nil
}

func CreateCard(workload *Workload, namespace string) (*component.Card, bool, error) {
	supportsMetrics := true

	workloadSummary, err := CreateWorkloadSummary(workload, component.DonutChartSizeMedium)
	if err != nil {
		return nil, false, fmt.Errorf("create workload summary: %w", err)
	}

	var section component.FlexLayoutSection

	if workloadSummary.MetricsEnabled {
		section = component.FlexLayoutSection{
			{
				Width: component.WidthThird,
				View:  workloadSummary.Summary,
			},
			{
				Width: component.WidthThird,
				View:  workloadSummary.Memory,
			},
			{
				Width: component.WidthThird,
				View:  workloadSummary.CPU,
			},
		}
	} else {
		section = component.FlexLayoutSection{
			{
				Width: component.WidthFull,
				View:  workloadSummary.Summary,
			},
		}
		supportsMetrics = false
	}

	cardPath := path.Join("/workloads/namespace", namespace, "detail", workload.Name)
	cardTitle := component.NewLink("", workload.Name, cardPath)

	card := component.NewCard([]component.TitleComponent{cardTitle})

	layout := component.NewFlexLayout(workload.Name)

	layout.AddSections(section)
	card.SetBody(layout)

	return card, supportsMetrics, nil
}

// ClusterWorkloadLoaderOption is option for configuring ClusterWorkloadLoader.
type ClusterWorkloadLoaderOption func(wl *ClusterWorkloadLoader)

// ClusterWorkloadLoader loads workloads from a Kubernetes cluster.
type ClusterWorkloadLoader struct {
	ObjectStatuser   func(context.Context, runtime.Object, store.Store) (objectstatus.ObjectStatus, error)
	ObjectStore      store.Store
	PodMetricsLoader PodMetricsLoader
}

// NewWorkloadLoader creates an instance of ClusterWorkloadLoader.
func NewClusterWorkloadLoader(objectStore store.Store, pml PodMetricsLoader, options ...ClusterWorkloadLoaderOption) (*ClusterWorkloadLoader, error) {
	wl := &ClusterWorkloadLoader{
		ObjectStatuser:   objectstatus.Status,
		ObjectStore:      objectStore,
		PodMetricsLoader: pml,
	}

	for _, option := range options {
		option(wl)
	}

	if wl.ObjectStore == nil {
		return nil, fmt.Errorf("object store is nil")
	}

	if wl.PodMetricsLoader == nil {
		return nil, fmt.Errorf("pod metrics loader is nil")
	}

	return wl, nil
}

// Load loads workloads from a cluster.
func (wl *ClusterWorkloadLoader) Load(ctx context.Context, namespace string) ([]Workload, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is blank")
	}

	podKey := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	podList, _, err := wl.ObjectStore.List(ctx, podKey)
	if err != nil {
		return nil, fmt.Errorf("unable to list pods in '%s': %w", namespace, err)
	}

	workloadTracker := make(map[types.UID]*Workload)

	for i := range podList.Items {
		object := &podList.Items[i]

		status, err := wl.ObjectStatuser(ctx, object, wl.ObjectStore)
		if err != nil {
			return nil, fmt.Errorf("get status for pod '%s': %w", object.GetName(), err)
		}

		owner, err := objectOwner(ctx, wl.ObjectStore, object)
		if err != nil {
			return nil, fmt.Errorf("find owner for pod '%s': %w", object.GetName(), err)
		}

		uid := owner.GetUID()

		if _, ok := workloadTracker[uid]; !ok {
			workloadTracker[uid] = NewWorkload(owner.GetName(), "application")
		}

		workload := workloadTracker[uid]
		workload.Owner = owner

		resourceList, err := wl.podMetrics(object)
		if err != nil {
			if IsPodMetricsNotSupported(err) {
				workload.SetPodMetricsDisabled()
				workload.AddPodStatus(status.Status(), object, corev1.ResourceList{})

				continue
			} else {
				return nil, fmt.Errorf("get metrics for pod '%s': %w", object.GetName(), err)
			}
		}

		workload.AddPodStatus(status.Status(), object, *resourceList)
	}

	var workloads []Workload
	for k := range workloadTracker {
		workloads = append(workloads, *workloadTracker[k])

	}

	sortWorkloads(workloads)
	return workloads, nil
}

func (wl *ClusterWorkloadLoader) podMetrics(pod *unstructured.Unstructured) (*corev1.ResourceList, error) {
	if pod == nil {
		return nil, fmt.Errorf("pod is nil")
	}

	supportsPodMetrics, err := wl.PodMetricsLoader.SupportsMetrics()
	if err != nil {
		return nil, fmt.Errorf("check for pod metrics support: %w", err)
	}

	if !supportsPodMetrics {
		return nil, &NoPodMetricsErr{}
	}

	namespace := pod.GetNamespace()
	name := pod.GetName()

	object, found, err := wl.PodMetricsLoader.Load(namespace, name)
	if err != nil {
		return nil, fmt.Errorf("get pod metrics: %w", err)
	}

	if !found {
		return nil, &NoPodMetricsErr{}
	}

	containersRaw, found, err := unstructured.NestedSlice(object.Object, "containers")
	if err != nil {
		return nil, fmt.Errorf("get containers in pod metrics")
	}
	if !found {
		return nil, fmt.Errorf("unable to find containers in pod metrics")
	}

	actual := corev1.ResourceList{}

	for i := range containersRaw {
		container, ok := containersRaw[i].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected container to be an object; it was '%T'", containersRaw[i])
		}

		usage, found, err := unstructured.NestedMap(container, "usage")
		if err != nil {
			return nil, fmt.Errorf("get container metric usage: %w", err)
		}

		if !found {
			return nil, fmt.Errorf("container metric usage was not found")
		}

		cpu, err := parseAsResourceList(usage, "cpu")
		if err != nil {
			return nil, fmt.Errorf("parse memory from resource list: %w", err)
		}
		memory, err := parseAsResourceList(usage, "memory")
		if err != nil {
			return nil, fmt.Errorf("parse memory from resource list: %w", err)
		}

		actual = addResourceList(actual, cpu, corev1.ResourceCPU)
		actual = addResourceList(actual, memory, corev1.ResourceMemory)
	}

	return &actual, nil
}

// PodCPUStat creates a single stat component for pod cpu. It will summarize all the pods in the workload.
func PodCPUStat(workload *Workload) (*component.SingleStat, error) {
	return podStat(workload, corev1.ResourceCPU, "CPU(cores)")
}

// PodMemoryStat creates a single stats component for pod memory. It will summarize all the pods in the workload.
func PodMemoryStat(workload *Workload) (*component.SingleStat, error) {
	return podStat(workload, corev1.ResourceMemory, "Memory(bytes)")
}

func podStat(workload *Workload, resourceName corev1.ResourceName, title string) (*component.SingleStat, error) {
	resources, metrics, err := summarizePodWithMetric(workload)
	if err != nil {
		return nil, err
	}

	background := WorkloadStatusColorOK

	actual := metrics[resourceName]

	av := actual.Value()
	r := resources.Requests[resourceName]
	rv := r.Value()
	l := resources.Limits[resourceName]
	lv := l.Value()
	if resourceName == corev1.ResourceCPU {
		av = actual.MilliValue()
		rv = r.MilliValue()
		lv = l.MilliValue()
	}

	if !r.IsZero() && av > rv {
		background = WorkloadStatusColorWarning
	}

	if !l.IsZero() && av > lv {
		background = WorkloadStatusColorError
	}

	m := metrics[resourceName]
	stat := component.NewSingleStat(title, printResourceUsage(resourceName, m), background)
	return stat, nil
}

func summarizePodWithMetric(workload *Workload) (corev1.ResourceRequirements, corev1.ResourceList, error) {
	if workload == nil {
		return corev1.ResourceRequirements{}, corev1.ResourceList{}, fmt.Errorf("workload is nil")
	}

	resources := corev1.ResourceRequirements{}
	metrics := corev1.ResourceList{}

	for _, podWithMetric := range workload.PodsWithMetrics() {
		pr, err := summarizePodResources(podWithMetric.Pod)
		if err != nil {
			return corev1.ResourceRequirements{}, corev1.ResourceList{}, fmt.Errorf("summarize pod resources: %w", err)
		}

		resources = CombineResourceRequirements(resources, *pr)
		metrics = addResourceList(metrics, *podWithMetric.ResourceList.Memory(), corev1.ResourceMemory)
		metrics = addResourceList(metrics, *podWithMetric.ResourceList.Cpu(), corev1.ResourceCPU)
	}

	return resources, metrics, nil
}

// CombineResourceRequirements combines two resource requirements into a new resource requirement.
func CombineResourceRequirements(a, b corev1.ResourceRequirements) corev1.ResourceRequirements {
	out := corev1.ResourceRequirements{
		Limits:   corev1.ResourceList{},
		Requests: corev1.ResourceList{},
	}

	out.Limits = addResourceList(out.Limits, *a.Limits.Cpu(), corev1.ResourceCPU)
	out.Limits = addResourceList(out.Limits, *a.Limits.Memory(), corev1.ResourceMemory)
	out.Limits = addResourceList(out.Limits, *b.Limits.Cpu(), corev1.ResourceCPU)
	out.Limits = addResourceList(out.Limits, *b.Limits.Memory(), corev1.ResourceMemory)
	out.Requests = addResourceList(out.Requests, *a.Requests.Cpu(), corev1.ResourceCPU)
	out.Requests = addResourceList(out.Requests, *a.Requests.Memory(), corev1.ResourceMemory)
	out.Requests = addResourceList(out.Requests, *b.Requests.Cpu(), corev1.ResourceCPU)
	out.Requests = addResourceList(out.Requests, *b.Requests.Memory(), corev1.ResourceMemory)

	return out
}

func summarizePodResources(object *unstructured.Unstructured) (*corev1.ResourceRequirements, error) {
	if object == nil {
		return nil, fmt.Errorf("object is nil")
	}

	var pod corev1.Pod
	if err := scheme.Scheme.Convert(object, &pod, nil); err != nil {
		return nil, fmt.Errorf("object is not a pod: %w", err)
	}

	requirements := corev1.ResourceRequirements{
		Limits:   corev1.ResourceList{},
		Requests: corev1.ResourceList{},
	}

	containers := pod.Spec.Containers
	for i := range containers {
		limit := containers[i].Resources.Limits
		request := containers[i].Resources.Requests

		requirements.Limits = addResourceList(requirements.Limits, *limit.Cpu(), corev1.ResourceCPU)
		requirements.Limits = addResourceList(requirements.Limits, *limit.Memory(), corev1.ResourceMemory)
		requirements.Requests = addResourceList(requirements.Requests, *request.Cpu(), corev1.ResourceCPU)
		requirements.Requests = addResourceList(requirements.Requests, *request.Memory(), corev1.ResourceMemory)
	}

	return &requirements, nil
}

type WorkloadSummary struct {
	Summary        component.Component
	Memory         component.Component
	CPU            component.Component
	MetricsEnabled bool
}

func CreateWorkloadSummary(workload *Workload, summarySize component.DonutChartSize) (WorkloadSummary, error) {
	if workload == nil {
		return WorkloadSummary{}, fmt.Errorf("can't create summary for nil workload")
	}

	summaryChart, err := workload.DonutChart(summarySize)
	if err != nil {
		return WorkloadSummary{}, fmt.Errorf("create workload summary component: %w", err)
	}

	if !workload.PodMetricsEnabled() {
		return WorkloadSummary{
			Summary:        summaryChart,
			MetricsEnabled: false,
		}, nil

	}

	memoryStat, err := PodMemoryStat(workload)
	if err != nil {
		return WorkloadSummary{}, fmt.Errorf("create workload memory component: %w", err)
	}

	cpuStat, err := PodCPUStat(workload)
	if err != nil {
		return WorkloadSummary{}, fmt.Errorf("create workload cpu component: %w", err)
	}

	return WorkloadSummary{
		Summary:        summaryChart,
		Memory:         memoryStat,
		CPU:            cpuStat,
		MetricsEnabled: true,
	}, nil
}
