/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestPodTemplateHeader(t *testing.T) {
	labels := map[string]string{
		"key": "value",
	}

	pth := NewPodTemplateHeader(labels)
	got := pth.Create()

	assert.Len(t, got.Config.Labels, 1)

	expected := component.Title(component.NewText("Pod Template"))

	assert.Equal(t, expected, got.Metadata.Title)
}

func stubPodTemplateSection(name string) podTemplateFunc {
	return func(ctx context.Context, fl *flexlayout.FlexLayout, options podTemplateLayoutOptions) error {
		section := fl.AddSection()
		return section.Add(component.NewText(name), component.WidthFull)
	}
}

func stubPodTemplateSectionWithError() podTemplateFunc {
	return func(ctx context.Context, fl *flexlayout.FlexLayout, options podTemplateLayoutOptions) error {
		return errors.Errorf("failed")
	}
}

func TestPodTemplate_AddToFlexLayout(t *testing.T) {
	cases := []struct {
		name                            string
		flexlayout                      *flexlayout.FlexLayout
		podTemplateHeaderFunc           podTemplateFunc
		podTemplateInitContainersFunc   podTemplateFunc
		podTemplateContainersFunc       podTemplateFunc
		podTemplatePodConfigurationFunc podTemplateFunc
		isErr                           bool
		expected                        func() *component.FlexLayout
	}{
		{
			name:                            "in general",
			flexlayout:                      flexlayout.New(),
			podTemplateHeaderFunc:           stubPodTemplateSection("header"),
			podTemplateInitContainersFunc:   stubPodTemplateSection("init containers"),
			podTemplateContainersFunc:       stubPodTemplateSection("containers"),
			podTemplatePodConfigurationFunc: stubPodTemplateSection("configuration"),
			expected: func() *component.FlexLayout {
				expected := component.NewFlexLayout("Pod Template")
				expected.AddSections([]component.FlexLayoutSection{
					{
						{
							Width: component.WidthFull,
							View:  component.NewText("header"),
						},
					},
					{
						{
							Width: component.WidthFull,
							View:  component.NewText("init containers"),
						},
					},
					{
						{
							Width: component.WidthFull,
							View:  component.NewText("containers"),
						},
					},
					{
						{
							Width: component.WidthFull,
							View:  component.NewText("configuration"),
						},
					},
				}...)

				return expected
			},
		},
		{
			name:                            "header error",
			flexlayout:                      flexlayout.New(),
			podTemplateHeaderFunc:           stubPodTemplateSectionWithError(),
			podTemplateInitContainersFunc:   stubPodTemplateSection("init containers"),
			podTemplateContainersFunc:       stubPodTemplateSection("containers"),
			podTemplatePodConfigurationFunc: stubPodTemplateSection("configuration"),
			isErr:                           true,
		},
		{
			name:                            "init containers error",
			flexlayout:                      flexlayout.New(),
			podTemplateHeaderFunc:           stubPodTemplateSection("header"),
			podTemplateInitContainersFunc:   stubPodTemplateSectionWithError(),
			podTemplateContainersFunc:       stubPodTemplateSection("containers"),
			podTemplatePodConfigurationFunc: stubPodTemplateSection("configuration"),
			isErr:                           true,
		},
		{
			name:                            "containers error",
			flexlayout:                      flexlayout.New(),
			podTemplateHeaderFunc:           stubPodTemplateSection("header"),
			podTemplateInitContainersFunc:   stubPodTemplateSection("init containers"),
			podTemplateContainersFunc:       stubPodTemplateSectionWithError(),
			podTemplatePodConfigurationFunc: stubPodTemplateSection("configuration"),
			isErr:                           true,
		},
		{
			name:                            "pod configuration error",
			flexlayout:                      flexlayout.New(),
			podTemplateHeaderFunc:           stubPodTemplateSection("header"),
			podTemplateInitContainersFunc:   stubPodTemplateSection("init containers"),
			podTemplateContainersFunc:       stubPodTemplateSection("containers"),
			podTemplatePodConfigurationFunc: stubPodTemplateSectionWithError(),
			isErr:                           true,
		},
		{
			name:  "nil flexLayout",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			parent := testutil.CreateDeployment("deployment")
			spec := corev1.PodTemplateSpec{}

			pt := NewPodTemplate(parent, spec)

			pt.podTemplateHeaderFunc = tc.podTemplateHeaderFunc
			pt.podTemplateInitContainersFunc = tc.podTemplateInitContainersFunc
			pt.podTemplateContainersFunc = tc.podTemplateContainersFunc
			pt.podTemplatePodConfigurationFunc = tc.podTemplatePodConfigurationFunc

			options := Options{}

			ctx := context.Background()
			err := pt.AddToFlexLayout(ctx, tc.flexlayout, options)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			got := tc.flexlayout.ToComponent("Pod Template")

			assert.Equal(t, tc.expected(), got)
		})
	}
}

func Test_podTemplateHeader(t *testing.T) {
	fl := flexlayout.New()

	options := podTemplateLayoutOptions{
		podTemplateSpec: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"foo": "bar",
				},
			},
		},
	}

	ctx := context.Background()
	require.NoError(t, podTemplateHeader(ctx, fl, options))

	got := fl.ToComponent("Foo")

	expected := component.NewFlexLayout("Foo")
	expectedLabels := component.NewLabels(map[string]string{
		"foo": "bar",
	})
	expectedLabels.Metadata.SetTitleText("Pod Template")
	expected.AddSections([]component.FlexLayoutSection{
		{
			{
				Width: component.WidthFull,
				View:  expectedLabels,
			},
		},
	}...)

	assert.Equal(t, expected, got)
}
