package configuration

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/store/fake"
)

func TestObjectDeleter_ActionName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	objectStore := fake.NewMockStore(controller)

	logger := log.NopLogger()

	d := NewObjectDeleter(logger, objectStore)
	require.Equal(t, octant.ActionDeleteObject, d.ActionName())
}

func TestObjectDeleter_Handle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	objectStore := fake.NewMockStore(controller)

	pod := testutil.CreatePod("pod")
	key, err := store.KeyFromObject(pod)
	require.NoError(t, err)

	objectStore.EXPECT().
		Delete(gomock.Any(), key).
		Return(nil)

	logger := log.NopLogger()

	d := NewObjectDeleter(logger, objectStore)

	ctx := context.Background()

	err = d.Handle(ctx, key.ToActionPayload())
	require.NoError(t, err)
}
