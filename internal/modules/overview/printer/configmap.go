package printer

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/pkg/view/component"

	corev1 "k8s.io/api/core/v1"
)

// ConfigMapListHandler is a printFunc that prints ConfigMaps
func ConfigMapListHandler(_ context.Context, list *corev1.ConfigMapList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	// Data column
	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	tbl := component.NewTable("ConfigMaps", cols)

	for _, c := range list.Items {
		row := component.TableRow{}

		nameLink, err := opts.Link.ForObject(&c, c.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink

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
func ConfigMapHandler(ctx context.Context, cm *corev1.ConfigMap, options Options) (component.Component, error) {
	o := NewObject(cm)

	summary, err := describeConfigMapConfig(cm)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(summary)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return describeConfigMapData(cm)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, options)
}

// ConfigMapConfiguration generates config map configuration
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
