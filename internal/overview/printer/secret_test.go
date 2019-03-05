package printer

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_SecretListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := &corev1.SecretList{
		Items: []corev1.Secret{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
				Data: map[string][]byte{
					"key": []byte("value"),
				},
				Type: corev1.SecretTypeOpaque,
			},
		},
	}

	got, err := SecretListHandler(object, printOptions)
	require.NoError(t, err)

	expected := component.NewTable("Secrets", secretTableCols)
	expected.Add(component.TableRow{
		"Name":   component.NewLink("", "secret", "/content/overview/namespace/default/config-and-storage/secrets/secret"),
		"Labels": component.NewLabels(labels),
		"Type":   component.NewText("Opaque"),
		"Data":   component.NewText("1"),
		"Age":    component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}
