/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// RoleListHandler is a printFunc that prints roles
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

// RoleHandler is a printFunc that prints roles
func RoleHandler(ctx context.Context, role *rbacv1.Role, options Options) (component.Component, error) {
	o := NewObject(role)

	rh, err := newRoleHandler(role, o)
	if err != nil {
		return nil, err
	}

	if err := rh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print role configuration")
	}

	if err := rh.PolicyRules(options); err != nil {
		return nil, errors.Wrap(err, "print role policy rules")
	}

	return o.ToComponent(ctx, options)
}

// RoleConfiguration generates a role configuration
type RoleConfiguration struct {
	role *rbacv1.Role
}

// NewRoleConfiguration creates an instance of RoleConfiguration
func NewRoleConfiguration(role *rbacv1.Role) *RoleConfiguration {
	return &RoleConfiguration{
		role: role,
	}
}

// Create creates a role configuration summary
func (r *RoleConfiguration) Create(options Options) (*component.Summary, error) {
	if r.role == nil {
		return nil, errors.New("role is nil")
	}
	role := r.role

	sections := component.SummarySections{}
	sections.AddText("Name", role.Name)
	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func createRolePolicyRulesView(role *rbacv1.Role) (*component.Table, error) {
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

type roleObject interface {
	Config(options Options) error
	PolicyRules(options Options) error
}

type roleHandler struct {
	role            *rbacv1.Role
	configFunc      func(*rbacv1.Role, Options) (*component.Summary, error)
	policyRulesFunc func(*rbacv1.Role, Options) (*component.Table, error)
	object          *Object
}

func newRoleHandler(role *rbacv1.Role, object *Object) (*roleHandler, error) {
	if role == nil {
		return nil, errors.New("can't print a nil role")
	}

	if object == nil {
		return nil, errors.New("can't print role using a nil object printer")
	}

	rh := &roleHandler{
		role:            role,
		configFunc:      defaultRoleConfig,
		policyRulesFunc: defaultRolePolicyRules,
		object:          object,
	}
	return rh, nil
}

func (r *roleHandler) Config(options Options) error {
	out, err := r.configFunc(r.role, options)
	if err != nil {
		return err
	}
	r.object.RegisterConfig(out)
	return nil
}

func defaultRoleConfig(role *rbacv1.Role, options Options) (*component.Summary, error) {
	return NewRoleConfiguration(role).Create(options)
}

func (r *roleHandler) PolicyRules(options Options) error {
	if r.role == nil {
		return errors.New("can't display policy rules for nil role")
	}

	r.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return r.policyRulesFunc(r.role, options)
		},
	})

	return nil
}

func defaultRolePolicyRules(role *rbacv1.Role, options Options) (*component.Table, error) {
	return createRolePolicyRulesView(role)
}
