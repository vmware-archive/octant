package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/pkg/api/fake"

	"github.com/vmware-tanzu/octant/internal/errors"
)

func TestNavigationManager(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	es, _ := errors.NewErrorStore()

	manager := NewNotificationsStateManager(es)

	state := octantFake.NewMockState(controller)

	ctx := context.Background()
	octantClient := fake.NewMockOctantClient(controller)

	octantClient.EXPECT().Send(gomock.Any())
	manager.Start(ctx, state, octantClient)

	es.AddError(fmt.Errorf("test"))
}

func TestHandlers(t *testing.T) {
	es, _ := errors.NewErrorStore()
	manager := NewNotificationsStateManager(es)
	h := manager.Handlers()

	assert.Equal(t, h[0].RequestType, RequestNotifications)
}
