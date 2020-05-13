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
		"Name": component.NewLink("", object.Name, "/path",
			genObjectStatus(component.TextStatusOK, []string{
				"rbac.authorization.k8s.io/v1 ClusterRole is OK",
			})),
		"Age": component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
	})

	testutil.AssertJSONEqual(t, expected, got)
}

func Test_ClusterRoleConfiguration(t *testing.T) {
	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")

	cases := []struct {
		name        string
		clusterRole *rbacv1.ClusterRole
		isErr       bool
		expected    *component.Summary
	}{
		{
			name:        "general",
			clusterRole: clusterRole,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Name",
					Content: component.NewText("aggregate-cron-tabs-edit"),
				},
			}...),
		},
		{
			name:        "clusterrole is nil",
			clusterRole: nil,
			isErr:       true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			cc := NewClusterRoleConfiguration(test.clusterRole)

			summary, err := cc.Create(printOptions)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, test.expected, summary)
		})
	}
}

func Test_createClusterRolePolicyRulesView(t *testing.T) {
	now := testutil.Time()

	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	clusterRole.CreationTimestamp = metav1.Time{Time: now}

	observed, err := createClusterRolePolicyRulesView(clusterRole)
	require.NoError(t, err)

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	expected := component.NewTable("Policy Rules", "There are no policy rules!", cols)
	expected.Add([]component.TableRow{
		{
			"Resources":         component.NewText("crontabs.stable.example.com"),
			"Non-Resource URLs": component.NewText(""),
			"Resource Names":    component.NewText(""),
			"Verbs":             component.NewText("['get', 'list', 'watch', 'create', 'update', 'patch', 'delete']"),
		},
	}...)

	component.AssertEqual(t, expected, observed)
}
