package printer

import (
	"context"
	"fmt"
	"sort"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

// ConfigMapListHandler is a printFunc that prints ConfigMaps
func ConfigMapListHandler(ctx context.Context, list *corev1.ConfigMapList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	// Data column
	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	tbl := component.NewTable("ConfigMaps", cols)

	for _, c := range list.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&c, c.Name)
		row["Labels"] = component.NewLabels(c.Labels)

		data := fmt.Sprintf("%d", len(c.Data))
		row["Data"] = component.NewText(data)

		ts := c.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// ConfigMapHandler is a printFunc that prints a ConfigMap
func ConfigMapHandler(ctx context.Context, cm *corev1.ConfigMap, options Options) (component.ViewComponent, error) {
	o := NewObject(cm)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return describeConfigMapConfig(cm)
	}, 16)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return describeConfigMapData(cm)
		},
		Width: 24,
	})

	return o.ToComponent(ctx, options)
}

// ConfigMapConfiguration generates configmap configuration
type ConfigMapConfiguration struct {
	configmap *corev1.ConfigMap
}

// NewConfigMapConfiguration creates an instance of ConfigMapConfiguration
func NewConfigMapConfiguration(cm *corev1.ConfigMap) *ConfigMapConfiguration {
	return &ConfigMapConfiguration{
		configmap: cm,
	}
}

// Create a configmap configuration summary
func describeConfigMapConfig(cm *corev1.ConfigMap) (*component.Summary, error) {
	if cm == nil {
		return nil, errors.New("config map is nil")
	}

	sections := component.SummarySections{}

	creationTimestamp := cm.CreationTimestamp.Time
	sections = append(sections, component.SummarySection{
		Header:  "Age",
		Content: component.NewTimestamp(creationTimestamp),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

// describeDataTable returns a table containing configmap data
func describeConfigMapData(cm *corev1.ConfigMap) (*component.Table, error) {
	if cm == nil {
		return nil, errors.New("config map is nil")
	}

	cols := component.NewTableCols("Key", "Value")
	tbl := component.NewTable("Data", cols)

	rows, err := describeConfigMapDataRows(cm)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		tbl.Add(row)
	}

	return tbl, nil
}

// describeDataRows prints key value pairs from data
func describeConfigMapDataRows(cm *corev1.ConfigMap) ([]component.TableRow, error) {
	if cm == nil {
		return nil, errors.New("config map is nil")
	}

	rows := make([]component.TableRow, 0)
	data := cm.Data

	// Alpha sort keys so that output is consistent
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		row := component.TableRow{}
		rows = append(rows, row)

		row["Key"] = component.NewText(k)

		row["Value"] = component.NewText(data[k])
	}

	return rows, nil
}
