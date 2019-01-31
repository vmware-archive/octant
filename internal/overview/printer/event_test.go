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
				ObjectMeta: metav1.ObjectMeta{
					Name: "event-12345",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Event",
				},
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
		"Message":    component.NewLink("", "message", "/content/overview/events/event-12345"),
		"Reason":     component.NewText("", "Reason"),
		"Type":       component.NewText("", "Type"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
	})

	assert.Equal(t, expected, got)
}

func Test_ReplicasetEvents(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	now := time.Unix(1547211430, 0)

	object := &corev1.EventList{
		Items: []corev1.Event{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "frontend",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
				},
				Count:  1,
				Type:   corev1.EventTypeNormal,
				Reason: "SuccessfulCreate",
				Source: corev1.EventSource{
					Component: "replicaset-controller",
				},
				Message:        "Created pod: frontend-97k6z",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "frontend",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
				},
				Count:  1,
				Type:   corev1.EventTypeNormal,
				Reason: "SuccessfulCreate",
				Source: corev1.EventSource{
					Component: "replicaset-controller",
				},
				Message:        "Created pod: frontend-8n77p",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "frontend",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
				},
				Count:  1,
				Type:   corev1.EventTypeNormal,
				Reason: "SuccessfulCreate",
				Source: corev1.EventSource{
					Component: "replicaset-controller",
				},
				Message:        "Created pod: frontend-b7fxf",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
		},
	}

	got, err := printer.PrintEvents(object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Type", "Reason", "Age", "From", "Message")
	expected := component.NewTable("Events", cols)

	expected.Add(component.TableRow{
		"Message":    component.NewText("", "Created pod: frontend-97k6z"),
		"Reason":     component.NewText("", "SuccessfulCreate"),
		"Type":       component.NewText("", "Normal"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
		"From":       component.NewText("", "replicaset-controller"),
		"Count":      component.NewText("", "1"),
	})

	expected.Add(component.TableRow{
		"Message":    component.NewText("", "Created pod: frontend-8n77p"),
		"Reason":     component.NewText("", "SuccessfulCreate"),
		"Type":       component.NewText("", "Normal"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
		"From":       component.NewText("", "replicaset-controller"),
		"Count":      component.NewText("", "1"),
	})

	expected.Add(component.TableRow{
		"Message":    component.NewText("", "Created pod: frontend-b7fxf"),
		"Reason":     component.NewText("", "SuccessfulCreate"),
		"Type":       component.NewText("", "Normal"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
		"From":       component.NewText("", "replicaset-controller"),
		"Count":      component.NewText("", "1"),
	})

	assert.Equal(t, expected, got)
}

func Test_EventHandler(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	event := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name: "event-12345",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Event",
		},
		InvolvedObject: corev1.ObjectReference{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       "d1",
		},
		Count:          1234,
		Message:        "message",
		Reason:         "Reason",
		Type:           corev1.EventTypeNormal,
		FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
		LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
		Source: corev1.EventSource{
			Component: "component",
			Host:      "host",
		},
	}

	got, err := printer.EventHandler(event, printOptions)
	require.NoError(t, err)

	eventDetailSections := []component.SummarySection{
		{
			Header:  "Last Seen",
			Content: component.NewTimestamp(time.Unix(1548424410, 0)),
		},
		{
			Header:  "First Seen",
			Content: component.NewTimestamp(time.Unix(1548424410, 0)),
		},
		{
			Header:  "Count",
			Content: component.NewText("", "1234"),
		},
		{
			Header:  "Message",
			Content: component.NewText("", "message"),
		},
		{
			Header:  "Kind",
			Content: component.NewText("", "Deployment"),
		},
		{
			Header:  "Involved Object",
			Content: component.NewLink("", "d1", "/content/overview/workloads/deployments/d1"),
		},
		{
			Header:  "Type",
			Content: component.NewText("", "Normal"),
		},
		{
			Header:  "Reason",
			Content: component.NewText("", "Reason"),
		},
		{
			Header:  "Source",
			Content: component.NewText("", "component on host"),
		},
	}
	eventDetailView := component.NewSummary("Event Detail", eventDetailSections...)

	eventDetailPanel := component.NewPanel("", eventDetailView)
	eventDetailPanel.Position(0, 0, 24, 10)

	expected := component.NewGrid("", *eventDetailPanel)

	assert.Equal(t, expected, got)
}
