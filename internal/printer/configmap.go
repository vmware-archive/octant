/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/view/component"

	corev1 "k8s.io/api/core/v1"
)

// ConfigMapListHandler is a printFunc that prints ConfigMaps
func ConfigMapListHandler(_ context.Context, list *corev1.ConfigMapList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	// Data column
	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	tbl := component.NewTable("ConfigMaps", "We couldn't find any config maps!", cols)

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

	ch, err := newConfigMapHandler(cm, o)
	if err != nil {
		return nil, err
	}

	if err := ch.Config(options); err != nil {
		return nil, errors.Wrap(err, "print configmap configuration")
	}

	if err := ch.Data(options); err != nil {
		return nil, errors.Wrap(err, "print configmap data")
	}

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
func (c *ConfigMapConfiguration) Create(options Options) (*component.Summary, error) {
	if c.configmap == nil {
		return nil, errors.New("config map is nil")
	}
	configMap := c.configmap

	sections := component.SummarySections{}

	creationTimestamp := configMap.CreationTimestamp.Time
	sections = append(sections, component.SummarySection{
		Header:  "Age",
		Content: component.NewTimestamp(creationTimestamp),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

// describeDataTable returns a table containing config map data
func describeConfigMapData(cm *corev1.ConfigMap) (*component.Table, error) {
	if cm == nil {
		return nil, errors.New("config map is nil")
	}

	cols := component.NewTableCols("Key", "Value")
	table := component.NewTable("Data", "No data has been configured for this config map!", cols)

	rows, err := describeConfigMapDataRows(cm)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		table.Add(row)
	}

	table.Sort("Key", false)

	return table, nil
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

		if strings.Contains(data[k], "\n") {
			row["Value"] = component.NewCodeBlock(data[k])
		} else {
			row["Value"] = component.NewText(data[k])
		}
	}

	return rows, nil
}

type configMapObject interface {
	Config(options Options) error
	Data(option Options) error
}

type configMapHandler struct {
	configMap  *corev1.ConfigMap
	configFunc func(*corev1.ConfigMap, Options) (*component.Summary, error)
	dataFunc   func(*corev1.ConfigMap, Options) (*component.Table, error)
	object     *Object
}

var _ configMapObject = (*configMapHandler)(nil)

func newConfigMapHandler(configMap *corev1.ConfigMap, object *Object) (*configMapHandler, error) {
	if configMap == nil {
		return nil, errors.New("can't print a nil configmap")
	}

	if object == nil {
		return nil, errors.New("can't print configmap using a nil object printer")
	}

	ch := &configMapHandler{
		configMap:  configMap,
		configFunc: defaultConfigMapConfig,
		dataFunc:   defaultConfigMapData,
		object:     object,
	}

	return ch, nil
}

func (c *configMapHandler) Config(options Options) error {
	out, err := c.configFunc(c.configMap, options)
	if err != nil {
		return err
	}
	c.object.RegisterConfig(out)
	return nil
}

func defaultConfigMapConfig(configMap *corev1.ConfigMap, options Options) (*component.Summary, error) {
	return NewConfigMapConfiguration(configMap).Create(options)
}

func (c *configMapHandler) Data(options Options) error {
	if c.configMap == nil {
		return errors.New("can't display data for nil configmap")
	}

	c.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return c.dataFunc(c.configMap, options)
		},
	})

	return nil
}

func defaultConfigMapData(configMap *corev1.ConfigMap, options Options) (*component.Table, error) {
	return describeConfigMapData(configMap)
}
