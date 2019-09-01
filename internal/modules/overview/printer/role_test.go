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

func Test_RoleListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := testutil.Time()

	role := testutil.CreateRole("pod-reader")
	role.CreationTimestamp = metav1.Time{Time: now}
	roleList := &rbacv1.RoleList{
		Items: []rbacv1.Role{
			*role,
		},
	}

	tpo.PathForObject(role, role.Name, "/role")

	ctx := context.Background()
	observed, err := RoleListHandler(ctx, roleList, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age")
	expected := component.NewTable("Roles", "We couldn't find any roles!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", role.Name, "/role"),
		"Age":  component.NewTimestamp(role.CreationTimestamp.Time),
	})

	component.AssertEqual(t, expected, observed)
}

func Test_printRoleConfig(t *testing.T) {
	now := testutil.Time()

	role := testutil.CreateRole("pod-reader")
	role.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printRoleConfig(role)
	require.NoError(t, err)

	sections := component.SummarySections{}
	sections.AddText("Name", role.Name)
	expected := component.NewSummary("Configuration", sections...)

	component.AssertEqual(t, expected, observed)
}

func Test_printRolePolicyRules(t *testing.T) {
	now := testutil.Time()

	role := testutil.CreateRole("pod-reader")
	role.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printRolePolicyRules(role)
	require.NoError(t, err)

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	expected := component.NewTable("PolicyRules", "There are no policy rules!", cols)

	row := component.TableRow{}
	row["Resources"] = component.NewText("pods")
	row["Non-Resource URLs"] = component.NewText("")
	row["Resource Names"] = component.NewText("")
	row["Verbs"] = component.NewText("['get', 'watch', 'list']")

	expected.Add(row)

	component.AssertEqual(t, expected, observed)
}
