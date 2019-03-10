package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_Object_ToComponent(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")

	fnMetdata := func(o *Object) {
		o.MetadataGen = func(object runtime.Object, fl *flexlayout.FlexLayout) error {
			section := fl.AddSection()
			require.NoError(t, section.Add(component.NewText("metadata"), 12))
			return nil
		}
	}

	fnPodTemplate := func(o *Object) {
		o.PodTemplateGen = func(_ runtime.Object, _ corev1.PodTemplateSpec, fl *flexlayout.FlexLayout, options Options) error {
			section := fl.AddSection()
			require.NoError(t, section.Add(component.NewText("pod template"), 12))
			return nil
		}
	}

	fnEvent := func(o *Object) {
		o.EventsGen = func(_ context.Context, _ runtime.Object, fl *flexlayout.FlexLayout, _ Options) error {
			section := fl.AddSection()
			require.NoError(t, section.Add(component.NewText("events"), 12))
			return nil
		}
	}

	cases := []struct {
		name     string
		object   runtime.Object
		initFunc func(*Object)
		sections []component.FlexLayoutSection
		isErr    bool
	}{
		{
			name:   "in general",
			object: deployment,
			sections: []component.FlexLayoutSection{
				{
					{
						Width: 12,
						View:  component.NewText("config"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("metadata"),
					},
				},
			},
		},
		{
			name:   "extra summary items",
			object: deployment,
			initFunc: func(o *Object) {
				o.RegisterSummary(func() (component.ViewComponent, error) {
					return component.NewText("summary object 1"), nil
				}, 12)
				o.RegisterSummary(func() (component.ViewComponent, error) {
					return component.NewText("summary object 2"), nil
				}, 12)
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: 12,
						View:  component.NewText("config"),
					},
					{
						Width: 12,
						View:  component.NewText("summary object 1"),
					},
					{
						Width: 12,
						View:  component.NewText("summary object 2"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("metadata"),
					},
				},
			},
		},
		{
			name:   "enable pod template",
			object: deployment,
			initFunc: func(o *Object) {
				o.EnablePodTemplate(deployment.Spec.Template)
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: 12,
						View:  component.NewText("config"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("metadata"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("pod template"),
					},
				},
			},
		},
		{
			name:   "enable events",
			object: deployment,
			initFunc: func(o *Object) {
				o.EnableEvents()
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: 12,
						View:  component.NewText("config"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("metadata"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("events"),
					},
				},
			},
		},
		{
			name:   "register items",
			object: deployment,
			initFunc: func(o *Object) {
				o.RegisterItems([]ItemDescriptor{
					{
						Func: func() (component.ViewComponent, error) {
							return component.NewText("item1"), nil
						},
						Width: 12,
					},
					{
						Func: func() (component.ViewComponent, error) {
							return component.NewText("item2"), nil
						},
						Width: 12,
					},
				}...)
				o.RegisterItems(ItemDescriptor{
					Func: func() (component.ViewComponent, error) {
						return component.NewText("item3"), nil
					},
					Width: 12,
				})
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: 12,
						View:  component.NewText("config"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("metadata"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("item1"),
					},
					{
						Width: 12,
						View:  component.NewText("item2"),
					},
				},
				{
					{
						Width: 12,
						View:  component.NewText("item3"),
					},
				},
			},
		},
		{
			name:   "nil object",
			object: nil,
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			printOptions := Options{
				Cache: cachefake.NewMockCache(controller),
			}

			o := NewObject(tc.object, fnMetdata, fnPodTemplate, fnEvent)

			o.RegisterConfig(func() (component.ViewComponent, error) {
				return component.NewText("config"), nil
			}, 12)

			if tc.initFunc != nil {
				tc.initFunc(o)
			}

			ctx := context.Background()
			got, err := o.ToComponent(ctx, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected := component.NewFlexLayout("Summary")
			expected.AddSections(tc.sections...)

			assert.Equal(t, expected, got)

		})
	}

}
