package printer

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ClusterRoleListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	object.CreationTimestamp = metav1.Time{Time: now}

	list := &rbacv1.ClusterRoleList{
		Items: []rbacv1.ClusterRole{*object},
	}

	got, err := ClusterRoleListHandler(list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age")
	expected := component.NewTable("Cluster Roles", cols)
	expected.Add(component.TableRow{
		"Name": link.ForObject(object, object.Name),
		"Age":  component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}

func Test_printClusterRoleConfig(t *testing.T) {
	now := time.Unix(1547211430, 0)

	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	clusterRole.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRoleConfig(clusterRole)
	require.NoError(t, err)

	var sections component.SummarySections
	sections.AddText("Name", clusterRole.Name)
	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, observed)
}

func Test_printClusterRolePolicyRule(t *testing.T) {
	now := time.Unix(1547211430, 0)

	clusterRole := testutil.CreateClusterRole("aggregate-cron-tabs-edit")
	clusterRole.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRolePolicyRules(clusterRole)
	require.NoError(t, err)

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	expected := component.NewTable("PolicyRules", cols)

	row := component.TableRow{}
	row["Resources"] = component.NewText("crontabs.stable.example.com")
	row["Non-Resource URLs"] = component.NewText("")
	row["Resource Names"] = component.NewText("")
	row["Verbs"] = component.NewText("['get', 'list', 'watch', 'create', 'update', 'patch', 'delete']")

	expected.Add(row)

	assert.Equal(t, expected, observed)
}
