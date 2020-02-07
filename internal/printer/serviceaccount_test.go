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
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ServiceAccountListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	object := testutil.CreateServiceAccount("sa")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels
	object.Secrets = []corev1.ObjectReference{{Name: "secret"}}

	tpo.PathForObject(object, object.Name, "/path")

	list := &corev1.ServiceAccountList{
		Items: []corev1.ServiceAccount{*object},
	}

	ctx := context.Background()
	got, err := ServiceAccountListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Secrets", "Age")
	expected := component.NewTable("Service Accounts", "We couldn't find any service accounts!", cols)
	expected.Add(component.TableRow{
		"Name":    component.NewLink("", object.Name, "/path"),
		"Labels":  component.NewLabels(labels),
		"Secrets": component.NewText("1"),
		"Age":     component.NewTimestamp(now),
	})

	component.AssertEqual(t, expected, got)
}

func Test_ServiceAccountConfiguration(t *testing.T) {
	serviceAccount := testutil.CreateServiceAccount("sa")
	serviceAccount.CreationTimestamp = metav1.Time{Time: testutil.Time()}
	serviceAccount.Secrets = []corev1.ObjectReference{{Name: "secret"}}
	serviceAccount.ImagePullSecrets = []corev1.LocalObjectReference{{Name: "secret"}}

	secret := testutil.CreateSecret("secret")
	secret.Type = corev1.SecretTypeServiceAccountToken
	secret.Annotations = map[string]string{
		corev1.ServiceAccountNameKey: serviceAccount.Name,
		corev1.ServiceAccountUIDKey:  string(serviceAccount.UID),
	}

	key := store.Key{
		Namespace:  serviceAccount.Namespace,
		APIVersion: "v1",
		Kind:       "Secret",
	}

	cases := []struct {
		name           string
		namespace      string
		serviceaccount *corev1.ServiceAccount
		secret         *corev1.Secret
		isErr          bool
		expected       *component.Summary
	}{
		{
			name:           "serviceaccount",
			namespace:      serviceAccount.Namespace,
			serviceaccount: serviceAccount,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header: "Image Pull Secrets",
					Content: component.NewList("", []component.Component{
						component.NewLink("", "secret", "/secret"),
					}),
				},
				{
					Header: "Mountable Secrets",
					Content: component.NewList("", []component.Component{
						component.NewLink("", "secret", "/secret"),
					}),
				},
				{
					Header: "Tokens",
					Content: component.NewList("", []component.Component{
						component.NewLink("", "secret", "/secret"),
					}),
				},
			}...),
		},
		{
			name:           "serviceaccount is nil",
			serviceaccount: nil,
			isErr:          true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			ctx := context.Background()
			tpo := newTestPrinterOptions(controller)

			if !tc.isErr {
				tpo.objectStore.EXPECT().List(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructuredList(t, testutil.ToUnstructured(t, secret)), false, nil)
			}

			printOptions := tpo.ToOptions()

			tpo.PathForGVK(tc.namespace, "v1", "Secret", "secret", "secret", "/secret")

			sac := NewServiceAccountConfiguration(ctx, tc.serviceaccount, printOptions)
			summary, err := sac.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_ServiceAccountPolicyRules(t *testing.T) {
	serviceAccount := testutil.CreateServiceAccount("sa")

	roleBindingKey := store.Key{
		Namespace:  serviceAccount.Namespace,
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	}

	role1 := testutil.CreateRole("role1")
	role2 := testutil.CreateClusterRole("role2")

	subjects1 := []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      serviceAccount.Name,
			Namespace: serviceAccount.Namespace,
		},
	}

	roleBinding := testutil.CreateRoleBinding("rb1", role1.Name, subjects1)
	roleBindingObjects := testutil.ToUnstructuredList(t, roleBinding)

	clusterRoleBindingKey := store.Key{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRoleBinding",
	}

	subjects2 := []rbacv1.Subject{
		{
			Kind:     "Group",
			Name:     "system:serviceaccounts",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	clusterRoleBinding := testutil.CreateClusterRoleBinding("crb1", role2.Name, subjects2)
	clusterRoleBinding.RoleRef.Kind = "ClusterRole"
	clusterRoleBindingObjects := testutil.ToUnstructuredList(t, clusterRoleBinding)

	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	role1Key, err := store.KeyFromObject(role1)
	require.NoError(t, err)

	role2Key, err := store.KeyFromObject(role2)
	require.NoError(t, err)

	tpo.objectStore.EXPECT().
		List(gomock.Any(), roleBindingKey).
		Return(roleBindingObjects, false, nil)

	tpo.objectStore.EXPECT().
		List(gomock.Any(), clusterRoleBindingKey).
		Return(clusterRoleBindingObjects, false, nil)

	tpo.objectStore.EXPECT().
		Get(gomock.Any(), role1Key).
		Return(testutil.ToUnstructured(t, role1), nil)

	tpo.objectStore.EXPECT().
		Get(gomock.Any(), role2Key).
		Return(testutil.ToUnstructured(t, role2), nil)

	saph := NewServiceAccountPolicyRules(ctx, serviceAccount, printOptions)
	got, err := saph.Create()
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
		{
			"Resources":         component.NewText("pods"),
			"Non-Resource URLs": component.NewText(""),
			"Resource Names":    component.NewText(""),
			"Verbs":             component.NewText("['get', 'watch', 'list']"),
		},
	}...)

	component.AssertEqual(t, expected, got)
}
