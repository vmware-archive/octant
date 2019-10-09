/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"

	"github.com/vmware/octant/pkg/view/component"
)

type emptyComponent struct{}

var _ component.Component = (*emptyComponent)(nil)

func (c *emptyComponent) GetMetadata() component.Metadata {
	return component.Metadata{
		Type: "empty",
	}
}

func (c *emptyComponent) SetAccessor(string) {
	// no-op
}

func (c *emptyComponent) IsEmpty() bool {
	return true
}

func (c *emptyComponent) String() string {
	return ""
}

func (c *emptyComponent) LessThan(interface{}) bool {
	return false
}

func (c emptyComponent) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	return json.Marshal(m)
}

func createPodTable(pods ...corev1.Pod) *component.Table {
	tableCols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable("/v1, Kind=PodList", "placeholder", tableCols)
	for _, pod := range pods {
		table.Add(component.TableRow{
			"Age":    component.NewTimestamp(pod.CreationTimestamp.Time),
			"Labels": component.NewLabels(pod.Labels),
			"Name":   component.NewText(pod.Name),
		})
	}

	return table
}

func podListType() interface{} {
	return &corev1.PodList{}
}

func podObjectType() interface{} {
	return &corev1.Pod{}
}
