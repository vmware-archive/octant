/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/pkg/errors"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// HorizontalPodAutoscalerListHandler is a printFunc that lists horizontal pod autoscalers
func HorizontalPodAutoscalerListHandler(ctx context.Context, list *autoscalingv1.HorizontalPodAutoscalerList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("horizontalpod handler list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Targets", "Minimum Pods", "Maximum Pods", "Replicas", "Age")
	ot := NewObjectTable("Horizontal Pod Autoscalers",
		"We couldn't find any horizontal pod autoscalers", cols, options.DashConfig.ObjectStore())
	ot.EnablePluginStatus(options.DashConfig.PluginManager())
	for _, horizontalPodAutoscaler := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&horizontalPodAutoscaler, horizontalPodAutoscaler.Name)
		if err != nil {
			return nil, err
		}

		horizontalPodAutoscalerMetrics, horizontalPodAutoscalerCurrentMetrics, err := parseAnnotations(horizontalPodAutoscaler)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse annotations")
		}

		aggregatedMetricTargets, err := getCombinedMetrics(horizontalPodAutoscaler, horizontalPodAutoscalerMetrics, horizontalPodAutoscalerCurrentMetrics)
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

		if err := ot.AddRowForObject(ctx, &horizontalPodAutoscaler, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// HorizontalPodAutoscalerHandler is a printFunc that prints a HorizontalPodAutoscaler
func HorizontalPodAutoscalerHandler(ctx context.Context, horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, options Options) (component.Component, error) {
	o := NewObject(horizontalPodAutoscaler)
	o.EnableEvents()
	o.DisableConditions()

	hh, err := newHorizontalPodAutoscalerHandler(horizontalPodAutoscaler, o)
	if err != nil {
		return nil, err
	}

	if err := hh.Config(ctx, options); err != nil {
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

	horizontalPodAutoscalerMetrics, horizontalPodAutoscalerCurrentMetrics, err := parseAnnotations(*horizontalPodAutoscaler)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse annotations")
	}

	aggregatedMetricTargets, err := getCombinedMetrics(*horizontalPodAutoscaler, horizontalPodAutoscalerMetrics, horizontalPodAutoscalerCurrentMetrics)
	if err != nil {
		return nil, errors.Wrap(err, "can't combine metrics")
	}

	status := horizontalPodAutoscaler.Status

	summary := component.NewSummary("Status")

	sections := component.SummarySections{}

	sections.AddText("Targets", aggregatedMetricTargets)

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

func createHorizontalPodAutoscalerMetricsStatusView(metricStatus *autoscalingv1.MetricStatus, options Options) (*component.Summary, error) {
	if metricStatus == nil {
		return nil, errors.New("unable to generate metrics from a nil metric status")
	}

	sections := component.SummarySections{}
	summary := component.NewSummary(fmt.Sprintf("Metric"))

	sections.AddText("Type", fmt.Sprintf("%s", metricStatus.Type))

	if metricStatus.Object != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Name",
			Content: component.NewText(metricStatus.Object.MetricName),
		})

		if metricStatus.Object.Target.Name != "" {
			sections = append(sections, component.SummarySection{
				Header:  "Described Object Name",
				Content: component.NewText(metricStatus.Object.Target.Name),
			})
			sections = append(sections, component.SummarySection{
				Header:  "Described Object API Version",
				Content: component.NewText(metricStatus.Object.Target.APIVersion),
			})
			sections = append(sections, component.SummarySection{
				Header:  "Described Object Kind",
				Content: component.NewText(metricStatus.Object.Target.Kind),
			})
		}
	}

	if metricStatus.Pods != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Name",
			Content: component.NewText(metricStatus.Pods.MetricName),
		})
		sections = append(sections, component.SummarySection{
			Header:  "Average Utilization",
			Content: component.NewText(fmt.Sprint(&metricStatus.Pods.CurrentAverageValue)),
		})
	}

	if metricStatus.Resource != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Name",
			Content: component.NewText(string(metricStatus.Resource.Name)),
		})
		if metricStatus.Resource.CurrentAverageUtilization != nil {
			sections = append(sections, component.SummarySection{
				Header:  "Average Utilization",
				Content: component.NewText(fmt.Sprint(*metricStatus.Resource.CurrentAverageUtilization)),
			})
		}
		sections = append(sections, component.SummarySection{
			Header:  "Average Value",
			Content: component.NewText(metricStatus.Resource.CurrentAverageValue.String()),
		})
	}

	summary.Add(sections...)

	return summary, nil
}

var hpaConditionColumns = [][]string{
	{"Type", "type"},
	{"Reason", "reason"},
	{"Status", "status"},
	{"Message", "message"},
	{"Last Transition", "lastTransitionTime"},
}

func createHorizontalPodAutoscalerConditionsView(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) (*component.Table, error) {
	horizontalPodAutoscalerConditions := make([]interface{}, 0)

	if conditions, ok := horizontalPodAutoscaler.Annotations["autoscaling.alpha.kubernetes.io/conditions"]; ok {
		err := json.Unmarshal([]byte(conditions), &horizontalPodAutoscalerConditions)
		if err != nil {
			return nil, err
		}
	}

	object := map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": horizontalPodAutoscalerConditions,
		},
	}

	conditions, err := parseConditions(unstructured.Unstructured{Object: object})
	if err != nil {
		return nil, err
	}
	table := createConditionsTable(conditions, conditionType, hpaConditionColumns)
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
	Config(ctx context.Context, options Options) error
	Status() error
	Metrics(ctx context.Context, options Options) error
	Conditions() error
}

type horizontalPodAutoscalerHandler struct {
	horizontalPodAutoScaler *autoscalingv1.HorizontalPodAutoscaler
	configFunc              func(context.Context, *autoscalingv1.HorizontalPodAutoscaler, Options) (*component.Summary, error)
	statusFunc              func(*autoscalingv1.HorizontalPodAutoscaler) (*component.Summary, error)
	metricsFunc             func(context.Context, *autoscalingv1.MetricStatus, Options) (*component.Summary, error)
	conditionsFunc          func(*autoscalingv1.HorizontalPodAutoscaler) (*component.Table, error)
	object                  *Object
}

// Create creates a horizontalpodautoscaler configuration summary
func (hc *HorizontalPodAutoscalerConfiguration) Create(ctx context.Context, options Options) (*component.Summary, error) {
	if hc.horizontalPodAutoscaler == nil {
		return nil, errors.New("horizontalpodautoscaler is nil")
	}

	hpa := hc.horizontalPodAutoscaler

	sections := component.SummarySections{}

	scaleTarget, err := forScaleTarget(ctx, hpa, &hpa.Spec.ScaleTargetRef, options)
	if err != nil {
		return nil, err
	}

	sections = append(sections, component.SummarySection{
		Header:  "Scale target",
		Content: scaleTarget,
	})

	minReplicas := fmt.Sprintf("%d", *hpa.Spec.MinReplicas)
	maxReplicas := fmt.Sprintf("%d", hpa.Spec.MaxReplicas)
	sections.AddText("Min Replicas", minReplicas)
	sections.AddText("Max Replicas", maxReplicas)

	b := autoscalingv2beta2.HorizontalPodAutoscalerBehavior{}

	if behavior, ok := hpa.Annotations["autoscaling.alpha.kubernetes.io/behavior"]; ok {
		err := json.Unmarshal([]byte(behavior), &b)
		if err != nil {
			return nil, err
		}

		var upPolicies, downPolicies []string
		for _, policy := range b.ScaleUp.Policies {
			p := fmt.Sprintf("%d %s / %d seconds", policy.Value, policy.Type, policy.PeriodSeconds)
			upPolicies = append(upPolicies, p)
		}

		for _, policy := range b.ScaleDown.Policies {
			p := fmt.Sprintf("%d %s / %d seconds", policy.Value, policy.Type, policy.PeriodSeconds)
			downPolicies = append(downPolicies, p)
		}

		cols := component.NewTableCols("Stabilization Window", "Select Policies", "Policies")
		scaleUpTbl := component.NewTableWithRows("", "There are no scale up policies!", cols,
			[]component.TableRow{
				{
					"Stabilization Window": component.NewText(fmt.Sprint(*b.ScaleUp.StabilizationWindowSeconds) + " seconds"),
					"Select Policies":      component.NewText(fmt.Sprint(*b.ScaleUp.SelectPolicy)),
					"Policies":             component.NewText(strings.Join(upPolicies, ", ")),
				},
			})
		sections.Add("Scale Up", scaleUpTbl)

		scaleDownTbl := component.NewTableWithRows("", "There are no scale down policies!", cols,
			[]component.TableRow{
				{
					"Stabilization Window": component.NewText(fmt.Sprint(*b.ScaleDown.StabilizationWindowSeconds) + " seconds"),
					"Select Policies":      component.NewText(fmt.Sprint(*b.ScaleDown.SelectPolicy)),
					"Policies":             component.NewText(strings.Join(downPolicies, ", ")),
				},
			})
		sections.Add("Scale Down", scaleDownTbl)
	}

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

func (h *horizontalPodAutoscalerHandler) Config(ctx context.Context, options Options) error {
	out, err := h.configFunc(ctx, h.horizontalPodAutoScaler, options)
	if err != nil {
		return err
	}

	h.object.RegisterConfig(out)
	return nil
}

func defaultHorizontalPodAutoscalerConfig(ctx context.Context, horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, options Options) (*component.Summary, error) {
	return NewHorizontalPodAutoscalerConfiguration(horizontalPodAutoscaler).Create(ctx, options)
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

func (h *horizontalPodAutoscalerHandler) metrics(ctx context.Context, currentMetrics []autoscalingv1.MetricStatus, options Options) error {
	if h == nil || h.horizontalPodAutoScaler == nil {
		return errors.New("can't display metrics for nil horizontalpodautoscaler")
	}

	for i := range currentMetrics {
		metric := currentMetrics[i]

		if metric.Type == "" {
			continue
		}

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

	_, metricStatus, err := parseAnnotations(*h.horizontalPodAutoScaler)
	if err != nil {
		return errors.New("can't parse annotations for metrics")
	}

	return h.metrics(ctx, metricStatus, options)
}

func defaultHorizontalPodAutoscalerMetrics(ctx context.Context, metricStatus *autoscalingv1.MetricStatus, options Options) (*component.Summary, error) {
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
func forScaleTarget(ctx context.Context, object runtime.Object, scaleTarget *autoscalingv1.CrossVersionObjectReference, options Options) (*component.Link, error) {
	if scaleTarget == nil || object == nil {
		return component.NewLink("", "none", ""), nil
	}

	accessor := meta.NewAccessor()
	ns, err := accessor.Namespace(object)
	if err != nil {
		return component.NewLink("", "none", ""), nil
	}

	key := store.Key{
		Namespace:  ns,
		APIVersion: scaleTarget.APIVersion,
		Kind:       scaleTarget.Kind,
		Name:       scaleTarget.Name,
	}

	objectStore := options.DashConfig.ObjectStore()
	u, err := objectStore.Get(ctx, key)
	if err != nil || u == nil {
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

func getCombinedMetrics(horizontalPodAutoscaler autoscalingv1.HorizontalPodAutoscaler, metricSpec []autoscalingv1.MetricSpec, metricStatus []autoscalingv1.MetricStatus) (string, error) {
	var targets = make(map[string]string)
	var currents = make(map[string]string)

	for _, m := range metricSpec {
		switch m.Type {
		case autoscalingv1.ObjectMetricSourceType:
			if m.Object != nil {
				if m.Object.MetricName != "" && &m.Object.TargetValue != nil {
					targets[m.Object.MetricName] = m.Object.TargetValue.String()
				}
			}
		case autoscalingv1.PodsMetricSourceType:
			if m.Pods != nil {
				if m.Pods.MetricName != "" && &m.Pods.TargetAverageValue != nil {
					target := m.Pods.TargetAverageValue.String()
					targets[m.Pods.MetricName] = target
				}
			}

		case autoscalingv1.ResourceMetricSourceType:
			if m.Resource != nil {
				if m.Resource.Name != "" && m.Resource.TargetAverageValue != nil {
					targets[string(m.Resource.Name)] = m.Resource.TargetAverageValue.String()
				}
			}

		case autoscalingv1.ExternalMetricSourceType:
			if m.External != nil {
				if m.External.MetricName != "" && m.External.TargetAverageValue != nil {
					targets[string(m.External.MetricName)] = m.External.TargetAverageValue.String()
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
		case autoscalingv1.ObjectMetricSourceType:
			if m.Object != nil {
				if m.Object.MetricName != "" {
					currents[m.Object.MetricName] = current
				}
			}
		case autoscalingv1.PodsMetricSourceType:
			if m.Pods != nil {
				if m.Pods.MetricName != "" {
					currents[m.Pods.MetricName] = current
				}
			}
		case autoscalingv1.ResourceMetricSourceType:
			if m.Resource != nil {
				if m.Resource.Name != "" {
					currents[string(m.Resource.Name)] = current
				}
			}
		case autoscalingv1.ExternalMetricSourceType:
			if m.External != nil {
				if m.External.MetricName != "" {
					currents[string(m.External.MetricName)] = current
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

	if horizontalPodAutoscaler.Spec.TargetCPUUtilizationPercentage != nil && horizontalPodAutoscaler.Status.CurrentCPUUtilizationPercentage != nil {
		cpu := fmt.Sprintf("%d/%d", *horizontalPodAutoscaler.Status.CurrentCPUUtilizationPercentage, *horizontalPodAutoscaler.Spec.TargetCPUUtilizationPercentage) + "%"
		result = append(result, cpu)
	}

	return strings.Join(result, ", "), nil
}

func getMetricStatusValue(metricStatus *autoscalingv1.MetricStatus) (string, error) {
	var value string
	if metricStatus == nil {
		return "", errors.New("nil metric status")
	}

	switch metricStatus.Type {
	case autoscalingv1.ObjectMetricSourceType:
		if &metricStatus.Object.CurrentValue != nil {
			value = metricStatus.Object.CurrentValue.String()
		}
	case autoscalingv1.PodsMetricSourceType:
		if &metricStatus.Pods.CurrentAverageValue != nil {
			value = metricStatus.Pods.CurrentAverageValue.String()
		}
	case autoscalingv1.ResourceMetricSourceType:
		if &metricStatus.Resource.CurrentAverageValue != nil {
			value = metricStatus.Resource.CurrentAverageValue.String()
		}
	case autoscalingv1.ExternalMetricSourceType:
		if metricStatus.External.CurrentAverageValue != nil {
			value = metricStatus.External.CurrentAverageValue.String()
		}
	}

	return value, nil
}

func parseAnnotations(horizontalPodAutoscaler autoscalingv1.HorizontalPodAutoscaler) ([]autoscalingv1.MetricSpec, []autoscalingv1.MetricStatus, error) {
	horizontalPodAutoscalerMetrics := make([]autoscalingv1.MetricSpec, 0)
	horizontalPodAutoscalerCurrentMetrics := make([]autoscalingv1.MetricStatus, 0)

	if metrics, ok := horizontalPodAutoscaler.Annotations["autoscaling.alpha.kubernetes.io/metrics"]; ok {
		err := json.Unmarshal([]byte(metrics), &horizontalPodAutoscalerMetrics)
		if err != nil {
			return nil, nil, err
		}
	}
	if currentMetrics, ok := horizontalPodAutoscaler.Annotations["autoscaling.alpha.kubernetes.io/current-metrics"]; ok {
		err := json.Unmarshal([]byte(currentMetrics), &horizontalPodAutoscalerCurrentMetrics)
		if err != nil {
			return nil, nil, err
		}
	}

	return horizontalPodAutoscalerMetrics, horizontalPodAutoscalerCurrentMetrics, nil
}
