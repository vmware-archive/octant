/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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

	"github.com/heptio/developer-dash/pkg/store"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func ServiceAccountListHandler(_ context.Context, list *corev1.ServiceAccountList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("service account list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Secrets", "Age")
	table := component.NewTable("Service Accounts", cols)

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

func ServiceAccountHandler(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) (component.Component, error) {
	h, err := newServiceAccountHandler(ctx, serviceAccount, options)
	if err != nil {
		return nil, err
	}

	return h.run()
}

type serviceAccountHandler struct {
	ctx            context.Context
	serviceAccount corev1.ServiceAccount
	options        Options

	configFunc      func(ctx context.Context, serviceAccount corev1.ServiceAccount, options Options) (*component.Summary, error)
	policyRulesFunc func(ctx context.Context, serviceAccount corev1.ServiceAccount, appObjectStore store.Store) (*component.Table, error)
}

func newServiceAccountHandler(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) (*serviceAccountHandler, error) {
	if serviceAccount == nil {
		return nil, errors.New("service account is nil")
	}

	return &serviceAccountHandler{
		ctx:            ctx,
		serviceAccount: *serviceAccount,
		options:        options,
		configFunc:     printServiceAccountConfig,
		policyRulesFunc: func(ctx context.Context, serviceAccount corev1.ServiceAccount, appObjectStore store.Store) (*component.Table, error) {
			s := newServiceAccountPolicyRules(ctx, serviceAccount, appObjectStore)
			return s.run()
		},
	}, nil
}

func (h *serviceAccountHandler) run() (component.Component, error) {
	o := NewObject(&h.serviceAccount)

	objectStore := h.options.DashConfig.ObjectStore()

	configSummary, err := h.configFunc(h.ctx, h.serviceAccount, h.options)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(configSummary)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return h.policyRulesFunc(h.ctx, h.serviceAccount, objectStore)
		},
		Width: component.WidthFull,
	})

	o.EnableEvents()

	return o.ToComponent(h.ctx, h.options)
}

func printServiceAccountConfig(ctx context.Context, serviceAccount corev1.ServiceAccount, options Options) (*component.Summary, error) {
	objectStore := options.DashConfig.ObjectStore()

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

	tokens, err := serviceAccountTokens(ctx, serviceAccount, objectStore)

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
	secretList, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "find secrets for service account")
	}

	var tokens []string

	for _, u := range secretList {
		secret := &corev1.Secret{}

		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, secret); err != nil {
			return nil, errors.Wrap(err, "convert unstructured secret to structured")
		}

		if err := copyObjectMeta(secret, u); err != nil {
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

type serviceAccountPolicyRules struct {
	ctx            context.Context
	serviceAccount corev1.ServiceAccount
	appObjectStore store.Store

	printPolicyRulesFunc func([]rbacv1.PolicyRule) (*component.Table, error)
}

func newServiceAccountPolicyRules(ctx context.Context, serviceAccount corev1.ServiceAccount, appObjectStore store.Store) *serviceAccountPolicyRules {
	return &serviceAccountPolicyRules{
		ctx:                  ctx,
		serviceAccount:       serviceAccount,
		appObjectStore:       appObjectStore,
		printPolicyRulesFunc: printPolicyRules,
	}
}

func (s *serviceAccountPolicyRules) run() (*component.Table, error) {
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
			object, err := s.appObjectStore.Get(s.ctx, key)
			if err != nil {
				return nil, err
			}

			clusterRole := &rbacv1.ClusterRole{}
			if err := scheme.Scheme.Convert(object, clusterRole, nil); err != nil {
				return nil, err
			}

			policyRules = append(policyRules, clusterRole.Rules...)

		case "Role":
			key.Namespace = s.serviceAccount.Namespace

			object, err := s.appObjectStore.Get(s.ctx, key)
			if err != nil {
				return nil, err
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

	return s.printPolicyRulesFunc(policyRules)
}

func (s *serviceAccountPolicyRules) listRoleBindings() ([]rbacv1.RoleRef, error) {
	roleBindingKey := store.Key{
		Namespace:  s.serviceAccount.Namespace,
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	}

	objects, err := s.appObjectStore.List(s.ctx, roleBindingKey)
	if err != nil {
		return nil, err
	}

	var list []rbacv1.RoleRef

	for _, object := range objects {
		roleBinding := &rbacv1.RoleBinding{}
		if err := scheme.Scheme.Convert(object, roleBinding, nil); err != nil {
			return nil, err
		}

		if s.isMatchSubjects(roleBinding.Subjects) {
			list = append(list, roleBinding.RoleRef)
		}
	}

	return list, nil
}

func (s *serviceAccountPolicyRules) listClusterRoleBindings() ([]rbacv1.RoleRef, error) {
	roleBindingKey := store.Key{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRoleBinding",
	}

	objects, err := s.appObjectStore.List(s.ctx, roleBindingKey)
	if err != nil {
		return nil, err
	}

	var list []rbacv1.RoleRef

	for _, object := range objects {
		roleBinding := &rbacv1.RoleBinding{}
		if err := scheme.Scheme.Convert(object, roleBinding, nil); err != nil {
			return nil, err
		}

		if s.isMatchSubjects(roleBinding.Subjects) {
			list = append(list, roleBinding.RoleRef)
		}
	}

	return list, nil
}

func (s *serviceAccountPolicyRules) isMatchSubjects(subjects []rbacv1.Subject) bool {
	subjectMatched := false
	for _, subject := range subjects {
		if s.isSubject(subject) {
			subjectMatched = true
			break
		}
	}

	return subjectMatched
}

func (s *serviceAccountPolicyRules) isSubject(subject rbacv1.Subject) bool {
	namespace := s.serviceAccount.Namespace
	inNamespace := fmt.Sprintf("system:serviceaccounts:%s", namespace)

	apiGroup := "rbac.authorization.k8s.io"

	if subject.Kind == "ServiceAccount" && subject.Name == s.serviceAccount.Name {
		return true
	} else if subject.Kind == "Group" && subject.Name == inNamespace && subject.APIGroup == apiGroup {
		return true
	} else if subject.Kind == "Group" && subject.Name == "system:serviceaccounts" && subject.APIGroup == apiGroup {
		return true
	}

	return false
}
