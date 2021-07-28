/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

package workloads

import (
	"context"
	"fmt"
	"path"
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/resourceviewer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type DetailDescriber struct{}

var _ describer.Describer = (*DetailDescriber)(nil)

func NewDetailDescriber() (*DetailDescriber, error) {
	d := &DetailDescriber{}

	return d, nil
}

func (d *DetailDescriber) loadWorkloads(ctx context.Context, namespace string, options describer.Options) ([]octant.Workload, error) {
	pml, err := octant.NewClusterPodMetricsLoader(options.Dash.ClusterClient())
	if err != nil {
		return nil, fmt.Errorf("create pod metrics loader")
	}

	loader, err := octant.NewClusterWorkloadLoader(options.Dash.ObjectStore(), pml)
	if err != nil {
		return nil, fmt.Errorf("create workload loader")
	}

	workloads, err := loader.Load(ctx, namespace, options.Link)
	if err != nil {
		return nil, fmt.Errorf("load workloads")
	}

	return workloads, nil
}

func (d *DetailDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	name := options.Fields["name"]

	workloads, err := d.loadWorkloads(ctx, namespace, options)
	if err != nil {
		cr := d.createResponse(
			component.NewError(component.TitleFromString("Unable to load workloads"), err),
		)
		return cr, nil
	}

	found := false
	var cur octant.Workload
	for _, workload := range workloads {
		if workload.Name != name {
			continue
		}

		found = true
		cur = workload
	}

	if !found {
		cr := d.createResponse(
			component.NewError(component.TitleFromString("Workload not found"),
				fmt.Errorf("unable to find workload %s", name)),
		)
		return cr, nil
	}

	layout := component.NewFlexLayout("Workload layout")

	summary, err := octant.CreateWorkloadSummary(&cur, component.DonutChartSizeSmall)
	if err != nil {
		cr := d.createResponse(
			component.NewError(component.TitleFromString("Create summary for workload"), err),
		)
		return cr, nil
	}

	workloadName := fmt.Sprintf(`
### %s
_%s_
`, cur.Name, cur.Owner.GroupVersionKind())

	headerSection := component.FlexLayoutSection{
		{
			Width: component.WidthQuarter,
			View:  component.NewMarkdownText(workloadName),
		},
		{
			Width: component.WidthQuarter,
			View:  summary.Summary,
		},
	}

	if summary.MetricsEnabled {
		headerSection = append(headerSection, []component.FlexLayoutItem{
			{
				Width: component.WidthQuarter,
				View:  summary.Memory,
			},
			{
				Width: component.WidthQuarter,
				View:  summary.CPU,
			},
		}...)
	}

	var objects []*unstructured.Unstructured
	pods := cur.Pods().Items

	for i := 0; i < len(pods); i++ {
		objects = append(objects, &pods[i])
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].GetName() < objects[j].GetName()
	})

	var selection string
	for _, obj := range objects {
		if len(obj.GetOwnerReferences()) > 0 {
			for _, ref := range obj.GetOwnerReferences() {
				selection = fmt.Sprintf("%s pods", ref.Name)
				continue
			}
		} else {
			selection = string(obj.GetUID())
		}
	}

	rv, err := resourceviewer.Create(ctx, options.Dash, options.Queryer, selection, objects...)
	if err != nil {
		cr := d.createResponse(
			component.NewError(component.TitleFromString("Unable to create resource viewer"), err),
		)
		return cr, nil
	}

	viewerSection := component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  rv,
		},
	}
	layout.AddSections(headerSection, viewerSection)

	cr := d.createResponse(
		layout,
	)

	return cr, nil
}

func (d *DetailDescriber) createResponse(components ...component.Component) component.ContentResponse {
	cr := component.ContentResponse{
		Title:      component.TitleFromString(""),
		Components: components,
	}

	return cr
}

func (d *DetailDescriber) PathFilters() []describer.PathFilter {
	return []describer.PathFilter{
		*describer.NewPathFilter(path.Join("/detail", describer.ResourceNameRegex), d),
	}
}

func (d DetailDescriber) Reset(_ context.Context) error {
	return nil
}
