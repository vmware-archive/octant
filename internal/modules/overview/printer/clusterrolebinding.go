package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/heptio/developer-dash/pkg/view/component"
)

func ClusterRoleBindingListHandler(_ context.Context, clusterRoleBindingList *rbacv1.ClusterRoleBindingList, options Options) (component.Component, error) {
	if clusterRoleBindingList == nil {
		return nil, errors.New("cluster role binding list is nil")
	}

	columns := component.NewTableCols("Name", "Labels", "Age", "Role kind", "Role name")
	table := component.NewTable("Cluster Role Bindings", columns)

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

func roleLinkFromClusterRoleBinding(roleBinding *rbacv1.ClusterRoleBinding, options Options) (*component.Link, error) {
	roleRef := roleBinding.RoleRef

	namespace := ""
	if roleRef.Kind == "Role" {
		namespace = roleBinding.Namespace
	}

	apiVersion := fmt.Sprintf("%s/%s", roleRef.APIGroup, "v1")
	return options.Link.ForGVK(namespace, apiVersion, roleRef.Kind, roleRef.Name, roleRef.Name)
}

func ClusterRoleBindingHandler(ctx context.Context, roleBinding *rbacv1.ClusterRoleBinding, options Options) (component.Component, error) {
	o := NewObject(roleBinding)

	summary, err := printClusterRoleBindingConfig(roleBinding, options)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(summary)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return printClusterRoleBindingSubjects(roleBinding)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, options)
}

func printClusterRoleBindingConfig(roleBinding *rbacv1.ClusterRoleBinding, options Options) (*component.Summary, error) {
	if roleBinding == nil {
		return nil, errors.New("role binding is nil")
	}

	sections := component.SummarySections{}

	sections.AddText("Role kind", roleBinding.RoleRef.Kind)

	roleName, err := roleLinkFromClusterRoleBinding(roleBinding, options)
	if err != nil {
		return nil, err
	}

	sections.Add("Role name", roleName)

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func printClusterRoleBindingSubjects(roleBinding *rbacv1.ClusterRoleBinding) (component.Component, error) {
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
