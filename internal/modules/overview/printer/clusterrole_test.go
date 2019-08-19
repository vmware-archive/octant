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

func Test_ClusterRoleListHandler(t *testing.T) {
	object := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	object.CreationTimestamp = *testutil.CreateTimestamp()

	list := &rbacv1.ClusterRoleList{
		Items: []rbacv1.ClusterRole{*object},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(object, object.Name, "/path")

	now := testutil.Time()

	ctx := context.Background()
	got, err := ClusterRoleListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age")
	expected := component.NewTable("Cluster Roles", "We couldn't find any cluster roles!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", object.Name, "/path"),
		"Age":  component.NewTimestamp(now),
	})

	component.AssertEqual(t, expected, got)
}

func Test_printClusterRoleConfig(t *testing.T) {
	now := testutil.Time()

	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	clusterRole.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRoleConfig(clusterRole)
	require.NoError(t, err)

	sections := component.SummarySections{}
	sections.AddText("Name", clusterRole.Name)
	expected := component.NewSummary("Configuration", sections...)

	component.AssertEqual(t, expected, observed)
}

func Test_printClusterRolePolicyRule(t *testing.T) {
	now := testutil.Time()

	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	clusterRole.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRolePolicyRules(clusterRole)
	require.NoError(t, err)

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	expected := component.NewTable("Policy Rules", "There are no policy rules!", cols)

	row := component.TableRow{}
	row["Resources"] = component.NewText("crontabs.stable.example.com")
	row["Non-Resource URLs"] = component.NewText("")
	row["Resource Names"] = component.NewText("")
	row["Verbs"] = component.NewText("['get', 'list', 'watch', 'create', 'update', 'patch', 'delete']")

	expected.Add(row)

	component.AssertEqual(t, expected, observed)
}
