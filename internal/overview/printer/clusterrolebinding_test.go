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

func Test_ClusterRoleBindingListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	now := time.Unix(1547211430, 0)

	clusterRoleBinding := testutil.CreateClusterRoleBinding("read-pods")
	labels := map[string]string{"foo": "bar"}
	clusterRoleBinding.Labels = labels
	clusterRoleBinding.CreationTimestamp = metav1.Time{Time: now}
	roleBindingList := &rbacv1.ClusterRoleBindingList{
		Items: []rbacv1.ClusterRoleBinding{
			*clusterRoleBinding,
		},
	}

	c := cachefake.NewMockCache(controller)

	observed, err := ClusterRoleBindingListHandler(roleBindingList, Options{Cache: c})
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Age", "Role kind", "Role name")
	expected := component.NewTable("Cluster Role Bindings", cols)
	expected.Add(component.TableRow{
		"Name":      link.ForObject(clusterRoleBinding, "read-pods"),
		"Labels":    component.NewLabels(labels),
		"Age":       component.NewTimestamp(now),
		"Role kind": component.NewText("Role"),
		"Role name": component.NewLink("", "role-name", "/content/overview/namespace/namespace/rbac/roles/role-name"),
	})

	assert.Equal(t, expected, observed)
}

func Test_printClusterRoleBindingSubjects(t *testing.T) {
	now := time.Unix(1547211430, 0)

	clusterRoleBinding := testutil.CreateClusterRoleBinding("read-pods")
	labels := map[string]string{"foo": "bar"}
	clusterRoleBinding.Labels = labels
	clusterRoleBinding.CreationTimestamp = metav1.Time{Time: now}

	observed, err := printClusterRoleBindingSubjects(clusterRoleBinding)
	require.NoError(t, err)

	columns := component.NewTableCols("Kind", "Name", "Namespace")
	expected := component.NewTable("Subjects", columns)

	row := component.TableRow{}
	row["Kind"] = component.NewText("User")
	row["Name"] = component.NewText("test@example.com")
	row["Namespace"] = component.NewText("")

	expected.Add(row)

	assert.Equal(t, expected, observed)
}

func Test_printClusterRoleBindingConfig(t *testing.T) {
	subject := testutil.CreateRoleBindingSubject("User", "test@test.com")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})

	observed, err := printRoleBindingConfig(roleBinding)
	require.NoError(t, err)

	var sections component.SummarySections

	sections.AddText("Role kind", "Role")
	sections.Add("Role name", component.NewLink("", "pod-reader", "/content/overview/namespace/namespace/rbac/roles/pod-reader"))

	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, observed)
}
