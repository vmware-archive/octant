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

func RoleBindingHandler(ctx context.Context, roleBinding *rbacv1.RoleBinding, opts Options) (component.Component, error) {
	o := NewObject(roleBinding)

	configSummary, err := printRoleBindingConfig(ctx, roleBinding, opts)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(configSummary)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return printRoleBindingSubjects(ctx, roleBinding, opts)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, opts)
}

func printRoleBindingConfig(ctx context.Context, roleBinding *rbacv1.RoleBinding, options Options) (*component.Summary, error) {
	if roleBinding == nil {
		return nil, errors.New("role binding is nil")
	}

	sections := component.SummarySections{}

	sections.AddText("Role kind", roleBinding.RoleRef.Kind)
	roleName, err := roleLinkFromRoleBinding(ctx, roleBinding, options)
	if err != nil {
		return nil, err
	}

	sections.Add("Role name", roleName)

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func printRoleBindingSubjects(ctx context.Context, roleBinding *rbacv1.RoleBinding, options Options) (component.Component, error) {
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
