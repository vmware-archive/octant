/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/heptio/developer-dash/internal/octant"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/module/fake"
)

func TestNavigationGenerator_Event(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mod := fake.NewMockModule(controller)
	mod.EXPECT().
		ContentPath().Return("/module").AnyTimes()
	mod.EXPECT().
		Navigation(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, namespace, prefix string) ([]octant.Navigation, error) {
			nav := octant.Navigation{
				Path:  prefix,
				Title: "module",
			}

			return []octant.Navigation{nav}, nil
		}).
		AnyTimes()

	g := NavigationGenerator{
		Modules: []module.Module{mod},
	}

	var ctx = context.Background()
	event, err := g.Event(ctx)
	require.NoError(t, err)

	expectedResponse := navigationResponse{
		Sections: []octant.Navigation{
			{
				Path:  "/content/module",
				Title: "module",
			},
		},
	}
	expectedData, err := json.Marshal(&expectedResponse)
	require.NoError(t, err)

	assert.Equal(t, octant.EventTypeNavigation, event.Type)
	assert.JSONEq(t, string(expectedData), string(event.Data))
	assert.Equal(t, expectedData, event.Data)
}

func TestNavigationGenerator_ScheduleDelay(t *testing.T) {
	g := NavigationGenerator{
		RunEvery: DefaultScheduleDelay,
	}

	assert.Equal(t, DefaultScheduleDelay, g.ScheduleDelay())
}

func TestNavigationGenerator_Name(t *testing.T) {
	g := NavigationGenerator{}
	assert.Equal(t, "navigation", g.Name())
}
