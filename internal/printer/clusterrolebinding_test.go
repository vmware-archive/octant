/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
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

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
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
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, clusterRoleBinding)
	observed, err := ClusterRoleBindingListHandler(ctx, roleBindingList, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Age", "Role kind", "Role name")
	expected := component.NewTable("Cluster Role Bindings", "We couldn't find any cluster role bindings!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", clusterRoleBinding.Name, "/cluster-role-binding-path",
			genObjectStatus(component.TextStatusOK, []string{
				"rbac.authorization.k8s.io/v1 ClusterRoleBinding is OK",
			})),
		"Labels":    component.NewLabels(labels),
		"Age":       component.NewTimestamp(now),
		"Role kind": component.NewText("Role"),
		"Role name": component.NewLink("", "role-name", "/cluster-role-path"),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, clusterRoleBinding),
		}),
	})

	component.AssertEqual(t, expected, observed)
}

func Test_ClusterRoleBindingConfiguration(t *testing.T) {
	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	clusterRoleBinding := testutil.CreateClusterRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})
	clusterRoleBinding.Namespace = "namespace"

	cases := []struct {
		name              string
		clusteRoleBinding *rbacv1.ClusterRoleBinding
		isErr             bool
		expected          *component.Summary
	}{
		{
			name:              "general",
			clusteRoleBinding: clusterRoleBinding,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Role kind",
					Content: component.NewText("Role"),
				},
				{
					Header:  "Role name",
					Content: component.NewLink("", "pod-reader", "/role-path"),
				},
			}...),
		},
		{
			name:              "clusterrolebinding is nil",
			clusteRoleBinding: nil,
			isErr:             true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			// tpo.PathForGVK(subject.Namespace, "rbac.authorization.k8s.io/v1", "Role", "pod-reader", "pod-reader")
			roleLink := component.NewLink("", "pod-reader", "/role-path")
			tpo.link.EXPECT().
				ForGVK(subject.Namespace, "rbac.authorization.k8s.io/v1", "Role", "pod-reader", "pod-reader").
				Return(roleLink, nil).AnyTimes()

			ctx := context.Background()

			rc := NewClusterRoleBindingConfiguration(tc.clusteRoleBinding)

			summary, err := rc.Create(ctx, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createClusterRoleBindingSubjectsView(t *testing.T) {
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

	observed, err := createClusterRoleBindingSubjectsView(clusterRoleBinding)
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
