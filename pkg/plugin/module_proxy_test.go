package plugin_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func TestModuleProxy_Name(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := fake.NewMockModuleService(controller)

	metadata := &plugin.Metadata{
		Name: "Test Plugin",
	}

	moduleProxy, err := plugin.NewModuleProxy("plugin-name", metadata, service)
	require.NoError(t, err)

	assert.Equal(t, metadata.Name, moduleProxy.Name())
}

func TestModuleProxy_ContentPath(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := fake.NewMockModuleService(controller)

	response := component.ContentResponse{}
	service.EXPECT().
		Content(gomock.Any(), "/path").
		Return(response, nil)

	metadata := &plugin.Metadata{
		Name: "Test Plugin",
	}

	moduleProxy, err := plugin.NewModuleProxy("plugin-name", metadata, service)
	require.NoError(t, err)

	ctx := context.Background()
	got, err := moduleProxy.Content(ctx, "/path", module.ContentOptions{})
	require.NoError(t, err)

	assert.Equal(t, response, got)
}

func TestModuleProxy_Navigation(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := fake.NewMockModuleService(controller)

	nav := navigation.Navigation{}
	service.EXPECT().
		Navigation(gomock.Any()).
		Return(nav, nil)

	metadata := &plugin.Metadata{
		Name: "Test Plugin",
	}

	moduleProxy, err := plugin.NewModuleProxy("plugin-name", metadata, service)
	require.NoError(t, err)

	ctx := context.Background()
	got, err := moduleProxy.Navigation(ctx, "", "")
	require.NoError(t, err)

	expected := []navigation.Navigation{
		nav,
	}

	assert.Equal(t, expected, got)
}
