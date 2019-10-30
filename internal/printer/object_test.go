/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/fake"
	"github.com/vmware-tanzu/octant/pkg/store"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
)

func Test_Object_ToComponent(t *testing.T) {
	type initOptions struct {
		Options       *Options
		PluginPrinter *fake.MockManagerInterface
	}

	deployment := testutil.CreateDeployment("deployment")

	defaultConfig := component.NewSummary("Configuration",
		component.SummarySection{Header: "local"})

	defaultConfigSection := component.FlexLayoutSection{
		{
			Width: component.WidthHalf,
			View:  defaultConfig,
		},
	}

	metadataSection := component.FlexLayoutSection{
		{
			Width: component.WidthHalf,
			View:  component.NewText("metadata"),
		},
	}

	fnMetadata := func(o *Object) {
		o.MetadataGen = func(object runtime.Object, fl *flexlayout.FlexLayout, options Options) error {
			section := fl.AddSection()
			require.NoError(t, section.Add(component.NewText("metadata"), 12))
			return nil
		}
	}

	fnPodTemplate := func(o *Object) {
		o.PodTemplateGen = func(_ context.Context, _ runtime.Object, _ corev1.PodTemplateSpec, fl *flexlayout.FlexLayout, options Options) error {
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

	stubPlugins := func(pluginPrinter *fake.MockManagerInterface) {
		printResponse := &plugin.PrintResponse{}
		pluginPrinter.EXPECT().
			Print(gomock.Any(), gomock.Any()).Return(printResponse, nil)
	}

	cases := []struct {
		name     string
		object   runtime.Object
		initFunc func(*Object, *initOptions)
		sections []component.FlexLayoutSection
		buttons  []component.Button
		isErr    bool
	}{
		{
			name:   "in general",
			object: deployment,
			initFunc: func(o *Object, options *initOptions) {
				stubPlugins(options.PluginPrinter)
			},
			sections: []component.FlexLayoutSection{
				defaultConfigSection,
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
					Print(gomock.Any(), gomock.Any()).Return(&printResponse, nil)
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
				stubPlugins(options.PluginPrinter)
			},
			sections: []component.FlexLayoutSection{
				defaultConfigSection,
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
				stubPlugins(options.PluginPrinter)
			},
			sections: []component.FlexLayoutSection{
				defaultConfigSection,
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
			name:   "add button",
			object: deployment,
			initFunc: func(object *Object, options *initOptions) {
				stubPlugins(options.PluginPrinter)
				object.AddButton("button name", action.Payload{})
			},
			buttons: []component.Button{
				component.NewButton("button name", action.Payload{}),
			},
			sections: []component.FlexLayoutSection{
				defaultConfigSection,
				metadataSection,
			},
		},
		{
			name:   "register items",
			object: deployment,
			initFunc: func(o *Object, options *initOptions) {
				stubPlugins(options.PluginPrinter)
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
				defaultConfigSection,
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

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			o := NewObject(tc.object, fnMetadata, fnPodTemplate, fnEvent)

			o.RegisterConfig(defaultConfig)

			if tc.initFunc != nil {
				options := &initOptions{
					Options:       &printOptions,
					PluginPrinter: tpo.pluginManager,
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
			buttonGroup := component.NewButtonGroup()
			if len(tc.buttons) > 0 {
				for _, button := range tc.buttons {
					buttonGroup.AddButton(button)
				}
			}

			key, err := store.KeyFromObject(deployment)
			require.NoError(t, err)

			buttonGroup.AddButton(
				component.NewButton("Delete",
					action.CreatePayload(octant.ActionDeleteObject, key.ToActionPayload()),
					component.WithButtonConfirmation(
						"Delete Deployment",
						"Are you sure you want to delete *Deployment* **deployment**? This action is permanent and cannot be recovered.",
					)))
			expected.SetButtonGroup(buttonGroup)

			component.AssertEqual(t, expected, got)
		})
	}
}

func Test_deleteObjectConfirmation(t *testing.T) {
	pod := testutil.CreatePod("pod")
	option, err := deleteObjectConfirmation(pod)
	require.NoError(t, err)

	button := component.Button{}
	option(&button)

	expected := component.Button{
		Confirmation: &component.Confirmation{
			Title: "Delete Pod",
			Body:  "Are you sure you want to delete *Pod* **pod**? This action is permanent and cannot be recovered.",
		},
	}

	assert.Equal(t, expected, button)
}
