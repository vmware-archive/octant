/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestApplyYamlDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	namespace := "default"

	dashConfig := configFake.NewMockDash(controller)

	p := NewApplyYamlDescriber()

	options := describer.Options{
		Dash: dashConfig,
	}

	cResponse, err := p.Describe(context.TODO(), namespace, options)
	require.NoError(t, err)

	list := component.NewList(append([]component.TitleComponent{}, component.NewText("Apply YAML")), nil)

	editor := component.NewEditor(component.TitleFromString("YAML"), "", false)
	editor.Config.SubmitAction = "action.octant.dev/apply"
	editor.Config.SubmitLabel = "Apply"
	list.Add(editor)

	require.Len(t, cResponse.Components, 1)
	component.AssertEqual(t, list, cResponse.Components[0])

	pf := p.PathFilters()
	require.Equal(t, "/apply", pf[0].String())

	err = p.Reset(context.TODO())
	require.NoError(t, err)
}
