package printer_test

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_Resource_Print(t *testing.T) {
	cases := []struct {
		name         string
		printFunc    interface{}
		object       runtime.Object
		isErr        bool
		expectedType string
	}{
		{
			name: "print known object",
			printFunc: func(deployment *appsv1.Deployment, options printer.Options) (component.ViewComponent, error) {
				return &stubViewComponent{Type: "type1"}, nil
			},
			object:       &appsv1.Deployment{},
			expectedType: "type1",
		},
		{
			name:         "print unregistered type returns text component",
			object:       &appsv1.Deployment{},
			expectedType: "text",
		},
		{
			name: "print handler returned error",
			printFunc: func(deployment *appsv1.Deployment, options printer.Options) (component.ViewComponent, error) {
				return nil, errors.New("failed")
			},
			object: &appsv1.Deployment{},
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := cache.NewMemoryCache()
			p := printer.NewResource(c)

			if tc.printFunc != nil {
				err := p.Handler(tc.printFunc)
				require.NoError(t, err)
			}

			got, err := p.Print(tc.object)
			if tc.isErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedType, got.GetMetadata().Type)

		})
	}

}

func Test_Resource_Handler(t *testing.T) {
	cases := []struct {
		name      string
		printFunc interface{}
		isErr     bool
	}{
		{
			name: "valid printer",
			printFunc: func(int, printer.Options) (component.ViewComponent, error) {
				return &stubViewComponent{Type: "type1"}, nil
			},
		},
		{
			name:      "non function printer",
			printFunc: nil,
			isErr:     true,
		},
		{
			name:      "print func invalid in/out count",
			printFunc: func() {},
			isErr:     true,
		},
		{
			name:      "print func invalid signature",
			printFunc: func(int) (int, error) { return 0, nil },
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := cache.NewMemoryCache()
			p := printer.NewResource(c)

			err := p.Handler(tc.printFunc)

			if tc.isErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_Resource_DuplicateHandler(t *testing.T) {
	printFunc := func(int, printer.Options) (component.ViewComponent, error) {
		return &stubViewComponent{Type: "type1"}, nil
	}

	c := cache.NewMemoryCache()
	p := printer.NewResource(c)

	err := p.Handler(printFunc)
	require.NoError(t, err)

	err = p.Handler(printFunc)
	assert.Error(t, err)

}

type stubViewComponent struct {
	Type string
}

var _ component.ViewComponent = (*stubViewComponent)(nil)

func (v *stubViewComponent) GetMetadata() component.Metadata {
	return component.Metadata{
		Type: v.Type,
	}
}

func (v *stubViewComponent) IsEmpty() bool {
	return false
}

func Test_DefaultPrinter(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
			},
		},
	}

	got, err := printer.DefaultPrintFunc(object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Age")
	expected := component.NewTable("*v1.DeploymentList", cols)
	expected.Add(component.TableRow{
		"Name":   component.NewText("", "deployment"),
		"Labels": component.NewLabels(labels),
		"Age":    component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}

func Test_DefaultPrinter_invalid_object(t *testing.T) {
	cases := []struct {
		name   string
		object runtime.Object
	}{
		{
			name:   "nil object",
			object: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			printOptions := printer.Options{
				Cache: cache.NewMemoryCache(),
			}

			_, err := printer.DefaultPrintFunc(tc.object, printOptions)
			assert.Error(t, err)
		})
	}

}
