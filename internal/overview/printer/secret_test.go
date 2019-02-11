package printer

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_SecretListHandler(t *testing.T) {
	printOptions := Options{
		Cache: cache.NewMemoryCache(),
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
					Name: "secret",
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
		"Name":   component.NewLink("", "secret", "/content/overview/config-and-storage/secrets/secret"),
		"Labels": component.NewLabels(labels),
		"Type":   component.NewText("Opaque"),
		"Data":   component.NewText("1"),
		"Age":    component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}

func Test_SecretHandler(t *testing.T) {
	printOptions := Options{
		Cache: cache.NewMemoryCache(),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "secret",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			Labels: labels,
		},
		Data: map[string][]byte{
			"key": []byte("value"),
		},
		Type: corev1.SecretTypeOpaque,
	}

	got, err := SecretHandler(secret, printOptions)
	require.NoError(t, err)

	config := component.NewSummary("Configuration", []component.SummarySection{
		{
			Header:  "Type",
			Content: component.NewText("Opaque"),
		},
	}...)
	configPanel := component.NewPanel("", config)
	configPanel.Position(0, 0, 12, 8)

	data := component.NewTable("Data", secretDataCols)
	data.Add(component.TableRow{
		"Key": component.NewText("key"),
	})

	dataPanel := component.NewPanel("", data)
	dataPanel.Position(0, 9, 24, 8)

	panels := []component.Panel{
		*configPanel,
		*dataPanel,
	}

	expected := component.NewGrid("Summary", panels...)

	assert.Equal(t, expected, got)
}
