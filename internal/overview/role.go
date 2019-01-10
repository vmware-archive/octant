package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

type RoleSummary struct{}

var _ View = (*RoleSummary)(nil)

func NewRoleSummary(prefix, namespace string, c clock.Clock) View {
	return &RoleSummary{}
}

func (js *RoleSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	secret, err := retrieveRole(object)
	if err != nil {
		return nil, err
	}

	detail, err := printRoleSummary(secret)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type RoleRule struct{}

var _ View = (*RoleRule)(nil)

func NewRoleRule(prefix, namespace string, c clock.Clock) View {
	return &RoleRule{}
}

func (js *RoleRule) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	secret, err := retrieveRole(object)
	if err != nil {
		return nil, err
	}

	rulesTable, err := printRoleRule(secret)
	if err != nil {
		return nil, err
	}

	return []content.Content{
		&rulesTable,
	}, nil
}

func retrieveRole(object runtime.Object) (*rbacv1.Role, error) {
	rc, ok := object.(*rbacv1.Role)
	if !ok {
		return nil, errors.Errorf("expected object to be a Role, it was %T", object)
	}

	return rc, nil
}

func getRole(namespace, name string, c cache.Cache) (*rbacv1.Role, error) {
	key := cache.Key{
		Namespace:  namespace,
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
		Name:       name,
	}

	roles, err := loadRoles(key, c)
	if err != nil {
		return nil, err
	}

	if len(roles) != 1 {
		return nil, errors.Errorf("expected exactly one Role; got %d", len(roles))
	}

	return roles[0], nil
}

func loadRoles(key cache.Key, c cache.Cache) ([]*rbacv1.Role, error) {
	objects, err := c.Retrieve(key)
	if err != nil {
		return nil, err
	}

	var list []*rbacv1.Role

	for _, object := range objects {
		e := &rbacv1.Role{}
		if err := scheme.Scheme.Convert(object, e, runtime.InternalGroupVersioner); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(e, object); err != nil {
			return nil, err
		}

		list = append(list, e)
	}

	return list, nil
}
