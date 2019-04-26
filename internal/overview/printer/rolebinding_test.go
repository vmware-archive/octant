package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
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

	o := storefake.NewMockObjectStore(controller)

	ctx := context.Background()
	observed, err := RoleBindingListHandler(ctx, roleBindingList, Options{ObjectStore: o})
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
				"Name":      component.NewLink("", "svc-auto", "/content/overview/namespace/namespace/config-and-storage/service-accounts/svc-auto"),
				"Namespace": component.NewText("namespace"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			roleBinding := testutil.CreateRoleBinding("read-pods", "pod-reader", []rbacv1.Subject{*tc.subject})

			ctx := context.Background()
			observed, err := printRoleBindingSubjects(ctx, roleBinding)
			require.NoError(t, err)

			expected := component.NewTableWithRows("Subjects",
				component.NewTableCols("Kind", "Name", "Namespace"),
				[]component.TableRow{tc.expected})

			assert.Equal(t, expected, observed)
		})
	}
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
