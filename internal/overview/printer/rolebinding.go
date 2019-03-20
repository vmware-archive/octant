package printer

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
)

func RoleBindingListHandler(ctx context.Context, roleBindingList *rbacv1.RoleBindingList, opts Options) (component.Component, error) {
	if roleBindingList == nil {
		return nil, errors.New("role binding list is nil")
	}

	columns := component.NewTableCols("Name", "Age", "Role kind", "Role name")
	table := component.NewTable("Role Bindings", columns)

	for _, roleBinding := range roleBindingList.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&roleBinding, roleBinding.Name)
		row["Age"] = component.NewTimestamp(roleBinding.CreationTimestamp.Time)
		row["Role kind"] = component.NewText(roleBinding.RoleRef.Kind)
		row["Role name"] = roleLinkFromRoleBinding(ctx, &roleBinding)

		table.Add(row)
	}

	return table, nil
}

func roleLinkFromRoleBinding(ctx context.Context, roleBinding *rbacv1.RoleBinding) *component.Link {
	roleRef := roleBinding.RoleRef

	namespace := ""
	if roleRef.Kind == "Role" {
		namespace = roleBinding.Namespace
	}

	apiVersion := fmt.Sprintf("%s/%s", roleRef.APIGroup, "v1")
	return link.ForGVK(namespace, apiVersion, roleRef.Kind, roleRef.Name, roleRef.Name)
}

func RoleBindingHandler(ctx context.Context, roleBinding *rbacv1.RoleBinding, opts Options) (component.Component, error) {
	o := NewObject(roleBinding)

	o.RegisterConfig(func() (component.Component, error) {
		return printRoleBindingConfig(ctx, roleBinding)
	}, 16)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return printRoleBindingSubjects(roleBinding)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, opts)
}

func printRoleBindingConfig(ctx context.Context, roleBinding *rbacv1.RoleBinding) (component.Component, error) {
	if roleBinding == nil {
		return nil, errors.New("role binding is nil")
	}

	var sections component.SummarySections

	sections.AddText("Role kind", roleBinding.RoleRef.Kind)
	sections.Add("Role name", roleLinkFromRoleBinding(ctx, roleBinding))

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func printRoleBindingSubjects(roleBinding *rbacv1.RoleBinding) (component.Component, error) {
	if roleBinding == nil {
		return nil, errors.New("role binding is nil")
	}

	columns := component.NewTableCols("Kind", "Name", "Namespace")
	table := component.NewTable("Subjects", columns)

	for _, subject := range roleBinding.Subjects {
		row := component.TableRow{}
		row["Kind"] = component.NewText(subject.Kind)
		row["Name"] = component.NewText(subject.Name)
		row["Namespace"] = component.NewText(subject.Namespace)

		table.Add(row)
	}

	return table, nil
}
