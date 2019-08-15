/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_ClusterRoleBindingListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := testutil.Time()

	subjects := []rbacv1.Subject{
		{
			Kind: "User",
			Name: "test@example.com",
		},
	}
	clusterRoleBinding := testutil.CreateClusterRoleBinding("read-pods", "role-name", subjects)
	labels := map[string]string{"foo": "bar"}
	clusterRoleBinding.Labels = labels
	clusterRoleBinding.CreationTimestamp = metav1.Time{Time: now}
	roleBindingList := &rbacv1.ClusterRoleBindingList{
		Items: []rbacv1.ClusterRoleBinding{
			*clusterRoleBinding,
		},
	}

	tpo.PathForObject(clusterRoleBinding, clusterRoleBinding.Name, "/cluster-role-binding-path")
	tpo.PathForGVK("", "rbac.authorization.k8s.io/v1", "Role", "role-name", "role-name", "/cluster-role-path")

	ctx := context.Background()
	observed, err := ClusterRoleBindingListHandler(ctx, roleBindingList, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Age", "Role kind", "Role name")
	expected := component.NewTable("Cluster Role Bindings", "We couldn't find any cluster role bindings!", cols)
	expected.Add(component.TableRow{
		"Name":      component.NewLink("", clusterRoleBinding.Name, "/cluster-role-binding-path"),
		"Labels":    component.NewLabels(labels),
		"Age":       component.NewTimestamp(now),
		"Role kind": component.NewText("Role"),
		"Role name": component.NewLink("", "role-name", "/cluster-role-path"),
	})

	component.AssertEqual(t, expected, observed)
}

func Test_printClusterRoleBindingSubjects(t *testing.T) {
	now := testutil.Time()

	subjects := []rbacv1.Subject{
		{
			Kind: "User",
			Name: "test@example.com",
		},
	}
	clusterRoleBinding := testutil.CreateClusterRoleBinding("read-pods", "role-name", subjects)
	labels := map[string]string{"foo": "bar"}
	clusterRoleBinding.Labels = labels
	clusterRoleBinding.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRoleBindingSubjects(clusterRoleBinding)
	require.NoError(t, err)

	columns := component.NewTableCols("Kind", "Name", "Namespace")
	expected := component.NewTable("Subjects", "There are no subjects!", columns)

	row := component.TableRow{}
	row["Kind"] = component.NewText("User")
	row["Name"] = component.NewText("test@example.com")
	row["Namespace"] = component.NewText("")

	expected.Add(row)

	component.AssertEqual(t, expected, observed)
}

func Test_printClusterRoleBindingConfig(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})

	tpo.PathForGVK(subject.Namespace, "rbac.authorization.k8s.io/v1", "Role", "pod-reader", "pod-reader", "/role-path")

	ctx := context.Background()
	observed, err := printRoleBindingConfig(ctx, roleBinding, printOptions)
	require.NoError(t, err)

	sections := component.SummarySections{}

	sections.AddText("Role kind", "Role")
	sections.Add("Role name", component.NewLink("", "pod-reader", "/role-path"))

	expected := component.NewSummary("Configuration", sections...)

	component.AssertEqual(t, expected, observed)
}
