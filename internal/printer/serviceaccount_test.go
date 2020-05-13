/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
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

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	linkFake "github.com/vmware-tanzu/octant/internal/link/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
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
		"Name": component.NewLink("", object.Name, "/path",
			genObjectStatus(component.TextStatusOK, []string{
				"v1 ServiceAccount is OK",
			})),
		"Labels":  component.NewLabels(labels),
		"Secrets": component.NewText("1"),
		"Age":     component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
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

	tests := []struct {
		name           string
		namespace      string
		serviceAccount *corev1.ServiceAccount
		secret         *corev1.Secret
		isErr          bool
		expected       *component.Summary
	}{
		{
			name:           "in general",
			namespace:      serviceAccount.Namespace,
			serviceAccount: serviceAccount,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header: "Image Pull Secrets",
					Content: component.NewList([]component.TitleComponent{}, []component.Component{
						component.NewLink("", "secret", "/secret"),
					}),
				},
				{
					Header: "Tokens",
					Content: component.NewList([]component.TitleComponent{}, []component.Component{
						component.NewLink("", "secret", "/secret"),
					}),
				},
			}...),
		},
		{
			name:           "service account is nil",
			serviceAccount: nil,
			isErr:          true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			ctx := context.Background()
			tpo := newTestPrinterOptions(controller)

			if !test.isErr {
				tpo.objectStore.EXPECT().List(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructuredList(t, testutil.ToUnstructured(t, secret)), false, nil)
			}

			printOptions := tpo.ToOptions()

			tpo.PathForGVK(test.namespace, "v1", "Secret", "secret", "secret", "/secret")

			sac := NewServiceAccountConfiguration(ctx, test.serviceAccount, printOptions)
			summary, err := sac.Create(printOptions)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, test.expected, summary)
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

	policyRules := NewServiceAccountPolicyRules(ctx, serviceAccount, printOptions)
	got, err := policyRules.Create()
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

func TestServiceAccountSecrets(t *testing.T) {
	tests := []struct {
		name        string
		account     *corev1.ServiceAccount
		initOptions func(controller *gomock.Controller) Options
		want        *component.Table
		wantError   bool
	}{
		{
			name:    "service account with no secrets",
			account: testutil.CreateServiceAccount("sa"),
			initOptions: func(controller *gomock.Controller) Options {
				return Options{}
			},
			want: nil,
		},
		{
			name: "service account with secret",
			account: testutil.CreateServiceAccount("sa", func(account *corev1.ServiceAccount) {
				account.Secrets = []corev1.ObjectReference{
					{
						Name: "secret",
					},
				}
			}),
			initOptions: func(ctrl *gomock.Controller) Options {
				secret := testutil.ToUnstructured(t, testutil.CreateSecret("secret", func(secret *corev1.Secret) {
					secret.Type = corev1.SecretTypeOpaque
				}))
				secretKey, err := store.KeyFromObject(secret)
				require.NoError(t, err)

				objectStore := storeFake.NewMockStore(ctrl)
				objectStore.EXPECT().
					Get(gomock.Any(), secretKey).
					Return(secret, nil)

				secretLink := component.NewLink("", "secret", "/secret")
				link := linkFake.NewMockInterface(ctrl)
				link.EXPECT().
					ForObject(secret, secret.GetName()).
					Return(secretLink, nil)

				dashConfig := configFake.NewMockDash(ctrl)
				dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

				options := Options{
					DashConfig: dashConfig,
					Link:       link,
				}

				return options
			},
			want: component.NewTableWithRows(
				"Secrets",
				ServiceAccountSecretPlaceholder,
				ServiceAccountSecretCols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "secret", "/secret"),
						"Type": component.NewText(string(corev1.SecretTypeOpaque)),
					},
				}),
		},
		{
			name:      "with nil service account",
			account:   nil,
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			var options Options
			if test.initOptions != nil {
				options = test.initOptions(ctrl)
			}

			actual, err := ServiceAccountSecrets(ctx, test.account, options)
			if test.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			testutil.AssertJSONEqual(t, test.want, actual)
		})
	}
}
