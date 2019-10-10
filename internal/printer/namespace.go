/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"path"
	"sort"

	"github.com/pkg/errors"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	namespaceListCols           = component.NewTableCols("Name", "Labels", "Status", "Age")
	namespaceResourceQuotasCols = component.NewTableCols("Resource", "Used", "Limit")
	namespaceResourceLimitsCols = component.NewTableCols("Type", "Resource", "Min", "Max", "Default Request", "Default Limit", "Limit/Request Ratio")
)

func NamespaceListHandler(_ context.Context, list *corev1.NamespaceList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("namespace list is nil")
	}

	table := component.NewTable("Namespaces", "We couldn't find any namespaces!", namespaceListCols)

	for _, ns := range list.Items {
		row := component.TableRow{}
		p := path.Join("/cluster-overview/namespaces", ns.Name)
		row["Name"] = component.NewLink("", ns.Name, p)
		row["Labels"] = component.NewLabels(ns.Labels)
		row["Status"] = component.NewText(namespaceStatus(ns))
		row["Age"] = component.NewTimestamp(ns.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

func NamespaceHandler(ctx context.Context, namespace *corev1.Namespace, options Options) (component.Component, error) {
	o := NewObject(namespace)

	nh, err := newNamespaceHandler(namespace, o)
	if err != nil {
		return nil, err
	}

	if err := nh.Status(options); err != nil {
		return nil, errors.Wrap(err, "print namespace status")
	}
	if err := nh.ResourceLimits(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print namespace resource limits")
	}
	if err := nh.ResourceQuotas(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print namespace resource quotas")
	}
	return o.ToComponent(ctx, options)
}

func namespaceStatus(namespace corev1.Namespace) string {
	return fmt.Sprintf("%s", namespace.Status.Phase)
}

type namespaceObject interface {
	Status(options Options) error
	ResourceQuotas(ctx context.Context, options Options) error
	ResourceLimits(ctx context.Context, options Options) error
}

type namespaceHandler struct {
	namespace          *corev1.Namespace
	statusFunc         func(*corev1.Namespace, Options) (*component.Summary, error)
	resourceQuotasFunc func(context.Context, *corev1.Namespace, Options) (*component.FlexLayout, error)
	resourceLimitsFunc func(context.Context, *corev1.Namespace, Options) (*component.Table, error)
	object             *Object
}

var _ namespaceObject = (*namespaceHandler)(nil)

func newNamespaceHandler(namespace *corev1.Namespace, object *Object) (*namespaceHandler, error) {
	if namespace == nil {
		return nil, errors.New("can't print a nil namespace")
	}

	if object == nil {
		return nil, errors.New("can't print node using a nil object printer")
	}

	nh := &namespaceHandler{
		namespace:          namespace,
		statusFunc:         defaultNamespaceStatus,
		resourceQuotasFunc: defaultNamespaceResourceQuotas,
		resourceLimitsFunc: defaultNamespaceResourceLimits,
		object:             object,
	}
	return nh, nil
}

func (n *namespaceHandler) Status(options Options) error {
	out, err := n.statusFunc(n.namespace, options)
	if err != nil {
		return err
	}
	n.object.RegisterSummary(out)
	return nil
}

func (n *namespaceHandler) ResourceQuotas(ctx context.Context, options Options) error {
	n.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return n.resourceQuotasFunc(ctx, n.namespace, options)
		},
	})
	return nil
}

func (n *namespaceHandler) ResourceLimits(ctx context.Context, options Options) error {
	n.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return n.resourceLimitsFunc(ctx, n.namespace, options)
		},
	})
	return nil
}

// NamespaceStatus creates a namespace status component.
type NamespaceStatus struct {
	namespace *corev1.Namespace
}

// NewNamespaceStatus creates an instance of NamespaceStatus
func NewNamespaceStatus(namespace *corev1.Namespace) *NamespaceStatus {
	return &NamespaceStatus{
		namespace: namespace,
	}
}

// Create creates a namespace configuration summary
func (n *NamespaceStatus) Create(options Options) (*component.Summary, error) {
	if n == nil || n.namespace == nil {
		return nil, errors.New("cannot generate status for nil node")
	}

	summary := component.NewSummary("Status", []component.SummarySection{
		{
			Header:  "Phase",
			Content: component.NewText(fmt.Sprintf("%s", n.namespace.Status.Phase)),
		},
	}...)

	return summary, nil
}

func defaultNamespaceStatus(namespace *corev1.Namespace, options Options) (*component.Summary, error) {
	return NewNamespaceStatus(namespace).Create(options)
}

// NamespaceResourceQuotas creates a namespace resource quota component.
type NamespaceResourceQuotas struct {
	namespace *corev1.Namespace
}

func NewNamespaceResourceQuotas(namespace *corev1.Namespace) *NamespaceResourceQuotas {
	return &NamespaceResourceQuotas{
		namespace: namespace,
	}
}

func (n *NamespaceResourceQuotas) Create(ctx context.Context, options Options) (*component.FlexLayout, error) {
	if n == nil || n.namespace == nil {
		return nil, errors.New("cannot generate resources for nil node")
	}
	return printNamespaceResourceQuotas(ctx, n.namespace, options)
}

func printNamespaceResourceQuotas(ctx context.Context, namespace *corev1.Namespace, options Options) (*component.FlexLayout, error) {
	objectStore := options.DashConfig.ObjectStore()
	key := store.Key{
		Namespace:  namespace.Name,
		APIVersion: "v1",
		Kind:       "ResourceQuota",
	}

	list, _, err := objectStore.List(ctx, key)
	items := make(map[string]component.FlexLayoutItem, len(list.Items))
	for i := range list.Items {
		rq := &corev1.ResourceQuota{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(list.Items[i].Object, rq)
		if err != nil {
			return nil, err
		}

		table := component.NewTable(rq.Name, "There are no resource quotas", namespaceResourceQuotasCols)
		quotas := paresResourceQuotas(rq)
		for _, resource := range resourceQuotaKeys(rq) {
			row := component.TableRow{}
			row["Resource"] = component.NewText(resource)
			row["Used"] = component.NewText(quotas["used"][resource])
			row["Limit"] = component.NewText(quotas["hard"][resource])
			table.Add(row)
		}
		table.Sort("Resource", false)
		items[rq.Name] = component.FlexLayoutItem{Width: component.WidthHalf, View: table}
	}

	fl := component.NewFlexLayout("Resource Quotas")
	fl.AddSections(createSortedResourceQuotaSections("Resource Quotas", items))
	return fl, err
}

func createSortedResourceQuotaSections(title string, sectionMap map[string]component.FlexLayoutItem) []component.FlexLayoutItem {
	length := len(sectionMap)
	// length + 1 = title section + items
	sections := make([]component.FlexLayoutItem, 0, length+1)

	// Add section Title as first element.
	sections = append(sections, createTitleItem(title, map[string]string{}))

	// Sort the Resource Quotas by name so they don't shift around during poller runs.
	keys := make([]string, 0, length)
	for k := range sectionMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		sections = append(sections, sectionMap[k])
	}
	return sections
}

func createTitleItem(title string, labels map[string]string) component.FlexLayoutItem {
	label := component.NewLabels(labels)
	label.Metadata.SetTitleText(title)
	return component.FlexLayoutItem{
		Width: component.WidthFull,
		View:  label,
	}
}

func resourceQuotaKeys(rq *corev1.ResourceQuota) []string {
	var keys []string
	for name, _ := range rq.Spec.Hard {
		keys = append(keys, name.String())
	}
	sort.Strings(keys)
	return keys
}

func paresResourceQuotas(rq *corev1.ResourceQuota) map[string]map[string]string {
	quotas := map[string]map[string]string{
		"hard": make(map[string]string, len(rq.Status.Hard)),
		"used": make(map[string]string, len(rq.Status.Used)),
	}
	for name, v := range rq.Status.Hard {
		quotas["hard"][name.String()] = v.String()
	}

	for name, v := range rq.Status.Used {
		quotas["used"][name.String()] = v.String()
	}
	return quotas
}

func defaultNamespaceResourceQuotas(ctx context.Context, namespace *corev1.Namespace, options Options) (*component.FlexLayout, error) {
	return NewNamespaceResourceQuotas(namespace).Create(ctx, options)
}

// NamespaceResourceLimits creates a namespace resource limit component.
type NamespaceResourceLimits struct {
	namespace *corev1.Namespace
}

func NewNamespaceResourceLimits(namespace *corev1.Namespace) *NamespaceResourceLimits {
	return &NamespaceResourceLimits{
		namespace: namespace,
	}
}

// Create creates a namespace limit component.
func (n *NamespaceResourceLimits) Create(ctx context.Context, options Options) (*component.Table, error) {
	if n == nil || n.namespace == nil {
		return nil, errors.New("cannot generate resources for nil namespace")
	}
	return printNamespaceResourceLimits(ctx, n.namespace, options)
}

func printNamespaceResourceLimits(ctx context.Context, namespace *corev1.Namespace, options Options) (*component.Table, error) {
	objectStore := options.DashConfig.ObjectStore()
	key := store.Key{
		Namespace:  namespace.Name,
		APIVersion: "v1",
		Kind:       "LimitRange",
	}

	table := component.NewTable("Resource Limits", "There are no resource limits", namespaceResourceLimitsCols)
	rows := map[string]component.TableRow{}

	list, _, err := objectStore.List(ctx, key)
	for i := range list.Items {
		lr := &corev1.LimitRange{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(list.Items[i].Object, lr)
		if err != nil {
			return nil, err
		}

		for _, item := range lr.Spec.Limits {
			sortKey, row, created := createResourceLimitCpuRow(item)
			if created {
				rows[sortKey] = row
			}

			sortKey, row, created = createResourceLimitMemoryRow(item)
			if created {
				rows[sortKey] = row
			}
		}
	}

	addResourceLimitRowsSorted(table, rows)
	return table, err
}

func addResourceLimitRowsSorted(table *component.Table, rows map[string]component.TableRow) {
	keys := make([]string, 0, len(rows))
	for k := range rows {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		table.Add(rows[key])
	}
}

func createResourceLimitCpuRow(item corev1.LimitRangeItem) (sortKey string, row component.TableRow, created bool) {
	if !item.Max.Cpu().IsZero() || !item.Min.Cpu().IsZero() || !item.Default.Cpu().IsZero() {
		limitType := fmt.Sprintf("%s", item.Type)
		minCpu := item.Min.Cpu().String()
		maxCpu := item.Max.Cpu().String()
		sortKey := fmt.Sprintf("%s-%s-%s", minCpu, maxCpu, limitType)

		row := component.TableRow{}
		row["Type"] = component.NewText(limitType)
		row["Resource"] = component.NewText("cpu")
		row["Min"] = component.NewText(minCpu)
		row["Max"] = component.NewText(maxCpu)
		row["Default Request"] = component.NewText(item.DefaultRequest.Cpu().String())
		row["Default Limit"] = component.NewText(item.Default.Cpu().String())
		row["Limit/Request Ratio"] = component.NewText(item.MaxLimitRequestRatio.Cpu().String())
		return sortKey, row, true
	}
	return "", nil, false
}

func createResourceLimitMemoryRow(item corev1.LimitRangeItem) (sortKey string, row component.TableRow, created bool) {
	if !item.Max.Memory().IsZero() || !item.Min.Memory().IsZero() || !item.Default.Memory().IsZero() {
		limitType := fmt.Sprintf("%s", item.Type)
		minMem := item.Min.Memory().String()
		maxMem := item.Max.Memory().String()
		sortKey := fmt.Sprintf("%s-%s-%s", minMem, maxMem, limitType)

		row := component.TableRow{}
		row["Type"] = component.NewText(limitType)
		row["Resource"] = component.NewText("memory")
		row["Min"] = component.NewText(minMem)
		row["Max"] = component.NewText(maxMem)
		row["Default Request"] = component.NewText(item.DefaultRequest.Memory().String())
		row["Default Limit"] = component.NewText(item.Default.Memory().String())
		row["Limit/Request Ratio"] = component.NewText(item.MaxLimitRequestRatio.Memory().String())
		return sortKey, row, true
	}
	return "", nil, false
}

func defaultNamespaceResourceLimits(ctx context.Context, namespace *corev1.Namespace, options Options) (*component.Table, error) {
	return NewNamespaceResourceLimits(namespace).Create(ctx, options)
}
