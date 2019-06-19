/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"sort"

	"github.com/pkg/errors"

	"github.com/vmware/octant/pkg/view/component"

	rbacv1 "k8s.io/api/rbac/v1"
)

// ClusterRoleListHandler is a printFunc that prints cluster roles
func ClusterRoleListHandler(_ context.Context, list *rbacv1.ClusterRoleList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("cluster role list is nil")
	}

	cols := component.NewTableCols("Name", "Age")
	tbl := component.NewTable("Cluster Roles", cols)

	for _, clusterRole := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&clusterRole, clusterRole.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		ts := clusterRole.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// ClusterRoleHandler is a printFunc that prints a cluster role
func ClusterRoleHandler(ctx context.Context, clusterRole *rbacv1.ClusterRole, options Options) (component.Component, error) {
	o := NewObject(clusterRole)

	summary, err := printClusterRoleConfig(clusterRole)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(summary)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return printClusterRolePolicyRules(clusterRole)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, options)
}

func printClusterRoleConfig(clusterRole *rbacv1.ClusterRole) (*component.Summary, error) {
	if clusterRole == nil {
		return nil, errors.New("cluster role is nil")
	}

	var sections component.SummarySections

	if clusterRoleAggregation := clusterRole.AggregationRule; clusterRoleAggregation != nil {
		if clusterRoleSelectors := clusterRoleAggregation.ClusterRoleSelectors; clusterRoleSelectors != nil {
			var selectors []component.Selector

			for _, selector := range clusterRoleSelectors {
				for k, v := range selector.MatchLabels {
					ls := component.NewLabelSelector(k, v)
					selectors = append(selectors, ls)
				}

				sections = append(sections, component.SummarySection{
					Header:  "Selectors",
					Content: component.NewSelectors(selectors),
				})
			}
		}
	}

	sections.AddText("Name", clusterRole.Name)
	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func printClusterRolePolicyRules(clusterRole *rbacv1.ClusterRole) (*component.Table, error) {
	if clusterRole == nil {
		return nil, errors.New("cluster role is nil")
	}

	return printPolicyRules(clusterRole.Rules)
}

func printPolicyRules(rules []rbacv1.PolicyRule) (*component.Table, error) {
	breakdownRules := []rbacv1.PolicyRule{}
	for _, rule := range rules {
		breakdownRules = append(breakdownRules, BreakdownRule(rule)...)
	}

	rules, err := compactRules(breakdownRules)
	if err != nil {
		return nil, errors.New("cannot compact rules")
	}

	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].String() < rules[j].String()
	})

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	tbl := component.NewTable("Policy Rules", cols)

	for _, r := range rules {
		row := component.TableRow{}
		row["Resources"] = component.NewText(CombineResourceGroup(r.Resources, r.APIGroups))
		row["Non-Resource URLs"] = component.NewText(printSlice(r.NonResourceURLs))
		row["Resource Names"] = component.NewText(printSlice(r.ResourceNames))
		row["Verbs"] = component.NewText(printSlice(r.Verbs))

		tbl.Add(row)
	}

	return tbl, nil
}
