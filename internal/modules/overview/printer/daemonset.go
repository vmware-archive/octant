/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"

	"github.com/heptio/developer-dash/pkg/view/component"
)

// DaemonSetListHandler is a printFunc that lists daemon sets
func DaemonSetListHandler(_ context.Context, list *appsv1.DaemonSetList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("daemon set list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Ready",
		"Up-To-Date", "Age", "Node Selector")
	table := component.NewTable("Daemon Sets", cols)

	for _, daemonSet := range list.Items {
		row := component.TableRow{}
		nameLink, err := opts.Link.ForObject(&daemonSet, daemonSet.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(daemonSet.Labels)
		row["Desired"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.DesiredNumberScheduled))
		row["Current"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.CurrentNumberScheduled))
		row["Ready"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.NumberReady))
		row["Up-To-Date"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.UpdatedNumberScheduled))
		row["Age"] = component.NewTimestamp(daemonSet.ObjectMeta.CreationTimestamp.Time)
		row["Node Selector"] = printSelectorMap(daemonSet.Spec.Template.Spec.NodeSelector)

		table.Add(row)
	}

	return table, nil
}

// DaemonSetHandler is a printFunc that prints a daemon set
func DaemonSetHandler(ctx context.Context, daemonSet *appsv1.DaemonSet, options Options) (component.Component, error) {
	o := NewObject(daemonSet)

	configSummary, err := printDaemonSetConfig(daemonSet)
	if err != nil {
		return nil, err
	}

	statusSummary, err := printDaemonSetStatus(daemonSet)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(configSummary)
	o.RegisterSummary(statusSummary)
	o.EnablePodTemplate(daemonSet.Spec.Template)
	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return createPodListView(ctx, daemonSet, options)
		},
		Width: component.WidthFull,
	})
	o.EnableEvents()

	return o.ToComponent(ctx, options)
}

func printDaemonSetConfig(daemonSet *appsv1.DaemonSet) (*component.Summary, error) {
	if daemonSet == nil {
		return nil, errors.New("daemon set is nil")
	}

	var sections component.SummarySections

	rollingUpdate := daemonSet.Spec.UpdateStrategy.RollingUpdate
	if rollingUpdate != nil {
		rollingUpdateText := fmt.Sprintf("Max Unavailable %s",
			rollingUpdate.MaxUnavailable.String(),
		)
		sections.AddText("Update Strategy", rollingUpdateText)
	}

	if historyLimit := daemonSet.Spec.RevisionHistoryLimit; historyLimit != nil {
		sections.AddText("Revision History Limit", fmt.Sprint(*historyLimit))
	}

	if selector := daemonSet.Spec.Selector; selector != nil {
		sections.Add("Selectors", printSelector(selector))
	}

	if selector := daemonSet.Spec.Template.Spec.NodeSelector; selector != nil {
		sections.Add("Node Selectors", printSelectorMap(selector))
	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func printDaemonSetStatus(daemonSet *appsv1.DaemonSet) (*component.Summary, error) {
	if daemonSet == nil {
		return nil, errors.New("daemon set is nil")
	}

	var sections component.SummarySections

	status := daemonSet.Status
	sections.AddText("Current Number Scheduled", fmt.Sprint(status.CurrentNumberScheduled))
	sections.AddText("Desired Number Scheduled", fmt.Sprint(status.DesiredNumberScheduled))
	sections.AddText("Number Available", fmt.Sprint(status.NumberAvailable))
	sections.AddText("Number Mis-scheduled", fmt.Sprint(status.NumberMisscheduled))
	sections.AddText("Number Ready", fmt.Sprint(status.NumberReady))
	sections.AddText("Updated Number Scheduled", fmt.Sprint(status.UpdatedNumberScheduled))

	summary := component.NewSummary("Status", sections...)

	return summary, nil
}
