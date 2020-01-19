/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package workloads

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// HomeDescriberOption is an option for configuring home describer.
type HomeDescriberOption func(hd *HomeDescriber)

// HomeDescriber describe the home page for workloads module.
type HomeDescriber struct{}

var _ describer.Describer = (*HomeDescriber)(nil)

// NewHomeDescriber creates an instance of HomeDescriber.
func NewHomeDescriber(options ...HomeDescriberOption) (*HomeDescriber, error) {
	d := &HomeDescriber{}

	for _, option := range options {
		option(d)
	}

	return d, nil
}

func loadCards(ctx context.Context, namespace string, options describer.Options) ([]*component.Card, bool, error) {
	pml, err := octant.NewClusterPodMetricsLoader(options.Dash.ClusterClient())
	if err != nil {
		return nil, false, fmt.Errorf("create pod metrics loader")
	}

	loader, err := octant.NewClusterWorkloadLoader(options.Dash.ObjectStore(), pml)
	if err != nil {
		return nil, false, fmt.Errorf("create workload loader")
	}

	collector, err := octant.NewWorkloadCardCollector(loader)
	if err != nil {
		return nil, false, fmt.Errorf("create card collector")
	}

	cards, fullMetrics, err := collector.Collect(ctx, namespace)
	if err != nil {
		return nil, false, fmt.Errorf("collect workload cards: %w", err)

	}

	return cards, fullMetrics, nil
}

// Describe creates a content response for workloads.
func (h *HomeDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {

	cards, fullMetrics, err := loadCards(ctx, namespace, options)
	if err != nil {
		return component.EmptyContentResponse, fmt.Errorf("load cards cards: %w", err)
	}

	cardWidth := component.WidthHalf
	if !fullMetrics {
		cardWidth = component.WidthQuarter
	}

	layout := component.NewFlexLayout("Workloads")
	section := component.FlexLayoutSection{}

	for _, card := range cards {
		section = append(section, component.FlexLayoutItem{
			Width: cardWidth,
			View:  card,
		})
	}

	layout.AddSections(section)

	cr := component.ContentResponse{
		Title:      component.TitleFromString("Workloads"),
		Components: []component.Component{layout},
		IconName:   "",
		IconSource: "",
	}

	return cr, nil
}

// PathFilters returns a path filter for the root path.
func (h *HomeDescriber) PathFilters() []describer.PathFilter {
	return []describer.PathFilter{
		*describer.NewPathFilter("/", h),
	}
}

// Reset is a no-op.
func (h HomeDescriber) Reset(ctx context.Context) error {
	return nil
}
