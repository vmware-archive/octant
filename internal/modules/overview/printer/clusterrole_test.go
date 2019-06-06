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

func Test_ClusterRoleListHandler(t *testing.T) {
	object := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	object.CreationTimestamp = metav1.Time{Time: now}

	list := &rbacv1.ClusterRoleList{
		Items: []rbacv1.ClusterRole{*object},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(object, object.Name, "/path")

	now := time.Unix(1547211430, 0)

	ctx := context.Background()
	got, err := ClusterRoleListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age")
	expected := component.NewTable("Cluster Roles", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", object.Name, "/path"),
		"Age":  component.NewTimestamp(now),
	})


	assertComponentEqual(t, expected, got)
}

func Test_printClusterRoleConfig(t *testing.T) {
	now := time.Unix(1547211430, 0)

	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	clusterRole.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRoleConfig(clusterRole)
	require.NoError(t, err)

	sections := component.SummarySections{}
	sections.AddText("Name", clusterRole.Name)
	expected := component.NewSummary("Configuration", sections...)

	assertComponentEqual(t, expected, observed)
}

func Test_printClusterRolePolicyRule(t *testing.T) {
	now := time.Unix(1547211430, 0)

	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	clusterRole.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRolePolicyRules(clusterRole)
	require.NoError(t, err)

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	expected := component.NewTable("Policy Rules", cols)

	row := component.TableRow{}
	row["Resources"] = component.NewText("crontabs.stable.example.com")
	row["Non-Resource URLs"] = component.NewText("")
	row["Resource Names"] = component.NewText("")
	row["Verbs"] = component.NewText("['get', 'list', 'watch', 'create', 'update', 'patch', 'delete']")

	expected.Add(row)

	assert.Equal(t, expected, observed)
}
