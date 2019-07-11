/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/store/fake"
)

func TestConfigurationEditor(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()

	deployment := testutil.CreateDeployment("deployment")
	deployment.Namespace = "default"

	objectStore := fake.NewMockStore(controller)

	key, err := store.KeyFromObject(deployment)
	require.NoError(t, err)

	updatedDeployment := deployment.DeepCopy()
	updatedDeployment.Spec.Replicas = pointer.Int32Ptr(5)

	objectStore.EXPECT().
		Update(gomock.Any(), key, gomock.Any()).
		DoAndReturn(func(ctx context.Context, key store.Key, fn func(object *unstructured.Unstructured) error) error {
			return nil
		})

	configurationEditor := NewConfigurationEditor(logger, objectStore)
	assert.Equal(t, configurationEditorAction, configurationEditor.ActionName())

	ctx := context.Background()

	payload := action.Payload{
		"group":     "apps",
		"version":   "v1",
		"kind":      "Deployment",
		"namespace": "default",
		"name":      "deployment",
		"replicas":  "5",
	}

	require.NoError(t, configurationEditor.Handle(ctx, payload))

}
