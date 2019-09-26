/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	autoscalingapiv1 "k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	"k8s.io/kubernetes/pkg/apis/core"

	"github.com/vmware/octant/pkg/view/component"
)

// HorizontalPodAutoscalerListHandler is a printFunc that lists horizontal pod autoscalers
func HorizontalPodAutoscalerListHandler(_ context.Context, list *autoscalingv1.HorizontalPodAutoscalerList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("horizontalpod handler list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Targets", "Minimum Pods", "Maximum Pods", "Replicas", "Age")
	tbl := component.NewTable("Horizontal Pod Autoscalers",
		"We couldn't find any horizontal pod autoscalers", cols)

	for _, horizontalPodAutoscaler := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&horizontalPodAutoscaler, horizontalPodAutoscaler.Name)
		if err != nil {
			return nil, err
		}

		convertedHPA, err := convertToAutoscaling(&horizontalPodAutoscaler)
		if err != nil {
			return nil, errors.Wrap(err, "can't convert hpa")
		}

		aggregatedMetricTargets, err := getMetricsOverview(convertedHPA.Spec.Metrics, convertedHPA.Status.CurrentMetrics)
		if err != nil {
			return nil, errors.Wrap(err, "can't combine metrics")
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(horizontalPodAutoscaler.Labels)
		row["Targets"] = component.NewText(aggregatedMetricTargets)
		row["Minimum Pods"] = component.NewText(fmt.Sprintf("%d", *horizontalPodAutoscaler.Spec.MinReplicas))
		row["Maximum Pods"] = component.NewText(fmt.Sprintf("%d", horizontalPodAutoscaler.Spec.MaxReplicas))
		row["Replicas"] = component.NewText(fmt.Sprintf("%d", horizontalPodAutoscaler.Status.CurrentReplicas))
		row["Age"] = component.NewTimestamp(horizontalPodAutoscaler.CreationTimestamp.Time)

		tbl.Add(row)
	}
	return tbl, nil
}

// HorizontalPodAutoscalerHandler is a printFunc that prints a HorizontalPodAutoscaler
func HorizontalPodAutoscalerHandler(ctx context.Context, horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, options Options) (component.Component, error) {
	o := NewObject(horizontalPodAutoscaler)
	o.EnableEvents()

	hh, err := newHorizontalPodAutoscalerHandler(horizontalPodAutoscaler, o)
	if err != nil {
		return nil, err
	}

	if err := hh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler configuration")
	}

	if err := hh.Status(); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler status")
	}

	if err := hh.Metrics(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler metrics")
	}

	if err := hh.Conditions(); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler conditions")
	}

	return o.ToComponent(ctx, options)
}

func createHorizontalPodAutoscalerSummaryStatus(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) (*component.Summary, error) {
	if horizontalPodAutoscaler == nil {
		return nil, errors.New("unable to generate status for a nil horizontalpodautoscaler")
	}

	status := horizontalPodAutoscaler.Status

	summary := component.NewSummary("Status")

	sections := component.SummarySections{}

	if status.ObservedGeneration != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Observed Generation",
			Content: component.NewText(fmt.Sprintf("%d", *status.ObservedGeneration)),
		})
	}

	if status.LastScaleTime != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Last Scale Time",
			Content: component.NewTimestamp(status.LastScaleTime.Time),
		})
	}

	sections.AddText("Current Replicas", fmt.Sprintf("%d", status.CurrentReplicas))
	sections.AddText("Desired Replicas", fmt.Sprintf("%d", status.DesiredReplicas))

	if status.CurrentCPUUtilizationPercentage != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Current CPU Utilization Percentage",
			Content: component.NewText(fmt.Sprintf("%d", *status.CurrentCPUUtilizationPercentage)),
		})
	}

	summary.Add(sections...)

	return summary, nil
}

func createHorizontalPodAutoscalerMetricsStatusView(metricStatus *autoscaling.MetricStatus, options Options) (*component.Summary, error) {
	if metricStatus == nil {
		return nil, errors.New("unable to generate metrics from a nil metric status")
	}

	sections := component.SummarySections{}
	summary := component.NewSummary(fmt.Sprintf("Metric"))
	sections.AddText("Type", fmt.Sprintf("%s", metricStatus.Type))

	if metricStatus.Object != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Name",
			Content: component.NewText(metricStatus.Object.Metric.Name),
		})

		if metricStatus.Object.DescribedObject.Name != "" {
			sections = append(sections, component.SummarySection{
				Header:  "Described Object Name",
				Content: component.NewText(metricStatus.Object.DescribedObject.Name),
			})
			sections = append(sections, component.SummarySection{
				Header:  "Described Object API Version",
				Content: component.NewText(metricStatus.Object.DescribedObject.APIVersion),
			})
			sections = append(sections, component.SummarySection{
				Header:  "Described Object Kind",
				Content: component.NewText(metricStatus.Object.DescribedObject.Kind),
			})
		}
	}

	if metricStatus.Pods != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Name",
			Content: component.NewText(metricStatus.Pods.Metric.Name),
		})
		if metricStatus.Pods.Current.AverageUtilization != nil {
			sections = append(sections, component.SummarySection{
				Header:  "Average Utilization",
				Content: component.NewText(fmt.Sprintf("%d%%", *metricStatus.Pods.Current.AverageUtilization)),
			})
		}
		if metricStatus.Pods.Current.AverageValue != nil {
			sections = append(sections, component.SummarySection{
				Header:  "Average Value",
				Content: component.NewText(metricStatus.Pods.Current.AverageValue.String()),
			})
		}
		if metricStatus.Pods.Current.Value != nil {
			sections = append(sections, component.SummarySection{
				Header:  "Value",
				Content: component.NewText(metricStatus.Pods.Current.Value.String()),
			})
		}
	}

	if metricStatus.Resource != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Name",
			Content: component.NewText(string(metricStatus.Resource.Name)),
		})
		if metricStatus.Resource.Current.AverageUtilization != nil {
			sections = append(sections, component.SummarySection{
				Header:  "Average Utilization",
				Content: component.NewText(fmt.Sprintf("%d%%", *metricStatus.Resource.Current.AverageUtilization)),
			})
		}
		if metricStatus.Resource.Current.AverageValue != nil {
			sections = append(sections, component.SummarySection{
				Header:  "Average Value",
				Content: component.NewText(metricStatus.Resource.Current.AverageValue.String()),
			})
		}
		if metricStatus.Resource.Current.Value != nil {
			sections = append(sections, component.SummarySection{
				Header:  "Value",
				Content: component.NewText(metricStatus.Resource.Current.Value.String()),
			})
		}
	}

	summary.Add(sections...)

	return summary, nil
}

func createHorizontalPodAutoscalerConditionsView(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) (*component.Table, error) {
	if horizontalPodAutoscaler == nil {
		return nil, errors.New("unable to generate conditions from a nil horizontalpodautoscaler")
	}

	convertedHPA, err := convertToAutoscaling(horizontalPodAutoscaler)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert hpa")
	}

	cols := component.NewTableCols("Type", "Reason", "Status", "Message", "Last Transition")
	table := component.NewTable("Conditions", "There are no horizontalpodautoscaler conditions!", cols)

	for _, condition := range convertedHPA.Status.Conditions {
		row := component.TableRow{
			"Type":            component.NewText(string(condition.Type)),
			"Reason":          component.NewText(condition.Reason),
			"Status":          component.NewText(string(condition.Status)),
			"Message":         component.NewText(condition.Message),
			"Last Transition": component.NewTimestamp(condition.LastTransitionTime.Time),
		}

		table.Add(row)
	}

	table.Sort("Type", false)

	return table, nil
}

// HorizontalPodAutoscalerConfiguration generates a horizontalpodautoscaler configuration
type HorizontalPodAutoscalerConfiguration struct {
	horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler
}

// NewHorizontalPodAutoscalerConfiguration creates an instance of HorizontalPodAutoscalerConfiguration
func NewHorizontalPodAutoscalerConfiguration(hpa *autoscalingv1.HorizontalPodAutoscaler) *HorizontalPodAutoscalerConfiguration {
	return &HorizontalPodAutoscalerConfiguration{
		horizontalPodAutoscaler: hpa,
	}
}

type horizontalPodAutoscalerObject interface {
	Config(options Options) error
	Status() error
	Metrics(ctx context.Context, options Options) error
	Conditions() error
}

type horizontalPodAutoscalerHandler struct {
	horizontalPodAutoScaler *autoscalingv1.HorizontalPodAutoscaler
	configFunc              func(*autoscalingv1.HorizontalPodAutoscaler, Options) (*component.Summary, error)
	statusFunc              func(*autoscalingv1.HorizontalPodAutoscaler) (*component.Summary, error)
	metricsFunc             func(context.Context, *autoscaling.MetricStatus, Options) (*component.Summary, error)
	conditionsFunc          func(*autoscalingv1.HorizontalPodAutoscaler) (*component.Table, error)
	object                  *Object
}

// Create creates a horizontalpodautoscaler configuration sumamry
func (hc *HorizontalPodAutoscalerConfiguration) Create(options Options) (*component.Summary, error) {
	if hc.horizontalPodAutoscaler == nil {
		return nil, errors.New("horizontalpodautoscaler is nil")
	}

	hpa := hc.horizontalPodAutoscaler

	sections := component.SummarySections{}

	scaleTarget, err := forScaleTarget(hpa, &hpa.Spec.ScaleTargetRef, options)
	if err != nil {
		return nil, err
	}

	sections = append(sections, component.SummarySection{
		Header:  "Reference",
		Content: scaleTarget,
	})

	minReplicas := fmt.Sprintf("%d", *hpa.Spec.MinReplicas)
	maxReplicas := fmt.Sprintf("%d", hpa.Spec.MaxReplicas)
	sections.AddText("Min Replicas", minReplicas)
	sections.AddText("Max Replicas", maxReplicas)

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

var _ horizontalPodAutoscalerObject = (*horizontalPodAutoscalerHandler)(nil)

func newHorizontalPodAutoscalerHandler(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, object *Object) (*horizontalPodAutoscalerHandler, error) {
	if horizontalPodAutoscaler == nil {
		return nil, errors.New("can't print a nil horizontalpodautoscaler")
	}

	if object == nil {
		return nil, errors.New("can't print horizontalpodautoscaler using a nil object printer")
	}

	hh := &horizontalPodAutoscalerHandler{
		horizontalPodAutoScaler: horizontalPodAutoscaler,
		configFunc:              defaultHorizontalPodAutoscalerConfig,
		statusFunc:              defaultHorizontalPodAutoscalerStatus,
		metricsFunc:             defaultHorizontalPodAutoscalerMetrics,
		conditionsFunc:          defaultHorizontalPodAutoscalerConditions,
		object:                  object,
	}

	return hh, nil
}

func (h *horizontalPodAutoscalerHandler) Config(options Options) error {
	out, err := h.configFunc(h.horizontalPodAutoScaler, options)
	if err != nil {
		return err
	}

	h.object.RegisterConfig(out)
	return nil
}

func defaultHorizontalPodAutoscalerConfig(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, options Options) (*component.Summary, error) {
	return NewHorizontalPodAutoscalerConfiguration(horizontalPodAutoscaler).Create(options)
}

func (h *horizontalPodAutoscalerHandler) Status() error {
	out, err := h.statusFunc(h.horizontalPodAutoScaler)
	if err != nil {
		return err
	}

	h.object.RegisterSummary(out)
	return nil
}

func defaultHorizontalPodAutoscalerStatus(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) (*component.Summary, error) {
	return createHorizontalPodAutoscalerSummaryStatus(horizontalPodAutoscaler)
}

func (h *horizontalPodAutoscalerHandler) metrics(ctx context.Context, currentMetrics []autoscaling.MetricStatus, options Options) error {
	if h == nil || h.horizontalPodAutoScaler == nil {
		return errors.New("can't display metrics for nil horizontalpodautoscaler")
	}

	for i := range currentMetrics {
		metric := currentMetrics[i]

		h.object.RegisterItems(ItemDescriptor{
			Width: component.WidthFull,
			Func: func() (component.Component, error) {
				return h.metricsFunc(ctx, &metric, options)
			},
		})
	}

	return nil
}

func (h *horizontalPodAutoscalerHandler) Metrics(ctx context.Context, options Options) error {
	if h.horizontalPodAutoScaler == nil {
		return errors.New("can't display metrics for nil horizontalpodautoscaler")
	}

	convertedHPA, err := convertToAutoscaling(h.horizontalPodAutoScaler)
	if err != nil {
		return errors.New("can't convert hpa")
	}

	return h.metrics(ctx, convertedHPA.Status.CurrentMetrics, options)
}

func defaultHorizontalPodAutoscalerMetrics(ctx context.Context, metricStatus *autoscaling.MetricStatus, options Options) (*component.Summary, error) {
	return createHorizontalPodAutoscalerMetricsStatusView(metricStatus, options)
}

func (h *horizontalPodAutoscalerHandler) Conditions() error {
	if h.horizontalPodAutoScaler == nil {
		return errors.New("can't display conditions for nil horizontalpodautoscaler")
	}

	h.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return h.conditionsFunc(h.horizontalPodAutoScaler)
		},
	})

	return nil
}

func defaultHorizontalPodAutoscalerConditions(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) (*component.Table, error) {
	return createHorizontalPodAutoscalerConditionsView(horizontalPodAutoscaler)
}

// forScaleTarget returns a scale target for a cross version object reference
func forScaleTarget(object runtime.Object, scaleTarget *autoscalingv1.CrossVersionObjectReference, options Options) (*component.Link, error) {
	if scaleTarget == nil || object == nil {
		return component.NewLink("", "none", ""), nil
	}

	accessor := meta.NewAccessor()
	ns, err := accessor.Namespace(object)
	if err != nil {
		return component.NewLink("", "none", ""), nil
	}

	return options.Link.ForGVK(
		ns,
		scaleTarget.APIVersion,
		scaleTarget.Kind,
		scaleTarget.Name,
		scaleTarget.Name,
	)
}

// convertToAutoscaling converts default v1 to an the internal type for getting metrics and conditions
func convertToAutoscaling(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) (*autoscaling.HorizontalPodAutoscaler, error) {
	convertedHPA := &autoscaling.HorizontalPodAutoscaler{}
	hpav1 := horizontalPodAutoscaler.DeepCopy()
	if err := autoscalingapiv1.Convert_v1_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler(hpav1, convertedHPA, nil); err != nil {
		return nil, err
	}
	return convertedHPA, nil
}

// convertToV1 converts an the internal type to v1
func convertToV1(horizontalPodAutoscaler *autoscaling.HorizontalPodAutoscaler) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	convertedHPA := &autoscalingv1.HorizontalPodAutoscaler{}
	hpa := horizontalPodAutoscaler.DeepCopy()
	if err := autoscalingapiv1.Convert_autoscaling_HorizontalPodAutoscaler_To_v1_HorizontalPodAutoscaler(hpa, convertedHPA, nil); err != nil {
		return nil, err
	}
	return convertedHPA, nil
}

func getMetricsOverview(metricSpec []autoscaling.MetricSpec, metricStatus []autoscaling.MetricStatus) (string, error) {
	var targets = make(map[string]string)
	var currents = make(map[string]string)

	for _, m := range metricSpec {
		switch m.Type {
		case autoscaling.ObjectMetricSourceType:
			if m.Object != nil {
				if m.Object.Metric.Name != "" {
					target, err := getMetricValue(&m.Object.Target)
					if err != nil {
						return "", err
					}
					targets[m.Object.Metric.Name] = target
				}
			}
		case autoscaling.PodsMetricSourceType:
			if m.Pods != nil {
				if m.Pods.Metric.Name != "" {
					target := m.Pods.Target.AverageValue.String()
					targets[m.Pods.Metric.Name] = target
				}
			}
		case autoscaling.ResourceMetricSourceType:
			if m.Resource != nil {
				if m.Resource.Name == core.ResourceCPU || m.Resource.Name == core.ResourceMemory {
					target, err := getMetricValue(&m.Resource.Target)
					if err != nil {
						return "", err
					}
					targets[string(m.Resource.Name)] = target
				}
			}
		}
	}

	for _, m := range metricStatus {
		current, err := getMetricStatusValue(&m)
		if err != nil {
			return "", err
		}

		switch m.Type {
		case autoscaling.ObjectMetricSourceType:
			if m.Object != nil {
				if m.Object.Metric.Name != "" {
					currents[m.Object.Metric.Name] = current
				}
			}
		case autoscaling.PodsMetricSourceType:
			if m.Pods != nil {
				if m.Pods.Metric.Name != "" {
					currents[m.Pods.Metric.Name] = current
				}
			}
		case autoscaling.ResourceMetricSourceType:
			if m.Resource != nil {
				if m.Resource.Name == core.ResourceCPU || m.Resource.Name == core.ResourceMemory {
					currents[string(m.Resource.Name)] = current
				}
			}
		}
	}

	var result []string
	keys := make([]string, 0, len(targets))
	for k := range targets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if currents[k] == "" {
			currents[k] = "<unknown>"
		}
		result = append(result, currents[k]+"/"+targets[k])
	}

	return strings.Join(result, ", "), nil
}

// getMetricValue returns a string representation of Value, AverageValue, or AverageUtilization
func getMetricValue(metricTarget *autoscaling.MetricTarget) (string, error) {
	var value string
	if metricTarget == nil {
		return "", errors.New("nil metric target")
	}

	switch metricTarget.Type {
	case autoscaling.UtilizationMetricType:
		value = fmt.Sprintf("%d%%", *metricTarget.AverageUtilization)
	case autoscaling.AverageValueMetricType:
		value = metricTarget.AverageValue.String()
	case autoscaling.ValueMetricType:
		value = metricTarget.Value.String()
	}

	return value, nil
}

func getMetricStatusValue(metricStatus *autoscaling.MetricStatus) (string, error) {
	var value string
	if metricStatus == nil {
		return "", errors.New("nil metric status")
	}

	switch metricStatus.Type {
	case autoscaling.ObjectMetricSourceType:
		value = metricValueStatusStringer(metricStatus.Object.Current)
	case autoscaling.PodsMetricSourceType:
		value = metricValueStatusStringer(metricStatus.Pods.Current)
	case autoscaling.ResourceMetricSourceType:
		value = metricValueStatusStringer(metricStatus.Resource.Current)
	case autoscaling.ExternalMetricSourceType:
		value = metricValueStatusStringer(metricStatus.External.Current)
	}

	return value, nil
}

func metricValueStatusStringer(status autoscaling.MetricValueStatus) string {
	if status.AverageUtilization != nil {
		return fmt.Sprintf("%d", *status.AverageUtilization)
	}
	if status.AverageValue != nil {
		return status.AverageValue.String()
	}
	if status.Value != nil {
		return status.Value.String()
	}
	return ""
}
