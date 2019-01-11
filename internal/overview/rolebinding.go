package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type RoleBindingSummary struct{}

var _ view.View = (*RoleBindingSummary)(nil)

func NewRoleBindingSummary(prefix, namespace string, c clock.Clock) view.View {
	return &RoleBindingSummary{}
}

func (js *RoleBindingSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	roleBinding, err := retrieveRoleBinding(object)
	if err != nil {
		return nil, err
	}

	role, err := getRole(roleBinding.GetNamespace(), roleBinding.RoleRef.Name, c)
	if err != nil {
		return nil, err
	}

	detail, err := printRoleBindingSummary(roleBinding, role)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type RoleBindingSubjects struct{}

var _ view.View = (*RoleBindingSubjects)(nil)

func NewRoleBindingSubjects(prefix, namespace string, c clock.Clock) view.View {
	return &RoleBindingSubjects{}
}

func (js *RoleBindingSubjects) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	roleBinding, err := retrieveRoleBinding(object)
	if err != nil {
		return nil, err
	}

	subjectsTable, err := printRoleBindingSubjects(roleBinding)
	if err != nil {
		return nil, err
	}

	return []content.Content{
		&subjectsTable,
	}, nil
}

func retrieveRoleBinding(object runtime.Object) (*rbacv1.RoleBinding, error) {
	rc, ok := object.(*rbacv1.RoleBinding)
	if !ok {
		return nil, errors.Errorf("expected object to be a RoleBinding, it was %T", object)
	}

	return rc, nil
}
