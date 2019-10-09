/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
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
	tbl := component.NewTable("Cluster Roles", "We couldn't find any cluster roles!", cols)

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

	ch, err := newClusterRoleHandler(clusterRole, o)
	if err != nil {
		return nil, err
	}

	if err := ch.Config(options); err != nil {
		return nil, errors.Wrap(err, "print clusterrole configuration")
	}

	if err := ch.PolicyRules(options); err != nil {
		return nil, errors.Wrap(err, "print clusterrole policy rules")
	}

	return o.ToComponent(ctx, options)
}

// ClusterRoleConfiguration generates a clusterrole configuration
type ClusterRoleConfiguration struct {
	clusterRole *rbacv1.ClusterRole
}

// NewClusterRoleConfiguration creates an instance of ClusterRoleConfiguration
func NewClusterRoleConfiguration(clusterRole *rbacv1.ClusterRole) *ClusterRoleConfiguration {
	return &ClusterRoleConfiguration{
		clusterRole: clusterRole,
	}
}

// Create creates a clusterrole configuration summary
func (c *ClusterRoleConfiguration) Create(options Options) (*component.Summary, error) {
	if c == nil || c.clusterRole == nil {
		return nil, errors.New("clusterrole is nil")
	}

	var sections component.SummarySections

	if clusterRoleAggregation := c.clusterRole.AggregationRule; clusterRoleAggregation != nil {
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

	sections.AddText("Name", c.clusterRole.Name)
	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func createClusterRolePolicyRulesView(clusterRole *rbacv1.ClusterRole) (*component.Table, error) {
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
	tbl := component.NewTable("Policy Rules", "There are no policy rules!", cols)

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

type clusterRoleObject interface {
	Config(options Options) error
	PolicyRules(options Options) error
}

type clusterRoleHandler struct {
	clusterRole     *rbacv1.ClusterRole
	configFunc      func(*rbacv1.ClusterRole, Options) (*component.Summary, error)
	policyRulesFunc func(*rbacv1.ClusterRole, Options) (*component.Table, error)
	object          *Object
}

var _ clusterRoleObject = (*clusterRoleHandler)(nil)

func newClusterRoleHandler(clusterRole *rbacv1.ClusterRole, object *Object) (*clusterRoleHandler, error) {
	if clusterRole == nil {
		return nil, errors.New("can't print a nil clusterrole")
	}

	if object == nil {
		return nil, errors.New("can't print a clusterrole using an nil object printer")
	}

	ch := &clusterRoleHandler{
		clusterRole:     clusterRole,
		configFunc:      defaultClusterRoleConfig,
		policyRulesFunc: defaultClusterRolePolicyRules,
		object:          object,
	}
	return ch, nil
}

func (c *clusterRoleHandler) Config(options Options) error {
	out, err := c.configFunc(c.clusterRole, options)
	if err != nil {
		return err
	}
	c.object.RegisterConfig(out)
	return nil
}

func defaultClusterRoleConfig(clusterRole *rbacv1.ClusterRole, options Options) (*component.Summary, error) {
	return NewClusterRoleConfiguration(clusterRole).Create(options)
}

func (c *clusterRoleHandler) PolicyRules(options Options) error {
	if c.clusterRole == nil {
		return errors.New("can't display policy rules for nil clusterrole")
	}

	c.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return c.policyRulesFunc(c.clusterRole, options)
		},
	})

	return nil
}

func defaultClusterRolePolicyRules(clusterRole *rbacv1.ClusterRole, options Options) (*component.Table, error) {
	return createClusterRolePolicyRulesView(clusterRole)
}
