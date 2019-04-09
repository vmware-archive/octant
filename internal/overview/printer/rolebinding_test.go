package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_RoleBindingListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	now := time.Unix(1547211430, 0)

	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})
	roleBinding.CreationTimestamp = metav1.Time{Time: now}
	roleBindingList := &rbacv1.RoleBindingList{
		Items: []rbacv1.RoleBinding{
			*roleBinding,
		},
	}


	c := cachefake.NewMockCache(controller)

	ctx := context.Background()
	observed, err := RoleBindingListHandler(ctx, roleBindingList, Options{Cache: c})
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age", "Role kind", "Role name")
	expected := component.NewTable("Role Bindings", cols)
	expected.Add(component.TableRow{
		"Name":      link.ForObject(roleBinding, "read-pods"),
		"Age":       component.NewTimestamp(now),
		"Role kind": component.NewText("Role"),
		"Role name": component.NewLink("", "pod-reader", "/content/overview/namespace/namespace/rbac/roles/pod-reader"),
	})

	assert.Equal(t, expected, observed)
}

func Test_printRoleBindingSubjects(t *testing.T) {
	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})

	observed, err := printRoleBindingSubjects(roleBinding)
	require.NoError(t, err)

	columns := component.NewTableCols("Kind", "Name", "Namespace")
	expected := component.NewTable("Subjects", columns)

	row := component.TableRow{}
	row["Kind"] = component.NewText("User")
	row["Name"] = component.NewText("test@test.com")
	row["Namespace"] = component.NewText("")

	expected.Add(row)

	assert.Equal(t, expected, observed)
}

func Test_printRoleBindingConfig(t *testing.T) {
	subject := testutil.CreateRoleBindingSubject("User", "test@test.com", "namespace")
	roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*subject})

	ctx := context.Background()
	observed, err := printRoleBindingConfig(ctx, roleBinding)
	require.NoError(t, err)

	var sections component.SummarySections

	sections.AddText("Role kind", "Role")
	sections.Add("Role name", component.NewLink("", "pod-reader", "/content/overview/namespace/namespace/rbac/roles/pod-reader"))

	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, observed)
}
