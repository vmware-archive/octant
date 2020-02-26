/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// RoleBindingListHandler is a printFunc that prints RoleBindings
func RoleBindingListHandler(ctx context.Context, roleBindingList *rbacv1.RoleBindingList, opts Options) (component.Component, error) {
	if roleBindingList == nil {
		return nil, errors.New("role binding list is nil")
	}

	columns := component.NewTableCols("Name", "Age", "Role kind", "Role name")
	table := component.NewTable("Role Bindings", "We couldn't find any role bindings!", columns)

	for _, roleBinding := range roleBindingList.Items {
		row := component.TableRow{}
		nameLink, err := opts.Link.ForObject(&roleBinding, roleBinding.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Age"] = component.NewTimestamp(roleBinding.CreationTimestamp.Time)
		row["Role kind"] = component.NewText(roleBinding.RoleRef.Kind)
		roleName, err := roleLinkFromRoleBinding(ctx, &roleBinding, opts)
		if err != nil {
			return nil, err
		}
		row["Role name"] = roleName

		table.Add(row)
	}

	return table, nil
}

func roleLinkFromRoleBinding(ctx context.Context, roleBinding *rbacv1.RoleBinding, options Options) (*component.Link, error) {
	roleRef := roleBinding.RoleRef

	namespace := ""
	if roleRef.Kind == "Role" {
		namespace = roleBinding.Namespace
	}

	apiVersion := fmt.Sprintf("%s/%s", roleRef.APIGroup, "v1")
	roleLink, err := options.Link.ForGVK(namespace, apiVersion, roleRef.Kind, roleRef.Name, roleRef.Name)
	if err != nil {
		return nil, err
	}

	return roleLink, nil
}

// RoleBindingHandler is a printfunc that prints a RoleBinding
func RoleBindingHandler(ctx context.Context, roleBinding *rbacv1.RoleBinding, options Options) (component.Component, error) {
	o := NewObject(roleBinding)

	rh, err := newRoleBindingHandler(roleBinding, o)
	if err != nil {
		return nil, err
	}

	if err := rh.Config(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print rolebinding configuration")
	}

	if err := rh.Subjects(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print rolebinding subjects")
	}

	return o.ToComponent(ctx, options)
}

// RoleBindingConfiguration generates a rolebinding configuration
type RoleBindingConfiguration struct {
	roleBinding *rbacv1.RoleBinding
}

// NewRoleBindingConfiguration creates an instance of RoleBindingConfiguration
func NewRoleBindingConfiguration(roleBinding *rbacv1.RoleBinding) *RoleBindingConfiguration {
	return &RoleBindingConfiguration{
		roleBinding: roleBinding,
	}
}

// Create creates a rolebinding configuration summary
func (r *RoleBindingConfiguration) Create(ctx context.Context, options Options) (*component.Summary, error) {
	if r == nil || r.roleBinding == nil {
		return nil, errors.New("role binding is nil")
	}

	sections := component.SummarySections{}

	sections.AddText("Role kind", r.roleBinding.RoleRef.Kind)
	roleName, err := roleLinkFromRoleBinding(ctx, r.roleBinding, options)
	if err != nil {
		return nil, err
	}

	sections.Add("Role name", roleName)

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func createRoleBindingSubjectsView(ctx context.Context, roleBinding *rbacv1.RoleBinding, options Options) (component.Component, error) {
	if roleBinding == nil {
		return nil, errors.New("role binding is nil")
	}

	columns := component.NewTableCols("Kind", "Name", "Namespace")
	table := component.NewTable("Subjects", "There are no subjects!", columns)

	for _, subject := range roleBinding.Subjects {
		row := component.TableRow{}

		row["Kind"] = component.NewText(subject.Kind)

		if subject.Kind == "ServiceAccount" {
			name, err := serviceAccountLinkFromSubjects(ctx, &subject, options)
			if err != nil {
				return nil, err
			}
			row["Name"] = name
		} else {
			row["Name"] = component.NewText(subject.Name)
		}

		row["Namespace"] = component.NewText(subject.Namespace)

		table.Add(row)
	}
	return table, nil
}

func serviceAccountLinkFromSubjects(_ context.Context, subject *rbacv1.Subject, options Options) (*component.Link, error) {
	namespace := ""
	if subject.Kind == "ServiceAccount" {
		namespace = subject.Namespace
	}

	return options.Link.ForGVK(namespace, "v1", subject.Kind, subject.Name, subject.Name)
}

type roleBindingObject interface {
	Config(ctx context.Context, options Options) error
	Subjects(ctx context.Context, options Options) error
}

type roleBindingHandler struct {
	roleBinding  *rbacv1.RoleBinding
	configFunc   func(context.Context, *rbacv1.RoleBinding, Options) (*component.Summary, error)
	subjectsFunc func(context.Context, *rbacv1.RoleBinding, Options) (component.Component, error)
	object       *Object
}

var _ roleBindingObject = (*roleBindingHandler)(nil)

func newRoleBindingHandler(roleBinding *rbacv1.RoleBinding, object *Object) (*roleBindingHandler, error) {
	if roleBinding == nil {
		return nil, errors.New("can't print a nil rolebinding")
	}

	if object == nil {
		return nil, errors.New("can't print rolebinding using a nil object printer")
	}

	rh := &roleBindingHandler{
		roleBinding:  roleBinding,
		configFunc:   defaultRoleBindingConfig,
		subjectsFunc: defaultRoleBindingSubjects,
		object:       object,
	}

	return rh, nil
}

func (r *roleBindingHandler) Config(ctx context.Context, options Options) error {
	out, err := r.configFunc(ctx, r.roleBinding, options)
	if err != nil {
		return err
	}
	r.object.RegisterConfig(out)
	return nil
}

func defaultRoleBindingConfig(ctx context.Context, roleBinding *rbacv1.RoleBinding, options Options) (*component.Summary, error) {
	return NewRoleBindingConfiguration(roleBinding).Create(ctx, options)
}

func (r *roleBindingHandler) Subjects(ctx context.Context, options Options) error {
	if r.roleBinding == nil {
		return errors.New("can't display subjects for nil rolebinding")
	}

	r.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return r.subjectsFunc(ctx, r.roleBinding, options)
		},
	})
	return nil
}

func defaultRoleBindingSubjects(ctx context.Context, roleBinding *rbacv1.RoleBinding, options Options) (component.Component, error) {
	return createRoleBindingSubjectsView(ctx, roleBinding, options)
}
