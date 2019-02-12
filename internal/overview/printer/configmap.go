package printer

import (
	"fmt"
	"sort"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

// ConfigMapListHandler is a printFunc that prints ConfigMaps
func ConfigMapListHandler(list *corev1.ConfigMapList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	// Data column
	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	tbl := component.NewTable("ConfigMaps", cols)

	for _, c := range list.Items {
		row := component.TableRow{}
		configmapPath := gvkPath(c.TypeMeta.APIVersion, c.TypeMeta.Kind, c.Name)
		row["Name"] = component.NewLink("", c.Name, configmapPath)
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
func ConfigMapHandler(cm *corev1.ConfigMap, options Options) (component.ViewComponent, error) {
	fl := flexlayout.New()

	configSection := fl.AddSection()

	cmConfigGen := NewConfigMapConfiguration(cm)
	configView, err := cmConfigGen.Create()
	if err != nil {
		return nil, err
	}

	if err := configSection.Add(configView, 16); err != nil {
		return nil, errors.Wrap(err, "add configmap config to layout")
	}

	dataSection := fl.AddSection()
	dataTable, err := cmConfigGen.describeConfigMapData()
	if err != nil {
		return nil, err
	}

	if err := dataSection.Add(dataTable, 16); err != nil {
		return nil, errors.Wrap(err, "add configmap data to layout")
	}

	view := fl.ToComponent("Summary")

	return view, nil
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
func (cc *ConfigMapConfiguration) Create() (*component.Summary, error) {
	if cc == nil || cc.configmap == nil {
		return nil, errors.New("configmap is nil")
	}

	cm := cc.configmap

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
func (cc *ConfigMapConfiguration) describeConfigMapData() (*component.Table, error) {
	if cc == nil || cc.configmap == nil {
		return nil, errors.New("configmap is nil")
	}

	cm := cc.configmap

	cols := component.NewTableCols("Key", "Value")
	tbl := component.NewTable("Data", cols)

	tbl.Add(describeConfigMapDataRows(cm)...)

	return tbl, nil
}

// describeDataRows prints key value pairs from data
func describeConfigMapDataRows(cm *corev1.ConfigMap) []component.TableRow {
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

	return rows
}
