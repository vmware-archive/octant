/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_RoleBindingListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := time.Unix(1547211430, 0)

	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})
	roleBinding.CreationTimestamp = metav1.Time{Time: now}
	roleBindingList := &rbacv1.RoleBindingList{
		Items: []rbacv1.RoleBinding{
			*roleBinding,
		},
	}

	tpo.PathForObject(roleBinding, roleBinding.Name, "/role-binding")
	tpo.PathForGVK(roleBinding.Namespace, rbacAPIVersion, "Role", "pod-reader", "pod-reader", "/role")

	ctx := context.Background()
	observed, err := RoleBindingListHandler(ctx, roleBindingList, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age", "Role kind", "Role name")
	expected := component.NewTable("Role Bindings", cols)
	expected.Add(component.TableRow{
		"Name":      component.NewLink("", roleBinding.Name, "/role-binding"),
		"Age":       component.NewTimestamp(now),
		"Role kind": component.NewText("Role"),
		"Role name": component.NewLink("", "pod-reader", "/role"),
	})

	component.AssertEqual(t, expected, observed)
}

func Test_printRoleBindingSubjects(t *testing.T) {
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
			observed, err := printRoleBindingSubjects(ctx, roleBinding, printOptions)
			require.NoError(t, err)

			expected := component.NewTableWithRows("Subjects",
				component.NewTableCols("Kind", "Name", "Namespace"),
				[]component.TableRow{tc.expected})

			component.AssertEqual(t, expected, observed)
		})
	}
}

func Test_printRoleBindingConfig(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})

	tpo.PathForGVK("namespace", rbacAPIVersion, "Role", "pod-reader", "pod-reader", "/role")

	ctx := context.Background()
	observed, err := printRoleBindingConfig(ctx, roleBinding, printOptions)
	require.NoError(t, err)

	sections := component.SummarySections{}

	sections.AddText("Role kind", "Role")
	sections.Add("Role name", component.NewLink("", "pod-reader", "/role"))

	expected := component.NewSummary("Configuration", sections...)

	component.AssertEqual(t, expected, observed)
}
