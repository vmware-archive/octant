package printer

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/golang/mock/gomock"

	"github.com/heptio/developer-dash/internal/cache"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	fakecache "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ServiceAccountListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
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

	got, err := ServiceAccountListHandler(list, printOptions)
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

func Test_printServiceAccountConfig(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	c := fakecache.NewMockCache(controller)

	now := time.Unix(1547211430, 0)

	object := testutil.CreateServiceAccount("sa")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Secrets = []corev1.ObjectReference{{Name: "secret"}}
	object.ImagePullSecrets = []corev1.LocalObjectReference{{Name: "secret"}}

	key := cache.Key{
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

	c.EXPECT().List(gomock.Eq(key)).
		Return([]*unstructured.Unstructured{testutil.ToUnstructured(t, secret)}, nil)

	got, err := printServiceAccountConfig(object, c)
	require.NoError(t, err)

	var sections component.SummarySections
	pullSecretsList := component.NewList("", []component.ViewComponent{
		link.ForGVK(object.Namespace, "v1", "Secret", "secret", "secret"),
	})
	sections.Add("Image Pull Secrets", pullSecretsList)

	mountSecretsList := component.NewList("", []component.ViewComponent{
		link.ForGVK(object.Namespace, "v1", "Secret", "secret", "secret"),
	})
	sections.Add("Mountable Secrets", mountSecretsList)

	tokenSecretsList := component.NewList("", []component.ViewComponent{
		link.ForGVK(object.Namespace, "v1", "Secret", "secret", "secret"),
	})
	sections.Add("Tokens", tokenSecretsList)

	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, got)
}
