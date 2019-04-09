package printer

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/objectstore"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/overview/link"
	printerfake "github.com/heptio/developer-dash/internal/overview/printer/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/cacheutil"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_ServiceAccountListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		ObjectStore: storefake.NewMockObjectStore(controller),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreateServiceAccount("sa")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels
	object.Secrets = []corev1.ObjectReference{{Name: "secret"}}

	list := &corev1.ServiceAccountList{
		Items: []corev1.ServiceAccount{*object},
	}

	ctx := context.Background()
	got, err := ServiceAccountListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Secrets", "Age")
	expected := component.NewTable("Service Accounts", cols)
	expected.Add(component.TableRow{
		"Name":    link.ForObject(object, object.Name),
		"Labels":  component.NewLabels(labels),
		"Secrets": component.NewText("1"),
		"Age":     component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}

func Test_serviceAccountHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	now := time.Unix(1547211430, 0)
	object := testutil.CreateServiceAccount("sa")
	object.CreationTimestamp = metav1.Time{Time: now}

	pluginPrinter := printerfake.NewMockPluginPrinter(controller)
	pluginPrinter.EXPECT().
		Print(gomock.Any()).
		Return(&plugin.PrintResponse{}, nil)

	appObjectStore := storefake.NewMockObjectStore(controller)
	mockObjectsEvents(t, appObjectStore, object.Namespace)

	printOptions := Options{
		ObjectStore:   appObjectStore,
		PluginPrinter: pluginPrinter,
	}

	ctx := context.Background()

	h, err := newServiceAccountHandler(ctx, object, printOptions)
	require.NoError(t, err)

	summaryConfig := component.NewSummary("config", component.SummarySection{
		Header: "foo", Content: component.NewText("bar")})

	h.configFunc = func(ctx context.Context, serviceAccount corev1.ServiceAccount, o objectstore.ObjectStore) (*component.Summary, error) {
		return summaryConfig, nil
	}

	policyTable := component.NewTable("policyTable", component.NewTableCols("col1"))
	h.policyRulesFunc = func(ctx context.Context, serviceAccount corev1.ServiceAccount, appObjectStore objectstore.ObjectStore) (*component.Table, error) {
		return policyTable, nil
	}

	got, err := h.run()
	require.NoError(t, err)

	fl := flexlayout.New()
	summarySection := fl.AddSection()
	require.NoError(t, summarySection.Add(summaryConfig, component.WidthHalf))

	stubMetadataForObject(t, object, fl)

	policyRulesSection := fl.AddSection()
	require.NoError(t, policyRulesSection.Add(policyTable, component.WidthFull))

	expected := fl.ToComponent("Summary")

	assert.Equal(t, expected, got)
}

func Test_printServiceAccountConfig(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)

	now := time.Unix(1547211430, 0)

	object := testutil.CreateServiceAccount("sa")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Secrets = []corev1.ObjectReference{{Name: "secret"}}
	object.ImagePullSecrets = []corev1.LocalObjectReference{{Name: "secret"}}

	key := cacheutil.Key{
		Namespace:  object.Namespace,
		APIVersion: "v1",
		Kind:       "Secret",
	}

	secret := testutil.CreateSecret("secret")
	secret.Type = corev1.SecretTypeServiceAccountToken
	secret.Annotations = map[string]string{
		corev1.ServiceAccountNameKey: object.Name,
		corev1.ServiceAccountUIDKey:  string(object.UID),
	}

	o.EXPECT().List(gomock.Any(), gomock.Eq(key)).
		Return([]*unstructured.Unstructured{testutil.ToUnstructured(t, secret)}, nil)

	ctx := context.Background()
	got, err := printServiceAccountConfig(ctx, *object, o)
	require.NoError(t, err)

	var sections component.SummarySections
	pullSecretsList := component.NewList("", []component.Component{
		link.ForGVK(object.Namespace, "v1", "Secret", "secret", "secret"),
	})
	sections.Add("Image Pull Secrets", pullSecretsList)

	mountSecretsList := component.NewList("", []component.Component{
		link.ForGVK(object.Namespace, "v1", "Secret", "secret", "secret"),
	})
	sections.Add("Mountable Secrets", mountSecretsList)

	tokenSecretsList := component.NewList("", []component.Component{
		link.ForGVK(object.Namespace, "v1", "Secret", "secret", "secret"),
	})
	sections.Add("Tokens", tokenSecretsList)

	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, got)
}

func Test_serviceAccountPolicyRules(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	serviceAccount := testutil.CreateServiceAccount("sa")

	appObjectStore := storefake.NewMockObjectStore(controller)

	roleBindingKey := cacheutil.Key{
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

	appObjectStore.EXPECT().
		List(gomock.Any(), roleBindingKey).
		Return(roleBindingObjects, nil)

	clusterRoleBindingKey := cacheutil.Key{
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

	appObjectStore.EXPECT().
		List(gomock.Any(), clusterRoleBindingKey).
		Return(clusterRoleBindingObjects, nil)

	role1Key, err := cacheutil.KeyFromObject(role1)
	require.NoError(t, err)

	appObjectStore.EXPECT().
		Get(gomock.Any(), role1Key).
		Return(testutil.ToUnstructured(t, role1), nil)

	role2Key, err := cacheutil.KeyFromObject(role2)
	require.NoError(t, err)

	spew.Dump(role2Key)

	appObjectStore.EXPECT().
		Get(gomock.Any(), role2Key).
		Return(testutil.ToUnstructured(t, role2), nil)

	ctx := context.Background()

	s := newServiceAccountPolicyRules(ctx, *serviceAccount, appObjectStore)

	var policyRules []rbacv1.PolicyRule
	s.printPolicyRulesFunc = func(rules []rbacv1.PolicyRule) (*component.Table, error) {
		policyRules = rules
		return nil, nil
	}

	_, err = s.run()
	require.NoError(t, err)

	var expected []rbacv1.PolicyRule
	expected = append(expected, role1.Rules...)
	expected = append(expected, role2.Rules...)

	assert.Equal(t, expected, policyRules)
}
