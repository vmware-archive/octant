/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/pkg/view/component"
)

func TestContentGenerator_Event(t *testing.T) {
	expectedLabelSet := &labels.Set{}
	expectedPath := "/path"
	expectedPrefix := "/prefix"
	expectedNamespace := "namespace"

	contentResponse := component.ContentResponse{}

	expectedResponse := dashResponse{
		contentResponse,
	}
	expectedData, err := json.Marshal(&expectedResponse)
	require.NoError(t, err)

	g := ContentGenerator{
		ResponseFactory: func(ctx context.Context, path, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
			assert.Equal(t, expectedPath, path)
			assert.Equal(t, expectedPrefix, prefix)
			assert.Equal(t, expectedNamespace, namespace)
			assert.Equal(t, expectedLabelSet, opts.LabelSet)

			return contentResponse, nil
		},
		Path:      expectedPath,
		Prefix:    expectedPrefix,
		Namespace: expectedNamespace,
		LabelSet:  expectedLabelSet,
		RunEvery:  DefaultScheduleDelay,
	}

	ctx := context.Background()
	event, err := g.Event(ctx)
	require.NoError(t, err)

	assert.Equal(t, octant.EventTypeContent, event.Type)

	assert.Equal(t, expectedData, event.Data)
}

func TestContentGenerator_ScheduleDelay(t *testing.T) {
	g := ContentGenerator{
		RunEvery: DefaultScheduleDelay,
	}

	assert.Equal(t, DefaultScheduleDelay, g.ScheduleDelay())
}

func TestContentGenerator_Name(t *testing.T) {
	g := ContentGenerator{}
	assert.Equal(t, "content", g.Name())
}
