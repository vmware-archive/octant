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

	"github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/octant"
)

func TestNamespacesGenerator_Event(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	namespaceClient := fake.NewMockNamespaceInterface(controller)

	namespaceClient.EXPECT().
		Names().
		Return([]string{"ns1", "ns2"}, nil)

	g := NamespacesGenerator{
		NamespaceClient: namespaceClient,
	}

	ctx := context.Background()
	event, err := g.Event(ctx)
	require.NoError(t, err)

	expectedResponse := namespacesResponse{
		Namespaces: []string{"ns1", "ns2"},
	}
	expectedData, err := json.Marshal(&expectedResponse)
	require.NoError(t, err)

	assert.Equal(t, octant.EventTypeNamespaces, event.Type)
	assert.Equal(t, expectedData, event.Data)
}

func TestNamespacesGenerator_ScheduleDelay(t *testing.T) {
	g := NamespacesGenerator{
	}

	assert.Equal(t, DefaultScheduleDelay, g.ScheduleDelay())
}

func TestNamespacesGenerator_Name(t *testing.T) {
	g := NamespacesGenerator{}
	assert.Equal(t, "namespaces", g.Name())
}
