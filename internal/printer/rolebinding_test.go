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

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_RoleBindingListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := testutil.Time()

	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})
	roleBinding.CreationTimestamp = *testutil.CreateTimestamp()
	roleBindingList := &rbacv1.RoleBindingList{
		Items: []rbacv1.RoleBinding{
			*roleBinding,
		},
	}

	tpo.PathForObject(roleBinding, roleBinding.Name, "/role-binding")
	tpo.PathForGVK(roleBinding.Namespace, rbacAPIVersion, "Role", "pod-reader", "pod-reader", "/role")

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, roleBinding)
	observed, err := RoleBindingListHandler(ctx, roleBindingList, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age", "Role kind", "Role name")
	expected := component.NewTable("Role Bindings", "We couldn't find any role bindings!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", roleBinding.Name, "/role-binding",
			genObjectStatus(component.TextStatusOK, []string{
				"rbac.authorization.k8s.io/v1 RoleBinding is OK",
			})),
		"Age":       component.NewTimestamp(now),
		"Role kind": component.NewText("Role"),
		"Role name": component.NewLink("", "pod-reader", "/role"),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, roleBinding),
		}),
	})

	component.AssertEqual(t, expected, observed)
}

func Test_RoleBindingConfiguration(t *testing.T) {
	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})

	cases := []struct {
		name        string
		roleBinding *rbacv1.RoleBinding
		isErr       bool
		expected    *component.Summary
	}{
		{
			name:        "general",
			roleBinding: roleBinding,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Role kind",
					Content: component.NewText("Role"),
				},
				{
					Header:  "Role name",
					Content: component.NewLink("", "pod-reader", "/role"),
				},
			}...),
		},
		{
			name:        "rolebinding is nil",
			roleBinding: nil,
			isErr:       true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			tpo.PathForGVK("namespace", rbacAPIVersion, "Role", "pod-reader", "pod-reader", "/role")

			ctx := context.Background()

			rc := NewRoleBindingConfiguration(tc.roleBinding)

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

func Test_createRoleBindingSubjectsView(t *testing.T) {
	cases := []struct {
		name     string
		subject  *rbacv1.Subject
		expected component.TableRow
	}{
		{
			name:    "User",
			subject: testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace"),
			expected: component.TableRow{
				"Kind":      component.NewText("User"),
				"Name":      component.NewText("test@test.com"),
				"Namespace": component.NewText("namespace"),
			},
		},
		{
			name:    "Service Account",
			subject: testutil.CreateRoleBindingSubject("ServiceAccount", "svc-auto", "namespace"),
			expected: component.TableRow{
				"Kind":      component.NewText("ServiceAccount"),
				"Name":      component.NewLink("", "serviceAccount", "/service-account"),
				"Namespace": component.NewText("namespace"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*tc.subject})

			if tc.subject != nil {
				serviceAccountLink := component.NewLink("", "serviceAccount", "/service-account")
				tpo.link.EXPECT().
					ForGVK(gomock.Any(), "v1", "ServiceAccount", gomock.Any(), gomock.Any()).
					Return(serviceAccountLink, nil).
					AnyTimes()
			}

			ctx := context.Background()
			observed, err := createRoleBindingSubjectsView(ctx, roleBinding, printOptions)
			require.NoError(t, err)

			expected := component.NewTableWithRows("Subjects", "There are no subjects!",
				component.NewTableCols("Kind", "Name", "Namespace"),
				[]component.TableRow{tc.expected})

			component.AssertEqual(t, expected, observed)
		})
	}
}
