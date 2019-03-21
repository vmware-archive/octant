package printer

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
)

func ClusterRoleBindingListHandler(ctx context.Context, clusterRoleBindingList *rbacv1.ClusterRoleBindingList, opts Options) (component.Component, error) {
	if clusterRoleBindingList == nil {
		return nil, errors.New("cluster role binding list is nil")
	}

	columns := component.NewTableCols("Name", "Labels", "Age", "Role kind", "Role name")
	table := component.NewTable("Cluster Role Bindings", columns)

	for _, roleBinding := range clusterRoleBindingList.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&roleBinding, roleBinding.Name)
		row["Labels"] = component.NewLabels(roleBinding.Labels)
		row["Age"] = component.NewTimestamp(roleBinding.CreationTimestamp.Time)
		row["Role kind"] = component.NewText(roleBinding.RoleRef.Kind)
		row["Role name"] = roleLinkFromClusterRoleBinding(&roleBinding)

		table.Add(row)
	}

	return table, nil
}

func roleLinkFromClusterRoleBinding(roleBinding *rbacv1.ClusterRoleBinding) *component.Link {
	roleRef := roleBinding.RoleRef

	namespace := ""
	if roleRef.Kind == "Role" {
		namespace = roleBinding.Namespace
	}

	apiVersion := fmt.Sprintf("%s/%s", roleRef.APIGroup, "v1")
	return link.ForGVK(namespace, apiVersion, roleRef.Kind, roleRef.Name, roleRef.Name)
}

func ClusterRoleBindingHandler(ctx context.Context, roleBinding *rbacv1.ClusterRoleBinding, opts Options) (component.Component, error) {
	o := NewObject(roleBinding)

	summary, err := printClusterRoleBindingConfig(roleBinding)
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

	return o.ToComponent(ctx, opts)
}

func printClusterRoleBindingConfig(roleBinding *rbacv1.ClusterRoleBinding) (*component.Summary, error) {
	if roleBinding == nil {
		return nil, errors.New("role binding is nil")
	}

	var sections component.SummarySections

	sections.AddText("Role kind", roleBinding.RoleRef.Kind)
	sections.Add("Role name", roleLinkFromClusterRoleBinding(roleBinding))

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
