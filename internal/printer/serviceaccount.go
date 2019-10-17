/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ServiceAccountListHandler is a printFunc that prints service accounts
func ServiceAccountListHandler(_ context.Context, list *corev1.ServiceAccountList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("service account list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Secrets", "Age")
	table := component.NewTable("Service Accounts",
		"We couldn't find any service accounts!", cols)

	for _, serviceAccount := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&serviceAccount, serviceAccount.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(serviceAccount.Labels)
		row["Secrets"] = component.NewText(fmt.Sprint(len(serviceAccount.Secrets)))
		row["Age"] = component.NewTimestamp(serviceAccount.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

type serviceAccountObject interface {
	Config(ctx context.Context, options Options) error
	PolicyRules(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) error
}

type serviceAccountHandler struct {
	serviceAccount  *corev1.ServiceAccount
	configFunc      func(context.Context, *corev1.ServiceAccount, Options) (*component.Summary, error)
	policyRulesFunc func(context.Context, *corev1.ServiceAccount, Options) (*component.Table, error)
	object          *Object
}

var _serviceAccountObject = (*serviceAccountHandler)(nil)

func newServiceAccountHandler(serviceAccount *corev1.ServiceAccount, object *Object) (*serviceAccountHandler, error) {
	if serviceAccount == nil {
		return nil, errors.New("can't print a nil service account")
	}

	if object == nil {
		return nil, errors.New("can't print service account using a nil object printer")
	}

	s := &serviceAccountHandler{
		serviceAccount:  serviceAccount,
		configFunc:      defaultServiceAccountConfig,
		policyRulesFunc: defaultServiceAccountPolicyRules,
		object:          object,
	}
	return s, nil
}

func (s *serviceAccountHandler) Config(ctx context.Context, options Options) error {
	out, err := s.configFunc(ctx, s.serviceAccount, options)
	if err != nil {
		return err
	}
	s.object.RegisterConfig(out)
	return nil
}

// ServiceAccountHandler is a printFunc that prints ServiceAccounts
func ServiceAccountHandler(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) (component.Component, error) {
	o := NewObject(serviceAccount)
	o.EnableEvents()

	s, err := newServiceAccountHandler(serviceAccount, o)
	if err != nil {
		return nil, err
	}

	if err := s.Config(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print service account configuration")
	}

	return o.ToComponent(ctx, options)
}

// ServiceAccountConfiguration generates a service account configuration
type ServiceAccountConfiguration struct {
	context        context.Context
	serviceAccount *corev1.ServiceAccount
	objectStore    store.Store
}

// NewServiceAccountConfiguration creates an instance of ServiceAccountConfiguration
func NewServiceAccountConfiguration(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) *ServiceAccountConfiguration {
	return &ServiceAccountConfiguration{
		context:        ctx,
		serviceAccount: serviceAccount,
		objectStore:    options.DashConfig.ObjectStore(),
	}
}

// Create creates a service account configuration summary
func (s *ServiceAccountConfiguration) Create(options Options) (*component.Summary, error) {
	if s == nil || s.serviceAccount == nil {
		return nil, errors.New("service account is nil")
	}

	serviceAccount := s.serviceAccount

	sections := component.SummarySections{}

	var pullSecrets []string

	for _, s := range serviceAccount.ImagePullSecrets {
		pullSecrets = append(pullSecrets, s.Name)
	}

	if len(pullSecrets) > 0 {
		view, err := generateServiceAccountSecretsList(serviceAccount.Namespace, pullSecrets, options)
		if err != nil {
			return nil, err
		}

		sections.Add("Image Pull Secrets", view)
	}

	var mountSecrets []string
	for _, s := range serviceAccount.Secrets {
		mountSecrets = append(mountSecrets, s.Name)
	}

	if len(mountSecrets) > 0 {
		view, err := generateServiceAccountSecretsList(serviceAccount.Namespace, mountSecrets, options)
		if err != nil {
			return nil, err
		}
		sections.Add("Mountable Secrets", view)
	}

	tokens, err := serviceAccountTokens(s.context, *serviceAccount, s.objectStore)

	if err != nil {
		return nil, errors.Wrap(err, "get tokens for service account")
	}

	if len(tokens) > 0 {
		view, err := generateServiceAccountSecretsList(serviceAccount.Namespace, tokens, options)
		if err != nil {
			return nil, err
		}
		sections.Add("Tokens", view)
	}

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func generateServiceAccountSecretsList(namespace string, secretNames []string, options Options) (*component.List, error) {
	var items []component.Component
	for _, name := range secretNames {
		nameLink, err := options.Link.ForGVK(namespace, "v1", "Secret", name, name)
		if err != nil {
			return nil, err
		}
		items = append(items, nameLink)
	}
	return component.NewList("", items), nil
}

func serviceAccountTokens(ctx context.Context, serviceAccount corev1.ServiceAccount, o store.Store) ([]string, error) {
	key := store.Key{
		Namespace:  serviceAccount.Namespace,
		APIVersion: "v1",
		Kind:       "Secret",
	}
	secretList, _, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "find secrets for service account")
	}

	var tokens []string

	for i := range secretList.Items {
		secret := &corev1.Secret{}

		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(secretList.Items[i].Object, secret); err != nil {
			return nil, errors.Wrap(err, "convert unstructured secret to structured")
		}

		if err := copyObjectMeta(secret, &secretList.Items[i]); err != nil {
			return nil, errors.Wrap(err, "copy object metadata to secret")
		}

		if secret.Type == corev1.SecretTypeServiceAccountToken {
			name := secret.Annotations[corev1.ServiceAccountNameKey]
			uid := secret.Annotations[corev1.ServiceAccountUIDKey]

			if name == serviceAccount.Name && uid == string(serviceAccount.UID) {
				tokens = append(tokens, secret.Name)
			}
		}
	}

	return tokens, nil
}

// ServiceAccountPolicyRules generates a service account policy rule
type ServiceAccountPolicyRules struct {
	context     context.Context
	namespace   string
	name        string
	objectStore store.Store
}

// NewServiceAccountPolicyRules creates an instance of ServiceAccountPolicyRules
func NewServiceAccountPolicyRules(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) *ServiceAccountPolicyRules {
	if err := options.DashConfig.Validate(); err != nil {
		return nil
	}

	return &ServiceAccountPolicyRules{
		context:     ctx,
		namespace:   serviceAccount.Namespace,
		name:        serviceAccount.Name,
		objectStore: options.DashConfig.ObjectStore(),
	}
}

// Create generates a service account policy rule table
func (s *ServiceAccountPolicyRules) Create() (*component.Table, error) {
	if s.namespace == "" || s.name == "" {
		return nil, errors.New("serviceaccount is nil")
	}

	var roleRefs []rbacv1.RoleRef

	roleBindingRoleRefs, err := s.listRoleBindings()
	if err != nil {
		return nil, err
	}

	roleRefs = append(roleRefs, roleBindingRoleRefs...)

	clusterRoleBindingRefs, err := s.listClusterRoleBindings()
	if err != nil {
		return nil, err
	}

	roleRefs = append(roleRefs, clusterRoleBindingRefs...)

	var policyRules []rbacv1.PolicyRule

	for _, roleRef := range roleRefs {
		key := store.Key{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       roleRef.Kind,
			Name:       roleRef.Name,
		}
		switch kind := roleRef.Kind; kind {
		case "ClusterRole":
			object, found, err := s.objectStore.Get(s.context, key)
			if err != nil {
				return nil, err
			}

			if !found {
				return nil, errors.Errorf("unable to find %s", key)
			}

			clusterRole := &rbacv1.ClusterRole{}
			if err := scheme.Scheme.Convert(object, clusterRole, nil); err != nil {
				return nil, err
			}

			policyRules = append(policyRules, clusterRole.Rules...)

		case "Role":
			key.Namespace = s.namespace

			object, found, err := s.objectStore.Get(s.context, key)
			if err != nil {
				return nil, err
			}

			if !found {
				return nil, errors.Errorf("unable to find %s", key)
			}

			role := &rbacv1.Role{}
			if err := scheme.Scheme.Convert(object, role, nil); err != nil {
				return nil, err
			}

			policyRules = append(policyRules, role.Rules...)

		default:
			return nil, errors.Errorf("unable to handle role ref kind %q", kind)
		}
	}

	return printPolicyRules(policyRules)
}

func (s *ServiceAccountPolicyRules) listRoleBindings() ([]rbacv1.RoleRef, error) {
	roleBindingKey := store.Key{
		Namespace:  s.namespace,
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	}

	objects, _, err := s.objectStore.List(s.context, roleBindingKey)
	if err != nil {
		return nil, err
	}

	var list []rbacv1.RoleRef

	for i := range objects.Items {
		roleBinding := &rbacv1.RoleBinding{}
		if err := scheme.Scheme.Convert(&objects.Items[i], roleBinding, nil); err != nil {
			return nil, err
		}

		if s.isMatchSubjects(roleBinding.Subjects) {
			list = append(list, roleBinding.RoleRef)
		}
	}

	return list, nil
}

func (s *ServiceAccountPolicyRules) listClusterRoleBindings() ([]rbacv1.RoleRef, error) {
	roleBindingKey := store.Key{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRoleBinding",
	}

	objects, _, err := s.objectStore.List(s.context, roleBindingKey)
	if err != nil {
		return nil, err
	}

	var list []rbacv1.RoleRef

	for _, object := range objects.Items {
		roleBinding := &rbacv1.RoleBinding{}
		if err := scheme.Scheme.Convert(&object, roleBinding, nil); err != nil {
			return nil, err
		}

		if s.isMatchSubjects(roleBinding.Subjects) {
			list = append(list, roleBinding.RoleRef)
		}
	}

	return list, nil
}

func (s *ServiceAccountPolicyRules) isMatchSubjects(subjects []rbacv1.Subject) bool {
	subjectMatched := false
	for _, subject := range subjects {
		if s.isSubject(subject) {
			subjectMatched = true
			break
		}
	}

	return subjectMatched
}

func (s *ServiceAccountPolicyRules) isSubject(subject rbacv1.Subject) bool {
	inNamespace := fmt.Sprintf("system:serviceaccounts:%s", s.namespace)

	apiGroup := "rbac.authorization.k8s.io"

	if subject.Kind == "ServiceAccount" && subject.Name == s.name {
		return true
	} else if subject.Kind == "Group" && subject.Name == inNamespace && subject.APIGroup == apiGroup {
		return true
	} else if subject.Kind == "Group" && subject.Name == "system:serviceaccounts" && subject.APIGroup == apiGroup {
		return true
	}

	return false
}

func defaultServiceAccountConfig(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) (*component.Summary, error) {
	return NewServiceAccountConfiguration(ctx, serviceAccount, options).Create(options)
}

func (s *serviceAccountHandler) PolicyRules(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) error {
	s.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return s.policyRulesFunc(ctx, serviceAccount, options)
		},
	})

	return nil
}

func defaultServiceAccountPolicyRules(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) (*component.Table, error) {
	return NewServiceAccountPolicyRules(ctx, serviceAccount, options).Create()
}
