package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	printerfake "github.com/heptio/developer-dash/internal/overview/printer/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_Object_ToComponent(t *testing.T) {
	type initOptions struct {
		Options       *Options
		PluginPrinter *printerfake.MockPluginPrinter
	}

	deployment := testutil.CreateDeployment("deployment")

	defaultConfig := component.NewSummary("Configuration",
		component.SummarySection{Header: "local"})

	metadataSection := component.FlexLayoutSection{
		{
			Width: component.WidthHalf,
			View:  component.NewText("metadata"),
		},
	}

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

	mockNoPlugins := func(pluginPrinter *printerfake.MockPluginPrinter) {
		printResponse := &plugin.PrintResponse{}
		pluginPrinter.EXPECT().
			Print(gomock.Any()).Return(printResponse, nil)
	}

	cases := []struct {
		name     string
		object   runtime.Object
		initFunc func(*Object, *initOptions)
		sections []component.FlexLayoutSection
		isErr    bool
	}{
		{
			name:   "in general",
			object: deployment,
			initFunc: func(o *Object, options *initOptions) {
				mockNoPlugins(options.PluginPrinter)
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: component.WidthHalf,
						View:  defaultConfig,
					},
				},
				metadataSection,
			},
		},
		{
			name:   "config data from plugin",
			object: deployment,
			initFunc: func(o *Object, options *initOptions) {
				printResponse := plugin.PrintResponse{
					Config: []component.SummarySection{
						{Header: "from plugin"},
					},
				}

				options.PluginPrinter.EXPECT().
					Print(gomock.Any()).Return(&printResponse, nil)
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: component.WidthHalf,
						View: component.NewSummary("Configuration",
							[]component.SummarySection{
								{Header: "local"},
								{Header: "from plugin"},
							}...),
					},
				},
				metadataSection,
			},
		},
		{
			name:   "enable pod template",
			object: deployment,
			initFunc: func(o *Object, options *initOptions) {
				o.EnablePodTemplate(deployment.Spec.Template)
				mockNoPlugins(options.PluginPrinter)
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: component.WidthHalf,
						View:  defaultConfig,
					},
				},
				metadataSection,
				{
					{
						Width: component.WidthHalf,
						View:  component.NewText("pod template"),
					},
				},
			},
		},
		{
			name:   "enable events",
			object: deployment,
			initFunc: func(o *Object, options *initOptions) {
				o.EnableEvents()
				mockNoPlugins(options.PluginPrinter)
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: component.WidthHalf,
						View:  defaultConfig,
					},
				},
				metadataSection,
				{
					{
						Width: component.WidthHalf,
						View:  component.NewText("events"),
					},
				},
			},
		},
		{
			name:   "register items",
			object: deployment,
			initFunc: func(o *Object, options *initOptions) {
				mockNoPlugins(options.PluginPrinter)
				o.RegisterItems([]ItemDescriptor{
					{
						Func: func() (component.Component, error) {
							return component.NewText("item1"), nil
						},
						Width: component.WidthHalf,
					},
					{
						Func: func() (component.Component, error) {
							return component.NewText("item2"), nil
						},
						Width: component.WidthHalf,
					},
				}...)
				o.RegisterItems(ItemDescriptor{
					Func: func() (component.Component, error) {
						return component.NewText("item3"), nil
					},
					Width: component.WidthHalf,
				})
			},
			sections: []component.FlexLayoutSection{
				{
					{
						Width: component.WidthHalf,
						View:  defaultConfig,
					},
				},
				metadataSection,
				{
					{
						Width: component.WidthHalf,
						View:  component.NewText("item1"),
					},
					{
						Width: component.WidthHalf,
						View:  component.NewText("item2"),
					},
				},
				{
					{
						Width: component.WidthHalf,
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

			pluginPrinter := printerfake.NewMockPluginPrinter(controller)

			printOptions := Options{
				ObjectStore:   storefake.NewMockObjectStore(controller),
				PluginPrinter: pluginPrinter,
			}

			o := NewObject(tc.object, fnMetdata, fnPodTemplate, fnEvent)

			o.RegisterConfig(defaultConfig)

			if tc.initFunc != nil {
				options := &initOptions{
					Options:       &printOptions,
					PluginPrinter: pluginPrinter,
				}
				tc.initFunc(o, options)
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
