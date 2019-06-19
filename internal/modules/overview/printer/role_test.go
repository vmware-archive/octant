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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_RoleListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := time.Unix(1547211430, 0)

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
	expected := component.NewTable("Roles", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", role.Name, "/role"),
		"Age":  component.NewTimestamp(role.CreationTimestamp.Time),
	})

	assert.Equal(t, expected, observed)
}

func Test_printRoleConfig(t *testing.T) {
	now := time.Unix(1547211430, 0)

	role := testutil.CreateRole("pod-reader")
	role.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printRoleConfig(role)
	require.NoError(t, err)

	sections := component.SummarySections{}
	sections.AddText("Name", role.Name)
	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, observed)
}

func Test_printRolePolicyRules(t *testing.T) {
	now := time.Unix(1547211430, 0)

	role := testutil.CreateRole("pod-reader")
	role.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printRolePolicyRules(role)
	require.NoError(t, err)

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	expected := component.NewTable("PolicyRules", cols)

	row := component.TableRow{}
	row["Resources"] = component.NewText("pods")
	row["Non-Resource URLs"] = component.NewText("")
	row["Resource Names"] = component.NewText("")
	row["Verbs"] = component.NewText("['get', 'watch', 'list']")

	expected.Add(row)

	assert.Equal(t, expected, observed)
}
