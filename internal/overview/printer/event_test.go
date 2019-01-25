package printer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
)

func Test_EventListHandler(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	object := &corev1.EventList{
		Items: []corev1.Event{
			{
				InvolvedObject: corev1.ObjectReference{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "d1",
				},
				Count:          1234,
				Message:        "message",
				Reason:         "Reason",
				Type:           "Type",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
		},
	}

	got, err := printer.EventListHandler(object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Kind", "Message", "Reason", "Type",
		"First Seen", "Last Seen")
	expected := component.NewTable("Events", cols)
	expected.Add(component.TableRow{
		"Kind": component.NewList("", []component.ViewComponent{
			component.NewLink("", "d1", "/content/overview/workloads/deployments/d1"),
			component.NewText("", "1234"),
		}),
		"Message":    component.NewText("", "message"),
		"Reason":     component.NewText("", "Reason"),
		"Type":       component.NewText("", "Type"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
	})

	assert.Equal(t, expected, got)
}
