/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/vmware/octant/pkg/view/component"
)

func RoleListHandler(_ context.Context, roleList *rbacv1.RoleList, options Options) (component.Component, error) {
	if roleList == nil {
		return nil, errors.New("role list is nil")
	}

	columns := component.NewTableCols("Name", "Age")
	table := component.NewTable("Roles", "We couldn't find any roles!", columns)

	for _, role := range roleList.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&role, role.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Age"] = component.NewTimestamp(role.CreationTimestamp.Time)
		table.Add(row)
	}

	return table, nil
}

func RoleHandler(ctx context.Context, role *rbacv1.Role, opts Options) (component.Component, error) {
	o := NewObject(role)

	configSummary, err := printRoleConfig(role)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(configSummary)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return printRolePolicyRules(role)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, opts)
}

func printRoleConfig(role *rbacv1.Role) (*component.Summary, error) {
	if role == nil {
		return nil, errors.New("role is nil")
	}

	sections := component.SummarySections{}
	sections.AddText("Name", role.Name)
	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func printRolePolicyRules(role *rbacv1.Role) (*component.Table, error) {
	if role == nil {
		return nil, errors.New("role is nil")
	}

	rules := role.Rules
	var breakdownRules []rbacv1.PolicyRule
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
	tbl := component.NewTable("PolicyRules", "There are no policy rules!", cols)

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
