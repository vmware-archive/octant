package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	printerfake "github.com/heptio/developer-dash/internal/overview/printer/fake"
	pffake "github.com/heptio/developer-dash/internal/portforward/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
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
		isNil        bool
		expectedType string
	}{
		{
			name: "print known object",
			printFunc: func(ctx context.Context, deployment *appsv1.Deployment, options Options) (component.Component, error) {
				return &stubComponent{Type: "type1"}, nil
			},
			object:       &appsv1.Deployment{},
			expectedType: "type1",
		},
		{
			name:   "print unregistered type returns error",
			object: &appsv1.Deployment{},
			isNil:  true,
		},
		{
			name:   "print unregistered list type runs a nil",
			object: &appsv1.DeploymentList{},
			isNil:  true,
		},
		{
			name: "print handler returns error",
			printFunc: func(ctx context.Context, deployment *appsv1.Deployment, options Options) (component.Component, error) {
				return nil, errors.New("failed")
			},
			object: &appsv1.Deployment{},
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			pf := pffake.NewMockPortForwarder(controller)
			pluginPrinter := printerfake.NewMockPluginPrinter(controller)

			p := NewResource(o, pf)

			if tc.printFunc != nil {
				err := p.Handler(tc.printFunc)
				require.NoError(t, err)
			}

			ctx := context.Background()
			got, err := p.Print(ctx, tc.object, pluginPrinter)
			if tc.isErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tc.isNil {
				assert.Nil(t, got)
				return
			}
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
			printFunc: func(context.Context, int, Options) (component.Component, error) {
				return &stubComponent{Type: "type1"}, nil
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
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			pf := pffake.NewMockPortForwarder(controller)

			p := NewResource(o, pf)

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
	printFunc := func(context.Context, int, Options) (component.Component, error) {
		return &stubComponent{Type: "type1"}, nil
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)
	pf := pffake.NewMockPortForwarder(controller)

	p := NewResource(o, pf)

	err := p.Handler(printFunc)
	require.NoError(t, err)

	err = p.Handler(printFunc)
	assert.Error(t, err)

}

type stubComponent struct {
	Type string
}

var _ component.Component = (*stubComponent)(nil)

func (v *stubComponent) GetMetadata() component.Metadata {
	return component.Metadata{
		Type: v.Type,
	}
}

func (v *stubComponent) IsEmpty() bool {
	return false
}

func (v *stubComponent) SetAccessor(string) {}

func (v *stubComponent) String() string { return "" }

func Test_DefaultPrinter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		ObjectStore: storefake.NewMockObjectStore(controller),
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

	ctx := context.Background()
	got, err := DefaultPrintFunc(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Age")
	expected := component.NewTable("/v1, Kind=DeploymentList", cols)
	expected.Add(component.TableRow{
		"Name":   component.NewText("deployment"),
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
			controller := gomock.NewController(t)
			defer controller.Finish()

			printOptions := Options{
				ObjectStore: storefake.NewMockObjectStore(controller),
			}

			ctx := context.Background()
			_, err := DefaultPrintFunc(ctx, tc.object, printOptions)
			assert.Error(t, err)
		})
	}

}
