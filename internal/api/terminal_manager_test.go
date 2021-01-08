package api_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/api/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
)

func Test_TerminalStateManager(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	state := octantFake.NewMockState(controller)
	octantClient := fake.NewMockOctantClient(controller)

	tsm := api.NewTerminalStateManager(dashConfig)

	ctx := context.Background()
	tsm.Start(ctx, state, octantClient)
}
