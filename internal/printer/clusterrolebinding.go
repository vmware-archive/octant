/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/vmware/octant/pkg/view/component"
)

// ClusterRoleBindingListHandler is a printFunc that prints ClusterRoldBindings
func ClusterRoleBindingListHandler(_ context.Context, clusterRoleBindingList *rbacv1.ClusterRoleBindingList, options Options) (component.Component, error) {
	if clusterRoleBindingList == nil {
		return nil, errors.New("cluster role binding list is nil")
	}

	columns := component.NewTableCols("Name", "Labels", "Age", "Role kind", "Role name")
	table := component.NewTable("Cluster Role Bindings", "We couldn't find any cluster role bindings!", columns)

	for _, roleBinding := range clusterRoleBindingList.Items {
		row := component.TableRow{}

		nameLink, err := options.Link.ForObject(&roleBinding, roleBinding.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink

		row["Labels"] = component.NewLabels(roleBinding.Labels)
		row["Age"] = component.NewTimestamp(roleBinding.CreationTimestamp.Time)
		row["Role kind"] = component.NewText(roleBinding.RoleRef.Kind)

		roleName, err := roleLinkFromClusterRoleBinding(&roleBinding, options)
		if err != nil {
			return nil, err
		}

		row["Role name"] = roleName

		table.Add(row)
	}

	return table, nil
}

func roleLinkFromClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding, options Options) (*component.Link, error) {
	roleRef := clusterRoleBinding.RoleRef

	namespace := ""
	if roleRef.Kind == "Role" {
		namespace = clusterRoleBinding.Namespace
	}

	apiVersion := fmt.Sprintf("%s/%s", roleRef.APIGroup, "v1")
	return options.Link.ForGVK(namespace, apiVersion, roleRef.Kind, roleRef.Name, roleRef.Name)
}

// ClusterRoleBindingHandler is a printFunc that prints a ClusterRoleBinding
func ClusterRoleBindingHandler(ctx context.Context, clusterRoleBinding *rbacv1.ClusterRoleBinding, options Options) (component.Component, error) {
	o := NewObject(clusterRoleBinding)

	ch, err := newClusterRoleBindingHandler(clusterRoleBinding, o)
	if err != nil {
		return nil, err
	}

	if err := ch.Config(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print clusterrolebinding configuration")
	}

	if err := ch.Subjects(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print clusterrolebinding subjects")
	}

	return o.ToComponent(ctx, options)
}

// ClusterRoleBindingConfiguration generates a clusterrolebinding configuration
type ClusterRoleBindingConfiguration struct {
	clusterRoleBinding *rbacv1.ClusterRoleBinding
}

// NewClusterRoleBindingConfiguration creates an instance of ClusterRoleBindingConfiguration
func NewClusterRoleBindingConfiguration(clusterRoleBinding *rbacv1.ClusterRoleBinding) *ClusterRoleBindingConfiguration {
	return &ClusterRoleBindingConfiguration{
		clusterRoleBinding: clusterRoleBinding,
	}
}

// Create creates a clusterrolebinding configuration summary
func (c *ClusterRoleBindingConfiguration) Create(ctx context.Context, options Options) (*component.Summary, error) {
	if c == nil || c.clusterRoleBinding == nil {
		return nil, errors.New("clusterrolebinding is nil")
	}

	sections := component.SummarySections{}

	sections.AddText("Role kind", c.clusterRoleBinding.RoleRef.Kind)

	roleName, err := roleLinkFromClusterRoleBinding(c.clusterRoleBinding, options)
	if err != nil {
		return nil, err
	}

	sections.Add("Role name", roleName)

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func createClusterRoleBindingSubjectsView(clusterRoleBinding *rbacv1.ClusterRoleBinding) (component.Component, error) {
	if clusterRoleBinding == nil {
		return nil, errors.New("cluster role binding is nil")
	}

	columns := component.NewTableCols("Kind", "Name", "Namespace")
	table := component.NewTable("Subjects", "There are no subjects!", columns)

	for _, subject := range clusterRoleBinding.Subjects {
		row := component.TableRow{}
		row["Kind"] = component.NewText(subject.Kind)
		row["Name"] = component.NewText(subject.Name)
		row["Namespace"] = component.NewText(subject.Namespace)

		table.Add(row)
	}

	return table, nil
}

type clusterRoleBindingObject interface {
	Config(ctx context.Context, options Options) error
	Subjects(ctx context.Context, options Options) error
}

type clusterRoleBindingHandler struct {
	clusterRoleBinding *rbacv1.ClusterRoleBinding
	configFunc         func(context.Context, *rbacv1.ClusterRoleBinding, Options) (*component.Summary, error)
	subjectsFunc       func(context.Context, *rbacv1.ClusterRoleBinding, Options) (component.Component, error)
	object             *Object
}

var _ clusterRoleBindingObject = (*clusterRoleBindingHandler)(nil)

func newClusterRoleBindingHandler(clusterRoleBinding *rbacv1.ClusterRoleBinding, object *Object) (*clusterRoleBindingHandler, error) {
	if clusterRoleBinding == nil {
		return nil, errors.New("can't print a nil rolebinding")
	}

	if object == nil {
		return nil, errors.New("can't print rolebinding using a nil object printer")
	}

	ch := &clusterRoleBindingHandler{
		clusterRoleBinding: clusterRoleBinding,
		configFunc:         defaultClusterRoleBindingConfig,
		subjectsFunc:       defaultClusterRoleBindingSubjects,
		object:             object,
	}

	return ch, nil
}

func (c *clusterRoleBindingHandler) Config(ctx context.Context, options Options) error {
	out, err := c.configFunc(ctx, c.clusterRoleBinding, options)
	if err != nil {
		return err
	}
	c.object.RegisterConfig(out)
	return nil
}

func defaultClusterRoleBindingConfig(ctx context.Context, clusterRoleBinding *rbacv1.ClusterRoleBinding, options Options) (*component.Summary, error) {
	return NewClusterRoleBindingConfiguration(clusterRoleBinding).Create(ctx, options)
}

func (c *clusterRoleBindingHandler) Subjects(ctx context.Context, options Options) error {
	if c.clusterRoleBinding == nil {
		return errors.New("can't display subjects for nil rolebinding")
	}

	c.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return c.subjectsFunc(ctx, c.clusterRoleBinding, options)
		},
	})
	return nil
}

func defaultClusterRoleBindingSubjects(ctx context.Context, clusterRoleBinding *rbacv1.ClusterRoleBinding, options Options) (component.Component, error) {
	return createClusterRoleBindingSubjectsView(clusterRoleBinding)
}
